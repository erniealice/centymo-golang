package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// SerialFormLabels holds i18n labels for the serial drawer form.
type SerialFormLabels struct {
	SerialNumber  string
	IMEI          string
	Status        string
	WarrantyStart string
	WarrantyEnd   string
	PurchaseOrder string
	SoldReference string
}

// SelectOption represents a select dropdown option.
type SelectOption struct {
	Value string
	Label string
}

// SerialFormData is the template data for the serial drawer form.
type SerialFormData struct {
	FormAction    string
	IsEdit        bool
	ID            string
	SerialNumber  string
	IMEI          string
	Status        string
	WarrantyStart string
	WarrantyEnd   string
	PurchaseOrder string
	SoldReference string
	Labels        SerialFormLabels
	StatusOptions []SelectOption
	CommonLabels  any
}

func serialFormLabels(t func(string) string) SerialFormLabels {
	return SerialFormLabels{
		SerialNumber:  t("inventory.serial.serialNumber"),
		IMEI:          t("inventory.serial.imei"),
		Status:        t("inventory.serial.status"),
		WarrantyStart: t("inventory.serial.warrantyStart"),
		WarrantyEnd:   t("inventory.serial.warrantyEnd"),
		PurchaseOrder: t("inventory.serial.purchaseOrder"),
		SoldReference: t("inventory.serial.soldReference"),
	}
}

func serialStatusOptions(t func(string) string) []SelectOption {
	return []SelectOption{
		{Value: "available", Label: t("inventory.serial.statusAvailable")},
		{Value: "sold", Label: t("inventory.serial.statusSold")},
		{Value: "reserved", Label: t("inventory.serial.statusReserved")},
		{Value: "defective", Label: t("inventory.serial.statusDefective")},
		{Value: "returned", Label: t("inventory.serial.statusReturned")},
	}
}

// NewSerialAssignAction creates the serial assign action (GET = form, POST = create).
func NewSerialAssignAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("serial-drawer-form", &SerialFormData{
				FormAction:    "/action/inventory/detail/" + inventoryItemID + "/serials/assign",
				Status:        "available",
				Labels:        serialFormLabels(viewCtx.T),
				StatusOptions: serialStatusOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - create serial
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		data := map[string]any{
			"inventory_item_id": inventoryItemID,
			"serial_number":     r.FormValue("serial_number"),
			"imei":              r.FormValue("imei"),
			"status":            r.FormValue("status"),
			"warranty_start":    r.FormValue("warranty_start"),
			"warranty_end":      r.FormValue("warranty_end"),
			"purchase_order":    r.FormValue("purchase_order"),
			"sold_reference":    r.FormValue("sold_reference"),
		}

		_, err := deps.DB.Create(ctx, "inventory_serial", data)
		if err != nil {
			log.Printf("Failed to create serial: %v", err)
			return centymo.HTMXError("Failed to create serial")
		}

		return centymo.HTMXSuccess("serial-table")
	})
}

// NewSerialEditAction creates the serial edit action (GET = form, POST = update).
func NewSerialEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		inventoryItemID := viewCtx.Request.PathValue("id")
		serialID := viewCtx.Request.PathValue("sid")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "inventory_serial", serialID)
			if err != nil {
				log.Printf("Failed to read serial %s: %v", serialID, err)
				return centymo.HTMXError("Serial not found")
			}

			serialNumber, _ := record["serial_number"].(string)
			imei, _ := record["imei"].(string)
			status, _ := record["status"].(string)
			warrantyStart, _ := record["warranty_start"].(string)
			warrantyEnd, _ := record["warranty_end"].(string)
			po, _ := record["purchase_order"].(string)
			soldRef, _ := record["sold_reference"].(string)

			return view.OK("serial-drawer-form", &SerialFormData{
				FormAction:    "/action/inventory/detail/" + inventoryItemID + "/serials/edit/" + serialID,
				IsEdit:        true,
				ID:            serialID,
				SerialNumber:  serialNumber,
				IMEI:          imei,
				Status:        status,
				WarrantyStart: warrantyStart,
				WarrantyEnd:   warrantyEnd,
				PurchaseOrder: po,
				SoldReference: soldRef,
				Labels:        serialFormLabels(viewCtx.T),
				StatusOptions: serialStatusOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - update serial
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		data := map[string]any{
			"serial_number":  r.FormValue("serial_number"),
			"imei":           r.FormValue("imei"),
			"status":         r.FormValue("status"),
			"warranty_start": r.FormValue("warranty_start"),
			"warranty_end":   r.FormValue("warranty_end"),
			"purchase_order": r.FormValue("purchase_order"),
			"sold_reference": r.FormValue("sold_reference"),
		}

		_, err := deps.DB.Update(ctx, "inventory_serial", serialID, data)
		if err != nil {
			log.Printf("Failed to update serial %s: %v", serialID, err)
			return centymo.HTMXError("Failed to update serial")
		}

		return centymo.HTMXSuccess("serial-table")
	})
}

// NewSerialRemoveAction creates the serial remove action (POST only).
func NewSerialRemoveAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Serial ID is required")
		}

		err := deps.DB.Delete(ctx, "inventory_serial", id)
		if err != nil {
			log.Printf("Failed to delete serial %s: %v", id, err)
			return centymo.HTMXError("Failed to delete serial")
		}

		return centymo.HTMXSuccess("serial-table")
	})
}

// NewSerialTableAction returns the serial table partial for HTMX refresh.
func NewSerialTableAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// Trigger table refresh on the client side
		return centymo.HTMXSuccess("serial-table")
	})
}
