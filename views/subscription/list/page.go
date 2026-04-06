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
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes                      centymo.SubscriptionRoutes
	GetSubscriptionListPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
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

var subscriptionAllowedSortCols = []string{
	"date_created", "date_start", "date_end", "name",
}

var subscriptionSearchFields = []string{"name"}

// NewView creates the subscription list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, subscriptionAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

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
			return view.Error(fmt.Errorf("failed to load subscriptions: %w", err))
		}

		l := deps.Labels
		columns := subscriptionColumns(l)
		rows := buildTableRows(resp.GetSubscriptionList(), status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		refreshURL := route.ResolveURL(deps.Routes.ListURL, "status", status)

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
			DefaultSortColumn:    "customer",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:           l.Buttons.AddSubscription,
				ActionURL:       deps.Routes.AddURL,
				Icon:            "icon-plus",
				Disabled:        !perms.Can("subscription", "create"),
				DisabledTooltip: l.Errors.NoPermission,
			},
			ServerPagination: sp,
		}
		types.ApplyTableSettings(tableConfig)

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

func subscriptionColumns(l centymo.SubscriptionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "customer", Label: l.Columns.Customer, Sortable: true},
		{Key: "plan", Label: l.Columns.Plan, Sortable: true},
		{Key: "start_date", Label: l.Columns.StartDate, Sortable: true, Width: "150px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(subscriptions []*subscriptionpb.Subscription, status string, l centymo.SubscriptionLabels, routes centymo.SubscriptionRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, s := range subscriptions {
		active := s.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := s.GetId()

		// Build customer display name: prefer company_name, fallback to user name
		customer := s.GetName()
		if c := s.GetClient(); c != nil {
			if companyName := c.GetName(); companyName != "" {
				customer = companyName
			} else if u := c.GetUser(); u != nil {
				firstName := u.GetFirstName()
				lastName := u.GetLastName()
				if firstName != "" || lastName != "" {
					customer = firstName + " " + lastName
				}
			}
		}

		// Get plan name from nested price_plan → plan
		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			if p := pp.GetPlan(); p != nil {
				planName = p.GetName()
			}
			if planName == "" {
				planName = pp.GetName()
			}
		}

		startDate := s.GetDateStart()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: customer},
				{Type: "text", Value: planName},
				{Type: "text", Value: startDate},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"customer":   customer,
				"plan":       planName,
				"start_date": startDate,
				"status":     recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("subscription", "update"), DisabledTooltip: l.Errors.NoPermission},
				{Type: "delete", Label: l.Actions.Cancel, Action: "delete", URL: routes.DeleteURL, ItemName: customer, Disabled: !perms.Can("subscription", "delete"), DisabledTooltip: l.Errors.NoPermission},
			},
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
