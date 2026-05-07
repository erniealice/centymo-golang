package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// NewDeleteAction creates the supplier_product_plan delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "delete") {
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
		if _, err := deps.DeleteSupplierProductPlan(ctx, &supplierproductplanpb.DeleteSupplierProductPlanRequest{
			Data: &supplierproductplanpb.SupplierProductPlan{Id: id},
		}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}

// NewBulkDeleteAction creates the supplier_product_plan bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeleteSupplierProductPlan(ctx, &supplierproductplanpb.DeleteSupplierProductPlanRequest{
					Data: &supplierproductplanpb.SupplierProductPlan{Id: id},
				})
			}
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}

// NewSetStatusAction creates the supplier_product_plan activate/deactivate action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		if deps.SetSupplierProductPlanActive != nil {
			if err := deps.SetSupplierProductPlanActive(ctx, id, status == "active"); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}

// NewBulkSetStatusAction creates the supplier_product_plan bulk activate/deactivate action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		active := viewCtx.Request.FormValue("target_status") == "active"
		if deps.SetSupplierProductPlanActive != nil {
			for _, id := range ids {
				if id != "" {
					_ = deps.SetSupplierProductPlanActive(ctx, id, active)
				}
			}
		}
		return centymo.HTMXSuccess("supplier-product-plans-table")
	})
}
