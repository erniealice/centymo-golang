package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	purchaseorderlineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order_line_item"
)

// LineItemDeps holds dependencies for the PO line item table view.
type LineItemDeps struct {
	Routes      centymo.ExpenditureRoutes
	Labels      centymo.ExpenditureLabels
	TableLabels types.TableLabels

	ReadPurchaseOrder          func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	ListPurchaseOrderLineItems func(ctx context.Context, req *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error)
}

// DetailViewDeps holds view dependencies for the purchase order detail page.
type DetailViewDeps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadPurchaseOrder          func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	ListPurchaseOrderLineItems func(ctx context.Context, req *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error)
}

// PageData holds the data for the purchase order detail page.
type PageData struct {
	types.PageData
	ContentTemplate    string
	PurchaseOrder      map[string]any
	Labels             centymo.ExpenditureLabels
	ActiveTab          string
	TabItems           []pyeza.TabItem
	LineItemTable      *types.TableConfig
	LineItemAddURL     string
	TotalAmount        string
	SetStatusURL       string
	ConfirmReceiptURL  string
}

// purchaseOrderToMap converts a PurchaseOrder proto to a map for template use.
func purchaseOrderToMap(po *purchaseorderpb.PurchaseOrder) map[string]any {
	supplierName := ""
	if supplier := po.GetSupplier(); supplier != nil {
		supplierName = supplier.GetName()
	}
	if supplierName == "" {
		supplierName = po.GetSupplierId()
	}

	locationName := ""
	if location := po.GetLocation(); location != nil {
		locationName = location.GetName()
	}

	currency := po.GetCurrency()
	return map[string]any{
		"id":                            po.GetId(),
		"po_number":                     po.GetPoNumber(),
		"po_type":                       po.GetPoType(),
		"status":                        po.GetStatus(),
		"supplier_id":                   po.GetSupplierId(),
		"supplier_name":                 supplierName,
		"location_id":                   po.GetLocationId(),
		"location_name":                 locationName,
		"order_date_string":             po.GetOrderDateString(),
		"expected_delivery_date_string": po.GetExpectedDeliveryDateString(),
		"currency":                      currency,
		"subtotal":                      types.MoneyCell(float64(po.GetSubtotal()), currency, true),
		"tax_amount":                    types.MoneyCell(float64(po.GetTaxAmount()), currency, true),
		"total_amount":                  types.MoneyCell(float64(po.GetTotalAmount()), currency, true),
		"payment_terms":                 po.GetPaymentTerms(),
		"shipping_terms":                po.GetShippingTerms(),
		"approved_by":                   po.GetApprovedBy(),
		"approved_date_string":          po.GetApprovedDateString(),
		"reference_number":              po.GetReferenceNumber(),
		"notes":                         po.GetNotes(),
		"active":                        po.GetActive(),
		"date_created_string":           po.GetDateCreatedString(),
		"date_modified_string":          po.GetDateModifiedString(),
	}
}

// buildTabItems builds the tab navigation for the purchase order detail page.
func buildTabItems(l centymo.ExpenditureLabels, id string, routes centymo.ExpenditureRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.PurchaseOrderDetailURL, "id", id)
	action := route.ResolveURL(routes.PurchaseOrderTabActionURL, "id", id, "tab", "")
	tabDetails := l.PurchaseOrder.Detail.TabBasicInfo
	if tabDetails == "" {
		tabDetails = "Details"
	}
	tabLineItems := l.PurchaseOrder.Detail.TabLineItems
	if tabLineItems == "" {
		tabLineItems = "Line Items"
	}
	return []pyeza.TabItem{
		{Key: "info", Label: tabDetails, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "items", Label: tabLineItems, Href: base + "?tab=items", HxGet: action + "items", Icon: "icon-list"},
	}
}

