package expenditure

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	cscaction "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/action"
	csctable "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/table"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
)

// CostSourceComponentModuleDeps holds the dependencies of the recoverable cost
// line module (usage-and-pass-through S1). UseCases is nil when the espyna use
// cases are not wired; the module then fails closed (routes answer 503).
type CostSourceComponentModuleDeps struct {
	Routes       csc.Routes
	Labels       csc.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *csc.UseCases

	// AllocationRoutes feed the per-row Allocate / View allocation actions.
	AllocationRoutes ab.Routes

	// ReadExpenditure reads the parent expenditure (required: the Add drawer proves the URL parent
	// exists in the workspace through it).
	ReadExpenditure func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
}

// CostSourceComponentModule holds the constructed views.
type CostSourceComponentModule struct {
	routes csc.Routes

	Add, Edit, Delete, Table view.View
}

// NewCostSourceComponentModule wires the add / edit / delete drawers and the table refresh.
func NewCostSourceComponentModule(deps *CostSourceComponentModuleDeps) *CostSourceComponentModule {
	if deps == nil {
		deps = &CostSourceComponentModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &csc.UseCases{} // fail closed: every closure nil
	}
	readCurrency := func(ctx context.Context, expenditureID string) string {
		if deps.ReadExpenditure == nil {
			return ""
		}
		resp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{Data: &expenditurepb.Expenditure{Id: expenditureID}})
		if err != nil || len(resp.GetData()) == 0 {
			return ""
		}
		return resp.GetData()[0].GetCurrency()
	}
	readParent := func(ctx context.Context, expenditureID string) (string, bool) {
		if deps.ReadExpenditure == nil || expenditureID == "" {
			return "", false
		}
		resp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{Data: &expenditurepb.Expenditure{Id: expenditureID}})
		if err != nil || len(resp.GetData()) == 0 || resp.GetData()[0] == nil {
			return "", false
		}
		return resp.GetData()[0].GetCurrency(), true
	}
	actionDeps := &cscaction.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, UseCases: ucs, ReadParent: readParent}
	return &CostSourceComponentModule{
		routes: deps.Routes,
		Add:    cscaction.NewAddAction(actionDeps),
		Edit:   cscaction.NewEditAction(actionDeps),
		Delete: cscaction.NewDeleteAction(actionDeps),
		Table: csctable.NewView(&csctable.Deps{
			Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs,
			AllocateURL: deps.AllocationRoutes.AllocateURL, ViewAllocationURL: deps.AllocationRoutes.ViewURL, ReadCurrency: readCurrency,
		}),
	}
}

// RegisterRoutes registers every recoverable cost line route.
func (m *CostSourceComponentModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.GET(m.routes.TableURL, m.Table)
}
