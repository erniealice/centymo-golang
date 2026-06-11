package purchaseorder

// ---------------------------------------------------------------------------
// Purchase Order labels
// ---------------------------------------------------------------------------

// ErrorLabels holds error messages for the purchase order action handlers.
type ErrorLabels struct {
	NoPermission string `json:"noPermission"`
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
	POTypes   POTypeLabels   `json:"poTypes"`
	LineTypes LineTypeLabels `json:"lineTypes"`
	Actions   ActionLabels   `json:"actions"`
	Bulk      BulkLabels     `json:"bulkActions"`
	Detail    DetailLabels   `json:"detail"`
	LineItems LineItemLabels `json:"lineItems"`
	Receipt   ReceiptLabels  `json:"receipt"`
	Errors    ErrorLabels    `json:"errors"`
}

type LabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"namePlural"`
	LineItem       string `json:"lineItem"`
	LineItemPlural string `json:"lineItemPlural"`
}

type PageLabels struct {
	Heading                  string `json:"heading"`
	Caption                  string `json:"caption"`
	HeadingDraft             string `json:"headingDraft"`
	HeadingPendingApproval   string `json:"headingPendingApproval"`
	HeadingApproved          string `json:"headingApproved"`
	HeadingPartiallyReceived string `json:"headingPartiallyReceived"`
	HeadingFullyReceived     string `json:"headingFullyReceived"`
	HeadingBilled            string `json:"headingBilled"`
	HeadingClosed            string `json:"headingClosed"`
	HeadingCancelled         string `json:"headingCancelled"`
	Dashboard                string `json:"dashboard"`
}

type ButtonLabels struct {
	Add         string `json:"add"`
	AddLineItem string `json:"addLineItem"`
}

type ColumnLabels struct {
	PONumber        string `json:"poNumber"`
	POType          string `json:"poType"`
	Supplier        string `json:"supplier"`
	Location        string `json:"location"`
	OrderDate       string `json:"orderDate"`
	Status          string `json:"status"`
	Currency        string `json:"currency"`
	Subtotal        string `json:"subtotal"`
	TaxAmount       string `json:"taxAmount"`
	TotalAmount     string `json:"totalAmount"`
	PaymentTerms    string `json:"paymentTerms"`
	ShippingTerms   string `json:"shippingTerms"`
	ApprovedBy      string `json:"approvedBy"`
	ReferenceNumber string `json:"referenceNumber"`
	Notes           string `json:"notes"`
}

type EmptyLabels struct {
	Title                    string `json:"title"`
	Message                  string `json:"message"`
	DraftTitle               string `json:"draftTitle"`
	DraftMessage             string `json:"draftMessage"`
	PendingApprovalTitle     string `json:"pendingApprovalTitle"`
	PendingApprovalMessage   string `json:"pendingApprovalMessage"`
	ApprovedTitle            string `json:"approvedTitle"`
	ApprovedMessage          string `json:"approvedMessage"`
	PartiallyReceivedTitle   string `json:"partiallyReceivedTitle"`
	PartiallyReceivedMessage string `json:"partiallyReceivedMessage"`
	FullyReceivedTitle       string `json:"fullyReceivedTitle"`
	FullyReceivedMessage     string `json:"fullyReceivedMessage"`
	BilledTitle              string `json:"billedTitle"`
	BilledMessage            string `json:"billedMessage"`
	ClosedTitle              string `json:"closedTitle"`
	ClosedMessage            string `json:"closedMessage"`
	CancelledTitle           string `json:"cancelledTitle"`
	CancelledMessage         string `json:"cancelledMessage"`
}

type FormLabels struct {
	PONumber                   string `json:"poNumber"`
	PONumberPlaceholder        string `json:"poNumberPlaceholder"`
	POType                     string `json:"poType"`
	SelectPOType               string `json:"selectPoType"`
	Supplier                   string `json:"supplier"`
	SelectSupplier             string `json:"selectSupplier"`
	Location                   string `json:"location"`
	SelectLocation             string `json:"selectLocation"`
	OrderDate                  string `json:"orderDate"`
	Currency                   string `json:"currency"`
	Subtotal                   string `json:"subtotal"`
	TaxAmount                  string `json:"taxAmount"`
	TotalAmount                string `json:"totalAmount"`
	PaymentTerms               string `json:"paymentTerms"`
	ShippingTerms              string `json:"shippingTerms"`
	ApprovedBy                 string `json:"approvedBy"`
	ReferenceNumber            string `json:"referenceNumber"`
	ReferenceNumberPlaceholder string `json:"referenceNumberPlaceholder"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	SectionInfo                string `json:"sectionInfo"`
	SectionSupplier            string `json:"sectionSupplier"`
	SectionFinancials          string `json:"sectionFinancials"`
	SectionNotes               string `json:"sectionNotes"`

	// Field-level info text surfaced via an info button beside each label.
	PONumberInfo         string `json:"poNumberInfo"`
	POTypeInfo           string `json:"poTypeInfo"`
	SupplierInfo         string `json:"supplierInfo"`
	OrderDateInfo        string `json:"orderDateInfo"`
	ExpectedDeliveryInfo string `json:"expectedDeliveryInfo"`
	CurrencyInfo         string `json:"currencyInfo"`
	PaymentTermsInfo     string `json:"paymentTermsInfo"`
	ShippingTermsInfo    string `json:"shippingTermsInfo"`
	ReferenceNumberInfo  string `json:"referenceNumberInfo"`
	NotesInfo            string `json:"notesInfo"`
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
	ConfirmReceipt string `json:"confirmReceipt"`
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
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"supplierInfo"`
	Supplier             string `json:"supplier"`
	Location             string `json:"location"`
	OrderDate            string `json:"orderDate"`
	PONumber             string `json:"poNumber"`
	POType               string `json:"poType"`
	Status               string `json:"status"`
	Currency             string `json:"currency"`
	Subtotal             string `json:"subtotal"`
	TaxAmount            string `json:"taxAmount"`
	TotalAmount          string `json:"totalAmount"`
	PaymentTerms         string `json:"paymentTerms"`
	ShippingTerms        string `json:"shippingTerms"`
	ApprovedBy           string `json:"approvedBy"`
	ReferenceNumber      string `json:"referenceNumber"`
	Notes                string `json:"notes"`
	LineItems            string `json:"lineItems"`
	Description          string `json:"description"`
	LineType             string `json:"lineType"`
	LineNumber           string `json:"lineNumber"`
	QuantityOrdered      string `json:"quantityOrdered"`
	QuantityReceived     string `json:"quantityReceived"`
	QuantityBilled       string `json:"quantityBilled"`
	UnitPrice            string `json:"unitPrice"`
	TotalPrice           string `json:"totalPrice"`
	SubTotal             string `json:"subTotal"`
	GrandTotal           string `json:"grandTotal"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabLineItems         string `json:"tabLineItems"`
	TabReceiving         string `json:"tabReceiving"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	Total                string `json:"total"`
	AddLineItem          string `json:"addLineItem"`
	NoLineItems          string `json:"noLineItems"`
	ConfirmReceiptBtn    string `json:"confirmReceiptBtn"`
	TabAttachments       string `json:"tabAttachments"`
}

// LineItemLabels holds translatable strings for the PO line item drawer form.
type LineItemLabels struct {
	AddItem                string `json:"addItem"`
	AddLineItem            string `json:"addLineItem"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	EditItem               string `json:"editItem"`
	EditLineItem           string `json:"editLineItem"`
	InventoryItem          string `json:"inventoryItem"`
	LineNumber             string `json:"lineNumber"`
	LineType               string `json:"lineType"`
	Location               string `json:"location"`
	Locked                 string `json:"locked"`
	NoItems                string `json:"noItems"`
	Notes                  string `json:"notes"`
	Product                string `json:"product"`
	QtyOrdered             string `json:"qtyOrdered"`
	QuantityBilled         string `json:"quantityBilled"`
	QuantityOrdered        string `json:"quantityOrdered"`
	QuantityReceived       string `json:"quantityReceived"`
	RemoveItem             string `json:"removeItem"`
	RemoveLineItem         string `json:"removeLineItem"`
	SelectItem             string `json:"selectItem"`
	TotalPrice             string `json:"totalPrice"`
	TypeExpense            string `json:"typeExpense"`
	TypeGoods              string `json:"typeGoods"`
	TypeService            string `json:"typeService"`
	UnitPrice              string `json:"unitPrice"`
	Type                   string `json:"type"`
	ProductID              string `json:"productId"`
	InventoryItemID        string `json:"inventoryItemId"`
	LocationID             string `json:"locationId"`
	Save                   string `json:"save"`
	Cancel                 string `json:"cancel"`
}

// ReceiptLabels holds translatable strings for the confirm receipt drawer form.
type ReceiptLabels struct {
	AutoConfirmed     string `json:"autoConfirmed"`
	NoLines           string `json:"noLines"`
	OverReceiptError  string `json:"overReceiptError"`
	PartialSuccess    string `json:"partialSuccess"`
	QtyToReceive      string `json:"qtyToReceive"`
	ReceiptDate       string `json:"receiptDate"`
	ReceivingLocation string `json:"receivingLocation"`
	ServiceRendered   string `json:"serviceRendered"`
	Success           string `json:"success"`
	Title             string `json:"title"`
	AllReceived       string `json:"allReceived"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Ordered           string `json:"ordered"`
	Received          string `json:"received"`
	Remaining         string `json:"remaining"`
	ConfirmButton     string `json:"confirmButton"`
	Cancel            string `json:"cancel"`
}
