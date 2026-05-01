package action

import (
	"context"
	"log"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — operator-driven CTA for
// AD_HOC subscriptions: Request Usage. Pool-Generate-Invoice reuses the existing
// /action/subscription/recognize-revenue/{id} endpoint (espyna's executeCore
// dispatches to executeAdHoc → executeAdHocPool automatically when the
// PricePlan kind is AD_HOC × TOTAL_PACKAGE; no new handler needed).
//
// Extend-Pool (writing Subscription.entitled_occurrences_override) is tracked
// as a v1.5.5 follow-up — needs a fresh espyna use case + consumer surface +
// adapter wiring; meaningfully more code than the other two CTAs.

// NewRequestUsageAction handles POST /action/subscription/{subscriptionId}/request-usage.
//
// Body fields:
//
//	usage_request_date  ISO 8601 YYYY-MM-DD (defaults to today UTC server-side)
//
// Reuses the existing MaterializeInstanceJobsForSubscription adapter wired in
// block.go for the cyclic spawn-cycle-jobs handler — espyna routes AD_HOC
// PricePlans through executeAdHoc inside the use case, so the centymo surface
// stays a thin RPC.
//
// On success, fires `subscription-operations-tab` HX-target so the Operations
// tab re-renders inline with the new usage row.
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

// _ keeps gofmt happy when context is the only import we need.
var _ = context.Background
