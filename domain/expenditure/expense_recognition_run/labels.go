package expense_recognition_run

// ---------------------------------------------------------------------------
// Expense Recognition Run labels (Plan A buying-side mirror of Revenue Run)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Expense
// Recognition Run module. Lyngua root key: "expenseRecognitionRun".
// Naming: expenseRecognitionRun / expense_recognition_run / ExpenseRecognitionRun
// / expense-recognition-run everywhere except the user-visible VALUE
// "Expense Run" (supplied by lyngua).
type Labels struct {
	AppLabel                 string            `json:"app_label"`
	Labels                   EntityLabels      `json:"labels"`
	Page                     PageLabels        `json:"page"`
	Buttons                  ButtonLabels      `json:"buttons"`
	Search                   SearchLabels      `json:"search"`
	Filters                  FilterLabels      `json:"filters"`
	Columns                  ColumnLabels      `json:"columns"`
	Queue                    QueueLabels       `json:"queue"`
	List                     ListLabels        `json:"list"`
	Detail                   DetailLabels      `json:"detail"`
	Drawer                   DrawerLabels      `json:"drawer"`
	StatusBadges             StatusBadgeLabels `json:"status_badges"`
	Actions                  ActionLabels      `json:"actions"`
	ScopeKind                ScopeKindLabels   `json:"scope_kind"`
	SourceKind               SourceKindLabels  `json:"source_kind"`
	AttemptOutcome           OutcomeLabels     `json:"attempt_outcome"`
	Outcome                  OutcomeLabels     `json:"outcome"`
	LinkedAdvanceSuppression SuppressionLabels `json:"linked_advance_suppression"`
	Empty                    EmptyLabels       `json:"empty"`
	Toast                    ToastLabels       `json:"toast"`
	Errors                   ErrorLabels       `json:"errors"`
}

// EntityLabels holds entity-level labels.
type EntityLabels struct {
	NameSingular string `json:"name_singular"`
	NamePlural   string `json:"name_plural"`
	ModuleTitle  string `json:"module_title"`
}

// PageLabels holds top-level page titles.
type PageLabels struct {
	QueueTitle    string `json:"queue_title"`
	QueueSubtitle string `json:"queue_subtitle"`
	ListTitle     string `json:"list_title"`
	ListSubtitle  string `json:"list_subtitle"`
	DetailTitle   string `json:"detail_title"`
}

// ButtonLabels holds button copy.
type ButtonLabels struct {
	Generate                 string `json:"generate"`
	RunForSelected           string `json:"run_for_selected"`
	RunForAllMatching        string `json:"run_for_all_matching"`
	Cancel                   string `json:"cancel"`
	ViewAttempts             string `json:"view_attempts"`
	ViewRun                  string `json:"view_run"`
	ViewSupplier             string `json:"view_supplier"`
	ViewSupplierSubscription string `json:"view_supplier_subscription"`
	ViewAdvanceDisbursement  string `json:"view_advance_disbursement"`
	ReRunFailed              string `json:"re_run_failed"`
	RunRecognitions          string `json:"run_recognitions"`
}

// SearchLabels holds search-input copy.
type SearchLabels struct {
	Placeholder string `json:"placeholder"`
}

// FilterLabels holds filter chip labels.
type FilterLabels struct {
	AsOfDate string `json:"as_of_date"`
	Supplier string `json:"supplier"`
	Status   string `json:"status"`
	Pending  string `json:"pending"`
	Complete string `json:"complete"`
	Failed   string `json:"failed"`
}

