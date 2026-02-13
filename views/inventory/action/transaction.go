package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
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
	FormAction  string
	Labels      TransactionFormLabels
	TypeOptions []SelectOption
	Today       string
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
				FormAction:  "/action/inventory/detail/" + inventoryItemID + "/transactions/assign",
				Labels:      transactionFormLabels(viewCtx.T),
				TypeOptions: transactionTypeOptions(viewCtx.T),
				Today:       time.Now().Format("2006-01-02"),
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

		data := map[string]any{
			"inventory_item_id": inventoryItemID,
			"transaction_type":  txType,
			"quantity":          qty,
			"transaction_date":  r.FormValue("transaction_date"),
			"reference":         r.FormValue("reference"),
			"serial_number":     r.FormValue("serial_number"),
			"notes":             r.FormValue("notes"),
			"performed_by":      "system", // TODO: current user
		}

		_, err := deps.DB.Create(ctx, "inventory_transaction", data)
		if err != nil {
			log.Printf("Failed to create transaction: %v", err)
			return centymo.HTMXError("Failed to record stock movement")
		}

		// Update inventory quantities based on transaction type
		item, err := deps.DB.Read(ctx, "inventory_item", inventoryItemID)
		if err == nil {
			currentOnHand := toFloat64FromAny(item["quantity_on_hand"])
			switch txType {
			case "received", "returned":
				currentOnHand += qty
			case "sold", "transferred", "write_off":
				currentOnHand -= qty
				if currentOnHand < 0 {
					currentOnHand = 0
				}
			}
			_, _ = deps.DB.Update(ctx, "inventory_item", inventoryItemID, map[string]any{
				"quantity_on_hand": currentOnHand,
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

func toFloat64FromAny(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}
