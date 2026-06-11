package collectionmodule

import (
	"context"

	shared "github.com/erniealice/centymo-golang/domain/treasury/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"

	epkg "github.com/erniealice/centymo-golang/domain/treasury/collection"
	collectionaction "github.com/erniealice/centymo-golang/domain/treasury/collection/action"
	collectiondashboard "github.com/erniealice/centymo-golang/domain/treasury/collection/dashboard"
	collectiondetail "github.com/erniealice/centymo-golang/domain/treasury/collection/detail"
	collectionlist "github.com/erniealice/centymo-golang/domain/treasury/collection/list"
)

// ModuleDeps holds all dependencies for the collection module.
type ModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
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
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Cash dashboard data callback (Phase 5 — nil-safe; degrades to zero values).
	// Orchestrator wraps the espyna treasury/collection/dashboard use case
	// here, projecting workspace_id from the request context.
	GetCashDashboardPageData func(ctx context.Context, req *collectiondashboard.Request) (*collectiondashboard.Response, error)

	// GetFunctionalCurrency returns the workspace ISO 4217 currency code for
	// money display in the dashboard. Nil-safe — when absent, money strings
	// omit the currency prefix.
	GetFunctionalCurrency func(ctx context.Context) string

	// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow
	// closures + label tables. Nil-safe: when SettleUnscheduled / RefundUnscheduled
	// / Cancel are unset, the drawer GET still renders so the operator can see
	// the field shape, but the POST returns a permission-denied error.
	AdvanceLabels            shared.TreasuryAdvanceLabels
	AdvanceEnumLabels        shared.AdvanceEnumLabels
	SettleUnscheduledAdvance func(ctx context.Context, in shared.AdvanceSettleViewInput) (*shared.AdvanceSettleViewOutput, error)
	RefundUnscheduledAdvance func(ctx context.Context, in shared.AdvanceRefundViewInput) (*shared.AdvanceRefundViewOutput, error)
	CancelAdvance            func(ctx context.Context, in shared.AdvanceCancelViewInput) (*shared.AdvanceCancelViewOutput, error)
}

// Module holds all constructed collection views.
type Module struct {
	routes           epkg.Routes
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
	// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow drawers.
	AdvanceSettle view.View
	AdvanceRefund view.View
	AdvanceCancel view.View
}

// NewModule creates the collection module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &collectionaction.Deps{
		Routes:            deps.Routes,
		Labels:            deps.Labels,
		CreateCollection:  deps.CreateCollection,
		ReadCollection:    deps.ReadCollection,
		UpdateCollection:  deps.UpdateCollection,
		DeleteCollection:  deps.DeleteCollection,
		AdvanceEnumLabels: deps.AdvanceEnumLabels,
	}

	detailDeps := &collectiondetail.DetailViewDeps{
		Routes:            deps.Routes,
		ReadCollection:    deps.ReadCollection,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
		AdvanceLabels:     deps.AdvanceLabels,
		AdvanceEnumLabels: deps.AdvanceEnumLabels,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewID

	listView := collectionlist.NewView(&collectionlist.ListViewDeps{
		Routes:          deps.Routes,
		ListCollections: deps.ListCollections,
		RefreshURL:      deps.Routes.ListURL,
		Labels:          deps.Labels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
	})

	dashboardView := collectiondashboard.NewView(&collectiondashboard.Deps{
		Routes:                deps.Routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		GetPageData:           deps.GetCashDashboardPageData,
		GetFunctionalCurrency: deps.GetFunctionalCurrency,
	})

	// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow.
	advanceActionDeps := &collectionaction.AdvanceActionDeps{
		Routes:            deps.Routes,
		Labels:            deps.Labels,
		AdvanceLabels:     deps.AdvanceLabels,
		EnumLabels:        deps.AdvanceEnumLabels,
		CommonLabels:      deps.CommonLabels,
		SettleUnscheduled: deps.SettleUnscheduledAdvance,
		RefundUnscheduled: deps.RefundUnscheduledAdvance,
		Cancel:            deps.CancelAdvance,
	}

	return &Module{
		routes:           deps.Routes,
		Dashboard:        dashboardView,
		List:             listView,
		Detail:           collectiondetail.NewView(detailDeps),
		TabAction:        collectiondetail.NewTabAction(detailDeps),
		Add:              collectionaction.NewAddAction(actionDeps),
		Edit:             collectionaction.NewEditAction(actionDeps),
		Delete:           collectionaction.NewDeleteAction(actionDeps),
		BulkDelete:       collectionaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        collectionaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    collectionaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: collectiondetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: collectiondetail.NewAttachmentDeleteAction(detailDeps),
		AdvanceSettle:    collectionaction.NewSettleAction(advanceActionDeps),
		AdvanceRefund:    collectionaction.NewRefundAction(advanceActionDeps),
		AdvanceCancel:    collectionaction.NewCancelAction(advanceActionDeps),
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
	// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow drawers.
	if m.routes.SettleURL != "" {
		r.GET(m.routes.SettleURL, m.AdvanceSettle)
		r.POST(m.routes.SettleURL, m.AdvanceSettle)
	}
	if m.routes.RefundURL != "" {
		r.GET(m.routes.RefundURL, m.AdvanceRefund)
		r.POST(m.routes.RefundURL, m.AdvanceRefund)
	}
	if m.routes.CancelURL != "" {
		r.GET(m.routes.CancelURL, m.AdvanceCancel)
		r.POST(m.routes.CancelURL, m.AdvanceCancel)
	}
}
