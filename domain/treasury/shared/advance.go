// Package shared holds the treasury-domain (and cross-domain) shared advance
// contract: the AdvanceKind/Status/ProrationPolicy enum labels, the
// TreasuryAdvance label structs, and the view-typed Settle/Refund/Cancel/
// RecognizeMilestone input/output shapes used by BOTH the treasury_collection
// (selling) and treasury_disbursement (buying) entities — and referenced
// cross-domain by subscription + expenditure billing-event recognize actions.
// It is a charter'd domain/<d>/shared/ leaf (rule-of-three: >=2 entity owners),
// imported DIRECTLY by sibling entity packages (never via the treasury facade —
// that would be an import cycle).
package shared

// advance.go — treasury-domain shared advance contract (centymo W5).
//
// Consolidates the two formerly-stranded root files into the treasury domain
// package per the domain-first restructure (Decision D, structure-mantra §8):
//   - advance_actions.go : view-typed Settle/Refund/Cancel/RecognizeMilestone
//                          input/output shapes shared by treasury_collection
//                          and treasury_disbursement action handlers.
//   - advance_labels.go (L1-181 + their Default* constructors) : AdvanceKind /
//                          AdvanceStatus / AdvanceProrationPolicy enum labels,
//                          the "Advance Schedule" tab (TreasuryAdvanceLabels),
//                          and the cash-app "Advances Dashboard" labels.
//
// The SupplierBillingEvent* labels that historically also lived in
// advance_labels.go are expenditure-domain (esqyma domainOf) and remain at the
// centymo root pending W6.
//
// Pure structural move — no behaviour change.

// advance_actions.go — view-typed input/output shapes for the UNSCHEDULED
// advance workflow drawers (Settle / Refund / Cancel) shared by the
// treasury_collection and treasury_disbursement action handlers.
//
// These shapes live at the centymo package root (not under block/) because
// the per-package view modules (collection/action, disbursement/action)
// import centymo for Routes + Labels + helpers and would otherwise need to
// import the block sub-package — which is the wrong direction in the centymo
// dep graph. The block layer's UseCases struct holds matching
// AdvanceSettleInput / AdvanceRefundInput / AdvanceCancelInput types; the
// service-admin adapter is the single place that translates between the two.

// AdvanceSettleViewInput is the per-drawer Settle submit shape.
type AdvanceSettleViewInput struct {
	AdvanceID       string
	Amount          int64
	TargetAccountID string
	Reason          string
}

// AdvanceSettleViewOutput is the Settle result the drawer renders into a
// toast / hx-trigger event.
type AdvanceSettleViewOutput struct {
	NewRemainingAmount  int64
	NewRecognizedAmount int64
	NewStatus           string
}

// AdvanceRefundViewInput is the per-drawer Refund submit shape.
type AdvanceRefundViewInput struct {
	AdvanceID          string
	Amount             int64
	RefundMethod       string
	DestinationAccount string
	Reason             string
}

// AdvanceRefundViewOutput is the Refund result.
type AdvanceRefundViewOutput struct {
	NewRemainingAmount int64
	NewStatus          string
}

// AdvanceCancelViewInput is the per-drawer Cancel submit shape.
type AdvanceCancelViewInput struct {
	AdvanceID string
	Reason    string
}

// AdvanceCancelViewOutput is the Cancel result.
type AdvanceCancelViewOutput struct {
	NewStatus string
}

// AdvanceRecognizeMilestoneInput captures the operator-supplied fields the
// MILESTONE recognize POST sends to the workflow closure.
//
// Symmetric across selling (BillingEvent) and buying (SupplierBillingEvent)
// sides — EventID semantics flip per side. The view-typed closure on the
// block.TreasuryAdvancesUseCases is responsible for routing to the right
// espyna use case based on the closure binding.
//
// 20260517-advance-cash-events Plan B Phase 7.
type AdvanceRecognizeMilestoneInput struct {
	// AdvanceID is the treasury_collection_id (selling) or
	// treasury_disbursement_id (buying).
	AdvanceID string
	// EventID is the billing_event_id (selling) or supplier_billing_event_id
	// (buying).
	EventID string
}

