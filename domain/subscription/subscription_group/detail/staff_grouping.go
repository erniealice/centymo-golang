package detail

// staff_grouping.go — the Subjects tab's umbrella grouping (one rendered row
// per PRODUCT). Split from staff.go so the grouping surface stays a focused
// unit; the row types + cell/action builders it composes live in staff.go.

import (
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"strings"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/pyeza-golang/types"
)

// =============================================================================
// Umbrella grouping — one rendered row per PRODUCT (owner call, 2026-07-25).
// A product whose plan carries several variant offerings (each its own class
// row) collapses into ONE row: chips per variant, merged Teachers text, ONE
// View (the lead member's class page), a whole-product Remove. A group of one
// is a strict no-op (sgppTableRow renders it untouched).
//
// Since the strand collapse (docs/plan/20260725-arts-design-strand-collapse)
// every product carries exactly one class per section, so multi-member groups
// are the defensive path, not the steady state — which is why the grouped row
// renders a single View: with one class per product the per-member View fan
// (and its per-variant action labels) is dead machinery. The strand signal on
// a collapsed class is the Teachers column's phase parentheticals.
// =============================================================================

// sgppGroup is one rendered row of the grouped tab: the display label plus
// its member class rows in DETERMINISTIC order (offering label, then variant
// label, then id — never map order). Members[0] is the group's lead row: its
// class id is the row identity (testids, data-id) so a group of one keeps
// today's exact identity.
type sgppGroup struct {
	Key     string
	Label   string
	Members []SectionSGPPRow
}

// groupSGPPRows groups class rows by product identity. Only MATERIALIZED rows
// with a resolved ProductID group; fallback (unmaterialized) rows and rows
// whose offering did not resolve stay singleton groups — an unresolved
// product must never accidentally merge with another unresolved one. Group
// order mirrors today's row order (case-insensitive label, id tiebreak) and
// is independent of input order.
func groupSGPPRows(rows []SectionSGPPRow) []sgppGroup {
	keyOrder := make([]string, 0, len(rows))
	byKey := map[string][]SectionSGPPRow{}
	for _, r := range rows {
		key := "row:" + sgppRowKey(r)
		if r.Materialized && r.ProductID != "" {
			key = "product:" + r.ProductID
		}
		if _, seen := byKey[key]; !seen {
			keyOrder = append(keyOrder, key)
		}
		byKey[key] = append(byKey[key], r)
	}
	groups := make([]sgppGroup, 0, len(keyOrder))
	for _, key := range keyOrder {
		members := byKey[key]
		sort.SliceStable(members, func(i, j int) bool {
			a, b := strings.ToLower(members[i].OfferingLabel), strings.ToLower(members[j].OfferingLabel)
			if a != b {
				return a < b
			}
			va, vb := strings.ToLower(members[i].VariantLabel), strings.ToLower(members[j].VariantLabel)
			if va != vb {
				return va < vb
			}
			return sgppRowKey(members[i]) < sgppRowKey(members[j])
		})
		groups = append(groups, sgppGroup{Key: key, Label: productDisplayLabel(members), Members: members})
	}
	sort.SliceStable(groups, func(i, j int) bool {
		a, b := strings.ToLower(groups[i].Label), strings.ToLower(groups[j].Label)
		if a == b {
			return groups[i].Key < groups[j].Key
		}
		return a < b
	})
	return groups
}

// productDisplayLabel derives a group's display label. Precedence: (1) a
// single member keeps its offering label verbatim (the no-op contract);
// (2) the hydrated product name when the adapter supplied product_plan.product;
// (3) the shared offering name when every member agrees (the end-state once
// offering names are the bare product name); (4) the members' longest common
// label prefix, trimmed of trailing separator punctuation — a RULE, never a
// hand-typed list; (5) the lead member's label as the deterministic fallback.
func productDisplayLabel(members []SectionSGPPRow) string {
	if len(members) == 0 {
		return ""
	}
	first := members[0].OfferingLabel
	if len(members) == 1 {
		return first
	}
	for _, m := range members {
		if l := strings.TrimSpace(m.ProductLabel); l != "" {
			return l
		}
	}
	same := true
	for _, m := range members[1:] {
		if m.OfferingLabel != first {
			same = false
			break
		}
	}
	if same {
		return first
	}
	prefix := []rune(first)
	for _, m := range members[1:] {
		r := []rune(m.OfferingLabel)
		n := 0
		for n < len(prefix) && n < len(r) && prefix[n] == r[n] {
			n++
		}
		prefix = prefix[:n]
	}
	if trimmed := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(string(prefix)), ":-–—·/|")); trimmed != "" {
		return trimmed
	}
	return first
}