// buildLineItemTable builds the line items table config for a purchase order.
func buildLineItemTable(items []map[string]any, tableLabels types.TableLabels, currency string, purchaseOrderID string, routes centymo.ExpenditureRoutes, isDraft bool, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "line_number", Label: "Line #", Sortable: false, WidthClass: "col-md"},
		{Key: "description", Label: "Description", Sortable: false},
		{Key: "line_type", Label: "Type", Sortable: false, WidthClass: "col-lg"},
		{Key: "quantity_ordered", Label: "Qty Ordered", Sortable: false, WidthClass: "col-xl"},
		{Key: "quantity_received", Label: "Qty Received", Sortable: false, WidthClass: "col-2xl"},
		{Key: "quantity_billed", Label: "Qty Billed", Sortable: false, WidthClass: "col-xl"},
		{Key: "unit_price", Label: "Unit Price", Sortable: false, WidthClass: "col-3xl"},
		{Key: "total", Label: "Total", Sortable: false, WidthClass: "col-3xl"},
	}

	canEdit := isDraft && perms != nil && perms.Can("purchase_order", "update")

	rows := []types.TableRow{}
	for _, item := range items {
		id, _ := item["id"].(string)
		lineNumber, _ := item["line_number"].(string)
		description, _ := item["description"].(string)
		lineType, _ := item["line_type"].(string)
		qtyOrdered, _ := item["quantity_ordered"].(string)
		qtyReceived, _ := item["quantity_received"].(string)
		qtyBilled, _ := item["quantity_billed"].(string)
		unitPriceCell, _ := item["unit_price"].(types.TableCell)
		totalCell, _ := item["total"].(types.TableCell)

		var actions []types.TableAction
		if isDraft {
			actions = []types.TableAction{
				{
					Type:            "edit",
					Label:           "Edit",
					Action:          "edit",
					URL:             route.ResolveURL(routes.PurchaseOrderLineItemEditURL, "id", purchaseOrderID, "itemId", id),
					DrawerTitle:     "Edit Line Item",
					Disabled:        !canEdit,
					DisabledTooltip: "No permission",
				},
				{
					Type:            "delete",
					Label:           "Remove",
					Action:          "delete",
					URL:             route.ResolveURL(routes.PurchaseOrderLineItemRemoveURL, "id", purchaseOrderID) + "?itemId=" + id,
					ItemName:        description,
					Disabled:        !canEdit,
					DisabledTooltip: "No permission",
				},
			}
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: lineNumber},
				{Type: "text", Value: description},
				{Type: "text", Value: lineType},
				{Type: "text", Value: qtyOrdered},
				{Type: "text", Value: qtyReceived},
				{Type: "text", Value: qtyBilled},
				unitPriceCell,
				totalCell,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "po-line-items-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No line items",
			Message: "This purchase order has no line items yet.",
		},
	}
}

// listLineItemMaps lists line items for a purchase order and returns as maps.
func listLineItemMaps(ctx context.Context, listFn func(context.Context, *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error), purchaseOrderID string, currency string) []map[string]any {
	resp, err := listFn(ctx, &purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest{
		PurchaseOrderId: &purchaseOrderID,
	})
	if err != nil {
		log.Printf("Failed to list line items for purchase order %s: %v", purchaseOrderID, err)
		return []map[string]any{}
	}
	items := []map[string]any{}
	for _, item := range resp.GetData() {
		if item.GetPurchaseOrderId() == purchaseOrderID {
			items = append(items, map[string]any{
				"id":                item.GetId(),
				"purchase_order_id": item.GetPurchaseOrderId(),
				"description":       item.GetDescription(),
				"line_number":       fmt.Sprintf("%d", item.GetLineNumber()),
				"line_type":         item.GetLineType(),
				"quantity_ordered":  fmt.Sprintf("%.0f", item.GetQuantityOrdered()),
				"quantity_received": fmt.Sprintf("%.0f", item.GetQuantityReceived()),
				"quantity_billed":   fmt.Sprintf("%.0f", item.GetQuantityBilled()),
				"unit_price":        types.MoneyCell(float64(item.GetUnitPrice()), currency, true),
				"total":             types.MoneyCell(float64(item.GetTotalPrice()), currency, true),
				"notes":             item.GetNotes(),
			})
		}
	}
	return items
}

