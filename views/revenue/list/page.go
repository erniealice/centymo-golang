package list

import (
	"context"
	"fmt"
	"log"
	"math"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes          centymo.RevenueRoutes
	GetListPageData func(ctx context.Context, req *revenuepb.GetRevenueListPageDataRequest) (*revenuepb.GetRevenueListPageDataResponse, error)
	Labels          centymo.RevenueLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the sales list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var revenueAllowedSortCols = []string{
	"revenue_date_string", "date_created", "date_modified", "total_amount", "status",
}

var revenueSearchFields = []string{"reference_number", "client_name"}

// NewView creates the sales list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "draft"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, revenueAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sale",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sales-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "sale"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("sales-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "draft"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, revenueAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches revenue data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, revenueSearchFields)

	// Inject status filter for server-side pagination
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "rv.status",
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{
				Value:    status,
				Operator: commonpb.StringOperator_STRING_EQUALS,
			},
		},
	})

	resp, err := deps.GetListPageData(ctx, &revenuepb.GetRevenueListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list sales: %v", err)
		return nil, fmt.Errorf("failed to load sales: %w", err)
	}

	l := deps.Labels
	columns := revenueColumns(l)
	rows := buildTableRows(resp.GetRevenueList(), status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	// Check if any revenue in list has a treasury collection (blocks bulk revert)
	hasAnyCollection := false
	for _, r := range resp.GetRevenueList() {
		if r.GetFulfillmentStatus() == "has_collection" {
			hasAnyCollection = true
			break
		}
	}

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(deps.CommonLabels, l, status, deps.Routes, hasAnyCollection)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	// Build ServerPagination
	totalRows := int(resp.GetPagination().GetTotalItems())
	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "offset",
		CurrentPage:   p.Page,
		PageSize:      p.PageSize,
		TotalRows:     totalRows,
		TotalPages:    int(math.Ceil(float64(totalRows) / float64(p.PageSize))),
		SearchQuery:   p.Search,
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refreshURL,
	}
	sp.BuildDisplay()

	tableConfig := &types.TableConfig{
		ID:                   "sales-table",
		RefreshURL:           refreshURL,
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
		DefaultSortColumn:    "revenue_date_string",
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
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func revenueColumns(l centymo.RevenueLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "reference_number", Label: l.Columns.Reference, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "client_name", Label: l.Columns.Customer, Sortable: true, Filterable: true, FilterType: types.FilterTypeString, WidthClass: "col-9xl"},
		{Key: "revenue_date_string", Label: l.Form.Date, Sortable: true, Filterable: true, FilterType: types.FilterTypeDate, WidthClass: "col-3xl"},
		{Key: "total_amount", Label: l.Columns.Amount, Sortable: true, Filterable: true, FilterType: types.FilterTypeMoney, WidthClass: "col-3xl", Align: "right"},
		{Key: "due_date", Label: l.Form.DueDate, Sortable: true, WidthClass: "col-3xl"},
		{Key: "payment_term", Label: l.Form.PaymentTerms, Sortable: false, WidthClass: "col-3xl"},
	}
}

