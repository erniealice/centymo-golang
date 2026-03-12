package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialhistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Customer             string
	Date                 string
	Currency             string
	Reference            string
	ReferencePlaceholder string
	Status               string
	Notes                string
	NotesPlaceholder     string
	Location             string
}

// FormData is the template data for the sales drawer form.
type FormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	Name            string
	ReferenceNumber string
	Date            string
	Currency        string
	Status          string
	Notes           string
	LocationID      string
	Locations       []map[string]string
	Labels          FormLabels
	CommonLabels    any
}

// Deps holds dependencies for sales action handlers.
type Deps struct {
	Routes centymo.RevenueRoutes
	Labels centymo.RevenueLabels
	DB     centymo.DataSource // KEEP — used for location, revenue_payment, and collection_method operations

	// Typed revenue operations
	CreateRevenue func(ctx context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
	ReadRevenue   func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	UpdateRevenue func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
	DeleteRevenue func(ctx context.Context, req *revenuepb.DeleteRevenueRequest) (*revenuepb.DeleteRevenueResponse, error)

	// Typed line item operations
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Typed inventory operations
	ReadInventoryItem            func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem          func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	ListInventoryItems           func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	UpdateInventorySerial        func(ctx context.Context, req *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	CreateInventorySerialHistory func(ctx context.Context, req *serialhistorypb.CreateInventorySerialHistoryRequest) (*serialhistorypb.CreateInventorySerialHistoryResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Customer:             t("sales.form.customer"),
		Date:                 t("sales.form.date"),
		Currency:             t("sales.form.currency"),
		Reference:            t("sales.form.reference"),
		ReferencePlaceholder: t("sales.form.referencePlaceholder"),
		Status:               t("sales.form.status"),
		Notes:                t("sales.form.notes"),
		NotesPlaceholder:     t("sales.form.notesPlaceholder"),
		Location:             t("sales.form.location"),
	}
}

// loadLocationOptions loads active locations for the dropdown.
func loadLocationOptions(ctx context.Context, db centymo.DataSource) []map[string]string {
	records, err := db.ListSimple(ctx, "location")
	if err != nil {
		log.Printf("Failed to list locations: %v", err)
		return nil
	}

	options := []map[string]string{}
	for _, r := range records {
		active, _ := r["active"].(bool)
		if !active {
			continue
		}
		id, _ := r["id"].(string)
		name, _ := r["name"].(string)
		if id == "" {
			continue
		}
		options = append(options, map[string]string{
			"Value": id,
			"Label": name,
		})
	}
	return options
}

// NewAddAction creates the sales add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("sales-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Currency:     "PHP",
				Status:       "ongoing",
				Locations:    loadLocationOptions(ctx, deps.DB),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		resp, err := deps.CreateRevenue(ctx, &revenuepb.CreateRevenueRequest{
			Data: &revenuepb.Revenue{
				Name:              r.FormValue("name"),
				ReferenceNumber:   strPtr(r.FormValue("reference_number")),
				RevenueDateString: strPtr(r.FormValue("revenue_date_string")),
				Currency:          r.FormValue("currency"),
				Status:            r.FormValue("status"),
				Notes:             strPtr(r.FormValue("notes")),
				LocationId:        r.FormValue("location_id"),
			},
		})
		if err != nil {
			log.Printf("Failed to create sale: %v", err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to new sale detail with Items tab
		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", newID) + "?tab=items",
				},
			}
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewEditAction creates the sales edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
				Data: &revenuepb.Revenue{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read sale %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			return view.OK("sales-drawer-form", &FormData{
				FormAction:      route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:          true,
				ID:              id,
				Name:            record.GetName(),
				ReferenceNumber: record.GetReferenceNumber(),
				Date:            record.GetRevenueDateString(),
				Currency:        record.GetCurrency(),
				Status:          record.GetStatus(),
				Notes:           record.GetNotes(),
				LocationID:      record.GetLocationId(),
				Locations:       loadLocationOptions(ctx, deps.DB),
				Labels:          formLabels(viewCtx.T),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — update sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
			Data: &revenuepb.Revenue{
				Id:                id,
				Name:              r.FormValue("name"),
				ReferenceNumber:   strPtr(r.FormValue("reference_number")),
				RevenueDateString: strPtr(r.FormValue("revenue_date_string")),
				Currency:          r.FormValue("currency"),
				Status:            r.FormValue("status"),
				Notes:             strPtr(r.FormValue("notes")),
				LocationId:        r.FormValue("location_id"),
			},
		})
		if err != nil {
			log.Printf("Failed to update sale %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to detail page (preserves current tab)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", id),
			},
		}
	})
}

