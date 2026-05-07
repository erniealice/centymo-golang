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

// NewEditAction creates the supplier_product_cost_plan edit action (inline within CostPlan detail).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		costPlanID := viewCtx.Request.PathValue("id")
		pcid := viewCtx.Request.PathValue("pcid")
		editURL := route.ResolveURL(
			route.ResolveURL(deps.CostPlanRoutes.ProductCostEditURL, "id", costPlanID),
			"pcid", pcid,
		)
		if viewCtx.Request.Method == http.MethodGet {
			var record *supplierproductcostplanpb.SupplierProductCostPlan
			if deps.GetSupplierProductCostPlanItemPageData != nil {
				resp, err := deps.GetSupplierProductCostPlanItemPageData(ctx, &supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest{
					SupplierProductCostPlanId: pcid,
				})
				if err != nil || resp == nil || resp.GetSupplierProductCostPlan() == nil {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSupplierProductCostPlan()
			} else {
				resp, err := deps.ReadSupplierProductCostPlan(ctx, &supplierproductcostplanpb.ReadSupplierProductCostPlanRequest{
					Data: &supplierproductcostplanpb.SupplierProductCostPlan{Id: pcid},
				})
				if err != nil || len(resp.GetData()) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetData()[0]
			}

			sppLabel := record.GetSupplierProductPlanId()
			if spp := record.GetSupplierProductPlan(); spp != nil && spp.GetName() != "" {
				sppLabel = spp.GetName()
			}

			return view.OK("supplier-product-cost-plan-drawer-form", &form.Data{
				FormAction:                   editURL,
				IsEdit:                       true,
				ID:                           pcid,
				CostPlanID:                   costPlanID,
				SupplierProductPlanID:        record.GetSupplierProductPlanId(),
				SupplierProductPlanLabel:     sppLabel,
				BillingTreatment:             record.GetBillingTreatment().String(),
				Amount:                       formatAmount(record.GetBillingAmount()),
				Active:                       record.GetActive(),
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
			Id:                    pcid,
			CostPlanId:            costPlanID,
			SupplierProductPlanId: supplierProductPlanID,
			BillingAmount:         amount,
			Active:                active,
		}
		if billingTreatment != "" {
			if bt, ok := supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment_value[billingTreatment]; ok {
				spcp.BillingTreatment = supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment(bt)
			}
		}

		if _, err := deps.UpdateSupplierProductCostPlan(ctx, &supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest{Data: spcp}); err != nil {
			log.Printf("Failed to update supplier product cost plan %s: %v", pcid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-plan-lines-table")
	})
}
