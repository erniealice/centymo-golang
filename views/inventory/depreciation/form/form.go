// Package form owns the template data shape for the depreciation drawer
// (depreciation-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// Labels holds i18n labels for the depreciation drawer form.
type Labels struct {
	Method       string
	CostBasis    string
	SalvageValue string
	UsefulLife   string
	StartDate    string

	// Field-level info text surfaced via an info button beside each label.
	MethodInfo       string
	CostBasisInfo    string
	SalvageValueInfo string
	UsefulLifeInfo   string
	StartDateInfo    string
}

// Data is the template data for the depreciation drawer form.
type Data struct {
	FormAction    string
	IsEdit        bool
	ID            string
	Method        string
	CostBasis     string
	SalvageValue  string
	UsefulLife    string
	StartDate     string
	Labels        Labels
	MethodOptions []pyeza.SelectOption
	CommonLabels  any
}