// AdvanceRecognizeMilestoneOutput captures the post-state for the toast
// banner. RecognitionID is either Revenue.id (selling) or
// ExpenseRecognition.id (buying).
type AdvanceRecognizeMilestoneOutput struct {
	Outcome             string // "CREATED" | "SKIPPED" | "ERRORED"
	RecognitionID       string
	ConflictingID       string
	NewRemainingAmount  int64
	NewRecognizedAmount int64
	NewStatus           string // ACTIVE | FULLY_RECOGNIZED | FULLY_AMORTIZED
	TrancheAmount       int64
}

// === Enum labels (AdvanceKind / AdvanceStatus / AdvanceProrationPolicy) ===

// AdvanceKindLabels holds the 5 operator-facing strings for the
// esqyma AdvanceKind enum (NONE / TIME_BASED / BURN_DOWN / MILESTONE /
// UNSCHEDULED). BURN_DOWN is reserved-but-disabled in v1; keep the label so
// audit/import UIs can still display historical rows if any ever leak.
type AdvanceKindLabels struct {
	None        string `json:"none"`
	TimeBased   string `json:"time_based"`
	BurnDown    string `json:"burn_down"`
	Milestone   string `json:"milestone"`
	Unscheduled string `json:"unscheduled"`
}

// AdvanceStatusLabels holds operator-facing strings for the AdvanceStatus
// enum: covers both buying (fullyAmortized) and selling (fullyRecognized)
// terminal states plus UNSCHEDULED-specific settled / refunded / cancelled
// states.
type AdvanceStatusLabels struct {
	Active           string `json:"active"`
	FullyRecognized  string `json:"fully_recognized"`
	FullyAmortized   string `json:"fully_amortized"`
	FullyDrawn       string `json:"fully_drawn"`
	Settled          string `json:"settled"`
	PartiallySettled string `json:"partially_settled"`
	Refunded         string `json:"refunded"`
	Cancelled        string `json:"cancelled"`
	Expired          string `json:"expired"`
}

// AdvanceProrationPolicyLabels holds the 3 enabled AdvanceProrationPolicy
// values. UNSPECIFIED is normalized to FULL_TRANCHE at the view layer and
// therefore never rendered (see Decision 13 in the plan).
type AdvanceProrationPolicyLabels struct {
	DayProrated     string `json:"day_prorated"`
	FullTranche     string `json:"full_tranche"`
	NextPeriodStart string `json:"next_period_start"`
}

// AdvanceEnumLabels bundles the three enum label structs so a view can
// pass one field downstream. JSON key `labels` matches the root structure
// of `advance_kind.json` once the `advanceKind.` prefix is stripped by the
// lyngua loader.
type AdvanceEnumLabels struct {
	Kind            AdvanceKindLabels            `json:"kind"`
	Status          AdvanceStatusLabels          `json:"status"`
	ProrationPolicy AdvanceProrationPolicyLabels `json:"proration_policy"`
}

// AdvanceKindRootLabels matches the JSON root in advance_kind.json
// (`{"advanceKind":{"labels":{...}}}`). The lyngua loader is pointed at the
// `advanceKind` key, leaving `labels` as the single field on this struct.
type AdvanceKindRootLabels struct {
	Labels AdvanceEnumLabels `json:"labels"`
}

// === Treasury Advance Schedule tab (selling + buying) ===

// TreasuryAdvanceActionLabels — operator-facing labels for the UNSCHEDULED
// Settle / Refund / Cancel drawers that live alongside the Advance Schedule
// tab. JSON shape mirrors `advance.actions.*` in treasury_collection.json
// and treasury_disbursement.json.
type TreasuryAdvanceActionLabels struct {
	Settle                  string `json:"settle"`
	Refund                  string `json:"refund"`
	Cancel                  string `json:"cancel"`
	SettleConfirm           string `json:"settle_confirm"`
	RefundConfirm           string `json:"refund_confirm"`
	CancelConfirm           string `json:"cancel_confirm"`
	ReasonField             string `json:"reason_field"`
	AmountField             string `json:"amount_field"`
	TargetAccountField      string `json:"target_account_field"`
	RefundMethodField       string `json:"refund_method_field"`
	DestinationAccountField string `json:"destination_account_field"`
}

