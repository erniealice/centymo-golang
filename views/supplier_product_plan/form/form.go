package form

// Labels holds all UI labels for the supplier_product_plan drawer form.
type Labels struct {
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string

	SupplierPlan              string
	SupplierPlanPlaceholder   string
	Product                   string
	ProductPlaceholder        string
	ProductVariant            string
	ProductVariantPlaceholder string
	SupplierSKU               string
	SupplierSKUPlaceholder    string
	SupplierUnit              string
	SupplierUnitPlaceholder   string
	Active                    string
}

// Data holds the form state for the supplier_product_plan drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	SupplierPlanID    string
	SupplierPlanLabel string
	ProductID         string
	ProductLabel      string
	ProductVariantID  string
	SupplierSKU       string
	SupplierUnit      string
	Active            bool

	SearchSupplierPlanURL string
	SearchProductURL      string

	Labels Labels
}
