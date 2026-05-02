package form

import (
	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the settlement drawer form.
//
// Used by both the inline Add/Edit drawer and the parent-level "Settle"
// drawer (the latter posts to AccruedExpenseSettleURL — see
// detail.NewSettleAction). The form fields are identical; only the action
// URL differs.
type Data struct {
	FormAction       string
	IsEdit           bool
	ID               string
	AccruedExpenseID string
	ExpenditureID    string
	AmountSettled    string
	Currency         string
	FxRate           string
	ReversalReason   string
	Expenditures     []types.SelectOption
	CommonLabels     any
	Labels           centymo.AccruedExpenseSettlementLabels
}
