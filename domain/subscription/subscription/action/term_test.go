package action

import (
	"context"
	"errors"
	"fmt"
	subform "github.com/erniealice/centymo-golang/domain/subscription/subscription/form"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/pyeza-golang/view"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPlanTermCalendarBoundaries(t *testing.T) {
	for _, tc := range []struct {
		start, unit, want string
		n                 int
	}{
		{"2026-09-13", "month", "2027-09-13", 12},
		{"2026-01-31", "month", "2026-02-28", 1},
		{"2024-02-29", "year", "2025-02-28", 1},
		{"2026-09-13", "week", "2026-09-27", 2},
		{"2026-09-13", "day", "2026-09-14", 1},
	} {
		start, _ := time.Parse("2006-01-02", tc.start)
		got, err := addPlanTerm(start, tc.n, tc.unit)
		if err != nil || got.Format("2006-01-02") != tc.want {
			t.Errorf("%+v: got %v, %v", tc, got, err)
		}
	}
}

func termTestDeps(saved **subscriptionpb.Subscription) *Deps {
	n, unit := int32(12), "month"
	return &Deps{Labels: testLabels(), Routes: testRoutes(), CreateOptions: subscription.CreateOptions{CommencementPricing: true},
		ReadPricePlan: func(context.Context, *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error) {
			return &priceplanpb.ReadPricePlanResponse{Data: []*priceplanpb.PricePlan{{Id: "plan-1", Active: true, Name: strPtr("Annual"), PriceScheduleId: strPtr("schedule-1"), BillingCurrency: "PHP", DefaultTermValue: &n, DefaultTermUnit: &unit}}}, nil
		},
		ReadPriceSchedule: func(context.Context, *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error) {
			return &priceschedulepb.ReadPriceScheduleResponse{Data: []*priceschedulepb.PriceSchedule{{Id: "schedule-1", Active: true, DateTimeStart: timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), DateTimeEnd: timestamppb.New(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC))}}}, nil
		},
		ReadClient: func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
			return &clientpb.ReadClientResponse{Data: []*clientpb.Client{{Id: "client-1", Active: true, BillingCurrency: strPtr("PHP")}}}, nil
		},
		CreateSubscription: func(_ context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
			*saved = req.Data
			return &subscriptionpb.CreateSubscriptionResponse{}, nil
		},
	}
}

func TestCreateCommencementAndExpiry(t *testing.T) {
	for _, tc := range []struct {
		name, start, end string
		accepted         bool
	}{
		{"override beyond price schedule", "2026-09-13", "2030-01-01", true},
		{"open ended", "2026-09-13", "", true},
		{"outside schedule", "2027-01-02", "2028-01-01", false},
		{"before schedule", "2025-12-01", "2026-06-01", false},
		{"end before start", "2026-09-13", "2026-09-12", false},
		{"missing start", "", "", false},
		{"malformed expiry", "2026-09-13", "not-a-date", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var saved *subscriptionpb.Subscription
			deps := termTestDeps(&saved)
			req := postForm("/action/subscription/add", url.Values{"client_id": {"client-1"}, "price_plan_id": {"plan-1"}, "date_start_date": {tc.start}, "date_end_date": {tc.end}, "date_time_start_iso": {"2000-01-01T00:00:00Z"}})
			NewAddAction(deps).Handle(ctxWithPerms("subscription:create"), &view.ViewContext{Request: req, BusinessType: "general"})
			if (saved != nil) != tc.accepted {
				t.Fatalf("saved=%v accepted=%v", saved, tc.accepted)
			}
			if saved != nil {
				tz, _ := time.LoadLocation("Asia/Manila")
				if tc.end == "" {
					if saved.DateTimeEnd != nil {
						t.Fatal("open-ended expiry overwritten")
					}
				} else if got := saved.DateTimeEnd.AsTime().In(tz).Format("2006-01-02"); got != tc.end {
					t.Fatalf("override lost: %s", got)
				}
			}
		})
	}
}

func TestSuggestedPlanExpiryInclusive(t *testing.T) {
	var saved *subscriptionpb.Subscription
	tz, _ := time.LoadLocation("Asia/Manila")
	start := timestamppb.New(time.Date(2026, 9, 13, 0, 0, 0, 0, tz))
	end, has, err := suggestedPlanExpiry(context.Background(), termTestDeps(&saved), "plan-1", start, tz)
	if err != nil || !has || end.AsTime().In(tz).Format("2006-01-02") != "2027-09-12" {
		t.Fatalf("end=%v has=%v err=%v", end, has, err)
	}
}

