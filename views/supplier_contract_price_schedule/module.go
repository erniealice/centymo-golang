// Package supplier_contract_price_schedule renders the views for the
// SupplierContractPriceSchedule master entity (date-windowed pricing layered
// on top of a supplier contract).
//
// MAIN-THREAD WIRING NOTE: this module's RegisterRoutes call must be added to
// packages/centymo-golang/block/block.go alongside the existing
// supplier_contract module. Mirror the supplier_contract_line wiring pattern:
//
//   if cfg.wantSupplierContractPriceSchedule() {
//       deps := &supplier_contract_price_schedule.ModuleDeps{ ... }
//       mod := supplier_contract_price_schedule.NewModule(deps)
//       mod.RegisterRoutes(ctx.Routes)
//   }
//
// The block.go edit is intentionally deferred to the main thread because a
// sister Wave 3 CYC agent is currently editing block.go.
package supplier_contract_price_schedule

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	scpsaction "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule/action"
	scpsdetail "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule/detail"
	scpslist "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule/list"

	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
	scpslpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the supplier_contract_price_schedule module.
type ModuleDeps struct {
	Routes       centymo.SupplierContractPriceScheduleRoutes
	Labels       centymo.SupplierContractPriceScheduleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Core CRUD
	ListSupplierContractPriceSchedules  func(ctx context.Context, req *scpspb.ListSupplierContractPriceSchedulesRequest) (*scpspb.ListSupplierContractPriceSchedulesResponse, error)
	ReadSupplierContractPriceSchedule   func(ctx context.Context, req *scpspb.ReadSupplierContractPriceScheduleRequest) (*scpspb.ReadSupplierContractPriceScheduleResponse, error)
	CreateSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.CreateSupplierContractPriceScheduleRequest) (*scpspb.CreateSupplierContractPriceScheduleResponse, error)
	UpdateSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.UpdateSupplierContractPriceScheduleRequest) (*scpspb.UpdateSupplierContractPriceScheduleResponse, error)
	DeleteSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.DeleteSupplierContractPriceScheduleRequest) (*scpspb.DeleteSupplierContractPriceScheduleResponse, error)

	// Workflow (closures injected by block.go)
	ActivateSupplierContractPriceSchedule  func(ctx context.Context, id string) error
	SupersedeSupplierContractPriceSchedule func(ctx context.Context, id, reason string) error
	SetSupplierContractPriceScheduleStatus func(ctx context.Context, id, status string) error

	// Child entity — schedule lines
	ListSupplierContractPriceScheduleLines func(ctx context.Context, req *scpslpb.ListSupplierContractPriceScheduleLinesRequest) (*scpslpb.ListSupplierContractPriceScheduleLinesResponse, error)

	// Related entities for dropdowns + linked tabs
	ListSupplierContracts     func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
	ListSupplierContractLines func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)
}

// Module holds all constructed supplier_contract_price_schedule views.
type Module struct {
	routes centymo.SupplierContractPriceScheduleRoutes

	List          view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	SetStatus     view.View
	BulkSetStatus view.View
	Activate      view.View
	Supersede     view.View
}

// NewModule creates the supplier_contract_price_schedule module with all
// views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &scpsaction.Deps{
		Routes:                                 deps.Routes,
		Labels:                                 deps.Labels,
		CommonLabels:                           deps.CommonLabels,
		CreateSupplierContractPriceSchedule:    deps.CreateSupplierContractPriceSchedule,
		ReadSupplierContractPriceSchedule:      deps.ReadSupplierContractPriceSchedule,
		UpdateSupplierContractPriceSchedule:    deps.UpdateSupplierContractPriceSchedule,
		DeleteSupplierContractPriceSchedule:    deps.DeleteSupplierContractPriceSchedule,
		ActivateSupplierContractPriceSchedule:  deps.ActivateSupplierContractPriceSchedule,
		SupersedeSupplierContractPriceSchedule: deps.SupersedeSupplierContractPriceSchedule,
		SetSupplierContractPriceScheduleStatus: deps.SetSupplierContractPriceScheduleStatus,
		ListSupplierContracts:                  deps.ListSupplierContracts,
	}

	listDeps := &scpslist.ListViewDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
		ListSupplierContractPriceSchedules: deps.ListSupplierContractPriceSchedules,
		ListSupplierContracts:              deps.ListSupplierContracts,
	}

	detailDeps := &scpsdetail.DetailViewDeps{
		Routes:                                 deps.Routes,
		Labels:                                 deps.Labels,
		CommonLabels:                           deps.CommonLabels,
		TableLabels:                            deps.TableLabels,
		ReadSupplierContractPriceSchedule:      deps.ReadSupplierContractPriceSchedule,
		ListSupplierContractPriceScheduleLines: deps.ListSupplierContractPriceScheduleLines,
		ListSupplierContractLines:              deps.ListSupplierContractLines,
	}

	m := &Module{
		routes:        deps.Routes,
		Add:           scpsaction.NewAddAction(actionDeps),
		Edit:          scpsaction.NewEditAction(actionDeps),
		Delete:        scpsaction.NewDeleteAction(actionDeps),
		SetStatus:     scpsaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: scpsaction.NewBulkSetStatusAction(actionDeps),
		Activate:      scpsaction.NewActivateAction(actionDeps),
		Supersede:     scpsaction.NewSupersedeAction(actionDeps),
		List:          scpslist.NewView(listDeps),
	}

	if deps.ReadSupplierContractPriceSchedule != nil {
		m.Detail = scpsdetail.NewView(detailDeps)
		m.TabAction = scpsdetail.NewTabAction(detailDeps)
	}

	return m
}

// RegisterRoutes registers all SupplierContractPriceSchedule routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	r.POST(m.routes.ActivateURL, m.Activate)
	r.POST(m.routes.SupersedeURL, m.Supersede)

	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
}
