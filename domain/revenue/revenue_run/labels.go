package revenuerun

// ---------------------------------------------------------------------------
// Revenue Run labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Revenue Run
// (invoice-run) module. Lyngua root key: "revenueRun".
// D13: naming is revenueRun / revenue_run / RevenueRun / revenue-run everywhere
// except the user-visible VALUE "Invoice Run" (supplied by lyngua).
type Labels struct {
	AppLabel       string            `json:"appLabel"`
	Queue          QueueLabels       `json:"queue"`
	List           ListLabels        `json:"list"`
	Detail         DetailLabels      `json:"detail"`
	StatusBadges   StatusBadgeLabels `json:"statusBadges"`
	Actions        ActionLabels      `json:"actions"`
	ScopeKind      ScopeKindLabels   `json:"scopeKind"`
	AttemptOutcome OutcomeLabels     `json:"attemptOutcome"`
	Errors         ErrorLabels       `json:"errors"`
	// ToastBatchSuccess is the message shown after a Surface B batch-run
	// submission. Supports the standard {{.Created}}/{{.Skipped}}/{{.Errored}}
	// placeholders, substituted Go-side before the toast is dispatched.
	ToastBatchSuccess string `json:"toastBatchSuccess"`
	// ViewRunLink is the link label used on toasts whose batch produced
	// exactly one run. Multi-run batches omit the link.
	ViewRunLink string `json:"viewRunLink"`
}

// QueueLabels holds copy for the workspace-queue page (Surface B).
type QueueLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// AsOfDateLabel is the label for the AsOfDate date picker above the table.
	AsOfDateLabel string            `json:"asOfDateLabel"`
	Columns       QueueColumnLabels `json:"columns"`
	Empty         QueueEmptyLabels  `json:"empty"`
	Bulk          QueueBulkLabels   `json:"bulk"`
}

type QueueColumnLabels struct {
	Client         string `json:"client"`
	Subscriptions  string `json:"subscriptions"`
	PendingPeriods string `json:"pendingPeriods"`
	Total          string `json:"total"`
	Currency       string `json:"currency"`
	Actions        string `json:"actions"`
	Run            string `json:"run"`
}

type QueueEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type QueueBulkLabels struct {
	RunSelected        string `json:"runSelected"`
	RunAllMatching     string `json:"runAllMatching"`
	CapExceededMessage string `json:"capExceededMessage"`
}

// ListLabels holds copy for the run history list page (Surface D).
type ListLabels struct {
	Title    string           `json:"title"`
	Subtitle string           `json:"subtitle"`
	Columns  ListColumnLabels `json:"columns"`
	Empty    ListEmptyLabels  `json:"empty"`
	Filters  ListFilterLabels `json:"filterLabels"`
}

type ListColumnLabels struct {
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
	Title      string              `json:"title"`
	Tabs       DetailTabLabels     `json:"tabs"`
	Summary    SummaryLabels       `json:"summary"`
	Selections SelectionsTabLabels `json:"selections"`
	Results    ResultsTabLabels    `json:"results"`
	Invoices   InvoicesTabLabels   `json:"invoices"`
}

// SelectionsTabLabels holds column headers and empty-state copy for
// the Selections tab on the run detail page.
type SelectionsTabLabels struct {
	ColSubscription string `json:"colSubscription"`
	ColPeriodStart  string `json:"colPeriodStart"`
	ColPeriodEnd    string `json:"colPeriodEnd"`
	ColPeriodMarker string `json:"colPeriodMarker"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

// ResultsTabLabels holds column headers and empty-state copy for
// the Results tab on the run detail page.
type ResultsTabLabels struct {
	ColSubscription string `json:"colSubscription"`
	ColPeriodStart  string `json:"colPeriodStart"`
	ColPeriodEnd    string `json:"colPeriodEnd"`
	ColOutcome      string `json:"colOutcome"`
	ColErrorCode    string `json:"colErrorCode"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

// InvoicesTabLabels holds column headers and empty-state copy for
// the Invoices tab on the run detail page. Also holds the coming-soon label
// for the Audit History tab.
type InvoicesTabLabels struct {
	ColReference string `json:"colReference"`
	ColDate      string `json:"colDate"`
	ColAmount    string `json:"colAmount"`
	ColStatus    string `json:"colStatus"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	// AuditHistoryComingSoon is the coming-soon message for the Audit History tab.
	AuditHistoryComingSoon string `json:"auditHistoryComingSoon"`
}

type DetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Invoices     string `json:"invoices"`
	AuditHistory string `json:"auditHistory"`
	Attachments  string `json:"attachments"`
}

type SummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"asOfDate"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiatedAt"`
	CompletedAt             string `json:"completedAt"`
	Status                  string `json:"status"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possiblyInterruptedNote"`
}

// StatusBadgeLabels holds display labels for each run status value.
type StatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possiblyInterrupted"`
}

