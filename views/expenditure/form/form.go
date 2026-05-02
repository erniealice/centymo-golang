// Package form owns the template data shape for the primary expense drawer
// (expense-drawer-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id      string
	Name    string
	NetDays int32
}

// PurchaseOrderOption is a minimal struct for rendering purchase order options in the form.
type PurchaseOrderOption struct {
	Id           string
	PoNumber     string
	SupplierName string
}

// Labels holds flat i18n labels for the expense drawer form template.
type Labels struct {
	Name                string
	NamePlaceholder     string
	Category            string
	Supplier            string
	Date                string
	Amount              string
	Currency            string
	ReferenceNumber     string
	Notes               string
	NotesPlaceholder    string
	Status              string
	ExpenditureType     string
	TypeExpense         string
	TypePurchase        string
	StatusPending       string
	StatusApproved      string
	StatusPaid          string
	StatusCancelled     string
	CurrencyPlaceholder string
	PaymentTerms        string
	SelectPaymentTerm   string
	DueDate             string
	LinkToPurchaseOrder string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo            string
	ExpenditureTypeInfo string
	CategoryInfo        string
	DateInfo            string
	AmountInfo          string
	CurrencyInfo        string
	ReferenceNumberInfo string
	SupplierInfo        string
	NotesInfo           string
}

// Data is the template data for the expense drawer form.
type Data struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	ExpenditureType       string
	ExpenditureCategoryID string
	SupplierID            string
	Date                  string
	TotalAmount           string
	Currency              string
	Status                string
	ReferenceNumber       string
	Notes                 string
	Categories            []pyeza.SelectOption
	Suppliers             []pyeza.SelectOption
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	PurchaseOrders        []*PurchaseOrderOption
	PurchaseOrderID       string
	Labels                Labels
	CommonLabels          any
}
