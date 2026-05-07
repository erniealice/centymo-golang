package action

import (
	"context"
	"math/rand"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_subscription/form"
	pyeza "github.com/erniealice/pyeza-golang"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Deps holds dependencies for supplier_subscription action handlers.
type Deps struct {
	Routes       centymo.SupplierSubscriptionRoutes
	Labels       centymo.SupplierSubscriptionLabels
	CommonLabels pyeza.CommonLabels

	CreateSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.CreateSupplierSubscriptionRequest) (*suppliersubscriptionpb.CreateSupplierSubscriptionResponse, error)
	ReadSupplierSubscription            func(ctx context.Context, req *suppliersubscriptionpb.ReadSupplierSubscriptionRequest) (*suppliersubscriptionpb.ReadSupplierSubscriptionResponse, error)
	UpdateSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.UpdateSupplierSubscriptionRequest) (*suppliersubscriptionpb.UpdateSupplierSubscriptionResponse, error)
	DeleteSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.DeleteSupplierSubscriptionRequest) (*suppliersubscriptionpb.DeleteSupplierSubscriptionResponse, error)
	GetSupplierSubscriptionItemPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)

	// SetSupplierSubscriptionActive performs a raw DB update for active toggling.
	// Required so proto3's bool=false is not silently omitted.
	SetSupplierSubscriptionActive func(ctx context.Context, id string, active bool) error
}

// generateCode returns a random 7-character uppercase alphanumeric code
// using visually-unambiguous chars (no O, I, 0, 1).
func generateCode() string {
	const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// splitTimestampForInputs splits a *timestamppb.Timestamp into date, time, iso strings
// for the two-row date+time input grid.
func splitTimestampForInputs(ts *timestamppb.Timestamp, tz *time.Location) (date, t, iso string) {
	if ts == nil || !ts.IsValid() {
		return "", "", ""
	}
	moment := ts.AsTime().In(tz)
	return moment.Format(pyezatypes.DateInputLayout), moment.Format(pyezatypes.TimeInputLayout), moment.Format(time.RFC3339)
}

// parseFormDateTime parses date+time+ISO form values into a *timestamppb.Timestamp.
// The hidden ISO field (set by JS) wins when present. Falls back to date+time-in-tz.
// Empty all → nil. isEnd controls default time (23:59:59 for end, 00:00:00 for start).
func parseFormDateTime(date, t, iso string, tz *time.Location, isEnd bool) *timestamppb.Timestamp {
	if iso != "" {
		if parsed, err := time.Parse(time.RFC3339, iso); err == nil {
			return timestamppb.New(parsed.UTC())
		}
	}
	if date == "" {
		return nil
	}
	if t == "" {
		if isEnd {
			t = "23:59:59"
		} else {
			t = "00:00:00"
		}
	} else if len(t) == 5 {
		t = t + ":00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+t, tz)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed.UTC())
}

// buildFormLabels converts centymo.SupplierSubscriptionLabels into form.Labels.
func buildFormLabels(l centymo.SupplierSubscriptionLabels) form.Labels {
	return form.Labels{
		SectionIdentification: l.Form.SectionIdentification,
		SectionRelationships:  l.Form.SectionRelationships,
		SectionConfiguration:  l.Form.SectionConfiguration,
		SectionSchedule:       l.Form.SectionSchedule,
		SectionNotes:          l.Form.SectionNotes,
		Name:                  l.Form.Name,
		NamePlaceholder:       l.Form.NamePlaceholder,
		Code:                  l.Form.Code,
		CodePlaceholder:       l.Form.CodePlaceholder,
		CostPlan:              l.Form.CostPlan,
		CostPlanPlaceholder:   l.Form.CostPlanPlaceholder,
		CostPlanSearch:        l.Form.CostPlanSearch,
		CostPlanNoResults:     l.Form.CostPlanNoResults,
		Supplier:              l.Form.Supplier,
		SupplierPlaceholder:   l.Form.SupplierPlaceholder,
		SupplierSearch:        l.Form.SupplierSearch,
		SupplierNoResults:     l.Form.SupplierNoResults,
		AutoRenew:             l.Form.AutoRenew,
		Active:                l.Form.Active,
		StartDate:             l.Form.StartDate,
		StartTime:             l.Form.StartTime,
		EndDate:               l.Form.EndDate,
		EndTime:               l.Form.EndTime,
		TimePlaceholder:       l.Form.TimePlaceholder,

		Notes:                 l.Form.Notes,
		NotesPlaceholder:      l.Form.NotesPlaceholder,
	}
}

// strPtr returns a pointer to a string value.
func strPtr(s string) *string {
	return &s
}
