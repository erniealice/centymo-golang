package treasuryadvancesdashboard

// routes.go — advances-dashboard route constants + Routes config struct (centymo W5).
//
// Extracted from the treasury-domain routes.go into the per-view
// advancesdashboard package per the domain-first restructure. Pure structural
// move — no behaviour change; the route strings are byte-identical.
//
// NOTE: DefaultRoutes inlines the three SupplierBillingEvent* route strings
// (list / detail / recognize). Those URL constants are expenditure-domain
// (esqyma domainOf == expenditure) and remain at the centymo root pending W6;
// inlining the literals here avoids an advancesdashboard->root import while
// preserving the exact same values.

const (
	// ---------------------------------------------------------------------------
	// 20260517-advance-cash-events Plan B Phase 3 — Advance Cash Events routes.
	// "Advances" is a Cash-app section that surfaces TreasuryCollection /
	// TreasuryDisbursement rows whose advance_kind != NONE plus a workspace
	// dashboard. These are first-class operator actions (Settle / Refund /
	// Cancel) anchored on the existing TreasuryCollection / TreasuryDisbursement
	// detail pages — there is no separate "advance" entity.
	// ---------------------------------------------------------------------------

	// Advances Dashboard — workspace-level summary (both sides).
	AdvancesDashboardURL = "/cash/advances/dashboard"

	// Filtered list URLs (advance_kind != NONE) — point at the existing
	// Collection / Disbursement list pages with the chip pre-applied via a
	// query string the list page interprets. These are sidebar Href targets,
	// NOT ServeMux patterns — the list pages are registered at the underlying
	// pattern (CollectionListURL / DisbursementListURL) and read advance_kind
	// from the request query string.
	AdvanceCollectionListURL   = "/collections/list/pending?advance_kind=any"
	AdvanceDisbursementListURL = "/disbursements/list/pending?advance_kind=any"
)

// Routes holds all route paths for the cash-app "Advances" section: the
// workspace-level dashboard plus the filtered-list URLs for advance Collections
// / advance Disbursements (which point at the existing list pages with the
// `advance_kind` filter chip pre-applied).
//
// The Settle / Refund / Cancel drawer routes live on collection.Routes /
// disbursement.Routes because they are anchored on the existing detail pages
// — there is no separate "advance" entity.
type Routes struct {
	// ActiveNav is the sidebar navigation context; the Advances section sits
	// inside the Cash app so this remains "cash".
	ActiveNav string `json:"active_nav"`

	// DashboardURL is the workspace-level Advances Dashboard.
	DashboardURL string `json:"dashboard_url"`

	// AdvanceCollectionListURL / AdvanceDisbursementListURL are deep-links
	// into the existing Collection / Disbursement list pages with the
	// `advance_kind=any` chip pre-applied via the query string.
	AdvanceCollectionListURL   string `json:"advance_collection_list_url"`
	AdvanceDisbursementListURL string `json:"advance_disbursement_list_url"`

	// SupplierBillingEvent surfaces (buying-side MILESTONE anchor). Listed
	// here so the cash-app sidebar can deep-link to them, even though the
	// entity is buying-side. (Plan B Phase 7 wires the Recognize button.)
	SupplierBillingEventListURL      string `json:"supplier_billing_event_list_url"`
	SupplierBillingEventDetailURL    string `json:"supplier_billing_event_detail_url"`
	SupplierBillingEventRecognizeURL string `json:"supplier_billing_event_recognize_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:                  "cash",
		DashboardURL:               AdvancesDashboardURL,
		AdvanceCollectionListURL:   AdvanceCollectionListURL,
		AdvanceDisbursementListURL: AdvanceDisbursementListURL,
		// SupplierBillingEvent* URLs are expenditure-domain consts at the
		// centymo root (W6). Inlined here as literals to avoid an
		// advancesdashboard->root import; values are byte-identical to the root
		// routes.go consts.
		SupplierBillingEventListURL:      "/supplier-billing-events/list/{status}",
		SupplierBillingEventDetailURL:    "/supplier-billing-events/detail/{id}",
		SupplierBillingEventRecognizeURL: "/action/supplier-billing-event/recognize/{id}",
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// treasury-advances routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"treasury_advances.dashboard":                 r.DashboardURL,
		"treasury_advances.advance_collection_list":   r.AdvanceCollectionListURL,
		"treasury_advances.advance_disbursement_list": r.AdvanceDisbursementListURL,
		"supplier_billing_event.list":                 r.SupplierBillingEventListURL,
		"supplier_billing_event.detail":               r.SupplierBillingEventDetailURL,
		"supplier_billing_event.recognize":            r.SupplierBillingEventRecognizeURL,
	}
}
