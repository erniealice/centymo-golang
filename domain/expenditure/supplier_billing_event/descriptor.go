package supplier_billing_event

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the supplier_billing_event
// entity. This entity has no Routes struct (only URL constants in routes.go), so
// Routes and RouteJSON are omitted.
func Describe() compose.Unit {
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.supplier_billing_event",
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "advances_dashboard.json", Key: "supplierBillingEvent"},
		LabelName: "SupplierBillingEventLabels",
		Templates: TemplatesFS,
	}
}
