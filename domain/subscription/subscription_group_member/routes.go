package subscription_group_member

// subscription_group_member route constants.
const (
	DashboardURL     = "/subscription-group-members/dashboard"
	ListURL          = "/subscription-group-members/list/{status}"
	TableURL         = "/action/subscription-group-member/table/{status}"
	DetailURL        = "/subscription-group-members/detail/{id}"
	AddURL           = "/action/subscription-group-member/add"
	EditURL          = "/action/subscription-group-member/edit/{id}"
	DeleteURL        = "/action/subscription-group-member/delete"
	BulkDeleteURL    = "/action/subscription-group-member/bulk-delete"
	SetStatusURL     = "/action/subscription-group-member/set-status"
	BulkSetStatusURL = "/action/subscription-group-member/bulk-set-status"
	TabActionURL     = "/action/subscription-group-member/{id}/tab/{tab}"
)

// Routes holds all route paths for subscription_group_member views and actions.
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
		ActiveSubNav:     "subscription-group-members",
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
		"subscription_group_member.dashboard":       r.DashboardURL,
		"subscription_group_member.list":            r.ListURL,
		"subscription_group_member.table":           r.TableURL,
		"subscription_group_member.detail":          r.DetailURL,
		"subscription_group_member.add":             r.AddURL,
		"subscription_group_member.edit":            r.EditURL,
		"subscription_group_member.delete":          r.DeleteURL,
		"subscription_group_member.bulk_delete":     r.BulkDeleteURL,
		"subscription_group_member.set_status":      r.SetStatusURL,
		"subscription_group_member.bulk_set_status": r.BulkSetStatusURL,
		"subscription_group_member.tab_action":      r.TabActionURL,
	}
}
