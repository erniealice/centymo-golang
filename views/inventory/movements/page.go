package movements

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
)

// Deps holds view dependencies.
type Deps struct {
	SqlDB                     *sql.DB // raw DB for filtered queries (nil = fallback to proto functions)
	ListInventoryItems        func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ListInventoryTransactions func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	ListLocations             func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	Labels                    centymo.InventoryLabels
	CommonLabels              pyeza.CommonLabels
	TableLabels               types.TableLabels
}

// LocationOption represents a location for the filter dropdown.
type LocationOption struct {
	ID   string
	Name string
}

// TransactionTypeOption represents a transaction type for the filter dropdown.
type TransactionTypeOption struct {
	Value string
	Label string
}

// PageData holds the data for the movements page.
type PageData struct {
	types.PageData
	ContentTemplate  string
	Table            *types.TableConfig
	Labels           centymo.InventoryMovementsLabels
	DateFrom         string
	DateTo           string
	LocationFilter   string
	TypeFilter       string
	Search           string
	Locations        []LocationOption
	TransactionTypes []TransactionTypeOption
}

// transactionTypes returns the list of known transaction types for filters.
func transactionTypes(l centymo.InventoryTransactionLabels) []TransactionTypeOption {
	return []TransactionTypeOption{
		{Value: "received", Label: l.TypeReceived},
		{Value: "sold", Label: l.TypeSold},
		{Value: "adjusted", Label: l.TypeAdjusted},
		{Value: "transferred", Label: l.TypeTransferred},
		{Value: "returned", Label: l.TypeReturned},
		{Value: "write_off", Label: l.TypeWriteOff},
	}
}

// NewView creates the global inventory movements page.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels

		// Default date range: 1st of current month to today
		now := time.Now()
		dateFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		dateTo := now.Format("2006-01-02")

		// Load locations for dropdown
		locations := loadLocations(ctx, deps)

		// Query transactions (filtered by default date range if SQL available)
		tableConfig := buildFilteredTable(ctx, deps, dateFrom, dateTo, "", "", "")

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Transactions",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				ActiveSubNav:   "movements",
				HeaderTitle:    "Transactions",
				HeaderSubtitle: l.Movements.Subtitle,
				HeaderIcon:     "icon-repeat",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "inventory-movements-content",
			Table:            tableConfig,
			Labels:           l.Movements,
			DateFrom:         dateFrom,
			DateTo:           dateTo,
			Locations:        locations,
			TransactionTypes: transactionTypes(l.Transaction),
		}

		return view.OK("inventory-movements", pageData)
	})
}

func txTypeVariant(txType string) string {
	switch txType {
	case "received":
		return "success"
	case "sold":
		return "default"
	case "adjusted":
		return "info"
	case "transferred":
		return "warning"
	case "returned":
		return "danger"
	case "write_off":
		return "danger"
	default:
		return "default"
	}
}

func formatQuantity(qty any, txType string) string {
	val := fmt.Sprintf("%v", qty)
	switch txType {
	case "received", "returned":
		return "+" + val
	case "sold", "transferred", "write_off":
		return "-" + val
	default:
		return val
	}
}

// loadLocations fetches locations from the DB for filter dropdowns.
func loadLocations(ctx context.Context, deps *Deps) []LocationOption {
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		log.Printf("Failed to list locations for movements filter: %v", err)
		return nil
	}
	locs := resp.GetData()
	options := make([]LocationOption, 0, len(locs))
	for _, loc := range locs {
		id := loc.GetId()
		name := loc.GetName()
		if id != "" && name != "" {
			options = append(options, LocationOption{ID: id, Name: name})
		}
	}
	return options
}

