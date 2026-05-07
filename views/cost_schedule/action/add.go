package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/cost_schedule/form"
	"github.com/erniealice/pyeza-golang/view"

	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
)

// NewAddAction creates the cost_schedule add action.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_schedule", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("cost-schedule-drawer-form", &form.Data{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Labels:     buildFormLabels(deps.Labels),
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		name := r.FormValue("name")
		description := r.FormValue("description")
		active := r.FormValue("active") != "false"
		req := &costschedulepb.CreateCostScheduleRequest{
			Data: &costschedulepb.CostSchedule{
				Name:        name,
				Description: &description,
				Active:      active,
			},
		}
		if _, err := deps.CreateCostSchedule(ctx, req); err != nil {
			log.Printf("Failed to create cost schedule: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-schedules-table")
	})
}
