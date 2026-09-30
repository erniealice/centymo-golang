package detail

import (
	"context"

	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// chargeTermsEnabled reports whether the opt-in Charge terms tab is wired for this app.
func chargeTermsEnabled(deps *DetailViewDeps) bool {
	return deps.ChargeTerms != nil && deps.ChargeTerms.ListAgreementLineTerms != nil
}

// chargeTermsTabLabel is the tab label; its English default lives in subscription.DefaultLabels (general
// tier key subscription.tabs.charge_terms).
func chargeTermsTabLabel(l subscription.Labels) string {
	return l.Tabs.ChargeTerms
}

// applyChargeTermsTabData loads the read-only agreement terms of the subscription.
// A no-op when the tab is not wired (the tab item is not offered either).
func applyChargeTermsTabData(ctx context.Context, deps *DetailViewDeps, pageData *PageData, id string, sub *subscriptionpb.Subscription) {
	if !chargeTermsEnabled(deps) {
		return
	}
	pageData.ChargeTerms = deps.ChargeTerms.Load(ctx, deps.ChargeTermLabels, deps.TableLabels, id, sub.GetPricePlanId())
}
