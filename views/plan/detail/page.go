package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes           centymo.PlanRoutes
	ReadPlan         func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	Labels           centymo.PlanLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels
	ListProductPlans func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListPricePlans   func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListLocations    func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)

	attachment.AttachmentOps
	auditlog.AuditOps
}

// PageData holds the data for the plan detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Plan            *planpb.Plan
	Labels          centymo.PlanLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	ID              string
	PlanName        string
	PlanDesc        string
	FulfillmentType string
	PlanStatus      string
	StatusVariant   string
	CreatedDate     string
	ModifiedDate    string
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
	fulfillmentType := plan.GetFulfillmentType()
	if fulfillmentType == "" {
		fulfillmentType = "—"
	}

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
		FulfillmentType: fulfillmentType,
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

	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.Product, Sortable: true},
		{Key: "price", Label: l.Detail.Price, Sortable: true, Width: "150px"},
		{Key: "currency", Label: l.Detail.Currency, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
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
				price := fmt.Sprintf("%.2f", pp.GetPrice())
				currency := pp.GetCurrency()
				if currency == "" {
					currency = "PHP"
				}

				active := pp.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						{Type: "text", Value: price},
						{Type: "text", Value: currency},
						{Type: "badge", Value: status, Variant: statusVariant(status)},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "plan-products-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          false,
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
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Price Lists tab table
// ---------------------------------------------------------------------------

func buildPriceListsTable(ctx context.Context, deps *DetailViewDeps, planID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.PricePlan, Sortable: true},
		{Key: "amount", Label: l.Detail.Price, Sortable: true, Width: "150px"},
		{Key: "duration", Label: l.Columns.Duration, Sortable: true, Width: "150px"},
		{Key: "location", Label: l.Columns.Location, Sortable: true, Width: "180px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}

	// Build a location ID → name map for display.
	locationNames := map[string]string{}
	if deps.ListLocations != nil {
		locResp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
		if err != nil {
			log.Printf("Failed to list locations for pricelists table: %v", err)
		} else {
			for _, loc := range locResp.GetData() {
				locationNames[loc.GetId()] = loc.GetName()
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
				amount := fmt.Sprintf("%.2f", pp.GetAmount())
				duration := fmt.Sprintf("%d %s", pp.GetDurationValue(), pp.GetDurationUnit())

				locationName := "—"
				if locID := pp.GetLocationId(); locID != "" {
					if n, ok := locationNames[locID]; ok && n != "" {
						locationName = n
					} else {
						locationName = locID
					}
				}

				active := pp.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						{Type: "text", Value: amount},
						{Type: "text", Value: duration},
						{Type: "text", Value: locationName},
						{Type: "badge", Value: status, Variant: statusVariant(status)},
					},
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
		ShowActions:          false,
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
