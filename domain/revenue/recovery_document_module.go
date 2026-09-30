package revenue

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rdaction "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/action"
	rddetail "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/detail"
	rdlist "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/list"
	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
)

// RecoveryDocumentModuleDeps holds all dependencies for the recovery document
// module. UseCases is nil when the espyna use cases are not wired; the module
// then fails closed (routes answer 503, never fake data).
type RecoveryDocumentModuleDeps struct {
	Routes       rd.Routes
	Labels       rd.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *rd.UseCases

	// ApplicationLabels / ReverseURL / ReceiveApplyURL come from the sibling
	// collection_application unit (applications tab + payment actions).
	ApplicationLabels ra.Labels
	ReverseURL        string
	ReceiveApplyURL   string

	// ApplicationsTableID is the collection_application unit's applications table id.
	ApplicationsTableID string
}

// RecoveryDocumentModule holds all constructed recovery document views.
type RecoveryDocumentModule struct {
	routes rd.Routes

	List, Table, Detail, TabAction view.View
	Void                           view.View
}

// NewRecoveryDocumentModule wires every recovery document view.
func NewRecoveryDocumentModule(deps *RecoveryDocumentModuleDeps) *RecoveryDocumentModule {
	if deps == nil {
		deps = &RecoveryDocumentModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &rd.UseCases{} // fail closed: every closure nil
	}
	listDeps := &rdlist.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs}
	detailDeps := &rddetail.Deps{
		Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs,
		ApplicationLabels: deps.ApplicationLabels, ReverseURL: deps.ReverseURL, ReceiveApplyURL: deps.ReceiveApplyURL,
		ApplicationsTableID: deps.ApplicationsTableID,
	}
	actionDeps := &rdaction.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, UseCases: ucs}
	return &RecoveryDocumentModule{
		routes:    deps.Routes,
		List:      rdlist.NewView(listDeps),
		Table:     rdlist.NewTableView(listDeps),
		Detail:    rddetail.NewView(detailDeps),
		TabAction: rddetail.NewTabAction(detailDeps),
		Void:      rdaction.NewVoidAction(actionDeps),
	}
}

// RegisterRoutes registers every recovery document route.
func (m *RecoveryDocumentModule) RegisterRoutes(r view.RouteRegistrar) {
	rt := m.routes
	r.GET(rt.ListURL, m.List)
	r.GET(rt.TableURL, m.Table)
	r.GET(rt.DetailURL, m.Detail)
	r.GET(rt.TabActionURL, m.TabAction)
	r.GET(rt.VoidURL, m.Void)
	r.POST(rt.VoidURL, m.Void)
}
