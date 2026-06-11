// Package form owns the template data shape for the PO receipt form
// (po-receipt-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

import "github.com/erniealice/centymo-golang/domain/expenditure"

// LineRow is a single line item row shown in the receipt form.
type LineRow struct {
	ID               string
	Description      string
	LineType         string
	QuantityOrdered  string
	QuantityReceived string
	Remaining        string
}

// Data is the template data for the PO receipt form.
type Data struct {
	FormAction      string
	WorkspaceID     string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	PurchaseOrderID string
	Lines           []LineRow
	Today           string
	LocationID      string
	Labels          expenditure.ExpenditureLabels
	CommonLabels    any
}
