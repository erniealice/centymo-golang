// Package form owns the template data shape for the pay drawer form
// (expense-pay-drawer-form, defined in detail/templates/detail.html).
// Pure types only — no Deps, no context.Context, no repository imports.
package form

import "github.com/erniealice/centymo-golang/domain/treasury"

// Data is the template data for the pay drawer form.
type Data struct {
	FormAction       string
	WorkspaceID      string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	ExpenditureID    string
	Name             string
	Amount           string
	Currency         string
	DisbursementType string
	Labels           treasury.DisbursementFormLabels
	CommonLabels     any
}
