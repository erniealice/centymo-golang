package plan_group

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.plan_group",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "plan_group"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "plan_group.json", Key: "planGroup"},
		LabelName: "PlanGroupLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "plan_group:list",
			Items: []compose.NavItem{
				{Key: "plan-groups-active", Route: "plan_group.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-layers", Permission: "plan_group:list"},
				{Key: "plan-groups-inactive", Route: "plan_group.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "plan_group:list"},
			},
		},
	}
}
