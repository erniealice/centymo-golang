// Package payment owns the handler and dep-bearing helpers for the revenue
// payment drawer (revenue-payment-drawer-form.html).
package payment

import (
	"context"
	"fmt"
	"log"
	"net/http"

	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	"github.com/erniealice/centymo-golang/domain/revenue/revenue/payment/form"
	revenuepaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_payment"
	collectionmethodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds dependencies for payment action handlers.
//
// 20260612-datasource-typed-path W5: the DataSource duck was replaced with
// typed proto closures. Every closure may be nil during W5 (bound by
// service-admin's adapter in W7) — each call site nil-checks and degrades to
// an empty/zero result, never panics.
type Deps struct {
	Routes revenuedomain.Routes
	Labels revenuedomain.Labels

	// Typed revenue_payment CRUD (replaces DataSource on "revenue_payment").
	CreateRevenuePayment func(ctx context.Context, req *revenuepaymentpb.CreateRevenuePaymentRequest) (*revenuepaymentpb.CreateRevenuePaymentResponse, error)
	ReadRevenuePayment   func(ctx context.Context, req *revenuepaymentpb.ReadRevenuePaymentRequest) (*revenuepaymentpb.ReadRevenuePaymentResponse, error)
	UpdateRevenuePayment func(ctx context.Context, req *revenuepaymentpb.UpdateRevenuePaymentRequest) (*revenuepaymentpb.UpdateRevenuePaymentResponse, error)
	DeleteRevenuePayment func(ctx context.Context, req *revenuepaymentpb.DeleteRevenuePaymentRequest) (*revenuepaymentpb.DeleteRevenuePaymentResponse, error)

	// Typed collection_method reads (replaces DataSource on "collection_method").
	ReadCollectionMethod  func(ctx context.Context, req *collectionmethodpb.ReadCollectionMethodRequest) (*collectionmethodpb.ReadCollectionMethodResponse, error)
	ListCollectionMethods func(ctx context.Context, req *collectionmethodpb.ListCollectionMethodsRequest) (*collectionmethodpb.ListCollectionMethodsResponse, error)
}

// strPtr returns a pointer to s for optional proto string fields. Returns nil
// for the empty string so optional fields stay unset rather than "" (matches
// the proto oneof semantics).
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// clearablePtr always returns a pointer (even for ""), so an emptied optional
// text field is sent as "" on update — clearing it — rather than being omitted
// (proto nil → key absent → the generic Update leaves the old DB value). This
// preserves the pre-typed map-write behavior where every form field was written.
func clearablePtr(s string) *string {
	return &s
}

// formatCentavosDecimal renders int64 centavos back to a plain decimal-pesos
// string (e.g. 15000 → "150.00") for the edit-form's amount input.
func formatCentavosDecimal(centavos int64) string {
	return fmt.Sprintf("%.2f", float64(centavos)/100)
}

// loadCollectionMethods loads collection methods via the typed use case and
// returns them as select options. Nil-safe: returns an empty slice when the
// closure is unwired (W5 — bound by service-admin in W7).
func loadCollectionMethods(ctx context.Context, list func(ctx context.Context, req *collectionmethodpb.ListCollectionMethodsRequest) (*collectionmethodpb.ListCollectionMethodsResponse, error)) []pyeza.SelectOption {
	if list == nil {
		return []pyeza.SelectOption{}
	}
	resp, err := list(ctx, &collectionmethodpb.ListCollectionMethodsRequest{})
	if err != nil {
		log.Printf("Failed to list collection methods: %v", err)
		return []pyeza.SelectOption{}
	}

	data := resp.GetData()
	options := make([]pyeza.SelectOption, 0, len(data))
	for _, m := range data {
		id := m.GetId()
		if id == "" {
			continue
		}
		name := m.GetName()
		if name == "" {
			name = id
		}
		options = append(options, pyeza.SelectOption{Value: id, Label: name})
	}
	return options
}

// resolveMethodName looks up the collection method's display name for the
// payment_method column. Nil-safe: falls back to the raw method id when the
// closure is unwired or the lookup fails.
func resolveMethodName(ctx context.Context, read func(ctx context.Context, req *collectionmethodpb.ReadCollectionMethodRequest) (*collectionmethodpb.ReadCollectionMethodResponse, error), methodID string) string {
	if methodID == "" || read == nil {
		return methodID
	}
	resp, err := read(ctx, &collectionmethodpb.ReadCollectionMethodRequest{
		Data: &collectionmethodpb.CollectionMethod{Id: methodID},
	})
	if err != nil {
		return methodID
	}
	for _, m := range resp.GetData() {
		if m.GetId() == methodID {
			if name := m.GetName(); name != "" {
				return name
			}
		}
	}
	return methodID
}

// NewAddAction creates the payment add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			methods := loadCollectionMethods(ctx, deps.ListCollectionMethods)
			return view.OK("revenue-payment-drawer-form", &form.Data{
				FormAction:     route.ResolveURL(deps.Routes.PaymentAddURL, "id", revenueID),
				RevenueID:      revenueID,
				Currency:       "PHP",
				PaymentMethods: methods,
				CommonLabels:   nil, // injected by ViewAdapter
				Labels:         deps.Labels,
			})
		}

		// POST — create payment
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		if deps.CreateRevenuePayment == nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		collectionMethodID := r.FormValue("collection_method_id")
		methodName := resolveMethodName(ctx, deps.ReadCollectionMethod, collectionMethodID)

		// Decimal-pesos form input → int64 centavos (Rule #1). Reject invalid
		// input rather than silently persisting 0 (the data-loss class this wave fixes).
		amountCentavos, perr := pyeza.ParseCentavos(r.FormValue("amount"))
		if perr != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		_, err := deps.CreateRevenuePayment(ctx, &revenuepaymentpb.CreateRevenuePaymentRequest{
			Data: &revenuepaymentpb.RevenuePayment{
				RevenueId:          revenueID,
				PaymentMethod:      strPtr(methodName),
				Amount:             amountCentavos,
				Currency:           r.FormValue("currency"),
				CollectionMethodId: strPtr(collectionMethodID),
				ReferenceNumber:    strPtr(r.FormValue("reference_number")),
				ReceivedBy:         strPtr(r.FormValue("received_by")),
				ReceivedRole:       strPtr(r.FormValue("received_role")),
				CollectionType:     strPtr("sale"),
				Status:             strPtr("completed"),
				Notes:              strPtr(r.FormValue("notes")),
			},
		})
		if err != nil {
			log.Printf("Failed to create payment for revenue %s: %v", revenueID, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("payment-table")
	})
}

