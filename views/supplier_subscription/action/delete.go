package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// NewDeleteAction creates the supplier_subscription single-record delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		if _, err := deps.DeleteSupplierSubscription(ctx, &suppliersubscriptionpb.DeleteSupplierSubscriptionRequest{
			Data: &suppliersubscriptionpb.SupplierSubscription{Id: id},
		}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}

// NewBulkDeleteAction creates the supplier_subscription bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeleteSupplierSubscription(ctx, &suppliersubscriptionpb.DeleteSupplierSubscriptionRequest{
					Data: &suppliersubscriptionpb.SupplierSubscription{Id: id},
				})
			}
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}

// NewSetStatusAction creates the supplier_subscription activate/deactivate action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		active := status == "active"
		if deps.SetSupplierSubscriptionActive != nil {
			if err := deps.SetSupplierSubscriptionActive(ctx, id, active); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}

// NewBulkSetStatusAction creates the supplier_subscription bulk activate/deactivate action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		active := status == "active"
		if deps.SetSupplierSubscriptionActive != nil {
			for _, id := range ids {
				if id != "" {
					_ = deps.SetSupplierSubscriptionActive(ctx, id, active)
				}
			}
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}
