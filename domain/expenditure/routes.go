package expenditure

// routes.go — expenditure-domain route constants + Routes config structs
// (centymo W6). Extracted verbatim from the root routes.go (URL consts) and
// routes_config.go (Expenditure/SupplierContract/ProcurementRequest/
// SupplierContractPriceSchedule/ExpenseRecognition/AccruedExpense/
// ExpenseRecognitionRun Routes types + Default* constructors + RouteMap methods)
// per the domain-first restructure. Pure structural move — route strings are
// byte-identical.

// Default route constants for expenditure views.
const (
	// SupplierBillingEvent (buying-side MILESTONE anchor).
	SupplierBillingEventListURL      = "/supplier-billing-events/list/{status}"
	SupplierBillingEventDetailURL    = "/supplier-billing-events/detail/{id}"
	SupplierBillingEventRecognizeURL = "/action/supplier-billing-event/recognize/{id}"

	// Expense Recognition Run (buying-side) routes — Plan A 20260517-expense-run.
	ExpenseRecognitionRunQueueURL                 = "/expense-recognition-run/queue"
	ExpenseRecognitionRunQueueTableURL            = "/action/expense-recognition-run/queue/table"
	ExpenseRecognitionRunListURL                  = "/expense-recognition-run/list/{status}"
	ExpenseRecognitionRunListTableURL             = "/action/expense-recognition-run/table/{status}"
	ExpenseRecognitionRunDetailURL                = "/expense-recognition-run/detail/{id}"
	ExpenseRecognitionRunDetailTabActionURL       = "/action/expense-recognition-run/detail/{id}/tab/{tab}"
	ExpenseRecognitionRunNewURL                   = "/expense-recognition-run/new"
	ExpenseRecognitionRunGenerateURL              = "/action/expense-recognition-run/generate"
	ExpenseRecognitionRunSubmitBatchURL           = "/action/expense-recognition-run/submit-batch"
	ExpenseRecognitionRunPerSupplierDrawerURL     = "/action/supplier/expense-recognition-run/{id}"
	ExpenseRecognitionRunPerSubscriptionDrawerURL = "/action/supplier-subscription/expense-recognition-run/{id}"

	// Expenditure (purchase + expense) routes
	ExpenditurePurchaseListURL      = "/purchases/list/{status}"
	ExpenditurePurchaseDashboardURL = "/purchases/dashboard"
	ExpenditureExpenseListURL       = "/expenses/list/{status}"
	ExpenditureExpenseDashboardURL  = "/expenses/dashboard"

	// Expenditure expense CRUD action routes
	ExpenditureExpenseAddURL       = "/action/expense/add"
	ExpenditureExpenseEditURL      = "/action/expense/edit/{id}"
	ExpenditureExpenseDeleteURL    = "/action/expense/delete"
	ExpenditureExpenseSetStatusURL = "/action/expense/set-status"
	ExpenditureExpenseDetailURL    = "/expenses/detail/{id}"
	ExpenditureExpenseTableURL     = "/action/expense/table/{status}"
	ExpenditureExpenseTabActionURL = "/action/expense/detail/{id}/tab/{tab}"
	ExpenditureAttachmentUploadURL = "/action/expense/detail/{id}/attachments/upload"
	ExpenditureAttachmentDeleteURL = "/action/expense/detail/{id}/attachments/delete"

	// Expenditure expense line item action routes
	ExpenditureExpenseLineItemAddURL    = "/action/expense/detail/{id}/items/add"
	ExpenditureExpenseLineItemEditURL   = "/action/expense/detail/{id}/items/edit/{itemId}"
	ExpenditureExpenseLineItemRemoveURL = "/action/expense/detail/{id}/items/remove"
	ExpenditureExpenseLineItemTableURL  = "/action/expense/detail/{id}/items/table"

	// Expenditure pay action route (creates a pre-linked disbursement)
	ExpenditureExpensePayURL = "/action/expense/detail/{id}/pay"

	// Expenditure report routes
	PurchasesSummaryURL = "/purchases/reports/purchases-summary"
	ExpensesSummaryURL  = "/expenses/reports/expenses-summary"

	// Expenditure settings (template management) routes
	ExpenditureSettingsTemplatesURL       = "/purchases/settings/templates"
	ExpenditureSettingsTemplateUploadURL  = "/action/purchase/settings/templates/upload"
	ExpenditureSettingsTemplateDeleteURL  = "/action/purchase/settings/templates/delete"
	ExpenditureSettingsTemplateDefaultURL = "/action/purchase/settings/templates/set-default/{id}"

	// Purchase Order routes
	PurchaseOrderListURL             = "/purchase-orders/list/{status}"
	PurchaseOrderDetailURL           = "/purchase-orders/detail/{id}"
	PurchaseOrderAddURL              = "/action/purchase-order/add"
	PurchaseOrderEditURL             = "/action/purchase-order/edit/{id}"
	PurchaseOrderDeleteURL           = "/action/purchase-order/delete"
	PurchaseOrderSetStatusURL        = "/action/purchase-order/set-status"
	PurchaseOrderTableURL            = "/action/purchase-order/table/{status}"
	PurchaseOrderTabActionURL        = "/action/purchase-order/detail/{id}/tab/{tab}"
	PurchaseOrderAttachmentUploadURL = "/action/purchase-order/detail/{id}/attachments/upload"
	PurchaseOrderAttachmentDeleteURL = "/action/purchase-order/detail/{id}/attachments/delete"

	// Purchase Order line item routes (within PO detail)
	PurchaseOrderLineItemTableURL  = "/action/purchase-order/detail/{id}/items/table"
	PurchaseOrderLineItemAddURL    = "/action/purchase-order/detail/{id}/items/add"
	PurchaseOrderLineItemEditURL   = "/action/purchase-order/detail/{id}/items/edit/{itemId}"
	PurchaseOrderLineItemRemoveURL = "/action/purchase-order/detail/{id}/items/remove"

	// Purchase Order receipt action
	PurchaseOrderConfirmReceiptURL = "/action/purchase-order/{id}/confirm-receipt"

	// Expense category settings routes
	ExpenditureExpenseCategoryListURL   = "/expenses/categories/list"
	ExpenditureExpenseCategoryAddURL    = "/action/expense/categories/add"
	ExpenditureExpenseCategoryEditURL   = "/action/expense/categories/edit/{id}"
	ExpenditureExpenseCategoryDeleteURL = "/action/expense/categories/delete"
	ExpenditureExpenseCategoryTableURL  = "/action/expense/categories/table"

	// ---------------------------------------------------------------------------
	// P3a — SupplierContract + SupplierContractLine route constants
	// ---------------------------------------------------------------------------

	// SupplierContract master routes
	SupplierContractListURL             = "/supplier-contracts/list/{status}"
	SupplierContractDetailURL           = "/supplier-contracts/detail/{id}"
	SupplierContractAddURL              = "/action/supplier-contract/add"
	SupplierContractEditURL             = "/action/supplier-contract/edit/{id}"
	SupplierContractDeleteURL           = "/action/supplier-contract/delete"
	SupplierContractSetStatusURL        = "/action/supplier-contract/set-status"
	SupplierContractBulkSetStatusURL    = "/action/supplier-contract/bulk-set-status"
	SupplierContractTabActionURL        = "/action/supplier-contract/detail/{id}/tab/{tab}"
	SupplierContractAttachmentUploadURL = "/action/supplier-contract/detail/{id}/attachments/upload"
	SupplierContractAttachmentDeleteURL = "/action/supplier-contract/detail/{id}/attachments/delete"
	SupplierContractApproveURL          = "/action/supplier-contract/approve/{id}"
	SupplierContractTerminateURL        = "/action/supplier-contract/terminate/{id}"

	// SupplierContractLine routes (child of contract detail)
	SupplierContractLineAddURL    = "/action/supplier-contract/{id}/lines/add"
	SupplierContractLineEditURL   = "/action/supplier-contract/{id}/lines/edit/{lid}"
	SupplierContractLineDeleteURL = "/action/supplier-contract/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// P3a — ProcurementRequest + ProcurementRequestLine route constants
	// ---------------------------------------------------------------------------

	// ProcurementRequest routes
	ProcurementRequestListURL             = "/procurement-requests/list/{status}"
	ProcurementRequestDetailURL           = "/procurement-requests/detail/{id}"
	ProcurementRequestAddURL              = "/action/procurement-request/add"
	ProcurementRequestEditURL             = "/action/procurement-request/edit/{id}"
	ProcurementRequestDeleteURL           = "/action/procurement-request/delete"
	ProcurementRequestSetStatusURL        = "/action/procurement-request/set-status"
	ProcurementRequestBulkSetStatusURL    = "/action/procurement-request/bulk-set-status"
	ProcurementRequestTabActionURL        = "/action/procurement-request/detail/{id}/tab/{tab}"
	ProcurementRequestAttachmentUploadURL = "/action/procurement-request/detail/{id}/attachments/upload"
	ProcurementRequestAttachmentDeleteURL = "/action/procurement-request/detail/{id}/attachments/delete"
	ProcurementRequestSubmitURL           = "/action/procurement-request/submit/{id}"
	ProcurementRequestApproveURL          = "/action/procurement-request/approve/{id}"
	ProcurementRequestRejectURL           = "/action/procurement-request/reject/{id}"
	ProcurementRequestSpawnPOURL          = "/action/procurement-request/spawn-po/{id}"

	// ProcurementRequestLine routes (child of request detail)
	ProcurementRequestLineAddURL    = "/action/procurement-request/{id}/lines/add"
	ProcurementRequestLineEditURL   = "/action/procurement-request/{id}/lines/edit/{lid}"
	ProcurementRequestLineDeleteURL = "/action/procurement-request/{id}/lines/delete"

	// SPS Wave 3 — CRIT-3 spawn-retry placeholder route. Wired into the line-row
	// "Retry" button; the actual retry use case lands in a later wave so the
	// handler is currently a no-op redirect (see action/action.go::NewRetrySpawnAction).
	// NOTE: pattern uses `/retry-spawn/{lid}` (not `/{lid}/retry-spawn`) to avoid
	// stdlib ServeMux conflict with the existing `/lines/edit/{lid}` pattern.
	ProcurementRequestLineRetrySpawnURL = "/action/procurement-request/{id}/lines/retry-spawn/{lid}"

	// ---------------------------------------------------------------------------
	// SPS P7 — SupplierContractPriceSchedule + SupplierContractPriceScheduleLine
	// ---------------------------------------------------------------------------

	// SupplierContractPriceSchedule master routes
	SupplierContractPriceScheduleListURL             = "/supplier-contract-price-schedules/list/{status}"
	SupplierContractPriceScheduleDetailURL           = "/supplier-contract-price-schedules/detail/{id}"
	SupplierContractPriceScheduleAddURL              = "/action/supplier-contract-price-schedule/add"
	SupplierContractPriceScheduleEditURL             = "/action/supplier-contract-price-schedule/edit/{id}"
	SupplierContractPriceScheduleDeleteURL           = "/action/supplier-contract-price-schedule/delete"
	SupplierContractPriceScheduleSetStatusURL        = "/action/supplier-contract-price-schedule/set-status"
	SupplierContractPriceScheduleBulkSetStatusURL    = "/action/supplier-contract-price-schedule/bulk-set-status"
	SupplierContractPriceScheduleTabActionURL        = "/action/supplier-contract-price-schedule/detail/{id}/tab/{tab}"
	SupplierContractPriceScheduleAttachmentUploadURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/upload"
	SupplierContractPriceScheduleAttachmentDeleteURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/delete"
	SupplierContractPriceScheduleActivateURL         = "/action/supplier-contract-price-schedule/activate/{id}"
	SupplierContractPriceScheduleSupersedeURL        = "/action/supplier-contract-price-schedule/supersede/{id}"

	// SupplierContractPriceScheduleLine routes (child of schedule detail)
	SupplierContractPriceScheduleLineAddURL    = "/action/supplier-contract-price-schedule/{id}/lines/add"
	SupplierContractPriceScheduleLineEditURL   = "/action/supplier-contract-price-schedule/{id}/lines/edit/{lid}"
	SupplierContractPriceScheduleLineDeleteURL = "/action/supplier-contract-price-schedule/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// SPS P10 — ExpenseRecognition + ExpenseRecognitionLine route constants
	// ---------------------------------------------------------------------------

	// ExpenseRecognition master routes (no add/edit drawer — created BY use case)
	ExpenseRecognitionListURL                     = "/expense-recognitions/list/{status}"
	ExpenseRecognitionDetailURL                   = "/expense-recognitions/detail/{id}"
	ExpenseRecognitionDeleteURL                   = "/action/expense-recognition/delete"
	ExpenseRecognitionTabActionURL                = "/action/expense-recognition/detail/{id}/tab/{tab}"
	ExpenseRecognitionAttachmentUploadURL         = "/action/expense-recognition/detail/{id}/attachments/upload"
	ExpenseRecognitionAttachmentDeleteURL         = "/action/expense-recognition/detail/{id}/attachments/delete"
	ExpenseRecognitionReverseURL                  = "/action/expense-recognition/reverse/{id}"
	ExpenseRecognitionRecognizeFromExpenditureURL = "/action/expense-recognition/recognize-from-expenditure"
	ExpenseRecognitionRecognizeFromContractURL    = "/action/expense-recognition/recognize-from-contract"

	// ExpenseRecognitionLine routes (child of recognition detail — inline CRUD)
	ExpenseRecognitionLineAddURL    = "/action/expense-recognition/{id}/lines/add"
	ExpenseRecognitionLineEditURL   = "/action/expense-recognition/{id}/lines/edit/{lid}"
	ExpenseRecognitionLineDeleteURL = "/action/expense-recognition/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// SPS P10 — AccruedExpense + AccruedExpenseSettlement route constants
	// ---------------------------------------------------------------------------

	// AccruedExpense master routes (manual create drawer is secondary — primary path is AccrueFromContract use case)
	AccruedExpenseListURL               = "/accrued-expenses/list/{status}"
	AccruedExpenseDetailURL             = "/accrued-expenses/detail/{id}"
	AccruedExpenseAddURL                = "/action/accrued-expense/add"
	AccruedExpenseEditURL               = "/action/accrued-expense/edit/{id}"
	AccruedExpenseDeleteURL             = "/action/accrued-expense/delete"
	AccruedExpenseSetStatusURL          = "/action/accrued-expense/set-status"
	AccruedExpenseBulkSetStatusURL      = "/action/accrued-expense/bulk-set-status"
	AccruedExpenseTabActionURL          = "/action/accrued-expense/detail/{id}/tab/{tab}"
	AccruedExpenseAttachmentUploadURL   = "/action/accrued-expense/detail/{id}/attachments/upload"
	AccruedExpenseAttachmentDeleteURL   = "/action/accrued-expense/detail/{id}/attachments/delete"
	AccruedExpenseSettleURL             = "/action/accrued-expense/settle/{id}"
	AccruedExpenseReverseURL            = "/action/accrued-expense/reverse/{id}"
	AccruedExpenseAccrueFromContractURL = "/action/accrued-expense/accrue-from-contract"

	// AccruedExpenseSettlement routes (child of accrual detail — inline CRUD)
	AccruedExpenseSettlementAddURL    = "/action/accrued-expense/{id}/settlements/add"
	AccruedExpenseSettlementEditURL   = "/action/accrued-expense/{id}/settlements/edit/{sid}"
	AccruedExpenseSettlementDeleteURL = "/action/accrued-expense/{id}/settlements/delete"
)

