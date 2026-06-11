package expenditure

// ---------------------------------------------------------------------------
// Expenditure labels
// ---------------------------------------------------------------------------

// ExpenditureLabels holds all translatable strings for the expenditure module
// (purchase + expense views).
type ExpenditureLabels struct {
	Labels               ExpenditureLabelNames                 `json:"labels"`
	Page                 ExpenditurePageLabels                 `json:"page"`
	Buttons              ExpenditureButtonLabels               `json:"buttons"`
	Columns              ExpenditureColumnLabels               `json:"columns"`
	Empty                ExpenditureEmptyLabels                `json:"empty"`
	Form                 ExpenditureFormLabels                 `json:"form"`
	Status               ExpenditureStatusLabels               `json:"status"`
	Types                ExpenditureTypeLabels                 `json:"types"`
	Actions              ExpenditureActionLabels               `json:"actions"`
	Bulk                 ExpenditureBulkLabels                 `json:"bulkActions"`
	Detail               ExpenditureDetailLabels               `json:"detail"`
	Errors               ExpenditureErrorLabels                `json:"errors"`
	Category             ExpenditureCategoryLabels             `json:"category"`
	PaymentMethod        ExpenditurePaymentMethodLabels        `json:"paymentMethod"`
	DisbursementCategory ExpenditureDisbursementCategoryLabels `json:"disbursementCategory"`
	Schedule             ExpenditureScheduleLabels             `json:"schedule"`
	LineItemForm         ExpenditureLineItemFormLabels         `json:"lineItemForm"`
	DisbursementForm     ExpenditureDisbursementFormLabels     `json:"disbursementForm"`
	PurchaseOrder        PurchaseOrderLabels                   `json:"purchaseOrder"`

	// Dashboard labels — Phase 5. One block per surface (purchase/expense).
	PurchaseDashboard PurchaseDashboardLabels `json:"purchaseDashboard"`
	ExpenseDashboard  ExpenseDashboardLabels  `json:"expenseDashboard"`
}

// PurchaseDashboardLabels holds translatable strings for the purchase
// dashboard (expenditure_type=purchase surface).
type PurchaseDashboardLabels struct {
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	StatOpenPOs       string `json:"statOpenPOs"`
	StatAwaiting      string `json:"statAwaiting"`
	StatSpentMTD      string `json:"statSpentMTD"`
	StatTopSupplier   string `json:"statTopSupplier"`
	WidgetMonthly     string `json:"widgetMonthly"`
	WidgetTopSupplier string `json:"widgetTopSupplier"`
	WidgetRecent      string `json:"widgetRecent"`
	QuickNew          string `json:"quickNew"`
	QuickReceive      string `json:"quickReceive"`
	QuickMatch        string `json:"quickMatch"`
	QuickSuppliers    string `json:"quickSuppliers"`
	ViewAll           string `json:"viewAll"`
	EmptyRecentTitle  string `json:"emptyRecentTitle"`
	EmptyRecentDesc   string `json:"emptyRecentDesc"`
	EmptySuppliers    string `json:"emptySuppliers"`
	NewPurchase       string `json:"newPurchase"`
	ColSupplier       string `json:"colSupplier"`
	ColTotal          string `json:"colTotal"`
}

// ExpenseDashboardLabels holds translatable strings for the expense
// dashboard (expenditure_type=expense surface).
type ExpenseDashboardLabels struct {
	Title                 string `json:"title"`
	Subtitle              string `json:"subtitle"`
	StatPendingApproval   string `json:"statPendingApproval"`
	StatApprovedMTD       string `json:"statApprovedMTD"`
	StatReimbursable      string `json:"statReimbursable"`
	StatCategoriesUsed    string `json:"statCategoriesUsed"`
	WidgetByCategory      string `json:"widgetByCategory"`
	WidgetTopCategory     string `json:"widgetTopCategory"`
	WidgetRecent          string `json:"widgetRecent"`
	QuickNew              string `json:"quickNew"`
	QuickApprove          string `json:"quickApprove"`
	QuickReimburse        string `json:"quickReimburse"`
	QuickCategorySettings string `json:"quickCategorySettings"`
	ViewAll               string `json:"viewAll"`
	EmptyRecentTitle      string `json:"emptyRecentTitle"`
	EmptyRecentDesc       string `json:"emptyRecentDesc"`
	EmptyCategories       string `json:"emptyCategories"`
	NewExpense            string `json:"newExpense"`
	ColCategory           string `json:"colCategory"`
	ColTotal              string `json:"colTotal"`
}

