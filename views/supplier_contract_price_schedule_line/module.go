// Package supplier_contract_price_schedule_line renders the inline-child
// CRUD modals for SupplierContractPriceScheduleLine rows under a parent
// SupplierContractPriceSchedule detail page. Mirrors the supplier_contract_line
// inline-child convention.
//
// MAIN-THREAD WIRING NOTE: this module's RegisterRoutes call must be added to
// packages/centymo-golang/block/block.go alongside the existing
// supplier_contract_line module. Mirror the supplier_contract_line wiring:
//
//   if cfg.wantSupplierContractPriceScheduleLine() {
//       deps := &supplier_contract_price_schedule_line.ModuleDeps{ ... }
//       mod := supplier_contract_price_schedule_line.NewModule(deps)
//       mod.RegisterRoutes(ctx.Routes)
//   }
//
// The block.go edit is intentionally deferred to the main thread because a
// sister Wave 3 CYC agent is currently editing block.go.
package supplier_contract_price_schedule_line

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	scpslaction "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule_line/action"

	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	scpslpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the SCPSL module.
type ModuleDeps struct {
	Routes       centymo.SupplierContractPriceScheduleRoutes
	Labels       centymo.SupplierContractPriceScheduleLabels
	CommonLabels pyeza.CommonLabels

	CreateSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.CreateSupplierContractPriceScheduleLineRequest) (*scpslpb.CreateSupplierContractPriceScheduleLineResponse, error)
	ReadSupplierContractPriceScheduleLine   func(ctx context.Context, req *scpslpb.ReadSupplierContractPriceScheduleLineRequest) (*scpslpb.ReadSupplierContractPriceScheduleLineResponse, error)
	UpdateSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.UpdateSupplierContractPriceScheduleLineRequest) (*scpslpb.UpdateSupplierContractPriceScheduleLineResponse, error)
	DeleteSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.DeleteSupplierContractPriceScheduleLineRequest) (*scpslpb.DeleteSupplierContractPriceScheduleLineResponse, error)

	ListSupplierContractLines func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)
}

// Module holds all constructed SCPSL views.
type Module struct {
	routes centymo.SupplierContractPriceScheduleRoutes
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the SCPSL module.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &scpslaction.Deps{
		Routes:                                  deps.Routes,
		Labels:                                  deps.Labels,
		CommonLabels:                            deps.CommonLabels,
		CreateSupplierContractPriceScheduleLine: deps.CreateSupplierContractPriceScheduleLine,
		ReadSupplierContractPriceScheduleLine:   deps.ReadSupplierContractPriceScheduleLine,
		UpdateSupplierContractPriceScheduleLine: deps.UpdateSupplierContractPriceScheduleLine,
		DeleteSupplierContractPriceScheduleLine: deps.DeleteSupplierContractPriceScheduleLine,
		ListSupplierContractLines:               deps.ListSupplierContractLines,
	}

	return &Module{
		routes: deps.Routes,
		Add:    scpslaction.NewAddAction(actionDeps),
		Edit:   scpslaction.NewEditAction(actionDeps),
		Delete: scpslaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all SCPSL action routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.Add != nil {
		r.GET(m.routes.LineAddURL, m.Add)
		r.POST(m.routes.LineAddURL, m.Add)
	}
	if m.Edit != nil {
		r.GET(m.routes.LineEditURL, m.Edit)
		r.POST(m.routes.LineEditURL, m.Edit)
	}
	if m.Delete != nil {
		r.POST(m.routes.LineDeleteURL, m.Delete)
	}
}
