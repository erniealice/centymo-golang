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
// route constants defined in routes.go.
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
	// Inventory mount has its own dashboard module — clear the service one.
	r.DashboardURL = ""
	// shift matches both pre-P4 (`/app/products/*`) and post-P4
	// (`/products/*`) constant shapes. P4 (2026-05-22) stripped the
	// `/app/` prefix from the URL constants in routes.go but did not
	// update this shift fn, which produced a duplicate ServeMux
	// registration for URLs the shift no longer matched (e.g.
	// OptionDetailURL `/products/detail/{id}/option/{oid}`). Adding the
	// post-P4 patterns alongside the pre-P4 ones keeps the function
	// idempotent across the parallel-mux window and survives the P12
	// hard cutover.
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
// the "Supplies" surface — the third Product module mount, scoped to
// product_kind = 'consumable'. Kept distinct from the inventory mount so the
// sidebar can surface "Products" (resold goods) and "Supplies" (used-in-
// service-delivery consumables) as separate entries without a ServeMux
// duplicate-registration panic.
//
// Shift rules:
//   - "/app/products/*"  → "/app/inventory/supplies/*"
//   - "/action/product/*" → "/action/inventory-supplies/*"
//
// Lyngua `product_supplies` route blocks can still override individual URLs.
func DefaultProductSuppliesRoutes() ProductRoutes {
	r := DefaultProductRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "supplies"
	// Supplies mount has no dashboard.
	r.DashboardURL = ""
	// shift matches both pre-P4 (`/app/products/*`) and post-P4
	// (`/products/*`) constant shapes — see DefaultProductInventoryRoutes
	// shift comment for the P4 regression context.
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
	// shift matches both pre-P4 (`/app/product-lines/*`) and post-P4
	// (`/product-lines/*`) constant shapes — see DefaultProductInventoryRoutes
	// shift comment for the P4 regression context.
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
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	DetailURL           string `json:"detail_url"`
	TableURL            string `json:"table_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

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
	PurchaseOrderListURL             string `json:"purchase_order_list_url"`
	PurchaseOrderDetailURL           string `json:"purchase_order_detail_url"`
	PurchaseOrderAddURL              string `json:"purchase_order_add_url"`
	PurchaseOrderEditURL             string `json:"purchase_order_edit_url"`
	PurchaseOrderDeleteURL           string `json:"purchase_order_delete_url"`
	PurchaseOrderSetStatusURL        string `json:"purchase_order_set_status_url"`
	PurchaseOrderTableURL            string `json:"purchase_order_table_url"`
	PurchaseOrderTabActionURL        string `json:"purchase_order_tab_action_url"`
	PurchaseOrderAttachmentUploadURL string `json:"purchase_order_attachment_upload_url"`
	PurchaseOrderAttachmentDeleteURL string `json:"purchase_order_attachment_delete_url"`

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

		AddURL:              ExpenditureExpenseAddURL,
		EditURL:             ExpenditureExpenseEditURL,
		DeleteURL:           ExpenditureExpenseDeleteURL,
		SetStatusURL:        ExpenditureExpenseSetStatusURL,
		DetailURL:           ExpenditureExpenseDetailURL,
		TableURL:            ExpenditureExpenseTableURL,
		TabActionURL:        ExpenditureExpenseTabActionURL,
		AttachmentUploadURL: ExpenditureAttachmentUploadURL,
		AttachmentDeleteURL: ExpenditureAttachmentDeleteURL,

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

		PurchaseOrderListURL:             PurchaseOrderListURL,
		PurchaseOrderDetailURL:           PurchaseOrderDetailURL,
		PurchaseOrderAddURL:              PurchaseOrderAddURL,
		PurchaseOrderEditURL:             PurchaseOrderEditURL,
		PurchaseOrderDeleteURL:           PurchaseOrderDeleteURL,
		PurchaseOrderSetStatusURL:        PurchaseOrderSetStatusURL,
		PurchaseOrderTableURL:            PurchaseOrderTableURL,
		PurchaseOrderTabActionURL:        PurchaseOrderTabActionURL,
		PurchaseOrderAttachmentUploadURL: PurchaseOrderAttachmentUploadURL,
		PurchaseOrderAttachmentDeleteURL: PurchaseOrderAttachmentDeleteURL,

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

		"expenditure.expense.add":               r.AddURL,
		"expenditure.expense.edit":              r.EditURL,
		"expenditure.expense.delete":            r.DeleteURL,
		"expenditure.expense.set_status":        r.SetStatusURL,
		"expenditure.expense.detail":            r.DetailURL,
		"expenditure.expense.table":             r.TableURL,
		"expenditure.expense.pay":               r.PayURL,
		"expenditure.expense.attachment.upload": r.AttachmentUploadURL,
		"expenditure.expense.attachment.delete": r.AttachmentDeleteURL,

		"expenditure.expense_category.list":   r.ExpenseCategoryListURL,
		"expenditure.expense_category.add":    r.ExpenseCategoryAddURL,
		"expenditure.expense_category.edit":   r.ExpenseCategoryEditURL,
		"expenditure.expense_category.delete": r.ExpenseCategoryDeleteURL,
		"expenditure.expense_category.table":  r.ExpenseCategoryTableURL,

		"expenditure.purchase_order.list":              r.PurchaseOrderListURL,
		"expenditure.purchase_order.detail":            r.PurchaseOrderDetailURL,
		"expenditure.purchase_order.add":               r.PurchaseOrderAddURL,
		"expenditure.purchase_order.edit":              r.PurchaseOrderEditURL,
		"expenditure.purchase_order.delete":            r.PurchaseOrderDeleteURL,
		"expenditure.purchase_order.set_status":        r.PurchaseOrderSetStatusURL,
		"expenditure.purchase_order.table":             r.PurchaseOrderTableURL,
		"expenditure.purchase_order.tab_action":        r.PurchaseOrderTabActionURL,
		"expenditure.purchase_order.attachment.upload": r.PurchaseOrderAttachmentUploadURL,
		"expenditure.purchase_order.attachment.delete": r.PurchaseOrderAttachmentDeleteURL,
		"expenditure.purchase_order.line_item.table":   r.PurchaseOrderLineItemTableURL,
		"expenditure.purchase_order.line_item.add":     r.PurchaseOrderLineItemAddURL,
		"expenditure.purchase_order.line_item.edit":    r.PurchaseOrderLineItemEditURL,
		"expenditure.purchase_order.line_item.remove":  r.PurchaseOrderLineItemRemoveURL,
		"expenditure.purchase_order.confirm_receipt":   r.PurchaseOrderConfirmReceiptURL,
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

	// Plan-scoped PricePlan detail (mirrors PriceSchedulePlanRoutes.PlanDetailURL
	// but anchored under /app/plans/detail/{id}/price/{ppid}). Lets the package
	// detail's package-prices tab keep ActiveNav on Services > Packages instead
	// of jumping to the rate-cards namespace.
	PricePlanDetailURL             string `json:"price_plan_detail_url"`
	PricePlanTabActionURL          string `json:"price_plan_tab_action_url"`
	PricePlanProductPriceAddURL    string `json:"price_plan_product_price_add_url"`
	PricePlanProductPriceEditURL   string `json:"price_plan_product_price_edit_url"`
	PricePlanProductPriceDeleteURL string `json:"price_plan_product_price_delete_url"`

	// ProductPlan CRUD routes (within plan context)
	ProductPlanAddURL    string `json:"product_plan_add_url"`
	ProductPlanEditURL   string `json:"product_plan_edit_url"`
	ProductPlanDeleteURL string `json:"product_plan_delete_url"`
	ProductPlanPickerURL string `json:"product_plan_picker_url"`
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

		PricePlanDetailURL:             PlanPricePlanDetailURL,
		PricePlanTabActionURL:          PlanPricePlanTabActionURL,
		PricePlanProductPriceAddURL:    PlanPricePlanProductPriceAddURL,
		PricePlanProductPriceEditURL:   PlanPricePlanProductPriceEditURL,
		PricePlanProductPriceDeleteURL: PlanPricePlanProductPriceDeleteURL,

		ProductPlanAddURL:    PlanProductPlanAddURL,
		ProductPlanEditURL:   PlanProductPlanEditURL,
		ProductPlanDeleteURL: PlanProductPlanDeleteURL,
		ProductPlanPickerURL: PlanProductPlanPickerURL,
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
	// shift matches both pre-P4 (`/app/plans/*`) and post-P4 (`/plans/*`)
	// constant shapes — see DefaultProductInventoryRoutes shift comment
	// for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/plans/", "/app/inventory/bundles/", 1)
		s = strings.Replace(s, "/plans/", "/inventory/bundles/", 1)
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
	r.PricePlanDetailURL = shift(r.PricePlanDetailURL)
	r.PricePlanTabActionURL = shift(r.PricePlanTabActionURL)
	r.PricePlanProductPriceAddURL = shift(r.PricePlanProductPriceAddURL)
	r.PricePlanProductPriceEditURL = shift(r.PricePlanProductPriceEditURL)
	r.PricePlanProductPriceDeleteURL = shift(r.PricePlanProductPriceDeleteURL)
	r.ProductPlanAddURL = shift(r.ProductPlanAddURL)
	r.ProductPlanEditURL = shift(r.ProductPlanEditURL)
	r.ProductPlanDeleteURL = shift(r.ProductPlanDeleteURL)
	r.ProductPlanPickerURL = shift(r.ProductPlanPickerURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// plan routes.
func (r PlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"plan.list":            r.ListURL,
		"plan.table":           r.TableURL,
		"plan.detail":          r.DetailURL,
		"plan.add":             r.AddURL,
		"plan.edit":            r.EditURL,
		"plan.delete":          r.DeleteURL,
		"plan.bulk_delete":     r.BulkDeleteURL,
		"plan.set_status":      r.SetStatusURL,
		"plan.bulk_set_status": r.BulkSetStatusURL,
		"plan.tab_action":      r.TabActionURL,

		"plan.attachment.upload": r.AttachmentUploadURL,
		"plan.attachment.delete": r.AttachmentDeleteURL,

		"plan.pricelist.add":    r.PricePlanAddURL,
		"plan.pricelist.edit":   r.PricePlanEditURL,
		"plan.pricelist.delete": r.PricePlanDeleteURL,

		"plan.price.detail":               r.PricePlanDetailURL,
		"plan.price.tab_action":           r.PricePlanTabActionURL,
		"plan.price.product_price.add":    r.PricePlanProductPriceAddURL,
		"plan.price.product_price.edit":   r.PricePlanProductPriceEditURL,
		"plan.price.product_price.delete": r.PricePlanProductPriceDeleteURL,

		"plan.product_plan.add":    r.ProductPlanAddURL,
		"plan.product_plan.edit":   r.ProductPlanEditURL,
		"plan.product_plan.delete": r.ProductPlanDeleteURL,
		"plan.product_plan.picker": r.ProductPlanPickerURL,
	}
}

// PricePlanRoutes holds all route paths for price plan (rate card) views and actions.
type PricePlanRoutes struct {
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

// PriceScheduleRoutes holds all route paths for price schedule views and actions.
type PriceScheduleRoutes struct {
	ActiveNav                 string `json:"active_nav"`
	ActiveSubNav              string `json:"active_sub_nav"`
	DashboardURL              string `json:"dashboard_url"`
	ListURL                   string `json:"list_url"`
	TableURL                  string `json:"table_url"`
	DetailURL                 string `json:"detail_url"`
	AddURL                    string `json:"add_url"`
	EditURL                   string `json:"edit_url"`
	DeleteURL                 string `json:"delete_url"`
	BulkDeleteURL             string `json:"bulk_delete_url"`
	SetStatusURL              string `json:"set_status_url"`
	BulkSetStatusURL          string `json:"bulk_set_status_url"`
	TabActionURL              string `json:"tab_action_url"`
	AttachmentUploadURL       string `json:"attachment_upload_url"`
	AttachmentDeleteURL       string `json:"attachment_delete_url"`
	PlanAddURL                string `json:"plan_add_url"`
	PlanDetailURL             string `json:"plan_detail_url"`
	PlanTabActionURL          string `json:"plan_tab_action_url"`
	PlanEditURL               string `json:"plan_edit_url"`
	PlanDeleteURL             string `json:"plan_delete_url"`
	PlanProductPriceAddURL    string `json:"plan_product_price_add_url"`
	PlanProductPriceEditURL   string `json:"plan_product_price_edit_url"`
	PlanProductPriceDeleteURL string `json:"plan_product_price_delete_url"`
	PlanAttachmentUploadURL   string `json:"plan_attachment_upload_url"`
	PlanAttachmentDeleteURL   string `json:"plan_attachment_delete_url"`
	// 2026-05-04 — Subscription detail nested under the
	// schedule-scoped price_plan path. Activates the rate-card → plan →
	// subscription breadcrumb in the subscription detail view. Empty string
	// disables the nested route. See
	// docs/plan/20260504-price-plan-engagements-tab/.
	PlanSubscriptionDetailURL string `json:"plan_subscription_detail_url"`
}

// DefaultPriceScheduleRoutes returns a PriceScheduleRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPriceScheduleRoutes() PriceScheduleRoutes {
	return PriceScheduleRoutes{
		ActiveNav:                 "service",
		ActiveSubNav:              "price-schedules",
		DashboardURL:              PriceScheduleDashboardURL,
		ListURL:                   PriceScheduleListURL,
		TableURL:                  PriceScheduleTableURL,
		DetailURL:                 PriceScheduleDetailURL,
		AddURL:                    PriceScheduleAddURL,
		EditURL:                   PriceScheduleEditURL,
		DeleteURL:                 PriceScheduleDeleteURL,
		BulkDeleteURL:             PriceScheduleBulkDeleteURL,
		SetStatusURL:              PriceScheduleSetStatusURL,
		BulkSetStatusURL:          PriceScheduleBulkSetStatusURL,
		TabActionURL:              PriceScheduleTabActionURL,
		AttachmentUploadURL:       PriceScheduleAttachmentUploadURL,
		AttachmentDeleteURL:       PriceScheduleAttachmentDeleteURL,
		PlanAddURL:                PriceSchedulePlanAddURL,
		PlanDetailURL:             PriceSchedulePlanDetailURL,
		PlanTabActionURL:          PriceSchedulePlanTabActionURL,
		PlanEditURL:               PriceSchedulePlanEditURL,
		PlanDeleteURL:             PriceSchedulePlanDeleteURL,
		PlanProductPriceAddURL:    PriceSchedulePlanProductPriceAddURL,
		PlanProductPriceEditURL:   PriceSchedulePlanProductPriceEditURL,
		PlanProductPriceDeleteURL: PriceSchedulePlanProductPriceDeleteURL,
		PlanAttachmentUploadURL:   PriceSchedulePlanAttachmentUploadURL,
		PlanAttachmentDeleteURL:   PriceSchedulePlanAttachmentDeleteURL,
		PlanSubscriptionDetailURL: PriceSchedulePlanSubscriptionDetailURL,
	}
}

// DefaultPriceScheduleInventoryRoutes returns a PriceScheduleRoutes with every
// URL namespace-shifted from the services namespace onto the inventory accordion
// namespace. Used as the route base for the PriceSchedule inventory-mount
// registration in block.go; a lyngua `price_schedule_inventory` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/price-schedules/"   → "/app/inventory/price-schedules/"
//   - "/action/price-schedule/" → "/action/inventory-price-schedule/"
func DefaultPriceScheduleInventoryRoutes() PriceScheduleRoutes {
	r := DefaultPriceScheduleRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "inventory-price-schedules-active"
	// shift matches both pre-P4 (`/app/price-schedules/*`) and post-P4
	// (`/price-schedules/*`) constant shapes — see DefaultProductInventoryRoutes
	// shift comment for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/price-schedules/", "/app/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/price-schedules/", "/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/action/price-schedule/", "/action/inventory-price-schedule/", 1)
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
	r.PlanAddURL = shift(r.PlanAddURL)
	r.PlanDetailURL = shift(r.PlanDetailURL)
	r.PlanTabActionURL = shift(r.PlanTabActionURL)
	r.PlanEditURL = shift(r.PlanEditURL)
	r.PlanDeleteURL = shift(r.PlanDeleteURL)
	r.PlanProductPriceAddURL = shift(r.PlanProductPriceAddURL)
	r.PlanProductPriceEditURL = shift(r.PlanProductPriceEditURL)
	r.PlanProductPriceDeleteURL = shift(r.PlanProductPriceDeleteURL)
	r.PlanAttachmentUploadURL = shift(r.PlanAttachmentUploadURL)
	r.PlanAttachmentDeleteURL = shift(r.PlanAttachmentDeleteURL)
	r.PlanSubscriptionDetailURL = shift(r.PlanSubscriptionDetailURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price schedule routes.
func (r PriceScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_schedule.dashboard":                 r.DashboardURL,
		"price_schedule.list":                      r.ListURL,
		"price_schedule.table":                     r.TableURL,
		"price_schedule.detail":                    r.DetailURL,
		"price_schedule.add":                       r.AddURL,
		"price_schedule.edit":                      r.EditURL,
		"price_schedule.delete":                    r.DeleteURL,
		"price_schedule.bulk_delete":               r.BulkDeleteURL,
		"price_schedule.set_status":                r.SetStatusURL,
		"price_schedule.bulk_set_status":           r.BulkSetStatusURL,
		"price_schedule.tab_action":                r.TabActionURL,
		"price_schedule.attachment.upload":         r.AttachmentUploadURL,
		"price_schedule.attachment.delete":         r.AttachmentDeleteURL,
		"price_schedule.plan.add":                  r.PlanAddURL,
		"price_schedule.plan.detail":               r.PlanDetailURL,
		"price_schedule.plan.tab_action":           r.PlanTabActionURL,
		"price_schedule.plan.edit":                 r.PlanEditURL,
		"price_schedule.plan.delete":               r.PlanDeleteURL,
		"price_schedule.plan.product_price.add":    r.PlanProductPriceAddURL,
		"price_schedule.plan.product_price.edit":   r.PlanProductPriceEditURL,
		"price_schedule.plan.product_price.delete": r.PlanProductPriceDeleteURL,
		"price_schedule.plan.attachment.upload":    r.PlanAttachmentUploadURL,
		"price_schedule.plan.attachment.delete":    r.PlanAttachmentDeleteURL,
		"price_schedule.plan.subscription.detail":   r.PlanSubscriptionDetailURL,
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

	ListURL              string `json:"list_url"`
	TableURL             string `json:"table_url"`
	DetailURL            string `json:"detail_url"`
	UnderClientDetailURL string `json:"under_client_detail_url"`
	AddURL               string `json:"add_url"`
	EditURL              string `json:"edit_url"`
	DeleteURL            string `json:"delete_url"`
	BulkDeleteURL        string `json:"bulk_delete_url"`
	SetStatusURL         string `json:"set_status_url"`
	BulkSetStatusURL     string `json:"bulk_set_status_url"`
	TabActionURL         string `json:"tab_action_url"`
	SearchPlanURL        string `json:"search_plan_url"`
	SearchClientURL      string `json:"search_client_url"`

	// RecognizeURL opens the "Recognize Revenue" drawer for a subscription.
	// GET = preview drawer (dry_run); POST = generate the Revenue.
	RecognizeURL string `json:"recognize_url"`

	// CustomizePackageURL is the POST endpoint that drives the "Customize
	// this package for {ClientName}" CTA on the subscription detail's
	// Package tab. Server clones the source Plan + PricePlan into a
	// client-scoped copy and HX-redirects to the new package URL.
	CustomizePackageURL string `json:"customize_package_url"`

	// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready + waive
	// endpoints for BillingEvent rows on the subscription Package tab.
	MilestoneMarkReadyURL string `json:"milestone_mark_ready_url"`
	MilestoneWaiveURL     string `json:"milestone_waive_url"`

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize button on a
	// BillingEvent row when it is linked to a MILESTONE advance Collection.
	// POSTs through espyna RecognizeMilestoneAdvanceCollection.
	MilestoneRecognizeURL string `json:"milestone_recognize_url"`

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
	// retroactive spawn drawer URL + HTMX-driven partial URL for the
	// Spawn Jobs section on the create form.
	SpawnJobsURL        string `json:"spawn_jobs_url"`
	SpawnJobsPartialURL string `json:"spawn_jobs_partial_url"`

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — manual cycle
	// spawn + backfill triggers. Both POST through espyna's
	// MaterializeInstanceJobsForSubscription consumer. Backfill GET renders
	// a preview drawer; POST commits the spawn.
	SpawnCycleJobsURL    string `json:"spawn_cycle_jobs_url"`
	BackfillCycleJobsURL string `json:"backfill_cycle_jobs_url"`

	// 2026-05-01 ad-hoc-subscription-billing — operator-driven Request Usage CTA.
	RequestUsageURL string `json:"request_usage_url"`

	// 2026-05-06 revenue-run — per-subscription Invoice Run drawer (Surface C,
	// CYCLE billing_kind only). Empty string when revenue-run module is not wired.
	RevenueRunURL string `json:"revenue_run_url"`

	// Attachment routes
	AttachmentUploadURL   string `json:"attachment_upload_url"`
	AttachmentDeleteURL   string `json:"attachment_delete_url"`
	AttachmentDownloadURL string `json:"attachment_download_url"`
}

// DefaultSubscriptionRoutes returns a SubscriptionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultSubscriptionRoutes() SubscriptionRoutes {
	return SubscriptionRoutes{
		ActiveNav:    "client",
		ActiveSubNav: "subscriptions",

		ListURL:              SubscriptionListURL,
		TableURL:             SubscriptionTableURL,
		DetailURL:            SubscriptionDetailURL,
		UnderClientDetailURL: SubscriptionUnderClientDetailURL,
		AddURL:               SubscriptionAddURL,
		EditURL:              SubscriptionEditURL,
		DeleteURL:            SubscriptionDeleteURL,
		BulkDeleteURL:        SubscriptionBulkDeleteURL,
		SetStatusURL:         SubscriptionSetStatusURL,
		BulkSetStatusURL:     SubscriptionBulkSetStatusURL,
		TabActionURL:         SubscriptionTabActionURL,
		SearchPlanURL:        SubscriptionSearchPlanURL,
		SearchClientURL:      SubscriptionSearchClientURL,
		RecognizeURL:         SubscriptionRecognizeURL,
		CustomizePackageURL:  SubscriptionCustomizePackageURL,

		// 2026-04-29 milestone-billing.
		MilestoneMarkReadyURL: MilestoneMarkReadyURL,
		MilestoneWaiveURL:     MilestoneWaiveURL,

		// 20260517-advance-cash-events Plan B Phase 7.
		MilestoneRecognizeURL: MilestoneRecognizeURL,

		// 2026-04-29 auto-spawn-jobs-from-subscription.
		SpawnJobsURL:        SubscriptionSpawnJobsURL,
		SpawnJobsPartialURL: SubscriptionSpawnJobsPartialURL,

		// 2026-04-30 cyclic-subscription-jobs.
		SpawnCycleJobsURL:    SubscriptionSpawnCycleJobsURL,
		BackfillCycleJobsURL: SubscriptionBackfillCycleJobsURL,

		// 2026-05-01 ad-hoc-subscription-billing.
		RequestUsageURL: SubscriptionRequestUsageURL,

		// 2026-05-06 revenue-run — per-subscription drawer.
		RevenueRunURL: SubscriptionRevenueRunURL,

		AttachmentUploadURL:   SubscriptionAttachmentUploadURL,
		AttachmentDeleteURL:   SubscriptionAttachmentDeleteURL,
		AttachmentDownloadURL: SubscriptionAttachmentDownloadURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription routes.
func (r SubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"subscription.list":                r.ListURL,
		"subscription.table":               r.TableURL,
		"subscription.detail":              r.DetailURL,
		"subscription.under_client_detail": r.UnderClientDetailURL,
		"subscription.add":                 r.AddURL,
		"subscription.edit":                r.EditURL,
		"subscription.delete":              r.DeleteURL,
		"subscription.bulk_delete":         r.BulkDeleteURL,
		"subscription.set_status":          r.SetStatusURL,
		"subscription.bulk_set_status":     r.BulkSetStatusURL,
		"subscription.tab_action":          r.TabActionURL,
		"subscription.search_plan":         r.SearchPlanURL,
		"subscription.search_client":       r.SearchClientURL,
		"subscription.recognize":           r.RecognizeURL,
		"subscription.customize_package":   r.CustomizePackageURL,

		// 2026-04-29 milestone-billing routes.
		"milestone.mark_ready": r.MilestoneMarkReadyURL,
		"milestone.waive":      r.MilestoneWaiveURL,

		// 20260517-advance-cash-events Plan B Phase 7.
		"milestone.recognize": r.MilestoneRecognizeURL,

		// 2026-04-29 auto-spawn-jobs-from-subscription routes.
		"subscription.spawn_jobs":         r.SpawnJobsURL,
		"subscription.spawn_jobs_partial": r.SpawnJobsPartialURL,

		// 2026-04-30 cyclic-subscription-jobs routes.
		"subscription.spawn_cycle_jobs":    r.SpawnCycleJobsURL,
		"subscription.backfill_cycle_jobs": r.BackfillCycleJobsURL,

		// 2026-05-01 ad-hoc-subscription-billing routes.
		"subscription.request_usage": r.RequestUsageURL,

		// 2026-05-06 revenue-run per-subscription drawer.
		"subscription.revenue_run": r.RevenueRunURL,

		"subscription.attachment.upload":   r.AttachmentUploadURL,
		"subscription.attachment.delete":   r.AttachmentDeleteURL,
		"subscription.attachment.download": r.AttachmentDownloadURL,
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

	// Advance Cash Events (Plan B Phase 3) — UNSCHEDULED workflow drawers +
	// the Advance Schedule tab partial. Empty defaults render the actions as
	// disabled / hidden.
	AdvanceScheduleTabURL string `json:"advance_schedule_tab_url"`
	SettleURL             string `json:"settle_url"`
	RefundURL             string `json:"refund_url"`
	CancelURL             string `json:"cancel_url"`
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

		// 20260517-advance-cash-events Plan B Phase 3.
		AdvanceScheduleTabURL: TreasuryCollectionAdvanceScheduleTabURL,
		SettleURL:             TreasuryCollectionSettleURL,
		RefundURL:             TreasuryCollectionRefundURL,
		CancelURL:             TreasuryCollectionCancelURL,
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

		// 20260517-advance-cash-events Plan B Phase 3.
		"collection.advance_schedule_tab": r.AdvanceScheduleTabURL,
		"collection.settle":               r.SettleURL,
		"collection.refund":               r.RefundURL,
		"collection.cancel":               r.CancelURL,
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

	// Advance Cash Events (Plan B Phase 3) — UNSCHEDULED workflow drawers +
	// the Advance Schedule tab partial. Mirrors CollectionRoutes.
	AdvanceScheduleTabURL string `json:"advance_schedule_tab_url"`
	SettleURL             string `json:"settle_url"`
	RefundURL             string `json:"refund_url"`
	CancelURL             string `json:"cancel_url"`
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

		// 20260517-advance-cash-events Plan B Phase 3.
		AdvanceScheduleTabURL: TreasuryDisbursementAdvanceScheduleTabURL,
		SettleURL:             TreasuryDisbursementSettleURL,
		RefundURL:             TreasuryDisbursementRefundURL,
		CancelURL:             TreasuryDisbursementCancelURL,
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

		// 20260517-advance-cash-events Plan B Phase 3.
		"disbursement.advance_schedule_tab": r.AdvanceScheduleTabURL,
		"disbursement.settle":               r.SettleURL,
		"disbursement.refund":               r.RefundURL,
		"disbursement.cancel":               r.CancelURL,
	}
}

// ---------------------------------------------------------------------------
// 20260517-advance-cash-events Plan B Phase 3 — TreasuryAdvancesRoutes
// ---------------------------------------------------------------------------

// TreasuryAdvancesRoutes holds all route paths for the cash-app "Advances"
// section: the workspace-level dashboard plus the filtered-list URLs for
// advance Collections / advance Disbursements (which point at the existing
// list pages with the `advance_kind` filter chip pre-applied).
//
// The Settle / Refund / Cancel drawer routes live on CollectionRoutes /
// DisbursementRoutes because they are anchored on the existing detail pages
// — there is no separate "advance" entity.
type TreasuryAdvancesRoutes struct {
	// ActiveNav is the sidebar navigation context; the Advances section sits
	// inside the Cash app so this remains "cash".
	ActiveNav string `json:"active_nav"`

	// DashboardURL is the workspace-level Advances Dashboard.
	DashboardURL string `json:"dashboard_url"`

	// AdvanceCollectionListURL / AdvanceDisbursementListURL are deep-links
	// into the existing Collection / Disbursement list pages with the
	// `advance_kind=any` chip pre-applied via the query string.
	AdvanceCollectionListURL   string `json:"advance_collection_list_url"`
	AdvanceDisbursementListURL string `json:"advance_disbursement_list_url"`

	// SupplierBillingEvent surfaces (buying-side MILESTONE anchor). Listed
	// here so the cash-app sidebar can deep-link to them, even though the
	// entity is buying-side. (Plan B Phase 7 wires the Recognize button.)
	SupplierBillingEventListURL      string `json:"supplier_billing_event_list_url"`
	SupplierBillingEventDetailURL    string `json:"supplier_billing_event_detail_url"`
	SupplierBillingEventRecognizeURL string `json:"supplier_billing_event_recognize_url"`
}

// DefaultTreasuryAdvancesRoutes returns a TreasuryAdvancesRoutes populated
// from the package-level route constants defined in routes.go.
func DefaultTreasuryAdvancesRoutes() TreasuryAdvancesRoutes {
	return TreasuryAdvancesRoutes{
		ActiveNav:                        "cash",
		DashboardURL:                     AdvancesDashboardURL,
		AdvanceCollectionListURL:         AdvanceCollectionListURL,
		AdvanceDisbursementListURL:       AdvanceDisbursementListURL,
		SupplierBillingEventListURL:      SupplierBillingEventListURL,
		SupplierBillingEventDetailURL:    SupplierBillingEventDetailURL,
		SupplierBillingEventRecognizeURL: SupplierBillingEventRecognizeURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// treasury-advances routes.
func (r TreasuryAdvancesRoutes) RouteMap() map[string]string {
	return map[string]string{
		"treasury_advances.dashboard":                  r.DashboardURL,
		"treasury_advances.advance_collection_list":    r.AdvanceCollectionListURL,
		"treasury_advances.advance_disbursement_list":  r.AdvanceDisbursementListURL,
		"supplier_billing_event.list":                  r.SupplierBillingEventListURL,
		"supplier_billing_event.detail":                r.SupplierBillingEventDetailURL,
		"supplier_billing_event.recognize":             r.SupplierBillingEventRecognizeURL,
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

// ---------------------------------------------------------------------------
// SupplierContractRoutes — P3a
// ---------------------------------------------------------------------------

// SupplierContractRoutes holds all route paths for supplier_contract views.
type SupplierContractRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	ApproveURL   string `json:"approve_url"`
	TerminateURL string `json:"terminate_url"`

	// Line item actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultSupplierContractRoutes returns a SupplierContractRoutes using the
// package-level route constants.
func DefaultSupplierContractRoutes() SupplierContractRoutes {
	return SupplierContractRoutes{
		ActiveNav:           "supplier-contracts",
		ActiveSubNav:        "active",
		ListURL:             SupplierContractListURL,
		DetailURL:           SupplierContractDetailURL,
		AddURL:              SupplierContractAddURL,
		EditURL:             SupplierContractEditURL,
		DeleteURL:           SupplierContractDeleteURL,
		SetStatusURL:        SupplierContractSetStatusURL,
		BulkSetStatusURL:    SupplierContractBulkSetStatusURL,
		TabActionURL:        SupplierContractTabActionURL,
		AttachmentUploadURL: SupplierContractAttachmentUploadURL,
		AttachmentDeleteURL: SupplierContractAttachmentDeleteURL,
		ApproveURL:          SupplierContractApproveURL,
		TerminateURL:        SupplierContractTerminateURL,
		LineAddURL:          SupplierContractLineAddURL,
		LineEditURL:         SupplierContractLineEditURL,
		LineDeleteURL:       SupplierContractLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierContractRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_contract.list":              r.ListURL,
		"supplier_contract.detail":            r.DetailURL,
		"supplier_contract.add":               r.AddURL,
		"supplier_contract.edit":              r.EditURL,
		"supplier_contract.delete":            r.DeleteURL,
		"supplier_contract.set_status":        r.SetStatusURL,
		"supplier_contract.attachment.upload": r.AttachmentUploadURL,
		"supplier_contract.attachment.delete": r.AttachmentDeleteURL,
		"supplier_contract.approve":           r.ApproveURL,
		"supplier_contract.terminate":         r.TerminateURL,
		"supplier_contract.line.add":          r.LineAddURL,
		"supplier_contract.line.edit":         r.LineEditURL,
		"supplier_contract.line.delete":       r.LineDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// ProcurementRequestRoutes — P3a
// ---------------------------------------------------------------------------

// ProcurementRequestRoutes holds all route paths for procurement_request views.
type ProcurementRequestRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow actions
	SubmitURL  string `json:"submit_url"`
	ApproveURL string `json:"approve_url"`
	RejectURL  string `json:"reject_url"`
	SpawnPOURL string `json:"spawn_po_url"`

	// Line item actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`

	// SPS Wave 3 — CRIT-3 retry placeholder. Wired but the action use case
	// itself is intentionally out-of-scope; handler currently logs + redirects.
	LineRetrySpawnURL string `json:"line_retry_spawn_url"`
}

// DefaultProcurementRequestRoutes returns a ProcurementRequestRoutes using the
// package-level route constants.
func DefaultProcurementRequestRoutes() ProcurementRequestRoutes {
	return ProcurementRequestRoutes{
		ActiveNav:           "procurement",
		ActiveSubNav:        "draft",
		ListURL:             ProcurementRequestListURL,
		DetailURL:           ProcurementRequestDetailURL,
		AddURL:              ProcurementRequestAddURL,
		EditURL:             ProcurementRequestEditURL,
		DeleteURL:           ProcurementRequestDeleteURL,
		SetStatusURL:        ProcurementRequestSetStatusURL,
		BulkSetStatusURL:    ProcurementRequestBulkSetStatusURL,
		TabActionURL:        ProcurementRequestTabActionURL,
		AttachmentUploadURL: ProcurementRequestAttachmentUploadURL,
		AttachmentDeleteURL: ProcurementRequestAttachmentDeleteURL,
		SubmitURL:           ProcurementRequestSubmitURL,
		ApproveURL:          ProcurementRequestApproveURL,
		RejectURL:           ProcurementRequestRejectURL,
		SpawnPOURL:          ProcurementRequestSpawnPOURL,
		LineAddURL:          ProcurementRequestLineAddURL,
		LineEditURL:         ProcurementRequestLineEditURL,
		LineDeleteURL:       ProcurementRequestLineDeleteURL,
		LineRetrySpawnURL:   ProcurementRequestLineRetrySpawnURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ProcurementRequestRoutes) RouteMap() map[string]string {
	return map[string]string{
		"procurement_request.list":              r.ListURL,
		"procurement_request.detail":            r.DetailURL,
		"procurement_request.add":               r.AddURL,
		"procurement_request.edit":              r.EditURL,
		"procurement_request.delete":            r.DeleteURL,
		"procurement_request.set_status":        r.SetStatusURL,
		"procurement_request.attachment.upload": r.AttachmentUploadURL,
		"procurement_request.attachment.delete": r.AttachmentDeleteURL,
		"procurement_request.submit":            r.SubmitURL,
		"procurement_request.approve":           r.ApproveURL,
		"procurement_request.reject":            r.RejectURL,
		"procurement_request.spawn_po":          r.SpawnPOURL,
		"procurement_request.line.add":          r.LineAddURL,
		"procurement_request.line.edit":         r.LineEditURL,
		"procurement_request.line.delete":       r.LineDeleteURL,
		"procurement_request.line.retry_spawn":  r.LineRetrySpawnURL,
	}
}

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app routes
// (composition surface; no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementRoutes holds the URL constants for the Procurement Operations app.
// These are defined in the centymo package so service-admin composition (P3c)
// can wire them into SidebarRoutes.Operations.Procurement.
type ProcurementRoutes struct {
	// Dashboard
	DashboardURL string `json:"dashboard_url"`

	// Contract operations (views over SupplierContract)
	RenewalCalendarURL string `json:"renewal_calendar_url"`
	VarianceURL        string `json:"variance_url"`
	UtilizationURL     string `json:"utilization_url"`

	// Recurrence drafts queue (lights up when P5 ships the recurrence engine)
	RecurrenceDraftsURL string `json:"recurrence_drafts_url"`
}

// DefaultProcurementRoutes returns a ProcurementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultProcurementRoutes() ProcurementRoutes {
	return ProcurementRoutes{
		DashboardURL:        ProcurementDashboardURL,
		RenewalCalendarURL:  ProcurementRenewalCalendarURL,
		VarianceURL:         ProcurementVarianceURL,
		UtilizationURL:      ProcurementUtilizationURL,
		RecurrenceDraftsURL: ProcurementRecurrenceDraftsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// procurement operations app routes.
func (r ProcurementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"procurement.dashboard":         r.DashboardURL,
		"procurement.renewals":          r.RenewalCalendarURL,
		"procurement.variance":          r.VarianceURL,
		"procurement.utilization":       r.UtilizationURL,
		"procurement.recurrence_drafts": r.RecurrenceDraftsURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierContractPriceScheduleRoutes — SPS P7
// ---------------------------------------------------------------------------

// SupplierContractPriceScheduleRoutes holds all route paths for
// supplier_contract_price_schedule + child line views.
type SupplierContractPriceScheduleRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	ActivateURL  string `json:"activate_url"`
	SupersedeURL string `json:"supersede_url"`

	// Schedule line actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultSupplierContractPriceScheduleRoutes returns a
// SupplierContractPriceScheduleRoutes using the package-level URL constants.
func DefaultSupplierContractPriceScheduleRoutes() SupplierContractPriceScheduleRoutes {
	return SupplierContractPriceScheduleRoutes{
		ActiveNav:           "supplier-contract-price-schedules",
		ActiveSubNav:        "active",
		ListURL:             SupplierContractPriceScheduleListURL,
		DetailURL:           SupplierContractPriceScheduleDetailURL,
		AddURL:              SupplierContractPriceScheduleAddURL,
		EditURL:             SupplierContractPriceScheduleEditURL,
		DeleteURL:           SupplierContractPriceScheduleDeleteURL,
		SetStatusURL:        SupplierContractPriceScheduleSetStatusURL,
		BulkSetStatusURL:    SupplierContractPriceScheduleBulkSetStatusURL,
		TabActionURL:        SupplierContractPriceScheduleTabActionURL,
		AttachmentUploadURL: SupplierContractPriceScheduleAttachmentUploadURL,
		AttachmentDeleteURL: SupplierContractPriceScheduleAttachmentDeleteURL,
		ActivateURL:         SupplierContractPriceScheduleActivateURL,
		SupersedeURL:        SupplierContractPriceScheduleSupersedeURL,
		LineAddURL:          SupplierContractPriceScheduleLineAddURL,
		LineEditURL:         SupplierContractPriceScheduleLineEditURL,
		LineDeleteURL:       SupplierContractPriceScheduleLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierContractPriceScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_contract_price_schedule.list":              r.ListURL,
		"supplier_contract_price_schedule.detail":            r.DetailURL,
		"supplier_contract_price_schedule.add":               r.AddURL,
		"supplier_contract_price_schedule.edit":              r.EditURL,
		"supplier_contract_price_schedule.delete":            r.DeleteURL,
		"supplier_contract_price_schedule.set_status":        r.SetStatusURL,
		"supplier_contract_price_schedule.attachment.upload": r.AttachmentUploadURL,
		"supplier_contract_price_schedule.attachment.delete": r.AttachmentDeleteURL,
		"supplier_contract_price_schedule.activate":          r.ActivateURL,
		"supplier_contract_price_schedule.supersede":         r.SupersedeURL,
		"supplier_contract_price_schedule.line.add":          r.LineAddURL,
		"supplier_contract_price_schedule.line.edit":         r.LineEditURL,
		"supplier_contract_price_schedule.line.delete":       r.LineDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// ExpenseRecognitionRoutes — SPS P10
// ---------------------------------------------------------------------------

// ExpenseRecognitionRoutes holds all route paths for expense_recognition views.
// Note: no Add/Edit URLs — recognitions are created BY use case, not by user.
type ExpenseRecognitionRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	DeleteURL           string `json:"delete_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	ReverseURL                  string `json:"reverse_url"`
	RecognizeFromExpenditureURL string `json:"recognize_from_expenditure_url"`
	RecognizeFromContractURL    string `json:"recognize_from_contract_url"`

	// Recognition line actions (child entity — inline CRUD)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultExpenseRecognitionRoutes returns an ExpenseRecognitionRoutes using the
// package-level URL constants.
func DefaultExpenseRecognitionRoutes() ExpenseRecognitionRoutes {
	return ExpenseRecognitionRoutes{
		ActiveNav:                   "expense-recognitions",
		ActiveSubNav:                "posted",
		ListURL:                     ExpenseRecognitionListURL,
		DetailURL:                   ExpenseRecognitionDetailURL,
		DeleteURL:                   ExpenseRecognitionDeleteURL,
		TabActionURL:                ExpenseRecognitionTabActionURL,
		AttachmentUploadURL:         ExpenseRecognitionAttachmentUploadURL,
		AttachmentDeleteURL:         ExpenseRecognitionAttachmentDeleteURL,
		ReverseURL:                  ExpenseRecognitionReverseURL,
		RecognizeFromExpenditureURL: ExpenseRecognitionRecognizeFromExpenditureURL,
		RecognizeFromContractURL:    ExpenseRecognitionRecognizeFromContractURL,
		LineAddURL:                  ExpenseRecognitionLineAddURL,
		LineEditURL:                 ExpenseRecognitionLineEditURL,
		LineDeleteURL:               ExpenseRecognitionLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ExpenseRecognitionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition.list":                       r.ListURL,
		"expense_recognition.detail":                     r.DetailURL,
		"expense_recognition.delete":                     r.DeleteURL,
		"expense_recognition.attachment.upload":          r.AttachmentUploadURL,
		"expense_recognition.attachment.delete":          r.AttachmentDeleteURL,
		"expense_recognition.reverse":                    r.ReverseURL,
		"expense_recognition.recognize_from_expenditure": r.RecognizeFromExpenditureURL,
		"expense_recognition.recognize_from_contract":    r.RecognizeFromContractURL,
		"expense_recognition.line.add":                   r.LineAddURL,
		"expense_recognition.line.edit":                  r.LineEditURL,
		"expense_recognition.line.delete":                r.LineDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// AccruedExpenseRoutes — SPS P10
// ---------------------------------------------------------------------------

// AccruedExpenseRoutes holds all route paths for accrued_expense views.
type AccruedExpenseRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachments
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	SettleURL             string `json:"settle_url"`
	ReverseURL            string `json:"reverse_url"`
	AccrueFromContractURL string `json:"accrue_from_contract_url"`

	// Settlement actions (child entity — inline CRUD)
	SettlementAddURL    string `json:"settlement_add_url"`
	SettlementEditURL   string `json:"settlement_edit_url"`
	SettlementDeleteURL string `json:"settlement_delete_url"`
}

// DefaultAccruedExpenseRoutes returns an AccruedExpenseRoutes using the
// package-level URL constants.
func DefaultAccruedExpenseRoutes() AccruedExpenseRoutes {
	return AccruedExpenseRoutes{
		ActiveNav:             "accrued-expenses",
		ActiveSubNav:          "outstanding",
		ListURL:               AccruedExpenseListURL,
		DetailURL:             AccruedExpenseDetailURL,
		AddURL:                AccruedExpenseAddURL,
		EditURL:               AccruedExpenseEditURL,
		DeleteURL:             AccruedExpenseDeleteURL,
		SetStatusURL:          AccruedExpenseSetStatusURL,
		BulkSetStatusURL:      AccruedExpenseBulkSetStatusURL,
		TabActionURL:          AccruedExpenseTabActionURL,
		AttachmentUploadURL:   AccruedExpenseAttachmentUploadURL,
		AttachmentDeleteURL:   AccruedExpenseAttachmentDeleteURL,
		SettleURL:             AccruedExpenseSettleURL,
		ReverseURL:            AccruedExpenseReverseURL,
		AccrueFromContractURL: AccruedExpenseAccrueFromContractURL,
		SettlementAddURL:      AccruedExpenseSettlementAddURL,
		SettlementEditURL:     AccruedExpenseSettlementEditURL,
		SettlementDeleteURL:   AccruedExpenseSettlementDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AccruedExpenseRoutes) RouteMap() map[string]string {
	return map[string]string{
		"accrued_expense.list":                 r.ListURL,
		"accrued_expense.detail":               r.DetailURL,
		"accrued_expense.add":                  r.AddURL,
		"accrued_expense.edit":                 r.EditURL,
		"accrued_expense.delete":               r.DeleteURL,
		"accrued_expense.set_status":           r.SetStatusURL,
		"accrued_expense.attachment.upload":    r.AttachmentUploadURL,
		"accrued_expense.attachment.delete":    r.AttachmentDeleteURL,
		"accrued_expense.settle":               r.SettleURL,
		"accrued_expense.reverse":              r.ReverseURL,
		"accrued_expense.accrue_from_contract": r.AccrueFromContractURL,
		"accrued_expense.settlement.add":       r.SettlementAddURL,
		"accrued_expense.settlement.edit":      r.SettlementEditURL,
		"accrued_expense.settlement.delete":    r.SettlementDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierSubscription routes (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// SupplierSubscriptionRoutes holds all route paths for supplier_subscription views.
type SupplierSubscriptionRoutes struct {
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

	// Search autocomplete endpoints for the add/edit drawer
	SearchCostPlanURL string `json:"search_cost_plan_url"`
	SearchSupplierURL string `json:"search_supplier_url"`

	// Recognition CTA — POST; opens the recognize-expense drawer on the detail page.
	RecognizeExpenseURL string `json:"recognize_expense_url"`

	// ExpenseRecognitionRunURL — GET; opens the per-SupplierSubscription Expense
	// Recognition Run drawer (Surface C). Resolved by resolveRecognitionsPrimaryAction
	// for CostPlan.billing_kind RECURRING / CONTRACT-with-cycle.
	// Plan A 20260517-expense-run Phase 4 / Surface C.
	ExpenseRecognitionRunURL string `json:"expense_recognition_run_url"`
}

// DefaultSupplierSubscriptionRoutes returns a SupplierSubscriptionRoutes using route constants.
func DefaultSupplierSubscriptionRoutes() SupplierSubscriptionRoutes {
	return SupplierSubscriptionRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-subscriptions",

		ListURL:                  SupplierSubscriptionListURL,
		TableURL:                 SupplierSubscriptionTableURL,
		DetailURL:                SupplierSubscriptionDetailURL,
		AddURL:                   SupplierSubscriptionAddURL,
		EditURL:                  SupplierSubscriptionEditURL,
		DeleteURL:                SupplierSubscriptionDeleteURL,
		BulkDeleteURL:            SupplierSubscriptionBulkDeleteURL,
		SetStatusURL:             SupplierSubscriptionSetStatusURL,
		BulkSetStatusURL:         SupplierSubscriptionBulkSetStatusURL,
		TabActionURL:             SupplierSubscriptionTabActionURL,
		SearchCostPlanURL:        SupplierSubscriptionSearchCostPlanURL,
		SearchSupplierURL:        SupplierSubscriptionSearchSupplierURL,
		RecognizeExpenseURL:      SupplierSubscriptionRecognizeExpenseURL,
		ExpenseRecognitionRunURL: ExpenseRecognitionRunPerSubscriptionDrawerURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierSubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_subscription.list":                    r.ListURL,
		"supplier_subscription.table":                   r.TableURL,
		"supplier_subscription.detail":                  r.DetailURL,
		"supplier_subscription.add":                     r.AddURL,
		"supplier_subscription.edit":                    r.EditURL,
		"supplier_subscription.delete":                  r.DeleteURL,
		"supplier_subscription.bulk_delete":             r.BulkDeleteURL,
		"supplier_subscription.set_status":              r.SetStatusURL,
		"supplier_subscription.bulk_set_status":         r.BulkSetStatusURL,
		"supplier_subscription.tab_action":              r.TabActionURL,
		"supplier_subscription.search_cost_plan":        r.SearchCostPlanURL,
		"supplier_subscription.search_supplier":         r.SearchSupplierURL,
		"supplier_subscription.recognize_expense":       r.RecognizeExpenseURL,
		"supplier_subscription.expense_recognition_run": r.ExpenseRecognitionRunURL,
	}
}

// ---------------------------------------------------------------------------
// ExpenseRecognitionRunRoutes — Plan A 20260517-expense-run
// ---------------------------------------------------------------------------

// ExpenseRecognitionRunRoutes holds all route paths for the Expense Recognition
// Run (buying-side) module. Mirrors RevenueRunRoutes shape.
// Surfaces: A (per-supplier drawer — entydad), B (workspace queue),
// C (per-supplier-subscription drawer), D (run history list + detail).
type ExpenseRecognitionRunRoutes struct {
	// Sidebar navigation context.
	ActiveNav string `json:"active_nav"`

	// Surface B — workspace queue page.
	QueueURL      string `json:"queue_url"`
	QueueTableURL string `json:"queue_table_url"`

	// Surface D — run history list + detail.
	ListURL            string `json:"list_url"`
	ListTableURL       string `json:"list_table_url"`
	DetailURL          string `json:"detail_url"`
	DetailTabActionURL string `json:"detail_tab_action_url"`

	// Action endpoints.
	NewURL         string `json:"new_url"`
	GenerateURL    string `json:"generate_url"`
	SubmitBatchURL string `json:"submit_batch_url"`

	// Surface A — per-supplier drawer (entydad supplier statement tab).
	PerSupplierDrawerURL string `json:"per_supplier_drawer_url"`

	// Surface C — per-SupplierSubscription drawer.
	PerSubscriptionDrawerURL string `json:"per_subscription_drawer_url"`
}

// DefaultExpenseRecognitionRunRoutes returns ExpenseRecognitionRunRoutes
// populated from the package-level route constants.
func DefaultExpenseRecognitionRunRoutes() ExpenseRecognitionRunRoutes {
	return ExpenseRecognitionRunRoutes{
		ActiveNav:                "expense-recognition-run",
		QueueURL:                 ExpenseRecognitionRunQueueURL,
		QueueTableURL:            ExpenseRecognitionRunQueueTableURL,
		ListURL:                  ExpenseRecognitionRunListURL,
		ListTableURL:             ExpenseRecognitionRunListTableURL,
		DetailURL:                ExpenseRecognitionRunDetailURL,
		DetailTabActionURL:       ExpenseRecognitionRunDetailTabActionURL,
		NewURL:                   ExpenseRecognitionRunNewURL,
		GenerateURL:              ExpenseRecognitionRunGenerateURL,
		SubmitBatchURL:           ExpenseRecognitionRunSubmitBatchURL,
		PerSupplierDrawerURL:     ExpenseRecognitionRunPerSupplierDrawerURL,
		PerSubscriptionDrawerURL: ExpenseRecognitionRunPerSubscriptionDrawerURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// expense-recognition-run routes.
func (r ExpenseRecognitionRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition_run.queue":                   r.QueueURL,
		"expense_recognition_run.queue_table":             r.QueueTableURL,
		"expense_recognition_run.list":                    r.ListURL,
		"expense_recognition_run.list_table":              r.ListTableURL,
		"expense_recognition_run.detail":                  r.DetailURL,
		"expense_recognition_run.detail_tab_action":       r.DetailTabActionURL,
		"expense_recognition_run.new":                     r.NewURL,
		"expense_recognition_run.generate":                r.GenerateURL,
		"expense_recognition_run.submit_batch":            r.SubmitBatchURL,
		"expense_recognition_run.per_supplier_drawer":     r.PerSupplierDrawerURL,
		"expense_recognition_run.per_subscription_drawer": r.PerSubscriptionDrawerURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — CostSchedule routes
// ---------------------------------------------------------------------------

// CostScheduleRoutes holds all route paths for cost_schedule views.
type CostScheduleRoutes struct {
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

// DefaultCostScheduleRoutes returns a CostScheduleRoutes using route constants.
func DefaultCostScheduleRoutes() CostScheduleRoutes {
	return CostScheduleRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-schedules",

		ListURL:          CostScheduleListURL,
		TableURL:         CostScheduleTableURL,
		DetailURL:        CostScheduleDetailURL,
		AddURL:           CostScheduleAddURL,
		EditURL:          CostScheduleEditURL,
		DeleteURL:        CostScheduleDeleteURL,
		BulkDeleteURL:    CostScheduleBulkDeleteURL,
		SetStatusURL:     CostScheduleSetStatusURL,
		BulkSetStatusURL: CostScheduleBulkSetStatusURL,
		TabActionURL:     CostScheduleTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostScheduleRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// P3 — SupplierPlan routes
// ---------------------------------------------------------------------------

// SupplierPlanRoutes holds all route paths for supplier_plan views.
type SupplierPlanRoutes struct {
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

	// Autocomplete search URL for the supplier select in add/edit forms.
	SearchSupplierURL string `json:"search_supplier_url"`
}

// DefaultSupplierPlanRoutes returns a SupplierPlanRoutes using route constants.
func DefaultSupplierPlanRoutes() SupplierPlanRoutes {
	return SupplierPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-plans",

		ListURL:          SupplierPlanListURL,
		TableURL:         SupplierPlanTableURL,
		DetailURL:        SupplierPlanDetailURL,
		AddURL:           SupplierPlanAddURL,
		EditURL:          SupplierPlanEditURL,
		DeleteURL:        SupplierPlanDeleteURL,
		BulkDeleteURL:    SupplierPlanBulkDeleteURL,
		SetStatusURL:     SupplierPlanSetStatusURL,
		BulkSetStatusURL: SupplierPlanBulkSetStatusURL,
		TabActionURL:     SupplierPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_plan.list":            r.ListURL,
		"supplier_plan.table":           r.TableURL,
		"supplier_plan.detail":          r.DetailURL,
		"supplier_plan.add":             r.AddURL,
		"supplier_plan.edit":            r.EditURL,
		"supplier_plan.delete":          r.DeleteURL,
		"supplier_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_plan.set_status":      r.SetStatusURL,
		"supplier_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_plan.tab_action":      r.TabActionURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — CostPlan routes
// ---------------------------------------------------------------------------

// CostPlanRoutes holds all route paths for cost_plan views.
type CostPlanRoutes struct {
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

	// SupplierProductCostPlan inline CRUD within cost_plan detail
	ProductCostAddURL    string `json:"product_cost_add_url"`
	ProductCostEditURL   string `json:"product_cost_edit_url"`
	ProductCostDeleteURL string `json:"product_cost_delete_url"`

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL        string `json:"search_supplier_plan_url"`
	SearchCostScheduleURL        string `json:"search_cost_schedule_url"`
	SearchSupplierProductPlanURL string `json:"search_supplier_product_plan_url"`
}

// DefaultCostPlanRoutes returns a CostPlanRoutes using route constants.
func DefaultCostPlanRoutes() CostPlanRoutes {
	return CostPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-plans",

		ListURL:              CostPlanListURL,
		TableURL:             CostPlanTableURL,
		DetailURL:            CostPlanDetailURL,
		AddURL:               CostPlanAddURL,
		EditURL:              CostPlanEditURL,
		DeleteURL:            CostPlanDeleteURL,
		BulkDeleteURL:        CostPlanBulkDeleteURL,
		SetStatusURL:         CostPlanSetStatusURL,
		BulkSetStatusURL:     CostPlanBulkSetStatusURL,
		TabActionURL:         CostPlanTabActionURL,
		ProductCostAddURL:    CostPlanProductCostAddURL,
		ProductCostEditURL:   CostPlanProductCostEditURL,
		ProductCostDeleteURL: CostPlanProductCostDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"cost_plan.list":                r.ListURL,
		"cost_plan.table":               r.TableURL,
		"cost_plan.detail":              r.DetailURL,
		"cost_plan.add":                 r.AddURL,
		"cost_plan.edit":                r.EditURL,
		"cost_plan.delete":              r.DeleteURL,
		"cost_plan.bulk_delete":         r.BulkDeleteURL,
		"cost_plan.set_status":          r.SetStatusURL,
		"cost_plan.bulk_set_status":     r.BulkSetStatusURL,
		"cost_plan.tab_action":          r.TabActionURL,
		"cost_plan.product_cost.add":    r.ProductCostAddURL,
		"cost_plan.product_cost.edit":   r.ProductCostEditURL,
		"cost_plan.product_cost.delete": r.ProductCostDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan routes
// ---------------------------------------------------------------------------

// SupplierProductPlanRoutes holds all route paths for supplier_product_plan views.
type SupplierProductPlanRoutes struct {
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

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL string `json:"search_supplier_plan_url"`
	SearchProductURL      string `json:"search_product_url"`
}

// DefaultSupplierProductPlanRoutes returns a SupplierProductPlanRoutes using route constants.
func DefaultSupplierProductPlanRoutes() SupplierProductPlanRoutes {
	return SupplierProductPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-product-plans",

		ListURL:          SupplierProductPlanListURL,
		TableURL:         SupplierProductPlanTableURL,
		DetailURL:        SupplierProductPlanDetailURL,
		AddURL:           SupplierProductPlanAddURL,
		EditURL:          SupplierProductPlanEditURL,
		DeleteURL:        SupplierProductPlanDeleteURL,
		BulkDeleteURL:    SupplierProductPlanBulkDeleteURL,
		SetStatusURL:     SupplierProductPlanSetStatusURL,
		BulkSetStatusURL: SupplierProductPlanBulkSetStatusURL,
		TabActionURL:     SupplierProductPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierProductPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_product_plan.list":            r.ListURL,
		"supplier_product_plan.table":           r.TableURL,
		"supplier_product_plan.detail":          r.DetailURL,
		"supplier_product_plan.add":             r.AddURL,
		"supplier_product_plan.edit":            r.EditURL,
		"supplier_product_plan.delete":          r.DeleteURL,
		"supplier_product_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_product_plan.set_status":      r.SetStatusURL,
		"supplier_product_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_product_plan.tab_action":      r.TabActionURL,
	}
}

// MapTableLabels is a shared helper used across all centymo view modules to
// produce a types.TableLabels from pyeza CommonLabels. Defined here to avoid
// duplication; all block module wirings call this.
func mapTableLabelsFromStrings(search, searchPlaceholder, sortAsc, sortDesc, noResults, loading string) struct{} {
	// Placeholder — actual implementation lives in the block package; this
	// comment documents the cross-module convention.
	return struct{}{}
}
