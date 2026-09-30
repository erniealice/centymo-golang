package expenditure

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	abaction "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch/action"
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
)

// AllocationBatchModuleDeps holds the dependencies of the cost allocation module
// (usage-and-pass-through S1). UseCases is nil when the espyna use cases are not
// wired; the module then fails closed (routes answer 503).
type AllocationBatchModuleDeps struct {
	Routes          ab.Routes
	Labels          ab.Labels
	ComponentLabels csc.Labels
	CommonLabels    pyeza.CommonLabels
	UseCases        *ab.UseCases
}

// AllocationBatchModule holds the constructed views.
type AllocationBatchModule struct {
	routes ab.Routes

	Allocate, Preview, Publish, View view.View
}

// NewAllocationBatchModule wires the allocate / preview / publish / view drawers.
func NewAllocationBatchModule(deps *AllocationBatchModuleDeps) *AllocationBatchModule {
	if deps == nil {
		deps = &AllocationBatchModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &ab.UseCases{} // fail closed: every closure nil
	}
	actionDeps := &abaction.Deps{Routes: deps.Routes, Labels: deps.Labels, ComponentLabels: deps.ComponentLabels, CommonLabels: deps.CommonLabels, UseCases: ucs}
	return &AllocationBatchModule{
		routes:   deps.Routes,
		Allocate: abaction.NewAllocateAction(actionDeps),
		Preview:  abaction.NewPreviewAction(actionDeps),
		Publish:  abaction.NewPublishAction(actionDeps),
		View:     abaction.NewViewAction(actionDeps),
	}
}

// RegisterRoutes registers every cost allocation route.
func (m *AllocationBatchModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.AllocateURL, m.Allocate)
	r.POST(m.routes.AllocateURL, m.Allocate)
	r.POST(m.routes.PreviewURL, m.Preview)
	r.POST(m.routes.PublishURL, m.Publish)
	r.GET(m.routes.ViewURL, m.View)
}
