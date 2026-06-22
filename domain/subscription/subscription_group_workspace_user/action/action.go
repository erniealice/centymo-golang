package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	sgwu "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	sgwupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_workspace_user"
)

// Deps holds dependencies for the subscription_group_workspace_user action views.
type Deps struct {
	Routes                                    sgwu.Routes
	Labels                                    sgwu.Labels
	CreateSubscriptionGroupWorkspaceUser      func(ctx context.Context, req *sgwupb.CreateSubscriptionGroupWorkspaceUserRequest) (*sgwupb.CreateSubscriptionGroupWorkspaceUserResponse, error)
	ReadSubscriptionGroupWorkspaceUser        func(ctx context.Context, req *sgwupb.ReadSubscriptionGroupWorkspaceUserRequest) (*sgwupb.ReadSubscriptionGroupWorkspaceUserResponse, error)
	UpdateSubscriptionGroupWorkspaceUser      func(ctx context.Context, req *sgwupb.UpdateSubscriptionGroupWorkspaceUserRequest) (*sgwupb.UpdateSubscriptionGroupWorkspaceUserResponse, error)
	DeleteSubscriptionGroupWorkspaceUser      func(ctx context.Context, req *sgwupb.DeleteSubscriptionGroupWorkspaceUserRequest) (*sgwupb.DeleteSubscriptionGroupWorkspaceUserResponse, error)
	GetSubscriptionGroupWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData maps POST body fields onto a SubscriptionGroupWorkspaceUser.
// workspace_user_id and subscription_group_id are set only when non-empty so
// unset optional FKs stay null on the wire. scope and role are free-text strings.
func applyFormToData(r *http.Request) *sgwupb.SubscriptionGroupWorkspaceUser {
	data := &sgwupb.SubscriptionGroupWorkspaceUser{
		Scope:   strings.TrimSpace(r.FormValue("scope")),
		Role:    strings.TrimSpace(r.FormValue("role")),
		IsOwner: r.FormValue("is_owner") == "true",
		Active:  r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("workspace_user_id")); v != "" {
		data.WorkspaceUserId = v
	}
	if v := strings.TrimSpace(r.FormValue("subscription_group_id")); v != "" {
		data.SubscriptionGroupId = v
	}
	return data
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("subscription-group-workspace-user-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &sgwupb.CreateSubscriptionGroupWorkspaceUserRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateSubscriptionGroupWorkspaceUser(ctx, req); err != nil {
			log.Printf("Failed to create subscription_group_workspace_user: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}

func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("subscription_group_workspace_user", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSubscriptionGroupWorkspaceUser(ctx, &sgwupb.ReadSubscriptionGroupWorkspaceUserRequest{
				Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
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

			return view.OK("subscription-group-workspace-user-drawer-form", &form.Data{
				FormAction:          formAction,
				IsEdit:              !isClone,
				ID:                  formID,
				WorkspaceUserId:     record.GetWorkspaceUserId(),
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				Scope:               record.GetScope(),
				Role:                record.GetRole(),
				IsOwner:             record.GetIsOwner(),
				Active:              record.GetActive(),
				Labels:              deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateSubscriptionGroupWorkspaceUser(ctx, &sgwupb.UpdateSubscriptionGroupWorkspaceUserRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "delete") {
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
		if deps.GetSubscriptionGroupWorkspaceUserInUseIDs != nil {
			if inUse, _ := deps.GetSubscriptionGroupWorkspaceUserInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteSubscriptionGroupWorkspaceUser(ctx, &sgwupb.DeleteSubscriptionGroupWorkspaceUserRequest{
			Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
		}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "delete") {
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
		if deps.GetSubscriptionGroupWorkspaceUserInUseIDs != nil {
			inUse, _ = deps.GetSubscriptionGroupWorkspaceUserInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteSubscriptionGroupWorkspaceUser(ctx, &sgwupb.DeleteSubscriptionGroupWorkspaceUserRequest{
				Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
			}); err != nil {
				log.Printf("Failed to delete subscription_group_workspace_user %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadSubscriptionGroupWorkspaceUser(ctx, &sgwupb.ReadSubscriptionGroupWorkspaceUserRequest{
			Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
		})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateSubscriptionGroupWorkspaceUser(ctx, &sgwupb.UpdateSubscriptionGroupWorkspaceUserRequest{
			Data: &sgwupb.SubscriptionGroupWorkspaceUser{
				Id:                  id,
				WorkspaceUserId:     record.GetWorkspaceUserId(),
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				Scope:               record.GetScope(),
				Role:                record.GetRole(),
				IsOwner:             record.GetIsOwner(),
				Active:              status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadSubscriptionGroupWorkspaceUser(ctx, &sgwupb.ReadSubscriptionGroupWorkspaceUserRequest{
				Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateSubscriptionGroupWorkspaceUser(ctx, &sgwupb.UpdateSubscriptionGroupWorkspaceUserRequest{
				Data: &sgwupb.SubscriptionGroupWorkspaceUser{
					Id:                  id,
					WorkspaceUserId:     record.GetWorkspaceUserId(),
					SubscriptionGroupId: record.GetSubscriptionGroupId(),
					Scope:               record.GetScope(),
					Role:                record.GetRole(),
					IsOwner:             record.GetIsOwner(),
					Active:              status == "active",
				},
			})
		}
		return view.HTMXSuccess("subscription-group-workspace-users-table")
	})
}
