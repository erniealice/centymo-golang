package inventory

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the inventory entity.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "inventory.inventory",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "inventory"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "inventory.json", Key: "inventory"},
		LabelName: "InventoryLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "product:list",
			AppEntry: &compose.AppEntry{
				Key: "inventory", Route: "inventory.dashboard",
				Label: "Inventory", Icon: "icon-package",
				Permission: "product:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "inventory.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "product:list", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				{Key: "movements", Route: "inventory.movements",
					Label: "Movements", Icon: "icon-repeat", Permission: "product:list", LabelKey: "movements_label", IconKey: "movements_icon"},
			},
		},
	}
}
