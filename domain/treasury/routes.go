package treasury

// routes.go — treasury-domain route constants + Routes config structs (centymo W5).
//
// Extracted from the root routes.go (URL consts) + routes_config.go
// (Collection/Disbursement/TreasuryAdvances Routes types, Default* constructors,
// RouteMap methods) into the treasury domain package per the domain-first
// restructure. Pure structural move — no behaviour change; the route strings are
// byte-identical.
//
// NOTE: DefaultTreasuryAdvancesRoutes inlines the three SupplierBillingEvent*
// route strings (list / detail / recognize). Those URL constants are
// expenditure-domain (esqyma domainOf == expenditure) and remain at the centymo
// root pending W6; inlining the literals here avoids a treasury->root import
// while preserving the exact same values.

// Default route constants for treasury views.
// Consumer apps can use these or define their own.
const (
	// Collection (money IN) routes
	CollectionListURL             = "/collections/list/{status}"
	CollectionDetailURL           = "/collections/detail/{id}"
	CollectionDashboardURL        = "/collections/dashboard"
	CollectionAddURL              = "/action/collection/add"
	CollectionEditURL             = "/action/collection/edit/{id}"
	CollectionDeleteURL           = "/action/collection/delete"
	CollectionBulkDeleteURL       = "/action/collection/bulk-delete"
	CollectionSetStatusURL        = "/action/collection/set-status"
	CollectionBulkSetStatusURL    = "/action/collection/bulk-set-status"
	CollectionTabActionURL        = "/action/collection/detail/{id}/tab/{tab}"
	CollectionAttachmentUploadURL = "/action/collection/detail/{id}/attachments/upload"
	CollectionAttachmentDeleteURL = "/action/collection/detail/{id}/attachments/delete"

	// Disbursement (money OUT) routes
	DisbursementListURL             = "/disbursements/list/{status}"
	DisbursementDetailURL           = "/disbursements/detail/{id}"
	DisbursementDashboardURL        = "/disbursements/dashboard"
	DisbursementAddURL              = "/action/disbursement/add"
	DisbursementEditURL             = "/action/disbursement/edit/{id}"
	DisbursementDeleteURL           = "/action/disbursement/delete"
	DisbursementBulkDeleteURL       = "/action/disbursement/bulk-delete"
	DisbursementSetStatusURL        = "/action/disbursement/set-status"
	DisbursementBulkSetStatusURL    = "/action/disbursement/bulk-set-status"
	DisbursementTabActionURL        = "/action/disbursement/detail/{id}/tab/{tab}"
	DisbursementAttachmentUploadURL = "/action/disbursement/detail/{id}/attachments/upload"
	DisbursementAttachmentDeleteURL = "/action/disbursement/detail/{id}/attachments/delete"

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

	// TreasuryCollection / TreasuryDisbursement Advance Schedule tab partials
	// (loaded via HTMX, sit beside info / attachments / audit / advance-schedule
	// values in the detail page's tab switch).
	TreasuryCollectionAdvanceScheduleTabURL   = "/action/collection/detail/{id}/tab/advance-schedule"
	TreasuryDisbursementAdvanceScheduleTabURL = "/action/disbursement/detail/{id}/tab/advance-schedule"

	// UNSCHEDULED workflow drawers — Settle / Refund / Cancel on both sides.
	// Verb-first to avoid Go ServeMux ambiguity with the existing edit/{id}
	// patterns at the same depth (same rationale as SubscriptionRecognizeURL).
	TreasuryCollectionSettleURL   = "/action/collection/settle/{id}"
	TreasuryCollectionRefundURL   = "/action/collection/refund/{id}"
	TreasuryCollectionCancelURL   = "/action/collection/cancel/{id}"
	TreasuryDisbursementSettleURL = "/action/disbursement/settle/{id}"
	TreasuryDisbursementRefundURL = "/action/disbursement/refund/{id}"
	TreasuryDisbursementCancelURL = "/action/disbursement/cancel/{id}"
)

