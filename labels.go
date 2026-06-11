package centymo

// ---------------------------------------------------------------------------
// Expense Recognition Run labels (Plan A buying-side mirror of Revenue Run)
// ---------------------------------------------------------------------------

// ExpenseRecognitionRunLabels holds all translatable strings for the Expense
// Recognition Run module. Lyngua root key: "expenseRecognitionRun".
// Naming: expenseRecognitionRun / expense_recognition_run / ExpenseRecognitionRun
// / expense-recognition-run everywhere except the user-visible VALUE
// "Expense Run" (supplied by lyngua).
type ExpenseRecognitionRunLabels struct {
	AppLabel                 string                                 `json:"appLabel"`
	Labels                   ExpenseRecognitionRunEntityLabels      `json:"labels"`
	Page                     ExpenseRecognitionRunPageLabels        `json:"page"`
	Buttons                  ExpenseRecognitionRunButtonLabels      `json:"buttons"`
	Search                   ExpenseRecognitionRunSearchLabels      `json:"search"`
	Filters                  ExpenseRecognitionRunFilterLabels      `json:"filters"`
	Columns                  ExpenseRecognitionRunColumnLabels      `json:"columns"`
	Queue                    ExpenseRecognitionRunQueueLabels       `json:"queue"`
	List                     ExpenseRecognitionRunListLabels        `json:"list"`
	Detail                   ExpenseRecognitionRunDetailLabels      `json:"detail"`
	Drawer                   ExpenseRecognitionRunDrawerLabels      `json:"drawer"`
	StatusBadges             ExpenseRecognitionRunStatusBadgeLabels `json:"statusBadges"`
	Actions                  ExpenseRecognitionRunActionLabels      `json:"actions"`
	ScopeKind                ExpenseRecognitionRunScopeKindLabels   `json:"scopeKind"`
	SourceKind               ExpenseRecognitionRunSourceKindLabels  `json:"sourceKind"`
	AttemptOutcome           ExpenseRecognitionRunOutcomeLabels     `json:"attemptOutcome"`
	Outcome                  ExpenseRecognitionRunOutcomeLabels     `json:"outcome"`
	LinkedAdvanceSuppression ExpenseRecognitionRunSuppressionLabels `json:"linkedAdvanceSuppression"`
	Empty                    ExpenseRecognitionRunEmptyLabels       `json:"empty"`
	Toast                    ExpenseRecognitionRunToastLabels       `json:"toast"`
	Errors                   ExpenseRecognitionRunErrorLabels       `json:"errors"`
}

// ExpenseRecognitionRunEntityLabels holds entity-level labels.
type ExpenseRecognitionRunEntityLabels struct {
	NameSingular string `json:"nameSingular"`
	NamePlural   string `json:"namePlural"`
	ModuleTitle  string `json:"moduleTitle"`
}

// ExpenseRecognitionRunPageLabels holds top-level page titles.
type ExpenseRecognitionRunPageLabels struct {
	QueueTitle    string `json:"queueTitle"`
	QueueSubtitle string `json:"queueSubtitle"`
	ListTitle     string `json:"listTitle"`
	ListSubtitle  string `json:"listSubtitle"`
	DetailTitle   string `json:"detailTitle"`
}

// ExpenseRecognitionRunButtonLabels holds button copy.
type ExpenseRecognitionRunButtonLabels struct {
	Generate                 string `json:"generate"`
	RunForSelected           string `json:"runForSelected"`
	RunForAllMatching        string `json:"runForAllMatching"`
	Cancel                   string `json:"cancel"`
	ViewAttempts             string `json:"viewAttempts"`
	ViewRun                  string `json:"viewRun"`
	ViewSupplier             string `json:"viewSupplier"`
	ViewSupplierSubscription string `json:"viewSupplierSubscription"`
	ViewAdvanceDisbursement  string `json:"viewAdvanceDisbursement"`
	ReRunFailed              string `json:"reRunFailed"`
	RunRecognitions          string `json:"runRecognitions"`
}

// ExpenseRecognitionRunSearchLabels holds search-input copy.
type ExpenseRecognitionRunSearchLabels struct {
	Placeholder string `json:"placeholder"`
}

// ExpenseRecognitionRunFilterLabels holds filter chip labels.
type ExpenseRecognitionRunFilterLabels struct {
	AsOfDate string `json:"asOfDate"`
	Supplier string `json:"supplier"`
	Status   string `json:"status"`
	Pending  string `json:"pending"`
	Complete string `json:"complete"`
	Failed   string `json:"failed"`
}

// ExpenseRecognitionRunColumnLabels holds the top-level column labels.
type ExpenseRecognitionRunColumnLabels struct {
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

// ExpenseRecognitionRunQueueLabels holds copy for the workspace-queue page
// (Surface B).
type ExpenseRecognitionRunQueueLabels struct {
	Title         string                                 `json:"title"`
	Subtitle      string                                 `json:"subtitle"`
	AsOfDateLabel string                                 `json:"asOfDateLabel"`
	Columns       ExpenseRecognitionRunQueueColumnLabels `json:"columns"`
	Empty         ExpenseRecognitionRunQueueEmptyLabels  `json:"empty"`
	Bulk          ExpenseRecognitionRunQueueBulkLabels   `json:"bulk"`
}

type ExpenseRecognitionRunQueueColumnLabels struct {
	Supplier             string `json:"supplier"`
	Subscriptions        string `json:"subscriptions"`
	AdvanceDisbursements string `json:"advanceDisbursements"`
	PendingPeriods       string `json:"pendingPeriods"`
	Total                string `json:"total"`
	Currency             string `json:"currency"`
	Actions              string `json:"actions"`
	Run                  string `json:"run"`
}

type ExpenseRecognitionRunQueueEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ExpenseRecognitionRunQueueBulkLabels struct {
	RunSelected        string `json:"runSelected"`
	RunAllMatching     string `json:"runAllMatching"`
	CapExceededMessage string `json:"capExceededMessage"`
}

// ExpenseRecognitionRunListLabels holds copy for the run history list page
// (Surface D).
type ExpenseRecognitionRunListLabels struct {
	Title    string                                `json:"title"`
	Subtitle string                                `json:"subtitle"`
	Columns  ExpenseRecognitionRunListColumnLabels `json:"columns"`
	Empty    ExpenseRecognitionRunListEmptyLabels  `json:"empty"`
	Filters  ExpenseRecognitionRunListFilterLabels `json:"filterLabels"`
}

type ExpenseRecognitionRunListColumnLabels struct {
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

type ExpenseRecognitionRunListEmptyLabels struct {
	Pending  ExpenseRecognitionRunListEmptyStateLabels `json:"pending"`
	Complete ExpenseRecognitionRunListEmptyStateLabels `json:"complete"`
	Failed   ExpenseRecognitionRunListEmptyStateLabels `json:"failed"`
}

type ExpenseRecognitionRunListEmptyStateLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ExpenseRecognitionRunListFilterLabels struct {
	Pending  string `json:"pending"`
	Complete string `json:"complete"`
	Failed   string `json:"failed"`
}

// ExpenseRecognitionRunDetailLabels holds copy for the run detail page (Surface D).
type ExpenseRecognitionRunDetailLabels struct {
	Title                  string                                     `json:"title"`
	Tabs                   ExpenseRecognitionRunDetailTabLabels       `json:"tabs"`
	TabHints               ExpenseRecognitionRunDetailTabHintLabels   `json:"tabHints"`
	Summary                ExpenseRecognitionRunSummaryLabels         `json:"summary"`
	Selections             ExpenseRecognitionRunSelectionsTabLabels   `json:"selections"`
	Results                ExpenseRecognitionRunResultsTabLabels      `json:"results"`
	Bills                  ExpenseRecognitionRunBillsTabLabels        `json:"bills"`
	Recognitions           ExpenseRecognitionRunRecognitionsTabLabels `json:"recognitions"`
	AuditHistoryComingSoon string                                     `json:"auditHistoryComingSoon"`
}

type ExpenseRecognitionRunDetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Bills        string `json:"bills"`
	Recognitions string `json:"recognitions"`
	AuditHistory string `json:"auditHistory"`
	Attachments  string `json:"attachments"`
}

type ExpenseRecognitionRunDetailTabHintLabels struct {
	BillsHint        string `json:"billsHint"`
	RecognitionsHint string `json:"recognitionsHint"`
}

type ExpenseRecognitionRunSummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"asOfDate"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiatedAt"`
	CompletedAt             string `json:"completedAt"`
	Status                  string `json:"status"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possiblyInterruptedNote"`
}

type ExpenseRecognitionRunSelectionsTabLabels struct {
	ColSource               string `json:"colSource"`
	ColSupplierSubscription string `json:"colSupplierSubscription"`
	ColAdvanceDisbursement  string `json:"colAdvanceDisbursement"`
	ColPeriodStart          string `json:"colPeriodStart"`
	ColPeriodEnd            string `json:"colPeriodEnd"`
	ColPeriodMarker         string `json:"colPeriodMarker"`
	EmptyTitle              string `json:"emptyTitle"`
	EmptyMessage            string `json:"emptyMessage"`
}

type ExpenseRecognitionRunResultsTabLabels struct {
	ColSource               string `json:"colSource"`
	ColSupplierSubscription string `json:"colSupplierSubscription"`
	ColAdvanceDisbursement  string `json:"colAdvanceDisbursement"`
	ColPeriodStart          string `json:"colPeriodStart"`
	ColPeriodEnd            string `json:"colPeriodEnd"`
	ColOutcome              string `json:"colOutcome"`
	ColErrorCode            string `json:"colErrorCode"`
	EmptyTitle              string `json:"emptyTitle"`
	EmptyMessage            string `json:"emptyMessage"`
}

type ExpenseRecognitionRunBillsTabLabels struct {
	ColReference string `json:"colReference"`
	ColDate      string `json:"colDate"`
	ColAmount    string `json:"colAmount"`
	ColStatus    string `json:"colStatus"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	Hint         string `json:"hint"`
}

type ExpenseRecognitionRunRecognitionsTabLabels struct {
	ColReference  string `json:"colReference"`
	ColDate       string `json:"colDate"`
	ColAmount     string `json:"colAmount"`
	ColSourceKind string `json:"colSourceKind"`
	ColStatus     string `json:"colStatus"`
	EmptyTitle    string `json:"emptyTitle"`
	EmptyMessage  string `json:"emptyMessage"`
	Hint          string `json:"hint"`
}

// ExpenseRecognitionRunDrawerLabels holds drawer-form labels for Surface A
// (per-supplier), Surface C (per-supplier-subscription), and the
// generate-confirmation modal.
type ExpenseRecognitionRunDrawerLabels struct {
	Supplier     ExpenseRecognitionRunSupplierDrawerLabels     `json:"supplier"`
	Subscription ExpenseRecognitionRunSubscriptionDrawerLabels `json:"subscription"`
	Confirmation ExpenseRecognitionRunConfirmationLabels       `json:"confirmation"`
}

type ExpenseRecognitionRunSupplierDrawerLabels struct {
	Title                           string `json:"title"`
	SubtitleTemplate                string `json:"subtitleTemplate"`
	AsOfDateLabel                   string `json:"asOfDateLabel"`
	AsOfDateHint                    string `json:"asOfDateHint"`
	ColumnSource                    string `json:"columnSource"`
	ColumnPeriod                    string `json:"columnPeriod"`
	ColumnAmount                    string `json:"columnAmount"`
	ColumnLines                     string `json:"columnLines"`
	ColumnRemaining                 string `json:"columnRemaining"`
	GroupSubscriptionCycle          string `json:"groupSubscriptionCycle"`
	GroupAdvanceDisbursementTranche string `json:"groupAdvanceDisbursementTranche"`
	GroupNoPending                  string `json:"groupNoPending"`
	GroupCurrencyMismatch           string `json:"groupCurrencyMismatch"`
	EmptyTitle                      string `json:"emptyTitle"`
	EmptyMessage                    string `json:"emptyMessage"`
	GenerateButton                  string `json:"generateButton"`
	GenerateButtonCount             string `json:"generateButtonCount"`
	CancelButton                    string `json:"cancelButton"`
	ViewRunLink                     string `json:"viewRunLink"`
}

type ExpenseRecognitionRunSubscriptionDrawerLabels struct {
	Title                        string `json:"title"`
	SubtitleTemplate             string `json:"subtitleTemplate"`
	AsOfDateLabel                string `json:"asOfDateLabel"`
	AsOfDateHint                 string `json:"asOfDateHint"`
	ColumnPeriod                 string `json:"columnPeriod"`
	ColumnAmount                 string `json:"columnAmount"`
	ColumnLines                  string `json:"columnLines"`
	EmptyTitle                   string `json:"emptyTitle"`
	EmptyMessage                 string `json:"emptyMessage"`
	SuppressedByAdvanceTitle     string `json:"suppressedByAdvanceTitle"`
	SuppressedByAdvanceExplainer string `json:"suppressedByAdvanceExplainer"`
	ViewAdvanceLink              string `json:"viewAdvanceLink"`
	GenerateButton               string `json:"generateButton"`
	GenerateButtonCount          string `json:"generateButtonCount"`
	CancelButton                 string `json:"cancelButton"`
	ViewRunLink                  string `json:"viewRunLink"`
}

type ExpenseRecognitionRunConfirmationLabels struct {
	Title          string `json:"title"`
	BodyTemplate   string `json:"bodyTemplate"`
	NoteIdempotent string `json:"noteIdempotent"`
	ConfirmButton  string `json:"confirmButton"`
	CancelButton   string `json:"cancelButton"`
}

// ExpenseRecognitionRunStatusBadgeLabels holds badge copy for each run status.
type ExpenseRecognitionRunStatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possiblyInterrupted"`
}

// ExpenseRecognitionRunActionLabels holds labels for interactive actions on
// run rows / pages / drawer triggers.
type ExpenseRecognitionRunActionLabels struct {
	Run                      string `json:"run"`
	RunRecognitions          string `json:"runRecognitions"`
	ReRunFailed              string `json:"reRunFailed"`
	ReRunFailedComingSoon    string `json:"reRunFailedComingSoon"`
	ViewRun                  string `json:"viewRun"`
	ViewSupplier             string `json:"viewSupplier"`
	ViewSupplierSubscription string `json:"viewSupplierSubscription"`
	ViewAdvanceDisbursement  string `json:"viewAdvanceDisbursement"`
	RunAriaLabel             string `json:"runAriaLabel"`
}

// ExpenseRecognitionRunScopeKindLabels holds display labels for each Run
// scope kind: supplier (Surface A), subscription (Surface C), workspace
// (Surface B).
type ExpenseRecognitionRunScopeKindLabels struct {
	Supplier     string `json:"supplier"`
	Subscription string `json:"subscription"`
	Workspace    string `json:"workspace"`
}

// ExpenseRecognitionRunSourceKindLabels holds display labels for each
// candidate source kind on a run attempt: subscription cycle vs advance
// disbursement tranche.
type ExpenseRecognitionRunSourceKindLabels struct {
	SubscriptionCycle   string `json:"subscriptionCycle"`
	AdvanceDisbursement string `json:"advanceDisbursement"`
}

// ExpenseRecognitionRunOutcomeLabels holds display labels for per-attempt
// outcome values.
type ExpenseRecognitionRunOutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// ExpenseRecognitionRunSuppressionLabels holds copy for the
// linked-advance-suppression banner + row chip rendered on Surface A and
// Surface C when a SupplierSubscription cycle is covered by a TIME_BASED
// advance TreasuryDisbursement (Plan B decision A).
type ExpenseRecognitionRunSuppressionLabels struct {
	BannerTitle     string `json:"bannerTitle"`
	BannerMessage   string `json:"bannerMessage"`
	RowChip         string `json:"rowChip"`
	RowExplainer    string `json:"rowExplainer"`
	ViewAdvanceLink string `json:"viewAdvanceLink"`
	AriaLabel       string `json:"ariaLabel"`
}

// ExpenseRecognitionRunEmptyLabels holds empty-state copy for every surface.
type ExpenseRecognitionRunEmptyLabels struct {
	QueueTitle          string `json:"queueTitle"`
	QueueMessage        string `json:"queueMessage"`
	ListTitle           string `json:"listTitle"`
	ListMessage         string `json:"listMessage"`
	SelectionsTitle     string `json:"selectionsTitle"`
	SelectionsMessage   string `json:"selectionsMessage"`
	ResultsTitle        string `json:"resultsTitle"`
	ResultsMessage      string `json:"resultsMessage"`
	BillsTitle          string `json:"billsTitle"`
	BillsMessage        string `json:"billsMessage"`
	RecognitionsTitle   string `json:"recognitionsTitle"`
	RecognitionsMessage string `json:"recognitionsMessage"`
}

// ExpenseRecognitionRunToastLabels holds toast / notification copy.
type ExpenseRecognitionRunToastLabels struct {
	Success              string `json:"success"`
	BatchSuccess         string `json:"batchSuccess"`
	BatchSuccessMultiRun string `json:"batchSuccessMultiRun"`
	ViewRunLink          string `json:"viewRunLink"`
	GenerateFailed       string `json:"generateFailed"`
	PermissionDenied     string `json:"permissionDenied"`
}

// ExpenseRecognitionRunErrorLabels holds error-message strings for the module.
type ExpenseRecognitionRunErrorLabels struct {
	CapExceeded                  string `json:"capExceeded"`
	PermissionDenied             string `json:"permissionDenied"`
	UseCaseUnavailable           string `json:"useCaseUnavailable"`
	InvalidSelection             string `json:"invalidSelection"`
	IdempotencyConflict          string `json:"idempotencyConflict"`
	SupplierMismatch             string `json:"supplierMismatch"`
	WorkspaceMismatch            string `json:"workspaceMismatch"`
	TamperedPeriod               string `json:"tamperedPeriod"`
	RunAllMatchingNotImplemented string `json:"runAllMatchingNotImplemented"`
	CrossWorkspace               string `json:"crossWorkspace"`
	MissingCostPlan              string `json:"missingCostPlan"`
	NoPendingPeriods             string `json:"noPendingPeriods"`
	CurrencyMismatch             string `json:"currencyMismatch"`
	GenerationFailed             string `json:"generationFailed"`
	SuppressedByAdvance          string `json:"suppressedByAdvance"`
}

