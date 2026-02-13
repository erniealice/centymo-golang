package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
)

// PaymentMethodOption represents a selectable payment/collection method.
type PaymentMethodOption struct {
	Value string
	Label string
}

// PaymentFormData is the template data for the payment drawer form.
type PaymentFormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	RevenueID          string
	CollectionMethodID string
	AmountPaid         string
	Currency           string
	ReferenceNumber    string
	Notes              string
	ReceivedBy         string
	ReceivedRole       string
	PaymentMethods     []PaymentMethodOption
	CommonLabels       any
}

// PaymentDeps holds dependencies for payment action handlers.
type PaymentDeps struct {
	DB centymo.DataSource
}

// loadCollectionMethods loads collection methods from the DB and returns them
// as select options.
func loadCollectionMethods(ctx context.Context, db centymo.DataSource) []PaymentMethodOption {
	methods, err := db.ListSimple(ctx, "collection_method")
	if err != nil {
		log.Printf("Failed to list collection methods: %v", err)
		return []PaymentMethodOption{}
	}

	options := make([]PaymentMethodOption, 0, len(methods))
	for _, m := range methods {
		id, _ := m["id"].(string)
		name, _ := m["name"].(string)
		if id == "" {
			continue
		}
		if name == "" {
			name = id
		}
		options = append(options, PaymentMethodOption{Value: id, Label: name})
	}
	return options
}

// NewPaymentAddAction creates the payment add action (GET = form, POST = create).
func NewPaymentAddAction(deps *PaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			methods := loadCollectionMethods(ctx, deps.DB)
			return view.OK("sales-payment-drawer-form", &PaymentFormData{
				FormAction:     fmt.Sprintf("/action/sales/detail/%s/payment/add", revenueID),
				RevenueID:      revenueID,
				Currency:       "PHP",
				PaymentMethods: methods,
				CommonLabels:   nil, // injected by ViewAdapter
			})
		}

		// POST — create payment
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		collectionMethodID := r.FormValue("collection_method_id")

		// Look up the method name for the payment_method column
		methodName := collectionMethodID
		if collectionMethodID != "" {
			method, err := deps.DB.Read(ctx, "collection_method", collectionMethodID)
			if err == nil {
				if name, ok := method["name"].(string); ok {
					methodName = name
				}
			}
		}

		data := map[string]any{
			"revenue_id":           revenueID,
			"payment_method":       methodName,
			"amount_paid":          r.FormValue("amount_paid"),
			"currency":             r.FormValue("currency"),
			"collection_method_id": collectionMethodID,
			"reference_number":     r.FormValue("reference_number"),
			"received_by":         r.FormValue("received_by"),
			"received_role":       r.FormValue("received_role"),
			"collection_type":     "sale",
			"status":              "completed",
			"notes":               r.FormValue("notes"),
		}

		_, err := deps.DB.Create(ctx, "revenue_payment", data)
		if err != nil {
			log.Printf("Failed to create payment for revenue %s: %v", revenueID, err)
			return centymo.HTMXError("Failed to record payment")
		}

		return centymo.HTMXSuccess("payment-table")
	})
}

// NewPaymentEditAction creates the payment edit action (GET = form, POST = update).
func NewPaymentEditAction(deps *PaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		revenueID := viewCtx.Request.PathValue("id")
		paymentID := viewCtx.Request.PathValue("pid")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "revenue_payment", paymentID)
			if err != nil {
				log.Printf("Failed to read payment %s: %v", paymentID, err)
				return centymo.HTMXError("Payment not found")
			}

			collectionMethodID, _ := record["collection_method_id"].(string)
			amountPaid, _ := record["amount_paid"].(string)
			currency, _ := record["currency"].(string)
			referenceNumber, _ := record["reference_number"].(string)
			notes, _ := record["notes"].(string)
			receivedBy, _ := record["received_by"].(string)
			receivedRole, _ := record["received_role"].(string)

			methods := loadCollectionMethods(ctx, deps.DB)
			return view.OK("sales-payment-drawer-form", &PaymentFormData{
				FormAction:         fmt.Sprintf("/action/sales/detail/%s/payment/edit/%s", revenueID, paymentID),
				IsEdit:             true,
				ID:                 paymentID,
				RevenueID:          revenueID,
				CollectionMethodID: collectionMethodID,
				AmountPaid:         amountPaid,
				Currency:           currency,
				ReferenceNumber:    referenceNumber,
				Notes:              notes,
				ReceivedBy:         receivedBy,
				ReceivedRole:       receivedRole,
				PaymentMethods:     methods,
				CommonLabels:       nil, // injected by ViewAdapter
			})
		}

		// POST — update payment
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		collectionMethodID := r.FormValue("collection_method_id")

		// Look up the method name for the payment_method column
		methodName := collectionMethodID
		if collectionMethodID != "" {
			method, err := deps.DB.Read(ctx, "collection_method", collectionMethodID)
			if err == nil {
				if name, ok := method["name"].(string); ok {
					methodName = name
				}
			}
		}

		data := map[string]any{
			"payment_method":       methodName,
			"amount_paid":          r.FormValue("amount_paid"),
			"currency":             r.FormValue("currency"),
			"collection_method_id": collectionMethodID,
			"reference_number":     r.FormValue("reference_number"),
			"received_role":        r.FormValue("received_role"),
			"notes":                r.FormValue("notes"),
			// Never change received_by on edit
		}

		_, err := deps.DB.Update(ctx, "revenue_payment", paymentID, data)
		if err != nil {
			log.Printf("Failed to update payment %s: %v", paymentID, err)
			return centymo.HTMXError("Failed to update payment")
		}

		return centymo.HTMXSuccess("payment-table")
	})
}

// NewPaymentRemoveAction creates the payment remove action (POST only).
func NewPaymentRemoveAction(deps *PaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Payment ID is required")
		}

		err := deps.DB.Delete(ctx, "revenue_payment", id)
		if err != nil {
			log.Printf("Failed to delete payment %s: %v", id, err)
			return centymo.HTMXError("Failed to remove payment")
		}

		return centymo.HTMXSuccess("payment-table")
	})
}

// NewPaymentTableAction returns a payment table refresh trigger for HTMX.
func NewPaymentTableAction(deps *PaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return centymo.HTMXSuccess("payment-table")
	})
}
