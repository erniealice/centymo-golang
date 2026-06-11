package expenditure

// routes.go — expenditure entity route constants + Routes config struct.
// Extracted from the domain-level routes.go during the per-entity restructure.
// Entity-local naming (Expenditure prefix stripped). PurchaseOrder* consts/fields
// are surfaced here because purchase_order views consume expenditure.Routes
// (the combined purchase+expense Routes). Pure structural move — route strings
// are byte-identical.

const (

	// Expenditure (purchase + expense) routes
	PurchaseListURL      = "/purchases/list/{status}"
	PurchaseDashboardURL = "/purchases/dashboard"
	ExpenseListURL       = "/expenses/list/{status}"
	ExpenseDashboardURL  = "/expenses/dashboard"

	// Expenditure expense CRUD action routes
	ExpenseAddURL       = "/action/expense/add"
	ExpenseEditURL      = "/action/expense/edit/{id}"
	ExpenseDeleteURL    = "/action/expense/delete"
	ExpenseSetStatusURL = "/action/expense/set-status"
	ExpenseDetailURL    = "/expenses/detail/{id}"
	ExpenseTableURL     = "/action/expense/table/{status}"
	ExpenseTabActionURL = "/action/expense/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/expense/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/expense/detail/{id}/attachments/delete"

	// Expenditure expense line item action routes
	ExpenseLineItemAddURL    = "/action/expense/detail/{id}/items/add"
	ExpenseLineItemEditURL   = "/action/expense/detail/{id}/items/edit/{itemId}"
	ExpenseLineItemRemoveURL = "/action/expense/detail/{id}/items/remove"
	ExpenseLineItemTableURL  = "/action/expense/detail/{id}/items/table"

	// Expenditure pay action route (creates a pre-linked disbursement)
	ExpensePayURL = "/action/expense/detail/{id}/pay"

	// Expenditure report routes
	PurchasesSummaryURL = "/purchases/reports/purchases-summary"
	ExpensesSummaryURL  = "/expenses/reports/expenses-summary"

	// Expenditure settings (template management) routes
	SettingsTemplatesURL       = "/purchases/settings/templates"
	SettingsTemplateUploadURL  = "/action/purchase/settings/templates/upload"
	SettingsTemplateDeleteURL  = "/action/purchase/settings/templates/delete"
	SettingsTemplateDefaultURL = "/action/purchase/settings/templates/set-default/{id}"

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
	ExpenseCategoryListURL   = "/expenses/categories/list"
	ExpenseCategoryAddURL    = "/action/expense/categories/add"
	ExpenseCategoryEditURL   = "/action/expense/categories/edit/{id}"
	ExpenseCategoryDeleteURL = "/action/expense/categories/delete"
	ExpenseCategoryTableURL  = "/action/expense/categories/table"
)

// Routes holds all route paths for expenditure views (purchase + expense).
type Routes struct {
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

// DefaultRoutes returns an Routes populated from the
// package-level route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		PurchaseListURL:      PurchaseListURL,
		PurchaseDashboardURL: PurchaseDashboardURL,
		ExpenseListURL:       ExpenseListURL,
		ExpenseDashboardURL:  ExpenseDashboardURL,

		PurchasesSummaryURL: PurchasesSummaryURL,
		ExpensesSummaryURL:  ExpensesSummaryURL,

		SettingsTemplatesURL:       SettingsTemplatesURL,
		SettingsTemplateUploadURL:  SettingsTemplateUploadURL,
		SettingsTemplateDeleteURL:  SettingsTemplateDeleteURL,
		SettingsTemplateDefaultURL: SettingsTemplateDefaultURL,

		AddURL:              ExpenseAddURL,
		EditURL:             ExpenseEditURL,
		DeleteURL:           ExpenseDeleteURL,
		SetStatusURL:        ExpenseSetStatusURL,
		DetailURL:           ExpenseDetailURL,
		TableURL:            ExpenseTableURL,
		TabActionURL:        ExpenseTabActionURL,
		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		LineItemAddURL:    ExpenseLineItemAddURL,
		LineItemEditURL:   ExpenseLineItemEditURL,
		LineItemRemoveURL: ExpenseLineItemRemoveURL,
		LineItemTableURL:  ExpenseLineItemTableURL,

		PayURL: ExpensePayURL,

		ExpenseCategoryListURL:   ExpenseCategoryListURL,
		ExpenseCategoryAddURL:    ExpenseCategoryAddURL,
		ExpenseCategoryEditURL:   ExpenseCategoryEditURL,
		ExpenseCategoryDeleteURL: ExpenseCategoryDeleteURL,
		ExpenseCategoryTableURL:  ExpenseCategoryTableURL,

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
func (r Routes) RouteMap() map[string]string {
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
