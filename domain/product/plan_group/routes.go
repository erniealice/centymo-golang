package plan_group

// PlanGroup route constants. A plan_group is a stable taxonomy node that
// groups related plans across periods (e.g. "Junior High" groups Grade 7/8/9).
const (
	DashboardURL     = "/plan-groups/dashboard"
	ListURL          = "/plan-groups/list/{status}"
	TableURL         = "/action/plan-group/table/{status}"
	DetailURL        = "/plan-groups/detail/{id}"
	AddURL           = "/action/plan-group/add"
	EditURL          = "/action/plan-group/edit/{id}"
	DeleteURL        = "/action/plan-group/delete"
	BulkDeleteURL    = "/action/plan-group/bulk-delete"
	SetStatusURL     = "/action/plan-group/set-status"
	BulkSetStatusURL = "/action/plan-group/bulk-set-status"
	TabActionURL     = "/action/plan-group/{id}/tab/{tab}"
)

// Routes holds all route paths for plan group views and actions.
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
		ActiveSubNav:     "plan-groups",
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
// plan group routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"plan_group.dashboard":       r.DashboardURL,
		"plan_group.list":            r.ListURL,
		"plan_group.table":           r.TableURL,
		"plan_group.detail":          r.DetailURL,
		"plan_group.add":             r.AddURL,
		"plan_group.edit":            r.EditURL,
		"plan_group.delete":          r.DeleteURL,
		"plan_group.bulk_delete":     r.BulkDeleteURL,
		"plan_group.set_status":      r.SetStatusURL,
		"plan_group.bulk_set_status": r.BulkSetStatusURL,
		"plan_group.tab_action":      r.TabActionURL,
	}
}
