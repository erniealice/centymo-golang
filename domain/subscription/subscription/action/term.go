package action

import (
	"context"
	"fmt"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"net/url"
	"strings"
	"time"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// suggestedPlanExpiry reads the selected plan through its scoped use case. The
// expiry is inclusive, matching revenue-run date windows. Calendar months
// clamp to the last valid day; expiry is the end of the preceding date.
func suggestedPlanExpiry(ctx context.Context, deps *Deps, id string, start *timestamppb.Timestamp, tz *time.Location) (*timestamppb.Timestamp, bool, error) {
	if id == "" {
		return nil, false, nil
	}
	if deps.ReadPricePlan == nil {
		return nil, false, fmt.Errorf("price plan reader unavailable")
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
	if err != nil {
		return nil, false, err
	}
	if resp == nil || len(resp.GetData()) == 0 {
		return nil, false, fmt.Errorf("price plan not found")
	}
	p := resp.GetData()[0]
	n := int(p.GetDefaultTermValue())
	if p.DefaultTermValue == nil {
		return nil, false, nil
	}
	if n <= 0 {
		return nil, true, fmt.Errorf("invalid price plan term")
	}
	if start == nil || start.CheckValid() != nil {
		return nil, true, fmt.Errorf("commencement date required")
	}
	local := start.AsTime().In(tz)
	startDate := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, tz)
	end, err := addPlanTerm(startDate, n, p.GetDefaultTermUnit())
	if err != nil {
		return nil, true, err
	}
	stamp := timestamppb.New(end.Add(-time.Nanosecond))
	if err := stamp.CheckValid(); err != nil {
		return nil, true, err
	}
	return stamp, true, nil
}

func addPlanTerm(start time.Time, n int, unit string) (time.Time, error) {
	if n <= 0 || n > 100000 {
		return time.Time{}, fmt.Errorf("invalid price plan term")
	}
	switch unit {
	case "day":
		return start.AddDate(0, 0, n), nil
	case "week":
		return start.AddDate(0, 0, n*7), nil
	case "month", "year":
		months := n
		if unit == "year" {
			months *= 12
		}
		first := time.Date(start.Year(), start.Month(), 1, start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), start.Location()).AddDate(0, months, 0)
		last := time.Date(first.Year(), first.Month()+1, 0, 0, 0, 0, 0, start.Location()).Day()
		day := start.Day()
		if day > last {
			day = last
		}
		return time.Date(first.Year(), first.Month(), day, start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), start.Location()), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported price plan term unit")
	}
}

// commencementSearchURL intentionally omits subscription expiry: price-schedule validity
// governs selection at commencement, not the full duration of the subscription.
func commencementSearchURL(base string, q url.Values, tz *time.Location) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	v := u.Query()
	v.Set("commencement_only", "1")
	if start := parseFormDateTime(q.Get("date_start_date"), q.Get("date_start_time"), "", tz, false); start != nil {
		v.Set("date_time_start_iso", start.AsTime().Format(time.RFC3339))
	}
	for _, k := range []string{"client_id", "billing_currency", "client_name"} {
		v.Set(k, q.Get(k))
	}
	u.RawQuery = v.Encode()
	return u.String()
}

func validateCommencementPricing(ctx context.Context, deps *Deps, id, clientID string, start *timestamppb.Timestamp) error {
	if start == nil || deps.ReadPricePlan == nil || deps.ReadPriceSchedule == nil || deps.ReadClient == nil {
		return fmt.Errorf("commencement pricing validation unavailable")
	}
	pp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
	if err != nil {
		return err
	}
	if pp == nil || len(pp.GetData()) == 0 {
		return fmt.Errorf("pricing not found")
	}
	p := pp.GetData()[0]
	if !p.GetActive() || p.GetPriceScheduleId() == "" || (p.GetClientId() != "" && p.GetClientId() != clientID) {
		return fmt.Errorf("pricing unavailable")
	}
	sr, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: p.GetPriceScheduleId()}})
	if err != nil {
		return err
	}
	if sr == nil || len(sr.GetData()) == 0 {
		return fmt.Errorf("price schedule not found")
	}
	s := sr.GetData()[0]
	if !s.GetActive() || s.GetDateTimeStart() == nil || s.GetDateTimeStart().AsTime().After(start.AsTime()) ||
		(s.GetDateTimeEnd() != nil && start.AsTime().After(s.GetDateTimeEnd().AsTime())) ||
		(s.GetClientId() != "" && s.GetClientId() != clientID) {
		return fmt.Errorf("price schedule invalid at commencement")
	}
	clients, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{Data: &clientpb.Client{Id: clientID}})
	if err != nil {
		return err
	}
	for _, c := range clients.GetData() {
		if c.GetId() == clientID && c.GetActive() {
			if c.GetBillingCurrency() != "" && !strings.EqualFold(c.GetBillingCurrency(), p.GetBillingCurrency()) {
				return fmt.Errorf("billing currency mismatch")
			}
			return nil
		}
	}
	return fmt.Errorf("client unavailable")
}
