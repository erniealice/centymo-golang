package detail

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

// DetailViewDeps holds view dependencies for the price schedule detail page.
type DetailViewDeps struct {
	Routes       centymo.PriceScheduleRoutes
	Labels       centymo.PriceScheduleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadPriceSchedule func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	ListLocations     func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPricePlans    func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPlans         func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	CreatePricePlan   func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)

	// Reference checker: returns a map of price_plan_id → true for plans in use by active subscriptions.
	// Delete is disabled for in-use plans; Edit remains enabled (Pricing fields lock inside the drawer).
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PlanFormData is rendered as the drawer form for adding a PricePlan under this schedule.
type PlanFormData struct {
	FormAction   string
	ScheduleID   string
	ScheduleName string
	PlanOptions  []map[string]any
	CommonLabels pyeza.CommonLabels
	Labels       centymo.PriceScheduleLabels
}

// PageData holds the data for the price schedule detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Schedule        *priceschedulepb.PriceSchedule
	Labels          centymo.PriceScheduleLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID             string
	Name           string
	Description    string
	DateStart      string
	DateEnd        string
	LocationName   string
	Status         string
	StatusVariant  string
	CreatedDate    string
	ModifiedDate   string

	PlansTable *types.TableConfig
}

// NewView creates the price schedule detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("price-schedule-detail", pageData)
	})
}

// NewTabAction handles GET /action/price-schedule/{id}/tab/{tab}.
// Returns only the tab partial template, for HTMX tab switching.
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

		return view.OK("price-schedule-tab-"+tab, pageData)
	})
}

// NewPlanAddAction handles GET/POST /action/price-schedule/{id}/plan/add.
// GET renders a drawer with price_schedule_id pre-locked; POST creates the PricePlan.
func NewPlanAddAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if deps.CreatePricePlan == nil {
			return centymo.HTMXError("Price plan create is not available")
		}
		scheduleID := viewCtx.Request.PathValue("id")
		if scheduleID == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		if viewCtx.Request.Method == http.MethodGet {
			scheduleName := lookupScheduleName(ctx, deps, scheduleID)
			planOptions := buildPlanOptions(ctx, deps)
			return view.OK("price-schedule-plan-drawer-form", &PlanFormData{
				FormAction:   route.ResolveURL(deps.Routes.PlanAddURL, "id", scheduleID),
				ScheduleID:   scheduleID,
				ScheduleName: scheduleName,
				PlanOptions:  planOptions,
				CommonLabels: deps.CommonLabels,
				Labels:       deps.Labels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		r := viewCtx.Request
		planID := r.FormValue("plan_id")
		if planID == "" {
			return centymo.HTMXError(deps.Labels.Detail.PlanRequired)
		}
		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		currency := r.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}

		ppName := r.FormValue("name")
		ppDesc := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			PlanId:        planID,
			Name:          &ppName,
			Description:   &ppDesc,
			Amount:        amount,
			Currency:      currency,
			DurationValue: int32(dv),
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        true,
		}
		pp.PriceScheduleId = &scheduleID

		if _, err := deps.CreatePricePlan(ctx, &priceplanpb.CreatePricePlanRequest{Data: pp}); err != nil {
			log.Printf("Failed to create price plan from schedule %s: %v", scheduleID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plans-table")
	})
}

func lookupScheduleName(ctx context.Context, deps *DetailViewDeps, scheduleID string) string {
	if deps.ReadPriceSchedule == nil {
		return scheduleID
	}
	resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
		Data: &priceschedulepb.PriceSchedule{Id: scheduleID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return scheduleID
	}
	if name := resp.GetData()[0].GetName(); name != "" {
		return name
	}
	return scheduleID
}

func buildPlanOptions(ctx context.Context, deps *DetailViewDeps) []map[string]any {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	opts := make([]map[string]any, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		if p == nil || !p.GetActive() {
			continue
		}
		opts = append(opts, map[string]any{
			"Value":       p.GetId(),
			"Label":       p.GetName(),
			"Description": p.GetDescription(),
		})
	}
	return opts
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
		Data: &priceschedulepb.PriceSchedule{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read price schedule %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	ps := data[0]

	locationName := deps.Labels.Detail.NoLocation
	if locID := ps.GetLocationId(); locID != "" {
		if n := lookupLocationName(ctx, deps, locID); n != "" {
			locationName = n
		} else {
			locationName = locID
		}
	}

	status := "active"
	statusVariant := "success"
	if !ps.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	l := deps.Labels
	planCount := countPlansForSchedule(ctx, deps, ps)
	tabItems := buildTabItems(id, l, planCount, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          ps.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    ps.GetName(),
			HeaderSubtitle: ps.GetDescription(),
			HeaderIcon:     "icon-calendar",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "price-schedule-detail-content",
		Schedule:        ps,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		Name:            ps.GetName(),
		Description:     ps.GetDescription(),
		DateStart:       ps.GetDateStart(),
		DateEnd:         ps.GetDateEnd(),
		LocationName:    locationName,
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     ps.GetDateCreatedString(),
		ModifiedDate:    ps.GetDateModifiedString(),
	}

	if activeTab == "plans" {
		pageData.PlansTable = buildPlansTable(ctx, deps, ps)
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.PriceScheduleLabels, planCount int, routes centymo.PriceScheduleRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "plans", Label: l.Tabs.Plans, Href: base + "?tab=plans", HxGet: action + "plans", Icon: "icon-layers", Count: planCount},
	}
}

