package treasury

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
	TimeBased   string `json:"timeBased"`
	BurnDown    string `json:"burnDown"`
	Milestone   string `json:"milestone"`
	Unscheduled string `json:"unscheduled"`
}

// AdvanceStatusLabels holds operator-facing strings for the AdvanceStatus
// enum: covers both buying (fullyAmortized) and selling (fullyRecognized)
// terminal states plus UNSCHEDULED-specific settled / refunded / cancelled
// states.
type AdvanceStatusLabels struct {
	Active           string `json:"active"`
	FullyRecognized  string `json:"fullyRecognized"`
	FullyAmortized   string `json:"fullyAmortized"`
	FullyDrawn       string `json:"fullyDrawn"`
	Settled          string `json:"settled"`
	PartiallySettled string `json:"partiallySettled"`
	Refunded         string `json:"refunded"`
	Cancelled        string `json:"cancelled"`
	Expired          string `json:"expired"`
}

// AdvanceProrationPolicyLabels holds the 3 enabled AdvanceProrationPolicy
// values. UNSPECIFIED is normalized to FULL_TRANCHE at the view layer and
// therefore never rendered (see Decision 13 in the plan).
type AdvanceProrationPolicyLabels struct {
	DayProrated     string `json:"dayProrated"`
	FullTranche     string `json:"fullTranche"`
	NextPeriodStart string `json:"nextPeriodStart"`
}

// AdvanceEnumLabels bundles the three enum label structs so a view can
// pass one field downstream. JSON key `labels` matches the root structure
// of `advance_kind.json` once the `advanceKind.` prefix is stripped by the
// lyngua loader.
type AdvanceEnumLabels struct {
	Kind            AdvanceKindLabels            `json:"kind"`
	Status          AdvanceStatusLabels          `json:"status"`
	ProrationPolicy AdvanceProrationPolicyLabels `json:"prorationPolicy"`
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
	SettleConfirm           string `json:"settleConfirm"`
	RefundConfirm           string `json:"refundConfirm"`
	CancelConfirm           string `json:"cancelConfirm"`
	ReasonField             string `json:"reasonField"`
	AmountField             string `json:"amountField"`
	TargetAccountField      string `json:"targetAccountField"`
	RefundMethodField       string `json:"refundMethodField"`
	DestinationAccountField string `json:"destinationAccountField"`
}

// TreasuryAdvanceLabels holds the strings rendered inside the "Advance
// Schedule" tab that the TreasuryCollection (selling-side) and
// TreasuryDisbursement (buying-side) detail pages share. The same struct
// shape is used for both sides; the two sides supply different defaults
// for BalanceAccount + TargetAccount + DashboardCard via
// DefaultTreasuryCollectionAdvanceLabels / DefaultTreasuryDisbursementAdvanceLabels.
type TreasuryAdvanceLabels struct {
	Tab                  string `json:"tab"`
	DashboardCard        string `json:"dashboardCard"`
	TotalLabel           string `json:"totalLabel"`
	RemainingLabel       string `json:"remainingLabel"`
	RecognizedLabel      string `json:"recognizedLabel"`
	StartDate            string `json:"startDate"`
	EndDate              string `json:"endDate"`
	PeriodCount          string `json:"periodCount"`
	PeriodUnit           string `json:"periodUnit"`
	Tranches             string `json:"tranches"`
	BalanceAccount       string `json:"balanceAccount"`
	TargetAccount        string `json:"targetAccount"`
	KindField            string `json:"kindField"`
	StatusField          string `json:"statusField"`
	ProrationPolicyField string `json:"prorationPolicyField"`
	// 20260517 — short metadata-grid labels (vs the longer "Advance kind" /
	// "Advance status" form labels above) + linked-milestones table + empty
	// state + actions section heading.
	KindShort                  string                      `json:"kindShort"`
	StatusShort                string                      `json:"statusShort"`
	ActionsSection             string                      `json:"actionsSection"`
	LinkedMilestones           string                      `json:"linkedMilestones"`
	TrancheColumn              string                      `json:"trancheColumn"`
	CurrencyColumn             string                      `json:"currencyColumn"`
	BillingEventColumn         string                      `json:"billingEventColumn"`
	SupplierBillingEventColumn string                      `json:"supplierBillingEventColumn"`
	RevenueColumn              string                      `json:"revenueColumn"`
	ExpenseRecognitionColumn   string                      `json:"expenseRecognitionColumn"`
	RecognizeButton            string                      `json:"recognizeButton"`
	EmptyTranchesTitle         string                      `json:"emptyTranchesTitle"`
	EmptyTranchesDesc          string                      `json:"emptyTranchesDesc"`
	Actions                    TreasuryAdvanceActionLabels `json:"actions"`
}

