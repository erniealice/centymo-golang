package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	poform "github.com/erniealice/centymo-golang/views/purchase_order/form"

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
)

// Deps holds dependencies for purchase order action handlers.
type Deps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	// Typed purchase order CRUD operations
	CreatePurchaseOrder func(ctx context.Context, req *purchaseorderpb.CreatePurchaseOrderRequest) (*purchaseorderpb.CreatePurchaseOrderResponse, error)
	ReadPurchaseOrder   func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	UpdatePurchaseOrder func(ctx context.Context, req *purchaseorderpb.UpdatePurchaseOrderRequest) (*purchaseorderpb.UpdatePurchaseOrderResponse, error)
	DeletePurchaseOrder func(ctx context.Context, req *purchaseorderpb.DeletePurchaseOrderRequest) (*purchaseorderpb.DeletePurchaseOrderResponse, error)
}

// formLabels maps ExpenditureLabels into the flat Labels struct for the template.
// Kept in action/ (not deleted) because it performs real transformation:
// hardcoded "PO Number", "Supplier", "PO Type" strings + draws from two label sources.
func formLabels(l centymo.ExpenditureLabels) poform.Labels {
	return poform.Labels{
		PoNumber:         "PO Number",
		SupplierID:       "Supplier",
		PoType:           "PO Type",
		OrderDate:        l.Form.ExpenditureDate,
		Currency:         l.Form.Currency,
		PaymentTerms:     l.Form.PaymentTerms,
		Notes:            l.Form.Notes,
		NotesPlaceholder: l.Form.NotesPlaceholder,
		Status:           l.Form.Status,
		// Info fields sourced from centymo.PurchaseOrderFormLabels (populated from lyngua JSON + defaults).
		PoNumberInfo:     l.PurchaseOrder.Form.PONumberInfo,
		SupplierIDInfo:   l.PurchaseOrder.Form.SupplierInfo,
		PoTypeInfo:       l.PurchaseOrder.Form.POTypeInfo,
		OrderDateInfo:    l.PurchaseOrder.Form.OrderDateInfo,
		CurrencyInfo:     l.PurchaseOrder.Form.CurrencyInfo,
		PaymentTermsInfo: l.PurchaseOrder.Form.PaymentTermsInfo,
		NotesInfo:        l.PurchaseOrder.Form.NotesInfo,
	}
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// NewAddAction creates the purchase order add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("purchase-order-drawer-form", &poform.Data{
				FormAction:   deps.Routes.PurchaseOrderAddURL,
				PoType:       "standard",
				Currency:     "PHP",
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create purchase order
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		resp, err := deps.CreatePurchaseOrder(ctx, &purchaseorderpb.CreatePurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{
				PoNumber:        r.FormValue("po_number"),
				SupplierId:      r.FormValue("supplier_id"),
				PoType:          r.FormValue("po_type"),
				OrderDateString: strPtr(r.FormValue("order_date_string")),
				Currency:        r.FormValue("currency"),
				PaymentTerms:    strPtr(r.FormValue("payment_terms")),
				Notes:           strPtr(r.FormValue("notes")),
			},
		})
		if err != nil {
			log.Printf("Failed to create purchase order: %v", err)
			return centymo.HTMXError(err.Error())
		}

		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.Routes.PurchaseOrderDetailURL, "id", newID),
				},
			}
		}

		return centymo.HTMXSuccess("purchase-orders-table")
	})
}

// NewEditAction creates the purchase order edit action (GET = pre-filled form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadPurchaseOrder(ctx, &purchaseorderpb.ReadPurchaseOrderRequest{
				Data: &purchaseorderpb.PurchaseOrder{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read purchase order %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			return view.OK("purchase-order-drawer-form", &poform.Data{
				FormAction:   route.ResolveURL(deps.Routes.PurchaseOrderEditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				PoNumber:     record.GetPoNumber(),
				SupplierID:   record.GetSupplierId(),
				PoType:       record.GetPoType(),
				OrderDate:    record.GetOrderDateString(),
				Currency:     record.GetCurrency(),
				PaymentTerms: record.GetPaymentTerms(),
				Notes:        record.GetNotes(),
				Labels:       formLabels(deps.Labels),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update purchase order
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdatePurchaseOrder(ctx, &purchaseorderpb.UpdatePurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{
				Id:              id,
				PoNumber:        r.FormValue("po_number"),
				SupplierId:      r.FormValue("supplier_id"),
				PoType:          r.FormValue("po_type"),
				OrderDateString: strPtr(r.FormValue("order_date_string")),
				Currency:        r.FormValue("currency"),
				PaymentTerms:    strPtr(r.FormValue("payment_terms")),
				Notes:           strPtr(r.FormValue("notes")),
			},
		})
		if err != nil {
			log.Printf("Failed to update purchase order %s: %v", id, err)
			return centymo.HTMXError(err.Error())
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

// NewDeleteAction creates the purchase order delete action (POST only).
// The record ID comes via query param (?id=xxx) or form field.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeletePurchaseOrder(ctx, &purchaseorderpb.DeletePurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete purchase order %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("purchase-orders-table")
	})
}

// NewSetStatusAction creates the purchase order status update action (POST only).
// Expects query params: ?id={poId}&status={approved|partially_received|fully_received|closed|cancelled}
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("purchase_order", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		validStatuses := map[string]bool{
			"approved":           true,
			"partially_received": true,
			"fully_received":     true,
			"closed":             true,
			"cancelled":          true,
		}
		if !validStatuses[targetStatus] {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		_, err := deps.UpdatePurchaseOrder(ctx, &purchaseorderpb.UpdatePurchaseOrderRequest{
			Data: &purchaseorderpb.PurchaseOrder{Id: id, Status: targetStatus},
		})
		if err != nil {
			log.Printf("Failed to update purchase order status %s -> %s: %v", id, targetStatus, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("purchase-orders-table")
	})
}
