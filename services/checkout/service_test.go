package checkout

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	inventoryItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	serialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialHistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	lineItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	paymentpb "github.com/erniealice/esqyma/pkg/schema/v1/integration/payment"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// mockDeps returns a CheckoutDeps with all fields wired to minimal happy-path
// stubs. Callers can override individual fields before passing to NewService.
func mockDeps() CheckoutDeps {
	return CheckoutDeps{
		CreateRevenue: func(_ context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			return &revenuepb.CreateRevenueResponse{
				Success: true,
				Data:    []*revenuepb.Revenue{{Id: "rev-001"}},
			}, nil
		},
		UpdateRevenue: func(_ context.Context, _ *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error) {
			return &revenuepb.UpdateRevenueResponse{Success: true}, nil
		},
		ListRevenues: func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			ref := "ORD-abcd-ef01"
			return &revenuepb.ListRevenuesResponse{
				Success: true,
				Data: []*revenuepb.Revenue{{
					Id:              "rev-001",
					ReferenceNumber: &ref,
					TotalAmount:     500.00, // unit form (= 50000 centavos)
					Status:          "paid",
				}},
			}, nil
		},
		CreateLineItem: func(_ context.Context, _ *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error) {
			return &lineItempb.CreateRevenueLineItemResponse{Success: true}, nil
		},
		ListLineItems: func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			pid := "prod-001"
			return &lineItempb.ListRevenueLineItemsResponse{
				Success: true,
				Data: []*lineItempb.RevenueLineItem{{
					ProductId:   &pid,
					Description: "Widget",
					Quantity:    2,
					UnitPrice:   100.00,
					TotalPrice:  200.00,
				}},
			}, nil
		},
		ListInventoryItems: func(_ context.Context, _ *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error) {
			return &inventoryItempb.ListInventoryItemsResponse{
				Success: true,
				Data: []*inventoryItempb.InventoryItem{{
					Id:                "inv-001",
					QuantityAvailable: 10,
					QuantityReserved:  0,
					QuantityOnHand:    10,
					Active:            true,
				}},
			}, nil
		},
		UpdateInventoryItem: func(_ context.Context, _ *inventoryItempb.UpdateInventoryItemRequest) (*inventoryItempb.UpdateInventoryItemResponse, error) {
			return &inventoryItempb.UpdateInventoryItemResponse{Success: true}, nil
		},
		ListSerials: func(_ context.Context, _ *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error) {
			return &serialpb.ListInventorySerialsResponse{Success: true}, nil
		},
		UpdateSerial: func(_ context.Context, _ *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error) {
			return &serialpb.UpdateInventorySerialResponse{Success: true}, nil
		},
		CreateSerialHistory: func(_ context.Context, _ *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error) {
			return &serialHistorypb.CreateInventorySerialHistoryResponse{Success: true}, nil
		},
		ProcessWebhook: func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{
				Success: true,
				Data: []*paymentpb.WebhookResult{{
					PaymentId: "rev-001",
					Status:    paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCESS,
					Action:    "payment.success",
				}},
			}, nil
		},
	}
}

func sampleRequest() CheckoutRequest {
	return CheckoutRequest{
		ClientID:        "client-001",
		CustomerEmail:   "test@example.com",
		CustomerName:    "Test User",
		CustomerPhone:   "+639001234567",
		Currency:        "PHP",
		FulfillmentType: "store_pickup",
		LocationID:      "loc-001",
		TotalAmount:     20000, // 200.00 PHP
		Items: []CheckoutItem{
			{
				ProductID:    "prod-001",
				ProductName:  "Widget",
				VariantID:    "var-001",
				VariantLabel: "Red",
				LocationID:   "loc-001",
				Quantity:     2,
				UnitPrice:    10000,
				TotalPrice:   20000,
				CostPrice:    50.00,
				PriceListID:  "pl-001",
			},
		},
	}
}

// ---------------------------------------------------------------------------
// generateRefNumber
// ---------------------------------------------------------------------------

func TestGenerateRefNumber(t *testing.T) {
	t.Parallel()

	re := regexp.MustCompile(`^ORD-[0-9a-f]{4}-[0-9a-f]{4}$`)

	for i := 0; i < 20; i++ {
		ref, err := generateRefNumber()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !re.MatchString(ref) {
			t.Errorf("generateRefNumber() = %q, does not match ORD-XXXX-XXXX pattern", ref)
		}
	}

	// Verify uniqueness across a batch (probabilistic but extremely unlikely to collide with 4 random bytes).
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		ref, _ := generateRefNumber()
		if seen[ref] {
			t.Errorf("duplicate ref number: %s", ref)
		}
		seen[ref] = true
	}
}

