package collection_method

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	cmaction "github.com/erniealice/centymo-golang/views/collection_method/action"
	cmdetail "github.com/erniealice/centymo-golang/views/collection_method/detail"
	cmlist "github.com/erniealice/centymo-golang/views/collection_method/list"

	cmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the collection_method module
// (treasury-domain-rebuild Stage 1, pages.md §B-5).
//
// The CRUD closures are all OPTIONAL — the espyna collection_method use cases
// are NOT yet implemented (W1/W2 added proto + migration only). When nil the
// module still mounts: list renders empty, detail redirects, and the drawer-
// form + HTMX fragment swap render correctly (the Stage-1 "compile + structural
// correctness" bar). They get wired once the espyna layer ships.
type ModuleDeps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListCollectionMethods  func(ctx context.Context, req *cmpb.ListCollectionMethodsRequest) (*cmpb.ListCollectionMethodsResponse, error)
	ReadCollectionMethod   func(ctx context.Context, req *cmpb.ReadCollectionMethodRequest) (*cmpb.ReadCollectionMethodResponse, error)
	CreateCollectionMethod func(ctx context.Context, req *cmpb.CreateCollectionMethodRequest) (*cmpb.CreateCollectionMethodResponse, error)
	UpdateCollectionMethod func(ctx context.Context, req *cmpb.UpdateCollectionMethodRequest) (*cmpb.UpdateCollectionMethodResponse, error)
	DeleteCollectionMethod func(ctx context.Context, req *cmpb.DeleteCollectionMethodRequest) (*cmpb.DeleteCollectionMethodResponse, error)
}

// Module holds all constructed collection_method views.
type Module struct {
	routes    centymo.CollectionMethodRoutes
	List      view.View
	Detail    view.View
	TabAction view.View
	Add       view.View
	Edit      view.View
	Fragment  view.View
	Delete    view.View
}

// NewModule creates the collection_method module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &cmaction.Deps{
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		CreateCollectionMethod: deps.CreateCollectionMethod,
		ReadCollectionMethod:   deps.ReadCollectionMethod,
		UpdateCollectionMethod: deps.UpdateCollectionMethod,
		DeleteCollectionMethod: deps.DeleteCollectionMethod,
	}
	detailDeps := &cmdetail.DetailViewDeps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ReadCollectionMethod: deps.ReadCollectionMethod,
	}
	listDeps := &cmlist.ListViewDeps{
		Routes:                deps.Routes,
		ListCollectionMethods: deps.ListCollectionMethods,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
	}

	return &Module{
		routes:    deps.Routes,
		List:      cmlist.NewView(listDeps),
		Detail:    cmdetail.NewView(detailDeps),
		TabAction: cmdetail.NewTabAction(detailDeps),
		Add:       cmaction.NewAddAction(actionDeps),
		Edit:      cmaction.NewEditAction(actionDeps),
		Fragment:  cmaction.NewFragmentAction(actionDeps),
		Delete:    cmaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all collection_method routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)

	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.GET(m.routes.FragmentURL, m.Fragment)
	r.POST(m.routes.DeleteURL, m.Delete)
}
