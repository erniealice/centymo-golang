package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_product_plan/form"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// NewAddAction creates the supplier_product_plan add action.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-product-plan-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				Active:                true,
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
		supplierUnit := r.FormValue("supplier_unit")
		active := r.FormValue("active") != "false"

		spp := &supplierproductplanpb.SupplierProductPlan{
			SupplierPlanId: supplierPlanID,
			ProductId:      productID,
			Active:         active,
		}
		if productVariantID != "" {
			spp.ProductVariantId = strPtr(productVariantID)
		}
		// Name is auto-derived server-side; send SKU as the name stub if supplied.
		if supplierSKU != "" {
			spp.Name = supplierSKU
		}
		_ = supplierUnit // stored via use-case if the proto adds it later

		if _, err := deps.CreateSupplierProductPlan(ctx, &supplierproductplanpb.CreateSupplierProductPlanRequest{Data: spp}); err != nil {
			log.Printf("Failed to create supplier product plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}