// === Advances Dashboard (cash-app workspace view) ===

// AdvancesDashboardTableLabels — column headers for the per-side
// (outflow / inflow) table on the Advances Dashboard. The selling-side
// counterparty header reads "Customer" while the buying-side reads
// "Supplier"; the Defaults functions provide both.
type AdvancesDashboardTableLabels struct {
	ID           string `json:"id"`
	Counterparty string `json:"counterparty"`
	Kind         string `json:"kind"`
	Total        string `json:"total"`
	Remaining    string `json:"remaining"`
	Status       string `json:"status"`
}

// AdvancesDashboardSectionLabels — labels for one half (outflow OR inflow)
// of the Advances Dashboard.
type AdvancesDashboardSectionLabels struct {
	CardTitle    string                       `json:"cardTitle"`
	Table        AdvancesDashboardTableLabels `json:"table"`
	EmptyTitle   string                       `json:"emptyTitle"`
	EmptyMessage string                       `json:"emptyMessage"`
}

// AdvancesDashboardLabels — root struct for the Advances Dashboard page.
// JSON shape mirrors `advancesDashboard.*` in advances_dashboard.json.
type AdvancesDashboardLabels struct {
	Title                string                         `json:"title"`
	AsOfLabel            string                         `json:"asOfLabel"`
	TotalOutflow         string                         `json:"totalOutflow"`
	TotalInflow          string                         `json:"totalInflow"`
	UtilizationLabel     string                         `json:"utilizationLabel"`
	ActiveCount          string                         `json:"activeCount"`
	FullyRecognizedCount string                         `json:"fullyRecognizedCount"`
	Outflow              AdvancesDashboardSectionLabels `json:"outflow"`
	Inflow               AdvancesDashboardSectionLabels `json:"inflow"`
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

// DefaultAdvancesDashboardLabels returns English defaults for the
// cash-app Advances Dashboard.
func DefaultAdvancesDashboardLabels() AdvancesDashboardLabels {
	return AdvancesDashboardLabels{
		Title:                "Advances Dashboard",
		AsOfLabel:            "As of",
		TotalOutflow:         "Total prepaid (asset)",
		TotalInflow:          "Total deferred (liability)",
		UtilizationLabel:     "Utilization",
		ActiveCount:          "Active",
		FullyRecognizedCount: "Fully recognized",
		Outflow: AdvancesDashboardSectionLabels{
			CardTitle: "Outflows (Advance Disbursements)",
			Table: AdvancesDashboardTableLabels{
				ID:           "Advance",
				Counterparty: "Supplier",
				Kind:         "Kind",
				Total:        "Total",
				Remaining:    "Remaining",
				Status:       "Status",
			},
			EmptyTitle:   "No outflow advances",
			EmptyMessage: "Advance disbursements appear here as they are recorded.",
		},
		Inflow: AdvancesDashboardSectionLabels{
			CardTitle: "Inflows (Advance Collections)",
			Table: AdvancesDashboardTableLabels{
				ID:           "Advance",
				Counterparty: "Customer",
				Kind:         "Kind",
				Total:        "Total",
				Remaining:    "Remaining",
				Status:       "Status",
			},
			EmptyTitle:   "No inflow advances",
			EmptyMessage: "Advance collections appear here as they are recorded.",
		},
	}
}
