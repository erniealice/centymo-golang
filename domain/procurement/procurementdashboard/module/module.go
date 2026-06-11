// Package procurementdashboard is a composition surface for procurement operations.
// It owns NO proto entity — it composes views over existing entities:
// SupplierContract, ProcurementRequest, and Expenditure.
//
// Pattern mirrors cyta-golang/views (schedule app): pure workflow surface,
// no "procurement" proto, no new domain layer.
package procurementdashboardmodule

import (
	"context"

	epkg "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard"
	procurementdashboard "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard/dashboard"
	procurementrecurrence "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard/recurrence_drafts"
	procurementrenewals "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard/renewals"
	procurementutilization "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard/utilization"
	procurementvariance "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard/variance"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the Procurement Operations composition app.
// ALL use cases are optional — the app degrades gracefully when not wired.
type ModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels

	// SupplierContract queries (used by renewals, variance, utilization, dashboard)
	// nil-safe: views render empty state when not provided
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)

	// ProcurementRequest queries (used by dashboard pending-approvals widget)
	// nil-safe: widget shows 0 when not provided
	ListProcurementRequests func(ctx context.Context, req *procurementrequestpb.ListProcurementRequestsRequest) (*procurementrequestpb.ListProcurementRequestsResponse, error)

	// Expenditure queries (used by recurrence-drafts queue)
	// nil-safe: recurrence queue shows empty state until P5 ships
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// Module holds all constructed Procurement Operations views.
type Module struct {
	routes           epkg.Routes
	Dashboard        view.View
	Renewals         view.View
	Variance         view.View
	Utilization      view.View
	RecurrenceDrafts view.View
}

// NewModule constructs the Procurement Operations composition module.
// All view constructors are nil-safe — missing use cases result in empty-state renders.
func NewModule(deps *ModuleDeps) *Module {
	routes := deps.Routes
	if routes.DashboardURL == "" {
		routes = epkg.DefaultRoutes()
	}

	dashDeps := &procurementdashboard.Deps{
		Routes:                  routes,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		ListSupplierContracts:   deps.ListSupplierContracts,
		ListProcurementRequests: deps.ListProcurementRequests,
	}

	renewalsDeps := &procurementrenewals.Deps{
		Routes:                routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		ListSupplierContracts: deps.ListSupplierContracts,
	}

	varianceDeps := &procurementvariance.Deps{
		Routes:                routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		ListSupplierContracts: deps.ListSupplierContracts,
	}

	utilizationDeps := &procurementutilization.Deps{
		Routes:                routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		ListSupplierContracts: deps.ListSupplierContracts,
	}

	recurrenceDeps := &procurementrecurrence.Deps{
		Routes:           routes,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		ListExpenditures: deps.ListExpenditures,
	}

	return &Module{
		routes:           routes,
		Dashboard:        procurementdashboard.NewView(dashDeps),
		Renewals:         procurementrenewals.NewView(renewalsDeps),
		Variance:         procurementvariance.NewView(varianceDeps),
		Utilization:      procurementutilization.NewView(utilizationDeps),
		RecurrenceDrafts: procurementrecurrence.NewView(recurrenceDeps),
	}
}

// RegisterRoutes registers all Procurement Operations GET routes.
// P3c (service-admin composition) calls this after mounting the module.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.RenewalCalendarURL, m.Renewals)
	r.GET(m.routes.VarianceURL, m.Variance)
	r.GET(m.routes.UtilizationURL, m.Utilization)
	r.GET(m.routes.RecurrenceDraftsURL, m.RecurrenceDrafts)
}
