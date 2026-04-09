package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	purchaseorderlineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order_line_item"
)

// LineItemFormData is the template data for the PO line item drawer form.
type LineItemFormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	PurchaseOrderID string
	LineType        string
	Description     string
	ProductID       string
	InventoryItemID string
	LocationID      string
	QuantityOrdered string
	UnitPrice       string
	Notes           string
	Labels          centymo.ExpenditureLabels
	CommonLabels    any
}

// LineItemDeps holds dependencies for PO line item action handlers.
type LineItemDeps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	CreatePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.CreatePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.CreatePurchaseOrderLineItemResponse, error)
	ReadPurchaseOrderLineItem   func(ctx context.Context, req *purchaseorderlineitempb.ReadPurchaseOrderLineItemRequest) (*purchaseorderlineitempb.ReadPurchaseOrderLineItemResponse, error)
	UpdatePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.UpdatePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.UpdatePurchaseOrderLineItemResponse, error)
	DeletePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.DeletePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.DeletePurchaseOrderLineItemResponse, error)
}

// lineItemHTMXSuccess returns a success HTMX response that refreshes the PO line items table.
func lineItemHTMXSuccess() view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": fmt.Sprintf(`{"formSuccess":true,"refreshTable":"%s"}`, "po-line-items-table"),
		},
	}
}

// lineItemHTMXError returns an error HTMX response.
func lineItemHTMXError(message string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusUnprocessableEntity,
		Headers: map[string]string{
			"HX-Error-Message": message,
		},
	}
}

// NewLineItemAddAction creates the PO line item add action (GET = form, POST = create).
func NewLineItemAddAction(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		purchaseOrderID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("po-line-item-drawer-form", &LineItemFormData{
				FormAction:      route.ResolveURL(deps.Routes.PurchaseOrderLineItemAddURL, "id", purchaseOrderID),
				PurchaseOrderID: purchaseOrderID,
				LineType:        "goods",
				QuantityOrdered: "1",
				Labels:          deps.Labels,
				CommonLabels:    nil,
			})
		}

		// POST — create line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		quantityStr := r.FormValue("quantity_ordered")
		unitPriceStr := r.FormValue("unit_price")

		quantityF, _ := strconv.ParseFloat(quantityStr, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPriceStr, 64)
		if quantityF == 0 {
			quantityF = 1
		}
		totalPrice := quantityF * unitPriceF

		notes := r.FormValue("notes")
		productID := r.FormValue("product_id")
		inventoryItemID := r.FormValue("inventory_item_id")
		locationID := r.FormValue("location_id")

		lineItem := &purchaseorderlineitempb.PurchaseOrderLineItem{
			PurchaseOrderId:  purchaseOrderID,
			LineType:         r.FormValue("line_type"),
			Description:      r.FormValue("description"),
			QuantityOrdered:  quantityF,
			QuantityReceived: 0,
			QuantityBilled:   0,
			UnitPrice:        int64(math.Round(unitPriceF * 100)),
			TotalPrice:       int64(math.Round(totalPrice * 100)),
			Notes:            &notes,
		}
		if productID != "" {
			lineItem.ProductId = &productID
		}
		if inventoryItemID != "" {
			lineItem.InventoryItemId = &inventoryItemID
		}
		if locationID != "" {
			lineItem.LocationId = &locationID
		}

		_, err := deps.CreatePurchaseOrderLineItem(ctx, &purchaseorderlineitempb.CreatePurchaseOrderLineItemRequest{
			Data: lineItem,
		})
		if err != nil {
			log.Printf("Failed to create PO line item for PO %s: %v", purchaseOrderID, err)
			return lineItemHTMXError(err.Error())
		}

		return lineItemHTMXSuccess()
	})
}

// NewLineItemEditAction creates the PO line item edit action (GET = pre-filled form, POST = update).
func NewLineItemEditAction(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		purchaseOrderID := viewCtx.Request.PathValue("id")
		itemID := viewCtx.Request.PathValue("itemId")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadPurchaseOrderLineItem(ctx, &purchaseorderlineitempb.ReadPurchaseOrderLineItemRequest{
				Data: &purchaseorderlineitempb.PurchaseOrderLineItem{Id: itemID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				return lineItemHTMXError(deps.Labels.Errors.NotFound)
			}
			item := readResp.GetData()[0]

			return view.OK("po-line-item-drawer-form", &LineItemFormData{
				FormAction:      route.ResolveURL(deps.Routes.PurchaseOrderLineItemEditURL, "id", purchaseOrderID, "itemId", itemID),
				IsEdit:          true,
				ID:              itemID,
				PurchaseOrderID: purchaseOrderID,
				LineType:        item.GetLineType(),
				Description:     item.GetDescription(),
				ProductID:       item.GetProductId(),
				InventoryItemID: item.GetInventoryItemId(),
				LocationID:      item.GetLocationId(),
				QuantityOrdered: fmt.Sprintf("%.0f", item.GetQuantityOrdered()),
				UnitPrice:       fmt.Sprintf("%.2f", float64(item.GetUnitPrice())/100.0),
				Notes:           item.GetNotes(),
				Labels:          deps.Labels,
				CommonLabels:    nil,
			})
		}

		// POST — update line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		quantityStr := r.FormValue("quantity_ordered")
		unitPriceStr := r.FormValue("unit_price")

		quantityF, _ := strconv.ParseFloat(quantityStr, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPriceStr, 64)
		if quantityF == 0 {
			quantityF = 1
		}
		totalPrice := quantityF * unitPriceF

		notes := r.FormValue("notes")
		productID := r.FormValue("product_id")
		inventoryItemID := r.FormValue("inventory_item_id")
		locationID := r.FormValue("location_id")

		lineItem := &purchaseorderlineitempb.PurchaseOrderLineItem{
			Id:              itemID,
			LineType:        r.FormValue("line_type"),
			Description:     r.FormValue("description"),
			QuantityOrdered: quantityF,
			UnitPrice:       int64(math.Round(unitPriceF * 100)),
			TotalPrice:      int64(math.Round(totalPrice * 100)),
			Notes:           &notes,
		}
		if productID != "" {
			lineItem.ProductId = &productID
		}
		if inventoryItemID != "" {
			lineItem.InventoryItemId = &inventoryItemID
		}
		if locationID != "" {
			lineItem.LocationId = &locationID
		}

		_, err := deps.UpdatePurchaseOrderLineItem(ctx, &purchaseorderlineitempb.UpdatePurchaseOrderLineItemRequest{
			Data: lineItem,
		})
		if err != nil {
			log.Printf("Failed to update PO line item %s: %v", itemID, err)
			return lineItemHTMXError(err.Error())
		}

		return lineItemHTMXSuccess()
	})
}

// NewLineItemRemoveAction creates the PO line item remove action (POST only).
func NewLineItemRemoveAction(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		itemID := viewCtx.Request.URL.Query().Get("itemId")
		if itemID == "" {
			_ = viewCtx.Request.ParseForm()
			itemID = viewCtx.Request.FormValue("itemId")
		}
		if itemID == "" {
			return lineItemHTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeletePurchaseOrderLineItem(ctx, &purchaseorderlineitempb.DeletePurchaseOrderLineItemRequest{
			Data: &purchaseorderlineitempb.PurchaseOrderLineItem{Id: itemID},
		})
		if err != nil {
			log.Printf("Failed to delete PO line item %s: %v", itemID, err)
			return lineItemHTMXError(err.Error())
		}

		return lineItemHTMXSuccess()
	})
}
