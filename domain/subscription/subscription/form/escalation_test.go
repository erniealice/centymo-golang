package form

import (
	"net/url"
	"testing"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

func TestApplyEscalationCreateBlankUsesPricePlanDefault(t *testing.T) {
	record := &subscriptionpb.Subscription{}
	if err := ApplyEscalation(record, url.Values{}, false); err != nil {
		t.Fatal(err)
	}
	if record.EscalationMode != nil {
		t.Fatalf("mode = %v; want nil so create snapshots the Price Plan default", record.EscalationMode)
	}
}

func TestApplyEscalationEditBlankClearsSnapshot(t *testing.T) {
	record := &subscriptionpb.Subscription{}
	if err := ApplyEscalation(record, url.Values{}, true); err != nil {
		t.Fatal(err)
	}
	if record.EscalationMode == nil || record.GetEscalationMode() != priceplanpb.EscalationMode_ESCALATION_MODE_UNSPECIFIED {
		t.Fatalf("mode = %v; want explicit unspecified", record.EscalationMode)
	}
}

func TestApplyEscalationFixedPercentage(t *testing.T) {
	record := &subscriptionpb.Subscription{}
	err := ApplyEscalation(record, url.Values{
		"escalation_mode":               {"ESCALATION_MODE_FIXED_PERCENTAGE"},
		"escalation_scope":              {"ESCALATION_SCOPE_WITHIN_AGREEMENT"},
		"escalation_rate":               {"4.25"},
		"escalation_first_after_months": {"12"},
		"escalation_every_months":       {"12"},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if record.GetEscalationRateBps() != 425 || record.GetEscalationFirstAfterMonths() != 12 || record.GetEscalationEveryMonths() != 12 {
		t.Fatalf("unexpected escalation mapping: %+v", record)
	}
}
