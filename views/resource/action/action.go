package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	resourcepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/resource"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name            string
	NamePlaceholder string
	Description     string
	DescPlaceholder string
	ProductId       string
	UserId          string
}

// FormData is the template data for the resource drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	ProductId    string
	UserId       string
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for resource action handlers.
type Deps struct {
	Routes         centymo.ResourceRoutes
	Labels         centymo.ResourceLabels
	CreateResource func(ctx context.Context, req *resourcepb.CreateResourceRequest) (*resourcepb.CreateResourceResponse, error)
	ReadResource   func(ctx context.Context, req *resourcepb.ReadResourceRequest) (*resourcepb.ReadResourceResponse, error)
	UpdateResource func(ctx context.Context, req *resourcepb.UpdateResourceRequest) (*resourcepb.UpdateResourceResponse, error)
	DeleteResource func(ctx context.Context, req *resourcepb.DeleteResourceRequest) (*resourcepb.DeleteResourceResponse, error)
}

func formLabels(l centymo.ResourceLabels) FormLabels {
	return FormLabels{
		Name:            l.Form.Name,
		NamePlaceholder: l.Form.NamePlaceholder,
		Description:     l.Form.Description,
		DescPlaceholder: l.Form.DescPlaceholder,
		ProductId:       l.Form.ProductId,
		UserId:          l.Form.UserId,
	}
}

func strPtr(s string) *string {
	return &s
}

// NewAddAction creates the resource add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("resource-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil,
			})
		}

		// POST — create resource
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		description := r.FormValue("description")
		userId := r.FormValue("user_id")

		req := &resourcepb.CreateResourceRequest{
			Data: &resourcepb.Resource{
				Name:      r.FormValue("name"),
				ProductId: r.FormValue("product_id"),
				Active:    true,
			},
		}
		if description != "" {
			req.Data.Description = strPtr(description)
		}
		if userId != "" {
			req.Data.UserId = strPtr(userId)
		}

		if _, err := deps.CreateResource(ctx, req); err != nil {
			log.Printf("Failed to create resource: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("resources-table")
	})
}

// NewEditAction creates the resource edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadResource(ctx, &resourcepb.ReadResourceRequest{
				Data: &resourcepb.Resource{Id: id},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				log.Printf("Failed to read resource %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readResp.GetData()[0]

			return view.OK("resource-drawer-form", &FormData{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         record.GetName(),
				Description:  record.GetDescription(),
				ProductId:    record.GetProductId(),
				UserId:       record.GetUserId(),
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil,
			})
		}

		// POST — update resource
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		description := r.FormValue("description")
		userId := r.FormValue("user_id")

		req := &resourcepb.UpdateResourceRequest{
			Data: &resourcepb.Resource{
				Id:        id,
				Name:      r.FormValue("name"),
				ProductId: r.FormValue("product_id"),
			},
		}
		if description != "" {
			req.Data.Description = strPtr(description)
		}
		if userId != "" {
			req.Data.UserId = strPtr(userId)
		}

		if _, err := deps.UpdateResource(ctx, req); err != nil {
			log.Printf("Failed to update resource %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("resources-table")
	})
}

// NewDeleteAction creates the resource delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		if _, err := deps.DeleteResource(ctx, &resourcepb.DeleteResourceRequest{
			Data: &resourcepb.Resource{Id: id},
		}); err != nil {
			log.Printf("Failed to delete resource %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("resources-table")
	})
}

// NewSetStatusAction creates the resource activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		readResp, err := deps.ReadResource(ctx, &resourcepb.ReadResourceRequest{
			Data: &resourcepb.Resource{Id: id},
		})
		if err != nil || len(readResp.GetData()) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]

		_, err = deps.UpdateResource(ctx, &resourcepb.UpdateResourceRequest{
			Data: &resourcepb.Resource{
				Id:          id,
				Name:        record.GetName(),
				Description: record.Description,
				ProductId:   record.GetProductId(),
				UserId:      record.UserId,
				Active:      targetStatus == "active",
			},
		})
		if err != nil {
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("resources-table")
	})
}

// NewBulkSetStatusAction creates the resource bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadResource(ctx, &resourcepb.ReadResourceRequest{
				Data: &resourcepb.Resource{Id: id},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateResource(ctx, &resourcepb.UpdateResourceRequest{
				Data: &resourcepb.Resource{
					Id:          id,
					Name:        record.GetName(),
					Description: record.Description,
					ProductId:   record.GetProductId(),
					UserId:      record.UserId,
					Active:      targetStatus == "active",
				},
			})
		}

		return centymo.HTMXSuccess("resources-table")
	})
}

// NewBulkDeleteAction creates the resource bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("resource", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeleteResource(ctx, &resourcepb.DeleteResourceRequest{
					Data: &resourcepb.Resource{Id: id},
				})
			}
		}

		return centymo.HTMXSuccess("resources-table")
	})
}
