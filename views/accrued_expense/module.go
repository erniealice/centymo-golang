// Package accrued_expense is the centymo views package for AccruedExpense.
//
// MAIN-THREAD WIRING NOTE (block.go):
//   Block-level wiring is intentionally DEFERRED to the main thread. The
//   integrator must:
//     1. Register routes/labels in apps/service-admin sidebar.
//     2. Add a centymo block.go entry that constructs ModuleDeps from the
//        domain providers (espyna AccruedExpense + AccruedExpenseSettlement
//        use case groups) and calls NewModule + RegisterRoutes.
//     3. Wire the per-tier translation file via translations.go using the
//        lyngua key root "accruedExpense".
//   See plan §7 (Phase P10) for the full integrator checklist.
package accrued_expense

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	accruedexpenseaction "github.com/erniealice/centymo-golang/views/accrued_expense/action"
	accruedexpensedetail "github.com/erniealice/centymo-golang/views/accrued_expense/detail"
	accruedexpenselist "github.com/erniealice/centymo-golang/views/accrued_expense/list"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the accrued_expense module.
type ModuleDeps struct {
	Routes       centymo.AccruedExpenseRoutes
	Labels       centymo.AccruedExpenseLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Core CRUD
	ListAccruedExpenses     func(ctx context.Context, req *accruedexpensepb.ListAccruedExpensesRequest) (*accruedexpensepb.ListAccruedExpensesResponse, error)
	ReadAccruedExpense      func(ctx context.Context, req *accruedexpensepb.ReadAccruedExpenseRequest) (*accruedexpensepb.ReadAccruedExpenseResponse, error)
	CreateAccruedExpense    func(ctx context.Context, req *accruedexpensepb.CreateAccruedExpenseRequest) (*accruedexpensepb.CreateAccruedExpenseResponse, error)
	UpdateAccruedExpense    func(ctx context.Context, req *accruedexpensepb.UpdateAccruedExpenseRequest) (*accruedexpensepb.UpdateAccruedExpenseResponse, error)
	DeleteAccruedExpense    func(ctx context.Context, req *accruedexpensepb.DeleteAccruedExpenseRequest) (*accruedexpensepb.DeleteAccruedExpenseResponse, error)
	SetAccruedExpenseStatus func(ctx context.Context, id, status string) error

	// Inline child — settlements
	ListAccruedExpenseSettlements func(ctx context.Context, req *accruedexpensepb.ListAccruedExpenseSettlementsRequest) (*accruedexpensepb.ListAccruedExpenseSettlementsResponse, error)

	// Workflow — espyna closures
	SettleAccrual  func(ctx context.Context, req *accruedexpensepb.SettleAccrualRequest) error
	ReverseAccrual func(ctx context.Context, id, reason string) error

	// AccrueFromContract — bulk accrual workflow (espyna). Optional; when nil
	// the action returns 422 indicating it isn't wired (still 200/422, never 405).
	AccrueFromContract accruedexpenseaction.AccrueFromContractFunc

	// Dropdowns
	ListSuppliers         func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string
}

// Module holds all constructed accrued_expense views.
type Module struct {
	routes             centymo.AccruedExpenseRoutes
	List               view.View
	Detail             view.View
	TabAction          view.View
	Add                view.View
	Edit               view.View
	Delete             view.View
	SetStatus          view.View
	BulkSetStatus      view.View
	Settle             view.View
	Reverse            view.View
	AccrueFromContract view.View
	AttachmentUpload   view.View
	AttachmentDelete   view.View
}

// NewModule creates the accrued_expense module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &accruedexpenselist.ListViewDeps{
		Routes:              deps.Routes,
		ListAccruedExpenses: deps.ListAccruedExpenses,
		Labels:              deps.Labels,
		CommonLabels:        deps.CommonLabels,
		TableLabels:         deps.TableLabels,
	}

	detailDeps := &accruedexpensedetail.DetailViewDeps{
		Routes:                        deps.Routes,
		Labels:                        deps.Labels,
		CommonLabels:                  deps.CommonLabels,
		TableLabels:                   deps.TableLabels,
		ReadAccruedExpense:            deps.ReadAccruedExpense,
		ListAccruedExpenseSettlements: deps.ListAccruedExpenseSettlements,
		SettleAccrual:                 deps.SettleAccrual,
		ReverseAccrual:                deps.ReverseAccrual,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewAttachmentID

	actionDeps := &accruedexpenseaction.Deps{
		Routes:                  deps.Routes,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		CreateAccruedExpense:    deps.CreateAccruedExpense,
		ReadAccruedExpense:      deps.ReadAccruedExpense,
		UpdateAccruedExpense:    deps.UpdateAccruedExpense,
		DeleteAccruedExpense:    deps.DeleteAccruedExpense,
		SetAccruedExpenseStatus: deps.SetAccruedExpenseStatus,
		ListSuppliers:           deps.ListSuppliers,
		ListSupplierContracts:   deps.ListSupplierContracts,
	}

	m := &Module{
		routes:        deps.Routes,
		List:          accruedexpenselist.NewView(listDeps),
		Add:           accruedexpenseaction.NewAddAction(actionDeps),
		Edit:          accruedexpenseaction.NewEditAction(actionDeps),
		Delete:        accruedexpenseaction.NewDeleteAction(actionDeps),
		SetStatus:     accruedexpenseaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: accruedexpenseaction.NewBulkSetStatusAction(actionDeps),
	}
	if deps.ReadAccruedExpense != nil {
		m.Detail = accruedexpensedetail.NewView(detailDeps)
		m.TabAction = accruedexpensedetail.NewTabAction(detailDeps)
		m.Settle = accruedexpensedetail.NewSettleAction(detailDeps)
		m.Reverse = accruedexpensedetail.NewReverseAction(detailDeps)
	}
	if deps.UploadFile != nil {
		m.AttachmentUpload = accruedexpensedetail.NewAttachmentUploadAction(detailDeps)
		m.AttachmentDelete = accruedexpensedetail.NewAttachmentDeleteAction(detailDeps)
	}
	// AccrueFromContract: register unconditionally so the route returns 422 (not
	// 405) even before the espyna closure is wired.
	m.AccrueFromContract = accruedexpenseaction.NewAccrueFromContractAction(deps.AccrueFromContract)
	return m
}

// RegisterRoutes registers all accrued_expense routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)

	if m.Detail != nil {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.Settle != nil {
		r.POST(m.routes.SettleURL, m.Settle)
	}
	if m.Reverse != nil {
		r.POST(m.routes.ReverseURL, m.Reverse)
	}
	if m.AccrueFromContract != nil {
		r.POST(m.routes.AccrueFromContractURL, m.AccrueFromContract)
	}
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
