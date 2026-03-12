package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Deps holds view dependencies.
type Deps struct {
	Routes           centymo.SubscriptionRoutes
	ReadSubscription func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	Labels           centymo.SubscriptionLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels
}

// PageData holds the data for the subscription detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Subscription    map[string]any
	Labels          centymo.SubscriptionLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
}

// subscriptionToMap converts a Subscription protobuf to a map[string]any for template use.
func subscriptionToMap(s *subscriptionpb.Subscription) map[string]any {
	// Build customer display name: prefer company_name, fallback to user name
	customer := s.GetName()
	if c := s.GetClient(); c != nil {
		if companyName := c.GetCompanyName(); companyName != "" {
			customer = companyName
		} else if u := c.GetUser(); u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			if first != "" || last != "" {
				customer = first + " " + last
			}
		}
	}

	// Get plan name from nested price_plan → plan
	planName := ""
	if pp := s.GetPricePlan(); pp != nil {
		if p := pp.GetPlan(); p != nil {
			planName = p.GetName()
		}
		if planName == "" {
			planName = pp.GetName()
		}
	}

	status := "active"
	if !s.GetActive() {
		status = "inactive"
	}

	return map[string]any{
		"id":                   s.GetId(),
		"name":                 s.GetName(),
		"customer":             customer,
		"plan":                 planName,
		"price_plan_id":        s.GetPricePlanId(),
		"client_id":            s.GetClientId(),
		"date_start_string":    s.GetDateStartString(),
		"date_end_string":      s.GetDateEndString(),
		"status":               status,
		"active":               s.GetActive(),
		"date_created_string":  s.GetDateCreatedString(),
		"date_modified_string": s.GetDateModifiedString(),
		"quantity":             s.GetQuantity(),
		"assigned_count":       s.GetAssignedCount(),
		"available_count":      s.GetAvailableCount(),
	}
}

// NewView creates the subscription detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
			Data: &subscriptionpb.Subscription{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read subscription %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load subscription: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Subscription %s not found", id)
			return view.Error(fmt.Errorf("subscription not found"))
		}
		subscription := subscriptionToMap(data[0])

		subName, _ := subscription["name"].(string)
		customer, _ := subscription["customer"].(string)
		headerTitle := subName
		if customer != "" {
			headerTitle = customer
		}

		l := deps.Labels

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "services",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "subscription-detail-content",
			Subscription:    subscription,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
		}

		return view.OK("subscription-detail", pageData)
	})
}

func buildTabItems(l centymo.SubscriptionLabels, id string, routes centymo.SubscriptionRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "audit", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
			Data: &subscriptionpb.Subscription{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read subscription %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load subscription: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Subscription %s not found", id)
			return view.Error(fmt.Errorf("subscription not found"))
		}
		subscription := subscriptionToMap(data[0])

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Subscription: subscription,
			Labels:       l,
			ActiveTab:    tab,
			TabItems:     buildTabItems(l, id, deps.Routes),
		}

		templateName := "subscription-tab-" + tab
		return view.OK(templateName, pageData)
	})
}
