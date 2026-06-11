// Package shared holds view-typed data shapes used by both the list and detail
// sub-packages of the revenue-run view module.
// Having a dedicated leaf package (no intra-module imports) breaks the import
// cycle that would otherwise form if list/page.go or detail/page.go imported
// from the parent views/revenue_run package.
package shared

// RevenueRunRow is the view-layer representation of a single revenue run.
type RevenueRunRow struct {
	ID               string
	ScopeKind        string // "subscription" | "client" | "workspace"
	ScopeLabel       string // human-readable scope display name
	ClientID         string
	ClientName       string
	SubscriptionID   string
	SubscriptionName string
	AsOfDate         string // YYYY-MM-DD
	InitiatedAt      string // RFC3339
	CompletedAt      string // RFC3339 or ""
	Initiator        string // workspace_user_id
	InitiatorName    string
	Status           string // "pending" | "complete" | "failed"
	SelectionCount   int32
	CreatedCount     int32
	SkippedCount     int32
	ErroredCount     int32
	// IsStalePending is true when status=pending AND now()-initiated_at > stale threshold.
	// Computed by block.go shim using REVENUE_RUN_PENDING_STALE_MINUTES env (default 5).
	IsStalePending bool
	Notes          string
}

// RevenueRunWithAttempts bundles a run and its attempt list for the detail page.
type RevenueRunWithAttempts struct {
	Run      RevenueRunRow
	Attempts []RevenueRunAttemptRow
}

// RevenueRunAttemptRow is the view-layer representation of a single run attempt.
type RevenueRunAttemptRow struct {
	ID               string
	RunID            string
	SubscriptionID   string
	SubscriptionName string
	PeriodStart      string // YYYY-MM-DD
	PeriodEnd        string // YYYY-MM-DD
	PeriodMarker     string
	AttemptedAt      string // RFC3339 or ""
	Outcome          string // "created" | "skipped" | "errored"
	RevenueID        string
	RevenueReference string
	ErrorCode        string
	ErrorMessage     string
}

// RevenueRow is a minimal invoice row for the Invoices tab. The block.go shim
// populates this from the existing revenue list use case filtered by run_id.
type RevenueRow struct {
	ID              string
	ReferenceNumber string
	RevenueDate     string
	TotalAmount     int64
	Currency        string
	Status          string
	// DetailURL is the pre-built href to the revenue detail page (resolved by block.go).
	DetailURL string
}

// ListRevenueRunsScope carries filter parameters for the list page.
type ListRevenueRunsScope struct {
	WorkspaceID    string
	Status         string // "" = all
	ClientID       string // "" = all
	SubscriptionID string // "" = all
	Cursor         string
	Limit          int32
}
