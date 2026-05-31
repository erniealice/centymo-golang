// Package form contains template-facing data types for the per-subscription
// Invoice Run drawer (Surface C). Per the drawer-form-subpackage-convention,
// each secondary feature drawer contributes its Data + Labels types here.
// No repository imports; no Deps structs. Local types — NOT imported from
// espyna or entydad. Duplicated per package (plan rule D12).
package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
)

// Data is the top-level template context for the
// subscription-revenue-run-drawer-form template.
type Data struct {
	// FormAction is the POST URL (same URL as the GET — single endpoint).
	FormAction string
	WorkspaceID string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce      string // CSP nonce; populated by ViewAdapter.injectPageData (NonceFromContext) for inline <script nonce>
	// FragmentURL is the GET URL used by the HTMX inner-swap partial when the
	// as_of_date changes. Typically FormAction + ?partial=candidates&as_of_date=.
	FragmentURL string
	// SubscriptionID is the subscription being operated on.
	SubscriptionID string
	// SubscriptionName is the display name for the subscription.
	SubscriptionName string
	// ClientHint is the pre-formatted hint string shown beneath the
	// read-only Subscription field, e.g. "Client: Acme Corp".
	ClientHint string
	// PlanName is the price-plan name (read-only context row).
	PlanName string
	// AsOfDate is the current as-of date value (YYYY-MM-DD).
	AsOfDate string
	// MaxAsOfDate caps the date picker to today (YYYY-MM-DD).
	MaxAsOfDate string
	// EligibleCount is the number of periods eligible for invoicing.
	EligibleCount int
	// Periods is the flat list of candidate billing periods.
	// Per-sub scope means there is no group layer (single subscription).
	Periods []Period
	// AdvanceCollectionRows holds advance-Collection tranche rows for this
	// client (Plan B Phase 5b). Rendered as a separate source-kind section.
	AdvanceCollectionRows []AdvanceRow
	// CurrencyMismatch is true when the subscription currency differs from the
	// client's billing currency; triggers the mismatch alert.
	CurrencyMismatch bool
	// Currency is the ISO currency code for this subscription's periods.
	Currency string
	// Labels carries all user-facing strings for this drawer.
	Labels centymo.SubscriptionRevenueRunLabels
	// CommonLabels carries shared UI strings (Save / Cancel / etc.).
	CommonLabels pyeza.CommonLabels
}

// AdvanceRow is one advance-Collection tranche row rendered in the per-sub
// drawer's "Advance Collections" section. Plan B Phase 5b.
type AdvanceRow struct {
	AdvanceCollectionID string
	Currency            string
	PeriodStart         string
	PeriodEnd           string
	PeriodMarker        string
	PeriodLabel         string
	Amount              int64
	AmountDisplay       string
	Eligible            bool
	BlockerReason       string
	// SelectionValue encoding: "{AdvanceID}|{start}|{end}|{marker}|ADVANCE_COLLECTION"
	SelectionValue string
}

// Period is one candidate billing period row in the drawer's period table.
type Period struct {
	// SubscriptionID is repeated here for the checkbox value encoding.
	SubscriptionID string
	// PeriodStart is YYYY-MM-DD.
	PeriodStart string
	// PeriodEnd is YYYY-MM-DD.
	PeriodEnd string
	// PeriodMarker is the canonical idempotency anchor.
	PeriodMarker string
	// PeriodLabel is the human-readable range (e.g. "Jan 1 – Jan 31").
	PeriodLabel string
	// Amount is the period amount in centavos.
	Amount int64
	// AmountDisplay is the pre-formatted display string (centavos ÷ 100).
	AmountDisplay string
	// LineItemCount is the number of line items for this period.
	LineItemCount int
	// Eligible indicates the period can be invoiced.
	Eligible bool
	// BlockerReason is the human-readable explanation when Eligible=false.
	BlockerReason string
	// SelectionValue is the composite checkbox value encoding:
	// "{SubscriptionID}|{PeriodStart}|{PeriodEnd}|{PeriodMarker}"
	SelectionValue string
	// SuppressingAdvanceCollectionID is set when this cycle is overlapped by
	// an active TIME_BASED advance Collection (Decision A; Plan B Phase 5b).
	// The drawer renders the row as a greyed info-only block with a
	// "View advance" link.
	SuppressingAdvanceCollectionID string
}
