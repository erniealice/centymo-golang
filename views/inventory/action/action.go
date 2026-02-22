package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
)

// FormLabels holds i18n labels for the inventory drawer form template.
type FormLabels struct {
	Product          string
	SKU              string
	SKUPlaceholder   string
	OnHand           string
	Reserved         string
	ReorderLevel     string
	UnitOfMeasure    string
	Notes            string
	NotesPlaceholder string
	Active           string
}

// FormData is the template data for the inventory drawer form.
type FormData struct {
	FormAction    string
	IsEdit        bool
	ID            string
	Name          string
	SKU           string
	OnHand        string
	Reserved      string
	ReorderLevel  string
	UnitOfMeasure string
	LocationID    string
	Notes         string
	Active        bool
	Labels        FormLabels
	CommonLabels  any
}

// Deps holds dependencies for inventory action handlers.
type Deps struct {
	CreateInventoryItem        func(ctx context.Context, req *inventoryitempb.CreateInventoryItemRequest) (*inventoryitempb.CreateInventoryItemResponse, error)
	ReadInventoryItem          func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem        func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	DeleteInventoryItem        func(ctx context.Context, req *inventoryitempb.DeleteInventoryItemRequest) (*inventoryitempb.DeleteInventoryItemResponse, error)
	CreateInventorySerial      func(ctx context.Context, req *inventoryserialpb.CreateInventorySerialRequest) (*inventoryserialpb.CreateInventorySerialResponse, error)
	ReadInventorySerial        func(ctx context.Context, req *inventoryserialpb.ReadInventorySerialRequest) (*inventoryserialpb.ReadInventorySerialResponse, error)
	UpdateInventorySerial      func(ctx context.Context, req *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	DeleteInventorySerial      func(ctx context.Context, req *inventoryserialpb.DeleteInventorySerialRequest) (*inventoryserialpb.DeleteInventorySerialResponse, error)
	CreateInventoryTransaction func(ctx context.Context, req *inventorytransactionpb.CreateInventoryTransactionRequest) (*inventorytransactionpb.CreateInventoryTransactionResponse, error)
	CreateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.CreateInventoryDepreciationRequest) (*inventorydepreciationpb.CreateInventoryDepreciationResponse, error)
	ReadInventoryDepreciation  func(ctx context.Context, req *inventorydepreciationpb.ReadInventoryDepreciationRequest) (*inventorydepreciationpb.ReadInventoryDepreciationResponse, error)
	UpdateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.UpdateInventoryDepreciationRequest) (*inventorydepreciationpb.UpdateInventoryDepreciationResponse, error)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Product:          t("inventory.form.product"),
		SKU:              t("inventory.form.sku"),
		SKUPlaceholder:   t("inventory.form.skuPlaceholder"),
		OnHand:           t("inventory.form.onHand"),
		Reserved:         t("inventory.form.reserved"),
		ReorderLevel:     t("inventory.form.reorderLevel"),
		UnitOfMeasure:    t("inventory.form.unitOfMeasure"),
		Notes:            t("inventory.form.notes"),
		NotesPlaceholder: t("inventory.form.notesPlaceholder"),
		Active:           t("inventory.form.active"),
	}
}

// NewAddAction creates the inventory add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("inventory-drawer-form", &FormData{
				FormAction:    "/action/inventory/add",
				Active:        true,
				UnitOfMeasure: "pcs",
				Labels:        formLabels(viewCtx.T),
				CommonLabels:  nil, // injected by ViewAdapter
			})
		}

		// POST - create inventory item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		onHand, _ := strconv.ParseFloat(r.FormValue("quantity_on_hand"), 64)
		reserved, _ := strconv.ParseFloat(r.FormValue("quantity_reserved"), 64)
		reorderLevel, _ := strconv.ParseFloat(r.FormValue("reorder_level"), 64)

		data := &inventoryitempb.InventoryItem{
			Name:             r.FormValue("product_name"),
			Sku:              strPtr(r.FormValue("sku")),
			QuantityOnHand:   onHand,
			QuantityReserved: reserved,
			ReorderLevel:     &reorderLevel,
			UnitOfMeasure:    r.FormValue("unit_of_measure"),
			LocationId:       strPtr(r.FormValue("location_id")),
			Notes:            strPtr(r.FormValue("notes")),
			Active:           active,
		}

		_, err := deps.CreateInventoryItem(ctx, &inventoryitempb.CreateInventoryItemRequest{Data: data})
		if err != nil {
			log.Printf("Failed to create inventory item: %v", err)
			return centymo.HTMXError("Failed to create inventory item")
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewEditAction creates the inventory edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read inventory item %s: %v", id, err)
				return centymo.HTMXError("Inventory item not found")
			}
			items := resp.GetData()
			if len(items) == 0 {
				return centymo.HTMXError("Inventory item not found")
			}
			item := items[0]

			return view.OK("inventory-drawer-form", &FormData{
				FormAction:    "/action/inventory/edit/" + id,
				IsEdit:        true,
				ID:            id,
				Name:          item.GetName(),
				SKU:           item.GetSku(),
				OnHand:        formatFloat(item.GetQuantityOnHand()),
				Reserved:      formatFloat(item.GetQuantityReserved()),
				ReorderLevel:  formatFloat(item.GetReorderLevel()),
				UnitOfMeasure: item.GetUnitOfMeasure(),
				LocationID:    item.GetLocationId(),
				Notes:         item.GetNotes(),
				Active:        item.GetActive(),
				Labels:        formLabels(viewCtx.T),
				CommonLabels:  nil, // injected by ViewAdapter
			})
		}

		// POST - update inventory item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		onHand, _ := strconv.ParseFloat(r.FormValue("quantity_on_hand"), 64)
		reserved, _ := strconv.ParseFloat(r.FormValue("quantity_reserved"), 64)
		reorderLevel, _ := strconv.ParseFloat(r.FormValue("reorder_level"), 64)

		data := &inventoryitempb.InventoryItem{
			Id:               id,
			Name:             r.FormValue("product_name"),
			Sku:              strPtr(r.FormValue("sku")),
			QuantityOnHand:   onHand,
			QuantityReserved: reserved,
			ReorderLevel:     &reorderLevel,
			UnitOfMeasure:    r.FormValue("unit_of_measure"),
			LocationId:       strPtr(r.FormValue("location_id")),
			Notes:            strPtr(r.FormValue("notes")),
			Active:           active,
		}

		_, err := deps.UpdateInventoryItem(ctx, &inventoryitempb.UpdateInventoryItemRequest{Data: data})
		if err != nil {
			log.Printf("Failed to update inventory item %s: %v", id, err)
			return centymo.HTMXError("Failed to update inventory item")
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewDeleteAction creates the inventory delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Inventory item ID is required")
		}

		_, err := deps.DeleteInventoryItem(ctx, &inventoryitempb.DeleteInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete inventory item %s: %v", id, err)
			return centymo.HTMXError("Failed to delete inventory item")
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewBulkDeleteAction creates the inventory bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError("No inventory item IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteInventoryItem(ctx, &inventoryitempb.DeleteInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete inventory item %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

func formatFloat(f float64) string {
	if f == 0 {
		return "0"
	}
	return fmt.Sprintf("%g", f)
}
