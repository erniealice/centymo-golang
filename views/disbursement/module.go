package disbursement

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"

	disbursementaction "github.com/erniealice/centymo-golang/views/disbursement/action"
	disbursementdetail "github.com/erniealice/centymo-golang/views/disbursement/detail"
	disbursementlist "github.com/erniealice/centymo-golang/views/disbursement/list"
)

// ModuleDeps holds all dependencies for the disbursement module.
type ModuleDeps struct {
	Routes       centymo.DisbursementRoutes
	Labels       centymo.DisbursementLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed disbursement use case functions
	CreateDisbursement func(ctx context.Context, req *disbursementpb.CreateDisbursementRequest) (*disbursementpb.CreateDisbursementResponse, error)
	ReadDisbursement   func(ctx context.Context, req *disbursementpb.ReadDisbursementRequest) (*disbursementpb.ReadDisbursementResponse, error)
	UpdateDisbursement func(ctx context.Context, req *disbursementpb.UpdateDisbursementRequest) (*disbursementpb.UpdateDisbursementResponse, error)
	DeleteDisbursement func(ctx context.Context, req *disbursementpb.DeleteDisbursementRequest) (*disbursementpb.DeleteDisbursementResponse, error)
	ListDisbursements  func(ctx context.Context, req *disbursementpb.ListDisbursementsRequest) (*disbursementpb.ListDisbursementsResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// Module holds all constructed disbursement views.
type Module struct {
	routes           centymo.DisbursementRoutes
	Dashboard        view.View
	List             view.View
	Detail           view.View
	TabAction        view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

// NewModule creates the disbursement module with all views.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &disbursementlist.ListViewDeps{
		Routes:            deps.Routes,
		ListDisbursements: deps.ListDisbursements,
		RefreshURL:        deps.Routes.ListURL,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
	}

	detailDeps := &disbursementdetail.DetailViewDeps{
		Routes:           deps.Routes,
		ReadDisbursement: deps.ReadDisbursement,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewID

	actionDeps := &disbursementaction.Deps{
		Routes:             deps.Routes,
		Labels:             deps.Labels,
		CreateDisbursement: deps.CreateDisbursement,
		ReadDisbursement:   deps.ReadDisbursement,
		UpdateDisbursement: deps.UpdateDisbursement,
		DeleteDisbursement: deps.DeleteDisbursement,
	}

	return &Module{
		routes:        deps.Routes,
		Dashboard:     disbursementlist.NewView(listDeps),
		List:          disbursementlist.NewView(listDeps),
		Detail:        disbursementdetail.NewView(detailDeps),
		TabAction:     disbursementdetail.NewTabAction(detailDeps),
		Add:           disbursementaction.NewAddAction(actionDeps),
		Edit:          disbursementaction.NewEditAction(actionDeps),
		Delete:        disbursementaction.NewDeleteAction(actionDeps),
		BulkDelete:    disbursementaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        disbursementaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    disbursementaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: disbursementdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: disbursementdetail.NewAttachmentDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all disbursement routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)

	// Action routes (GET + POST for form-based)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)

	// Delete + status (POST only)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
