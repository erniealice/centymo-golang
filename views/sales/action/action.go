package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
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
	DB centymo.DataSource
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
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("sales-drawer-form", &FormData{
				FormAction:   "/action/sales/add",
				Currency:     "PHP",
				Status:       "ongoing",
				Locations:    loadLocationOptions(ctx, deps.DB),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request

		data := map[string]any{
			"name":                r.FormValue("name"),
			"reference_number":    r.FormValue("reference_number"),
			"revenue_date_string": r.FormValue("revenue_date_string"),
			"currency":            r.FormValue("currency"),
			"status":              r.FormValue("status"),
			"notes":               r.FormValue("notes"),
			"location_id":         r.FormValue("location_id"),
		}

		created, err := deps.DB.Create(ctx, "revenue", data)
		if err != nil {
			log.Printf("Failed to create sale: %v", err)
			return centymo.HTMXError("Failed to create sale")
		}

		// Redirect to new sale detail with Items tab
		newID, _ := created["id"].(string)
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": "/app/sales/detail/" + newID + "?tab=items",
				},
			}
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewEditAction creates the sales edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "revenue", id)
			if err != nil {
				log.Printf("Failed to read sale %s: %v", id, err)
				return centymo.HTMXError("Sale not found")
			}

			name, _ := record["name"].(string)
			refNumber, _ := record["reference_number"].(string)
			date, _ := record["revenue_date_string"].(string)
			currency, _ := record["currency"].(string)
			status, _ := record["status"].(string)
			notes, _ := record["notes"].(string)
			locationID, _ := record["location_id"].(string)

			return view.OK("sales-drawer-form", &FormData{
				FormAction:      "/action/sales/edit/" + id,
				IsEdit:          true,
				ID:              id,
				Name:            name,
				ReferenceNumber: refNumber,
				Date:            date,
				Currency:        currency,
				Status:          status,
				Notes:           notes,
				LocationID:      locationID,
				Locations:       loadLocationOptions(ctx, deps.DB),
				Labels:          formLabels(viewCtx.T),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — update sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request

		data := map[string]any{
			"name":                r.FormValue("name"),
			"reference_number":    r.FormValue("reference_number"),
			"revenue_date_string": r.FormValue("revenue_date_string"),
			"currency":            r.FormValue("currency"),
			"status":              r.FormValue("status"),
			"notes":               r.FormValue("notes"),
			"location_id":         r.FormValue("location_id"),
		}

		_, err := deps.DB.Update(ctx, "revenue", id, data)
		if err != nil {
			log.Printf("Failed to update sale %s: %v", id, err)
			return centymo.HTMXError("Failed to update sale")
		}

		// Redirect to detail page (preserves current tab)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": "/app/sales/detail/" + id,
			},
		}
	})
}

// NewDeleteAction creates the sales delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Sale ID is required")
		}

		err := deps.DB.Delete(ctx, "revenue", id)
		if err != nil {
			log.Printf("Failed to delete sale %s: %v", id, err)
			return centymo.HTMXError("Failed to delete sale")
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// NewBulkDeleteAction creates the sales bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError("No sale IDs provided")
		}

		for _, id := range ids {
			err := deps.DB.Delete(ctx, "revenue", id)
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
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError("Sale ID is required")
		}
		if targetStatus != "ongoing" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError("Invalid status")
		}

		// D20: Block completion with zero items
		if targetStatus == "complete" {
			lineItems, err := getLineItemsForRevenue(ctx, deps.DB, id)
			if err != nil {
				log.Printf("Failed to list line items for sale %s: %v", id, err)
				return centymo.HTMXError("Failed to check sale items")
			}
			if len(lineItems) == 0 {
				return centymo.HTMXError("Cannot complete a sale with no items. Add items first.")
			}

			// Update status
			if _, err := deps.DB.Update(ctx, "revenue", id, map[string]any{"status": targetStatus}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError("Failed to update sale status")
			}

			// D5: Deduct stock on completion
			deductStockForLineItems(ctx, deps.DB, id, lineItems)

			return centymo.HTMXSuccess("sales-table")
		}

		// D21: Block cancellation if payments exist
		if targetStatus == "cancelled" {
			payments, err := getPaymentsForRevenue(ctx, deps.DB, id)
			if err != nil {
				log.Printf("Failed to list payments for sale %s: %v", id, err)
				return centymo.HTMXError("Failed to check sale payments")
			}
			if len(payments) > 0 {
				return centymo.HTMXError("Cannot cancel a sale with recorded payments. Remove payments first.")
			}

			// Update status
			if _, err := deps.DB.Update(ctx, "revenue", id, map[string]any{"status": targetStatus}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError("Failed to update sale status")
			}

			// D6: Release serials on cancellation
			lineItems, err := getLineItemsForRevenue(ctx, deps.DB, id)
			if err != nil {
				log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
			} else {
				releaseSerialsForLineItems(ctx, deps.DB, id, lineItems)
			}

			return centymo.HTMXSuccess("sales-table")
		}

		// Default: ongoing — just update status
		if _, err := deps.DB.Update(ctx, "revenue", id, map[string]any{"status": targetStatus}); err != nil {
			log.Printf("Failed to update sale status %s: %v", id, err)
			return centymo.HTMXError("Failed to update sale status")
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
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError("No sale IDs provided")
		}
		if targetStatus != "ongoing" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError("Invalid target status")
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
					"%d of %d selected sales have recorded payments. Remove payments first.",
					withPayments, len(ids),
				))
			}
		}

		// D20: Block bulk completion if any sale has zero line items
		if targetStatus == "complete" {
			emptyCount := 0
			for _, id := range ids {
				lineItems, err := getLineItemsForRevenue(ctx, deps.DB, id)
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
					"%d of %d selected sales have no items. Add items first.",
					emptyCount, len(ids),
				))
			}
		}

		// Update all statuses and apply side-effects
		for _, id := range ids {
			if _, err := deps.DB.Update(ctx, "revenue", id, map[string]any{"status": targetStatus}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				continue
			}

			// D5: Deduct stock on completion
			if targetStatus == "complete" {
				lineItems, err := getLineItemsForRevenue(ctx, deps.DB, id)
				if err != nil {
					log.Printf("Failed to list line items for stock deduction on sale %s: %v", id, err)
					continue
				}
				deductStockForLineItems(ctx, deps.DB, id, lineItems)
			}

			// D6: Release serials on cancellation
			if targetStatus == "cancelled" {
				lineItems, err := getLineItemsForRevenue(ctx, deps.DB, id)
				if err != nil {
					log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
					continue
				}
				releaseSerialsForLineItems(ctx, deps.DB, id, lineItems)
			}
		}

		return centymo.HTMXSuccess("sales-table")
	})
}

