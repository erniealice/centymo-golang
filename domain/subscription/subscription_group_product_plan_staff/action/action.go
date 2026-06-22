package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	sgpps "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// Deps holds all action-layer dependencies for
// subscription_group_product_plan_staff CRUD.
type Deps struct {
	Routes                                       sgpps.Routes
	Labels                                       sgpps.Labels
	CreateSubscriptionGroupProductPlanStaff      func(ctx context.Context, req *sgppspb.CreateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.CreateSubscriptionGroupProductPlanStaffResponse, error)
	ReadSubscriptionGroupProductPlanStaff        func(ctx context.Context, req *sgppspb.ReadSubscriptionGroupProductPlanStaffRequest) (*sgppspb.ReadSubscriptionGroupProductPlanStaffResponse, error)
	UpdateSubscriptionGroupProductPlanStaff      func(ctx context.Context, req *sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.UpdateSubscriptionGroupProductPlanStaffResponse, error)
	DeleteSubscriptionGroupProductPlanStaff      func(ctx context.Context, req *sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest) (*sgppspb.DeleteSubscriptionGroupProductPlanStaffResponse, error)
	GetSubscriptionGroupProductPlanStaffInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Optional FK pickers — nil disables the picker (field shows free-text).
	ListSubscriptionGroupOptions func(ctx context.Context) []form.Pair
	ListProductPlanOptions       func(ctx context.Context) []form.Pair
	ListStaffOptions             func(ctx context.Context) []form.Pair
}

// applyFormToData writes the POST body onto a
// SubscriptionGroupProductPlanStaff. FK fields (subscription_group_id,
// product_plan_id, staff_id) are required strings in the proto — set only when
// the form value is non-empty so unset fields stay empty rather than overwriting
// with blank.
func applyFormToData(r *http.Request) *sgppspb.SubscriptionGroupProductPlanStaff {
	data := &sgppspb.SubscriptionGroupProductPlanStaff{
		Role:   strings.TrimSpace(r.FormValue("role")),
		Active: r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("subscription_group_id")); v != "" {
		data.SubscriptionGroupId = v
	}
	if v := strings.TrimSpace(r.FormValue("product_plan_id")); v != "" {
		data.ProductPlanId = v
	}
	if v := strings.TrimSpace(r.FormValue("staff_id")); v != "" {
		data.StaffId = v
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

func loadStaffOpts(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListStaffOptions == nil {
		return nil
	}
	return deps.ListStaffOptions(ctx)
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			groupPairs := loadGroupOpts(ctx, deps)
			planPairs := loadPlanOpts(ctx, deps)
			staffPairs := loadStaffOpts(ctx, deps)
			return view.OK("sgpps-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				Active:                true,
				SubscriptionGroupOpts: form.BuildAutoCompleteOptions(groupPairs, ""),
				ProductPlanOpts:       form.BuildAutoCompleteOptions(planPairs, ""),
				StaffOpts:             form.BuildAutoCompleteOptions(staffPairs, ""),
				Labels:                deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &sgppspb.CreateSubscriptionGroupProductPlanStaffRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateSubscriptionGroupProductPlanStaff(ctx, req); err != nil {
			log.Printf("Failed to create subscription_group_product_plan_staff: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("sgpps-table")
	})
}

// NewEditAction creates the edit action (GET = form, POST = update). Supports
// ?clone=1 to pre-populate with source data wired to AddURL.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("subscription_group_product_plan_staff", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSubscriptionGroupProductPlanStaff(ctx, &sgppspb.ReadSubscriptionGroupProductPlanStaffRequest{Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id}})
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
			staffPairs := loadStaffOpts(ctx, deps)

			selectedGroupID := record.GetSubscriptionGroupId()
			selectedPlanID := record.GetProductPlanId()
			selectedStaffID := record.GetStaffId()

			return view.OK("sgpps-drawer-form", &form.Data{
				FormAction:             formAction,
				IsEdit:                 !isClone,
				ID:                     formID,
				SubscriptionGroupID:    selectedGroupID,
				SubscriptionGroupLabel: form.FindLabel(groupPairs, selectedGroupID),
				SubscriptionGroupOpts:  form.BuildAutoCompleteOptions(groupPairs, selectedGroupID),
				ProductPlanID:          selectedPlanID,
				ProductPlanLabel:       form.FindLabel(planPairs, selectedPlanID),
				ProductPlanOpts:        form.BuildAutoCompleteOptions(planPairs, selectedPlanID),
				StaffID:                selectedStaffID,
				StaffLabel:             form.FindLabel(staffPairs, selectedStaffID),
				StaffOpts:              form.BuildAutoCompleteOptions(staffPairs, selectedStaffID),
				Role:                   record.GetRole(),
				Active:                 record.GetActive(),
				Labels:                 deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateSubscriptionGroupProductPlanStaff(ctx, &sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("sgpps-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "delete") {
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
		if deps.GetSubscriptionGroupProductPlanStaffInUseIDs != nil {
			if inUse, _ := deps.GetSubscriptionGroupProductPlanStaffInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteSubscriptionGroupProductPlanStaff(ctx, &sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest{Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("sgpps-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "delete") {
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
		if deps.GetSubscriptionGroupProductPlanStaffInUseIDs != nil {
			inUse, _ = deps.GetSubscriptionGroupProductPlanStaffInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteSubscriptionGroupProductPlanStaff(ctx, &sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest{Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id}}); err != nil {
				log.Printf("Failed to delete subscription_group_product_plan_staff %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("sgpps-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadSubscriptionGroupProductPlanStaff(ctx, &sgppspb.ReadSubscriptionGroupProductPlanStaffRequest{Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateSubscriptionGroupProductPlanStaff(ctx, &sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest{
			Data: &sgppspb.SubscriptionGroupProductPlanStaff{
				Id:                  id,
				SubscriptionGroupId: record.GetSubscriptionGroupId(),
				ProductPlanId:       record.GetProductPlanId(),
				StaffId:             record.GetStaffId(),
				Role:                record.GetRole(),
				Active:              status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("sgpps-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadSubscriptionGroupProductPlanStaff(ctx, &sgppspb.ReadSubscriptionGroupProductPlanStaffRequest{Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateSubscriptionGroupProductPlanStaff(ctx, &sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest{
				Data: &sgppspb.SubscriptionGroupProductPlanStaff{
					Id:                  id,
					SubscriptionGroupId: record.GetSubscriptionGroupId(),
					ProductPlanId:       record.GetProductPlanId(),
					StaffId:             record.GetStaffId(),
					Role:                record.GetRole(),
					Active:              status == "active",
				},
			})
		}
		return view.HTMXSuccess("sgpps-table")
	})
}