// ActionLabels holds labels for interactive actions on run rows/pages.
type ActionLabels struct {
	Run                   string `json:"run"`
	ReRunFailed           string `json:"reRunFailed"`
	ReRunFailedComingSoon string `json:"reRunFailedComingSoon"`
	ViewRun               string `json:"viewRun"`
	ViewClient            string `json:"viewClient"`
	ViewSubscription      string `json:"viewSubscription"`
}

// ScopeKindLabels holds display labels for each scope kind value.
type ScopeKindLabels struct {
	Subscription string `json:"subscription"`
	Client       string `json:"client"`
	Workspace    string `json:"workspace"`
}

// OutcomeLabels holds display labels for per-attempt outcome values.
type OutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// ErrorLabels holds error message strings for the revenue-run module.
type ErrorLabels struct {
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

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		AppLabel: "Invoice Run",
		Queue: QueueLabels{
			Title:         "Invoice Run Queue",
			Subtitle:      "Clients with pending billing periods ready to invoice",
			AsOfDateLabel: "As of date",
			Columns: QueueColumnLabels{
				Client:         "Client",
				Subscriptions:  "Subscriptions",
				PendingPeriods: "Pending periods",
				Total:          "Total",
				Currency:       "Currency",
				Actions:        "Actions",
				Run:            "Run",
			},
			Empty: QueueEmptyLabels{
				Title:   "Queue is empty",
				Message: "No clients have pending billing periods at this time.",
			},
			Bulk: QueueBulkLabels{
				RunSelected:        "Run for selected",
				RunAllMatching:     "Run for all matching",
				CapExceededMessage: "Capped at 50 clients per batch. Narrow the filter to run the rest.",
			},
		},
		List: ListLabels{
			Title:    "Invoice Runs",
			Subtitle: "History of invoice run batches",
			Columns: ListColumnLabels{
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
			Empty: ListEmptyLabels{
				Pending: ListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "There are no invoice runs currently in progress.",
				},
				Complete: ListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "No invoice runs have completed yet.",
				},
				Failed: ListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "No invoice runs have failed.",
				},
			},
			Filters: ListFilterLabels{
				Pending:  "Pending",
				Complete: "Complete",
				Failed:   "Failed",
			},
		},
		Detail: DetailLabels{
			Title: "Invoice Run",
			Tabs: DetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Invoices:     "Invoices",
				AuditHistory: "Audit History",
				Attachments:  "Attachments",
			},
			Summary: SummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some invoices may be missing.",
			},
			Selections: SelectionsTabLabels{
				ColSubscription: "Subscription",
				ColPeriodStart:  "Period start",
				ColPeriodEnd:    "Period end",
				ColPeriodMarker: "Period marker",
				EmptyTitle:      "No selections",
				EmptyMessage:    "This run has no attempt records.",
			},
			Results: ResultsTabLabels{
				ColSubscription: "Subscription",
				ColPeriodStart:  "Period start",
				ColPeriodEnd:    "Period end",
				ColOutcome:      "Outcome",
				ColErrorCode:    "Error code",
				EmptyTitle:      "No results",
				EmptyMessage:    "This run has no attempt records.",
			},
			Invoices: InvoicesTabLabels{
				ColReference:           "Reference",
				ColDate:                "Date",
				ColAmount:              "Amount",
				ColStatus:              "Status",
				EmptyTitle:             "No invoices",
				EmptyMessage:           "No invoices were created by this run.",
				AuditHistoryComingSoon: "Audit history is coming soon.",
			},
		},
		StatusBadges: StatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		Actions: ActionLabels{
			Run:                   "Run",
			ReRunFailed:           "Re-run failed",
			ReRunFailedComingSoon: "Re-run failed (coming soon)",
			ViewRun:               "View run",
			ViewClient:            "View client",
			ViewSubscription:      "View subscription",
		},
		ScopeKind: ScopeKindLabels{
			Subscription: "Subscription",
			Client:       "Client",
			Workspace:    "Workspace",
		},
		AttemptOutcome: OutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Errors: ErrorLabels{
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
