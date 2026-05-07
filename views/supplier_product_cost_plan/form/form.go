package form

// Labels holds all UI labels for the supplier_product_cost_plan inline drawer form.
type Labels struct {
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

// Data holds the form state for the supplier_product_cost_plan inline drawer form.
type Data struct {
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

	Labels Labels
}
