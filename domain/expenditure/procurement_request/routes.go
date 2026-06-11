package procurement_request

// routes.go — procurement_request entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (ProcurementRequest prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// ---------------------------------------------------------------------------
	// P3a — ProcurementRequest + ProcurementRequestLine route constants
	// ---------------------------------------------------------------------------

	// ProcurementRequest routes
	ListURL             = "/procurement-requests/list/{status}"
	DetailURL           = "/procurement-requests/detail/{id}"
	AddURL              = "/action/procurement-request/add"
	EditURL             = "/action/procurement-request/edit/{id}"
	DeleteURL           = "/action/procurement-request/delete"
	SetStatusURL        = "/action/procurement-request/set-status"
	BulkSetStatusURL    = "/action/procurement-request/bulk-set-status"
	TabActionURL        = "/action/procurement-request/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/procurement-request/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/procurement-request/detail/{id}/attachments/delete"
	SubmitURL           = "/action/procurement-request/submit/{id}"
	ApproveURL          = "/action/procurement-request/approve/{id}"
	RejectURL           = "/action/procurement-request/reject/{id}"
	SpawnPOURL          = "/action/procurement-request/spawn-po/{id}"

	// ProcurementRequestLine routes (child of request detail)
	LineAddURL    = "/action/procurement-request/{id}/lines/add"
	LineEditURL   = "/action/procurement-request/{id}/lines/edit/{lid}"
	LineDeleteURL = "/action/procurement-request/{id}/lines/delete"

	// SPS Wave 3 — CRIT-3 spawn-retry placeholder route. Wired into the line-row
	// "Retry" button; the actual retry use case lands in a later wave so the
	// handler is currently a no-op redirect (see action/action.go::NewRetrySpawnAction).
	// NOTE: pattern uses `/retry-spawn/{lid}` (not `/{lid}/retry-spawn`) to avoid
	// stdlib ServeMux conflict with the existing `/lines/edit/{lid}` pattern.
	LineRetrySpawnURL = "/action/procurement-request/{id}/lines/retry-spawn/{lid}"
)

type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow actions
	SubmitURL  string `json:"submit_url"`
	ApproveURL string `json:"approve_url"`
	RejectURL  string `json:"reject_url"`
	SpawnPOURL string `json:"spawn_po_url"`

	// Line item actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`

	// SPS Wave 3 — CRIT-3 retry placeholder. Wired but the action use case
	// itself is intentionally out-of-scope; handler currently logs + redirects.
	LineRetrySpawnURL string `json:"line_retry_spawn_url"`
}

// DefaultRoutes returns a Routes using the
// package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:           "procurement",
		ActiveSubNav:        "draft",
		ListURL:             ListURL,
		DetailURL:           DetailURL,
		AddURL:              AddURL,
		EditURL:             EditURL,
		DeleteURL:           DeleteURL,
		SetStatusURL:        SetStatusURL,
		BulkSetStatusURL:    BulkSetStatusURL,
		TabActionURL:        TabActionURL,
		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
		SubmitURL:           SubmitURL,
		ApproveURL:          ApproveURL,
		RejectURL:           RejectURL,
		SpawnPOURL:          SpawnPOURL,
		LineAddURL:          LineAddURL,
		LineEditURL:         LineEditURL,
		LineDeleteURL:       LineDeleteURL,
		LineRetrySpawnURL:   LineRetrySpawnURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"procurement_request.list":              r.ListURL,
		"procurement_request.detail":            r.DetailURL,
		"procurement_request.add":               r.AddURL,
		"procurement_request.edit":              r.EditURL,
		"procurement_request.delete":            r.DeleteURL,
		"procurement_request.set_status":        r.SetStatusURL,
		"procurement_request.attachment.upload": r.AttachmentUploadURL,
		"procurement_request.attachment.delete": r.AttachmentDeleteURL,
		"procurement_request.submit":            r.SubmitURL,
		"procurement_request.approve":           r.ApproveURL,
		"procurement_request.reject":            r.RejectURL,
		"procurement_request.spawn_po":          r.SpawnPOURL,
		"procurement_request.line.add":          r.LineAddURL,
		"procurement_request.line.edit":         r.LineEditURL,
		"procurement_request.line.delete":       r.LineDeleteURL,
		"procurement_request.line.retry_spawn":  r.LineRetrySpawnURL,
	}
}