// ExpenditureCategoryLabels holds translatable strings for the expenditure
// category settings list and CRUD drawer.
type ExpenditureCategoryLabels struct {
	Page    ExpenditureCategoryPageLabels    `json:"page"`
	Columns ExpenditureCategoryColumnLabels  `json:"columns"`
	Empty   ExpenditureCategoryEmptyLabels   `json:"empty"`
	Form    ExpenditureCategoryFormLabels    `json:"form"`
	Actions ExpenditureCategoryActionLabels  `json:"actions"`
	Errors  ExpenditureCategoryErrorLabels   `json:"errors"`
	Confirm ExpenditureCategoryConfirmLabels `json:"confirm"`
	Buttons ExpenditureCategoryButtonLabels  `json:"buttons"`
}

type ExpenditureCategoryPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ExpenditureCategoryButtonLabels struct {
	AddCategory string `json:"addCategory"`
}

type ExpenditureCategoryColumnLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ExpenditureCategoryEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ExpenditureCategoryFormLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Field-level info text surfaced via an info button beside each label.
	CodeInfo        string `json:"codeInfo"`
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
}

type ExpenditureCategoryActionLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type ExpenditureCategoryErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	InvalidFormData  string `json:"invalidFormData"`
}

type ExpenditureCategoryConfirmLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
}

// ExpenditureErrorLabels holds error messages for the expenditure action handlers.
type ExpenditureErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	NoPermission     string `json:"noPermission"`
}

type ExpenditureLabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"namePlural"`
	Purchase       string `json:"purchase"`
	PurchasePlural string `json:"purchasePlural"`
	PurchaseOrder  string `json:"purchaseOrder"`
	Expense        string `json:"expense"`
	ExpensePlural  string `json:"expensePlural"`
}

type ExpenditurePageLabels struct {
	PurchaseHeading          string `json:"purchaseHeading"`
	PurchaseCaption          string `json:"purchaseCaption"`
	PurchaseHeadingDraft     string `json:"purchaseHeadingDraft"`
	PurchaseHeadingPending   string `json:"purchaseHeadingPending"`
	PurchaseHeadingApproved  string `json:"purchaseHeadingApproved"`
	PurchaseHeadingPaid      string `json:"purchaseHeadingPaid"`
	PurchaseHeadingCancelled string `json:"purchaseHeadingCancelled"`
	PurchaseHeadingOverdue   string `json:"purchaseHeadingOverdue"`
	ExpenseHeading           string `json:"expenseHeading"`
	ExpenseCaption           string `json:"expenseCaption"`
	ExpenseHeadingDraft      string `json:"expenseHeadingDraft"`
	ExpenseHeadingPending    string `json:"expenseHeadingPending"`
	ExpenseHeadingApproved   string `json:"expenseHeadingApproved"`
	ExpenseHeadingPaid       string `json:"expenseHeadingPaid"`
	ExpenseHeadingCancelled  string `json:"expenseHeadingCancelled"`
	ExpenseHeadingOverdue    string `json:"expenseHeadingOverdue"`
	DashboardPurchase        string `json:"dashboardPurchase"`
	DashboardExpense         string `json:"dashboardExpense"`
}

type ExpenditureButtonLabels struct {
	AddPurchase string `json:"addPurchase"`
	AddExpense  string `json:"addExpense"`
}

type ExpenditureColumnLabels struct {
	Reference string `json:"reference"`
	Vendor    string `json:"vendor"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Category  string `json:"category"`
}

