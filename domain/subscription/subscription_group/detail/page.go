package detail

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	subscription_group "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
)

// DetailViewDeps holds view dependencies for the subscription group detail page.
type DetailViewDeps struct {
	Routes                subscription_group.Routes
	Labels                subscription_group.Labels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
	ReadSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListPlans             func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules    func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
}

// PageData holds the data for the subscription group detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Group           *subscriptiongrouppb.SubscriptionGroup
	Labels          subscription_group.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	Name          string
	Kind          string
	PlanName      string
	ScheduleName  string
	Capacity      string
	Status        string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string
}

// NewView creates the subscription group detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "read") {
			return view.Forbidden("subscription_group:read")
		}
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("subscription-group-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group/{id}/tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}
		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("subscription-group-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
		Data: &subscriptiongrouppb.SubscriptionGroup{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read subscription group %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	sg := data[0]

	l := deps.Labels

	planName := l.Detail.NoPlan
	if pid := sg.GetPlanId(); pid != "" {
		if n := lookupPlanName(ctx, deps, pid); n != "" {
			planName = n
		} else {
			planName = pid
		}
	}
	scheduleName := l.Detail.NoSchedule
	if sid := sg.GetPriceScheduleId(); sid != "" {
		if n := lookupScheduleName(ctx, deps, sid); n != "" {
			scheduleName = n
		} else {
			scheduleName = sid
		}
	}

	kind := sg.GetKind()
	if kind == "" {
		kind = l.Detail.NoKind
	}

	status := "active"
	statusVariant := "success"
	if !sg.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	tabItems := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info,
			Href:  route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=info",
			HxGet: route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "info"),
			Icon:  "icon-info"},
	}

	tz := types.LocationFromContext(ctx)
	createdDate := ""
	if ms := sg.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := sg.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	headerSubtitle := strings.TrimSpace(planName)
	if headerSubtitle == "" || planName == l.Detail.NoPlan {
		headerSubtitle = l.Detail.NoSubtitle
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          sg.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    sg.GetName(),
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-users",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "subscription-group-detail-content",
		Group:           sg,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		Name:            sg.GetName(),
		Kind:            kind,
		PlanName:        planName,
		ScheduleName:    scheduleName,
		Capacity:        formatCapacity(sg, l),
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}
	return pageData, nil
}

// formatCapacity renders the capacity summary from capacity_mode (+ max_capacity
// when CAPPED). UNSPECIFIED is treated as Unlimited per the proto contract.
func formatCapacity(sg *subscriptiongrouppb.SubscriptionGroup, l subscription_group.Labels) string {
	switch sg.GetCapacityMode() {
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_CAPPED:
		return fmt.Sprintf(l.Detail.CapacityValue, sg.GetMaxCapacity())
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_CLOSED:
		return l.Form.CapClosed
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_UNLIMITED:
		return l.Form.CapUnlimited
	default:
		return l.Detail.CapacityModeNF
	}
}

func lookupPlanName(ctx context.Context, deps *DetailViewDeps, planID string) string {
	if deps.ListPlans == nil || planID == "" {
		return ""
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return ""
	}
	for _, p := range resp.GetData() {
		if p.GetId() == planID {
			return p.GetName()
		}
	}
	return ""
}

func lookupScheduleName(ctx context.Context, deps *DetailViewDeps, scheduleID string) string {
	if deps.ListPriceSchedules == nil || scheduleID == "" {
		return ""
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		return ""
	}
	for _, ps := range resp.GetData() {
		if ps.GetId() == scheduleID {
			return ps.GetName()
		}
	}
	return ""
}
