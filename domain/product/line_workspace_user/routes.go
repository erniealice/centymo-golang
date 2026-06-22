package line_workspace_user

// LineWorkspaceUser route constants. A line_workspace_user pins an operator
// (workspace_user) at a line node for tier-2 group visibility.
const (
	DashboardURL     = "/line-workspace-users/dashboard"
	ListURL          = "/line-workspace-users/list/{status}"
	TableURL         = "/action/line-workspace-user/table/{status}"
	DetailURL        = "/line-workspace-users/detail/{id}"
	AddURL           = "/action/line-workspace-user/add"
	EditURL          = "/action/line-workspace-user/edit/{id}"
	DeleteURL        = "/action/line-workspace-user/delete"
	BulkDeleteURL    = "/action/line-workspace-user/bulk-delete"
	SetStatusURL     = "/action/line-workspace-user/set-status"
	BulkSetStatusURL = "/action/line-workspace-user/bulk-set-status"
	TabActionURL     = "/action/line-workspace-user/{id}/tab/{tab}"
)

// Routes holds all route paths for line_workspace_user views and actions.
type Routes struct {
	ActiveNav        string `json:"active_nav"`
	ActiveSubNav     string `json:"active_sub_nav"`
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined above.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:        "product",
		ActiveSubNav:     "line-workspace-users",
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		TabActionURL:     TabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// line_workspace_user routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"line_workspace_user.dashboard":       r.DashboardURL,
		"line_workspace_user.list":            r.ListURL,
		"line_workspace_user.table":           r.TableURL,
		"line_workspace_user.detail":          r.DetailURL,
		"line_workspace_user.add":             r.AddURL,
		"line_workspace_user.edit":            r.EditURL,
		"line_workspace_user.delete":          r.DeleteURL,
		"line_workspace_user.bulk_delete":     r.BulkDeleteURL,
		"line_workspace_user.set_status":      r.SetStatusURL,
		"line_workspace_user.bulk_set_status": r.BulkSetStatusURL,
		"line_workspace_user.tab_action":      r.TabActionURL,
	}
}