type ExpenditureEmptyLabels struct {
	PurchaseTitle            string `json:"purchaseTitle"`
	PurchaseMessage          string `json:"purchaseMessage"`
	PurchaseDraftTitle       string `json:"purchaseDraftTitle"`
	PurchaseDraftMessage     string `json:"purchaseDraftMessage"`
	PurchasePendingTitle     string `json:"purchasePendingTitle"`
	PurchasePendingMessage   string `json:"purchasePendingMessage"`
	PurchaseApprovedTitle    string `json:"purchaseApprovedTitle"`
	PurchaseApprovedMessage  string `json:"purchaseApprovedMessage"`
	PurchasePaidTitle        string `json:"purchasePaidTitle"`
	PurchasePaidMessage      string `json:"purchasePaidMessage"`
	PurchaseCancelledTitle   string `json:"purchaseCancelledTitle"`
	PurchaseCancelledMessage string `json:"purchaseCancelledMessage"`
	PurchaseOverdueTitle     string `json:"purchaseOverdueTitle"`
	PurchaseOverdueMessage   string `json:"purchaseOverdueMessage"`
	ExpenseTitle             string `json:"expenseTitle"`
	ExpenseMessage           string `json:"expenseMessage"`
	ExpenseDraftTitle        string `json:"expenseDraftTitle"`
	ExpenseDraftMessage      string `json:"expenseDraftMessage"`
	ExpensePendingTitle      string `json:"expensePendingTitle"`
	ExpensePendingMessage    string `json:"expensePendingMessage"`
	ExpenseApprovedTitle     string `json:"expenseApprovedTitle"`
	ExpenseApprovedMessage   string `json:"expenseApprovedMessage"`
	ExpensePaidTitle         string `json:"expensePaidTitle"`
	ExpensePaidMessage       string `json:"expensePaidMessage"`
	ExpenseCancelledTitle    string `json:"expenseCancelledTitle"`
	ExpenseCancelledMessage  string `json:"expenseCancelledMessage"`
	ExpenseOverdueTitle      string `json:"expenseOverdueTitle"`
	ExpenseOverdueMessage    string `json:"expenseOverdueMessage"`
}

type ExpenditureFormLabels struct {
	VendorName                 string `json:"vendorName"`
	VendorNamePlaceholder      string `json:"vendorNamePlaceholder"`
	ExpenditureDate            string `json:"expenditureDate"`
	TotalAmount                string `json:"totalAmount"`
	Currency                   string `json:"currency"`
	Status                     string `json:"status"`
	ReferenceNumber            string `json:"referenceNumber"`
	ReferenceNumberPlaceholder string `json:"referenceNumberPlaceholder"`
	PaymentTerms               string `json:"paymentTerms"`
	DueDate                    string `json:"dueDate"`
	ApprovedBy                 string `json:"approvedBy"`
	ExpenditureType            string `json:"expenditureType"`
	ExpenditureCategory        string `json:"expenditureCategory"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	SectionInfo                string `json:"sectionInfo"`
	SectionVendor              string `json:"sectionVendor"`
	SectionPayment             string `json:"sectionPayment"`
	SectionNotes               string `json:"sectionNotes"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo            string `json:"nameInfo"`
	ExpenditureTypeInfo string `json:"expenditureTypeInfo"`
	CategoryInfo        string `json:"categoryInfo"`
	DateInfo            string `json:"dateInfo"`
	AmountInfo          string `json:"amountInfo"`
	CurrencyInfo        string `json:"currencyInfo"`
	ReferenceNumberInfo string `json:"referenceNumberInfo"`
	SupplierInfo        string `json:"supplierInfo"`
	NotesInfo           string `json:"notesInfo"`
}

type ExpenditureStatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
	Overdue   string `json:"overdue"`
}

type ExpenditureTypeLabels struct {
	Purchase string `json:"purchase"`
	Expense  string `json:"expense"`
	Refund   string `json:"refund"`
	Payroll  string `json:"payroll"`
}

type ExpenditureActionLabels struct {
	Add            string `json:"add"`
	Edit           string `json:"edit"`
	Delete         string `json:"delete"`
	Approve        string `json:"approve"`
	Reject         string `json:"reject"`
	MarkPaid       string `json:"markPaid"`
	ViewPurchase   string `json:"viewPurchase"`
	EditPurchase   string `json:"editPurchase"`
	DeletePurchase string `json:"deletePurchase"`
	ViewExpense    string `json:"viewExpense"`
	EditExpense    string `json:"editExpense"`
	DeleteExpense  string `json:"deleteExpense"`
}

type ExpenditureBulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"markPaid"`
}

type ExpenditureDetailLabels struct {
	PurchasePageTitle    string `json:"purchasePageTitle"`
	ExpensePageTitle     string `json:"expensePageTitle"`
	VendorInfo           string `json:"vendorInfo"`
	VendorName           string `json:"vendorName"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Type                 string `json:"type"`
	Category             string `json:"category"`
	ReferenceNumber      string `json:"referenceNumber"`
	PaymentTerms         string `json:"paymentTerms"`
	DueDate              string `json:"dueDate"`
	ApprovedBy           string `json:"approvedBy"`
	Notes                string `json:"notes"`
	LineItems            string `json:"lineItems"`
	Description          string `json:"description"`
	Quantity             string `json:"quantity"`
	UnitPrice            string `json:"unitPrice"`
	Total                string `json:"total"`
	SubTotal             string `json:"subTotal"`
	GrandTotal           string `json:"grandTotal"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabLineItems         string `json:"tabLineItems"`
	TabPayment           string `json:"tabPayment"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	// Additional fields used in the expense detail template
	Title          string `json:"title"`
	InfoSection    string `json:"infoSection"`
	Name           string `json:"name"`
	PaymentSummary string `json:"paymentSummary"`
	TotalAmount    string `json:"totalAmount"`
	Paid           string `json:"paid"`
	Outstanding    string `json:"outstanding"`
	PaymentStatus  string `json:"paymentStatus"`
	UpdateStatus   string `json:"updateStatus"`
	SaveStatus     string `json:"saveStatus"`
	Payment        string `json:"payment"`
	Pay            string `json:"pay"`
	AddItem        string `json:"addItem"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	TabDetails     string `json:"tabDetails"`
	TabPayments    string `json:"tabPayments"`
	// SPS P10 — Recognition + Accrual tabs on expenditure detail
	TabRecognition          string `json:"tabRecognition"`
	TabAccrual              string `json:"tabAccrual"`
	RecognitionEmptyTitle   string `json:"recognitionEmptyTitle"`
	RecognitionEmptyMessage string `json:"recognitionEmptyMessage"`
	RecognitionRecognizeCTA string `json:"recognitionRecognizeCta"`
	AccrualEmptyTitle       string `json:"accrualEmptyTitle"`
	AccrualEmptyMessage     string `json:"accrualEmptyMessage"`
	TabAttachments          string `json:"tabAttachments"`
}

// ExpenditurePaymentMethodLabels holds translatable strings for disbursement payment methods.
type ExpenditurePaymentMethodLabels struct {
	Cash         string `json:"cash"`
	BankTransfer string `json:"bankTransfer"`
	Check        string `json:"check"`
	GCash        string `json:"gcash"`
	Other        string `json:"other"`
}

// ExpenditureDisbursementCategoryLabels holds translatable strings for disbursement categories.
type ExpenditureDisbursementCategoryLabels struct {
	SupplierPayment string `json:"supplierPayment"`
	Payroll         string `json:"payroll"`
	Rent            string `json:"rent"`
	Utilities       string `json:"utilities"`
	Other           string `json:"other"`
}

// ExpenditureScheduleLabels holds translatable strings for the payment schedule tab.
type ExpenditureScheduleLabels struct {
	Scheduled    string `json:"scheduled"`
	Paid         string `json:"paid"`
	Remaining    string `json:"remaining"`
	DueDate      string `json:"dueDate"`
	AmountDue    string `json:"amountDue"`
	PaidAmount   string `json:"paidAmount"`
	PaidDate     string `json:"paidDate"`
	Reference    string `json:"reference"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ExpenditureLineItemFormLabels holds translatable strings for the line item drawer form.
type ExpenditureLineItemFormLabels struct {
	EditTitle              string `json:"editTitle"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Quantity               string `json:"quantity"`
	UnitPrice              string `json:"unitPrice"`
	Notes                  string `json:"notes"`
	Save                   string `json:"save"`
	Cancel                 string `json:"cancel"`
}

// ExpenditureDisbursementFormLabels holds translatable strings for the pay (disbursement) drawer form.
type ExpenditureDisbursementFormLabels struct {
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	Payee                string `json:"payee"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	CurrencyPlaceholder  string `json:"currencyPlaceholder"`
	PaymentMethod        string `json:"paymentMethod"`
	Category             string `json:"category"`
	ApprovedBy           string `json:"approvedBy"`
	ApproverPlaceholder  string `json:"approverPlaceholder"`
}
