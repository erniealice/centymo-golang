package disbursement

// routes.go — disbursement-entity route constants + Routes config struct (centymo W5).
//
// Extracted from the treasury-domain routes.go into the per-entity disbursement
// package per the domain-first restructure. Pure structural move — no behaviour
// change; the route strings are byte-identical.

// Default route constants for disbursement views.
// Consumer apps can use these or define their own.
const (
	// Disbursement (money OUT) routes
	ListURL             = "/disbursements/list/{status}"
	DetailURL           = "/disbursements/detail/{id}"
	DashboardURL        = "/disbursements/dashboard"
	AddURL              = "/action/disbursement/add"
	EditURL             = "/action/disbursement/edit/{id}"
	DeleteURL           = "/action/disbursement/delete"
	BulkDeleteURL       = "/action/disbursement/bulk-delete"
	SetStatusURL        = "/action/disbursement/set-status"
	BulkSetStatusURL    = "/action/disbursement/bulk-set-status"
	TabActionURL        = "/action/disbursement/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/disbursement/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/disbursement/detail/{id}/attachments/delete"

	// 20260517-advance-cash-events Plan B Phase 3 — buying-side advance routes.

	// TreasuryDisbursement Advance Schedule tab partial (loaded via HTMX, sits
	// beside info / attachments / audit values in the detail page's tab switch).
	TreasuryDisbursementAdvanceScheduleTabURL = "/action/disbursement/detail/{id}/tab/advance-schedule"

	// UNSCHEDULED workflow drawers — Settle / Refund / Cancel. Verb-first to
	// avoid Go ServeMux ambiguity with the existing edit/{id} patterns at the
	// same depth (same rationale as SubscriptionRecognizeURL).
	TreasuryDisbursementSettleURL = "/action/disbursement/settle/{id}"
	TreasuryDisbursementRefundURL = "/action/disbursement/refund/{id}"
	TreasuryDisbursementCancelURL = "/action/disbursement/cancel/{id}"
)

// Routes holds all route paths for disbursement (money OUT) views and actions.
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
	// the Advance Schedule tab partial. Mirrors collection Routes.
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
		AdvanceScheduleTabURL: TreasuryDisbursementAdvanceScheduleTabURL,
		SettleURL:             TreasuryDisbursementSettleURL,
		RefundURL:             TreasuryDisbursementRefundURL,
		CancelURL:             TreasuryDisbursementCancelURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// disbursement routes.
func (r Routes) RouteMap() map[string]string {
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
