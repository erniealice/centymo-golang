package treasury

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
	raaction "github.com/erniealice/centymo-golang/domain/treasury/collection_application/action"
)

// CollectionApplicationModuleDeps holds all dependencies for the receive-and-apply
// module. UseCases is nil when the espyna use cases are not wired; the module
// then fails closed (routes answer 503, never fake data).
type CollectionApplicationModuleDeps struct {
	Routes       ra.Routes
	Labels       ra.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *ra.UseCases

	// CollectionDetailURL is the collection detail route ({id}) the success
	// redirect targets; DefaultCurrency preselects the currency select.
	CollectionDetailURL string
	DefaultCurrency     string

	// SearchClientURL is the async client search route the drawer's client
	// picker uses (empty = a select of every client, paged through).
	SearchClientURL string
}

// CollectionApplicationModule holds all constructed receive-and-apply views.
type CollectionApplicationModule struct {
	routes ra.Routes

	ReceiveApply, Preview, Reverse view.View
}

// NewCollectionApplicationModule wires the receive-and-apply views.
func NewCollectionApplicationModule(deps *CollectionApplicationModuleDeps) *CollectionApplicationModule {
	if deps == nil {
		deps = &CollectionApplicationModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &ra.UseCases{} // fail closed: every closure nil
	}
	actionDeps := &raaction.Deps{
		Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, UseCases: ucs,
		CollectionDetailURL: deps.CollectionDetailURL, DefaultCurrency: deps.DefaultCurrency,
		SearchClientURL: deps.SearchClientURL,
	}
	return &CollectionApplicationModule{
		routes:       deps.Routes,
		ReceiveApply: raaction.NewReceiveApplyAction(actionDeps),
		Preview:      raaction.NewPreviewAction(actionDeps),
		Reverse:      raaction.NewReverseAction(actionDeps),
	}
}

// RegisterRoutes registers the receive-and-apply routes.
func (m *CollectionApplicationModule) RegisterRoutes(r view.RouteRegistrar) {
	rt := m.routes
	r.GET(rt.ReceiveApplyURL, m.ReceiveApply)
	r.POST(rt.ReceiveApplyURL, m.ReceiveApply)
	r.GET(rt.PreviewURL, m.Preview)
	r.POST(rt.ReverseURL, m.Reverse)
}
