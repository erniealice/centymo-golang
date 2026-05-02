package dashboard

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
)

// Deps holds view dependencies.
type Deps struct {
	Routes                     centymo.InventoryRoutes
	ListInventoryItems         func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ListInventorySerials       func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)
	ListInventoryTransactions  func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	ListInventoryDepreciations func(ctx context.Context, req *inventorydepreciationpb.ListInventoryDepreciationsRequest) (*inventorydepreciationpb.ListInventoryDepreciationsResponse, error)
	Labels                     centymo.InventoryLabels
	CommonLabels               pyeza.CommonLabels
}

// PageData holds the data for the inventory dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// aggregates is the projected aggregate set used by both the dashboard view
// and the legacy HTMX partial actions. Phase 1 refactor (2026-05-02) keeps
// the source data flowing through the existing list-based aggregation; only
// the rendering surface changed.
type aggregates struct {
	totalItems      int
	totalStockValue float64
	lowStockCount   int
	categoryCount   map[string]int
	serialAvailable int
	totalSerials    int
	totalCostBasis  int64
	totalBookValue  int64
}

// NewView creates the inventory dashboard view.
//
// Phase 1 refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// The same in-memory aggregation that the legacy widget grid used now feeds
// typed Stats / Widgets / QuickActions on DashboardData. No new aggregate
// methods — this is refactor-only.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard
		ag := loadAggregates(ctx, deps)

		// Categorical bar chart: items per tracking-mode bucket.
		var (
			categoryLabels []string
			categoryValues []float64
		)
		for k, v := range ag.categoryCount {
			categoryLabels = append(categoryLabels, k)
			categoryValues = append(categoryValues, float64(v))
		}
		if len(categoryLabels) == 0 {
			// Render an empty bar with one bucket so the chart still draws axes.
			categoryLabels = []string{"-"}
			categoryValues = []float64{0}
		}
		categoryChart := &types.ChartData{
			Labels: categoryLabels,
			Series: []types.ChartSeries{{
				Name:   l.CategoryDistribution,
				Values: categoryValues,
				Color:  "amber",
			}},
		}
		categoryChart.AutoScale()

		// Recent activity list, derived from the most recent transactions.
		recentItems := buildRecentTransactionsList(ctx, deps)

		// Low-stock alert list, derived from items below reorder level.
		alertItems := buildLowStockAlertsList(ctx, deps)

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewItem, Href: deps.Routes.AddURL, Variant: "primary", TestID: "inventory-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "inventory-action-list"},
				{Icon: "icon-activity", Label: l.QuickMovements, Href: deps.Routes.MovementsURL, TestID: "inventory-action-movements"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-dollar-sign", Value: fmt.Sprintf("%.0f", ag.totalStockValue), Label: l.TotalStockValue, Color: "terracotta", TestID: "inventory-stat-value"},
				{Icon: "icon-alert-triangle", Value: fmt.Sprintf("%d", ag.lowStockCount), Label: l.LowStockAlerts, Color: "amber", TestID: "inventory-stat-low-stock"},
				{Icon: "icon-repeat", Value: fmt.Sprintf("%d", ag.totalItems), Label: l.StockTurnover, Color: "sage", TestID: "inventory-stat-items"},
				{Icon: "icon-map-pin", Value: fmt.Sprintf("%d", len(centymo.LocationMap)), Label: l.ItemsByLocation, Color: "navy", TestID: "inventory-stat-locations"},
				{Icon: "icon-trending-down", Value: fmt.Sprintf("%.0f / %.0f", float64(ag.totalCostBasis)/100.0, float64(ag.totalBookValue)/100.0), Label: l.DepreciationSummary, Color: "terracotta", TestID: "inventory-stat-depreciation"},
				{Icon: "icon-hash", Value: fmt.Sprintf("%d / %d", ag.serialAvailable, ag.totalSerials), Label: l.SerialUnitStatus, Color: "sage", TestID: "inventory-stat-serials"},
				{Icon: "icon-pie-chart", Value: fmt.Sprintf("%d %s", len(ag.categoryCount), l.TypesUnit), Label: l.CategoryDistribution, Color: "amber", TestID: "inventory-stat-categories"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "categories", Title: l.CategoryDistribution,
					Type: "chart", ChartKind: "bar",
					ChartData: categoryChart, Span: 2,
				},
				{
					ID: "recent", Title: l.RecentActivity, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.MovementsURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-activity",
						Title: l.RecentActivity,
						Desc:  l.RecentMovements,
					},
				},
				{
					ID: "alerts", Title: l.LowStockAlerts, Type: "list", Span: 2,
					ListItems: alertItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-alert-triangle",
						Title: l.LowStockAlerts,
						Desc:  l.LowStockAlerts,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "inventory",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-briefcase",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "inventory-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("inventory-dashboard", pageData)
	})
}

