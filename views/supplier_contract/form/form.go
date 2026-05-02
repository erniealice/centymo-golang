package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// KindOption is a select option for the contract kind enum.
type KindOption struct {
	Value    string
	Label    string
	Selected bool
}

// Data is the template data for the supplier contract drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Section 1 — Company / Identity Details
	Name            string
	ReferenceNumber string
	Kind            string
	KindOptions     []KindOption
	SupplierID      string
	SupplierName    string

	// Section 2 — Validity & Recurrence
	StartDate         string
	EndDate           string
	BillingCycleValue string
	BillingCycleUnit  string
	AutoRenew         bool
	RenewalNoticeDays string

	// Section 3 — Money & Approval
	Currency        string
	CommittedAmount string
	CycleAmount     string
	PaymentTermID   string
	ApprovedBy      string
	ApprovedDate    string
	RequestedBy     string

	// Section 4 — Categorization
	ExpenditureCategoryID string
	ExpenseAccountID      string
	LocationID            string

	// Section 5 — Others
	Notes  string
	Active bool
	Status string

	// Dropdown options
	Suppliers     []types.SelectOption
	PaymentTerms  []types.SelectOption
	StatusOptions []types.SelectOption

	Labels       centymo.SupplierContractFormLabels
	CommonLabels pyeza.CommonLabels
}
