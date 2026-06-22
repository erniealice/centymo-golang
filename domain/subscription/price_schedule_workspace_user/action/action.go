package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	pswu "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user"
	"github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	pswupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule_workspace_user"
)

type Deps struct {
	Routes                                pswu.Routes
	Labels                                pswu.Labels
	CreatePriceScheduleWorkspaceUser      func(ctx context.Context, req *pswupb.CreatePriceScheduleWorkspaceUserRequest) (*pswupb.CreatePriceScheduleWorkspaceUserResponse, error)
	ReadPriceScheduleWorkspaceUser        func(ctx context.Context, req *pswupb.ReadPriceScheduleWorkspaceUserRequest) (*pswupb.ReadPriceScheduleWorkspaceUserResponse, error)
	UpdatePriceScheduleWorkspaceUser      func(ctx context.Context, req *pswupb.UpdatePriceScheduleWorkspaceUserRequest) (*pswupb.UpdatePriceScheduleWorkspaceUserResponse, error)
	DeletePriceScheduleWorkspaceUser      func(ctx context.Context, req *pswupb.DeletePriceScheduleWorkspaceUserRequest) (*pswupb.DeletePriceScheduleWorkspaceUserResponse, error)
	GetPriceScheduleWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData writes the POST body onto a PriceScheduleWorkspaceUser.
// Shared by Add (no id) and Edit (id set by caller). FK fields are required
// scalars on the proto — always set them from the form value (empty string is
// allowed to clear, consistent with proto zero value).
func applyFormToData(r *http.Request) *pswupb.PriceScheduleWorkspaceUser {
	data := &pswupb.PriceScheduleWorkspaceUser{
		PriceScheduleId: strings.TrimSpace(r.FormValue("price_schedule_id")),
		WorkspaceUserId: strings.TrimSpace(r.FormValue("workspace_user_id")),
		Scope:           strings.TrimSpace(r.FormValue("scope")),
		Role:            strings.TrimSpace(r.FormValue("role")),
		IsOwner:         r.FormValue("is_owner") == "true",
		Active:          r.FormValue("active") == "true",
	}
	return data
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("price-schedule-workspace-user-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &pswupb.CreatePriceScheduleWorkspaceUserRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreatePriceScheduleWorkspaceUser(ctx, req); err != nil {
			log.Printf("Failed to create price_schedule_workspace_user: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}

func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")

		if !perms.Can("price_schedule_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPriceScheduleWorkspaceUser(ctx, &pswupb.ReadPriceScheduleWorkspaceUserRequest{Data: &pswupb.PriceScheduleWorkspaceUser{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]

			return view.OK("price-schedule-workspace-user-drawer-form", &form.Data{
				FormAction:      route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:          true,
				ID:              id,
				PriceScheduleId: record.GetPriceScheduleId(),
				WorkspaceUserId: record.GetWorkspaceUserId(),
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
		if _, err := deps.UpdatePriceScheduleWorkspaceUser(ctx, &pswupb.UpdatePriceScheduleWorkspaceUserRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "delete") {
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
		if deps.GetPriceScheduleWorkspaceUserInUseIDs != nil {
			if inUse, _ := deps.GetPriceScheduleWorkspaceUserInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeletePriceScheduleWorkspaceUser(ctx, &pswupb.DeletePriceScheduleWorkspaceUserRequest{Data: &pswupb.PriceScheduleWorkspaceUser{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "delete") {
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
		if deps.GetPriceScheduleWorkspaceUserInUseIDs != nil {
			inUse, _ = deps.GetPriceScheduleWorkspaceUserInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeletePriceScheduleWorkspaceUser(ctx, &pswupb.DeletePriceScheduleWorkspaceUserRequest{Data: &pswupb.PriceScheduleWorkspaceUser{Id: id}}); err != nil {
				log.Printf("Failed to delete price_schedule_workspace_user %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPriceScheduleWorkspaceUser(ctx, &pswupb.ReadPriceScheduleWorkspaceUserRequest{Data: &pswupb.PriceScheduleWorkspaceUser{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdatePriceScheduleWorkspaceUser(ctx, &pswupb.UpdatePriceScheduleWorkspaceUserRequest{
			Data: &pswupb.PriceScheduleWorkspaceUser{
				Id:              id,
				PriceScheduleId: record.GetPriceScheduleId(),
				WorkspaceUserId: record.GetWorkspaceUserId(),
				Scope:           record.GetScope(),
				Role:            record.GetRole(),
				IsOwner:         record.GetIsOwner(),
				Active:          status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPriceScheduleWorkspaceUser(ctx, &pswupb.ReadPriceScheduleWorkspaceUserRequest{Data: &pswupb.PriceScheduleWorkspaceUser{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdatePriceScheduleWorkspaceUser(ctx, &pswupb.UpdatePriceScheduleWorkspaceUserRequest{
				Data: &pswupb.PriceScheduleWorkspaceUser{
					Id:              id,
					PriceScheduleId: record.GetPriceScheduleId(),
					WorkspaceUserId: record.GetWorkspaceUserId(),
					Scope:           record.GetScope(),
					Role:            record.GetRole(),
					IsOwner:         record.GetIsOwner(),
					Active:          status == "active",
				},
			})
		}
		return view.HTMXSuccess("price-schedule-workspace-users-table")
	})
}
