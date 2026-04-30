// Package expense_recognition_line is the centymo views package for
// ExpenseRecognitionLine. It is an INLINE child of ExpenseRecognition's
// detail page Lines tab — there is no standalone list view.
//
// MAIN-THREAD WIRING NOTE (block.go):
//   Block-level wiring is intentionally DEFERRED to the main thread. The
//   integrator must construct ModuleDeps from the espyna use cases and call
//   NewModule + RegisterRoutes inside a centymo block.go entry. Routes use
//   the parent's ExpenseRecognitionRoutes — the integrator should pass the
//   same routes struct used to register the parent module.
//   See plan §7 (Phase P10) for the integrator checklist.
package expense_recognition_line

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	expenserecognitionlineaction "github.com/erniealice/centymo-golang/views/expense_recognition_line/action"

	expenserecognitionlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the expense_recognition_line module.
type ModuleDeps struct {
	Routes       centymo.ExpenseRecognitionRoutes
	Labels       centymo.ExpenseRecognitionLabels
	CommonLabels pyeza.CommonLabels

	CreateExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.CreateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.CreateExpenseRecognitionLineResponse, error)
	ReadExpenseRecognitionLine   func(ctx context.Context, req *expenserecognitionlinepb.ReadExpenseRecognitionLineRequest) (*expenserecognitionlinepb.ReadExpenseRecognitionLineResponse, error)
	UpdateExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.UpdateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.UpdateExpenseRecognitionLineResponse, error)
	DeleteExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.DeleteExpenseRecognitionLineRequest) (*expenserecognitionlinepb.DeleteExpenseRecognitionLineResponse, error)
}

// Module holds all constructed line views.
type Module struct {
	routes centymo.ExpenseRecognitionRoutes
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the expense_recognition_line module with all action handlers.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &expenserecognitionlineaction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		CreateExpenseRecognitionLine: deps.CreateExpenseRecognitionLine,
		ReadExpenseRecognitionLine:   deps.ReadExpenseRecognitionLine,
		UpdateExpenseRecognitionLine: deps.UpdateExpenseRecognitionLine,
		DeleteExpenseRecognitionLine: deps.DeleteExpenseRecognitionLine,
	}
	return &Module{
		routes: deps.Routes,
		Add:    expenserecognitionlineaction.NewAddAction(actionDeps),
		Edit:   expenserecognitionlineaction.NewEditAction(actionDeps),
		Delete: expenserecognitionlineaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all expense_recognition_line action routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.LineAddURL, m.Add)
	r.POST(m.routes.LineAddURL, m.Add)
	r.GET(m.routes.LineEditURL, m.Edit)
	r.POST(m.routes.LineEditURL, m.Edit)
	r.POST(m.routes.LineDeleteURL, m.Delete)
}
