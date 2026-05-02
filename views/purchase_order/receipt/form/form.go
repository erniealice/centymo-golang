// Package form owns the template data shape for the PO receipt form
// (po-receipt-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

import centymo "github.com/erniealice/centymo-golang"

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
	PurchaseOrderID string
	Lines           []LineRow
	Today           string
	LocationID      string
	Labels          centymo.ExpenditureLabels
	CommonLabels    any
}
