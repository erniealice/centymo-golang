package detail

// staff_grouping_test.go — unit coverage for the umbrella-grouping contract
// (owner call, 2026-07-25): the section tab renders ONE row per PRODUCT, with
// the product's variant offerings as chips and every member's teachers merged
// through the one shared joiner.
//
// The pins, in order of blast radius:
//  1. A product with ONE offering renders EXACTLY as before grouping — a
//     group of one is a byte-identical no-op (if this breaks, every subject
//     on every section changes).
//  2. An un-varianted single offering (the whole-product class shape) gets no
//     chip, no empty badge, no panic.
//  3. Grouping is by product identity, deterministic regardless of input
//     order, and never merges rows whose product did not resolve.
//  4. The grouped row: chips in deterministic order, per-variant EXCLUDED
//     signalling, merged Teachers text via joinAssignments, one direct View
//     per member class (ActionVariant-labeled), and a whole-product Remove
//     that submits every member id exactly once via the bulk-delete route.

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/pyeza-golang/types"
)

// groupingDeps wires the cross-domain URL closures with predictable shapes so
// action wiring is exercised (not the nil-degraded branches).
func groupingDeps() *DetailViewDeps {
	return &DetailViewDeps{
		Routes: subscription_group.DefaultRoutes(),
		Labels: subscription_group.DefaultLabels(),
		SGPPDetailURL: func(sectionID, sgppID string) string {
			return "/detail/" + sectionID + "/" + sgppID
		},
		SGPPAssignURL:     func(sgppID string) string { return "/assign/" + sgppID },
		SGPPSetStatusURL:  func(sgppID, status string) string { return "/set-status/" + sgppID + "?status=" + status },
		SGPPDeleteURL:     "/delete",
		SGPPBulkDeleteURL: "/bulk-delete",
	}
}

func groupingPerms(codes ...string) *types.UserPermissions {
	return types.NewUserPermissions(codes)
}

// allSGPPPerms grants every verb the row actions gate on.
func allSGPPPerms() *types.UserPermissions {
	return groupingPerms(
		sgppEntity+":list", sgppEntity+":create", sgppEntity+":update", sgppEntity+":delete",
		sgppsCreate, sgppsUpdate, sgppsDelete,
	)
}

// singleRow builds one materialized single-offering row (its own product).
func singleRow(id, label, productID string) SectionSGPPRow {
	return SectionSGPPRow{
		ID:            "c-" + id,
		ProductPlanID: "pp-" + id,
		OfferingLabel: label,
		Materialized:  true,
		ProductID:     productID,
	}
}

// TestBuildSGPPTable_SingleOfferingGroupsAreNoOp is pin #1: with every
// product carrying exactly one offering (incl. an un-varianted one and an
// EXCLUDED one), the grouped table's rows must be DEEP-EQUAL to what the
// ungrouped per-row builder produces for the same rows in the same order.
func TestBuildSGPPTable_SingleOfferingGroupsAreNoOp(t *testing.T) {
	deps := groupingDeps()
	l := deps.Labels
	perms := allSGPPPerms()

	rows := []SectionSGPPRow{
		singleRow("m", "Course A", "prod-a"),
		func() SectionSGPPRow { // un-varianted, whole-product class (pin #2)
			r := singleRow("k", "Course B", "prod-b")
			r.AssignmentsText = "Alma Reyes"
			return r
		}(),
		func() SectionSGPPRow { // EXCLUDED single offering keeps today's tag + actions
			r := singleRow("x", "Course C", "prod-c")
			r.Excluded = true
			return r
		}(),
		func() SectionSGPPRow { // varianted but the product's ONLY offering
			r := singleRow("v", "Course D: Strand One", "prod-d")
			r.VariantLabel = "Strand One"
			return r
		}(),
	}

	data := &SectionSGPPTabData{SubscriptionGroupID: "sec1", Rows: rows}
	cfg := buildSGPPTable(deps, data, perms, l)

	want := make([]types.TableRow, 0, len(rows))
	for _, r := range rows { // input already in label order
		want = append(want, sgppTableRow(deps, "sec1", r, perms, l))
	}
	// buildSGPPTable runs ApplyColumnStyles over its rows (copying column
	// label/width into cells); normalize the expected rows identically so the
	// deep-equal compares the ROW CONTENT, not the styling pass.
	types.ApplyColumnStyles([]types.TableColumn{
		{Key: "subject", Label: l.Staff.ColumnSubject, NoSort: true, NoFilter: true, Width: "40%"},
		{Key: "assignments", Label: l.Staff.ColumnServicer, NoSort: true, NoFilter: true},
	}, want)
	if !reflect.DeepEqual(cfg.Rows, want) {
		t.Fatalf("single-offering groups changed the rendered rows:\n got %+v\nwant %+v", cfg.Rows, want)
	}
	for _, tr := range cfg.Rows {
		if html := string(tr.Cells[0].HTML); strings.Contains(html, `badge--muted"></span>`) || strings.Contains(html, `badge--muted">​`) {
			t.Errorf("empty variant badge rendered: %q", html)
		}
	}
}

