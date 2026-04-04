package checkout

import (
	"context"
	"fmt"
	"testing"

	inventoryItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	serialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialHistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
)

func TestReserveSerials(t *testing.T) {
	t.Parallel()

	t.Run("nil deps returns nil without panic", func(t *testing.T) {
		t.Parallel()

		deps := CheckoutDeps{} // all nil
		svc := NewService(deps)

		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("no serials available returns nil", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{Success: true, Data: nil}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   2,
		}})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("filters by available status and active flag", func(t *testing.T) {
		t.Parallel()

		var updatedIDs []string

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "reserved", Active: true},
					{Id: "s-2", InventoryItemId: "inv-001", SerialNumber: "SN-002", Status: "available", Active: false}, // inactive
					{Id: "s-3", InventoryItemId: "inv-001", SerialNumber: "SN-003", Status: "available", Active: true},  // should be picked
					{Id: "s-4", InventoryItemId: "inv-001", SerialNumber: "SN-004", Status: "sold", Active: true},
					{Id: "s-5", InventoryItemId: "inv-001", SerialNumber: "SN-005", Status: "available", Active: true}, // should be picked
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, req *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updatedIDs = append(updatedIDs, req.GetData().GetId())
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = func(_ context.Context, _ *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   2,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(updatedIDs) != 2 {
			t.Fatalf("expected 2 serial updates, got %d", len(updatedIDs))
		}
		if updatedIDs[0] != "s-3" {
			t.Errorf("first updated serial = %q, want %q", updatedIDs[0], "s-3")
		}
		if updatedIDs[1] != "s-5" {
			t.Errorf("second updated serial = %q, want %q", updatedIDs[1], "s-5")
		}
	})

	t.Run("fewer available serials than quantity reserves what is available", func(t *testing.T) {
		t.Parallel()

		var updateCount int

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "available", Active: true},
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updateCount++
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = func(_ context.Context, _ *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   5, // want 5 but only 1 available
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updateCount != 1 {
			t.Errorf("expected 1 serial update (only 1 available), got %d", updateCount)
		}
	})

	t.Run("no inventory item found skips serials silently", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListInventoryItems = func(_ context.Context, _ *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error) {
			return &inventoryItempb.ListInventoryItemsResponse{Success: true, Data: nil}, nil
		}

		serialCalled := false
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			serialCalled = true
			return &serialpb.ListInventorySerialsResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if serialCalled {
			t.Error("ListSerials should not be called when no inventory item found")
		}
	})

	t.Run("all serials non-available results in zero updates", func(t *testing.T) {
		t.Parallel()

		var updateCount int

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "sold", Active: true},
					{Id: "s-2", InventoryItemId: "inv-001", SerialNumber: "SN-002", Status: "reserved", Active: true},
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updateCount++
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   2,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updateCount != 0 {
			t.Errorf("expected 0 updates when no serials are available, got %d", updateCount)
		}
	})

	t.Run("creates serial history entries", func(t *testing.T) {
		t.Parallel()

		var historyCount int

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "available", Active: true},
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = func(_ context.Context, req *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			historyCount++
			data := req.GetData()
			if data.GetFromStatus() != "available" {
				t.Errorf("FromStatus = %q, want %q", data.GetFromStatus(), "available")
			}
			if data.GetToStatus() != "reserved" {
				t.Errorf("ToStatus = %q, want %q", data.GetToStatus(), "reserved")
			}
			if data.GetReferenceType() != "sale" {
				t.Errorf("ReferenceType = %q, want %q", data.GetReferenceType(), "sale")
			}
			if data.GetReferenceId() != "rev-001" {
				t.Errorf("ReferenceId = %q, want %q", data.GetReferenceId(), "rev-001")
			}
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if historyCount != 1 {
			t.Errorf("expected 1 history entry, got %d", historyCount)
		}
	})
}

// ---------------------------------------------------------------------------
// reserveSerials — defensive / edge cases
// ---------------------------------------------------------------------------

func TestReserveSerials_Defensive(t *testing.T) {
	t.Parallel()

	t.Run("empty serial list returns nil", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data:    []*serialpb.InventorySerial{}, // empty slice (not nil)
			}, nil
		}

		updateCalled := false
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updateCalled = true
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   2,
		}})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if updateCalled {
			t.Error("UpdateSerial should not be called when serial list is empty")
		}
	})

	t.Run("serial with empty ID is still processed", func(t *testing.T) {
		t.Parallel()

		var capturedID string
		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "available", Active: true},
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, req *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			capturedID = req.GetData().GetId()
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = func(_ context.Context, _ *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedID != "" {
			t.Errorf("expected empty ID to be passed through, got %q", capturedID)
		}
	})

	t.Run("serial with empty status is not picked as available", func(t *testing.T) {
		t.Parallel()

		updateCalled := false
		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "", Active: true},
				},
			}, nil
		}
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updateCalled = true
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updateCalled {
			t.Error("serial with empty status should not be treated as available")
		}
	})

	t.Run("ListSerials error continues silently", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return nil, fmt.Errorf("serial service unavailable")
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("expected nil error (best-effort), got %v", err)
		}
	})

	t.Run("UpdateSerial failure continues with next serial", func(t *testing.T) {
		t.Parallel()

		var historyCount int
		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "available", Active: true},
					{Id: "s-2", InventoryItemId: "inv-001", SerialNumber: "SN-002", Status: "available", Active: true},
				},
			}, nil
		}
		callCount := 0
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			callCount++
			if callCount == 1 {
				return nil, fmt.Errorf("first update failed")
			}
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = func(_ context.Context, _ *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			historyCount++
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   2,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if callCount != 2 {
			t.Errorf("expected 2 update attempts, got %d", callCount)
		}
		// Only 1 history should be created (second serial succeeded)
		if historyCount != 1 {
			t.Errorf("expected 1 history entry (skipping failed), got %d", historyCount)
		}
	})

	t.Run("empty items list returns nil immediately", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		serialCalled := false
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			serialCalled = true
			return &serialpb.ListInventorySerialsResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", nil)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if serialCalled {
			t.Error("ListSerials should not be called with empty items list")
		}
	})

	t.Run("ListInventoryItems error skips silently", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListInventoryItems = func(_ context.Context, _ *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error) {
			return nil, fmt.Errorf("inventory service down")
		}

		serialCalled := false
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			serialCalled = true
			return &serialpb.ListInventorySerialsResponse{Success: true}, nil
		}

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if serialCalled {
			t.Error("ListSerials should not be called when ListInventoryItems fails")
		}
	})

	t.Run("nil CreateSerialHistory skips history creation", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListSerials = func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{
				Success: true,
				Data: []*serialpb.InventorySerial{
					{Id: "s-1", InventoryItemId: "inv-001", SerialNumber: "SN-001", Status: "available", Active: true},
				},
			}, nil
		}
		updateCalled := false
		deps.UpdateSerial = func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			updateCalled = true
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		}
		deps.CreateSerialHistory = nil // nil history creator

		svc := NewService(deps)
		err := svc.reserveSerials(context.Background(), "rev-001", []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !updateCalled {
			t.Error("serial should still be updated even without history creation")
		}
	})
}
