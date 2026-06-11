package cost_schedule

// routes.go — cost_schedule route constants + Routes type (centymo W7).
//
// Extracted from the procurement-domain routes.go per the domain-first
// restructure. Pure structural move — no behaviour change; route strings are
// byte-identical.

// Default route constants for the cost_schedule module.
const (
	ListURL          = "/cost-schedules/list/{status}"
	TableURL         = "/action/cost-schedule/table/{status}"
	DetailURL        = "/cost-schedules/detail/{id}"
	AddURL           = "/action/cost-schedule/add"
	EditURL          = "/action/cost-schedule/edit/{id}"
	DeleteURL        = "/action/cost-schedule/delete"
	BulkDeleteURL    = "/action/cost-schedule/bulk-delete"
	SetStatusURL     = "/action/cost-schedule/set-status"
	BulkSetStatusURL = "/action/cost-schedule/bulk-set-status"
	TabActionURL     = "/action/cost-schedule/detail/{id}/tab/{tab}"
)

// Routes holds all route paths for cost_schedule views.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

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

// DefaultRoutes returns a Routes using route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-schedules",

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
		"cost_schedule.list":            r.ListURL,
		"cost_schedule.table":           r.TableURL,
		"cost_schedule.detail":          r.DetailURL,
		"cost_schedule.add":             r.AddURL,
		"cost_schedule.edit":            r.EditURL,
		"cost_schedule.delete":          r.DeleteURL,
		"cost_schedule.bulk_delete":     r.BulkDeleteURL,
		"cost_schedule.set_status":      r.SetStatusURL,
		"cost_schedule.bulk_set_status": r.BulkSetStatusURL,
		"cost_schedule.tab_action":      r.TabActionURL,
	}
}
