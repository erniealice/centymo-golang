// Package supplier_billing_event — module wiring for the buying-side
// SupplierBillingEvent list / detail / Recognize surfaces.
//
// 20260517-advance-cash-events Plan B Phase 7.
package supplier_billing_eventmodule

import (
	"context"

	supplierbillingeventaction "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event/action"
	supplierbillingeventdetail "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event/detail"
	supplierbillingeventlist "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event/list"
	sib_treasury_advancesdashboard "github.com/erniealice/centymo-golang/domain/treasury/advancesdashboard"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event"
	supplierbillingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_billing_event"
)

// ModuleDeps holds all dependencies for the supplier_billing_event module.
type ModuleDeps struct {
	Routes       sib_treasury_advancesdashboard.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSupplierBillingEvents func(ctx context.Context, req *supplierbillingeventpb.ListSupplierBillingEventsRequest) (*supplierbillingeventpb.ListSupplierBillingEventsResponse, error)
	ReadSupplierBillingEvent  func(ctx context.Context, req *supplierbillingeventpb.ReadSupplierBillingEventRequest) (*supplierbillingeventpb.ReadSupplierBillingEventResponse, error)

	// 20260517-advance-cash-events Plan B Phase 7 — buying-side MILESTONE
	// recognize closure. Nil-safe — Recognize button surfaces a disabled
	// state when unwired.
	Recognize supplierbillingeventaction.RecognizeMilestoneAdvanceDisbursementFn
}

// Module holds all constructed supplier_billing_event views.
type Module struct {
	routes    sib_treasury_advancesdashboard.Routes
	List      view.View
	Detail    view.View
	Recognize view.View
}

// NewModule creates the supplier_billing_event module with all views wired.
func NewModule(deps ModuleDeps) *Module {
	listView := supplierbillingeventlist.NewView(&supplierbillingeventlist.ListViewDeps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		TableLabels:               deps.TableLabels,
		ListSupplierBillingEvents: deps.ListSupplierBillingEvents,
	})
	detailView := supplierbillingeventdetail.NewView(&supplierbillingeventdetail.DetailViewDeps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		ReadSupplierBillingEvent: deps.ReadSupplierBillingEvent,
	})
	var recognizeView view.View
	if deps.Recognize != nil {
		recognizeView = supplierbillingeventaction.NewRecognizeAction(deps.Recognize, deps.Labels.Errors)
	}
	return &Module{
		routes:    deps.Routes,
		List:      listView,
		Detail:    detailView,
		Recognize: recognizeView,
	}
}

// RegisterRoutes mounts the module's views on the given route registrar.
func (m *Module) RegisterRoutes(routes view.RouteRegistrar) {
	if m.List != nil && m.routes.SupplierBillingEventListURL != "" {
		routes.GET(m.routes.SupplierBillingEventListURL, m.List)
	}
	if m.Detail != nil && m.routes.SupplierBillingEventDetailURL != "" {
		routes.GET(m.routes.SupplierBillingEventDetailURL, m.Detail)
	}
	if m.Recognize != nil && m.routes.SupplierBillingEventRecognizeURL != "" {
		routes.POST(m.routes.SupplierBillingEventRecognizeURL, m.Recognize)
	}
}
