// Package list renders the billable charges list (Open · Issued · Cancelled).
package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	billablechargepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billable_charge"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
)

// Deps holds the list view dependencies.
type Deps struct {
	Routes       bc.Routes
	Labels       bc.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *bc.UseCases
}

// PageData is the list page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	StatusTabs      []pyeza.TabItem
	ActiveStatus    string
}

var searchFields = []string{"obligation_key"}

// NewView renders the full page (or the HTMX content partial).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("billable_charge", "list") {
			return view.Forbidden("billable_charge:list")
		}
		if deps.UseCases == nil || deps.UseCases.ListBillableCharges == nil {
			return unavailable(deps)
		}
		status := bc.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, tabs, err := build(ctx, deps, viewCtx, status, perms)
		if err != nil {
			return view.Error(err)
		}
		title := statusTitle(deps.Labels, status)
		return view.OK("billable-charge-list", &PageData{
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
			ContentTemplate: "billable-charge-list-content",
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
		if !perms.Can("billable_charge", "list") {
			return view.Forbidden("billable_charge:list")
		}
		if deps.UseCases == nil || deps.UseCases.ListBillableCharges == nil {
			return unavailable(deps)
		}
		status := bc.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, _, err := build(ctx, deps, viewCtx, status, perms)
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

func statusEnum(status string) billablechargepb.BillableChargeStatus {
	switch status {
	case "issued":
		return billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_ISSUED
	case "cancelled":
		return billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_CANCELLED
	default:
		return billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_OPEN
	}
}

func statusTitle(l bc.Labels, status string) string {
	switch status {
	case "issued":
		return l.Page.TitleIssued
	case "cancelled":
		return l.Page.TitleCancelled
	default:
		return l.Page.TitleOpen
	}
}

func statusTabLabel(l bc.Labels, status string) string {
	switch status {
	case "issued":
		return l.Tabs.Issued
	case "cancelled":
		return l.Tabs.Cancelled
	default:
		return l.Tabs.Open
	}
}

func emptyTitle(l bc.Labels, status string) string {
	switch status {
	case "issued":
		return l.Empty.IssuedTitle
	case "cancelled":
		return l.Empty.CancelledTitle
	default:
		return l.Empty.OpenTitle
	}
}

func emptyMessage(l bc.Labels, status string) string {
	if status == "open" {
		return l.Empty.OpenMessage
	}
	return ""
}

func statusLabelVariant(l bc.Labels, st billablechargepb.BillableChargeStatus) (string, string) {
	switch st {
	case billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_ISSUED:
		return l.Enums.StatusIssued, "success"
	case billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_CANCELLED:
		return l.Enums.StatusCancelled, "muted"
	default:
		return l.Enums.StatusOpen, "warning"
	}
}

func kindLabel(l bc.Labels, k billablechargepb.BillableChargeKind) string {
	if k == billablechargepb.BillableChargeKind_BILLABLE_CHARGE_KIND_CORRECTION {
		return l.Enums.ChargeKindCorrection
	}
	return l.Enums.ChargeKindOriginal
}

func columns(l bc.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "client", Label: l.Columns.Client, NoSort: true, NoFilter: true},
		{Key: "subscription", Label: l.Columns.Subscription, NoSort: true, NoFilter: true},
		{Key: "charge_kind", Label: l.Columns.ChargeKind, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "service_from", Label: l.Columns.ServicePeriod, NoFilter: true, WidthClass: "col-4xl"},
		{Key: "amount", Label: l.Columns.Amount, NoFilter: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: l.Columns.Status, NoFilter: true, WidthClass: "col-2xl"},
	}
}

