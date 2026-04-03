package action

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// FormLabels holds i18n labels for the subscription drawer form template.
type FormLabels struct {
	Customer                  string
	CustomerPlaceholder       string
	Plan                      string
	PlanPlaceholder           string
	StartDate                 string
	EndDate                   string
	Notes                     string
	NotesPlaceholder          string
	CustomerSearchPlaceholder string
	PlanSearchPlaceholder     string
	CustomerNoResults         string
	PlanNoResults             string
	Code                      string
	CodePlaceholder           string
}

// FormData is the template data for the subscription drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Code         string
	ClientID     string
	PricePlanID  string
	DateStart    string
	DateEnd      string
	Notes        string
	Clients         []map[string]string
	PricePlans      []map[string]string
	SearchClientURL string
	SearchPlanURL   string
	ClientLabel     string
	PlanLabel       string
	ClientLocked    bool
	Labels          FormLabels
	CommonLabels    any
}

// Deps holds dependencies for subscription action handlers.
type Deps struct {
	Routes centymo.SubscriptionRoutes
	Labels centymo.SubscriptionLabels

	CreateSubscription  func(ctx context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error)
	ReadSubscription    func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	UpdateSubscription  func(ctx context.Context, req *subscriptionpb.UpdateSubscriptionRequest) (*subscriptionpb.UpdateSubscriptionResponse, error)
	DeleteSubscription  func(ctx context.Context, req *subscriptionpb.DeleteSubscriptionRequest) (*subscriptionpb.DeleteSubscriptionResponse, error)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListPlans           func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
	SearchPlansByName   func(ctx context.Context, req *planpb.SearchPlansByNameRequest) (*planpb.SearchPlansByNameResponse, error)
	ListPricePlans      func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
}

func formLabels(l centymo.SubscriptionLabels) FormLabels {
	return FormLabels{
		Customer:                  l.Form.Customer,
		CustomerPlaceholder:       l.Form.CustomerPlaceholder,
		Plan:                      l.Form.Plan,
		PlanPlaceholder:           l.Form.PlanPlaceholder,
		StartDate:                 l.Form.StartDate,
		EndDate:                   l.Form.EndDate,
		Notes:                     l.Form.Notes,
		NotesPlaceholder:          l.Form.NotesPlaceholder,
		CustomerSearchPlaceholder: l.Form.CustomerSearchPlaceholder,
		PlanSearchPlaceholder:     l.Form.PlanSearchPlaceholder,
		CustomerNoResults:         l.Form.CustomerNoResults,
		PlanNoResults:             l.Form.PlanNoResults,
		Code:                      l.Form.Code,
		CodePlaceholder:           l.Form.CodePlaceholder,
	}
}

// generateCode returns a random 7-character uppercase alphanumeric code,
// using chars that are visually unambiguous (no O, I, 0, 1).
func generateCode() string {
	const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// loadClientOptions fetches the client list and converts to select options.
func loadClientOptions(ctx context.Context, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) []map[string]string {
	if listClients == nil {
		return nil
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		log.Printf("Failed to load clients for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, c := range resp.GetData() {
		label := c.GetId()
		if u := c.GetUser(); u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			if first != "" || last != "" {
				label = first + " " + last
			}
		}
		options = append(options, map[string]string{
			"Value": c.GetId(),
			"Label": label,
		})
	}
	return options
}

// loadPlanOptions fetches the plan list and converts to select options.
func loadPlanOptions(ctx context.Context, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) []map[string]string {
	if listPlans == nil {
		return nil
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("Failed to load plans for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		options = append(options, map[string]string{
			"Value": p.GetId(),
			"Label": p.GetName(),
		})
	}
	return options
}

// resolveClientLabel finds the display name for a client by ID.
func resolveClientLabel(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			if cn := c.GetName(); cn != "" {
				return cn
			}
			if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					return strings.TrimSpace(first + " " + last)
				}
			}
			return clientID
		}
	}
	return clientID
}