// ---------------------------------------------------------------------------
// PlaceOrder
// ---------------------------------------------------------------------------

func TestPlaceOrder(t *testing.T) {
	t.Parallel()

	t.Run("happy path without payment provider", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		svc := NewService(deps)
		req := sampleRequest()
		req.PaymentProvider = "" // no payment

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RevenueID != "rev-001" {
			t.Errorf("RevenueID = %q, want %q", result.RevenueID, "rev-001")
		}
		if result.Status != "pending" {
			t.Errorf("Status = %q, want %q", result.Status, "pending")
		}
		if result.TotalAmount != 20000 {
			t.Errorf("TotalAmount = %d, want %d", result.TotalAmount, 20000)
		}
		if result.CheckoutURL != "" {
			t.Errorf("CheckoutURL should be empty without payment provider, got %q", result.CheckoutURL)
		}

		re := regexp.MustCompile(`^ORD-[0-9a-f]{4}-[0-9a-f]{4}$`)
		if !re.MatchString(result.ReferenceNumber) {
			t.Errorf("ReferenceNumber = %q, does not match expected pattern", result.ReferenceNumber)
		}
	})

	t.Run("happy path with payment provider", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateCheckoutSession = func(_ context.Context, req *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error) {
			return &paymentpb.CreateCheckoutSessionResponse{
				Success: true,
				Data: []*paymentpb.CheckoutSession{{
					Id:          "session-001",
					CheckoutUrl: "https://pay.example.com/checkout/session-001",
				}},
			}, nil
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.PaymentProvider = "maya"
		req.SuccessURL = "https://shop.example.com/success"
		req.FailureURL = "https://shop.example.com/fail"
		req.CancelURL = "https://shop.example.com/cancel"

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.CheckoutURL != "https://pay.example.com/checkout/session-001" {
			t.Errorf("CheckoutURL = %q, want payment URL", result.CheckoutURL)
		}
		if result.CheckoutID != "session-001" {
			t.Errorf("CheckoutID = %q, want %q", result.CheckoutID, "session-001")
		}
	})

	t.Run("create revenue fails", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateRevenue = func(_ context.Context, _ *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			return nil, fmt.Errorf("db connection lost")
		}
		svc := NewService(deps)

		_, err := svc.PlaceOrder(context.Background(), sampleRequest())
		if err == nil {
			t.Fatal("expected error when CreateRevenue fails")
		}
	})

	t.Run("create revenue returns empty data", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateRevenue = func(_ context.Context, _ *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			return &revenuepb.CreateRevenueResponse{Success: true, Data: nil}, nil
		}
		svc := NewService(deps)

		_, err := svc.PlaceOrder(context.Background(), sampleRequest())
		if err == nil {
			t.Fatal("expected error when CreateRevenue returns empty data")
		}
	})

	t.Run("payment session failure propagates error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateCheckoutSession = func(_ context.Context, _ *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error) {
			return nil, fmt.Errorf("payment gateway down")
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.PaymentProvider = "maya"

		_, err := svc.PlaceOrder(context.Background(), req)
		if err == nil {
			t.Fatal("expected error when CreateCheckoutSession fails")
		}
	})

	t.Run("line item creation failure does not block checkout", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateLineItem = func(_ context.Context, _ *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error) {
			return nil, fmt.Errorf("line item write failed")
		}
		svc := NewService(deps)

		result, err := svc.PlaceOrder(context.Background(), sampleRequest())
		if err != nil {
			t.Fatalf("line item failure should not block checkout: %v", err)
		}
		if result.RevenueID == "" {
			t.Error("expected non-empty RevenueID despite line item failure")
		}
	})
}

// ---------------------------------------------------------------------------
// reserveStock
// ---------------------------------------------------------------------------

