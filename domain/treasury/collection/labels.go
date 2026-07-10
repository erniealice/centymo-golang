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
	Bulk      BulkLabels          `json:"bulk_actions"`
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
	StatPending        string `json:"stat_pending"`
	StatOverdue        string `json:"stat_overdue"`
	StatCollectedToday string `json:"stat_collected_today"`
	StatCollectedWeek  string `json:"stat_collected_week"`
	WidgetDailyTrend   string `json:"widget_daily_trend"`
	WidgetByMode       string `json:"widget_by_mode"`
	WidgetRecent       string `json:"widget_recent"`
	QuickRecord        string `json:"quick_record"`
	QuickReconcile     string `json:"quick_reconcile"`
	QuickAging         string `json:"quick_aging"`
	QuickMarkCleared   string `json:"quick_mark_cleared"`
	ViewAll            string `json:"view_all"`
	EmptyRecentTitle   string `json:"empty_recent_title"`
	EmptyRecentDesc    string `json:"empty_recent_desc"`
	NewCollection      string `json:"new_collection"`
	CollectionUpdated  string `json:"collection_updated"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingPending   string `json:"heading_pending"`
	HeadingCompleted string `json:"heading_completed"`
	HeadingFailed    string `json:"heading_failed"`
	Caption          string `json:"caption"`
	CaptionPending   string `json:"caption_pending"`
	CaptionCompleted string `json:"caption_completed"`
	CaptionFailed    string `json:"caption_failed"`
	Dashboard        string `json:"dashboard"`
}

type ButtonLabels struct {
	AddCollection string `json:"add_collection"`
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
	PendingTitle     string `json:"pending_title"`
	PendingMessage   string `json:"pending_message"`
	CompletedTitle   string `json:"completed_title"`
	CompletedMessage string `json:"completed_message"`
	FailedTitle      string `json:"failed_title"`
	FailedMessage    string `json:"failed_message"`
}

type FormLabels struct {
	Customer                string `json:"customer"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"reference_placeholder"`
	PaymentMethod           string `json:"payment_method"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notes_placeholder"`
	CustomerNamePlaceholder string `json:"customer_name_placeholder"`
	AmountPlaceholder       string `json:"amount_placeholder"`
	CurrencyPlaceholder     string `json:"currency_placeholder"`
	MethodCash              string `json:"method_cash"`
	MethodBankTransfer      string `json:"method_bank_transfer"`
	MethodCheck             string `json:"method_check"`
	MethodGCash             string `json:"method_gcash"`
	MethodMaya              string `json:"method_maya"`
	MethodCard              string `json:"method_card"`
	MethodOther             string `json:"method_other"`
	StatusPending           string `json:"status_pending"`
	StatusCompleted         string `json:"status_completed"`
	StatusFailed            string `json:"status_failed"`

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo     string `json:"reference_info"`
	CustomerInfo      string `json:"customer_info"`
	AmountInfo        string `json:"amount_info"`
	CurrencyInfo      string `json:"currency_info"`
	PaymentMethodInfo string `json:"payment_method_info"`
	DateInfo          string `json:"date_info"`
	StatusInfo        string `json:"status_info"`
	NotesInfo         string `json:"notes_info"`

	// 20260517-advance-cash-events Plan B Phase 4 — advance metadata fields
	// rendered conditionally in the collection drawer form. The enum option
	// labels (None / Time-based / Milestone / Unscheduled / Full tranche /
	// Day-prorated / Next period start) are sourced from AdvanceEnumLabels
	// (loaded from advance_kind.json) rather than duplicated here.
	AdvanceMetadata        string `json:"advance_metadata"`
	AdvanceKind            string `json:"advance_kind"`
	AdvanceProrationPolicy string `json:"advance_proration_policy"`
}

type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"mark_complete"`
	Reactivate   string `json:"reactivate"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle            string `json:"page_title"`
	TitlePrefix          string `json:"title_prefix"`
	PaymentInfo          string `json:"payment_info"`
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Method               string `json:"method"`
	Reference            string `json:"reference"`
	Notes                string `json:"notes"`
	TabBasicInfo         string `json:"tab_basic_info"`
	TabAttachments       string `json:"tab_attachments"`
	TabAuditTrail        string `json:"tab_audit_trail"`
	TabAuditHistory      string `json:"tab_audit_history"`
	AuditAction          string `json:"audit_action"`
	AuditUser            string `json:"audit_user"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditTrailDesc       string `json:"audit_trail_desc"`
}

type StatusLabels struct {
	Pending   string `json:"pending"`
	Completed string `json:"completed"`
	Failed    string `json:"failed"`
}

type ConfirmLabels struct {
	MarkComplete          string `json:"mark_complete"`
	MarkCompleteMessage   string `json:"mark_complete_message"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivate_message"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	BulkComplete          string `json:"bulk_complete"`
	BulkCompleteMessage   string `json:"bulk_complete_message"`
	BulkReactivate        string `json:"bulk_reactivate"`
	BulkReactivateMessage string `json:"bulk_reactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
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
