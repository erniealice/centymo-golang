package product

import "strings"

// Three-level routing system for centymo product-domain views:
//
// Level 1: Generic defaults from Go consts (this file).
//   DefaultXxxRoutes() constructors return structs populated from the route
//   constants defined below. These are sensible defaults that work
//   out of the box for any app.
//
// Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//   Consumer apps can load a JSON config that partially overrides the
//   default routes. Struct fields carry json tags for unmarshalling.
//
// Level 3: App-specific overrides via Go field assignment (optional).
//   After loading defaults and/or JSON, consumer apps can programmatically
//   set individual fields to further customize routing.
//
// Each route struct also exposes a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "product.list"),
// useful for template rendering, URL resolution, and debugging.

// ---------------------------------------------------------------------------
// Product route constants
// ---------------------------------------------------------------------------

const (
	ListURL   = "/products/list/{status}"
	TableURL  = "/action/product/table/{status}"
	DetailURL = "/products/detail/{id}"

	// Service dashboard — services are products filtered to product_kind="service".
	// The dashboard sits ABOVE the service-mount product list at this URL.
	ServiceDashboardURL = "/services/dashboard"

	// Product action routes
	AddURL        = "/action/product/add"
	EditURL       = "/action/product/edit/{id}"
	DeleteURL     = "/action/product/delete"
	BulkDeleteURL = "/action/product/bulk-delete"

	// Product status routes
	SetStatusURL     = "/action/product/set-status"
	BulkSetStatusURL = "/action/product/bulk-set-status"

	// Product detail tab action route
	TabActionURL          = "/action/product/detail/{id}/tab/{tab}"
	AttachmentUploadURL   = "/action/product/detail/{id}/attachments/upload"
	AttachmentDeleteURL   = "/action/product/detail/{id}/attachments/delete"
	AttachmentDownloadURL = "/action/product/detail/{id}/attachments/download"

	// Product variant routes (within product detail)
	VariantTableURL  = "/action/product/detail/{id}/variants/table"
	VariantAssignURL = "/action/product/detail/{id}/variants/assign"
	VariantEditURL   = "/action/product/detail/{id}/variants/edit/{vid}"
	VariantRemoveURL = "/action/product/detail/{id}/variants/remove"

	// Product attribute routes (within product detail)
	AttributeTableURL  = "/action/product/detail/{id}/attributes/table"
	AttributeAssignURL = "/action/product/detail/{id}/attributes/assign"
	AttributeRemoveURL = "/action/product/detail/{id}/attributes/remove"

	// Product option routes (within product detail)
	OptionTableURL  = "/action/product/detail/{id}/options/table"
	OptionAddURL    = "/action/product/detail/{id}/options/add"
	OptionEditURL   = "/action/product/detail/{id}/options/edit/{oid}"
	OptionDeleteURL = "/action/product/detail/{id}/options/delete"

	// Product option detail page (option values management)
	OptionDetailURL = "/products/detail/{id}/option/{oid}"

	// Product variant detail page (variant info, pricing, stock, audit, images)
	VariantDetailURL    = "/products/detail/{id}/variant/{vid}"
	VariantTabActionURL = "/action/product/detail/{id}/variant/{vid}/tab/{tab}"

	// Product variant image routes (upload/delete within variant detail)
	VariantImageUploadURL = "/action/product/detail/{id}/variant/{vid}/images/upload"
	VariantImageDeleteURL = "/action/product/detail/{id}/variant/{vid}/images/delete"

	// Product variant attachment routes
	VariantAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/attachments/upload"
	VariantAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/attachments/delete"

	// Product variant stock detail (inventory item within variant context)
	VariantStockDetailURL    = "/products/detail/{id}/variant/{vid}/stock/{iid}"
	VariantStockTabActionURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/tab/{tab}"

	// Product variant stock attachment routes
	VariantStockAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/upload"
	VariantStockAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/delete"

	// Inventory serial detail (individual serial within inventory item)
	VariantSerialDetailURL = "/products/detail/{id}/variant/{vid}/stock/{iid}/serial/{sid}"

	// Product option value routes (within product option)
	OptionValueTableURL  = "/action/product/detail/{id}/options/{oid}/values/table"
	OptionValueAddURL    = "/action/product/detail/{id}/options/{oid}/values/add"
	OptionValueEditURL   = "/action/product/detail/{id}/options/{oid}/values/edit/{vid}"
	OptionValueDeleteURL = "/action/product/detail/{id}/options/{oid}/values/delete"

	// Product line routes
	LineDashboardURL        = "/product-lines/dashboard"
	LineListURL             = "/product-lines/list/{status}"
	LineTableURL            = "/action/product-line/table/{status}"
	LineDetailURL           = "/product-lines/detail/{id}"
	LineAddURL              = "/action/product-line/add"
	LineEditURL             = "/action/product-line/edit/{id}"
	LineDeleteURL           = "/action/product-line/delete"
	LineBulkDeleteURL       = "/action/product-line/bulk-delete"
	LineSetStatusURL        = "/action/product-line/set-status"
	LineBulkSetStatusURL    = "/action/product-line/bulk-set-status"
	LineTabActionURL        = "/action/product-line/{id}/tab/{tab}"
	LineAttachmentUploadURL = "/action/product-line/{id}/attachments/upload"
	LineAttachmentDeleteURL = "/action/product-line/{id}/attachments/delete"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds all route paths for product views and actions,
// including variant, option, attribute, image, stock, and serial sub-routes.
type Routes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	// Phase 5 — dashboard URL for the service mount (product_kind=service).
	// Empty for non-service mounts; the inventory and supplies mounts route
	// dashboards from elsewhere.
	DashboardURL string `json:"dashboard_url"`

	ListURL       string `json:"list_url"`
	TableURL      string `json:"table_url"`
	DetailURL     string `json:"detail_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`

	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	TabActionURL string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL   string `json:"attachment_upload_url"`
	AttachmentDeleteURL   string `json:"attachment_delete_url"`
	AttachmentDownloadURL string `json:"attachment_download_url"`

	// Variant routes
	VariantTableURL  string `json:"variant_table_url"`
	VariantAssignURL string `json:"variant_assign_url"`
	VariantEditURL   string `json:"variant_edit_url"`
	VariantRemoveURL string `json:"variant_remove_url"`

	// Variant detail routes
	VariantDetailURL    string `json:"variant_detail_url"`
	VariantTabActionURL string `json:"variant_tab_action_url"`

	// Variant image routes
	VariantImageUploadURL string `json:"variant_image_upload_url"`
	VariantImageDeleteURL string `json:"variant_image_delete_url"`

	// Variant attachment routes
	VariantAttachmentUploadURL string `json:"variant_attachment_upload_url"`
	VariantAttachmentDeleteURL string `json:"variant_attachment_delete_url"`

	// Variant stock routes
	VariantStockDetailURL    string `json:"variant_stock_detail_url"`
	VariantStockTabActionURL string `json:"variant_stock_tab_action_url"`

	// Variant stock attachment routes
	VariantStockAttachmentUploadURL string `json:"variant_stock_attachment_upload_url"`
	VariantStockAttachmentDeleteURL string `json:"variant_stock_attachment_delete_url"`

	// Variant serial routes
	VariantSerialDetailURL string `json:"variant_serial_detail_url"`

	// Attribute routes
	AttributeTableURL  string `json:"attribute_table_url"`
	AttributeAssignURL string `json:"attribute_assign_url"`
	AttributeRemoveURL string `json:"attribute_remove_url"`

	// Option routes
	OptionTableURL  string `json:"option_table_url"`
	OptionAddURL    string `json:"option_add_url"`
	OptionEditURL   string `json:"option_edit_url"`
	OptionDeleteURL string `json:"option_delete_url"`
	OptionDetailURL string `json:"option_detail_url"`

	// Option value routes
	OptionValueTableURL  string `json:"option_value_table_url"`
	OptionValueAddURL    string `json:"option_value_add_url"`
	OptionValueEditURL   string `json:"option_value_edit_url"`
	OptionValueDeleteURL string `json:"option_value_delete_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in this file.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "inventory",
		ActiveSubNav: "masterlist",

		// Default to the service dashboard URL — only meaningful for the
		// service-flavoured mount. Inventory/supplies mounts overwrite or
		// ignore this field.
		DashboardURL: ServiceDashboardURL,

		ListURL:       ListURL,
		TableURL:      TableURL,
		DetailURL:     DetailURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,

		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,

		TabActionURL: TabActionURL,

		AttachmentUploadURL:   AttachmentUploadURL,
		AttachmentDeleteURL:   AttachmentDeleteURL,
		AttachmentDownloadURL: AttachmentDownloadURL,

		VariantTableURL:  VariantTableURL,
		VariantAssignURL: VariantAssignURL,
		VariantEditURL:   VariantEditURL,
		VariantRemoveURL: VariantRemoveURL,

		VariantDetailURL:    VariantDetailURL,
		VariantTabActionURL: VariantTabActionURL,

		VariantImageUploadURL: VariantImageUploadURL,
		VariantImageDeleteURL: VariantImageDeleteURL,

		VariantAttachmentUploadURL: VariantAttachmentUploadURL,
		VariantAttachmentDeleteURL: VariantAttachmentDeleteURL,

		VariantStockDetailURL:    VariantStockDetailURL,
		VariantStockTabActionURL: VariantStockTabActionURL,

		VariantStockAttachmentUploadURL: VariantStockAttachmentUploadURL,
		VariantStockAttachmentDeleteURL: VariantStockAttachmentDeleteURL,

		VariantSerialDetailURL: VariantSerialDetailURL,

		AttributeTableURL:  AttributeTableURL,
		AttributeAssignURL: AttributeAssignURL,
		AttributeRemoveURL: AttributeRemoveURL,

		OptionTableURL:  OptionTableURL,
		OptionAddURL:    OptionAddURL,
		OptionEditURL:   OptionEditURL,
		OptionDeleteURL: OptionDeleteURL,
		OptionDetailURL: OptionDetailURL,

		OptionValueTableURL:  OptionValueTableURL,
		OptionValueAddURL:    OptionValueAddURL,
		OptionValueEditURL:   OptionValueEditURL,
		OptionValueDeleteURL: OptionValueDeleteURL,
	}
}

// DefaultInventoryRoutes returns a Routes with every URL
// namespace-shifted from the service/product surface onto the inventory
// surface.
func DefaultInventoryRoutes() Routes {
	r := DefaultRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "masterlist"
	r.DashboardURL = ""
	shift := func(s string) string {
		s = strings.Replace(s, "/app/products/", "/app/inventory/products/", 1)
		s = strings.Replace(s, "/products/", "/inventory/products/", 1)
		s = strings.Replace(s, "/action/product/", "/action/inventory-product/", 1)
		return s
	}
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	r.AttachmentDownloadURL = shift(r.AttachmentDownloadURL)
	r.VariantTableURL = shift(r.VariantTableURL)
	r.VariantAssignURL = shift(r.VariantAssignURL)
	r.VariantEditURL = shift(r.VariantEditURL)
	r.VariantRemoveURL = shift(r.VariantRemoveURL)
	r.VariantDetailURL = shift(r.VariantDetailURL)
	r.VariantTabActionURL = shift(r.VariantTabActionURL)
	r.VariantImageUploadURL = shift(r.VariantImageUploadURL)
	r.VariantImageDeleteURL = shift(r.VariantImageDeleteURL)
	r.VariantAttachmentUploadURL = shift(r.VariantAttachmentUploadURL)
	r.VariantAttachmentDeleteURL = shift(r.VariantAttachmentDeleteURL)
	r.VariantStockDetailURL = shift(r.VariantStockDetailURL)
	r.VariantStockTabActionURL = shift(r.VariantStockTabActionURL)
	r.VariantStockAttachmentUploadURL = shift(r.VariantStockAttachmentUploadURL)
	r.VariantStockAttachmentDeleteURL = shift(r.VariantStockAttachmentDeleteURL)
	r.VariantSerialDetailURL = shift(r.VariantSerialDetailURL)
	r.AttributeTableURL = shift(r.AttributeTableURL)
	r.AttributeAssignURL = shift(r.AttributeAssignURL)
	r.AttributeRemoveURL = shift(r.AttributeRemoveURL)
	r.OptionTableURL = shift(r.OptionTableURL)
	r.OptionAddURL = shift(r.OptionAddURL)
	r.OptionEditURL = shift(r.OptionEditURL)
	r.OptionDeleteURL = shift(r.OptionDeleteURL)
	r.OptionDetailURL = shift(r.OptionDetailURL)
	r.OptionValueTableURL = shift(r.OptionValueTableURL)
	r.OptionValueAddURL = shift(r.OptionValueAddURL)
	r.OptionValueEditURL = shift(r.OptionValueEditURL)
	r.OptionValueDeleteURL = shift(r.OptionValueDeleteURL)
	return r
}

// DefaultSuppliesRoutes returns a Routes namespace-shifted onto
// the "Supplies" surface.
func DefaultSuppliesRoutes() Routes {
	r := DefaultRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "supplies"
	r.DashboardURL = ""
	shift := func(s string) string {
		s = strings.Replace(s, "/app/products/", "/app/inventory/supplies/", 1)
		s = strings.Replace(s, "/products/", "/inventory/supplies/", 1)
		s = strings.Replace(s, "/action/product/", "/action/inventory-supplies/", 1)
		return s
	}
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	r.AttachmentDownloadURL = shift(r.AttachmentDownloadURL)
	r.VariantTableURL = shift(r.VariantTableURL)
	r.VariantAssignURL = shift(r.VariantAssignURL)
	r.VariantEditURL = shift(r.VariantEditURL)
	r.VariantRemoveURL = shift(r.VariantRemoveURL)
	r.VariantDetailURL = shift(r.VariantDetailURL)
	r.VariantTabActionURL = shift(r.VariantTabActionURL)
	r.VariantImageUploadURL = shift(r.VariantImageUploadURL)
	r.VariantImageDeleteURL = shift(r.VariantImageDeleteURL)
	r.VariantAttachmentUploadURL = shift(r.VariantAttachmentUploadURL)
	r.VariantAttachmentDeleteURL = shift(r.VariantAttachmentDeleteURL)
	r.VariantStockDetailURL = shift(r.VariantStockDetailURL)
	r.VariantStockTabActionURL = shift(r.VariantStockTabActionURL)
	r.VariantStockAttachmentUploadURL = shift(r.VariantStockAttachmentUploadURL)
	r.VariantStockAttachmentDeleteURL = shift(r.VariantStockAttachmentDeleteURL)
	r.VariantSerialDetailURL = shift(r.VariantSerialDetailURL)
	r.AttributeTableURL = shift(r.AttributeTableURL)
	r.AttributeAssignURL = shift(r.AttributeAssignURL)
	r.AttributeRemoveURL = shift(r.AttributeRemoveURL)
	r.OptionTableURL = shift(r.OptionTableURL)
	r.OptionAddURL = shift(r.OptionAddURL)
	r.OptionEditURL = shift(r.OptionEditURL)
	r.OptionDeleteURL = shift(r.OptionDeleteURL)
	r.OptionDetailURL = shift(r.OptionDetailURL)
	r.OptionValueTableURL = shift(r.OptionValueTableURL)
	r.OptionValueAddURL = shift(r.OptionValueAddURL)
	r.OptionValueEditURL = shift(r.OptionValueEditURL)
	r.OptionValueDeleteURL = shift(r.OptionValueDeleteURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// product routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"product.dashboard":   r.DashboardURL,
		"product.list":        r.ListURL,
		"product.table":       r.TableURL,
		"product.detail":      r.DetailURL,
		"product.add":         r.AddURL,
		"product.edit":        r.EditURL,
		"product.delete":      r.DeleteURL,
		"product.bulk_delete": r.BulkDeleteURL,

		"product.set_status":      r.SetStatusURL,
		"product.bulk_set_status": r.BulkSetStatusURL,

		"product.tab_action": r.TabActionURL,

		"product.attachment.upload":   r.AttachmentUploadURL,
		"product.attachment.delete":   r.AttachmentDeleteURL,
		"product.attachment.download": r.AttachmentDownloadURL,

		"product.variant.table":  r.VariantTableURL,
		"product.variant.assign": r.VariantAssignURL,
		"product.variant.edit":   r.VariantEditURL,
		"product.variant.remove": r.VariantRemoveURL,

		"product.variant.detail":     r.VariantDetailURL,
		"product.variant.tab_action": r.VariantTabActionURL,

		"product.variant.image.upload": r.VariantImageUploadURL,
		"product.variant.image.delete": r.VariantImageDeleteURL,

		"product.variant.attachment.upload": r.VariantAttachmentUploadURL,
		"product.variant.attachment.delete": r.VariantAttachmentDeleteURL,

		"product.variant.stock.detail":     r.VariantStockDetailURL,
		"product.variant.stock.tab_action": r.VariantStockTabActionURL,

		"product.variant.stock.attachment.upload": r.VariantStockAttachmentUploadURL,
		"product.variant.stock.attachment.delete": r.VariantStockAttachmentDeleteURL,

		"product.variant.serial.detail": r.VariantSerialDetailURL,

		"product.attribute.table":  r.AttributeTableURL,
		"product.attribute.assign": r.AttributeAssignURL,
		"product.attribute.remove": r.AttributeRemoveURL,

		"product.option.table":  r.OptionTableURL,
		"product.option.add":    r.OptionAddURL,
		"product.option.edit":   r.OptionEditURL,
		"product.option.delete": r.OptionDeleteURL,
		"product.option.detail": r.OptionDetailURL,

		"product.option_value.table":  r.OptionValueTableURL,
		"product.option_value.add":    r.OptionValueAddURL,
		"product.option_value.edit":   r.OptionValueEditURL,
		"product.option_value.delete": r.OptionValueDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// LineRoutes
// ---------------------------------------------------------------------------

// LineRoutes holds all route paths for product line views and actions.
type LineRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

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

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultLineRoutes returns a LineRoutes populated from the
// package-level route constants defined in this file.
func DefaultLineRoutes() LineRoutes {
	return LineRoutes{
		ActiveNav:    "service",
		ActiveSubNav: "product-lines",

		DashboardURL:     LineDashboardURL,
		ListURL:          LineListURL,
		TableURL:         LineTableURL,
		DetailURL:        LineDetailURL,
		AddURL:           LineAddURL,
		EditURL:          LineEditURL,
		DeleteURL:        LineDeleteURL,
		BulkDeleteURL:    LineBulkDeleteURL,
		SetStatusURL:     LineSetStatusURL,
		BulkSetStatusURL: LineBulkSetStatusURL,
		TabActionURL:     LineTabActionURL,

		AttachmentUploadURL: LineAttachmentUploadURL,
		AttachmentDeleteURL: LineAttachmentDeleteURL,
	}
}

// DefaultLineInventoryRoutes returns a LineRoutes with every URL
// namespace-shifted from the service surface onto the inventory accordion
// namespace.
func DefaultLineInventoryRoutes() LineRoutes {
	r := DefaultLineRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "product-lines-active"
	shift := func(s string) string {
		s = strings.Replace(s, "/app/product-lines/", "/app/inventory/product-lines/", 1)
		s = strings.Replace(s, "/product-lines/", "/inventory/product-lines/", 1)
		s = strings.Replace(s, "/action/product-line/", "/action/inventory-product-line/", 1)
		return s
	}
	r.DashboardURL = shift(r.DashboardURL)
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// product line routes.
func (r LineRoutes) RouteMap() map[string]string {
	return map[string]string{
		"product_line.dashboard":         r.DashboardURL,
		"product_line.list":              r.ListURL,
		"product_line.table":             r.TableURL,
		"product_line.detail":            r.DetailURL,
		"product_line.add":               r.AddURL,
		"product_line.edit":              r.EditURL,
		"product_line.delete":            r.DeleteURL,
		"product_line.bulk_delete":       r.BulkDeleteURL,
		"product_line.set_status":        r.SetStatusURL,
		"product_line.bulk_set_status":   r.BulkSetStatusURL,
		"product_line.tab_action":        r.TabActionURL,
		"product_line.attachment.upload": r.AttachmentUploadURL,
		"product_line.attachment.delete": r.AttachmentDeleteURL,
	}
}