// DefaultExpenseRecognitionRunLabels returns ExpenseRecognitionRunLabels with
// sensible English defaults. Lyngua overlays via the "expenseRecognitionRun"
// root key in general/expense_recognition_run.json.
func DefaultExpenseRecognitionRunLabels() ExpenseRecognitionRunLabels {
	return ExpenseRecognitionRunLabels{
		AppLabel: "Expense Run",
		Labels: ExpenseRecognitionRunEntityLabels{
			NameSingular: "Expense Run",
			NamePlural:   "Expense Runs",
			ModuleTitle:  "Expense Recognition Run",
		},
		Page: ExpenseRecognitionRunPageLabels{
			QueueTitle:    "Expense Run Queue",
			QueueSubtitle: "Suppliers with pending subscription cycles or advance disbursement tranches ready to recognize",
			ListTitle:     "Expense Runs",
			ListSubtitle:  "History of expense run batches",
			DetailTitle:   "Expense Run",
		},
		Buttons: ExpenseRecognitionRunButtonLabels{
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
		Search: ExpenseRecognitionRunSearchLabels{
			Placeholder: "Search expense runs",
		},
		Filters: ExpenseRecognitionRunFilterLabels{
			AsOfDate: "As of date",
			Supplier: "Supplier",
			Status:   "Status",
			Pending:  "Pending",
			Complete: "Complete",
			Failed:   "Failed",
		},
		Columns: ExpenseRecognitionRunColumnLabels{
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
		Queue: ExpenseRecognitionRunQueueLabels{
			Title:         "Expense Run Queue",
			Subtitle:      "Suppliers with pending subscription cycles or advance disbursement tranches ready to recognize",
			AsOfDateLabel: "As of date",
			Columns: ExpenseRecognitionRunQueueColumnLabels{
				Supplier:             "Supplier",
				Subscriptions:        "Subscriptions",
				AdvanceDisbursements: "Advance Disbursements",
				PendingPeriods:       "Pending periods",
				Total:                "Total",
				Currency:             "Currency",
				Actions:              "Actions",
				Run:                  "Run",
			},
			Empty: ExpenseRecognitionRunQueueEmptyLabels{
				Title:   "Nothing to recognize",
				Message: "No suppliers have pending subscription cycles or advance disbursement tranches as of this date.",
			},
			Bulk: ExpenseRecognitionRunQueueBulkLabels{
				RunSelected:        "Run for selected",
				RunAllMatching:     "Run for all matching",
				CapExceededMessage: "Capped at 50 suppliers per batch. Narrow the filter to run the rest.",
			},
		},
		List: ExpenseRecognitionRunListLabels{
			Title:    "Expense Runs",
			Subtitle: "History of expense run batches",
			Columns: ExpenseRecognitionRunListColumnLabels{
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
			Empty: ExpenseRecognitionRunListEmptyLabels{
				Pending: ExpenseRecognitionRunListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "Runs in flight will appear here.",
				},
				Complete: ExpenseRecognitionRunListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "Completed runs will appear here.",
				},
				Failed: ExpenseRecognitionRunListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "Failed runs will appear here.",
				},
			},
			Filters: ExpenseRecognitionRunListFilterLabels{
				Pending:  "Pending",
				Complete: "Complete",
				Failed:   "Failed",
			},
		},
		Detail: ExpenseRecognitionRunDetailLabels{
			Title: "Expense Run",
			Tabs: ExpenseRecognitionRunDetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Bills:        "Draft Bills",
				Recognitions: "Recognitions",
				AuditHistory: "Audit History",
				Attachments:  "Attachments",
			},
			TabHints: ExpenseRecognitionRunDetailTabHintLabels{
				BillsHint:        "Draft Expenditure rows created by this run. AP edits or cancels them when actual vendor bills arrive.",
				RecognitionsHint: "ExpenseRecognition rows created or amortized by this run.",
			},
			Summary: ExpenseRecognitionRunSummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some recognitions may be missing.",
			},
			Selections: ExpenseRecognitionRunSelectionsTabLabels{
				ColSource:               "Source",
				ColSupplierSubscription: "Supplier Subscription",
				ColAdvanceDisbursement:  "Advance Disbursement",
				ColPeriodStart:          "Period start",
				ColPeriodEnd:            "Period end",
				ColPeriodMarker:         "Period marker",
				EmptyTitle:              "No selections",
				EmptyMessage:            "This run has no attempt records.",
			},
			Results: ExpenseRecognitionRunResultsTabLabels{
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
			Bills: ExpenseRecognitionRunBillsTabLabels{
				ColReference: "Reference",
				ColDate:      "Date",
				ColAmount:    "Amount",
				ColStatus:    "Status",
				EmptyTitle:   "No draft bills",
				EmptyMessage: "No draft Expenditure rows were created by this run.",
				Hint:         "Draft Expenditure rows created by this run. AP edits or cancels them when actual vendor bills arrive.",
			},
			Recognitions: ExpenseRecognitionRunRecognitionsTabLabels{
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
		Drawer: ExpenseRecognitionRunDrawerLabels{
			Supplier: ExpenseRecognitionRunSupplierDrawerLabels{
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
			Subscription: ExpenseRecognitionRunSubscriptionDrawerLabels{
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
			Confirmation: ExpenseRecognitionRunConfirmationLabels{
				Title:          "Confirm Generate",
				BodyTemplate:   "Generate {{.Count}} recognitions for {{.ScopeName}} as of {{.AsOfDate}}?",
				NoteIdempotent: "Re-running the same period for the same source skips already-recognized rows; no duplicates will be created.",
				ConfirmButton:  "Generate",
				CancelButton:   "Cancel",
			},
		},
		StatusBadges: ExpenseRecognitionRunStatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		Actions: ExpenseRecognitionRunActionLabels{
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
		ScopeKind: ExpenseRecognitionRunScopeKindLabels{
			Supplier:     "Supplier",
			Subscription: "Supplier Subscription",
			Workspace:    "Workspace",
		},
		SourceKind: ExpenseRecognitionRunSourceKindLabels{
			SubscriptionCycle:   "Subscription cycle",
			AdvanceDisbursement: "Advance disbursement",
		},
		AttemptOutcome: ExpenseRecognitionRunOutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Outcome: ExpenseRecognitionRunOutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		LinkedAdvanceSuppression: ExpenseRecognitionRunSuppressionLabels{
			BannerTitle:     "Cycle suppressed by linked advance",
			BannerMessage:   "An advance disbursement covers this subscription cycle. Recognition flows through the advance amortization.",
			RowChip:         "Suppressed",
			RowExplainer:    "Covered by linked advance disbursement",
			ViewAdvanceLink: "View advance disbursement",
			AriaLabel:       "Row suppressed by linked advance disbursement",
		},
		Empty: ExpenseRecognitionRunEmptyLabels{
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
		Toast: ExpenseRecognitionRunToastLabels{
			Success:              "Recognized {{.Created}} of {{.Total}} ({{.Skipped}} skipped, {{.Errored}} errored)",
			BatchSuccess:         "Expense batch run — {{.Created}} created, {{.Skipped}} skipped, {{.Errored}} failed.",
			BatchSuccessMultiRun: "Ran {{.RunCount}} expense runs ({{.Created}} created, {{.Skipped}} skipped, {{.Errored}} errored)",
			ViewRunLink:          "View run",
			GenerateFailed:       "Failed to generate expense run.",
			PermissionDenied:     "You do not have permission to run expense recognition.",
		},
		Errors: ExpenseRecognitionRunErrorLabels{
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

// ---------------------------------------------------------------------------
// Purchase Order labels
// ---------------------------------------------------------------------------

// PurchaseOrderErrorLabels holds error messages for the purchase order action handlers.
type PurchaseOrderErrorLabels struct {
	NoPermission string `json:"noPermission"`
}

// PurchaseOrderLabels holds all translatable strings for the purchase order module.
type PurchaseOrderLabels struct {
	Labels    PurchaseOrderLabelNames     `json:"labels"`
	Page      PurchaseOrderPageLabels     `json:"page"`
	Buttons   PurchaseOrderButtonLabels   `json:"buttons"`
	Columns   PurchaseOrderColumnLabels   `json:"columns"`
	Empty     PurchaseOrderEmptyLabels    `json:"empty"`
	Form      PurchaseOrderFormLabels     `json:"form"`
	Status    PurchaseOrderStatusLabels   `json:"status"`
	POTypes   PurchaseOrderPOTypeLabels   `json:"poTypes"`
	LineTypes PurchaseOrderLineTypeLabels `json:"lineTypes"`
	Actions   PurchaseOrderActionLabels   `json:"actions"`
	Bulk      PurchaseOrderBulkLabels     `json:"bulkActions"`
	Detail    PurchaseOrderDetailLabels   `json:"detail"`
	LineItems PurchaseOrderLineItemLabels `json:"lineItems"`
	Receipt   PurchaseOrderReceiptLabels  `json:"receipt"`
	Errors    PurchaseOrderErrorLabels    `json:"errors"`
}

type PurchaseOrderLabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"namePlural"`
	LineItem       string `json:"lineItem"`
	LineItemPlural string `json:"lineItemPlural"`
}

type PurchaseOrderPageLabels struct {
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

type PurchaseOrderButtonLabels struct {
	Add         string `json:"add"`
	AddLineItem string `json:"addLineItem"`
}

type PurchaseOrderColumnLabels struct {
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

type PurchaseOrderEmptyLabels struct {
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

type PurchaseOrderFormLabels struct {
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

type PurchaseOrderStatusLabels struct {
	Draft             string `json:"draft"`
	PendingApproval   string `json:"pending_approval"`
	Approved          string `json:"approved"`
	PartiallyReceived string `json:"partially_received"`
	FullyReceived     string `json:"fully_received"`
	Billed            string `json:"billed"`
	Closed            string `json:"closed"`
	Cancelled         string `json:"cancelled"`
}

type PurchaseOrderPOTypeLabels struct {
	Standard string `json:"standard"`
	Blanket  string `json:"blanket"`
	Contract string `json:"contract"`
}

type PurchaseOrderLineTypeLabels struct {
	Goods   string `json:"goods"`
	Service string `json:"service"`
	Expense string `json:"expense"`
}

type PurchaseOrderActionLabels struct {
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

type PurchaseOrderBulkLabels struct {
	Delete  string `json:"delete"`
	Approve string `json:"approve"`
	Close   string `json:"close"`
}

// PurchaseOrderDetailLabels holds translatable strings for the PO detail page.
type PurchaseOrderDetailLabels struct {
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

// PurchaseOrderLineItemLabels holds translatable strings for the PO line item drawer form.
type PurchaseOrderLineItemLabels struct {
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

// PurchaseOrderReceiptLabels holds translatable strings for the confirm receipt drawer form.
type PurchaseOrderReceiptLabels struct {
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

// ---------------------------------------------------------------------------
// Collection labels (money IN — payment collections, receivables)
// ---------------------------------------------------------------------------

// CollectionLabels holds all translatable strings for the collection module.
type CollectionLabels struct {
	Page      CollectionPageLabels    `json:"page"`
	Buttons   CollectionButtonLabels  `json:"buttons"`
	Columns   CollectionColumnLabels  `json:"columns"`
	Empty     CollectionEmptyLabels   `json:"empty"`
	Form      CollectionFormLabels    `json:"form"`
	Actions   CollectionActionLabels  `json:"actions"`
	Bulk      CollectionBulkLabels    `json:"bulkActions"`
	Detail    CollectionDetailLabels  `json:"detail"`
	Status    CollectionStatusLabels  `json:"status"`
	Confirm   CollectionConfirmLabels `json:"confirm"`
	Errors    CollectionErrorLabels   `json:"errors"`
	Dashboard CashDashboardLabels     `json:"dashboard"`
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

type CollectionPageLabels struct {
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

type CollectionButtonLabels struct {
	AddCollection string `json:"addCollection"`
}

type CollectionColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
}

type CollectionEmptyLabels struct {
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	CompletedTitle   string `json:"completedTitle"`
	CompletedMessage string `json:"completedMessage"`
	FailedTitle      string `json:"failedTitle"`
	FailedMessage    string `json:"failedMessage"`
}

type CollectionFormLabels struct {
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

type CollectionActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"markComplete"`
	Reactivate   string `json:"reactivate"`
}

type CollectionBulkLabels struct {
	Delete string `json:"delete"`
}

type CollectionDetailLabels struct {
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

type CollectionStatusLabels struct {
	Pending   string `json:"pending"`
	Completed string `json:"completed"`
	Failed    string `json:"failed"`
}

type CollectionConfirmLabels struct {
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

type CollectionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
}

// DefaultCollectionLabels returns CollectionLabels with sensible English defaults.
func DefaultCollectionLabels() CollectionLabels {
	return CollectionLabels{
		Page: CollectionPageLabels{
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
		Buttons: CollectionButtonLabels{
			AddCollection: "Add Collection",
		},
		Columns: CollectionColumnLabels{
			Reference: "Reference",
			Customer:  "Customer",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
		},
		Empty: CollectionEmptyLabels{
			PendingTitle:     "No pending collections",
			PendingMessage:   "No pending collections to display.",
			CompletedTitle:   "No completed collections",
			CompletedMessage: "No completed collections to display.",
			FailedTitle:      "No failed collections",
			FailedMessage:    "No failed collections to display.",
		},
		Form: CollectionFormLabels{
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
		Actions: CollectionActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			MarkComplete: "Mark Complete",
			Reactivate:   "Reactivate",
		},
		Bulk: CollectionBulkLabels{
			Delete: "Delete Selected",
		},
		Detail: CollectionDetailLabels{
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
		Status: CollectionStatusLabels{
			Pending:   "Pending",
			Completed: "Completed",
			Failed:    "Failed",
		},
		Confirm: CollectionConfirmLabels{
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
		Errors: CollectionErrorLabels{
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

// ---------------------------------------------------------------------------
// Disbursement labels (money OUT — payments, refunds, payouts)
// ---------------------------------------------------------------------------

// DisbursementLabels holds all translatable strings for the disbursement module.
type DisbursementLabels struct {
	Page    DisbursementPageLabels    `json:"page"`
	Buttons DisbursementButtonLabels  `json:"buttons"`
	Columns DisbursementColumnLabels  `json:"columns"`
	Empty   DisbursementEmptyLabels   `json:"empty"`
	Form    DisbursementFormLabels    `json:"form"`
	Actions DisbursementActionLabels  `json:"actions"`
	Bulk    DisbursementBulkLabels    `json:"bulkActions"`
	Detail  DisbursementDetailLabels  `json:"detail"`
	Status  DisbursementStatusLabels  `json:"status"`
	Confirm DisbursementConfirmLabels `json:"confirm"`
	Errors  DisbursementErrorLabels   `json:"errors"`
}

type DisbursementPageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingPending   string `json:"headingPending"`
	HeadingApproved  string `json:"headingApproved"`
	HeadingPaid      string `json:"headingPaid"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionPending   string `json:"captionPending"`
	CaptionApproved  string `json:"captionApproved"`
	CaptionPaid      string `json:"captionPaid"`
	CaptionCancelled string `json:"captionCancelled"`
	Dashboard        string `json:"dashboard"`
}

type DisbursementButtonLabels struct {
	AddDisbursement string `json:"addDisbursement"`
}

type DisbursementColumnLabels struct {
	Reference string `json:"reference"`
	Payee     string `json:"payee"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
	Category  string `json:"category"`
}

type DisbursementEmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	ApprovedTitle    string `json:"approvedTitle"`
	ApprovedMessage  string `json:"approvedMessage"`
	PaidTitle        string `json:"paidTitle"`
	PaidMessage      string `json:"paidMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type DisbursementFormLabels struct {
	Payee                   string `json:"payee"`
	PayeePlaceholder        string `json:"payeePlaceholder"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"referencePlaceholder"`
	PaymentMethod           string `json:"paymentMethod"`
	Category                string `json:"category"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
	ApprovedBy              string `json:"approvedBy"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	MethodCash              string `json:"methodCash"`
	MethodBankTransfer      string `json:"methodBankTransfer"`
	MethodCheck             string `json:"methodCheck"`
	MethodGCash             string `json:"methodGCash"`
	MethodOther             string `json:"methodOther"`
	StatusDraft             string `json:"statusDraft"`
	StatusPending           string `json:"statusPending"`
	StatusApproved          string `json:"statusApproved"`
	StatusPaid              string `json:"statusPaid"`
	StatusCancelled         string `json:"statusCancelled"`
	TypeSupplierPayment     string `json:"typeSupplierPayment"`
	TypePayroll             string `json:"typePayroll"`
	TypeRent                string `json:"typeRent"`
	TypeUtilities           string `json:"typeUtilities"`
	TypeOther               string `json:"typeOther"`
	ApproverNamePlaceholder string `json:"approverNamePlaceholder"`
	LinkToBill              string `json:"linkToBill"`
	NoBillOption            string `json:"noBillOption"`

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo     string `json:"referenceInfo"`
	DateInfo          string `json:"dateInfo"`
	PayeeInfo         string `json:"payeeInfo"`
	AmountInfo        string `json:"amountInfo"`
	CurrencyInfo      string `json:"currencyInfo"`
	PaymentMethodInfo string `json:"paymentMethodInfo"`
	StatusInfo        string `json:"statusInfo"`
	CategoryInfo      string `json:"categoryInfo"`
	ApprovedByInfo    string `json:"approvedByInfo"`
	NotesInfo         string `json:"notesInfo"`
}

type DisbursementActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Approve    string `json:"approve"`
	MarkPaid   string `json:"markPaid"`
	Cancel     string `json:"cancel"`
	Submit     string `json:"submit"`
	Reactivate string `json:"reactivate"`
}

type DisbursementBulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"markPaid"`
}

type DisbursementDetailLabels struct {
	PageTitle         string `json:"pageTitle"`
	TitlePrefix       string `json:"titlePrefix"`
	PaymentInfo       string `json:"paymentInfo"`
	Payee             string `json:"payee"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	Method            string `json:"method"`
	Category          string `json:"category"`
	Reference         string `json:"reference"`
	ApprovedBy        string `json:"approvedBy"`
	Notes             string `json:"notes"`
	TabBasicInfo      string `json:"tabBasicInfo"`
	TabAttachments    string `json:"tabAttachments"`
	TabAuditTrail     string `json:"tabAuditTrail"`
	TabAuditHistory   string `json:"tabAuditHistory"`
	AuditAction       string `json:"auditAction"`
	AuditUser         string `json:"auditUser"`
	AuditEmptyTitle   string `json:"auditEmptyTitle"`
	AuditEmptyMessage string `json:"auditEmptyMessage"`
}

type DisbursementStatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
}

type DisbursementConfirmLabels struct {
	Submit                string `json:"submit"`
	SubmitMessage         string `json:"submitMessage"`
	Approve               string `json:"approve"`
	ApproveMessage        string `json:"approveMessage"`
	MarkPaid              string `json:"markPaid"`
	MarkPaidMessage       string `json:"markPaidMessage"`
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancelMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkSubmit            string `json:"bulkSubmit"`
	BulkSubmitMessage     string `json:"bulkSubmitMessage"`
	BulkApprove           string `json:"bulkApprove"`
	BulkApproveMessage    string `json:"bulkApproveMessage"`
	BulkMarkPaid          string `json:"bulkMarkPaid"`
	BulkMarkPaidMessage   string `json:"bulkMarkPaidMessage"`
	BulkCancel            string `json:"bulkCancel"`
	BulkCancelMessage     string `json:"bulkCancelMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type DisbursementErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	InvalidFormData   string `json:"invalidFormData"`
	NotFound          string `json:"notFound"`
	IDRequired        string `json:"idRequired"`
	NoIDsProvided     string `json:"noIDsProvided"`
	InvalidStatus     string `json:"invalidStatus"`
	InvalidTransition string `json:"invalidTransition"`
}

// DefaultDisbursementLabels returns DisbursementLabels with sensible English defaults.
func DefaultDisbursementLabels() DisbursementLabels {
	return DisbursementLabels{
		Page: DisbursementPageLabels{
			Heading:          "Disbursements",
			HeadingDraft:     "Draft Disbursements",
			HeadingPending:   "Pending Disbursements",
			HeadingApproved:  "Approved Disbursements",
			HeadingPaid:      "Paid Disbursements",
			HeadingCancelled: "Cancelled Disbursements",
			Caption:          "Manage disbursements and payouts",
			CaptionDraft:     "Draft disbursements awaiting submission",
			CaptionPending:   "Disbursements awaiting approval",
			CaptionApproved:  "Approved disbursements ready for payment",
			CaptionPaid:      "Completed disbursement payments",
			CaptionCancelled: "Cancelled disbursements",
			Dashboard:        "Disbursements Dashboard",
		},
		Buttons: DisbursementButtonLabels{
			AddDisbursement: "Add Disbursement",
		},
		Columns: DisbursementColumnLabels{
			Reference: "Reference",
			Payee:     "Payee",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
			Category:  "Category",
		},
		Empty: DisbursementEmptyLabels{
			DraftTitle:       "No draft disbursements",
			DraftMessage:     "No draft disbursements to display.",
			PendingTitle:     "No pending disbursements",
			PendingMessage:   "No pending disbursements to display.",
			ApprovedTitle:    "No approved disbursements",
			ApprovedMessage:  "No approved disbursements to display.",
			PaidTitle:        "No paid disbursements",
			PaidMessage:      "No paid disbursements to display.",
			CancelledTitle:   "No cancelled disbursements",
			CancelledMessage: "No cancelled disbursements to display.",
		},
		Form: DisbursementFormLabels{
			Payee:                   "Payee",
			PayeePlaceholder:        "Enter payee name",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. DISB-001",
			PaymentMethod:           "Payment Method",
			Category:                "Category",
			Status:                  "Status",
			Notes:                   "Notes",
			NotesPlaceholder:        "Additional notes...",
			ApprovedBy:              "Approved By",
			AmountPlaceholder:       "0.00",
			CurrencyPlaceholder:     "PHP",
			MethodCash:              "Cash",
			MethodBankTransfer:      "Bank Transfer",
			MethodCheck:             "Check",
			MethodGCash:             "GCash",
			MethodOther:             "Other",
			StatusDraft:             "Draft",
			StatusPending:           "Pending",
			StatusApproved:          "Approved",
			StatusPaid:              "Paid",
			StatusCancelled:         "Cancelled",
			TypeSupplierPayment:     "Supplier Payment",
			TypePayroll:             "Payroll",
			TypeRent:                "Rent",
			TypeUtilities:           "Utilities",
			TypeOther:               "Other",
			ApproverNamePlaceholder: "Approver name",
			LinkToBill:              "Link to Bill",
			NoBillOption:            "— No Bill —",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ReferenceInfo:     "Unique reference number for this disbursement.",
			DateInfo:          "Date the disbursement was issued.",
			PayeeInfo:         "Name of the recipient (supplier, payroll, etc.).",
			AmountInfo:        "Total amount disbursed (in centavos; displayed as amount ÷ 100).",
			CurrencyInfo:      "Currency of the disbursed amount.",
			PaymentMethodInfo: "How the payment was made.",
			StatusInfo:        "Current state of this disbursement.",
			CategoryInfo:      "Type of disbursement for categorisation and reporting.",
			ApprovedByInfo:    "Name of the person who authorised this disbursement.",
			NotesInfo:         "Internal remarks — not shown on supplier-facing documents.",
		},
		Actions: DisbursementActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Approve:    "Approve",
			MarkPaid:   "Mark as Paid",
			Cancel:     "Cancel",
			Submit:     "Submit",
			Reactivate: "Reactivate",
		},
		Bulk: DisbursementBulkLabels{
			Delete:   "Delete Selected",
			Approve:  "Approve Selected",
			MarkPaid: "Mark Selected as Paid",
		},
		Detail: DisbursementDetailLabels{
			PageTitle:         "Disbursement Details",
			TitlePrefix:       "Disbursement #",
			PaymentInfo:       "Payment Information",
			Payee:             "Payee",
			Date:              "Date",
			Amount:            "Amount",
			Currency:          "Currency",
			Status:            "Status",
			Method:            "Payment Method",
			Category:          "Category",
			Reference:         "Reference",
			ApprovedBy:        "Approved By",
			Notes:             "Notes",
			TabBasicInfo:      "Basic Info",
			TabAttachments:    "Attachments",
			TabAuditTrail:     "Audit Trail",
			TabAuditHistory:   "History",
			AuditAction:       "Action",
			AuditUser:         "User",
			AuditEmptyTitle:   "No audit records",
			AuditEmptyMessage: "No audit trail entries yet.",
		},
		Status: DisbursementStatusLabels{
			Draft:     "Draft",
			Pending:   "Pending",
			Approved:  "Approved",
			Paid:      "Paid",
			Cancelled: "Cancelled",
		},
		Confirm: DisbursementConfirmLabels{
			Submit:                "Submit",
			SubmitMessage:         "Are you sure you want to submit {{count}} disbursement(s)?",
			Approve:               "Approve",
			ApproveMessage:        "Are you sure you want to approve {{count}} disbursement(s)?",
			MarkPaid:              "Mark as Paid",
			MarkPaidMessage:       "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			Cancel:                "Cancel",
			CancelMessage:         "Are you sure you want to cancel {{count}} disbursement(s)?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate {{count}} disbursement(s)?",
			Delete:                "Delete",
			DeleteMessage:         "Are you sure you want to delete {{count}} disbursement(s)?",
			BulkSubmit:            "Submit Disbursements",
			BulkSubmitMessage:     "Are you sure you want to submit {{count}} disbursement(s)?",
			BulkApprove:           "Approve Disbursements",
			BulkApproveMessage:    "Are you sure you want to approve {{count}} disbursement(s)?",
			BulkMarkPaid:          "Mark as Paid",
			BulkMarkPaidMessage:   "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			BulkCancel:            "Cancel Disbursements",
			BulkCancelMessage:     "Are you sure you want to cancel {{count}} disbursement(s)?",
			BulkReactivate:        "Reactivate Disbursements",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} disbursement(s)?",
			BulkDelete:            "Delete Disbursements",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} disbursement(s)?",
		},
		Errors: DisbursementErrorLabels{
			PermissionDenied:  "Permission denied",
			InvalidFormData:   "Invalid form data",
			NotFound:          "Disbursement not found",
			IDRequired:        "Disbursement ID is required",
			NoIDsProvided:     "No disbursement IDs provided",
			InvalidStatus:     "Invalid target status",
			InvalidTransition: "Cannot transition from %s to %s",
		},
	}
}

// ---------------------------------------------------------------------------
// SupplierContract labels  (P3a)
// ---------------------------------------------------------------------------

// SupplierContractLabels holds all translatable strings for the supplier_contract module.
type SupplierContractLabels struct {
	Page               SupplierContractPageLabels              `json:"page"`
	Columns            SupplierContractColumnLabels            `json:"columns"`
	Tabs               SupplierContractTabLabels               `json:"tabs"`
	Detail             SupplierContractDetailLabels            `json:"detail"`
	Lines              SupplierContractLineLabels              `json:"lines"`
	LinkedPOs          SupplierContractLinkedPOLabels          `json:"linkedPos"`
	LinkedExpenditures SupplierContractLinkedExpenditureLabels `json:"linkedExpenditures"`
	Form               SupplierContractFormLabels              `json:"form"`
	Empty              SupplierContractEmptyLabels             `json:"empty"`
}

type SupplierContractPageLabels struct {
	Heading           string `json:"heading"`
	HeadingDraft      string `json:"headingDraft"`
	HeadingActive     string `json:"headingActive"`
	HeadingExpiring   string `json:"headingExpiring"`
	HeadingExpired    string `json:"headingExpired"`
	HeadingTerminated string `json:"headingTerminated"`
	Caption           string `json:"caption"`
	AddButton         string `json:"addButton"`
	DetailSubtitle    string `json:"detailSubtitle"`
}

type SupplierContractColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	Validity  string `json:"validity"`
	Committed string `json:"committed"`
	Released  string `json:"released"`
	Billed    string `json:"billed"`
	Remaining string `json:"remaining"`
}

type SupplierContractTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedPOs           string `json:"linkedPos"`
	LinkedExpenditures  string `json:"linkedExpenditures"`
	PriceSchedules      string `json:"priceSchedules"`
	Activity            string `json:"activity"`
	ActivityEmpty       string `json:"activityEmpty"`
	PriceSchedulesEmpty string `json:"priceSchedulesEmpty"`
}

type SupplierContractDetailLabels struct {
	InfoSection     string `json:"infoSection"`
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Status          string `json:"status"`
	Supplier        string `json:"supplier"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	AutoRenew       string `json:"autoRenew"`
	Currency        string `json:"currency"`
	CommittedAmount string `json:"committedAmount"`
	ReleasedAmount  string `json:"releasedAmount"`
	BilledAmount    string `json:"billedAmount"`
	RemainingAmount string `json:"remainingAmount"`
	Notes           string `json:"notes"`
	TabAttachments  string `json:"tabAttachments"`
}

type SupplierContractLineLabels struct {
	// Column labels
	Description  string `json:"description"`
	LineType     string `json:"lineType"`
	Quantity     string `json:"quantity"`
	UnitPrice    string `json:"unitPrice"`
	Total        string `json:"total"`
	Treatment    string `json:"treatment"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	AddLine      string `json:"addLine"`

	// Enum label values for treatment
	TreatmentRecurring         string `json:"treatmentRecurring"`
	TreatmentOneTime           string `json:"treatmentOneTime"`
	TreatmentUsageBased        string `json:"treatmentUsageBased"`
	TreatmentMinimumCommitment string `json:"treatmentMinimumCommitment"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription               string `json:"formDescription"`
	FormDescriptionPlaceholder    string `json:"formDescriptionPlaceholder"`
	FormLineType                  string `json:"formLineType"`
	FormLineTypeInfo              string `json:"formLineTypeInfo"`
	FormTreatment                 string `json:"formTreatment"`
	FormTreatmentInfo             string `json:"formTreatmentInfo"`
	FormProduct                   string `json:"formProduct"`
	FormProductPlaceholder        string `json:"formProductPlaceholder"`
	FormQuantity                  string `json:"formQuantity"`
	FormQuantityInfo              string `json:"formQuantityInfo"`
	FormUnitPrice                 string `json:"formUnitPrice"`
	FormUnitPriceInfo             string `json:"formUnitPriceInfo"`
	FormExpenseAccount            string `json:"formExpenseAccount"`
	FormExpenseAccountPlaceholder string `json:"formExpenseAccountPlaceholder"`
	FormStartDate                 string `json:"formStartDate"`
	FormStartDateHint             string `json:"formStartDateHint"`
	FormEndDate                   string `json:"formEndDate"`
	FormLineNumber                string `json:"formLineNumber"`
}

type SupplierContractLinkedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

type SupplierContractLinkedExpenditureLabels struct {
	Reference    string `json:"reference"`
	Status       string `json:"status"`
	Amount       string `json:"amount"`
	Date         string `json:"date"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// SupplierContractFormLabels holds all form-level labels for the drawer form.
type SupplierContractFormLabels struct {
	// Section headers (5-section parity layout)
	SectionIdentity       string `json:"sectionIdentity"`
	SectionValidity       string `json:"sectionValidity"`
	SectionMoney          string `json:"sectionMoney"`
	SectionCategorization string `json:"sectionCategorization"`
	SectionOthers         string `json:"sectionOthers"`

	// §1 Identity
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"namePlaceholder"`
	NameInfo                  string `json:"nameInfo"`
	ContractNumber            string `json:"contractNumber"`
	ContractNumberPlaceholder string `json:"contractNumberPlaceholder"`
	Kind                      string `json:"kind"`
	KindInfo                  string `json:"kindInfo"`
	KindSubscription          string `json:"kindSubscription"`
	KindRetainer              string `json:"kindRetainer"`
	KindLease                 string `json:"kindLease"`
	KindUtility               string `json:"kindUtility"`
	KindFramework             string `json:"kindFramework"`
	KindBlanket               string `json:"kindBlanket"`
	KindOneTime               string `json:"kindOneTime"`
	KindOther                 string `json:"kindOther"`
	Supplier                  string `json:"supplier"`
	SupplierPlaceholder       string `json:"supplierPlaceholder"`
	SupplierInfo              string `json:"supplierInfo"`

	// §2 Validity & Recurrence
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	EndDateHint           string `json:"endDateHint"`
	BillingCycleValue     string `json:"billingCycleValue"`
	BillingCycleUnit      string `json:"billingCycleUnit"`
	BillingCycleInfo      string `json:"billingCycleInfo"`
	CycleUnitDay          string `json:"cycleUnitDay"`
	CycleUnitWeek         string `json:"cycleUnitWeek"`
	CycleUnitMonth        string `json:"cycleUnitMonth"`
	CycleUnitYear         string `json:"cycleUnitYear"`
	AutoRenew             string `json:"autoRenew"`
	RenewalNoticeDays     string `json:"renewalNoticeDays"`
	RenewalNoticeDaysHint string `json:"renewalNoticeDaysHint"`

	// §3 Money & Approval
	Currency               string `json:"currency"`
	CurrencyInfo           string `json:"currencyInfo"`
	Status                 string `json:"status"`
	StatusInfo             string `json:"statusInfo"`
	StatusDraft            string `json:"statusDraft"`
	StatusRequested        string `json:"statusRequested"`
	StatusPendingApproval  string `json:"statusPendingApproval"`
	StatusApproved         string `json:"statusApproved"`
	StatusActive           string `json:"statusActive"`
	StatusExpiring         string `json:"statusExpiring"`
	StatusSuspended        string `json:"statusSuspended"`
	StatusExpired          string `json:"statusExpired"`
	StatusTerminated       string `json:"statusTerminated"`
	StatusRejected         string `json:"statusRejected"`
	CommittedAmount        string `json:"committedAmount"`
	CommittedAmountInfo    string `json:"committedAmountInfo"`
	CycleAmount            string `json:"cycleAmount"`
	CycleAmountHint        string `json:"cycleAmountHint"`
	PaymentTerm            string `json:"paymentTerm"`
	PaymentTermPlaceholder string `json:"paymentTermPlaceholder"`
	ApprovedBy             string `json:"approvedBy"`
	ApprovedDate           string `json:"approvedDate"`

	// §4 Categorization
	ExpenditureCategory            string `json:"expenditureCategory"`
	ExpenditureCategoryPlaceholder string `json:"expenditureCategoryPlaceholder"`
	ExpenseAccount                 string `json:"expenseAccount"`
	ExpenseAccountPlaceholder      string `json:"expenseAccountPlaceholder"`
	Location                       string `json:"location"`
	LocationPlaceholder            string `json:"locationPlaceholder"`

	// §5 Others
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	Active           string `json:"active"`

	// Action buttons on detail page
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Approve   string `json:"approve"`
	Terminate string `json:"terminate"`
}

type SupplierContractEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultSupplierContractLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultSupplierContractLabels() SupplierContractLabels {
	return SupplierContractLabels{
		Page: SupplierContractPageLabels{
			Heading:           "Supplier Contracts",
			HeadingDraft:      "Draft Contracts",
			HeadingActive:     "Active Contracts",
			HeadingExpiring:   "Expiring Contracts",
			HeadingExpired:    "Expired Contracts",
			HeadingTerminated: "Terminated Contracts",
			Caption:           "Standing agreements with suppliers",
			AddButton:         "New Contract",
			DetailSubtitle:    "Contract details",
		},
		Columns: SupplierContractColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			Kind:      "Kind",
			Status:    "Status",
			Validity:  "Validity",
			Committed: "Committed",
			Released:  "Released",
			Billed:    "Billed",
			Remaining: "Remaining",
		},
		Tabs: SupplierContractTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedPOs:           "Linked POs",
			LinkedExpenditures:  "Linked Expenditures",
			PriceSchedules:      "Price Schedules",
			Activity:            "Activity",
			ActivityEmpty:       "No activity recorded yet.",
			PriceSchedulesEmpty: "No price schedules yet. Add a schedule to layer multi-year pricing on this contract.",
		},
		Detail: SupplierContractDetailLabels{
			InfoSection:     "Contract Information",
			Name:            "Name",
			Kind:            "Kind",
			Status:          "Status",
			Supplier:        "Supplier",
			StartDate:       "Start Date",
			EndDate:         "End Date",
			AutoRenew:       "Auto Renew",
			Currency:        "Currency",
			CommittedAmount: "Committed Amount",
			ReleasedAmount:  "Released Amount",
			BilledAmount:    "Billed Amount",
			RemainingAmount: "Remaining Amount",
			Notes:           "Notes",
			TabAttachments:  "Attachments",
		},
		Lines: SupplierContractLineLabels{
			Description:                   "Description",
			LineType:                      "Line Type",
			Quantity:                      "Quantity",
			UnitPrice:                     "Unit Price",
			Total:                         "Total",
			Treatment:                     "Treatment",
			EmptyTitle:                    "No lines yet",
			EmptyMessage:                  "Add a line to this contract.",
			AddLine:                       "Add Line",
			TreatmentRecurring:            "Recurring",
			TreatmentOneTime:              "One Time",
			TreatmentUsageBased:           "Usage Based",
			TreatmentMinimumCommitment:    "Minimum Commitment",
			LineTypeGoods:                 "Goods",
			LineTypeService:               "Service",
			LineTypeExpense:               "Expense",
			FormDescription:               "Description",
			FormDescriptionPlaceholder:    "e.g. Cloud hosting — 50 seats",
			FormLineType:                  "Line Type",
			FormLineTypeInfo:              "Goods = physical items; Service = intangible; Expense = direct cost",
			FormTreatment:                 "Treatment",
			FormTreatmentInfo:             "How this line is billed: recurring, one-time, usage-based, or minimum commitment",
			FormProduct:                   "Product",
			FormProductPlaceholder:        "Select a product (optional)",
			FormQuantity:                  "Quantity",
			FormQuantityInfo:              "For recurring lines, this is the per-cycle quantity.",
			FormUnitPrice:                 "Unit Price",
			FormUnitPriceInfo:             "Amount in centavos ÷ 100 for display.",
			FormExpenseAccount:            "Expense Account",
			FormExpenseAccountPlaceholder: "GL account ID",
			FormStartDate:                 "Start Date",
			FormStartDateHint:             "Leave empty to inherit from contract.",
			FormEndDate:                   "End Date",
			FormLineNumber:                "Line Number",
		},
		LinkedPOs: SupplierContractLinkedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No linked purchase orders",
			EmptyMessage: "POs created against this contract will appear here.",
		},
		LinkedExpenditures: SupplierContractLinkedExpenditureLabels{
			Reference:    "Reference",
			Status:       "Status",
			Amount:       "Amount",
			Date:         "Date",
			EmptyTitle:   "No linked expenditures",
			EmptyMessage: "Expenditures linked to this contract will appear here.",
		},
		Form: SupplierContractFormLabels{
			SectionIdentity:                "Identity Details",
			SectionValidity:                "Validity & Recurrence",
			SectionMoney:                   "Money & Approval",
			SectionCategorization:          "Categorization",
			SectionOthers:                  "Others",
			Name:                           "Contract Name",
			NamePlaceholder:                "e.g. AWS Hosting MSA 2026",
			NameInfo:                       "A short descriptive name for this contract.",
			ContractNumber:                 "Contract Number",
			ContractNumberPlaceholder:      "Supplier's reference number",
			Kind:                           "Kind",
			KindInfo:                       "Subscription = recurring time-based; Blanket = quantity-based commitment; Framework = pricing agreement only.",
			KindSubscription:               "Subscription",
			KindRetainer:                   "Retainer",
			KindLease:                      "Lease",
			KindUtility:                    "Utility",
			KindFramework:                  "Framework",
			KindBlanket:                    "Blanket",
			KindOneTime:                    "One Time",
			KindOther:                      "Other",
			Supplier:                       "Supplier",
			SupplierPlaceholder:            "Select supplier",
			SupplierInfo:                   "The vendor or service provider you are committing to.",
			StartDate:                      "Start Date",
			EndDate:                        "End Date",
			EndDateHint:                    "Leave empty for open-ended.",
			BillingCycleValue:              "Billing Cycle",
			BillingCycleUnit:               "Cycle Unit",
			BillingCycleInfo:               "How often this contract generates an expenditure (for recurring kinds).",
			CycleUnitDay:                   "Day",
			CycleUnitWeek:                  "Week",
			CycleUnitMonth:                 "Month",
			CycleUnitYear:                  "Year",
			AutoRenew:                      "Auto Renew",
			RenewalNoticeDays:              "Renewal Notice (days)",
			RenewalNoticeDaysHint:          "How many days before expiry to send a renewal reminder.",
			Currency:                       "Currency",
			CurrencyInfo:                   "ISO 4217 currency code (e.g. PHP, USD).",
			Status:                         "Status",
			StatusInfo:                     "Lifecycle stage. draft → requested → pending_approval → approved → active → expiring/expired/terminated.",
			StatusDraft:                    "Draft",
			StatusRequested:                "Requested",
			StatusPendingApproval:          "Pending Approval",
			StatusApproved:                 "Approved",
			StatusActive:                   "Active",
			StatusExpiring:                 "Expiring",
			StatusSuspended:                "Suspended",
			StatusExpired:                  "Expired",
			StatusTerminated:               "Terminated",
			StatusRejected:                 "Rejected",
			CommittedAmount:                "Committed Amount",
			CommittedAmountInfo:            "Total value committed at signing (centavos). Immutable after approval.",
			CycleAmount:                    "Cycle Amount",
			CycleAmountHint:                "Expected per-cycle charge for recurring contracts.",
			PaymentTerm:                    "Payment Term",
			PaymentTermPlaceholder:         "Select payment term",
			ApprovedBy:                     "Approved By",
			ApprovedDate:                   "Approved Date",
			ExpenditureCategory:            "Expenditure Category",
			ExpenditureCategoryPlaceholder: "Select category",
			ExpenseAccount:                 "Expense Account",
			ExpenseAccountPlaceholder:      "GL account ID",
			Location:                       "Location",
			LocationPlaceholder:            "Branch or cost center",
			Notes:                          "Notes",
			NotesPlaceholder:               "Additional notes or context",
			Active:                         "Active",
			Edit:                           "Edit",
			EditTitle:                      "Edit Supplier Contract",
			Approve:                        "Approve",
			Terminate:                      "Terminate",
		},
		Empty: SupplierContractEmptyLabels{
			Title:   "No supplier contracts",
			Message: "Create your first supplier contract to start tracking commitments.",
		},
	}
}

// ---------------------------------------------------------------------------
// ProcurementRequest labels  (P3a)
// ---------------------------------------------------------------------------

// ProcurementRequestLabels holds all translatable strings for the procurement_request module.
type ProcurementRequestLabels struct {
	Page       ProcurementRequestPageLabels      `json:"page"`
	Columns    ProcurementRequestColumnLabels    `json:"columns"`
	Tabs       ProcurementRequestTabLabels       `json:"tabs"`
	Detail     ProcurementRequestDetailLabels    `json:"detail"`
	Lines      ProcurementRequestLineLabels      `json:"lines"`
	SpawnedPOs ProcurementRequestSpawnedPOLabels `json:"spawnedPos"`
	Form       ProcurementRequestFormLabels      `json:"form"`
	Empty      ProcurementRequestEmptyLabels     `json:"empty"`

	// SPS Wave 3 — F1/F2/F3 + CRIT-3 spawn lifecycle
	Filters              ProcurementRequestFilterLabels              `json:"filters"`
	FulfillmentStrategy  ProcurementRequestFulfillmentStrategyLabels `json:"fulfillmentStrategy"`
	FulfillmentMode      ProcurementRequestFulfillmentModeLabels     `json:"fulfillmentMode"`
	FulfillmentModeHints ProcurementRequestFulfillmentModeHintLabels `json:"fulfillmentModeHints"`
	Spawn                ProcurementRequestSpawnLabels               `json:"spawn"`
	PolicyDecision       ProcurementRequestPolicyDecisionLabels      `json:"policyDecision"`
}

// ProcurementRequestFilterLabels — F3 filter chips on the list page.
type ProcurementRequestFilterLabels struct {
	All                    string `json:"all"`
	Status                 string `json:"status"`
	FulfillmentStrategy    string `json:"fulfillmentStrategy"`
	FulfillmentMode        string `json:"fulfillmentMode"`
	AnyStatus              string `json:"anyStatus"`
	AnyFulfillmentStrategy string `json:"anyFulfillmentStrategy"`
	AnyFulfillmentMode     string `json:"anyFulfillmentMode"`
}

// ProcurementRequestFulfillmentStrategyLabels — F3 strategy values for header-level rollup.
type ProcurementRequestFulfillmentStrategyLabels struct {
	UniformOutright  string `json:"uniformOutright"`
	UniformStockable string `json:"uniformStockable"`
	UniformRecurring string `json:"uniformRecurring"`
	UniformPetty     string `json:"uniformPetty"`
	Mixed            string `json:"mixed"`
	Hint             string `json:"hint"`
}

// ProcurementRequestFulfillmentModeLabels — F1 line-level mode values.
type ProcurementRequestFulfillmentModeLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// ProcurementRequestFulfillmentModeHintLabels — F1 short hints rendered under each radio choice.
type ProcurementRequestFulfillmentModeHintLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// ProcurementRequestSpawnLabels — CRIT-3 spawn lifecycle UI strings.
type ProcurementRequestSpawnLabels struct {
	StatusColumn      string `json:"statusColumn"`
	StatusPending     string `json:"statusPending"`
	StatusSpawning    string `json:"statusSpawning"`
	StatusSpawned     string `json:"statusSpawned"`
	StatusFailed      string `json:"statusFailed"`
	StatusUnspecified string `json:"statusUnspecified"`
	ModeColumn        string `json:"modeColumn"`
	SpawnedColumn     string `json:"spawnedColumn"`
	LinkPO            string `json:"linkPo"`
	LinkContract      string `json:"linkContract"`
	LinkExpenditure   string `json:"linkExpenditure"`
	NotApplicable     string `json:"notApplicable"`
	ErrorPrefix       string `json:"errorPrefix"`
	RetryButton       string `json:"retryButton"`
	RetryConfirm      string `json:"retryConfirm"`
}

// ProcurementRequestPolicyDecisionLabels — policy_decision_log section on Info tab.
type ProcurementRequestPolicyDecisionLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Toggle       string `json:"toggle"`
	EmptyMessage string `json:"emptyMessage"`
	Info         string `json:"info"`
}

type ProcurementRequestPageLabels struct {
	Heading                string `json:"heading"`
	HeadingDraft           string `json:"headingDraft"`
	HeadingSubmitted       string `json:"headingSubmitted"`
	HeadingPendingApproval string `json:"headingPendingApproval"`
	HeadingApproved        string `json:"headingApproved"`
	HeadingRejected        string `json:"headingRejected"`
	HeadingFulfilled       string `json:"headingFulfilled"`
	HeadingCancelled       string `json:"headingCancelled"`
	Caption                string `json:"caption"`
	AddButton              string `json:"addButton"`
	DetailSubtitle         string `json:"detailSubtitle"`
}

type ProcurementRequestColumnLabels struct {
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
}

type ProcurementRequestTabLabels struct {
	Info          string `json:"info"`
	Lines         string `json:"lines"`
	SpawnedPOs    string `json:"spawnedPos"`
	Activity      string `json:"activity"`
	ActivityEmpty string `json:"activityEmpty"`
}

type ProcurementRequestDetailLabels struct {
	InfoSection    string `json:"infoSection"`
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	Currency       string `json:"currency"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
	ApprovedBy     string `json:"approvedBy"`
	Justification  string `json:"justification"`
	TabAttachments string `json:"tabAttachments"`
}

type ProcurementRequestLineLabels struct {
	// Column labels
	Description         string `json:"description"`
	LineType            string `json:"lineType"`
	Quantity            string `json:"quantity"`
	EstimatedUnitPrice  string `json:"estimatedUnitPrice"`
	EstimatedTotalPrice string `json:"estimatedTotalPrice"`
	EmptyTitle          string `json:"emptyTitle"`
	EmptyMessage        string `json:"emptyMessage"`
	AddLine             string `json:"addLine"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription                    string `json:"formDescription"`
	FormDescriptionPlaceholder         string `json:"formDescriptionPlaceholder"`
	FormLineType                       string `json:"formLineType"`
	FormLineTypeInfo                   string `json:"formLineTypeInfo"`
	FormProduct                        string `json:"formProduct"`
	FormProductPlaceholder             string `json:"formProductPlaceholder"`
	FormQuantity                       string `json:"formQuantity"`
	FormQuantityInfo                   string `json:"formQuantityInfo"`
	FormEstimatedUnitPrice             string `json:"formEstimatedUnitPrice"`
	FormEstimatedUnitPriceInfo         string `json:"formEstimatedUnitPriceInfo"`
	FormEstimatedTotalPrice            string `json:"formEstimatedTotalPrice"`
	FormEstimatedTotalPriceHint        string `json:"formEstimatedTotalPriceHint"`
	FormExpenditureCategory            string `json:"formExpenditureCategory"`
	FormExpenditureCategoryPlaceholder string `json:"formExpenditureCategoryPlaceholder"`
	FormLocation                       string `json:"formLocation"`
	FormLocationPlaceholder            string `json:"formLocationPlaceholder"`
	FormLineNumber                     string `json:"formLineNumber"`

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields + PETTY hint
	FormFulfillmentMode     string `json:"formFulfillmentMode"`
	FormFulfillmentModeInfo string `json:"formFulfillmentModeInfo"`
	FormFulfillmentModeHint string `json:"formFulfillmentModeHint"`

	FormRecurringSection    string `json:"formRecurringSection"`
	FormRecurringCycleValue string `json:"formRecurringCycleValue"`
	FormRecurringCycleUnit  string `json:"formRecurringCycleUnit"`
	FormRecurringTermValue  string `json:"formRecurringTermValue"`
	FormRecurringTermUnit   string `json:"formRecurringTermUnit"`
	FormRecurringCycleHint  string `json:"formRecurringCycleHint"`
	FormRecurringTermHint   string `json:"formRecurringTermHint"`
	FormRecurringUnitDay    string `json:"formRecurringUnitDay"`
	FormRecurringUnitWeek   string `json:"formRecurringUnitWeek"`
	FormRecurringUnitMonth  string `json:"formRecurringUnitMonth"`
	FormRecurringUnitYear   string `json:"formRecurringUnitYear"`

	FormPettyHint string `json:"formPettyHint"`

	// CRIT-3 spawn lifecycle column on the lines table
	ModeBadgeColumn string `json:"modeBadgeColumn"`
}

type ProcurementRequestSpawnedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ProcurementRequestFormLabels holds all form-level labels for the drawer form.
type ProcurementRequestFormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionFinancial string `json:"sectionFinancial"`
	SectionApproval  string `json:"sectionApproval"`
	SectionOthers    string `json:"sectionOthers"`

	// §1 Identity
	RequestNumber            string `json:"requestNumber"`
	RequestNumberPlaceholder string `json:"requestNumberPlaceholder"`
	RequestNumberInfo        string `json:"requestNumberInfo"`
	RequesterUser            string `json:"requesterUser"`
	RequesterUserPlaceholder string `json:"requesterUserPlaceholder"`
	Supplier                 string `json:"supplier"`
	SupplierPlaceholder      string `json:"supplierPlaceholder"`
	SupplierHint             string `json:"supplierHint"`
	Location                 string `json:"location"`
	LocationPlaceholder      string `json:"locationPlaceholder"`

	// §2 Financial
	Currency           string `json:"currency"`
	CurrencyInfo       string `json:"currencyInfo"`
	EstimatedTotal     string `json:"estimatedTotal"`
	EstimatedTotalInfo string `json:"estimatedTotalInfo"`

	// §3 Timing & Approval
	NeededByDate               string `json:"neededByDate"`
	NeededByDateInfo           string `json:"neededByDateInfo"`
	Status                     string `json:"status"`
	StatusInfo                 string `json:"statusInfo"`
	StatusDraft                string `json:"statusDraft"`
	StatusSubmitted            string `json:"statusSubmitted"`
	StatusPendingApproval      string `json:"statusPendingApproval"`
	StatusApproved             string `json:"statusApproved"`
	StatusApprovedPendingSpawn string `json:"statusApprovedPendingSpawn"`
	StatusRejected             string `json:"statusRejected"`
	StatusFulfilled            string `json:"statusFulfilled"`
	StatusCancelled            string `json:"statusCancelled"`
	ApprovedBy                 string `json:"approvedBy"`

	// §4 Others
	Justification            string `json:"justification"`
	JustificationPlaceholder string `json:"justificationPlaceholder"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	Active                   string `json:"active"`

	// Action buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Submit    string `json:"submit"`
	Approve   string `json:"approve"`
	Reject    string `json:"reject"`
	SpawnPO   string `json:"spawnPo"`
}

type ProcurementRequestEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultProcurementRequestLabels returns English fallback labels.
func DefaultProcurementRequestLabels() ProcurementRequestLabels {
	return ProcurementRequestLabels{
		Page: ProcurementRequestPageLabels{
			Heading:                "Procurement Requests",
			HeadingDraft:           "Draft Requests",
			HeadingSubmitted:       "Submitted Requests",
			HeadingPendingApproval: "Pending Approval",
			HeadingApproved:        "Approved Requests",
			HeadingRejected:        "Rejected Requests",
			HeadingFulfilled:       "Fulfilled Requests",
			HeadingCancelled:       "Cancelled Requests",
			Caption:                "Internal purchase intent records",
			AddButton:              "New Request",
			DetailSubtitle:         "Procurement request details",
		},
		Columns: ProcurementRequestColumnLabels{
			RequestNumber:  "Request #",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
		},
		Tabs: ProcurementRequestTabLabels{
			Info:          "Info",
			Lines:         "Lines",
			SpawnedPOs:    "Spawned POs",
			Activity:      "Activity",
			ActivityEmpty: "No activity recorded yet.",
		},
		Detail: ProcurementRequestDetailLabels{
			InfoSection:    "Request Information",
			RequestNumber:  "Request Number",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			Currency:       "Currency",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
			ApprovedBy:     "Approved By",
			Justification:  "Justification",
			TabAttachments: "Attachments",
		},
		Lines: ProcurementRequestLineLabels{
			Description:                        "Description",
			LineType:                           "Line Type",
			Quantity:                           "Quantity",
			EstimatedUnitPrice:                 "Est. Unit Price",
			EstimatedTotalPrice:                "Est. Total",
			EmptyTitle:                         "No lines yet",
			EmptyMessage:                       "Add a line to this request.",
			AddLine:                            "Add Line",
			LineTypeGoods:                      "Goods",
			LineTypeService:                    "Service",
			LineTypeExpense:                    "Expense",
			FormDescription:                    "Description",
			FormDescriptionPlaceholder:         "e.g. 50 laptop units",
			FormLineType:                       "Line Type",
			FormLineTypeInfo:                   "Goods = physical items; Service = intangible; Expense = direct cost",
			FormProduct:                        "Product",
			FormProductPlaceholder:             "Select a product (optional)",
			FormQuantity:                       "Quantity",
			FormQuantityInfo:                   "Number of units requested.",
			FormEstimatedUnitPrice:             "Estimated Unit Price",
			FormEstimatedUnitPriceInfo:         "Best estimate in centavos ÷ 100.",
			FormEstimatedTotalPrice:            "Estimated Total Price",
			FormEstimatedTotalPriceHint:        "Auto-calculated. Override if needed.",
			FormExpenditureCategory:            "Expenditure Category",
			FormExpenditureCategoryPlaceholder: "Select category",
			FormLocation:                       "Location",
			FormLocationPlaceholder:            "Branch or cost center",
			FormLineNumber:                     "Line Number",
			FormFulfillmentMode:                "Fulfillment Mode",
			FormFulfillmentModeInfo:            "How this line will be sourced after approval. Drives the downstream artifact created when the request is approved.",
			FormFulfillmentModeHint:            "Pick one — the spawn cascade dispatches per-line based on this choice.",
			FormRecurringSection:               "Recurring Schedule",
			FormRecurringCycleValue:            "Cycle Every",
			FormRecurringCycleUnit:             "Cycle Unit",
			FormRecurringTermValue:             "Term Length",
			FormRecurringTermUnit:              "Term Unit",
			FormRecurringCycleHint:             "Billing/delivery cadence (e.g. every 1 month).",
			FormRecurringTermHint:              "Total contract horizon (e.g. 24 months).",
			FormRecurringUnitDay:               "Day",
			FormRecurringUnitWeek:              "Week",
			FormRecurringUnitMonth:             "Month",
			FormRecurringUnitYear:              "Year",
			FormPettyHint:                      "Petty mode auto-approves under threshold and posts a direct expenditure. No PO, no contract.",
			ModeBadgeColumn:                    "Mode",
		},
		SpawnedPOs: ProcurementRequestSpawnedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No purchase orders yet",
			EmptyMessage: "POs spawned from this request will appear here after approval.",
		},
		Form: ProcurementRequestFormLabels{
			SectionIdentity:            "Identity",
			SectionFinancial:           "Financial",
			SectionApproval:            "Timing & Approval",
			SectionOthers:              "Others",
			RequestNumber:              "Request Number",
			RequestNumberPlaceholder:   "e.g. PR-2026-001",
			RequestNumberInfo:          "A unique identifier for this procurement request.",
			RequesterUser:              "Requester",
			RequesterUserPlaceholder:   "User ID of requester",
			Supplier:                   "Supplier",
			SupplierPlaceholder:        "Select supplier (optional for RFQ)",
			SupplierHint:               "Leave empty if supplier is not yet selected (RFQ flow).",
			Location:                   "Location",
			LocationPlaceholder:        "Branch or cost center",
			Currency:                   "Currency",
			CurrencyInfo:               "ISO 4217 currency code (e.g. PHP, USD).",
			EstimatedTotal:             "Estimated Total",
			EstimatedTotalInfo:         "Best estimate of total spend (centavos ÷ 100 for display).",
			NeededByDate:               "Needed By",
			NeededByDateInfo:           "When the goods or services are required.",
			Status:                     "Status",
			StatusInfo:                 "Lifecycle stage. draft → submitted → pending_approval → approved/rejected → fulfilled/cancelled.",
			StatusDraft:                "Draft",
			StatusSubmitted:            "Submitted",
			StatusPendingApproval:      "Pending Approval",
			StatusApproved:             "Approved",
			StatusApprovedPendingSpawn: "Approved — Pending Spawn",
			StatusRejected:             "Rejected",
			StatusFulfilled:            "Fulfilled",
			StatusCancelled:            "Cancelled",
			ApprovedBy:                 "Approved By",
			Justification:              "Justification",
			JustificationPlaceholder:   "Business reason for this request",
			Notes:                      "Notes",
			NotesPlaceholder:           "Additional notes or context",
			Active:                     "Active",
			Edit:                       "Edit",
			EditTitle:                  "Edit Procurement Request",
			Submit:                     "Submit for Approval",
			Approve:                    "Approve",
			Reject:                     "Reject",
			SpawnPO:                    "Create PO",
		},
		Empty: ProcurementRequestEmptyLabels{
			Title:   "No procurement requests",
			Message: "Create a procurement request to start the approval workflow.",
		},
		Filters: ProcurementRequestFilterLabels{
			All:                    "All",
			Status:                 "Status",
			FulfillmentStrategy:    "Fulfillment",
			FulfillmentMode:        "Mode",
			AnyStatus:              "Any Status",
			AnyFulfillmentStrategy: "Any Fulfillment",
			AnyFulfillmentMode:     "Any Mode",
		},
		FulfillmentStrategy: ProcurementRequestFulfillmentStrategyLabels{
			UniformOutright:  "Uniform — Outright",
			UniformStockable: "Uniform — Stockable",
			UniformRecurring: "Uniform — Recurring",
			UniformPetty:     "Uniform — Petty",
			Mixed:            "Mixed Modes",
			Hint:             "Auto-derived from per-line fulfillment modes. Mixed = lines split across multiple modes.",
		},
		FulfillmentMode: ProcurementRequestFulfillmentModeLabels{
			Outright:  "Outright",
			Stockable: "Stockable",
			Recurring: "Recurring",
			Petty:     "Petty",
		},
		FulfillmentModeHints: ProcurementRequestFulfillmentModeHintLabels{
			Outright:  "One-shot purchase. Spawns a single purchase order on approval; no recurrence, no inventory side-effect.",
			Stockable: "Replenishment buy. Spawns a purchase order; received goods credit inventory on receipt.",
			Recurring: "Standing agreement. Spawns a supplier contract on approval; the recurrence engine emits cycle bills.",
			Petty:     "Cash-out. Spawns an expenditure directly against petty cash. No PO, no contract.",
		},
		Spawn: ProcurementRequestSpawnLabels{
			StatusColumn:      "Spawn Status",
			StatusPending:     "Pending",
			StatusSpawning:    "Spawning",
			StatusSpawned:     "Spawned",
			StatusFailed:      "Failed",
			StatusUnspecified: "—",
			ModeColumn:        "Mode",
			SpawnedColumn:     "Spawned Artifact",
			LinkPO:            "View PO line",
			LinkContract:      "View contract",
			LinkExpenditure:   "View expenditure",
			NotApplicable:     "—",
			ErrorPrefix:       "Error",
			RetryButton:       "Retry spawn",
			RetryConfirm:      "Retry spawning the downstream artifact for this line?",
		},
		PolicyDecision: ProcurementRequestPolicyDecisionLabels{
			SectionTitle: "Approval Policy Log",
			Toggle:       "Show / Hide",
			EmptyMessage: "No policy decisions logged yet.",
			Info:         "Audit trail of approval policy decisions taken on this request (auto-approve, escalation, override). Read-only.",
		},
	}
}

// LocationMap maps location slugs to display names.
var LocationMap = map[string]string{
	"ayala-central-bloc": "Ayala Central Bloc",
	"sm-city-cebu":       "SM City Cebu",
	"ayala-center-cebu":  "Ayala Center Cebu",
	"robinsons-galleria": "Robinsons Galleria",
}

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string {
	if name, ok := LocationMap[slug]; ok {
		return name
	}
	return slug
}

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app labels
// (composition surface, no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementLabels holds all translatable strings for the Procurement
// Operations composition app. Populated via lyngua (P4). These keys are
// intentionally generic so they render without overrides when lyngua has not
// yet supplied values.
type ProcurementLabels struct {
	AppLabel              string `json:"app_label"`
	DashboardTitle        string `json:"dashboard_title"`
	PendingApprovalsTitle string `json:"pending_approvals_title"`
	ExpiringTitle         string `json:"expiring_title"`
	VarianceTitle         string `json:"variance_title"`
	RecurrenceTitle       string `json:"recurrence_title"`
	RenewalsTitle         string `json:"renewals_title"`
	UtilizationTitle      string `json:"utilization_title"`
	EmptyRenewals         string `json:"empty_renewals"`
	EmptyVariance         string `json:"empty_variance"`
	EmptyUtilization      string `json:"empty_utilization"`
	EmptyRecurrence       string `json:"empty_recurrence"`
	DaysUntilExpiry       string `json:"days_until_expiry"`
	UtilizationPercent    string `json:"utilization_percent"`
	BudgetPressureLabel   string `json:"budget_pressure_label"`
}

// ---------------------------------------------------------------------------
// SPS P7 — SupplierContractPriceSchedule labels
// (mirrors lyngua/translations/en/general/supplier_contract_price_schedule.json
//  root key "supplierContractPriceSchedule")
// ---------------------------------------------------------------------------

// SupplierContractPriceScheduleLabels holds all translatable strings for the
// supplier_contract_price_schedule + child line views.
type SupplierContractPriceScheduleLabels struct {
	Labels  SupplierContractPriceScheduleNounLabels   `json:"labels"`
	Page    SupplierContractPriceSchedulePageLabels   `json:"page"`
	Buttons SupplierContractPriceScheduleButtonLabels `json:"buttons"`
	Filters SupplierContractPriceScheduleFilterLabels `json:"filters"`
	Columns SupplierContractPriceScheduleColumnLabels `json:"columns"`
	Empty   SupplierContractPriceScheduleEmptyLabels  `json:"empty"`
	Form    SupplierContractPriceScheduleFormLabels   `json:"form"`
	Status  SupplierContractPriceScheduleStatusLabels `json:"status"`
	Tabs    SupplierContractPriceScheduleTabLabels    `json:"tabs"`
	Lines   SupplierContractPriceScheduleLinesLabels  `json:"lines"`
	Detail  SupplierContractPriceScheduleDetailLabels `json:"detail"`
	Errors  SupplierContractPriceScheduleErrorLabels  `json:"errors"`
}

type SupplierContractPriceScheduleNounLabels struct {
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
	Line       string `json:"line"`
	LinePlural string `json:"linePlural"`
}

type SupplierContractPriceSchedulePageLabels struct {
	Heading           string `json:"heading"`
	Caption           string `json:"caption"`
	HeadingScheduled  string `json:"headingScheduled"`
	HeadingActive     string `json:"headingActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	HeadingCancelled  string `json:"headingCancelled"`
	TabTitle          string `json:"tabTitle"`
}

type SupplierContractPriceScheduleButtonLabels struct {
	Add       string `json:"add"`
	AddLine   string `json:"addLine"`
	Activate  string `json:"activate"`
	Supersede string `json:"supersede"`
	Cancel    string `json:"cancel"`
}

type SupplierContractPriceScheduleFilterLabels struct {
	All                 string `json:"all"`
	Status              string `json:"status"`
	AnyStatus           string `json:"anyStatus"`
	SupplierContract    string `json:"supplierContract"`
	AnySupplierContract string `json:"anySupplierContract"`
	DateRange           string `json:"dateRange"`
}

type SupplierContractPriceScheduleColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	SupplierContract string `json:"supplierContract"`
	SequenceNumber   string `json:"sequenceNumber"`
	DateStart        string `json:"dateStart"`
	DateEnd          string `json:"dateEnd"`
	Status           string `json:"status"`
	Currency         string `json:"currency"`
	LineCount        string `json:"lineCount"`
	Total            string `json:"total"`
}

type SupplierContractPriceScheduleEmptyLabels struct {
	Title             string `json:"title"`
	Message           string `json:"message"`
	ScheduledTitle    string `json:"scheduledTitle"`
	ScheduledMessage  string `json:"scheduledMessage"`
	ActiveTitle       string `json:"activeTitle"`
	ActiveMessage     string `json:"activeMessage"`
	SupersededTitle   string `json:"supersededTitle"`
	SupersededMessage string `json:"supersededMessage"`
	CancelledTitle    string `json:"cancelledTitle"`
	CancelledMessage  string `json:"cancelledMessage"`
}

type SupplierContractPriceScheduleFormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionValidity  string `json:"sectionValidity"`
	SectionScoping   string `json:"sectionScoping"`
	SectionLifecycle string `json:"sectionLifecycle"`
	SectionNotes     string `json:"sectionNotes"`

	// Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	NameInfo               string `json:"nameInfo"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID             string `json:"internalId"`
	InternalIDPlaceholder  string `json:"internalIdPlaceholder"`
	InternalIDInfo         string `json:"internalIdInfo"`

	// Scoping
	SupplierContract       string `json:"supplierContract"`
	SelectSupplierContract string `json:"selectSupplierContract"`
	SupplierContractInfo   string `json:"supplierContractInfo"`

	// Validity
	DateStart          string `json:"dateStart"`
	DateStartInfo      string `json:"dateStartInfo"`
	DateEnd            string `json:"dateEnd"`
	DateEndPlaceholder string `json:"dateEndPlaceholder"`
	DateEndInfo        string `json:"dateEndInfo"`

	// Currency / location
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	CurrencyInfo        string `json:"currencyInfo"`
	Location            string `json:"location"`
	SelectLocation      string `json:"selectLocation"`
	LocationInfo        string `json:"locationInfo"`

	// Lifecycle
	Status                    string `json:"status"`
	SelectStatus              string `json:"selectStatus"`
	StatusInfo                string `json:"statusInfo"`
	SequenceNumber            string `json:"sequenceNumber"`
	SequenceNumberPlaceholder string `json:"sequenceNumberPlaceholder"`
	SequenceNumberInfo        string `json:"sequenceNumberInfo"`

	// Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`
}

type SupplierContractPriceScheduleStatusLabels struct {
	Scheduled  string `json:"scheduled"`
	Active     string `json:"active"`
	Superseded string `json:"superseded"`
	Cancelled  string `json:"cancelled"`
}

type SupplierContractPriceScheduleTabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Activity string `json:"activity"`
}

type SupplierContractPriceScheduleLinesLabels struct {
	Title               string                                      `json:"title"`
	Empty               string                                      `json:"empty"`
	AddLine             string                                      `json:"addLine"`
	ColumnContractLine  string                                      `json:"columnContractLine"`
	ColumnUnitPrice     string                                      `json:"columnUnitPrice"`
	ColumnQuantity      string                                      `json:"columnQuantity"`
	ColumnMinimumAmount string                                      `json:"columnMinimumAmount"`
	ColumnCurrency      string                                      `json:"columnCurrency"`
	ColumnCycleOverride string                                      `json:"columnCycleOverride"`
	LineForm            SupplierContractPriceScheduleLineFormLabels `json:"lineForm"`
}

type SupplierContractPriceScheduleLineFormLabels struct {
	SectionLink                   string `json:"sectionLink"`
	SectionPricing                string `json:"sectionPricing"`
	SectionCycle                  string `json:"sectionCycle"`
	SupplierContractLine          string `json:"supplierContractLine"`
	SelectSupplierContractLine    string `json:"selectSupplierContractLine"`
	SupplierContractLineInfo      string `json:"supplierContractLineInfo"`
	UnitPrice                     string `json:"unitPrice"`
	UnitPricePlaceholder          string `json:"unitPricePlaceholder"`
	UnitPriceInfo                 string `json:"unitPriceInfo"`
	MinimumAmount                 string `json:"minimumAmount"`
	MinimumAmountPlaceholder      string `json:"minimumAmountPlaceholder"`
	MinimumAmountInfo             string `json:"minimumAmountInfo"`
	Quantity                      string `json:"quantity"`
	QuantityPlaceholder           string `json:"quantityPlaceholder"`
	QuantityInfo                  string `json:"quantityInfo"`
	Currency                      string `json:"currency"`
	CurrencyPlaceholder           string `json:"currencyPlaceholder"`
	CycleValueOverride            string `json:"cycleValueOverride"`
	CycleValueOverridePlaceholder string `json:"cycleValueOverridePlaceholder"`
	CycleValueOverrideInfo        string `json:"cycleValueOverrideInfo"`
	CycleUnitOverride             string `json:"cycleUnitOverride"`
	CycleUnitOverridePlaceholder  string `json:"cycleUnitOverridePlaceholder"`
	CycleUnitOverrideInfo         string `json:"cycleUnitOverrideInfo"`
}

type SupplierContractPriceScheduleDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	LinesSection         string `json:"linesSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`
}

type SupplierContractPriceScheduleErrorLabels struct {
	PermissionDenied    string `json:"permissionDenied"`
	InvalidFormData     string `json:"invalidFormData"`
	NotFound            string `json:"notFound"`
	IDRequired          string `json:"idRequired"`
	NoPermission        string `json:"noPermission"`
	CannotDelete        string `json:"cannotDelete"`
	InUse               string `json:"inUse"`
	CreationFailed      string `json:"creation_failed"`
	UpdateFailed        string `json:"update_failed"`
	DeletionFailed      string `json:"deletion_failed"`
	ListFailed          string `json:"list_failed"`
	AuthorizationFailed string `json:"authorization_failed"`
	ActivationFailed    string `json:"activation_failed"`
	SupersedeFailed     string `json:"supersede_failed"`
	OverlapDetected     string `json:"overlap_detected"`
	LoadFailed          string `json:"loadFailed"`
}

// DefaultSupplierContractPriceScheduleLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultSupplierContractPriceScheduleLabels() SupplierContractPriceScheduleLabels {
	return SupplierContractPriceScheduleLabels{
		Labels: SupplierContractPriceScheduleNounLabels{
			Name:       "Price Schedule",
			NamePlural: "Price Schedules",
			Line:       "Schedule Line",
			LinePlural: "Schedule Lines",
		},
		Page: SupplierContractPriceSchedulePageLabels{
			Heading:           "Contract Price Schedules",
			Caption:           "Date-windowed pricing layered on top of a supplier contract for multi-year escalation",
			HeadingScheduled:  "Scheduled Periods",
			HeadingActive:     "Active Periods",
			HeadingSuperseded: "Superseded Periods",
			HeadingCancelled:  "Cancelled Periods",
			TabTitle:          "Price Schedules",
		},
		Buttons: SupplierContractPriceScheduleButtonLabels{
			Add:       "New Schedule",
			AddLine:   "Add Schedule Line",
			Activate:  "Activate",
			Supersede: "Supersede",
			Cancel:    "Cancel Schedule",
		},
		Filters: SupplierContractPriceScheduleFilterLabels{
			All:                 "All",
			Status:              "Status",
			AnyStatus:           "Any Status",
			SupplierContract:    "Contract",
			AnySupplierContract: "Any Contract",
			DateRange:           "Effective Window",
		},
		Columns: SupplierContractPriceScheduleColumnLabels{
			InternalID:       "ID",
			Name:             "Schedule Name",
			SupplierContract: "Contract",
			SequenceNumber:   "Seq.",
			DateStart:        "Start",
			DateEnd:          "End",
			Status:           "Status",
			Currency:         "Currency",
			LineCount:        "Lines",
			Total:            "Total",
		},
		Empty: SupplierContractPriceScheduleEmptyLabels{
			Title:             "No price schedules yet",
			Message:           "Add a schedule period to layer multi-year pricing onto this contract.",
			ScheduledTitle:    "No upcoming schedules",
			ScheduledMessage:  "Future-dated schedules will appear here once added.",
			ActiveTitle:       "No active schedule",
			ActiveMessage:     "The contract is using header pricing — no schedule is in effect right now.",
			SupersededTitle:   "No past schedules",
			SupersededMessage: "Schedules whose window has passed will be archived here.",
			CancelledTitle:    "No cancelled schedules",
			CancelledMessage:  "Schedules cancelled before activation will appear here.",
		},
		Form: SupplierContractPriceScheduleFormLabels{
			SectionIdentity:           "Schedule Identity",
			SectionValidity:           "Validity Window",
			SectionScoping:            "Scoping",
			SectionLifecycle:          "Lifecycle",
			SectionNotes:              "Notes",
			Name:                      "Schedule Name",
			NamePlaceholder:           "e.g. Year 1 (2026)",
			NameInfo:                  "Human-readable label for this pricing window. Often a year or renewal label.",
			Description:               "Description",
			DescriptionPlaceholder:    "Optional details about this pricing window...",
			InternalID:                "Internal ID",
			InternalIDPlaceholder:     "Auto-generated",
			InternalIDInfo:            "Auto-generated unique identifier.",
			SupplierContract:          "Contract",
			SelectSupplierContract:    "Select contract...",
			SupplierContractInfo:      "The supplier contract this pricing window applies to. Schedules are scoped to a single contract.",
			DateStart:                 "Start Date",
			DateStartInfo:             "When this pricing window takes effect. Window is half-open: start is inclusive.",
			DateEnd:                   "End Date",
			DateEndPlaceholder:        "Leave empty for open-ended",
			DateEndInfo:               "When this pricing window ends. Leave blank for the open-ended last bucket. Window is half-open: end is exclusive.",
			Currency:                  "Currency",
			CurrencyPlaceholder:       "PHP",
			CurrencyInfo:              "ISO 4217 currency for prices in this schedule.",
			Location:                  "Location",
			SelectLocation:            "Select location...",
			LocationInfo:              "Optional location override. Defaults to the parent contract's location.",
			Status:                    "Status",
			SelectStatus:              "Select status...",
			StatusInfo:                "Lifecycle state. New schedules default to Scheduled and progress to Active when their window arrives.",
			SequenceNumber:            "Sequence",
			SequenceNumberPlaceholder: "1",
			SequenceNumberInfo:        "Ordering position within the contract (1, 2, 3...).",
			Notes:                     "Notes",
			NotesPlaceholder:          "Internal notes about this pricing window...",
			NotesInfo:                 "Internal remarks only. Not visible to the supplier.",
		},
		Status: SupplierContractPriceScheduleStatusLabels{
			Scheduled:  "Scheduled",
			Active:     "Active",
			Superseded: "Superseded",
			Cancelled:  "Cancelled",
		},
		Tabs: SupplierContractPriceScheduleTabLabels{
			Info:     "Information",
			Lines:    "Schedule Lines",
			Activity: "Activity",
		},
		Lines: SupplierContractPriceScheduleLinesLabels{
			Title:               "Schedule Lines",
			Empty:               "No schedule lines yet. Add lines to override per-line pricing for this window.",
			AddLine:             "Add Line",
			ColumnContractLine:  "Contract Line",
			ColumnUnitPrice:     "Unit Price",
			ColumnQuantity:      "Qty",
			ColumnMinimumAmount: "Minimum",
			ColumnCurrency:      "Currency",
			ColumnCycleOverride: "Cycle Override",
			LineForm: SupplierContractPriceScheduleLineFormLabels{
				SectionLink:                   "Contract Line",
				SectionPricing:                "Pricing",
				SectionCycle:                  "Cycle Override",
				SupplierContractLine:          "Contract Line",
				SelectSupplierContractLine:    "Select contract line...",
				SupplierContractLineInfo:      "The line whose unit price is overridden during this window.",
				UnitPrice:                     "Unit Price",
				UnitPricePlaceholder:          "0.00",
				UnitPriceInfo:                 "Per-unit price during this window, in the schedule's currency.",
				MinimumAmount:                 "Minimum Amount",
				MinimumAmountPlaceholder:      "0.00",
				MinimumAmountInfo:             "For Minimum Commitment lines: the floor charged per cycle.",
				Quantity:                      "Quantity",
				QuantityPlaceholder:           "0",
				QuantityInfo:                  "Optional committed quantity for blanket or minimum-commitment lines.",
				Currency:                      "Currency",
				CurrencyPlaceholder:           "PHP",
				CycleValueOverride:            "Cycle Value Override",
				CycleValueOverridePlaceholder: "e.g. 1",
				CycleValueOverrideInfo:        "Optional cycle-length override. Most lines inherit the contract cycle.",
				CycleUnitOverride:             "Cycle Unit Override",
				CycleUnitOverridePlaceholder:  "month",
				CycleUnitOverrideInfo:         "Optional cycle-unit override (day, week, month, year).",
			},
		},
		Detail: SupplierContractPriceScheduleDetailLabels{
			PageTitle:            "Schedule Details",
			Title:                "Schedule Detail",
			InfoSection:          "Schedule Information",
			LinesSection:         "Per-Line Pricing",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this schedule will appear here.",
			TabAttachments:       "Attachments",
		},
		Errors: SupplierContractPriceScheduleErrorLabels{
			PermissionDenied:    "You do not have permission to perform this action.",
			InvalidFormData:     "Invalid form data. Please check your inputs and try again.",
			NotFound:            "Price schedule not found.",
			IDRequired:          "Schedule ID is required.",
			NoPermission:        "No permission.",
			CannotDelete:        "Cannot delete — this schedule is currently active or has dependent lines.",
			InUse:               "Cannot delete — this schedule is referenced by existing records.",
			CreationFailed:      "Schedule creation failed",
			UpdateFailed:        "Schedule update failed",
			DeletionFailed:      "Schedule deletion failed",
			ListFailed:          "Failed to retrieve price schedules",
			AuthorizationFailed: "Authorization failed for price schedules",
			ActivationFailed:    "Schedule activation failed",
			SupersedeFailed:     "Schedule supersede failed",
			OverlapDetected:     "Schedule windows overlap; adjust dates and retry",
			LoadFailed:          "Failed to load price schedule",
		},
	}
}

// ---------------------------------------------------------------------------
// ExpenseRecognition labels  (SPS P10)
// ---------------------------------------------------------------------------

// ExpenseRecognitionLabels holds all translatable strings for the
// expense_recognition module. Loaded from lyngua key root "expenseRecognition".
type ExpenseRecognitionLabels struct {
	Page    ExpenseRecognitionPageLabels    `json:"page"`
	Buttons ExpenseRecognitionButtonLabels  `json:"buttons"`
	Columns ExpenseRecognitionColumnLabels  `json:"columns"`
	Tabs    ExpenseRecognitionTabLabels     `json:"tabs"`
	Detail  ExpenseRecognitionDetailLabels  `json:"detail"`
	Lines   ExpenseRecognitionLineLabels    `json:"lines"`
	Source  ExpenseRecognitionSourceLabels  `json:"source"`
	Status  ExpenseRecognitionStatusLabels  `json:"status"`
	Actions ExpenseRecognitionActionLabels  `json:"actions"`
	Confirm ExpenseRecognitionConfirmLabels `json:"confirm"`
	Empty   ExpenseRecognitionEmptyLabels   `json:"empty"`
	Errors  ExpenseRecognitionErrorLabels   `json:"errors"`
}

type ExpenseRecognitionPageLabels struct {
	Heading         string `json:"heading"`
	Caption         string `json:"caption"`
	HeadingDraft    string `json:"headingDraft"`
	HeadingPosted   string `json:"headingPosted"`
	HeadingReversed string `json:"headingReversed"`
	Dashboard       string `json:"dashboard"`
}

type ExpenseRecognitionButtonLabels struct {
	Add                      string `json:"add"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	Reverse                  string `json:"reverse"`
}

type ExpenseRecognitionColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	Expenditure      string `json:"expenditure"`
	Currency         string `json:"currency"`
	TotalAmount      string `json:"totalAmount"`
	Status           string `json:"status"`
	Source           string `json:"source"`
	IdempotencyKey   string `json:"idempotencyKey"`
}

type ExpenseRecognitionTabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Source   string `json:"source"`
	Activity string `json:"activity"`
}

type ExpenseRecognitionDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	SourceSection        string `json:"sourceSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`

	// Info-tab + source-tab field labels (4.4)
	Notes           string `json:"notes"`
	SourceContract  string `json:"sourceContract"`
	SourceBill      string `json:"sourceBill"`
	DeferredExpense string `json:"deferredExpense"`
	SourceAccrual   string `json:"sourceAccrual"`
	ReversalOf      string `json:"reversalOf"`
}

type ExpenseRecognitionLineLabels struct {
	Description    string `json:"description"`
	Quantity       string `json:"quantity"`
	UnitAmount     string `json:"unitAmount"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Product        string `json:"product"`
	ExpenseAccount string `json:"expenseAccount"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	AddLine        string `json:"addLine"`

	// Drawer form labels
	FormDescription            string `json:"formDescription"`
	FormDescriptionPlaceholder string `json:"formDescriptionPlaceholder"`
	FormQuantity               string `json:"formQuantity"`
	FormUnitAmount             string `json:"formUnitAmount"`
	FormAmount                 string `json:"formAmount"`
	FormCurrency               string `json:"formCurrency"`
}

type ExpenseRecognitionSourceLabels struct {
	Recurrence  string `json:"recurrence"`
	Expenditure string `json:"expenditure"`
	Manual      string `json:"manual"`
	Reversal    string `json:"reversal"`
}

type ExpenseRecognitionStatusLabels struct {
	Draft    string `json:"draft"`
	Posted   string `json:"posted"`
	Reversed string `json:"reversed"`
}

type ExpenseRecognitionActionLabels struct {
	View                     string `json:"view"`
	Edit                     string `json:"edit"`
	Delete                   string `json:"delete"`
	Reverse                  string `json:"reverse"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	NoPermission             string `json:"noPermission"`
}

type ExpenseRecognitionConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
}

type ExpenseRecognitionEmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	DraftTitle      string `json:"draftTitle"`
	DraftMessage    string `json:"draftMessage"`
	PostedTitle     string `json:"postedTitle"`
	PostedMessage   string `json:"postedMessage"`
	ReversedTitle   string `json:"reversedTitle"`
	ReversedMessage string `json:"reversedMessage"`
}

type ExpenseRecognitionErrorLabels struct {
	PermissionDenied     string `json:"permissionDenied"`
	InvalidFormData      string `json:"invalidFormData"`
	NotFound             string `json:"notFound"`
	IDRequired           string `json:"idRequired"`
	NoPermission         string `json:"noPermission"`
	CreationFailed       string `json:"creation_failed"`
	UpdateFailed         string `json:"update_failed"`
	DeletionFailed       string `json:"deletion_failed"`
	ListFailed           string `json:"list_failed"`
	ReverseFailed        string `json:"reverse_failed"`
	IdempotencyCollision string `json:"idempotency_collision"`
	LoadFailed           string `json:"load_failed"`
}

// DefaultExpenseRecognitionLabels returns English fallback labels.
// Tier overrides belong in lyngua JSON.
func DefaultExpenseRecognitionLabels() ExpenseRecognitionLabels {
	return ExpenseRecognitionLabels{
		Page: ExpenseRecognitionPageLabels{
			Heading:         "Expense Recognition",
			Caption:         "Period in which a supplier cost is recognized",
			HeadingDraft:    "Draft Recognitions",
			HeadingPosted:   "Posted Recognitions",
			HeadingReversed: "Reversed Recognitions",
			Dashboard:       "Expense Recognition Dashboard",
		},
		Buttons: ExpenseRecognitionButtonLabels{
			Add:                      "New Recognition",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			Reverse:                  "Reverse Recognition",
		},
		Columns: ExpenseRecognitionColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			Expenditure:      "Source Bill",
			Currency:         "Currency",
			TotalAmount:      "Amount",
			Status:           "Status",
			Source:           "Source",
			IdempotencyKey:   "Idempotency Key",
		},
		Tabs: ExpenseRecognitionTabLabels{
			Info:     "Information",
			Lines:    "Recognition Lines",
			Source:   "Source",
			Activity: "Activity",
		},
		Detail: ExpenseRecognitionDetailLabels{
			PageTitle:            "Recognition Details",
			Title:                "Recognition Detail",
			InfoSection:          "Recognition Information",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this recognition will appear here.",
			TabAttachments:       "Attachments",
		},
		Lines: ExpenseRecognitionLineLabels{
			Description:                "Description",
			Quantity:                   "Quantity",
			UnitAmount:                 "Unit Amount",
			Amount:                     "Amount",
			Currency:                   "Currency",
			Product:                    "Product",
			ExpenseAccount:             "Expense Account",
			EmptyTitle:                 "No recognition lines",
			EmptyMessage:               "Lines breaking down this recognition will appear here.",
			AddLine:                    "Add Line",
			FormDescription:            "Description",
			FormDescriptionPlaceholder: "e.g. Cloud hosting — May 2026",
			FormQuantity:               "Quantity",
			FormUnitAmount:             "Unit Amount",
			FormAmount:                 "Amount",
			FormCurrency:               "Currency",
		},
		Source: ExpenseRecognitionSourceLabels{
			Recurrence:  "Recurrence Engine",
			Expenditure: "From Bill",
			Manual:      "Manual",
			Reversal:    "Reversal",
		},
		Status: ExpenseRecognitionStatusLabels{
			Draft:    "Draft",
			Posted:   "Posted",
			Reversed: "Reversed",
		},
		Actions: ExpenseRecognitionActionLabels{
			View:                     "View Recognition",
			Edit:                     "Edit Recognition",
			Delete:                   "Delete Recognition",
			Reverse:                  "Reverse",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			NoPermission:             "No permission",
		},
		Confirm: ExpenseRecognitionConfirmLabels{
			Delete:         "Delete Recognition",
			DeleteMessage:  "Are you sure you want to delete this recognition? Only Draft recognitions can be deleted.",
			Reverse:        "Reverse Recognition",
			ReverseMessage: "Reversing creates a counter-entry and marks this recognition as Reversed. Continue?",
		},
		Empty: ExpenseRecognitionEmptyLabels{
			Title:           "No recognitions yet",
			Message:         "Recognize an expense from a posted bill or directly from a contract cycle.",
			DraftTitle:      "No draft recognitions",
			DraftMessage:    "Draft recognitions will appear here.",
			PostedTitle:     "No posted recognitions",
			PostedMessage:   "Posted recognitions will appear here.",
			ReversedTitle:   "No reversed recognitions",
			ReversedMessage: "Reversed recognitions will appear here.",
		},
		Errors: ExpenseRecognitionErrorLabels{
			PermissionDenied:     "You do not have permission to perform this action.",
			InvalidFormData:      "Invalid form data. Please check your inputs and try again.",
			NotFound:             "Recognition not found.",
			IDRequired:           "Recognition ID is required.",
			NoPermission:         "No permission.",
			CreationFailed:       "Recognition creation failed",
			UpdateFailed:         "Recognition update failed",
			DeletionFailed:       "Recognition deletion failed",
			ListFailed:           "Failed to retrieve recognitions",
			ReverseFailed:        "Recognition reversal failed",
			IdempotencyCollision: "A recognition for this source and period already exists",
			LoadFailed:           "Failed to load recognition",
		},
	}
}

// ---------------------------------------------------------------------------
// AccruedExpense labels  (SPS P10)
// ---------------------------------------------------------------------------

// AccruedExpenseLabels holds all translatable strings for the accrued_expense
// + accrued_expense_settlement modules. Loaded from lyngua key root
// "accruedExpense" with settlement subkeys merged in via composition.
type AccruedExpenseLabels struct {
	Page        AccruedExpensePageLabels       `json:"page"`
	Buttons     AccruedExpenseButtonLabels     `json:"buttons"`
	Columns     AccruedExpenseColumnLabels     `json:"columns"`
	Tabs        AccruedExpenseTabLabels        `json:"tabs"`
	Detail      AccruedExpenseDetailLabels     `json:"detail"`
	Settlements AccruedExpenseSettlementLabels `json:"settlements"`
	Form        AccruedExpenseFormLabels       `json:"form"`
	Status      AccruedExpenseStatusLabels     `json:"status"`
	Actions     AccruedExpenseActionLabels     `json:"actions"`
	Confirm     AccruedExpenseConfirmLabels    `json:"confirm"`
	Balances    AccruedExpenseBalanceLabels    `json:"balances"`
	Empty       AccruedExpenseEmptyLabels      `json:"empty"`
	Errors      AccruedExpenseErrorLabels      `json:"errors"`
}

type AccruedExpensePageLabels struct {
	Heading            string `json:"heading"`
	Caption            string `json:"caption"`
	HeadingOutstanding string `json:"headingOutstanding"`
	HeadingPartial     string `json:"headingPartial"`
	HeadingSettled     string `json:"headingSettled"`
	HeadingReversed    string `json:"headingReversed"`
	Dashboard          string `json:"dashboard"`
}

type AccruedExpenseButtonLabels struct {
	Add                string `json:"add"`
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
}

type AccruedExpenseColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Currency         string `json:"currency"`
	AccruedAmount    string `json:"accruedAmount"`
	SettledAmount    string `json:"settledAmount"`
	RemainingAmount  string `json:"remainingAmount"`
	Status           string `json:"status"`
}

type AccruedExpenseTabLabels struct {
	Info        string `json:"info"`
	Settlements string `json:"settlements"`
	Source      string `json:"source"`
	Activity    string `json:"activity"`
}

type AccruedExpenseDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	SettlementsSection   string `json:"settlementsSection"`
	SourceSection        string `json:"sourceSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`

	// Info-tab + source-tab field labels (4.4)
	Notes          string `json:"notes"`
	SourceContract string `json:"sourceContract"`
	Supplier       string `json:"supplier"`
	ExpenseAccount string `json:"expenseAccount"`
	AccrualAccount string `json:"accrualAccount"`
}

type AccruedExpenseSettlementLabels struct {
	Expenditure        string `json:"expenditure"`
	AmountSettled      string `json:"amountSettled"`
	Currency           string `json:"currency"`
	FxRate             string `json:"fxRate"`
	FxAdjustmentAmount string `json:"fxAdjustmentAmount"`
	SettledAt          string `json:"settledAt"`
	Reversal           string `json:"reversal"`
	EmptyTitle         string `json:"emptyTitle"`
	EmptyMessage       string `json:"emptyMessage"`
	AddSettlement      string `json:"addSettlement"`

	// Drawer form labels
	FormExpenditure            string `json:"formExpenditure"`
	FormExpenditurePlaceholder string `json:"formExpenditurePlaceholder"`
	FormAmountSettled          string `json:"formAmountSettled"`
	FormCurrency               string `json:"formCurrency"`
	FormFxRate                 string `json:"formFxRate"`
	FormFxRateInfo             string `json:"formFxRateInfo"`
	FormReversalReason         string `json:"formReversalReason"`
}

type AccruedExpenseFormLabels struct {
	// Section headers
	SectionIdentity   string `json:"sectionIdentity"`
	SectionSource     string `json:"sectionSource"`
	SectionPeriod     string `json:"sectionPeriod"`
	SectionMoney      string `json:"sectionMoney"`
	SectionAccounting string `json:"sectionAccounting"`
	SectionLifecycle  string `json:"sectionLifecycle"`
	SectionNotes      string `json:"sectionNotes"`

	// §1 Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	NameInfo               string `json:"nameInfo"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID             string `json:"internalId"`
	InternalIDPlaceholder  string `json:"internalIdPlaceholder"`

	// §2 Source
	SupplierContract       string `json:"supplierContract"`
	SelectSupplierContract string `json:"selectSupplierContract"`
	SupplierContractInfo   string `json:"supplierContractInfo"`
	Supplier               string `json:"supplier"`
	SelectSupplier         string `json:"selectSupplier"`
	SupplierInfo           string `json:"supplierInfo"`

	// §3 Period
	RecognitionDate      string `json:"recognitionDate"`
	RecognitionDateInfo  string `json:"recognitionDateInfo"`
	PeriodStart          string `json:"periodStart"`
	PeriodStartInfo      string `json:"periodStartInfo"`
	PeriodEnd            string `json:"periodEnd"`
	PeriodEndInfo        string `json:"periodEndInfo"`
	CycleDate            string `json:"cycleDate"`
	CycleDatePlaceholder string `json:"cycleDatePlaceholder"`
	CycleDateInfo        string `json:"cycleDateInfo"`

	// §4 Money
	Currency                 string `json:"currency"`
	CurrencyPlaceholder      string `json:"currencyPlaceholder"`
	CurrencyInfo             string `json:"currencyInfo"`
	AccruedAmount            string `json:"accruedAmount"`
	AccruedAmountPlaceholder string `json:"accruedAmountPlaceholder"`
	AccruedAmountInfo        string `json:"accruedAmountInfo"`
	SettledAmount            string `json:"settledAmount"`
	SettledAmountInfo        string `json:"settledAmountInfo"`
	RemainingAmount          string `json:"remainingAmount"`
	RemainingAmountInfo      string `json:"remainingAmountInfo"`

	// §5 Lifecycle
	Status            string `json:"status"`
	SelectStatus      string `json:"selectStatus"`
	StatusInfo        string `json:"statusInfo"`
	StatusOutstanding string `json:"statusOutstanding"`
	StatusPartial     string `json:"statusPartial"`
	StatusSettled     string `json:"statusSettled"`
	StatusReversed    string `json:"statusReversed"`

	// §6 Accounting
	ExpenseAccount       string `json:"expenseAccount"`
	SelectExpenseAccount string `json:"selectExpenseAccount"`
	ExpenseAccountInfo   string `json:"expenseAccountInfo"`
	AccrualAccount       string `json:"accrualAccount"`
	SelectAccrualAccount string `json:"selectAccrualAccount"`
	AccrualAccountInfo   string `json:"accrualAccountInfo"`

	// §7 Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`

	// Buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Active    string `json:"active"`
}

type AccruedExpenseStatusLabels struct {
	Outstanding string `json:"outstanding"`
	Partial     string `json:"partial"`
	Settled     string `json:"settled"`
	Reversed    string `json:"reversed"`
}

type AccruedExpenseActionLabels struct {
	View               string `json:"view"`
	Edit               string `json:"edit"`
	Delete             string `json:"delete"`
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
	NoPermission       string `json:"noPermission"`
}

type AccruedExpenseConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Settle         string `json:"settle"`
	SettleMessage  string `json:"settleMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
}

type AccruedExpenseBalanceLabels struct {
	Title       string `json:"title"`
	Accrued     string `json:"accrued"`
	Settled     string `json:"settled"`
	Remaining   string `json:"remaining"`
	Utilization string `json:"utilization"`
}

type AccruedExpenseEmptyLabels struct {
	Title              string `json:"title"`
	Message            string `json:"message"`
	OutstandingTitle   string `json:"outstandingTitle"`
	OutstandingMessage string `json:"outstandingMessage"`
	PartialTitle       string `json:"partialTitle"`
	PartialMessage     string `json:"partialMessage"`
	SettledTitle       string `json:"settledTitle"`
	SettledMessage     string `json:"settledMessage"`
	ReversedTitle      string `json:"reversedTitle"`
	ReversedMessage    string `json:"reversedMessage"`
}

type AccruedExpenseErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	CreationFailed   string `json:"creation_failed"`
	UpdateFailed     string `json:"update_failed"`
	DeletionFailed   string `json:"deletion_failed"`
	ListFailed       string `json:"list_failed"`
	SettleFailed     string `json:"settle_failed"`
	ReverseFailed    string `json:"reverse_failed"`
	BalanceDrift     string `json:"balance_drift"`
	LoadFailed       string `json:"load_failed"`
}

// DefaultAccruedExpenseLabels returns English fallback labels.
func DefaultAccruedExpenseLabels() AccruedExpenseLabels {
	return AccruedExpenseLabels{
		Page: AccruedExpensePageLabels{
			Heading:            "Accrued Expenses",
			Caption:            "Recognized supplier obligations awaiting the actual bill",
			HeadingOutstanding: "Outstanding Accruals",
			HeadingPartial:     "Partially Settled",
			HeadingSettled:     "Settled Accruals",
			HeadingReversed:    "Reversed Accruals",
			Dashboard:          "Accrued Expense Dashboard",
		},
		Buttons: AccruedExpenseButtonLabels{
			Add:                "New Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle",
			Reverse:            "Reverse",
			AddSettlement:      "Record Settlement",
		},
		Columns: AccruedExpenseColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Currency:         "Currency",
			AccruedAmount:    "Accrued",
			SettledAmount:    "Settled",
			RemainingAmount:  "Remaining",
			Status:           "Status",
		},
		Tabs: AccruedExpenseTabLabels{
			Info:        "Information",
			Settlements: "Settlements",
			Source:      "Source",
			Activity:    "Activity",
		},
		Detail: AccruedExpenseDetailLabels{
			PageTitle:            "Accrual Details",
			Title:                "Accrual Detail",
			InfoSection:          "Accrual Information",
			SettlementsSection:   "Settlements",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this accrual will appear here.",
			TabAttachments:       "Attachments",
		},
		Settlements: AccruedExpenseSettlementLabels{
			Expenditure:                "Bill",
			AmountSettled:              "Amount Settled",
			Currency:                   "Currency",
			FxRate:                     "FX Rate",
			FxAdjustmentAmount:         "FX Adjustment",
			SettledAt:                  "Settled At",
			Reversal:                   "Reversal",
			EmptyTitle:                 "No settlements yet",
			EmptyMessage:               "Settlements applied against this accrual will appear here.",
			AddSettlement:              "Record Settlement",
			FormExpenditure:            "Bill",
			FormExpenditurePlaceholder: "Select bill...",
			FormAmountSettled:          "Amount Settled",
			FormCurrency:               "Currency",
			FormFxRate:                 "FX Rate",
			FormFxRateInfo:             "Bill currency to accrual currency conversion rate.",
			FormReversalReason:         "Reversal Reason",
		},
		Form: AccruedExpenseFormLabels{
			SectionIdentity:          "Accrual Identity",
			SectionSource:            "Source",
			SectionPeriod:            "Period",
			SectionMoney:             "Money",
			SectionAccounting:        "Accounting",
			SectionLifecycle:         "Lifecycle",
			SectionNotes:             "Notes",
			Name:                     "Accrual Name",
			NamePlaceholder:          "e.g. Utilities — May 2026 (estimate)",
			NameInfo:                 "Descriptive label for the accrued obligation.",
			Description:              "Description",
			DescriptionPlaceholder:   "Optional details about this accrual...",
			InternalID:               "Internal ID",
			InternalIDPlaceholder:    "Auto-generated",
			SupplierContract:         "Source Contract",
			SelectSupplierContract:   "Select contract...",
			SupplierContractInfo:     "The contract whose cycle is being accrued.",
			Supplier:                 "Supplier",
			SelectSupplier:           "Select supplier...",
			SupplierInfo:             "Supplier the obligation is owed to.",
			RecognitionDate:          "Recognition Date",
			RecognitionDateInfo:      "Date the obligation is booked into the period.",
			PeriodStart:              "Period Start",
			PeriodStartInfo:          "Start of the period the accrual covers.",
			PeriodEnd:                "Period End",
			PeriodEndInfo:            "End of the period the accrual covers.",
			CycleDate:                "Cycle Date",
			CycleDatePlaceholder:     "YYYY-MM-DD",
			CycleDateInfo:            "Cycle bucket used for idempotency.",
			Currency:                 "Currency",
			CurrencyPlaceholder:      "PHP",
			CurrencyInfo:             "ISO 4217 currency for the accrual.",
			AccruedAmount:            "Accrued Amount",
			AccruedAmountPlaceholder: "0.00",
			AccruedAmountInfo:        "Estimated obligation for the period.",
			SettledAmount:            "Settled Amount",
			SettledAmountInfo:        "Sum of settlements applied so far. Read-only.",
			RemainingAmount:          "Remaining",
			RemainingAmountInfo:      "Accrued minus settled. Read-only.",
			Status:                   "Status",
			SelectStatus:             "Select status...",
			StatusInfo:               "Lifecycle state.",
			StatusOutstanding:        "Outstanding",
			StatusPartial:            "Partially Settled",
			StatusSettled:            "Settled",
			StatusReversed:           "Reversed",
			ExpenseAccount:           "Expense Account",
			SelectExpenseAccount:     "Select expense account...",
			ExpenseAccountInfo:       "GL account to debit on recognition.",
			AccrualAccount:           "Accrual Account",
			SelectAccrualAccount:     "Select accrual account...",
			AccrualAccountInfo:       "GL account to credit while outstanding.",
			Notes:                    "Notes",
			NotesPlaceholder:         "Internal notes about this accrual...",
			NotesInfo:                "Internal remarks only.",
			Edit:                     "Edit",
			EditTitle:                "Edit Accrual",
			Active:                   "Active",
		},
		Status: AccruedExpenseStatusLabels{
			Outstanding: "Outstanding",
			Partial:     "Partially Settled",
			Settled:     "Settled",
			Reversed:    "Reversed",
		},
		Actions: AccruedExpenseActionLabels{
			View:               "View Accrual",
			Edit:               "Edit Accrual",
			Delete:             "Delete Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle Accrual",
			Reverse:            "Reverse Accrual",
			AddSettlement:      "Record Settlement",
			NoPermission:       "No permission",
		},
		Confirm: AccruedExpenseConfirmLabels{
			Delete:         "Delete Accrual",
			DeleteMessage:  "Are you sure you want to delete this accrual?",
			Settle:         "Settle Accrual",
			SettleMessage:  "Apply a settlement against this accrual?",
			Reverse:        "Reverse Accrual",
			ReverseMessage: "Reversing flips this accrual to Reversed and posts a reversing journal. Continue?",
		},
		Balances: AccruedExpenseBalanceLabels{
			Title:       "Settlement Progress",
			Accrued:     "Accrued",
			Settled:     "Settled",
			Remaining:   "Remaining",
			Utilization: "Settlement Progress",
		},
		Empty: AccruedExpenseEmptyLabels{
			Title:              "No accrued expenses yet",
			Message:            "Accrue from a contract cycle when the period closes before the supplier bill arrives.",
			OutstandingTitle:   "No outstanding accruals",
			OutstandingMessage: "Outstanding accruals waiting for a supplier bill will appear here.",
			PartialTitle:       "No partially settled accruals",
			PartialMessage:     "Accruals with at least one settlement will appear here.",
			SettledTitle:       "No settled accruals",
			SettledMessage:     "Fully reconciled accruals will appear here.",
			ReversedTitle:      "No reversed accruals",
			ReversedMessage:    "Reversed accruals will appear here.",
		},
		Errors: AccruedExpenseErrorLabels{
			PermissionDenied: "You do not have permission to perform this action.",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Accrued expense not found.",
			IDRequired:       "Accrual ID is required.",
			NoPermission:     "No permission.",
			CreationFailed:   "Accrual creation failed",
			UpdateFailed:     "Accrual update failed",
			DeletionFailed:   "Accrual deletion failed",
			ListFailed:       "Failed to retrieve accruals",
			SettleFailed:     "Settlement recording failed",
			ReverseFailed:    "Accrual reversal failed",
			BalanceDrift:     "Accrual balance drift detected; recompute required",
			LoadFailed:       "Failed to load accrual",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierSubscription labels (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// SupplierSubscriptionLabels holds all translatable strings for the supplier_subscription module.
type SupplierSubscriptionLabels struct {
	Page    SupplierSubscriptionPageLabels    `json:"page"`
	Columns SupplierSubscriptionColumnLabels  `json:"columns"`
	Tabs    SupplierSubscriptionTabLabels     `json:"tabs"`
	Detail  SupplierSubscriptionDetailLabels  `json:"detail"`
	Form    SupplierSubscriptionFormLabels    `json:"form"`
	Actions SupplierSubscriptionActionLabels  `json:"actions"`
	Confirm SupplierSubscriptionConfirmLabels `json:"confirm"`
	Buttons SupplierSubscriptionButtonLabels  `json:"buttons"`
	Bulk    SupplierSubscriptionBulkLabels    `json:"bulk"`
	Status  SupplierSubscriptionStatusLabels  `json:"status"`
	Empty   SupplierSubscriptionEmptyLabels   `json:"empty"`
	Errors  SupplierSubscriptionErrorLabels   `json:"errors"`
}

type SupplierSubscriptionPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierSubscriptionColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	CostPlan  string `json:"costPlan"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Active    string `json:"active"`
	AutoRenew string `json:"autoRenew"`
	Code      string `json:"code"`
}

type SupplierSubscriptionTabLabels struct {
	Info                 string `json:"info"`
	CostPlan             string `json:"costPlan"`
	LinkedExpenditures   string `json:"linkedExpenditures"`
	LinkedPurchaseOrders string `json:"linkedPurchaseOrders"`
	LinkedRecognitions   string `json:"linkedRecognitions"`
	Activity             string `json:"activity"`
}

type SupplierSubscriptionDetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	Supplier    string `json:"supplier"`
	CostPlan    string `json:"costPlan"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
	AutoRenew   string `json:"autoRenew"`
	Location    string `json:"location"`
	Notes       string `json:"notes"`

	// Linked-recognitions tab (4.4)
	Recognitions SupplierSubscriptionRecognitionsLabels `json:"recognitions"`
}

// SupplierSubscriptionRecognitionsLabels labels the linked-recognitions tab
// table headers and empty state on the supplier_subscription detail page.
type SupplierSubscriptionRecognitionsLabels struct {
	Name            string `json:"name"`
	Status          string `json:"status"`
	RecognitionDate string `json:"recognitionDate"`
	Amount          string `json:"amount"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

type SupplierSubscriptionFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Code                string `json:"code"`
	CodePlaceholder     string `json:"codePlaceholder"`
	Supplier            string `json:"supplier"`
	SupplierPlaceholder string `json:"supplierPlaceholder"`
	SupplierSearch      string `json:"supplierSearch"`
	SupplierNoResults   string `json:"supplierNoResults"`
	CostPlan            string `json:"costPlan"`
	CostPlanPlaceholder string `json:"costPlanPlaceholder"`
	CostPlanSearch      string `json:"costPlanSearch"`
	CostPlanNoResults   string `json:"costPlanNoResults"`
	AutoRenew           string `json:"autoRenew"`
	Active              string `json:"active"`
	StartDate           string `json:"startDate"`
	StartTime           string `json:"startTime"`
	EndDate             string `json:"endDate"`
	EndTime             string `json:"endTime"`
	TimePlaceholder     string `json:"timePlaceholder"`
	Notes               string `json:"notes"`
	NotesPlaceholder    string `json:"notesPlaceholder"`
	CurrencyError       string `json:"currencyError"`
	EditLockedReason    string `json:"editLockedReason"`
}

type SupplierSubscriptionActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierSubscriptionConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type SupplierSubscriptionButtonLabels struct {
	AddSupplierSubscription string `json:"addSupplierSubscription"`
	RecognizeExpense        string `json:"recognizeExpense"`
}

type SupplierSubscriptionBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierSubscriptionStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierSubscriptionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierSubscriptionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierSubscriptionLabels returns English fallback labels for the supplier_subscription module.
func DefaultSupplierSubscriptionLabels() SupplierSubscriptionLabels {
	return SupplierSubscriptionLabels{
		Page: SupplierSubscriptionPageLabels{
			Heading:         "Supplier Subscriptions",
			HeadingActive:   "Active Supplier Subscriptions",
			HeadingInactive: "Inactive Supplier Subscriptions",
			Caption:         "Recurring supplier commitments",
			CaptionActive:   "Active recurring supplier commitments",
			CaptionInactive: "Inactive recurring supplier commitments",
			PageTitle:       "Supplier Subscription",
		},
		Columns: SupplierSubscriptionColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			CostPlan:  "Cost Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Active:    "Status",
			Code:      "Code",
		},
		Tabs: SupplierSubscriptionTabLabels{
			Info:                 "Info",
			CostPlan:             "Cost Plan",
			LinkedExpenditures:   "Expenditures",
			LinkedPurchaseOrders: "Purchase Orders",
			LinkedRecognitions:   "Recognitions",
			Activity:             "Activity",
		},
		Detail: SupplierSubscriptionDetailLabels{
			InfoSection: "Subscription Details",
			Name:        "Name",
			Supplier:    "Supplier",
			CostPlan:    "Cost Plan",
			Code:        "Code",
			StartDate:   "Start Date",
			EndDate:     "End Date",
			Active:      "Active",
			Inactive:    "Inactive",
			AutoRenew:   "Auto-renew",
			Location:    "Location",
			Notes:       "Notes",
		},
		Form: SupplierSubscriptionFormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. Cloud Hosting — AWS",
			Code:                  "Code",
			CodePlaceholder:       "e.g. SUB-2026-001",
			Supplier:              "Supplier",
			SupplierPlaceholder:   "Search supplier…",
			SupplierSearch:        "Search suppliers",
			SupplierNoResults:     "No suppliers found",
			CostPlan:              "Cost Plan",
			CostPlanPlaceholder:   "Search cost plan…",
			CostPlanSearch:        "Search cost plans",
			CostPlanNoResults:     "No cost plans found",
			AutoRenew:             "Auto-renew",
			Active:                "Active",
			StartDate:             "Start Date",
			StartTime:             "Start Time",
			EndDate:               "End Date",
			EndTime:               "End Time",
			TimePlaceholder:       "HH:MM",
			Notes:                 "Notes",
			NotesPlaceholder:      "Internal notes about this subscription",
			CurrencyError:         "The selected cost plan's billing currency does not match the workspace functional currency.",
			EditLockedReason:      "This subscription has linked expenditures and cannot be fully edited.",
		},
		Actions: SupplierSubscriptionActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierSubscriptionConfirmLabels{
			Delete:                "Delete Supplier Subscription",
			DeleteMessage:         "Are you sure you want to delete this supplier subscription?",
			Activate:              "Activate Supplier Subscription",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Subscription",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Subscriptions",
			BulkDeleteMessage:     "Delete selected supplier subscriptions?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier subscriptions?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier subscriptions?",
		},
		Buttons: SupplierSubscriptionButtonLabels{
			AddSupplierSubscription: "Add Supplier Subscription",
			RecognizeExpense:        "Recognize Expense",
		},
		Bulk: SupplierSubscriptionBulkLabels{
			Delete: "Delete",
		},
		Status: SupplierSubscriptionStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierSubscriptionEmptyLabels{
			Title:   "No supplier subscriptions yet",
			Message: "Add a supplier subscription to start tracking recurring vendor commitments.",
		},
		Errors: SupplierSubscriptionErrorLabels{
			PermissionDenied: "You do not have permission to perform this action.",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Supplier subscription not found.",
			IDRequired:       "Supplier subscription ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This subscription is in use and cannot be deleted.",
			LoadFailed:       "Failed to load supplier subscription.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — CostSchedule labels
// ---------------------------------------------------------------------------

// CostScheduleLabels holds all translatable strings for the cost_schedule module.
type CostScheduleLabels struct {
	Page    CostSchedulePageLabels    `json:"page"`
	Columns CostScheduleColumnLabels  `json:"columns"`
	Tabs    CostScheduleTabLabels     `json:"tabs"`
	Detail  CostScheduleDetailLabels  `json:"detail"`
	Form    CostScheduleFormLabels    `json:"form"`
	Actions CostScheduleActionLabels  `json:"actions"`
	Confirm CostScheduleConfirmLabels `json:"confirm"`
	Buttons CostScheduleButtonLabels  `json:"buttons"`
	Bulk    CostScheduleBulkLabels    `json:"bulk"`
	Status  CostScheduleStatusLabels  `json:"status"`
	Empty   CostScheduleEmptyLabels   `json:"empty"`
	Errors  CostScheduleErrorLabels   `json:"errors"`
}

type CostSchedulePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type CostScheduleColumnLabels struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Location  string `json:"location"`
	Active    string `json:"active"`
}

type CostScheduleTabLabels struct {
	Info      string `json:"info"`
	CostPlans string `json:"costPlans"`
	Activity  string `json:"activity"`
}

type CostScheduleDetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type CostScheduleFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descPlaceholder"`
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	Active              string `json:"active"`
}

type CostScheduleActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type CostScheduleConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type CostScheduleButtonLabels struct {
	AddCostSchedule string `json:"addCostSchedule"`
}

type CostScheduleBulkLabels struct {
	Delete string `json:"delete"`
}

type CostScheduleStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type CostScheduleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CostScheduleErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultCostScheduleLabels returns English fallback labels.
func DefaultCostScheduleLabels() CostScheduleLabels {
	return CostScheduleLabels{
		Page: CostSchedulePageLabels{
			Heading:         "Cost Schedules",
			HeadingActive:   "Active Cost Schedules",
			HeadingInactive: "Inactive Cost Schedules",
			Caption:         "Date-bounded supplier pricing windows",
			CaptionActive:   "Active pricing windows",
			CaptionInactive: "Inactive pricing windows",
			PageTitle:       "Cost Schedule",
		},
		Columns: CostScheduleColumnLabels{
			Name:      "Name",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Location:  "Location",
			Active:    "Status",
		},
		Tabs: CostScheduleTabLabels{
			Info:      "Info",
			CostPlans: "Cost Plans",
			Activity:  "Activity",
		},
		Detail: CostScheduleDetailLabels{
			InfoSection: "Schedule Details",
			Name:        "Name",
			StartDate:   "Start Date",
			EndDate:     "End Date",
			Location:    "Location",
			Description: "Description",
			Active:      "Active",
			Inactive:    "Inactive",
		},
		Form: CostScheduleFormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. Q1 2026 Supplier Rates",
			Description:           "Description",
			DescPlaceholder:       "Internal notes about this cost schedule",
			StartDate:             "Start Date",
			EndDate:               "End Date",
			Location:              "Location",
			LocationPlaceholder:   "Select location",
			Active:                "Active",
		},
		Actions: CostScheduleActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: CostScheduleConfirmLabels{
			Delete:                "Delete Cost Schedule",
			DeleteMessage:         "Are you sure you want to delete this cost schedule?",
			Activate:              "Activate Cost Schedule",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Cost Schedule",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Cost Schedules",
			BulkDeleteMessage:     "Delete selected cost schedules?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected cost schedules?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected cost schedules?",
		},
		Buttons: CostScheduleButtonLabels{
			AddCostSchedule: "Add Cost Schedule",
		},
		Bulk: CostScheduleBulkLabels{Delete: "Delete"},
		Status: CostScheduleStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: CostScheduleEmptyLabels{
			Title:   "No cost schedules yet",
			Message: "Add a cost schedule to group supplier cost plans by date range.",
		},
		Errors: CostScheduleErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost schedule not found.",
			IDRequired:       "Cost schedule ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This cost schedule has linked cost plans and cannot be deleted.",
			LoadFailed:       "Failed to load cost schedule.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierPlan labels
// ---------------------------------------------------------------------------

// SupplierPlanLabels holds all translatable strings for the supplier_plan module.
type SupplierPlanLabels struct {
	Page    SupplierPlanPageLabels    `json:"page"`
	Columns SupplierPlanColumnLabels  `json:"columns"`
	Tabs    SupplierPlanTabLabels     `json:"tabs"`
	Detail  SupplierPlanDetailLabels  `json:"detail"`
	Form    SupplierPlanFormLabels    `json:"form"`
	Actions SupplierPlanActionLabels  `json:"actions"`
	Confirm SupplierPlanConfirmLabels `json:"confirm"`
	Buttons SupplierPlanButtonLabels  `json:"buttons"`
	Bulk    SupplierPlanBulkLabels    `json:"bulk"`
	Status  SupplierPlanStatusLabels  `json:"status"`
	Empty   SupplierPlanEmptyLabels   `json:"empty"`
	Errors  SupplierPlanErrorLabels   `json:"errors"`
}

type SupplierPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierPlanColumnLabels struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Supplier string `json:"supplier"`
	Active   string `json:"active"`
}

type SupplierPlanTabLabels struct {
	Info         string `json:"info"`
	CostPlans    string `json:"costPlans"`
	ProductPlans string `json:"productPlans"`
	Activity     string `json:"activity"`
}

type SupplierPlanDetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Supplier    string `json:"supplier"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type SupplierPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Code                string `json:"code"`
	CodePlaceholder     string `json:"codePlaceholder"`
	Supplier            string `json:"supplier"`
	SupplierPlaceholder string `json:"supplierPlaceholder"`
	Active              string `json:"active"`
}

type SupplierPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierPlanConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type SupplierPlanButtonLabels struct {
	AddSupplierPlan string `json:"addSupplierPlan"`
}

type SupplierPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierPlanLabels returns English fallback labels.
func DefaultSupplierPlanLabels() SupplierPlanLabels {
	return SupplierPlanLabels{
		Page: SupplierPlanPageLabels{
			Heading:         "Supplier Plans",
			HeadingActive:   "Active Supplier Plans",
			HeadingInactive: "Inactive Supplier Plans",
			Caption:         "Supplier product and pricing plans",
			CaptionActive:   "Active supplier plans",
			CaptionInactive: "Inactive supplier plans",
			PageTitle:       "Supplier Plan",
		},
		Columns: SupplierPlanColumnLabels{
			Name:     "Name",
			Code:     "Code",
			Supplier: "Supplier",
			Active:   "Status",
		},
		Tabs: SupplierPlanTabLabels{
			Info:         "Info",
			CostPlans:    "Cost Plans",
			ProductPlans: "Product Plans",
			Activity:     "Activity",
		},
		Detail: SupplierPlanDetailLabels{
			InfoSection: "Plan Details",
			Name:        "Name",
			Code:        "Code",
			Supplier:    "Supplier",
			Active:      "Active",
			Inactive:    "Inactive",
		},
		Form: SupplierPlanFormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. AWS Standard Plan",
			Code:                  "Code",
			CodePlaceholder:       "e.g. PLAN-AWS-001",
			Supplier:              "Supplier",
			SupplierPlaceholder:   "Select supplier",
			Active:                "Active",
		},
		Actions: SupplierPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierPlanConfirmLabels{
			Delete:                "Delete Supplier Plan",
			DeleteMessage:         "Are you sure you want to delete this supplier plan?",
			Activate:              "Activate Supplier Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Plans",
			BulkDeleteMessage:     "Delete selected supplier plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier plans?",
		},
		Buttons: SupplierPlanButtonLabels{AddSupplierPlan: "Add Supplier Plan"},
		Bulk:    SupplierPlanBulkLabels{Delete: "Delete"},
		Status: SupplierPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierPlanEmptyLabels{
			Title:   "No supplier plans yet",
			Message: "Add a supplier plan to group cost plans and product plans for a vendor.",
		},
		Errors: SupplierPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Supplier plan not found.",
			IDRequired:       "Supplier plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This supplier plan has linked cost plans or product plans and cannot be deleted.",
			LoadFailed:       "Failed to load supplier plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — CostPlan labels
// ---------------------------------------------------------------------------

// CostPlanLabels holds all translatable strings for the cost_plan module.
type CostPlanLabels struct {
	Page    CostPlanPageLabels    `json:"page"`
	Columns CostPlanColumnLabels  `json:"columns"`
	Tabs    CostPlanTabLabels     `json:"tabs"`
	Detail  CostPlanDetailLabels  `json:"detail"`
	Form    CostPlanFormLabels    `json:"form"`
	Actions CostPlanActionLabels  `json:"actions"`
	Confirm CostPlanConfirmLabels `json:"confirm"`
	Buttons CostPlanButtonLabels  `json:"buttons"`
	Bulk    CostPlanBulkLabels    `json:"bulk"`
	Status  CostPlanStatusLabels  `json:"status"`
	Empty   CostPlanEmptyLabels   `json:"empty"`
	Errors  CostPlanErrorLabels   `json:"errors"`
}

type CostPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type CostPlanColumnLabels struct {
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
}

type CostPlanTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedSubscriptions string `json:"linkedSubscriptions"`
	Activity            string `json:"activity"`
}

type CostPlanDetailLabels struct {
	InfoSection  string `json:"infoSection"`
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	AmountBasis  string `json:"amountBasis"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	BillingCycle string `json:"billingCycle"`
	DefaultTerm  string `json:"defaultTerm"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
	Inactive     string `json:"inactive"`
}

type CostPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"descPlaceholder"`
	SupplierPlan            string `json:"supplierPlan"`
	SupplierPlanPlaceholder string `json:"supplierPlanPlaceholder"`
	CostSchedule            string `json:"costSchedule"`
	CostSchedulePlaceholder string `json:"costSchedulePlaceholder"`
	BillingKind             string `json:"billingKind"`
	AmountBasis             string `json:"amountBasis"`
	Amount                  string `json:"amount"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	Currency                string `json:"currency"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	BillingCycle            string `json:"billingCycle"`
	BillingCyclePlaceholder string `json:"billingCyclePlaceholder"`
	DefaultTerm             string `json:"defaultTerm"`
	DefaultTermPlaceholder  string `json:"defaultTermPlaceholder"`
	Active                  string `json:"active"`

	// BillingKind option labels
	BillingKindOneTime    string `json:"billingKindOneTime"`
	BillingKindRecurring  string `json:"billingKindRecurring"`
	BillingKindContract   string `json:"billingKindContract"`
	BillingKindUsageBased string `json:"billingKindUsageBased"`
	BillingKindAdHoc      string `json:"billingKindAdHoc"`

	// AmountBasis option labels
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`

	// Duration unit option labels (shared by billing_cycle_unit and default_term_unit)
	DurationUnitDay   string `json:"durationUnitDay"`
	DurationUnitWeek  string `json:"durationUnitWeek"`
	DurationUnitMonth string `json:"durationUnitMonth"`
	DurationUnitYear  string `json:"durationUnitYear"`
}

type CostPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type CostPlanConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type CostPlanButtonLabels struct {
	AddCostPlan string `json:"addCostPlan"`
}

type CostPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type CostPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type CostPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CostPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultCostPlanLabels returns English fallback labels.
func DefaultCostPlanLabels() CostPlanLabels {
	return CostPlanLabels{
		Page: CostPlanPageLabels{
			Heading:         "Cost Plans",
			HeadingActive:   "Active Cost Plans",
			HeadingInactive: "Inactive Cost Plans",
			Caption:         "Supplier pricing plans and billing schedules",
			CaptionActive:   "Active cost plans",
			CaptionInactive: "Inactive cost plans",
			PageTitle:       "Cost Plan",
		},
		Columns: CostPlanColumnLabels{
			Name:         "Name",
			BillingKind:  "Billing Kind",
			Amount:       "Amount",
			Currency:     "Currency",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Status",
		},
		Tabs: CostPlanTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedSubscriptions: "Subscriptions",
			Activity:            "Activity",
		},
		Detail: CostPlanDetailLabels{
			InfoSection:  "Cost Plan Details",
			Name:         "Name",
			BillingKind:  "Billing Kind",
			AmountBasis:  "Amount Basis",
			Amount:       "Amount",
			Currency:     "Currency",
			BillingCycle: "Billing Cycle",
			DefaultTerm:  "Default Term",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Active",
			Inactive:     "Inactive",
		},
		Form: CostPlanFormLabels{
			SectionIdentification:       "Identification",
			SectionRelationships:        "Relationships",
			SectionConfiguration:        "Configuration",
			SectionSchedule:             "Schedule",
			SectionNotes:                "Notes",
			Name:                        "Name",
			NamePlaceholder:             "e.g. AWS EC2 Monthly",
			Description:                 "Description",
			DescPlaceholder:             "Internal notes about this cost plan",
			SupplierPlan:                "Supplier Plan",
			SupplierPlanPlaceholder:     "Select supplier plan",
			CostSchedule:                "Cost Schedule",
			CostSchedulePlaceholder:     "Select cost schedule",
			BillingKind:                 "Billing Kind",
			AmountBasis:                 "Amount Basis",
			Amount:                      "Amount",
			AmountPlaceholder:           "0.00",
			Currency:                    "Currency",
			CurrencyPlaceholder:         "e.g. PHP",
			BillingCycle:                "Billing Cycle",
			BillingCyclePlaceholder:     "e.g. 1",
			DefaultTerm:                 "Default Term",
			DefaultTermPlaceholder:      "e.g. 12",
			Active:                      "Active",
			BillingKindOneTime:          "One Time",
			BillingKindRecurring:        "Recurring",
			BillingKindContract:         "Contract",
			BillingKindUsageBased:       "Usage Based",
			BillingKindAdHoc:            "Ad Hoc",
			AmountBasisPerCycle:         "Per Cycle",
			AmountBasisTotalPackage:     "Total Package",
			AmountBasisDerivedFromLines: "Derived From Lines",
			AmountBasisPerOccurrence:    "Per Occurrence",
			DurationUnitDay:             "Day",
			DurationUnitWeek:            "Week",
			DurationUnitMonth:           "Month",
			DurationUnitYear:            "Year",
		},
		Actions: CostPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: CostPlanConfirmLabels{
			Delete:                "Delete Cost Plan",
			DeleteMessage:         "Are you sure you want to delete this cost plan?",
			Activate:              "Activate Cost Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Cost Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Cost Plans",
			BulkDeleteMessage:     "Delete selected cost plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected cost plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected cost plans?",
		},
		Buttons: CostPlanButtonLabels{AddCostPlan: "Add Cost Plan"},
		Bulk:    CostPlanBulkLabels{Delete: "Delete"},
		Status: CostPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: CostPlanEmptyLabels{
			Title:   "No cost plans yet",
			Message: "Add a cost plan to define billing terms for a supplier engagement.",
		},
		Errors: CostPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan not found.",
			IDRequired:       "Cost plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This cost plan has linked subscriptions and cannot be deleted.",
			LoadFailed:       "Failed to load cost plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan labels
// ---------------------------------------------------------------------------

// SupplierProductPlanLabels holds all translatable strings for the supplier_product_plan module.
type SupplierProductPlanLabels struct {
	Page    SupplierProductPlanPageLabels    `json:"page"`
	Columns SupplierProductPlanColumnLabels  `json:"columns"`
	Tabs    SupplierProductPlanTabLabels     `json:"tabs"`
	Detail  SupplierProductPlanDetailLabels  `json:"detail"`
	Form    SupplierProductPlanFormLabels    `json:"form"`
	Actions SupplierProductPlanActionLabels  `json:"actions"`
	Confirm SupplierProductPlanConfirmLabels `json:"confirm"`
	Buttons SupplierProductPlanButtonLabels  `json:"buttons"`
	Bulk    SupplierProductPlanBulkLabels    `json:"bulk"`
	Status  SupplierProductPlanStatusLabels  `json:"status"`
	Empty   SupplierProductPlanEmptyLabels   `json:"empty"`
	Errors  SupplierProductPlanErrorLabels   `json:"errors"`
}

type SupplierProductPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierProductPlanColumnLabels struct {
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
}

type SupplierProductPlanTabLabels struct {
	Info          string `json:"info"`
	CostPlanLines string `json:"costPlanLines"`
	Activity      string `json:"activity"`
}

type SupplierProductPlanDetailLabels struct {
	InfoSection    string `json:"infoSection"`
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
}

type SupplierProductPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	SupplierPlan              string `json:"supplierPlan"`
	SupplierPlanPlaceholder   string `json:"supplierPlanPlaceholder"`
	Product                   string `json:"product"`
	ProductPlaceholder        string `json:"productPlaceholder"`
	ProductVariant            string `json:"productVariant"`
	ProductVariantPlaceholder string `json:"productVariantPlaceholder"`
	SupplierSKU               string `json:"supplierSku"`
	SupplierSKUPlaceholder    string `json:"supplierSkuPlaceholder"`
	SupplierUnit              string `json:"supplierUnit"`
	SupplierUnitPlaceholder   string `json:"supplierUnitPlaceholder"`
	Active                    string `json:"active"`
}

type SupplierProductPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierProductPlanConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type SupplierProductPlanButtonLabels struct {
	AddSupplierProductPlan string `json:"addSupplierProductPlan"`
}

type SupplierProductPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierProductPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierProductPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierProductPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierProductPlanLabels returns English fallback labels.
func DefaultSupplierProductPlanLabels() SupplierProductPlanLabels {
	return SupplierProductPlanLabels{
		Page: SupplierProductPlanPageLabels{
			Heading:         "Supplier Product Plans",
			HeadingActive:   "Active Supplier Product Plans",
			HeadingInactive: "Inactive Supplier Product Plans",
			Caption:         "Supplier product catalogue line items",
			CaptionActive:   "Active supplier product plans",
			CaptionInactive: "Inactive supplier product plans",
			PageTitle:       "Supplier Product Plan",
		},
		Columns: SupplierProductPlanColumnLabels{
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Status",
		},
		Tabs: SupplierProductPlanTabLabels{
			Info:          "Info",
			CostPlanLines: "Cost Plan Lines",
			Activity:      "Activity",
		},
		Detail: SupplierProductPlanDetailLabels{
			InfoSection:    "Product Plan Details",
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Active",
			Inactive:       "Inactive",
		},
		Form: SupplierProductPlanFormLabels{
			SectionIdentification:     "Identification",
			SectionRelationships:      "Relationships",
			SectionConfiguration:      "Configuration",
			SectionSchedule:           "Schedule",
			SectionNotes:              "Notes",
			SupplierPlan:              "Supplier Plan",
			SupplierPlanPlaceholder:   "Select supplier plan",
			Product:                   "Product",
			ProductPlaceholder:        "Select product",
			ProductVariant:            "Variant (optional)",
			ProductVariantPlaceholder: "Select variant",
			SupplierSKU:               "Supplier SKU",
			SupplierSKUPlaceholder:    "Supplier's internal SKU code",
			SupplierUnit:              "Supplier Unit",
			SupplierUnitPlaceholder:   "e.g. vCPU·hour",
			Active:                    "Active",
		},
		Actions: SupplierProductPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierProductPlanConfirmLabels{
			Delete:                "Delete Supplier Product Plan",
			DeleteMessage:         "Are you sure you want to delete this supplier product plan?",
			Activate:              "Activate Supplier Product Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Product Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Product Plans",
			BulkDeleteMessage:     "Delete selected supplier product plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier product plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier product plans?",
		},
		Buttons: SupplierProductPlanButtonLabels{AddSupplierProductPlan: "Add Supplier Product Plan"},
		Bulk:    SupplierProductPlanBulkLabels{Delete: "Delete"},
		Status: SupplierProductPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierProductPlanEmptyLabels{
			Title:   "No supplier product plans yet",
			Message: "Add a supplier product plan to map vendor catalogue items to your internal products.",
		},
		Errors: SupplierProductPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Supplier product plan not found.",
			IDRequired:       "Supplier product plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This supplier product plan has linked cost plan lines and cannot be deleted.",
			LoadFailed:       "Failed to load supplier product plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierProductCostPlan labels (inline editor, no full module)
// ---------------------------------------------------------------------------

// SupplierProductCostPlanLabels holds translatable strings for the inline cost plan line editor.
type SupplierProductCostPlanLabels struct {
	Form    SupplierProductCostPlanFormLabels   `json:"form"`
	Columns SupplierProductCostPlanColumnLabels `json:"columns"`
	Empty   SupplierProductCostPlanEmptyLabels  `json:"empty"`
	Actions SupplierProductCostPlanActionLabels `json:"actions"`
	Errors  SupplierProductCostPlanErrorLabels  `json:"errors"`
}

type SupplierProductCostPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	SupplierProductPlan            string `json:"supplierProductPlan"`
	SupplierProductPlanPlaceholder string `json:"supplierProductPlanPlaceholder"`
	BillingTreatment               string `json:"billingTreatment"`
	Amount                         string `json:"amount"`
	AmountPlaceholder              string `json:"amountPlaceholder"`
	MinimumCommitment              string `json:"minimumCommitment"`
	MinimumCommitmentPlaceholder   string `json:"minimumCommitmentPlaceholder"`
	Active                         string `json:"active"`

	// BillingTreatment option labels
	TreatmentRecurring         string `json:"treatmentRecurring"`
	TreatmentOneTimeInitial    string `json:"treatmentOneTimeInitial"`
	TreatmentUsageBased        string `json:"treatmentUsageBased"`
	TreatmentMinimumCommitment string `json:"treatmentMinimumCommitment"`
}

type SupplierProductCostPlanColumnLabels struct {
	SupplierProductPlan string `json:"supplierProductPlan"`
	BillingTreatment    string `json:"billingTreatment"`
	Amount              string `json:"amount"`
	Active              string `json:"active"`
}

type SupplierProductCostPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	AddLine string `json:"addLine"`
}

type SupplierProductCostPlanActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Add          string `json:"add"`
	NoPermission string `json:"noPermission"`
}

type SupplierProductCostPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
}

// DefaultSupplierProductCostPlanLabels returns English fallback labels.
func DefaultSupplierProductCostPlanLabels() SupplierProductCostPlanLabels {
	return SupplierProductCostPlanLabels{
		Form: SupplierProductCostPlanFormLabels{
			SectionIdentification:          "Identification",
			SectionRelationships:           "Relationships",
			SectionConfiguration:           "Configuration",
			SectionSchedule:                "Schedule",
			SectionNotes:                   "Notes",
			SupplierProductPlan:            "Supplier Product Plan",
			SupplierProductPlanPlaceholder: "Select product plan",
			BillingTreatment:               "Billing Treatment",
			Amount:                         "Amount",
			AmountPlaceholder:              "0.00",
			MinimumCommitment:              "Minimum Commitment",
			MinimumCommitmentPlaceholder:   "0.00",
			Active:                         "Active",
			TreatmentRecurring:             "Recurring",
			TreatmentOneTimeInitial:        "One-Time Initial",
			TreatmentUsageBased:            "Usage Based",
			TreatmentMinimumCommitment:     "Minimum Commitment",
		},
		Columns: SupplierProductCostPlanColumnLabels{
			SupplierProductPlan: "Product Plan",
			BillingTreatment:    "Treatment",
			Amount:              "Amount",
			Active:              "Status",
		},
		Empty: SupplierProductCostPlanEmptyLabels{
			Title:   "No cost plan lines yet",
			Message: "Add product-level cost lines to this cost plan.",
			AddLine: "Add Line",
		},
		Actions: SupplierProductCostPlanActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Add:          "Add Line",
			NoPermission: "No permission",
		},
		Errors: SupplierProductCostPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan line not found.",
			IDRequired:       "Cost plan line ID is required.",
		},
	}
}
