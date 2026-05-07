package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/cost_plan/form"
	"github.com/erniealice/pyeza-golang/view"

	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
)

// NewAddAction creates the cost_plan add action.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("cost_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("cost-plan-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				Active:                true,
				SearchSupplierPlanURL: deps.SearchSupplierPlanURL,
				SearchCostScheduleURL: deps.SearchCostScheduleURL,
				Labels:                buildFormLabels(deps.Labels),
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		name := r.FormValue("name")
		supplierPlanID := r.FormValue("supplier_plan_id")
		costScheduleID := r.FormValue("cost_schedule_id")
		billingKind := r.FormValue("billing_kind")
		amountBasis := r.FormValue("amount_basis")
		amount := parseAmount(r.FormValue("amount"))
		currency := r.FormValue("currency")
		billingCycleValue := r.FormValue("billing_cycle_value")
		billingCycleUnit := r.FormValue("billing_cycle_unit")
		defaultTermValue := r.FormValue("default_term_value")
		defaultTermUnit := r.FormValue("default_term_unit")
		description := r.FormValue("description")
		active := r.FormValue("active") != "false"

		cp := &costplanpb.CostPlan{
			SupplierPlanId:  supplierPlanID,
			BillingAmount:   amount,
			BillingCurrency: currency,
			Active:          active,
		}
		if name != "" {
			cp.Name = strPtr(name)
		}
		if description != "" {
			cp.Description = strPtr(description)
		}
		if costScheduleID != "" {
			cp.CostScheduleId = strPtr(costScheduleID)
		}
		if billingKind != "" {
			if bk, ok := costplanpb.CostPlanBillingKind_value[billingKind]; ok {
				cp.BillingKind = costplanpb.CostPlanBillingKind(bk)
			}
		}
		if amountBasis != "" {
			if ab, ok := costplanpb.CostPlanAmountBasis_value[amountBasis]; ok {
				cp.AmountBasis = costplanpb.CostPlanAmountBasis(ab)
			}
		}
		if v, err := strconv.ParseInt(billingCycleValue, 10, 32); err == nil {
			v32 := int32(v)
			cp.BillingCycleValue = &v32
		}
		if billingCycleUnit != "" {
			cp.BillingCycleUnit = strPtr(billingCycleUnit)
		}
		if v, err := strconv.ParseInt(defaultTermValue, 10, 32); err == nil {
			v32 := int32(v)
			cp.DefaultTermValue = &v32
		}
		if defaultTermUnit != "" {
			cp.DefaultTermUnit = strPtr(defaultTermUnit)
		}

		if _, err := deps.CreateCostPlan(ctx, &costplanpb.CreateCostPlanRequest{Data: cp}); err != nil {
			log.Printf("Failed to create cost plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("cost-plans-table")
	})
}
