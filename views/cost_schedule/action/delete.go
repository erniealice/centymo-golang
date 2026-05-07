package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
)

// NewDeleteAction creates the cost_schedule delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "delete") {
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
		if _, err := deps.DeleteCostSchedule(ctx, &costschedulepb.DeleteCostScheduleRequest{
			Data: &costschedulepb.CostSchedule{Id: id},
		}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-schedules-table")
	})
}

// NewBulkDeleteAction creates the cost_schedule bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeleteCostSchedule(ctx, &costschedulepb.DeleteCostScheduleRequest{
					Data: &costschedulepb.CostSchedule{Id: id},
				})
			}
		}
		return centymo.HTMXSuccess("cost-schedules-table")
	})
}

// NewSetStatusAction creates the cost_schedule activate/deactivate action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "update") {
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
		if deps.SetCostScheduleActive != nil {
			if err := deps.SetCostScheduleActive(ctx, id, status == "active"); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("cost-schedules-table")
	})
}

// NewBulkSetStatusAction creates the cost_schedule bulk activate/deactivate action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		active := viewCtx.Request.FormValue("target_status") == "active"
		if deps.SetCostScheduleActive != nil {
			for _, id := range ids {
				if id != "" {
					_ = deps.SetCostScheduleActive(ctx, id, active)
				}
			}
		}
		return centymo.HTMXSuccess("cost-schedules-table")
	})
}
