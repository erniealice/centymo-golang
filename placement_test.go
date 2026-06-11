// placement_test.go — centymo's adoption of the CANONICAL Wave-T placement gate.
//
// COPIED verbatim from
// docs/orchestrate/20260610-package-cleanup/placement_test.template.go, with the
// `//go:build ignore` tag stripped and ONLY the per-package config block edited
// (crossCutting / legacyAllow / charterViews / the package clause). Everything
// below the config block is the source-of-truth logic and must NOT be edited
// per-package — that is the whole point of a single parameterized gate
// (manifest §6 TT).
//
// It is pure stdlib (go/parser, go/ast, os, path/filepath, strings, testing) and
// derives the esqyma domain set + the entity→domain map LIVE from
// packages/esqyma/proto/v1/domain/ at test time, so the rules can never drift
// from proto. See docs/plan/20260610-centymo-restructuring/placement-test.md.
//
// The four rules (domain variant):
//   R1 Empty root      — no package .go at the module root (only *_test.go).
//   R2 Canonical dirs  — every domain/<d> is an esqyma proto domain; the only
//                        other first-level dirs are infra surfaces.
//   R3 Entity placement— each exported XxxLabels/XxxRoutes lives under the
//                        domain dir that esqyma says owns its entity.
//   R4 No god-files    — no .go > 1200 lines.
// Cross-cutting variant (crossCutting==true, e.g. hybra): skips R1/R2/R3, asserts
// views/<x> ∈ charterViews and no framework-leak files at root; keeps R4.
//
// (The body below is duplicated verbatim from the template. Keep it
//  byte-identical except for the config block.)

package centymo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// ── per-package config (the ONLY part that differs between packages) ──────────
var crossCutting = false // true ONLY for a cross-cutting pkg (hybra)

// legacyAllow — the shrinking migration ledger (W0 → W9). Each wave DELETES its
// entries; a non-empty set is remaining debt (printed by -v). Matched on a path's
// FIRST segment OR its basename. W9 empties this and the gate goes strict.
var legacyAllow = map[string]string{
	"labels.go":             "W1-W7: dissolve per-domain into domain/<d>/labels.go (10199 LoC, 434 *Labels structs across 33 sections)",
	"routes_config.go":      "W1-W7: dissolve per-domain into domain/<d>/routes.go (2730 LoC, 32 Default*Routes funcs)",
	"routes.go":             "W1-W7: dissolve per-domain into domain/<d>/routes.go (770 LoC URL consts)",
	"routes_config_test.go": "W1-W7: moves/splits alongside routes_config.go (153 LoC)",
	"advance_actions.go":    "W5: -> domain/treasury/advance.go (86 LoC, 8 Advance*ViewInput/Output structs)",
	"advance_labels.go":     "W5+W6: SPLIT -> treasury/advance.go (L1-181 Advance/AdvancesDashboard) + expenditure/supplier_billing_event_labels.go (L182-498 SupplierBillingEvent*) per esqyma domainOf",
	"views":                 "W1-W7: regroup 35 entity-flat dirs under domain/<d>/views/<entity>/",
	"services":              "W8 (DEFERRED): -> espyna (checkout) + keep serial.go/service.go/types.go pending; out of scope for W1-W7",
	"datasource.go":         "DEFERRAL (D2/TD): DataSource legacy view-data port — keep at root, do not touch in W1-W7",
	"assets.go":             "DEFERRAL: AssetsFS //go:embed (post-Wave-P) — keep at root strict-root caveat, do not touch in W1-W7",
	"docs":                  "PERMANENT (not debt): centymo/docs/ holds plan markdown, no Go. Resolve at W9 by adding \"docs\" to the SHARED allowedFirstLevelDirs across all adoptions (cyta/hybra/centymo) — a template-level change — NOT by deleting this line.",
	// ── W4 (subscription) — R3 mechanical-longest-match false positives ──────
	// These 3 files hold types whose NAMES resolve to another esqyma domain
	// (ProductPlanForm/ProductKindOption -> product; ClientPackages -> entity)
	// but which are functionally subscription-aggregate sub-types / projections.
	// They stay in package subscription for cohesion (moving them would force a
	// subscription->product / ->entity sibling import). Renaming is out of scope
	// for a pure structural wave. W9 resolves the naming, then deletes these.
	"plan_product_subtypes.go":  "W9: rename ProductPlanForm/ProductKindOption sub-types (R3 longest-match collides with product domain); functionally Plan sub-types nested in PlanLabels",
	"client_packages_labels.go": "W9: rename ClientPackages projection (R3 longest-match collides with entity/client); subscription-domain client-Packages-tab labels",
	"product_plan.go":           "W9: view-local ProductPlanFormLabels in domain/subscription/views/plan/action (R3 longest-match collides with product domain); single file, single type",
	// ── W4 (subscription) — R4 pre-existing oversized view handlers ──────────
	// detail/page.go for subscription (1916) + price_plan (1720) were ALREADY
	// >1200 at HEAD, excused under the "views" first-segment before W4 re-rooted
	// them under domain/. Splitting view handlers is a separate (view-split)
	// wave, not this pure label/route relocation. Only these 2 page.go exceed
	// the threshold today, so this basename excuses no other current violation.
	"page.go": "VIEW-SPLIT wave: subscription/detail/page.go (1916) + price_plan/detail/page.go (1720) pre-existing >1200 view handlers, re-rooted from views/ in W4; split per concern in a dedicated wave",
}
var charterViews = []string{} // crossCutting only: allowed views/<x> concern groups

