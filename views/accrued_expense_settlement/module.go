// Package accrued_expense_settlement is the centymo views package for
// AccruedExpenseSettlement (inline child of AccruedExpense). Settlements
// render inside the parent's Settlements tab — no standalone list view.
//
// MAIN-THREAD WIRING NOTE (block.go):
//   Block-level wiring is intentionally DEFERRED to the main thread. Routes
//   share AccruedExpenseRoutes with the parent module. The integrator must
//   construct ModuleDeps from the espyna AccruedExpenseSettlement use case
//   group and call NewModule + RegisterRoutes inside the centymo block.go
//   entry that already wires the parent accrued_expense module.
//   See plan §7 (Phase P10) for the integrator checklist.
package accrued_expense_settlement

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	settlementaction "github.com/erniealice/centymo-golang/views/accrued_expense_settlement/action"

	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the accrued_expense_settlement module.
type ModuleDeps struct {
	Routes       centymo.AccruedExpenseRoutes
	Labels       centymo.AccruedExpenseLabels
	CommonLabels pyeza.CommonLabels

	CreateAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.CreateAccruedExpenseSettlementRequest) (*accruedexpensepb.CreateAccruedExpenseSettlementResponse, error)
	ReadAccruedExpenseSettlement   func(ctx context.Context, req *accruedexpensepb.ReadAccruedExpenseSettlementRequest) (*accruedexpensepb.ReadAccruedExpenseSettlementResponse, error)
	UpdateAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.UpdateAccruedExpenseSettlementRequest) (*accruedexpensepb.UpdateAccruedExpenseSettlementResponse, error)
	DeleteAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.DeleteAccruedExpenseSettlementRequest) (*accruedexpensepb.DeleteAccruedExpenseSettlementResponse, error)
	ListExpenditures               func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// Module holds all constructed settlement views.
type Module struct {
	routes centymo.AccruedExpenseRoutes
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the settlement module with all action handlers.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &settlementaction.Deps{
		Routes:                         deps.Routes,
		Labels:                         deps.Labels,
		CommonLabels:                   deps.CommonLabels,
		CreateAccruedExpenseSettlement: deps.CreateAccruedExpenseSettlement,
		ReadAccruedExpenseSettlement:   deps.ReadAccruedExpenseSettlement,
		UpdateAccruedExpenseSettlement: deps.UpdateAccruedExpenseSettlement,
		DeleteAccruedExpenseSettlement: deps.DeleteAccruedExpenseSettlement,
		ListExpenditures:               deps.ListExpenditures,
	}
	return &Module{
		routes: deps.Routes,
		Add:    settlementaction.NewAddAction(actionDeps),
		Edit:   settlementaction.NewEditAction(actionDeps),
		Delete: settlementaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all settlement action routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.SettlementAddURL, m.Add)
	r.POST(m.routes.SettlementAddURL, m.Add)
	r.GET(m.routes.SettlementEditURL, m.Edit)
	r.POST(m.routes.SettlementEditURL, m.Edit)
	r.POST(m.routes.SettlementDeleteURL, m.Delete)
}
