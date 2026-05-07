package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// NewDeleteAction creates the supplier_product_cost_plan delete action (inline within CostPlan detail).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		if _, err := deps.DeleteSupplierProductCostPlan(ctx, &supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest{
			Data: &supplierproductcostplanpb.SupplierProductCostPlan{Id: id},
		}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-plan-lines-table")
	})
}
