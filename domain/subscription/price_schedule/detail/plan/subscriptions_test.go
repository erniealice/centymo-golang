package plan

import (
	"context"
	"net/url"
	"strings"
	"testing"

	priceplan "github.com/erniealice/centymo-golang/domain/subscription/price_plan"
	schedule "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

func TestSubscriptionPrimaryAction(t *testing.T) {
	for _, routeCase := range []struct{ name, path, label, slug string }{
		{"default_routes", subscription.DefaultRoutes().AddURL, subscription.DefaultLabels().Buttons.AddSubscription, "subscriptions"},
		{"custom_override", "/test/custom-subscription/add", "Custom Create", "custom-items"},
	} {
		t.Run(routeCase.name, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				perms    *types.UserPermissions
				disabled bool
			}{
				{"allowed", types.NewUserPermissions([]string{"subscription:create"}), false},
				{"denied", types.NewEmptyUserPermissions(), true},
				{"missing", nil, true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					deps := &DetailViewDeps{Routes: schedule.DefaultRoutes(), ScheduleLabels: schedule.DefaultLabels(), PlanLabels: priceplan.DefaultLabels(), CommonLabels: pyeza.CommonLabels{}, SubscriptionAddURL: routeCase.path, SubscriptionAddLabel: routeCase.label,
						ListClientNames: func(context.Context) map[string]string { return map[string]string{"client-1": "Sample Client"} }}
					deps.ScheduleLabels.Tabs.SubscriptionsSlug = routeCase.slug
					deps.CommonLabels.Errors.MissingPermission = "Missing permission: %s"
					clientID := "client-1"
					pp := &priceplanpb.PricePlan{ClientId: &clientID, BillingCurrency: "PHP"}
					table := buildSubscriptionsTable(view.WithUserPermissions(context.Background(), tc.perms), deps, "schedule-1", "price-1", pp, "Sample Plan", nil)
					action := table.PrimaryAction
					if action == nil || action.Label != routeCase.label || action.Disabled != tc.disabled {
						t.Fatalf("unexpected primary action: %+v", action)
					}
					if tc.disabled && !strings.Contains(action.DisabledTooltip, "subscription:create") {
						t.Fatalf("missing permission tooltip: %q", action.DisabledTooltip)
					}
					u, err := url.Parse(action.ActionURL)
					if err != nil {
						t.Fatal(err)
					}
					if u.Path != routeCase.path {
						t.Fatalf("action path = %q, want %q", u.Path, routeCase.path)
					}
					for key, want := range map[string]string{"price_plan_id": "price-1", "client_id": "client-1", "client_name": "Sample Client", "billing_currency": "PHP", "plan_label": "Sample Plan"} {
						if u.Query().Get(key) != want {
							t.Errorf("%s = %q", key, u.Query().Get(key))
						}
					}
					if !strings.HasSuffix(table.RefreshURL, "/"+routeCase.slug) {
						t.Errorf("refresh = %s", table.RefreshURL)
					}
				})
			}
		})
	}
}
