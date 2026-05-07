package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_plan/form"
	"github.com/erniealice/pyeza-golang/view"

	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// NewAddAction creates the supplier_plan add action.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-plan-drawer-form", &form.Data{
				FormAction:        deps.Routes.AddURL,
				Active:            true,
				SearchSupplierURL: deps.SearchSupplierURL,
				Labels:            buildFormLabels(deps.Labels),
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		name := r.FormValue("name")
		supplierID := r.FormValue("supplier_id")
		active := r.FormValue("active") != "false"
		req := &supplierplanpb.CreateSupplierPlanRequest{
			Data: &supplierplanpb.SupplierPlan{
				Name:       name,
				SupplierId: supplierID,
				Active:     active,
			},
		}
		if _, err := deps.CreateSupplierPlan(ctx, req); err != nil {
			log.Printf("Failed to create supplier plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-plans-table")
	})
}
