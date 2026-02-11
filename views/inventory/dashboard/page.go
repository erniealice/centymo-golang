package dashboard

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB               centymo.DataSource
	Labels           centymo.InventoryLabels
	CommonLabels     pyeza.CommonLabels
}

// WidgetData holds a single KPI widget's data.
type WidgetData struct {
	Icon    string
	Value   string
	Label   string
	Trend   string
	TrendUp bool
	Color   string
}

// PageData holds the data for the inventory dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.InventoryLabels
	Widgets         []WidgetData
}

// NewView creates the inventory dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels

		// Load dashboard data
		widgets := buildDashboardWidgets(ctx, deps.DB, l)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Inventory Dashboard",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "inventory",
				ActiveSubNav: "dashboard",
				HeaderTitle:  "Inventory Dashboard",
				HeaderIcon:   "icon-briefcase",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "inventory-dashboard-content",
			Labels:          l,
			Widgets:         widgets,
		}

		return view.OK("inventory-dashboard", pageData)
	})
}

// NewDashboardStatsAction returns stats widget partials via HTMX.
func NewDashboardStatsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		widgets := buildDashboardWidgets(ctx, deps.DB, deps.Labels)
		return view.OK("inventory-dashboard-stats", map[string]any{
			"Widgets": widgets,
			"Labels":  deps.Labels,
		})
	})
}

// NewDashboardChartAction returns the chart widget partial via HTMX.
func NewDashboardChartAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return view.OK("inventory-dashboard-chart", map[string]any{
			"Labels": deps.Labels,
		})
	})
}

// NewDashboardMovementsAction returns the recent movements widget partial.
func NewDashboardMovementsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		txns, err := deps.DB.ListSimple(ctx, "inventory_transaction")
		if err != nil {
			log.Printf("Failed to list transactions for dashboard: %v", err)
			txns = []map[string]any{}
		}

		// Take last 10 sorted by date (most recent first)
		if len(txns) > 10 {
			txns = txns[len(txns)-10:]
		}

		return view.OK("inventory-dashboard-movements", map[string]any{
			"Transactions": txns,
			"Labels":       deps.Labels,
		})
	})
}

// NewDashboardAlertsAction returns the low stock alerts widget partial.
func NewDashboardAlertsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		items, err := deps.DB.ListSimple(ctx, "inventory_item")
		if err != nil {
			log.Printf("Failed to list items for alerts: %v", err)
			items = []map[string]any{}
		}

		alerts := []map[string]any{}
		for _, item := range items {
			onHand := toFloat64(item["quantity_on_hand"])
			reserved := toFloat64(item["quantity_reserved"])
			reorderLvl := toFloat64(item["reorder_level"])
			available := onHand - reserved
			if reorderLvl > 0 && available <= reorderLvl {
				alerts = append(alerts, item)
			}
		}

		return view.OK("inventory-dashboard-alerts", map[string]any{
			"Alerts": alerts,
			"Labels": deps.Labels,
		})
	})
}

func buildDashboardWidgets(ctx context.Context, db centymo.DataSource, l centymo.InventoryLabels) []WidgetData {
	items, err := db.ListSimple(ctx, "inventory_item")
	if err != nil {
		log.Printf("Dashboard: Failed to list inventory items: %v", err)
		items = []map[string]any{}
	}

	// Calculate stats
	totalItems := len(items)
	var totalStockValue float64
	lowStockCount := 0
	categoryCount := map[string]int{}

	for _, item := range items {
		onHand := toFloat64(item["quantity_on_hand"])
		reserved := toFloat64(item["quantity_reserved"])
		reorderLvl := toFloat64(item["reorder_level"])
		available := onHand - reserved

		// Approximate stock value (quantity * a base price)
		totalStockValue += onHand * 100 // placeholder unit cost

		if reorderLvl > 0 && available <= reorderLvl {
			lowStockCount++
		}

		itemType, _ := item["item_type"].(string)
		if itemType == "" {
			itemType = "non_serialized"
		}
		categoryCount[itemType]++
	}

	// Load serials for status distribution
	serials, err := db.ListSimple(ctx, "inventory_serial")
	if err != nil {
		log.Printf("Dashboard: Failed to list serials: %v", err)
		serials = []map[string]any{}
	}
	serialAvailable := 0
	for _, s := range serials {
		status, _ := s["status"].(string)
		if status == "available" {
			serialAvailable++
		}
	}

	// Load depreciation data
	depreciations, err := db.ListSimple(ctx, "inventory_depreciation")
	if err != nil {
		log.Printf("Dashboard: Failed to list depreciations: %v", err)
		depreciations = []map[string]any{}
	}
	var totalCostBasis, totalBookValue float64
	for _, d := range depreciations {
		totalCostBasis += toFloat64(d["cost_basis"])
		totalBookValue += toFloat64(d["book_value"])
	}

	return []WidgetData{
		{Icon: "icon-dollar-sign", Value: fmt.Sprintf("%.0f", totalStockValue), Label: l.Dashboard.TotalStockValue, Color: "terracotta"},
		{Icon: "icon-alert-triangle", Value: fmt.Sprintf("%d", lowStockCount), Label: l.Dashboard.LowStockAlerts, Color: "amber"},
		{Icon: "icon-repeat", Value: fmt.Sprintf("%d", totalItems), Label: l.Dashboard.StockTurnover, Color: "sage"},
		{Icon: "icon-map-pin", Value: fmt.Sprintf("%d", len(centymo.LocationMap)), Label: l.Dashboard.ItemsByLocation, Color: "navy"},
		{Icon: "icon-trending-down", Value: fmt.Sprintf("%.0f / %.0f", totalCostBasis, totalBookValue), Label: l.Dashboard.DepreciationSummary, Color: "terracotta"},
		{Icon: "icon-hash", Value: fmt.Sprintf("%d / %d", serialAvailable, len(serials)), Label: l.Dashboard.SerialUnitStatus, Color: "sage"},
		{Icon: "icon-activity", Value: "—", Label: l.Dashboard.RecentMovements, Color: "navy"},
		{Icon: "icon-pie-chart", Value: fmt.Sprintf("%d types", len(categoryCount)), Label: l.Dashboard.CategoryDistribution, Color: "amber"},
	}
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}
