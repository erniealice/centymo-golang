package disbursement_method

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	dmaction "github.com/erniealice/centymo-golang/views/disbursement_method/action"
	dmdetail "github.com/erniealice/centymo-golang/views/disbursement_method/detail"
	dmlist "github.com/erniealice/centymo-golang/views/disbursement_method/list"

	dmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_method"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the disbursement_method module
// (treasury-domain-rebuild Stage 1, pages.md §C-5). Buying-side mirror of the
// collection_method module. CRUD closures are OPTIONAL (espyna use cases pending).
type ModuleDeps struct {
	Routes       centymo.DisbursementMethodRoutes
	Labels       centymo.DisbursementMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListDisbursementMethods  func(ctx context.Context, req *dmpb.ListDisbursementMethodsRequest) (*dmpb.ListDisbursementMethodsResponse, error)
	ReadDisbursementMethod   func(ctx context.Context, req *dmpb.ReadDisbursementMethodRequest) (*dmpb.ReadDisbursementMethodResponse, error)
	CreateDisbursementMethod func(ctx context.Context, req *dmpb.CreateDisbursementMethodRequest) (*dmpb.CreateDisbursementMethodResponse, error)
	UpdateDisbursementMethod func(ctx context.Context, req *dmpb.UpdateDisbursementMethodRequest) (*dmpb.UpdateDisbursementMethodResponse, error)
	DeleteDisbursementMethod func(ctx context.Context, req *dmpb.DeleteDisbursementMethodRequest) (*dmpb.DeleteDisbursementMethodResponse, error)
}

// Module holds all constructed disbursement_method views.
type Module struct {
	routes    centymo.DisbursementMethodRoutes
	List      view.View
	Detail    view.View
	TabAction view.View
	Add       view.View
	Edit      view.View
	Fragment  view.View
	Delete    view.View
}

// NewModule creates the disbursement_method module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &dmaction.Deps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		CreateDisbursementMethod: deps.CreateDisbursementMethod,
		ReadDisbursementMethod:   deps.ReadDisbursementMethod,
		UpdateDisbursementMethod: deps.UpdateDisbursementMethod,
		DeleteDisbursementMethod: deps.DeleteDisbursementMethod,
	}
	detailDeps := &dmdetail.DetailViewDeps{
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ReadDisbursementMethod: deps.ReadDisbursementMethod,
	}
	listDeps := &dmlist.ListViewDeps{
		Routes:                  deps.Routes,
		ListDisbursementMethods: deps.ListDisbursementMethods,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		TableLabels:             deps.TableLabels,
	}

	return &Module{
		routes:    deps.Routes,
		List:      dmlist.NewView(listDeps),
		Detail:    dmdetail.NewView(detailDeps),
		TabAction: dmdetail.NewTabAction(detailDeps),
		Add:       dmaction.NewAddAction(actionDeps),
		Edit:      dmaction.NewEditAction(actionDeps),
		Fragment:  dmaction.NewFragmentAction(actionDeps),
		Delete:    dmaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all disbursement_method routes.
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
