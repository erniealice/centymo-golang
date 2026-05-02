package action

// spawn_cycle_jobs_wrapper.go provides backward-compatible shim constructors
// that keep block.go's subscriptionaction.NewSpawnCycleJobsAction and
// NewBackfillCyclesAction call sites unchanged while the implementation now
// lives in the spawn_cycle_jobs/ sub-package.

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	spawncyclepkg "github.com/erniealice/centymo-golang/views/subscription/spawn_cycle_jobs"
)

// NewSpawnCycleJobsAction is the backward-compatible shim for block.go.
func NewSpawnCycleJobsAction(deps *Deps) view.View {
	return spawncyclepkg.NewSpawnCycleJobsAction(adaptSpawnCycleDeps(deps))
}

// NewBackfillCyclesAction is the backward-compatible shim for block.go.
func NewBackfillCyclesAction(deps *Deps) view.View {
	return spawncyclepkg.NewBackfillCyclesAction(adaptSpawnCycleDeps(deps))
}

// adaptSpawnCycleDeps builds spawn_cycle_jobs.Deps from action.Deps.
// The MaterializeInstanceJobsForSubscription adapter is bridged by converting
// between the action-package request/response types and the sub-package types.
func adaptSpawnCycleDeps(deps *Deps) *spawncyclepkg.Deps {
	var adapter spawncyclepkg.MaterializeInstanceJobsForSubscriptionAdapter
	if deps.MaterializeInstanceJobsForSubscription != nil {
		adapter = func(ctx context.Context, req *spawncyclepkg.MaterializeInstanceJobsRequest) (*spawncyclepkg.MaterializeInstanceJobsResponse, error) {
			resp, err := deps.MaterializeInstanceJobsForSubscription(ctx, &MaterializeInstanceJobsRequest{
				SubscriptionID:   req.SubscriptionID,
				CyclePeriodStart: req.CyclePeriodStart,
				Backfill:         req.Backfill,
				UsageRequestDate: req.UsageRequestDate,
			})
			if err != nil || resp == nil {
				return nil, err
			}
			return &spawncyclepkg.MaterializeInstanceJobsResponse{
				SpawnedCycleCount:         resp.SpawnedCycleCount,
				SpawnedJobCount:           resp.SpawnedJobCount,
				OnceAtStartJobCount:       resp.OnceAtStartJobCount,
				EngagementWasNewlyCreated: resp.EngagementWasNewlyCreated,
				SkippedReason:             resp.SkippedReason,
				BackfillCappedAt:          resp.BackfillCappedAt,
			}, nil
		}
	}
	return &spawncyclepkg.Deps{
		Routes: deps.Routes,
		Labels: deps.Labels,
		ResolveSubscriptionLabel: func(ctx context.Context, subscriptionID string) string {
			return resolveSubscriptionLabel(ctx, deps, subscriptionID)
		},
		MaterializeInstanceJobsForSubscription: adapter,
	}
}
