package list

import (
	"context"
	"fmt"
	"log"
	"math"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes                      centymo.SubscriptionRoutes
	GetSubscriptionListPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	GetInUseIDs                 func(ctx context.Context, ids []string) (map[string]bool, error)
	Labels                      centymo.SubscriptionLabels
	CommonLabels                pyeza.CommonLabels
	TableLabels                 types.TableLabels
}

// PageData holds the data for the subscription list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// SubscriptionSortSpec is the canonical sort specification for the subscription
// list page. It is exported so that the sort-consistency guard test can register
// it without duplicating the declaration.
var SubscriptionSortSpec = espynahttp.SortSpec{
	AllowedCols: []string{"name", "date_created", "date_start", "date_end", "client"},
	DefaultCol:  "date_created",
	DefaultDir:  "desc",
	// ColMap bridges view-facing column keys to their SQL counterparts.
	// "client" maps to "client_name" which is the alias projected by the enriched CTE.
	ColMap: map[string]string{
		"date_start": "date_time_start",
		"date_end":   "date_time_end",
		"client":     "client_name",
	},
}

// SubscriptionTableDefaults exposes the table config defaults for the sort-consistency guard test.
var SubscriptionTableDefaults = struct {
	DefaultSortColumn    string
	DefaultSortDirection string
}{
	DefaultSortColumn:    SubscriptionSortSpec.DefaultCol,
	DefaultSortDirection: SubscriptionSortSpec.DefaultDir,
}

var subscriptionSearchFields = []string{"name"}

// buildTableConfig fetches subscription data and builds the table configuration.
// Shared by NewView (full page render) and NewTableView (HTMX partial swap target).
func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, subscriptionSearchFields)

	// Inject status filter for server-side pagination
	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "s.active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.GetSubscriptionListPageData(ctx, &subscriptionpb.GetSubscriptionListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list subscriptions: %v", err)
		return nil, fmt.Errorf("failed to load subscriptions: %w", err)
	}

	// Collect IDs and check which are in use (referenced by dependent tables).
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, s := range resp.GetSubscriptionList() {
			itemIDs = append(itemIDs, s.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	columns := subscriptionColumns(l)
	rows := buildTableRows(ctx, resp.GetSubscriptionList(), status, l, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	// data-refresh-url MUST point at the table-only endpoint so HTMX swaps
	// just the table-card partial; pointing at the full /app/list URL would
	// re-render the entire page (including app-shell) and confuse the swap.
	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)
	if deps.Routes.TableURL == "" {
		refreshURL = route.ResolveURL(deps.Routes.ListURL, "status", status)
	}

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

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:              "activate",
			Label:            l.Status.Activate,
			Icon:             "icon-check-circle",
			Variant:          "success",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"active"}`,
			ConfirmTitle:     l.Confirm.BulkActivate,
			ConfirmMessage:   l.Confirm.BulkActivateMessage,
			RequiresDataAttr: "activatable",
		},
		{
			Key:              "deactivate",
			Label:            l.Status.Deactivate,
			Icon:             "icon-x-circle",
			Variant:          "warning",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"inactive"}`,
			ConfirmTitle:     l.Confirm.BulkDeactivate,
			ConfirmMessage:   l.Confirm.BulkDeactivateMessage,
			RequiresDataAttr: "deactivatable",
		},
		{
			Key:              "delete",
			Label:            l.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Confirm.BulkDelete,
			ConfirmMessage:   l.Confirm.BulkDeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	tableConfig := &types.TableConfig{
		ID:                   "subscriptions-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    SubscriptionSortSpec.DefaultCol,
		DefaultSortDirection: SubscriptionSortSpec.DefaultDir,
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		ServerPagination: sp,
		BulkActions:      &bulkCfg,
	}
	// Add button is only meaningful on the active list — new engagements
	// always start active. Mirrors the plan list's behavior at
	// /app/services/list/inactive.
	if status == "active" {
		tableConfig.PrimaryAction = &types.PrimaryAction{
			Label:           l.Buttons.AddSubscription,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("subscription", "create"),
			DisabledTooltip: l.Errors.NoPermission,
		}
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

// NewView creates the subscription list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, SubscriptionSortSpec)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav + "-" + status,
				HeaderTitle:    statusTitle(l, status),
				HeaderSubtitle: statusSubtitle(l, status),
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "subscription-list-content",
			Table:           tableConfig,
		}

		return view.OK("subscription-list", pageData)
	})
}

