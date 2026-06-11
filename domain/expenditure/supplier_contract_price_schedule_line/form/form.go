package form

import (
	sib_expenditure_supplier_contract_price_schedule "github.com/erniealice/centymo-golang/domain/expenditure/supplier_contract_price_schedule"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the SCPSL drawer form.
type Data struct {
	FormAction                      string
	WorkspaceID                     string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
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

	Labels       sib_expenditure_supplier_contract_price_schedule.LineFormLabels
	NounLabels   sib_expenditure_supplier_contract_price_schedule.LinesLabels
	CommonLabels pyeza.CommonLabels
}
