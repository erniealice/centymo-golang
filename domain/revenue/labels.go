package revenue

// ---------------------------------------------------------------------------
// Revenue labels
// ---------------------------------------------------------------------------

// RevenueLabels holds all translatable strings for the revenue module.
type RevenueLabels struct {
	Page      RevenuePageLabels      `json:"page"`
	Buttons   RevenueButtonLabels    `json:"buttons"`
	Columns   RevenueColumnLabels    `json:"columns"`
	Empty     RevenueEmptyLabels     `json:"empty"`
	Form      RevenueFormLabels      `json:"form"`
	Actions   RevenueActionLabels    `json:"actions"`
	Bulk      RevenueBulkLabels      `json:"bulkActions"`
	Detail    RevenueDetailLabels    `json:"detail"`
	Confirm   RevenueConfirmLabels   `json:"confirm"`
	Errors    RevenueErrorLabels     `json:"errors"`
	Dashboard RevenueDashboardLabels `json:"dashboard"`
	Settings  RevenueSettingsLabels  `json:"settings"`
}

type RevenuePageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingComplete  string `json:"headingComplete"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionComplete  string `json:"captionComplete"`
	CaptionCancelled string `json:"captionCancelled"`
}

type RevenueButtonLabels struct {
	AddSale string `json:"addSale"`
}

type RevenueColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
}

type RevenueEmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
	CompleteTitle    string `json:"completeTitle"`
	CompleteMessage  string `json:"completeMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type RevenueFormLabels struct {
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

