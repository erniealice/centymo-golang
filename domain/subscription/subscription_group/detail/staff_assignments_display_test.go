package detail

// staff_assignments_display_test.go — unit coverage for the I-1 headline
// display contract (plan.md increment I-1 / centymo.md §3, audit gap R-G1):
// the section Subjects tab's Teachers column renders ALL of a class's active
// assignments, comma-separated, each phase-scoped entry carrying its phase
// label in a parenthetical, ordered by phase order (phase-NULL first) then by
// staff name.
//
// Before this file the contract was asserted only by a live-DB Playwright
// case. A revert to "first active edge wins", a dropped parenthetical, a
// hardcoded ", " separator or a lost aria-label all passed the Go suite
// untouched. joinAssignments and sgppAssignmentsCell are pure string
// composition over a row set, so they are asserted directly here — no DB, no
// deps wiring.

import (
	"strings"
	"testing"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
)

// dispLabels returns DefaultLabels with the two display-contract strings
// overridden to values that CANNOT be confused with a hardcoded fallback, so
// every assertion below proves the string was sourced from lyngua labels
// rather than baked into the composer.
func dispLabels(sep, aria string) subscription_group.Labels {
	l := subscription_group.DefaultLabels()
	l.Staff.ListSeparator = sep
	l.Staff.ListAria = aria
	return l
}

// disp is a terse constructor for one resolved assignment entry.
func disp(name, phase string, order int32) sgppAssignmentDisplay {
	return sgppAssignmentDisplay{staffName: name, phaseLabel: phase, order: order}
}

// phaseNullOrder is the order value listSGPPAssignments stamps on an
// assignment with no job_template_phase_id — it must sort ahead of every real
// phase, including a phase whose phase_order is 0.
const phaseNullOrder = int32(-1)