// NewView creates the purchase order detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadPurchaseOrder(ctx, &purchaseorderpb.ReadPurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read purchase order %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load purchase order: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("purchase order not found"))
		}
		po := purchaseOrderToMap(data[0])

		poNumber, _ := po["po_number"].(string)
		headerTitle := poNumber
		if headerTitle == "" {
			headerTitle, _ = po["reference_number"].(string)
		}
		if headerTitle == "" {
			headerTitle = id
		}

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(deps.Labels, id, deps.Routes)

		poStatus, _ := po["status"].(string)
		confirmReceiptURL := ""
		if poStatus == "approved" || poStatus == "partially_received" {
			confirmReceiptURL = route.ResolveURL(deps.Routes.PurchaseOrderConfirmReceiptURL, "id", id)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "purchase",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: deps.Labels.Page.PurchaseHeading,
				HeaderIcon:     "icon-shopping-cart",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "purchase-order-detail-content",
			PurchaseOrder:     po,
			Labels:            deps.Labels,
			ActiveTab:         activeTab,
			TabItems:          tabItems,
			SetStatusURL:      deps.Routes.PurchaseOrderSetStatusURL,
			ConfirmReceiptURL: confirmReceiptURL,
		}

		switch activeTab {
		case "info":
			// po map has everything
		case "items":
			if deps.ListPurchaseOrderLineItems != nil {
				perms := view.GetUserPermissions(ctx)
				currency, _ := po["currency"].(string)
				status, _ := po["status"].(string)
				isDraft := status == "draft"
				lineItems := listLineItemMaps(ctx, deps.ListPurchaseOrderLineItems, id, currency)
				pageData.LineItemTable = buildLineItemTable(lineItems, deps.TableLabels, currency, id, deps.Routes, isDraft, perms)
				if isDraft {
					pageData.LineItemAddURL = route.ResolveURL(deps.Routes.PurchaseOrderLineItemAddURL, "id", id)
				}
				if cell, ok := po["total_amount"].(types.TableCell); ok {
					pageData.TotalAmount = cell.Currency + " " + cell.Value
				}
			}
		}

		return view.OK("purchase-order-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadPurchaseOrder(ctx, &purchaseorderpb.ReadPurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read purchase order %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load purchase order: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("purchase order not found"))
		}
		po := purchaseOrderToMap(data[0])

		tabPoStatus, _ := po["status"].(string)
		tabConfirmReceiptURL := ""
		if tabPoStatus == "approved" || tabPoStatus == "partially_received" {
			tabConfirmReceiptURL = route.ResolveURL(deps.Routes.PurchaseOrderConfirmReceiptURL, "id", id)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			PurchaseOrder:     po,
			Labels:            deps.Labels,
			ActiveTab:         tab,
			TabItems:          buildTabItems(deps.Labels, id, deps.Routes),
			SetStatusURL:      deps.Routes.PurchaseOrderSetStatusURL,
			ConfirmReceiptURL: tabConfirmReceiptURL,
		}

		switch tab {
		case "info":
			// po map has everything
		case "items":
			if deps.ListPurchaseOrderLineItems != nil {
				perms := view.GetUserPermissions(ctx)
				currency, _ := po["currency"].(string)
				status, _ := po["status"].(string)
				isDraft := status == "draft"
				lineItems := listLineItemMaps(ctx, deps.ListPurchaseOrderLineItems, id, currency)
				pageData.LineItemTable = buildLineItemTable(lineItems, deps.TableLabels, currency, id, deps.Routes, isDraft, perms)
				if isDraft {
					pageData.LineItemAddURL = route.ResolveURL(deps.Routes.PurchaseOrderLineItemAddURL, "id", id)
				}
				if cell, ok := po["total_amount"].(types.TableCell); ok {
					pageData.TotalAmount = cell.Currency + " " + cell.Value
				}
			}
		}

		templateName := "purchase-order-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// NewLineItemTableView returns a view that renders only the PO line items table (HTMX refresh).
func NewLineItemTableView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		purchaseOrderID := viewCtx.Request.PathValue("id")

		// Re-read PO to get status and currency
		var currency string
		var isDraft bool
		if deps.ReadPurchaseOrder != nil {
			resp, err := deps.ReadPurchaseOrder(ctx, &purchaseorderpb.ReadPurchaseOrderRequest{
				Data: &purchaseorderpb.PurchaseOrder{Id: purchaseOrderID},
			})
			if err == nil && len(resp.GetData()) > 0 {
				po := resp.GetData()[0]
				currency = po.GetCurrency()
				isDraft = po.GetStatus() == "draft"
			}
		}

		lineItems := listLineItemMaps(ctx, deps.ListPurchaseOrderLineItems, purchaseOrderID, currency)
		perms := view.GetUserPermissions(ctx)
		table := buildLineItemTable(lineItems, deps.TableLabels, currency, purchaseOrderID, deps.Routes, isDraft, perms)
		return view.OK("table-card", table)
	})
}
