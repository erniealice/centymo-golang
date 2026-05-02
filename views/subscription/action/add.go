package action

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/erniealice/pyeza-golang/view"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/subscription/form"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// NewAddAction creates the subscription add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			clientID := viewCtx.Request.URL.Query().Get("client_id")
			clientName := viewCtx.Request.URL.Query().Get("client_name")
			clientBillingCurrency := viewCtx.Request.URL.Query().Get("billing_currency")
			clientLocked := clientID != ""

			tz := pyezatypes.LocationFromContext(ctx)
			// Default new engagement to "today, 00:00" in the operator's TZ.
			today := time.Now().In(tz)
			defaultDate := today.Format(pyezatypes.DateInputLayout)
			defaultISO := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, tz).Format(time.RFC3339)
			labels := buildFormLabels(deps.Labels)
			return view.OK("subscription-drawer-form", &form.Data{
				FormAction:            deps.Routes.AddURL,
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchPlanURL:         deps.Routes.SearchPlanURL,
				ClientID:              clientID,
				ClientLabel:           clientName,
				ClientLocked:          clientLocked,
				ClientBillingCurrency: clientBillingCurrency,
				Code:                  generateCode(),
				DateStartDate:         defaultDate,
				DateStartISO:          defaultISO,
				DefaultTZ:             tz.String(),
				PlanOptionGroups:      form.LoadPricePlanOptionGroups(ctx, deps.ListPricePlans, deps.ListPriceSchedules, clientID, clientName, labels),
				// Spawn Jobs section starts hidden on add (no PricePlan
				// selected yet); the HTMX partial fills it after selection.
				SpawnJobsAvailable:  false,
				SpawnJobsDefault:    true,
				SpawnJobsPartialURL: deps.Routes.SpawnJobsPartialURL,
				Labels:              labels,
				CommonLabels:        nil, // injected by ViewAdapter
			})
		}

		// POST — create subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		tz := pyezatypes.LocationFromContext(ctx)
		dateTimeStart := parseFormDateTime(
			r.FormValue("date_start_date"),
			r.FormValue("date_start_time"),
			r.FormValue("date_time_start_iso"),
			tz,
			false,
		)
		dateTimeEnd := parseFormDateTime(
			r.FormValue("date_end_date"),
			r.FormValue("date_end_time"),
			r.FormValue("date_time_end_iso"),
			tz,
			true,
		)

		pricePlanID := r.FormValue("price_plan_id")

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}

		// Resolve plan name for auto-generated subscription name. The drawer
		// submits a price_plan_id, so look up the PricePlan (not the Plan).
		planName := resolvePricePlanName(ctx, pricePlanID, deps)
		name := planName
		if code != "" {
			name = planName + " [" + code + "]"
		}

		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 — propagate the
		// operator's "Spawn Jobs on Create" toggle decision through context so
		// CreateSubscriptionUseCase → JobTemplateInstantiator can honor opt-out.
		//
		// Tri-state via the `spawn_jobs_field_present` hidden marker emitted
		// by the form template only when the section was rendered (i.e. the
		// selected Plan resolved to a JobTemplate):
		//
		//   - marker absent           → section not rendered → don't override
		//                                 (espyna falls back to its legacy
		//                                  default-on, which will short-circuit
		//                                  with no_template_found anyway).
		//   - marker present + spawn_jobs truthy → operator opted in.
		//   - marker present + spawn_jobs absent  → operator unchecked the box.
		//
		// Plain-string key mirrors espyna's exported constant
		// SpawnJobsOverrideKey (espyna's internal pkg is not importable from
		// centymo, but a string-keyed context value crosses the module boundary
		// the same way "businessType" already does).
		spawnCtx := ctx
		if r.Form.Get("spawn_jobs_field_present") != "" {
			rawVal := strings.ToLower(strings.TrimSpace(r.FormValue("spawn_jobs")))
			val := rawVal == "true" || rawVal == "on" || rawVal == "1" || rawVal == "yes"
			v := val
			spawnCtx = context.WithValue(ctx, "spawn_jobs_override", &v)
		}

		resp, err := deps.CreateSubscription(spawnCtx, &subscriptionpb.CreateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Name:          name,
				ClientId:      r.FormValue("client_id"),
				PricePlanId:   pricePlanID,
				Code:          strPtr(code),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
				Active:        true,
			},
		})
		if err != nil {
			log.Printf("Failed to create subscription: %v", err)
			return centymo.HTMXError(err.Error())
		}

		_ = resp
		return centymo.HTMXSuccess("subscriptions-table")
	})
}