// ColumnLabels holds the top-level column labels.
type ColumnLabels struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AsOfDate    string `json:"as_of_date"`
	Initiator   string `json:"initiator"`
	InitiatedAt string `json:"initiated_at"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Skipped     string `json:"skipped"`
	Errored     string `json:"errored"`
	Actions     string `json:"actions"`
}

// QueueLabels holds copy for the workspace-queue page
// (Surface B).
type QueueLabels struct {
	Title         string            `json:"title"`
	Subtitle      string            `json:"subtitle"`
	AsOfDateLabel string            `json:"as_of_date_label"`
	Columns       QueueColumnLabels `json:"columns"`
	Empty         QueueEmptyLabels  `json:"empty"`
	Bulk          QueueBulkLabels   `json:"bulk"`
}

type QueueColumnLabels struct {
	Supplier             string `json:"supplier"`
	Subscriptions        string `json:"subscriptions"`
	AdvanceDisbursements string `json:"advance_disbursements"`
	PendingPeriods       string `json:"pending_periods"`
	Total                string `json:"total"`
	Currency             string `json:"currency"`
	Actions              string `json:"actions"`
	Run                  string `json:"run"`
}

type QueueEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type QueueBulkLabels struct {
	RunSelected        string `json:"run_selected"`
	RunAllMatching     string `json:"run_all_matching"`
	CapExceededMessage string `json:"cap_exceeded_message"`
}

// ListLabels holds copy for the run history list page
// (Surface D).
type ListLabels struct {
	Title    string           `json:"title"`
	Subtitle string           `json:"subtitle"`
	Columns  ListColumnLabels `json:"columns"`
	Empty    ListEmptyLabels  `json:"empty"`
	Filters  ListFilterLabels `json:"filter_labels"`
}

type ListColumnLabels struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AsOfDate    string `json:"as_of_date"`
	Initiator   string `json:"initiator"`
	InitiatedAt string `json:"initiated_at"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Skipped     string `json:"skipped"`
	Errored     string `json:"errored"`
	Actions     string `json:"actions"`
}

type ListEmptyLabels struct {
	Pending  ListEmptyStateLabels `json:"pending"`
	Complete ListEmptyStateLabels `json:"complete"`
	Failed   ListEmptyStateLabels `json:"failed"`
}

type ListEmptyStateLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ListFilterLabels struct {
	Pending  string `json:"pending"`
	Complete string `json:"complete"`
	Failed   string `json:"failed"`
}

// DetailLabels holds copy for the run detail page (Surface D).
type DetailLabels struct {
	Title                  string                `json:"title"`
	Tabs                   DetailTabLabels       `json:"tabs"`
	TabHints               DetailTabHintLabels   `json:"tab_hints"`
	Summary                SummaryLabels         `json:"summary"`
	Selections             SelectionsTabLabels   `json:"selections"`
	Results                ResultsTabLabels      `json:"results"`
	Bills                  BillsTabLabels        `json:"bills"`
	Recognitions           RecognitionsTabLabels `json:"recognitions"`
	AuditHistoryComingSoon string                `json:"audit_history_coming_soon"`
}

type DetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Bills        string `json:"bills"`
	Recognitions string `json:"recognitions"`
	AuditHistory string `json:"audit_history"`
	Attachments  string `json:"attachments"`
}

type DetailTabHintLabels struct {
	BillsHint        string `json:"bills_hint"`
	RecognitionsHint string `json:"recognitions_hint"`
}

type SummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"as_of_date"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiated_at"`
	CompletedAt             string `json:"completed_at"`
	Status                  string `json:"status"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possibly_interrupted_note"`
}

type SelectionsTabLabels struct {
	ColSource               string `json:"col_source"`
	ColSupplierSubscription string `json:"col_supplier_subscription"`
	ColAdvanceDisbursement  string `json:"col_advance_disbursement"`
	ColPeriodStart          string `json:"col_period_start"`
	ColPeriodEnd            string `json:"col_period_end"`
	ColPeriodMarker         string `json:"col_period_marker"`
	EmptyTitle              string `json:"empty_title"`
	EmptyMessage            string `json:"empty_message"`
}

type ResultsTabLabels struct {
	ColSource               string `json:"col_source"`
	ColSupplierSubscription string `json:"col_supplier_subscription"`
	ColAdvanceDisbursement  string `json:"col_advance_disbursement"`
	ColPeriodStart          string `json:"col_period_start"`
	ColPeriodEnd            string `json:"col_period_end"`
	ColOutcome              string `json:"col_outcome"`
	ColErrorCode            string `json:"col_error_code"`
	EmptyTitle              string `json:"empty_title"`
	EmptyMessage            string `json:"empty_message"`
}

