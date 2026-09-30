package subscription

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	bcaction "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/action"
	bclist "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/list"
)

// BillableChargeModuleDeps holds all dependencies for the billable charge
// module. UseCases is nil when the espyna use cases are not wired; the module
// then fails closed (routes answer 503, never fake data).
type BillableChargeModuleDeps struct {
	Routes         bc.Routes
	Labels         bc.Labels
	RecoveryLabels rd.Labels
	CommonLabels   pyeza.CommonLabels
	TableLabels    types.TableLabels
	UseCases       *bc.UseCases
}

// BillableChargeModule holds all constructed billable charge views.
type BillableChargeModule struct {
	routes bc.Routes

	List, Table   view.View
	Issue, Adjust view.View
}

// NewBillableChargeModule wires every billable charge view.
func NewBillableChargeModule(deps *BillableChargeModuleDeps) *BillableChargeModule {
	if deps == nil {
		deps = &BillableChargeModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &bc.UseCases{} // fail closed: every closure nil
	}
	listDeps := &bclist.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs}
	actionDeps := &bcaction.Deps{Routes: deps.Routes, Labels: deps.Labels, RecoveryLabels: deps.RecoveryLabels, CommonLabels: deps.CommonLabels, UseCases: ucs}
	return &BillableChargeModule{
		routes: deps.Routes,
		List:   bclist.NewView(listDeps),
		Table:  bclist.NewTableView(listDeps),
		Issue:  bcaction.NewIssueAction(actionDeps),
		Adjust: bcaction.NewAdjustAction(actionDeps),
	}
}

// RegisterRoutes registers every billable charge route.
func (m *BillableChargeModule) RegisterRoutes(r view.RouteRegistrar) {
	rt := m.routes
	r.GET(rt.ListURL, m.List)
	r.GET(rt.TableURL, m.Table)
	r.GET(rt.IssueURL, m.Issue)
	r.POST(rt.IssueURL, m.Issue)
	r.GET(rt.AdjustURL, m.Adjust)
	r.POST(rt.AdjustURL, m.Adjust)
}
