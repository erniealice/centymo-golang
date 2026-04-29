package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// BuildDurationUnitOptions returns the select options for billing_cycle_unit
// and default_term_unit, sourced from the shared CommonLabels.DurationUnit
// translations (day(s) / week(s) / month(s) / year(s)).
func BuildDurationUnitOptions(cl pyeza.CommonLabels) []types.SelectOption {
	du := cl.DurationUnit
	return []types.SelectOption{
		{Value: "day", Label: du.DaySelect},
		{Value: "week", Label: du.WeekSelect},
		{Value: "month", Label: du.MonthSelect},
		{Value: "year", Label: du.YearSelect},
	}
}

// BuildBillingKindOptions builds the select options for the BillingKind enum
// on the PricePlan form. Values match the proto enum string representation
// returned by BillingKind.String() — e.g. "BILLING_KIND_ONE_TIME".
//
// 2026-04-29 milestone-billing plan §2.2 / Phase D — adds MILESTONE option.
// When MILESTONE is selected, the drawer's JS clears + disables the cycle
// inputs (cycle has no meaning for milestones); the use case defends against
// stale values by coercing to nil server-side.
func BuildBillingKindOptions(labels centymo.PricePlanFormLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "BILLING_KIND_ONE_TIME", Label: labels.BillingKindOneTime},
		{Value: "BILLING_KIND_RECURRING", Label: labels.BillingKindRecurring},
		{Value: "BILLING_KIND_CONTRACT", Label: labels.BillingKindContract},
		{Value: "BILLING_KIND_MILESTONE", Label: labels.BillingKindMilestone},
	}
}

// BuildAmountBasisOptions builds the select options for the AmountBasis enum
// on the PricePlan form. Values match the proto enum string representation
// returned by AmountBasis.String() — e.g. "AMOUNT_BASIS_PER_CYCLE".
func BuildAmountBasisOptions(labels centymo.PricePlanFormLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "AMOUNT_BASIS_PER_CYCLE", Label: labels.AmountBasisPerCycle},
		{Value: "AMOUNT_BASIS_TOTAL_PACKAGE", Label: labels.AmountBasisTotalPackage},
		{Value: "AMOUNT_BASIS_DERIVED_FROM_LINES", Label: labels.AmountBasisDerivedFromLines},
	}
}

// BuildBillingTreatmentOptions builds the select options for the BillingTreatment
// enum on the ProductPricePlan form. Values match the proto enum string representation
// returned by BillingTreatment.String() — e.g. "BILLING_TREATMENT_RECURRING".
func BuildBillingTreatmentOptions(labels centymo.ProductPricePlanFormLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "BILLING_TREATMENT_RECURRING", Label: labels.BillingTreatmentRecurring},
		{Value: "BILLING_TREATMENT_ONE_TIME_INITIAL", Label: labels.BillingTreatmentOneTimeInitial},
		{Value: "BILLING_TREATMENT_USAGE_BASED", Label: labels.BillingTreatmentUsageBased},
	}
}