// CollectionRoutes holds all route paths for collection (money IN) views
// and actions.
type CollectionRoutes struct {
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	DashboardURL     string `json:"dashboard_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Advance Cash Events (Plan B Phase 3) — UNSCHEDULED workflow drawers +
	// the Advance Schedule tab partial. Empty defaults render the actions as
	// disabled / hidden.
	AdvanceScheduleTabURL string `json:"advance_schedule_tab_url"`
	SettleURL             string `json:"settle_url"`
	RefundURL             string `json:"refund_url"`
	CancelURL             string `json:"cancel_url"`
}

// DefaultCollectionRoutes returns a CollectionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultCollectionRoutes() CollectionRoutes {
	return CollectionRoutes{
		ListURL:          CollectionListURL,
		DetailURL:        CollectionDetailURL,
		DashboardURL:     CollectionDashboardURL,
		AddURL:           CollectionAddURL,
		EditURL:          CollectionEditURL,
		DeleteURL:        CollectionDeleteURL,
		BulkDeleteURL:    CollectionBulkDeleteURL,
		SetStatusURL:     CollectionSetStatusURL,
		BulkSetStatusURL: CollectionBulkSetStatusURL,
		TabActionURL:     CollectionTabActionURL,

		AttachmentUploadURL: CollectionAttachmentUploadURL,
		AttachmentDeleteURL: CollectionAttachmentDeleteURL,

		// 20260517-advance-cash-events Plan B Phase 3.
		AdvanceScheduleTabURL: TreasuryCollectionAdvanceScheduleTabURL,
		SettleURL:             TreasuryCollectionSettleURL,
		RefundURL:             TreasuryCollectionRefundURL,
		CancelURL:             TreasuryCollectionCancelURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// collection routes.
func (r CollectionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"collection.list":            r.ListURL,
		"collection.detail":          r.DetailURL,
		"collection.dashboard":       r.DashboardURL,
		"collection.add":             r.AddURL,
		"collection.edit":            r.EditURL,
		"collection.delete":          r.DeleteURL,
		"collection.bulk_delete":     r.BulkDeleteURL,
		"collection.set_status":      r.SetStatusURL,
		"collection.bulk_set_status": r.BulkSetStatusURL,
		"collection.tab_action":      r.TabActionURL,

		"collection.attachment.upload": r.AttachmentUploadURL,
		"collection.attachment.delete": r.AttachmentDeleteURL,

		// 20260517-advance-cash-events Plan B Phase 3.
		"collection.advance_schedule_tab": r.AdvanceScheduleTabURL,
		"collection.settle":               r.SettleURL,
		"collection.refund":               r.RefundURL,
		"collection.cancel":               r.CancelURL,
	}
}

// DisbursementRoutes holds all route paths for disbursement (money OUT) views
// and actions.
type DisbursementRoutes struct {
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	DashboardURL     string `json:"dashboard_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Advance Cash Events (Plan B Phase 3) — UNSCHEDULED workflow drawers +
	// the Advance Schedule tab partial. Mirrors CollectionRoutes.
	AdvanceScheduleTabURL string `json:"advance_schedule_tab_url"`
	SettleURL             string `json:"settle_url"`
	RefundURL             string `json:"refund_url"`
	CancelURL             string `json:"cancel_url"`
}

