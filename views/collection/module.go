package collection

import (
	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	collectionaction "github.com/erniealice/centymo-golang/views/collection/action"
	collectiondetail "github.com/erniealice/centymo-golang/views/collection/detail"
	collectionlist "github.com/erniealice/centymo-golang/views/collection/list"
)

// ModuleDeps holds all dependencies for the collection module.
type ModuleDeps struct {
	Routes       centymo.CollectionRoutes
	DB           centymo.DataSource
	Labels       centymo.CollectionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// Module holds all constructed collection views.
type Module struct {
	routes        centymo.CollectionRoutes
	Dashboard     view.View
	List          view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

// NewModule creates the collection module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &collectionaction.Deps{
		Routes: deps.Routes,
		DB:     deps.DB,
	}

	detailDeps := &collectiondetail.Deps{
		Routes:       deps.Routes,
		DB:           deps.DB,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}

	listView := collectionlist.NewView(&collectionlist.Deps{
		Routes:       deps.Routes,
		DB:           deps.DB,
		RefreshURL:   deps.Routes.ListURL,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
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
		SetStatus:     collectionaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: collectionaction.NewBulkSetStatusAction(actionDeps),
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
}
