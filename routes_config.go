package centymo

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
	DetailURL     string `json:"detail_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`

	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	TabActionURL string `json:"tab_action_url"`

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

	// Variant stock routes
	VariantStockDetailURL    string `json:"variant_stock_detail_url"`
	VariantStockTabActionURL string `json:"variant_stock_tab_action_url"`

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
		DetailURL:     ProductDetailURL,
		AddURL:        ProductAddURL,
		EditURL:       ProductEditURL,
		DeleteURL:     ProductDeleteURL,
		BulkDeleteURL: ProductBulkDeleteURL,

		SetStatusURL:     ProductSetStatusURL,
		BulkSetStatusURL: ProductBulkSetStatusURL,

		TabActionURL: ProductTabActionURL,

		VariantTableURL:  ProductVariantTableURL,
		VariantAssignURL: ProductVariantAssignURL,
		VariantEditURL:   ProductVariantEditURL,
		VariantRemoveURL: ProductVariantRemoveURL,

		VariantDetailURL:    ProductVariantDetailURL,
		VariantTabActionURL: ProductVariantTabActionURL,

		VariantImageUploadURL: ProductVariantImageUploadURL,
		VariantImageDeleteURL: ProductVariantImageDeleteURL,

		VariantStockDetailURL:    ProductVariantStockDetailURL,
		VariantStockTabActionURL: ProductVariantStockTabActionURL,

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

