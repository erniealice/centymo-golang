package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB           centymo.DataSource
	Labels       centymo.SalesLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the sales detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Revenue         map[string]any
	LineItems       []map[string]any
	Labels          centymo.SalesLabels
	LineItemTable   *types.TableConfig
}

// NewView creates the sales detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		revenue, err := deps.DB.Read(ctx, "revenue", id)
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load sale: %w", err))
		}

		// Fetch line items and filter by revenue_id
		allLineItems, err := deps.DB.ListSimple(ctx, "revenue_line_item")
		if err != nil {
			log.Printf("Failed to list line items for revenue %s: %v", id, err)
			allLineItems = []map[string]any{}
		}

		lineItems := filterLineItems(allLineItems, id)
		lineItemTable := buildLineItemTable(lineItems, deps.Labels, deps.TableLabels)

		refNumber, _ := revenue["reference_number"].(string)
		headerTitle := "Sale #" + refNumber

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sales",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: deps.Labels.Detail.PageTitle,
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sales-detail-content",
			Revenue:         revenue,
			LineItems:       lineItems,
			Labels:          deps.Labels,
			LineItemTable:   lineItemTable,
		}

		return view.OK("sales-detail", pageData)
	})
}

func filterLineItems(all []map[string]any, revenueID string) []map[string]any {
	items := []map[string]any{}
	for _, item := range all {
		rid, _ := item["revenue_id"].(string)
		if rid == revenueID {
			items = append(items, item)
		}
	}
	return items
}

func buildLineItemTable(items []map[string]any, l centymo.SalesLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "description", Label: l.Detail.Description, Sortable: false},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: false, Width: "100px"},
		{Key: "unit_price", Label: l.Detail.UnitPrice, Sortable: false, Width: "140px"},
		{Key: "discount", Label: l.Detail.Discount, Sortable: false, Width: "120px"},
		{Key: "total", Label: l.Detail.Total, Sortable: false, Width: "140px"},
	}

	rows := []types.TableRow{}
	for _, item := range items {
		id, _ := item["id"].(string)
		description, _ := item["description"].(string)
		quantity, _ := item["quantity"].(string)
		unitPrice, _ := item["unit_price"].(string)
		discount, _ := item["discount"].(string)
		total, _ := item["total"].(string)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: description},
				{Type: "text", Value: quantity},
				{Type: "text", Value: unitPrice},
				{Type: "text", Value: discount},
				{Type: "text", Value: total},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "line-items-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No line items",
			Message: "This sale has no line items.",
		},
	}
}
