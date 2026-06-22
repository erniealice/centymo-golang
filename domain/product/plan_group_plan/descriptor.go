package plan_group_plan

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.plan_group_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "plan_group_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "plan_group_plan.json", Key: "planGroupPlan"},
		LabelName: "PlanGroupPlanLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "plan_group_plan:list",
			Items: []compose.NavItem{
				{Key: "plan-group-plans-active", Route: "plan_group_plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "plan_group_plan:list"},
				{Key: "plan-group-plans-inactive", Route: "plan_group_plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "plan_group_plan:list"},
			},
		},
	}
}
