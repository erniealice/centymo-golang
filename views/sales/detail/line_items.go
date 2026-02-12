package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
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
	CommonLabels     any
	Labels           centymo.SalesDetailLabels
}

// DiscountFormData is the template data for the discount drawer form.
type DiscountFormData struct {
	FormAction   string
	RevenueID    string
	Description  string
	Amount       string
	CommonLabels any
	Labels       centymo.SalesDetailLabels
}

// SelectOption represents an option in a select dropdown.
type SelectOption struct {
	Value string
	Label string
}

// LineItemDeps holds dependencies for line item action handlers.
type LineItemDeps struct {
	DB           centymo.DataSource
	Labels       centymo.SalesLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// NewLineItemTableView returns a view that renders only the line items table (for HTMX refresh).
func NewLineItemTableView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		revenue, err := deps.DB.Read(ctx, "revenue", revenueID)
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", revenueID, err)
			return lineItemHTMXError("Failed to load sale")
		}

		allLineItems, err := deps.DB.ListSimple(ctx, "revenue_line_item")
		if err != nil {
			log.Printf("Failed to list line items: %v", err)
			allLineItems = []map[string]any{}
		}
		lineItems := filterLineItems(allLineItems, revenueID)
		currency, _ := revenue["currency"].(string)
		table := buildLineItemTableWithActions(lineItems, deps.Labels, deps.TableLabels, currency, revenueID)
		return view.OK("table-card", table)
	})
}

