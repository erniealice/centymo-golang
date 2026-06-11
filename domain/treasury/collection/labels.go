package collection

// labels.go — collection-entity label structs (centymo W5).
//
// Collection (money IN) and the Cash dashboard labels, extracted verbatim from
// the treasury-domain labels.go into the per-entity collection package per the
// domain-first restructure. Pure structural move — no behaviour change. Lyngua
// JSON load paths are unchanged.

// ---------------------------------------------------------------------------
// Collection labels (money IN — payment collections, receivables)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the collection module.
type Labels struct {
	Page      PageLabels          `json:"page"`
	Buttons   ButtonLabels        `json:"buttons"`
	Columns   ColumnLabels        `json:"columns"`
	Empty     EmptyLabels         `json:"empty"`
	Form      FormLabels          `json:"form"`
	Actions   ActionLabels        `json:"actions"`
	Bulk      BulkLabels          `json:"bulkActions"`
	Detail    DetailLabels        `json:"detail"`
	Status    StatusLabels        `json:"status"`
	Confirm   ConfirmLabels       `json:"confirm"`
	Errors    ErrorLabels         `json:"errors"`
	Dashboard CashDashboardLabels `json:"dashboard"`
}

// CashDashboardLabels holds translatable strings for the cash (collection)
// dashboard page. The "Cash" wording is preferred at the dashboard surface
// because the sidebar key is "cash"; underlying entity is still Collection.
type CashDashboardLabels struct {
	Title              string `json:"title"`
	Subtitle           string `json:"subtitle"`
	StatPending        string `json:"statPending"`
	StatOverdue        string `json:"statOverdue"`
	StatCollectedToday string `json:"statCollectedToday"`
	StatCollectedWeek  string `json:"statCollectedWeek"`
	WidgetDailyTrend   string `json:"widgetDailyTrend"`
	WidgetByMode       string `json:"widgetByMode"`
	WidgetRecent       string `json:"widgetRecent"`
	QuickRecord        string `json:"quickRecord"`
	QuickReconcile     string `json:"quickReconcile"`
	QuickAging         string `json:"quickAging"`
	QuickMarkCleared   string `json:"quickMarkCleared"`
	ViewAll            string `json:"viewAll"`
	EmptyRecentTitle   string `json:"emptyRecentTitle"`
	EmptyRecentDesc    string `json:"emptyRecentDesc"`
	NewCollection      string `json:"newCollection"`
	CollectionUpdated  string `json:"collectionUpdated"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingPending   string `json:"headingPending"`
	HeadingCompleted string `json:"headingCompleted"`
	HeadingFailed    string `json:"headingFailed"`
	Caption          string `json:"caption"`
	CaptionPending   string `json:"captionPending"`
	CaptionCompleted string `json:"captionCompleted"`
	CaptionFailed    string `json:"captionFailed"`
	Dashboard        string `json:"dashboard"`
}

type ButtonLabels struct {
	AddCollection string `json:"addCollection"`
}

type ColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
}

type EmptyLabels struct {
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	CompletedTitle   string `json:"completedTitle"`
	CompletedMessage string `json:"completedMessage"`
	FailedTitle      string `json:"failedTitle"`
	FailedMessage    string `json:"failedMessage"`
}

type FormLabels struct {
	Customer                string `json:"customer"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"referencePlaceholder"`
	PaymentMethod           string `json:"paymentMethod"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
	CustomerNamePlaceholder string `json:"customerNamePlaceholder"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	MethodCash              string `json:"methodCash"`
	MethodBankTransfer      string `json:"methodBankTransfer"`
	MethodCheck             string `json:"methodCheck"`
	MethodGCash             string `json:"methodGCash"`
	MethodMaya              string `json:"methodMaya"`
	MethodCard              string `json:"methodCard"`
	MethodOther             string `json:"methodOther"`
	StatusPending           string `json:"statusPending"`
	StatusCompleted         string `json:"statusCompleted"`
	StatusFailed            string `json:"statusFailed"`

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo     string `json:"referenceInfo"`
	CustomerInfo      string `json:"customerInfo"`
	AmountInfo        string `json:"amountInfo"`
	CurrencyInfo      string `json:"currencyInfo"`
	PaymentMethodInfo string `json:"paymentMethodInfo"`
	DateInfo          string `json:"dateInfo"`
	StatusInfo        string `json:"statusInfo"`
	NotesInfo         string `json:"notesInfo"`

	// 20260517-advance-cash-events Plan B Phase 4 — advance metadata fields
	// rendered conditionally in the collection drawer form. The enum option
	// labels (None / Time-based / Milestone / Unscheduled / Full tranche /
	// Day-prorated / Next period start) are sourced from AdvanceEnumLabels
	// (loaded from advance_kind.json) rather than duplicated here.
	AdvanceMetadata        string `json:"advanceMetadata"`
	AdvanceKind            string `json:"advanceKind"`
	AdvanceProrationPolicy string `json:"advanceProrationPolicy"`
}

type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"markComplete"`
	Reactivate   string `json:"reactivate"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	TitlePrefix          string `json:"titlePrefix"`
	PaymentInfo          string `json:"paymentInfo"`
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Method               string `json:"method"`
	Reference            string `json:"reference"`
	Notes                string `json:"notes"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabAttachments       string `json:"tabAttachments"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	TabAuditHistory      string `json:"tabAuditHistory"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditTrailDesc       string `json:"auditTrailDesc"`
}

type StatusLabels struct {
	Pending   string `json:"pending"`
	Completed string `json:"completed"`
	Failed    string `json:"failed"`
}