type RevenueActionLabels struct {
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

type RevenueBulkLabels struct {
	Delete string `json:"delete"`
}

type RevenueDetailLabels struct {
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

type RevenueConfirmLabels struct {
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

type RevenueErrorLabels struct {
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

type RevenueDashboardLabels struct {
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

// RevenueSettingsLabels holds translatable strings for the revenue settings page
// (invoice template management).
type RevenueSettingsLabels struct {
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

// ---------------------------------------------------------------------------
// Revenue Run labels
// ---------------------------------------------------------------------------

// RevenueRunLabels holds all translatable strings for the Revenue Run
// (invoice-run) module. Lyngua root key: "revenueRun".
// D13: naming is revenueRun / revenue_run / RevenueRun / revenue-run everywhere
// except the user-visible VALUE "Invoice Run" (supplied by lyngua).
type RevenueRunLabels struct {
	AppLabel       string                      `json:"appLabel"`
	Queue          RevenueRunQueueLabels       `json:"queue"`
	List           RevenueRunListLabels        `json:"list"`
	Detail         RevenueRunDetailLabels      `json:"detail"`
	StatusBadges   RevenueRunStatusBadgeLabels `json:"statusBadges"`
	Actions        RevenueRunActionLabels      `json:"actions"`
	ScopeKind      RevenueRunScopeKindLabels   `json:"scopeKind"`
	AttemptOutcome RevenueRunOutcomeLabels     `json:"attemptOutcome"`
	Errors         RevenueRunErrorLabels       `json:"errors"`
	// ToastBatchSuccess is the message shown after a Surface B batch-run
	// submission. Supports the standard {{.Created}}/{{.Skipped}}/{{.Errored}}
	// placeholders, substituted Go-side before the toast is dispatched.
	ToastBatchSuccess string `json:"toastBatchSuccess"`
	// ViewRunLink is the link label used on toasts whose batch produced
	// exactly one run. Multi-run batches omit the link.
	ViewRunLink string `json:"viewRunLink"`
}

// RevenueRunQueueLabels holds copy for the workspace-queue page (Surface B).
type RevenueRunQueueLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// AsOfDateLabel is the label for the AsOfDate date picker above the table.
	AsOfDateLabel string                      `json:"asOfDateLabel"`
	Columns       RevenueRunQueueColumnLabels `json:"columns"`
	Empty         RevenueRunQueueEmptyLabels  `json:"empty"`
	Bulk          RevenueRunQueueBulkLabels   `json:"bulk"`
}

type RevenueRunQueueColumnLabels struct {
	Client         string `json:"client"`
	Subscriptions  string `json:"subscriptions"`
	PendingPeriods string `json:"pendingPeriods"`
	Total          string `json:"total"`
	Currency       string `json:"currency"`
	Actions        string `json:"actions"`
	Run            string `json:"run"`
}

type RevenueRunQueueEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RevenueRunQueueBulkLabels struct {
	RunSelected        string `json:"runSelected"`
	RunAllMatching     string `json:"runAllMatching"`
	CapExceededMessage string `json:"capExceededMessage"`
}

// RevenueRunListLabels holds copy for the run history list page (Surface D).
type RevenueRunListLabels struct {
	Title    string                     `json:"title"`
	Subtitle string                     `json:"subtitle"`
	Columns  RevenueRunListColumnLabels `json:"columns"`
	Empty    RevenueRunListEmptyLabels  `json:"empty"`
	Filters  RevenueRunListFilterLabels `json:"filterLabels"`
}

type RevenueRunListColumnLabels struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AsOfDate    string `json:"asOfDate"`
	Initiator   string `json:"initiator"`
	InitiatedAt string `json:"initiatedAt"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Skipped     string `json:"skipped"`
	Errored     string `json:"errored"`
	Actions     string `json:"actions"`
}

type RevenueRunListEmptyLabels struct {
	Pending  RevenueRunListEmptyStateLabels `json:"pending"`
	Complete RevenueRunListEmptyStateLabels `json:"complete"`
	Failed   RevenueRunListEmptyStateLabels `json:"failed"`
}

type RevenueRunListEmptyStateLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RevenueRunListFilterLabels struct {
	Pending  string `json:"pending"`
	Complete string `json:"complete"`
	Failed   string `json:"failed"`
}

// RevenueRunDetailLabels holds copy for the run detail page (Surface D).
type RevenueRunDetailLabels struct {
	Title      string                        `json:"title"`
	Tabs       RevenueRunDetailTabLabels     `json:"tabs"`
	Summary    RevenueRunSummaryLabels       `json:"summary"`
	Selections RevenueRunSelectionsTabLabels `json:"selections"`
	Results    RevenueRunResultsTabLabels    `json:"results"`
	Invoices   RevenueRunInvoicesTabLabels   `json:"invoices"`
}

// RevenueRunSelectionsTabLabels holds column headers and empty-state copy for
// the Selections tab on the run detail page.
type RevenueRunSelectionsTabLabels struct {
	ColSubscription string `json:"colSubscription"`
	ColPeriodStart  string `json:"colPeriodStart"`
	ColPeriodEnd    string `json:"colPeriodEnd"`
	ColPeriodMarker string `json:"colPeriodMarker"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

// RevenueRunResultsTabLabels holds column headers and empty-state copy for
// the Results tab on the run detail page.
type RevenueRunResultsTabLabels struct {
	ColSubscription string `json:"colSubscription"`
	ColPeriodStart  string `json:"colPeriodStart"`
	ColPeriodEnd    string `json:"colPeriodEnd"`
	ColOutcome      string `json:"colOutcome"`
	ColErrorCode    string `json:"colErrorCode"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

// RevenueRunInvoicesTabLabels holds column headers and empty-state copy for
// the Invoices tab on the run detail page. Also holds the coming-soon label
// for the Audit History tab.
type RevenueRunInvoicesTabLabels struct {
	ColReference string `json:"colReference"`
	ColDate      string `json:"colDate"`
	ColAmount    string `json:"colAmount"`
	ColStatus    string `json:"colStatus"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	// AuditHistoryComingSoon is the coming-soon message for the Audit History tab.
	AuditHistoryComingSoon string `json:"auditHistoryComingSoon"`
}

type RevenueRunDetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Invoices     string `json:"invoices"`
	AuditHistory string `json:"auditHistory"`
	Attachments  string `json:"attachments"`
}

type RevenueRunSummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"asOfDate"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiatedAt"`
	CompletedAt             string `json:"completedAt"`
	Status                  string `json:"status"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possiblyInterruptedNote"`
}

// RevenueRunStatusBadgeLabels holds display labels for each run status value.
type RevenueRunStatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possiblyInterrupted"`
}

// RevenueRunActionLabels holds labels for interactive actions on run rows/pages.
type RevenueRunActionLabels struct {
	Run                   string `json:"run"`
	ReRunFailed           string `json:"reRunFailed"`
	ReRunFailedComingSoon string `json:"reRunFailedComingSoon"`
	ViewRun               string `json:"viewRun"`
	ViewClient            string `json:"viewClient"`
	ViewSubscription      string `json:"viewSubscription"`
}

// RevenueRunScopeKindLabels holds display labels for each scope kind value.
type RevenueRunScopeKindLabels struct {
	Subscription string `json:"subscription"`
	Client       string `json:"client"`
	Workspace    string `json:"workspace"`
}

// RevenueRunOutcomeLabels holds display labels for per-attempt outcome values.
type RevenueRunOutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// RevenueRunErrorLabels holds error message strings for the revenue-run module.
type RevenueRunErrorLabels struct {
	CapExceeded         string `json:"capExceeded"`
	PermissionDenied    string `json:"permissionDenied"`
	UseCaseUnavailable  string `json:"useCaseUnavailable"`
	InvalidSelection    string `json:"invalidSelection"`
	IdempotencyConflict string `json:"idempotencyConflict"`
	ClientMismatch      string `json:"clientMismatch"`
	WorkspaceMismatch   string `json:"workspaceMismatch"`
	TamperedPeriod      string `json:"tamperedPeriod"`
	// RunAllMatchingNotImplemented is shown when the operator attempts
	// "run for all matching" before FilterToken signing is wired (Wave 3 stub).
	RunAllMatchingNotImplemented string `json:"runAllMatchingNotImplemented"`
}

// DefaultRevenueRunLabels returns RevenueRunLabels with sensible English defaults.
func DefaultRevenueRunLabels() RevenueRunLabels {
	return RevenueRunLabels{
		AppLabel: "Invoice Run",
		Queue: RevenueRunQueueLabels{
			Title:         "Invoice Run Queue",
			Subtitle:      "Clients with pending billing periods ready to invoice",
			AsOfDateLabel: "As of date",
			Columns: RevenueRunQueueColumnLabels{
				Client:         "Client",
				Subscriptions:  "Subscriptions",
				PendingPeriods: "Pending periods",
				Total:          "Total",
				Currency:       "Currency",
				Actions:        "Actions",
				Run:            "Run",
			},
			Empty: RevenueRunQueueEmptyLabels{
				Title:   "Queue is empty",
				Message: "No clients have pending billing periods at this time.",
			},
			Bulk: RevenueRunQueueBulkLabels{
				RunSelected:        "Run for selected",
				RunAllMatching:     "Run for all matching",
				CapExceededMessage: "Capped at 50 clients per batch. Narrow the filter to run the rest.",
			},
		},
		List: RevenueRunListLabels{
			Title:    "Invoice Runs",
			Subtitle: "History of invoice run batches",
			Columns: RevenueRunListColumnLabels{
				ID:          "Run ID",
				Scope:       "Scope",
				AsOfDate:    "As of date",
				Initiator:   "Initiator",
				InitiatedAt: "Initiated",
				Status:      "Status",
				Created:     "Created",
				Skipped:     "Skipped",
				Errored:     "Errored",
				Actions:     "Actions",
			},
			Empty: RevenueRunListEmptyLabels{
				Pending: RevenueRunListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "There are no invoice runs currently in progress.",
				},
				Complete: RevenueRunListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "No invoice runs have completed yet.",
				},
				Failed: RevenueRunListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "No invoice runs have failed.",
				},
			},
			Filters: RevenueRunListFilterLabels{
				Pending:  "Pending",
				Complete: "Complete",
				Failed:   "Failed",
			},
		},
		Detail: RevenueRunDetailLabels{
			Title: "Invoice Run",
			Tabs: RevenueRunDetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Invoices:     "Invoices",
				AuditHistory: "Audit History",
				Attachments:  "Attachments",
			},
			Summary: RevenueRunSummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some invoices may be missing.",
			},
			Selections: RevenueRunSelectionsTabLabels{
				ColSubscription: "Subscription",
				ColPeriodStart:  "Period start",
				ColPeriodEnd:    "Period end",
				ColPeriodMarker: "Period marker",
				EmptyTitle:      "No selections",
				EmptyMessage:    "This run has no attempt records.",
			},
			Results: RevenueRunResultsTabLabels{
				ColSubscription: "Subscription",
				ColPeriodStart:  "Period start",
				ColPeriodEnd:    "Period end",
				ColOutcome:      "Outcome",
				ColErrorCode:    "Error code",
				EmptyTitle:      "No results",
				EmptyMessage:    "This run has no attempt records.",
			},
			Invoices: RevenueRunInvoicesTabLabels{
				ColReference:           "Reference",
				ColDate:                "Date",
				ColAmount:              "Amount",
				ColStatus:              "Status",
				EmptyTitle:             "No invoices",
				EmptyMessage:           "No invoices were created by this run.",
				AuditHistoryComingSoon: "Audit history is coming soon.",
			},
		},
		StatusBadges: RevenueRunStatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		Actions: RevenueRunActionLabels{
			Run:                   "Run",
			ReRunFailed:           "Re-run failed",
			ReRunFailedComingSoon: "Re-run failed (coming soon)",
			ViewRun:               "View run",
			ViewClient:            "View client",
			ViewSubscription:      "View subscription",
		},
		ScopeKind: RevenueRunScopeKindLabels{
			Subscription: "Subscription",
			Client:       "Client",
			Workspace:    "Workspace",
		},
		AttemptOutcome: RevenueRunOutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Errors: RevenueRunErrorLabels{
			CapExceeded:                  "Batch cap exceeded — maximum 50 clients per run.",
			PermissionDenied:             "You do not have permission to run invoices.",
			UseCaseUnavailable:           "Invoice run is not available for this subscription type.",
			InvalidSelection:             "One or more selected subscriptions are invalid.",
			IdempotencyConflict:          "An invoice for one or more periods already exists.",
			ClientMismatch:               "Selected subscriptions belong to different clients.",
			WorkspaceMismatch:            "Selected subscriptions belong to a different workspace.",
			TamperedPeriod:               "A billing period was modified after selection. Please retry.",
			RunAllMatchingNotImplemented: "Run for all matching is not yet available. Please select individual clients.",
		},
		ToastBatchSuccess: "Invoice batch run — {{.Created}} created, {{.Skipped}} skipped, {{.Errored}} failed.",
		ViewRunLink:       "View run",
	}
}
