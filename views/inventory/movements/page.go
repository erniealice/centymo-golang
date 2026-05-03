package movements

import (
	"context"
	"fmt"
	"log"
	"time"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
)

// Deps holds view dependencies.
type Deps struct {
	GetInventoryMovementsListPageData func(ctx context.Context, req *inventorytransactionpb.GetInventoryMovementsListPageDataRequest) (*inventorytransactionpb.GetInventoryMovementsListPageDataResponse, error)
	ListInventoryItems                func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ListInventoryTransactions         func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	ListLocations                     func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	Labels                            centymo.InventoryLabels
	CommonLabels                      pyeza.CommonLabels
	TableLabels                       types.TableLabels
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

		// Query transactions (filtered by default date range if use case available)
		tableConfig := buildFilteredTable(ctx, deps, dateFrom, dateTo, "", "", "")

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Movements.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				ActiveSubNav:   "movements",
				HeaderTitle:    l.Movements.Title,
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

// buildFilteredTable builds the table config using the typed use case when available,
// falling back to in-memory ListSimple with no filters.
func buildFilteredTable(ctx context.Context, deps *Deps, dateFrom, dateTo, location, txType, search string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "transaction_date", Label: l.Detail.Date, WidthClass: "col-3xl"},
		{Key: "item_name", Label: l.Columns.ProductName},
		{Key: "product_name", Label: l.Movements.ProductColumn},
		{Key: "variant_sku", Label: l.Movements.VariantSKU, WidthClass: "col-3xl"},
		{Key: "sku", Label: l.Columns.SKU, WidthClass: "col-3xl"},
		{Key: "location", Label: l.Detail.Location, WidthClass: "col-5xl"},
		{Key: "transaction_type", Label: l.Detail.Type, WidthClass: "col-2xl"},
		{Key: "quantity", Label: l.Detail.Quantity, WidthClass: "col-lg"},
		{Key: "serial_number", Label: l.Detail.Serial, NoSort: true, WidthClass: "col-4xl"},
		{Key: "reference", Label: l.Detail.Reference, NoSort: true},
		{Key: "performed_by", Label: l.Detail.PerformedBy, NoSort: true, WidthClass: "col-4xl"},
	}

	var rows []types.TableRow

	if deps.GetInventoryMovementsListPageData != nil {
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

// queryFilteredRows calls the typed use case and maps the proto response to table rows.
func queryFilteredRows(ctx context.Context, deps *Deps, dateFrom, dateTo, location, txType, search string) []types.TableRow {
	req := &inventorytransactionpb.GetInventoryMovementsListPageDataRequest{}
	if dateFrom != "" {
		req.DateFrom = &dateFrom
	}
	if dateTo != "" {
		req.DateTo = &dateTo
	}
	if location != "" {
		req.LocationId = &location
	}
	if txType != "" {
		req.TransactionType = &txType
	}
	if search != "" {
		req.Search = &search
	}

	resp, err := deps.GetInventoryMovementsListPageData(ctx, req)
	if err != nil {
		log.Printf("Failed to query filtered movements: %v", err)
		return nil
	}

	movements := resp.GetData()
	rows := make([]types.TableRow, 0, len(movements))
	for _, m := range movements {
		txTypeVal := m.GetTransactionType()

		ref := m.GetReferenceType()
		if refID := m.GetReferenceId(); refID != "" {
			if ref != "" {
				ref += ": " + refID
			} else {
				ref = refID
			}
		}

		locationName := centymo.LocationDisplayName(m.GetLocationId())
		qtyStr := formatQuantity(m.GetQuantity(), txTypeVal)

		rows = append(rows, types.TableRow{
			ID: m.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: m.GetTransactionDate()},
				{Type: "text", Value: m.GetItemName()},
				{Type: "text", Value: m.GetProductName()},
				{Type: "text", Value: m.GetVariantSku()},
				{Type: "text", Value: m.GetItemSku()},
				{Type: "text", Value: locationName},
				{Type: "badge", Value: txTypeVal, Variant: txTypeVariant(txTypeVal)},
				{Type: "text", Value: qtyStr},
				{Type: "text", Value: m.GetSerialNumber()},
				{Type: "text", Value: ref},
				{Type: "text", Value: m.GetPerformedBy()},
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
				{Type: "text", Value: ""}, // product_name (not available without JOIN)
				{Type: "text", Value: ""}, // variant_sku (not available without JOIN)
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
