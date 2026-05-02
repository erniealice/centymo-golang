package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the procurement request drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Section 1 — Identity
	RequestNumber   string
	RequesterUserID string
	SupplierID      string
	LocationID      string

	// Section 2 — Financial
	Currency             string
	EstimatedTotalAmount string

	// Section 3 — Timing & Approval
	NeededByDate string
	ApprovedBy   string
	ApprovedAt   string
	Status       string

	// Section 4 — Others
	Justification string
	Notes         string
	Active        bool

	// Dropdown options
	Suppliers     []types.SelectOption
	StatusOptions []types.SelectOption

	Labels       centymo.ProcurementRequestFormLabels
	CommonLabels pyeza.CommonLabels
}
