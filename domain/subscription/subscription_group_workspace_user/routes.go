package subscription_group_workspace_user

// SubscriptionGroupWorkspaceUser route constants. Pins an operator
// (workspace_user) at a subscription_group (cohort) for group-level servicing.
const (
	DashboardURL     = "/subscription-group-workspace-users/dashboard"
	ListURL          = "/subscription-group-workspace-users/list/{status}"
	TableURL         = "/action/subscription-group-workspace-user/table/{status}"
	DetailURL        = "/subscription-group-workspace-users/detail/{id}"
	AddURL           = "/action/subscription-group-workspace-user/add"
	EditURL          = "/action/subscription-group-workspace-user/edit/{id}"
	DeleteURL        = "/action/subscription-group-workspace-user/delete"
	BulkDeleteURL    = "/action/subscription-group-workspace-user/bulk-delete"
	SetStatusURL     = "/action/subscription-group-workspace-user/set-status"
	BulkSetStatusURL = "/action/subscription-group-workspace-user/bulk-set-status"
	TabActionURL     = "/action/subscription-group-workspace-user/{id}/tab/{tab}"
)

// Routes holds all route paths for subscription_group_workspace_user views and actions.
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

// DefaultRoutes returns a Routes populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:        "service",
		ActiveSubNav:     "subscription-group-workspace-users",
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

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"subscription_group_workspace_user.dashboard":       r.DashboardURL,
		"subscription_group_workspace_user.list":            r.ListURL,
		"subscription_group_workspace_user.table":           r.TableURL,
		"subscription_group_workspace_user.detail":          r.DetailURL,
		"subscription_group_workspace_user.add":             r.AddURL,
		"subscription_group_workspace_user.edit":            r.EditURL,
		"subscription_group_workspace_user.delete":          r.DeleteURL,
		"subscription_group_workspace_user.bulk_delete":     r.BulkDeleteURL,
		"subscription_group_workspace_user.set_status":      r.SetStatusURL,
		"subscription_group_workspace_user.bulk_set_status": r.BulkSetStatusURL,
		"subscription_group_workspace_user.tab_action":      r.TabActionURL,
	}
}
