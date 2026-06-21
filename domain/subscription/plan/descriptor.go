package plan

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "plan.json", Key: "plan"},
		LabelName: "PlanLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "plan:list",
			Items: []compose.NavItem{
				// service app — "Plans" / "Packages" section
				{Key: "plans-active", Route: "plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "plan:list"},
				{Key: "plans-inactive", Route: "plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "plan:list"},
			},
		},
	}
}
