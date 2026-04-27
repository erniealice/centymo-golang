package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// NewCustomizePackageAction creates the POST handler that drives the
// "Customize this package for {ClientName}" CTA on the subscription detail's
// Package tab (plan §6.5 / §4.4.1).
//
// Algorithm (handler-side):
//   1. Read subscription → resolve price_plan_id, client_id, client name.
//   2. Read customClientPriceScheduleLabelSuffix from typed labels.
//   3. Build derivedName = "{ClientName} - {suffix}" (note: space-hyphen-space).
//   4. Call espyna CustomizePlanForClient with the source IDs + derivedName.
//   5. On success, respond with HX-Push-Url pointing at the new PricePlan's
//      package page and HX-Trigger refresh-package.
//
// Cross-package contract: Deps.CustomizePlanForClient is a function pointer
// matching the espyna use case's signature. The block wires it via the
// CustomizePlanForClientRequest/Response shape declared in action.go.
//
// 2026-04-27 plan-client-scope plan §6.5.
func NewCustomizePackageAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		// Authz: same gate the use case checks server-side. Surface the same
		// permission key set the espyna side validates.
		if !perms.Can("revenue", "create") && !perms.Can("plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		// POST only — the customize CTA is a single-button flow, no GET drawer.
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if deps.CustomizePlanForClient == nil {
			return centymo.HTMXError(deps.Labels.Errors.CustomizeFailed)
		}

		subscriptionID := viewCtx.Request.PathValue("id")
		if subscriptionID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		// Step 1 — load the subscription with its joined client + price plan.
		sub, err := loadSubscriptionForCustomize(ctx, deps, subscriptionID)
		if err != nil {
			log.Printf("Customize: failed to load subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		clientID := sub.GetClientId()
		pricePlanID := sub.GetPricePlanId()
		if clientID == "" || pricePlanID == "" {
			return centymo.HTMXError(deps.Labels.Errors.CustomizeFailed)
		}

		// Resolve source plan_id from the live PricePlan record (the join may
		// not include plan_id directly).
		sourcePlanID := ""
		if pp := sub.GetPricePlan(); pp != nil {
			sourcePlanID = pp.GetPlanId()
		}
		if sourcePlanID == "" && deps.ReadPricePlan != nil {
			if r, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
				Data: &priceplanpb.PricePlan{Id: pricePlanID},
			}); err == nil && len(r.GetData()) > 0 {
				sourcePlanID = r.GetData()[0].GetPlanId()
			}
		}
		if sourcePlanID == "" {
			return centymo.HTMXError(deps.Labels.Errors.CustomizeFailed)
		}

		// Step 1 (cont) — resolve client name. Fallback chain mirrors
		// resolveClientBreadcrumb in subscription/detail/page.go.
		clientName := resolveClientNameForCustomize(ctx, deps, sub, clientID)

		// Step 2 — read suffix from typed labels (plan §4.4.1 step 2).
		// The lyngua-resolved suffix is threaded through Deps by block.go so
		// the professional-tier "Rate Cards" override reaches this handler
		// without relying on the X-Client-Schedule-Suffix header pattern.
		suffix := deps.CustomClientPriceScheduleLabelSuffix

		// Step 3 — build derivedName per plan §4.4.1.
		// Note: space-hyphen-space separator ("{Client.name} - {suffix}").
		derivedName := buildDerivedScheduleName(clientName, suffix)

		// Step 4 — invoke the espyna use case.
		req := &CustomizePlanForClientRequest{
			SourcePlanID:      sourcePlanID,
			SourcePricePlanID: pricePlanID,
			ClientID:          clientID,
			SubscriptionID:    subscriptionID,
			NewScheduleName:   derivedName,
		}
		resp, err := deps.CustomizePlanForClient(ctx, req)
		if err != nil {
			log.Printf("Customize: espyna use case failed for sub %s: %v", subscriptionID, err)
			return centymo.HTMXError(err.Error())
		}
		if resp == nil || resp.NewPricePlanID == "" {
			return centymo.HTMXError(deps.Labels.Errors.CustomizeFailed)
		}

		// Step 5 — HX-redirect to the new package page + trigger.
		newURL := buildPackageURLForClient(clientID, subscriptionID, resp.NewPricePlanID)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Push-Url": newURL,
				"HX-Trigger":  "refresh-package",
			},
		}
	})
}

// loadSubscriptionForCustomize fetches the subscription with its joined
// client + price plan. Prefers GetSubscriptionItemPageData; falls back to
// ReadSubscription.
func loadSubscriptionForCustomize(ctx context.Context, deps *Deps, id string) (*subscriptionpb.Subscription, error) {
	if deps.GetSubscriptionItemPageData != nil {
		resp, err := deps.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{
			SubscriptionId: id,
		})
		if err != nil {
			return nil, err
		}
		if resp != nil && resp.GetSubscription() != nil {
			return resp.GetSubscription(), nil
		}
	}
	resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: id},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, errNotFound
	}
	return resp.GetData()[0], nil
}

// resolveClientNameForCustomize resolves a display name for the client,
// preferring the joined client on the subscription, then falling back to
// ListClients lookup, and finally the bare client_id.
func resolveClientNameForCustomize(ctx context.Context, deps *Deps, sub *subscriptionpb.Subscription, clientID string) string {
	if c := sub.GetClient(); c != nil {
		if name := c.GetName(); name != "" {
			return name
		}
		if u := c.GetUser(); u != nil {
			full := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			if full != "" {
				return full
			}
		}
	}
	if deps.ListClients != nil {
		resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{})
		if err == nil {
			for _, c := range resp.GetData() {
				if c.GetId() != clientID {
					continue
				}
				if name := c.GetName(); name != "" {
					return name
				}
				if u := c.GetUser(); u != nil {
					full := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
					if full != "" {
						return full
					}
				}
			}
		}
	}
	return clientID
}

// buildDerivedScheduleName mirrors plan §4.4.1: "{Client.name} - {suffix}".
// Note the space-hyphen-space separator. When the suffix is empty, falls
// back to just the client name (no trailing dash).
func buildDerivedScheduleName(clientName, suffix string) string {
	clientName = strings.TrimSpace(clientName)
	suffix = strings.TrimSpace(suffix)
	if clientName == "" && suffix == "" {
		return ""
	}
	if suffix == "" {
		return clientName
	}
	if clientName == "" {
		return suffix
	}
	return clientName + " - " + suffix
}

// buildPackageURLForClient builds the post-customize redirect target per
// plan §6.5: /app/clients/detail/{cid}/engagements/{sid}/package/{newPpid}.
//
// We don't pull the URL from a route map because that path is owned by
// entydad's client-detail page (plan §6.3) — at this seam we just construct
// the canonical convention. If/when the URL diverges, plumb it through Deps.
func buildPackageURLForClient(clientID, subscriptionID, pricePlanID string) string {
	return "/app/clients/detail/" + clientID + "/engagements/" + subscriptionID + "/package/" + pricePlanID
}

// errNotFound is the sentinel returned from loadSubscriptionForCustomize
// when the subscription read returns an empty data slice. Wrapped at the
// call site into the lyngua "not found" error.
var errNotFound = customizeError("subscription not found")

type customizeError string

func (e customizeError) Error() string { return string(e) }