func TestCreateCommencementPricingOptions(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, businessType := range []string{"", "general", "leasing"} {
			t.Run(fmt.Sprintf("enabled=%v/type=%s", enabled, businessType), func(t *testing.T) {
				for _, preview := range []string{"", "term_preview", "pricing_preview"} {
					var saved *subscriptionpb.Subscription
					deps := termTestDeps(&saved)
					deps.CreateOptions.CommencementPricing = enabled
					q := url.Values{"client_id": {"client-1"}, "billing_currency": {"USD"}, "price_plan_id": {"plan-1"}, "date_start_date": {"2026-09-13"}}
					if preview != "" {
						q.Set(preview, "1")
					}
					req := httptest.NewRequest("GET", deps.Routes.AddURL+"?"+q.Encode(), nil)
					result := NewAddAction(deps).Handle(ctxWithPerms("subscription:create"), &view.ViewContext{Request: req, BusinessType: businessType})
					data, ok := result.Data.(*subform.Data)
					if !ok {
						t.Fatalf("GET %s: %+v", preview, result)
					}
					if data.SuggestPlanTerm != enabled {
						t.Fatalf("mode lost: %+v", data)
					}
					wantTemplate := "subscription-drawer-form"
					if enabled && preview == "term_preview" {
						wantTemplate = "subscription-term-end-fields"
						if data.DateEndDate != "2027-09-12" {
							t.Fatalf("suggestion = %s", data.DateEndDate)
						}
					}
					if enabled && preview == "pricing_preview" {
						wantTemplate = "subscription-pricing-fields"
					}
					if result.Template != wantTemplate {
						t.Fatalf("template=%s want=%s", result.Template, wantTemplate)
					}
					u, _ := url.Parse(data.SearchPlanURL)
					if (u.Query().Get("commencement_only") == "1") != enabled {
						t.Fatalf("search mode=%s", data.SearchPlanURL)
					}
					if preview == "" || !enabled {
						wantCurrency := "USD"
						if enabled {
							wantCurrency = "PHP"
						}
						if data.ClientBillingCurrency != wantCurrency {
							t.Fatalf("currency=%s", data.ClientBillingCurrency)
						}
					}
				}
				var saved *subscriptionpb.Subscription
				deps := termTestDeps(&saved)
				deps.CreateOptions.CommencementPricing = enabled
				req := postForm(deps.Routes.AddURL, url.Values{"client_id": {"client-1"}, "price_plan_id": {"plan-1"}, "date_start_date": {"2027-01-02"}, "date_time_start_iso": {"2027-01-02T00:00:00Z"}})
				NewAddAction(deps).Handle(ctxWithPerms("subscription:create"), &view.ViewContext{Request: req, BusinessType: businessType})
				if (saved != nil) == enabled {
					t.Fatalf("outside schedule saved=%v enabled=%v", saved, enabled)
				}
			})
		}
	}
}

func TestCreateCommencementPricingFailures(t *testing.T) {
	for _, name := range []string{"missing_client_reader", "missing_price_reader", "missing_schedule_reader", "scoped_read_denied", "currency_mismatch", "client_mismatch", "permission_denied", "nil_permissions"} {
		t.Run(name, func(t *testing.T) {
			var saved *subscriptionpb.Subscription
			deps := termTestDeps(&saved)
			ctx := ctxWithPerms("subscription:create")
			switch name {
			case "missing_client_reader":
				deps.ReadClient = nil
			case "missing_price_reader":
				deps.ReadPricePlan = nil
			case "missing_schedule_reader":
				deps.ReadPriceSchedule = nil
			case "scoped_read_denied":
				deps.ReadPricePlan = func(context.Context, *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error) {
					return nil, errors.New("denied")
				}
			case "currency_mismatch":
				deps.ReadClient = func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
					return &clientpb.ReadClientResponse{Data: []*clientpb.Client{{Id: "client-1", Active: true, BillingCurrency: strPtr("USD")}}}, nil
				}
			case "client_mismatch":
				deps.ReadClient = func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
					return &clientpb.ReadClientResponse{Data: []*clientpb.Client{{Id: "other-client", Active: true}}}, nil
				}
			case "permission_denied":
				ctx = ctxWithPerms()
			case "nil_permissions":
				ctx = context.Background()
			}
			req := postForm(deps.Routes.AddURL, url.Values{"client_id": {"client-1"}, "price_plan_id": {"plan-1"}, "date_start_date": {"2026-09-13"}})
			result := NewAddAction(deps).Handle(ctx, &view.ViewContext{Request: req})
			if saved != nil || result.StatusCode < 400 {
				t.Fatalf("failure did not deny: saved=%v result=%+v", saved, result)
			}
		})
	}
}
