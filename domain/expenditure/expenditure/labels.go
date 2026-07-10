package expenditure

import purchaseorder "github.com/erniealice/centymo-golang/domain/expenditure/purchase_order"

// ---------------------------------------------------------------------------
// Expenditure labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the expenditure module
// (purchase + expense views).
type Labels struct {
	Labels               LabelNames                 `json:"labels"`
	Page                 PageLabels                 `json:"page"`
	Buttons              ButtonLabels               `json:"buttons"`
	Columns              ColumnLabels               `json:"columns"`
	Empty                EmptyLabels                `json:"empty"`
	Form                 FormLabels                 `json:"form"`
	Status               StatusLabels               `json:"status"`
	Types                TypeLabels                 `json:"types"`
	Actions              ActionLabels               `json:"actions"`
	Bulk                 BulkLabels                 `json:"bulk_actions"`
	Detail               DetailLabels               `json:"detail"`
	Errors               ErrorLabels                `json:"errors"`
	Category             CategoryLabels             `json:"category"`
	PaymentMethod        PaymentMethodLabels        `json:"payment_method"`
	DisbursementCategory DisbursementCategoryLabels `json:"disbursement_category"`
	Schedule             ScheduleLabels             `json:"schedule"`
	LineItemForm         LineItemFormLabels         `json:"line_item_form"`
	DisbursementForm     DisbursementFormLabels     `json:"disbursement_form"`
	PurchaseOrder        purchaseorder.Labels       `json:"purchase_order"`

	// Dashboard labels — Phase 5. One block per surface (purchase/expense).
	PurchaseDashboard PurchaseDashboardLabels `json:"purchase_dashboard"`
	ExpenseDashboard  ExpenseDashboardLabels  `json:"expense_dashboard"`
}

// PurchaseDashboardLabels holds translatable strings for the purchase
// dashboard (expenditure_type=purchase surface).
type PurchaseDashboardLabels struct {
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	StatOpenPOs       string `json:"stat_open_pos"`
	StatAwaiting      string `json:"stat_awaiting"`
	StatSpentMTD      string `json:"stat_spent_mtd"`
	StatTopSupplier   string `json:"stat_top_supplier"`
	WidgetMonthly     string `json:"widget_monthly"`
	WidgetTopSupplier string `json:"widget_top_supplier"`
	WidgetRecent      string `json:"widget_recent"`
	QuickNew          string `json:"quick_new"`
	QuickReceive      string `json:"quick_receive"`
	QuickMatch        string `json:"quick_match"`
	QuickSuppliers    string `json:"quick_suppliers"`
	ViewAll           string `json:"view_all"`
	EmptyRecentTitle  string `json:"empty_recent_title"`
	EmptyRecentDesc   string `json:"empty_recent_desc"`
	EmptySuppliers    string `json:"empty_suppliers"`
	NewPurchase       string `json:"new_purchase"`
	ColSupplier       string `json:"col_supplier"`
	ColTotal          string `json:"col_total"`
}

// ExpenseDashboardLabels holds translatable strings for the expense
// dashboard (expenditure_type=expense surface).
type ExpenseDashboardLabels struct {
	Title                 string `json:"title"`
	Subtitle              string `json:"subtitle"`
	StatPendingApproval   string `json:"stat_pending_approval"`
	StatApprovedMTD       string `json:"stat_approved_mtd"`
	StatReimbursable      string `json:"stat_reimbursable"`
	StatCategoriesUsed    string `json:"stat_categories_used"`
	WidgetByCategory      string `json:"widget_by_category"`
	WidgetTopCategory     string `json:"widget_top_category"`
	WidgetRecent          string `json:"widget_recent"`
	QuickNew              string `json:"quick_new"`
	QuickApprove          string `json:"quick_approve"`
	QuickReimburse        string `json:"quick_reimburse"`
	QuickCategorySettings string `json:"quick_category_settings"`
	ViewAll               string `json:"view_all"`
	EmptyRecentTitle      string `json:"empty_recent_title"`
	EmptyRecentDesc       string `json:"empty_recent_desc"`
	EmptyCategories       string `json:"empty_categories"`
	NewExpense            string `json:"new_expense"`
	ColCategory           string `json:"col_category"`
	ColTotal              string `json:"col_total"`
}

