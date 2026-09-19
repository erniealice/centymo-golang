package search

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCommencementPricingIgnoresSubscriptionExpiry(t *testing.T) {
	id, name := "schedule-1", "Annual"
	schedule := &priceschedulepb.PriceSchedule{Id: id, Active: true, DateTimeStart: timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), DateTimeEnd: timestamppb.New(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC))}
	deps := &Deps{
		ListPricePlans: func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error) {
			return &priceplanpb.ListPricePlansResponse{Data: []*priceplanpb.PricePlan{{Id: "pricing-1", Active: true, Name: &name, PriceScheduleId: &id}}}, nil
		},
		ListPriceSchedules: func(context.Context, *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error) {
			return &priceschedulepb.ListPriceSchedulesResponse{Data: []*priceschedulepb.PriceSchedule{schedule}}, nil
		},
	}
	for _, tc := range []struct {
		name, query  string
		active, want bool
	}{
		{"commencement valid", "commencement_only=1&date_time_start_iso=2026-09-13T00:00:00Z&date_time_end_iso=2030-01-01T00:00:00Z", true, true},
		{"other consumers retain range", "date_time_start_iso=2026-09-13T00:00:00Z&date_time_end_iso=2030-01-01T00:00:00Z", true, false},
		{"outside validity", "commencement_only=1&date_time_start_iso=2027-01-01T00:00:00Z", true, false},
		{"missing commencement", "commencement_only=1", true, false},
		{"inactive price schedule", "commencement_only=1&date_time_start_iso=2026-09-13T00:00:00Z", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			schedule.Active = tc.active
			w := httptest.NewRecorder()
			NewSearchPlansAction(deps)(w, httptest.NewRequest("GET", "/search?"+tc.query, nil))
			if strings.Contains(w.Body.String(), "pricing-1") != tc.want {
				t.Fatalf("response=%s", w.Body.String())
			}
		})
	}
}