// ── shared logic — DO NOT EDIT per package ───────────────────────────────────

const godFileThreshold = 1200

// allowedFirstLevelDirs are the non-domain first-level dirs a domain package may
// hold besides domain/. Note BOTH "service" and "services" are infra surfaces
// (cyta uses the plural for its private recurrence/availability helpers).
var allowedFirstLevelDirs = map[string]bool{
	"domain":   true,
	"block":    true,
	"assets":   true,
	"service":  true,
	"services": true,
	"scripts":  true,
	"internal": true,
	"tests":    true,
	"web":      true,
}

// frameworkLeakFiles must never appear at a cross-cutting package's root — these
// concerns belong in pyeza (manifest §3 Wave P).
var frameworkLeakFiles = map[string]bool{
	"htmx.go":        true,
	"assets.go":      true,
	"datasource.go":  true,
	"package_dir.go": true,
	"pkgdir.go":      true,
}

// locateEsqymaDomain walks up from the CWD looking for the esqyma proto domain
// dir, then derives the domain set + entity→domain map from its subdirs.
func locateEsqymaDomain(t *testing.T) (root string, domainSet map[string]bool, entityDomain map[string]string) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("placement: cannot get cwd: %v", err)
	}
	rel := filepath.Join("proto", "v1", "domain")
	var candidate string
	// Prefer a sibling esqyma checkout (packages/<pkg>/.. == packages/), then
	// walk up looking for packages/esqyma/proto/v1/domain or esqyma/proto/...
	dir := cwd
	for {
		for _, c := range []string{
			filepath.Join(dir, "..", "esqyma", rel),
			filepath.Join(dir, "packages", "esqyma", rel),
			filepath.Join(dir, "esqyma", rel),
		} {
			if fi, err := os.Stat(c); err == nil && fi.IsDir() {
				candidate = c
				break
			}
		}
		if candidate != "" {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	if candidate == "" {
		t.Fatalf("placement: could not locate packages/esqyma/proto/v1/domain from %s — the gate cannot run without the esqyma source of truth", cwd)
	}
	root = filepath.Clean(candidate)

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("placement: cannot read esqyma domain dir %s: %v", root, err)
	}
	domainSet = map[string]bool{}
	entityDomain = map[string]string{}
	for _, de := range entries {
		if !de.IsDir() {
			continue
		}
		domain := de.Name()
		domainSet[domain] = true
		ents, err := os.ReadDir(filepath.Join(root, domain))
		if err != nil {
			continue
		}
		for _, e := range ents {
			if e.IsDir() {
				entityDomain[e.Name()] = domain
			}
		}
	}
	if len(domainSet) == 0 {
		t.Fatalf("placement: esqyma domain dir %s has no domains", root)
	}
	return root, domainSet, entityDomain
}

