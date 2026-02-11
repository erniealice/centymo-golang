package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"

	"github.com/erniealice/centymo-golang"
)

// Deps holds view dependencies.
type Deps struct {
	ListProducts func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	RefreshURL   string
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the product list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the product list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err != nil {
			log.Printf("Failed to list products: %v", err)
			return view.Error(fmt.Errorf("failed to load products: %w", err))
		}

		l := deps.Labels
		columns := productColumns(l)
		rows := buildTableRows(resp.GetData(), status, l)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:             "activate",
				Label:           l.Status.Activate,
				Icon:            "icon-check-circle",
				Variant:         "success",
				Endpoint:        centymo.ProductBulkSetStatusURL,
				ExtraParamsJSON: `{"target_status":"active"}`,
				ConfirmTitle:    l.Status.Activate,
				ConfirmMessage:  "Are you sure you want to activate {{count}} product(s)?",
			},
			{
				Key:             "deactivate",
				Label:           l.Status.Deactivate,
				Icon:            "icon-x-circle",
				Variant:         "warning",
				Endpoint:        centymo.ProductBulkSetStatusURL,
				ExtraParamsJSON: `{"target_status":"inactive"}`,
				ConfirmTitle:    l.Status.Deactivate,
				ConfirmMessage:  "Are you sure you want to deactivate {{count}} product(s)?",
			},
			{
				Key:            "delete",
				Label:          l.Bulk.Delete,
				Icon:           "icon-trash-2",
				Variant:        "danger",
				Endpoint:       centymo.ProductBulkDeleteURL,
				ConfirmTitle:   l.Bulk.Delete,
				ConfirmMessage: "Are you sure you want to delete {{count}} product(s)? This action cannot be undone.",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "products-table",
			RefreshURL:           deps.RefreshURL,
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
			DefaultSortColumn:    "name",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   statusEmptyTitle(l, status),
				Message: statusEmptyMessage(l, status),
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddProduct,
				ActionURL: centymo.ProductAddURL,
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "products",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "product-list-content",
			Table:           tableConfig,
		}

		return view.OK("product-list", pageData)
	})
}

func productColumns(l centymo.ProductLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
		{Key: "price", Label: l.Columns.Price, Sortable: true, Width: "150px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(products []*productpb.Product, status string, l centymo.ProductLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range products {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := p.GetId()
		name := p.GetName()
		description := p.GetDescription()
		price := formatPrice(p.GetCurrency(), p.GetPrice())

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: price},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"price":  price,
				"status": recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: "/app/products/detail/" + id},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/products/edit/" + id, DrawerTitle: l.Actions.Edit},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: "/action/products/delete", ItemName: name},
			},
		})
	}
	return rows
}

func formatPrice(currency string, price float64) string {
	if currency == "" {
		currency = "PHP"
	}
	// Format with 2 decimal places, then insert commas for thousands
	raw := fmt.Sprintf("%.2f", price)
	parts := splitDecimal(raw)
	intPart := parts[0]
	decPart := parts[1]

	// Insert commas
	n := len(intPart)
	if n <= 3 {
		return currency + " " + intPart + "." + decPart
	}
	var result []byte
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return currency + " " + string(result) + "." + decPart
}

func splitDecimal(s string) [2]string {
	for i, c := range s {
		if c == '.' {
			return [2]string{s[:i], s[i+1:]}
		}
	}
	return [2]string{s, "00"}
}

func statusPageTitle(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}
