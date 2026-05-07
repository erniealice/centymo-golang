package action

import (
	"context"
	"log"
	"net/http"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_subscription/form"
	"github.com/erniealice/pyeza-golang/view"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// NewAddAction creates the supplier_subscription add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			tz := pyezatypes.LocationFromContext(ctx)
			today := time.Now().In(tz)
			defaultDate := today.Format(pyezatypes.DateInputLayout)
			defaultISO := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, tz).Format(time.RFC3339)
			labels := buildFormLabels(deps.Labels)
			return view.OK("supplier-subscription-drawer-form", &form.Data{
				FormAction:        deps.Routes.AddURL,
				Active:            true,
				Code:              generateCode(),
				DateStartDate:     defaultDate,
				DateStartISO:      defaultISO,
				DefaultTZ:         tz.String(),
				SearchCostPlanURL: deps.Routes.SearchCostPlanURL,
				SearchSupplierURL: deps.Routes.SearchSupplierURL,
				Labels:            labels,
			})
		}

		// POST — create supplier subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		tz := pyezatypes.LocationFromContext(ctx)

		dateTimeStart := parseFormDateTime(
			r.FormValue("date_start_date"),
			r.FormValue("date_start_time"),
			r.FormValue("date_time_start_iso"),
			tz, false,
		)
		dateTimeEnd := parseFormDateTime(
			r.FormValue("date_end_date"),
			r.FormValue("date_end_time"),
			r.FormValue("date_time_end_iso"),
			tz, true,
		)

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}
		name := r.FormValue("name")
		if name == "" {
			name = code
		}
		autoRenew := r.FormValue("auto_renew") == "true" || r.FormValue("auto_renew") == "on"
		active := r.FormValue("active") != "false"

		costPlanID := r.FormValue("cost_plan_id")
		supplierID := r.FormValue("supplier_id")

		req := &suppliersubscriptionpb.CreateSupplierSubscriptionRequest{
			Data: &suppliersubscriptionpb.SupplierSubscription{
				Name:          name,
				Code:          strPtr(code),
				CostPlanId:    costPlanID,
				SupplierId:    supplierID,
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
				AutoRenew:     autoRenew,
				Active:        active,
			},
		}

		if notes := r.FormValue("notes"); notes != "" {
			if req.Data.Metadata == nil {
				req.Data.Metadata = map[string]string{}
			}
			req.Data.Metadata["notes"] = notes
		}

		if _, err := deps.CreateSupplierSubscription(ctx, req); err != nil {
			log.Printf("Failed to create supplier subscription: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}
