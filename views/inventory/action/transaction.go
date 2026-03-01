package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
)

// TransactionFormLabels holds i18n labels for the transaction drawer form.
type TransactionFormLabels struct {
	Type      string
	Quantity  string
	Date      string
	Reference string
}

// TransactionFormData is the template data for the transaction drawer form.
type TransactionFormData struct {
	FormAction   string
	Labels       TransactionFormLabels
	TypeOptions  []SelectOption
	Today        string
	CommonLabels any
}

func transactionFormLabels(t func(string) string) TransactionFormLabels {
	return TransactionFormLabels{
		Type:      t("inventory.transaction.type"),
		Quantity:  t("inventory.transaction.quantity"),
		Date:      t("inventory.transaction.date"),
		Reference: t("inventory.transaction.reference"),
	}
}

func transactionTypeOptions(t func(string) string) []SelectOption {
	return []SelectOption{
		{Value: "received", Label: t("inventory.transaction.typeReceived")},
		{Value: "sold", Label: t("inventory.transaction.typeSold")},
		{Value: "adjusted", Label: t("inventory.transaction.typeAdjusted")},
		{Value: "transferred", Label: t("inventory.transaction.typeTransferred")},
		{Value: "returned", Label: t("inventory.transaction.typeReturned")},
		{Value: "write_off", Label: t("inventory.transaction.typeWriteOff")},
	}
}

// NewTransactionAssignAction creates the stock movement action (GET = form, POST = create).
func NewTransactionAssignAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("transaction-drawer-form", &TransactionFormData{
				FormAction:   route.ResolveURL(deps.Routes.TransactionAssignURL, "id", inventoryItemID),
				Labels:       transactionFormLabels(viewCtx.T),
				TypeOptions:  transactionTypeOptions(viewCtx.T),
				Today:        time.Now().Format("2006-01-02"),
				CommonLabels: nil,
			})
		}

		// POST - create transaction + update quantities
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
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

// NewTransactionTableAction returns the transaction table partial for HTMX refresh.
func NewTransactionTableAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return centymo.HTMXSuccess("transaction-table")
	})
}
