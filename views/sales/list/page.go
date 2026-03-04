package list

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.SalesRoutes
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
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "ongoing"
		}

		records, err := deps.DB.ListSimple(ctx, "revenue")
		if err != nil {
			log.Printf("Failed to list sales: %v", err)
			return view.Error(fmt.Errorf("failed to load sales: %w", err))
		}

		l := deps.Labels
		columns := salesColumns(l)
		rows := buildTableRows(records, status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(deps.CommonLabels, status, deps.Routes)

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
				Label:           l.Buttons.AddSale,
				ActionURL:       deps.Routes.AddURL,
				Icon:            "icon-plus",
				Disabled:        !perms.Can("invoice", "create"),
				DisabledTooltip: "No permission",
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

func buildTableRows(records []map[string]any, status string, l centymo.SalesLabels, routes centymo.SalesRoutes, perms *types.UserPermissions) []types.TableRow {
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

		detailURL := route.ResolveURL(routes.DetailURL, "id", id)
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("invoice", "update"), DisabledTooltip: "No permission"},
		}
		if recordStatus == "ongoing" {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: "Complete", Action: "deactivate",
				URL: routes.SetStatusURL + "?status=complete", ItemName: refNumber,
				ConfirmTitle:    "Complete",
				ConfirmMessage:  fmt.Sprintf("Are you sure you want to mark %s as complete?", refNumber),
				Disabled:        !perms.Can("invoice", "update"),
				DisabledTooltip: "No permission",
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: "Reactivate", Action: "activate",
				URL: routes.SetStatusURL + "?status=ongoing", ItemName: refNumber,
				ConfirmTitle:    "Reactivate",
				ConfirmMessage:  fmt.Sprintf("Are you sure you want to reactivate %s?", refNumber),
				Disabled:        !perms.Can("invoice", "update"),
				DisabledTooltip: "No permission",
			})
		}
		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
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
			Actions: actions,
		})
	}
	return rows
}

func statusPageTitle(l centymo.SalesLabels, status string) string {
	switch status {
	case "ongoing":
		return l.Page.HeadingOngoing
	case "complete":
		return l.Page.HeadingComplete
	case "cancelled":
		return l.Page.HeadingCancelled
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.SalesLabels, status string) string {
	switch status {
	case "ongoing":
		return l.Page.CaptionOngoing
	case "complete":
		return l.Page.CaptionComplete
	case "cancelled":
		return l.Page.CaptionCancelled
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.SalesLabels, status string) string {
	switch status {
	case "ongoing":
		return l.Empty.OngoingTitle
	case "complete":
		return l.Empty.CompleteTitle
	case "cancelled":
		return l.Empty.CancelledTitle
	default:
		return l.Empty.OngoingTitle
	}
}

func statusEmptyMessage(l centymo.SalesLabels, status string) string {
	switch status {
	case "ongoing":
		return l.Empty.OngoingMessage
	case "complete":
		return l.Empty.CompleteMessage
	case "cancelled":
		return l.Empty.CancelledMessage
	default:
		return l.Empty.OngoingMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "ongoing":
		return "info"
	case "complete":
		return "success"
	case "cancelled":
		return "warning"
	default:
		return "default"
	}
}

func buildBulkActions(common pyeza.CommonLabels, status string, routes centymo.SalesRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "ongoing":
		actions = append(actions, types.BulkAction{
			Key:             "complete",
			Label:           "Mark Complete",
			Icon:            "icon-check-circle",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    "Mark Complete",
			ConfirmMessage:  "Are you sure you want to mark {{count}} sale(s) as complete?",
			ExtraParamsJSON: `{"target_status":"complete"}`,
		})
	case "complete", "cancelled":
		actions = append(actions, types.BulkAction{
			Key:             "reactivate",
			Label:           "Reactivate",
			Icon:            "icon-play",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    "Reactivate",
			ConfirmMessage:  "Are you sure you want to reactivate {{count}} sale(s)?",
			ExtraParamsJSON: `{"target_status":"ongoing"}`,
		})
	}

	return actions
}
