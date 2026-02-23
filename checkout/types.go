package checkout

import (
	"context"

	inventoryItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	serialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialHistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	lineItempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	paymentpb "github.com/erniealice/esqyma/pkg/schema/v1/integration/payment"
)

// CheckoutDeps holds typed function references to espyna use cases.
// Consumer apps (retail-client) wire these from espyna container.GetUseCases().
type CheckoutDeps struct {
	// Revenue
	CreateRevenue func(ctx context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
	UpdateRevenue func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
	ReadRevenue   func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	ListRevenues  func(ctx context.Context, req *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error)

	// Revenue Line Items
	CreateLineItem func(ctx context.Context, req *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error)
	ListLineItems  func(ctx context.Context, req *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error)

	// Inventory (for stock reservation)
	UpdateInventoryItem func(ctx context.Context, req *inventoryItempb.UpdateInventoryItemRequest) (*inventoryItempb.UpdateInventoryItemResponse, error)
	ListInventoryItems  func(ctx context.Context, req *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error)

	// Inventory Serial (for serial assignment)
	ListSerials         func(ctx context.Context, req *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error)
	UpdateSerial        func(ctx context.Context, req *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error)
	CreateSerialHistory func(ctx context.Context, req *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error)

	// Payment (Maya integration)
	CreateCheckoutSession func(ctx context.Context, req *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error)
	ProcessWebhook        func(ctx context.Context, req *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error)
}

// CheckoutItem represents a single item in the checkout request.
type CheckoutItem struct {
	ProductID    string
	ProductName  string
	VariantID    string
	VariantLabel string
	LocationID   string
	Quantity     int
	UnitPrice    int // centavos
	TotalPrice   int // centavos
	CostPrice    float64
	PriceListID  string
}

// CheckoutRequest holds all data needed for PlaceOrder.
type CheckoutRequest struct {
	// Customer
	ClientID      string
	CustomerEmail string
	CustomerName  string
	CustomerPhone string
	// Order
	Items           []CheckoutItem
	TotalAmount     int    // centavos
	Currency        string // "PHP"
	FulfillmentType string // "store_pickup" or "home_delivery"
	LocationID      string // pickup branch
	DeliveryAddress string
	PaymentProvider string // "maya"
	// Payment redirect URLs
	SuccessURL string
	FailureURL string
	CancelURL  string
}

// CheckoutResult holds the result of PlaceOrder.
type CheckoutResult struct {
	RevenueID       string
	ReferenceNumber string
	CheckoutURL     string // Maya checkout redirect URL (empty if no payment provider)
	CheckoutID      string // Maya session ID
	TotalAmount     int    // centavos
	Status          string // "pending"
}

// WebhookResult holds the result of HandlePaymentWebhook.
type WebhookResult struct {
	RevenueID string
	Status    string // "paid", "cancelled"
	PaymentID string
	Action    string
}

// OrderData holds full order data for GetOrder.
type OrderData struct {
	RevenueID       string
	ReferenceNumber string
	TotalAmount     int // centavos
	Status          string
	Items           []OrderLineItem
}

// OrderLineItem represents a line item in order data.
type OrderLineItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   int // centavos
	TotalPrice  int // centavos
}