// resolvePlanLabel finds the display name for a plan by ID.
func resolvePlanLabel(ctx context.Context, planID string, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) string {
	if planID == "" || listPlans == nil {
		return ""
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return planID
	}
	for _, p := range resp.GetData() {
		if p.GetId() == planID {
			return p.GetName()
		}
	}
	return planID
}

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
			clientLocked := clientID != ""

			return view.OK("subscription-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				SearchClientURL: deps.Routes.SearchClientURL,
				SearchPlanURL:   deps.Routes.SearchPlanURL,
				ClientID:        clientID,
				ClientLabel:     clientName,
				ClientLocked:    clientLocked,
				Code:            generateCode(),
				Labels:          formLabels(deps.Labels),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — create subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		dateStart := r.FormValue("date_start_string")
		dateEnd := r.FormValue("date_end_string")

		pricePlanID := r.FormValue("price_plan_id")

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}

		// Resolve plan name for auto-generated subscription name
		planName := ""
		if deps.ListPlans != nil && pricePlanID != "" {
			planName = resolvePlanLabel(ctx, pricePlanID, deps.ListPlans)
		}
		name := planName
		if code != "" {
			name = planName + " [" + code + "]"
		}

		// Validate: start date must not be in the past
		if dateStart != "" {
			startTime, parseErr := time.Parse("2006-01-02", dateStart)
			if parseErr == nil {
				today := time.Now().Truncate(24 * time.Hour)
				if startTime.Before(today) {
					return centymo.HTMXError("Start date cannot be in the past")
				}
			}
		}

		resp, err := deps.CreateSubscription(ctx, &subscriptionpb.CreateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Name:            name,
				ClientId:        r.FormValue("client_id"),
				PricePlanId:     pricePlanID,
				Code:            strPtr(code),
				DateStartString: strPtr(dateStart),
				DateEndString:   strPtr(dateEnd),
				Active:          true,
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

// NewEditAction creates the subscription edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
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
			record := readData[0]

			clientLabel := resolveClientLabel(ctx, record.GetClientId(), deps.ListClients)
			planLabel := resolvePlanLabel(ctx, record.GetPricePlanId(), deps.ListPlans)

			// Lock client field when opened from client detail page
			clientLocked := viewCtx.Request.URL.Query().Get("client_id") != ""

			return view.OK("subscription-drawer-form", &FormData{
				FormAction:      route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:          true,
				ID:              id,
				Code:            record.GetCode(),
				ClientID:        record.GetClientId(),
				PricePlanID:     record.GetPricePlanId(),
				DateStart:       record.GetDateStartString(),
				DateEnd:         record.GetDateEndString(),
				SearchClientURL: deps.Routes.SearchClientURL,
				SearchPlanURL:   deps.Routes.SearchPlanURL,
				ClientLabel:     clientLabel,
				ClientLocked:    clientLocked,
				PlanLabel:       planLabel,
				Labels:          formLabels(deps.Labels),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — update subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		dateStart := r.FormValue("date_start_string")
		dateEnd := r.FormValue("date_end_string")

		pricePlanID := r.FormValue("price_plan_id")
		if pricePlanID == "" {
			pricePlanID = r.FormValue("plan_id")
		}

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}

		// Resolve plan name for auto-generated subscription name
		planName := ""
		if deps.ListPlans != nil && pricePlanID != "" {
			planName = resolvePlanLabel(ctx, pricePlanID, deps.ListPlans)
		}
		name := planName
		if code != "" {
			name = planName + " [" + code + "]"
		}

		// Validate: start date must not be in the past
		if dateStart != "" {
			startTime, parseErr := time.Parse("2006-01-02", dateStart)
			if parseErr == nil {
				today := time.Now().Truncate(24 * time.Hour)
				if startTime.Before(today) {
					return centymo.HTMXError("Start date cannot be in the past")
				}
			}
		}

		_, err := deps.UpdateSubscription(ctx, &subscriptionpb.UpdateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Id:              id,
				Name:            name,
				ClientId:        r.FormValue("client_id"),
				PricePlanId:     pricePlanID,
				Code:            strPtr(code),
				DateStartString: strPtr(dateStart),
				DateEndString:   strPtr(dateEnd),
			},
		})
		if err != nil {
			log.Printf("Failed to update subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// NewDeleteAction creates the subscription delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "delete") {
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

		_, err := deps.DeleteSubscription(ctx, &subscriptionpb.DeleteSubscriptionRequest{
			Data: &subscriptionpb.Subscription{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}
