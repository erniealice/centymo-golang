// Package customize handles the "Customize Package for Client" feature.
// This is a POST-only flow — no own drawer template.
// 2026-04-27 plan-client-scope plan §6.5.
package customize

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

// Request mirrors the espyna use-case request shape (plan §4.1).
type Request struct {
	SourcePlanID      string
	SourcePricePlanID string
	ClientID          string
	SubscriptionID    string
	NewScheduleName   string
}

// Response mirrors the espyna use-case response shape (plan §4.1).
type Response struct {
	NewPlanID      string
	NewPricePlanID string
	NewScheduleID  string
	Reused         bool
}

// Deps is the dependency subset needed by the customize feature.
type Deps struct {
	Labels centymo.SubscriptionLabels

	// CustomClientPriceScheduleLabelSuffix carries the lyngua-resolved suffix
	// appended to a client's name when constructing the custom PriceSchedule name.
	CustomClientPriceScheduleLabelSuffix string

	CustomizePlanForClient      func(ctx context.Context, req *Request) (*Response, error)
	GetSubscriptionItemPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
	ReadSubscription            func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	ReadPricePlan               func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListClients                 func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
}

// NewAction creates the POST handler for "Customize this package for {ClientName}".
// Algorithm (handler-side):
//  1. Read subscription → resolve price_plan_id, client_id, client name.
//  2. Read customClientPriceScheduleLabelSuffix from typed labels.
//  3. Build derivedName = "{ClientName} - {suffix}".
//  4. Call espyna CustomizePlanForClient with the source IDs + derivedName.
//  5. On success, respond with HX-Push-Url pointing at the new PricePlan's
//     package page and HX-Trigger refresh-package.
func NewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("revenue", "create") && !perms.Can("plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

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

		sub, err := loadSubscription(ctx, deps, subscriptionID)
		if err != nil {
			log.Printf("Customize: failed to load subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		clientID := sub.GetClientId()
		pricePlanID := sub.GetPricePlanId()
		if clientID == "" || pricePlanID == "" {
			return centymo.HTMXError(deps.Labels.Errors.CustomizeFailed)
		}

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

		clientName := resolveClientName(ctx, deps, sub, clientID)
		suffix := deps.CustomClientPriceScheduleLabelSuffix
		derivedName := buildScheduleName(clientName, suffix)

		req := &Request{
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

		newURL := "/app/clients/detail/" + clientID + "/engagements/" + subscriptionID + "/package/" + resp.NewPricePlanID
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Push-Url": newURL,
				"HX-Trigger":  "refresh-package",
			},
		}
	})
}

// loadSubscription fetches the subscription with joined client + price plan.
func loadSubscription(ctx context.Context, deps *Deps, id string) (*subscriptionpb.Subscription, error) {
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

// resolveClientName resolves a display name for the client.
func resolveClientName(ctx context.Context, deps *Deps, sub *subscriptionpb.Subscription, clientID string) string {
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

// buildScheduleName mirrors plan §4.4.1: "{Client.name} - {suffix}".
func buildScheduleName(clientName, suffix string) string {
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

var errNotFound = customizeError("subscription not found")

type customizeError string

func (e customizeError) Error() string { return string(e) }