// TestGroupSGPPRows_ByProductDeterministic is pin #3: same-product rows merge
// into one group, member and group order is deterministic regardless of input
// order, unresolved-product rows stay singletons, and unmaterialized fallback
// rows never merge.
func TestGroupSGPPRows_ByProductDeterministic(t *testing.T) {
	strandA := SectionSGPPRow{ID: "c1", ProductPlanID: "pp1", OfferingLabel: "Course R: One", VariantLabel: "One", Materialized: true, ProductID: "prod-r"}
	strandB := SectionSGPPRow{ID: "c2", ProductPlanID: "pp2", OfferingLabel: "Course R: Two", VariantLabel: "Two", Materialized: true, ProductID: "prod-r"}
	single := singleRow("s", "Course A", "prod-a")
	noProduct1 := SectionSGPPRow{ID: "c3", ProductPlanID: "pp3", OfferingLabel: "pp3", Materialized: true}
	noProduct2 := SectionSGPPRow{ID: "c4", ProductPlanID: "pp4", OfferingLabel: "pp4", Materialized: true}
	fallback1 := SectionSGPPRow{ProductPlanID: "pp5", OfferingLabel: "Course F: One", Materialized: false, ProductID: "prod-f"}
	fallback2 := SectionSGPPRow{ProductPlanID: "pp6", OfferingLabel: "Course F: Two", Materialized: false, ProductID: "prod-f"}

	inputs := [][]SectionSGPPRow{
		{strandA, strandB, single, noProduct1, noProduct2, fallback1, fallback2},
		{fallback2, noProduct2, strandB, single, fallback1, strandA, noProduct1},
	}
	var first []sgppGroup
	for i, in := range inputs {
		got := groupSGPPRows(in)
		if i == 0 {
			first = got
			continue
		}
		if !reflect.DeepEqual(got, first) {
			t.Fatalf("grouping depends on input order:\n got %+v\nwant %+v", got, first)
		}
	}

	byLabel := map[string]sgppGroup{}
	for _, g := range first {
		byLabel[g.Label] = g
	}
	r, ok := byLabel["Course R"]
	if !ok || len(r.Members) != 2 {
		t.Fatalf("same-product strands did not merge into one 'Course R' group: %+v", first)
	}
	if r.Members[0].ID != "c1" || r.Members[1].ID != "c2" {
		t.Errorf("member order not deterministic by label: %+v", r.Members)
	}
	if g := byLabel["pp3"]; len(g.Members) != 1 {
		t.Errorf("unresolved-product row merged: %+v", first)
	}
	if g := byLabel["pp4"]; len(g.Members) != 1 {
		t.Errorf("unresolved-product row merged: %+v", first)
	}
	// Unmaterialized fallback rows share prod-f but must NOT merge.
	fallbackGroups := 0
	for _, g := range first {
		if len(g.Members) == 1 && !g.Members[0].Materialized {
			fallbackGroups++
		}
	}
	if fallbackGroups != 2 {
		t.Errorf("fallback rows merged (want 2 singleton groups, got %d): %+v", fallbackGroups, first)
	}
}

