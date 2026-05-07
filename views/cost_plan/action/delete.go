package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
)

// NewDeleteAction creates the cost_plan delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_plan", "delete") {
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
		if _, err := deps.DeleteCostPlan(ctx, &costplanpb.DeleteCostPlanRequest{
			Data: &costplanpb.CostPlan{Id: id},
		}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-plans-table")
	})
}

// NewBulkDeleteAction creates the cost_plan bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeleteCostPlan(ctx, &costplanpb.DeleteCostPlanRequest{
					Data: &costplanpb.CostPlan{Id: id},
				})
			}
		}
		return centymo.HTMXSuccess("cost-plans-table")
	})
}

// NewSetStatusAction creates the cost_plan activate/deactivate action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_plan", "update") {
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
		if deps.SetCostPlanActive != nil {
			if err := deps.SetCostPlanActive(ctx, id, status == "active"); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("cost-plans-table")
	})
}

// NewBulkSetStatusAction creates the cost_plan bulk activate/deactivate action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		active := viewCtx.Request.FormValue("target_status") == "active"
		if deps.SetCostPlanActive != nil {
			for _, id := range ids {
				if id != "" {
					_ = deps.SetCostPlanActive(ctx, id, active)
				}
			}
		}
		return centymo.HTMXSuccess("cost-plans-table")
	})
}
