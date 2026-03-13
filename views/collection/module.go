package collection

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"

	collectionaction "github.com/erniealice/centymo-golang/views/collection/action"
	collectiondetail "github.com/erniealice/centymo-golang/views/collection/detail"
	collectionlist "github.com/erniealice/centymo-golang/views/collection/list"
)

// ModuleDeps holds all dependencies for the collection module.
type ModuleDeps struct {
	Routes       centymo.CollectionRoutes
	Labels       centymo.CollectionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed collection use case functions
	CreateCollection func(ctx context.Context, req *collectionpb.CreateCollectionRequest) (*collectionpb.CreateCollectionResponse, error)
	ReadCollection   func(ctx context.Context, req *collectionpb.ReadCollectionRequest) (*collectionpb.ReadCollectionResponse, error)
	UpdateCollection func(ctx context.Context, req *collectionpb.UpdateCollectionRequest) (*collectionpb.UpdateCollectionResponse, error)
	DeleteCollection func(ctx context.Context, req *collectionpb.DeleteCollectionRequest) (*collectionpb.DeleteCollectionResponse, error)
	ListCollections  func(ctx context.Context, req *collectionpb.ListCollectionsRequest) (*collectionpb.ListCollectionsResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, entityType, entityID string) ([]map[string]any, error)
	CreateAttachment func(ctx context.Context, data map[string]any) error
	DeleteAttachment func(ctx context.Context, id string) error
	NewID            func() string
}

// Module holds all constructed collection views.
type Module struct {
	routes           centymo.CollectionRoutes
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

// NewModule creates the collection module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &collectionaction.Deps{
		Routes:           deps.Routes,
		Labels:           deps.Labels,
		CreateCollection: deps.CreateCollection,
		ReadCollection:   deps.ReadCollection,
		UpdateCollection: deps.UpdateCollection,
		DeleteCollection: deps.DeleteCollection,
	}

	detailDeps := &collectiondetail.Deps{
		Routes:           deps.Routes,
		ReadCollection:   deps.ReadCollection,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		UploadFile:       deps.UploadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		DeleteAttachment: deps.DeleteAttachment,
		NewID:            deps.NewID,
	}

	listView := collectionlist.NewView(&collectionlist.Deps{
		Routes:          deps.Routes,
		ListCollections: deps.ListCollections,
		RefreshURL:      deps.Routes.ListURL,
		Labels:          deps.Labels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	})

	return &Module{
		routes:    deps.Routes,
		Dashboard: listView, // Dashboard reuses list view for now
		List:      listView,
		Detail:    collectiondetail.NewView(detailDeps),
		TabAction: collectiondetail.NewTabAction(detailDeps),
		Add:           collectionaction.NewAddAction(actionDeps),
		Edit:          collectionaction.NewEditAction(actionDeps),
		Delete:        collectionaction.NewDeleteAction(actionDeps),
		BulkDelete:    collectionaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        collectionaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    collectionaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: collectiondetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: collectiondetail.NewAttachmentDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all collection routes.
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
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