// countPlansForSchedule counts price_plans linked to this schedule via price_schedule_id FK.
func countPlansForSchedule(ctx context.Context, deps *DetailViewDeps, ps *priceschedulepb.PriceSchedule) int {
	if deps.ListPricePlans == nil {
		return 0
	}
	resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
	if err != nil {
		return 0
	}
	schedID := ps.GetId()
	count := 0
	for _, pp := range resp.GetData() {
		if pp != nil && pp.GetPriceScheduleId() == schedID {
			count++
		}
	}
	return count
}

func buildPlansTable(ctx context.Context, deps *DetailViewDeps, ps *priceschedulepb.PriceSchedule) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "name", Label: l.Detail.PlanColumnPlan, Sortable: true},
		{Key: "amount", Label: l.Detail.PlanColumnAmount, Sortable: true, WidthClass: "col-4xl", Align: "right"},
		{Key: "duration", Label: l.Detail.PlanColumnDuration, Sortable: false, WidthClass: "col-3xl"},
		{Key: "status", Label: l.Detail.PlanColumnStatus, Sortable: false, WidthClass: "col-2xl"},
	}

	perms := view.GetUserPermissions(ctx)
	rows := []types.TableRow{}
	if deps.ListPricePlans != nil {
		resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list price plans for schedule %s: %v", ps.GetId(), err)
		} else {
			schedID := ps.GetId()

			// Collect IDs for the reference checker (one batch call for the whole table).
			var ppIDs []string
			for _, pp := range resp.GetData() {
				if pp != nil && pp.GetPriceScheduleId() == schedID {
					ppIDs = append(ppIDs, pp.GetId())
				}
			}
			inUseIDs := map[string]bool{}
			if deps.GetPricePlanInUseIDs != nil && len(ppIDs) > 0 {
				inUseIDs, _ = deps.GetPricePlanInUseIDs(ctx, ppIDs)
			}

			for _, pp := range resp.GetData() {
				if pp == nil || pp.GetPriceScheduleId() != schedID {
					continue
				}
				ppID := pp.GetId()
				name := pp.GetName()
				currency := pp.GetCurrency()
				if currency == "" {
					currency = "PHP"
				}
				duration := ""
				if dv := pp.GetDurationValue(); dv > 0 {
					duration = fmt.Sprintf("%d %s", dv, pp.GetDurationUnit())
				}
				planStatus := "active"
				planVariant := "success"
				if !pp.GetActive() {
					planStatus = "inactive"
					planVariant = "warning"
				}

				inUse := inUseIDs[ppID]
				deleteDisabled := !perms.Can("price_plan", "delete") || inUse
				deleteTooltip := l.Errors.Unauthorized
				if inUse {
					deleteTooltip = l.Detail.PlanInUseTooltip
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						types.MoneyCell(float64(pp.GetAmount()), currency, true),
						{Type: "text", Value: duration},
						{Type: "badge", Value: planStatus, Variant: planVariant},
					},
					Actions: []types.TableAction{
						{
							Type:   "view",
							Label:  l.Detail.PlanView,
							Action: "view",
							Href:   route.ResolveURL(deps.Routes.PlanDetailURL, "id", schedID, "ppid", ppID),
						},
						{
							Type:            "edit",
							Label:           l.Detail.PlanEdit,
							Action:          "edit",
							URL:             route.ResolveURL(deps.Routes.PlanEditURL, "id", schedID, "ppid", ppID),
							DrawerTitle:     l.Detail.PlanEditDrawerTitle,
							Disabled:        !perms.Can("price_plan", "update"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
						{
							Type:            "delete",
							Label:           l.Detail.PlanDelete,
							Action:          "delete",
							URL:             route.ResolveURL(deps.Routes.PlanDeleteURL, "id", schedID, "ppid", ppID),
							ItemName:        name,
							ConfirmTitle:    l.Detail.PlanDeleteTitle,
							ConfirmMessage:  fmt.Sprintf(l.Detail.PlanDeleteMsg, name),
							Disabled:        deleteDisabled,
							DisabledTooltip: deleteTooltip,
						},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)


	refreshURL := route.ResolveURL(deps.Routes.TabActionURL, "id", ps.GetId(), "tab", "plans")
	tableConfig := &types.TableConfig{
		ID:                   "price-schedule-plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.PlansEmptyTitle,
			Message: l.Detail.PlansEmptyMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Detail.PlanAdd,
			ActionURL:       route.ResolveURL(deps.Routes.PlanAddURL, "id", ps.GetId()),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_plan", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func lookupLocationName(ctx context.Context, deps *DetailViewDeps, locationID string) string {
	if deps.ListLocations == nil || locationID == "" {
		return ""
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		return ""
	}
	for _, loc := range resp.GetData() {
		if loc.GetId() == locationID {
			return loc.GetName()
		}
	}
	return ""
}
