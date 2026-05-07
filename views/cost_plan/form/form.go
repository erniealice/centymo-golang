package form

// Labels holds all UI labels for the cost_plan drawer form.
type Labels struct {
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string

	Name                    string
	NamePlaceholder         string
	Description             string
	DescPlaceholder         string
	SupplierPlan            string
	SupplierPlanPlaceholder string
	CostSchedule            string
	CostSchedulePlaceholder string
	BillingKind             string
	AmountBasis             string
	Amount                  string
	AmountPlaceholder       string
	Currency                string
	CurrencyPlaceholder     string
	BillingCycle            string
	BillingCyclePlaceholder string
	DefaultTerm             string
	DefaultTermPlaceholder  string
	Active                  string

	// BillingKind option labels
	BillingKindOneTime    string
	BillingKindRecurring  string
	BillingKindContract   string
	BillingKindUsageBased string
	BillingKindAdHoc      string

	// AmountBasis option labels
	AmountBasisPerCycle         string
	AmountBasisTotalPackage     string
	AmountBasisDerivedFromLines string
	AmountBasisPerOccurrence    string

	// Duration unit option labels
	DurationUnitDay   string
	DurationUnitWeek  string
	DurationUnitMonth string
	DurationUnitYear  string
}

// Data holds the form state for the cost_plan drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	Name            string
	Description     string
	SupplierPlanID  string
	SupplierPlanLabel string
	CostScheduleID  string
	CostScheduleLabel string
	BillingKind     string
	AmountBasis     string
	Amount          string
	Currency        string
	BillingCycleValue string
	BillingCycleUnit  string
	DefaultTermValue  string
	DefaultTermUnit   string
	Active          bool

	SearchSupplierPlanURL  string
	SearchCostScheduleURL  string

	Labels Labels
}