// buildFilteredTable builds the table config with optional SQL filtering.
// If SqlDB is nil, falls back to in-memory ListSimple with no filters.
func buildFilteredTable(ctx context.Context, deps *Deps, dateFrom, dateTo, location, txType, search string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "transaction_date", Label: l.Detail.Date, Sortable: true, Width: "130px"},
		{Key: "item_name", Label: l.Columns.ProductName, Sortable: true},
		{Key: "product_name", Label: "Product", Sortable: true},
		{Key: "variant_sku", Label: "Variant SKU", Sortable: true, Width: "130px"},
		{Key: "sku", Label: l.Columns.SKU, Sortable: true, Width: "130px"},
		{Key: "location", Label: l.Detail.Location, Sortable: true, Width: "160px"},
		{Key: "transaction_type", Label: l.Detail.Type, Sortable: true, Width: "120px"},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: true, Width: "100px"},
		{Key: "serial_number", Label: l.Detail.Serial, Sortable: false, Width: "150px"},
		{Key: "reference", Label: l.Detail.Reference, Sortable: false},
		{Key: "performed_by", Label: l.Detail.PerformedBy, Sortable: false, Width: "150px"},
	}

	var rows []types.TableRow

	if deps.SqlDB != nil {
		rows = queryFilteredRows(ctx, deps, dateFrom, dateTo, location, txType, search)
	} else {
		rows = queryFallbackRows(ctx, deps)
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "movements-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "transaction_date",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.TransactionEmptyTitle,
			Message: l.Detail.TransactionEmptyMessage,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

// queryFilteredRows uses raw SQL with JOINs and filters.
func queryFilteredRows(ctx context.Context, deps *Deps, dateFrom, dateTo, location, txType, search string) []types.TableRow {
	query := `
		SELECT it.id, it.transaction_date, it.transaction_type, it.quantity,
		       it.serial_number, it.reference_type, it.reference_id, it.performed_by,
		       COALESCE(ii.name, '') as item_name,
		       COALESCE(ii.location_id, '') as location_id,
		       COALESCE(ii.sku, '') as item_sku,
		       COALESCE(pv.sku, '') as variant_sku,
		       COALESCE(p.name, '') as product_name
		FROM inventory_transaction it
		LEFT JOIN inventory_item ii ON it.inventory_item_id = ii.id
		LEFT JOIN product_variant pv ON ii.product_variant_id = pv.id
		LEFT JOIN product p ON pv.product_id = p.id
		WHERE it.active = true
		  AND ($1 = '' OR it.transaction_date >= $1::timestamptz)
		  AND ($2 = '' OR it.transaction_date <= ($2::date + interval '1 day')::timestamptz)
		  AND ($3 = '' OR ii.location_id = $3)
		  AND ($4 = '' OR it.transaction_type = $4)
		  AND ($5 = '' OR (
		       p.name ILIKE '%' || $5 || '%'
		    OR pv.sku ILIKE '%' || $5 || '%'
		    OR ii.sku ILIKE '%' || $5 || '%'
		    OR ii.name ILIKE '%' || $5 || '%'
		  ))
		ORDER BY it.transaction_date DESC
	`

	sqlRows, err := deps.SqlDB.QueryContext(ctx, query, dateFrom, dateTo, location, txType, search)
	if err != nil {
		log.Printf("Failed to query filtered movements: %v", err)
		return nil
	}
	defer sqlRows.Close()

	var rows []types.TableRow
	for sqlRows.Next() {
		var id, txTypeVal string
		var itemName, locationID, itemSKU, variantSKU, productName string
		var qty float64
		var txDateNullable sql.NullTime
		var serial, refType, refID, performer sql.NullString

		if err := sqlRows.Scan(
			&id, &txDateNullable, &txTypeVal, &qty,
			&serial, &refType, &refID, &performer,
			&itemName, &locationID, &itemSKU, &variantSKU, &productName,
		); err != nil {
			log.Printf("Failed to scan movement row: %v", err)
			continue
		}

		txDate := ""
		if txDateNullable.Valid {
			txDate = txDateNullable.Time.Format("2006-01-02")
		}

		ref := refType.String
		if refID.String != "" {
			if ref != "" {
				ref += ": " + refID.String
			} else {
				ref = refID.String
			}
		}

		locationName := centymo.LocationDisplayName(locationID)
		qtyStr := formatQuantity(qty, txTypeVal)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: txDate},
				{Type: "text", Value: itemName},
				{Type: "text", Value: productName},
				{Type: "text", Value: variantSKU},
				{Type: "text", Value: itemSKU},
				{Type: "text", Value: locationName},
				{Type: "badge", Value: txTypeVal, Variant: txTypeVariant(txTypeVal)},
				{Type: "text", Value: qtyStr},
				{Type: "text", Value: serial.String},
				{Type: "text", Value: ref},
				{Type: "text", Value: performer.String},
			},
		})
	}
	return rows
}

// queryFallbackRows uses the typed proto functions (no SQL filters, no JOINs).
func queryFallbackRows(ctx context.Context, deps *Deps) []types.TableRow {
	txnResp, err := deps.ListInventoryTransactions(ctx, &inventorytransactionpb.ListInventoryTransactionsRequest{})
	if err != nil {
		log.Printf("Failed to list inventory_transaction: %v", err)
	}
	var transactions []*inventorytransactionpb.InventoryTransaction
	if txnResp != nil {
		transactions = txnResp.GetData()
	}

	itemResp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Failed to list inventory_item for movements: %v", err)
	}
	var items []*inventoryitempb.InventoryItem
	if itemResp != nil {
		items = itemResp.GetData()
	}
	itemMap := map[string]*inventoryitempb.InventoryItem{}
	for _, item := range items {
		itemMap[item.GetId()] = item
	}

	var rows []types.TableRow
	for _, t := range transactions {
		id := t.GetId()
		txDate := t.GetTransactionDateString()
		txType := t.GetTransactionType()
		qty := formatQuantity(t.GetQuantity(), txType)
		ref := t.GetReferenceType()
		serial := t.GetSerialNumber()
		performer := t.GetPerformedBy()

		inventoryItemID := t.GetInventoryItemId()
		itemName := inventoryItemID
		locationName := ""
		itemSKU := ""
		if item, ok := itemMap[inventoryItemID]; ok {
			name := item.GetName()
			if name != "" {
				itemName = name
			}
			locationName = centymo.LocationDisplayName(item.GetLocationId())
			itemSKU = item.GetSku()
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: txDate},
				{Type: "text", Value: itemName},
				{Type: "text", Value: ""},     // product_name (not available without JOIN)
				{Type: "text", Value: ""},     // variant_sku (not available without JOIN)
				{Type: "text", Value: itemSKU},
				{Type: "text", Value: locationName},
				{Type: "badge", Value: txType, Variant: txTypeVariant(txType)},
				{Type: "text", Value: qty},
				{Type: "text", Value: serial},
				{Type: "text", Value: ref},
				{Type: "text", Value: performer},
			},
		})
	}
	return rows
}
