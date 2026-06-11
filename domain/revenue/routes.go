package revenue

// Default route constants for revenue views.
// Consumer apps can use these or define their own via lyngua route.json overrides.
const (
	// Revenue routes
	RevenueDashboardURL     = "/sales/dashboard"
	RevenueListURL          = "/sales/list/{status}"
	RevenueTableURL         = "/action/revenue/table/{status}"
	RevenueDetailURL        = "/sales/detail/{id}"
	RevenueAddURL           = "/action/revenue/add"
	RevenueEditURL          = "/action/revenue/edit/{id}"
	RevenueDeleteURL        = "/action/revenue/delete"
	RevenueBulkDeleteURL    = "/action/revenue/bulk-delete"
	RevenueSetStatusURL     = "/action/revenue/set-status"
	RevenueBulkSetStatusURL = "/action/revenue/bulk-set-status"

	// Revenue tab action route
	RevenueTabActionURL        = "/action/revenue/detail/{id}/tab/{tab}"
	RevenueAttachmentUploadURL = "/action/revenue/detail/{id}/attachments/upload"
	RevenueAttachmentDeleteURL = "/action/revenue/detail/{id}/attachments/delete"

	// Revenue line item routes (within revenue detail)
	RevenueLineItemTableURL    = "/action/revenue/detail/{id}/items/table"
	RevenueLineItemAddURL      = "/action/revenue/detail/{id}/items/add"
	RevenueLineItemEditURL     = "/action/revenue/detail/{id}/items/edit/{itemId}"
	RevenueLineItemRemoveURL   = "/action/revenue/detail/{id}/items/remove"
	RevenueLineItemDiscountURL = "/action/revenue/detail/{id}/items/add-discount"

	// Revenue payment routes (within revenue detail)
	RevenuePaymentTableURL  = "/action/revenue/detail/{id}/payment/table"
	RevenuePaymentAddURL    = "/action/revenue/detail/{id}/payment/add"
	RevenuePaymentEditURL   = "/action/revenue/detail/{id}/payment/edit/{pid}"
	RevenuePaymentRemoveURL = "/action/revenue/detail/{id}/payment/remove"

	// Revenue report routes
	RevenueSummaryURL = "/sales/reports/sales-summary"

	// Revenue invoice document routes
	RevenueInvoiceDownloadURL = "/action/revenue/detail/{id}/invoice/download"
	RevenueEmailURL           = "/action/revenue/detail/{id}/invoice/send-email"

	// Revenue settings routes (template management)
	RevenueSettingsTemplatesURL       = "/sales/settings/templates"
	RevenueSettingsTemplateUploadURL  = "/action/revenue/settings/templates/upload"
	RevenueSettingsTemplateDeleteURL  = "/action/revenue/settings/templates/delete"
	RevenueSettingsTemplateDefaultURL = "/action/revenue/settings/templates/set-default/{id}"
	RevenueSearchClientURL            = "/action/revenue/search/clients"
	RevenueSearchSubscriptionURL      = "/action/revenue/search/subscriptions"
	RevenueSearchLocationURL          = "/action/revenue/search/locations"
	RevenueSearchProductURL           = "/action/revenue/search/products"
	RevenuePriceLookupURL             = "/action/revenue/price-lookup"
	RevenueRecomputeTaxesURL          = "/action/revenue/detail/{id}/taxes/recompute"

	// Revenue Run (invoice-run) routes
	RevenueRunQueueURL            = "/revenue-run/queue"
	RevenueRunQueueTableURL       = "/action/revenue-run/queue/table"
	RevenueRunListURL             = "/revenue-run/list/{status}"
	RevenueRunListTableURL        = "/action/revenue-run/table/{status}"
	RevenueRunDetailURL           = "/revenue-run/detail/{id}"
	RevenueRunDetailTabActionURL  = "/action/revenue-run/detail/{id}/tab/{tab}"
	RevenueRunAttachmentUploadURL = "/action/revenue-run/detail/{id}/attachments/upload"
	RevenueRunAttachmentDeleteURL = "/action/revenue-run/detail/{id}/attachments/delete"
	RevenueRunSubmitBatchURL      = "/action/revenue-run/submit-batch"
)

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

	// Tax recompute (Phase 4 wiring — stub until ComputeTaxesForRevenue is available)
	RecomputeTaxesURL string `json:"recompute_taxes_url"`
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
		RecomputeTaxesURL:          RevenueRecomputeTaxesURL,
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

// RevenueRunRoutes holds all route paths for the Revenue Run (invoice-run) module.
// Surface B = workspace queue page; Surface D = run history list + detail pages.
type RevenueRunRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override.
	ActiveNav string `json:"active_nav"`

	QueueURL            string `json:"queue_url"`
	QueueTableURL       string `json:"queue_table_url"`
	ListURL             string `json:"list_url"`
	ListTableURL        string `json:"list_table_url"`
	DetailURL           string `json:"detail_url"`
	DetailTabActionURL  string `json:"detail_tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
	SubmitBatchURL      string `json:"submit_batch_url"`
}

// DefaultRevenueRunRoutes returns a RevenueRunRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultRevenueRunRoutes() RevenueRunRoutes {
	return RevenueRunRoutes{
		ActiveNav:           "revenue-run",
		QueueURL:            RevenueRunQueueURL,
		QueueTableURL:       RevenueRunQueueTableURL,
		ListURL:             RevenueRunListURL,
		ListTableURL:        RevenueRunListTableURL,
		DetailURL:           RevenueRunDetailURL,
		DetailTabActionURL:  RevenueRunDetailTabActionURL,
		AttachmentUploadURL: RevenueRunAttachmentUploadURL,
		AttachmentDeleteURL: RevenueRunAttachmentDeleteURL,
		SubmitBatchURL:      RevenueRunSubmitBatchURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// revenue-run routes.
func (r RevenueRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"revenue_run.queue":             r.QueueURL,
		"revenue_run.queue_table":       r.QueueTableURL,
		"revenue_run.list":              r.ListURL,
		"revenue_run.list_table":        r.ListTableURL,
		"revenue_run.detail":            r.DetailURL,
		"revenue_run.detail_tab_action": r.DetailTabActionURL,
		"revenue_run.attachment.upload": r.AttachmentUploadURL,
		"revenue_run.attachment.delete": r.AttachmentDeleteURL,
		"revenue_run.submit_batch":      r.SubmitBatchURL,
	}
}
