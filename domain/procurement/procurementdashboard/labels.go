package procurementdashboard

// labels.go — Procurement Operations composition-surface labels (centymo W7).
//
// The procurementdashboard package owns NO esqyma proto entity — it is a
// composition surface over the procurement domain's entities (mirrors the
// schedule/cyta pattern). It lives in domain/procurement/procurementdashboard/
// by domain-name match.
//
// Extracted verbatim from the procurement-domain procurement.go (ProcurementLabels)
// per the domain-first restructure. Pure structural move — no behaviour change;
// identifiers are entity-local (ProcurementLabels -> Labels).

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app labels
// (composition surface, no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Procurement Operations
// composition app. Populated via lyngua (P4). These keys are intentionally
// generic so they render without overrides when lyngua has not yet supplied
// values.
type Labels struct {
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
