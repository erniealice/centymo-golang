// Package form owns the template data shape for the payment drawer
// (revenue-payment-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import (
	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the payment drawer form.
type Data struct {
	FormAction         string
	WorkspaceID        string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit             bool
	ID                 string
	RevenueID          string
	CollectionMethodID string
	AmountPaid         string
	Currency           string
	ReferenceNumber    string
	Notes              string
	ReceivedBy         string
	ReceivedRole       string
	PaymentMethods     []pyeza.SelectOption
	CommonLabels       any
	Labels             revenuedomain.Labels
}
