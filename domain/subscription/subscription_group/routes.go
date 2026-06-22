package subscription_group

// SubscriptionGroup-domain route constants. A subscription_group is the
// education "section/cohort" — a per-period roster anchored to a plan (the
// program) and a price_schedule (the period). Mirrors the price_schedule core
// CRUD route shape minus the nested-plan feature.
const (
	DashboardURL     = "/subscription-groups/dashboard"
	ListURL          = "/subscription-groups/list/{status}"
	TableURL         = "/action/subscription-group/table/{status}"
	DetailURL        = "/subscription-groups/detail/{id}"
	AddURL           = "/action/subscription-group/add"
	EditURL          = "/action/subscription-group/edit/{id}"
	DeleteURL        = "/action/subscription-group/delete"
	BulkDeleteURL    = "/action/subscription-group/bulk-delete"
	SetStatusURL     = "/action/subscription-group/set-status"
	BulkSetStatusURL = "/action/subscription-group/bulk-set-status"
	TabActionURL     = "/action/subscription-group/{id}/tab/{tab}"
)

// Routes holds all route paths for subscription group views and actions.
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
		ActiveSubNav:     "subscription-groups",
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
// subscription group routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"subscription_group.dashboard":       r.DashboardURL,
		"subscription_group.list":            r.ListURL,
		"subscription_group.table":           r.TableURL,
		"subscription_group.detail":          r.DetailURL,
		"subscription_group.add":             r.AddURL,
		"subscription_group.edit":            r.EditURL,
		"subscription_group.delete":          r.DeleteURL,
		"subscription_group.bulk_delete":     r.BulkDeleteURL,
		"subscription_group.set_status":      r.SetStatusURL,
		"subscription_group.bulk_set_status": r.BulkSetStatusURL,
		"subscription_group.tab_action":      r.TabActionURL,
	}
}
