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

	// Teaching-staff tab (detail) — the section-centric class-edge upsert.
	// {id} is the section (subscription_group); it is authoritative from the
	// signed path, never a request body (no cross-tenant IDOR). The upsert
	// branches create/update/clear on the (section, product_plan) active edge.
	// Verb-first like edit/{id} — a "{id}/assign" shape is ambiguous with
	// "edit/{id}" under ServeMux precedence and panics at route registration.
	AssignURL = "/action/subscription-group/assign/{id}"

	// Attachments tab (detail) — upload/delete handlers behind the drawer.
	AttachmentUploadURL = "/action/subscription-group/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/subscription-group/detail/{id}/attachments/delete"
)

// Routes holds all route paths for subscription group views and actions.
type Routes struct {
	ActiveNav           string `json:"active_nav"`
	ActiveSubNav        string `json:"active_sub_nav"`
	DashboardURL        string `json:"dashboard_url"`
	ListURL             string `json:"list_url"`
	TableURL            string `json:"table_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	BulkDeleteURL       string `json:"bulk_delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AssignURL           string `json:"assign_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined above.
func DefaultRoutes() Routes {
	return Routes{
		// The sidebar mounts subscription groups in the "job" app (both apps,
		// sidebar_staff.go Operations group, right after Engagements), so
		// pages must highlight that app — not "service". ActiveSubNav stays a
		// BARE prefix: the list page composes "<base>-<status>" to match the
		// sidebar item keys (subscription-groups-active/-inactive).
		ActiveNav:           "job",
		ActiveSubNav:        "subscription-groups",
		DashboardURL:        DashboardURL,
		ListURL:             ListURL,
		TableURL:            TableURL,
		DetailURL:           DetailURL,
		AddURL:              AddURL,
		EditURL:             EditURL,
		DeleteURL:           DeleteURL,
		BulkDeleteURL:       BulkDeleteURL,
		SetStatusURL:        SetStatusURL,
		BulkSetStatusURL:    BulkSetStatusURL,
		TabActionURL:        TabActionURL,
		AssignURL:           AssignURL,
		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription group routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"subscription_group.dashboard":         r.DashboardURL,
		"subscription_group.list":              r.ListURL,
		"subscription_group.table":             r.TableURL,
		"subscription_group.detail":            r.DetailURL,
		"subscription_group.add":               r.AddURL,
		"subscription_group.edit":              r.EditURL,
		"subscription_group.delete":            r.DeleteURL,
		"subscription_group.bulk_delete":       r.BulkDeleteURL,
		"subscription_group.set_status":        r.SetStatusURL,
		"subscription_group.bulk_set_status":   r.BulkSetStatusURL,
		"subscription_group.tab_action":        r.TabActionURL,
		"subscription_group.assign":            r.AssignURL,
		"subscription_group.attachment.upload": r.AttachmentUploadURL,
		"subscription_group.attachment.delete": r.AttachmentDeleteURL,
	}
}