// moduleRoot returns the directory holding the package's go.mod (the placement
// gate's module root). It walks up from CWD; CWD is already the module root for a
// root-level placement_test.go, but walking up keeps the test correct if it is
// hosted in a subdir (e.g. internal/structure/).
func moduleRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("placement: cannot get cwd: %v", err)
	}
	dir := cwd
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("placement: no go.mod found walking up from %s", cwd)
		}
		dir = parent
	}
}

// inLegacyAllow reports whether a path (relative to module root) is excused by
// legacyAllow — matched on its FIRST path segment OR its basename. That is the
// shrinking migration ledger: each wave deletes entries.
func inLegacyAllow(relPath string) bool {
	if len(legacyAllow) == 0 {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	first := relPath
	if i := strings.IndexByte(relPath, '/'); i >= 0 {
		first = relPath[:i]
	}
	if _, ok := legacyAllow[first]; ok {
		return true
	}
	base := filepath.Base(relPath)
	_, ok := legacyAllow[base]
	return ok
}

// camelToSnake converts a CamelCase identifier prefix to snake_case
// (EventTagButton -> event_tag_button, ProductPricePlan -> product_price_plan).
func camelToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r - 'A' + 'a')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// resolveEntity maps a CamelCase type prefix (the part before Labels/Routes) to
// the LONGEST known esqyma entity that is a segment-aligned snake_case prefix of
// it. EventTagButton -> event_tag (not event); ProductPricePlan ->
// product_price_plan (not plan). Returns "" if no known entity matches (the
// prefix may be a Surface/projection/page type — R3 does not fail on those).
func resolveEntity(prefix string, entityDomain map[string]string) string {
	snake := camelToSnake(prefix)
	segs := strings.Split(snake, "_")
	best := ""
	// Try the longest segment-aligned prefix first: full -> shorter.
	for n := len(segs); n >= 1; n-- {
		cand := strings.Join(segs[:n], "_")
		if _, ok := entityDomain[cand]; ok {
			best = cand
			break // first hit is the longest by construction
		}
	}
	return best
}

// labelsRoutesTypes parses a .go file and returns the names of exported type
// decls whose name ends in Labels or Routes.
func labelsRoutesTypes(path string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || !ts.Name.IsExported() {
				continue
			}
			name := ts.Name.Name
			if strings.HasSuffix(name, "Labels") || strings.HasSuffix(name, "Routes") {
				out = append(out, name)
			}
		}
	}
	return out, nil
}

func countLines(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	n := strings.Count(string(data), "\n")
	if data[len(data)-1] != '\n' {
		n++ // last line without trailing newline
	}
	return n, nil
}

func isGoFile(name string) bool   { return strings.HasSuffix(name, ".go") }
func isTestFile(name string) bool { return strings.HasSuffix(name, "_test.go") }

func TestPlacement(t *testing.T) {
	root := moduleRoot(t)
	_, domainSet, entityDomain := locateEsqymaDomain(t)

	// R4 (all variants): no god-files anywhere (excl. *_test.go), unless excused.
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !isGoFile(info.Name()) || isTestFile(info.Name()) {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if inLegacyAllow(rel) {
			return nil
		}
		n, cerr := countLines(path)
		if cerr != nil {
			return nil
		}
		if n > godFileThreshold {
			t.Errorf("%s: %d lines exceeds the %d god-file threshold — split per entity", rel, n, godFileThreshold)
		}
		return nil
	})

	if crossCutting {
		runCrossCutting(t, root)
	} else {
		runDomainVariant(t, root, domainSet, entityDomain)
	}

	if testing.Verbose() {
		if len(legacyAllow) == 0 {
			t.Logf("placement: legacyAllow EMPTY — STRICT gate (no remaining migration debt)")
		} else {
			keys := make([]string, 0, len(legacyAllow))
			for k := range legacyAllow {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			t.Logf("placement: legacyAllow remaining debt (%d):", len(keys))
			for _, k := range keys {
				t.Logf("  - %s: %s", k, legacyAllow[k])
			}
		}
	}
}