// TestProductDisplayLabel covers the label-derivation rule ladder.
func TestProductDisplayLabel(t *testing.T) {
	tests := []struct {
		name    string
		members []SectionSGPPRow
		want    string
	}{
		{
			name:    "single member keeps its offering label verbatim",
			members: []SectionSGPPRow{{OfferingLabel: "Course A — Level 10"}},
			want:    "Course A — Level 10",
		},
		{
			name: "hydrated product name wins",
			members: []SectionSGPPRow{
				{OfferingLabel: "Course R: One", ProductLabel: "Course R"},
				{OfferingLabel: "Course R: Two"},
			},
			want: "Course R",
		},
		{
			name: "identical member names (the post-rename end-state)",
			members: []SectionSGPPRow{
				{OfferingLabel: "Course R"},
				{OfferingLabel: "Course R"},
			},
			want: "Course R",
		},
		{
			name: "colon-suffixed strands derive the shared prefix",
			members: []SectionSGPPRow{
				{OfferingLabel: "Course R: One"},
				{OfferingLabel: "Course R: Two"},
			},
			want: "Course R",
		},
		{
			name: "em-dash-suffixed members derive the shared prefix",
			members: []SectionSGPPRow{
				{OfferingLabel: "Course R — Alpha"},
				{OfferingLabel: "Course R — Beta"},
			},
			want: "Course R",
		},
		{
			name: "no common prefix falls back to the lead member's label",
			members: []SectionSGPPRow{
				{OfferingLabel: "Alpha"},
				{OfferingLabel: "Beta"},
			},
			want: "Alpha",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := productDisplayLabel(tc.members); got != tc.want {
				t.Errorf("productDisplayLabel() = %q, want %q", got, tc.want)
			}
		})
	}
}

// groupedFixture returns a two-strand group (both materialized, same product)
// with per-member assignments, ready for sgppGroupTableRow.
func groupedFixture() sgppGroup {
	m1 := SectionSGPPRow{
		ID: "c1", ProductPlanID: "pp1", OfferingLabel: "Course R: One", VariantLabel: "One",
		Materialized: true, ProductID: "prod-r",
		assignments: []sgppAssignmentDisplay{{staffName: "Bea Cruz", phaseLabel: "Phase 1", order: 1}},
	}
	m2 := SectionSGPPRow{
		ID: "c2", ProductPlanID: "pp2", OfferingLabel: "Course R: Two", VariantLabel: "Two",
		Materialized: true, ProductID: "prod-r",
		assignments: []sgppAssignmentDisplay{{staffName: "Alma Reyes", phaseLabel: "Phase 2", order: 2}},
	}
	return sgppGroup{Key: "product:prod-r", Label: "Course R", Members: []SectionSGPPRow{m1, m2}}
}

