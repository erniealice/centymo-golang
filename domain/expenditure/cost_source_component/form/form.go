// Package form holds the drawer template data of the recoverable cost line
// add/edit drawer.
package form

import (
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	pyeza "github.com/erniealice/pyeza-golang"
)

// Option is one <select> option.
type Option struct {
	Value       string
	Label       string
	Description string // read by the pyeza form-group select (optional)
	Selected    bool
}

// Data is the template data of "cost-source-component-drawer-form".
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for the action workspace guard
	IsEdit      bool
	ID          string

	ExpenditureID string
	Currency      string

	Description string
	BasisUnit   string
	Quantity    string
	Amount      string
	ServiceFrom string
	ServiceTo   string

	Kinds     []Option
	TaxFacts  []Option
	BillLines []Option // empty hides the picker

	// AmountLabel is the amount label with the currency code; SubmitLabel the drawer's submit text.
	AmountLabel string
	SubmitLabel string

	Labels       csc.Labels
	CommonLabels pyeza.CommonLabels
}
