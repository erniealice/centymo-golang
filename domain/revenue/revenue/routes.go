package revenue

// Default route constants for revenue views.
// Consumer apps can use these or define their own via lyngua route.json overrides.
const (
	// Revenue routes
	DashboardURL     = "/sales/dashboard"
	ListURL          = "/sales/list/{status}"
	TableURL         = "/action/revenue/table/{status}"
	DetailURL        = "/sales/detail/{id}"
	AddURL           = "/action/revenue/add"
	EditURL          = "/action/revenue/edit/{id}"
	DeleteURL        = "/action/revenue/delete"
	BulkDeleteURL    = "/action/revenue/bulk-delete"
	SetStatusURL     = "/action/revenue/set-status"
	BulkSetStatusURL = "/action/revenue/bulk-set-status"

	// Revenue tab action route
	TabActionURL        = "/action/revenue/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/revenue/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/revenue/detail/{id}/attachments/delete"

	// Revenue line item routes (within revenue detail)
	LineItemTableURL    = "/action/revenue/detail/{id}/items/table"
	LineItemAddURL      = "/action/revenue/detail/{id}/items/add"
	LineItemEditURL     = "/action/revenue/detail/{id}/items/edit/{itemId}"
	LineItemRemoveURL   = "/action/revenue/detail/{id}/items/remove"
	LineItemDiscountURL = "/action/revenue/detail/{id}/items/add-discount"

	// Revenue payment routes (within revenue detail)
	PaymentTableURL  = "/action/revenue/detail/{id}/payment/table"
	PaymentAddURL    = "/action/revenue/detail/{id}/payment/add"
	PaymentEditURL   = "/action/revenue/detail/{id}/payment/edit/{pid}"
	PaymentRemoveURL = "/action/revenue/detail/{id}/payment/remove"

	// Revenue report routes
	SummaryURL = "/sales/reports/sales-summary"

	// Revenue invoice document routes
	InvoiceDownloadURL = "/action/revenue/detail/{id}/invoice/download"
	EmailURL           = "/action/revenue/detail/{id}/invoice/send-email"

	// Revenue settings routes (template management)
	SettingsTemplatesURL       = "/sales/settings/templates"
	SettingsTemplateUploadURL  = "/action/revenue/settings/templates/upload"
	SettingsTemplateDeleteURL  = "/action/revenue/settings/templates/delete"
	SettingsTemplateDefaultURL = "/action/revenue/settings/templates/set-default/{id}"
	SearchClientURL            = "/action/revenue/search/clients"
	SearchSubscriptionURL      = "/action/revenue/search/subscriptions"
	SearchLocationURL          = "/action/revenue/search/locations"
	SearchProductURL           = "/action/revenue/search/products"
	PriceLookupURL             = "/action/revenue/price-lookup"
	RecomputeTaxesURL          = "/action/revenue/detail/{id}/taxes/recompute"
)

// Routes holds all route paths for revenue views and actions,
// including line item and payment sub-routes.
type Routes struct {
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

	// Tax recompute (Phase 4 wiring — stub until ComputeTaxesForRevenue is available)
	RecomputeTaxesURL string `json:"recompute_taxes_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,

		TabActionURL: TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		LineItemTableURL:    LineItemTableURL,
		LineItemAddURL:      LineItemAddURL,
		LineItemEditURL:     LineItemEditURL,
		LineItemRemoveURL:   LineItemRemoveURL,
		LineItemDiscountURL: LineItemDiscountURL,

		PaymentTableURL:  PaymentTableURL,
		PaymentAddURL:    PaymentAddURL,
		PaymentEditURL:   PaymentEditURL,
		PaymentRemoveURL: PaymentRemoveURL,

		RevenueSummaryURL:          SummaryURL,
		InvoiceDownloadURL:         InvoiceDownloadURL,
		SendEmailURL:               EmailURL,
		SettingsTemplatesURL:       SettingsTemplatesURL,
		SettingsTemplateUploadURL:  SettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  SettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: SettingsTemplateDefaultURL,
		SearchClientURL:            SearchClientURL,
		SearchSubscriptionURL:      SearchSubscriptionURL,
		SearchLocationURL:          SearchLocationURL,
		SearchProductURL:           SearchProductURL,
		PriceLookupURL:             PriceLookupURL,
		RecomputeTaxesURL:          RecomputeTaxesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// revenue routes.
func (r Routes) RouteMap() map[string]string {
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

		"revenue.summary":                   r.RevenueSummaryURL,
		"revenue.invoice_download":          r.InvoiceDownloadURL,
		"revenue.send_email":                r.SendEmailURL,
		"revenue.settings.templates":        r.SettingsTemplatesURL,
		"revenue.settings.template_upload":  r.SettingsTemplateUploadURL,
		"revenue.settings.template_delete":  r.SettingsTemplateDeleteURL,
		"revenue.settings.template_default": r.SettingsTemplateDefaultURL,
		"revenue.search_client":             r.SearchClientURL,
		"revenue.search.subscriptions":      r.SearchSubscriptionURL,
		"revenue.search.locations":          r.SearchLocationURL,
		"revenue.search.products":           r.SearchProductURL,
		"revenue.price_lookup":              r.PriceLookupURL,
		"revenue.taxes.recompute":           r.RecomputeTaxesURL,
	}
}
