package form

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Data is the template data for the recognition line drawer form.
type Data struct {
	FormAction           string
	IsEdit               bool
	ID                   string
	ExpenseRecognitionID string
	Description          string
	Quantity             string
	UnitAmount           string
	Amount               string
	Currency             string
	CommonLabels         any
	Labels               centymo.ExpenseRecognitionLineLabels
}
