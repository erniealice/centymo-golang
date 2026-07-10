package revenue

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the revenue entity.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "revenue.revenue",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "revenue"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "revenue.json", Key: "revenue"},
		LabelName: "RevenueLabels",
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
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "invoice:list", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				// Revenue (invoices) by status
				{Key: "draft", Route: "revenue.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "invoice:list", LabelKey: "revenue_draft_label", IconKey: "revenue_draft_icon"},
				{Key: "complete", Route: "revenue.list", Params: map[string]string{"status": "complete"},
					Label: "Complete", Icon: "icon-check-circle", Permission: "invoice:list", LabelKey: "revenue_complete_label", IconKey: "revenue_complete_icon"},
				{Key: "cancelled", Route: "revenue.list", Params: map[string]string{"status": "cancelled"},
					Label: "Cancelled", Icon: "icon-x-circle", Permission: "invoice:list", LabelKey: "revenue_cancelled_label", IconKey: "revenue_cancelled_icon"},
				// Note: invoice templates URL (SettingsTemplatesURL) is not in the
				// revenue RouteMap — it will be added in Phase 2 sidebar skeleton.
			},
		},
	}
}