// CategoryLabels holds translatable strings for the expenditure
// category settings list and CRUD drawer.
type CategoryLabels struct {
	Page    CategoryPageLabels    `json:"page"`
	Columns CategoryColumnLabels  `json:"columns"`
	Empty   CategoryEmptyLabels   `json:"empty"`
	Form    CategoryFormLabels    `json:"form"`
	Actions CategoryActionLabels  `json:"actions"`
	Errors  CategoryErrorLabels   `json:"errors"`
	Confirm CategoryConfirmLabels `json:"confirm"`
	Buttons CategoryButtonLabels  `json:"buttons"`
}

type CategoryPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type CategoryButtonLabels struct {
	AddCategory string `json:"add_category"`
}

type CategoryColumnLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CategoryEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CategoryFormLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Field-level info text surfaced via an info button beside each label.
	CodeInfo        string `json:"code_info"`
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
}

type CategoryActionLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type CategoryErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	InvalidFormData  string `json:"invalid_form_data"`
}

type CategoryConfirmLabels struct {
	DeleteTitle   string `json:"delete_title"`
	DeleteMessage string `json:"delete_message"`
}

// ErrorLabels holds error messages for the expenditure action handlers.
type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
	NoPermission     string `json:"no_permission"`
}

type LabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"name_plural"`
	Purchase       string `json:"purchase"`
	PurchasePlural string `json:"purchase_plural"`
	PurchaseOrder  string `json:"purchase_order"`
	Expense        string `json:"expense"`
	ExpensePlural  string `json:"expense_plural"`
}

type PageLabels struct {
	PurchaseHeading          string `json:"purchase_heading"`
	PurchaseCaption          string `json:"purchase_caption"`
	PurchaseHeadingDraft     string `json:"purchase_heading_draft"`
	PurchaseHeadingPending   string `json:"purchase_heading_pending"`
	PurchaseHeadingApproved  string `json:"purchase_heading_approved"`
	PurchaseHeadingPaid      string `json:"purchase_heading_paid"`
	PurchaseHeadingCancelled string `json:"purchase_heading_cancelled"`
	PurchaseHeadingOverdue   string `json:"purchase_heading_overdue"`
	ExpenseHeading           string `json:"expense_heading"`
	ExpenseCaption           string `json:"expense_caption"`
	ExpenseHeadingDraft      string `json:"expense_heading_draft"`
	ExpenseHeadingPending    string `json:"expense_heading_pending"`
	ExpenseHeadingApproved   string `json:"expense_heading_approved"`
	ExpenseHeadingPaid       string `json:"expense_heading_paid"`
	ExpenseHeadingCancelled  string `json:"expense_heading_cancelled"`
	ExpenseHeadingOverdue    string `json:"expense_heading_overdue"`
	DashboardPurchase        string `json:"dashboard_purchase"`
	DashboardExpense         string `json:"dashboard_expense"`
}

type ButtonLabels struct {
	AddPurchase string `json:"add_purchase"`
	AddExpense  string `json:"add_expense"`
}

type ColumnLabels struct {
	Reference string `json:"reference"`
	Vendor    string `json:"vendor"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Category  string `json:"category"`
}

