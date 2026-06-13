package cost_plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.cost_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "cost_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "cost_plan.json", Key: "costPlan"},
		LabelName: "CostPlanLabels",
		Templates: templateFS,
		Nav: compose.NavContrib{
			Permission: "cost_plan:list",
			Items: []compose.NavItem{
				// supplier app — Cost Plans section
				{Key: "cost-plans-active", Route: "cost_plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "cost_plan:list"},
				{Key: "cost-plans-inactive", Route: "cost_plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "cost_plan:list"},
			},
		},
	}
}
