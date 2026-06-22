package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	line_workspace_user "github.com/erniealice/centymo-golang/domain/product/line_workspace_user"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	lineworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line_workspace_user"
)

// DetailViewDeps holds view dependencies for the line_workspace_user detail page.
type DetailViewDeps struct {
	Routes                line_workspace_user.Routes
	Labels                line_workspace_user.Labels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
	ReadLineWorkspaceUser func(ctx context.Context, req *lineworkspaceuserpb.ReadLineWorkspaceUserRequest) (*lineworkspaceuserpb.ReadLineWorkspaceUserResponse, error)
}

// PageData holds the data for the line_workspace_user detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *lineworkspaceuserpb.LineWorkspaceUser
	Labels          line_workspace_user.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID              string
	WorkspaceUserId string
	LineId          string
	Scope           string
	Role            string
	IsOwner         bool
	Status          string
	StatusVariant   string
	CreatedDate     string
	ModifiedDate    string
}

// NewView creates the line_workspace_user detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("line_workspace_user", "read") {
			return view.Forbidden("line_workspace_user:read")
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
		return view.OK("line-workspace-user-detail", pageData)
	})
}

// NewTabAction handles GET /action/line-workspace-user/{id}/tab/{tab}.
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
		return view.OK("line-workspace-user-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadLineWorkspaceUser(ctx, &lineworkspaceuserpb.ReadLineWorkspaceUserRequest{
		Data: &lineworkspaceuserpb.LineWorkspaceUser{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read line_workspace_user %s: %v", id, err)
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

	headerTitle := record.GetWorkspaceUserId()
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
			HeaderSubtitle: record.GetLineId(),
			HeaderIcon:     "icon-user",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "line-workspace-user-detail-content",
		Record:          record,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		WorkspaceUserId: record.GetWorkspaceUserId(),
		LineId:          record.GetLineId(),
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
