package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes            centymo.DisbursementRoutes
	ListDisbursements func(ctx context.Context, req *disbursementpb.ListDisbursementsRequest) (*disbursementpb.ListDisbursementsResponse, error)
	RefreshURL        string
	Labels            centymo.DisbursementLabels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
}

// PageData holds the data for the disbursement list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the disbursement list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		resp, err := deps.ListDisbursements(ctx, &disbursementpb.ListDisbursementsRequest{})
		if err != nil {
			log.Printf("Failed to list disbursements: %v", err)
			return view.Error(fmt.Errorf("failed to load disbursements: %w", err))
		}

		// Filter by status
		var filtered []*disbursementpb.Disbursement
		for _, d := range resp.GetData() {
			if d.GetStatus() == status {
				filtered = append(filtered, d)
			}
		}

		l := deps.Labels
		columns := disbursementColumns(l)
		rows := buildTableRows(filtered, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(l, status, deps.Routes)

		tableConfig := &types.TableConfig{
			ID:                   "disbursements-table",
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
				Label:           l.Buttons.AddDisbursement,
				ActionURL:       deps.Routes.AddURL,
				Icon:            "icon-plus",
				Disabled:        !perms.Can("disbursement", "create"),
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
				ActiveNav:      "cash",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-arrow-up-right",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "disbursement-list-content",
			Table:           tableConfig,
		}

		return view.OK("disbursement-list", pageData)
	})
}

