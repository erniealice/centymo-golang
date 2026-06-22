package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	sgwu "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	sgwupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_workspace_user"
)

// DetailViewDeps holds view dependencies for the detail page.
type DetailViewDeps struct {
	Routes                             sgwu.Routes
	Labels                             sgwu.Labels
	CommonLabels                       pyeza.CommonLabels
	TableLabels                        types.TableLabels
	ReadSubscriptionGroupWorkspaceUser func(ctx context.Context, req *sgwupb.ReadSubscriptionGroupWorkspaceUserRequest) (*sgwupb.ReadSubscriptionGroupWorkspaceUserResponse, error)
}

// PageData holds the data for the subscription_group_workspace_user detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *sgwupb.SubscriptionGroupWorkspaceUser
	Labels          sgwu.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID                  string
	WorkspaceUserId     string
	SubscriptionGroupId string
	Scope               string
	Role                string
	IsOwner             string
	Status              string
	StatusVariant       string
	CreatedDate         string
	ModifiedDate        string
}

// NewView creates the subscription_group_workspace_user detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "read") {
			return view.Forbidden("subscription_group_workspace_user:read")
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
		return view.OK("subscription-group-workspace-user-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group-workspace-user/{id}/tab/{tab}.
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
		return view.OK("subscription-group-workspace-user-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroupWorkspaceUser(ctx, &sgwupb.ReadSubscriptionGroupWorkspaceUserRequest{
		Data: &sgwupb.SubscriptionGroupWorkspaceUser{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read subscription_group_workspace_user %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	rec := data[0]

	l := deps.Labels

	scopeVal := rec.GetScope()
	if scopeVal == "" {
		scopeVal = l.Detail.NoScope
	}
	roleVal := rec.GetRole()
	if roleVal == "" {
		roleVal = l.Detail.NoRole
	}
	isOwnerLabel := l.Detail.OwnerNo
	if rec.GetIsOwner() {
		isOwnerLabel = l.Detail.OwnerYes
	}

	status := "active"
	statusVariant := "success"
	if !rec.GetActive() {
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
	if ms := rec.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := rec.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	// Use workspace_user_id as the page header title.
	headerTitle := rec.GetWorkspaceUserId()
	if headerTitle == "" {
		headerTitle = id
	}
	headerSubtitle := rec.GetSubscriptionGroupId()
	if headerSubtitle == "" {
		headerSubtitle = l.Detail.NoGroup
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          headerTitle,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    headerTitle,
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-users",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:     "subscription-group-workspace-user-detail-content",
		Record:              rec,
		Labels:              l,
		ActiveTab:           activeTab,
		TabItems:            tabItems,
		ID:                  id,
		WorkspaceUserId:     rec.GetWorkspaceUserId(),
		SubscriptionGroupId: rec.GetSubscriptionGroupId(),
		Scope:               scopeVal,
		Role:                roleVal,
		IsOwner:             isOwnerLabel,
		Status:              status,
		StatusVariant:       statusVariant,
		CreatedDate:         createdDate,
		ModifiedDate:        modifiedDate,
	}
	return pageData, nil
}
