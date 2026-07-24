package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
)

// sgppEntity is the permission entity for the class row. See labels.go /
// plan.md §5 (D-5): subscription_group_product_plan:create/update/delete are
// distinct new grants from the assignment edge's own subscription_group_
// product_plan_staff:* grants.
const sgppEntity = "subscription_group_product_plan"

// sgppTableID is the S7 admin list's table-card id (the HTMXSuccess refresh
// target).
const sgppTableID = "subscription-group-product-plans-table"

// Deps holds all action-layer dependencies for subscription_group_product_plan
// CRUD (S7 admin add/edit/delete/bulk-delete).
type Deps struct {
	Routes                             sgpp.Routes
	Labels                             sgpp.Labels
	CreateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.CreateSubscriptionGroupProductPlanRequest) (*sgpppb.CreateSubscriptionGroupProductPlanResponse, error)
	ReadSubscriptionGroupProductPlan   func(ctx context.Context, req *sgpppb.ReadSubscriptionGroupProductPlanRequest) (*sgpppb.ReadSubscriptionGroupProductPlanResponse, error)
	UpdateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.UpdateSubscriptionGroupProductPlanRequest) (*sgpppb.UpdateSubscriptionGroupProductPlanResponse, error)
	DeleteSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.DeleteSubscriptionGroupProductPlanRequest) (*sgpppb.DeleteSubscriptionGroupProductPlanResponse, error)
	// GetSubscriptionGroupProductPlanInUseIDs is the fail-closed in-use gate
	// for the delete row action (plan.md §2, the section-delete-guard
	// precedent) — the server-side guard also lives inside the Delete/Update
	// use cases (espyna), this closure only drives the UI's disabled state.
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Optional FK pickers — nil disables the picker (field shows free-text).
	ListSubscriptionGroupOptions func(ctx context.Context) []form.Pair
	ListProductPlanOptions       func(ctx context.Context) []form.Pair
	ListJobTemplateOptions       func(ctx context.Context) []form.Pair
}

// applyFormToData writes the POST body onto a SubscriptionGroupProductPlan.
func applyFormToData(r *http.Request) *sgpppb.SubscriptionGroupProductPlan {
	data := &sgpppb.SubscriptionGroupProductPlan{
		Active: r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("subscription_group_id")); v != "" {
		data.SubscriptionGroupId = v
	}
	if v := strings.TrimSpace(r.FormValue("product_plan_id")); v != "" {
		data.ProductPlanId = v
	}
	if v := strings.TrimSpace(r.FormValue("job_template_id")); v != "" {
		data.JobTemplateId = v
	}
	if r.FormValue("status") == "excluded" {
		data.Status = sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED
	} else {
		data.Status = sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_ACTIVE
	}
	return data
}

func loadGroupOpts(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListSubscriptionGroupOptions == nil {
		return nil
	}
	return deps.ListSubscriptionGroupOptions(ctx)
}

func loadPlanOpts(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListProductPlanOptions == nil {
		return nil
	}
	return deps.ListProductPlanOptions(ctx)
}

func loadTemplateOpts(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListJobTemplateOptions == nil {
		return nil
	}
	return deps.ListJobTemplateOptions(ctx)
}

// NewAddAction creates the subscription_group_product_plan add action (GET =
// form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			groupPairs := loadGroupOpts(ctx, deps)
			planPairs := loadPlanOpts(ctx, deps)
			templatePairs := loadTemplateOpts(ctx, deps)
			return view.OK("sgpp-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				Active:                true,
				SubscriptionGroupOpts: form.BuildAutoCompleteOptions(groupPairs, ""),
				ProductPlanOpts:       form.BuildAutoCompleteOptions(planPairs, ""),
				JobTemplateOpts:       form.BuildAutoCompleteOptions(templatePairs, ""),
				Labels:                deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &sgpppb.CreateSubscriptionGroupProductPlanRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateSubscriptionGroupProductPlan(ctx, req); err != nil {
			log.Printf("Failed to create subscription_group_product_plan: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sgppTableID)
	})
}

// NewEditAction creates the subscription_group_product_plan edit action (GET
// = form, POST = update). Supports ?clone=1 to pre-populate from source data
// wired to AddURL.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can(sgppEntity, requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSubscriptionGroupProductPlan(ctx, &sgpppb.ReadSubscriptionGroupProductPlanRequest{
				Data: &sgpppb.SubscriptionGroupProductPlan{Id: id},
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

			groupPairs := loadGroupOpts(ctx, deps)
			planPairs := loadPlanOpts(ctx, deps)
			templatePairs := loadTemplateOpts(ctx, deps)

			selectedGroupID := record.GetSubscriptionGroupId()
			selectedPlanID := record.GetProductPlanId()
			selectedTemplateID := record.GetJobTemplateId()

			return view.OK("sgpp-drawer-form", &form.Data{
				FormAction:             formAction,
				IsEdit:                 !isClone,
				ID:                     formID,
				SubscriptionGroupID:    selectedGroupID,
				SubscriptionGroupLabel: form.FindLabel(groupPairs, selectedGroupID),
				SubscriptionGroupOpts:  form.BuildAutoCompleteOptions(groupPairs, selectedGroupID),
				ProductPlanID:          selectedPlanID,
				ProductPlanLabel:       form.FindLabel(planPairs, selectedPlanID),
				ProductPlanOpts:        form.BuildAutoCompleteOptions(planPairs, selectedPlanID),
				JobTemplateID:          selectedTemplateID,
				JobTemplateLabel:       form.FindLabel(templatePairs, selectedTemplateID),
				JobTemplateOpts:        form.BuildAutoCompleteOptions(templatePairs, selectedTemplateID),
				StatusExcluded:         record.GetStatus() == sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED,
				Active:                 record.GetActive(),
				Labels:                 deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateSubscriptionGroupProductPlan(ctx, &sgpppb.UpdateSubscriptionGroupProductPlanRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sgppTableID)
	})
}

// NewDeleteAction creates the subscription_group_product_plan delete action.
// In-use gating (plan.md §2) is defense-in-depth here (UI hint); the espyna
// Delete use case fail-closes the same guard server-side regardless.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "delete") {
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
		if deps.GetSubscriptionGroupProductPlanInUseIDs != nil {
			if inUse, _ := deps.GetSubscriptionGroupProductPlanInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteSubscriptionGroupProductPlan(ctx, &sgpppb.DeleteSubscriptionGroupProductPlanRequest{
			Data: &sgpppb.SubscriptionGroupProductPlan{Id: id},
		}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sgppTableID)
	})
}

// NewBulkDeleteAction creates the subscription_group_product_plan bulk-delete
// action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "delete") {
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
		if deps.GetSubscriptionGroupProductPlanInUseIDs != nil {
			inUse, _ = deps.GetSubscriptionGroupProductPlanInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteSubscriptionGroupProductPlan(ctx, &sgpppb.DeleteSubscriptionGroupProductPlanRequest{
				Data: &sgpppb.SubscriptionGroupProductPlan{Id: id},
			}); err != nil {
				log.Printf("Failed to delete subscription_group_product_plan %s during bulk: %v", id, err)
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
		return view.HTMXSuccess(sgppTableID)
	})
}
