package form

import (
	"errors"
	"net/url"
	"strconv"

	priceplanform "github.com/erniealice/centymo-golang/domain/subscription/price_plan/form"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// PopulateEscalation prepares the agreement's saved escalation snapshot for
// the edit drawer. New agreements intentionally remain blank so the use case
// copies the selected Price Plan defaults exactly once at creation time.
func PopulateEscalation(data *Data, subscription *subscriptionpb.Subscription) {
	if data == nil || subscription == nil {
		return
	}
	if subscription.EscalationMode != nil {
		data.EscalationMode = subscription.GetEscalationMode().String()
	}
	if subscription.EscalationScope != nil {
		data.EscalationScope = subscription.GetEscalationScope().String()
	}
	if subscription.EscalationRateBps != nil {
		data.EscalationRate = priceplanform.FormatEscalationPercentBPS(subscription.GetEscalationRateBps())
	}
	if subscription.EscalationFirstAfterMonths != nil {
		data.EscalationFirstAfterMonths = strconv.FormatInt(int64(subscription.GetEscalationFirstAfterMonths()), 10)
	}
	if subscription.EscalationEveryMonths != nil {
		data.EscalationEveryMonths = strconv.FormatInt(int64(subscription.GetEscalationEveryMonths()), 10)
	}
}

// ApplyEscalation maps drawer fields onto an agreement. A blank create value
// leaves every pointer nil, which tells espyna to snapshot the Price Plan
// defaults. A blank edit value explicitly clears an existing snapshot.
func ApplyEscalation(subscription *subscriptionpb.Subscription, values url.Values, isEdit bool) error {
	if subscription == nil {
		return errors.New("subscription is required")
	}
	modeRaw := values.Get("escalation_mode")
	if modeRaw == "" && !isEdit {
		return nil
	}

	translated := url.Values{
		"default_escalation_mode":               {modeRaw},
		"default_escalation_scope":              {values.Get("escalation_scope")},
		"default_escalation_rate":               {values.Get("escalation_rate")},
		"default_escalation_first_after_months": {values.Get("escalation_first_after_months")},
		"default_escalation_every_months":       {values.Get("escalation_every_months")},
	}
	temporary := &priceplanpb.PricePlan{}
	if err := priceplanform.ApplyDefaultEscalation(temporary, translated); err != nil {
		return err
	}
	subscription.EscalationMode = temporary.DefaultEscalationMode
	subscription.EscalationScope = temporary.DefaultEscalationScope
	subscription.EscalationRateBps = temporary.DefaultEscalationRateBps
	subscription.EscalationFirstAfterMonths = temporary.DefaultEscalationFirstAfterMonths
	subscription.EscalationEveryMonths = temporary.DefaultEscalationEveryMonths
	return nil
}