type EmptyLabels struct {
	PurchaseTitle            string `json:"purchase_title"`
	PurchaseMessage          string `json:"purchase_message"`
	PurchaseDraftTitle       string `json:"purchase_draft_title"`
	PurchaseDraftMessage     string `json:"purchase_draft_message"`
	PurchasePendingTitle     string `json:"purchase_pending_title"`
	PurchasePendingMessage   string `json:"purchase_pending_message"`
	PurchaseApprovedTitle    string `json:"purchase_approved_title"`
	PurchaseApprovedMessage  string `json:"purchase_approved_message"`
	PurchasePaidTitle        string `json:"purchase_paid_title"`
	PurchasePaidMessage      string `json:"purchase_paid_message"`
	PurchaseCancelledTitle   string `json:"purchase_cancelled_title"`
	PurchaseCancelledMessage string `json:"purchase_cancelled_message"`
	PurchaseOverdueTitle     string `json:"purchase_overdue_title"`
	PurchaseOverdueMessage   string `json:"purchase_overdue_message"`
	ExpenseTitle             string `json:"expense_title"`
	ExpenseMessage           string `json:"expense_message"`
	ExpenseDraftTitle        string `json:"expense_draft_title"`
	ExpenseDraftMessage      string `json:"expense_draft_message"`
	ExpensePendingTitle      string `json:"expense_pending_title"`
	ExpensePendingMessage    string `json:"expense_pending_message"`
	ExpenseApprovedTitle     string `json:"expense_approved_title"`
	ExpenseApprovedMessage   string `json:"expense_approved_message"`
	ExpensePaidTitle         string `json:"expense_paid_title"`
	ExpensePaidMessage       string `json:"expense_paid_message"`
	ExpenseCancelledTitle    string `json:"expense_cancelled_title"`
	ExpenseCancelledMessage  string `json:"expense_cancelled_message"`
	ExpenseOverdueTitle      string `json:"expense_overdue_title"`
	ExpenseOverdueMessage    string `json:"expense_overdue_message"`
}

type FormLabels struct {
	VendorName                 string `json:"vendor_name"`
	VendorNamePlaceholder      string `json:"vendor_name_placeholder"`
	ExpenditureDate            string `json:"expenditure_date"`
	TotalAmount                string `json:"total_amount"`
	Currency                   string `json:"currency"`
	Status                     string `json:"status"`
	ReferenceNumber            string `json:"reference_number"`
	ReferenceNumberPlaceholder string `json:"reference_number_placeholder"`
	PaymentTerms               string `json:"payment_terms"`
	DueDate                    string `json:"due_date"`
	ApprovedBy                 string `json:"approved_by"`
	ExpenditureType            string `json:"expenditure_type"`
	ExpenditureCategory        string `json:"expenditure_category"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notes_placeholder"`
	SectionInfo                string `json:"section_info"`
	SectionVendor              string `json:"section_vendor"`
	SectionPayment             string `json:"section_payment"`
	SectionNotes               string `json:"section_notes"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo            string `json:"name_info"`
	ExpenditureTypeInfo string `json:"expenditure_type_info"`
	CategoryInfo        string `json:"category_info"`
	DateInfo            string `json:"date_info"`
	AmountInfo          string `json:"amount_info"`
	CurrencyInfo        string `json:"currency_info"`
	ReferenceNumberInfo string `json:"reference_number_info"`
	SupplierInfo        string `json:"supplier_info"`
	NotesInfo           string `json:"notes_info"`
}

type StatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
	Overdue   string `json:"overdue"`
}

type TypeLabels struct {
	Purchase string `json:"purchase"`
	Expense  string `json:"expense"`
	Refund   string `json:"refund"`
	Payroll  string `json:"payroll"`
}

type ActionLabels struct {
	Add            string `json:"add"`
	Edit           string `json:"edit"`
	Delete         string `json:"delete"`
	Approve        string `json:"approve"`
	Reject         string `json:"reject"`
	MarkPaid       string `json:"mark_paid"`
	ViewPurchase   string `json:"view_purchase"`
	EditPurchase   string `json:"edit_purchase"`
	DeletePurchase string `json:"delete_purchase"`
	ViewExpense    string `json:"view_expense"`
	EditExpense    string `json:"edit_expense"`
	DeleteExpense  string `json:"delete_expense"`
}

type BulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"mark_paid"`
}

