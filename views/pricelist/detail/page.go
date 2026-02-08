package detail

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"

	"github.com/erniealice/centymo-golang"
)

// Deps holds view dependencies.
type Deps struct {
	ReadPriceList     func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	ListPriceProducts func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
	Labels            centymo.PriceListLabels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
}

// PageData holds the data for the price list detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	PriceList       *pricelistpb.PriceList
	ActiveTab       string
	Tabs            []TabConfig
	PricesTable     *types.TableConfig
	Labels          centymo.PriceListLabels
}

// TabConfig defines a single tab in the detail view.
type TabConfig struct {
	Key    string
	Label  string
	Active bool
	URL    string
}

// NewView creates the price list detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		tab := viewCtx.Request.URL.Query().Get("tab")
		if tab == "" {
			tab = "basic"
		}

		resp, err := deps.ReadPriceList(ctx, &pricelistpb.ReadPriceListRequest{
			Data: &pricelistpb.PriceList{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read price list %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load price list: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("price list not found"))
		}
		priceList := data[0]

		name := priceList.GetName()
		description := priceList.GetDescription()

		tabs := []TabConfig{
			{Key: "basic", Label: deps.Labels.Detail.BasicInfo, Active: tab == "basic", URL: fmt.Sprintf("/app/price-lists/%s?tab=basic", id)},
			{Key: "prices", Label: deps.Labels.Detail.Prices, Active: tab == "prices", URL: fmt.Sprintf("/app/price-lists/%s?tab=prices", id)},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          name,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "price-lists",
				HeaderTitle:    name,
				HeaderSubtitle: description,
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "pricelist-detail-content",
			PriceList:       priceList,
			ActiveTab:       tab,
			Tabs:            tabs,
			Labels:          deps.Labels,
		}

		// Populate prices table on the "prices" tab
		if tab == "prices" {
			pricesTable, err := buildPricesTable(ctx, deps, id)
			if err != nil {
				log.Printf("Failed to load price products for price list %s: %v", id, err)
			}
			pageData.PricesTable = pricesTable
		}

		return view.OK("pricelist-detail", pageData)
	})
}

func buildPricesTable(ctx context.Context, deps *Deps, priceListID string) (*types.TableConfig, error) {
	resp, err := deps.ListPriceProducts(ctx, &priceproductpb.ListPriceProductsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list price products: %w", err)
	}

	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "product_name", Label: l.Detail.ProductName, Sortable: true},
		{Key: "amount", Label: l.Detail.Amount, Sortable: true, Width: "150px"},
		{Key: "currency", Label: l.Detail.Currency, Sortable: true, Width: "120px"},
	}

	rows := []types.TableRow{}
	for _, pp := range resp.GetData() {
		// Filter price products belonging to this price list
		// The price_product proto doesn't have a direct price_list_id field in its
		// ListPriceProductsRequest, so we filter client-side for now
		id := pp.GetId()
		productName := pp.GetName()
		amount := fmt.Sprintf("%d", pp.GetAmount())
		currency := pp.GetCurrency()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: productName},
				{Type: "text", Value: amount},
				{Type: "text", Value: currency},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "price-products-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		DefaultSortColumn:    "product_name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No prices configured",
			Message: "Add products to this price list to configure pricing.",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}
