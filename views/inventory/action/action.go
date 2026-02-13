package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
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
	DB centymo.DataSource
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

		data := map[string]any{
			"name":               r.FormValue("product_name"),
			"sku":                r.FormValue("sku"),
			"quantity_on_hand":   onHand,
			"quantity_reserved":  reserved,
			"reorder_level":      reorderLevel,
			"unit_of_measure":    r.FormValue("unit_of_measure"),
			"location_id":        r.FormValue("location_id"),
			"notes":              r.FormValue("notes"),
			"active":             active,
		}

		_, err := deps.DB.Create(ctx, "inventory_item", data)
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
			record, err := deps.DB.Read(ctx, "inventory_item", id)
			if err != nil {
				log.Printf("Failed to read inventory item %s: %v", id, err)
				return centymo.HTMXError("Inventory item not found")
			}

			name, _ := record["name"].(string)
			sku, _ := record["sku"].(string)
			locationID, _ := record["location_id"].(string)
			notes, _ := record["notes"].(string)
			unitOfMeasure, _ := record["unit_of_measure"].(string)
			active, _ := record["active"].(bool)

			return view.OK("inventory-drawer-form", &FormData{
				FormAction:    "/action/inventory/edit/" + id,
				IsEdit:        true,
				ID:            id,
				Name:          name,
				SKU:           sku,
				OnHand:        anyToString(record["quantity_on_hand"]),
				Reserved:      anyToString(record["quantity_reserved"]),
				ReorderLevel:  anyToString(record["reorder_level"]),
				UnitOfMeasure: unitOfMeasure,
				LocationID:    locationID,
				Notes:         notes,
				Active:        active,
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

		data := map[string]any{
			"name":               r.FormValue("product_name"),
			"sku":                r.FormValue("sku"),
			"quantity_on_hand":   onHand,
			"quantity_reserved":  reserved,
			"reorder_level":      reorderLevel,
			"unit_of_measure":    r.FormValue("unit_of_measure"),
			"location_id":        r.FormValue("location_id"),
			"notes":              r.FormValue("notes"),
			"active":             active,
		}

		_, err := deps.DB.Update(ctx, "inventory_item", id, data)
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

		err := deps.DB.Delete(ctx, "inventory_item", id)
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
			err := deps.DB.Delete(ctx, "inventory_item", id)
			if err != nil {
				log.Printf("Failed to delete inventory item %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("inventory-table")
	})
}

func anyToString(v any) string {
	if v == nil {
		return "0"
	}
	return fmt.Sprintf("%v", v)
}
