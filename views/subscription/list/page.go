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

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Deps holds view dependencies.
type Deps struct {
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

// NewView creates the subscription list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.GetSubscriptionListPageData(ctx, &subscriptionpb.GetSubscriptionListPageDataRequest{})
		if err != nil {
			log.Printf("Failed to list subscriptions: %v", err)
			return view.Error(fmt.Errorf("failed to load subscriptions: %w", err))
		}

		l := deps.Labels
		columns := subscriptionColumns(l)
		rows := buildTableRows(resp.GetSubscriptionList(), status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "subscriptions-table",
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
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "services",
				ActiveSubNav:   "subscriptions-" + status,
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
		if recordStatus != status {
			continue
		}

		id := s.GetId()

		// Build customer display name from nested client → user
		customer := s.GetName()
		if c := s.GetClient(); c != nil {
			if u := c.GetUser(); u != nil {
				firstName := u.GetFirstName()
				lastName := u.GetLastName()
				if firstName != "" || lastName != "" {
					customer = firstName + " " + lastName
				}
			}
		}

		// Get plan name from nested price plan
		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			planName = pp.GetName()
		}

		startDate := s.GetDateStartString()

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
