package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/centymo-golang/domain/procurement/views/cost_schedule/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
)

// NewEditAction creates the cost_schedule edit action.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			var record *costschedulepb.CostSchedule
			if deps.GetCostScheduleItemPageData != nil {
				resp, err := deps.GetCostScheduleItemPageData(ctx, &costschedulepb.GetCostScheduleItemPageDataRequest{
					CostScheduleId: id,
				})
				if err != nil || resp == nil || resp.GetCostSchedule() == nil {
					return view.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetCostSchedule()
			} else {
				resp, err := deps.ReadCostSchedule(ctx, &costschedulepb.ReadCostScheduleRequest{
					Data: &costschedulepb.CostSchedule{Id: id},
				})
				if err != nil || len(resp.GetData()) == 0 {
					return view.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetData()[0]
			}
			return view.OK("cost-schedule-drawer-form", &form.Data{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         record.GetName(),
				Description:  record.GetDescription(),
				Active:       record.GetActive(),
				Labels:       buildFormLabels(deps.Labels),
				CommonLabels: deps.CommonLabels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		name := r.FormValue("name")
		description := r.FormValue("description")
		active := r.FormValue("active") != "false"
		req := &costschedulepb.UpdateCostScheduleRequest{
			Data: &costschedulepb.CostSchedule{
				Id:          id,
				Name:        name,
				Description: &description,
				Active:      active,
			},
		}
		if _, err := deps.UpdateCostSchedule(ctx, req); err != nil {
			log.Printf("Failed to update cost schedule %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("cost-schedules-table")
	})
}
