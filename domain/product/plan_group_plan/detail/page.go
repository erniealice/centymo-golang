package detail

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	plan_group_plan "github.com/erniealice/centymo-golang/domain/product/plan_group_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	plangroupplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group_plan"
)

// DetailViewDeps holds view dependencies for the plan_group_plan detail page.
type DetailViewDeps struct {
	Routes            plan_group_plan.Routes
	Labels            plan_group_plan.Labels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
	ReadPlanGroupPlan func(ctx context.Context, req *plangroupplanpb.ReadPlanGroupPlanRequest) (*plangroupplanpb.ReadPlanGroupPlanResponse, error)
}

// PageData holds the data for the plan_group_plan detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *plangroupplanpb.PlanGroupPlan
	Labels          plan_group_plan.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	PlanGroupID   string
	PlanID        string
	SequenceOrder string
	Status        string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string
}

// NewView creates the plan_group_plan detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan_group_plan", "read") {
			return view.Forbidden("plan_group_plan:read")
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
		return view.OK("plan-group-plan-detail", pageData)
	})
}

// NewTabAction handles GET /action/plan-group-plan/{id}/tab/{tab}.
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
		return view.OK("plan-group-plan-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPlanGroupPlan(ctx, &plangroupplanpb.ReadPlanGroupPlanRequest{
		Data: &plangroupplanpb.PlanGroupPlan{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read plan group plan %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	rec := data[0]

	l := deps.Labels

	status := "active"
	statusVariant := "success"
	if !rec.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	seqOrder := ""
	if rec.SequenceOrder != nil {
		seqOrder = strconv.FormatInt(int64(rec.GetSequenceOrder()), 10)
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

	headerTitle := rec.GetPlanGroupId()
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
			HeaderSubtitle: rec.GetPlanId(),
			HeaderIcon:     "icon-layers",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "plan-group-plan-detail-content",
		Record:          rec,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		PlanGroupID:     rec.GetPlanGroupId(),
		PlanID:          rec.GetPlanId(),
		SequenceOrder:   seqOrder,
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}
	return pageData, nil
}
