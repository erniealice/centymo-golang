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

// PageData holds the data for the plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the plan list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		records, err := deps.DB.ListSimple(ctx, "plan")
		if err != nil {
			log.Printf("Failed to list plans: %v", err)
			return view.Error(fmt.Errorf("failed to load plans: %w", err))
		}

		columns := planColumns()
		rows := buildTableRows(records, status)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "plans",
				ActiveSubNav:   status,
				HeaderTitle:    statusTitle(status),
				HeaderSubtitle: statusSubtitle(status),
				HeaderIcon:     "icon-file-text",
			},
			ContentTemplate: "plan-list-content",
			Table: &types.TableConfig{
				ID:          "plans-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   "No plans found",
					Message: "No " + status + " plans to display.",
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     "Add Plan",
					ActionURL: "/action/plans/add",
					Icon:      "icon-plus",
				},
			},
		}

		return view.OK("plan-list", pageData)
	})
}

func planColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: "Name", Sortable: true},
		{Key: "interval", Label: "Interval", Sortable: true, Width: "150px"},
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
		interval, _ := record["interval"].(string)
		price, _ := record["price"].(string)
		if recordStatus == "" {
			recordStatus = status
		}
		if interval == "" {
			interval = "monthly"
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "badge", Value: interval, Variant: intervalVariant(interval)},
				{Type: "text", Value: price},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":     name,
				"interval": interval,
				"price":    price,
				"status":   recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View Plan", Action: "view", Href: "/app/plans/" + id},
				{Type: "edit", Label: "Edit Plan", Action: "edit", URL: "/action/plans/edit/" + id, DrawerTitle: "Edit Plan"},
				{Type: "delete", Label: "Delete Plan", Action: "delete", URL: "/action/plans/delete", ItemName: name},
			},
		})
	}
	return rows
}

func statusTitle(status string) string {
	switch status {
	case "active":
		return "Active Plans"
	case "inactive":
		return "Inactive Plans"
	default:
		return "Plans"
	}
}

func statusSubtitle(status string) string {
	switch status {
	case "active":
		return "Manage your active billing plans"
	case "inactive":
		return "View inactive billing plans"
	default:
		return "Plan management"
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

func intervalVariant(interval string) string {
	switch interval {
	case "monthly":
		return "info"
	case "annual":
		return "primary"
	default:
		return "default"
	}
}
