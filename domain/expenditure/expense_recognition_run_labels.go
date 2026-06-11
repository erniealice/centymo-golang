package expenditure

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
