package billable_charge

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the billable charge pages. The
// unit is mounted by centymo's BillableChargeUnit only when the app opts in
// (block.WithRecoveryCharges); it is not part of the default centymo unit set.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.billable_charge",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "billable_charge"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "billable_charge.json", Key: "billable_charge"},
		LabelName: "BillableChargeLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "billable_charge:list",
			Items: []compose.NavItem{
				{Key: "billable-charges", Route: "billable_charge.list", Params: map[string]string{"status": "open"},
					Label: "Billable Charges", Icon: "icon-file-text", Permission: "billable_charge:list",
					LabelKey: "billable_charges_label", IconKey: "billable_charges_icon"},
			},
		},
	}
}
