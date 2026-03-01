package expenditure

import (
	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenditurelist "github.com/erniealice/centymo-golang/views/expenditure/list"
)

// ModuleDeps holds all dependencies for the expenditure module.
type ModuleDeps struct {
	Routes       centymo.ExpenditureRoutes
	DB           centymo.DataSource
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// Module holds all constructed expenditure views.
type Module struct {
	routes            centymo.ExpenditureRoutes
	PurchaseList      view.View
	PurchaseDashboard view.View
	ExpenseList       view.View
	ExpenseDashboard  view.View
}

// NewModule creates the expenditure module with purchase and expense views.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{
		routes: deps.Routes,
		PurchaseList: expenditurelist.NewView(&expenditurelist.Deps{
			DB:              deps.DB,
			RefreshURL:      deps.Routes.PurchaseListURL,
			ExpenditureType: "purchase",
			Labels:          deps.Labels,
			CommonLabels:    deps.CommonLabels,
			TableLabels:     deps.TableLabels,
		}),
		ExpenseList: expenditurelist.NewView(&expenditurelist.Deps{
			DB:              deps.DB,
			RefreshURL:      deps.Routes.ExpenseListURL,
			ExpenditureType: "expense",
			Labels:          deps.Labels,
			CommonLabels:    deps.CommonLabels,
			TableLabels:     deps.TableLabels,
		}),
		// Dashboards use same list view for now (will be enhanced later)
		PurchaseDashboard: expenditurelist.NewView(&expenditurelist.Deps{
			DB:              deps.DB,
			RefreshURL:      deps.Routes.PurchaseListURL,
			ExpenditureType: "purchase",
			Labels:          deps.Labels,
			CommonLabels:    deps.CommonLabels,
			TableLabels:     deps.TableLabels,
		}),
		ExpenseDashboard: expenditurelist.NewView(&expenditurelist.Deps{
			DB:              deps.DB,
			RefreshURL:      deps.Routes.ExpenseListURL,
			ExpenditureType: "expense",
			Labels:          deps.Labels,
			CommonLabels:    deps.CommonLabels,
			TableLabels:     deps.TableLabels,
		}),
	}
}

// RegisterRoutes registers all expenditure routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.PurchaseListURL, m.PurchaseList)
	r.GET(m.routes.PurchaseDashboardURL, m.PurchaseDashboard)
	r.GET(m.routes.ExpenseListURL, m.ExpenseList)
	r.GET(m.routes.ExpenseDashboardURL, m.ExpenseDashboard)
}
