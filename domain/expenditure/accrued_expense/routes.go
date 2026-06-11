package accrued_expense

// routes.go — accrued_expense entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (AccruedExpense prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// ---------------------------------------------------------------------------
	// SPS P10 — AccruedExpense + AccruedExpenseSettlement route constants
	// ---------------------------------------------------------------------------

	// AccruedExpense master routes (manual create drawer is secondary — primary path is AccrueFromContract use case)
	ListURL               = "/accrued-expenses/list/{status}"
	DetailURL             = "/accrued-expenses/detail/{id}"
	AddURL                = "/action/accrued-expense/add"
	EditURL               = "/action/accrued-expense/edit/{id}"
	DeleteURL             = "/action/accrued-expense/delete"
	SetStatusURL          = "/action/accrued-expense/set-status"
	BulkSetStatusURL      = "/action/accrued-expense/bulk-set-status"
	TabActionURL          = "/action/accrued-expense/detail/{id}/tab/{tab}"
	AttachmentUploadURL   = "/action/accrued-expense/detail/{id}/attachments/upload"
	AttachmentDeleteURL   = "/action/accrued-expense/detail/{id}/attachments/delete"
	SettleURL             = "/action/accrued-expense/settle/{id}"
	ReverseURL            = "/action/accrued-expense/reverse/{id}"
	AccrueFromContractURL = "/action/accrued-expense/accrue-from-contract"

	// AccruedExpenseSettlement routes (child of accrual detail — inline CRUD)
	SettlementAddURL    = "/action/accrued-expense/{id}/settlements/add"
	SettlementEditURL   = "/action/accrued-expense/{id}/settlements/edit/{sid}"
	SettlementDeleteURL = "/action/accrued-expense/{id}/settlements/delete"
)

type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachments
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	SettleURL             string `json:"settle_url"`
	ReverseURL            string `json:"reverse_url"`
	AccrueFromContractURL string `json:"accrue_from_contract_url"`

	// Settlement actions (child entity — inline CRUD)
	SettlementAddURL    string `json:"settlement_add_url"`
	SettlementEditURL   string `json:"settlement_edit_url"`
	SettlementDeleteURL string `json:"settlement_delete_url"`
}

// DefaultRoutes returns an Routes using the
// package-level URL constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:             "accrued-expenses",
		ActiveSubNav:          "outstanding",
		ListURL:               ListURL,
		DetailURL:             DetailURL,
		AddURL:                AddURL,
		EditURL:               EditURL,
		DeleteURL:             DeleteURL,
		SetStatusURL:          SetStatusURL,
		BulkSetStatusURL:      BulkSetStatusURL,
		TabActionURL:          TabActionURL,
		AttachmentUploadURL:   AttachmentUploadURL,
		AttachmentDeleteURL:   AttachmentDeleteURL,
		SettleURL:             SettleURL,
		ReverseURL:            ReverseURL,
		AccrueFromContractURL: AccrueFromContractURL,
		SettlementAddURL:      SettlementAddURL,
		SettlementEditURL:     SettlementEditURL,
		SettlementDeleteURL:   SettlementDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"accrued_expense.list":                 r.ListURL,
		"accrued_expense.detail":               r.DetailURL,
		"accrued_expense.add":                  r.AddURL,
		"accrued_expense.edit":                 r.EditURL,
		"accrued_expense.delete":               r.DeleteURL,
		"accrued_expense.set_status":           r.SetStatusURL,
		"accrued_expense.attachment.upload":    r.AttachmentUploadURL,
		"accrued_expense.attachment.delete":    r.AttachmentDeleteURL,
		"accrued_expense.settle":               r.SettleURL,
		"accrued_expense.reverse":              r.ReverseURL,
		"accrued_expense.accrue_from_contract": r.AccrueFromContractURL,
		"accrued_expense.settlement.add":       r.SettlementAddURL,
		"accrued_expense.settlement.edit":      r.SettlementEditURL,
		"accrued_expense.settlement.delete":    r.SettlementDeleteURL,
	}
}
