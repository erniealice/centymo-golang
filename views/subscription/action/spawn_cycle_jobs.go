package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// MaterializeInstanceJobsRequest mirrors the espyna consumer-surface request
// (consumer.MaterializeInstanceJobsForSubscriptionRequest). Keeping a centymo-
// local struct so this package does not import espyna directly — the centymo
// block.go wiring builds an adapter func that translates between the two.
//
// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D.
type MaterializeInstanceJobsRequest struct {
	SubscriptionID   string
	CyclePeriodStart string
	Backfill         bool
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
// the centymo block.go wires once the espyna consumer is available. Adapter
// translates centymo's MaterializeInstanceJobsRequest → espyna's consumer
// request shape (cross-module; centymo cannot import espyna directly).
//
// nil-safe: the action handler reports a permission-style error when this is
// unwired, mirroring the SpawnJobs handler's posture.
type MaterializeInstanceJobsForSubscriptionAdapter func(
	ctx context.Context, req *MaterializeInstanceJobsRequest,
) (*MaterializeInstanceJobsResponse, error)

// BackfillDrawerData is the template shape for
// `subscription-backfill-cycles-drawer-form.html` — see plan §7.1 backfill CTA.
type BackfillDrawerData struct {
	FormAction        string
	SubscriptionID    string
	SubscriptionLabel string

	// MaxCycles caps the number input (see plan §15 risk mitigation —
	// 24 cycles per request).
	MaxCycles int
	// DefaultCycles is the prefilled value (1 = "spawn the next missing").
	DefaultCycles int

	Labels       BackfillDrawerLabels
	CommonLabels any
}

// BackfillDrawerLabels carries the typed strings consumed by the drawer
// template. Keeps the template free from optional-chain `.Labels.Backfill.*`
// nav (mirrors the SpawnDrawerLabels pattern in spawn_jobs.go).
type BackfillDrawerLabels struct {
	Title        string
	Description  string
	CountLabel   string
	Confirm      string
	Cancel       string
	MaxWarning   string
}

// NewSpawnCycleJobsAction handles POST /action/subscription/spawn-cycle-jobs/{subscriptionId}.
//
// Calls MaterializeInstanceJobsForSubscription with backfill=false and the
// optional `cycle_period_start` form field (when blank, espyna spawns the
// next un-spawned cycle from sub.date_time_start). On success, emits an
// HX-Trigger that the Operations tab listens for to refresh inline; on
// non-fatal skip (e.g. "no_pending_cycles") the same trigger fires and the
// tab simply re-renders with no new accordions.
//
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
			// Adapter unwired — report a permission-style failure so the
			// drawer toast is consistent with the existing spawn-jobs flow.
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
		// Refresh the Operations tab inline. The tab listens for
		// `refresh-subscription-operations-tab` (added in operations-tab
		// template).
		return centymo.HTMXSuccess("subscription-operations-tab")
	})
}

// NewBackfillCyclesAction handles GET (drawer) / POST (commit) at
// /action/subscription/backfill-cycle-jobs/{subscriptionId}.
//
// GET renders the drawer with a count input capped at 24 (plan §15 risk).
// POST invokes the adapter with backfill=true; the espyna use case caps at
// the same value and surfaces BackfillCappedAt when the request exceeded it.
//
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
			subLabel := resolveSubscriptionLabel(ctx, deps, subscriptionID)
			drawerLabels := BackfillDrawerLabels{
				Title:       deps.Labels.Backfill.DrawerTitle,
				Description: deps.Labels.Backfill.DrawerDescription,
				CountLabel:  deps.Labels.Backfill.CountLabel,
				Confirm:     deps.Labels.Backfill.Confirm,
				Cancel:      deps.Labels.Backfill.Cancel,
				MaxWarning:  deps.Labels.Backfill.MaxWarning,
			}
			return view.OK("subscription-backfill-cycles-drawer-form", &BackfillDrawerData{
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
		// "count" is informational on the centymo side — the espyna use case
		// caps internally at 24 and returns BackfillCappedAt when the calling
		// window exceeded it. We still validate the input so a hostile POST
		// can't trick the form into "0 cycles" (no-op pass).
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
		_ = count // espyna caps; centymo just sends backfill=true and the
		// use case figures out the cycle window itself.

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