type BillsTabLabels struct {
	ColReference string `json:"col_reference"`
	ColDate      string `json:"col_date"`
	ColAmount    string `json:"col_amount"`
	ColStatus    string `json:"col_status"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
	Hint         string `json:"hint"`
}

type RecognitionsTabLabels struct {
	ColReference  string `json:"col_reference"`
	ColDate       string `json:"col_date"`
	ColAmount     string `json:"col_amount"`
	ColSourceKind string `json:"col_source_kind"`
	ColStatus     string `json:"col_status"`
	EmptyTitle    string `json:"empty_title"`
	EmptyMessage  string `json:"empty_message"`
	Hint          string `json:"hint"`
}

// DrawerLabels holds drawer-form labels for Surface A
// (per-supplier), Surface C (per-supplier-subscription), and the
// generate-confirmation modal.
type DrawerLabels struct {
	Supplier     SupplierDrawerLabels     `json:"supplier"`
	Subscription SubscriptionDrawerLabels `json:"subscription"`
	Confirmation ConfirmationLabels       `json:"confirmation"`
}

type SupplierDrawerLabels struct {
	Title                           string `json:"title"`
	SubtitleTemplate                string `json:"subtitle_template"`
	AsOfDateLabel                   string `json:"as_of_date_label"`
	AsOfDateHint                    string `json:"as_of_date_hint"`
	ColumnSource                    string `json:"column_source"`
	ColumnPeriod                    string `json:"column_period"`
	ColumnAmount                    string `json:"column_amount"`
	ColumnLines                     string `json:"column_lines"`
	ColumnRemaining                 string `json:"column_remaining"`
	GroupSubscriptionCycle          string `json:"group_subscription_cycle"`
	GroupAdvanceDisbursementTranche string `json:"group_advance_disbursement_tranche"`
	GroupNoPending                  string `json:"group_no_pending"`
	GroupCurrencyMismatch           string `json:"group_currency_mismatch"`
	EmptyTitle                      string `json:"empty_title"`
	EmptyMessage                    string `json:"empty_message"`
	GenerateButton                  string `json:"generate_button"`
	GenerateButtonCount             string `json:"generate_button_count"`
	CancelButton                    string `json:"cancel_button"`
	ViewRunLink                     string `json:"view_run_link"`
}

type SubscriptionDrawerLabels struct {
	Title                        string `json:"title"`
	SubtitleTemplate             string `json:"subtitle_template"`
	AsOfDateLabel                string `json:"as_of_date_label"`
	AsOfDateHint                 string `json:"as_of_date_hint"`
	ColumnPeriod                 string `json:"column_period"`
	ColumnAmount                 string `json:"column_amount"`
	ColumnLines                  string `json:"column_lines"`
	EmptyTitle                   string `json:"empty_title"`
	EmptyMessage                 string `json:"empty_message"`
	SuppressedByAdvanceTitle     string `json:"suppressed_by_advance_title"`
	SuppressedByAdvanceExplainer string `json:"suppressed_by_advance_explainer"`
	ViewAdvanceLink              string `json:"view_advance_link"`
	GenerateButton               string `json:"generate_button"`
	GenerateButtonCount          string `json:"generate_button_count"`
	CancelButton                 string `json:"cancel_button"`
	ViewRunLink                  string `json:"view_run_link"`
}

type ConfirmationLabels struct {
	Title          string `json:"title"`
	BodyTemplate   string `json:"body_template"`
	NoteIdempotent string `json:"note_idempotent"`
	ConfirmButton  string `json:"confirm_button"`
	CancelButton   string `json:"cancel_button"`
}

// StatusBadgeLabels holds badge copy for each run status.
type StatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possibly_interrupted"`
}

