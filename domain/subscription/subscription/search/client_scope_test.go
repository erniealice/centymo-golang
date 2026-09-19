package search

import (
	"context"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientScopedPlanEffectiveCommencement(t *testing.T) {
	scheduleID, clientID, name := "schedule-1", "client-1", "Custom Plan"
	start, err := time.Parse(time.RFC3339Nano, "2026-09-16T22:56:17.283308+08:00")
	if err != nil {
		t.Fatal(err)
	}
	deps := &Deps{
		ListPricePlans: func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error) {
			return &priceplanpb.ListPricePlansResponse{Data: []*priceplanpb.PricePlan{{Id: "custom-pricing", Name: &name, PriceScheduleId: &scheduleID, Active: true, BillingCurrency: "PHP"}}}, nil
		},
		ListPriceSchedules: func(context.Context, *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error) {
			return &priceschedulepb.ListPriceSchedulesResponse{Data: []*priceschedulepb.PriceSchedule{{Id: scheduleID, ClientId: &clientID, Active: true, DateTimeStart: timestamppb.New(start)}}}, nil
		},
	}
	for _, tc := range []struct {
		name, client, start string
		want                bool
	}{
		{"same day midnight precedes schedule", "client-1", "2026-09-15T16:00:00Z", false},
		{"next day includes matching client", "client-1", "2026-09-16T16:00:00Z", true},
		{"other client excluded by parent schedule", "client-2", "2026-09-16T16:00:00Z", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/search?commencement_only=1&billing_currency=PHP&client_id="+tc.client+"&date_time_start_iso="+tc.start, nil)
			w := httptest.NewRecorder()
			NewSearchPlansAction(deps)(w, req)
			if strings.Contains(w.Body.String(), "custom-pricing") != tc.want {
				t.Fatalf("response %s", w.Body.String())
			}
		})
	}
}
