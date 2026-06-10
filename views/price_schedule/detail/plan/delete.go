package plan

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// NewDeleteAction handles POST /action/price-schedule/{id}/plan/{ppid}/delete.
// Hard delete — PricePlan rows are removed permanently (matches price_schedule's delete semantics).
func NewDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		ppid := viewCtx.Request.PathValue("ppid")
		if ppid == "" {
			_ = viewCtx.Request.ParseForm()
			ppid = viewCtx.Request.FormValue("id")
		}
		if ppid == "" {
			return view.HTMXError(deps.PlanLabels.Errors.NotFound)
		}
		if _, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("price-schedule-plans-table")
	})
}
