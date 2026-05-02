package action

// ad_hoc_wrapper.go provides backward-compatible shim constructors that keep
// block.go's subscriptionaction.NewRequestUsageAction call site unchanged
// while the implementation now lives in the ad_hoc_actions/ sub-package.

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	adhocpkg "github.com/erniealice/centymo-golang/views/subscription/ad_hoc_actions"
)

// NewRequestUsageAction is the backward-compatible shim for block.go.
func NewRequestUsageAction(deps *Deps) view.View {
	var adapter adhocpkg.MaterializeInstanceJobsForSubscriptionAdapter
	if deps.MaterializeInstanceJobsForSubscription != nil {
		adapter = func(ctx context.Context, req *adhocpkg.MaterializeInstanceJobsRequest) (*adhocpkg.MaterializeInstanceJobsResponse, error) {
			resp, err := deps.MaterializeInstanceJobsForSubscription(ctx, &MaterializeInstanceJobsRequest{
				SubscriptionID:   req.SubscriptionID,
				CyclePeriodStart: req.CyclePeriodStart,
				Backfill:         req.Backfill,
				UsageRequestDate: req.UsageRequestDate,
			})
			if err != nil || resp == nil {
				return nil, err
			}
			return &adhocpkg.MaterializeInstanceJobsResponse{
				SpawnedCycleCount:         resp.SpawnedCycleCount,
				SpawnedJobCount:           resp.SpawnedJobCount,
				OnceAtStartJobCount:       resp.OnceAtStartJobCount,
				EngagementWasNewlyCreated: resp.EngagementWasNewlyCreated,
				SkippedReason:             resp.SkippedReason,
				BackfillCappedAt:          resp.BackfillCappedAt,
			}, nil
		}
	}
	return adhocpkg.NewRequestUsageAction(&adhocpkg.Deps{
		Labels:                                 deps.Labels,
		MaterializeInstanceJobsForSubscription: adapter,
	})
}
