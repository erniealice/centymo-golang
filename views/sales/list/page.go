package list

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
	RefreshURL   string
	Labels       centymo.SalesLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the sales list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the sales list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		records, err := deps.DB.ListSimple(ctx, "revenue")
		if err != nil {
			log.Printf("Failed to list sales: %v", err)
			return view.Error(fmt.Errorf("failed to load sales: %w", err))
		}

		l := deps.Labels
		columns := salesColumns(l)
		rows := buildTableRows(records, status, l)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:            "delete",
				Label:          deps.CommonLabels.Bulk.Delete,
				Icon:           "icon-trash-2",
				Variant:        "danger",
				Endpoint:       "/action/sales/bulk-delete",
				ConfirmTitle:   deps.CommonLabels.Bulk.Delete,
				ConfirmMessage: "Are you sure you want to delete {{count}} sale(s)? This action cannot be undone.",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "sales-table",
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
			DefaultSortColumn:    "date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   statusEmptyTitle(l, status),
				Message: statusEmptyMessage(l, status),
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddSale,
				ActionURL: "/action/sales/add",
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
				ActiveNav:      "sales",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sales-list-content",
			Table:           tableConfig,
		}

		return view.OK("sales-list", pageData)
	})
}

func salesColumns(l centymo.SalesLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "reference", Label: l.Columns.Reference, Sortable: true},
		{Key: "customer", Label: l.Columns.Customer, Sortable: true},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "140px"},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(records []map[string]any, status string, l centymo.SalesLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, record := range records {
		recordStatus, _ := record["status"].(string)
		if recordStatus != status {
			continue
		}

		id, _ := record["id"].(string)
		refNumber, _ := record["reference_number"].(string)
		name, _ := record["name"].(string)
		date, _ := record["revenue_date_string"].(string)
		amount, _ := record["total_amount"].(string)
		currency, _ := record["currency"].(string)

		amountDisplay := currency + " " + amount

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: name},
				{Type: "text", Value: date},
				{Type: "text", Value: amountDisplay},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"customer":  name,
				"date":      date,
				"amount":    amountDisplay,
				"status":    recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: "/app/sales/" + id},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/sales/edit/" + id, DrawerTitle: l.Actions.Edit},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: "/action/sales/delete", ItemName: refNumber},
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.SalesLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "completed":
		return l.Page.HeadingCompleted
	case "cancelled":
		return l.Page.HeadingCancelled
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.SalesLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "completed":
		return l.Page.CaptionCompleted
	case "cancelled":
		return l.Page.CaptionCancelled
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.SalesLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "completed":
		return l.Empty.CompletedTitle
	case "cancelled":
		return l.Empty.CancelledTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l centymo.SalesLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "completed":
		return l.Empty.CompletedMessage
	case "cancelled":
		return l.Empty.CancelledMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "info"
	case "completed":
		return "success"
	case "cancelled":
		return "warning"
	default:
		return "default"
	}
}