// RouteMap returns a map of dot-notation keys to route paths for all
// product routes.
func (r ProductRoutes) RouteMap() map[string]string {
	return map[string]string{
		"product.list":        r.ListURL,
		"product.detail":      r.DetailURL,
		"product.add":         r.AddURL,
		"product.edit":        r.EditURL,
		"product.delete":      r.DeleteURL,
		"product.bulk_delete": r.BulkDeleteURL,

		"product.set_status":      r.SetStatusURL,
		"product.bulk_set_status": r.BulkSetStatusURL,

		"product.tab_action": r.TabActionURL,

		"product.variant.table":  r.VariantTableURL,
		"product.variant.assign": r.VariantAssignURL,
		"product.variant.edit":   r.VariantEditURL,
		"product.variant.remove": r.VariantRemoveURL,

		"product.variant.detail":     r.VariantDetailURL,
		"product.variant.tab_action": r.VariantTabActionURL,

		"product.variant.image.upload": r.VariantImageUploadURL,
		"product.variant.image.delete": r.VariantImageDeleteURL,

		"product.variant.stock.detail":     r.VariantStockDetailURL,
		"product.variant.stock.tab_action": r.VariantStockTabActionURL,

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

// SalesRoutes holds all route paths for sales (revenue) views and actions,
// including line item and payment sub-routes.
type SalesRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	TabActionURL string `json:"tab_action_url"`

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
	SalesSummaryURL string `json:"sales_summary_url"`

	// Document generation routes
	InvoiceDownloadURL string `json:"invoice_download_url"`

	// Settings routes (template management)
	SettingsTemplatesURL       string `json:"settings_templates_url"`
	SettingsTemplateUploadURL  string `json:"settings_template_upload_url"`
	SettingsTemplateDeleteURL  string `json:"settings_template_delete_url"`
	SettingsTemplateDefaultURL string `json:"settings_template_default_url"`
}

// DefaultSalesRoutes returns a SalesRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultSalesRoutes() SalesRoutes {
	return SalesRoutes{
		DashboardURL:     SalesDashboardURL,
		ListURL:          SalesListURL,
		DetailURL:        SalesDetailURL,
		AddURL:           SalesAddURL,
		EditURL:          SalesEditURL,
		DeleteURL:        SalesDeleteURL,
		BulkDeleteURL:    SalesBulkDeleteURL,
		SetStatusURL:     SalesSetStatusURL,
		BulkSetStatusURL: SalesBulkSetStatusURL,

		TabActionURL: SalesTabActionURL,

		LineItemTableURL:    SalesLineItemTableURL,
		LineItemAddURL:      SalesLineItemAddURL,
		LineItemEditURL:     SalesLineItemEditURL,
		LineItemRemoveURL:   SalesLineItemRemoveURL,
		LineItemDiscountURL: SalesLineItemDiscountURL,

		PaymentTableURL:  SalesPaymentTableURL,
		PaymentAddURL:    SalesPaymentAddURL,
		PaymentEditURL:   SalesPaymentEditURL,
		PaymentRemoveURL: SalesPaymentRemoveURL,

		SalesSummaryURL: SalesSummaryURL,
		InvoiceDownloadURL:         SalesInvoiceDownloadURL,
		SettingsTemplatesURL:       SalesSettingsTemplatesURL,
		SettingsTemplateUploadURL:  SalesSettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  SalesSettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: SalesSettingsTemplateDefaultURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// sales routes.
func (r SalesRoutes) RouteMap() map[string]string {
	return map[string]string{
		"sales.dashboard":       r.DashboardURL,
		"sales.list":            r.ListURL,
		"sales.detail":          r.DetailURL,
		"sales.add":             r.AddURL,
		"sales.edit":            r.EditURL,
		"sales.delete":          r.DeleteURL,
		"sales.bulk_delete":     r.BulkDeleteURL,
		"sales.set_status":      r.SetStatusURL,
		"sales.bulk_set_status": r.BulkSetStatusURL,

		"sales.tab_action": r.TabActionURL,

		"sales.line_item.table":    r.LineItemTableURL,
		"sales.line_item.add":      r.LineItemAddURL,
		"sales.line_item.edit":     r.LineItemEditURL,
		"sales.line_item.remove":   r.LineItemRemoveURL,
		"sales.line_item.discount": r.LineItemDiscountURL,

		"sales.payment.table":  r.PaymentTableURL,
		"sales.payment.add":    r.PaymentAddURL,
		"sales.payment.edit":   r.PaymentEditURL,
		"sales.payment.remove": r.PaymentRemoveURL,

		"sales.sales_summary": r.SalesSummaryURL,
		"sales.invoice_download":           r.InvoiceDownloadURL,
		"sales.settings.templates":         r.SettingsTemplatesURL,
		"sales.settings.template_upload":   r.SettingsTemplateUploadURL,
		"sales.settings.template_delete":   r.SettingsTemplateDeleteURL,
		"sales.settings.template_default":  r.SettingsTemplateDefaultURL,
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
	}
}

// PlanRoutes holds all route paths for plan views and actions.
type PlanRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
	TabActionURL string `json:"tab_action_url"`
}

// DefaultPlanRoutes returns a PlanRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPlanRoutes() PlanRoutes {
	return PlanRoutes{
		ActiveNav:    "services",
		ActiveSubNav: "plans",

		ListURL:      PlanListURL,
		DetailURL:    PlanDetailURL,
		AddURL:       PlanAddURL,
		EditURL:      PlanEditURL,
		DeleteURL:    PlanDeleteURL,
		TabActionURL: PlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// plan routes.
func (r PlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"plan.list":       r.ListURL,
		"plan.detail":     r.DetailURL,
		"plan.add":        r.AddURL,
		"plan.edit":       r.EditURL,
		"plan.delete":     r.DeleteURL,
		"plan.tab_action": r.TabActionURL,
	}
}

// SubscriptionRoutes holds all route paths for subscription views and actions.
type SubscriptionRoutes struct {
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
	TabActionURL string `json:"tab_action_url"`
}

// DefaultSubscriptionRoutes returns a SubscriptionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultSubscriptionRoutes() SubscriptionRoutes {
	return SubscriptionRoutes{
		ListURL:      SubscriptionListURL,
		DetailURL:    SubscriptionDetailURL,
		AddURL:       SubscriptionAddURL,
		EditURL:      SubscriptionEditURL,
		DeleteURL:    SubscriptionDeleteURL,
		TabActionURL: SubscriptionTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription routes.
func (r SubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"subscription.list":       r.ListURL,
		"subscription.detail":     r.DetailURL,
		"subscription.add":        r.AddURL,
		"subscription.edit":       r.EditURL,
		"subscription.delete":     r.DeleteURL,
		"subscription.tab_action": r.TabActionURL,
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
	}
}

// PriceListRoutes holds all route paths for price list views and actions,
// including price product sub-routes.
type PriceListRoutes struct {
	ListURL       string `json:"list_url"`
	DetailURL     string `json:"detail_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`

	TabActionURL string `json:"tab_action_url"`

	// Price product routes
	PriceProductAddURL    string `json:"price_product_add_url"`
	PriceProductDeleteURL string `json:"price_product_delete_url"`
}

// DefaultPriceListRoutes returns a PriceListRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultPriceListRoutes() PriceListRoutes {
	return PriceListRoutes{
		ListURL:       PriceListListURL,
		DetailURL:     PriceListDetailURL,
		AddURL:        PriceListAddURL,
		EditURL:       PriceListEditURL,
		DeleteURL:     PriceListDeleteURL,
		BulkDeleteURL: PriceListBulkDeleteURL,

		TabActionURL: PriceListTabActionURL,

		PriceProductAddURL:    PriceProductAddURL,
		PriceProductDeleteURL: PriceProductDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price list routes.
func (r PriceListRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_list.list":        r.ListURL,
		"price_list.detail":      r.DetailURL,
		"price_list.add":         r.AddURL,
		"price_list.edit":        r.EditURL,
		"price_list.delete":      r.DeleteURL,
		"price_list.bulk_delete": r.BulkDeleteURL,

		"price_list.tab_action": r.TabActionURL,

		"price_list.price_product.add":    r.PriceProductAddURL,
		"price_list.price_product.delete": r.PriceProductDeleteURL,
	}
}
