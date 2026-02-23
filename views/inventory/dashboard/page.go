package dashboard

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
)

// Deps holds view dependencies.
type Deps struct {
	ListInventoryItems         func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ListInventorySerials       func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)
	ListInventoryTransactions  func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	ListInventoryDepreciations func(ctx context.Context, req *inventorydepreciationpb.ListInventoryDepreciationsRequest) (*inventorydepreciationpb.ListInventoryDepreciationsResponse, error)
	Labels                     centymo.InventoryLabels
	CommonLabels               pyeza.CommonLabels
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
		widgets := buildDashboardWidgets(ctx, deps, l)

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
		widgets := buildDashboardWidgets(ctx, deps, deps.Labels)
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
				"performed_by":     t.GetPerformedBy(),
				"inventory_item_id": t.GetInventoryItemId(),
			})
		}

		return view.OK("inventory-dashboard-movements", map[string]any{
			"Transactions": txnMaps,
			"Labels":       deps.Labels,
		})
	})
}

// NewDashboardAlertsAction returns the low stock alerts widget partial.
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

		// Convert low-stock items to map[string]any for template backward compat
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
					"item_type":         item.GetProduct().GetItemType(),
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

func buildDashboardWidgets(ctx context.Context, deps *Deps, l centymo.InventoryLabels) []WidgetData {
	resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list inventory items: %v", err)
	}
	var items []*inventoryitempb.InventoryItem
	if resp != nil {
		items = resp.GetData()
	}

	// Calculate stats
	totalItems := len(items)
	var totalStockValue float64
	lowStockCount := 0
	categoryCount := map[string]int{}

	for _, item := range items {
		onHand := item.GetQuantityOnHand()
		reserved := item.GetQuantityReserved()
		reorderLvl := item.GetReorderLevel()
		available := onHand - reserved

		// Approximate stock value (quantity * a base price)
		totalStockValue += onHand * 100 // placeholder unit cost

		if reorderLvl > 0 && available <= reorderLvl {
			lowStockCount++
		}

		itemType := item.GetProduct().GetItemType()
		if itemType == "" {
			itemType = "non_serialized"
		}
		categoryCount[itemType]++
	}

	// Load serials for status distribution
	serialResp, err := deps.ListInventorySerials(ctx, &inventoryserialpb.ListInventorySerialsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list serials: %v", err)
	}
	var serials []*inventoryserialpb.InventorySerial
	if serialResp != nil {
		serials = serialResp.GetData()
	}
	serialAvailable := 0
	for _, s := range serials {
		if s.GetStatus() == "available" {
			serialAvailable++
		}
	}

	// Load depreciation data
	depResp, err := deps.ListInventoryDepreciations(ctx, &inventorydepreciationpb.ListInventoryDepreciationsRequest{})
	if err != nil {
		log.Printf("Dashboard: Failed to list depreciations: %v", err)
	}
	var depreciations []*inventorydepreciationpb.InventoryDepreciation
	if depResp != nil {
		depreciations = depResp.GetData()
	}
	var totalCostBasis, totalBookValue float64
	for _, d := range depreciations {
		totalCostBasis += d.GetCostBasis()
		totalBookValue += d.GetBookValue()
	}

	return []WidgetData{
		{Icon: "icon-dollar-sign", Value: fmt.Sprintf("%.0f", totalStockValue), Label: l.Dashboard.TotalStockValue, Color: "terracotta"},
		{Icon: "icon-alert-triangle", Value: fmt.Sprintf("%d", lowStockCount), Label: l.Dashboard.LowStockAlerts, Color: "amber"},
		{Icon: "icon-repeat", Value: fmt.Sprintf("%d", totalItems), Label: l.Dashboard.StockTurnover, Color: "sage"},
		{Icon: "icon-map-pin", Value: fmt.Sprintf("%d", len(centymo.LocationMap)), Label: l.Dashboard.ItemsByLocation, Color: "navy"},
		{Icon: "icon-trending-down", Value: fmt.Sprintf("%.0f / %.0f", totalCostBasis, totalBookValue), Label: l.Dashboard.DepreciationSummary, Color: "terracotta"},
		{Icon: "icon-hash", Value: fmt.Sprintf("%d / %d", serialAvailable, len(serials)), Label: l.Dashboard.SerialUnitStatus, Color: "sage"},
		{Icon: "icon-activity", Value: "\u2014", Label: l.Dashboard.RecentMovements, Color: "navy"},
		{Icon: "icon-pie-chart", Value: fmt.Sprintf("%d types", len(categoryCount)), Label: l.Dashboard.CategoryDistribution, Color: "amber"},
	}
}