// TreasuryAdvanceLabels holds the strings rendered inside the "Advance
// Schedule" tab that the TreasuryCollection (selling-side) and
// TreasuryDisbursement (buying-side) detail pages share. The same struct
// shape is used for both sides; the two sides supply different defaults
// for BalanceAccount + TargetAccount + DashboardCard via
// DefaultTreasuryCollectionAdvanceLabels / DefaultTreasuryDisbursementAdvanceLabels.
type TreasuryAdvanceLabels struct {
	Tab                  string `json:"tab"`
	DashboardCard        string `json:"dashboard_card"`
	TotalLabel           string `json:"total_label"`
	RemainingLabel       string `json:"remaining_label"`
	RecognizedLabel      string `json:"recognized_label"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	PeriodCount          string `json:"period_count"`
	PeriodUnit           string `json:"period_unit"`
	Tranches             string `json:"tranches"`
	BalanceAccount       string `json:"balance_account"`
	TargetAccount        string `json:"target_account"`
	KindField            string `json:"kind_field"`
	StatusField          string `json:"status_field"`
	ProrationPolicyField string `json:"proration_policy_field"`
	// 20260517 — short metadata-grid labels (vs the longer "Advance kind" /
	// "Advance status" form labels above) + linked-milestones table + empty
	// state + actions section heading.
	KindShort                  string                      `json:"kind_short"`
	StatusShort                string                      `json:"status_short"`
	ActionsSection             string                      `json:"actions_section"`
	LinkedMilestones           string                      `json:"linked_milestones"`
	TrancheColumn              string                      `json:"tranche_column"`
	CurrencyColumn             string                      `json:"currency_column"`
	BillingEventColumn         string                      `json:"billing_event_column"`
	SupplierBillingEventColumn string                      `json:"supplier_billing_event_column"`
	RevenueColumn              string                      `json:"revenue_column"`
	ExpenseRecognitionColumn   string                      `json:"expense_recognition_column"`
	RecognizeButton            string                      `json:"recognize_button"`
	EmptyTranchesTitle         string                      `json:"empty_tranches_title"`
	EmptyTranchesDesc          string                      `json:"empty_tranches_desc"`
	Actions                    TreasuryAdvanceActionLabels `json:"actions"`
}

// === Defaults ===

// DefaultAdvanceKindLabels returns English defaults for AdvanceKind.
func DefaultAdvanceKindLabels() AdvanceKindLabels {
	return AdvanceKindLabels{
		None:        "None",
		TimeBased:   "Time-based",
		BurnDown:    "Burn-down",
		Milestone:   "Milestone",
		Unscheduled: "Unscheduled",
	}
}

// DefaultAdvanceStatusLabels returns English defaults for AdvanceStatus.
func DefaultAdvanceStatusLabels() AdvanceStatusLabels {
	return AdvanceStatusLabels{
		Active:           "Active",
		FullyRecognized:  "Fully recognized",
		FullyAmortized:   "Fully amortized",
		FullyDrawn:       "Fully drawn",
		Settled:          "Settled",
		PartiallySettled: "Partially settled",
		Refunded:         "Refunded",
		Cancelled:        "Cancelled",
		Expired:          "Expired",
	}
}

// DefaultAdvanceProrationPolicyLabels returns English defaults. UNSPECIFIED
// is intentionally absent (normalized to FullTranche at the view layer).
func DefaultAdvanceProrationPolicyLabels() AdvanceProrationPolicyLabels {
	return AdvanceProrationPolicyLabels{
		DayProrated:     "Day-prorated",
		FullTranche:     "Full tranche",
		NextPeriodStart: "Next period start",
	}
}

// DefaultAdvanceEnumLabels bundles the three enum-label defaults.
func DefaultAdvanceEnumLabels() AdvanceEnumLabels {
	return AdvanceEnumLabels{
		Kind:            DefaultAdvanceKindLabels(),
		Status:          DefaultAdvanceStatusLabels(),
		ProrationPolicy: DefaultAdvanceProrationPolicyLabels(),
	}
}

// DefaultAdvanceKindRootLabels returns the labels in their JSON-root shape.
// Pair with lyngua's LoadPath("...", "advance_kind.json", "advanceKind", ...).
func DefaultAdvanceKindRootLabels() AdvanceKindRootLabels {
	return AdvanceKindRootLabels{Labels: DefaultAdvanceEnumLabels()}
}

func defaultTreasuryAdvanceActionLabels() TreasuryAdvanceActionLabels {
	return TreasuryAdvanceActionLabels{
		Settle:                  "Settle",
		Refund:                  "Refund",
		Cancel:                  "Cancel",
		SettleConfirm:           "Settle advance",
		RefundConfirm:           "Refund advance",
		CancelConfirm:           "Cancel advance",
		ReasonField:             "Reason",
		AmountField:             "Amount",
		TargetAccountField:      "Target account",
		RefundMethodField:       "Refund method",
		DestinationAccountField: "Destination account",
	}
}

// DefaultTreasuryCollectionAdvanceLabels — selling-side defaults
// (liability balance account / revenue target account).
func DefaultTreasuryCollectionAdvanceLabels() TreasuryAdvanceLabels {
	return TreasuryAdvanceLabels{
		Tab:                        "Advance Schedule",
		DashboardCard:              "Advance Collections",
		TotalLabel:                 "Total",
		RemainingLabel:             "Remaining",
		RecognizedLabel:            "Recognized",
		StartDate:                  "Start",
		EndDate:                    "End",
		PeriodCount:                "Periods",
		PeriodUnit:                 "Unit",
		Tranches:                   "Tranches",
		BalanceAccount:             "Liability account",
		TargetAccount:              "Revenue account",
		KindField:                  "Advance kind",
		StatusField:                "Advance status",
		ProrationPolicyField:       "Proration policy",
		KindShort:                  "Kind",
		StatusShort:                "Status",
		ActionsSection:             "Actions",
		LinkedMilestones:           "Linked Milestones",
		TrancheColumn:              "Tranche",
		CurrencyColumn:             "Currency",
		BillingEventColumn:         "Billing event",
		SupplierBillingEventColumn: "Supplier billing event",
		RevenueColumn:              "Revenue",
		ExpenseRecognitionColumn:   "Expense recognition",
		RecognizeButton:            "Recognize",
		EmptyTranchesTitle:         "No tranches recognized yet",
		EmptyTranchesDesc:          "Tranches appear here as each period is recognized.",
		Actions:                    defaultTreasuryAdvanceActionLabels(),
	}
}

// DefaultTreasuryDisbursementAdvanceLabels — buying-side defaults
// (asset balance account / expense target account).
func DefaultTreasuryDisbursementAdvanceLabels() TreasuryAdvanceLabels {
	return TreasuryAdvanceLabels{
		Tab:                        "Advance Schedule",
		DashboardCard:              "Advance Disbursements",
		TotalLabel:                 "Total",
		RemainingLabel:             "Remaining",
		RecognizedLabel:            "Recognized",
		StartDate:                  "Start",
		EndDate:                    "End",
		PeriodCount:                "Periods",
		PeriodUnit:                 "Unit",
		Tranches:                   "Tranches",
		BalanceAccount:             "Asset (prepaid) account",
		TargetAccount:              "Expense account",
		KindField:                  "Advance kind",
		StatusField:                "Advance status",
		ProrationPolicyField:       "Proration policy",
		KindShort:                  "Kind",
		StatusShort:                "Status",
		ActionsSection:             "Actions",
		LinkedMilestones:           "Linked Milestones",
		TrancheColumn:              "Tranche",
		CurrencyColumn:             "Currency",
		BillingEventColumn:         "Billing event",
		SupplierBillingEventColumn: "Supplier billing event",
		RevenueColumn:              "Revenue",
		ExpenseRecognitionColumn:   "Expense recognition",
		RecognizeButton:            "Recognize",
		EmptyTranchesTitle:         "No tranches recognized yet",
		EmptyTranchesDesc:          "Tranches appear here as each period is recognized.",
		Actions:                    defaultTreasuryAdvanceActionLabels(),
	}
}
