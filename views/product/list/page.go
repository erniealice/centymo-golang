package list

import (
	"context"
	"fmt"
	"log"

	"leapfor.xyz/centymo"

	"github.com/erniealice/pyeza-golang/types"
)

// Deps holds view dependencies.
type Deps struct {
	DB centymo.DataSource
}

// PageData holds the data for the product list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the product list view.
func NewView(deps *Deps) centymo.View {
	return centymo.ViewFunc(func(ctx context.Context, viewCtx *centymo.ViewContext) centymo.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		records, err := deps.DB.ListSimple(ctx, "product")
		if err != nil {
			log.Printf("Failed to list products: %v", err)
			return centymo.Error(fmt.Errorf("failed to load products: %w", err))
		}

		columns := productColumns()
		rows := buildTableRows(records, status)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "products",
				ActiveSubNav:   status,
				HeaderTitle:    statusTitle(status),
				HeaderSubtitle: statusSubtitle(status),
				HeaderIcon:     "icon-package",
			},
			ContentTemplate: "product-list-content",
			Table: &types.TableConfig{
				ID:          "products-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   "No products found",
					Message: "No " + status + " products to display.",
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     "Add Product",
					ActionURL: "/action/products/add",
					Icon:      "icon-plus",
				},
			},
		}

		return centymo.OK("product-list", pageData)
	})
}

func productColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: "Name", Sortable: true},
		{Key: "sku", Label: "SKU", Sortable: true, Width: "150px"},
		{Key: "price", Label: "Price", Sortable: true, Width: "120px"},
		{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
	}
}

func buildTableRows(records []map[string]any, status string) []types.TableRow {
	rows := []types.TableRow{}
	for _, record := range records {
		recordStatus, _ := record["status"].(string)
		if recordStatus != "" && recordStatus != status {
			continue
		}

		id, _ := record["id"].(string)
		name, _ := record["name"].(string)
		sku, _ := record["sku"].(string)
		price, _ := record["price"].(string)
		if recordStatus == "" {
			recordStatus = status
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: sku},
				{Type: "text", Value: price},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"sku":    sku,
				"price":  price,
				"status": recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View Product", Action: "view", Href: "/app/products/" + id},
				{Type: "edit", Label: "Edit Product", Action: "edit", URL: "/action/products/edit/" + id, DrawerTitle: "Edit Product"},
				{Type: "delete", Label: "Delete Product", Action: "delete", URL: "/action/products/delete", ItemName: name},
			},
		})
	}
	return rows
}

func statusTitle(status string) string {
	switch status {
	case "active":
		return "Active Products"
	case "inactive":
		return "Inactive Products"
	default:
		return "Products"
	}
}

func statusSubtitle(status string) string {
	switch status {
	case "active":
		return "Manage your active products"
	case "inactive":
		return "View discontinued products"
	default:
		return "Product management"
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
