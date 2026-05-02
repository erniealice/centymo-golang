package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the SCPSL drawer form.
type Data struct {
	FormAction                      string
	IsEdit                          bool
	ID                              string
	SupplierContractPriceScheduleID string

	// Fields
	SupplierContractLineID string
	Currency               string
	UnitPrice              string
	MinimumAmount          string
	Quantity               string
	CycleValueOverride     string
	CycleUnitOverride      string
	Notes                  string

	// Options
	ContractLines []types.SelectOption

	Labels       centymo.SupplierContractPriceScheduleLineFormLabels
	NounLabels   centymo.SupplierContractPriceScheduleLinesLabels
	CommonLabels pyeza.CommonLabels
}
