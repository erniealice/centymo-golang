package checkout

import (
	"context"
	"fmt"
	"log"
	"time"

	inventoryItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	serialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialHistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
)

// reserveSerials attempts to reserve inventory serials for the given items.
// This is a best-effort operation -- if serial tracking is not available for a product,
// it falls back to quantity-only reservation (handled by reserveStock).
func (s *Service) reserveSerials(ctx context.Context, revenueID string, items []CheckoutItem) error {
	if s.deps.ListSerials == nil || s.deps.UpdateSerial == nil {
		return nil
	}

	for _, item := range items {
		if err := s.reserveSerialsForItem(ctx, revenueID, item); err != nil {
			log.Printf("checkout: reserve serials for product %s: %v", item.ProductID, err)
			// Continue with other items -- serial reservation failures don't block checkout
		}
	}
	return nil
}

// reserveSerialsForItem reserves serials for a single checkout item.
func (s *Service) reserveSerialsForItem(ctx context.Context, revenueID string, item CheckoutItem) error {
	// List available serials for this product at the given location.
	// We use the inventory_item_id filter indirectly -- first we need serials
	// that belong to inventory items for this product+location.
	// The ListSerials endpoint filters by inventory_item_id, so we first
	// need to find the inventory item ID.
	if s.deps.ListInventoryItems == nil {
		return nil
	}

	invResp, err := s.deps.ListInventoryItems(ctx, &inventoryItempb.ListInventoryItemsRequest{
		ProductId:  &item.ProductID,
		LocationId: &item.LocationID,
	})
	if err != nil || !invResp.GetSuccess() || len(invResp.GetData()) == 0 {
		return nil // No inventory item found -- product may not have inventory tracking
	}

	invItemID := invResp.GetData()[0].GetId()

	// List available serials for this inventory item
	serialResp, err := s.deps.ListSerials(ctx, &serialpb.ListInventorySerialsRequest{
		InventoryItemId: &invItemID,
	})
	if err != nil || !serialResp.GetSuccess() || len(serialResp.GetData()) == 0 {
		return nil // No serials available -- product may not use serial tracking
	}

	// Filter to available serials only and select up to quantity
	var availableSerials []*serialpb.InventorySerial
	for _, serial := range serialResp.GetData() {
		if serial.GetStatus() == "available" && serial.GetActive() {
			availableSerials = append(availableSerials, serial)
			if len(availableSerials) >= item.Quantity {
				break
			}
		}
	}

	if len(availableSerials) == 0 {
		return nil
	}

	// Reserve each selected serial
	now := time.Now()
	nowMillis := now.UnixMilli()
	nowStr := now.Format(time.RFC3339)

	for _, serial := range availableSerials {
		// Update serial status to "reserved"
		_, err := s.deps.UpdateSerial(ctx, &serialpb.UpdateInventorySerialRequest{
			Data: &serialpb.InventorySerial{
				Id:                 serial.GetId(),
				InventoryItemId:    serial.GetInventoryItemId(),
				SerialNumber:       serial.GetSerialNumber(),
				Status:             "reserved",
				Active:             true,
				DateModified:       &nowMillis,
				DateModifiedString: &nowStr,
			},
		})
		if err != nil {
			log.Printf("checkout: update serial %s to reserved: %v", serial.GetId(), err)
			continue
		}

		// Create serial history entry
		if s.deps.CreateSerialHistory != nil {
			historyID, idErr := generateID()
			if idErr != nil {
				log.Printf("checkout: generate history id: %v", idErr)
				continue
			}

			_, err = s.deps.CreateSerialHistory(ctx, &serialHistorypb.CreateInventorySerialHistoryRequest{
				Data: &serialHistorypb.InventorySerialHistory{
					Id:                historyID,
					InventorySerialId: serial.GetId(),
					InventoryItemId:   serial.GetInventoryItemId(),
					FromStatus:        "available",
					ToStatus:          "reserved",
					ReferenceType:     "sale",
					ReferenceId:       revenueID,
					Notes:             fmt.Sprintf("Reserved for order %s", revenueID),
					DateCreated:       &nowMillis,
					DateCreatedString: &nowStr,
				},
			})
			if err != nil {
				log.Printf("checkout: create serial history for %s: %v", serial.GetId(), err)
			}
		}
	}

	return nil
}
