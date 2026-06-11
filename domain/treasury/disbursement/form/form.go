package form

import (
	disbursement "github.com/erniealice/centymo-golang/domain/treasury/disbursement"
)

// ExpenditureOption is a minimal struct for rendering expenditure (bill) options in the form.
type ExpenditureOption struct {
	Id     string
	Name   string
	Amount string
}

// Data is the template data for the disbursement drawer form.
type Data struct {
	FormAction       string
	WorkspaceID      string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit           bool
	ID               string
	ReferenceNumber  string
	Payee            string
	Amount           string
	Currency         string
	Method           string
	Date             string
	ApprovedBy       string
	ApprovedRole     string
	Notes            string
	DisbursementType string
	ExpenditureID    string
	Expenditures     []*ExpenditureOption
	Status           string
	Labels           disbursement.FormLabels
	CommonLabels     any
}
