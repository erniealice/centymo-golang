package plan_group_plan

// PlanGroupPlan domain route constants. A plan_group_plan links a plan_group
// to a specific plan, with an optional sequence_order.
const (
	DashboardURL     = "/plan-group-plans/dashboard"
	ListURL          = "/plan-group-plans/list/{status}"
	TableURL         = "/action/plan-group-plan/table/{status}"
	DetailURL        = "/plan-group-plans/detail/{id}"
	AddURL           = "/action/plan-group-plan/add"
	EditURL          = "/action/plan-group-plan/edit/{id}"
	DeleteURL        = "/action/plan-group-plan/delete"
	BulkDeleteURL    = "/action/plan-group-plan/bulk-delete"
	SetStatusURL     = "/action/plan-group-plan/set-status"
	BulkSetStatusURL = "/action/plan-group-plan/bulk-set-status"
	TabActionURL     = "/action/plan-group-plan/{id}/tab/{tab}"
)

// Routes holds all route paths for plan_group_plan views and actions.
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
		ActiveSubNav:     "plan-group-plans",
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
// plan_group_plan routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"plan_group_plan.dashboard":       r.DashboardURL,
		"plan_group_plan.list":            r.ListURL,
		"plan_group_plan.table":           r.TableURL,
		"plan_group_plan.detail":          r.DetailURL,
		"plan_group_plan.add":             r.AddURL,
		"plan_group_plan.edit":            r.EditURL,
		"plan_group_plan.delete":          r.DeleteURL,
		"plan_group_plan.bulk_delete":     r.BulkDeleteURL,
		"plan_group_plan.set_status":      r.SetStatusURL,
		"plan_group_plan.bulk_set_status": r.BulkSetStatusURL,
		"plan_group_plan.tab_action":      r.TabActionURL,
	}
}
