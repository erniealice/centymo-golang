// Package form owns the template data shape for the supplier-plan-drawer-form.html template.
package form

// Labels holds i18n labels for the supplier_plan drawer form template.
type Labels struct {
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string
	Name                  string
	NamePlaceholder       string
	Supplier              string
	SupplierPlaceholder   string
	Active                string
}

// Data is the template data for the supplier_plan drawer form.
type Data struct {
	FormAction        string
	IsEdit            bool
	ID                string
	Name              string
	SupplierID        string
	SupplierLabel     string
	Active            bool
	SearchSupplierURL string
	Labels            Labels
	CommonLabels      any
}