func runDomainVariant(t *testing.T, root string, domainSet map[string]bool, entityDomain map[string]string) {
	// R1 Empty root: no package .go directly at module root (only *_test.go).
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("placement: cannot read module root %s: %v", root, err)
	}
	for _, de := range rootEntries {
		if de.IsDir() || !isGoFile(de.Name()) || isTestFile(de.Name()) {
			continue
		}
		if inLegacyAllow(de.Name()) {
			continue
		}
		t.Errorf("%s: root holds no package code — re-home (→ domain/<d>/, → pyeza, or owning pkg)", de.Name())
	}

	// R2 Canonical domains: every first-level dir is an allowed infra dir or
	// `domain`; every subdir of domain/ is an esqyma proto domain.
	for _, de := range rootEntries {
		if !de.IsDir() {
			continue
		}
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if inLegacyAllow(name) {
			continue
		}
		if name == "domain" || allowedFirstLevelDirs[name] {
			continue
		}
		t.Errorf("%s/: not an esqyma proto domain — fold into the owning domain or a service/ surface", name)
	}
	domainDir := filepath.Join(root, "domain")
	if domEntries, err := os.ReadDir(domainDir); err == nil {
		for _, de := range domEntries {
			if !de.IsDir() {
				continue
			}
			d := de.Name()
			rel := filepath.ToSlash(filepath.Join("domain", d))
			if inLegacyAllow(rel) || inLegacyAllow(d) {
				continue
			}
			if !domainSet[d] {
				t.Errorf("%s/: not an esqyma proto domain — fold into the owning domain or a service/ surface", rel)
			}
		}
	}

	// R3 Entity placement: each exported XxxLabels/XxxRoutes under domain/<d>/
	// must have domainOf(entity) == <d> (or unknown entity → report only).
	if domEntries, err := os.ReadDir(domainDir); err == nil {
		for _, de := range domEntries {
			if !de.IsDir() {
				continue
			}
			d := de.Name()
			if !domainSet[d] {
				continue // R2 already flagged it
			}
			ddir := filepath.Join(domainDir, d)
			_ = filepath.Walk(ddir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || !isGoFile(info.Name()) || isTestFile(info.Name()) {
					return nil
				}
				rel, _ := filepath.Rel(root, path)
				if inLegacyAllow(rel) || inLegacyAllow(info.Name()) {
					return nil
				}
				types, perr := labelsRoutesTypes(path)
				if perr != nil {
					return nil
				}
				for _, typeName := range types {
					prefix := strings.TrimSuffix(strings.TrimSuffix(typeName, "Labels"), "Routes")
					if prefix == "" {
						continue
					}
					entity := resolveEntity(prefix, entityDomain)
					if entity == "" {
						if testing.Verbose() {
							t.Logf("placement: %s:%s — prefix %q maps to no esqyma entity (Surface/projection?), skipped by R3", rel, typeName, camelToSnake(prefix))
						}
						continue
					}
					owner := entityDomain[entity]
					if owner != "" && owner != d {
						t.Errorf("%s:%s: entity %s belongs to domain/%s/, found in domain/%s/", rel, typeName, entity, owner, d)
					}
				}
				return nil
			})
		}
	}
}

func runCrossCutting(t *testing.T, root string) {
	// No framework-leak files at root.
	if rootEntries, err := os.ReadDir(root); err == nil {
		for _, de := range rootEntries {
			if de.IsDir() {
				continue
			}
			if frameworkLeakFiles[de.Name()] && !inLegacyAllow(de.Name()) {
				t.Errorf("%s: framework concern leaked to root — belongs in pyeza (Wave P)", de.Name())
			}
		}
	}
	// Every subdir of views/ must be a chartered concern group.
	charter := map[string]bool{}
	for _, c := range charterViews {
		charter[c] = true
	}
	viewsDir := filepath.Join(root, "views")
	if vEntries, err := os.ReadDir(viewsDir); err == nil {
		for _, de := range vEntries {
			if !de.IsDir() {
				continue
			}
			name := de.Name()
			if inLegacyAllow(filepath.ToSlash(filepath.Join("views", name))) || inLegacyAllow(name) {
				continue
			}
			if !charter[name] {
				t.Errorf("views/%s/: not a chartered cross-cutting concern group — expected one of %v", name, charterViews)
			}
		}
	}
}
