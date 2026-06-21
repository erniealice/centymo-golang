package inventory

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the inventory entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "inventory.inventory",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "inventory"},
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
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "product:list"},
				{Key: "movements", Route: "inventory.movements",
					Label: "Movements", Icon: "icon-repeat", Permission: "product:list"},
			},
		},
	}
}
