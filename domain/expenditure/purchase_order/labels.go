package purchaseorder

// ---------------------------------------------------------------------------
// Purchase Order labels
// ---------------------------------------------------------------------------

// ErrorLabels holds error messages for the purchase order action handlers.
type ErrorLabels struct {
	NoPermission string `json:"no_permission"`
}

// Labels holds all translatable strings for the purchase order module.
type Labels struct {
	Labels    LabelNames     `json:"labels"`
	Page      PageLabels     `json:"page"`
	Buttons   ButtonLabels   `json:"buttons"`
	Columns   ColumnLabels   `json:"columns"`
	Empty     EmptyLabels    `json:"empty"`
	Form      FormLabels     `json:"form"`
	Status    StatusLabels   `json:"status"`
	POTypes   POTypeLabels   `json:"po_types"`
	LineTypes LineTypeLabels `json:"line_types"`
	Actions   ActionLabels   `json:"actions"`
	Bulk      BulkLabels     `json:"bulk_actions"`
	Detail    DetailLabels   `json:"detail"`
	LineItems LineItemLabels `json:"line_items"`
	Receipt   ReceiptLabels  `json:"receipt"`
	Errors    ErrorLabels    `json:"errors"`
}

type LabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"name_plural"`
	LineItem       string `json:"line_item"`
	LineItemPlural string `json:"line_item_plural"`
}

type PageLabels struct {
	Heading                  string `json:"heading"`
	Caption                  string `json:"caption"`
	HeadingDraft             string `json:"heading_draft"`
	HeadingPendingApproval   string `json:"heading_pending_approval"`
	HeadingApproved          string `json:"heading_approved"`
	HeadingPartiallyReceived string `json:"heading_partially_received"`
	HeadingFullyReceived     string `json:"heading_fully_received"`
	HeadingBilled            string `json:"heading_billed"`
	HeadingClosed            string `json:"heading_closed"`
	HeadingCancelled         string `json:"heading_cancelled"`
	Dashboard                string `json:"dashboard"`
}

type ButtonLabels struct {
	Add         string `json:"add"`
	AddLineItem string `json:"add_line_item"`
}

type ColumnLabels struct {
	PONumber        string `json:"po_number"`
	POType          string `json:"po_type"`
	Supplier        string `json:"supplier"`
	Location        string `json:"location"`
	OrderDate       string `json:"order_date"`
	Status          string `json:"status"`
	Currency        string `json:"currency"`
	Subtotal        string `json:"subtotal"`
	TaxAmount       string `json:"tax_amount"`
	TotalAmount     string `json:"total_amount"`
	PaymentTerms    string `json:"payment_terms"`
	ShippingTerms   string `json:"shipping_terms"`
	ApprovedBy      string `json:"approved_by"`
	ReferenceNumber string `json:"reference_number"`
	Notes           string `json:"notes"`
}

type EmptyLabels struct {
	Title                    string `json:"title"`
	Message                  string `json:"message"`
	DraftTitle               string `json:"draft_title"`
	DraftMessage             string `json:"draft_message"`
	PendingApprovalTitle     string `json:"pending_approval_title"`
	PendingApprovalMessage   string `json:"pending_approval_message"`
	ApprovedTitle            string `json:"approved_title"`
	ApprovedMessage          string `json:"approved_message"`
	PartiallyReceivedTitle   string `json:"partially_received_title"`
	PartiallyReceivedMessage string `json:"partially_received_message"`
	FullyReceivedTitle       string `json:"fully_received_title"`
	FullyReceivedMessage     string `json:"fully_received_message"`
	BilledTitle              string `json:"billed_title"`
	BilledMessage            string `json:"billed_message"`
	ClosedTitle              string `json:"closed_title"`
	ClosedMessage            string `json:"closed_message"`
	CancelledTitle           string `json:"cancelled_title"`
	CancelledMessage         string `json:"cancelled_message"`
}

type FormLabels struct {
	PONumber                   string `json:"po_number"`
	PONumberPlaceholder        string `json:"po_number_placeholder"`
	POType                     string `json:"po_type"`
	SelectPOType               string `json:"select_po_type"`
	Supplier                   string `json:"supplier"`
	SelectSupplier             string `json:"select_supplier"`
	Location                   string `json:"location"`
	SelectLocation             string `json:"select_location"`
	OrderDate                  string `json:"order_date"`
	Currency                   string `json:"currency"`
	Subtotal                   string `json:"subtotal"`
	TaxAmount                  string `json:"tax_amount"`
	TotalAmount                string `json:"total_amount"`
	PaymentTerms               string `json:"payment_terms"`
	ShippingTerms              string `json:"shipping_terms"`
	ApprovedBy                 string `json:"approved_by"`
	ReferenceNumber            string `json:"reference_number"`
	ReferenceNumberPlaceholder string `json:"reference_number_placeholder"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notes_placeholder"`
	SectionInfo                string `json:"section_info"`
	SectionSupplier            string `json:"section_supplier"`
	SectionFinancials          string `json:"section_financials"`
	SectionNotes               string `json:"section_notes"`

	// Field-level info text surfaced via an info button beside each label.
	PONumberInfo         string `json:"po_number_info"`
	POTypeInfo           string `json:"po_type_info"`
	SupplierInfo         string `json:"supplier_info"`
	OrderDateInfo        string `json:"order_date_info"`
	ExpectedDeliveryInfo string `json:"expected_delivery_info"`
	CurrencyInfo         string `json:"currency_info"`
	PaymentTermsInfo     string `json:"payment_terms_info"`
	ShippingTermsInfo    string `json:"shipping_terms_info"`
	ReferenceNumberInfo  string `json:"reference_number_info"`
	NotesInfo            string `json:"notes_info"`
}

