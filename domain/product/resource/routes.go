package resource

// ---------------------------------------------------------------------------
// Resource route constants
// ---------------------------------------------------------------------------

const (
	// Resource routes (person or equipment linked to a Product for billing)
	ListURL          = "/resources/list/{status}"
	TableURL         = "/action/resource/table/{status}"
	DetailURL        = "/resources/detail/{id}"
	AddURL           = "/action/resource/add"
	EditURL          = "/action/resource/edit/{id}"
	DeleteURL        = "/action/resource/delete"
	BulkDeleteURL    = "/action/resource/bulk-delete"
	SetStatusURL     = "/action/resource/set-status"
	BulkSetStatusURL = "/action/resource/bulk-set-status"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds all route paths for resource views and actions.
// A resource links a person or equipment to a Product for billing purposes.
type Routes struct {
	ActiveNav        string `json:"active_nav"`
	ActiveSubNav     string `json:"active_sub_nav"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in this file.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:        "service",
		ActiveSubNav:     "resources-active",
		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// resource routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"resource.list":            r.ListURL,
		"resource.table":           r.TableURL,
		"resource.detail":          r.DetailURL,
		"resource.add":             r.AddURL,
		"resource.edit":            r.EditURL,
		"resource.delete":          r.DeleteURL,
		"resource.bulk_delete":     r.BulkDeleteURL,
		"resource.set_status":      r.SetStatusURL,
		"resource.bulk_set_status": r.BulkSetStatusURL,
	}
}
