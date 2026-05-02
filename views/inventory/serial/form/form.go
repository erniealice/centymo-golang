// Package form owns the template data shape for the serial drawer
// (serial-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// Labels holds i18n labels for the serial drawer form.
type Labels struct {
	SerialNumber  string
	IMEI          string
	Status        string
	WarrantyStart string
	WarrantyEnd   string
	PurchaseOrder string
	SoldReference string

	// Field-level info text surfaced via an info button beside each label.
	SerialNumberInfo  string
	IMEIInfo          string
	StatusInfo        string
	WarrantyStartInfo string
	WarrantyEndInfo   string
	PurchaseOrderInfo string
	SoldReferenceInfo string
}

// Data is the template data for the serial drawer form.
type Data struct {
	FormAction    string
	IsEdit        bool
	ID            string
	SerialNumber  string
	IMEI          string
	Status        string
	WarrantyStart string
	WarrantyEnd   string
	PurchaseOrder string
	SoldReference string
	Labels        Labels
	StatusOptions []pyeza.SelectOption
	CommonLabels  any
}
