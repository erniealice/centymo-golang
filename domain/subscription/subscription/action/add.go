package action

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription/form"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// NewAddAction creates the subscription add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "create") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			q := viewCtx.Request.URL.Query()
			if deps.CreateOptions.CommencementPricing && q.Get("client_id") != "" && deps.ReadClient != nil {
				clients, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{Data: &clientpb.Client{Id: q.Get("client_id")}})
				if err != nil {
					return view.HTMXError(deps.Labels.Errors.InvalidFormData)
				}
				q.Set("billing_currency", "")
				for _, c := range clients.GetData() {
					if c.GetId() == q.Get("client_id") {
						q.Set("billing_currency", c.GetBillingCurrency())
						break
					}
				}
			}
			if deps.CreateOptions.CommencementPricing && (q.Get("term_preview") == "1" || q.Get("pricing_preview") == "1") {
				tz := pyezatypes.LocationFromContext(ctx)
				data := &form.Data{Labels: buildFormLabels(deps.Labels), FormAction: deps.Routes.AddURL, SuggestPlanTerm: true,
					SearchPlanURL: commencementSearchURL(deps.Routes.SearchPlanURL, q, tz)}
				if q.Get("pricing_preview") == "1" {
					return view.OK("subscription-pricing-fields", data)
				}
				start := parseFormDateTime(q.Get("date_start_date"), q.Get("date_start_time"), "", tz, false)
				end, _, err := suggestedPlanExpiry(ctx, deps, q.Get("price_plan_id"), start, tz)
				if err != nil {
					return view.HTMXError(deps.Labels.Errors.InvalidFormData)
				}
				if end != nil {
					data.DateEndDate = end.AsTime().In(tz).Format("2006-01-02")
				}
				return view.OK("subscription-term-end-fields", data)
			}
			clientID := q.Get("client_id")
			clientName := q.Get("client_name")
			clientBillingCurrency := q.Get("billing_currency")
			clientLocked := clientID != ""

			// 2026-05-11 — opposite-context entrypoint (price-plan detail
			// "Add Subscription" primary action). When price_plan_id is
			// present, lock the Plan picker and pre-fill its display label.
			pricePlanID := q.Get("price_plan_id")
			pricePlanLabel := q.Get("plan_label")
			pricePlanLocked := pricePlanID != ""

			tz := pyezatypes.LocationFromContext(ctx)
			// Default new engagement to "today, 00:00" in the operator's TZ.
			today := time.Now().In(tz)
			defaultDate := today.Format(pyezatypes.DateInputLayout)
			defaultISO := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, tz).Format(time.RFC3339)
			searchURL := deps.Routes.SearchPlanURL
			if deps.CreateOptions.CommencementPricing {
				q.Set("date_start_date", defaultDate)
				searchURL = commencementSearchURL(searchURL, q, tz)
			}
			labels := buildFormLabels(deps.Labels)
			// 2026-05-03 — Substitute {{.Currency}} placeholder in the
			// PlanClientScopeNotice with the client's billing currency code so
			// the banner reads "Plans below match this client's billing
			// currency (PHP)." Falls back to a no-currency variant when the
			// client has no billing currency set.
			labels.PlanClientScopeNotice = resolvePlanClientScopeNotice(labels.PlanClientScopeNotice, clientBillingCurrency)
			return view.OK("subscription-drawer-form", &form.Data{
				SuggestPlanTerm:       deps.CreateOptions.CommencementPricing,
				FormAction:            deps.Routes.AddURL,
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchPlanURL:         searchURL,
				ClientID:              clientID,
				ClientLabel:           clientName,
				ClientLocked:          clientLocked,
				ClientBillingCurrency: clientBillingCurrency,
				PricePlanID:           pricePlanID,
				PlanLabel:             pricePlanLabel,
				PricePlanLocked:       pricePlanLocked,
				Code:                  generateCode(),
				DateStartDate:         defaultDate,
				DateStartISO:          defaultISO,
				DefaultTZ:             tz.String(),
				// 2026-05-03 — pre-render dropped. The autocomplete is in action
				// mode (SearchPlanURL is set), so the inline-rendered options
				// were ignored anyway and the search endpoint already returns
				// the same client-scoped grouping on first dropdown open.
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
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
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
		if deps.CreateOptions.CommencementPricing {
			// Visible inputs remain authoritative when optional inline timezone JS is unavailable.
			dateTimeStart = parseFormDateTime(r.FormValue("date_start_date"), r.FormValue("date_start_time"), "", tz, false)
			dateTimeEnd = parseFormDateTime(r.FormValue("date_end_date"), r.FormValue("date_end_time"), "", tz, true)
			if dateTimeStart == nil || (r.FormValue("date_end_date") != "" && dateTimeEnd == nil) ||
				(dateTimeEnd != nil && dateTimeEnd.AsTime().Before(dateTimeStart.AsTime())) {
				return view.HTMXError(deps.Labels.Errors.InvalidFormData)
			}
			if err := validateCommencementPricing(ctx, deps, pricePlanID, r.FormValue("client_id"), dateTimeStart); err != nil {
				log.Printf("subscription commencement validation: %v", err)
				return view.HTMXError(deps.Labels.Errors.InvalidFormData)
			}

		}

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

		// Require-spawn-success toggle (Q-GSE-8) — sibling of the Spawn Jobs
		// toggle above, but a direct proto field on CreateSubscriptionRequest
		// (not a context override): espyna reads req.RequireSpawnSuccess
		// directly. Default unchecked -> false pointer -> byte-identical to
		// the legacy best-effort path when the operator leaves it alone.
		requireSpawnRaw := strings.ToLower(strings.TrimSpace(r.FormValue("require_spawn_success")))
		requireSpawnSuccess := requireSpawnRaw == "true" || requireSpawnRaw == "on" || requireSpawnRaw == "1" || requireSpawnRaw == "yes"

		request := &subscriptionpb.CreateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Name:          name,
				ClientId:      r.FormValue("client_id"),
				PricePlanId:   pricePlanID,
				Code:          strPtr(code),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
				Active:        true,
			},
			RequireSpawnSuccess: &requireSpawnSuccess,
		}
		if err := form.ApplyEscalation(request.Data, r.PostForm, false); err != nil {
			return view.HTMXError(deps.Labels.Form.EscalationInvalid)
		}

		resp, err := deps.CreateSubscription(spawnCtx, request)
		if err != nil {
			log.Printf("Failed to create subscription: %v", err)
			return view.HTMXError(err.Error())
		}

		_ = resp
		return view.HTMXSuccess("subscriptions-table")
	})
}
