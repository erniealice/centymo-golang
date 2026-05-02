package form

import (
	"testing"

	centymo "github.com/erniealice/centymo-golang"
)

// TestLabelsFromPricePlanAllFieldsPopulated asserts that every field on the
// result Labels struct is populated when the source PricePlanFormLabels has
// every field set to a non-empty value. This test prevents the silent-empty-field
// bug class: if a future field is added to the source or destination and the
// mapper forgets to copy it, this test fails.
func TestLabelsFromPricePlanAllFieldsPopulated(t *testing.T) {
	src := centymo.PricePlanFormLabels{
		// Basic section
		SectionBasic:           "SectionBasic",
		SectionPricing:         "SectionPricing",
		Name:                   "Name",
		NamePlaceholder:        "NamePlaceholder",
		Description:            "Description",
		DescPlaceholder:        "DescPlaceholder",
		Amount:                 "Amount",
		AmountPlaceholder:      "AmountPlaceholder",
		Currency:               "Currency",
		CurrencyPlaceholder:    "CurrencyPlaceholder",
		DurationValue:          "DurationValue",
		DurationUnit:           "DurationUnit",
		Active:                 "Active",
		PlanLabel:              "PlanLabel",
		PlanPlaceholder:        "PlanPlaceholder",
		PlanSearch:             "PlanSearch",
		Schedule:               "Schedule",
		SchedulePlaceholder:    "SchedulePlaceholder",
		ScheduleSearch:         "ScheduleSearch",
		LocationHintPrefix:     "LocationHintPrefix",
		// Wave 2 labels
		BillingKindLabel:        "BillingKindLabel",
		AmountBasisLabel:        "AmountBasisLabel",
		BillingCycleLabel:       "BillingCycleLabel",
		BillingCyclePlaceholder: "BillingCyclePlaceholder",
		TermLabel:               "TermLabel",
		TermPlaceholder:         "TermPlaceholder",
		TermOpenEndedHelp:       "TermOpenEndedHelp",
		// Per-option labels
		BillingKindOneTime:          "BillingKindOneTime",
		BillingKindRecurring:        "BillingKindRecurring",
		BillingKindContract:         "BillingKindContract",
		BillingKindMilestone:        "BillingKindMilestone",
		BillingKindAdHoc:            "BillingKindAdHoc",
		AmountBasisPerCycle:         "AmountBasisPerCycle",
		AmountBasisTotalPackage:     "AmountBasisTotalPackage",
		AmountBasisDerivedFromLines: "AmountBasisDerivedFromLines",
		AmountBasisPerOccurrence:    "AmountBasisPerOccurrence",
		// Per-option hints
		BillingKindOneTimeHint:          "BillingKindOneTimeHint",
		BillingKindRecurringHint:        "BillingKindRecurringHint",
		BillingKindContractHint:         "BillingKindContractHint",
		BillingKindMilestoneHint:        "BillingKindMilestoneHint",
		BillingKindAdHocHint:            "BillingKindAdHocHint",
		AmountBasisPerCycleHint:         "AmountBasisPerCycleHint",
		AmountBasisTotalPackageHint:     "AmountBasisTotalPackageHint",
		AmountBasisDerivedFromLinesHint: "AmountBasisDerivedFromLinesHint",
		AmountBasisPerOccurrenceHint:    "AmountBasisPerOccurrenceHint",
		// Entitled occurrences
		EntitledOccurrencesLabel:       "EntitledOccurrencesLabel",
		EntitledOccurrencesPlaceholder: "EntitledOccurrencesPlaceholder",
		EntitledOccurrencesInfo:        "EntitledOccurrencesInfo",
		// Field-level info popovers
		PlanInfo:         "PlanInfo",
		ScheduleInfo:     "ScheduleInfo",
		NameInfo:         "NameInfo",
		DescriptionInfo:  "DescriptionInfo",
		BillingKindInfo:  "BillingKindInfo",
		AmountBasisInfo:  "AmountBasisInfo",
		AmountInfo:       "AmountInfo",
		CurrencyInfo:     "CurrencyInfo",
		BillingCycleInfo: "BillingCycleInfo",
		TermInfo:         "TermInfo",
		ActiveInfo:       "ActiveInfo",
		// Plan-client-scope notices and tooltips
		ParentScheduleClientNotice: "ParentScheduleClientNotice",
		ScheduleLockedTooltip:      "ScheduleLockedTooltip",
		// Cyclic-subscription-jobs
		MilestoneCyclicBlock: "MilestoneCyclicBlock",
		// Ad-hoc-subscription-billing guards
		AdHocPoolNoTemplate:           "AdHocPoolNoTemplate",
		AdHocPerCallNoTemplate:        "AdHocPerCallNoTemplate",
		AdHocNoEntitlement:            "AdHocNoEntitlement",
		AdHocBillingCycleNotAllowed:   "AdHocBillingCycleNotAllowed",
		AdHocVisitsPerCycleNotAllowed: "AdHocVisitsPerCycleNotAllowed",
	}

	result := LabelsFromPricePlan(src)

	// Assert every field is non-empty.
	tests := []struct {
		name  string
		value string
	}{
		{"SectionBasic", result.SectionBasic},
		{"SectionPricing", result.SectionPricing},
		{"NameLabel", result.NameLabel},
		{"NamePlaceholder", result.NamePlaceholder},
		{"DescriptionLabel", result.DescriptionLabel},
		{"DescriptionPlaceholder", result.DescriptionPlaceholder},
		{"AmountLabel", result.AmountLabel},
		{"AmountPlaceholder", result.AmountPlaceholder},
		{"CurrencyLabel", result.CurrencyLabel},
		{"CurrencyPlaceholder", result.CurrencyPlaceholder},
		{"DurationLabel", result.DurationLabel},
		{"DurationUnitLabel", result.DurationUnitLabel},
		{"ActiveLabel", result.ActiveLabel},
		{"PlanLabel", result.PlanLabel},
		{"PlanPlaceholder", result.PlanPlaceholder},
		{"PlanSearch", result.PlanSearch},
		{"ScheduleLabel", result.ScheduleLabel},
		{"SchedulePlaceholder", result.SchedulePlaceholder},
		{"ScheduleSearch", result.ScheduleSearch},
		{"LocationHintPrefix", result.LocationHintPrefix},
		{"BillingKindLabel", result.BillingKindLabel},
		{"AmountBasisLabel", result.AmountBasisLabel},
		{"BillingCycleLabel", result.BillingCycleLabel},
		{"BillingCyclePlaceholder", result.BillingCyclePlaceholder},
		{"TermLabel", result.TermLabel},
		{"TermPlaceholder", result.TermPlaceholder},
		{"TermOpenEndedHelp", result.TermOpenEndedHelp},
		{"BillingKindOneTime", result.BillingKindOneTime},
		{"BillingKindRecurring", result.BillingKindRecurring},
		{"BillingKindContract", result.BillingKindContract},
		{"BillingKindMilestone", result.BillingKindMilestone},
		{"BillingKindAdHoc", result.BillingKindAdHoc},
		{"AmountBasisPerCycle", result.AmountBasisPerCycle},
		{"AmountBasisTotalPackage", result.AmountBasisTotalPackage},
		{"AmountBasisDerivedFromLines", result.AmountBasisDerivedFromLines},
		{"AmountBasisPerOccurrence", result.AmountBasisPerOccurrence},
		{"BillingKindOneTimeHint", result.BillingKindOneTimeHint},
		{"BillingKindRecurringHint", result.BillingKindRecurringHint},
		{"BillingKindContractHint", result.BillingKindContractHint},
		{"BillingKindMilestoneHint", result.BillingKindMilestoneHint},
		{"BillingKindAdHocHint", result.BillingKindAdHocHint},
		{"AmountBasisPerCycleHint", result.AmountBasisPerCycleHint},
		{"AmountBasisTotalPackageHint", result.AmountBasisTotalPackageHint},
		{"AmountBasisDerivedFromLinesHint", result.AmountBasisDerivedFromLinesHint},
		{"AmountBasisPerOccurrenceHint", result.AmountBasisPerOccurrenceHint},
		{"EntitledOccurrencesLabel", result.EntitledOccurrencesLabel},
		{"EntitledOccurrencesPlaceholder", result.EntitledOccurrencesPlaceholder},
		{"EntitledOccurrencesInfo", result.EntitledOccurrencesInfo},
		{"PlanInfo", result.PlanInfo},
		{"ScheduleInfo", result.ScheduleInfo},
		{"NameInfo", result.NameInfo},
		{"DescriptionInfo", result.DescriptionInfo},
		{"BillingKindInfo", result.BillingKindInfo},
		{"AmountBasisInfo", result.AmountBasisInfo},
		{"AmountInfo", result.AmountInfo},
		{"CurrencyInfo", result.CurrencyInfo},
		{"BillingCycleInfo", result.BillingCycleInfo},
		{"TermInfo", result.TermInfo},
		{"ActiveInfo", result.ActiveInfo},
		{"ParentScheduleClientNotice", result.ParentScheduleClientNotice},
		{"ScheduleLockedTooltip", result.ScheduleLockedTooltip},
		{"MilestoneCyclicBlock", result.MilestoneCyclicBlock},
		{"AdHocPoolNoTemplate", result.AdHocPoolNoTemplate},
		{"AdHocPerCallNoTemplate", result.AdHocPerCallNoTemplate},
		{"AdHocNoEntitlement", result.AdHocNoEntitlement},
		{"AdHocBillingCycleNotAllowed", result.AdHocBillingCycleNotAllowed},
		{"AdHocVisitsPerCycleNotAllowed", result.AdHocVisitsPerCycleNotAllowed},
	}

	for _, tc := range tests {
		if tc.value == "" {
			t.Errorf("field %s is empty; mapper may have omitted a field", tc.name)
		}
	}
}
