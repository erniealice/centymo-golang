package checkout

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"time"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	inventoryItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	lineItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	paymentpb "github.com/erniealice/esqyma/pkg/schema/v1/integration/payment"
)

// Service orchestrates the checkout flow using typed espyna use case functions.
type Service struct {
	deps CheckoutDeps
}

// NewService creates a new checkout service.
func NewService(deps CheckoutDeps) *Service {
	return &Service{deps: deps}
}

// generateRefNumber generates a reference number in ORD-XXXX-XXXX format.
func generateRefNumber() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate reference number: %w", err)
	}
	h := hex.EncodeToString(b)
	return fmt.Sprintf("ORD-%s-%s", h[:4], h[4:]), nil
}

// generateID generates a random ID string using crypto/rand.
func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	// Format as UUID-like: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	h := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s-%s-%s-%s", h[:8], h[8:12], h[12:16], h[16:20], h[20:32]), nil
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

// PlaceOrder orchestrates the full checkout flow:
// 1. Generate reference number
// 2. Create Revenue record
// 3. Create RevenueLineItems
// 4. Reserve inventory stock
// 5. Reserve serials (best-effort)
// 6. Create payment session (if payment provider set)
// 7. Return checkout result
func (s *Service) PlaceOrder(ctx context.Context, req CheckoutRequest) (*CheckoutResult, error) {
	// 1. Generate reference number
	refNum, err := generateRefNumber()
	if err != nil {
		return nil, fmt.Errorf("checkout: %w", err)
	}

	now := time.Now()
	nowMillis := now.UnixMilli()
	nowStr := now.Format(time.RFC3339)

	// 2. Create Revenue record
	createRevenueResp, err := s.deps.CreateRevenue(ctx, &revenuepb.CreateRevenueRequest{
		Data: &revenuepb.Revenue{
			Active:             true,
			Name:               "Order " + refNum,
			ClientId:           req.ClientID,
			RevenueDate:        &nowMillis,
			RevenueDateString:  &nowStr,
			DateCreated:        &nowMillis,
			DateCreatedString:  &nowStr,
			DateModified:       &nowMillis,
			DateModifiedString: &nowStr,
			TotalAmount:        float64(req.TotalAmount) / 100.0,
			Currency:           req.Currency,
			Status:             "pending",
			ReferenceNumber:    &refNum,
			LocationId:         req.LocationID,
			PaymentProvider:    &req.PaymentProvider,
			FulfillmentType:    &req.FulfillmentType,
			DeliveryAddress:    &req.DeliveryAddress,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("checkout: create revenue: %w", err)
	}
	if !createRevenueResp.GetSuccess() || len(createRevenueResp.GetData()) == 0 {
		return nil, fmt.Errorf("checkout: create revenue: unexpected empty response")
	}

	revenueID := createRevenueResp.GetData()[0].GetId()

	// 3. Create RevenueLineItem for each CheckoutItem
	for _, item := range req.Items {
		_, err := s.deps.CreateLineItem(ctx, &lineItempb.CreateRevenueLineItemRequest{
			Data: &lineItempb.RevenueLineItem{
				Active:             true,
				RevenueId:          revenueID,
				ProductId:          &item.ProductID,
				Description:        item.ProductName,
				Quantity:           float64(item.Quantity),
				UnitPrice:          float64(item.UnitPrice) / 100.0,
				TotalPrice:         float64(item.TotalPrice) / 100.0,
				LineItemType:       "item",
				PriceListId:        &item.PriceListID,
				VariantId:          &item.VariantID,
				VariantLabel:       &item.VariantLabel,
				LocationId:         &item.LocationID,
				CostPrice:          &item.CostPrice,
				DateCreated:        &nowMillis,
				DateCreatedString:  &nowStr,
				DateModified:       &nowMillis,
				DateModifiedString: &nowStr,
			},
		})
		if err != nil {
			log.Printf("checkout: create line item for product %s: %v", item.ProductID, err)
			// Continue with other items — line item creation failures are logged but don't block checkout
		}
	}

	// 4. Reserve stock via UpdateInventoryItem
	s.reserveStock(ctx, req.Items)

	// 5. Reserve serials (best-effort)
	if err := s.reserveSerials(ctx, revenueID, req.Items); err != nil {
		log.Printf("checkout: reserve serials: %v", err)
	}

	result := &CheckoutResult{
		RevenueID:       revenueID,
		ReferenceNumber: refNum,
		TotalAmount:     req.TotalAmount,
		Status:          "pending",
	}

	// 6. Create payment session if provider is set
	if s.deps.CreateCheckoutSession != nil && req.PaymentProvider != "" {
		sessionResp, err := s.deps.CreateCheckoutSession(ctx, &paymentpb.CreateCheckoutSessionRequest{
			Data: &paymentpb.CheckoutSessionData{
				Amount:      float64(req.TotalAmount),
				Currency:    req.Currency,
				Description: "Order " + refNum,
				PaymentId:   revenueID,
				OrderRef:    refNum,
				SuccessUrl:  req.SuccessURL,
				FailureUrl:  req.FailureURL,
				CancelUrl:   req.CancelURL,
				Customer: &paymentpb.CustomerInfo{
					Email: req.CustomerEmail,
					Name:  req.CustomerName,
					Phone: req.CustomerPhone,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("checkout: create payment session: %w", err)
		}

		if sessionResp.GetSuccess() && len(sessionResp.GetData()) > 0 {
			session := sessionResp.GetData()[0]
			result.CheckoutURL = session.GetCheckoutUrl()
			result.CheckoutID = session.GetId()

			// Update revenue with checkout session ID
			checkoutSessionID := session.GetId()
			_, updateErr := s.deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{
					Id:                revenueID,
					CheckoutSessionId: &checkoutSessionID,
					DateModified:      ptr(time.Now().UnixMilli()),
					DateModifiedString: ptr(time.Now().Format(time.RFC3339)),
				},
			})
			if updateErr != nil {
				log.Printf("checkout: update revenue with checkout session: %v", updateErr)
			}
		}
	}

	return result, nil
}

// reserveStock decrements quantity_available and increments quantity_reserved
// for each checkout item's inventory.
func (s *Service) reserveStock(ctx context.Context, items []CheckoutItem) {
	if s.deps.ListInventoryItems == nil || s.deps.UpdateInventoryItem == nil {
		return
	}

	for _, item := range items {
		// List inventory items by product_id + location_id
		listResp, err := s.deps.ListInventoryItems(ctx, &inventoryItempb.ListInventoryItemsRequest{
			ProductId:  &item.ProductID,
			LocationId: &item.LocationID,
		})
		if err != nil {
			log.Printf("checkout: list inventory for product %s at location %s: %v", item.ProductID, item.LocationID, err)
			continue
		}
		if !listResp.GetSuccess() || len(listResp.GetData()) == 0 {
			continue
		}

		invItem := listResp.GetData()[0]
		newAvailable := invItem.GetQuantityAvailable() - float64(item.Quantity)
		newReserved := invItem.GetQuantityReserved() + float64(item.Quantity)

		_, err = s.deps.UpdateInventoryItem(ctx, &inventoryItempb.UpdateInventoryItemRequest{
			Data: &inventoryItempb.InventoryItem{
				Id:                invItem.GetId(),
				QuantityAvailable: newAvailable,
				QuantityReserved:  newReserved,
				QuantityOnHand:    invItem.GetQuantityOnHand(),
				Active:            invItem.GetActive(),
				DateModified:      ptr(time.Now().UnixMilli()),
				DateModifiedString: ptr(time.Now().Format(time.RFC3339)),
			},
		})
		if err != nil {
			log.Printf("checkout: update inventory for product %s: %v", item.ProductID, err)
		}
	}
}

// HandlePaymentWebhook processes a payment webhook and updates the revenue status.
func (s *Service) HandlePaymentWebhook(ctx context.Context, webhookReq *paymentpb.ProcessWebhookRequest) (*WebhookResult, error) {
	if s.deps.ProcessWebhook == nil {
		return nil, fmt.Errorf("checkout: webhook processing not configured")
	}

	resp, err := s.deps.ProcessWebhook(ctx, webhookReq)
	if err != nil {
		return nil, fmt.Errorf("checkout: process webhook: %w", err)
	}

	if !resp.GetSuccess() || len(resp.GetData()) == 0 {
		errMsg := "unknown error"
		if resp.GetError() != nil {
			errMsg = resp.GetError().GetMessage()
		}
		return nil, fmt.Errorf("checkout: webhook processing failed: %s", errMsg)
	}

	webhookData := resp.GetData()[0]
	revenueID := webhookData.GetPaymentId()

	// Map payment status to revenue status
	var revenueStatus string
	switch webhookData.GetStatus() {
	case paymentpb.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		revenueStatus = "paid"
	case paymentpb.PaymentStatus_PAYMENT_STATUS_FAILED:
		revenueStatus = "cancelled"
	case paymentpb.PaymentStatus_PAYMENT_STATUS_CANCELLED:
		revenueStatus = "cancelled"
	case paymentpb.PaymentStatus_PAYMENT_STATUS_EXPIRED:
		revenueStatus = "cancelled"
	default:
		revenueStatus = "pending"
	}

	// Update revenue status
	if revenueID != "" && revenueStatus != "pending" {
		_, err := s.deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
			Data: &revenuepb.Revenue{
				Id:                 revenueID,
				Status:             revenueStatus,
				DateModified:       ptr(time.Now().UnixMilli()),
				DateModifiedString: ptr(time.Now().Format(time.RFC3339)),
			},
		})
		if err != nil {
			log.Printf("checkout: update revenue status for %s: %v", revenueID, err)
		}
	}

	return &WebhookResult{
		RevenueID: revenueID,
		Status:    revenueStatus,
		PaymentID: webhookData.GetPaymentId(),
		Action:    webhookData.GetAction(),
	}, nil
}

// GetOrder retrieves a full order by reference number.
func (s *Service) GetOrder(ctx context.Context, referenceNumber string) (*OrderData, error) {
	// List revenues filtered by reference_number
	listResp, err := s.deps.ListRevenues(ctx, &revenuepb.ListRevenuesRequest{
		Filters: &commonpb.FilterRequest{
			Logic: commonpb.FilterLogic_AND,
			Filters: []*commonpb.TypedFilter{
				{
					Field: "reference_number",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    referenceNumber,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("checkout: list revenues: %w", err)
	}
	if !listResp.GetSuccess() || len(listResp.GetData()) == 0 {
		return nil, fmt.Errorf("checkout: order not found: %s", referenceNumber)
	}

	revenue := listResp.GetData()[0]
	revenueID := revenue.GetId()

	// List line items for this revenue
	lineItemsResp, err := s.deps.ListLineItems(ctx, &lineItempb.ListRevenueLineItemsRequest{
		RevenueId: &revenueID,
	})
	if err != nil {
		return nil, fmt.Errorf("checkout: list line items: %w", err)
	}

	var items []OrderLineItem
	if lineItemsResp.GetSuccess() {
		for _, li := range lineItemsResp.GetData() {
			items = append(items, OrderLineItem{
				ProductID:   li.GetProductId(),
				ProductName: li.GetDescription(),
				Quantity:    int(math.Round(li.GetQuantity())),
				UnitPrice:   int(math.Round(li.GetUnitPrice() * 100)),
				TotalPrice:  int(math.Round(li.GetTotalPrice() * 100)),
			})
		}
	}

	return &OrderData{
		RevenueID:       revenueID,
		ReferenceNumber: revenue.GetReferenceNumber(),
		TotalAmount:     int(math.Round(revenue.GetTotalAmount() * 100)),
		Status:          revenue.GetStatus(),
		Items:           items,
	}, nil
}
