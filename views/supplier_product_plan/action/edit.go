package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_product_plan/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// NewEditAction creates the supplier_product_plan edit action.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			var record *supplierproductplanpb.SupplierProductPlan
			if deps.GetSupplierProductPlanItemPageData != nil {
				resp, err := deps.GetSupplierProductPlanItemPageData(ctx, &supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest{
					SupplierProductPlanId: id,
				})
				if err != nil || resp == nil || resp.GetSupplierProductPlan() == nil {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSupplierProductPlan()
			} else {
				resp, err := deps.ReadSupplierProductPlan(ctx, &supplierproductplanpb.ReadSupplierProductPlanRequest{
					Data: &supplierproductplanpb.SupplierProductPlan{Id: id},
				})
				if err != nil || len(resp.GetData()) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetData()[0]
			}

			productLabel := record.GetProductId()
			if p := record.GetProduct(); p != nil && p.GetName() != "" {
				productLabel = p.GetName()
			}

			return view.OK("supplier-product-plan-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				SupplierPlanID:        record.GetSupplierPlanId(),
				SupplierPlanLabel:     record.GetSupplierPlanId(),
				ProductID:             record.GetProductId(),
				ProductLabel:          productLabel,
				ProductVariantID:      record.GetProductVariantId(),
				SupplierSKU:           record.GetName(),
				Active:                record.GetActive(),
				SearchSupplierPlanURL: deps.SearchSupplierPlanURL,
				SearchProductURL:      deps.SearchProductURL,
				Labels:                buildFormLabels(deps.Labels),
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		supplierPlanID := r.FormValue("supplier_plan_id")
		productID := r.FormValue("product_id")
		productVariantID := r.FormValue("product_variant_id")
		supplierSKU := r.FormValue("supplier_sku")
		active := r.FormValue("active") != "false"

		spp := &supplierproductplanpb.SupplierProductPlan{
			Id:             id,
			SupplierPlanId: supplierPlanID,
			ProductId:      productID,
			Active:         active,
		}
		if supplierSKU != "" {
			spp.Name = supplierSKU
		}
		if productVariantID != "" {
			spp.ProductVariantId = strPtr(productVariantID)
		}

		if _, err := deps.UpdateSupplierProductPlan(ctx, &supplierproductplanpb.UpdateSupplierProductPlanRequest{Data: spp}); err != nil {
			log.Printf("Failed to update supplier product plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}
