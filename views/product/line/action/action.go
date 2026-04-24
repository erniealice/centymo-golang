package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
)

// FormLabels holds i18n labels for the line drawer form template.
type FormLabels struct {
	Name            string
	Description     string
	DescPlaceholder string
	Active          string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string
	DescriptionInfo string
	ActiveInfo      string
}

// FormData is the template data for the line drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for line action handlers.
type Deps struct {
	Routes     centymo.ProductLineRoutes
	Labels     centymo.ProductLineLabels
	CreateLine func(ctx context.Context, req *linepb.CreateLineRequest) (*linepb.CreateLineResponse, error)
	ReadLine   func(ctx context.Context, req *linepb.ReadLineRequest) (*linepb.ReadLineResponse, error)
	UpdateLine func(ctx context.Context, req *linepb.UpdateLineRequest) (*linepb.UpdateLineResponse, error)
	DeleteLine func(ctx context.Context, req *linepb.DeleteLineRequest) (*linepb.DeleteLineResponse, error)
}

func formLabels(labels centymo.ProductLineLabels) FormLabels {
	return FormLabels{
		Name:            labels.Form.Name,
		Description:     labels.Form.Description,
		DescPlaceholder: labels.Form.DescPlaceholder,
		Active:          labels.Form.Active,
		// Info fields sourced from centymo.ProductLineFormLabels (populated from lyngua JSON + defaults).
		NameInfo:        labels.Form.NameInfo,
		DescriptionInfo: labels.Form.DescriptionInfo,
		ActiveInfo:      labels.Form.ActiveInfo,
	}
}

// NewAddAction creates the line add action.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("product-line-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		req := &linepb.CreateLineRequest{
			Data: &linepb.Line{
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Active:      active,
			},
		}

		if _, err := deps.CreateLine(ctx, req); err != nil {
			log.Printf("Failed to create line: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}

// NewEditAction creates the line edit action.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadLine(ctx, &linepb.ReadLineRequest{Data: &linepb.Line{Id: id}})
			if err != nil {
				log.Printf("Failed to read line %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := data[0]

			return view.OK("product-line-drawer-form", &FormData{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         record.GetName(),
				Description:  record.GetDescription(),
				Active:       record.GetActive(),
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		req := &linepb.UpdateLineRequest{
			Data: &linepb.Line{
				Id:          id,
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Active:      active,
			},
		}

		if _, err := deps.UpdateLine(ctx, req); err != nil {
			log.Printf("Failed to update line %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}

// NewDeleteAction creates the line delete action.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "delete") {
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

		if _, err := deps.DeleteLine(ctx, &linepb.DeleteLineRequest{Data: &linepb.Line{Id: id}}); err != nil {
			log.Printf("Failed to delete line %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}

// NewBulkDeleteAction creates the line bulk delete action.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			if id == "" {
				continue
			}
			if _, err := deps.DeleteLine(ctx, &linepb.DeleteLineRequest{Data: &linepb.Line{Id: id}}); err != nil {
				log.Printf("Failed to delete line %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}

// NewSetStatusAction toggles line active status.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if status != "active" && status != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		readResp, err := deps.ReadLine(ctx, &linepb.ReadLineRequest{Data: &linepb.Line{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		if _, err := deps.UpdateLine(ctx, &linepb.UpdateLineRequest{
			Data: &linepb.Line{
				Id:          id,
				Name:        record.GetName(),
				Description: record.GetDescription(),
				Active:      status == "active",
			},
		}); err != nil {
			log.Printf("Failed to update line status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}

// NewBulkSetStatusAction toggles line active status for multiple rows.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if status != "active" && status != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadLine(ctx, &linepb.ReadLineRequest{Data: &linepb.Line{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			if _, err := deps.UpdateLine(ctx, &linepb.UpdateLineRequest{
				Data: &linepb.Line{
					Id:          id,
					Name:        record.GetName(),
					Description: record.GetDescription(),
					Active:      status == "active",
				},
			}); err != nil {
				log.Printf("Failed to update line status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("product-lines-table")
	})
}
