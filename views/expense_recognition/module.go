// Package expense_recognition is the centymo views package for ExpenseRecognition.
//
// MAIN-THREAD WIRING NOTE (block.go):
//   Block-level wiring (NewModule call + dependency injection from espyna use
//   cases) is intentionally DEFERRED to the main thread. This package
//   exposes Module + ModuleDeps + RegisterRoutes following the same
//   convention as views/supplier_contract/. The integrator must:
//     1. Register routes/labels in apps/service-admin sidebar.
//     2. Add a centymo block.go entry that constructs ModuleDeps from the
//        domain providers and calls NewModule + RegisterRoutes.
//     3. Wire the per-tier translation file via translations.go.
//   See plan §7 (Phase P10) for the full integrator checklist.
package expense_recognition

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	expenserecognitionaction "github.com/erniealice/centymo-golang/views/expense_recognition/action"
	expenserecognitiondetail "github.com/erniealice/centymo-golang/views/expense_recognition/detail"
	expenserecognitionlist "github.com/erniealice/centymo-golang/views/expense_recognition/list"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	expenserecognitionlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the expense_recognition module.
type ModuleDeps struct {
	Routes       centymo.ExpenseRecognitionRoutes
	Labels       centymo.ExpenseRecognitionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Core CRUD (no Create — recognitions are created BY use case)
	ListExpenseRecognitions  func(ctx context.Context, req *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)
	ReadExpenseRecognition   func(ctx context.Context, req *expenserecognitionpb.ReadExpenseRecognitionRequest) (*expenserecognitionpb.ReadExpenseRecognitionResponse, error)
	DeleteExpenseRecognition func(ctx context.Context, req *expenserecognitionpb.DeleteExpenseRecognitionRequest) (*expenserecognitionpb.DeleteExpenseRecognitionResponse, error)

	// Inline child — lines
	ListExpenseRecognitionLines func(ctx context.Context, req *expenserecognitionlinepb.ListExpenseRecognitionLinesRequest) (*expenserecognitionlinepb.ListExpenseRecognitionLinesResponse, error)

	// Workflow — Reverse use case (espyna)
	ReverseExpenseRecognition func(ctx context.Context, id, reason string) error

	// Workflow — Recognition use cases (espyna). Optional; when nil the
	// corresponding action returns a 422 indicating it isn't wired.
	RecognizeFromExpenditure expenserecognitionaction.RecognizeFromExpenditureFunc
	RecognizeFromContract    expenserecognitionaction.RecognizeFromContractFunc
}

// Module holds all constructed expense_recognition views.
type Module struct {
	routes                   centymo.ExpenseRecognitionRoutes
	List                     view.View
	Detail                   view.View
	TabAction                view.View
	Delete                   view.View
	Reverse                  view.View
	RecognizeFromExpenditure view.View
	RecognizeFromContract    view.View
}

// NewModule creates the expense_recognition module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &expenserecognitionlist.ListViewDeps{
		Routes:                  deps.Routes,
		ListExpenseRecognitions: deps.ListExpenseRecognitions,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		TableLabels:             deps.TableLabels,
	}

	detailDeps := &expenserecognitiondetail.DetailViewDeps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ReadExpenseRecognition:      deps.ReadExpenseRecognition,
		ListExpenseRecognitionLines: deps.ListExpenseRecognitionLines,
		ReverseExpenseRecognition:   deps.ReverseExpenseRecognition,
	}

	actionDeps := &expenserecognitionaction.Deps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		DeleteExpenseRecognition: deps.DeleteExpenseRecognition,
	}

	m := &Module{
		routes: deps.Routes,
		List:   expenserecognitionlist.NewView(listDeps),
		Delete: expenserecognitionaction.NewDeleteAction(actionDeps),
	}
	if deps.ReadExpenseRecognition != nil {
		m.Detail = expenserecognitiondetail.NewView(detailDeps)
		m.TabAction = expenserecognitiondetail.NewTabAction(detailDeps)
		m.Reverse = expenserecognitiondetail.NewReverseAction(detailDeps)
	}
	// Recognition action handlers: register unconditionally so the routes
	// respond with 422 (not 405) even when use cases aren't wired yet.
	m.RecognizeFromExpenditure = expenserecognitionaction.NewRecognizeFromExpenditureAction(deps.RecognizeFromExpenditure)
	m.RecognizeFromContract = expenserecognitionaction.NewRecognizeFromContractAction(deps.RecognizeFromContract)
	return m
}

// RegisterRoutes registers all expense_recognition routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.POST(m.routes.DeleteURL, m.Delete)
	if m.Detail != nil {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.Reverse != nil {
		r.POST(m.routes.ReverseURL, m.Reverse)
	}
	if m.RecognizeFromExpenditure != nil {
		r.POST(m.routes.RecognizeFromExpenditureURL, m.RecognizeFromExpenditure)
	}
	if m.RecognizeFromContract != nil {
		r.POST(m.routes.RecognizeFromContractURL, m.RecognizeFromContract)
	}
}
