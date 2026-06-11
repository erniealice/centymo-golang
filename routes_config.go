package centymo

import (
	"github.com/erniealice/centymo-golang/domain/subscription"
	"github.com/erniealice/centymo-golang/domain/treasury"
)

// ── centymo W4 subscription-domain compatibility shim ────────────────────────
// The Subscription/PriceSchedule route types + their Default* constructors moved
// to domain/subscription/ (centymo W4). entydad-golang/block/route_loading.go is
// an EXTERNAL consumer (outside this wave's edit scope) that still references
// centymo.SubscriptionRoutes / centymo.PriceScheduleRoutes and their Default*
// constructors. These thin aliases + forwarders keep that consumer compiling
// with ZERO behaviour change (pure type-identity aliases). Remove once entydad
// is re-pointed to domain/subscription directly (W9 / entydad-coordinated).
type SubscriptionRoutes = subscription.SubscriptionRoutes
type PriceScheduleRoutes = subscription.PriceScheduleRoutes

func DefaultSubscriptionRoutes() SubscriptionRoutes   { return subscription.DefaultSubscriptionRoutes() }
func DefaultPriceScheduleRoutes() PriceScheduleRoutes { return subscription.DefaultPriceScheduleRoutes() }

// ── centymo W5 treasury-domain compatibility shim ────────────────────────────
// Treasury types (Collection/Disbursement labels+routes, TreasuryAdvancesRoutes,
// the AdvanceRecognizeMilestone view I/O) moved to domain/treasury/ (centymo W5).
// The not-yet-migrated W6 view packages still reference a subset of them via the
// centymo root:
//   - views/expenditure/*            -> DisbursementRoutes / DisbursementLabels /
//                                       DisbursementFormLabels (the expense "pay"
//                                       flow creates a pre-linked disbursement)
//   - views/supplier_billing_event/* -> TreasuryAdvancesRoutes (+ its Default*)
//                                       and AdvanceRecognizeMilestoneInput/Output
//   - domain/subscription/views/...  -> AdvanceRecognizeMilestoneInput/Output
//                                       (already-migrated W4 billing-event action)
// These thin aliases + forwarders keep those consumers compiling with ZERO
// behaviour change. Removed as each consuming domain migrates (W6 / W9).
type DisbursementRoutes = treasury.DisbursementRoutes
type DisbursementLabels = treasury.DisbursementLabels
type DisbursementFormLabels = treasury.DisbursementFormLabels
type TreasuryAdvancesRoutes = treasury.TreasuryAdvancesRoutes
type AdvanceRecognizeMilestoneInput = treasury.AdvanceRecognizeMilestoneInput
type AdvanceRecognizeMilestoneOutput = treasury.AdvanceRecognizeMilestoneOutput

func DefaultTreasuryAdvancesRoutes() TreasuryAdvancesRoutes {
	return treasury.DefaultTreasuryAdvancesRoutes()
}

// Three-level routing system for centymo views:
//
// Level 1: Generic defaults from Go consts (this file).
//   DefaultXxxRoutes() constructors return structs populated from the route
//   constants defined in routes.go. These are sensible defaults that work
//   out of the box for any app.
//
// Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//   Consumer apps can load a JSON config that partially overrides the
//   default routes. Struct fields carry json tags for unmarshalling.
//
// Level 3: App-specific overrides via Go field assignment (optional).
//   After loading defaults and/or JSON, consumer apps can programmatically
//   set individual fields to further customize routing.
//
// Each route struct also exposes a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "product.list"),
// useful for template rendering, URL resolution, and debugging.

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
// P3b — Procurement Operations app routes
// (composition surface; no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementRoutes holds the URL constants for the Procurement Operations app.
// These are defined in the centymo package so service-admin composition (P3c)
// can wire them into SidebarRoutes.Operations.Procurement.
type ProcurementRoutes struct {
	// Dashboard
	DashboardURL string `json:"dashboard_url"`

	// Contract operations (views over SupplierContract)
	RenewalCalendarURL string `json:"renewal_calendar_url"`
	VarianceURL        string `json:"variance_url"`
	UtilizationURL     string `json:"utilization_url"`

	// Recurrence drafts queue (lights up when P5 ships the recurrence engine)
	RecurrenceDraftsURL string `json:"recurrence_drafts_url"`
}

