package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	plan_group "github.com/erniealice/centymo-golang/domain/product/plan_group"
	"github.com/erniealice/centymo-golang/domain/product/plan_group/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	plangroupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group"
)

type Deps struct {
	Routes          plan_group.Routes
	Labels          plan_group.Labels
	CreatePlanGroup func(ctx context.Context, req *plangroupb.CreatePlanGroupRequest) (*plangroupb.CreatePlanGroupResponse, error)
	ReadPlanGroup   func(ctx context.Context, req *plangroupb.ReadPlanGroupRequest) (*plangroupb.ReadPlanGroupResponse, error)
	UpdatePlanGroup func(ctx context.Context, req *plangroupb.UpdatePlanGroupRequest) (*plangroupb.UpdatePlanGroupResponse, error)
	DeletePlanGroup func(ctx context.Context, req *plangroupb.DeletePlanGroupRequest) (*plangroupb.DeletePlanGroupResponse, error)
	// ListPlanGroups is used to populate the parent-group picker.
	ListPlanGroups       func(ctx context.Context, req *plangroupb.ListPlanGroupsRequest) (*plangroupb.ListPlanGroupsResponse, error)
	GetPlanGroupInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

func loadParentPairs(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListPlanGroups == nil {
		return nil
	}
	resp, err := deps.ListPlanGroups(ctx, &plangroupb.ListPlanGroupsRequest{})
	if err != nil {
		return nil
	}
	pairs := make([]form.Pair, 0, len(resp.GetData()))
	for _, pg := range resp.GetData() {
		if pg == nil || !pg.GetActive() {
			continue
		}
		label := pg.GetName()
		if label == "" {
			label = pg.GetId()
		}
		pairs = append(pairs, form.Pair{ID: pg.GetId(), Label: label})
	}
	return pairs
}

// applyFormToData writes the POST body onto a PlanGroup. Shared by Add (no id)
// and Edit (id set by caller). Optional FK pointer parent_id and optional
// scalar code are set only when present so unset fields stay null.
func applyFormToData(r *http.Request) *plangroupb.PlanGroup {
	data := &plangroupb.PlanGroup{
		Name:   r.FormValue("name"),
		Active: r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("code")); v != "" {
		data.Code = &v
	}
	if v := strings.TrimSpace(r.FormValue("parent_id")); v != "" {
		data.ParentId = &v
	}
	return data
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			parentPairs := loadParentPairs(ctx, deps)
			return view.OK("plan-group-drawer-form", &form.Data{
				FormAction:    deps.Routes.AddURL,
				Active:        true,
				ParentOptions: form.BuildParentOptions(parentPairs, ""),
				Labels:        deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &plangroupb.CreatePlanGroupRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreatePlanGroup(ctx, req); err != nil {
			log.Printf("Failed to create plan group: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-groups-table")
	})
}

// NewEditAction creates the plan-group edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer
// pre-populated from the source record but wired to AddURL.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("plan_group", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPlanGroup(ctx, &plangroupb.ReadPlanGroupRequest{Data: &plangroupb.PlanGroup{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]

			name := record.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			parentPairs := loadParentPairs(ctx, deps)
			selectedParentID := record.GetParentId()

			return view.OK("plan-group-drawer-form", &form.Data{
				FormAction:    formAction,
				IsEdit:        !isClone,
				ID:            formID,
				Name:          name,
				Code:          record.GetCode(),
				ParentID:      selectedParentID,
				ParentLabel:   form.FindLabel(parentPairs, selectedParentID),
				ParentOptions: form.BuildParentOptions(parentPairs, selectedParentID),
				Active:        record.GetActive(),
				Labels:        deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdatePlanGroup(ctx, &plangroupb.UpdatePlanGroupRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-groups-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "delete") {
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
		if deps.GetPlanGroupInUseIDs != nil {
			if inUse, _ := deps.GetPlanGroupInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeletePlanGroup(ctx, &plangroupb.DeletePlanGroupRequest{Data: &plangroupb.PlanGroup{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-groups-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "delete") {
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
		if deps.GetPlanGroupInUseIDs != nil {
			inUse, _ = deps.GetPlanGroupInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeletePlanGroup(ctx, &plangroupb.DeletePlanGroupRequest{Data: &plangroupb.PlanGroup{Id: id}}); err != nil {
				log.Printf("Failed to delete plan group %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("plan-groups-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPlanGroup(ctx, &plangroupb.ReadPlanGroupRequest{Data: &plangroupb.PlanGroup{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdatePlanGroup(ctx, &plangroupb.UpdatePlanGroupRequest{
			Data: &plangroupb.PlanGroup{
				Id:       id,
				Name:     record.GetName(),
				Code:     record.Code,
				ParentId: record.ParentId,
				Active:   status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-groups-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPlanGroup(ctx, &plangroupb.ReadPlanGroupRequest{Data: &plangroupb.PlanGroup{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdatePlanGroup(ctx, &plangroupb.UpdatePlanGroupRequest{
				Data: &plangroupb.PlanGroup{
					Id:       id,
					Name:     record.GetName(),
					Code:     record.Code,
					ParentId: record.ParentId,
					Active:   status == "active",
				},
			})
		}
		return view.HTMXSuccess("plan-groups-table")
	})
}
