package category

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
)

// FormLabels holds flat i18n labels for the category drawer form template.
type FormLabels struct {
	Code        string
	Name        string
	Description string
}

// FormData is the template data for the category drawer form.
type FormData struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Code        string
	Name        string
	Description string
	Labels      FormLabels
	CommonLabels any
}

// ActionDeps holds dependencies for category action handlers.
type ActionDeps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	// Typed expenditure category CRUD operations
	CreateExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.CreateExpenditureCategoryRequest) (*expenditurecategorypb.CreateExpenditureCategoryResponse, error)
	ReadExpenditureCategory   func(ctx context.Context, req *expenditurecategorypb.ReadExpenditureCategoryRequest) (*expenditurecategorypb.ReadExpenditureCategoryResponse, error)
	UpdateExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.UpdateExpenditureCategoryRequest) (*expenditurecategorypb.UpdateExpenditureCategoryResponse, error)
	DeleteExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.DeleteExpenditureCategoryRequest) (*expenditurecategorypb.DeleteExpenditureCategoryResponse, error)
}

// formLabels maps ExpenditureCategoryLabels into the flat FormLabels struct for the template.
func formLabels(l centymo.ExpenditureCategoryLabels) FormLabels {
	code := l.Form.Code
	name := l.Form.Name
	description := l.Form.Description
	if code == "" {
		code = "Code"
	}
	if name == "" {
		name = "Name"
	}
	if description == "" {
		description = "Description"
	}
	return FormLabels{
		Code:        code,
		Name:        name,
		Description: description,
	}
}

// errLabels is a convenience helper to get category error labels.
func errLabels(l centymo.ExpenditureLabels) centymo.ExpenditureCategoryErrorLabels {
	e := l.Category.Errors
	if e.PermissionDenied == "" {
		e.PermissionDenied = "Permission denied"
	}
	if e.NotFound == "" {
		e.NotFound = "Category not found"
	}
	if e.IDRequired == "" {
		e.IDRequired = "ID is required"
	}
	if e.InvalidFormData == "" {
		e.InvalidFormData = "Invalid form data"
	}
	return e
}

// NewAddAction creates the category add action (GET = form, POST = create).
func NewAddAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		errs := errLabels(deps.Labels)
		if !perms.Can("expenditure_category", "create") {
			return centymo.HTMXError(errs.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("category-drawer-form", &FormData{
				FormAction:   deps.Routes.ExpenseCategoryAddURL,
				Labels:       formLabels(deps.Labels.Category),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create category
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(errs.InvalidFormData)
		}

		r := viewCtx.Request
		active := true

		resp, err := deps.CreateExpenditureCategory(ctx, &expenditurecategorypb.CreateExpenditureCategoryRequest{
			Data: &expenditurecategorypb.ExpenditureCategory{
				Code:        r.FormValue("code"),
				Name:        r.FormValue("name"),
				Description: strPtr(r.FormValue("description")),
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create expenditure category: %v", err)
			return centymo.HTMXError(err.Error())
		}
		_ = resp

		return centymo.HTMXSuccess("expenditure-categories-table")
	})
}

// NewEditAction creates the category edit action (GET = pre-filled form, POST = update).
func NewEditAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		errs := errLabels(deps.Labels)
		if !perms.Can("expenditure_category", "update") {
			return centymo.HTMXError(errs.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" || id == "{id}" {
			id = viewCtx.Request.URL.Query().Get("id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadExpenditureCategory(ctx, &expenditurecategorypb.ReadExpenditureCategoryRequest{
				Data: &expenditurecategorypb.ExpenditureCategory{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read expenditure category %s: %v", id, err)
				return centymo.HTMXError(errs.NotFound)
			}
			data := readResp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(errs.NotFound)
			}
			rec := data[0]

			return view.OK("category-drawer-form", &FormData{
				FormAction:   route.ResolveURL(deps.Routes.ExpenseCategoryEditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Code:         rec.GetCode(),
				Name:         rec.GetName(),
				Description:  rec.GetDescription(),
				Labels:       formLabels(deps.Labels.Category),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update category
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(errs.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdateExpenditureCategory(ctx, &expenditurecategorypb.UpdateExpenditureCategoryRequest{
			Data: &expenditurecategorypb.ExpenditureCategory{
				Id:          id,
				Code:        r.FormValue("code"),
				Name:        r.FormValue("name"),
				Description: strPtr(r.FormValue("description")),
			},
		})
		if err != nil {
			log.Printf("Failed to update expenditure category %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"formSuccess":true}`,
			},
		}
	})
}

// NewDeleteAction creates the category delete action (POST only).
// The row ID comes via query param (?id=xxx) or form field.
func NewDeleteAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		errs := errLabels(deps.Labels)
		if !perms.Can("expenditure_category", "delete") {
			return centymo.HTMXError(errs.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(errs.IDRequired)
		}

		_, err := deps.DeleteExpenditureCategory(ctx, &expenditurecategorypb.DeleteExpenditureCategoryRequest{
			Data: &expenditurecategorypb.ExpenditureCategory{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete expenditure category %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expenditure-categories-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}
