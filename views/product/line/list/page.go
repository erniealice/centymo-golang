package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

type ListViewDeps struct {
	Routes       centymo.ProductLineRoutes
	ListLines    func(ctx context.Context, req *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)
	Labels       centymo.ProductLineLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var productLineAllowedSortCols = []string{"date_created", "date_modified", "name", "status"}
var productLineSearchFields = []string{"name", "description"}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParams(viewCtx.Request, productLineAllowedSortCols)
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
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   statusSubNav(deps.Routes.ActiveSubNav, status),
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-layers",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "product-line-list-content",
			Table:           tableConfig,
		}

		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "product-line"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("product-line-list", pageData)
	})
}

func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParams(viewCtx.Request, productLineAllowedSortCols)
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

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	listParams := espynahttp.ToListParams(p, productLineSearchFields)
	resp, err := deps.ListLines(ctx, &linepb.ListLinesRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list lines: %v", err)
		return nil, fmt.Errorf("failed to load lines: %w", err)
	}

	l := deps.Labels
	columns := productLineColumns(l)
	rows := buildTableRows(resp.GetData(), status, l, deps.Routes, perms)
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
			ConfirmMessage:   l.Confirm.BulkDeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)
	tableConfig := &types.TableConfig{
		ID:                   "product-lines-table",
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
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddProductLine,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("line", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func productLineColumns(l centymo.ProductLineLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
		{Key: "date_created", Label: l.Columns.DateCreated, Sortable: true, Filterable: true, FilterType: types.FilterTypeDate},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Filterable: false, Width: "120px"},
	}
}

func buildTableRows(lines []*linepb.Line, status string, l centymo.ProductLineLabels, routes centymo.ProductLineRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, line := range lines {
		recordStatus := "active"
		if !line.GetActive() {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}
		id := line.GetId()
		name := line.GetName()
		description := line.GetDescription()
		if description == "" {
			description = "—"
		}
		dateCreated := line.GetDateCreatedString()

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if !perms.Can("line", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.PermissionDenied
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: dateCreated},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(true),
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("line", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				deleteAction,
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.ProductLineLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.ProductLineLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.ProductLineLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l centymo.ProductLineLabels, status string) string {
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

func statusSubNav(base, status string) string {
	if base == "" {
		return status
	}
	return base + "-" + status
}
