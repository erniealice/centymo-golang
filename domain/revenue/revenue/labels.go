package revenue

// ---------------------------------------------------------------------------
// Revenue labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the revenue module.
type Labels struct {
	Page      PageLabels      `json:"page"`
	Buttons   ButtonLabels    `json:"buttons"`
	Columns   ColumnLabels    `json:"columns"`
	Empty     EmptyLabels     `json:"empty"`
	Form      FormLabels      `json:"form"`
	Actions   ActionLabels    `json:"actions"`
	Bulk      BulkLabels      `json:"bulk_actions"`
	Detail    DetailLabels    `json:"detail"`
	Confirm   ConfirmLabels   `json:"confirm"`
	Errors    ErrorLabels     `json:"errors"`
	Dashboard DashboardLabels `json:"dashboard"`
	Settings  SettingsLabels  `json:"settings"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"heading_draft"`
	HeadingComplete  string `json:"heading_complete"`
	HeadingCancelled string `json:"heading_cancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"caption_draft"`
	CaptionComplete  string `json:"caption_complete"`
	CaptionCancelled string `json:"caption_cancelled"`
}

type ButtonLabels struct {
	AddSale string `json:"add_sale"`
}

type ColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	DraftTitle       string `json:"draft_title"`
	DraftMessage     string `json:"draft_message"`
	CompleteTitle    string `json:"complete_title"`
	CompleteMessage  string `json:"complete_message"`
	CancelledTitle   string `json:"cancelled_title"`
	CancelledMessage string `json:"cancelled_message"`
}

type FormLabels struct {
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"reference_placeholder"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	NotesPlaceholder     string `json:"notes_placeholder"`
	Active               string `json:"active"`
	Location             string `json:"location"`

	// Payment terms and client search labels
	PaymentTerms              string `json:"payment_terms"`
	SelectPaymentTerm         string `json:"select_payment_term"`
	DueDate                   string `json:"due_date"`
	CustomerSearchPlaceholder string `json:"customer_search_placeholder"`
	CustomerNoResults         string `json:"customer_no_results"`

	// Subscription search labels
	Subscription          string `json:"subscription"`
	SubscriptionNoResults string `json:"subscription_no_results"`

	// Placeholders and translated option labels
	CurrencyPlaceholder            string `json:"currency_placeholder"`
	CustomerNamePlaceholder        string `json:"customer_name_placeholder"`
	StatusDraft                    string `json:"status_draft"`
	StatusComplete                 string `json:"status_complete"`
	StatusCancelled                string `json:"status_cancelled"`
	PaymentMethod                  string `json:"payment_method"`
	ReferenceNumber                string `json:"reference_number"`
	TransactionIdPlaceholder       string `json:"transaction_id_placeholder"`
	ReceivedBy                     string `json:"received_by"`
	Role                           string `json:"role"`
	SelectInventoryItem            string `json:"select_inventory_item"`
	ItemDescriptionPlaceholder     string `json:"item_description_placeholder"`
	DiscountDescriptionPlaceholder string `json:"discount_description_placeholder"`

	// Field-level info text for the payment drawer form.
	PaymentMethodInfo   string `json:"payment_method_info"`
	AmountInfo          string `json:"amount_info"`
	CurrencyInfo        string `json:"currency_info"`
	ReferenceNumberInfo string `json:"reference_number_info"`
	ReceivedByInfo      string `json:"received_by_info"`
	RoleInfo            string `json:"role_info"`
	NotesInfo           string `json:"notes_info"`
}

type ActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Complete          string `json:"complete"`
	Reactivate        string `json:"reactivate"`
	DownloadInvoice   string `json:"download_invoice"`
	SendEmail         string `json:"send_email"`
	Cancel            string `json:"cancel"`
	ReclassifyToDraft string `json:"reclassify_to_draft"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle   string `json:"page_title"`
	TitlePrefix string `json:"title_prefix"`
	InvoiceInfo string `json:"invoice_info"`
	LineItems   string `json:"line_items"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	CostPrice   string `json:"cost_price"`
	GrossProfit string `json:"gross_profit"`
	Total       string `json:"total"`
	Discount    string `json:"discount"`
	SubTotal    string `json:"sub_total"`
	GrandTotal  string `json:"grand_total"`

	// Tab labels
	TabBasicInfo    string `json:"tab_basic_info"`
	TabLineItems    string `json:"tab_line_items"`
	TabPayment      string `json:"tab_payment"`
	TabAttachments  string `json:"tab_attachments"`
	TabAuditTrail   string `json:"tab_audit_trail"`
	TabAuditHistory string `json:"tab_audit_history"`

	// Basic info fields
	Customer     string `json:"customer"`
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	PaymentTerms string `json:"payment_terms"`
	DueDate      string `json:"due_date"`

	// Payment fields
	PaymentMethod string `json:"payment_method"`
	AmountPaid    string `json:"amount_paid"`
	CardDetails   string `json:"card_details"`
	PaymentDate   string `json:"payment_date"`
	ReceivedBy    string `json:"received_by"`
	PaymentInfo   string `json:"payment_info"`

	// Audit trail
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditAction          string `json:"audit_action"`
	AuditUser            string `json:"audit_user"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`

	// Totals
	TotalGrossProfit string `json:"total_gross_profit"`

	// Payment empty/table
	Reference           string `json:"reference"`
	PaymentEmptyTitle   string `json:"payment_empty_title"`
	PaymentEmptyMessage string `json:"payment_empty_message"`

	// Line item management
	AddItem                    string `json:"add_item"`
	AddDiscount                string `json:"add_discount"`
	EditItem                   string `json:"edit_item"`
	RemoveItem                 string `json:"remove_item"`
	ItemType                   string `json:"item_type"`
	ItemTypeItem               string `json:"item_type_item"`
	ItemTypeDiscount           string `json:"item_type_discount"`
	InventoryItem              string `json:"inventory_item"`
	SelectInventoryItem        string `json:"select_inventory_item"`
	ItemDescriptionPlaceholder string `json:"item_description_placeholder"`
	NotesPlaceholder           string `json:"notes_placeholder"`
	SerialNumber               string `json:"serial_number"`
	Product                    string `json:"product"`
	ProductNoResults           string `json:"product_no_results"`
	ProductPlaceholder         string `json:"product_placeholder"`
	ItemEmptyTitle             string `json:"item_empty_title"`
	ItemEmptyMessage           string `json:"item_empty_message"`

	// Field-level info text for the line-item drawer form.
	ProductInfo     string `json:"product_info"`
	DescriptionInfo string `json:"description_info"`
	QuantityInfo    string `json:"quantity_info"`
	UnitPriceInfo   string `json:"unit_price_info"`
	CostPriceInfo   string `json:"cost_price_info"`
	DiscountInfo    string `json:"discount_info"`
	NotesInfo       string `json:"notes_info"`

	// Payment tab
	TotalPaid                  string `json:"total_paid"`
	Remaining                  string `json:"remaining"`
	RecordPayment              string `json:"record_payment"`
	NoPaymentInfo              string `json:"no_payment_info"`
	PaymentDetailsNotAvailable string `json:"payment_details_not_available"`
}

