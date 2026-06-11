package price_plan

// PricePlan-domain route constants. Relocated from the subscription routes.go
// god-file (entity-local extraction). Pure structural move.
const (
	// PricePlan standalone routes (rate cards as independent entity)
	DashboardURL        = "/price-plans/dashboard"
	ListURL             = "/price-plans/list/{status}"
	TableURL            = "/action/price-plan/table/{status}"
	DetailURL           = "/price-plans/detail/{id}"
	StandaloneAddURL    = "/action/price-plan/add"
	StandaloneEditURL   = "/action/price-plan/edit/{id}"
	StandaloneDeleteURL = "/action/price-plan/delete"
	BulkDeleteURL       = "/action/price-plan/bulk-delete"
	SetStatusURL        = "/action/price-plan/set-status"
	BulkSetStatusURL    = "/action/price-plan/bulk-set-status"
	TabActionURL        = "/action/price-plan/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/price-plan/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/price-plan/{id}/attachments/delete"

	// ProductPricePlan CRUD routes (within price plan / rate card detail)
	ProductPriceAddURL    = "/action/price-plan/{id}/product-prices/add"
	ProductPriceEditURL   = "/action/price-plan/{id}/product-prices/edit/{ppid}"
	ProductPriceDeleteURL = "/action/price-plan/{id}/product-prices/delete"
)

// Routes holds all route paths for price plan (rate card) views and actions.
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
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// ProductPricePlan CRUD routes (within rate card detail)
	ProductPriceAddURL    string `json:"product_price_add_url"`
	ProductPriceEditURL   string `json:"product_price_edit_url"`
	ProductPriceDeleteURL string `json:"product_price_delete_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:             "service",
		ActiveSubNav:          "rate-cards",
		DashboardURL:          DashboardURL,
		ListURL:               ListURL,
		TableURL:              TableURL,
		DetailURL:             DetailURL,
		AddURL:                StandaloneAddURL,
		EditURL:               StandaloneEditURL,
		DeleteURL:             StandaloneDeleteURL,
		BulkDeleteURL:         BulkDeleteURL,
		SetStatusURL:          SetStatusURL,
		BulkSetStatusURL:      BulkSetStatusURL,
		TabActionURL:          TabActionURL,
		AttachmentUploadURL:   AttachmentUploadURL,
		AttachmentDeleteURL:   AttachmentDeleteURL,
		ProductPriceAddURL:    ProductPriceAddURL,
		ProductPriceEditURL:   ProductPriceEditURL,
		ProductPriceDeleteURL: ProductPriceDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price plan routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"price_plan.dashboard":            r.DashboardURL,
		"price_plan.list":                 r.ListURL,
		"price_plan.table":                r.TableURL,
		"price_plan.detail":               r.DetailURL,
		"price_plan.add":                  r.AddURL,
		"price_plan.edit":                 r.EditURL,
		"price_plan.delete":               r.DeleteURL,
		"price_plan.bulk_delete":          r.BulkDeleteURL,
		"price_plan.set_status":           r.SetStatusURL,
		"price_plan.bulk_set_status":      r.BulkSetStatusURL,
		"price_plan.tab_action":           r.TabActionURL,
		"price_plan.attachment.upload":    r.AttachmentUploadURL,
		"price_plan.attachment.delete":    r.AttachmentDeleteURL,
		"price_plan.product_price.add":    r.ProductPriceAddURL,
		"price_plan.product_price.edit":   r.ProductPriceEditURL,
		"price_plan.product_price.delete": r.ProductPriceDeleteURL,
	}
}
