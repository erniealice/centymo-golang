package action

import (
	"context"
	"strings"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	"github.com/erniealice/pyeza-golang/view"

	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
)

// SetStatusDeps holds dependencies for the S6 exclude/restore action — the
// ACTIVE<->EXCLUDED status-enum transition (distinct from the entity's own
// active/inactive soft-delete axis; centymo.md §2 SetStatusURL).
type SetStatusDeps struct {
	Routes                             sgpp.Routes
	Labels                             sgpp.Labels
	ReadSubscriptionGroupProductPlan   func(ctx context.Context, req *sgpppb.ReadSubscriptionGroupProductPlanRequest) (*sgpppb.ReadSubscriptionGroupProductPlanResponse, error)
	UpdateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.UpdateSubscriptionGroupProductPlanRequest) (*sgpppb.UpdateSubscriptionGroupProductPlanResponse, error)
}

// NewSetStatusAction handles POST /action/subscription-group-product-plan/
// set-status/{sgppId}?status=excluded|active. The ACTIVE->EXCLUDED transition
// is guarded server-side (espyna's UpdateSubscriptionGroupProductPlanUseCase,
// plan.md §2: refuse while active assignments or a live realized course
// reference the class) — this action surfaces that guard's error verbatim,
// fail-closed on the read/permission gates.
func NewSetStatusAction(deps *SetStatusDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		sgppID := viewCtx.Request.PathValue("sgppId")
		if sgppID == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		status := viewCtx.Request.URL.Query().Get("status")
		if status == "" {
			_ = viewCtx.Request.ParseForm()
			status = viewCtx.Request.FormValue("status")
		}
		status = strings.TrimSpace(status)

		readResp, err := deps.ReadSubscriptionGroupProductPlan(ctx, &sgpppb.ReadSubscriptionGroupProductPlanRequest{
			Data: &sgpppb.SubscriptionGroupProductPlan{Id: sgppID},
		})
		if err != nil || readResp == nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]

		newStatus := sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_ACTIVE
		if status == "excluded" {
			newStatus = sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED
		}

		_, err = deps.UpdateSubscriptionGroupProductPlan(ctx, &sgpppb.UpdateSubscriptionGroupProductPlanRequest{
			Data: &sgpppb.SubscriptionGroupProductPlan{
				Id:                  sgppID,
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				ProductPlanId:       record.GetProductPlanId(),
				JobTemplateId:       record.GetJobTemplateId(),
				Status:              newStatus,
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sgpp.InfoTabPanelID)
	})
}