type ConfirmLabels struct {
	Complete                 string `json:"complete"`
	CompleteMessage          string `json:"complete_message"`
	Reactivate               string `json:"reactivate"`
	ReactivateMessage        string `json:"reactivate_message"`
	BulkComplete             string `json:"bulk_complete"`
	BulkCompleteMessage      string `json:"bulk_complete_message"`
	BulkReactivate           string `json:"bulk_reactivate"`
	BulkReactivateMessage    string `json:"bulk_reactivate_message"`
	SendEmail                string `json:"send_email"`
	SendEmailMessage         string `json:"send_email_message"`
	Cancel                   string `json:"cancel"`
	CancelMessage            string `json:"cancel_message"`
	ReclassifyToDraft        string `json:"reclassify_to_draft"`
	ReclassifyToDraftMessage string `json:"reclassify_to_draft_message"`
}

type ErrorLabels struct {
	PermissionDenied        string `json:"permission_denied"`
	InvalidFormData         string `json:"invalid_form_data"`
	NotFound                string `json:"not_found"`
	IDRequired              string `json:"id_required"`
	NoIDsProvided           string `json:"no_ids_provided"`
	InvalidStatus           string `json:"invalid_status"`
	InvalidTargetStatus     string `json:"invalid_target_status"`
	NoItemsCannotComplete   string `json:"no_items_cannot_complete"`
	HasPaymentsCannotCancel string `json:"has_payments_cannot_cancel"`
	BulkHasPayments         string `json:"bulk_has_payments"`
	BulkNoItems             string `json:"bulk_no_items"`
	PaymentNotFound         string `json:"payment_not_found"`
	InvalidDiscount         string `json:"invalid_discount"`
	// RecomputeUnavailable is the 501 body returned by the RecomputeTaxes stub
	// until Phase 4 wires ComputeTaxesForRevenue (Phase 5 M2).
	RecomputeUnavailable string `json:"recompute_unavailable"`
}

type DashboardLabels struct {
	Title             string `json:"title"`
	TotalRevenue      string `json:"total_revenue"`
	Revenue           string `json:"revenue"`
	Completed         string `json:"completed"`
	Active            string `json:"active"`
	RevenueTrend      string `json:"revenue_trend"`
	Week              string `json:"week"`
	Month             string `json:"month"`
	Year              string `json:"year"`
	RecentRevenue     string `json:"recent_revenue"`
	ViewAll           string `json:"view_all"`
	NewRevenueCreated string `json:"new_revenue_created"`
	RevenueCompleted  string `json:"revenue_completed"`
	RevenueUpdated    string `json:"revenue_updated"`
	RevenueCancelled  string `json:"revenue_cancelled"`
	QuickNewRevenue   string `json:"quick_new_revenue"`
	QuickViewAll      string `json:"quick_view_all"`
}

// SettingsLabels holds translatable strings for the revenue settings page
// (invoice template management).
type SettingsLabels struct {
	PageTitle      string `json:"page_title"`
	Caption        string `json:"caption"`
	UploadTemplate string `json:"upload_template"`
	TemplateName   string `json:"template_name"`
	TemplateType   string `json:"template_type"`
	Purpose        string `json:"purpose"`
	SetDefault     string `json:"set_default"`
	Delete         string `json:"delete"`
	DefaultBadge   string `json:"default_badge"`
	EmptyTitle     string `json:"empty_title"`
	EmptyMessage   string `json:"empty_message"`
	UploadSuccess  string `json:"upload_success"`
	DeleteConfirm  string `json:"delete_confirm"`
}

// DefaultLabels returns the zero-value label set. Every rendered string for
// this entity must come from the lyngua cascade (general -> business-type
// tier); there are no Go-side default strings to fall back on.
func DefaultLabels() Labels { return Labels{} }