func buildTableRows(revenues []*revenuepb.Revenue, status string, l centymo.RevenueLabels, routes centymo.RevenueRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range revenues {
		recordStatus := r.GetStatus()

		id := r.GetId()
		refNumber := r.GetReferenceNumber()
		name := r.GetName()
		revenueDate := r.GetRevenueDate()
		amount := centymo.FormatCentavoAmount(r.GetTotalAmount(), r.GetCurrency())
		dueDate := r.GetDueDate()
		paymentTermName := ""
		if pt := r.GetPaymentTerm(); pt != nil {
			paymentTermName = pt.GetName()
		}

		detailURL := route.ResolveURL(routes.DetailURL, "id", id)
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
		}
		switch recordStatus {
		case "draft":
			actions = append(actions,
				types.TableAction{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				types.TableAction{
					Type: "check", Label: l.Actions.Complete, Action: "deactivate",
					URL: routes.SetStatusURL + "?status=complete", ItemName: refNumber,
					ConfirmTitle: l.Confirm.Complete, ConfirmMessage: fmt.Sprintf(l.Confirm.CompleteMessage, refNumber),
					Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied,
				},
				types.TableAction{
					Type: "delete", Label: l.Actions.Cancel, Action: "delete",
					URL: routes.SetStatusURL + "?status=cancelled", ItemName: refNumber,
					ConfirmTitle: l.Confirm.Cancel, ConfirmMessage: fmt.Sprintf(l.Confirm.CancelMessage, refNumber),
					Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied,
				},
				types.TableAction{Type: "download", Label: l.Actions.DownloadInvoice, Action: "download", URL: route.ResolveURL(routes.InvoiceDownloadURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Actions.DownloadInvoice, ConfirmMessage: fmt.Sprintf("Download invoice for %s?", refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
				types.TableAction{Type: "mail", Label: l.Actions.SendEmail, Action: "send-email", URL: route.ResolveURL(routes.SendEmailURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Confirm.SendEmail, ConfirmMessage: fmt.Sprintf(l.Confirm.SendEmailMessage, refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
			)
		case "complete":
			hasCollection := r.GetFulfillmentStatus() == "has_collection"
			undoDisabled := !perms.Can("invoice", "update") || hasCollection
			undoTooltip := l.Errors.PermissionDenied
			if hasCollection {
				undoTooltip = l.Errors.HasPaymentsCannotCancel
			}
			actions = append(actions,
				types.TableAction{
					Type: "undo", Label: l.Actions.ReclassifyToDraft, Action: "undo",
					URL: routes.SetStatusURL + "?status=draft", ItemName: refNumber,
					ConfirmTitle: l.Confirm.ReclassifyToDraft, ConfirmMessage: fmt.Sprintf(l.Confirm.ReclassifyToDraftMessage, refNumber),
					Disabled: undoDisabled, DisabledTooltip: undoTooltip,
				},
				types.TableAction{Type: "download", Label: l.Actions.DownloadInvoice, Action: "download", URL: route.ResolveURL(routes.InvoiceDownloadURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Actions.DownloadInvoice, ConfirmMessage: fmt.Sprintf("Download invoice for %s?", refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
				types.TableAction{Type: "mail", Label: l.Actions.SendEmail, Action: "send-email", URL: route.ResolveURL(routes.SendEmailURL, "id", id), ItemName: refNumber, ConfirmTitle: l.Confirm.SendEmail, ConfirmMessage: fmt.Sprintf(l.Confirm.SendEmailMessage, refNumber), Disabled: !perms.Can("invoice", "read"), DisabledTooltip: l.Errors.PermissionDenied},
			)
		case "cancelled":
			// view only — no other actions
		}
		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: name},
				types.DateTimeCell(revenueDate, types.DateReadable),
				{Type: "text", Value: amount},
				types.DateTimeCell(dueDate, types.DateReadable),
				{Type: "text", Value: paymentTermName},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"customer":  name,
				"date":      revenueDate,
				"amount":    amount,
				"undoable":  func() string { if r.GetFulfillmentStatus() == "has_collection" { return "false" }; return "true" }(),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusPageTitle(l centymo.RevenueLabels, status string) string {
	switch status {
	case "draft":
		return l.Page.HeadingDraft
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
	case "draft":
		return l.Page.CaptionDraft
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
	case "draft":
		return l.Empty.DraftTitle
	case "complete":
		return l.Empty.CompleteTitle
	case "cancelled":
		return l.Empty.CancelledTitle
	default:
		return l.Empty.DraftTitle
	}
}

func statusEmptyMessage(l centymo.RevenueLabels, status string) string {
	switch status {
	case "draft":
		return l.Empty.DraftMessage
	case "complete":
		return l.Empty.CompleteMessage
	case "cancelled":
		return l.Empty.CancelledMessage
	default:
		return l.Empty.DraftMessage
	}
}

func buildBulkActions(common pyeza.CommonLabels, l centymo.RevenueLabels, status string, routes centymo.RevenueRoutes, hasAnyCollection bool) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "draft":
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
	case "complete":
		actions = append(actions, types.BulkAction{
			Key:              "revert",
			Label:            l.Confirm.BulkReactivate,
			Icon:             "icon-undo",
			Variant:          "warning",
			Endpoint:         routes.BulkSetStatusURL,
			ConfirmTitle:     l.Confirm.BulkReactivate,
			ConfirmMessage:   l.Confirm.BulkReactivateMessage,
			ExtraParamsJSON:  `{"target_status":"draft"}`,
			RequiresDataAttr: "undoable",
		})
	case "cancelled":
		// view only — no bulk actions
	}

	return actions
}
