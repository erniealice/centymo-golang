// Package form owns the template data shape for the supplier-subscription-drawer-form.html
// template. Pure types only — no Deps, no context.Context, no repository imports.
package form

// Labels holds i18n labels for the supplier subscription drawer form template.
type Labels struct {
	// Section headings (5-section layout)
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string

	// Field labels
	Name                string
	NamePlaceholder     string
	Code                string
	CodePlaceholder     string
	CostPlan            string
	CostPlanPlaceholder string
	CostPlanSearch      string
	CostPlanNoResults   string
	Supplier            string
	SupplierPlaceholder string
	SupplierSearch      string
	SupplierNoResults   string
	AutoRenew           string
	Active              string
	StartDate           string
	StartTime           string
	EndDate             string
	EndTime             string
	TimePlaceholder     string
	Notes               string
	NotesPlaceholder    string
}

// Data is the template data for the supplier subscription drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Field values
	Name          string
	Code          string
	CostPlanID    string
	CostPlanLabel string
	SupplierID    string
	SupplierLabel string
	AutoRenew     bool
	Active        bool

	// Date/Time form values split for two-row grid.
	DateStartDate string
	DateStartTime string
	DateEndDate   string
	DateEndTime   string
	DateStartISO  string
	DateEndISO    string
	DefaultTZ     string

	Notes string

	// Search endpoints
	SearchCostPlanURL string
	SearchSupplierURL string

	Labels       Labels
	CommonLabels any
}