// ExpenditureRoutes holds all route paths for expenditure views (purchase + expense).
type ExpenditureRoutes struct {
	PurchaseListURL      string `json:"purchase_list_url"`
	PurchaseDashboardURL string `json:"purchase_dashboard_url"`
	ExpenseListURL       string `json:"expense_list_url"`
	ExpenseDashboardURL  string `json:"expense_dashboard_url"`

	// Report routes
	PurchasesSummaryURL string `json:"purchases_summary_url"`
	ExpensesSummaryURL  string `json:"expenses_summary_url"`

	// Settings (template management)
	SettingsTemplatesURL       string `json:"settings_templates_url"`
	SettingsTemplateUploadURL  string `json:"settings_template_upload_url"`
	SettingsTemplateDeleteURL  string `json:"settings_template_delete_url"`
	SettingsTemplateDefaultURL string `json:"settings_template_default_url"`

	// Expense CRUD action routes
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	DetailURL           string `json:"detail_url"`
	TableURL            string `json:"table_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Expense line item action routes
	LineItemAddURL    string `json:"line_item_add_url"`
	LineItemEditURL   string `json:"line_item_edit_url"`
	LineItemRemoveURL string `json:"line_item_remove_url"`
	LineItemTableURL  string `json:"line_item_table_url"`

	// Pay action route (creates pre-linked disbursement)
	PayURL string `json:"pay_url"`

	// Expense category CRUD routes
	ExpenseCategoryListURL   string `json:"expense_category_list_url"`
	ExpenseCategoryAddURL    string `json:"expense_category_add_url"`
	ExpenseCategoryEditURL   string `json:"expense_category_edit_url"`
	ExpenseCategoryDeleteURL string `json:"expense_category_delete_url"`
	ExpenseCategoryTableURL  string `json:"expense_category_table_url"`

	// Purchase Order routes
	PurchaseOrderListURL             string `json:"purchase_order_list_url"`
	PurchaseOrderDetailURL           string `json:"purchase_order_detail_url"`
	PurchaseOrderAddURL              string `json:"purchase_order_add_url"`
	PurchaseOrderEditURL             string `json:"purchase_order_edit_url"`
	PurchaseOrderDeleteURL           string `json:"purchase_order_delete_url"`
	PurchaseOrderSetStatusURL        string `json:"purchase_order_set_status_url"`
	PurchaseOrderTableURL            string `json:"purchase_order_table_url"`
	PurchaseOrderTabActionURL        string `json:"purchase_order_tab_action_url"`
	PurchaseOrderAttachmentUploadURL string `json:"purchase_order_attachment_upload_url"`
	PurchaseOrderAttachmentDeleteURL string `json:"purchase_order_attachment_delete_url"`

	// Purchase Order line item routes (within PO detail)
	PurchaseOrderLineItemTableURL  string `json:"purchase_order_line_item_table_url"`
	PurchaseOrderLineItemAddURL    string `json:"purchase_order_line_item_add_url"`
	PurchaseOrderLineItemEditURL   string `json:"purchase_order_line_item_edit_url"`
	PurchaseOrderLineItemRemoveURL string `json:"purchase_order_line_item_remove_url"`

	// Purchase Order receipt action
	PurchaseOrderConfirmReceiptURL string `json:"purchase_order_confirm_receipt_url"`
}

// DefaultExpenditureRoutes returns an ExpenditureRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultExpenditureRoutes() ExpenditureRoutes {
	return ExpenditureRoutes{
		PurchaseListURL:      ExpenditurePurchaseListURL,
		PurchaseDashboardURL: ExpenditurePurchaseDashboardURL,
		ExpenseListURL:       ExpenditureExpenseListURL,
		ExpenseDashboardURL:  ExpenditureExpenseDashboardURL,

		PurchasesSummaryURL: PurchasesSummaryURL,
		ExpensesSummaryURL:  ExpensesSummaryURL,

		SettingsTemplatesURL:       ExpenditureSettingsTemplatesURL,
		SettingsTemplateUploadURL:  ExpenditureSettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  ExpenditureSettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: ExpenditureSettingsTemplateDefaultURL,

		AddURL:              ExpenditureExpenseAddURL,
		EditURL:             ExpenditureExpenseEditURL,
		DeleteURL:           ExpenditureExpenseDeleteURL,
		SetStatusURL:        ExpenditureExpenseSetStatusURL,
		DetailURL:           ExpenditureExpenseDetailURL,
		TableURL:            ExpenditureExpenseTableURL,
		TabActionURL:        ExpenditureExpenseTabActionURL,
		AttachmentUploadURL: ExpenditureAttachmentUploadURL,
		AttachmentDeleteURL: ExpenditureAttachmentDeleteURL,

		LineItemAddURL:    ExpenditureExpenseLineItemAddURL,
		LineItemEditURL:   ExpenditureExpenseLineItemEditURL,
		LineItemRemoveURL: ExpenditureExpenseLineItemRemoveURL,
		LineItemTableURL:  ExpenditureExpenseLineItemTableURL,

		PayURL: ExpenditureExpensePayURL,

		ExpenseCategoryListURL:   ExpenditureExpenseCategoryListURL,
		ExpenseCategoryAddURL:    ExpenditureExpenseCategoryAddURL,
		ExpenseCategoryEditURL:   ExpenditureExpenseCategoryEditURL,
		ExpenseCategoryDeleteURL: ExpenditureExpenseCategoryDeleteURL,
		ExpenseCategoryTableURL:  ExpenditureExpenseCategoryTableURL,

		PurchaseOrderListURL:             PurchaseOrderListURL,
		PurchaseOrderDetailURL:           PurchaseOrderDetailURL,
		PurchaseOrderAddURL:              PurchaseOrderAddURL,
		PurchaseOrderEditURL:             PurchaseOrderEditURL,
		PurchaseOrderDeleteURL:           PurchaseOrderDeleteURL,
		PurchaseOrderSetStatusURL:        PurchaseOrderSetStatusURL,
		PurchaseOrderTableURL:            PurchaseOrderTableURL,
		PurchaseOrderTabActionURL:        PurchaseOrderTabActionURL,
		PurchaseOrderAttachmentUploadURL: PurchaseOrderAttachmentUploadURL,
		PurchaseOrderAttachmentDeleteURL: PurchaseOrderAttachmentDeleteURL,

		PurchaseOrderLineItemTableURL:  PurchaseOrderLineItemTableURL,
		PurchaseOrderLineItemAddURL:    PurchaseOrderLineItemAddURL,
		PurchaseOrderLineItemEditURL:   PurchaseOrderLineItemEditURL,
		PurchaseOrderLineItemRemoveURL: PurchaseOrderLineItemRemoveURL,

		PurchaseOrderConfirmReceiptURL: PurchaseOrderConfirmReceiptURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// expenditure routes.
func (r ExpenditureRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expenditure.purchase.list":      r.PurchaseListURL,
		"expenditure.purchase.dashboard": r.PurchaseDashboardURL,
		"expenditure.expense.list":       r.ExpenseListURL,
		"expenditure.expense.dashboard":  r.ExpenseDashboardURL,

		"expenditure.purchases_summary": r.PurchasesSummaryURL,
		"expenditure.expenses_summary":  r.ExpensesSummaryURL,

		"purchases.settings.templates":        r.SettingsTemplatesURL,
		"purchases.settings.template_upload":  r.SettingsTemplateUploadURL,
		"purchases.settings.template_delete":  r.SettingsTemplateDeleteURL,
		"purchases.settings.template_default": r.SettingsTemplateDefaultURL,

		"expenditure.expense.add":               r.AddURL,
		"expenditure.expense.edit":              r.EditURL,
		"expenditure.expense.delete":            r.DeleteURL,
		"expenditure.expense.set_status":        r.SetStatusURL,
		"expenditure.expense.detail":            r.DetailURL,
		"expenditure.expense.table":             r.TableURL,
		"expenditure.expense.pay":               r.PayURL,
		"expenditure.expense.attachment.upload": r.AttachmentUploadURL,
		"expenditure.expense.attachment.delete": r.AttachmentDeleteURL,

		"expenditure.expense_category.list":   r.ExpenseCategoryListURL,
		"expenditure.expense_category.add":    r.ExpenseCategoryAddURL,
		"expenditure.expense_category.edit":   r.ExpenseCategoryEditURL,
		"expenditure.expense_category.delete": r.ExpenseCategoryDeleteURL,
		"expenditure.expense_category.table":  r.ExpenseCategoryTableURL,

		"expenditure.purchase_order.list":              r.PurchaseOrderListURL,
		"expenditure.purchase_order.detail":            r.PurchaseOrderDetailURL,
		"expenditure.purchase_order.add":               r.PurchaseOrderAddURL,
		"expenditure.purchase_order.edit":              r.PurchaseOrderEditURL,
		"expenditure.purchase_order.delete":            r.PurchaseOrderDeleteURL,
		"expenditure.purchase_order.set_status":        r.PurchaseOrderSetStatusURL,
		"expenditure.purchase_order.table":             r.PurchaseOrderTableURL,
		"expenditure.purchase_order.tab_action":        r.PurchaseOrderTabActionURL,
		"expenditure.purchase_order.attachment.upload": r.PurchaseOrderAttachmentUploadURL,
		"expenditure.purchase_order.attachment.delete": r.PurchaseOrderAttachmentDeleteURL,
		"expenditure.purchase_order.line_item.table":   r.PurchaseOrderLineItemTableURL,
		"expenditure.purchase_order.line_item.add":     r.PurchaseOrderLineItemAddURL,
		"expenditure.purchase_order.line_item.edit":    r.PurchaseOrderLineItemEditURL,
		"expenditure.purchase_order.line_item.remove":  r.PurchaseOrderLineItemRemoveURL,
		"expenditure.purchase_order.confirm_receipt":   r.PurchaseOrderConfirmReceiptURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierContractRoutes — P3a
// ---------------------------------------------------------------------------

// SupplierContractRoutes holds all route paths for supplier_contract views.
type SupplierContractRoutes struct {
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

// DefaultSupplierContractRoutes returns a SupplierContractRoutes using the
// package-level route constants.
func DefaultSupplierContractRoutes() SupplierContractRoutes {
	return SupplierContractRoutes{
		ActiveNav:           "supplier-contracts",
		ActiveSubNav:        "active",
		ListURL:             SupplierContractListURL,
		DetailURL:           SupplierContractDetailURL,
		AddURL:              SupplierContractAddURL,
		EditURL:             SupplierContractEditURL,
		DeleteURL:           SupplierContractDeleteURL,
		SetStatusURL:        SupplierContractSetStatusURL,
		BulkSetStatusURL:    SupplierContractBulkSetStatusURL,
		TabActionURL:        SupplierContractTabActionURL,
		AttachmentUploadURL: SupplierContractAttachmentUploadURL,
		AttachmentDeleteURL: SupplierContractAttachmentDeleteURL,
		ApproveURL:          SupplierContractApproveURL,
		TerminateURL:        SupplierContractTerminateURL,
		LineAddURL:          SupplierContractLineAddURL,
		LineEditURL:         SupplierContractLineEditURL,
		LineDeleteURL:       SupplierContractLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierContractRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// ProcurementRequestRoutes — P3a
// ---------------------------------------------------------------------------

// ProcurementRequestRoutes holds all route paths for procurement_request views.
type ProcurementRequestRoutes struct {
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

// DefaultProcurementRequestRoutes returns a ProcurementRequestRoutes using the
// package-level route constants.
func DefaultProcurementRequestRoutes() ProcurementRequestRoutes {
	return ProcurementRequestRoutes{
		ActiveNav:           "procurement",
		ActiveSubNav:        "draft",
		ListURL:             ProcurementRequestListURL,
		DetailURL:           ProcurementRequestDetailURL,
		AddURL:              ProcurementRequestAddURL,
		EditURL:             ProcurementRequestEditURL,
		DeleteURL:           ProcurementRequestDeleteURL,
		SetStatusURL:        ProcurementRequestSetStatusURL,
		BulkSetStatusURL:    ProcurementRequestBulkSetStatusURL,
		TabActionURL:        ProcurementRequestTabActionURL,
		AttachmentUploadURL: ProcurementRequestAttachmentUploadURL,
		AttachmentDeleteURL: ProcurementRequestAttachmentDeleteURL,
		SubmitURL:           ProcurementRequestSubmitURL,
		ApproveURL:          ProcurementRequestApproveURL,
		RejectURL:           ProcurementRequestRejectURL,
		SpawnPOURL:          ProcurementRequestSpawnPOURL,
		LineAddURL:          ProcurementRequestLineAddURL,
		LineEditURL:         ProcurementRequestLineEditURL,
		LineDeleteURL:       ProcurementRequestLineDeleteURL,
		LineRetrySpawnURL:   ProcurementRequestLineRetrySpawnURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ProcurementRequestRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// SupplierContractPriceScheduleRoutes — SPS P7
// ---------------------------------------------------------------------------

// SupplierContractPriceScheduleRoutes holds all route paths for
// supplier_contract_price_schedule + child line views.
type SupplierContractPriceScheduleRoutes struct {
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

// DefaultSupplierContractPriceScheduleRoutes returns a
// SupplierContractPriceScheduleRoutes using the package-level URL constants.
func DefaultSupplierContractPriceScheduleRoutes() SupplierContractPriceScheduleRoutes {
	return SupplierContractPriceScheduleRoutes{
		ActiveNav:           "supplier-contract-price-schedules",
		ActiveSubNav:        "active",
		ListURL:             SupplierContractPriceScheduleListURL,
		DetailURL:           SupplierContractPriceScheduleDetailURL,
		AddURL:              SupplierContractPriceScheduleAddURL,
		EditURL:             SupplierContractPriceScheduleEditURL,
		DeleteURL:           SupplierContractPriceScheduleDeleteURL,
		SetStatusURL:        SupplierContractPriceScheduleSetStatusURL,
		BulkSetStatusURL:    SupplierContractPriceScheduleBulkSetStatusURL,
		TabActionURL:        SupplierContractPriceScheduleTabActionURL,
		AttachmentUploadURL: SupplierContractPriceScheduleAttachmentUploadURL,
		AttachmentDeleteURL: SupplierContractPriceScheduleAttachmentDeleteURL,
		ActivateURL:         SupplierContractPriceScheduleActivateURL,
		SupersedeURL:        SupplierContractPriceScheduleSupersedeURL,
		LineAddURL:          SupplierContractPriceScheduleLineAddURL,
		LineEditURL:         SupplierContractPriceScheduleLineEditURL,
		LineDeleteURL:       SupplierContractPriceScheduleLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierContractPriceScheduleRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// ExpenseRecognitionRoutes — SPS P10
// ---------------------------------------------------------------------------

// ExpenseRecognitionRoutes holds all route paths for expense_recognition views.
// Note: no Add/Edit URLs — recognitions are created BY use case, not by user.
type ExpenseRecognitionRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	DeleteURL           string `json:"delete_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	ReverseURL                  string `json:"reverse_url"`
	RecognizeFromExpenditureURL string `json:"recognize_from_expenditure_url"`
	RecognizeFromContractURL    string `json:"recognize_from_contract_url"`

	// Recognition line actions (child entity — inline CRUD)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultExpenseRecognitionRoutes returns an ExpenseRecognitionRoutes using the
// package-level URL constants.
func DefaultExpenseRecognitionRoutes() ExpenseRecognitionRoutes {
	return ExpenseRecognitionRoutes{
		ActiveNav:                   "expense-recognitions",
		ActiveSubNav:                "posted",
		ListURL:                     ExpenseRecognitionListURL,
		DetailURL:                   ExpenseRecognitionDetailURL,
		DeleteURL:                   ExpenseRecognitionDeleteURL,
		TabActionURL:                ExpenseRecognitionTabActionURL,
		AttachmentUploadURL:         ExpenseRecognitionAttachmentUploadURL,
		AttachmentDeleteURL:         ExpenseRecognitionAttachmentDeleteURL,
		ReverseURL:                  ExpenseRecognitionReverseURL,
		RecognizeFromExpenditureURL: ExpenseRecognitionRecognizeFromExpenditureURL,
		RecognizeFromContractURL:    ExpenseRecognitionRecognizeFromContractURL,
		LineAddURL:                  ExpenseRecognitionLineAddURL,
		LineEditURL:                 ExpenseRecognitionLineEditURL,
		LineDeleteURL:               ExpenseRecognitionLineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ExpenseRecognitionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition.list":                       r.ListURL,
		"expense_recognition.detail":                     r.DetailURL,
		"expense_recognition.delete":                     r.DeleteURL,
		"expense_recognition.attachment.upload":          r.AttachmentUploadURL,
		"expense_recognition.attachment.delete":          r.AttachmentDeleteURL,
		"expense_recognition.reverse":                    r.ReverseURL,
		"expense_recognition.recognize_from_expenditure": r.RecognizeFromExpenditureURL,
		"expense_recognition.recognize_from_contract":    r.RecognizeFromContractURL,
		"expense_recognition.line.add":                   r.LineAddURL,
		"expense_recognition.line.edit":                  r.LineEditURL,
		"expense_recognition.line.delete":                r.LineDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// AccruedExpenseRoutes — SPS P10
// ---------------------------------------------------------------------------

// AccruedExpenseRoutes holds all route paths for accrued_expense views.
type AccruedExpenseRoutes struct {
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

// DefaultAccruedExpenseRoutes returns an AccruedExpenseRoutes using the
// package-level URL constants.
func DefaultAccruedExpenseRoutes() AccruedExpenseRoutes {
	return AccruedExpenseRoutes{
		ActiveNav:             "accrued-expenses",
		ActiveSubNav:          "outstanding",
		ListURL:               AccruedExpenseListURL,
		DetailURL:             AccruedExpenseDetailURL,
		AddURL:                AccruedExpenseAddURL,
		EditURL:               AccruedExpenseEditURL,
		DeleteURL:             AccruedExpenseDeleteURL,
		SetStatusURL:          AccruedExpenseSetStatusURL,
		BulkSetStatusURL:      AccruedExpenseBulkSetStatusURL,
		TabActionURL:          AccruedExpenseTabActionURL,
		AttachmentUploadURL:   AccruedExpenseAttachmentUploadURL,
		AttachmentDeleteURL:   AccruedExpenseAttachmentDeleteURL,
		SettleURL:             AccruedExpenseSettleURL,
		ReverseURL:            AccruedExpenseReverseURL,
		AccrueFromContractURL: AccruedExpenseAccrueFromContractURL,
		SettlementAddURL:      AccruedExpenseSettlementAddURL,
		SettlementEditURL:     AccruedExpenseSettlementEditURL,
		SettlementDeleteURL:   AccruedExpenseSettlementDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AccruedExpenseRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// ExpenseRecognitionRunRoutes — Plan A 20260517-expense-run
// ---------------------------------------------------------------------------

// ExpenseRecognitionRunRoutes holds all route paths for the Expense Recognition
// Run (buying-side) module. Mirrors RevenueRunRoutes shape.
// Surfaces: A (per-supplier drawer — entydad), B (workspace queue),
// C (per-supplier-subscription drawer), D (run history list + detail).
type ExpenseRecognitionRunRoutes struct {
	// Sidebar navigation context.
	ActiveNav string `json:"active_nav"`

	// Surface B — workspace queue page.
	QueueURL      string `json:"queue_url"`
	QueueTableURL string `json:"queue_table_url"`

	// Surface D — run history list + detail.
	ListURL            string `json:"list_url"`
	ListTableURL       string `json:"list_table_url"`
	DetailURL          string `json:"detail_url"`
	DetailTabActionURL string `json:"detail_tab_action_url"`

	// Action endpoints.
	NewURL         string `json:"new_url"`
	GenerateURL    string `json:"generate_url"`
	SubmitBatchURL string `json:"submit_batch_url"`

	// Surface A — per-supplier drawer (entydad supplier statement tab).
	PerSupplierDrawerURL string `json:"per_supplier_drawer_url"`

	// Surface C — per-SupplierSubscription drawer.
	PerSubscriptionDrawerURL string `json:"per_subscription_drawer_url"`
}

// DefaultExpenseRecognitionRunRoutes returns ExpenseRecognitionRunRoutes
// populated from the package-level route constants.
func DefaultExpenseRecognitionRunRoutes() ExpenseRecognitionRunRoutes {
	return ExpenseRecognitionRunRoutes{
		ActiveNav:                "expense-recognition-run",
		QueueURL:                 ExpenseRecognitionRunQueueURL,
		QueueTableURL:            ExpenseRecognitionRunQueueTableURL,
		ListURL:                  ExpenseRecognitionRunListURL,
		ListTableURL:             ExpenseRecognitionRunListTableURL,
		DetailURL:                ExpenseRecognitionRunDetailURL,
		DetailTabActionURL:       ExpenseRecognitionRunDetailTabActionURL,
		NewURL:                   ExpenseRecognitionRunNewURL,
		GenerateURL:              ExpenseRecognitionRunGenerateURL,
		SubmitBatchURL:           ExpenseRecognitionRunSubmitBatchURL,
		PerSupplierDrawerURL:     ExpenseRecognitionRunPerSupplierDrawerURL,
		PerSubscriptionDrawerURL: ExpenseRecognitionRunPerSubscriptionDrawerURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// expense-recognition-run routes.
func (r ExpenseRecognitionRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition_run.queue":                   r.QueueURL,
		"expense_recognition_run.queue_table":             r.QueueTableURL,
		"expense_recognition_run.list":                    r.ListURL,
		"expense_recognition_run.list_table":              r.ListTableURL,
		"expense_recognition_run.detail":                  r.DetailURL,
		"expense_recognition_run.detail_tab_action":       r.DetailTabActionURL,
		"expense_recognition_run.new":                     r.NewURL,
		"expense_recognition_run.generate":                r.GenerateURL,
		"expense_recognition_run.submit_batch":            r.SubmitBatchURL,
		"expense_recognition_run.per_supplier_drawer":     r.PerSupplierDrawerURL,
		"expense_recognition_run.per_subscription_drawer": r.PerSubscriptionDrawerURL,
	}
}
