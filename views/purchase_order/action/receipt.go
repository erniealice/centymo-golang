package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/enums"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorymovementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/inventory_movement"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	purchaseorderlineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order_line_item"
)

// ReceiptLineRow is a single line item row shown in the receipt form.
type ReceiptLineRow struct {
	ID              string
	Description     string
	LineType        string
	QuantityOrdered string
	QuantityReceived string
	Remaining       string
}

// ReceiptFormData is the template data for the PO receipt form.
type ReceiptFormData struct {
	FormAction      string
	PurchaseOrderID string
	Lines           []ReceiptLineRow
	Today           string
	LocationID      string
	Labels          centymo.ExpenditureLabels
	CommonLabels    any
}

// ReceiptActionDeps holds dependencies for the confirm-receipt action handler.
type ReceiptActionDeps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	ReadPurchaseOrder           func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	UpdatePurchaseOrder         func(ctx context.Context, req *purchaseorderpb.UpdatePurchaseOrderRequest) (*purchaseorderpb.UpdatePurchaseOrderResponse, error)
	ListPurchaseOrderLineItems  func(ctx context.Context, req *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error)
	ReadPurchaseOrderLineItem   func(ctx context.Context, req *purchaseorderlineitempb.ReadPurchaseOrderLineItemRequest) (*purchaseorderlineitempb.ReadPurchaseOrderLineItemResponse, error)
	UpdatePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.UpdatePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.UpdatePurchaseOrderLineItemResponse, error)
	CreateInventoryMovement     func(ctx context.Context, req *inventorymovementpb.CreateInventoryMovementRequest) (*inventorymovementpb.CreateInventoryMovementResponse, error)
	ReadInventoryItem           func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem         func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
}

