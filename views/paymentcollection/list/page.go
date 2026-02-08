package list

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB centymo.DataSource
}

// PageData holds the data for the payment collection list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the payment collection list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		records, err := deps.DB.ListSimple(ctx, "payment_collection")
		if err != nil {
			log.Printf("Failed to list payment collections: %v", err)
			return view.Error(fmt.Errorf("failed to load payment collections: %w", err))
		}

		columns := paymentCollectionColumns()
		rows := buildTableRows(records, status)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "payment-collections",
				ActiveSubNav:   status,
				HeaderTitle:    statusTitle(status),
				HeaderSubtitle: statusSubtitle(status),
				HeaderIcon:     "icon-credit-card",
			},
			ContentTemplate: "paymentcollection-list-content",
			Table: &types.TableConfig{
				ID:          "payment-collections-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   "No payment collections found",
					Message: "No " + status + " payment collections to display.",
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     "Add Payment Collection",
					ActionURL: "/action/payment-collections/add",
					Icon:      "icon-plus",
				},
			},
		}

		return view.OK("paymentcollection-list", pageData)
	})
}

func paymentCollectionColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "customer", Label: "Customer", Sortable: true},
		{Key: "amount", Label: "Amount", Sortable: true, Width: "120px"},
		{Key: "date", Label: "Date", Sortable: true, Width: "150px"},
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
		customer, _ := record["customer"].(string)
		amount, _ := record["amount"].(string)
		date, _ := record["date"].(string)
		if recordStatus == "" {
			recordStatus = status
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: customer},
				{Type: "text", Value: amount},
				{Type: "text", Value: date},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"customer": customer,
				"amount":   amount,
				"date":     date,
				"status":   recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View Payment", Action: "view", Href: "/app/payment-collections/" + id},
				{Type: "edit", Label: "Edit Payment", Action: "edit", URL: "/action/payment-collections/edit/" + id, DrawerTitle: "Edit Payment Collection"},
				{Type: "delete", Label: "Delete Payment", Action: "delete", URL: "/action/payment-collections/delete", ItemName: customer},
			},
		})
	}
	return rows
}

func statusTitle(status string) string {
	switch status {
	case "pending":
		return "Pending Payments"
	case "completed":
		return "Completed Payments"
	case "failed":
		return "Failed Payments"
	default:
		return "Payment Collections"
	}
}

func statusSubtitle(status string) string {
	switch status {
	case "pending":
		return "Payments awaiting collection"
	case "completed":
		return "Successfully collected payments"
	case "failed":
		return "Failed payment attempts"
	default:
		return "Payment collection management"
	}
}

func statusVariant(status string) string {
	switch status {
	case "pending":
		return "warning"
	case "completed":
		return "success"
	case "failed":
		return "danger"
	default:
		return "default"
	}
}
