package supplier_contract

// routes.go — supplier_contract entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (SupplierContract prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// ---------------------------------------------------------------------------
	// P3a — SupplierContract + SupplierContractLine route constants
	// ---------------------------------------------------------------------------

	// SupplierContract master routes
	ListURL             = "/supplier-contracts/list/{status}"
	DetailURL           = "/supplier-contracts/detail/{id}"
	AddURL              = "/action/supplier-contract/add"
	EditURL             = "/action/supplier-contract/edit/{id}"
	DeleteURL           = "/action/supplier-contract/delete"
	SetStatusURL        = "/action/supplier-contract/set-status"
	BulkSetStatusURL    = "/action/supplier-contract/bulk-set-status"
	TabActionURL        = "/action/supplier-contract/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/supplier-contract/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/supplier-contract/detail/{id}/attachments/delete"
	ApproveURL          = "/action/supplier-contract/approve/{id}"
	TerminateURL        = "/action/supplier-contract/terminate/{id}"

	// SupplierContractLine routes (child of contract detail)
	LineAddURL    = "/action/supplier-contract/{id}/lines/add"
	LineEditURL   = "/action/supplier-contract/{id}/lines/edit/{lid}"
	LineDeleteURL = "/action/supplier-contract/{id}/lines/delete"
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

	// Workflow
	ApproveURL   string `json:"approve_url"`
	TerminateURL string `json:"terminate_url"`

	// Line item actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultRoutes returns a Routes using the
// package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:           "supplier-contracts",
		ActiveSubNav:        "active",
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
		ApproveURL:          ApproveURL,
		TerminateURL:        TerminateURL,
		LineAddURL:          LineAddURL,
		LineEditURL:         LineEditURL,
		LineDeleteURL:       LineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_contract.list":              r.ListURL,
		"supplier_contract.detail":            r.DetailURL,
		"supplier_contract.add":               r.AddURL,
		"supplier_contract.edit":              r.EditURL,
		"supplier_contract.delete":            r.DeleteURL,
		"supplier_contract.set_status":        r.SetStatusURL,
		"supplier_contract.attachment.upload": r.AttachmentUploadURL,
		"supplier_contract.attachment.delete": r.AttachmentDeleteURL,
		"supplier_contract.approve":           r.ApproveURL,
		"supplier_contract.terminate":         r.TerminateURL,
		"supplier_contract.line.add":          r.LineAddURL,
		"supplier_contract.line.edit":         r.LineEditURL,
		"supplier_contract.line.delete":       r.LineDeleteURL,
	}
}
