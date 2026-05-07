package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_plan/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// NewEditAction creates the supplier_plan edit action.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			var record *supplierplanpb.SupplierPlan
			if deps.GetSupplierPlanItemPageData != nil {
				resp, err := deps.GetSupplierPlanItemPageData(ctx, &supplierplanpb.GetSupplierPlanItemPageDataRequest{
					SupplierPlanId: id,
				})
				if err != nil || resp == nil || resp.GetSupplierPlan() == nil {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSupplierPlan()
			} else {
				resp, err := deps.ReadSupplierPlan(ctx, &supplierplanpb.ReadSupplierPlanRequest{
					Data: &supplierplanpb.SupplierPlan{Id: id},
				})
				if err != nil || len(resp.GetData()) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetData()[0]
			}
			return view.OK("supplier-plan-drawer-form", &form.Data{
				FormAction:        route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:            true,
				ID:                id,
				Name:              record.GetName(),
				SupplierID:        record.GetSupplierId(),
				SupplierLabel:     record.GetSupplierId(),
				Active:            record.GetActive(),
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
		req := &supplierplanpb.UpdateSupplierPlanRequest{
			Data: &supplierplanpb.SupplierPlan{
				Id:         id,
				Name:       name,
				SupplierId: supplierID,
				Active:     active,
			},
		}
		if _, err := deps.UpdateSupplierPlan(ctx, req); err != nil {
			log.Printf("Failed to update supplier plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-plans-table")
	})
}
