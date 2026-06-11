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
	Bulk      BulkLabels      `json:"bulkActions"`
	Detail    DetailLabels    `json:"detail"`
	Confirm   ConfirmLabels   `json:"confirm"`
	Errors    ErrorLabels     `json:"errors"`
	Dashboard DashboardLabels `json:"dashboard"`
	Settings  SettingsLabels  `json:"settings"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingComplete  string `json:"headingComplete"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionComplete  string `json:"captionComplete"`
	CaptionCancelled string `json:"captionCancelled"`
}

type ButtonLabels struct {
	AddSale string `json:"addSale"`
}

type ColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
	CompleteTitle    string `json:"completeTitle"`
	CompleteMessage  string `json:"completeMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type FormLabels struct {
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	NotesPlaceholder     string `json:"notesPlaceholder"`
	Active               string `json:"active"`
	Location             string `json:"location"`

	// Payment terms and client search labels
	PaymentTerms              string `json:"paymentTerms"`
	SelectPaymentTerm         string `json:"selectPaymentTerm"`
	DueDate                   string `json:"dueDate"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`

	// Subscription search labels
	Subscription          string `json:"subscription"`
	SubscriptionNoResults string `json:"subscriptionNoResults"`

	// Placeholders and translated option labels
	CurrencyPlaceholder            string `json:"currencyPlaceholder"`
	CustomerNamePlaceholder        string `json:"customerNamePlaceholder"`
	StatusDraft                    string `json:"statusDraft"`
	StatusComplete                 string `json:"statusComplete"`
	StatusCancelled                string `json:"statusCancelled"`
	PaymentMethod                  string `json:"paymentMethod"`
	ReferenceNumber                string `json:"referenceNumber"`
	TransactionIdPlaceholder       string `json:"transactionIdPlaceholder"`
	ReceivedBy                     string `json:"receivedBy"`
	Role                           string `json:"role"`
	SelectInventoryItem            string `json:"selectInventoryItem"`
	ItemDescriptionPlaceholder     string `json:"itemDescriptionPlaceholder"`
	DiscountDescriptionPlaceholder string `json:"discountDescriptionPlaceholder"`

	// Field-level info text for the payment drawer form.
	PaymentMethodInfo   string `json:"paymentMethodInfo"`
	AmountInfo          string `json:"amountInfo"`
	CurrencyInfo        string `json:"currencyInfo"`
	ReferenceNumberInfo string `json:"referenceNumberInfo"`
	ReceivedByInfo      string `json:"receivedByInfo"`
	RoleInfo            string `json:"roleInfo"`
	NotesInfo           string `json:"notesInfo"`
}

type ActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Complete          string `json:"complete"`
	Reactivate        string `json:"reactivate"`
	DownloadInvoice   string `json:"downloadInvoice"`
	SendEmail         string `json:"sendEmail"`
	Cancel            string `json:"cancel"`
	ReclassifyToDraft string `json:"reclassifyToDraft"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle   string `json:"pageTitle"`
	TitlePrefix string `json:"titlePrefix"`
	InvoiceInfo string `json:"invoiceInfo"`
	LineItems   string `json:"lineItems"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unitPrice"`
	CostPrice   string `json:"costPrice"`
	GrossProfit string `json:"grossProfit"`
	Total       string `json:"total"`
	Discount    string `json:"discount"`
	SubTotal    string `json:"subTotal"`
	GrandTotal  string `json:"grandTotal"`

	// Tab labels
	TabBasicInfo    string `json:"tabBasicInfo"`
	TabLineItems    string `json:"tabLineItems"`
	TabPayment      string `json:"tabPayment"`
	TabAttachments  string `json:"tabAttachments"`
	TabAuditTrail   string `json:"tabAuditTrail"`
	TabAuditHistory string `json:"tabAuditHistory"`

	// Basic info fields
	Customer     string `json:"customer"`
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	PaymentTerms string `json:"paymentTerms"`
	DueDate      string `json:"dueDate"`

	// Payment fields
	PaymentMethod string `json:"paymentMethod"`
	AmountPaid    string `json:"amountPaid"`
	CardDetails   string `json:"cardDetails"`
	PaymentDate   string `json:"paymentDate"`
	ReceivedBy    string `json:"receivedBy"`
	PaymentInfo   string `json:"paymentInfo"`

	// Audit trail
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`

	// Totals
	TotalGrossProfit string `json:"totalGrossProfit"`

	// Payment empty/table
	Reference           string `json:"reference"`
	PaymentEmptyTitle   string `json:"paymentEmptyTitle"`
	PaymentEmptyMessage string `json:"paymentEmptyMessage"`

	// Line item management
	AddItem                    string `json:"addItem"`
	AddDiscount                string `json:"addDiscount"`
	EditItem                   string `json:"editItem"`
	RemoveItem                 string `json:"removeItem"`
	ItemType                   string `json:"itemType"`
	ItemTypeItem               string `json:"itemTypeItem"`
	ItemTypeDiscount           string `json:"itemTypeDiscount"`
	InventoryItem              string `json:"inventoryItem"`
	SelectInventoryItem        string `json:"selectInventoryItem"`
	ItemDescriptionPlaceholder string `json:"itemDescriptionPlaceholder"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	SerialNumber               string `json:"serialNumber"`
	Product                    string `json:"product"`
	ProductNoResults           string `json:"productNoResults"`
	ProductPlaceholder         string `json:"productPlaceholder"`
	ItemEmptyTitle             string `json:"itemEmptyTitle"`
	ItemEmptyMessage           string `json:"itemEmptyMessage"`

	// Field-level info text for the line-item drawer form.
	ProductInfo     string `json:"productInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	QuantityInfo    string `json:"quantityInfo"`
	UnitPriceInfo   string `json:"unitPriceInfo"`
	CostPriceInfo   string `json:"costPriceInfo"`
	DiscountInfo    string `json:"discountInfo"`
	NotesInfo       string `json:"notesInfo"`

	// Payment tab
	TotalPaid                  string `json:"totalPaid"`
	Remaining                  string `json:"remaining"`
	RecordPayment              string `json:"recordPayment"`
	NoPaymentInfo              string `json:"noPaymentInfo"`
	PaymentDetailsNotAvailable string `json:"paymentDetailsNotAvailable"`
}

type ConfirmLabels struct {
	Complete                 string `json:"complete"`
	CompleteMessage          string `json:"completeMessage"`
	Reactivate               string `json:"reactivate"`
	ReactivateMessage        string `json:"reactivateMessage"`
	BulkComplete             string `json:"bulkComplete"`
	BulkCompleteMessage      string `json:"bulkCompleteMessage"`
	BulkReactivate           string `json:"bulkReactivate"`
	BulkReactivateMessage    string `json:"bulkReactivateMessage"`
	SendEmail                string `json:"sendEmail"`
	SendEmailMessage         string `json:"sendEmailMessage"`
	Cancel                   string `json:"cancel"`
	CancelMessage            string `json:"cancelMessage"`
	ReclassifyToDraft        string `json:"reclassifyToDraft"`
	ReclassifyToDraftMessage string `json:"reclassifyToDraftMessage"`
}

type ErrorLabels struct {
	PermissionDenied        string `json:"permissionDenied"`
	InvalidFormData         string `json:"invalidFormData"`
	NotFound                string `json:"notFound"`
	IDRequired              string `json:"idRequired"`
	NoIDsProvided           string `json:"noIDsProvided"`
	InvalidStatus           string `json:"invalidStatus"`
	InvalidTargetStatus     string `json:"invalidTargetStatus"`
	NoItemsCannotComplete   string `json:"noItemsCannotComplete"`
	HasPaymentsCannotCancel string `json:"hasPaymentsCannotCancel"`
	BulkHasPayments         string `json:"bulkHasPayments"`
	BulkNoItems             string `json:"bulkNoItems"`
	PaymentNotFound         string `json:"paymentNotFound"`
	InvalidDiscount         string `json:"invalidDiscount"`
	// RecomputeUnavailable is the 501 body returned by the RecomputeTaxes stub
	// until Phase 4 wires ComputeTaxesForRevenue (Phase 5 M2).
	RecomputeUnavailable string `json:"recomputeUnavailable"`
}

type DashboardLabels struct {
	Title             string `json:"title"`
	TotalRevenue      string `json:"totalRevenue"`
	Revenue           string `json:"revenue"`
	Completed         string `json:"completed"`
	Active            string `json:"active"`
	RevenueTrend      string `json:"revenueTrend"`
	Week              string `json:"week"`
	Month             string `json:"month"`
	Year              string `json:"year"`
	RecentRevenue     string `json:"recentRevenue"`
	ViewAll           string `json:"viewAll"`
	NewRevenueCreated string `json:"newRevenueCreated"`
	RevenueCompleted  string `json:"revenueCompleted"`
	RevenueUpdated    string `json:"revenueUpdated"`
	RevenueCancelled  string `json:"revenueCancelled"`
	QuickNewRevenue   string `json:"quickNewRevenue"`
	QuickViewAll      string `json:"quickViewAll"`
}

// SettingsLabels holds translatable strings for the revenue settings page
// (invoice template management).
type SettingsLabels struct {
	PageTitle      string `json:"pageTitle"`
	Caption        string `json:"caption"`
	UploadTemplate string `json:"uploadTemplate"`
	TemplateName   string `json:"templateName"`
	TemplateType   string `json:"templateType"`
	Purpose        string `json:"purpose"`
	SetDefault     string `json:"setDefault"`
	Delete         string `json:"delete"`
	DefaultBadge   string `json:"defaultBadge"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	UploadSuccess  string `json:"uploadSuccess"`
	DeleteConfirm  string `json:"deleteConfirm"`
}
