package subscription_group_product_plan

// subscription_group_product_plan route constants (centymo.md §2 — the
// complete URL map). {id} is always a subscription_group (section) id;
// {sgppId} is always a subscription_group_product_plan (class) id.
const (
	DashboardURL  = "/subscription-group-product-plans/dashboard"
	ListURL       = "/subscription-group-product-plans/list/{status}"
	TableURL      = "/action/subscription-group-product-plan/table/{status}"
	AddURL        = "/action/subscription-group-product-plan/add"
	EditURL       = "/action/subscription-group-product-plan/edit/{id}"
	DeleteURL     = "/action/subscription-group-product-plan/delete"
	BulkDeleteURL = "/action/subscription-group-product-plan/bulk-delete"
	// DetailURL — the S4 class page. {id}=section, {sgppId}=class. The generic
	// path renders as-is; an education-tier route.json override renders the
	// same handler under /sections/detail/{id}/subject/{sgppId} (lyngua-only).
	DetailURL = "/subscription-groups/detail/{id}/offering/{sgppId}"
	// TabActionURL — S4 HTMX tab swap (info | staff canonical keys).
	TabActionURL = "/action/subscription-group-product-plan/{sgppId}/tab/{tab}"
	// PickerURL — S2 add-offerings picker. {id}=section id (GET candidates,
	// POST create-batch).
	PickerURL = "/action/subscription-group-product-plan/picker/{id}"
	// AssignURL — S3 (quick drawer, mounted from the section tab) / S5 (the
	// class page's Teachers-tab Add drawer) shared assign action. {sgppId}=
	// class id. GET renders the form; POST creates/updates one or more
	// subscription_group_product_plan_staff rows; clear=1+row id soft-deletes
	// one (the sgpps Clear idiom).
	AssignURL = "/action/subscription-group-product-plan/assign/{sgppId}"
	// SetStatusURL — S6 exclude/restore (the ACTIVE<->EXCLUDED status enum
	// transition; distinct from the entity's active/inactive soft-delete axis,
	// which the Edit form's Active toggle alone governs — S7 carries no
	// activate/deactivate row action by spec).
	SetStatusURL = "/action/subscription-group-product-plan/set-status/{sgppId}"
)

// Routes holds all route paths for subscription_group_product_plan views and
// actions.
type Routes struct {
	ActiveNav     string `json:"active_nav"`
	ActiveSubNav  string `json:"active_sub_nav"`
	DashboardURL  string `json:"dashboard_url"`
	ListURL       string `json:"list_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	DetailURL     string `json:"detail_url"`
	TabActionURL  string `json:"tab_action_url"`
	PickerURL     string `json:"picker_url"`
	AssignURL     string `json:"assign_url"`
	SetStatusURL  string `json:"set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:     "service",
		ActiveSubNav:  "subscription-group-product-plans",
		DashboardURL:  DashboardURL,
		ListURL:       ListURL,
		TableURL:      TableURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,
		DetailURL:     DetailURL,
		TabActionURL:  TabActionURL,
		PickerURL:     PickerURL,
		AssignURL:     AssignURL,
		SetStatusURL:  SetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"subscription_group_product_plan.dashboard":   r.DashboardURL,
		"subscription_group_product_plan.list":        r.ListURL,
		"subscription_group_product_plan.table":       r.TableURL,
		"subscription_group_product_plan.add":         r.AddURL,
		"subscription_group_product_plan.edit":        r.EditURL,
		"subscription_group_product_plan.delete":      r.DeleteURL,
		"subscription_group_product_plan.bulk_delete": r.BulkDeleteURL,
		"subscription_group_product_plan.detail":      r.DetailURL,
		"subscription_group_product_plan.tab_action":  r.TabActionURL,
		"subscription_group_product_plan.picker":      r.PickerURL,
		"subscription_group_product_plan.assign":      r.AssignURL,
		"subscription_group_product_plan.set_status":  r.SetStatusURL,
	}
}
