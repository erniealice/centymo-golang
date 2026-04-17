package centymo

import "strings"

// Three-level routing system for centymo views:
//
// Level 1: Generic defaults from Go consts (this file).
//   DefaultXxxRoutes() constructors return structs populated from the route
//   constants defined in routes.go. These are sensible defaults that work
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

// ProductRoutes holds all route paths for product views and actions,
// including variant, option, attribute, image, stock, and serial sub-routes.
type ProductRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

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
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

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
// route constants defined in routes.go.
func DefaultProductRoutes() ProductRoutes {
	return ProductRoutes{
		ActiveNav:    "inventory",
		ActiveSubNav: "masterlist",

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

		AttachmentUploadURL: ProductAttachmentUploadURL,
		AttachmentDeleteURL: ProductAttachmentDeleteURL,

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
// surface. Used for the inventory-flavoured Product list mount — keeps
// both mounts on distinct URLs so the stdlib ServeMux does not panic on
// duplicate registrations when both modules register against the same mux.
//
// Shift rules:
//   - "/app/products/*"  → "/app/inventory/products/*"
//   - "/action/product/*" → "/action/inventory-product/*"
//
// Lyngua `product_inventory` route blocks can still override individual
// URLs on top of this baseline.
func DefaultProductInventoryRoutes() ProductRoutes {
	r := DefaultProductRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "masterlist"
	shift := func(s string) string {
		s = strings.Replace(s, "/app/products/", "/app/inventory/products/", 1)
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

		"product.attachment.upload": r.AttachmentUploadURL,
		"product.attachment.delete": r.AttachmentDeleteURL,

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
// package-level route constants defined in routes.go.
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
// namespace. Inventory-mount variant of ProductLine routes — namespaces every
// URL under /app/inventory/product-lines/ so both accordions can carry
// ProductLine without collision.
//
// Shift rules:
//   - "/app/product-lines/"   → "/app/inventory/product-lines/"
//   - "/action/product-line/" → "/action/inventory-product-line/"
func DefaultProductLineInventoryRoutes() ProductLineRoutes {
	r := DefaultProductLineRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "product-lines-active"
	shift := func(s string) string {
		s = strings.Replace(s, "/app/product-lines/", "/app/inventory/product-lines/", 1)
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

// InventoryRoutes holds all route paths for inventory views and actions,
// including serial, transaction, depreciation, dashboard, and movement sub-routes.
type InventoryRoutes struct {
	DashboardURL  string `json:"dashboard_url"`
	ListURL       string `json:"list_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	DetailURL     string `json:"detail_url"`
	TableURL      string `json:"table_url"`

	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	TabActionURL string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Movement routes
	MovementsURL       string `json:"movements_url"`
	MovementsTableURL  string `json:"movements_table_url"`
	MovementsExportURL string `json:"movements_export_url"`

	// Serial routes
	SerialTableURL  string `json:"serial_table_url"`
	SerialAssignURL string `json:"serial_assign_url"`
	SerialEditURL   string `json:"serial_edit_url"`
	SerialRemoveURL string `json:"serial_remove_url"`

	// Transaction routes
	TransactionTableURL  string `json:"transaction_table_url"`
	TransactionAssignURL string `json:"transaction_assign_url"`

	// Depreciation routes
	DepreciationAssignURL string `json:"depreciation_assign_url"`
	DepreciationEditURL   string `json:"depreciation_edit_url"`

	// Attribute routes
	AttributeTableURL string `json:"attribute_table_url"`

	// Dashboard partial routes
	DashboardStatsURL     string `json:"dashboard_stats_url"`
	DashboardChartURL     string `json:"dashboard_chart_url"`
	DashboardMovementsURL string `json:"dashboard_movements_url"`
	DashboardAlertsURL    string `json:"dashboard_alerts_url"`

	// Product-context detail routes
	ProductDetailURL    string `json:"product_detail_url"`
	ProductTabActionURL string `json:"product_tab_action_url"`
}

// DefaultInventoryRoutes returns an InventoryRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultInventoryRoutes() InventoryRoutes {
	return InventoryRoutes{
		DashboardURL:  InventoryDashboardURL,
		ListURL:       InventoryListURL,
		AddURL:        InventoryAddURL,
		EditURL:       InventoryEditURL,
		DeleteURL:     InventoryDeleteURL,
		BulkDeleteURL: InventoryBulkDeleteURL,
		DetailURL:     InventoryDetailURL,
		TableURL:      InventoryTableURL,

		SetStatusURL:     InventorySetStatusURL,
		BulkSetStatusURL: InventoryBulkSetStatusURL,

		TabActionURL: InventoryTabActionURL,

		AttachmentUploadURL: InventoryAttachmentUploadURL,
		AttachmentDeleteURL: InventoryAttachmentDeleteURL,

		MovementsURL:       InventoryMovementsURL,
		MovementsTableURL:  InventoryMovementsTableURL,
		MovementsExportURL: InventoryMovementsExportURL,

		SerialTableURL:  InventorySerialTableURL,
		SerialAssignURL: InventorySerialAssignURL,
		SerialEditURL:   InventorySerialEditURL,
		SerialRemoveURL: InventorySerialRemoveURL,

		TransactionTableURL:  InventoryTransactionTableURL,
		TransactionAssignURL: InventoryTransactionAssignURL,

		DepreciationAssignURL: InventoryDepreciationAssignURL,
		DepreciationEditURL:   InventoryDepreciationEditURL,

		AttributeTableURL: InventoryAttributeTableURL,

		DashboardStatsURL:     InventoryDashboardStatsURL,
		DashboardChartURL:     InventoryDashboardChartURL,
		DashboardMovementsURL: InventoryDashboardMovementsURL,
		DashboardAlertsURL:    InventoryDashboardAlertsURL,

		ProductDetailURL:    InventoryProductDetailURL,
		ProductTabActionURL: InventoryProductTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// inventory routes.
func (r InventoryRoutes) RouteMap() map[string]string {
	return map[string]string{
		"inventory.dashboard":   r.DashboardURL,
		"inventory.list":        r.ListURL,
		"inventory.add":         r.AddURL,
		"inventory.edit":        r.EditURL,
		"inventory.delete":      r.DeleteURL,
		"inventory.bulk_delete": r.BulkDeleteURL,
		"inventory.detail":      r.DetailURL,
		"inventory.table":       r.TableURL,

		"inventory.set_status":      r.SetStatusURL,
		"inventory.bulk_set_status": r.BulkSetStatusURL,

		"inventory.tab_action": r.TabActionURL,

		"inventory.attachment.upload": r.AttachmentUploadURL,
		"inventory.attachment.delete": r.AttachmentDeleteURL,

		"inventory.movements":        r.MovementsURL,
		"inventory.movements.table":  r.MovementsTableURL,
		"inventory.movements.export": r.MovementsExportURL,

		"inventory.serial.table":  r.SerialTableURL,
		"inventory.serial.assign": r.SerialAssignURL,
		"inventory.serial.edit":   r.SerialEditURL,
		"inventory.serial.remove": r.SerialRemoveURL,

		"inventory.transaction.table":  r.TransactionTableURL,
		"inventory.transaction.assign": r.TransactionAssignURL,

		"inventory.depreciation.assign": r.DepreciationAssignURL,
		"inventory.depreciation.edit":   r.DepreciationEditURL,

		"inventory.attribute.table": r.AttributeTableURL,

		"inventory.dashboard.stats":     r.DashboardStatsURL,
		"inventory.dashboard.chart":     r.DashboardChartURL,
		"inventory.dashboard.movements": r.DashboardMovementsURL,
		"inventory.dashboard.alerts":    r.DashboardAlertsURL,

		"inventory.product.detail":     r.ProductDetailURL,
		"inventory.product.tab_action": r.ProductTabActionURL,
	}
}

// RevenueRoutes holds all route paths for revenue views and actions,
// including line item and payment sub-routes.
type RevenueRoutes struct {
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

	TabActionURL string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Line item routes
	LineItemTableURL    string `json:"line_item_table_url"`
	LineItemAddURL      string `json:"line_item_add_url"`
	LineItemEditURL     string `json:"line_item_edit_url"`
	LineItemRemoveURL   string `json:"line_item_remove_url"`
	LineItemDiscountURL string `json:"line_item_discount_url"`

	// Payment routes
	PaymentTableURL  string `json:"payment_table_url"`
	PaymentAddURL    string `json:"payment_add_url"`
	PaymentEditURL   string `json:"payment_edit_url"`
	PaymentRemoveURL string `json:"payment_remove_url"`

	// Report routes
	RevenueSummaryURL string `json:"revenue_summary_url"`

	// Document generation routes
	InvoiceDownloadURL string `json:"invoice_download_url"`

	// Send email route
	SendEmailURL string `json:"send_email_url"`

	// Settings routes (template management)
	SettingsTemplatesURL       string `json:"settings_templates_url"`
	SettingsTemplateUploadURL  string `json:"settings_template_upload_url"`
	SettingsTemplateDeleteURL  string `json:"settings_template_delete_url"`
	SettingsTemplateDefaultURL string `json:"settings_template_default_url"`

	// Client search for revenue form autocomplete
	SearchClientURL string `json:"search_client_url"`

	// Subscription search for revenue form autocomplete
	SearchSubscriptionURL string `json:"search_subscription_url"`

	// Location search for revenue form autocomplete
	SearchLocationURL string `json:"search_location_url"`

	// Product/service search for revenue line item autocomplete
	SearchProductURL string `json:"search_product_url"`

	// Price lookup for revenue line item (product_id + location_id + date → price)
	PriceLookupURL string `json:"price_lookup_url"`
}

// DefaultRevenueRoutes returns a RevenueRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultRevenueRoutes() RevenueRoutes {
	return RevenueRoutes{
		DashboardURL:     RevenueDashboardURL,
		ListURL:          RevenueListURL,
		TableURL:         RevenueTableURL,
		DetailURL:        RevenueDetailURL,
		AddURL:           RevenueAddURL,
		EditURL:          RevenueEditURL,
		DeleteURL:        RevenueDeleteURL,
		BulkDeleteURL:    RevenueBulkDeleteURL,
		SetStatusURL:     RevenueSetStatusURL,
		BulkSetStatusURL: RevenueBulkSetStatusURL,

		TabActionURL: RevenueTabActionURL,

		AttachmentUploadURL: RevenueAttachmentUploadURL,
		AttachmentDeleteURL: RevenueAttachmentDeleteURL,

		LineItemTableURL:    RevenueLineItemTableURL,
		LineItemAddURL:      RevenueLineItemAddURL,
		LineItemEditURL:     RevenueLineItemEditURL,
		LineItemRemoveURL:   RevenueLineItemRemoveURL,
		LineItemDiscountURL: RevenueLineItemDiscountURL,

		PaymentTableURL:  RevenuePaymentTableURL,
		PaymentAddURL:    RevenuePaymentAddURL,
		PaymentEditURL:   RevenuePaymentEditURL,
		PaymentRemoveURL: RevenuePaymentRemoveURL,

		RevenueSummaryURL:          RevenueSummaryURL,
		InvoiceDownloadURL:         RevenueInvoiceDownloadURL,
		SendEmailURL:               RevenueEmailURL,
		SettingsTemplatesURL:       RevenueSettingsTemplatesURL,
		SettingsTemplateUploadURL:  RevenueSettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  RevenueSettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: RevenueSettingsTemplateDefaultURL,
		SearchClientURL:            RevenueSearchClientURL,
		SearchSubscriptionURL:      RevenueSearchSubscriptionURL,
		SearchLocationURL:          RevenueSearchLocationURL,
		SearchProductURL:           RevenueSearchProductURL,
		PriceLookupURL:             RevenuePriceLookupURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// revenue routes.
func (r RevenueRoutes) RouteMap() map[string]string {
	return map[string]string{
		"revenue.dashboard":       r.DashboardURL,
		"revenue.list":            r.ListURL,
		"revenue.table":           r.TableURL,
		"revenue.detail":          r.DetailURL,
		"revenue.add":             r.AddURL,
		"revenue.edit":            r.EditURL,
		"revenue.delete":          r.DeleteURL,
		"revenue.bulk_delete":     r.BulkDeleteURL,
		"revenue.set_status":      r.SetStatusURL,
		"revenue.bulk_set_status": r.BulkSetStatusURL,

		"revenue.tab_action": r.TabActionURL,

		"revenue.attachment.upload": r.AttachmentUploadURL,
		"revenue.attachment.delete": r.AttachmentDeleteURL,

		"revenue.line_item.table":    r.LineItemTableURL,
		"revenue.line_item.add":      r.LineItemAddURL,
		"revenue.line_item.edit":     r.LineItemEditURL,
		"revenue.line_item.remove":   r.LineItemRemoveURL,
		"revenue.line_item.discount": r.LineItemDiscountURL,

		"revenue.payment.table":  r.PaymentTableURL,
		"revenue.payment.add":    r.PaymentAddURL,
		"revenue.payment.edit":   r.PaymentEditURL,
		"revenue.payment.remove": r.PaymentRemoveURL,

		"revenue.summary":                 r.RevenueSummaryURL,
		"revenue.invoice_download":        r.InvoiceDownloadURL,
		"revenue.send_email":              r.SendEmailURL,
		"revenue.settings.templates":        r.SettingsTemplatesURL,
		"revenue.settings.template_upload":  r.SettingsTemplateUploadURL,
		"revenue.settings.template_delete":  r.SettingsTemplateDeleteURL,
		"revenue.settings.template_default": r.SettingsTemplateDefaultURL,
		"revenue.search_client":        r.SearchClientURL,
		"revenue.search.subscriptions": r.SearchSubscriptionURL,
		"revenue.search.locations":     r.SearchLocationURL,
		"revenue.search.products":      r.SearchProductURL,
		"revenue.price_lookup":         r.PriceLookupURL,
	}
}

// ExpenditureRoutes holds all route paths for expenditure views (purchase + expense).
type ExpenditureRoutes struct {
	PurchaseListURL      string `json:"purchase_list_url"`
	PurchaseDashboardURL string `json:"purchase_dashboard_url"`
	ExpenseListURL       string `json:"expense_list_url"`
	ExpenseDashboardURL  string `json:"expense_dashboard_url"`

	// Report routes
	PurchasesSummaryURL string `json:"purchases_summary_url"`
	ExpensesSummaryURL  string `json:"expenses_summary_url"`

	// Settings (template management)
	SettingsTemplatesURL       string `json:"settings_templates_url"`
	SettingsTemplateUploadURL  string `json:"settings_template_upload_url"`
	SettingsTemplateDeleteURL  string `json:"settings_template_delete_url"`
	SettingsTemplateDefaultURL string `json:"settings_template_default_url"`

	// Expense CRUD action routes
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
	SetStatusURL string `json:"set_status_url"`
	DetailURL    string `json:"detail_url"`
	TableURL     string `json:"table_url"`
	TabActionURL string `json:"tab_action_url"`

	// Expense line item action routes
	LineItemAddURL    string `json:"line_item_add_url"`
	LineItemEditURL   string `json:"line_item_edit_url"`
	LineItemRemoveURL string `json:"line_item_remove_url"`
	LineItemTableURL  string `json:"line_item_table_url"`

	// Pay action route (creates pre-linked disbursement)
	PayURL string `json:"pay_url"`

	// Expense category CRUD routes
	ExpenseCategoryListURL   string `json:"expense_category_list_url"`
	ExpenseCategoryAddURL    string `json:"expense_category_add_url"`
	ExpenseCategoryEditURL   string `json:"expense_category_edit_url"`
	ExpenseCategoryDeleteURL string `json:"expense_category_delete_url"`
	ExpenseCategoryTableURL  string `json:"expense_category_table_url"`

	// Purchase Order routes
	PurchaseOrderListURL      string `json:"purchase_order_list_url"`
	PurchaseOrderDetailURL    string `json:"purchase_order_detail_url"`
	PurchaseOrderAddURL       string `json:"purchase_order_add_url"`
	PurchaseOrderEditURL      string `json:"purchase_order_edit_url"`
	PurchaseOrderDeleteURL    string `json:"purchase_order_delete_url"`
	PurchaseOrderSetStatusURL string `json:"purchase_order_set_status_url"`
	PurchaseOrderTableURL     string `json:"purchase_order_table_url"`
	PurchaseOrderTabActionURL string `json:"purchase_order_tab_action_url"`

	// Purchase Order line item routes (within PO detail)
	PurchaseOrderLineItemTableURL  string `json:"purchase_order_line_item_table_url"`
	PurchaseOrderLineItemAddURL    string `json:"purchase_order_line_item_add_url"`
	PurchaseOrderLineItemEditURL   string `json:"purchase_order_line_item_edit_url"`
	PurchaseOrderLineItemRemoveURL string `json:"purchase_order_line_item_remove_url"`

	// Purchase Order receipt action
	PurchaseOrderConfirmReceiptURL string `json:"purchase_order_confirm_receipt_url"`
}

// DefaultExpenditureRoutes returns an ExpenditureRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultExpenditureRoutes() ExpenditureRoutes {
	return ExpenditureRoutes{
		PurchaseListURL:      ExpenditurePurchaseListURL,
		PurchaseDashboardURL: ExpenditurePurchaseDashboardURL,
		ExpenseListURL:       ExpenditureExpenseListURL,
		ExpenseDashboardURL:  ExpenditureExpenseDashboardURL,

		PurchasesSummaryURL: PurchasesSummaryURL,
		ExpensesSummaryURL:  ExpensesSummaryURL,

		SettingsTemplatesURL:       ExpenditureSettingsTemplatesURL,
		SettingsTemplateUploadURL:  ExpenditureSettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  ExpenditureSettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: ExpenditureSettingsTemplateDefaultURL,

		AddURL:       ExpenditureExpenseAddURL,
		EditURL:      ExpenditureExpenseEditURL,
		DeleteURL:    ExpenditureExpenseDeleteURL,
		SetStatusURL: ExpenditureExpenseSetStatusURL,
		DetailURL:    ExpenditureExpenseDetailURL,
		TableURL:     ExpenditureExpenseTableURL,
		TabActionURL: ExpenditureExpenseTabActionURL,

		LineItemAddURL:    ExpenditureExpenseLineItemAddURL,
		LineItemEditURL:   ExpenditureExpenseLineItemEditURL,
		LineItemRemoveURL: ExpenditureExpenseLineItemRemoveURL,
		LineItemTableURL:  ExpenditureExpenseLineItemTableURL,

		PayURL: ExpenditureExpensePayURL,

		ExpenseCategoryListURL:   ExpenditureExpenseCategoryListURL,
		ExpenseCategoryAddURL:    ExpenditureExpenseCategoryAddURL,
		ExpenseCategoryEditURL:   ExpenditureExpenseCategoryEditURL,
		ExpenseCategoryDeleteURL: ExpenditureExpenseCategoryDeleteURL,
		ExpenseCategoryTableURL:  ExpenditureExpenseCategoryTableURL,

		PurchaseOrderListURL:      PurchaseOrderListURL,
		PurchaseOrderDetailURL:    PurchaseOrderDetailURL,
		PurchaseOrderAddURL:       PurchaseOrderAddURL,
		PurchaseOrderEditURL:      PurchaseOrderEditURL,
		PurchaseOrderDeleteURL:    PurchaseOrderDeleteURL,
		PurchaseOrderSetStatusURL: PurchaseOrderSetStatusURL,
		PurchaseOrderTableURL:     PurchaseOrderTableURL,
		PurchaseOrderTabActionURL: PurchaseOrderTabActionURL,

		PurchaseOrderLineItemTableURL:  PurchaseOrderLineItemTableURL,
		PurchaseOrderLineItemAddURL:    PurchaseOrderLineItemAddURL,
		PurchaseOrderLineItemEditURL:   PurchaseOrderLineItemEditURL,
		PurchaseOrderLineItemRemoveURL: PurchaseOrderLineItemRemoveURL,

		PurchaseOrderConfirmReceiptURL: PurchaseOrderConfirmReceiptURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// expenditure routes.
func (r ExpenditureRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expenditure.purchase.list":      r.PurchaseListURL,
		"expenditure.purchase.dashboard": r.PurchaseDashboardURL,
		"expenditure.expense.list":       r.ExpenseListURL,
		"expenditure.expense.dashboard":  r.ExpenseDashboardURL,

		"expenditure.purchases_summary": r.PurchasesSummaryURL,
		"expenditure.expenses_summary":  r.ExpensesSummaryURL,

		"purchases.settings.templates":        r.SettingsTemplatesURL,
		"purchases.settings.template_upload":  r.SettingsTemplateUploadURL,
		"purchases.settings.template_delete":  r.SettingsTemplateDeleteURL,
		"purchases.settings.template_default": r.SettingsTemplateDefaultURL,

		"expenditure.expense.add":        r.AddURL,
		"expenditure.expense.edit":       r.EditURL,
		"expenditure.expense.delete":     r.DeleteURL,
		"expenditure.expense.set_status": r.SetStatusURL,
		"expenditure.expense.detail":     r.DetailURL,
		"expenditure.expense.table":      r.TableURL,
		"expenditure.expense.pay":        r.PayURL,

		"expenditure.expense_category.list":   r.ExpenseCategoryListURL,
		"expenditure.expense_category.add":    r.ExpenseCategoryAddURL,
		"expenditure.expense_category.edit":   r.ExpenseCategoryEditURL,
		"expenditure.expense_category.delete": r.ExpenseCategoryDeleteURL,
		"expenditure.expense_category.table":  r.ExpenseCategoryTableURL,

		"expenditure.purchase_order.list":                  r.PurchaseOrderListURL,
		"expenditure.purchase_order.detail":                r.PurchaseOrderDetailURL,
		"expenditure.purchase_order.add":                   r.PurchaseOrderAddURL,
		"expenditure.purchase_order.edit":                  r.PurchaseOrderEditURL,
		"expenditure.purchase_order.delete":                r.PurchaseOrderDeleteURL,
		"expenditure.purchase_order.set_status":            r.PurchaseOrderSetStatusURL,
		"expenditure.purchase_order.table":                 r.PurchaseOrderTableURL,
		"expenditure.purchase_order.tab_action":            r.PurchaseOrderTabActionURL,
		"expenditure.purchase_order.line_item.table":       r.PurchaseOrderLineItemTableURL,
		"expenditure.purchase_order.line_item.add":         r.PurchaseOrderLineItemAddURL,
		"expenditure.purchase_order.line_item.edit":        r.PurchaseOrderLineItemEditURL,
		"expenditure.purchase_order.line_item.remove":      r.PurchaseOrderLineItemRemoveURL,
		"expenditure.purchase_order.confirm_receipt":       r.PurchaseOrderConfirmReceiptURL,
	}
}

// PlanRoutes holds all route paths for plan views and actions.
type PlanRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
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

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    string `json:"price_plan_add_url"`
	PricePlanEditURL   string `json:"price_plan_edit_url"`
	PricePlanDeleteURL string `json:"price_plan_delete_url"`

	// ProductPlan CRUD routes (within plan context)
	ProductPlanAddURL    string `json:"product_plan_add_url"`
	ProductPlanEditURL   string `json:"product_plan_edit_url"`
	ProductPlanDeleteURL string `json:"product_plan_delete_url"`
}

// DefaultPlanRoutes returns a PlanRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPlanRoutes() PlanRoutes {
	return PlanRoutes{
		ActiveNav:    "service",
		ActiveSubNav: "plans",

		ListURL:          PlanListURL,
		TableURL:         PlanTableURL,
		DetailURL:        PlanDetailURL,
		AddURL:           PlanAddURL,
		EditURL:          PlanEditURL,
		DeleteURL:        PlanDeleteURL,
		BulkDeleteURL:    PlanBulkDeleteURL,
		SetStatusURL:     PlanSetStatusURL,
		BulkSetStatusURL: PlanBulkSetStatusURL,
		TabActionURL:     PlanTabActionURL,

		AttachmentUploadURL: PlanAttachmentUploadURL,
		AttachmentDeleteURL: PlanAttachmentDeleteURL,

		PricePlanAddURL:    PricePlanAddURL,
		PricePlanEditURL:   PricePlanEditURL,
		PricePlanDeleteURL: PricePlanDeleteURL,

		ProductPlanAddURL:    PlanProductPlanAddURL,
		ProductPlanEditURL:   PlanProductPlanEditURL,
		ProductPlanDeleteURL: PlanProductPlanDeleteURL,
	}
}

// DefaultPlanBundleRoutes returns a PlanRoutes with every URL namespace-shifted
// from the services namespace onto the inventory accordion namespace. Used as
// the route base for the Plan inventory-mount registration in block.go; a lyngua
// `plan_bundle` override can layer additional tweaks on top.
//
// Bundle-mount variant of Plan routes — shifts every page + action URL from the
// services namespace (/app/plans/*) onto the inventory accordion namespace
// (/app/inventory/bundles/*). Used as the route base for the Plan
// inventory-mount registration in block.go; a lyngua `plan_bundle` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/plans/"   → "/app/inventory/bundles/"
//   - "/action/plan/" → "/action/inventory-bundle/"
func DefaultPlanBundleRoutes() PlanRoutes {
	r := DefaultPlanRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "bundles-active"
	shift := func(s string) string {
		s = strings.Replace(s, "/app/plans/", "/app/inventory/bundles/", 1)
		s = strings.Replace(s, "/action/plan/", "/action/inventory-bundle/", 1)
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
	r.PricePlanAddURL = shift(r.PricePlanAddURL)
	r.PricePlanEditURL = shift(r.PricePlanEditURL)
	r.PricePlanDeleteURL = shift(r.PricePlanDeleteURL)
	r.ProductPlanAddURL = shift(r.ProductPlanAddURL)
	r.ProductPlanEditURL = shift(r.ProductPlanEditURL)
	r.ProductPlanDeleteURL = shift(r.ProductPlanDeleteURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// plan routes.
func (r PlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"plan.list":             r.ListURL,
		"plan.table":            r.TableURL,
		"plan.detail":           r.DetailURL,
		"plan.add":              r.AddURL,
		"plan.edit":             r.EditURL,
		"plan.delete":           r.DeleteURL,
		"plan.bulk_delete":      r.BulkDeleteURL,
		"plan.set_status":       r.SetStatusURL,
		"plan.bulk_set_status":  r.BulkSetStatusURL,
		"plan.tab_action":       r.TabActionURL,

		"plan.attachment.upload": r.AttachmentUploadURL,
		"plan.attachment.delete": r.AttachmentDeleteURL,

		"plan.pricelist.add":    r.PricePlanAddURL,
		"plan.pricelist.edit":   r.PricePlanEditURL,
		"plan.pricelist.delete": r.PricePlanDeleteURL,

		"plan.product_plan.add":    r.ProductPlanAddURL,
		"plan.product_plan.edit":   r.ProductPlanEditURL,
		"plan.product_plan.delete": r.ProductPlanDeleteURL,
	}
}

// PricePlanRoutes holds all route paths for price plan (rate card) views and actions.
type PricePlanRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
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

// DefaultPricePlanRoutes returns a PricePlanRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPricePlanRoutes() PricePlanRoutes {
	return PricePlanRoutes{
		ActiveNav:             "service",
		ActiveSubNav:          "rate-cards",
		DashboardURL:          PricePlanDashboardURL,
		ListURL:               PricePlanListURL,
		TableURL:              PricePlanTableURL,
		DetailURL:             PricePlanDetailURL,
		AddURL:                PricePlanStandaloneAddURL,
		EditURL:               PricePlanStandaloneEditURL,
		DeleteURL:             PricePlanStandaloneDeleteURL,
		BulkDeleteURL:         PricePlanBulkDeleteURL,
		SetStatusURL:          PricePlanSetStatusURL,
		BulkSetStatusURL:      PricePlanBulkSetStatusURL,
		TabActionURL:          PricePlanTabActionURL,
		AttachmentUploadURL:   PricePlanAttachmentUploadURL,
		AttachmentDeleteURL:   PricePlanAttachmentDeleteURL,
		ProductPriceAddURL:    PricePlanProductPriceAddURL,
		ProductPriceEditURL:   PricePlanProductPriceEditURL,
		ProductPriceDeleteURL: PricePlanProductPriceDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price plan routes.
func (r PricePlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_plan.dashboard":         r.DashboardURL,
		"price_plan.list":              r.ListURL,
		"price_plan.table":             r.TableURL,
		"price_plan.detail":            r.DetailURL,
		"price_plan.add":               r.AddURL,
		"price_plan.edit":              r.EditURL,
		"price_plan.delete":            r.DeleteURL,
		"price_plan.bulk_delete":       r.BulkDeleteURL,
		"price_plan.set_status":        r.SetStatusURL,
		"price_plan.bulk_set_status":   r.BulkSetStatusURL,
		"price_plan.tab_action":        r.TabActionURL,
		"price_plan.attachment.upload":       r.AttachmentUploadURL,
		"price_plan.attachment.delete":       r.AttachmentDeleteURL,
		"price_plan.product_price.add":       r.ProductPriceAddURL,
		"price_plan.product_price.edit":      r.ProductPriceEditURL,
		"price_plan.product_price.delete":    r.ProductPriceDeleteURL,
	}
}

// PriceScheduleRoutes holds all route paths for price schedule views and actions.
type PriceScheduleRoutes struct {
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
	SetStatusURL               string `json:"set_status_url"`
	BulkSetStatusURL           string `json:"bulk_set_status_url"`
	TabActionURL               string `json:"tab_action_url"`
	PlanAddURL                 string `json:"plan_add_url"`
	PlanDetailURL              string `json:"plan_detail_url"`
	PlanTabActionURL           string `json:"plan_tab_action_url"`
	PlanEditURL                string `json:"plan_edit_url"`
	PlanDeleteURL              string `json:"plan_delete_url"`
	PlanProductPriceAddURL     string `json:"plan_product_price_add_url"`
	PlanProductPriceEditURL    string `json:"plan_product_price_edit_url"`
	PlanProductPriceDeleteURL  string `json:"plan_product_price_delete_url"`
}

// DefaultPriceScheduleRoutes returns a PriceScheduleRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPriceScheduleRoutes() PriceScheduleRoutes {
	return PriceScheduleRoutes{
		ActiveNav:        "service",
		ActiveSubNav:     "price-schedules",
		DashboardURL:     PriceScheduleDashboardURL,
		ListURL:          PriceScheduleListURL,
		TableURL:         PriceScheduleTableURL,
		DetailURL:        PriceScheduleDetailURL,
		AddURL:           PriceScheduleAddURL,
		EditURL:          PriceScheduleEditURL,
		DeleteURL:        PriceScheduleDeleteURL,
		BulkDeleteURL:    PriceScheduleBulkDeleteURL,
		SetStatusURL:              PriceScheduleSetStatusURL,
		BulkSetStatusURL:          PriceScheduleBulkSetStatusURL,
		TabActionURL:              PriceScheduleTabActionURL,
		PlanAddURL:                PriceSchedulePlanAddURL,
		PlanDetailURL:             PriceSchedulePlanDetailURL,
		PlanTabActionURL:          PriceSchedulePlanTabActionURL,
		PlanEditURL:               PriceSchedulePlanEditURL,
		PlanDeleteURL:             PriceSchedulePlanDeleteURL,
		PlanProductPriceAddURL:    PriceSchedulePlanProductPriceAddURL,
		PlanProductPriceEditURL:   PriceSchedulePlanProductPriceEditURL,
		PlanProductPriceDeleteURL: PriceSchedulePlanProductPriceDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price schedule routes.
func (r PriceScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_schedule.dashboard":       r.DashboardURL,
		"price_schedule.list":            r.ListURL,
		"price_schedule.table":           r.TableURL,
		"price_schedule.detail":          r.DetailURL,
		"price_schedule.add":             r.AddURL,
		"price_schedule.edit":            r.EditURL,
		"price_schedule.delete":          r.DeleteURL,
		"price_schedule.bulk_delete":     r.BulkDeleteURL,
		"price_schedule.set_status":      r.SetStatusURL,
		"price_schedule.bulk_set_status":                r.BulkSetStatusURL,
		"price_schedule.tab_action":                     r.TabActionURL,
		"price_schedule.plan.add":                       r.PlanAddURL,
		"price_schedule.plan.detail":                    r.PlanDetailURL,
		"price_schedule.plan.tab_action":                r.PlanTabActionURL,
		"price_schedule.plan.edit":                      r.PlanEditURL,
		"price_schedule.plan.delete":                    r.PlanDeleteURL,
		"price_schedule.plan.product_price.add":         r.PlanProductPriceAddURL,
		"price_schedule.plan.product_price.edit":        r.PlanProductPriceEditURL,
		"price_schedule.plan.product_price.delete":      r.PlanProductPriceDeleteURL,
	}
}

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
// route constants defined in routes.go.
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

// SubscriptionRoutes holds all route paths for subscription views and actions.
type SubscriptionRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL         string `json:"list_url"`
	DetailURL       string `json:"detail_url"`
	AddURL          string `json:"add_url"`
	EditURL         string `json:"edit_url"`
	DeleteURL       string `json:"delete_url"`
	TabActionURL    string `json:"tab_action_url"`
	SearchPlanURL   string `json:"search_plan_url"`
	SearchClientURL string `json:"search_client_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultSubscriptionRoutes returns a SubscriptionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultSubscriptionRoutes() SubscriptionRoutes {
	return SubscriptionRoutes{
		ActiveNav:    "client",
		ActiveSubNav: "subscriptions",

		ListURL:         SubscriptionListURL,
		DetailURL:       SubscriptionDetailURL,
		AddURL:          SubscriptionAddURL,
		EditURL:         SubscriptionEditURL,
		DeleteURL:       SubscriptionDeleteURL,
		TabActionURL:    SubscriptionTabActionURL,
		SearchPlanURL:   SubscriptionSearchPlanURL,
		SearchClientURL: SubscriptionSearchClientURL,

		AttachmentUploadURL: SubscriptionAttachmentUploadURL,
		AttachmentDeleteURL: SubscriptionAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription routes.
func (r SubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"subscription.list":          r.ListURL,
		"subscription.detail":        r.DetailURL,
		"subscription.add":           r.AddURL,
		"subscription.edit":          r.EditURL,
		"subscription.delete":        r.DeleteURL,
		"subscription.tab_action":    r.TabActionURL,
		"subscription.search_plan":   r.SearchPlanURL,
		"subscription.search_client": r.SearchClientURL,

		"subscription.attachment.upload": r.AttachmentUploadURL,
		"subscription.attachment.delete": r.AttachmentDeleteURL,
	}
}

// CollectionRoutes holds all route paths for collection (money IN) views
// and actions.
type CollectionRoutes struct {
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	DashboardURL     string `json:"dashboard_url"`
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

// DefaultCollectionRoutes returns a CollectionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultCollectionRoutes() CollectionRoutes {
	return CollectionRoutes{
		ListURL:          CollectionListURL,
		DetailURL:        CollectionDetailURL,
		DashboardURL:     CollectionDashboardURL,
		AddURL:           CollectionAddURL,
		EditURL:          CollectionEditURL,
		DeleteURL:        CollectionDeleteURL,
		BulkDeleteURL:    CollectionBulkDeleteURL,
		SetStatusURL:     CollectionSetStatusURL,
		BulkSetStatusURL: CollectionBulkSetStatusURL,
		TabActionURL:     CollectionTabActionURL,

		AttachmentUploadURL: CollectionAttachmentUploadURL,
		AttachmentDeleteURL: CollectionAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// collection routes.
func (r CollectionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"collection.list":            r.ListURL,
		"collection.detail":          r.DetailURL,
		"collection.dashboard":       r.DashboardURL,
		"collection.add":             r.AddURL,
		"collection.edit":            r.EditURL,
		"collection.delete":          r.DeleteURL,
		"collection.bulk_delete":     r.BulkDeleteURL,
		"collection.set_status":      r.SetStatusURL,
		"collection.bulk_set_status": r.BulkSetStatusURL,
		"collection.tab_action":      r.TabActionURL,

		"collection.attachment.upload": r.AttachmentUploadURL,
		"collection.attachment.delete": r.AttachmentDeleteURL,
	}
}

// DisbursementRoutes holds all route paths for disbursement (money OUT) views
// and actions.
type DisbursementRoutes struct {
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	DashboardURL     string `json:"dashboard_url"`
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

// DefaultDisbursementRoutes returns a DisbursementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultDisbursementRoutes() DisbursementRoutes {
	return DisbursementRoutes{
		ListURL:          DisbursementListURL,
		DetailURL:        DisbursementDetailURL,
		DashboardURL:     DisbursementDashboardURL,
		AddURL:           DisbursementAddURL,
		EditURL:          DisbursementEditURL,
		DeleteURL:        DisbursementDeleteURL,
		BulkDeleteURL:    DisbursementBulkDeleteURL,
		SetStatusURL:     DisbursementSetStatusURL,
		BulkSetStatusURL: DisbursementBulkSetStatusURL,
		TabActionURL:     DisbursementTabActionURL,

		AttachmentUploadURL: DisbursementAttachmentUploadURL,
		AttachmentDeleteURL: DisbursementAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// disbursement routes.
func (r DisbursementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"disbursement.list":            r.ListURL,
		"disbursement.detail":          r.DetailURL,
		"disbursement.dashboard":       r.DashboardURL,
		"disbursement.add":             r.AddURL,
		"disbursement.edit":            r.EditURL,
		"disbursement.delete":          r.DeleteURL,
		"disbursement.bulk_delete":     r.BulkDeleteURL,
		"disbursement.set_status":      r.SetStatusURL,
		"disbursement.bulk_set_status": r.BulkSetStatusURL,
		"disbursement.tab_action":      r.TabActionURL,

		"disbursement.attachment.upload": r.AttachmentUploadURL,
		"disbursement.attachment.delete": r.AttachmentDeleteURL,
	}
}

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
// package-level route constants defined in routes.go.
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