func disbursementColumns(l centymo.DisbursementLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "reference", Label: l.Columns.Reference, Sortable: true},
		{Key: "payee", Label: l.Columns.Payee, Sortable: true},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px", Align: "right"},
		{Key: "method", Label: l.Columns.Method, Sortable: true, Width: "140px"},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(disbursements []*disbursementpb.Disbursement, l centymo.DisbursementLabels, routes centymo.DisbursementRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, d := range disbursements {
		id := d.GetId()
		refNumber := d.GetReferenceNumber()
		payee := d.GetName()
		date := d.GetDateCreatedString()
		currency := d.GetCurrency()
		method := d.GetDisbursementMethodId()
		recordStatus := d.GetStatus()

		amountDisplay := centymo.FormatCentavoAmount(d.GetAmount(), currency)

		detailURL := route.ResolveURL(routes.DetailURL, "id", id)
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("disbursement", "update"), DisabledTooltip: l.Errors.PermissionDenied},
			{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: routes.DeleteURL, ItemName: refNumber, Disabled: !perms.Can("disbursement", "delete"), DisabledTooltip: l.Errors.PermissionDenied},
		}

		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: payee},
				{Type: "text", Value: amountDisplay},
				{Type: "text", Value: method},
				types.DateTimeCell(date, types.DateReadable),
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"payee":     payee,
				"amount":    amountDisplay,
				"method":    method,
				"date":      date,
				"status":    recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusPageTitle(l centymo.DisbursementLabels, status string) string {
	switch status {
	case "draft":
		return l.Page.HeadingDraft
	case "pending":
		return l.Page.HeadingPending
	case "approved":
		return l.Page.HeadingApproved
	case "paid":
		return l.Page.HeadingPaid
	case "cancelled":
		return l.Page.HeadingCancelled
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.DisbursementLabels, status string) string {
	switch status {
	case "draft":
		return l.Page.CaptionDraft
	case "pending":
		return l.Page.CaptionPending
	case "approved":
		return l.Page.CaptionApproved
	case "paid":
		return l.Page.CaptionPaid
	case "cancelled":
		return l.Page.CaptionCancelled
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.DisbursementLabels, status string) string {
	switch status {
	case "draft":
		return l.Empty.DraftTitle
	case "pending":
		return l.Empty.PendingTitle
	case "approved":
		return l.Empty.ApprovedTitle
	case "paid":
		return l.Empty.PaidTitle
	case "cancelled":
		return l.Empty.CancelledTitle
	default:
		return l.Empty.PendingTitle
	}
}

func statusEmptyMessage(l centymo.DisbursementLabels, status string) string {
	switch status {
	case "draft":
		return l.Empty.DraftMessage
	case "pending":
		return l.Empty.PendingMessage
	case "approved":
		return l.Empty.ApprovedMessage
	case "paid":
		return l.Empty.PaidMessage
	case "cancelled":
		return l.Empty.CancelledMessage
	default:
		return l.Empty.PendingMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "paid":
		return "success"
	case "cancelled":
		return "danger"
	case "overdue":
		return "danger"
	default:
		return "default"
	}
}

func buildBulkActions(l centymo.DisbursementLabels, status string, routes centymo.DisbursementRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "draft":
		actions = append(actions,
			types.BulkAction{
				Key:             "submit",
				Label:           l.Confirm.BulkSubmit,
				Icon:            "icon-send",
				Variant:         "primary",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkSubmit,
				ConfirmMessage:  l.Confirm.BulkSubmitMessage,
				ExtraParamsJSON: `{"target_status":"pending"}`,
			},
			types.BulkAction{
				Key:            "delete",
				Label:          l.Confirm.BulkDelete,
				Icon:           "icon-trash",
				Variant:        "danger",
				Endpoint:       routes.BulkDeleteURL,
				ConfirmTitle:   l.Confirm.BulkDelete,
				ConfirmMessage: l.Confirm.BulkDeleteMessage,
			},
		)
	case "pending":
		actions = append(actions,
			types.BulkAction{
				Key:             "approve",
				Label:           l.Confirm.BulkApprove,
				Icon:            "icon-check-circle",
				Variant:         "primary",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkApprove,
				ConfirmMessage:  l.Confirm.BulkApproveMessage,
				ExtraParamsJSON: `{"target_status":"approved"}`,
			},
			types.BulkAction{
				Key:             "cancel",
				Label:           l.Confirm.BulkCancel,
				Icon:            "icon-x-circle",
				Variant:         "warning",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkCancel,
				ConfirmMessage:  l.Confirm.BulkCancelMessage,
				ExtraParamsJSON: `{"target_status":"cancelled"}`,
			},
			types.BulkAction{
				Key:            "delete",
				Label:          l.Confirm.BulkDelete,
				Icon:           "icon-trash",
				Variant:        "danger",
				Endpoint:       routes.BulkDeleteURL,
				ConfirmTitle:   l.Confirm.BulkDelete,
				ConfirmMessage: l.Confirm.BulkDeleteMessage,
			},
		)
	case "approved":
		actions = append(actions,
			types.BulkAction{
				Key:             "mark_paid",
				Label:           l.Confirm.BulkMarkPaid,
				Icon:            "icon-check",
				Variant:         "success",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkMarkPaid,
				ConfirmMessage:  l.Confirm.BulkMarkPaidMessage,
				ExtraParamsJSON: `{"target_status":"paid"}`,
			},
			types.BulkAction{
				Key:             "cancel",
				Label:           l.Confirm.BulkCancel,
				Icon:            "icon-x-circle",
				Variant:         "warning",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkCancel,
				ConfirmMessage:  l.Confirm.BulkCancelMessage,
				ExtraParamsJSON: `{"target_status":"cancelled"}`,
			},
		)
	case "overdue":
		actions = append(actions,
			types.BulkAction{
				Key:             "mark_paid",
				Label:           l.Confirm.BulkMarkPaid,
				Icon:            "icon-check",
				Variant:         "success",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkMarkPaid,
				ConfirmMessage:  l.Confirm.BulkMarkPaidMessage,
				ExtraParamsJSON: `{"target_status":"paid"}`,
			},
			types.BulkAction{
				Key:             "cancel",
				Label:           l.Confirm.BulkCancel,
				Icon:            "icon-x-circle",
				Variant:         "warning",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkCancel,
				ConfirmMessage:  l.Confirm.BulkCancelMessage,
				ExtraParamsJSON: `{"target_status":"cancelled"}`,
			},
		)
	case "cancelled":
		actions = append(actions,
			types.BulkAction{
				Key:             "reactivate",
				Label:           l.Confirm.BulkReactivate,
				Icon:            "icon-play",
				Variant:         "primary",
				Endpoint:        routes.BulkSetStatusURL,
				ConfirmTitle:    l.Confirm.BulkReactivate,
				ConfirmMessage:  l.Confirm.BulkReactivateMessage,
				ExtraParamsJSON: `{"target_status":"draft"}`,
			},
			types.BulkAction{
				Key:            "delete",
				Label:          l.Confirm.BulkDelete,
				Icon:           "icon-trash",
				Variant:        "danger",
				Endpoint:       routes.BulkDeleteURL,
				ConfirmTitle:   l.Confirm.BulkDelete,
				ConfirmMessage: l.Confirm.BulkDeleteMessage,
			},
		)
	}

	return actions
}
