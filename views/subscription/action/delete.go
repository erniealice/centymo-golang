package action

import (
	"context"
	"log"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// NewDeleteAction creates the subscription delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteSubscription(ctx, &subscriptionpb.DeleteSubscriptionRequest{
			Data: &subscriptionpb.Subscription{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// NewBulkDeleteAction creates the subscription bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields.
// IDs that are in use (referenced by dependent records) are skipped silently;
// the remaining IDs are deleted.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		// Gate: skip IDs that have dependent records.
		var inUse map[string]bool
		if deps.GetInUseIDs != nil {
			var err error
			inUse, err = deps.GetInUseIDs(ctx, ids)
			if err != nil {
				log.Printf("Failed to check subscription in-use IDs: %v", err)
			}
		}

		for _, id := range ids {
			if inUse[id] {
				log.Printf("Skipping bulk delete for subscription %s — has dependent records", id)
				continue
			}
			idCopy := id
			if _, err := deps.DeleteSubscription(ctx, &subscriptionpb.DeleteSubscriptionRequest{
				Data: &subscriptionpb.Subscription{Id: idCopy},
			}); err != nil {
				log.Printf("Failed to bulk-delete subscription %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}
