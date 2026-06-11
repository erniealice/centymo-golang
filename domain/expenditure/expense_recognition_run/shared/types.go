// Package shared holds view-typed data shapes used across the
// expense_recognition_run view module (list, detail, queue sub-packages).
// Having a dedicated leaf package (no intra-module imports) breaks the import
// cycle that would otherwise form if list/page.go or detail/page.go imported
// from the parent views/expense_recognition_run package.
//
// Mirror of packages/centymo-golang/views/revenue_run/shared/types.go.
// Plan A 20260517-expense-run Phase 4.
package shared

// ExpenseRecognitionRunRow is the view-layer representation of a single
// expense_recognition_run row. Buying-side mirror of RevenueRunRow.
type ExpenseRecognitionRunRow struct {
	ID                       string
	ScopeKind                string // "supplier" | "subscription" | "workspace"
	ScopeLabel               string // human-readable scope display name
	SupplierID               string
	SupplierName             string
	SupplierSubscriptionID   string
	SupplierSubscriptionName string
	AsOfDate                 string // YYYY-MM-DD
	InitiatedAt              string // RFC3339
	CompletedAt              string // RFC3339 or ""
	Initiator                string // workspace_user_id
	InitiatorName            string
	Status                   string // "pending" | "complete" | "failed"
	SelectionCount           int32
	CreatedCount             int32
	SkippedCount             int32
	ErroredCount             int32
	// IsStalePending is true when status=pending AND now()-initiated_at > stale threshold.
	// Computed by block.go shim using EXPENSE_RUN_PENDING_STALE_MINUTES env (default 5).
	IsStalePending bool
	Notes          string
}

// ExpenseRecognitionRunWithAttempts bundles a run and its attempt list for the
// detail page.
type ExpenseRecognitionRunWithAttempts struct {
	Run      ExpenseRecognitionRunRow
	Attempts []ExpenseRecognitionRunAttemptRow
}

// ExpenseRecognitionRunAttemptRow is the view-layer representation of a single
// run attempt.
type ExpenseRecognitionRunAttemptRow struct {
	ID                       string
	RunID                    string
	SourceKind               string // "subscription" | "advance_disbursement"
	SupplierSubscriptionID   string
	SupplierSubscriptionName string
	AdvanceDisbursementID    string
	AdvanceDisbursementName  string
	PeriodStart              string // YYYY-MM-DD
	PeriodEnd                string // YYYY-MM-DD
	PeriodMarker             string
	AttemptedAt              string // RFC3339 or ""
	Outcome                  string // "created" | "skipped" | "errored"
	ExpenseRecognitionID     string
	ExpenditureID            string
	ErrorCode                string
	ErrorMessage             string
}

// ExpenseRecognitionRow is a minimal recognition row for the Recognitions tab.
// The block.go shim populates this from the existing expense_recognition list
// use case filtered by run_id.
type ExpenseRecognitionRow struct {
	ID              string
	ReferenceNumber string
	RecognitionDate string
	TotalAmount     int64
	Currency        string
	Status          string
	SourceKind      string // "subscription" | "advance_disbursement"
	// DetailURL is the pre-built href to the expense-recognition detail page
	// (resolved by block.go). Empty when not configured.
	DetailURL string
}

// ExpenditureRow is a minimal draft-bill row for the Bills tab.
// The block.go shim populates this from the expenditure list use case filtered
// by run_id.
type ExpenditureRow struct {
	ID              string
	ReferenceNumber string
	ExpenditureDate string
	TotalAmount     int64
	Currency        string
	Status          string
	// DetailURL is the pre-built href to the expenditure detail page
	// (resolved by block.go). Empty when not configured.
	DetailURL string
}

// ListExpenseRecognitionRunsScope carries filter parameters for the list page.
type ListExpenseRecognitionRunsScope struct {
	WorkspaceID            string
	Status                 string // "" = all
	SupplierID             string // "" = all
	SupplierSubscriptionID string // "" = all
	Cursor                 string
	Limit                  int32
}

// QueueSupplierRecord is a minimal supplier row used for queue population.
type QueueSupplierRecord struct {
	ID   string
	Name string
}

// QueueCandidateInput is the minimal shape the queue needs from a candidate.
// Populated by the ListExpenseRunCandidates callback shim in block.go.
type QueueCandidateInput struct {
	SourceKind             string // "subscription" | "advance_disbursement"
	SupplierSubscriptionID string
	AdvanceDisbursementID  string
	Currency               string
	Amount                 int64
	Eligible               bool
}

// BatchRunInput is the per-supplier input for GenerateExpenseRun in the batch handler.
type BatchRunInput struct {
	SupplierID string
	AsOfDate   string
}

// BatchRunOutput is the per-supplier output from GenerateExpenseRun.
type BatchRunOutput struct {
	RunID   string
	Created int
	Skipped int
	Errored int
}
