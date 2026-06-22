package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	product_plan_staff "github.com/erniealice/centymo-golang/domain/product/product_plan_staff"
	"github.com/erniealice/centymo-golang/domain/product/product_plan_staff/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
)

type Deps struct {
	Routes                      product_plan_staff.Routes
	Labels                      product_plan_staff.Labels
	CreateProductPlanStaff      func(ctx context.Context, req *productplanstaffpb.CreateProductPlanStaffRequest) (*productplanstaffpb.CreateProductPlanStaffResponse, error)
	ReadProductPlanStaff        func(ctx context.Context, req *productplanstaffpb.ReadProductPlanStaffRequest) (*productplanstaffpb.ReadProductPlanStaffResponse, error)
	UpdateProductPlanStaff      func(ctx context.Context, req *productplanstaffpb.UpdateProductPlanStaffRequest) (*productplanstaffpb.UpdateProductPlanStaffResponse, error)
	DeleteProductPlanStaff      func(ctx context.Context, req *productplanstaffpb.DeleteProductPlanStaffRequest) (*productplanstaffpb.DeleteProductPlanStaffResponse, error)
	GetProductPlanStaffInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// applyFormToData writes the POST body onto a ProductPlanStaff. Shared by
// Add (no id) and Edit (id set by caller). StaffId and ProductPlanId are
// required non-pointer FKs in the proto (plain string). Role is free-text.
func applyFormToData(r *http.Request) *productplanstaffpb.ProductPlanStaff {
	return &productplanstaffpb.ProductPlanStaff{
		StaffId:       strings.TrimSpace(r.FormValue("staff_id")),
		ProductPlanId: strings.TrimSpace(r.FormValue("product_plan_id")),
		Role:          strings.TrimSpace(r.FormValue("role")),
		Active:        r.FormValue("active") == "true",
	}
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("product-plan-staff-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &productplanstaffpb.CreateProductPlanStaffRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateProductPlanStaff(ctx, req); err != nil {
			log.Printf("Failed to create product plan staff: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("product-plan-staffs-table")
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
		if !perms.Can("product_plan_staff", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProductPlanStaff(ctx, &productplanstaffpb.ReadProductPlanStaffRequest{Data: &productplanstaffpb.ProductPlanStaff{Id: id}})
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

			return view.OK("product-plan-staff-drawer-form", &form.Data{
				FormAction:    formAction,
				IsEdit:        !isClone,
				ID:            formID,
				StaffID:       record.GetStaffId(),
				ProductPlanID: record.GetProductPlanId(),
				Role:          record.GetRole(),
				Active:        record.GetActive(),
				Labels:        deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateProductPlanStaff(ctx, &productplanstaffpb.UpdateProductPlanStaffRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("product-plan-staffs-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "delete") {
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
		if deps.GetProductPlanStaffInUseIDs != nil {
			if inUse, _ := deps.GetProductPlanStaffInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteProductPlanStaff(ctx, &productplanstaffpb.DeleteProductPlanStaffRequest{Data: &productplanstaffpb.ProductPlanStaff{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("product-plan-staffs-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "delete") {
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
		if deps.GetProductPlanStaffInUseIDs != nil {
			inUse, _ = deps.GetProductPlanStaffInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteProductPlanStaff(ctx, &productplanstaffpb.DeleteProductPlanStaffRequest{Data: &productplanstaffpb.ProductPlanStaff{Id: id}}); err != nil {
				log.Printf("Failed to delete product plan staff %s during bulk: %v", id, err)
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
		return view.HTMXSuccess("product-plan-staffs-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadProductPlanStaff(ctx, &productplanstaffpb.ReadProductPlanStaffRequest{Data: &productplanstaffpb.ProductPlanStaff{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateProductPlanStaff(ctx, &productplanstaffpb.UpdateProductPlanStaffRequest{
			Data: &productplanstaffpb.ProductPlanStaff{
				Id:            id,
				StaffId:       record.GetStaffId(),
				ProductPlanId: record.GetProductPlanId(),
				Role:          record.GetRole(),
				Active:        status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("product-plan-staffs-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadProductPlanStaff(ctx, &productplanstaffpb.ReadProductPlanStaffRequest{Data: &productplanstaffpb.ProductPlanStaff{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateProductPlanStaff(ctx, &productplanstaffpb.UpdateProductPlanStaffRequest{
				Data: &productplanstaffpb.ProductPlanStaff{
					Id:            id,
					StaffId:       record.GetStaffId(),
					ProductPlanId: record.GetProductPlanId(),
					Role:          record.GetRole(),
					Active:        status == "active",
				},
			})
		}
		return view.HTMXSuccess("product-plan-staffs-table")
	})
}
