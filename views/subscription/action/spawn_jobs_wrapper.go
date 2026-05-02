package action

// spawn_jobs_wrapper.go provides backward-compatible shim constructors that
// keep block.go's subscriptionaction.NewSpawnJobsAction call site unchanged
// while the implementation now lives in the spawn_jobs/ sub-package.

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	spawnjobspkg "github.com/erniealice/centymo-golang/views/subscription/spawn_jobs"
)

// NewSpawnJobsAction is the backward-compatible shim for block.go.
// Delegates to spawn_jobs.NewAction using a sub-set of action.Deps.
func NewSpawnJobsAction(deps *Deps) view.View {
	return spawnjobspkg.NewAction(&spawnjobspkg.Deps{
		Routes:                         deps.Routes,
		Labels:                         deps.Labels,
		ReadSubscription:               deps.ReadSubscription,
		MaterializeJobsForSubscription: deps.MaterializeJobsForSubscription,
		// Wire the detection callback by closing over action.Deps so the full
		// dep graph (ReadPricePlan, ReadPlan, ReadJobTemplate, etc.) is available.
		DetectSpawnJobs: func(ctx context.Context, pricePlanID string) spawnjobspkg.DetectionResult {
			det := detectSpawnJobs(ctx, deps, pricePlanID)
			return spawnjobspkg.DetectionResult{
				Available:     det.Available,
				TemplateNames: det.TemplateNames,
				JobCount:      det.JobCount,
				PhaseCount:    det.PhaseCount,
				TaskCount:     det.TaskCount,
			}
		},
	})
}
