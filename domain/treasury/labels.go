package treasury

// labels.go — treasury-domain label structs (centymo W5).
//
// The Collection (money IN), Disbursement (money OUT), and Advances Dashboard
// label structs that formerly lived here have been extracted into their
// per-entity packages under domain/treasury/<entity>/labels.go as part of the
// domain-first restructure:
//   - CollectionLabels        -> collection.Labels
//   - DisbursementLabels      -> disbursement.Labels
//   - AdvancesDashboardLabels -> advancesdashboard.Labels
//
// The shared advance contract (AdvanceKind / AdvanceStatus / AdvanceProrationPolicy
// enum labels, TreasuryAdvanceLabels, and the Settle/Refund/Cancel view-typed
// input/output shapes) remains in advance.go pending relocation to
// domain/treasury/shared/ by the finalize agent.
