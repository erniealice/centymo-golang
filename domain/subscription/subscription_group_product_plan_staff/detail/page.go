package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	sgpps "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// DetailViewDeps holds view dependencies for the
// subscription_group_product_plan_staff detail page.
type DetailViewDeps struct {
	Routes                                sgpps.Routes
	Labels                                sgpps.Labels
	CommonLabels                          pyeza.CommonLabels
	TableLabels                           types.TableLabels
	ReadSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.ReadSubscriptionGroupProductPlanStaffRequest) (*sgppspb.ReadSubscriptionGroupProductPlanStaffResponse, error)
}

// PageData holds the data for the detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *sgppspb.SubscriptionGroupProductPlanStaff
	Labels          sgpps.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID                  string
	SubscriptionGroupID string
	ProductPlanID       string
	StaffID             string
	Role                string
	Status              string
	StatusVariant       string
	CreatedDate         string
	ModifiedDate        string
}

// NewView creates the subscription_group_product_plan_staff detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_product_plan_staff", "read") {
			return view.Forbidden("subscription_group_product_plan_staff:read")
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
		return view.OK("sgpps-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group-product-plan-staff/{id}/tab/{tab}.
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
		return view.OK("sgpps-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroupProductPlanStaff(ctx, &sgppspb.ReadSubscriptionGroupProductPlanStaffRequest{
		Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read subscription_group_product_plan_staff %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	rec := data[0]

	l := deps.Labels

	role := rec.GetRole()
	if role == "" {
		role = l.Detail.NoRole
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

	// Use staff_id as the page header title (most identifying field)
	headerTitle := rec.GetStaffId()
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
		ContentTemplate:     "sgpps-detail-content",
		Record:              rec,
		Labels:              l,
		ActiveTab:           activeTab,
		TabItems:            tabItems,
		ID:                  id,
		SubscriptionGroupID: rec.GetSubscriptionGroupId(),
		ProductPlanID:       rec.GetProductPlanId(),
		StaffID:             rec.GetStaffId(),
		Role:                role,
		Status:              status,
		StatusVariant:       statusVariant,
		CreatedDate:         createdDate,
		ModifiedDate:        modifiedDate,
	}
	return pageData, nil
}
