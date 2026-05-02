// Package form owns the template data shape for the transaction drawer
// (transaction-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// Labels holds i18n labels for the transaction drawer form.
type Labels struct {
	Type      string
	Quantity  string
	Date      string
	Reference string

	// Field-level info text surfaced via an info button beside each label.
	TypeInfo      string
	QuantityInfo  string
	DateInfo      string
	ReferenceInfo string
}

// Data is the template data for the transaction drawer form.
type Data struct {
	FormAction   string
	Labels       Labels
	TypeOptions  []pyeza.SelectOption
	Today        string
	CommonLabels any
}