// sgppGroupTableRow renders one MULTI-member group. Row identity = the lead
// member's class id (deterministic, sgpp-row-{leadId} testid). Direct action:
// exactly ONE View, opening the LEAD member's class detail page (post-collapse
// a product has one class, so one View is the truthful shape; on a residual
// multi-member group the lead is the deterministic representative). Overflow:
// Remove, the whole-product delete ("remove the class rows" — nothing more),
// POSTed once through the bulk-delete route (non-lead ids ride the URL query;
// the shared table JS appends the lead id from data-id, so every member id is
// submitted exactly once). Assign/Exclude/Restore deliberately do NOT appear
// here — per-class management escalates to the class page (View), where the
// target is explicit; the grouped row still SIGNALS per-variant exceptions
// via the chips (sgppGroupSubjectCell).
func sgppGroupTableRow(deps *DetailViewDeps, groupID string, g sgppGroup, perms *types.UserPermissions, l subscription_group.Labels) types.TableRow {
	lead := g.Members[0]
	key := sgppRowKey(lead)

	merged := make([]sgppAssignmentDisplay, 0)
	allExcluded := true
	anyInUse := false
	for _, m := range g.Members {
		merged = append(merged, m.assignments...)
		if !m.Excluded {
			allExcluded = false
		}
		if m.InUse {
			anyInUse = true
		}
	}
	status := "active"
	if allExcluded {
		status = "excluded"
	}
	// The Teachers cell reuses the SAME composer + cell renderer as an
	// ungrouped row (one joiner, one lyngua separator/aria) over the merged
	// member assignments.
	display := SectionSGPPRow{ID: key, Materialized: true, AssignmentsText: joinAssignments(merged, l)}

	row := types.TableRow{
		ID: key,
		DataAttrs: map[string]string{
			"testid": "sgpp-row-" + key,
			"status": status,
		},
		Cells: []types.TableCell{sgppGroupSubjectCell(g, l), sgppAssignmentsCell(display, l)},
	}

	viewAction := types.TableAction{
		Type: "view", Action: "view",
		Label:  l.Staff.ViewAction,
		TestID: "sgpp-view-" + lead.ID,
	}
	if deps.SGPPDetailURL != nil {
		viewAction.Href = deps.SGPPDetailURL(groupID, lead.ID)
	} else {
		viewAction.Disabled = true
	}
	row.Actions = append(row.Actions, viewAction)

	removeAction := types.TableAction{
		Type: "delete", Action: "delete", Label: l.Staff.RemoveAction, Overflow: true,
		URL: sgppGroupRemoveURL(deps.SGPPBulkDeleteURL, g.Members), ItemName: g.Label,
		ConfirmTitle: l.Staff.RemoveConfirmTitle, ConfirmMessage: fmt.Sprintf(l.Staff.RemoveConfirmMessage, g.Label),
		TestID: "sgpp-remove-" + key,
	}
	switch {
	case anyInUse:
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Staff.RemoveBlockedInUse
	case !perms.Can(sgppEntity, "delete"):
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Staff.Unauthorized
	case removeAction.URL == "":
		removeAction.Disabled = true
	}
	row.Actions = append(row.Actions, removeAction)

	return row
}

// sgppGroupSubjectCell renders the grouped Subject cell: the product label
// plus one chip per member in member order — the variant name on a normal
// strand, the VariantExcludedTag rendering when THAT strand is EXCLUDED (the
// per-variant exception stays visible on the grouped row), and no chip at all
// for an un-varianted member (never an empty badge). Every interpolated value
// is HTML-escaped.
func sgppGroupSubjectCell(g sgppGroup, l subscription_group.Labels) types.TableCell {
	var b strings.Builder
	b.WriteString(`<span class="centymo-sgpp-subject-name">`)
	b.WriteString(template.HTMLEscapeString(g.Label))
	b.WriteString(`</span>`)
	for _, m := range g.Members {
		switch {
		case m.Excluded && m.VariantLabel != "":
			b.WriteString(` <span class="badge badge--warning" data-testid="sgpp-excluded-tag-` + template.HTMLEscapeString(m.ID) + `">`)
			b.WriteString(template.HTMLEscapeString(strings.ReplaceAll(l.Staff.VariantExcludedTag, "{{variant}}", m.VariantLabel)))
			b.WriteString(`</span>`)
		case m.Excluded:
			b.WriteString(` <span class="badge badge--warning" data-testid="sgpp-excluded-tag-` + template.HTMLEscapeString(m.ID) + `">`)
			b.WriteString(template.HTMLEscapeString(l.Staff.ExcludedTag))
			b.WriteString(`</span>`)
		case m.VariantLabel != "":
			b.WriteString(` <span class="badge badge--muted" data-testid="sgpp-variant-` + template.HTMLEscapeString(m.ID) + `">`)
			b.WriteString(template.HTMLEscapeString(m.VariantLabel))
			b.WriteString(`</span>`)
		}
	}
	return types.TableCell{Type: "html", HTML: template.HTML(b.String()), Value: g.Label}
}

// sgppGroupRemoveURL builds the grouped Remove's POST target on the
// bulk-delete route: every NON-lead member id is baked into the query; the
// table JS appends the lead id (`data-id`) on click, completing the set.
// The action_workspace_guard signs over the URL path only, so the baked query
// ids do not disturb the rowActionTokens signature. Returns "" (a disabled
// action) when the bulk route is unwired.
func sgppGroupRemoveURL(base string, members []SectionSGPPRow) string {
	if base == "" || len(members) < 2 {
		return ""
	}
	var b strings.Builder
	b.WriteString(base)
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	for _, m := range members[1:] {
		b.WriteString(sep)
		b.WriteString("id=")
		b.WriteString(url.QueryEscape(m.ID))
		sep = "&"
	}
	return b.String()
}
