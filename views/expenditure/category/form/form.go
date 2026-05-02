// Package form owns the template data shape for the category drawer form
// (category-drawer-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

// Labels holds flat i18n labels for the category drawer form template.
type Labels struct {
	Code        string
	Name        string
	Description string

	// Field-level info text surfaced via an info button beside each label.
	CodeInfo        string
	NameInfo        string
	DescriptionInfo string
}

// Data is the template data for the category drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Code         string
	Name         string
	Description  string
	Labels       Labels
	CommonLabels any
}
