package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	plan_group "github.com/erniealice/centymo-golang/domain/product/plan_group"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	plangroupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group"
)

// DetailViewDeps holds view dependencies for the plan group detail page.
type DetailViewDeps struct {
	Routes        plan_group.Routes
	Labels        plan_group.Labels
	CommonLabels  pyeza.CommonLabels
	TableLabels   types.TableLabels
	ReadPlanGroup func(ctx context.Context, req *plangroupb.ReadPlanGroupRequest) (*plangroupb.ReadPlanGroupResponse, error)
	// ListPlanGroups is used to resolve the parent group name.
	ListPlanGroups func(ctx context.Context, req *plangroupb.ListPlanGroupsRequest) (*plangroupb.ListPlanGroupsResponse, error)
}

// PageData holds the data for the plan group detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Group           *plangroupb.PlanGroup
	Labels          plan_group.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	Name          string
	Code          string
	ParentName    string
	Status        string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string
}

// NewView creates the plan group detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group", "read") {
			return view.Forbidden("plan_group:read")
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
		return view.OK("plan-group-detail", pageData)
	})
}

// NewTabAction handles GET /action/plan-group/{id}/tab/{tab}.
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
		return view.OK("plan-group-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPlanGroup(ctx, &plangroupb.ReadPlanGroupRequest{
		Data: &plangroupb.PlanGroup{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read plan group %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	pg := data[0]

	l := deps.Labels

	parentName := l.Detail.NoParent
	if pid := pg.GetParentId(); pid != "" {
		if n := lookupParentName(ctx, deps, pid); n != "" {
			parentName = n
		} else {
			parentName = pid
		}
	}

	code := pg.GetCode()
	if code == "" {
		code = l.Detail.NoCode
	}

	status := "active"
	statusVariant := "success"
	if !pg.GetActive() {
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
	if ms := pg.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := pg.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	headerSubtitle := parentName
	if headerSubtitle == "" || parentName == l.Detail.NoParent {
		headerSubtitle = l.Detail.NoSubtitle
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          pg.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    pg.GetName(),
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-layers",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "plan-group-detail-content",
		Group:           pg,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		Name:            pg.GetName(),
		Code:            code,
		ParentName:      parentName,
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}
	return pageData, nil
}

func lookupParentName(ctx context.Context, deps *DetailViewDeps, parentID string) string {
	if deps.ListPlanGroups == nil || parentID == "" {
		return ""
	}
	resp, err := deps.ListPlanGroups(ctx, &plangroupb.ListPlanGroupsRequest{})
	if err != nil {
		return ""
	}
	for _, pg := range resp.GetData() {
		if pg.GetId() == parentID {
			return pg.GetName()
		}
	}
	return ""
}