type ConfirmLabels struct {
	MarkComplete          string `json:"markComplete"`
	MarkCompleteMessage   string `json:"markCompleteMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkComplete          string `json:"bulkComplete"`
	BulkCompleteMessage   string `json:"bulkCompleteMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:          "Collections",
			HeadingPending:   "Pending Collections",
			HeadingCompleted: "Completed Collections",
			HeadingFailed:    "Failed Collections",
			Caption:          "Manage payment collections",
			CaptionPending:   "Payments awaiting collection",
			CaptionCompleted: "Successfully collected payments",
			CaptionFailed:    "Failed payment attempts",
			Dashboard:        "Collections Dashboard",
		},
		Buttons: ButtonLabels{
			AddCollection: "Add Collection",
		},
		Columns: ColumnLabels{
			Reference: "Reference",
			Customer:  "Customer",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
		},
		Empty: EmptyLabels{
			PendingTitle:     "No pending collections",
			PendingMessage:   "No pending collections to display.",
			CompletedTitle:   "No completed collections",
			CompletedMessage: "No completed collections to display.",
			FailedTitle:      "No failed collections",
			FailedMessage:    "No failed collections to display.",
		},
		Form: FormLabels{
			Customer:                "Customer",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. INV-001",
			PaymentMethod:           "Payment Method",
			Status:                  "Status",
			Notes:                   "Notes",
			NotesPlaceholder:        "Additional notes...",
			CustomerNamePlaceholder: "Customer name",
			AmountPlaceholder:       "0.00",
			CurrencyPlaceholder:     "PHP",
			MethodCash:              "Cash",
			MethodBankTransfer:      "Bank Transfer",
			MethodCheck:             "Check",
			MethodGCash:             "GCash",
			MethodMaya:              "Maya",
			MethodCard:              "Card",
			MethodOther:             "Other",
			StatusPending:           "Pending",
			StatusCompleted:         "Completed",
			StatusFailed:            "Failed",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ReferenceInfo:     "Unique reference number for this collection record.",
			CustomerInfo:      "Name of the customer or payer.",
			AmountInfo:        "Total amount collected (in centavos; displayed as amount ÷ 100).",
			CurrencyInfo:      "Currency of the collected amount.",
			PaymentMethodInfo: "How the payment was received.",
			DateInfo:          "Date the payment was collected.",
			StatusInfo:        "Current state of this collection record.",
			NotesInfo:         "Internal remarks — not shown on customer-facing documents.",
			// 20260517-advance-cash-events Plan B Phase 4 — advance metadata
			// section header + the two field labels rendered in the form.
			AdvanceMetadata:        "Advance metadata",
			AdvanceKind:            "Advance kind",
			AdvanceProrationPolicy: "Proration policy",
		},
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			MarkComplete: "Mark Complete",
			Reactivate:   "Reactivate",
		},
		Bulk: BulkLabels{
			Delete: "Delete Selected",
		},
		Detail: DetailLabels{
			PageTitle:            "Collection Details",
			TitlePrefix:          "Collection #",
			PaymentInfo:          "Payment Information",
			Customer:             "Customer",
			Date:                 "Date",
			Amount:               "Amount",
			Currency:             "Currency",
			Status:               "Status",
			Method:               "Payment Method",
			Reference:            "Reference",
			Notes:                "Notes",
			TabBasicInfo:         "Basic Info",
			TabAttachments:       "Attachments",
			TabAuditTrail:        "Audit Trail",
			TabAuditHistory:      "History",
			AuditAction:          "Action",
			AuditUser:            "User",
			AuditEmptyTitle:      "No audit records",
			AuditEmptyMessage:    "No audit trail entries yet.",
			AuditTrailComingSoon: "Audit trail coming soon.",
			AuditTrailDesc:       "Audit trail for collection changes is coming soon.",
		},
		Status: StatusLabels{
			Pending:   "Pending",
			Completed: "Completed",
			Failed:    "Failed",
		},
		Confirm: ConfirmLabels{
			MarkComplete:          "Mark Complete",
			MarkCompleteMessage:   "Are you sure you want to mark %s as complete?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate %s?",
			Delete:                "Delete",
			DeleteMessage:         "Are you sure you want to delete %s?",
			BulkComplete:          "Mark Complete",
			BulkCompleteMessage:   "Are you sure you want to mark {{count}} collection(s) as complete?",
			BulkReactivate:        "Reactivate",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} collection(s)?",
			BulkDelete:            "Delete Collections",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} collection(s)?",
		},
		Errors: ErrorLabels{
			PermissionDenied: "Permission denied",
			InvalidFormData:  "Invalid form data",
			NotFound:         "Collection not found",
			IDRequired:       "Collection ID is required",
			NoIDsProvided:    "No collection IDs provided",
			InvalidStatus:    "Invalid status",
		},
		Dashboard: CashDashboardLabels{
			Title:              "Cash",
			Subtitle:           "Track collected payments and outstanding balances",
			StatPending:        "Pending",
			StatOverdue:        "Overdue",
			StatCollectedToday: "Collected Today",
			StatCollectedWeek:  "Collected This Week",
			WidgetDailyTrend:   "Collected per day (30d)",
			WidgetByMode:       "By payment mode",
			WidgetRecent:       "Recent collections",
			QuickRecord:        "Record Collection",
			QuickReconcile:     "Reconcile",
			QuickAging:         "Aging Report",
			QuickMarkCleared:   "Mark Cleared",
			ViewAll:            "View All",
			EmptyRecentTitle:   "No recent collections",
			EmptyRecentDesc:    "Recent collections will appear here once payments are recorded.",
			NewCollection:      "New collection",
			CollectionUpdated:  "Collection updated",
		},
	}
}