// ActionLabels holds labels for interactive actions on
// run rows / pages / drawer triggers.
type ActionLabels struct {
	Run                      string `json:"run"`
	RunRecognitions          string `json:"run_recognitions"`
	ReRunFailed              string `json:"re_run_failed"`
	ReRunFailedComingSoon    string `json:"re_run_failed_coming_soon"`
	ViewRun                  string `json:"view_run"`
	ViewSupplier             string `json:"view_supplier"`
	ViewSupplierSubscription string `json:"view_supplier_subscription"`
	ViewAdvanceDisbursement  string `json:"view_advance_disbursement"`
	RunAriaLabel             string `json:"run_aria_label"`
}

// ScopeKindLabels holds display labels for each Run
// scope kind: supplier (Surface A), subscription (Surface C), workspace
// (Surface B).
type ScopeKindLabels struct {
	Supplier     string `json:"supplier"`
	Subscription string `json:"subscription"`
	Workspace    string `json:"workspace"`
}

// SourceKindLabels holds display labels for each
// candidate source kind on a run attempt: subscription cycle vs advance
// disbursement tranche.
type SourceKindLabels struct {
	SubscriptionCycle   string `json:"subscription_cycle"`
	AdvanceDisbursement string `json:"advance_disbursement"`
}

// OutcomeLabels holds display labels for per-attempt
// outcome values.
type OutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// SuppressionLabels holds copy for the
// linked-advance-suppression banner + row chip rendered on Surface A and
// Surface C when a SupplierSubscription cycle is covered by a TIME_BASED
// advance TreasuryDisbursement (Plan B decision A).
type SuppressionLabels struct {
	BannerTitle     string `json:"banner_title"`
	BannerMessage   string `json:"banner_message"`
	RowChip         string `json:"row_chip"`
	RowExplainer    string `json:"row_explainer"`
	ViewAdvanceLink string `json:"view_advance_link"`
	AriaLabel       string `json:"aria_label"`
}

// EmptyLabels holds empty-state copy for every surface.
type EmptyLabels struct {
	QueueTitle          string `json:"queue_title"`
	QueueMessage        string `json:"queue_message"`
	ListTitle           string `json:"list_title"`
	ListMessage         string `json:"list_message"`
	SelectionsTitle     string `json:"selections_title"`
	SelectionsMessage   string `json:"selections_message"`
	ResultsTitle        string `json:"results_title"`
	ResultsMessage      string `json:"results_message"`
	BillsTitle          string `json:"bills_title"`
	BillsMessage        string `json:"bills_message"`
	RecognitionsTitle   string `json:"recognitions_title"`
	RecognitionsMessage string `json:"recognitions_message"`
}

// ToastLabels holds toast / notification copy.
type ToastLabels struct {
	Success              string `json:"success"`
	BatchSuccess         string `json:"batch_success"`
	BatchSuccessMultiRun string `json:"batch_success_multi_run"`
	ViewRunLink          string `json:"view_run_link"`
	GenerateFailed       string `json:"generate_failed"`
	PermissionDenied     string `json:"permission_denied"`
}

// ErrorLabels holds error-message strings for the module.
type ErrorLabels struct {
	CapExceeded                  string `json:"cap_exceeded"`
	PermissionDenied             string `json:"permission_denied"`
	UseCaseUnavailable           string `json:"use_case_unavailable"`
	InvalidSelection             string `json:"invalid_selection"`
	IdempotencyConflict          string `json:"idempotency_conflict"`
	SupplierMismatch             string `json:"supplier_mismatch"`
	WorkspaceMismatch            string `json:"workspace_mismatch"`
	TamperedPeriod               string `json:"tampered_period"`
	RunAllMatchingNotImplemented string `json:"run_all_matching_not_implemented"`
	CrossWorkspace               string `json:"cross_workspace"`
	MissingCostPlan              string `json:"missing_cost_plan"`
	NoPendingPeriods             string `json:"no_pending_periods"`
	CurrencyMismatch             string `json:"currency_mismatch"`
	GenerationFailed             string `json:"generation_failed"`
	SuppressedByAdvance          string `json:"suppressed_by_advance"`
}

