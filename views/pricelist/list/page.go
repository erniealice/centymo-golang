package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"

	"github.com/erniealice/centymo-golang"
)

// Deps holds view dependencies.
type Deps struct {
	Routes         centymo.PriceListRoutes
	ListPriceLists func(ctx context.Context, req *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error)
	GetInUseIDs    func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL     string
	Labels         centymo.PriceListLabels
	CommonLabels   pyeza.CommonLabels
	TableLabels    types.TableLabels
}

// PageData holds the data for the price list list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the price list list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.ListPriceLists(ctx, &pricelistpb.ListPriceListsRequest{})
		if err != nil {
			log.Printf("Failed to list price lists: %v", err)
			return view.Error(fmt.Errorf("failed to load price lists: %w", err))
		}

		var inUseIDs map[string]bool
		if deps.GetInUseIDs != nil {
			var itemIDs []string
			for _, item := range resp.GetData() {
				itemIDs = append(itemIDs, item.GetId())
			}
			inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
		}

		l := deps.Labels
		columns := priceListColumns(l)
		rows := buildTableRows(resp.GetData(), status, l, deps.Routes, inUseIDs)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:              "delete",
				Label:            l.Bulk.Delete,
				Icon:             "icon-trash-2",
				Variant:          "danger",
				Endpoint:         deps.Routes.BulkDeleteURL,
				ConfirmTitle:     l.Bulk.Delete,
				ConfirmMessage:   "Are you sure you want to delete {{count}} price list(s)? This action cannot be undone.",
				RequiresDataAttr: "deletable",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "price-lists-table",
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
			DefaultSortColumn:    "name",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   statusEmptyTitle(l, status),
				Message: statusEmptyMessage(l, status),
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddPriceList,
				ActionURL: deps.Routes.AddURL,
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
				ActiveSubNav:   "price-lists-" + status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "pricelist-list-content",
			Table:           tableConfig,
		}

		return view.OK("pricelist-list", pageData)
	})
}

func priceListColumns(l centymo.PriceListLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "date_start", Label: l.Columns.DateStart, Sortable: true, Width: "150px"},
		{Key: "date_end", Label: l.Columns.DateEnd, Sortable: true, Width: "150px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(priceLists []*pricelistpb.PriceList, status string, l centymo.PriceListLabels, routes centymo.PriceListRoutes, inUseIDs map[string]bool) []types.TableRow {
	rows := []types.TableRow{}
	for _, pl := range priceLists {
		active := pl.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := pl.GetId()
		name := pl.GetName()
		dateStart := pl.GetDateStartString()
		dateEnd := pl.GetDateEndString()
		if dateEnd == "" {
			dateEnd = "—"
		}
		isInUse := inUseIDs[id]

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if isInUse {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = "Cannot delete: price list has products or is used in sales"
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: dateStart},
				{Type: "text", Value: dateEnd},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit},
				deleteAction,
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.PriceListLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.PriceListLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.PriceListLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l centymo.PriceListLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
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
