package line_workspace_user

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.line_workspace_user",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "line_workspace_user"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "line_workspace_user.json", Key: "lineWorkspaceUser"},
		LabelName: "LineWorkspaceUserLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "line_workspace_user:list",
			Items: []compose.NavItem{
				{Key: "line-workspace-users-active", Route: "line_workspace_user.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "line_workspace_user:list"},
				{Key: "line-workspace-users-inactive", Route: "line_workspace_user.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "line_workspace_user:list"},
			},
		},
	}
}