// NewDeleteAction creates the sales delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "delete") {
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

		_, err := deps.DeleteRevenue(ctx, &revenuepb.DeleteRevenueRequest{
			Data: &revenuepb.Revenue{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete sale %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewBulkDeleteAction creates the sales bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			_, err := deps.DeleteRevenue(ctx, &revenuepb.DeleteRevenueRequest{
				Data: &revenuepb.Revenue{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete sale %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewSetStatusAction creates the sales status update action (POST only).
// Expects query params: ?id={saleId}&status={ongoing|complete|cancelled}
//
// Business rules:
//   - D20: Block completion with zero line items
//   - D21: Block cancellation if payments exist
//   - D5: Deduct stock on completion
//   - D6: Release serials on cancellation
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
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
		if targetStatus != "ongoing" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		// D20: Block completion with zero items
		if targetStatus == "complete" {
			lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
			if err != nil {
				log.Printf("Failed to list line items for sale %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
			if len(lineItems) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NoItemsCannotComplete)
			}

			// Update status
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}

			// D5: Deduct stock on completion
			deductStockForLineItems(ctx, deps, id, lineItems)

			return centymo.HTMXSuccess("sales-table")
		}

		// D21: Block cancellation if payments exist
		if targetStatus == "cancelled" {
			payments, err := getPaymentsForRevenue(ctx, deps.DB, id)
			if err != nil {
				log.Printf("Failed to list payments for sale %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
			if len(payments) > 0 {
				return centymo.HTMXError(deps.Labels.Errors.HasPaymentsCannotCancel)
			}

			// Update status
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}

			// D6: Release serials on cancellation
			lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
			if err != nil {
				log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
			} else {
				releaseSerialsForLineItems(ctx, deps, id, lineItems)
			}

			return centymo.HTMXSuccess("sales-table")
		}

		// Default: ongoing — just update status
		if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
			Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
		}); err != nil {
			log.Printf("Failed to update sale status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewBulkSetStatusAction creates the sales bulk status update action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
//
// Business rules:
//   - D20: Block bulk completion if any sale has zero line items
//   - D21: Block bulk cancellation if any sale has payments
//   - D5: Deduct stock on completion for each sale
//   - D6: Release serials on cancellation for each sale
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "ongoing" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidTargetStatus)
		}

		// D21: Block bulk cancellation if any sale has payments
		if targetStatus == "cancelled" {
			withPayments := 0
			for _, id := range ids {
				payments, err := getPaymentsForRevenue(ctx, deps.DB, id)
				if err != nil {
					log.Printf("Failed to check payments for sale %s: %v", id, err)
					continue
				}
				if len(payments) > 0 {
					withPayments++
				}
			}
			if withPayments > 0 {
				return centymo.HTMXError(fmt.Sprintf(
					deps.Labels.Errors.BulkHasPayments,
					withPayments, len(ids),
				))
			}
		}

		// D20: Block bulk completion if any sale has zero line items
		if targetStatus == "complete" {
			emptyCount := 0
			for _, id := range ids {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to check line items for sale %s: %v", id, err)
					continue
				}
				if len(lineItems) == 0 {
					emptyCount++
				}
			}
			if emptyCount > 0 {
				return centymo.HTMXError(fmt.Sprintf(
					deps.Labels.Errors.BulkNoItems,
					emptyCount, len(ids),
				))
			}
		}

		// Update all statuses and apply side-effects
		for _, id := range ids {
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				continue
			}

			// D5: Deduct stock on completion
			if targetStatus == "complete" {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to list line items for stock deduction on sale %s: %v", id, err)
					continue
				}
				deductStockForLineItems(ctx, deps, id, lineItems)
			}

			// D6: Release serials on cancellation
			if targetStatus == "cancelled" {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
					continue
				}
				releaseSerialsForLineItems(ctx, deps, id, lineItems)
			}
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// ---------------------------------------------------------------------------
// Helpers for status change business rules
// ---------------------------------------------------------------------------

// getLineItemsForRevenueTyped returns all revenue_line_item records for a given revenue ID using typed use case.
func getLineItemsForRevenueTyped(
	ctx context.Context,
	listFn func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error),
	revenueID string,
) ([]map[string]any, error) {
	resp, err := listFn(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
		RevenueId: &revenueID,
	})
	if err != nil {
		return nil, err
	}
	var items []map[string]any
	for _, item := range resp.GetData() {
		if item.GetRevenueId() == revenueID {
			items = append(items, map[string]any{
				"id":                  item.GetId(),
				"revenue_id":          item.GetRevenueId(),
				"description":         item.GetDescription(),
				"quantity":            fmt.Sprintf("%.0f", item.GetQuantity()),
				"unit_price":          fmt.Sprintf("%.2f", item.GetUnitPrice()),
				"cost_price":          fmt.Sprintf("%.2f", item.GetCostPrice()),
				"total":               fmt.Sprintf("%.2f", item.GetTotalPrice()),
				"line_item_type":      item.GetLineItemType(),
				"inventory_item_id":   item.GetInventoryItemId(),
				"inventory_serial_id": item.GetInventorySerialId(),
				"notes":               item.GetNotes(),
			})
		}
	}
	return items, nil
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// getPaymentsForRevenue returns all revenue_payment records for a given revenue ID.
func getPaymentsForRevenue(ctx context.Context, db centymo.DataSource, revenueID string) ([]map[string]any, error) {
	all, err := db.ListSimple(ctx, "revenue_payment")
	if err != nil {
		return nil, err
	}
	var payments []map[string]any
	for _, r := range all {
		rid, _ := r["revenue_id"].(string)
		if rid == revenueID {
			payments = append(payments, r)
		}
	}
	return payments, nil
}