func TestReserveStock(t *testing.T) {
	t.Parallel()

	t.Run("decrements available and increments reserved", func(t *testing.T) {
		t.Parallel()

		var capturedUpdate *inventoryItempb.UpdateInventoryItemRequest

		deps := mockDeps()
		deps.ListInventoryItems = func(_ context.Context, _ *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error) {
			return &inventoryItempb.ListInventoryItemsResponse{
				Success: true,
				Data: []*inventoryItempb.InventoryItem{{
					Id:                "inv-001",
					QuantityAvailable: 10,
					QuantityReserved:  2,
					QuantityOnHand:    12,
					Active:            true,
				}},
			}, nil
		}
		deps.UpdateInventoryItem = func(_ context.Context, req *inventoryItempb.UpdateInventoryItemRequest) (*inventoryItempb.UpdateInventoryItemResponse, error) {
			capturedUpdate = req
			return &inventoryItempb.UpdateInventoryItemResponse{Success: true}, nil
		}

		svc := NewService(deps)
		items := []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   3,
		}}
		svc.reserveStock(context.Background(), items)

		if capturedUpdate == nil {
			t.Fatal("UpdateInventoryItem was not called")
		}
		data := capturedUpdate.GetData()
		if data.GetQuantityAvailable() != 7 {
			t.Errorf("QuantityAvailable = %v, want 7", data.GetQuantityAvailable())
		}
		if data.GetQuantityReserved() != 5 {
			t.Errorf("QuantityReserved = %v, want 5", data.GetQuantityReserved())
		}
	})

	t.Run("nil deps skips reservation", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListInventoryItems = nil
		deps.UpdateInventoryItem = nil
		svc := NewService(deps)

		// Should not panic.
		svc.reserveStock(context.Background(), []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})
	})

	t.Run("no inventory found continues silently", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListInventoryItems = func(_ context.Context, _ *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error) {
			return &inventoryItempb.ListInventoryItemsResponse{Success: true, Data: nil}, nil
		}

		updateCalled := false
		deps.UpdateInventoryItem = func(_ context.Context, _ *inventoryItempb.UpdateInventoryItemRequest) (*inventoryItempb.UpdateInventoryItemResponse, error) {
			updateCalled = true
			return &inventoryItempb.UpdateInventoryItemResponse{Success: true}, nil
		}

		svc := NewService(deps)
		svc.reserveStock(context.Background(), []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
		}})

		if updateCalled {
			t.Error("UpdateInventoryItem should not be called when no inventory items found")
		}
	})
}

// ---------------------------------------------------------------------------
// HandlePaymentWebhook
// ---------------------------------------------------------------------------

func TestHandlePaymentWebhook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		paymentStatus paymentpb.PaymentStatus
		wantStatus    string
	}{
		{name: "success maps to paid", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCESS, wantStatus: "paid"},
		{name: "failed maps to cancelled", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_FAILED, wantStatus: "cancelled"},
		{name: "cancelled maps to cancelled", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_CANCELLED, wantStatus: "cancelled"},
		{name: "expired maps to cancelled", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_EXPIRED, wantStatus: "cancelled"},
		{name: "unspecified maps to pending", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED, wantStatus: "pending"},
		{name: "processing maps to pending", paymentStatus: paymentpb.PaymentStatus_PAYMENT_STATUS_PROCESSING, wantStatus: "pending"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := mockDeps()
			deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
				return &paymentpb.ProcessWebhookResponse{
					Success: true,
					Data: []*paymentpb.WebhookResult{{
						PaymentId: "rev-001",
						Status:    tt.paymentStatus,
						Action:    "payment.test",
					}},
				}, nil
			}

			svc := NewService(deps)
			result, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", result.Status, tt.wantStatus)
			}
			if result.RevenueID != "rev-001" {
				t.Errorf("RevenueID = %q, want %q", result.RevenueID, "rev-001")
			}
		})
	}

	t.Run("nil ProcessWebhook returns error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = nil
		svc := NewService(deps)

		_, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err == nil {
			t.Fatal("expected error when ProcessWebhook is nil")
		}
	})

	t.Run("webhook processing failure returns error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return nil, fmt.Errorf("provider timeout")
		}
		svc := NewService(deps)

		_, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err == nil {
			t.Fatal("expected error on webhook processing failure")
		}
	})

	t.Run("unsuccessful response returns error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{Success: false, Data: nil}, nil
		}
		svc := NewService(deps)

		_, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err == nil {
			t.Fatal("expected error on unsuccessful webhook response")
		}
	})
}

// ---------------------------------------------------------------------------
// PlaceOrder — negative / defensive
// ---------------------------------------------------------------------------

