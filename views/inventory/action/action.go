package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	inventoryform "github.com/erniealice/centymo-golang/views/inventory/form"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
)

// Deps holds dependencies for the primary inventory CRUD action handlers.
// Feature-specific deps (serial, transaction, depreciation) live in their
// respective feature packages (inventory/serial, inventory/transaction,
// inventory/depreciation).
type Deps struct {
	Routes              centymo.InventoryRoutes
	Labels              centymo.InventoryLabels
	CreateInventoryItem func(ctx context.Context, req *inventoryitempb.CreateInventoryItemRequest) (*inventoryitempb.CreateInventoryItemResponse, error)
	ReadInventoryItem   func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	DeleteInventoryItem func(ctx context.Context, req *inventoryitempb.DeleteInventoryItemRequest) (*inventoryitempb.DeleteInventoryItemResponse, error)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func formLabels(t func(string) string, f centymo.InventoryFormLabels) inventoryform.Labels {
	return inventoryform.Labels{
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
		// Info fields sourced from centymo.InventoryFormLabels (populated from lyngua JSON + defaults).
		ProductInfo:       f.ProductInfo,
		SKUInfo:           f.SKUInfo,
		OnHandInfo:        f.OnHandInfo,
		ReservedInfo:      f.ReservedInfo,
		ReorderLevelInfo:  f.ReorderLevelInfo,
		UnitOfMeasureInfo: f.UnitOfMeasureInfo,
		NotesInfo:         f.NotesInfo,
		ActiveInfo:        f.ActiveInfo,
	}
}

// NewAddAction creates the inventory add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("inventory-drawer-form", &inventoryform.Data{
				FormAction:    deps.Routes.AddURL,
				Active:        true,
				UnitOfMeasure: "pcs",
				Labels:        formLabels(viewCtx.T, deps.Labels.Form),
				CommonLabels:  nil, // injected by ViewAdapter
			})
		}

		// POST - create inventory item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
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
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewEditAction creates the inventory edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read inventory item %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			items := resp.GetData()
			if len(items) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			item := items[0]

			return view.OK("inventory-drawer-form", &inventoryform.Data{
				FormAction:    route.ResolveURL(deps.Routes.EditURL, "id", id),
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
				Labels:        formLabels(viewCtx.T, deps.Labels.Form),
				CommonLabels:  nil, // injected by ViewAdapter
			})
		}

		// POST - update inventory item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
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
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewDeleteAction creates the inventory delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "delete") {
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

		_, err := deps.DeleteInventoryItem(ctx, &inventoryitempb.DeleteInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete inventory item %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

// NewBulkDeleteAction creates the inventory bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
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
