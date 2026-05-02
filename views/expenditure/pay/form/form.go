// Package form owns the template data shape for the pay drawer form
// (expense-pay-drawer-form, defined in detail/templates/detail.html).
// Pure types only — no Deps, no context.Context, no repository imports.
package form

import centymo "github.com/erniealice/centymo-golang"

// Data is the template data for the pay drawer form.
type Data struct {
	FormAction       string
	ExpenditureID    string
	Name             string
	Amount           string
	Currency         string
	DisbursementType string
	Labels           centymo.DisbursementFormLabels
	CommonLabels     any
}