// TestSGPPGroupTableRow_ChipsTeachersAndActions is pin #4's happy path.
func TestSGPPGroupTableRow_ChipsTeachersAndActions(t *testing.T) {
	deps := groupingDeps()
	l := deps.Labels
	g := groupedFixture()

	row := sgppGroupTableRow(deps, "sec1", g, allSGPPPerms(), l)

	if row.ID != "c1" || row.DataAttrs["testid"] != "sgpp-row-c1" {
		t.Errorf("grouped row identity must be the lead member's class id, got ID=%q testid=%q", row.ID, row.DataAttrs["testid"])
	}
	if row.DataAttrs["status"] != "active" {
		t.Errorf("data-status = %q, want active", row.DataAttrs["status"])
	}

	subject := string(row.Cells[0].HTML)
	if !strings.Contains(subject, ">Course R<") {
		t.Errorf("subject cell missing the product label: %q", subject)
	}
	iOne := strings.Index(subject, `data-testid="sgpp-variant-c1">One<`)
	iTwo := strings.Index(subject, `data-testid="sgpp-variant-c2">Two<`)
	if iOne < 0 || iTwo < 0 || iOne > iTwo {
		t.Errorf("variant chips missing or out of deterministic order: %q", subject)
	}
	if row.Cells[0].Value != "Course R" {
		t.Errorf("subject cell Value = %q, want the product label (CSV export)", row.Cells[0].Value)
	}

	// Merged Teachers text flows through the ONE joiner: phase order across
	// members, lyngua separator, aria template intact (sgppAssignmentsCell).
	teachers := row.Cells[1]
	wantText := "Bea Cruz (Phase 1), Alma Reyes (Phase 2)"
	if teachers.Value != wantText {
		t.Errorf("merged Teachers text = %q, want %q", teachers.Value, wantText)
	}
	if html := string(teachers.HTML); !strings.Contains(html, `data-testid="sgpp-row-assignments-c1"`) {
		t.Errorf("teachers cell missing the lead-keyed testid: %q", html)
	}

	if len(row.Actions) != 3 {
		t.Fatalf("grouped row actions = %d, want 2 View + 1 Remove: %+v", len(row.Actions), row.Actions)
	}
	v1, v2, remove := row.Actions[0], row.Actions[1], row.Actions[2]
	wantLabel := func(variant string) string {
		s := strings.ReplaceAll(l.Staff.ActionVariant, "{{action}}", l.Staff.ViewAction)
		return strings.ReplaceAll(s, "{{variant}}", variant)
	}
	if v1.Label != wantLabel("One") || v1.Href != "/detail/sec1/c1" || v1.TestID != "sgpp-view-c1" || v1.Disabled {
		t.Errorf("first View action wrong: %+v", v1)
	}
	if v2.Label != wantLabel("Two") || v2.Href != "/detail/sec1/c2" || v2.TestID != "sgpp-view-c2" || v2.Disabled {
		t.Errorf("second View action wrong: %+v", v2)
	}
	if remove.Action != "delete" || !remove.Overflow || remove.Disabled {
		t.Errorf("Remove action wrong: %+v", remove)
	}
	if remove.URL != "/bulk-delete?id=c2" {
		t.Errorf("Remove URL = %q, want the bulk route with the NON-lead ids baked (lead rides data-id)", remove.URL)
	}
	if remove.TestID != "sgpp-remove-c1" || remove.ItemName != "Course R" {
		t.Errorf("Remove identity wrong: %+v", remove)
	}
	if want := fmt.Sprintf(l.Staff.RemoveConfirmMessage, "Course R"); remove.ConfirmMessage != want {
		t.Errorf("Remove confirm = %q, want %q", remove.ConfirmMessage, want)
	}
}

// TestSGPPGroupTableRow_PerVariantExclusionSignalling: one EXCLUDED strand
// renders the VariantExcludedTag chip while the row stays active; all-excluded
// flips the row's data-status to excluded.
func TestSGPPGroupTableRow_PerVariantExclusionSignalling(t *testing.T) {
	deps := groupingDeps()
	l := deps.Labels

	g := groupedFixture()
	g.Members[1].Excluded = true
	row := sgppGroupTableRow(deps, "sec1", g, allSGPPPerms(), l)

	if row.DataAttrs["status"] != "active" {
		t.Errorf("mixed group data-status = %q, want active (the exception is per-variant)", row.DataAttrs["status"])
	}
	subject := string(row.Cells[0].HTML)
	wantChip := strings.ReplaceAll(l.Staff.VariantExcludedTag, "{{variant}}", "Two")
	if !strings.Contains(subject, `data-testid="sgpp-excluded-tag-c2"`) || !strings.Contains(subject, ">"+template2html(wantChip)+"<") {
		t.Errorf("excluded strand chip missing (want %q): %q", wantChip, subject)
	}
	if !strings.Contains(subject, `data-testid="sgpp-variant-c1">One<`) {
		t.Errorf("non-excluded strand chip lost: %q", subject)
	}

	g2 := groupedFixture()
	g2.Members[0].Excluded = true
	g2.Members[1].Excluded = true
	row2 := sgppGroupTableRow(deps, "sec1", g2, allSGPPPerms(), l)
	if row2.DataAttrs["status"] != "excluded" {
		t.Errorf("all-excluded group data-status = %q, want excluded", row2.DataAttrs["status"])
	}
}

// template2html mirrors the HTML escaping the cell builder applies, so label
// assertions survive labels containing markup-significant runes (the default
// VariantExcludedTag carries an em-dash, which escapes to itself).
func template2html(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&#34;", "'", "&#39;")
	return r.Replace(s)
}