type StatusLabels struct {
	Draft             string `json:"draft"`
	PendingApproval   string `json:"pending_approval"`
	Approved          string `json:"approved"`
	PartiallyReceived string `json:"partially_received"`
	FullyReceived     string `json:"fully_received"`
	Billed            string `json:"billed"`
	Closed            string `json:"closed"`
	Cancelled         string `json:"cancelled"`
}

type POTypeLabels struct {
	Standard string `json:"standard"`
	Blanket  string `json:"blanket"`
	Contract string `json:"contract"`
}

type LineTypeLabels struct {
	Goods   string `json:"goods"`
	Service string `json:"service"`
	Expense string `json:"expense"`
}

type ActionLabels struct {
	Cancel         string `json:"cancel"`
	Close          string `json:"close"`
	ConfirmReceipt string `json:"confirm_receipt"`
	Create         string `json:"create"`
	Delete         string `json:"delete"`
	Edit           string `json:"edit"`
	Approve        string `json:"approve"`
	Receive        string `json:"receive"`
	Reject         string `json:"reject"`
	View           string `json:"view"`
}

type BulkLabels struct {
	Delete  string `json:"delete"`
	Approve string `json:"approve"`
	Close   string `json:"close"`
}

// DetailLabels holds translatable strings for the PO detail page.
type DetailLabels struct {
	PageTitle            string `json:"page_title"`
	Title                string `json:"title"`
	InfoSection          string `json:"supplier_info"`
	Supplier             string `json:"supplier"`
	Location             string `json:"location"`
	OrderDate            string `json:"order_date"`
	PONumber             string `json:"po_number"`
	POType               string `json:"po_type"`
	Status               string `json:"status"`
	Currency             string `json:"currency"`
	Subtotal             string `json:"subtotal"`
	TaxAmount            string `json:"tax_amount"`
	TotalAmount          string `json:"total_amount"`
	PaymentTerms         string `json:"payment_terms"`
	ShippingTerms        string `json:"shipping_terms"`
	ApprovedBy           string `json:"approved_by"`
	ReferenceNumber      string `json:"reference_number"`
	Notes                string `json:"notes"`
	LineItems            string `json:"line_items"`
	Description          string `json:"description"`
	LineType             string `json:"line_type"`
	LineNumber           string `json:"line_number"`
	QuantityOrdered      string `json:"quantity_ordered"`
	QuantityReceived     string `json:"quantity_received"`
	QuantityBilled       string `json:"quantity_billed"`
	UnitPrice            string `json:"unit_price"`
	TotalPrice           string `json:"total_price"`
	SubTotal             string `json:"sub_total"`
	GrandTotal           string `json:"grand_total"`
	TabBasicInfo         string `json:"tab_basic_info"`
	TabLineItems         string `json:"tab_line_items"`
	TabReceiving         string `json:"tab_receiving"`
	TabAuditTrail        string `json:"tab_audit_trail"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditAction          string `json:"audit_action"`
	AuditUser            string `json:"audit_user"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	Total                string `json:"total"`
	AddLineItem          string `json:"add_line_item"`
	NoLineItems          string `json:"no_line_items"`
	ConfirmReceiptBtn    string `json:"confirm_receipt_btn"`
	TabAttachments       string `json:"tab_attachments"`
}

// LineItemLabels holds translatable strings for the PO line item drawer form.
type LineItemLabels struct {
	AddItem                string `json:"add_item"`
	AddLineItem            string `json:"add_line_item"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	EditItem               string `json:"edit_item"`
	EditLineItem           string `json:"edit_line_item"`
	InventoryItem          string `json:"inventory_item"`
	LineNumber             string `json:"line_number"`
	LineType               string `json:"line_type"`
	Location               string `json:"location"`
	Locked                 string `json:"locked"`
	NoItems                string `json:"no_items"`
	Notes                  string `json:"notes"`
	Product                string `json:"product"`
	QtyOrdered             string `json:"qty_ordered"`
	QuantityBilled         string `json:"quantity_billed"`
	QuantityOrdered        string `json:"quantity_ordered"`
	QuantityReceived       string `json:"quantity_received"`
	RemoveItem             string `json:"remove_item"`
	RemoveLineItem         string `json:"remove_line_item"`
	SelectItem             string `json:"select_item"`
	TotalPrice             string `json:"total_price"`
	TypeExpense            string `json:"type_expense"`
	TypeGoods              string `json:"type_goods"`
	TypeService            string `json:"type_service"`
	UnitPrice              string `json:"unit_price"`
	Type                   string `json:"type"`
	ProductID              string `json:"product_id"`
	InventoryItemID        string `json:"inventory_item_id"`
	LocationID             string `json:"location_id"`
	Save                   string `json:"save"`
	Cancel                 string `json:"cancel"`
}

// ReceiptLabels holds translatable strings for the confirm receipt drawer form.
type ReceiptLabels struct {
	AutoConfirmed     string `json:"auto_confirmed"`
	NoLines           string `json:"no_lines"`
	OverReceiptError  string `json:"over_receipt_error"`
	PartialSuccess    string `json:"partial_success"`
	QtyToReceive      string `json:"qty_to_receive"`
	ReceiptDate       string `json:"receipt_date"`
	ReceivingLocation string `json:"receiving_location"`
	ServiceRendered   string `json:"service_rendered"`
	Success           string `json:"success"`
	Title             string `json:"title"`
	AllReceived       string `json:"all_received"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Ordered           string `json:"ordered"`
	Received          string `json:"received"`
	Remaining         string `json:"remaining"`
	ConfirmButton     string `json:"confirm_button"`
	Cancel            string `json:"cancel"`
}
