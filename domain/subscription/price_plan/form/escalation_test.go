package form

import (
	"net/url"
	"testing"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

func TestEscalationPercentBasisPoints(t *testing.T) {
	tests := []struct {
		input   string
		want    int32
		wantErr bool
	}{
		{input: "5", want: 500},
		{input: "5.5", want: 550},
		{input: "5.25", want: 525},
		{input: "0.01", want: 1},
		{input: "100.00", want: 10000},
		{input: "5.001", wantErr: true},
		{input: "0", wantErr: true},
		{input: "100.01", wantErr: true},
		{input: "5%", wantErr: true},
	}
	for _, tc := range tests {
		got, err := ParseEscalationPercentBPS(tc.input)
		if (err != nil) != tc.wantErr || got != tc.want {
			t.Errorf("ParseEscalationPercentBPS(%q) = %d, %v; want %d, err=%v", tc.input, got, err, tc.want, tc.wantErr)
		}
	}
}

func TestApplyDefaultEscalation(t *testing.T) {
	pp := &priceplanpb.PricePlan{}
	err := ApplyDefaultEscalation(pp, url.Values{
		"default_escalation_mode":               {"ESCALATION_MODE_FIXED_PERCENTAGE"},
		"default_escalation_scope":              {"ESCALATION_SCOPE_WITHIN_AGREEMENT"},
		"default_escalation_rate":               {"5.25"},
		"default_escalation_first_after_months": {"12"},
		"default_escalation_every_months":       {"12"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if pp.GetDefaultEscalationRateBps() != 525 || pp.GetDefaultEscalationEveryMonths() != 12 {
		t.Fatalf("unexpected mapping: %+v", pp)
	}
}