// NewDashboardStatsAction returns stats widget partials via HTMX (legacy route).
func NewDashboardStatsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return view.OK("inventory-dashboard-stats", map[string]any{
			"Labels": deps.Labels,
		})
	})
}

// NewDashboardChartAction returns the chart widget partial via HTMX (legacy route).
func NewDashboardChartAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return view.OK("inventory-dashboard-chart", map[string]any{
			"Labels": deps.Labels,
		})
	})
}

// NewDashboardMovementsAction returns the recent movements widget partial (legacy route).
func NewDashboardMovementsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		resp, err := deps.ListInventoryTransactions(ctx, &inventorytransactionpb.ListInventoryTransactionsRequest{})
		if err != nil {
			log.Printf("Failed to list transactions for dashboard: %v", err)
		}

		var txns []*inventorytransactionpb.InventoryTransaction
		if resp != nil {
			txns = resp.GetData()
		}

		// Take last 10 (most recent)
		if len(txns) > 10 {
			txns = txns[len(txns)-10:]
		}

		// Convert to map[string]any for template backward compat
		txnMaps := make([]map[string]any, 0, len(txns))
		for _, t := range txns {
			txnMaps = append(txnMaps, map[string]any{
				"id":                t.GetId(),
				"transaction_type":  t.GetTransactionType(),
				"quantity":          t.GetQuantity(),
				"transaction_date":  t.GetTransactionDateString(),
				"serial_number":     t.GetSerialNumber(),
				"performed_by":      t.GetPerformedBy(),
				"inventory_item_id": t.GetInventoryItemId(),
			})
		}

		return view.OK("inventory-dashboard-movements", map[string]any{
			"Transactions": txnMaps,
			"Labels":       deps.Labels,
		})
	})
}

// NewDashboardAlertsAction returns the low stock alerts widget partial (legacy route).
func NewDashboardAlertsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
		if err != nil {
			log.Printf("Failed to list items for alerts: %v", err)
		}

		var items []*inventoryitempb.InventoryItem
		if resp != nil {
			items = resp.GetData()
		}

		alerts := []map[string]any{}
		for _, item := range items {
			onHand := item.GetQuantityOnHand()
			reserved := item.GetQuantityReserved()
			reorderLvl := item.GetReorderLevel()
			available := onHand - reserved
			if reorderLvl > 0 && available <= reorderLvl {
				alerts = append(alerts, map[string]any{
					"id":                item.GetId(),
					"name":              item.GetName(),
					"sku":               item.GetSku(),
					"quantity_on_hand":  onHand,
					"quantity_reserved": reserved,
					"reorder_level":     reorderLvl,
					"tracking_mode":     item.GetProduct().GetTrackingMode(),
					"location_id":       item.GetLocationId(),
				})
			}
		}

		return view.OK("inventory-dashboard-alerts", map[string]any{
			"Alerts": alerts,
			"Labels": deps.Labels,
		})
	})
}