// ---------------------------------------------------------------------------
// Helpers for status change business rules
// ---------------------------------------------------------------------------

// getLineItemsForRevenue returns all revenue_line_item records for a given revenue ID.
func getLineItemsForRevenue(ctx context.Context, db centymo.DataSource, revenueID string) ([]map[string]any, error) {
	all, err := db.ListSimple(ctx, "revenue_line_item")
	if err != nil {
		return nil, err
	}
	var items []map[string]any
	for _, r := range all {
		rid, _ := r["revenue_id"].(string)
		if rid == revenueID {
			items = append(items, r)
		}
	}
	return items, nil
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
func deductStockForLineItems(ctx context.Context, db centymo.DataSource, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		inventoryItemID, _ := item["inventory_item_id"].(string)
		serialID, _ := item["inventory_serial_id"].(string)

		// Deduct quantity from inventory item
		if inventoryItemID != "" {
			invItem, err := db.Read(ctx, "inventory_item", inventoryItemID)
			if err != nil {
				log.Printf("Failed to read inventory item %s for stock deduction: %v", inventoryItemID, err)
				continue
			}

			invQtyStr, _ := invItem["quantity"].(string)
			lineQtyStr, _ := item["quantity"].(string)

			invQty, _ := strconv.ParseFloat(invQtyStr, 64)
			lineQty, _ := strconv.ParseFloat(lineQtyStr, 64)

			newQty := invQty - lineQty
			if _, err := db.Update(ctx, "inventory_item", inventoryItemID, map[string]any{
				"quantity": strconv.FormatFloat(newQty, 'f', -1, 64),
			}); err != nil {
				log.Printf("Failed to deduct stock for inventory item %s: %v", inventoryItemID, err)
			}
		}

		// Mark serial as sold and create history
		if serialID != "" {
			if _, err := db.Update(ctx, "inventory_serial", serialID, map[string]any{
				"status": "sold",
			}); err != nil {
				log.Printf("Failed to mark serial %s as sold: %v", serialID, err)
			}

			if _, err := db.Create(ctx, "inventory_serial_history", map[string]any{
				"inventory_serial_id": serialID,
				"inventory_item_id":   inventoryItemID,
				"from_status":         "reserved",
				"to_status":           "sold",
				"reference_type":      "revenue",
				"reference_id":        saleID,
				"notes":               "Auto: sale completed",
				"changed_by":          "",
				"changed_by_role":     "",
			}); err != nil {
				log.Printf("Failed to create serial history for %s: %v", serialID, err)
			}
		}
	}
}

// releaseSerialsForLineItems marks serials as available and creates history records.
func releaseSerialsForLineItems(ctx context.Context, db centymo.DataSource, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		serialID, _ := item["inventory_serial_id"].(string)
		if serialID == "" {
			continue
		}

		inventoryItemID, _ := item["inventory_item_id"].(string)

		if _, err := db.Update(ctx, "inventory_serial", serialID, map[string]any{
			"status": "available",
		}); err != nil {
			log.Printf("Failed to release serial %s: %v", serialID, err)
		}

		if _, err := db.Create(ctx, "inventory_serial_history", map[string]any{
			"inventory_serial_id": serialID,
			"inventory_item_id":   inventoryItemID,
			"from_status":         "available",
			"to_status":           "available",
			"reference_type":      "revenue",
			"reference_id":        saleID,
			"notes":               "Auto: sale cancelled",
			"changed_by":          "",
			"changed_by_role":     "",
		}); err != nil {
			log.Printf("Failed to create serial history for %s: %v", serialID, err)
		}
	}
}