// deductStockForLineItems decrements inventory quantities and marks serials as sold.
func deductStockForLineItems(ctx context.Context, deps *Deps, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		inventoryItemID, _ := item["inventory_item_id"].(string)
		serialID, _ := item["inventory_serial_id"].(string)

		// Deduct quantity from inventory item
		if inventoryItemID != "" {
			resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{Id: inventoryItemID},
			})
			if err != nil {
				log.Printf("Failed to read inventory item %s for stock deduction: %v", inventoryItemID, err)
				continue
			}
			data := resp.GetData()
			if len(data) == 0 {
				log.Printf("Inventory item %s not found for stock deduction", inventoryItemID)
				continue
			}
			invItem := data[0]

			lineQtyStr, _ := item["quantity"].(string)
			lineQty, _ := strconv.ParseFloat(lineQtyStr, 64)

			newQty := invItem.GetQuantityOnHand() - lineQty
			if _, err := deps.UpdateInventoryItem(ctx, &inventoryitempb.UpdateInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{
					Id:             inventoryItemID,
					QuantityOnHand: newQty,
				},
			}); err != nil {
				log.Printf("Failed to deduct stock for inventory item %s: %v", inventoryItemID, err)
			}
		}

		// Mark serial as sold and create history
		if serialID != "" {
			if _, err := deps.UpdateInventorySerial(ctx, &inventoryserialpb.UpdateInventorySerialRequest{
				Data: &inventoryserialpb.InventorySerial{
					Id:     serialID,
					Status: "sold",
				},
			}); err != nil {
				log.Printf("Failed to mark serial %s as sold: %v", serialID, err)
			}

			if _, err := deps.CreateInventorySerialHistory(ctx, &serialhistorypb.CreateInventorySerialHistoryRequest{
				Data: &serialhistorypb.InventorySerialHistory{
					InventorySerialId: serialID,
					InventoryItemId:   inventoryItemID,
					FromStatus:        "reserved",
					ToStatus:          "sold",
					ReferenceType:     "revenue",
					ReferenceId:       saleID,
					Notes:             "Auto: sale completed",
				},
			}); err != nil {
				log.Printf("Failed to create serial history for %s: %v", serialID, err)
			}
		}
	}
}

// releaseSerialsForLineItems marks serials as available and creates history records.
func releaseSerialsForLineItems(ctx context.Context, deps *Deps, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		serialID, _ := item["inventory_serial_id"].(string)
		if serialID == "" {
			continue
		}

		inventoryItemID, _ := item["inventory_item_id"].(string)

		if _, err := deps.UpdateInventorySerial(ctx, &inventoryserialpb.UpdateInventorySerialRequest{
			Data: &inventoryserialpb.InventorySerial{
				Id:     serialID,
				Status: "available",
			},
		}); err != nil {
			log.Printf("Failed to release serial %s: %v", serialID, err)
		}

		if _, err := deps.CreateInventorySerialHistory(ctx, &serialhistorypb.CreateInventorySerialHistoryRequest{
			Data: &serialhistorypb.InventorySerialHistory{
				InventorySerialId: serialID,
				InventoryItemId:   inventoryItemID,
				FromStatus:        "available",
				ToStatus:          "available",
				ReferenceType:     "revenue",
				ReferenceId:       saleID,
				Notes:             "Auto: sale cancelled",
			},
		}); err != nil {
			log.Printf("Failed to create serial history for %s: %v", serialID, err)
		}
	}
}
