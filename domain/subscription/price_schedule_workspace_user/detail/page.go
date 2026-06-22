package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	pswu "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	pswupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule_workspace_user"
)

// DetailViewDeps holds view dependencies for the price_schedule_workspace_user detail page.
type DetailViewDeps struct {
	Routes                         pswu.Routes
	Labels                         pswu.Labels
	CommonLabels                   pyeza.CommonLabels
	TableLabels                    types.TableLabels
	ReadPriceScheduleWorkspaceUser func(ctx context.Context, req *pswupb.ReadPriceScheduleWorkspaceUserRequest) (*pswupb.ReadPriceScheduleWorkspaceUserResponse, error)
}

// PageData holds the data for the price_schedule_workspace_user detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *pswupb.PriceScheduleWorkspaceUser
	Labels          pswu.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID              string
	PriceScheduleId string
	WorkspaceUserId string
	Scope           string
	Role            string
	IsOwner         bool
	Status          string
	StatusVariant   string
	CreatedDate     string
	ModifiedDate    string
}

// NewView creates the price_schedule_workspace_user detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "read") {
			return view.Forbidden("price_schedule_workspace_user:read")
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
		return view.OK("price-schedule-workspace-user-detail", pageData)
	})
}

// NewTabAction handles GET /action/price-schedule-workspace-user/{id}/tab/{tab}.
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
		return view.OK("price-schedule-workspace-user-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPriceScheduleWorkspaceUser(ctx, &pswupb.ReadPriceScheduleWorkspaceUserRequest{
		Data: &pswupb.PriceScheduleWorkspaceUser{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read price_schedule_workspace_user %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	record := data[0]

	l := deps.Labels

	status := "active"
	statusVariant := "success"
	if !record.GetActive() {
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
	if ms := record.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := record.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	headerTitle := record.GetPriceScheduleId()
	if headerTitle == "" {
		headerTitle = id
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          headerTitle,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    headerTitle,
			HeaderSubtitle: record.GetWorkspaceUserId(),
			HeaderIcon:     "icon-users",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "price-schedule-workspace-user-detail-content",
		Record:          record,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		PriceScheduleId: record.GetPriceScheduleId(),
		WorkspaceUserId: record.GetWorkspaceUserId(),
		Scope:           record.GetScope(),
		Role:            record.GetRole(),
		IsOwner:         record.GetIsOwner(),
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}
	return pageData, nil
}