func TestPlaceOrder_Defensive(t *testing.T) {
	t.Parallel()

	t.Run("empty items list succeeds (order with no line items)", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		svc := NewService(deps)
		req := sampleRequest()
		req.Items = nil

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RevenueID == "" {
			t.Error("expected non-empty RevenueID")
		}
	})

	t.Run("zero quantity items still creates line items", func(t *testing.T) {
		t.Parallel()

		var lineItemCount int
		deps := mockDeps()
		deps.CreateLineItem = func(_ context.Context, req *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error) {
			lineItemCount++
			if req.GetData().GetQuantity() != 0 {
				t.Errorf("expected Quantity = 0, got %v", req.GetData().GetQuantity())
			}
			return &lineItempb.CreateRevenueLineItemResponse{Success: true}, nil
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.Items = []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   0,
			UnitPrice:  10000,
			TotalPrice: 0,
		}}

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RevenueID == "" {
			t.Error("expected non-empty RevenueID")
		}
		if lineItemCount != 1 {
			t.Errorf("expected 1 line item created, got %d", lineItemCount)
		}
	})

	t.Run("negative prices are passed through to revenue", func(t *testing.T) {
		t.Parallel()

		var capturedRevenue *revenuepb.Revenue
		deps := mockDeps()
		deps.CreateRevenue = func(_ context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			capturedRevenue = req.GetData()
			return &revenuepb.CreateRevenueResponse{
				Success: true,
				Data:    []*revenuepb.Revenue{{Id: "rev-neg"}},
			}, nil
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.TotalAmount = -5000 // negative centavos

		_, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedRevenue == nil {
			t.Fatal("CreateRevenue was not called")
		}
		if capturedRevenue.GetTotalAmount() != -50.0 {
			t.Errorf("TotalAmount = %v, want -50.0 (unit form of -5000 centavos)", capturedRevenue.GetTotalAmount())
		}
	})

	t.Run("negative unit price on item is passed through", func(t *testing.T) {
		t.Parallel()

		var capturedLineItem *lineItempb.RevenueLineItem
		deps := mockDeps()
		deps.CreateLineItem = func(_ context.Context, req *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error) {
			capturedLineItem = req.GetData()
			return &lineItempb.CreateRevenueLineItemResponse{Success: true}, nil
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.Items = []CheckoutItem{{
			ProductID:  "prod-001",
			LocationID: "loc-001",
			Quantity:   1,
			UnitPrice:  -500,
			TotalPrice: -500,
		}}

		_, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedLineItem == nil {
			t.Fatal("CreateLineItem was not called")
		}
		if capturedLineItem.GetUnitPrice() != -5.0 {
			t.Errorf("UnitPrice = %v, want -5.0", capturedLineItem.GetUnitPrice())
		}
	})

	t.Run("empty client and customer fields are accepted", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		svc := NewService(deps)
		req := CheckoutRequest{
			Items: []CheckoutItem{{
				ProductID:  "prod-001",
				LocationID: "loc-001",
				Quantity:   1,
				UnitPrice:  100,
				TotalPrice: 100,
			}},
		}

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RevenueID == "" {
			t.Error("expected non-empty RevenueID even with empty customer fields")
		}
	})

	t.Run("payment session returns empty data still succeeds", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateCheckoutSession = func(_ context.Context, _ *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error) {
			return &paymentpb.CreateCheckoutSessionResponse{
				Success: true,
				Data:    nil, // empty data
			}, nil
		}
		svc := NewService(deps)
		req := sampleRequest()
		req.PaymentProvider = "maya"

		result, err := svc.PlaceOrder(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.CheckoutURL != "" {
			t.Errorf("expected empty CheckoutURL when session data is empty, got %q", result.CheckoutURL)
		}
		if result.CheckoutID != "" {
			t.Errorf("expected empty CheckoutID when session data is empty, got %q", result.CheckoutID)
		}
	})

	t.Run("create revenue returns success false", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.CreateRevenue = func(_ context.Context, _ *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			return &revenuepb.CreateRevenueResponse{Success: false, Data: nil}, nil
		}
		svc := NewService(deps)

		_, err := svc.PlaceOrder(context.Background(), sampleRequest())
		if err == nil {
			t.Fatal("expected error when CreateRevenue returns success=false")
		}
	})
}

