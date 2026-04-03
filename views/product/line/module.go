package line

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	lineaction "github.com/erniealice/centymo-golang/views/product/line/action"
	linedetail "github.com/erniealice/centymo-golang/views/product/line/detail"
	linelist "github.com/erniealice/centymo-golang/views/product/line/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
)

// ModuleDeps holds all dependencies for the line module.
type ModuleDeps struct {
	Routes       centymo.ProductLineRoutes
	Labels       centymo.ProductLineLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	CreateLine func(ctx context.Context, req *linepb.CreateLineRequest) (*linepb.CreateLineResponse, error)
	ReadLine   func(ctx context.Context, req *linepb.ReadLineRequest) (*linepb.ReadLineResponse, error)
	UpdateLine func(ctx context.Context, req *linepb.UpdateLineRequest) (*linepb.UpdateLineResponse, error)
	DeleteLine func(ctx context.Context, req *linepb.DeleteLineRequest) (*linepb.DeleteLineResponse, error)
	ListLines  func(ctx context.Context, req *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)

	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// Module holds all constructed line views.
type Module struct {
	routes           centymo.ProductLineRoutes
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

// NewModule creates the line module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &lineaction.Deps{
		Routes:     deps.Routes,
		Labels:     deps.Labels,
		CreateLine: deps.CreateLine,
		ReadLine:   deps.ReadLine,
		UpdateLine: deps.UpdateLine,
		DeleteLine: deps.DeleteLine,
	}

	detailDeps := &linedetail.DetailViewDeps{
		Routes:       deps.Routes,
		ReadLine:     deps.ReadLine,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewID

	listView := linelist.NewView(&linelist.ListViewDeps{
		Routes:       deps.Routes,
		ListLines:    deps.ListLines,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	})

	return &Module{
		routes:           deps.Routes,
		Dashboard:        listView,
		List:             listView,
		Detail:           linedetail.NewView(detailDeps),
		TabAction:        linedetail.NewTabAction(detailDeps),
		Add:              lineaction.NewAddAction(actionDeps),
		Edit:             lineaction.NewEditAction(actionDeps),
		Delete:           lineaction.NewDeleteAction(actionDeps),
		BulkDelete:       lineaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        lineaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    lineaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: linedetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: linedetail.NewAttachmentDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all line routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