// DefaultLabels returns Labels with
// sensible English defaults. Lyngua overlays via the "expenseRecognitionRun"
// root key in general/expense_recognition_run.json.
func DefaultLabels() Labels {
	return Labels{
		AppLabel: "Expense Run",
		Labels: EntityLabels{
			NameSingular: "Expense Run",
			NamePlural:   "Expense Runs",
			ModuleTitle:  "Expense Recognition Run",
		},
		Page: PageLabels{
			QueueTitle:    "Expense Run Queue",
			QueueSubtitle: "Suppliers with pending subscription cycles or advance disbursement tranches ready to recognize",
			ListTitle:     "Expense Runs",
			ListSubtitle:  "History of expense run batches",
			DetailTitle:   "Expense Run",
		},
		Buttons: ButtonLabels{
			Generate:                 "Generate",
			RunForSelected:           "Run for selected",
			RunForAllMatching:        "Run for all matching",
			Cancel:                   "Cancel",
			ViewAttempts:             "View attempts",
			ViewRun:                  "View run",
			ViewSupplier:             "View supplier",
			ViewSupplierSubscription: "View supplier subscription",
			ViewAdvanceDisbursement:  "View advance disbursement",
			ReRunFailed:              "Re-run failed",
			RunRecognitions:          "Run Recognitions",
		},
		Search: SearchLabels{
			Placeholder: "Search expense runs",
		},
		Filters: FilterLabels{
			AsOfDate: "As of date",
			Supplier: "Supplier",
			Status:   "Status",
			Pending:  "Pending",
			Complete: "Complete",
			Failed:   "Failed",
		},
		Columns: ColumnLabels{
			ID:          "Run",
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
		Queue: QueueLabels{
			Title:         "Expense Run Queue",
			Subtitle:      "Suppliers with pending subscription cycles or advance disbursement tranches ready to recognize",
			AsOfDateLabel: "As of date",
			Columns: QueueColumnLabels{
				Supplier:             "Supplier",
				Subscriptions:        "Subscriptions",
				AdvanceDisbursements: "Advance Disbursements",
				PendingPeriods:       "Pending periods",
				Total:                "Total",
				Currency:             "Currency",
				Actions:              "Actions",
				Run:                  "Run",
			},
			Empty: QueueEmptyLabels{
				Title:   "Nothing to recognize",
				Message: "No suppliers have pending subscription cycles or advance disbursement tranches as of this date.",
			},
			Bulk: QueueBulkLabels{
				RunSelected:        "Run for selected",
				RunAllMatching:     "Run for all matching",
				CapExceededMessage: "Capped at 50 suppliers per batch. Narrow the filter to run the rest.",
			},
		},
		List: ListLabels{
			Title:    "Expense Runs",
			Subtitle: "History of expense run batches",
			Columns: ListColumnLabels{
				ID:          "Run",
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
			Empty: ListEmptyLabels{
				Pending: ListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "Runs in flight will appear here.",
				},
				Complete: ListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "Completed runs will appear here.",
				},
				Failed: ListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "Failed runs will appear here.",
				},
			},
			Filters: ListFilterLabels{
				Pending:  "Pending",
				Complete: "Complete",
				Failed:   "Failed",
			},
		},
		Detail: DetailLabels{
			Title: "Expense Run",
			Tabs: DetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Bills:        "Draft Bills",
				Recognitions: "Recognitions",
				AuditHistory: "Audit History",
				Attachments:  "Attachments",
			},
			TabHints: DetailTabHintLabels{
				BillsHint:        "Draft Expenditure rows created by this run. AP edits or cancels them when actual vendor bills arrive.",
				RecognitionsHint: "ExpenseRecognition rows created or amortized by this run.",
			},
			Summary: SummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some recognitions may be missing.",
			},
			Selections: SelectionsTabLabels{
				ColSource:               "Source",
				ColSupplierSubscription: "Supplier Subscription",
				ColAdvanceDisbursement:  "Advance Disbursement",
				ColPeriodStart:          "Period start",
				ColPeriodEnd:            "Period end",
				ColPeriodMarker:         "Period marker",
				EmptyTitle:              "No selections",
				EmptyMessage:            "This run has no attempt records.",
			},
			Results: ResultsTabLabels{
				ColSource:               "Source",
				ColSupplierSubscription: "Supplier Subscription",
				ColAdvanceDisbursement:  "Advance Disbursement",
				ColPeriodStart:          "Period start",
				ColPeriodEnd:            "Period end",
				ColOutcome:              "Outcome",
				ColErrorCode:            "Error code",
				EmptyTitle:              "No results",
				EmptyMessage:            "This run has no attempt records.",
			},
			Bills: BillsTabLabels{
				ColReference: "Reference",
				ColDate:      "Date",
				ColAmount:    "Amount",
				ColStatus:    "Status",
				EmptyTitle:   "No draft bills",
				EmptyMessage: "No draft Expenditure rows were created by this run.",
				Hint:         "Draft Expenditure rows created by this run. AP edits or cancels them when actual vendor bills arrive.",
			},
			Recognitions: RecognitionsTabLabels{
				ColReference:  "Reference",
				ColDate:       "Date",
				ColAmount:     "Amount",
				ColSourceKind: "Source",
				ColStatus:     "Status",
				EmptyTitle:    "No recognitions",
				EmptyMessage:  "No ExpenseRecognition rows were created by this run.",
				Hint:          "ExpenseRecognition rows created or amortized by this run.",
			},
			AuditHistoryComingSoon: "Audit history is coming soon.",
		},
		Drawer: DrawerLabels{
			Supplier: SupplierDrawerLabels{
				Title:                           "Run Recognitions",
				SubtitleTemplate:                "Pending periods for {{.SupplierName}} as of {{.AsOfDate}}",
				AsOfDateLabel:                   "As of date",
				AsOfDateHint:                    "Defaults to today in the workspace timezone.",
				ColumnSource:                    "Source",
				ColumnPeriod:                    "Period",
				ColumnAmount:                    "Amount",
				ColumnLines:                     "Lines",
				ColumnRemaining:                 "Remaining",
				GroupSubscriptionCycle:          "Subscription cycles",
				GroupAdvanceDisbursementTranche: "Advance disbursement tranches",
				GroupNoPending:                  "No pending periods",
				GroupCurrencyMismatch:           "Currency mismatch",
				EmptyTitle:                      "Nothing to recognize",
				EmptyMessage:                    "This supplier has no pending subscription cycles or advance disbursement tranches as of this date.",
				GenerateButton:                  "Generate",
				GenerateButtonCount:             "Generate ({{.Count}})",
				CancelButton:                    "Cancel",
				ViewRunLink:                     "View run",
			},
			Subscription: SubscriptionDrawerLabels{
				Title:                        "Run Recognitions",
				SubtitleTemplate:             "Pending periods for {{.SubscriptionName}} as of {{.AsOfDate}}",
				AsOfDateLabel:                "As of date",
				AsOfDateHint:                 "Defaults to today in the workspace timezone.",
				ColumnPeriod:                 "Period",
				ColumnAmount:                 "Amount",
				ColumnLines:                  "Lines",
				EmptyTitle:                   "Nothing to recognize",
				EmptyMessage:                 "This supplier subscription has no pending periods as of this date.",
				SuppressedByAdvanceTitle:     "Suppressed by linked advance",
				SuppressedByAdvanceExplainer: "Cycles for this subscription are recognized via the linked advance disbursement.",
				ViewAdvanceLink:              "View advance",
				GenerateButton:               "Generate",
				GenerateButtonCount:          "Generate ({{.Count}})",
				CancelButton:                 "Cancel",
				ViewRunLink:                  "View run",
			},
			Confirmation: ConfirmationLabels{
				Title:          "Confirm Generate",
				BodyTemplate:   "Generate {{.Count}} recognitions for {{.ScopeName}} as of {{.AsOfDate}}?",
				NoteIdempotent: "Re-running the same period for the same source skips already-recognized rows; no duplicates will be created.",
				ConfirmButton:  "Generate",
				CancelButton:   "Cancel",
			},
		},
		StatusBadges: StatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		Actions: ActionLabels{
			Run:                      "Run",
			RunRecognitions:          "Run Recognitions",
			ReRunFailed:              "Re-run failed",
			ReRunFailedComingSoon:    "Re-run failed (coming soon)",
			ViewRun:                  "View run",
			ViewSupplier:             "View supplier",
			ViewSupplierSubscription: "View supplier subscription",
			ViewAdvanceDisbursement:  "View advance disbursement",
			RunAriaLabel:             "Run expense recognition",
		},
		ScopeKind: ScopeKindLabels{
			Supplier:     "Supplier",
			Subscription: "Supplier Subscription",
			Workspace:    "Workspace",
		},
		SourceKind: SourceKindLabels{
			SubscriptionCycle:   "Subscription cycle",
			AdvanceDisbursement: "Advance disbursement",
		},
		AttemptOutcome: OutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Outcome: OutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		LinkedAdvanceSuppression: SuppressionLabels{
			BannerTitle:     "Cycle suppressed by linked advance",
			BannerMessage:   "An advance disbursement covers this subscription cycle. Recognition flows through the advance amortization.",
			RowChip:         "Suppressed",
			RowExplainer:    "Covered by linked advance disbursement",
			ViewAdvanceLink: "View advance disbursement",
			AriaLabel:       "Row suppressed by linked advance disbursement",
		},
		Empty: EmptyLabels{
			QueueTitle:          "Nothing to recognize",
			QueueMessage:        "No suppliers have pending subscription cycles or advance disbursement tranches as of this date.",
			ListTitle:           "No expense runs",
			ListMessage:         "Expense runs you initiate will appear here.",
			SelectionsTitle:     "No selections",
			SelectionsMessage:   "This run has no attempt records.",
			ResultsTitle:        "No results",
			ResultsMessage:      "This run has no attempt records.",
			BillsTitle:          "No draft bills",
			BillsMessage:        "No draft Expenditure rows were created by this run.",
			RecognitionsTitle:   "No recognitions",
			RecognitionsMessage: "No ExpenseRecognition rows were created by this run.",
		},
		Toast: ToastLabels{
			Success:              "Recognized {{.Created}} of {{.Total}} ({{.Skipped}} skipped, {{.Errored}} errored)",
			BatchSuccess:         "Expense batch run — {{.Created}} created, {{.Skipped}} skipped, {{.Errored}} failed.",
			BatchSuccessMultiRun: "Ran {{.RunCount}} expense runs ({{.Created}} created, {{.Skipped}} skipped, {{.Errored}} errored)",
			ViewRunLink:          "View run",
			GenerateFailed:       "Failed to generate expense run.",
			PermissionDenied:     "You do not have permission to run expense recognition.",
		},
		Errors: ErrorLabels{
			CapExceeded:                  "Run for selected is capped at 50 suppliers per batch.",
			PermissionDenied:             "You do not have permission to run expense recognition.",
			UseCaseUnavailable:           "Expense run is not available for this supplier subscription type.",
			InvalidSelection:             "One or more selected supplier subscriptions are invalid.",
			IdempotencyConflict:          "A recognition for one or more periods already exists.",
			SupplierMismatch:             "Selected supplier subscriptions belong to different suppliers.",
			WorkspaceMismatch:            "Selected supplier subscriptions belong to a different workspace.",
			TamperedPeriod:               "A period was modified after selection. Please retry.",
			RunAllMatchingNotImplemented: "Run for all matching is not yet available. Please select individual suppliers.",
			CrossWorkspace:               "Selection crosses workspace boundary.",
			MissingCostPlan:              "SupplierSubscription has no active CostPlan.",
			NoPendingPeriods:             "No pending periods for the chosen as-of date.",
			CurrencyMismatch:             "Currency mismatch between source and recognition.",
			GenerationFailed:             "Expense run generation failed.",
			SuppressedByAdvance:          "Subscription cycle suppressed by linked advance disbursement.",
		},
	}
}
