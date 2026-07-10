package product

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the product entity
// (services mount).
// Note: product has multiple mounts (services/inventory/supplies); each
// gets its own Unit key in catalog.go with per-mount route overrides applied.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.product",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "product"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "product.json", Key: "product"},
		LabelName: "ProductLabels",
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
					Label: "Active", Icon: "icon-check-circle", Permission: "service:list", LabelKey: "active_label", IconKey: "services_active_icon"},
				{Key: "services-inactive", Route: "product.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "service:list", LabelKey: "inactive_label", IconKey: "services_inactive_icon"},
			},
		},
	}
}
