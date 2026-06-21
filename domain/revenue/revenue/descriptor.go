package revenue

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the revenue entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "revenue.revenue",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "revenue"},
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "invoice:list",
			AppEntry: &compose.AppEntry{
				Key: "revenue", Route: "revenue.dashboard",
				Label: "Sales", Icon: "icon-trending-up",
				Permission: "invoice:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "revenue.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "invoice:list"},
				// Revenue (invoices) by status
				{Key: "draft", Route: "revenue.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "invoice:list"},
				{Key: "complete", Route: "revenue.list", Params: map[string]string{"status": "complete"},
					Label: "Complete", Icon: "icon-check-circle", Permission: "invoice:list"},
				{Key: "cancelled", Route: "revenue.list", Params: map[string]string{"status": "cancelled"},
					Label: "Cancelled", Icon: "icon-x-circle", Permission: "invoice:list"},
				// Note: invoice templates URL (SettingsTemplatesURL) is not in the
				// revenue RouteMap — it will be added in Phase 2 sidebar skeleton.
			},
		},
	}
}
