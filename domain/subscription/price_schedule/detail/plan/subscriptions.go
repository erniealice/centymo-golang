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

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
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
	Active     bool
}

// subscriptionPageLimit is the adapter's per-page cap (GetSubscriptionListPageData).
const subscriptionPageLimit = 100

// fetchSubscriptionsForPricePlan pages through the reverse-index use case for
// ONE status. The adapter's CTE is a hard `s.active = $x` equality — "both
// statuses" cannot be expressed in a single call (same constraint as
// shared.InactiveFilter) — and its default page size is 20, so the full set
// needs the page loop (a level-AY carries 100+ enrollments).
func fetchSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string, activeOnly bool) []*subscriptionpb.Subscription {
	var out []*subscriptionpb.Subscription
	for page := int32(1); ; page++ {
		resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
			PricePlanId: pricePlanID,
			ActiveOnly:  &activeOnly,
			Pagination: &commonpb.PaginationRequest{
				Limit:  subscriptionPageLimit,
				Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
			},
		})
		if err != nil {
			log.Printf("Failed to load subscriptions for price plan %s (active=%t): %v", pricePlanID, activeOnly, err)
			return out
		}
		batch := resp.GetSubscriptionList()
		out = append(out, batch...)
		if len(batch) < subscriptionPageLimit {
			return out
		}
	}
}

// countSubscriptionsForPricePlan returns the count of subscriptions (both
// statuses — a previous academic year's enrollments are all inactive rows)
// referencing the given PricePlan, via TotalItems on a limit-1 probe per
// status so the badge never materializes the roster. Returns 0 (no badge)
// when the dep is unwired.
func countSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) int {
	if deps.ListSubscriptionsByPricePlan == nil {
		return 0
	}
	total := 0
	for _, activeOnly := range []bool{true, false} {
		ao := activeOnly
		resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
			PricePlanId: pricePlanID,
			ActiveOnly:  &ao,
			Pagination: &commonpb.PaginationRequest{
				Limit:  1,
				Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: 1}},
			},
		})
		if err != nil {
			log.Printf("Failed to count subscriptions for price plan %s (active=%t): %v", pricePlanID, ao, err)
			continue
		}
		total += int(resp.GetPagination().GetTotalItems())
	}
	return total
}

// loadSubscriptionsForPricePlan fetches the subscriptions referencing the
// given PricePlan — BOTH statuses, all pages — and shapes them into
// SubscriptionRow values for the tab table. Hydration of Client + PricePlan +
// Plan is provided by the espyna use case (single CTE-based JOIN) — the view
// layer does not chain N+1 lookups.
func loadSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) []SubscriptionRow {
	if deps.ListSubscriptionsByPricePlan == nil {
		return nil
	}
	seen := map[string]bool{}
	var subs []*subscriptionpb.Subscription
	for _, activeOnly := range []bool{true, false} {
		for _, s := range fetchSubscriptionsForPricePlan(ctx, deps, pricePlanID, activeOnly) {
			if s == nil || seen[s.GetId()] {
				continue
			}
			seen[s.GetId()] = true
			subs = append(subs, s)
		}
	}

	tz := types.LocationFromContext(ctx)
	rows := make([]SubscriptionRow, 0, len(subs))
	for _, s := range subs {
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
			Active:     s.GetActive(),
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
		// The tab now lists BOTH statuses (historical academic years are
		// inactive rows) — the status chip is what tells them apart.
		{Key: "status", Label: deps.PlanLabels.Columns.Status, WidthClass: "col-2xl"},
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

		// Badge Value renders verbatim — use the lyngua status labels, not the
		// raw status key.
		statusLabel, statusVariant := deps.CommonLabels.Status.Active, "success"
		if !r.Active {
			statusLabel, statusVariant = deps.CommonLabels.Status.Inactive, "warning"
		}

		tableRows = append(tableRows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: r.Name},
				{Type: "text", Value: r.ClientName},
				{Type: "text", Value: r.Plan},
				{Type: "text", Value: r.DateStart},
				{Type: "text", Value: r.DateEnd},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
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
