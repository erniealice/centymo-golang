package resource

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.resource",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "resource"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "resource.json", Key: "resource"},
		LabelName: "ResourceLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "resource:list",
			Items: []compose.NavItem{
				// service app — "Resources" section
				{Key: "resources-active", Route: "resource.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "resource:list"},
				{Key: "resources-inactive", Route: "resource.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "resource:list"},
			},
		},
	}
}
