package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_product_cost_plan/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// NewAddAction creates the supplier_product_cost_plan add action (inline within CostPlan detail).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		costPlanID := viewCtx.Request.PathValue("id")
		addURL := route.ResolveURL(deps.CostPlanRoutes.ProductCostAddURL, "id", costPlanID)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-product-cost-plan-drawer-form", &form.Data{
				FormAction:                   addURL,
				CostPlanID:                   costPlanID,
				Active:                       true,
				SearchSupplierProductPlanURL: deps.SearchSupplierProductPlanURL,
				Labels:                       buildFormLabels(deps.Labels),
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		supplierProductPlanID := r.FormValue("supplier_product_plan_id")
		billingTreatment := r.FormValue("billing_treatment")
		amount := parseAmount(r.FormValue("amount"))
		active := r.FormValue("active") != "false"

		spcp := &supplierproductcostplanpb.SupplierProductCostPlan{
			CostPlanId:              costPlanID,
			SupplierProductPlanId:   supplierProductPlanID,
			BillingAmount:           amount,
			Active:                  active,
		}
		if billingTreatment != "" {
			if bt, ok := supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment_value[billingTreatment]; ok {
				spcp.BillingTreatment = supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment(bt)
			}
		}

		if _, err := deps.CreateSupplierProductCostPlan(ctx, &supplierproductcostplanpb.CreateSupplierProductCostPlanRequest{Data: spcp}); err != nil {
			log.Printf("Failed to create supplier product cost plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-plan-lines-table")
	})
}
