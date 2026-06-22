package price_schedule_workspace_user

// PriceScheduleWorkspaceUser-domain route constants. A
// price_schedule_workspace_user pins a workspace operator at a price_schedule
// (period / academic-year) node for period-level visibility — the "year
// coordinator" access record.
const (
	DashboardURL     = "/price-schedule-workspace-users/dashboard"
	ListURL          = "/price-schedule-workspace-users/list/{status}"
	TableURL         = "/action/price-schedule-workspace-user/table/{status}"
	DetailURL        = "/price-schedule-workspace-users/detail/{id}"
	AddURL           = "/action/price-schedule-workspace-user/add"
	EditURL          = "/action/price-schedule-workspace-user/edit/{id}"
	DeleteURL        = "/action/price-schedule-workspace-user/delete"
	BulkDeleteURL    = "/action/price-schedule-workspace-user/bulk-delete"
	SetStatusURL     = "/action/price-schedule-workspace-user/set-status"
	BulkSetStatusURL = "/action/price-schedule-workspace-user/bulk-set-status"
	TabActionURL     = "/action/price-schedule-workspace-user/{id}/tab/{tab}"
)

// Routes holds all route paths for price_schedule_workspace_user views and actions.
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
		ActiveNav:        "service",
		ActiveSubNav:     "price-schedule-workspace-users",
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
// price_schedule_workspace_user routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"price_schedule_workspace_user.dashboard":       r.DashboardURL,
		"price_schedule_workspace_user.list":            r.ListURL,
		"price_schedule_workspace_user.table":           r.TableURL,
		"price_schedule_workspace_user.detail":          r.DetailURL,
		"price_schedule_workspace_user.add":             r.AddURL,
		"price_schedule_workspace_user.edit":            r.EditURL,
		"price_schedule_workspace_user.delete":          r.DeleteURL,
		"price_schedule_workspace_user.bulk_delete":     r.BulkDeleteURL,
		"price_schedule_workspace_user.set_status":      r.SetStatusURL,
		"price_schedule_workspace_user.bulk_set_status": r.BulkSetStatusURL,
		"price_schedule_workspace_user.tab_action":      r.TabActionURL,
	}
}
