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

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.RevenueRoutes
	ListRevenues func(ctx context.Context, req *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error)
	RefreshURL   string
	Labels       centymo.RevenueLabels
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

		resp, err := deps.ListRevenues(ctx, &revenuepb.ListRevenuesRequest{})
		if err != nil {
			log.Printf("Failed to list sales: %v", err)
			return view.Error(fmt.Errorf("failed to load sales: %w", err))
		}

		l := deps.Labels
		columns := salesColumns(l)
		rows := buildTableRows(resp.GetData(), status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(deps.CommonLabels, l, status, deps.Routes)

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
				DisabledTooltip: l.Errors.PermissionDenied,
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

func salesColumns(l centymo.RevenueLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "reference", Label: l.Columns.Reference, Sortable: true},
		{Key: "customer", Label: l.Columns.Customer, Sortable: true},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "140px"},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(revenues []*revenuepb.Revenue, status string, l centymo.RevenueLabels, routes centymo.RevenueRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range revenues {
		recordStatus := r.GetStatus()
		if recordStatus != status {
			continue
		}

		id := r.GetId()
		refNumber := r.GetReferenceNumber()
		name := r.GetName()
		date := r.GetRevenueDateString()
		amount := formatAmount(r.GetCurrency(), r.GetTotalAmount())

		detailURL := route.ResolveURL(routes.DetailURL, "id", id)
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied},
			// Download invoice action
			{Type: "download", Label: l.Actions.DownloadInvoice, Action: "download", URL: route.ResolveURL(routes.InvoiceDownloadURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Actions.DownloadInvoice, ConfirmMessage: fmt.Sprintf("Download invoice for %s?", refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
			// Send email action
			{Type: "mail", Label: l.Actions.SendEmail, Action: "send-email", URL: route.ResolveURL(routes.SendEmailURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Confirm.SendEmail, ConfirmMessage: fmt.Sprintf(l.Confirm.SendEmailMessage, refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
		}
		if recordStatus == "ongoing" {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Complete, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=complete", ItemName: refNumber,
				ConfirmTitle:    l.Confirm.Complete,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.CompleteMessage, refNumber),
				Disabled:        !perms.Can("invoice", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Reactivate, Action: "activate",
				URL: routes.SetStatusURL + "?status=ongoing", ItemName: refNumber,
				ConfirmTitle:    l.Confirm.Reactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ReactivateMessage, refNumber),
				Disabled:        !perms.Can("invoice", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}
		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: name},
				{Type: "text", Value: date},
				{Type: "text", Value: amount},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"customer":  name,
				"date":      date,
				"amount":    amount,
				"status":    recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func formatAmount(currency string, amount float64) string {
	if currency == "" {
		currency = "PHP"
	}
	return currency + " " + fmt.Sprintf("%.2f", amount)
}

func statusPageTitle(l centymo.RevenueLabels, status string) string {
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

func statusPageCaption(l centymo.RevenueLabels, status string) string {
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

func statusEmptyTitle(l centymo.RevenueLabels, status string) string {
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

func statusEmptyMessage(l centymo.RevenueLabels, status string) string {
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

func buildBulkActions(common pyeza.CommonLabels, l centymo.RevenueLabels, status string, routes centymo.RevenueRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "ongoing":
		actions = append(actions, types.BulkAction{
			Key:             "complete",
			Label:           l.Confirm.BulkComplete,
			Icon:            "icon-check-circle",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Confirm.BulkComplete,
			ConfirmMessage:  l.Confirm.BulkCompleteMessage,
			ExtraParamsJSON: `{"target_status":"complete"}`,
		})
	case "complete", "cancelled":
		actions = append(actions, types.BulkAction{
			Key:             "reactivate",
			Label:           l.Confirm.BulkReactivate,
			Icon:            "icon-play",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Confirm.BulkReactivate,
			ConfirmMessage:  l.Confirm.BulkReactivateMessage,
			ExtraParamsJSON: `{"target_status":"ongoing"}`,
		})
	}

	return actions
}
