package list

import (
	"context"
	"fmt"
	"log"

	"leapfor.xyz/centymo"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB centymo.DataSource
}

// PageData holds the data for the subscription list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the subscription list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		records, err := deps.DB.ListSimple(ctx, "subscription")
		if err != nil {
			log.Printf("Failed to list subscriptions: %v", err)
			return view.Error(fmt.Errorf("failed to load subscriptions: %w", err))
		}

		columns := subscriptionColumns()
		rows := buildTableRows(records, status)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "subscriptions",
				ActiveSubNav:   status,
				HeaderTitle:    statusTitle(status),
				HeaderSubtitle: statusSubtitle(status),
				HeaderIcon:     "icon-refresh-cw",
			},
			ContentTemplate: "subscription-list-content",
			Table: &types.TableConfig{
				ID:          "subscriptions-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   "No subscriptions found",
					Message: "No " + status + " subscriptions to display.",
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     "Add Subscription",
					ActionURL: "/action/subscriptions/add",
					Icon:      "icon-plus",
				},
			},
		}

		return view.OK("subscription-list", pageData)
	})
}

func subscriptionColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "customer", Label: "Customer", Sortable: true},
		{Key: "plan", Label: "Plan", Sortable: true},
		{Key: "start_date", Label: "Start Date", Sortable: true, Width: "150px"},
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
		plan, _ := record["plan"].(string)
		startDate, _ := record["start_date"].(string)
		if recordStatus == "" {
			recordStatus = status
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: customer},
				{Type: "text", Value: plan},
				{Type: "text", Value: startDate},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"customer":   customer,
				"plan":       plan,
				"start_date": startDate,
				"status":     recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View Subscription", Action: "view", Href: "/app/subscriptions/" + id},
				{Type: "edit", Label: "Edit Subscription", Action: "edit", URL: "/action/subscriptions/edit/" + id, DrawerTitle: "Edit Subscription"},
				{Type: "delete", Label: "Cancel Subscription", Action: "delete", URL: "/action/subscriptions/delete", ItemName: customer},
			},
		})
	}
	return rows
}

func statusTitle(status string) string {
	switch status {
	case "active":
		return "Active Subscriptions"
	case "inactive":
		return "Inactive Subscriptions"
	default:
		return "Subscriptions"
	}
}

func statusSubtitle(status string) string {
	switch status {
	case "active":
		return "Manage your active subscriptions"
	case "inactive":
		return "View cancelled or expired subscriptions"
	default:
		return "Subscription management"
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
