package detail

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// LineItemFormData is the template data for the line item drawer form.
type LineItemFormData struct {
	FormAction       string
	IsEdit           bool
	ID               string
	RevenueID        string
	Description      string
	Quantity         string
	UnitPrice        string
	CostPrice        string
	Discount         string
	Notes            string
	LineItemType     string
	InventoryItemID  string
	InventoryItems   []SelectOption
	SearchProductURL string
	ProductID        string
	ProductLabel     string
	PriceListID      string
	PriceProductID   string
	CommonLabels     any
	Labels           centymo.RevenueDetailLabels
}

// DiscountFormData is the template data for the discount drawer form.
type DiscountFormData struct {
	FormAction   string
	RevenueID    string
	Description  string
	Amount       string
	CommonLabels any
	Labels       centymo.RevenueDetailLabels
}

// SelectOption represents an option in a select dropdown.
type SelectOption struct {
	Value string
	Label string
}

// LineItemDeps holds dependencies for line item action handlers.
type LineItemDeps struct {
	Routes           centymo.RevenueRoutes
	Labels           centymo.RevenueLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels
	SearchProductURL string

	// Typed inventory operations
	ListInventoryItems func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)

	// Typed revenue operations
	ReadRevenue   func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	UpdateRevenue func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)

	// Typed line item operations
	CreateRevenueLineItem func(ctx context.Context, req *revenuelineitempb.CreateRevenueLineItemRequest) (*revenuelineitempb.CreateRevenueLineItemResponse, error)
	ReadRevenueLineItem   func(ctx context.Context, req *revenuelineitempb.ReadRevenueLineItemRequest) (*revenuelineitempb.ReadRevenueLineItemResponse, error)
	UpdateRevenueLineItem func(ctx context.Context, req *revenuelineitempb.UpdateRevenueLineItemRequest) (*revenuelineitempb.UpdateRevenueLineItemResponse, error)
	DeleteRevenueLineItem func(ctx context.Context, req *revenuelineitempb.DeleteRevenueLineItemRequest) (*revenuelineitempb.DeleteRevenueLineItemResponse, error)
	ListRevenueLineItems  func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)
}

// NewLineItemTableView returns a view that renders only the line items table (for HTMX refresh).
func NewLineItemTableView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
			Data: &revenuepb.Revenue{Id: revenueID},
		})
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", revenueID, err)
			return lineItemHTMXError(err.Error())
		}
		rData := resp.GetData()
		if len(rData) == 0 {
			return lineItemHTMXError("sale not found")
		}
		currency := rData[0].GetCurrency()

		lineItems := listLineItemMaps(ctx, deps.ListRevenueLineItems, revenueID, currency)
		perms := view.GetUserPermissions(ctx)
		table := buildLineItemTableWithActions(lineItems, deps.Labels, deps.TableLabels, currency, revenueID, deps.Routes, perms)
		return view.OK("table-card", table)
	})
}

