package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	plan_group_plan "github.com/erniealice/centymo-golang/domain/product/plan_group_plan"
	"github.com/erniealice/centymo-golang/domain/product/plan_group_plan/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	plangroupplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group_plan"
)

// Deps holds all dependencies for plan_group_plan actions.
type Deps struct {
	Routes                   plan_group_plan.Routes
	Labels                   plan_group_plan.Labels
	CreatePlanGroupPlan      func(ctx context.Context, req *plangroupplanpb.CreatePlanGroupPlanRequest) (*plangroupplanpb.CreatePlanGroupPlanResponse, error)
	ReadPlanGroupPlan        func(ctx context.Context, req *plangroupplanpb.ReadPlanGroupPlanRequest) (*plangroupplanpb.ReadPlanGroupPlanResponse, error)
	UpdatePlanGroupPlan      func(ctx context.Context, req *plangroupplanpb.UpdatePlanGroupPlanRequest) (*plangroupplanpb.UpdatePlanGroupPlanResponse, error)
	DeletePlanGroupPlan      func(ctx context.Context, req *plangroupplanpb.DeletePlanGroupPlanRequest) (*plangroupplanpb.DeletePlanGroupPlanResponse, error)
	GetPlanGroupPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData writes the POST body onto a PlanGroupPlan. Shared by Add
// (no id) and Edit (id set by caller). SequenceOrder is optional; left nil
// when the form value is blank so the field stays NULL server-side.
func applyFormToData(r *http.Request) *plangroupplanpb.PlanGroupPlan {
	data := &plangroupplanpb.PlanGroupPlan{
		PlanGroupId: strings.TrimSpace(r.FormValue("plan_group_id")),
		PlanId:      strings.TrimSpace(r.FormValue("plan_id")),
		Active:      r.FormValue("active") == "true",
	}
	if s := strings.TrimSpace(r.FormValue("sequence_order")); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			data.SequenceOrder = &v32
		}
	}
	return data
}

// NewAddAction creates the plan_group_plan add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("plan-group-plan-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &plangroupplanpb.CreatePlanGroupPlanRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreatePlanGroupPlan(ctx, req); err != nil {
			log.Printf("Failed to create plan group plan: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-group-plans-table")
	})
}

// NewEditAction creates the plan_group_plan edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("plan_group_plan", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPlanGroupPlan(ctx, &plangroupplanpb.ReadPlanGroupPlanRequest{Data: &plangroupplanpb.PlanGroupPlan{Id: id}})
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

			seqOrder := ""
			if record.SequenceOrder != nil {
				seqOrder = strconv.FormatInt(int64(record.GetSequenceOrder()), 10)
			}

			return view.OK("plan-group-plan-drawer-form", &form.Data{
				FormAction:    formAction,
				IsEdit:        !isClone,
				ID:            formID,
				PlanGroupID:   record.GetPlanGroupId(),
				PlanID:        record.GetPlanId(),
				SequenceOrder: seqOrder,
				Active:        record.GetActive(),
				Labels:        deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdatePlanGroupPlan(ctx, &plangroupplanpb.UpdatePlanGroupPlanRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-group-plans-table")
	})
}

// NewDeleteAction creates the plan_group_plan delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "delete") {
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
		if deps.GetPlanGroupPlanInUseIDs != nil {
			if inUse, _ := deps.GetPlanGroupPlanInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeletePlanGroupPlan(ctx, &plangroupplanpb.DeletePlanGroupPlanRequest{Data: &plangroupplanpb.PlanGroupPlan{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-group-plans-table")
	})
}

// NewBulkDeleteAction creates the plan_group_plan bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "delete") {
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
		if deps.GetPlanGroupPlanInUseIDs != nil {
			inUse, _ = deps.GetPlanGroupPlanInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeletePlanGroupPlan(ctx, &plangroupplanpb.DeletePlanGroupPlanRequest{Data: &plangroupplanpb.PlanGroupPlan{Id: id}}); err != nil {
				log.Printf("Failed to delete plan group plan %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("plan-group-plans-table")
	})
}

// NewSetStatusAction creates the plan_group_plan set-status action.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPlanGroupPlan(ctx, &plangroupplanpb.ReadPlanGroupPlanRequest{Data: &plangroupplanpb.PlanGroupPlan{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdatePlanGroupPlan(ctx, &plangroupplanpb.UpdatePlanGroupPlanRequest{
			Data: &plangroupplanpb.PlanGroupPlan{
				Id:            id,
				PlanGroupId:   record.GetPlanGroupId(),
				PlanId:        record.GetPlanId(),
				SequenceOrder: record.SequenceOrder,
				Active:        status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("plan-group-plans-table")
	})
}

// NewBulkSetStatusAction creates the plan_group_plan bulk set-status action.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPlanGroupPlan(ctx, &plangroupplanpb.ReadPlanGroupPlanRequest{Data: &plangroupplanpb.PlanGroupPlan{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdatePlanGroupPlan(ctx, &plangroupplanpb.UpdatePlanGroupPlanRequest{
				Data: &plangroupplanpb.PlanGroupPlan{
					Id:            id,
					PlanGroupId:   record.GetPlanGroupId(),
					PlanId:        record.GetPlanId(),
					SequenceOrder: record.SequenceOrder,
					Active:        status == "active",
				},
			})
		}
		return view.HTMXSuccess("plan-group-plans-table")
	})
}
