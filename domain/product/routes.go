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
	ProductListURL   = "/products/list/{status}"
	ProductTableURL  = "/action/product/table/{status}"
	ProductDetailURL = "/products/detail/{id}"

	// Service dashboard — services are products filtered to product_kind="service".
	// The dashboard sits ABOVE the service-mount product list at this URL.
	ServiceDashboardURL = "/services/dashboard"

	// Product action routes
	ProductAddURL        = "/action/product/add"
	ProductEditURL       = "/action/product/edit/{id}"
	ProductDeleteURL     = "/action/product/delete"
	ProductBulkDeleteURL = "/action/product/bulk-delete"

	// Product status routes
	ProductSetStatusURL     = "/action/product/set-status"
	ProductBulkSetStatusURL = "/action/product/bulk-set-status"

	// Product detail tab action route
	ProductTabActionURL          = "/action/product/detail/{id}/tab/{tab}"
	ProductAttachmentUploadURL   = "/action/product/detail/{id}/attachments/upload"
	ProductAttachmentDeleteURL   = "/action/product/detail/{id}/attachments/delete"
	ProductAttachmentDownloadURL = "/action/product/detail/{id}/attachments/download"

	// Product variant routes (within product detail)
	ProductVariantTableURL  = "/action/product/detail/{id}/variants/table"
	ProductVariantAssignURL = "/action/product/detail/{id}/variants/assign"
	ProductVariantEditURL   = "/action/product/detail/{id}/variants/edit/{vid}"
	ProductVariantRemoveURL = "/action/product/detail/{id}/variants/remove"

	// Product attribute routes (within product detail)
	ProductAttributeTableURL  = "/action/product/detail/{id}/attributes/table"
	ProductAttributeAssignURL = "/action/product/detail/{id}/attributes/assign"
	ProductAttributeRemoveURL = "/action/product/detail/{id}/attributes/remove"

	// Product option routes (within product detail)
	ProductOptionTableURL  = "/action/product/detail/{id}/options/table"
	ProductOptionAddURL    = "/action/product/detail/{id}/options/add"
	ProductOptionEditURL   = "/action/product/detail/{id}/options/edit/{oid}"
	ProductOptionDeleteURL = "/action/product/detail/{id}/options/delete"

	// Product option detail page (option values management)
	ProductOptionDetailURL = "/products/detail/{id}/option/{oid}"

	// Product variant detail page (variant info, pricing, stock, audit, images)
	ProductVariantDetailURL    = "/products/detail/{id}/variant/{vid}"
	ProductVariantTabActionURL = "/action/product/detail/{id}/variant/{vid}/tab/{tab}"

	// Product variant image routes (upload/delete within variant detail)
	ProductVariantImageUploadURL = "/action/product/detail/{id}/variant/{vid}/images/upload"
	ProductVariantImageDeleteURL = "/action/product/detail/{id}/variant/{vid}/images/delete"

	// Product variant attachment routes
	ProductVariantAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/attachments/upload"
	ProductVariantAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/attachments/delete"

	// Product variant stock detail (inventory item within variant context)
	ProductVariantStockDetailURL    = "/products/detail/{id}/variant/{vid}/stock/{iid}"
	ProductVariantStockTabActionURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/tab/{tab}"

	// Product variant stock attachment routes
	ProductVariantStockAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/upload"
	ProductVariantStockAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/delete"

	// Inventory serial detail (individual serial within inventory item)
	ProductVariantSerialDetailURL = "/products/detail/{id}/variant/{vid}/stock/{iid}/serial/{sid}"

	// Product option value routes (within product option)
	ProductOptionValueTableURL  = "/action/product/detail/{id}/options/{oid}/values/table"
	ProductOptionValueAddURL    = "/action/product/detail/{id}/options/{oid}/values/add"
	ProductOptionValueEditURL   = "/action/product/detail/{id}/options/{oid}/values/edit/{vid}"
	ProductOptionValueDeleteURL = "/action/product/detail/{id}/options/{oid}/values/delete"

	// Product line routes
	ProductLineDashboardURL        = "/product-lines/dashboard"
	ProductLineListURL             = "/product-lines/list/{status}"
	ProductLineTableURL            = "/action/product-line/table/{status}"
	ProductLineDetailURL           = "/product-lines/detail/{id}"
	ProductLineAddURL              = "/action/product-line/add"
	ProductLineEditURL             = "/action/product-line/edit/{id}"
	ProductLineDeleteURL           = "/action/product-line/delete"
	ProductLineBulkDeleteURL       = "/action/product-line/bulk-delete"
	ProductLineSetStatusURL        = "/action/product-line/set-status"
	ProductLineBulkSetStatusURL    = "/action/product-line/bulk-set-status"
	ProductLineTabActionURL        = "/action/product-line/{id}/tab/{tab}"
	ProductLineAttachmentUploadURL = "/action/product-line/{id}/attachments/upload"
	ProductLineAttachmentDeleteURL = "/action/product-line/{id}/attachments/delete"

	// Resource routes (person or equipment linked to a Product for billing)
	ResourceListURL          = "/resources/list/{status}"
	ResourceTableURL         = "/action/resource/table/{status}"
	ResourceDetailURL        = "/resources/detail/{id}"
	ResourceAddURL           = "/action/resource/add"
	ResourceEditURL          = "/action/resource/edit/{id}"
	ResourceDeleteURL        = "/action/resource/delete"
	ResourceBulkDeleteURL    = "/action/resource/bulk-delete"
	ResourceSetStatusURL     = "/action/resource/set-status"
	ResourceBulkSetStatusURL = "/action/resource/bulk-set-status"

	// Price List routes — canonical home is the inventory accordion (/app/inventory/price-lists/*)
	PriceListListURL       = "/inventory/price-lists/list/{status}"
	PriceListTableURL      = "/action/inventory-price-list/table/{status}"
	PriceListDetailURL     = "/inventory/price-lists/detail/{id}"
	PriceListAddURL        = "/action/inventory-price-list/add"
	PriceListEditURL       = "/action/inventory-price-list/edit/{id}"
	PriceListDeleteURL     = "/action/inventory-price-list/delete"
	PriceListBulkDeleteURL = "/action/inventory-price-list/bulk-delete"

	PriceListTabActionURL        = "/action/inventory-price-list/{id}/tab/{tab}"
	PriceListAttachmentUploadURL = "/action/inventory-price-list/{id}/attachments/upload"
	PriceListAttachmentDeleteURL = "/action/inventory-price-list/{id}/attachments/delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/inventory-price-list/{id}/products/add"
	PriceProductDeleteURL = "/action/inventory-price-list/{id}/products/delete"
)