// NewLineItemAddView creates the line item add action (GET = form, POST = create).
func NewLineItemAddView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			inventoryItems := loadInventoryItems(ctx, deps.ListInventoryItems)
			return view.OK("revenue-line-item-drawer-form", &LineItemFormData{
				FormAction:       route.ResolveURL(deps.Routes.LineItemAddURL, "id", revenueID),
				RevenueID:        revenueID,
				Quantity:         "1",
				LineItemType:     "item",
				InventoryItems:   inventoryItems,
				SearchProductURL: deps.SearchProductURL,
				Labels:           deps.Labels.Detail,
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — create line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")
		costPrice := r.FormValue("cost_price")
		discount := r.FormValue("discount")

		total := calculateLineItemTotal(quantity, unitPrice, discount)

		quantityF, _ := strconv.ParseFloat(quantity, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPrice, 64)
		costPriceF, _ := strconv.ParseFloat(costPrice, 64)
		unitPriceCentavos := int64(math.Round(unitPriceF * 100))
		costPriceCentavos := int64(math.Round(costPriceF * 100))

		lineItemData := &revenuelineitempb.RevenueLineItem{
			RevenueId:       revenueID,
			Description:     r.FormValue("description"),
			Quantity:        quantityF,
			UnitPrice:       unitPriceCentavos,
			CostPrice:       &costPriceCentavos,
			TotalPrice:      total,
			LineItemType:    "item",
			InventoryItemId: r.FormValue("inventory_item_id"),
			Notes:           strPtr(r.FormValue("notes")),
		}
		if pid := r.FormValue("product_id"); pid != "" {
			lineItemData.ProductId = strPtr(pid)
		}
		if v := r.FormValue("price_product_id"); v != "" {
			lineItemData.PriceProductId = strPtr(v)
		}
		_, err := deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
			Data: lineItemData,
		})
		if err != nil {
			log.Printf("Failed to create line item: %v", err)
			return lineItemHTMXError(err.Error())
		}

		// Recalculate sale total
		recalculateRevenueTotalTyped(ctx, deps.ListRevenueLineItems, deps.UpdateRevenue, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// floatPtr returns a pointer to a float64.
func floatPtr(f float64) *float64 {
	return &f
}

// NewLineItemEditView creates the line item edit action (GET = form, POST = update).
func NewLineItemEditView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")
		itemID := viewCtx.Request.PathValue("itemId")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadRevenueLineItem(ctx, &revenuelineitempb.ReadRevenueLineItemRequest{
				Data: &revenuelineitempb.RevenueLineItem{Id: itemID},
			})
			if err != nil {
				log.Printf("Failed to read line item %s: %v", itemID, err)
				return lineItemHTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return lineItemHTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			inventoryItems := loadInventoryItems(ctx, deps.ListInventoryItems)

			productID := record.GetProductId()
			productLabel := ""
			if productID != "" {
				if p := record.GetProduct(); p != nil && p.GetName() != "" {
					productLabel = p.GetName()
				} else {
					productLabel = record.GetDescription()
				}
			}

			return view.OK("revenue-line-item-drawer-form", &LineItemFormData{
				FormAction:       route.ResolveURL(deps.Routes.LineItemEditURL, "id", revenueID, "itemId", itemID),
				IsEdit:           true,
				ID:               itemID,
				RevenueID:        revenueID,
				Description:      record.GetDescription(),
				Quantity:         fmt.Sprintf("%.0f", record.GetQuantity()),
				UnitPrice:        fmt.Sprintf("%.2f", float64(record.GetUnitPrice())/100.0),
				CostPrice:        fmt.Sprintf("%.2f", float64(record.GetCostPrice())/100.0),
				Discount:         "0",
				Notes:            record.GetNotes(),
				LineItemType:     "item",
				InventoryItemID:  record.GetInventoryItemId(),
				InventoryItems:   inventoryItems,
				SearchProductURL: deps.SearchProductURL,
				ProductID:        productID,
				ProductLabel:     productLabel,
				PriceListID:      record.GetPriceListId(),
				PriceProductID:   record.GetPriceProductId(),
				Labels:           deps.Labels.Detail,
				CommonLabels:     nil,
			})
		}

		// POST — update line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")
		discount := r.FormValue("discount")

		total := calculateLineItemTotal(quantity, unitPrice, discount)

		quantityF, _ := strconv.ParseFloat(quantity, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPrice, 64)
		costPriceF, _ := strconv.ParseFloat(r.FormValue("cost_price"), 64)
		unitPriceCentavos := int64(math.Round(unitPriceF * 100))
		costPriceCentavos := int64(math.Round(costPriceF * 100))

		updateData := &revenuelineitempb.RevenueLineItem{
			Id:              itemID,
			Description:     r.FormValue("description"),
			Quantity:        quantityF,
			UnitPrice:       unitPriceCentavos,
			CostPrice:       &costPriceCentavos,
			TotalPrice:      total,
			InventoryItemId: r.FormValue("inventory_item_id"),
			Notes:           strPtr(r.FormValue("notes")),
		}
		if pid := r.FormValue("product_id"); pid != "" {
			updateData.ProductId = strPtr(pid)
		}
		if v := r.FormValue("price_product_id"); v != "" {
			updateData.PriceProductId = strPtr(v)
		}
		_, err := deps.UpdateRevenueLineItem(ctx, &revenuelineitempb.UpdateRevenueLineItemRequest{
			Data: updateData,
		})
		if err != nil {
			log.Printf("Failed to update line item %s: %v", itemID, err)
			return lineItemHTMXError(err.Error())
		}

		recalculateRevenueTotalTyped(ctx, deps.ListRevenueLineItems, deps.UpdateRevenue, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemRemoveView creates the line item remove action (POST only).
func NewLineItemRemoveView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")

		itemID := viewCtx.Request.URL.Query().Get("itemId")
		if itemID == "" {
			_ = viewCtx.Request.ParseForm()
			itemID = viewCtx.Request.FormValue("itemId")
		}
		if itemID == "" {
			return lineItemHTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteRevenueLineItem(ctx, &revenuelineitempb.DeleteRevenueLineItemRequest{
			Data: &revenuelineitempb.RevenueLineItem{Id: itemID},
		})
		if err != nil {
			log.Printf("Failed to delete line item %s: %v", itemID, err)
			return lineItemHTMXError(err.Error())
		}

		recalculateRevenueTotalTyped(ctx, deps.ListRevenueLineItems, deps.UpdateRevenue, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemDiscountView creates the discount add action (GET = form, POST = create).
func NewLineItemDiscountView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("revenue-line-item-discount-form", &DiscountFormData{
				FormAction:   route.ResolveURL(deps.Routes.LineItemDiscountURL, "id", revenueID),
				RevenueID:    revenueID,
				Labels:       deps.Labels.Detail,
				CommonLabels: nil,
			})
		}

		// POST — create discount line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		amount := r.FormValue("amount")

		// Store discount amount as negative total
		amountF, err := strconv.ParseFloat(amount, 64)
		if err != nil || amountF <= 0 {
			return lineItemHTMXError(deps.Labels.Errors.InvalidDiscount)
		}

		_, err = deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
			Data: &revenuelineitempb.RevenueLineItem{
				RevenueId:    revenueID,
				Description:  r.FormValue("description"),
				Quantity:     1,
				UnitPrice:    0,
				TotalPrice:   -int64(math.Round(amountF * 100)),
				LineItemType: "discount",
			},
		})
		if err != nil {
			log.Printf("Failed to create discount line item: %v", err)
			return lineItemHTMXError(err.Error())
		}

		recalculateRevenueTotalTyped(ctx, deps.ListRevenueLineItems, deps.UpdateRevenue, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// buildLineItemTableWithActions builds the line items table with row actions.
func buildLineItemTableWithActions(items []map[string]any, l centymo.RevenueLabels, tableLabels types.TableLabels, currency string, revenueID string, routes centymo.RevenueRoutes, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "type", Label: l.Detail.ItemType, Sortable: false, WidthClass: "col-lg"},
		{Key: "description", Label: l.Detail.Description, Sortable: false},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: false, WidthClass: "col-md"},
		{Key: "unit_price", Label: l.Detail.UnitPrice, Sortable: false, WidthClass: "col-3xl"},
		{Key: "discount", Label: l.Detail.Discount, Sortable: false, WidthClass: "col-lg"},
		{Key: "total", Label: l.Detail.Total, Sortable: false, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, item := range items {
		id, _ := item["id"].(string)
		description, _ := item["description"].(string)
		quantity, _ := item["quantity"].(string)
		unitPrice, _ := item["unit_price"].(string)
		discount, _ := item["discount"].(string)
		total, _ := item["total"].(string)
		lineItemType, _ := item["line_item_type"].(string)

		typeBadge := l.Detail.ItemTypeItem
		typeVariant := "info"
		if lineItemType == "discount" {
			typeBadge = l.Detail.ItemTypeDiscount
			typeVariant = "warning"
		}

		var actions []types.TableAction
		if lineItemType != "discount" {
			actions = append(actions, types.TableAction{
				Type:            "edit",
				Label:           l.Detail.EditItem,
				Action:          "edit",
				URL:             route.ResolveURL(routes.LineItemEditURL, "id", revenueID, "itemId", id),
				DrawerTitle:     l.Detail.EditItem,
				Disabled:        !perms.Can("invoice", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}
		actions = append(actions, types.TableAction{
			Type:            "delete",
			Label:           l.Detail.RemoveItem,
			Action:          "delete",
			URL:             route.ResolveURL(routes.LineItemRemoveURL, "id", revenueID) + "?itemId=" + id,
			ItemName:        description,
			Disabled:        !perms.Can("invoice", "update"),
			DisabledTooltip: l.Errors.PermissionDenied,
		})

		row := types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "badge", Value: typeBadge, Variant: typeVariant},
				{Type: "text", Value: description},
				{Type: "text", Value: quantity},
				{Type: "text", Value: currency + " " + unitPrice},
				{Type: "text", Value: discount},
				{Type: "text", Value: currency + " " + total},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "line-items-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.ItemEmptyTitle,
			Message: l.Detail.ItemEmptyMessage,
		},
	}
}

// loadInventoryItems loads inventory items for the select dropdown.
func loadInventoryItems(ctx context.Context, listFn func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)) []SelectOption {
	resp, err := listFn(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Failed to list inventory items: %v", err)
		return nil
	}

	options := []SelectOption{}
	for _, item := range resp.GetData() {
		id := item.GetId()
		name := item.GetName()
		sku := item.GetSku()

		label := name
		if sku != "" {
			label = name + " (" + sku + ")"
		}
		if label == "" {
			label = id
		}

		options = append(options, SelectOption{
			Value: id,
			Label: label,
		})
	}

	return options
}

// calculateLineItemTotal calculates total from quantity * unit_price - discount.
// All monetary values are in centavos (int64). unitPrice and discount are parsed
// as decimal strings and converted to centavos internally.
func calculateLineItemTotal(quantityStr, unitPriceStr, discountStr string) int64 {
	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	unitPrice, _ := strconv.ParseFloat(unitPriceStr, 64)
	discount, _ := strconv.ParseFloat(discountStr, 64)

	if quantity == 0 {
		quantity = 1
	}

	total := quantity*unitPrice - discount
	return int64(math.Round(total * 100))
}

// recalculateRevenueTotalTyped recalculates the sale's total_amount from its line items using typed use cases.
func recalculateRevenueTotalTyped(
	ctx context.Context,
	listFn func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error),
	updateFn func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error),
	revenueID string,
) {
	resp, err := listFn(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
		RevenueId: &revenueID,
	})
	if err != nil {
		log.Printf("Failed to list line items for total recalculation: %v", err)
		return
	}

	var totalAmount int64
	for _, item := range resp.GetData() {
		if item.GetRevenueId() == revenueID {
			totalAmount += item.GetTotalPrice()
		}
	}

	_, err = updateFn(ctx, &revenuepb.UpdateRevenueRequest{
		Data: &revenuepb.Revenue{
			Id:          revenueID,
			TotalAmount: totalAmount,
		},
	})
	if err != nil {
		log.Printf("Failed to update revenue total: %v", err)
	}
}

// lineItemHTMXSuccess returns a success response that triggers table refresh.
func lineItemHTMXSuccess(tableID string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": fmt.Sprintf(`{"formSuccess":true,"refreshTable":"%s"}`, tableID),
		},
	}
}

// lineItemHTMXError returns an error response for HTMX.
func lineItemHTMXError(message string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusUnprocessableEntity,
		Headers: map[string]string{
			"HX-Error-Message": message,
		},
	}
}
