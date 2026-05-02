package form

import (
	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the accrued_expense drawer form.
//
// The manual create drawer is the SECONDARY path — the primary path is
// AccrueFromContract (the recurrence engine) — but a manual form is needed
// for one-off accruals (e.g., a utility estimate where no contract exists).
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	// §1 Identity
	Name        string
	Description string

	// §2 Source
	SupplierContractID string
	SupplierID         string

	// §3 Period
	RecognitionDate string
	PeriodStart     string
	PeriodEnd       string
	CycleDate       string

	// §4 Money
	Currency        string
	AccruedAmount   string
	SettledAmount   string
	RemainingAmount string

	// §5 Lifecycle
	Status string

	// §6 Accounting
	ExpenseAccountID string
	AccrualAccountID string

	// §7 Notes
	Notes  string
	Active bool

	// Dropdown options
	SupplierContracts []types.SelectOption
	Suppliers         []types.SelectOption
	StatusOptions     []types.SelectOption

	Labels       centymo.AccruedExpenseFormLabels
	CommonLabels any
}
