package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes             centymo.PlanRoutes
	ReadPlan           func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	Labels             centymo.PlanLabels
	CommonLabels       pyeza.CommonLabels
	TableLabels        types.TableLabels
	ListProductPlans   func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListPricePlans     func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListLocations      func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	attachment.AttachmentOps
	auditlog.AuditOps
}

// PageData holds the data for the plan detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	Plan                *planpb.Plan
	Labels              centymo.PlanLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	ID                  string
	PlanName            string
	PlanDesc            string
	PlanStatus          string
	StatusVariant       string
	CreatedDate         string
	ModifiedDate        string
	ProductsTable       *types.TableConfig
	PriceListsTable     *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the plan detail view (full page).
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

		return view.OK("plan-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/plans/detail/{id}/tab/{tab}
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

		// Return only the tab partial template
		templateName := "plan-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads plan data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{
		Data: &planpb.Plan{Id: &id},
	})
	if err != nil {
		log.Printf("Failed to read plan %s: %v", id, err)
		return nil, fmt.Errorf("failed to load plan: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("plan not found")
	}
	plan := data[0]

	name := plan.GetName()
	description := plan.GetDescription()

	planStatus := "active"
	if !plan.GetActive() {
		planStatus = "inactive"
	}
	statusVariant := "success"
	if planStatus == "inactive" {
		statusVariant = "warning"
	}

	createdDate := plan.GetDateCreatedString()
	modifiedDate := plan.GetDateModifiedString()

	// Get counts for tab badges — filter by plan_id so only this plan's products are counted
	productCount := 0
	if deps.ListProductPlans != nil {
		ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
			Filters: &commonpb.FilterRequest{
				Logic: commonpb.FilterLogic_AND,
				Filters: []*commonpb.TypedFilter{
					{
						Field: "plan_id",
						FilterType: &commonpb.TypedFilter_StringFilter{
							StringFilter: &commonpb.StringFilter{
								Value:    id,
								Operator: commonpb.StringOperator_STRING_EQUALS,
							},
						},
					},
				},
			},
		})
		if err == nil {
			productCount = len(ppResp.GetData())
		}
	}

	priceListCount := 0
	if deps.ListPricePlans != nil {
		plResp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err == nil {
			for _, pp := range plResp.GetData() {
				if pp.GetPlanId() == id {
					priceListCount++
				}
			}
		}
	}

	l := deps.Labels
	tabItems := buildTabItems(id, l, productCount, priceListCount, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    name,
			HeaderSubtitle: description,
			HeaderIcon:     "icon-layers",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "plan-detail-content",
		Plan:            plan,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		PlanName:        name,
		PlanDesc:        description,
		PlanStatus:      planStatus,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}

	// Load tab-specific data
	switch activeTab {
	case "products":
		tableConfig := buildProductsTable(ctx, deps, id)
		pageData.ProductsTable = tableConfig
	case "pricelists":
		tableConfig := buildPriceListsTable(ctx, deps, id)
		pageData.PriceListsTable = tableConfig
	case "attachments":
		if deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments: %v", err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}
		pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
	case "audit-history":
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType:  "plan",
				EntityID:    id,
				Limit:       20,
				CursorToken: cursor,
			})
			if err != nil {
				log.Printf("Failed to load audit history: %v", err)
			}
			if auditResp != nil {
				pageData.AuditEntries = auditResp.Entries
				pageData.AuditHasNext = auditResp.HasNext
				pageData.AuditNextCursor = auditResp.NextCursor
			}
		}
		pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.PlanLabels, productCount, priceListCount int, routes centymo.PlanRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "products", Label: l.Tabs.Products, Href: base + "?tab=products", HxGet: action + "products", Icon: "icon-package", Count: productCount, Disabled: false},
		{Key: "pricelists", Label: l.Tabs.PriceLists, Href: base + "?tab=pricelists", HxGet: action + "pricelists", Icon: "icon-tag", Count: priceListCount, Disabled: false},
		{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		{Key: "audit", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock", Count: 0, Disabled: false},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// ---------------------------------------------------------------------------
// Products tab table
// ---------------------------------------------------------------------------

func buildProductsTable(ctx context.Context, deps *DetailViewDeps, planID string) *types.TableConfig {
	l := deps.Labels
	perms := view.GetUserPermissions(ctx)

	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.Product, Sortable: true},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
	}

	rows := []types.TableRow{}

	if deps.ListProductPlans != nil {
		ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
			Filters: &commonpb.FilterRequest{
				Logic: commonpb.FilterLogic_AND,
				Filters: []*commonpb.TypedFilter{
					{
						Field: "plan_id",
						FilterType: &commonpb.TypedFilter_StringFilter{
							StringFilter: &commonpb.StringFilter{
								Value:    planID,
								Operator: commonpb.StringOperator_STRING_EQUALS,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("Failed to list product plans: %v", err)
		} else {
			for _, pp := range ppResp.GetData() {
				ppID := pp.GetId()
				name := pp.GetName()

				active := pp.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				rowActions := []types.TableAction{
					{
						Type:        "edit",
						Label:       l.Actions.Edit,
						Action:      "edit",
						URL:         route.ResolveURL(deps.Routes.ProductPlanEditURL, "id", planID, "ppid", ppID),
						DrawerTitle: l.Actions.Edit,
					},
					{
						Type:     "delete",
						Label:    l.Actions.Delete,
						Action:   "delete",
						URL:      route.ResolveURL(deps.Routes.ProductPlanDeleteURL, "id", planID),
						ItemName: name,
					},
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						{Type: "badge", Value: status, Variant: statusVariant(status)},
					},
					Actions: rowActions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	refreshURL := route.ResolveURL(deps.Routes.TabActionURL, "id", planID, "tab", "") + "products"

	tableConfig := &types.TableConfig{
		ID:                   "plan-products-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.NoProductsAssigned,
			Message: l.Detail.NoProductsAssignedMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddProduct,
			ActionURL:       route.ResolveURL(deps.Routes.ProductPlanAddURL, "id", planID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("product_plan", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Price Lists tab table
// ---------------------------------------------------------------------------

func buildPriceListsTable(ctx context.Context, deps *DetailViewDeps, planID string) *types.TableConfig {
	l := deps.Labels
	perms := view.GetUserPermissions(ctx)

	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.PricePlan, Sortable: true},
		{Key: "amount", Label: l.Detail.Price, Sortable: true, WidthClass: "col-4xl"},
		{Key: "duration", Label: l.Columns.Duration, Sortable: true, WidthClass: "col-4xl"},
		{Key: "schedule", Label: "Schedule", Sortable: true, WidthClass: "col-6xl"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
	}

	// Build a schedule ID → name map for display.
	scheduleNames := map[string]string{}
	if deps.ListPriceSchedules != nil {
		schedResp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
		if err != nil {
			log.Printf("Failed to list price schedules for pricelists table: %v", err)
		} else {
			for _, s := range schedResp.GetData() {
				scheduleNames[s.GetId()] = s.GetName()
			}
		}
	}

	rows := []types.TableRow{}

	if deps.ListPricePlans != nil {
		plResp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list price plans: %v", err)
		} else {
			for _, pp := range plResp.GetData() {
				if pp.GetPlanId() != planID {
					continue
				}

				ppID := pp.GetId()
				name := pp.GetName()
				ppCurrency := pp.GetCurrency()
				if ppCurrency == "" {
					ppCurrency = "PHP"
				}
				amountCell := types.MoneyCell(float64(pp.GetAmount()), ppCurrency, true)
				duration := fmt.Sprintf("%d %s", pp.GetDurationValue(), pp.GetDurationUnit())

				scheduleName := "—"
				if schedID := pp.GetPriceScheduleId(); schedID != "" {
					if n, ok := scheduleNames[schedID]; ok && n != "" {
						scheduleName = n
					} else {
						scheduleName = schedID
					}
				}

				active := pp.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				rowActions := []types.TableAction{
					{
						Type:        "edit",
						Label:       l.Actions.Edit,
						Action:      "edit",
						URL:         route.ResolveURL(deps.Routes.PricePlanEditURL, "id", planID, "ppid", ppID),
						DrawerTitle: l.Actions.Edit,
					},
					{
						Type:    "delete",
						Label:   l.Actions.Delete,
						Action:  "delete",
						URL:     route.ResolveURL(deps.Routes.PricePlanDeleteURL, "id", planID),
						ItemName: name,
					},
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						amountCell,
						{Type: "text", Value: duration},
						{Type: "text", Value: scheduleName},
						{Type: "badge", Value: status, Variant: statusVariant(status)},
					},
					Actions: rowActions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "plan-pricelists-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.NoPricePlans,
			Message: l.Detail.NoPricePlansMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddPricePlan,
			ActionURL:       route.ResolveURL(deps.Routes.PricePlanAddURL, "id", planID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_plan", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}
