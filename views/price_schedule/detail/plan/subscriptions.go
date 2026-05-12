package plan

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// SubscriptionRow is one row in the schedule-scoped plan detail
// "Engagements" / "Subscriptions" tab table. Names + dates are pre-formatted
// for the tier's display TZ; the view-layer table builder consumes this struct
// directly into pyeza TableRow cells.
type SubscriptionRow struct {
	ID         string
	Name       string
	ClientID   string
	ClientName string
	Plan       string
	DateStart  string
	DateEnd    string
}

// countSubscriptionsForPricePlan returns the count of active subscriptions
// referencing the given PricePlan. The use case is the same one the tab body
// renders against — so the badge count and row count cannot drift. Returns 0
// (no badge) when the dep is unwired.
func countSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) int {
	if deps.ListSubscriptionsByPricePlan == nil {
		return 0
	}
	activeOnly := true
	resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
		PricePlanId: pricePlanID,
		ActiveOnly:  &activeOnly,
	})
	if err != nil {
		log.Printf("Failed to count subscriptions for price plan %s: %v", pricePlanID, err)
		return 0
	}
	return len(resp.GetSubscriptionList())
}

// loadSubscriptionsForPricePlan fetches active subscriptions for the given
// PricePlan and shapes them into SubscriptionRow values for the tab table.
// Hydration of Client + PricePlan + Plan is provided by the espyna use case
// (single CTE-based JOIN) — the view layer does not chain N+1 lookups.
func loadSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) []SubscriptionRow {
	if deps.ListSubscriptionsByPricePlan == nil {
		return nil
	}
	activeOnly := true
	resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
		PricePlanId: pricePlanID,
		ActiveOnly:  &activeOnly,
	})
	if err != nil {
		log.Printf("Failed to load subscriptions for price plan %s: %v", pricePlanID, err)
		return nil
	}

	tz := types.LocationFromContext(ctx)
	rows := make([]SubscriptionRow, 0, len(resp.GetSubscriptionList()))
	for _, s := range resp.GetSubscriptionList() {
		if s == nil {
			continue
		}
		clientName := ""
		clientID := s.GetClientId()
		if c := s.GetClient(); c != nil {
			clientName = c.GetName()
			if clientName == "" {
				if u := c.GetUser(); u != nil {
					first := u.GetFirstName()
					last := u.GetLastName()
					if first != "" || last != "" {
						clientName = strings.TrimSpace(first + " " + last)
					}
					if clientName == "" {
						clientName = u.GetEmailAddress()
					}
				}
			}
			if clientID == "" {
				clientID = c.GetId()
			}
		}

		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			if p := pp.GetPlan(); p != nil {
				planName = p.GetName()
			}
			if planName == "" {
				planName = pp.GetName()
			}
		}

		rows = append(rows, SubscriptionRow{
			ID:         s.GetId(),
			Name:       s.GetName(),
			ClientID:   clientID,
			ClientName: clientName,
			Plan:       planName,
			DateStart:  types.FormatTimestampInTZ(s.GetDateTimeStart(), tz, types.DateTimeReadable),
			DateEnd:    types.FormatTimestampInTZ(s.GetDateTimeEnd(), tz, types.DateTimeReadable),
		})
	}
	return rows
}

