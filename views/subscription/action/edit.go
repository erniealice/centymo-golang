package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/subscription/form"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// NewEditAction creates the subscription edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Prefer the joined item-page-data path so Client (+ User) is populated
			// without a second ListClients-and-iterate roundtrip. Falls back to
			// ReadSubscription only if the dep is unwired.
			var record *subscriptionpb.Subscription
			if deps.GetSubscriptionItemPageData != nil {
				resp, err := deps.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{
					SubscriptionId: id,
				})
				if err != nil || resp == nil || resp.GetSubscription() == nil {
					log.Printf("Failed to read subscription %s: %v", id, err)
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSubscription()
			} else {
				readResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
					Data: &subscriptionpb.Subscription{Id: id},
				})
				if err != nil {
					log.Printf("Failed to read subscription %s: %v", id, err)
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				readData := readResp.GetData()
				if len(readData) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = readData[0]
			}

			// Prefer the joined client (populated by GetSubscriptionItemPageData);
			// fall back to the ListClients lookup for the legacy ReadSubscription path.
			clientLabel := ""
			if c := record.GetClient(); c != nil {
				if name := c.GetName(); name != "" {
					clientLabel = name
				} else if u := c.GetUser(); u != nil {
					clientLabel = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
				}
			}
			if clientLabel == "" {
				clientLabel = resolveClientLabel(ctx, record.GetClientId(), deps.ListClients)
			}
			clientBillingCurrency := ""
			if c := record.GetClient(); c != nil {
				clientBillingCurrency = c.GetBillingCurrency()
			}
			if clientBillingCurrency == "" {
				clientBillingCurrency = resolveClientBillingCurrency(ctx, record.GetClientId(), deps.ListClients)
			}
			// PricePlanID, not a plan_id — resolve via PricePlan so the selected
			// label matches the autocomplete dropdown's display.
			planLabel := resolvePricePlanName(ctx, record.GetPricePlanId(), deps)

			// Lock client field when opened from client detail page
			clientLocked := viewCtx.Request.URL.Query().Get("client_id") != ""

			tz := pyezatypes.LocationFromContext(ctx)
			startDate, startTime, startISO := splitTimestampForInputs(record.GetDateTimeStart(), tz)
			endDate, endTime, endISO := splitTimestampForInputs(record.GetDateTimeEnd(), tz)

			labels := buildFormLabels(deps.Labels)
			return view.OK("subscription-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Code:                  record.GetCode(),
				ClientID:              record.GetClientId(),
				PricePlanID:           record.GetPricePlanId(),
				DateStartDate:         startDate,
				DateStartTime:         startTime,
				DateStartISO:          startISO,
				DateEndDate:           endDate,
				DateEndTime:           endTime,
				DateEndISO:            endISO,
				DefaultTZ:             tz.String(),
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchPlanURL:         deps.Routes.SearchPlanURL,
				ClientLabel:           clientLabel,
				ClientLocked:          clientLocked,
				ClientBillingCurrency: clientBillingCurrency,
				PlanLabel:             planLabel,
				PlanOptionGroups:      form.LoadPricePlanOptionGroups(ctx, deps.ListPricePlans, deps.ListPriceSchedules, record.GetClientId(), clientLabel, labels),
				// Edit drawer never spawns Jobs — the toggle hides because
				// SpawnJobsAvailable defaults false. Operators trigger
				// retroactive spawn via the Operations tab CTA.
				SpawnJobsAvailable:  false,
				SpawnJobsPartialURL: deps.Routes.SpawnJobsPartialURL,
				Labels:              labels,
				CommonLabels:        nil, // injected by ViewAdapter
			})
		}

		// POST — update subscription
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
		if pricePlanID == "" {
			pricePlanID = r.FormValue("plan_id")
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

		_, err := deps.UpdateSubscription(ctx, &subscriptionpb.UpdateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Id:            id,
				Name:          name,
				ClientId:      r.FormValue("client_id"),
				PricePlanId:   pricePlanID,
				Code:          strPtr(code),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
			},
		})
		if err != nil {
			log.Printf("Failed to update subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}