// ---------------------------------------------------------------------------
// HandlePaymentWebhook — negative / defensive
// ---------------------------------------------------------------------------

func TestHandlePaymentWebhook_Defensive(t *testing.T) {
	t.Parallel()

	t.Run("empty webhook data returns error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{
				Success: true,
				Data:    nil, // no results
			}, nil
		}
		svc := NewService(deps)

		_, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err == nil {
			t.Fatal("expected error when webhook returns empty data")
		}
	})

	t.Run("webhook result with empty payment ID still returns result", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{
				Success: true,
				Data: []*paymentpb.WebhookResult{{
					PaymentId: "",
					Status:    paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCESS,
					Action:    "payment.success",
				}},
			}, nil
		}

		// When payment ID is empty and status is "paid", UpdateRevenue should be
		// skipped (guard: revenueID != "" && revenueStatus != "pending")
		updateCalled := false
		deps.UpdateRevenue = func(_ context.Context, _ *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error) {
			updateCalled = true
			return &revenuepb.UpdateRevenueResponse{Success: true}, nil
		}

		svc := NewService(deps)
		result, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RevenueID != "" {
			t.Errorf("expected empty RevenueID, got %q", result.RevenueID)
		}
		if result.Status != "paid" {
			t.Errorf("Status = %q, want %q", result.Status, "paid")
		}
		if updateCalled {
			t.Error("UpdateRevenue should not be called when payment ID is empty")
		}
	})

	t.Run("update revenue failure is non-fatal", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{
				Success: true,
				Data: []*paymentpb.WebhookResult{{
					PaymentId: "rev-fail",
					Status:    paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCESS,
					Action:    "payment.success",
				}},
			}, nil
		}
		deps.UpdateRevenue = func(_ context.Context, _ *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error) {
			return nil, fmt.Errorf("update failed")
		}

		svc := NewService(deps)
		result, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err != nil {
			t.Fatalf("update revenue failure should not block webhook result: %v", err)
		}
		if result.Status != "paid" {
			t.Errorf("Status = %q, want %q", result.Status, "paid")
		}
	})

	t.Run("pending status skips revenue update", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ProcessWebhook = func(_ context.Context, _ *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error) {
			return &paymentpb.ProcessWebhookResponse{
				Success: true,
				Data: []*paymentpb.WebhookResult{{
					PaymentId: "rev-pending",
					Status:    paymentpb.PaymentStatus_PAYMENT_STATUS_PROCESSING,
					Action:    "payment.processing",
				}},
			}, nil
		}
		updateCalled := false
		deps.UpdateRevenue = func(_ context.Context, _ *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error) {
			updateCalled = true
			return &revenuepb.UpdateRevenueResponse{Success: true}, nil
		}

		svc := NewService(deps)
		result, err := svc.HandlePaymentWebhook(context.Background(), &paymentpb.ProcessWebhookRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != "pending" {
			t.Errorf("Status = %q, want %q", result.Status, "pending")
		}
		if updateCalled {
			t.Error("UpdateRevenue should not be called for pending status")
		}
	})
}

// ---------------------------------------------------------------------------
// GetOrder
// ---------------------------------------------------------------------------

