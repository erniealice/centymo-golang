package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	subscription_group_member "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
)

// Deps holds all dependencies for subscription_group_member action views.
type Deps struct {
	Routes                             subscription_group_member.Routes
	Labels                             subscription_group_member.Labels
	CreateSubscriptionGroupMember      func(ctx context.Context, req *subscriptiongroupmemberpb.CreateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.CreateSubscriptionGroupMemberResponse, error)
	ReadSubscriptionGroupMember        func(ctx context.Context, req *subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.ReadSubscriptionGroupMemberResponse, error)
	UpdateSubscriptionGroupMember      func(ctx context.Context, req *subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberResponse, error)
	DeleteSubscriptionGroupMember      func(ctx context.Context, req *subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberResponse, error)
	GetSubscriptionGroupMemberInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData writes the POST body onto a SubscriptionGroupMember.
// Optional FK fields (subscription_group_id, subscription_id, client_id) are
// set only when non-empty so that unset values remain null on the server.
func applyFormToData(r *http.Request) *subscriptiongroupmemberpb.SubscriptionGroupMember {
	data := &subscriptiongroupmemberpb.SubscriptionGroupMember{
		Active: r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("subscription_group_id")); v != "" {
		data.SubscriptionGroupId = v
	}
	if v := strings.TrimSpace(r.FormValue("subscription_id")); v != "" {
		data.SubscriptionId = v
	}
	if v := strings.TrimSpace(r.FormValue("client_id")); v != "" {
		data.ClientId = v
	}
	return data
}

// NewAddAction creates the subscription_group_member add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("subscription-group-member-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &subscriptiongroupmemberpb.CreateSubscriptionGroupMemberRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateSubscriptionGroupMember(ctx, req); err != nil {
			log.Printf("Failed to create subscription group member: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}

// NewEditAction creates the subscription_group_member edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("subscription_group_member", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest{
				Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
			})
			if err != nil || len(resp.GetData()) == 0 {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]

			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				formAction = deps.Routes.AddURL
				formID = ""
			}

			return view.OK("subscription-group-member-drawer-form", &form.Data{
				FormAction:          formAction,
				IsEdit:              !isClone,
				ID:                  formID,
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				SubscriptionId:      record.GetSubscriptionId(),
				ClientId:            record.GetClientId(),
				Active:              record.GetActive(),
				Labels:              deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}

// NewDeleteAction creates the subscription_group_member delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "delete") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if deps.GetSubscriptionGroupMemberInUseIDs != nil {
			if inUse, _ := deps.GetSubscriptionGroupMemberInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberRequest{
			Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
		}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}

// NewBulkDeleteAction creates the subscription_group_member bulk-delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "delete") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		ids := viewCtx.Request.Form["id"]
		var attempted []string
		for _, id := range ids {
			if id != "" {
				attempted = append(attempted, id)
			}
		}
		if len(attempted) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		var inUse map[string]bool
		if deps.GetSubscriptionGroupMemberInUseIDs != nil {
			inUse, _ = deps.GetSubscriptionGroupMemberInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberRequest{
				Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
			}); err != nil {
				log.Printf("Failed to delete subscription group member %s during bulk: %v", id, err)
				failed++
				continue
			}
			deleted++
		}
		if deleted == 0 {
			if blocked > 0 && failed == 0 {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
			return view.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}

// NewSetStatusAction creates the subscription_group_member single set-status action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest{
			Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
		})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest{
			Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{
				Id:                  id,
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				SubscriptionId:      record.GetSubscriptionId(),
				ClientId:            record.GetClientId(),
				Active:              status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}

// NewBulkSetStatusAction creates the subscription_group_member bulk set-status action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest{
				Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest{
				Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{
					Id:                  id,
					SubscriptionGroupId: record.GetSubscriptionGroupId(),
					SubscriptionId:      record.GetSubscriptionId(),
					ClientId:            record.GetClientId(),
					Active:              status == "active",
				},
			})
		}
		return view.HTMXSuccess("subscription-group-members-table")
	})
}