func loadAggregates(ctx context.Context, deps *Deps) aggregates {
	ag := aggregates{categoryCount: map[string]int{}}

	resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list inventory items: %v", err)
	}
	var items []*inventoryitempb.InventoryItem
	if resp != nil {
		items = resp.GetData()
	}

	ag.totalItems = len(items)
	for _, item := range items {
		onHand := item.GetQuantityOnHand()
		reserved := item.GetQuantityReserved()
		reorderLvl := item.GetReorderLevel()
		available := onHand - reserved

		// Approximate stock value (quantity * a base price); preserved from
		// the legacy widget builder.
		ag.totalStockValue += onHand * 100

		if reorderLvl > 0 && available <= reorderLvl {
			ag.lowStockCount++
		}

		trackingMode := item.GetProduct().GetTrackingMode()
		if trackingMode == "" {
			trackingMode = "bulk"
		}
		ag.categoryCount[trackingMode]++
	}

	serialResp, err := deps.ListInventorySerials(ctx, &inventoryserialpb.ListInventorySerialsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list serials: %v", err)
	}
	var serials []*inventoryserialpb.InventorySerial
	if serialResp != nil {
		serials = serialResp.GetData()
	}
	ag.totalSerials = len(serials)
	for _, s := range serials {
		if s.GetStatus() == "available" {
			ag.serialAvailable++
		}
	}

	depResp, err := deps.ListInventoryDepreciations(ctx, &inventorydepreciationpb.ListInventoryDepreciationsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list depreciations: %v", err)
	}
	var depreciations []*inventorydepreciationpb.InventoryDepreciation
	if depResp != nil {
		depreciations = depResp.GetData()
	}
	for _, d := range depreciations {
		ag.totalCostBasis += d.GetCostBasis()
		ag.totalBookValue += d.GetBookValue()
	}

	return ag
}

func buildRecentTransactionsList(ctx context.Context, deps *Deps) []types.ActivityItem {
	resp, err := deps.ListInventoryTransactions(ctx, &inventorytransactionpb.ListInventoryTransactionsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list transactions: %v", err)
	}
	var txns []*inventorytransactionpb.InventoryTransaction
	if resp != nil {
		txns = resp.GetData()
	}
	if len(txns) > 5 {
		txns = txns[len(txns)-5:]
	}

	items := make([]types.ActivityItem, 0, len(txns))
	for i, t := range txns {
		items = append(items, types.ActivityItem{
			IconName:    iconForTransactionType(t.GetTransactionType()),
			IconVariant: variantForTransactionType(t.GetTransactionType()),
			Title:       t.GetTransactionType(),
			Description: fmt.Sprintf("%s — %s", t.GetSerialNumber(), t.GetPerformedBy()),
			Time:        t.GetTransactionDateString(),
			TestID:      fmt.Sprintf("inventory-activity-%d", i),
		})
	}
	return items
}

func buildLowStockAlertsList(ctx context.Context, deps *Deps) []types.ActivityItem {
	resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list items for alerts: %v", err)
	}
	var items []*inventoryitempb.InventoryItem
	if resp != nil {
		items = resp.GetData()
	}

	alerts := make([]types.ActivityItem, 0, 5)
	for _, item := range items {
		if len(alerts) >= 5 {
			break
		}
		onHand := item.GetQuantityOnHand()
		reserved := item.GetQuantityReserved()
		reorderLvl := item.GetReorderLevel()
		available := onHand - reserved
		if reorderLvl > 0 && available <= reorderLvl {
			alerts = append(alerts, types.ActivityItem{
				IconName:    "icon-alert-triangle",
				IconVariant: "quote",
				Title:       item.GetName(),
				Description: fmt.Sprintf("SKU %s — %.0f available (reorder %.0f)", item.GetSku(), available, reorderLvl),
				TestID:      fmt.Sprintf("inventory-alert-%s", item.GetId()),
			})
		}
	}
	return alerts
}

func iconForTransactionType(kind string) string {
	switch kind {
	case "received", "in":
		return "icon-plus"
	case "issued", "out":
		return "icon-minus"
	case "transfer":
		return "icon-map-pin"
	case "adjustment":
		return "icon-edit"
	default:
		return "icon-activity"
	}
}

func variantForTransactionType(kind string) string {
	switch kind {
	case "received", "in":
		return "client"
	case "issued", "out":
		return "integration"
	case "transfer":
		return "award"
	default:
		return "quote"
	}
}
