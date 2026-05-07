package procurement_request

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	procurementrequestaction "github.com/erniealice/centymo-golang/views/procurement_request/action"
	procurementrequestdetail "github.com/erniealice/centymo-golang/views/procurement_request/detail"
	procurementrequestlist "github.com/erniealice/centymo-golang/views/procurement_request/list"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the procurement_request module.
type ModuleDeps struct {
	Routes       centymo.ProcurementRequestRoutes
	Labels       centymo.ProcurementRequestLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Core CRUD
	ListProcurementRequests  func(ctx context.Context, req *procurementrequestpb.ListProcurementRequestsRequest) (*procurementrequestpb.ListProcurementRequestsResponse, error)
	ReadProcurementRequest   func(ctx context.Context, req *procurementrequestpb.ReadProcurementRequestRequest) (*procurementrequestpb.ReadProcurementRequestResponse, error)
	CreateProcurementRequest func(ctx context.Context, req *procurementrequestpb.CreateProcurementRequestRequest) (*procurementrequestpb.CreateProcurementRequestResponse, error)
	UpdateProcurementRequest func(ctx context.Context, req *procurementrequestpb.UpdateProcurementRequestRequest) (*procurementrequestpb.UpdateProcurementRequestResponse, error)
	DeleteProcurementRequest func(ctx context.Context, req *procurementrequestpb.DeleteProcurementRequestRequest) (*procurementrequestpb.DeleteProcurementRequestResponse, error)

	// Workflow actions
	SetProcurementRequestStatus func(ctx context.Context, id, status string) error

	// Child entity — lines
	ListProcurementRequestLines func(ctx context.Context, req *procurementrequestlinepb.ListProcurementRequestLinesRequest) (*procurementrequestlinepb.ListProcurementRequestLinesResponse, error)

	// Related entities
	ListSuppliers      func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)

	// Workflow invocations (block.go injects use-case-backed closures)
	SubmitProcurementRequest  func(ctx context.Context, id string) error
	ApproveProcurementRequest func(ctx context.Context, id string) error
	RejectProcurementRequest  func(ctx context.Context, id string, reason string) error
	SpawnPurchaseOrder        func(ctx context.Context, id string) (newPOID string, err error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string
}

// Module holds all constructed procurement_request views.
type Module struct {
	routes        centymo.ProcurementRequestRoutes
	List          view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	SetStatus     view.View
	BulkSetStatus view.View
	Submit           view.View
	Approve          view.View
	Reject           view.View
	SpawnPO          view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

// NewModule creates the procurement_request module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &procurementrequestaction.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		CreateProcurementRequest:    deps.CreateProcurementRequest,
		ReadProcurementRequest:      deps.ReadProcurementRequest,
		UpdateProcurementRequest:    deps.UpdateProcurementRequest,
		DeleteProcurementRequest:    deps.DeleteProcurementRequest,
		SetProcurementRequestStatus: deps.SetProcurementRequestStatus,
		ListSuppliers:               deps.ListSuppliers,
	}

	detailDeps := &procurementrequestdetail.DetailViewDeps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ReadProcurementRequest:      deps.ReadProcurementRequest,
		ListProcurementRequestLines: deps.ListProcurementRequestLines,
		ListPurchaseOrders:          deps.ListPurchaseOrders,
		SubmitProcurementRequest:    deps.SubmitProcurementRequest,
		ApproveProcurementRequest:   deps.ApproveProcurementRequest,
		RejectProcurementRequest:    deps.RejectProcurementRequest,
		SpawnPurchaseOrder:          deps.SpawnPurchaseOrder,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewAttachmentID

	listDeps := &procurementrequestlist.ListViewDeps{
		Routes:                  deps.Routes,
		ListProcurementRequests: deps.ListProcurementRequests,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		TableLabels:             deps.TableLabels,
	}

	m := &Module{
		routes:        deps.Routes,
		List:          procurementrequestlist.NewView(listDeps),
		Add:           procurementrequestaction.NewAddAction(actionDeps),
		Edit:          procurementrequestaction.NewEditAction(actionDeps),
		Delete:        procurementrequestaction.NewDeleteAction(actionDeps),
		SetStatus:     procurementrequestaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: procurementrequestaction.NewBulkSetStatusAction(actionDeps),
	}

	if deps.ReadProcurementRequest != nil {
		m.Detail = procurementrequestdetail.NewView(detailDeps)
		m.TabAction = procurementrequestdetail.NewTabAction(detailDeps)
		m.Submit = procurementrequestdetail.NewSubmitAction(detailDeps)
		m.Approve = procurementrequestdetail.NewApproveAction(detailDeps)
		m.Reject = procurementrequestdetail.NewRejectAction(detailDeps)
		m.SpawnPO = procurementrequestdetail.NewSpawnPOAction(detailDeps)
	}
	if deps.UploadFile != nil {
		m.AttachmentUpload = procurementrequestdetail.NewAttachmentUploadAction(detailDeps)
		m.AttachmentDelete = procurementrequestdetail.NewAttachmentDeleteAction(detailDeps)
	}

	return m
}

// RegisterRoutes registers all procurement_request routes.
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
	if m.Submit != nil {
		r.POST(m.routes.SubmitURL, m.Submit)
	}
	if m.Approve != nil {
		r.POST(m.routes.ApproveURL, m.Approve)
	}
	if m.Reject != nil {
		r.POST(m.routes.RejectURL, m.Reject)
	}
	if m.SpawnPO != nil {
		r.POST(m.routes.SpawnPOURL, m.SpawnPO)
	}
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
