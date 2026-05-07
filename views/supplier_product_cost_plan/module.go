// Package supplier_product_cost_plan provides the inline cost-plan line editor
// mounted within the CostPlan detail page. There is no standalone list/detail;
// lines are created/edited/deleted via drawers on the CostPlan Lines tab.
package supplier_product_cost_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	spcpaction "github.com/erniealice/centymo-golang/views/supplier_product_cost_plan/action"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// ModuleDeps holds all dependencies for the supplier_product_cost_plan inline module.
type ModuleDeps struct {
	// CostPlanRoutes is the parent entity routes — used to resolve add/edit/delete URLs.
	CostPlanRoutes centymo.CostPlanRoutes
	Labels         centymo.SupplierProductCostPlanLabels
	CommonLabels   pyeza.CommonLabels

	CreateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	ReadSupplierProductCostPlan            func(ctx context.Context, req *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)
	GetSupplierProductCostPlanItemPageData func(ctx context.Context, req *supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest) (*supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataResponse, error)

	// SearchSupplierProductPlanURL for the autocomplete in the drawer.
	SearchSupplierProductPlanURL string
}

// Module holds add/edit/delete views for the inline SPCP editor.
type Module struct {
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the supplier_product_cost_plan inline module.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &spcpaction.Deps{
		CostPlanRoutes:                         deps.CostPlanRoutes,
		Labels:                                 deps.Labels,
		CommonLabels:                           deps.CommonLabels,
		CreateSupplierProductCostPlan:          deps.CreateSupplierProductCostPlan,
		ReadSupplierProductCostPlan:            deps.ReadSupplierProductCostPlan,
		UpdateSupplierProductCostPlan:          deps.UpdateSupplierProductCostPlan,
		DeleteSupplierProductCostPlan:          deps.DeleteSupplierProductCostPlan,
		GetSupplierProductCostPlanItemPageData: deps.GetSupplierProductCostPlanItemPageData,
		SearchSupplierProductPlanURL:           deps.SearchSupplierProductPlanURL,
	}
	return &Module{
		Add:    spcpaction.NewAddAction(actionDeps),
		Edit:   spcpaction.NewEditAction(actionDeps),
		Delete: spcpaction.NewDeleteAction(actionDeps),
	}
}