func TestGetOrder(t *testing.T) {
	t.Parallel()

	t.Run("returns order with centavo conversion", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		ref := "ORD-abcd-ef01"
		deps.ListRevenues = func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			return &revenuepb.ListRevenuesResponse{
				Success: true,
				Data: []*revenuepb.Revenue{{
					Id:              "rev-100",
					ReferenceNumber: &ref,
					TotalAmount:     500.00, // unit form -> 50000 centavos
					Status:          "paid",
				}},
			}, nil
		}
		pid := "prod-A"
		deps.ListLineItems = func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			return &lineItempb.ListRevenueLineItemsResponse{
				Success: true,
				Data: []*lineItempb.RevenueLineItem{{
					ProductId:   &pid,
					Description: "Gadget",
					Quantity:    3,
					UnitPrice:   100.50, // unit form
					TotalPrice:  301.50, // unit form
				}},
			}, nil
		}

		svc := NewService(deps)
		order, err := svc.GetOrder(context.Background(), ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order.RevenueID != "rev-100" {
			t.Errorf("RevenueID = %q, want %q", order.RevenueID, "rev-100")
		}
		if order.ReferenceNumber != ref {
			t.Errorf("ReferenceNumber = %q, want %q", order.ReferenceNumber, ref)
		}
		if order.TotalAmount != 50000 {
			t.Errorf("TotalAmount = %d, want 50000 (centavos)", order.TotalAmount)
		}
		if order.Status != "paid" {
			t.Errorf("Status = %q, want %q", order.Status, "paid")
		}
		if len(order.Items) != 1 {
			t.Fatalf("len(Items) = %d, want 1", len(order.Items))
		}
		item := order.Items[0]
		if item.UnitPrice != 10050 {
			t.Errorf("item.UnitPrice = %d, want 10050 (centavos)", item.UnitPrice)
		}
		if item.TotalPrice != 30150 {
			t.Errorf("item.TotalPrice = %d, want 30150 (centavos)", item.TotalPrice)
		}
		if item.Quantity != 3 {
			t.Errorf("item.Quantity = %d, want 3", item.Quantity)
		}
	})

	t.Run("order not found", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListRevenues = func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			return &revenuepb.ListRevenuesResponse{Success: true, Data: nil}, nil
		}
		svc := NewService(deps)

		_, err := svc.GetOrder(context.Background(), "ORD-0000-0000")
		if err == nil {
			t.Fatal("expected error for missing order")
		}
	})

	t.Run("list revenues error propagates", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListRevenues = func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			return nil, fmt.Errorf("network error")
		}
		svc := NewService(deps)

		_, err := svc.GetOrder(context.Background(), "ORD-1111-2222")
		if err == nil {
			t.Fatal("expected error on list revenues failure")
		}
	})

	t.Run("empty line items returns empty slice", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListLineItems = func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			return &lineItempb.ListRevenueLineItemsResponse{Success: true, Data: nil}, nil
		}
		svc := NewService(deps)

		order, err := svc.GetOrder(context.Background(), "ORD-abcd-ef01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(order.Items))
		}
	})

	t.Run("line items with nil product ID uses zero value", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListLineItems = func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			return &lineItempb.ListRevenueLineItemsResponse{
				Success: true,
				Data: []*lineItempb.RevenueLineItem{{
					ProductId:   nil, // nil ProductId
					Description: "",
					Quantity:    0,
					UnitPrice:   0,
					TotalPrice:  0,
				}},
			}, nil
		}

		svc := NewService(deps)
		order, err := svc.GetOrder(context.Background(), "ORD-abcd-ef01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(order.Items))
		}
		if order.Items[0].ProductID != "" {
			t.Errorf("expected empty ProductID for nil ProductId, got %q", order.Items[0].ProductID)
		}
		if order.Items[0].Quantity != 0 {
			t.Errorf("expected 0 quantity, got %d", order.Items[0].Quantity)
		}
	})

	t.Run("line items error propagates", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListLineItems = func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			return nil, fmt.Errorf("line items service down")
		}

		svc := NewService(deps)
		_, err := svc.GetOrder(context.Background(), "ORD-abcd-ef01")
		if err == nil {
			t.Fatal("expected error when ListLineItems fails")
		}
	})

	t.Run("revenue with nil reference number returns empty string", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListRevenues = func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			return &revenuepb.ListRevenuesResponse{
				Success: true,
				Data: []*revenuepb.Revenue{{
					Id:              "rev-nil-ref",
					ReferenceNumber: nil,
					TotalAmount:     0,
					Status:          "pending",
				}},
			}, nil
		}
		deps.ListLineItems = func(_ context.Context, _ *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error) {
			return &lineItempb.ListRevenueLineItemsResponse{Success: true, Data: nil}, nil
		}

		svc := NewService(deps)
		order, err := svc.GetOrder(context.Background(), "ORD-abcd-ef01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order.ReferenceNumber != "" {
			t.Errorf("expected empty ReferenceNumber for nil pointer, got %q", order.ReferenceNumber)
		}
	})

	t.Run("list revenues success false returns error", func(t *testing.T) {
		t.Parallel()

		deps := mockDeps()
		deps.ListRevenues = func(_ context.Context, _ *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error) {
			return &revenuepb.ListRevenuesResponse{Success: false, Data: nil}, nil
		}

		svc := NewService(deps)
		_, err := svc.GetOrder(context.Background(), "ORD-1111-2222")
		if err == nil {
			t.Fatal("expected error when ListRevenues returns success=false")
		}
	})
}