// TestSGPPGroupTableRow_RemoveGates: the whole-product Remove disables on any
// in-use member, on a missing delete grant, and on an unwired bulk route.
func TestSGPPGroupTableRow_RemoveGates(t *testing.T) {
	l := subscription_group.DefaultLabels()

	t.Run("any member in use blocks Remove", func(t *testing.T) {
		deps := groupingDeps()
		g := groupedFixture()
		g.Members[1].InUse = true
		row := sgppGroupTableRow(deps, "sec1", g, allSGPPPerms(), l)
		remove := row.Actions[len(row.Actions)-1]
		if !remove.Disabled || remove.DisabledTooltip != l.Staff.RemoveBlockedInUse {
			t.Errorf("in-use group Remove not blocked: %+v", remove)
		}
	})

	t.Run("no delete grant disables Remove", func(t *testing.T) {
		deps := groupingDeps()
		row := sgppGroupTableRow(deps, "sec1", groupedFixture(), groupingPerms(sgppEntity+":list"), l)
		remove := row.Actions[len(row.Actions)-1]
		if !remove.Disabled || remove.DisabledTooltip != l.Staff.Unauthorized {
			t.Errorf("permission-less group Remove not disabled: %+v", remove)
		}
	})

	t.Run("unwired bulk route disables Remove", func(t *testing.T) {
		deps := groupingDeps()
		deps.SGPPBulkDeleteURL = ""
		row := sgppGroupTableRow(deps, "sec1", groupedFixture(), allSGPPPerms(), l)
		remove := row.Actions[len(row.Actions)-1]
		if !remove.Disabled || remove.URL != "" {
			t.Errorf("unwired bulk route left Remove enabled: %+v", remove)
		}
	})
}

// TestSGPPGroupRemoveURL: every non-lead member id rides the query exactly
// once, URL-escaped; the lead is deliberately absent (data-id supplies it).
func TestSGPPGroupRemoveURL(t *testing.T) {
	members := []SectionSGPPRow{{ID: "lead"}, {ID: "m 2"}, {ID: "m3"}}
	if got, want := sgppGroupRemoveURL("/bulk-delete", members), "/bulk-delete?id=m+2&id=m3"; got != want {
		t.Errorf("sgppGroupRemoveURL = %q, want %q", got, want)
	}
	if got := sgppGroupRemoveURL("", members); got != "" {
		t.Errorf("unwired base must return \"\", got %q", got)
	}
	if got := sgppGroupRemoveURL("/bulk-delete", members[:1]); got != "" {
		t.Errorf("a single-member group never builds a bulk URL, got %q", got)
	}
}

// TestSGPPVariantActionLabel: the ActionVariant template fills both
// placeholders; an un-varianted member falls back to its offering label; an
// empty template degrades to the bare action label.
func TestSGPPVariantActionLabel(t *testing.T) {
	l := subscription_group.DefaultLabels()
	m := SectionSGPPRow{VariantLabel: "One", OfferingLabel: "Course R: One"}
	want := strings.ReplaceAll(strings.ReplaceAll(l.Staff.ActionVariant, "{{action}}", l.Staff.ViewAction), "{{variant}}", "One")
	if got := sgppVariantActionLabel(l, l.Staff.ViewAction, m); got != want {
		t.Errorf("variant action label = %q, want %q", got, want)
	}
	m2 := SectionSGPPRow{OfferingLabel: "Course R: Two"}
	want2 := strings.ReplaceAll(strings.ReplaceAll(l.Staff.ActionVariant, "{{action}}", l.Staff.ViewAction), "{{variant}}", "Course R: Two")
	if got := sgppVariantActionLabel(l, l.Staff.ViewAction, m2); got != want2 {
		t.Errorf("un-varianted member label = %q, want offering-label fallback %q", got, want2)
	}
	l.Staff.ActionVariant = ""
	if got := sgppVariantActionLabel(l, l.Staff.ViewAction, m); got != l.Staff.ViewAction {
		t.Errorf("empty template must degrade to the bare action label, got %q", got)
	}
}
