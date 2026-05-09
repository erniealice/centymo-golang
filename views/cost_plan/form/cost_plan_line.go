package form

// CostPlanLineLabels holds all UI labels for the SupplierProductCostPlan
// inline drawer form (rendered within the CostPlan Lines tab).
type CostPlanLineLabels struct {
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string

	SupplierProductPlan            string
	SupplierProductPlanPlaceholder string
	BillingTreatment               string
	Amount                         string
	AmountPlaceholder              string
	Active                         string

	// BillingTreatment option labels
	TreatmentRecurring         string
	TreatmentOneTimeInitial    string
	TreatmentUsageBased        string
	TreatmentMinimumCommitment string
}

// CostPlanLineData holds the form state for the SupplierProductCostPlan
// inline drawer form.
type CostPlanLineData struct {
	FormAction string
	IsEdit     bool
	ID         string

	// The parent CostPlan ID — used to scope the form action URL.
	CostPlanID string

	SupplierProductPlanID    string
	SupplierProductPlanLabel string
	BillingTreatment         string
	Amount                   string
	Active                   bool

	SearchSupplierProductPlanURL string

	Labels       CostPlanLineLabels
	CommonLabels any
}
