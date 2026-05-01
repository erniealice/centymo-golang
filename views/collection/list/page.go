package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes          centymo.CollectionRoutes
	ListCollections func(ctx context.Context, req *collectionpb.ListCollectionsRequest) (*collectionpb.ListCollectionsResponse, error)
	RefreshURL      string
	Labels          centymo.CollectionLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the collection list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var collectionSearchFields = []string{"name", "reference_number", "status"}

// NewView creates the collection list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := collectionColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "date_created", "desc")
		if err != nil {
			return view.Error(err)
		}

		listParams := espynahttp.ToListParams(p, collectionSearchFields)

		// Inject status filter for server-side pagination
		if listParams.Filters == nil {
			listParams.Filters = &commonpb.FilterRequest{}
		}
		listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
			Field: "status",
			FilterType: &commonpb.TypedFilter_StringFilter{
				StringFilter: &commonpb.StringFilter{
					Value:    status,
					Operator: commonpb.StringOperator_STRING_EQUALS,
				},
			},
		})

		resp, err := deps.ListCollections(ctx, &collectionpb.ListCollectionsRequest{
			Search:     listParams.Search,
			Filters:    listParams.Filters,
			Sort:       listParams.Sort,
			Pagination: listParams.Pagination,
		})
		if err != nil {
			log.Printf("Failed to list collections: %v", err)
			return view.Error(fmt.Errorf("failed to load collections: %w", err))
		}

		l := deps.Labels
		rows := buildTableRows(resp.GetData(), status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(l, status, deps.Routes)

		tableConfig := &types.TableConfig{
			ID:                   "collections-table",
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
				Label:           l.Buttons.AddCollection,
				ActionURL:       deps.Routes.AddURL,
				Icon:            "icon-plus",
				Disabled:        !perms.Can("collection", "create"),
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
				ActiveSubNav:   "collections-" + status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "collection-list-content",
			Table:           tableConfig,
		}

		return view.OK("collection-list", pageData)
	})
}

func collectionColumns(l centymo.CollectionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "reference", Label: l.Columns.Reference},
		{Key: "customer", Label: l.Columns.Customer},
		{Key: "amount", Label: l.Columns.Amount, WidthClass: "col-3xl", Align: "right"},
		{Key: "method", Label: l.Columns.Method, WidthClass: "col-3xl"},
		{Key: "date", Label: l.Columns.Date, WidthClass: "col-3xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
	}
}

func buildTableRows(collections []*collectionpb.Collection, status string, l centymo.CollectionLabels, routes centymo.CollectionRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range collections {
		recordStatus := c.GetStatus()

		id := c.GetId()
		refNumber := c.GetReferenceNumber()
		customer := c.GetName()
		currency := c.GetCurrency()
		method := c.GetCollectionMethodId()
		date := c.GetDateCreatedString()

		detailURL := route.ResolveURL(routes.DetailURL, "id", id)
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("collection", "update"), DisabledTooltip: l.Errors.PermissionDenied},
		}

		switch recordStatus {
		case "pending":
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.MarkComplete, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=completed", ItemName: refNumber,
				ConfirmTitle:    l.Confirm.MarkComplete,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.MarkCompleteMessage, refNumber),
				Disabled:        !perms.Can("collection", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
			actions = append(actions, types.TableAction{
				Type: "delete", Label: l.Actions.Delete, Action: "delete",
				URL: routes.DeleteURL, ItemName: refNumber,
				Disabled:        !perms.Can("collection", "delete"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		case "completed":
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Reactivate, Action: "activate",
				URL: routes.SetStatusURL + "?status=pending", ItemName: refNumber,
				ConfirmTitle:    l.Confirm.Reactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ReactivateMessage, refNumber),
				Disabled:        !perms.Can("collection", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		case "failed":
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Reactivate, Action: "activate",
				URL: routes.SetStatusURL + "?status=pending", ItemName: refNumber,
				ConfirmTitle:    l.Confirm.Reactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ReactivateMessage, refNumber),
				Disabled:        !perms.Can("collection", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
			actions = append(actions, types.TableAction{
				Type: "delete", Label: l.Actions.Delete, Action: "delete",
				URL: routes.DeleteURL, ItemName: refNumber,
				Disabled:        !perms.Can("collection", "delete"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}

		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: customer},
				types.MoneyCell(float64(c.GetAmount()), currency, true),
				{Type: "text", Value: method},
				types.DateTimeCell(date, types.DateReadable),
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"customer":  customer,
				"amount":    fmt.Sprintf("%d", c.GetAmount()),
				"method":    method,
				"date":      date,
				"status":    recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusPageTitle(l centymo.CollectionLabels, status string) string {
	switch status {
	case "pending":
		return l.Page.HeadingPending
	case "completed":
		return l.Page.HeadingCompleted
	case "failed":
		return l.Page.HeadingFailed
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.CollectionLabels, status string) string {
	switch status {
	case "pending":
		return l.Page.CaptionPending
	case "completed":
		return l.Page.CaptionCompleted
	case "failed":
		return l.Page.CaptionFailed
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.CollectionLabels, status string) string {
	switch status {
	case "pending":
		return l.Empty.PendingTitle
	case "completed":
		return l.Empty.CompletedTitle
	case "failed":
		return l.Empty.FailedTitle
	default:
		return l.Empty.PendingTitle
	}
}

func statusEmptyMessage(l centymo.CollectionLabels, status string) string {
	switch status {
	case "pending":
		return l.Empty.PendingMessage
	case "completed":
		return l.Empty.CompletedMessage
	case "failed":
		return l.Empty.FailedMessage
	default:
		return l.Empty.PendingMessage
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

func buildBulkActions(l centymo.CollectionLabels, status string, routes centymo.CollectionRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "pending":
		actions = append(actions, types.BulkAction{
			Key:             "complete",
			Label:           l.Confirm.BulkComplete,
			Icon:            "icon-check-circle",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Confirm.BulkComplete,
			ConfirmMessage:  l.Confirm.BulkCompleteMessage,
			ExtraParamsJSON: `{"target_status":"completed"}`,
		})
		actions = append(actions, types.BulkAction{
			Key:            "delete",
			Label:          l.Confirm.BulkDelete,
			Icon:           "icon-trash",
			Variant:        "danger",
			Endpoint:       routes.BulkDeleteURL,
			ConfirmTitle:   l.Confirm.BulkDelete,
			ConfirmMessage: l.Confirm.BulkDeleteMessage,
		})
	case "completed":
		actions = append(actions, types.BulkAction{
			Key:             "reactivate",
			Label:           l.Confirm.BulkReactivate,
			Icon:            "icon-play",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Confirm.BulkReactivate,
			ConfirmMessage:  l.Confirm.BulkReactivateMessage,
			ExtraParamsJSON: `{"target_status":"pending"}`,
		})
	case "failed":
		actions = append(actions, types.BulkAction{
			Key:             "reactivate",
			Label:           l.Confirm.BulkReactivate,
			Icon:            "icon-play",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Confirm.BulkReactivate,
			ConfirmMessage:  l.Confirm.BulkReactivateMessage,
			ExtraParamsJSON: `{"target_status":"pending"}`,
		})
		actions = append(actions, types.BulkAction{
			Key:            "delete",
			Label:          l.Confirm.BulkDelete,
			Icon:           "icon-trash",
			Variant:        "danger",
			Endpoint:       routes.BulkDeleteURL,
			ConfirmTitle:   l.Confirm.BulkDelete,
			ConfirmMessage: l.Confirm.BulkDeleteMessage,
		})
	}

	return actions
}
