package form

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Data is the template data for the collection drawer form.
type Data struct {
	FormAction       string
	IsEdit           bool
	ID               string
	Customer         string
	ReferenceNumber  string
	Amount           string
	Currency         string
	CollectionMethod string
	Date             string
	ReceivedBy       string
	ReceivedRole     string
	Notes            string
	CollectionType   string
	Status           string
	Labels           centymo.CollectionFormLabels
	CommonLabels     any
}