func build(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, status string, perms *types.UserPermissions) (*types.TableConfig, []pyeza.TabItem, error) {
	l := deps.Labels
	cols := columns(l)
	p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(cols), types.FilterableKeys(cols), "date_created", "desc")
	if err != nil {
		return nil, nil, err
	}
	lp := espynahttp.ToListParams(p, searchFields)
	st := statusEnum(status)
	resp, err := deps.UseCases.ListBillableCharges(ctx, &billablechargepb.ListBillableChargesRequest{
		Search: lp.Search, Filters: lp.Filters, Sort: lp.Sort, Pagination: lp.Pagination, Status: &st,
	})
	if err != nil {
		log.Printf("billable_charge list: ListBillableCharges(%s): %v", status, err)
		return nil, nil, err
	}

	clients := deps.UseCases.ClientNames(ctx)
	subs := deps.UseCases.SubscriptionNames(ctx)
	rows := buildRows(deps, resp.GetData(), status, clients, subs, perms)
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
		ID:                   "billable-charges-table",
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
	// Issue is offered on the open tab; disabled (never hidden) without the permission.
	if status == "open" {
		table.PrimaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Issue,
			ActionURL:       deps.Routes.IssueURL,
			SheetTitle:      l.Buttons.Issue,
			Icon:            "icon-file-plus",
			TestID:          "billable-charge-issue",
			Disabled:        !perms.Can("recovery_document", "issue"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "recovery_document:issue"),
		}
	}
	types.ApplyTableSettings(table)

	tabs := make([]pyeza.TabItem, 0, 2)
	for _, s := range []string{"open", "issued"} {
		tabs = append(tabs, pyeza.TabItem{
			Key:   s,
			Label: statusTabLabel(l, s),
			Href:  route.ResolveURL(deps.Routes.ListURL, "status", s),
		})
	}
	if status == "cancelled" {
		tabs = append(tabs, pyeza.TabItem{Key: "cancelled", Label: l.Tabs.Cancelled, Href: route.ResolveURL(deps.Routes.ListURL, "status", "cancelled")})
	}
	return table, tabs, nil
}

func orID(names map[string]string, id string) string {
	if n := names[id]; n != "" {
		return n
	}
	return id
}

func period(c *billablechargepb.BillableCharge) string {
	switch {
	case c.GetServiceFrom() != "" && c.GetServiceTo() != "":
		return c.GetServiceFrom() + " – " + c.GetServiceTo()
	case c.GetServiceFrom() != "":
		return c.GetServiceFrom()
	default:
		return ""
	}
}

func buildRows(deps *Deps, charges []*billablechargepb.BillableCharge, status string, clients, subs map[string]string, perms *types.UserPermissions) []types.TableRow {
	l := deps.Labels
	cl := deps.CommonLabels
	rows := make([]types.TableRow, 0, len(charges))
	for _, c := range charges {
		id := c.GetId()
		stLabel, stVariant := statusLabelVariant(l, c.GetStatus())
		cells := []types.TableCell{
			{Type: "text", Value: orID(clients, c.GetClientId())},
			{Type: "text", Value: orID(subs, c.GetSubscriptionId())},
			{Type: "text", Value: kindLabel(l, c.GetChargeKind())},
			{Type: "text", Value: period(c)},
			{Type: "text", Value: types.FormatMoney(c.GetAmount(), c.GetCurrency())},
			{Type: "badge", Value: stLabel, Variant: stVariant},
		}
		var actions []types.TableAction
		switch {
		case c.GetStatus() == billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_OPEN:
			actions = append(actions, types.TableAction{
				Type: "edit", Label: l.Buttons.Issue, Action: "issue", TestID: "billable-charge-issue-row-" + id,
				HxGet: deps.Routes.IssueURL + "?id=" + id, HxTarget: "#sheetContent", HxSwap: "innerHTML", DrawerTitle: l.Buttons.Issue,
				Disabled:        !perms.Can("recovery_document", "issue"),
				DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "recovery_document:issue"),
			})
		case c.GetStatus() == billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_ISSUED &&
			c.GetChargeKind() == billablechargepb.BillableChargeKind_BILLABLE_CHARGE_KIND_ORIGINAL:
			actions = append(actions, types.TableAction{
				Type: "edit", Label: l.Buttons.Adjust, Action: "adjust", TestID: "billable-charge-adjust-row-" + id,
				HxGet: route.ResolveURL(deps.Routes.AdjustURL, "id", id), HxTarget: "#sheetContent", HxSwap: "innerHTML", DrawerTitle: l.Buttons.Adjust,
				Disabled:        !perms.Can("billable_charge", "adjust"),
				DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "billable_charge:adjust"),
			})
		}
		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"testid": "billable-charge-row-" + id,
				"status": status,
				"kind":   c.GetChargeKind().String(),
			},
			Actions: actions,
		})
	}
	return rows
}
