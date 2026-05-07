// Package form owns the template data shape for the cost-schedule-drawer-form.html template.
package form

// Labels holds i18n labels for the cost_schedule drawer form template.
type Labels struct {
	SectionIdentification string
	SectionRelationships  string
	SectionConfiguration  string
	SectionSchedule       string
	SectionNotes          string
	Name                  string
	NamePlaceholder       string
	Description           string
	DescPlaceholder       string
	StartDate             string
	EndDate               string
	Location              string
	LocationPlaceholder   string
	Active                string
}

// Data is the template data for the cost_schedule drawer form.
type Data struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Name        string
	Description string
	StartDate   string
	EndDate     string
	LocationID  string
	Active      bool
	Labels      Labels
	CommonLabels any
}
