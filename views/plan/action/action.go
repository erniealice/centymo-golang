package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name            string
	NamePlaceholder string
	Description     string
	DescPlaceholder string
	FulfillmentType string
	TypeSchedule    string
	TypeLicense     string
	TypeContent     string
	TypePhysical    string
}

// FormData is the template data for the plan drawer form.
type FormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	Name            string
	Description     string
	FulfillmentType string
	Labels          FormLabels
	CommonLabels    any
}

// Deps holds dependencies for plan action handlers.
type Deps struct {
	Routes centymo.PlanRoutes
	Labels centymo.PlanLabels

	// Typed plan operations
	CreatePlan    func(ctx context.Context, req *planpb.CreatePlanRequest) (*planpb.CreatePlanResponse, error)
	ReadPlan      func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	UpdatePlan    func(ctx context.Context, req *planpb.UpdatePlanRequest) (*planpb.UpdatePlanResponse, error)
	DeletePlan    func(ctx context.Context, req *planpb.DeletePlanRequest) (*planpb.DeletePlanResponse, error)
	SetPlanActive func(ctx context.Context, id string, active bool) error
}

func formLabels(l centymo.PlanLabels) FormLabels {
	return FormLabels{
		Name:            l.Form.Name,
		NamePlaceholder: l.Form.NamePlaceholder,
		Description:     l.Form.Description,
		DescPlaceholder: l.Form.DescPlaceholder,
		FulfillmentType: l.Form.FulfillmentType,
		TypeSchedule:    l.Form.TypeSchedule,
		TypeLicense:     l.Form.TypeLicense,
		TypeContent:     l.Form.TypeContent,
		TypePhysical:    l.Form.TypePhysical,
	}
}

// NewAddAction creates the plan add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("plan-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				FulfillmentType: "schedule",
				Labels:          formLabels(deps.Labels),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — create plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		fulfillmentType := r.FormValue("fulfillment_type")

		resp, err := deps.CreatePlan(ctx, &planpb.CreatePlanRequest{
			Data: &planpb.Plan{
				Name:            r.FormValue("name"),
				Description:     strPtr(r.FormValue("description")),
				FulfillmentType: strPtr(fulfillmentType),
				Active:          true,
			},
		})
		if err != nil {
			log.Printf("Failed to create plan: %v", err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to new plan detail
		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", newID),
				},
			}
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewEditAction creates the plan edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{
				Data: &planpb.Plan{Id: &id},
			})
			if err != nil {
				log.Printf("Failed to read plan %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			return view.OK("plan-drawer-form", &FormData{
				FormAction:      route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:          true,
				ID:              id,
				Name:            record.GetName(),
				Description:     record.GetDescription(),
				FulfillmentType: record.GetFulfillmentType(),
				Labels:          formLabels(deps.Labels),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — update plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		fulfillmentType := r.FormValue("fulfillment_type")

		_, err := deps.UpdatePlan(ctx, &planpb.UpdatePlanRequest{
			Data: &planpb.Plan{
				Id:              &id,
				Name:            r.FormValue("name"),
				Description:     strPtr(r.FormValue("description")),
				FulfillmentType: strPtr(fulfillmentType),
			},
		})
		if err != nil {
			log.Printf("Failed to update plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to detail page
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", id),
			},
		}
	})
}

// NewDeleteAction creates the plan delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "delete") {
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

		_, err := deps.DeletePlan(ctx, &planpb.DeletePlanRequest{
			Data: &planpb.Plan{Id: &id},
		})
		if err != nil {
			log.Printf("Failed to delete plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewSetStatusAction creates the plan activate/deactivate action (POST only).
// Expects query params: ?id={planId}&status={active|inactive}
//
// Uses SetPlanActive (raw map update) instead of protobuf because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "update") {
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
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if err := deps.SetPlanActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update plan status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewBulkSetStatusAction creates the plan bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetPlanActive(ctx, id, active); err != nil {
				log.Printf("Failed to update plan status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}
