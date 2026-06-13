package procurementdashboard

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the procurement dashboard.
// This is a composition-surface (no proto entity) so it has no Labels DefaultLabels
// factory yet — LabelJSON is left empty until one is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "procurement.procurementdashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "procurement"},
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "procurement_request:list|supplier_contract:list",
			AppEntry: &compose.AppEntry{
				Key: "procurement", Route: "procurement.dashboard",
				Label: "Procurement", Icon: "icon-clipboard",
				Permission: "procurement_request:list|supplier_contract:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "procurement.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard",
					Permission: "procurement_request:list|supplier_contract:list"},
			},
		},
	}
}