type DetailLabels struct {
	PurchasePageTitle    string `json:"purchase_page_title"`
	ExpensePageTitle     string `json:"expense_page_title"`
	VendorInfo           string `json:"vendor_info"`
	VendorName           string `json:"vendor_name"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Type                 string `json:"type"`
	Category             string `json:"category"`
	ReferenceNumber      string `json:"reference_number"`
	PaymentTerms         string `json:"payment_terms"`
	DueDate              string `json:"due_date"`
	ApprovedBy           string `json:"approved_by"`
	Notes                string `json:"notes"`
	LineItems            string `json:"line_items"`
	Description          string `json:"description"`
	Quantity             string `json:"quantity"`
	UnitPrice            string `json:"unit_price"`
	Total                string `json:"total"`
	SubTotal             string `json:"sub_total"`
	GrandTotal           string `json:"grand_total"`
	TabBasicInfo         string `json:"tab_basic_info"`
	TabLineItems         string `json:"tab_line_items"`
	TabPayment           string `json:"tab_payment"`
	TabAuditTrail        string `json:"tab_audit_trail"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditAction          string `json:"audit_action"`
	AuditUser            string `json:"audit_user"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	// Additional fields used in the expense detail template
	Title          string `json:"title"`
	InfoSection    string `json:"info_section"`
	Name           string `json:"name"`
	PaymentSummary string `json:"payment_summary"`
	TotalAmount    string `json:"total_amount"`
	Paid           string `json:"paid"`
	Outstanding    string `json:"outstanding"`
	PaymentStatus  string `json:"payment_status"`
	UpdateStatus   string `json:"update_status"`
	SaveStatus     string `json:"save_status"`
	Payment        string `json:"payment"`
	Pay            string `json:"pay"`
	AddItem        string `json:"add_item"`
	EmptyTitle     string `json:"empty_title"`
	EmptyMessage   string `json:"empty_message"`
	TabDetails     string `json:"tab_details"`
	TabPayments    string `json:"tab_payments"`
	// SPS P10 — Recognition + Accrual tabs on expenditure detail
	TabRecognition          string `json:"tab_recognition"`
	TabAccrual              string `json:"tab_accrual"`
	RecognitionEmptyTitle   string `json:"recognition_empty_title"`
	RecognitionEmptyMessage string `json:"recognition_empty_message"`
	RecognitionRecognizeCTA string `json:"recognition_recognize_cta"`
	AccrualEmptyTitle       string `json:"accrual_empty_title"`
	AccrualEmptyMessage     string `json:"accrual_empty_message"`
	TabAttachments          string `json:"tab_attachments"`
}

// PaymentMethodLabels holds translatable strings for disbursement payment methods.
type PaymentMethodLabels struct {
	Cash         string `json:"cash"`
	BankTransfer string `json:"bank_transfer"`
	Check        string `json:"check"`
	GCash        string `json:"gcash"`
	Other        string `json:"other"`
}

// DisbursementCategoryLabels holds translatable strings for disbursement categories.
type DisbursementCategoryLabels struct {
	SupplierPayment string `json:"supplier_payment"`
	Payroll         string `json:"payroll"`
	Rent            string `json:"rent"`
	Utilities       string `json:"utilities"`
	Other           string `json:"other"`
}

// ScheduleLabels holds translatable strings for the payment schedule tab.
type ScheduleLabels struct {
	Scheduled    string `json:"scheduled"`
	Paid         string `json:"paid"`
	Remaining    string `json:"remaining"`
	DueDate      string `json:"due_date"`
	AmountDue    string `json:"amount_due"`
	PaidAmount   string `json:"paid_amount"`
	PaidDate     string `json:"paid_date"`
	Reference    string `json:"reference"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
}

// LineItemFormLabels holds translatable strings for the line item drawer form.
type LineItemFormLabels struct {
	EditTitle              string `json:"edit_title"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	Quantity               string `json:"quantity"`
	UnitPrice              string `json:"unit_price"`
	Notes                  string `json:"notes"`
	Save                   string `json:"save"`
	Cancel                 string `json:"cancel"`
}

// DisbursementFormLabels holds translatable strings for the pay (disbursement) drawer form.
type DisbursementFormLabels struct {
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"reference_placeholder"`
	Payee                string `json:"payee"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	CurrencyPlaceholder  string `json:"currency_placeholder"`
	PaymentMethod        string `json:"payment_method"`
	Category             string `json:"category"`
	ApprovedBy           string `json:"approved_by"`
	ApproverPlaceholder  string `json:"approver_placeholder"`
}

// DefaultLabels returns the zero-value label set. Every rendered string for
// this entity must come from the lyngua cascade (general -> business-type
// tier); there are no Go-side default strings to fall back on.
func DefaultLabels() Labels { return Labels{} }