// NewTableView returns ONLY the table-card partial. Used as the data-refresh-url
// after row/bulk activate/deactivate/delete so HTMX swaps just the table.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, SubscriptionSortSpec)
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

func subscriptionColumns(l centymo.SubscriptionLabels) []types.TableColumn {
	nameLabel := l.Columns.Name
	if nameLabel == "" {
		nameLabel = "Engagement"
	}
	clientLabel := l.Columns.Client
	if clientLabel == "" {
		clientLabel = l.Columns.Customer
	}
	endDateLabel := l.Columns.EndDate
	if endDateLabel == "" {
		endDateLabel = "End Date"
	}
	// Status column omitted on purpose — the list page is already scoped
	// by /list/{status}, so a per-row badge would be redundant.
	// Column Keys match SubscriptionSortSpec.AllowedCols exactly so that the
	// sort parameter sent by the browser is whitelisted at parse time.
	// "plan" is not in AllowedCols (no SQL sort support); NoSort is not set so
	// the column header renders as sortable-looking, but clicking it falls back
	// to the default — acceptable until a plan-sort CASE WHEN branch is added.
	return []types.TableColumn{
		{Key: "name", Label: nameLabel},
		{Key: "client", Label: clientLabel},
		{Key: "plan", Label: l.Columns.Plan, NoSort: true},
		{Key: "date_start", Label: l.Columns.StartDate, WidthClass: "col-4xl"},
		{Key: "date_end", Label: endDateLabel, WidthClass: "col-4xl"},
	}
}

func buildTableRows(ctx context.Context, subscriptions []*subscriptionpb.Subscription, status string, l centymo.SubscriptionLabels, routes centymo.SubscriptionRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	tz := types.LocationFromContext(ctx)
	rows := []types.TableRow{}
	for _, s := range subscriptions {
		active := s.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := s.GetId()
		subName := s.GetName()

		// Client display name: prefer company_name, fallback to representative
		// full name; empty when the join is missing.
		clientName := ""
		if c := s.GetClient(); c != nil {
			if companyName := c.GetName(); companyName != "" {
				clientName = companyName
			} else if u := c.GetUser(); u != nil {
				firstName := u.GetFirstName()
				lastName := u.GetLastName()
				if firstName != "" || lastName != "" {
					clientName = firstName + " " + lastName
				}
			}
		}

		// Plan name from nested price_plan → plan, with PricePlan as fallback.
		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			if p := pp.GetPlan(); p != nil {
				planName = p.GetName()
			}
			if planName == "" {
				planName = pp.GetName()
			}
		}

		startDate := types.FormatTimestampInTZ(s.GetDateTimeStart(), tz, types.DateTimeReadable)
		endDate := types.FormatTimestampInTZ(s.GetDateTimeEnd(), tz, types.DateTimeReadable)

		// Build per-row actions — conditional on status and in-use state.
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		}

		if recordStatus == "active" {
			actions = append(actions, types.TableAction{
				Type:            "edit",
				Label:           l.Actions.Edit,
				Action:          "edit",
				URL:             route.ResolveURL(routes.EditURL, "id", id),
				DrawerTitle:     l.Actions.Edit,
				Disabled:        !perms.Can("subscription", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Actions.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        subName,
				ConfirmTitle:    l.Confirm.Deactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, subName),
				Disabled:        !perms.Can("subscription", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:            "activate",
				Label:           l.Actions.Activate,
				Action:          "activate",
				URL:             routes.SetStatusURL + "?status=active",
				ItemName:        subName,
				ConfirmTitle:    l.Confirm.Activate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ActivateMessage, subName),
				Disabled:        !perms.Can("subscription", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		}

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: subName,
		}
		if inUseIDs[id] {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.InUse
		}
		if !perms.Can("subscription", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.NoPermission
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: subName},
				{Type: "text", Value: clientName},
				{Type: "text", Value: planName},
				{Type: "datetime", Value: startDate},
				{Type: "datetime", Value: endDate},
			},
			DataAttrs: map[string]string{
				"name":          subName,
				"client":        clientName,
				"plan":          planName,
				"start_date":    startDate,
				"end_date":      endDate,
				"status":        recordStatus,
				"deletable":     boolAttr(!inUseIDs[id]),
				"activatable":   boolAttr(recordStatus == "inactive"),
				"deactivatable": boolAttr(recordStatus == "active"),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l centymo.SubscriptionLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l centymo.SubscriptionLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func boolAttr(v bool) string {
	if v {
		return "true"
	}
	return "false"
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
