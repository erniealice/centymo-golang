package supplier_contract_price_schedule

// routes.go — supplier_contract_price_schedule entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (SupplierContractPriceSchedule prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// ---------------------------------------------------------------------------
	// SPS P7 — SupplierContractPriceSchedule + SupplierContractPriceScheduleLine
	// ---------------------------------------------------------------------------

	// SupplierContractPriceSchedule master routes
	ListURL             = "/supplier-contract-price-schedules/list/{status}"
	DetailURL           = "/supplier-contract-price-schedules/detail/{id}"
	AddURL              = "/action/supplier-contract-price-schedule/add"
	EditURL             = "/action/supplier-contract-price-schedule/edit/{id}"
	DeleteURL           = "/action/supplier-contract-price-schedule/delete"
	SetStatusURL        = "/action/supplier-contract-price-schedule/set-status"
	BulkSetStatusURL    = "/action/supplier-contract-price-schedule/bulk-set-status"
	TabActionURL        = "/action/supplier-contract-price-schedule/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/delete"
	ActivateURL         = "/action/supplier-contract-price-schedule/activate/{id}"
	SupersedeURL        = "/action/supplier-contract-price-schedule/supersede/{id}"

	// SupplierContractPriceScheduleLine routes (child of schedule detail)
	LineAddURL    = "/action/supplier-contract-price-schedule/{id}/lines/add"
	LineEditURL   = "/action/supplier-contract-price-schedule/{id}/lines/edit/{lid}"
	LineDeleteURL = "/action/supplier-contract-price-schedule/{id}/lines/delete"
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
	ActivateURL  string `json:"activate_url"`
	SupersedeURL string `json:"supersede_url"`

	// Schedule line actions (child entity)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultRoutes returns a
// Routes using the package-level URL constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:           "supplier-contract-price-schedules",
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
		ActivateURL:         ActivateURL,
		SupersedeURL:        SupersedeURL,
		LineAddURL:          LineAddURL,
		LineEditURL:         LineEditURL,
		LineDeleteURL:       LineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_contract_price_schedule.list":              r.ListURL,
		"supplier_contract_price_schedule.detail":            r.DetailURL,
		"supplier_contract_price_schedule.add":               r.AddURL,
		"supplier_contract_price_schedule.edit":              r.EditURL,
		"supplier_contract_price_schedule.delete":            r.DeleteURL,
		"supplier_contract_price_schedule.set_status":        r.SetStatusURL,
		"supplier_contract_price_schedule.attachment.upload": r.AttachmentUploadURL,
		"supplier_contract_price_schedule.attachment.delete": r.AttachmentDeleteURL,
		"supplier_contract_price_schedule.activate":          r.ActivateURL,
		"supplier_contract_price_schedule.supersede":         r.SupersedeURL,
		"supplier_contract_price_schedule.line.add":          r.LineAddURL,
		"supplier_contract_price_schedule.line.edit":         r.LineEditURL,
		"supplier_contract_price_schedule.line.delete":       r.LineDeleteURL,
	}
}
