package form

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// BuildDurationUnitOptions returns the select options for billing_cycle_unit
// and default_term_unit, sourced from the shared CommonLabels.DurationUnit
// translations (day(s) / week(s) / month(s) / year(s)).
//
// 2026-04-30 enum-select-canonicalize — duration_unit is *not* a proto enum
// (it's a plain string column with the values "day" / "week" / "month" /
// "year"), so this helper is allowed to stay. The drift sweep targets the
// proto-backed enums (BillingKind, AmountBasis, BillingTreatment) whose
// option lists used to live in Go option-builders that drifted from the
// proto. Those now render inline in the drawer templates with a checked-in
// drift test; see docs/plan/20260430-enum-select-canonicalize/plan.md.
func BuildDurationUnitOptions(cl pyeza.CommonLabels) []types.SelectOption {
	du := cl.DurationUnit
	return []types.SelectOption{
		{Value: "day", Label: du.DaySelect},
		{Value: "week", Label: du.WeekSelect},
		{Value: "month", Label: du.MonthSelect},
		{Value: "year", Label: du.YearSelect},
	}
}
