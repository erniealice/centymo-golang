// Package spawn_cycle_jobs handles cyclic job spawning and backfill for subscriptions.
// Backfill drawer template: subscription-backfill-cycles-drawer-form.html (flat at view root).
package spawn_cycle_jobs

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	spawnform "github.com/erniealice/centymo-golang/views/subscription/spawn_cycle_jobs/form"
)

// MaterializeInstanceJobsRequest mirrors the espyna consumer-surface request.
// Keeping a centymo-local struct so this package does not import espyna directly.
// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D.
type MaterializeInstanceJobsRequest struct {
	SubscriptionID   string
	CyclePeriodStart string
	Backfill         bool
	// 2026-05-01 ad-hoc-subscription-billing plan §3.2 — operator-supplied
	// usage request date for AD_HOC plans. Empty defaults to today UTC.
	UsageRequestDate string
}

// MaterializeInstanceJobsResponse mirrors the espyna consumer-surface response.
type MaterializeInstanceJobsResponse struct {
	SpawnedCycleCount         int
	SpawnedJobCount           int
	OnceAtStartJobCount       int
	EngagementWasNewlyCreated bool
	SkippedReason             string
	BackfillCappedAt          int32
}

// MaterializeInstanceJobsForSubscriptionAdapter is the function-pointer type
// the centymo block.go wires once the espyna consumer is available.
type MaterializeInstanceJobsForSubscriptionAdapter func(
	ctx context.Context, req *MaterializeInstanceJobsRequest,
) (*MaterializeInstanceJobsResponse, error)

// Deps is the dependency subset needed by the spawn_cycle_jobs feature.
type Deps struct {
	Routes centymo.SubscriptionRoutes
	Labels centymo.SubscriptionLabels

	// ResolveSubscriptionLabel returns "code · name" for the drawer header.
	// Provided by the caller (wired through action.Deps).
	ResolveSubscriptionLabel func(ctx context.Context, subscriptionID string) string

	MaterializeInstanceJobsForSubscription MaterializeInstanceJobsForSubscriptionAdapter
}

// NewSpawnCycleJobsAction handles POST /action/subscription/spawn-cycle-jobs/{subscriptionId}.
// 2026-04-30 cyclic-subscription-jobs plan §5.3.
func NewSpawnCycleJobsAction(deps *Deps) view.View {
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
		cyclePeriodStart := strings.TrimSpace(viewCtx.Request.FormValue("cycle_period_start"))

		_, err := deps.MaterializeInstanceJobsForSubscription(ctx, &MaterializeInstanceJobsRequest{
			SubscriptionID:   subscriptionID,
			CyclePeriodStart: cyclePeriodStart,
			Backfill:         false,
		})
		if err != nil {
			log.Printf("Failed to spawn cycle jobs for subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("subscription-operations-tab")
	})
}

// NewBackfillCyclesAction handles GET (drawer) / POST (commit) at
// /action/subscription/backfill-cycle-jobs/{subscriptionId}.
// 2026-04-30 cyclic-subscription-jobs plan §5.3.
func NewBackfillCyclesAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		subscriptionID := viewCtx.Request.PathValue("subscriptionId")
		if subscriptionID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		formAction := strings.ReplaceAll(deps.Routes.BackfillCycleJobsURL, "{subscriptionId}", subscriptionID)

		if viewCtx.Request.Method == http.MethodGet {
			subLabel := ""
			if deps.ResolveSubscriptionLabel != nil {
				subLabel = deps.ResolveSubscriptionLabel(ctx, subscriptionID)
			}
			drawerLabels := spawnform.Labels{
				Title:       deps.Labels.Backfill.DrawerTitle,
				Description: deps.Labels.Backfill.DrawerDescription,
				CountLabel:  deps.Labels.Backfill.CountLabel,
				Confirm:     deps.Labels.Backfill.Confirm,
				Cancel:      deps.Labels.Backfill.Cancel,
				MaxWarning:  deps.Labels.Backfill.MaxWarning,
			}
			return view.OK("subscription-backfill-cycles-drawer-form", &spawnform.Data{
				FormAction:        formAction,
				SubscriptionID:    subscriptionID,
				SubscriptionLabel: subLabel,
				MaxCycles:         24,
				DefaultCycles:     1,
				Labels:            drawerLabels,
				CommonLabels:      nil,
			})
		}

		// POST — invoke the adapter.
		if deps.MaterializeInstanceJobsForSubscription == nil {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_ = viewCtx.Request.ParseForm()
		countStr := strings.TrimSpace(viewCtx.Request.FormValue("count"))
		if countStr == "" {
			countStr = "1"
		}
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		if count > 24 {
			count = 24
		}
		_ = count // espyna caps; centymo just sends backfill=true

		_, err = deps.MaterializeInstanceJobsForSubscription(ctx, &MaterializeInstanceJobsRequest{
			SubscriptionID: subscriptionID,
			Backfill:       true,
		})
		if err != nil {
			log.Printf("Failed to backfill cycle jobs for subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("subscription-operations-tab")
	})
}