// NewConfirmReceiptAction creates the PO confirm-receipt action (GET = form, POST = process).
func NewConfirmReceiptAction(deps *ReceiptActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return lineItemHTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		// Read the PO
		poResp, err := deps.ReadPurchaseOrder(ctx, &purchaseorderpb.ReadPurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{Id: id},
		})
		if err != nil || len(poResp.GetData()) == 0 {
			log.Printf("Failed to read purchase order %s for receipt: %v", id, err)
			return lineItemHTMXError(deps.Labels.Errors.NotFound)
		}
		po := poResp.GetData()[0]

		// Validate PO status
		status := po.GetStatus()
		if status != "approved" && status != "partially_received" {
			return lineItemHTMXError("Purchase order must be approved or partially received to confirm receipt.")
		}

		// List line items
		listResp, err := deps.ListPurchaseOrderLineItems(ctx, &purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest{
			PurchaseOrderId: strPtr(id),
		})
		if err != nil {
			log.Printf("Failed to list line items for PO %s: %v", id, err)
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Build unreceived lines (excluding expense and fully received)
			var lines []ReceiptLineRow
			for _, item := range listResp.GetData() {
				if item.GetPurchaseOrderId() != id {
					continue
				}
				if item.GetLineType() == "expense" {
					continue
				}
				remaining := item.GetQuantityOrdered() - item.GetQuantityReceived()
				if remaining <= 0 {
					continue
				}
				lines = append(lines, ReceiptLineRow{
					ID:               item.GetId(),
					Description:      item.GetDescription(),
					LineType:         item.GetLineType(),
					QuantityOrdered:  fmt.Sprintf("%.0f", item.GetQuantityOrdered()),
					QuantityReceived: fmt.Sprintf("%.0f", item.GetQuantityReceived()),
					Remaining:        fmt.Sprintf("%.0f", remaining),
				})
			}

			locationID := po.GetLocationId()

			return view.OK("po-receipt-form", &ReceiptFormData{
				FormAction:      route.ResolveURL(deps.Routes.PurchaseOrderConfirmReceiptURL, "id", id),
				PurchaseOrderID: id,
				Lines:           lines,
				Today:           time.Now().Format("2006-01-02"),
				LocationID:      locationID,
				Labels:          deps.Labels,
				CommonLabels:    nil,
			})
		}

		// POST — process the receipt
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		receiptLocationID := r.FormValue("location_id")
		receiptDate := r.FormValue("receipt_date")
		if receiptDate == "" {
			receiptDate = time.Now().Format("2006-01-02")
		}

		// Process each line item
		for _, item := range listResp.GetData() {
			if item.GetPurchaseOrderId() != id {
				continue
			}
			if item.GetLineType() == "expense" {
				continue
			}

			itemID := item.GetId()
			qtyStr := r.FormValue(fmt.Sprintf("line_%s_qty", itemID))
			if qtyStr == "" {
				continue
			}
			receivedQty, err := strconv.ParseFloat(qtyStr, 64)
			if err != nil || receivedQty <= 0 {
				continue
			}

			// For goods lines: create inventory movement + update inventory item
			if item.GetLineType() == "goods" && deps.CreateInventoryMovement != nil {
				// Determine location
				toLocationID := receiptLocationID
				if toLocationID == "" {
					toLocationID = item.GetLocationId()
				}

				// Create inventory movement
				referenceType := "purchase_order"
				referenceID := itemID
				performedBy := "system"
				movementData := &inventorymovementpb.InventoryMovement{
					MovementType:       enumspb.MovementType_MOVEMENT_TYPE_RECEIPT,
					ProductId:          item.GetProductId(),
					Quantity:           receivedQty,
					UnitCost:           item.GetUnitPrice(),
					MovementDateString: strPtr(receiptDate),
					ReferenceType:      &referenceType,
					ReferenceId:        &referenceID,
					PerformedBy:        &performedBy,
				}
				if inventoryItemID := item.GetInventoryItemId(); inventoryItemID != "" {
					movementData.InventoryItemId = &inventoryItemID
				}
				if toLocationID != "" {
					movementData.ToLocationId = &toLocationID
				}
				_, err := deps.CreateInventoryMovement(ctx, &inventorymovementpb.CreateInventoryMovementRequest{
					Data: movementData,
				})
				if err != nil {
					log.Printf("Failed to create inventory movement for line item %s: %v", itemID, err)
					return lineItemHTMXError(fmt.Sprintf("Failed to create inventory movement: %v", err))
				}

				// Update inventory item quantity
				if inventoryItemID := item.GetInventoryItemId(); inventoryItemID != "" && deps.ReadInventoryItem != nil {
					itemResp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
						Data: &inventoryitempb.InventoryItem{Id: inventoryItemID},
					})
					if err == nil && len(itemResp.GetData()) > 0 {
						invItem := itemResp.GetData()[0]
						newOnHand := invItem.GetQuantityOnHand() + receivedQty
						newAvailable := newOnHand - invItem.GetQuantityReserved()
						_, _ = deps.UpdateInventoryItem(ctx, &inventoryitempb.UpdateInventoryItemRequest{
							Data: &inventoryitempb.InventoryItem{
								Id:                 inventoryItemID,
								QuantityOnHand:     newOnHand,
								QuantityAvailable:  newAvailable,
							},
						})
					}
				}
			}

			// For goods and service lines: update poli.quantity_received
			newReceived := item.GetQuantityReceived() + receivedQty
			_, err = deps.UpdatePurchaseOrderLineItem(ctx, &purchaseorderlineitempb.UpdatePurchaseOrderLineItemRequest{
				Data: &purchaseorderlineitempb.PurchaseOrderLineItem{
					Id:               itemID,
					QuantityReceived: newReceived,
				},
			})
			if err != nil {
				log.Printf("Failed to update quantity_received for line item %s: %v", itemID, err)
				return lineItemHTMXError(fmt.Sprintf("Failed to update line item: %v", err))
			}
		}

		// Recalculate PO status based on updated line items
		updatedListResp, err := deps.ListPurchaseOrderLineItems(ctx, &purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest{
			PurchaseOrderId: strPtr(id),
		})
		if err == nil {
			allReceived := true
			anyReceived := false
			for _, item := range updatedListResp.GetData() {
				if item.GetPurchaseOrderId() != id {
					continue
				}
				if item.GetLineType() == "expense" {
					continue
				}
				if item.GetQuantityReceived() > 0 {
					anyReceived = true
				}
				if item.GetQuantityReceived() < item.GetQuantityOrdered() {
					allReceived = false
				}
			}

			newStatus := status
			if allReceived && anyReceived {
				newStatus = "fully_received"
			} else if anyReceived {
				newStatus = "partially_received"
			}

			if newStatus != status {
				_, err = deps.UpdatePurchaseOrder(ctx, &purchaseorderpb.UpdatePurchaseOrderRequest{
					Data: &purchaseorderpb.PurchaseOrder{Id: id, Status: newStatus},
				})
				if err != nil {
					log.Printf("Failed to update PO %s status to %s: %v", id, newStatus, err)
				}
			}
		}

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.PurchaseOrderDetailURL, "id", id),
			},
		}
	})
}
