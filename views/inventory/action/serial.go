package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
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
	StatusOptions []pyeza.SelectOption
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

func serialStatusOptions(t func(string) string) []pyeza.SelectOption {
	return []pyeza.SelectOption{
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
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("serial-drawer-form", &SerialFormData{
				FormAction:    route.ResolveURL(deps.Routes.SerialAssignURL, "id", inventoryItemID),
				Status:        "available",
				Labels:        serialFormLabels(viewCtx.T),
				StatusOptions: serialStatusOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - create serial
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		data := &inventoryserialpb.InventorySerial{
			InventoryItemId: inventoryItemID,
			SerialNumber:    r.FormValue("serial_number"),
			Imei:            strPtr(r.FormValue("imei")),
			Status:          r.FormValue("status"),
			WarrantyStart:   strPtr(r.FormValue("warranty_start")),
			WarrantyEnd:     strPtr(r.FormValue("warranty_end")),
			PurchaseOrder:   strPtr(r.FormValue("purchase_order")),
		}

		_, err := deps.CreateInventorySerial(ctx, &inventoryserialpb.CreateInventorySerialRequest{Data: data})
		if err != nil {
			log.Printf("Failed to create serial: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("serial-table")
	})
}

// NewSerialEditAction creates the serial edit action (GET = form, POST = update).
func NewSerialEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		inventoryItemID := viewCtx.Request.PathValue("id")
		serialID := viewCtx.Request.PathValue("sid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadInventorySerial(ctx, &inventoryserialpb.ReadInventorySerialRequest{
				Data: &inventoryserialpb.InventorySerial{Id: serialID},
			})
			if err != nil {
				log.Printf("Failed to read serial %s: %v", serialID, err)
				return centymo.HTMXError(deps.Labels.Errors.SerialNotFound)
			}
			records := resp.GetData()
			if len(records) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.SerialNotFound)
			}
			record := records[0]

			return view.OK("serial-drawer-form", &SerialFormData{
				FormAction:    route.ResolveURL(deps.Routes.SerialEditURL, "id", inventoryItemID, "sid", serialID),
				IsEdit:        true,
				ID:            serialID,
				SerialNumber:  record.GetSerialNumber(),
				IMEI:          record.GetImei(),
				Status:        record.GetStatus(),
				WarrantyStart: record.GetWarrantyStart(),
				WarrantyEnd:   record.GetWarrantyEnd(),
				PurchaseOrder: record.GetPurchaseOrder(),
				Labels:        serialFormLabels(viewCtx.T),
				StatusOptions: serialStatusOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - update serial
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		data := &inventoryserialpb.InventorySerial{
			Id:            serialID,
			SerialNumber:  r.FormValue("serial_number"),
			Imei:          strPtr(r.FormValue("imei")),
			Status:        r.FormValue("status"),
			WarrantyStart: strPtr(r.FormValue("warranty_start")),
			WarrantyEnd:   strPtr(r.FormValue("warranty_end")),
			PurchaseOrder: strPtr(r.FormValue("purchase_order")),
		}

		_, err := deps.UpdateInventorySerial(ctx, &inventoryserialpb.UpdateInventorySerialRequest{Data: data})
		if err != nil {
			log.Printf("Failed to update serial %s: %v", serialID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("serial-table")
	})
}

// NewSerialRemoveAction creates the serial remove action (POST only).
func NewSerialRemoveAction(deps *Deps) view.View {
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
			return centymo.HTMXError(deps.Labels.Errors.SerialIDRequired)
		}

		_, err := deps.DeleteInventorySerial(ctx, &inventoryserialpb.DeleteInventorySerialRequest{
			Data: &inventoryserialpb.InventorySerial{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete serial %s: %v", id, err)
			return centymo.HTMXError(err.Error())
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