// DefaultDisbursementRoutes returns a DisbursementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultDisbursementRoutes() DisbursementRoutes {
	return DisbursementRoutes{
		ListURL:          DisbursementListURL,
		DetailURL:        DisbursementDetailURL,
		DashboardURL:     DisbursementDashboardURL,
		AddURL:           DisbursementAddURL,
		EditURL:          DisbursementEditURL,
		DeleteURL:        DisbursementDeleteURL,
		BulkDeleteURL:    DisbursementBulkDeleteURL,
		SetStatusURL:     DisbursementSetStatusURL,
		BulkSetStatusURL: DisbursementBulkSetStatusURL,
		TabActionURL:     DisbursementTabActionURL,

		AttachmentUploadURL: DisbursementAttachmentUploadURL,
		AttachmentDeleteURL: DisbursementAttachmentDeleteURL,

		// 20260517-advance-cash-events Plan B Phase 3.
		AdvanceScheduleTabURL: TreasuryDisbursementAdvanceScheduleTabURL,
		SettleURL:             TreasuryDisbursementSettleURL,
		RefundURL:             TreasuryDisbursementRefundURL,
		CancelURL:             TreasuryDisbursementCancelURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// disbursement routes.
func (r DisbursementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"disbursement.list":            r.ListURL,
		"disbursement.detail":          r.DetailURL,
		"disbursement.dashboard":       r.DashboardURL,
		"disbursement.add":             r.AddURL,
		"disbursement.edit":            r.EditURL,
		"disbursement.delete":          r.DeleteURL,
		"disbursement.bulk_delete":     r.BulkDeleteURL,
		"disbursement.set_status":      r.SetStatusURL,
		"disbursement.bulk_set_status": r.BulkSetStatusURL,
		"disbursement.tab_action":      r.TabActionURL,

		"disbursement.attachment.upload": r.AttachmentUploadURL,
		"disbursement.attachment.delete": r.AttachmentDeleteURL,

		// 20260517-advance-cash-events Plan B Phase 3.
		"disbursement.advance_schedule_tab": r.AdvanceScheduleTabURL,
		"disbursement.settle":               r.SettleURL,
		"disbursement.refund":               r.RefundURL,
		"disbursement.cancel":               r.CancelURL,
	}
}

// ---------------------------------------------------------------------------
// 20260517-advance-cash-events Plan B Phase 3 — TreasuryAdvancesRoutes
// ---------------------------------------------------------------------------

// TreasuryAdvancesRoutes holds all route paths for the cash-app "Advances"
// section: the workspace-level dashboard plus the filtered-list URLs for
// advance Collections / advance Disbursements (which point at the existing
// list pages with the `advance_kind` filter chip pre-applied).
//
// The Settle / Refund / Cancel drawer routes live on CollectionRoutes /
// DisbursementRoutes because they are anchored on the existing detail pages
// — there is no separate "advance" entity.
type TreasuryAdvancesRoutes struct {
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

// DefaultTreasuryAdvancesRoutes returns a TreasuryAdvancesRoutes populated
// from the package-level route constants defined in routes.go.
func DefaultTreasuryAdvancesRoutes() TreasuryAdvancesRoutes {
	return TreasuryAdvancesRoutes{
		ActiveNav:                  "cash",
		DashboardURL:               AdvancesDashboardURL,
		AdvanceCollectionListURL:   AdvanceCollectionListURL,
		AdvanceDisbursementListURL: AdvanceDisbursementListURL,
		// SupplierBillingEvent* URLs are expenditure-domain consts at the
		// centymo root (W6). Inlined here as literals to avoid a treasury->root
		// import; values are byte-identical to the root routes.go consts.
		SupplierBillingEventListURL:      "/supplier-billing-events/list/{status}",
		SupplierBillingEventDetailURL:    "/supplier-billing-events/detail/{id}",
		SupplierBillingEventRecognizeURL: "/action/supplier-billing-event/recognize/{id}",
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// treasury-advances routes.
func (r TreasuryAdvancesRoutes) RouteMap() map[string]string {
	return map[string]string{
		"treasury_advances.dashboard":                 r.DashboardURL,
		"treasury_advances.advance_collection_list":   r.AdvanceCollectionListURL,
		"treasury_advances.advance_disbursement_list": r.AdvanceDisbursementListURL,
		"supplier_billing_event.list":                 r.SupplierBillingEventListURL,
		"supplier_billing_event.detail":               r.SupplierBillingEventDetailURL,
		"supplier_billing_event.recognize":            r.SupplierBillingEventRecognizeURL,
	}
}
