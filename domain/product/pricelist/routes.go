package pricelist

// ---------------------------------------------------------------------------
// Price List route constants
// ---------------------------------------------------------------------------

const (
	// Price List routes — canonical home is the inventory accordion (/app/inventory/price-lists/*)
	ListURL       = "/inventory/price-lists/list/{status}"
	TableURL      = "/action/inventory-price-list/table/{status}"
	DetailURL     = "/inventory/price-lists/detail/{id}"
	AddURL        = "/action/inventory-price-list/add"
	EditURL       = "/action/inventory-price-list/edit/{id}"
	DeleteURL     = "/action/inventory-price-list/delete"
	BulkDeleteURL = "/action/inventory-price-list/bulk-delete"

	TabActionURL        = "/action/inventory-price-list/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/inventory-price-list/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/inventory-price-list/{id}/attachments/delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/inventory-price-list/{id}/products/add"
	PriceProductDeleteURL = "/action/inventory-price-list/{id}/products/delete"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds all route paths for price list views and actions,
// including price product sub-routes.
type Routes struct {
	ListURL       string `json:"list_url"`
	TableURL      string `json:"table_url"`
	DetailURL     string `json:"detail_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`

	TabActionURL string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Price product routes
	PriceProductAddURL    string `json:"price_product_add_url"`
	PriceProductDeleteURL string `json:"price_product_delete_url"`
}

// DefaultRoutes returns a Routes populated from the
// package-level route constants defined in this file.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:       ListURL,
		TableURL:      TableURL,
		DetailURL:     DetailURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,

		TabActionURL: TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		PriceProductAddURL:    PriceProductAddURL,
		PriceProductDeleteURL: PriceProductDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price list routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"price_list.list":        r.ListURL,
		"price_list.table":       r.TableURL,
		"price_list.detail":      r.DetailURL,
		"price_list.add":         r.AddURL,
		"price_list.edit":        r.EditURL,
		"price_list.delete":      r.DeleteURL,
		"price_list.bulk_delete": r.BulkDeleteURL,

		"price_list.tab_action": r.TabActionURL,

		"price_list.attachment.upload": r.AttachmentUploadURL,
		"price_list.attachment.delete": r.AttachmentDeleteURL,

		"price_list.price_product.add":    r.PriceProductAddURL,
		"price_list.price_product.delete": r.PriceProductDeleteURL,
	}
}
