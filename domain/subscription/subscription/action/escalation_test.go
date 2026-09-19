package action

import (
	"context"
	"net/url"
	"testing"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/pyeza-golang/view"
)

func TestEscalationAgreementFormMapping(t *testing.T) {
	var saved *subscriptionpb.Subscription
	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		CreateSubscription: func(_ context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
			saved = req.Data
			return &subscriptionpb.CreateSubscriptionResponse{Success: true}, nil
		},
	}
	req := postForm("/subscriptions/add", url.Values{
		"client_id":                     {"client-1"},
		"price_plan_id":                 {"price-plan-1"},
		"escalation_mode":               {"ESCALATION_MODE_FIXED_PERCENTAGE"},
		"escalation_scope":              {"ESCALATION_SCOPE_WITHIN_AGREEMENT"},
		"escalation_rate":               {"4.25"},
		"escalation_first_after_months": {"12"},
		"escalation_every_months":       {"12"},
	})

	result := NewAddAction(deps).Handle(ctxWithPerms("subscription:create"), &view.ViewContext{Request: req})
	if result.StatusCode != 200 {
		t.Fatalf("status = %d, error = %q", result.StatusCode, result.Headers["HX-Error-Message"])
	}
	if saved == nil || saved.GetEscalationMode() != priceplanpb.EscalationMode_ESCALATION_MODE_FIXED_PERCENTAGE || saved.GetEscalationRateBps() != 425 {
		t.Fatalf("unexpected saved escalation: %+v", saved)
	}
}

func TestEscalationAgreementFormMappingRejectsMalformedRate(t *testing.T) {
	called := false
	labels := testLabels()
	labels.Form.EscalationInvalid = "invalid escalation"
	deps := &Deps{
		Routes: testRoutes(),
		Labels: labels,
		CreateSubscription: func(context.Context, *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
			called = true
			return &subscriptionpb.CreateSubscriptionResponse{}, nil
		},
	}
	req := postForm("/subscriptions/add", url.Values{
		"escalation_mode":  {"ESCALATION_MODE_FIXED_PERCENTAGE"},
		"escalation_scope": {"ESCALATION_SCOPE_ON_RENEWAL"},
		"escalation_rate":  {"4.999"},
	})

	result := NewAddAction(deps).Handle(ctxWithPerms("subscription:create"), &view.ViewContext{Request: req})
	if called || result.Headers["HX-Error-Message"] != "invalid escalation" {
		t.Fatalf("called=%v status=%d error=%q", called, result.StatusCode, result.Headers["HX-Error-Message"])
	}
}
