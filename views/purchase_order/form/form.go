// Package form owns the template data shape for the primary purchase order
// drawer (purchase-order-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

// Labels holds flat i18n labels for the purchase order drawer form template.
type Labels struct {
	PoNumber         string
	SupplierID       string
	PoType           string
	OrderDate        string
	Currency         string
	PaymentTerms     string
	Notes            string
	NotesPlaceholder string
	Status           string

	// Field-level info text surfaced via an info button beside each label.
	PoNumberInfo     string
	SupplierIDInfo   string
	PoTypeInfo       string
	OrderDateInfo    string
	CurrencyInfo     string
	PaymentTermsInfo string
	NotesInfo        string
}

// Data is the template data for the purchase order drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	PoNumber     string
	SupplierID   string
	PoType       string
	OrderDate    string
	Currency     string
	PaymentTerms string
	Notes        string
	Labels       Labels
	CommonLabels any
}
