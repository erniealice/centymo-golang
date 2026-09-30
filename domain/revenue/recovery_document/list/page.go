// Package list renders the recovery documents list (Issued · Void).
package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
)

// Deps holds the list view dependencies.
type Deps struct {
	Routes       rd.Routes
	Labels       rd.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *rd.UseCases
}

// PageData is the list page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	StatusTabs      []pyeza.TabItem
	ActiveStatus    string
}

var searchFields = []string{"document_number"}

// NewView renders the full page (or the HTMX content partial).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "list") {
			return view.Forbidden("recovery_document:list")
		}
		if deps.UseCases == nil || deps.UseCases.ListRecoveryDocuments == nil {
			return unavailable(deps)
		}
		status := rd.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, tabs, err := build(ctx, deps, viewCtx, status)
		if err != nil {
			return view.Error(err)
		}
		title := statusTitle(deps.Labels, status)
		return view.OK("recovery-document-list", &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    title,
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "recovery-document-list-content",
			Table:           table,
			StatusTabs:      tabs,
			ActiveStatus:    status,
		})
	})
}

// NewTableView renders only the table card (HTMX refresh / pagination).
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "list") {
			return view.Forbidden("recovery_document:list")
		}
		if deps.UseCases == nil || deps.UseCases.ListRecoveryDocuments == nil {
			return unavailable(deps)
		}
		status := rd.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, _, err := build(ctx, deps, viewCtx, status)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", table)
	})
}

func unavailable(deps *Deps) view.ViewResult {
	r := view.Error(fmt.Errorf("%s", deps.CommonLabels.Errors.General))
	r.StatusCode = http.StatusServiceUnavailable
	return r
}

func statusEnum(status string) recoverydocumentpb.RecoveryDocumentStatus {
	if status == "void" {
		return recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_VOID
	}
	return recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_ISSUED
}

func statusTitle(l rd.Labels, status string) string {
	if status == "void" {
		return l.Page.TitleVoid
	}
	return l.Page.TitleIssued
}

func statusTabLabel(l rd.Labels, status string) string {
	if status == "void" {
		return l.Tabs.Void
	}
	return l.Tabs.Issued
}

func emptyTitle(l rd.Labels, status string) string {
	if status == "void" {
		return l.Empty.VoidTitle
	}
	return l.Empty.IssuedTitle
}

func emptyMessage(l rd.Labels, status string) string {
	if status == "issued" {
		return l.Empty.IssuedMessage
	}
	return ""
}

func columns(l rd.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "document_number", Label: l.Columns.DocumentNumber, NoFilter: true, WidthClass: "col-4xl"},
		{Key: "document_type", Label: l.Columns.DocumentType, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "client", Label: l.Columns.Client, NoSort: true, NoFilter: true},
		{Key: "issue_date", Label: l.Columns.IssueDate, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "due_date", Label: l.Columns.DueDate, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "total_amount", Label: l.Columns.Total, NoFilter: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: l.Columns.Status, NoFilter: true, WidthClass: "col-2xl"},
	}
}

func typeLabel(l rd.Labels, t recoverydocumentpb.RecoveryDocumentType) string {
	if t == recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_CREDIT_NOTE {
		return l.Enums.TypeCreditNote
	}
	return l.Enums.TypeStatement
}

func build(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, status string) (*types.TableConfig, []pyeza.TabItem, error) {
	l := deps.Labels
	cols := columns(l)
	p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(cols), types.FilterableKeys(cols), "date_created", "desc")
	if err != nil {
		return nil, nil, err
	}
	lp := espynahttp.ToListParams(p, searchFields)
	st := statusEnum(status)
	resp, err := deps.UseCases.ListRecoveryDocuments(ctx, &recoverydocumentpb.ListRecoveryDocumentsRequest{
		Search: lp.Search, Filters: lp.Filters, Sort: lp.Sort, Pagination: lp.Pagination, Status: &st,
	})
	if err != nil {
		log.Printf("recovery_document list: ListRecoveryDocuments(%s): %v", status, err)
		return nil, nil, err
	}

	clients := deps.UseCases.ClientNames(ctx)
	rows := buildRows(deps, resp.GetData(), status, clients)
	types.ApplyColumnStyles(cols, rows)

	refresh := route.ResolveURL(deps.Routes.TableURL, "status", status)
	totalRows := int(resp.GetPagination().GetTotalItems())
	pageSize := p.PageSize
	if pageSize < 1 {
		pageSize = 25
	}
	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "offset",
		CurrentPage:   p.Page,
		PageSize:      pageSize,
		TotalRows:     totalRows,
		TotalPages:    int(math.Ceil(float64(totalRows) / float64(pageSize))),
		SearchQuery:   p.Search,
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refresh,
	}
	sp.BuildDisplay()

	table := &types.TableConfig{
		ID:                   "recovery-documents-table",
		RefreshURL:           refresh,
		Columns:              cols,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "date_created",
		DefaultSortDirection: "desc",
		ServerPagination:     sp,
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   emptyTitle(l, status),
			Message: emptyMessage(l, status),
		},
	}
	types.ApplyTableSettings(table)

	tabs := make([]pyeza.TabItem, 0, len(rd.Statuses))
	for _, s := range rd.Statuses {
		tabs = append(tabs, pyeza.TabItem{
			Key:   s,
			Label: statusTabLabel(l, s),
			Href:  route.ResolveURL(deps.Routes.ListURL, "status", s),
		})
	}
	return table, tabs, nil
}

func buildRows(deps *Deps, docs []*recoverydocumentpb.RecoveryDocument, status string, clients map[string]string) []types.TableRow {
	l := deps.Labels
	rows := make([]types.TableRow, 0, len(docs))
	for _, d := range docs {
		id := d.GetId()
		detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		client := clients[d.GetClientId()]
		if client == "" {
			client = d.GetClientId()
		}
		stLabel, stVariant := l.Enums.StatusIssued, "success"
		if d.GetStatus() == recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_VOID {
			stLabel, stVariant = l.Enums.StatusVoid, "muted"
		}
		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "link", Value: d.GetDocumentNumber(), Href: detailURL},
				{Type: "text", Value: typeLabel(l, d.GetDocumentType())},
				{Type: "text", Value: client},
				{Type: "text", Value: d.GetIssueDate()},
				{Type: "text", Value: d.GetDueDate()},
				{Type: "text", Value: types.FormatMoney(d.GetTotalAmount(), d.GetCurrency())},
				{Type: "badge", Value: stLabel, Variant: stVariant},
			},
			DataAttrs: map[string]string{
				"testid": "recovery-document-row-" + id,
				"status": status,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: deps.CommonLabels.Actions.View, Action: "view", Href: detailURL, TestID: "recovery-document-view-" + id},
			},
		})
	}
	return rows
}
