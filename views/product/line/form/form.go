// Package form owns the template data shape for the primary product line
// drawer (product-line-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import centymo "github.com/erniealice/centymo-golang"

// Data is the template data for the line drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Active       bool
	Labels       centymo.ProductLineFormLabels
	CommonLabels any
}