// NewEditAction creates the payment edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		revenueID := viewCtx.Request.PathValue("id")
		paymentID := viewCtx.Request.PathValue("pid")

		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadRevenuePayment == nil {
				return view.HTMXError(deps.Labels.Errors.PaymentNotFound)
			}
			resp, err := deps.ReadRevenuePayment(ctx, &revenuepaymentpb.ReadRevenuePaymentRequest{
				Data: &revenuepaymentpb.RevenuePayment{Id: paymentID},
			})
			if err != nil {
				log.Printf("Failed to read payment %s: %v", paymentID, err)
				return view.HTMXError(deps.Labels.Errors.PaymentNotFound)
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.HTMXError(deps.Labels.Errors.PaymentNotFound)
			}
			record := data[0]

			methods := loadCollectionMethods(ctx, deps.ListCollectionMethods)
			return view.OK("revenue-payment-drawer-form", &form.Data{
				FormAction:         route.ResolveURL(deps.Routes.PaymentEditURL, "id", revenueID, "pid", paymentID),
				IsEdit:             true,
				ID:                 paymentID,
				RevenueID:          revenueID,
				CollectionMethodID: record.GetCollectionMethodId(),
				AmountPaid:         formatCentavosDecimal(record.GetAmount()),
				Currency:           record.GetCurrency(),
				ReferenceNumber:    record.GetReferenceNumber(),
				Notes:              record.GetNotes(),
				ReceivedBy:         record.GetReceivedBy(),
				ReceivedRole:       record.GetReceivedRole(),
				PaymentMethods:     methods,
				CommonLabels:       nil, // injected by ViewAdapter
				Labels:             deps.Labels,
			})
		}

		// POST — update payment
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		if deps.UpdateRevenuePayment == nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		collectionMethodID := r.FormValue("collection_method_id")
		methodName := resolveMethodName(ctx, deps.ReadCollectionMethod, collectionMethodID)

		amountCentavos, perr := pyeza.ParseCentavos(r.FormValue("amount"))
		if perr != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		_, err := deps.UpdateRevenuePayment(ctx, &revenuepaymentpb.UpdateRevenuePaymentRequest{
			Data: &revenuepaymentpb.RevenuePayment{
				Id:                 paymentID,
				PaymentMethod:      strPtr(methodName),
				Amount:             amountCentavos,
				Currency:           r.FormValue("currency"),
				CollectionMethodId: strPtr(collectionMethodID),
				// clearablePtr (not strPtr): emptying these on edit must persist as
				// "" (clear), matching the pre-typed map-write behavior.
				ReferenceNumber: clearablePtr(r.FormValue("reference_number")),
				ReceivedRole:    clearablePtr(r.FormValue("received_role")),
				Notes:           clearablePtr(r.FormValue("notes")),
				// Never change received_by on edit (omitted from update payload).
			},
		})
		if err != nil {
			log.Printf("Failed to update payment %s: %v", paymentID, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("payment-table")
	})
}

// NewRemoveAction creates the payment remove action (POST only).
func NewRemoveAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.IDRequired)
		}

		if deps.DeleteRevenuePayment == nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		_, err := deps.DeleteRevenuePayment(ctx, &revenuepaymentpb.DeleteRevenuePaymentRequest{
			Data: &revenuepaymentpb.RevenuePayment{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete payment %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("payment-table")
	})
}

// NewTableAction returns a payment table refresh trigger for HTMX.
func NewTableAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return view.HTMXSuccess("payment-table")
	})
}