// TestJoinAssignments is the I-1 display contract itself: every active
// assignment appears (nothing is hidden), phase-scoped entries carry a
// parenthetical, phase-NULL entries render bare, and the order is phase order
// first then case-insensitive staff name — never the input order.
func TestJoinAssignments(t *testing.T) {
	l := subscription_group.DefaultLabels() // ListSeparator = ", "

	tests := []struct {
		name    string
		entries []sgppAssignmentDisplay
		want    string
	}{
		{
			name:    "zero staff — the unstaffed state renders an empty string",
			entries: nil,
			want:    "",
		},
		{
			name:    "zero staff — an allocated but empty slice is also empty",
			entries: []sgppAssignmentDisplay{},
			want:    "",
		},
		{
			name:    "one staff, no phase — renders bare, no parenthetical",
			entries: []sgppAssignmentDisplay{disp("Angeline Daoang", "", phaseNullOrder)},
			want:    "Angeline Daoang",
		},
		{
			name:    "one staff with a phase — label in a parenthetical",
			entries: []sgppAssignmentDisplay{disp("Angeline Daoang", "Semester 1", 1)},
			want:    "Angeline Daoang (Semester 1)",
		},
		{
			name: "two co-teachers on the SAME phase — both listed, name-ordered, phase repeated",
			entries: []sgppAssignmentDisplay{
				disp("Zoe Mercado", "Semester 1", 1),
				disp("Adam Ilagan", "Semester 1", 1),
			},
			want: "Adam Ilagan (Semester 1), Zoe Mercado (Semester 1)",
		},
		{
			name: "several staff across different phases — ordered by phase order, not input order or name",
			entries: []sgppAssignmentDisplay{
				disp("Carmen Uy", "Quarter 3", 3),
				disp("Alma Reyes", "Quarter 1", 1),
				disp("Bea Cruz", "Quarter 2", 2),
			},
			want: "Alma Reyes (Quarter 1), Bea Cruz (Quarter 2), Carmen Uy (Quarter 3)",
		},
		{
			name: "phase order wins over name — a Z teacher in phase 1 precedes an A teacher in phase 2",
			entries: []sgppAssignmentDisplay{
				disp("Aaron Yap", "Quarter 2", 2),
				disp("Zenaida Uy", "Quarter 1", 1),
			},
			want: "Zenaida Uy (Quarter 1), Aaron Yap (Quarter 2)",
		},
		{
			name: "a phase-NULL (whole-year) assignment sorts ahead of every phased one, including phase_order 0",
			entries: []sgppAssignmentDisplay{
				disp("Bea Cruz", "Quarter 1", 0),
				disp("Alma Reyes", "", phaseNullOrder),
			},
			want: "Alma Reyes, Bea Cruz (Quarter 1)",
		},
		{
			name: "mixed roster — phase-NULL first, then phases in order, co-teachers name-ordered within a phase",
			entries: []sgppAssignmentDisplay{
				disp("Carmen Uy", "Semester 2", 2),
				disp("Bea Cruz", "Semester 1", 1),
				disp("Delia Ong", "Semester 1", 1),
				disp("Alma Reyes", "", phaseNullOrder),
			},
			want: "Alma Reyes, Bea Cruz (Semester 1), Delia Ong (Semester 1), Carmen Uy (Semester 2)",
		},
		{
			name: "same phase order — the name tiebreak is case-insensitive (a byte compare would invert this)",
			entries: []sgppAssignmentDisplay{
				disp("Bea Cruz", "Semester 1", 1),
				disp("alonzo Reyes", "Semester 1", 1),
			},
			want: "alonzo Reyes (Semester 1), Bea Cruz (Semester 1)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := joinAssignments(tc.entries, l)
			if got != tc.want {
				t.Errorf("joinAssignments() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestJoinAssignments_SeparatorComesFromLabels proves the separator is a
// lyngua label and not a hardcoded ", ": a tier that ships " · " must get
// " · ", and an unset label must fall back to the documented ", " default.
func TestJoinAssignments_SeparatorComesFromLabels(t *testing.T) {
	entries := []sgppAssignmentDisplay{
		disp("Alma Reyes", "", phaseNullOrder),
		disp("Bea Cruz", "Semester 1", 1),
	}

	tests := []struct {
		name string
		sep  string
		want string
	}{
		{
			name: "tier-supplied separator is used verbatim",
			sep:  " | ",
			want: "Alma Reyes | Bea Cruz (Semester 1)",
		},
		{
			name: "unset separator falls back to the documented default",
			sep:  "",
			want: "Alma Reyes, Bea Cruz (Semester 1)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Fresh copies per case: joinAssignments sorts its argument in place.
			in := append([]sgppAssignmentDisplay(nil), entries...)
			got := joinAssignments(in, dispLabels(tc.sep, "Staff: {{names}}"))
			if got != tc.want {
				t.Errorf("joinAssignments() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestSGPPAssignmentsCell_Materialized asserts the rendered Teachers cell for
// a class row that has assignments: the stable per-row test id, the aria-label
// built from the lyngua template (placeholder substituted, not left literal),
// the visible text, and the raw text on Value for CSV export.
func TestSGPPAssignmentsCell_Materialized(t *testing.T) {
	l := dispLabels(", ", "Teachers: {{names}}.")
	row := SectionSGPPRow{
		ID:              "class-1",
		Materialized:    true,
		AssignmentsText: "Alma Reyes, Bea Cruz (Semester 1)",
	}

	cell := sgppAssignmentsCell(row, l)
	html := string(cell.HTML)

	if cell.Type != "html" {
		t.Errorf("cell.Type = %q, want \"html\"", cell.Type)
	}
	if want := `data-testid="sgpp-row-assignments-class-1"`; !strings.Contains(html, want) {
		t.Errorf("cell HTML = %q, want it to contain %q", html, want)
	}
	if want := `aria-label="Teachers: Alma Reyes, Bea Cruz (Semester 1)."`; !strings.Contains(html, want) {
		t.Errorf("cell HTML = %q, want the aria-label built from l.Staff.ListAria: %q", html, want)
	}
	if strings.Contains(html, "{{names}}") {
		t.Errorf("cell HTML = %q, want the {{names}} placeholder substituted", html)
	}
	if !strings.Contains(html, ">"+row.AssignmentsText+"<") {
		t.Errorf("cell HTML = %q, want it to render the assignments text %q", html, row.AssignmentsText)
	}
	if cell.Value != row.AssignmentsText {
		t.Errorf("cell.Value = %q, want the raw assignments text %q (CSV export)", cell.Value, row.AssignmentsText)
	}
}

// TestSGPPAssignmentsCell_States covers the two non-list states: a
// materialized class with zero assignments must show the lyngua Unstaffed
// placeholder (NOT an empty cell, NOT an em-dash), and a transitional
// fallback row (no class row exists yet) must show the em-dash.
func TestSGPPAssignmentsCell_States(t *testing.T) {
	l := dispLabels(", ", "Teachers: {{names}}.")

	tests := []struct {
		name string
		row  SectionSGPPRow
		want string
	}{
		{
			name: "materialized class with zero assignments — the unstaffed label",
			row:  SectionSGPPRow{ID: "class-1", Materialized: true, AssignmentsText: ""},
			want: l.Staff.Unstaffed,
		},
		{
			name: "transitional fallback row — em-dash, no class row exists yet",
			row:  SectionSGPPRow{ProductPlanID: "pp1", Materialized: false},
			want: emDash,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cell := sgppAssignmentsCell(tc.row, l)
			if cell.Type != "text" {
				t.Errorf("cell.Type = %q, want \"text\" (a plain placeholder, no markup)", cell.Type)
			}
			if cell.Value != tc.want {
				t.Errorf("cell.Value = %q, want %q", cell.Value, tc.want)
			}
			if cell.HTML != "" {
				t.Errorf("cell.HTML = %q, want empty for a text cell", cell.HTML)
			}
		})
	}
}

// TestSGPPAssignmentsCell_EscapesStaffNames: the cell is emitted as raw HTML,
// so a person name carrying markup must be escaped in BOTH the visible text
// and the aria-label. Names are operator-entered data, not trusted labels.
func TestSGPPAssignmentsCell_EscapesStaffNames(t *testing.T) {
	l := dispLabels(", ", "Teachers: {{names}}.")
	row := SectionSGPPRow{
		ID:              "class-1",
		Materialized:    true,
		AssignmentsText: "<script>x</script> & Co (Semester 1)",
	}

	html := string(sgppAssignmentsCell(row, l).HTML)

	if strings.Contains(html, "<script>") {
		t.Errorf("cell HTML = %q, want the staff name escaped (no live <script>)", html)
	}
	if n := strings.Count(html, "&lt;script&gt;x&lt;/script&gt; &amp; Co"); n != 2 {
		t.Errorf("escaped name occurs %d times in %q, want 2 (aria-label + visible text)", n, html)
	}
}
