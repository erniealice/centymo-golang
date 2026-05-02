// Package form owns the template data shape for the payment drawer
// (revenue-payment-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import (
	"github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the payment drawer form.
type Data struct {
	FormAction         string
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
	Labels             centymo.RevenueLabels
}
