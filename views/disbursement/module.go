package disbursement

import (
	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	disbursementaction "github.com/erniealice/centymo-golang/views/disbursement/action"
	disbursementdetail "github.com/erniealice/centymo-golang/views/disbursement/detail"
	disbursementlist "github.com/erniealice/centymo-golang/views/disbursement/list"
)

// ModuleDeps holds all dependencies for the disbursement module.
type ModuleDeps struct {
	Routes       centymo.DisbursementRoutes
	DB           centymo.DataSource
	Labels       centymo.DisbursementLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// Module holds all constructed disbursement views.
type Module struct {
	routes        centymo.DisbursementRoutes
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

// NewModule creates the disbursement module with all views.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &disbursementlist.Deps{
		Routes:       deps.Routes,
		DB:           deps.DB,
		RefreshURL:   deps.Routes.ListURL,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}

	detailDeps := &disbursementdetail.Deps{
		Routes:       deps.Routes,
		DB:           deps.DB,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}

	actionDeps := &disbursementaction.Deps{
		Routes: deps.Routes,
		DB:     deps.DB,
		Labels: deps.Labels,
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
		SetStatus:     disbursementaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: disbursementaction.NewBulkSetStatusAction(actionDeps),
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
}