// DefaultProcurementRoutes returns a ProcurementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultProcurementRoutes() ProcurementRoutes {
	return ProcurementRoutes{
		DashboardURL:        ProcurementDashboardURL,
		RenewalCalendarURL:  ProcurementRenewalCalendarURL,
		VarianceURL:         ProcurementVarianceURL,
		UtilizationURL:      ProcurementUtilizationURL,
		RecurrenceDraftsURL: ProcurementRecurrenceDraftsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// procurement operations app routes.
func (r ProcurementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"procurement.dashboard":         r.DashboardURL,
		"procurement.renewals":          r.RenewalCalendarURL,
		"procurement.variance":          r.VarianceURL,
		"procurement.utilization":       r.UtilizationURL,
		"procurement.recurrence_drafts": r.RecurrenceDraftsURL,
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
// P3 — SupplierSubscription routes (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// SupplierSubscriptionRoutes holds all route paths for supplier_subscription views.
type SupplierSubscriptionRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Search autocomplete endpoints for the add/edit drawer
	SearchCostPlanURL string `json:"search_cost_plan_url"`
	SearchSupplierURL string `json:"search_supplier_url"`

	// Recognition CTA — POST; opens the recognize-expense drawer on the detail page.
	RecognizeExpenseURL string `json:"recognize_expense_url"`

	// ExpenseRecognitionRunURL — GET; opens the per-SupplierSubscription Expense
	// Recognition Run drawer (Surface C). Resolved by resolveRecognitionsPrimaryAction
	// for CostPlan.billing_kind RECURRING / CONTRACT-with-cycle.
	// Plan A 20260517-expense-run Phase 4 / Surface C.
	ExpenseRecognitionRunURL string `json:"expense_recognition_run_url"`
}

// DefaultSupplierSubscriptionRoutes returns a SupplierSubscriptionRoutes using route constants.
func DefaultSupplierSubscriptionRoutes() SupplierSubscriptionRoutes {
	return SupplierSubscriptionRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-subscriptions",

		ListURL:                  SupplierSubscriptionListURL,
		TableURL:                 SupplierSubscriptionTableURL,
		DetailURL:                SupplierSubscriptionDetailURL,
		AddURL:                   SupplierSubscriptionAddURL,
		EditURL:                  SupplierSubscriptionEditURL,
		DeleteURL:                SupplierSubscriptionDeleteURL,
		BulkDeleteURL:            SupplierSubscriptionBulkDeleteURL,
		SetStatusURL:             SupplierSubscriptionSetStatusURL,
		BulkSetStatusURL:         SupplierSubscriptionBulkSetStatusURL,
		TabActionURL:             SupplierSubscriptionTabActionURL,
		SearchCostPlanURL:        SupplierSubscriptionSearchCostPlanURL,
		SearchSupplierURL:        SupplierSubscriptionSearchSupplierURL,
		RecognizeExpenseURL:      SupplierSubscriptionRecognizeExpenseURL,
		ExpenseRecognitionRunURL: ExpenseRecognitionRunPerSubscriptionDrawerURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierSubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_subscription.list":                    r.ListURL,
		"supplier_subscription.table":                   r.TableURL,
		"supplier_subscription.detail":                  r.DetailURL,
		"supplier_subscription.add":                     r.AddURL,
		"supplier_subscription.edit":                    r.EditURL,
		"supplier_subscription.delete":                  r.DeleteURL,
		"supplier_subscription.bulk_delete":             r.BulkDeleteURL,
		"supplier_subscription.set_status":              r.SetStatusURL,
		"supplier_subscription.bulk_set_status":         r.BulkSetStatusURL,
		"supplier_subscription.tab_action":              r.TabActionURL,
		"supplier_subscription.search_cost_plan":        r.SearchCostPlanURL,
		"supplier_subscription.search_supplier":         r.SearchSupplierURL,
		"supplier_subscription.recognize_expense":       r.RecognizeExpenseURL,
		"supplier_subscription.expense_recognition_run": r.ExpenseRecognitionRunURL,
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

// ---------------------------------------------------------------------------
// P3 — CostSchedule routes
// ---------------------------------------------------------------------------

// CostScheduleRoutes holds all route paths for cost_schedule views.
type CostScheduleRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`
}

// DefaultCostScheduleRoutes returns a CostScheduleRoutes using route constants.
func DefaultCostScheduleRoutes() CostScheduleRoutes {
	return CostScheduleRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-schedules",

		ListURL:          CostScheduleListURL,
		TableURL:         CostScheduleTableURL,
		DetailURL:        CostScheduleDetailURL,
		AddURL:           CostScheduleAddURL,
		EditURL:          CostScheduleEditURL,
		DeleteURL:        CostScheduleDeleteURL,
		BulkDeleteURL:    CostScheduleBulkDeleteURL,
		SetStatusURL:     CostScheduleSetStatusURL,
		BulkSetStatusURL: CostScheduleBulkSetStatusURL,
		TabActionURL:     CostScheduleTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"cost_schedule.list":            r.ListURL,
		"cost_schedule.table":           r.TableURL,
		"cost_schedule.detail":          r.DetailURL,
		"cost_schedule.add":             r.AddURL,
		"cost_schedule.edit":            r.EditURL,
		"cost_schedule.delete":          r.DeleteURL,
		"cost_schedule.bulk_delete":     r.BulkDeleteURL,
		"cost_schedule.set_status":      r.SetStatusURL,
		"cost_schedule.bulk_set_status": r.BulkSetStatusURL,
		"cost_schedule.tab_action":      r.TabActionURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierPlan routes
// ---------------------------------------------------------------------------

// SupplierPlanRoutes holds all route paths for supplier_plan views.
type SupplierPlanRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Autocomplete search URL for the supplier select in add/edit forms.
	SearchSupplierURL string `json:"search_supplier_url"`
}

// DefaultSupplierPlanRoutes returns a SupplierPlanRoutes using route constants.
func DefaultSupplierPlanRoutes() SupplierPlanRoutes {
	return SupplierPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-plans",

		ListURL:          SupplierPlanListURL,
		TableURL:         SupplierPlanTableURL,
		DetailURL:        SupplierPlanDetailURL,
		AddURL:           SupplierPlanAddURL,
		EditURL:          SupplierPlanEditURL,
		DeleteURL:        SupplierPlanDeleteURL,
		BulkDeleteURL:    SupplierPlanBulkDeleteURL,
		SetStatusURL:     SupplierPlanSetStatusURL,
		BulkSetStatusURL: SupplierPlanBulkSetStatusURL,
		TabActionURL:     SupplierPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_plan.list":            r.ListURL,
		"supplier_plan.table":           r.TableURL,
		"supplier_plan.detail":          r.DetailURL,
		"supplier_plan.add":             r.AddURL,
		"supplier_plan.edit":            r.EditURL,
		"supplier_plan.delete":          r.DeleteURL,
		"supplier_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_plan.set_status":      r.SetStatusURL,
		"supplier_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_plan.tab_action":      r.TabActionURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — CostPlan routes
// ---------------------------------------------------------------------------

// CostPlanRoutes holds all route paths for cost_plan views.
type CostPlanRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// SupplierProductCostPlan inline CRUD within cost_plan detail
	ProductCostAddURL    string `json:"product_cost_add_url"`
	ProductCostEditURL   string `json:"product_cost_edit_url"`
	ProductCostDeleteURL string `json:"product_cost_delete_url"`

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL        string `json:"search_supplier_plan_url"`
	SearchCostScheduleURL        string `json:"search_cost_schedule_url"`
	SearchSupplierProductPlanURL string `json:"search_supplier_product_plan_url"`
}

// DefaultCostPlanRoutes returns a CostPlanRoutes using route constants.
func DefaultCostPlanRoutes() CostPlanRoutes {
	return CostPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-plans",

		ListURL:              CostPlanListURL,
		TableURL:             CostPlanTableURL,
		DetailURL:            CostPlanDetailURL,
		AddURL:               CostPlanAddURL,
		EditURL:              CostPlanEditURL,
		DeleteURL:            CostPlanDeleteURL,
		BulkDeleteURL:        CostPlanBulkDeleteURL,
		SetStatusURL:         CostPlanSetStatusURL,
		BulkSetStatusURL:     CostPlanBulkSetStatusURL,
		TabActionURL:         CostPlanTabActionURL,
		ProductCostAddURL:    CostPlanProductCostAddURL,
		ProductCostEditURL:   CostPlanProductCostEditURL,
		ProductCostDeleteURL: CostPlanProductCostDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"cost_plan.list":                r.ListURL,
		"cost_plan.table":               r.TableURL,
		"cost_plan.detail":              r.DetailURL,
		"cost_plan.add":                 r.AddURL,
		"cost_plan.edit":                r.EditURL,
		"cost_plan.delete":              r.DeleteURL,
		"cost_plan.bulk_delete":         r.BulkDeleteURL,
		"cost_plan.set_status":          r.SetStatusURL,
		"cost_plan.bulk_set_status":     r.BulkSetStatusURL,
		"cost_plan.tab_action":          r.TabActionURL,
		"cost_plan.product_cost.add":    r.ProductCostAddURL,
		"cost_plan.product_cost.edit":   r.ProductCostEditURL,
		"cost_plan.product_cost.delete": r.ProductCostDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan routes
// ---------------------------------------------------------------------------

// SupplierProductPlanRoutes holds all route paths for supplier_product_plan views.
type SupplierProductPlanRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL string `json:"search_supplier_plan_url"`
	SearchProductURL      string `json:"search_product_url"`
}

// DefaultSupplierProductPlanRoutes returns a SupplierProductPlanRoutes using route constants.
func DefaultSupplierProductPlanRoutes() SupplierProductPlanRoutes {
	return SupplierProductPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-product-plans",

		ListURL:          SupplierProductPlanListURL,
		TableURL:         SupplierProductPlanTableURL,
		DetailURL:        SupplierProductPlanDetailURL,
		AddURL:           SupplierProductPlanAddURL,
		EditURL:          SupplierProductPlanEditURL,
		DeleteURL:        SupplierProductPlanDeleteURL,
		BulkDeleteURL:    SupplierProductPlanBulkDeleteURL,
		SetStatusURL:     SupplierProductPlanSetStatusURL,
		BulkSetStatusURL: SupplierProductPlanBulkSetStatusURL,
		TabActionURL:     SupplierProductPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierProductPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_product_plan.list":            r.ListURL,
		"supplier_product_plan.table":           r.TableURL,
		"supplier_product_plan.detail":          r.DetailURL,
		"supplier_product_plan.add":             r.AddURL,
		"supplier_product_plan.edit":            r.EditURL,
		"supplier_product_plan.delete":          r.DeleteURL,
		"supplier_product_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_product_plan.set_status":      r.SetStatusURL,
		"supplier_product_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_product_plan.tab_action":      r.TabActionURL,
	}
}

// MapTableLabels is a shared helper used across all centymo view modules to
// produce a types.TableLabels from pyeza CommonLabels. Defined here to avoid
// duplication; all block module wirings call this.
func mapTableLabelsFromStrings(search, searchPlaceholder, sortAsc, sortDesc, noResults, loading string) struct{} {
	// Placeholder — actual implementation lives in the block package; this
	// comment documents the cross-module convention.
	return struct{}{}
}
