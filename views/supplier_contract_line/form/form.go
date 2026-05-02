package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// TreatmentOption is a select option for the line treatment enum.
type TreatmentOption struct {
	Value    string
	Label    string
	Selected bool
}

// Data is the template data for the supplier contract line drawer form.
type Data struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	SupplierContractID string

	// Core fields
	Description      string
	LineType         string
	ProductID        string
	Quantity         string
	UnitPrice        string
	Treatment        string
	StartDate        string
	EndDate          string
	ExpenseAccountID string
	LineNumber       string

	// Options
	TreatmentOptions []TreatmentOption
	Products         []types.SelectOption

	Labels       centymo.SupplierContractLabels
	CommonLabels pyeza.CommonLabels
}
