// Package form owns the template data shape for the PO line item drawer
// (po-line-item-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import expenditureentity "github.com/erniealice/centymo-golang/domain/expenditure/expenditure"

// Data is the template data for the PO line item drawer form.
type Data struct {
	FormAction      string
	WorkspaceID     string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit          bool
	ID              string
	PurchaseOrderID string
	LineType        string
	Description     string
	ProductID       string
	InventoryItemID string
	LocationID      string
	QuantityOrdered string
	UnitPrice       string
	Notes           string
	Labels          expenditureentity.Labels
	CommonLabels    any
}