// NewLineItemAddView creates the line item add action (GET = form, POST = create).
func NewLineItemAddView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			inventoryItems := loadInventoryItems(ctx, deps.DB, revenueID)
			return view.OK("sales-line-item-drawer-form", &LineItemFormData{
				FormAction:      fmt.Sprintf("/action/sales/detail/%s/items/add", revenueID),
				RevenueID:       revenueID,
				Quantity:        "1",
				LineItemType:    "item",
				InventoryItems:  inventoryItems,
				Labels:          deps.Labels.Detail,
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — create line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError("Invalid form data")
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")
		costPrice := r.FormValue("cost_price")
		discount := r.FormValue("discount")

		total := calculateLineItemTotal(quantity, unitPrice, discount)

		data := map[string]any{
			"revenue_id":        revenueID,
			"description":       r.FormValue("description"),
			"quantity":          quantity,
			"unit_price":        unitPrice,
			"cost_price":        costPrice,
			"discount":          discount,
			"total":             fmt.Sprintf("%.2f", total),
			"line_item_type":    "item",
			"inventory_item_id": r.FormValue("inventory_item_id"),
			"notes":             r.FormValue("notes"),
		}

		_, err := deps.DB.Create(ctx, "revenue_line_item", data)
		if err != nil {
			log.Printf("Failed to create line item: %v", err)
			return lineItemHTMXError("Failed to add line item")
		}

		// Recalculate sale total
		recalculateSaleTotal(ctx, deps.DB, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemEditView creates the line item edit action (GET = form, POST = update).
func NewLineItemEditView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")
		itemID := viewCtx.Request.PathValue("itemId")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "revenue_line_item", itemID)
			if err != nil {
				log.Printf("Failed to read line item %s: %v", itemID, err)
				return lineItemHTMXError("Line item not found")
			}

			description, _ := record["description"].(string)
			quantity, _ := record["quantity"].(string)
			unitPrice, _ := record["unit_price"].(string)
			costPrice, _ := record["cost_price"].(string)
			discount, _ := record["discount"].(string)
			notes, _ := record["notes"].(string)
			inventoryItemID, _ := record["inventory_item_id"].(string)

			inventoryItems := loadInventoryItems(ctx, deps.DB, revenueID)

			return view.OK("sales-line-item-drawer-form", &LineItemFormData{
				FormAction:      fmt.Sprintf("/action/sales/detail/%s/items/edit/%s", revenueID, itemID),
				IsEdit:          true,
				ID:              itemID,
				RevenueID:       revenueID,
				Description:     description,
				Quantity:        quantity,
				UnitPrice:       unitPrice,
				CostPrice:       costPrice,
				Discount:        discount,
				Notes:           notes,
				LineItemType:    "item",
				InventoryItemID: inventoryItemID,
				InventoryItems:  inventoryItems,
				Labels:          deps.Labels.Detail,
				CommonLabels:    nil,
			})
		}

		// POST — update line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError("Invalid form data")
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")
		discount := r.FormValue("discount")

		total := calculateLineItemTotal(quantity, unitPrice, discount)

		data := map[string]any{
			"description":       r.FormValue("description"),
			"quantity":          quantity,
			"unit_price":        unitPrice,
			"cost_price":        r.FormValue("cost_price"),
			"discount":          discount,
			"total":             fmt.Sprintf("%.2f", total),
			"inventory_item_id": r.FormValue("inventory_item_id"),
			"notes":             r.FormValue("notes"),
		}

		_, err := deps.DB.Update(ctx, "revenue_line_item", itemID, data)
		if err != nil {
			log.Printf("Failed to update line item %s: %v", itemID, err)
			return lineItemHTMXError("Failed to update line item")
		}

		recalculateSaleTotal(ctx, deps.DB, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemRemoveView creates the line item remove action (POST only).
func NewLineItemRemoveView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		itemID := viewCtx.Request.URL.Query().Get("itemId")
		if itemID == "" {
			_ = viewCtx.Request.ParseForm()
			itemID = viewCtx.Request.FormValue("itemId")
		}
		if itemID == "" {
			return lineItemHTMXError("Line item ID is required")
		}

		err := deps.DB.Delete(ctx, "revenue_line_item", itemID)
		if err != nil {
			log.Printf("Failed to delete line item %s: %v", itemID, err)
			return lineItemHTMXError("Failed to remove line item")
		}

		recalculateSaleTotal(ctx, deps.DB, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemDiscountView creates the discount add action (GET = form, POST = create).
func NewLineItemDiscountView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("sales-line-item-discount-form", &DiscountFormData{
				FormAction:  fmt.Sprintf("/action/sales/detail/%s/items/add-discount", revenueID),
				RevenueID:   revenueID,
				Labels:      deps.Labels.Detail,
				CommonLabels: nil,
			})
		}

		// POST — create discount line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError("Invalid form data")
		}

		r := viewCtx.Request
		amount := r.FormValue("amount")

		// Store discount amount as negative total
		amountF, err := strconv.ParseFloat(amount, 64)
		if err != nil || amountF <= 0 {
			return lineItemHTMXError("Discount amount must be a positive number")
		}

		data := map[string]any{
			"revenue_id":     revenueID,
			"description":    r.FormValue("description"),
			"quantity":       "1",
			"unit_price":     "0",
			"cost_price":     "0",
			"discount":       "0",
			"total":          fmt.Sprintf("-%.2f", amountF),
			"line_item_type": "discount",
		}

		_, err = deps.DB.Create(ctx, "revenue_line_item", data)
		if err != nil {
			log.Printf("Failed to create discount line item: %v", err)
			return lineItemHTMXError("Failed to add discount")
		}

		recalculateSaleTotal(ctx, deps.DB, revenueID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// buildLineItemTableWithActions builds the line items table with row actions.
func buildLineItemTableWithActions(items []map[string]any, l centymo.SalesLabels, tableLabels types.TableLabels, currency string, revenueID string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "type", Label: l.Detail.ItemType, Sortable: false, Width: "90px"},
		{Key: "description", Label: l.Detail.Description, Sortable: false},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: false, Width: "80px"},
		{Key: "unit_price", Label: l.Detail.UnitPrice, Sortable: false, Width: "130px"},
		{Key: "discount", Label: l.Detail.Discount, Sortable: false, Width: "100px"},
		{Key: "total", Label: l.Detail.Total, Sortable: false, Width: "130px"},
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
				Type:        "edit",
				Label:       l.Detail.EditItem,
				Action:      "edit",
				URL:         fmt.Sprintf("/action/sales/detail/%s/items/edit/%s", revenueID, id),
				DrawerTitle: l.Detail.EditItem,
			})
		}
		actions = append(actions, types.TableAction{
			Type:     "delete",
			Label:    l.Detail.RemoveItem,
			Action:   "delete",
			URL:      fmt.Sprintf("/action/sales/detail/%s/items/remove?itemId=%s", revenueID, id),
			ItemName: description,
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
func loadInventoryItems(ctx context.Context, db centymo.DataSource, revenueID string) []SelectOption {
	allItems, err := db.ListSimple(ctx, "inventory_item")
	if err != nil {
		log.Printf("Failed to list inventory items: %v", err)
		return nil
	}

	options := []SelectOption{}
	for _, item := range allItems {
		id, _ := item["id"].(string)
		productName, _ := item["product_name"].(string)
		sku, _ := item["sku"].(string)

		label := productName
		if sku != "" {
			label = productName + " (" + sku + ")"
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
func calculateLineItemTotal(quantityStr, unitPriceStr, discountStr string) float64 {
	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	unitPrice, _ := strconv.ParseFloat(unitPriceStr, 64)
	discount, _ := strconv.ParseFloat(discountStr, 64)

	if quantity == 0 {
		quantity = 1
	}

	total := quantity*unitPrice - discount
	return total
}

// recalculateSaleTotal recalculates the sale's total_amount from its line items.
func recalculateSaleTotal(ctx context.Context, db centymo.DataSource, revenueID string) {
	allLineItems, err := db.ListSimple(ctx, "revenue_line_item")
	if err != nil {
		log.Printf("Failed to list line items for total recalculation: %v", err)
		return
	}

	var totalAmount float64
	for _, item := range allLineItems {
		rid, _ := item["revenue_id"].(string)
		if rid != revenueID {
			continue
		}
		totalStr, _ := item["total"].(string)
		totalF, _ := strconv.ParseFloat(totalStr, 64)
		totalAmount += totalF
	}

	_, err = db.Update(ctx, "revenue", revenueID, map[string]any{
		"total_amount": fmt.Sprintf("%.2f", totalAmount),
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
