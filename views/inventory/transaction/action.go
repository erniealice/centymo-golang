// Package transaction handles the stock movement feature for inventory items.
// Drawer template: transaction-drawer-form.html (stays flat at view root).
package transaction

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	transactionform "github.com/erniealice/centymo-golang/views/inventory/transaction/form"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
)

// Deps is the dependency subset needed by the transaction feature.
type Deps struct {
	Routes centymo.InventoryRoutes
	Labels centymo.InventoryLabels

	CreateInventoryTransaction func(ctx context.Context, req *inventorytransactionpb.CreateInventoryTransactionRequest) (*inventorytransactionpb.CreateInventoryTransactionResponse, error)
	ReadInventoryItem          func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem        func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
}

func formLabels(t func(string) string, tx centymo.InventoryTransactionLabels) transactionform.Labels {
	return transactionform.Labels{
		Type:      t("inventory.transaction.type"),
		Quantity:  t("inventory.transaction.quantity"),
		Date:      t("inventory.transaction.date"),
		Reference: t("inventory.transaction.reference"),
		// Info fields sourced from centymo.InventoryTransactionLabels (populated from lyngua JSON + defaults).
		TypeInfo:      tx.TypeInfo,
		QuantityInfo:  tx.QuantityInfo,
		DateInfo:      tx.DateInfo,
		ReferenceInfo: tx.ReferenceInfo,
	}
}

func typeOptions(t func(string) string) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "received", Label: t("inventory.transaction.typeReceived")},
		{Value: "sold", Label: t("inventory.transaction.typeSold")},
		{Value: "adjusted", Label: t("inventory.transaction.typeAdjusted")},
		{Value: "transferred", Label: t("inventory.transaction.typeTransferred")},
		{Value: "returned", Label: t("inventory.transaction.typeReturned")},
		{Value: "write_off", Label: t("inventory.transaction.typeWriteOff")},
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NewAssignAction creates the stock movement action (GET = form, POST = create).
func NewAssignAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("transaction-drawer-form", &transactionform.Data{
				FormAction:   route.ResolveURL(deps.Routes.TransactionAssignURL, "id", inventoryItemID),
				Labels:       formLabels(viewCtx.T, deps.Labels.Transaction),
				TypeOptions:  typeOptions(viewCtx.T),
				Today:        time.Now().Format("2006-01-02"),
				CommonLabels: nil,
			})
		}

		// POST - create transaction + update quantities
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		txType := r.FormValue("transaction_type")

		data := &inventorytransactionpb.InventoryTransaction{
			InventoryItemId:       inventoryItemID,
			TransactionType:       txType,
			Quantity:              qty,
			TransactionDateString: strPtr(r.FormValue("transaction_date")),
			ReferenceType:         strPtr(r.FormValue("reference")),
			SerialNumber:          strPtr(r.FormValue("serial_number")),
			Notes:                 strPtr(r.FormValue("notes")),
			PerformedBy:           strPtr("system"), // TODO: current user
		}

		_, err := deps.CreateInventoryTransaction(ctx, &inventorytransactionpb.CreateInventoryTransactionRequest{Data: data})
		if err != nil {
			log.Printf("Failed to create transaction: %v", err)
			return centymo.HTMXError(err.Error())
		}

		// Update inventory quantities based on transaction type
		itemResp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: inventoryItemID},
		})
		if err == nil && len(itemResp.GetData()) > 0 {
			item := itemResp.GetData()[0]
			currentOnHand := item.GetQuantityOnHand()
			switch txType {
			case "received", "returned":
				currentOnHand += qty
			case "sold", "transferred", "write_off":
				currentOnHand -= qty
				if currentOnHand < 0 {
					currentOnHand = 0
				}
			}
			_, _ = deps.UpdateInventoryItem(ctx, &inventoryitempb.UpdateInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{
					Id:             inventoryItemID,
					QuantityOnHand: currentOnHand,
				},
			})
		}

		return centymo.HTMXSuccess("transaction-table")
	})
}

// NewTableAction returns the transaction table partial for HTMX refresh.
func NewTableAction(_ *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return centymo.HTMXSuccess("transaction-table")
	})
}
