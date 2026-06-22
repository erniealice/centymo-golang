package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	line_workspace_user "github.com/erniealice/centymo-golang/domain/product/line_workspace_user"
	"github.com/erniealice/centymo-golang/domain/product/line_workspace_user/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	lineworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line_workspace_user"
)

// Deps holds the action-layer dependencies for line_workspace_user.
type Deps struct {
	Routes                       line_workspace_user.Routes
	Labels                       line_workspace_user.Labels
	CreateLineWorkspaceUser      func(ctx context.Context, req *lineworkspaceuserpb.CreateLineWorkspaceUserRequest) (*lineworkspaceuserpb.CreateLineWorkspaceUserResponse, error)
	ReadLineWorkspaceUser        func(ctx context.Context, req *lineworkspaceuserpb.ReadLineWorkspaceUserRequest) (*lineworkspaceuserpb.ReadLineWorkspaceUserResponse, error)
	UpdateLineWorkspaceUser      func(ctx context.Context, req *lineworkspaceuserpb.UpdateLineWorkspaceUserRequest) (*lineworkspaceuserpb.UpdateLineWorkspaceUserResponse, error)
	DeleteLineWorkspaceUser      func(ctx context.Context, req *lineworkspaceuserpb.DeleteLineWorkspaceUserRequest) (*lineworkspaceuserpb.DeleteLineWorkspaceUserResponse, error)
	GetLineWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData writes the POST body onto a LineWorkspaceUser. Shared by
// Add (no id) and Edit (id set by caller). Optional FK strings are set only
// when present so unset fields stay null.
func applyFormToData(r *http.Request) *lineworkspaceuserpb.LineWorkspaceUser {
	data := &lineworkspaceuserpb.LineWorkspaceUser{
		Scope:   strings.TrimSpace(r.FormValue("scope")),
		Role:    strings.TrimSpace(r.FormValue("role")),
		IsOwner: r.FormValue("is_owner") == "true",
		Active:  r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("workspace_user_id")); v != "" {
		data.WorkspaceUserId = v
	}
	if v := strings.TrimSpace(r.FormValue("line_id")); v != "" {
		data.LineId = v
	}
	return data
}

// NewAddAction creates the line_workspace_user add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("line-workspace-user-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &lineworkspaceuserpb.CreateLineWorkspaceUserRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateLineWorkspaceUser(ctx, req); err != nil {
			log.Printf("Failed to create line_workspace_user: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("line-workspace-users-table")
	})
}

// NewEditAction creates the line_workspace_user edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("line_workspace_user", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadLineWorkspaceUser(ctx, &lineworkspaceuserpb.ReadLineWorkspaceUserRequest{
				Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
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

			return view.OK("line-workspace-user-drawer-form", &form.Data{
				FormAction:      formAction,
				IsEdit:          !isClone,
				ID:              formID,
				WorkspaceUserId: record.GetWorkspaceUserId(),
				LineId:          record.GetLineId(),
				Scope:           record.GetScope(),
				Role:            record.GetRole(),
				IsOwner:         record.GetIsOwner(),
				Active:          record.GetActive(),
				Labels:          deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateLineWorkspaceUser(ctx, &lineworkspaceuserpb.UpdateLineWorkspaceUserRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("line-workspace-users-table")
	})
}

// NewDeleteAction creates the line_workspace_user delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "delete") {
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
		if deps.GetLineWorkspaceUserInUseIDs != nil {
			if inUse, _ := deps.GetLineWorkspaceUserInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteLineWorkspaceUser(ctx, &lineworkspaceuserpb.DeleteLineWorkspaceUserRequest{
			Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
		}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("line-workspace-users-table")
	})
}

// NewBulkDeleteAction creates the line_workspace_user bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "delete") {
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
		if deps.GetLineWorkspaceUserInUseIDs != nil {
			inUse, _ = deps.GetLineWorkspaceUserInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteLineWorkspaceUser(ctx, &lineworkspaceuserpb.DeleteLineWorkspaceUserRequest{
				Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
			}); err != nil {
				log.Printf("Failed to delete line_workspace_user %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("line-workspace-users-table")
	})
}

// NewSetStatusAction creates the line_workspace_user set-status action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadLineWorkspaceUser(ctx, &lineworkspaceuserpb.ReadLineWorkspaceUserRequest{
			Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
		})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateLineWorkspaceUser(ctx, &lineworkspaceuserpb.UpdateLineWorkspaceUserRequest{
			Data: &lineworkspaceuserpb.LineWorkspaceUser{
				Id:              id,
				WorkspaceUserId: record.GetWorkspaceUserId(),
				LineId:          record.GetLineId(),
				Scope:           record.GetScope(),
				Role:            record.GetRole(),
				IsOwner:         record.GetIsOwner(),
				Active:          status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("line-workspace-users-table")
	})
}

// NewBulkSetStatusAction creates the line_workspace_user bulk set-status action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadLineWorkspaceUser(ctx, &lineworkspaceuserpb.ReadLineWorkspaceUserRequest{
				Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateLineWorkspaceUser(ctx, &lineworkspaceuserpb.UpdateLineWorkspaceUserRequest{
				Data: &lineworkspaceuserpb.LineWorkspaceUser{
					Id:              id,
					WorkspaceUserId: record.GetWorkspaceUserId(),
					LineId:          record.GetLineId(),
					Scope:           record.GetScope(),
					Role:            record.GetRole(),
					IsOwner:         record.GetIsOwner(),
					Active:          status == "active",
				},
			})
		}
		return view.HTMXSuccess("line-workspace-users-table")
	})
}
