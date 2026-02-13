package action

import (
	"context"
	"log"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// NewSetStatusAction creates the inventory activate/deactivate action (POST only).
func NewSetStatusAction(setActive func(ctx context.Context, id string, active bool) error) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError("Inventory item ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError("Invalid status")
		}

		if err := setActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update inventory status %s: %v", id, err)
			return centymo.HTMXError("Failed to update inventory status")
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewBulkSetStatusAction creates the inventory bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(setActive func(ctx context.Context, id string, active bool) error) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError("No inventory item IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := setActive(ctx, id, active); err != nil {
				log.Printf("Failed to update inventory status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}
