package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	subscription_group_member "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
)

// DetailViewDeps holds view dependencies for the subscription_group_member detail page.
type DetailViewDeps struct {
	Routes                      subscription_group_member.Routes
	Labels                      subscription_group_member.Labels
	CommonLabels                pyeza.CommonLabels
	TableLabels                 types.TableLabels
	ReadSubscriptionGroupMember func(ctx context.Context, req *subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.ReadSubscriptionGroupMemberResponse, error)
}

// PageData holds the data for the subscription_group_member detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Member          *subscriptiongroupmemberpb.SubscriptionGroupMember
	Labels          subscription_group_member.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID                  string
	SubscriptionGroupId string
	SubscriptionId      string
	ClientId            string
	Status              string
	StatusVariant       string
	CreatedDate         string
	ModifiedDate        string
}

// NewView creates the subscription_group_member detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "read") {
			return view.Forbidden("subscription_group_member:read")
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
		return view.OK("subscription-group-member-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group-member/{id}/tab/{tab}.
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
		return view.OK("subscription-group-member-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroupMember(ctx, &subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest{
		Data: &subscriptiongroupmemberpb.SubscriptionGroupMember{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read subscription group member %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	m := data[0]

	l := deps.Labels

	status := "active"
	statusVariant := "success"
	if !m.GetActive() {
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
	if ms := m.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := m.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	headerTitle := m.GetId()
	if g := m.GetSubscriptionGroupId(); g != "" {
		headerTitle = g
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          l.Detail.Title,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    headerTitle,
			HeaderSubtitle: m.GetSubscriptionId(),
			HeaderIcon:     "icon-user",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:     "subscription-group-member-detail-content",
		Member:              m,
		Labels:              l,
		ActiveTab:           activeTab,
		TabItems:            tabItems,
		ID:                  id,
		SubscriptionGroupId: m.GetSubscriptionGroupId(),
		SubscriptionId:      m.GetSubscriptionId(),
		ClientId:            m.GetClientId(),
		Status:              status,
		StatusVariant:       statusVariant,
		CreatedDate:         createdDate,
		ModifiedDate:        modifiedDate,
	}
	return pageData, nil
}
