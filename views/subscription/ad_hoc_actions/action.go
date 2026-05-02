// Package ad_hoc_actions handles the AD_HOC subscription billing CTA
// (Request Usage endpoint). 2026-05-01 ad-hoc-subscription-billing plan §5.2.
package ad_hoc_actions

import (
	"context"
	"log"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// MaterializeInstanceJobsRequest mirrors the centymo block's request shape.
type MaterializeInstanceJobsRequest struct {
	SubscriptionID   string
	CyclePeriodStart string
	Backfill         bool
	UsageRequestDate string
}

// MaterializeInstanceJobsResponse mirrors the centymo block's response shape.
type MaterializeInstanceJobsResponse struct {
	SpawnedCycleCount         int
	SpawnedJobCount           int
	OnceAtStartJobCount       int
	EngagementWasNewlyCreated bool
	SkippedReason             string
	BackfillCappedAt          int32
}

// MaterializeInstanceJobsForSubscriptionAdapter is the function-pointer type.
type MaterializeInstanceJobsForSubscriptionAdapter func(
	ctx context.Context, req *MaterializeInstanceJobsRequest,
) (*MaterializeInstanceJobsResponse, error)

// Deps is the dependency subset for ad_hoc_actions.
type Deps struct {
	Labels                                 centymo.SubscriptionLabels
	MaterializeInstanceJobsForSubscription MaterializeInstanceJobsForSubscriptionAdapter
}

// NewRequestUsageAction handles POST /action/subscription/{subscriptionId}/request-usage.
//
// Reuses MaterializeInstanceJobsForSubscription wired by block.go — espyna routes
// AD_HOC PricePlans through executeAdHoc inside the use case.
func NewRequestUsageAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		subscriptionID := viewCtx.Request.PathValue("subscriptionId")
		if subscriptionID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if deps.MaterializeInstanceJobsForSubscription == nil {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseForm()
		usageRequestDate := strings.TrimSpace(viewCtx.Request.FormValue("usage_request_date"))

		_, err := deps.MaterializeInstanceJobsForSubscription(ctx, &MaterializeInstanceJobsRequest{
			SubscriptionID:   subscriptionID,
			UsageRequestDate: usageRequestDate,
			Backfill:         false,
		})
		if err != nil {
			log.Printf("Failed to spawn usage job for subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("subscription-operations-tab")
	})
}