// ---------------------------------------------------------------------------
// ProductRoutes
// ---------------------------------------------------------------------------

// ProductRoutes holds all route paths for product views and actions,
// including variant, option, attribute, image, stock, and serial sub-routes.
type ProductRoutes struct {
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

// DefaultProductRoutes returns a ProductRoutes populated from the package-level
// route constants defined in this file.
func DefaultProductRoutes() ProductRoutes {
	return ProductRoutes{
		ActiveNav:    "inventory",
		ActiveSubNav: "masterlist",

		// Default to the service dashboard URL — only meaningful for the
		// service-flavoured mount. Inventory/supplies mounts overwrite or
		// ignore this field.
		DashboardURL: ServiceDashboardURL,

		ListURL:       ProductListURL,
		TableURL:      ProductTableURL,
		DetailURL:     ProductDetailURL,
		AddURL:        ProductAddURL,
		EditURL:       ProductEditURL,
		DeleteURL:     ProductDeleteURL,
		BulkDeleteURL: ProductBulkDeleteURL,

		SetStatusURL:     ProductSetStatusURL,
		BulkSetStatusURL: ProductBulkSetStatusURL,

		TabActionURL: ProductTabActionURL,

		AttachmentUploadURL:   ProductAttachmentUploadURL,
		AttachmentDeleteURL:   ProductAttachmentDeleteURL,
		AttachmentDownloadURL: ProductAttachmentDownloadURL,

		VariantTableURL:  ProductVariantTableURL,
		VariantAssignURL: ProductVariantAssignURL,
		VariantEditURL:   ProductVariantEditURL,
		VariantRemoveURL: ProductVariantRemoveURL,

		VariantDetailURL:    ProductVariantDetailURL,
		VariantTabActionURL: ProductVariantTabActionURL,

		VariantImageUploadURL: ProductVariantImageUploadURL,
		VariantImageDeleteURL: ProductVariantImageDeleteURL,

		VariantAttachmentUploadURL: ProductVariantAttachmentUploadURL,
		VariantAttachmentDeleteURL: ProductVariantAttachmentDeleteURL,

		VariantStockDetailURL:    ProductVariantStockDetailURL,
		VariantStockTabActionURL: ProductVariantStockTabActionURL,

		VariantStockAttachmentUploadURL: ProductVariantStockAttachmentUploadURL,
		VariantStockAttachmentDeleteURL: ProductVariantStockAttachmentDeleteURL,

		VariantSerialDetailURL: ProductVariantSerialDetailURL,

		AttributeTableURL:  ProductAttributeTableURL,
		AttributeAssignURL: ProductAttributeAssignURL,
		AttributeRemoveURL: ProductAttributeRemoveURL,

		OptionTableURL:  ProductOptionTableURL,
		OptionAddURL:    ProductOptionAddURL,
		OptionEditURL:   ProductOptionEditURL,
		OptionDeleteURL: ProductOptionDeleteURL,
		OptionDetailURL: ProductOptionDetailURL,

		OptionValueTableURL:  ProductOptionValueTableURL,
		OptionValueAddURL:    ProductOptionValueAddURL,
		OptionValueEditURL:   ProductOptionValueEditURL,
		OptionValueDeleteURL: ProductOptionValueDeleteURL,
	}
}

// DefaultProductInventoryRoutes returns a ProductRoutes with every URL
// namespace-shifted from the service/product surface onto the inventory
// surface.
func DefaultProductInventoryRoutes() ProductRoutes {
	r := DefaultProductRoutes()
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

// DefaultProductSuppliesRoutes returns a ProductRoutes namespace-shifted onto
// the "Supplies" surface.
func DefaultProductSuppliesRoutes() ProductRoutes {
	r := DefaultProductRoutes()
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
func (r ProductRoutes) RouteMap() map[string]string {
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
// ProductLineRoutes
// ---------------------------------------------------------------------------

// ProductLineRoutes holds all route paths for product line views and actions.
type ProductLineRoutes struct {
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

// DefaultProductLineRoutes returns a ProductLineRoutes populated from the
// package-level route constants defined in this file.
func DefaultProductLineRoutes() ProductLineRoutes {
	return ProductLineRoutes{
		ActiveNav:    "service",
		ActiveSubNav: "product-lines",

		DashboardURL:     ProductLineDashboardURL,
		ListURL:          ProductLineListURL,
		TableURL:         ProductLineTableURL,
		DetailURL:        ProductLineDetailURL,
		AddURL:           ProductLineAddURL,
		EditURL:          ProductLineEditURL,
		DeleteURL:        ProductLineDeleteURL,
		BulkDeleteURL:    ProductLineBulkDeleteURL,
		SetStatusURL:     ProductLineSetStatusURL,
		BulkSetStatusURL: ProductLineBulkSetStatusURL,
		TabActionURL:     ProductLineTabActionURL,

		AttachmentUploadURL: ProductLineAttachmentUploadURL,
		AttachmentDeleteURL: ProductLineAttachmentDeleteURL,
	}
}

// DefaultProductLineInventoryRoutes returns a ProductLineRoutes with every URL
// namespace-shifted from the service surface onto the inventory accordion
// namespace.
func DefaultProductLineInventoryRoutes() ProductLineRoutes {
	r := DefaultProductLineRoutes()
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
func (r ProductLineRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// ResourceRoutes
// ---------------------------------------------------------------------------

// ResourceRoutes holds all route paths for resource views and actions.
// A resource links a person or equipment to a Product for billing purposes.
type ResourceRoutes struct {
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

// DefaultResourceRoutes returns a ResourceRoutes populated from the package-level
// route constants defined in this file.
func DefaultResourceRoutes() ResourceRoutes {
	return ResourceRoutes{
		ActiveNav:        "service",
		ActiveSubNav:     "resources-active",
		ListURL:          ResourceListURL,
		TableURL:         ResourceTableURL,
		DetailURL:        ResourceDetailURL,
		AddURL:           ResourceAddURL,
		EditURL:          ResourceEditURL,
		DeleteURL:        ResourceDeleteURL,
		BulkDeleteURL:    ResourceBulkDeleteURL,
		SetStatusURL:     ResourceSetStatusURL,
		BulkSetStatusURL: ResourceBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// resource routes.
func (r ResourceRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// PriceListRoutes
// ---------------------------------------------------------------------------

// PriceListRoutes holds all route paths for price list views and actions,
// including price product sub-routes.
type PriceListRoutes struct {
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

// DefaultPriceListRoutes returns a PriceListRoutes populated from the
// package-level route constants defined in this file.
func DefaultPriceListRoutes() PriceListRoutes {
	return PriceListRoutes{
		ListURL:       PriceListListURL,
		TableURL:      PriceListTableURL,
		DetailURL:     PriceListDetailURL,
		AddURL:        PriceListAddURL,
		EditURL:       PriceListEditURL,
		DeleteURL:     PriceListDeleteURL,
		BulkDeleteURL: PriceListBulkDeleteURL,

		TabActionURL: PriceListTabActionURL,

		AttachmentUploadURL: PriceListAttachmentUploadURL,
		AttachmentDeleteURL: PriceListAttachmentDeleteURL,

		PriceProductAddURL:    PriceProductAddURL,
		PriceProductDeleteURL: PriceProductDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price list routes.
func (r PriceListRoutes) RouteMap() map[string]string {
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
