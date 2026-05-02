package action

import (
	"context"
	"log"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// NewSetStatusAction creates the subscription activate/deactivate action (POST only).
// Expects query params: ?id={subscriptionId}&status={active|inactive}
//
// Uses SetSubscriptionActive (raw map update) instead of protobuf because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if deps.SetSubscriptionActive == nil {
			return centymo.HTMXError("set-status not configured")
		}

		if err := deps.SetSubscriptionActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update subscription status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// NewBulkSetStatusAction creates the subscription bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if deps.SetSubscriptionActive == nil {
			return centymo.HTMXError("set-status not configured")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetSubscriptionActive(ctx, id, active); err != nil {
				log.Printf("Failed to update subscription status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

