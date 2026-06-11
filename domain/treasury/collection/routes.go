package collection

// routes.go — collection-entity route constants + Routes config struct (centymo W5).
//
// Extracted from the treasury-domain routes.go into the per-entity collection
// package per the domain-first restructure. Pure structural move — no behaviour
// change; the route strings are byte-identical.

// Default route constants for collection views.
// Consumer apps can use these or define their own.
const (
	// Collection (money IN) routes
	ListURL             = "/collections/list/{status}"
	DetailURL           = "/collections/detail/{id}"
	DashboardURL        = "/collections/dashboard"
	AddURL              = "/action/collection/add"
	EditURL             = "/action/collection/edit/{id}"
	DeleteURL           = "/action/collection/delete"
	BulkDeleteURL       = "/action/collection/bulk-delete"
	SetStatusURL        = "/action/collection/set-status"
	BulkSetStatusURL    = "/action/collection/bulk-set-status"
	TabActionURL        = "/action/collection/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/collection/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/collection/detail/{id}/attachments/delete"

	// 20260517-advance-cash-events Plan B Phase 3 — selling-side advance routes.

	// TreasuryCollection Advance Schedule tab partial (loaded via HTMX, sits
	// beside info / attachments / audit values in the detail page's tab switch).
	TreasuryCollectionAdvanceScheduleTabURL = "/action/collection/detail/{id}/tab/advance-schedule"

	// UNSCHEDULED workflow drawers — Settle / Refund / Cancel. Verb-first to
	// avoid Go ServeMux ambiguity with the existing edit/{id} patterns at the
	// same depth (same rationale as SubscriptionRecognizeURL).
	TreasuryCollectionSettleURL = "/action/collection/settle/{id}"
	TreasuryCollectionRefundURL = "/action/collection/refund/{id}"
	TreasuryCollectionCancelURL = "/action/collection/cancel/{id}"
)

// Routes holds all route paths for collection (money IN) views and actions.
type Routes struct {
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

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:          ListURL,
		DetailURL:        DetailURL,
		DashboardURL:     DashboardURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		TabActionURL:     TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		// 20260517-advance-cash-events Plan B Phase 3.
		AdvanceScheduleTabURL: TreasuryCollectionAdvanceScheduleTabURL,
		SettleURL:             TreasuryCollectionSettleURL,
		RefundURL:             TreasuryCollectionRefundURL,
		CancelURL:             TreasuryCollectionCancelURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// collection routes.
func (r Routes) RouteMap() map[string]string {
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
