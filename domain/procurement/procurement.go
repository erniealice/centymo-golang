package procurement

// procurement.go — procurement-domain composition surface (centymo W7).
//
// Holds the Procurement Operations composition app's label + route contract.
// The bare "procurement" view dir owns NO esqyma proto entity — it is a
// composition surface over the procurement domain's entities (mirrors the
// schedule/cyta pattern). It lives in domain/procurement/ by domain-name match.
//
// Extracted verbatim from the root labels.go (ProcurementLabels) and
// routes_config.go (ProcurementRoutes + Default*/RouteMap) per the domain-first
// restructure. Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app labels
// (composition surface, no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementLabels holds all translatable strings for the Procurement
// Operations composition app. Populated via lyngua (P4). These keys are
// intentionally generic so they render without overrides when lyngua has not
// yet supplied values.
type ProcurementLabels struct {
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

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app routes
// (composition surface; no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementRoutes holds the URL constants for the Procurement Operations app.
// service-admin composition (P3c) wires them into
// SidebarRoutes.Operations.Procurement.
type ProcurementRoutes struct {
	// Dashboard
	DashboardURL string `json:"dashboard_url"`

	// Contract operations (views over SupplierContract)
	RenewalCalendarURL string `json:"renewal_calendar_url"`
	VarianceURL        string `json:"variance_url"`
	UtilizationURL     string `json:"utilization_url"`

	// Recurrence drafts queue (lights up when P5 ships the recurrence engine)
	RecurrenceDraftsURL string `json:"recurrence_drafts_url"`
}

// DefaultProcurementRoutes returns a ProcurementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultProcurementRoutes() ProcurementRoutes {
	return ProcurementRoutes{
		DashboardURL:        ProcurementDashboardURL,
		RenewalCalendarURL:  ProcurementRenewalCalendarURL,
		VarianceURL:         ProcurementVarianceURL,
		UtilizationURL:      ProcurementUtilizationURL,
		RecurrenceDraftsURL: ProcurementRecurrenceDraftsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// procurement operations app routes.
func (r ProcurementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"procurement.dashboard":         r.DashboardURL,
		"procurement.renewals":          r.RenewalCalendarURL,
		"procurement.variance":          r.VarianceURL,
		"procurement.utilization":       r.UtilizationURL,
		"procurement.recurrence_drafts": r.RecurrenceDraftsURL,
	}
}
