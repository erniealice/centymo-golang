package procurementdashboard

// routes.go — Procurement Operations composition-surface route constants +
// Routes type (centymo W7).
//
// Extracted from the procurement-domain routes.go (Procurement Operations URL
// consts) and procurement.go (ProcurementRoutes type, DefaultProcurementRoutes()
// constructor, RouteMap() method) per the domain-first restructure. Pure
// structural move — no behaviour change; route strings are byte-identical and
// identifiers are entity-local (ProcurementRoutes -> Routes).

// Default route constants for the Procurement Operations composition app.
// Consumer apps can use these or define their own.
const (
	// Procurement Operations app — all GET, read-only views
	DashboardURL        = "/procurement/dashboard"
	RenewalCalendarURL  = "/procurement/renewals"
	VarianceURL         = "/procurement/variance"
	UtilizationURL      = "/procurement/utilization"
	RecurrenceDraftsURL = "/procurement/recurrence-drafts/list/{status}"
)

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app routes
// (composition surface; no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// Routes holds the URL constants for the Procurement Operations app.
// service-admin composition (P3c) wires them into
// SidebarRoutes.Operations.Procurement.
type Routes struct {
	// Dashboard
	DashboardURL string `json:"dashboard_url"`

	// Contract operations (views over SupplierContract)
	RenewalCalendarURL string `json:"renewal_calendar_url"`
	VarianceURL        string `json:"variance_url"`
	UtilizationURL     string `json:"utilization_url"`

	// Recurrence drafts queue (lights up when P5 ships the recurrence engine)
	RecurrenceDraftsURL string `json:"recurrence_drafts_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined above.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:        DashboardURL,
		RenewalCalendarURL:  RenewalCalendarURL,
		VarianceURL:         VarianceURL,
		UtilizationURL:      UtilizationURL,
		RecurrenceDraftsURL: RecurrenceDraftsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// procurement operations app routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"procurement.dashboard":         r.DashboardURL,
		"procurement.renewals":          r.RenewalCalendarURL,
		"procurement.variance":          r.VarianceURL,
		"procurement.utilization":       r.UtilizationURL,
		"procurement.recurrence_drafts": r.RecurrenceDraftsURL,
	}
}
