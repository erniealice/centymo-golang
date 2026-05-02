package form

import (
	centymo "github.com/erniealice/centymo-golang"
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
	Labels           centymo.DisbursementFormLabels
	CommonLabels     any
}