// buildSubscriptionsTable assembles the TableConfig for the schedule-scoped
// plan detail's "Engagements"/"Subscriptions" tab. Columns:
// Name → Client → Plan → Start Date → End Date. The View action targets the
// nested engagement URL when configured so the breadcrumb chains
// rate-card → plan → engagement.
func buildSubscriptionsTable(ctx context.Context, deps *DetailViewDeps, sid, ppid string, pp *priceplanpb.PricePlan, planLabel string, rows []SubscriptionRow) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	subLabels := deps.PlanLabels.Detail.Subscriptions

	columns := []types.TableColumn{
		{Key: "name", Label: subLabels.ColumnName},
		{Key: "client", Label: subLabels.ColumnClient},
		{Key: "plan", Label: subLabels.ColumnPlan},
		{Key: "start_date", Label: subLabels.ColumnStartDate, WidthClass: "col-3xl"},
		{Key: "end_date", Label: subLabels.ColumnEndDate, WidthClass: "col-3xl"},
	}

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		viewURL := ""
		if deps.PlanSubscriptionDetailURL != "" {
			viewURL = route.ResolveURL(deps.PlanSubscriptionDetailURL, "id", sid, "ppid", ppid, "eid", r.ID)
		} else if deps.SubscriptionDetailURL != "" {
			viewURL = route.ResolveURL(deps.SubscriptionDetailURL, "id", r.ID)
		}

		actions := []types.TableAction{}
		if viewURL != "" {
			actions = append(actions, types.TableAction{Type: "view", Label: deps.CommonLabels.Actions.View, Action: "view", Href: viewURL})
		}
		if perms.Can("subscription", "update") && deps.SubscriptionEditURL != "" {
			editURL := route.ResolveURL(deps.SubscriptionEditURL, "id", r.ID)
			actions = append(actions, types.TableAction{Type: "edit", Label: deps.CommonLabels.Actions.Edit, Action: "edit", URL: editURL, DrawerTitle: r.Name})
		}
		if perms.Can("subscription", "delete") && deps.SubscriptionDeleteURL != "" {
			actions = append(actions, types.TableAction{
				Type:           "delete",
				Label:          deps.CommonLabels.Actions.Delete,
				Action:         "delete",
				URL:            deps.SubscriptionDeleteURL,
				ItemName:       r.Name,
				ConfirmTitle:   subLabels.ConfirmDeleteTitle,
				ConfirmMessage: fmt.Sprintf(subLabels.ConfirmDeleteMessage, r.Name),
			})
		}

		tableRows = append(tableRows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: r.Name},
				{Type: "text", Value: r.ClientName},
				{Type: "text", Value: r.Plan},
				{Type: "text", Value: r.DateStart},
				{Type: "text", Value: r.DateEnd},
			},
			DataAttrs: map[string]string{
				"name":   r.Name,
				"client": r.ClientName,
				"plan":   r.Plan,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	refreshURL := route.ResolveURL(deps.Routes.PlanTabActionURL, "id", sid, "ppid", ppid, "tab", deps.ScheduleLabels.Tabs.ResolveTabSlug("subscriptions"))

	// Mirror of the client-detail "Add Engagement" CTA, but with the price
	// plan as the locked context. The drawer's GET handler reads these query
	// params (see views/subscription/action/add.go) and pre-fills + locks
	// the Plan picker. When the price plan is client-scoped, we also pass
	// client_id and billing_currency so the Customer field locks too.
	var primaryAction *types.PrimaryAction
	if deps.SubscriptionAddURL != "" && perms.Can("subscription", "create") && pp != nil {
		q := url.Values{}
		q.Set("price_plan_id", ppid)
		if planLabel != "" {
			q.Set("plan_label", planLabel)
		}
		if cid := pp.GetClientId(); cid != "" {
			q.Set("client_id", cid)
		}
		if cur := pp.GetBillingCurrency(); cur != "" {
			q.Set("billing_currency", cur)
		}
		actionURL := deps.SubscriptionAddURL + "?" + q.Encode()

		label := deps.CommonLabels.Buttons.Add
		if label == "" {
			label = "Add"
		}
		// Compose "Add <Subscription>" using whatever the tier label calls
		// engagements (subscriptionsLabel resolution lives in the page builder).
		if subscriptionsLabel := deps.ScheduleLabels.Tabs.Subscriptions; subscriptionsLabel != "" {
			label = label + " " + subscriptionsLabel
		} else if deps.PlanLabels.Tabs.Subscriptions != "" {
			label = label + " " + deps.PlanLabels.Tabs.Subscriptions
		}
		primaryAction = &types.PrimaryAction{
			Label:     label,
			Icon:      "icon-plus",
			ActionURL: actionURL,
		}
	}

	tc := &types.TableConfig{
		ID:                   "subscriptions-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 tableRows,
		Labels:               deps.TableLabels,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   subLabels.EmptyTitle,
			Message: subLabels.EmptyMessage,
		},
		PrimaryAction: primaryAction,
	}
	types.ApplyTableSettings(tc)
	return tc
}
