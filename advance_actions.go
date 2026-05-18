package centymo

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
