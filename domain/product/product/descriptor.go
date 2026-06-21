package product

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the product entity
// (services mount). Labels are not yet exposed via DefaultLabels() — the
// LabelJSON binding is left empty until a DefaultLabels factory is added.
// Note: product has multiple mounts (services/inventory/supplies); each
// gets its own Unit key in catalog.go with per-mount route overrides applied.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "product.product",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "product"},
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "service:list",
			AppEntry: &compose.AppEntry{
				Key: "service", Route: "product.list", Params: map[string]string{"status": "active"},
				Label: "Services", Icon: "icon-briefcase",
				Permission: "service:list",
			},
			Items: []compose.NavItem{
				{Key: "services-active", Route: "product.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "service:list"},
				{Key: "services-inactive", Route: "product.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "service:list"},
			},
		},
	}
}
