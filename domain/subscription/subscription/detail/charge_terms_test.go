package detail

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	alt "github.com/erniealice/centymo-golang/domain/subscription/agreement_line_term"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

func keysOf(deps *DetailViewDeps, jobs bool) []string {
	var out []string
	for _, it := range buildTabItems(subscription.DefaultLabels(), "s1", subscription.DefaultRoutes(), jobs, chargeTermsEnabled(deps)) {
		out = append(out, it.Key)
	}
	return out
}

func has(keys []string, k string) bool {
	for _, x := range keys {
		if x == k {
			return true
		}
	}
	return false
}

// The Charge terms tab exists only where the app wired it (opt-in): apps that do not
// opt in keep exactly their tab set.
func TestChargeTermsTabIsOptIn(t *testing.T) {
	off := keysOf(&DetailViewDeps{}, false)
	if has(off, "charge-terms") {
		t.Fatalf("unwired app offers the tab: %v", off)
	}
	wired := &DetailViewDeps{ChargeTerms: &alt.UseCases{ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
		return &agreementlinetermpb.ListAgreementLineTermsResponse{}, nil
	}}}
	on := keysOf(wired, false)
	if !has(on, "charge-terms") || len(on) != len(off)+1 {
		t.Fatalf("wired tabs = %v, unwired = %v", on, off)
	}
	// only ListAgreementLineTerms gates the tab (name reads are optional)
	if chargeTermsEnabled(&DetailViewDeps{ChargeTerms: &alt.UseCases{}}) {
		t.Fatal("a use-case set without the list closure must not enable the tab")
	}
	if !knownTabs["charge-terms"] || resolveTab(subscription.DefaultRoutes(), "charge-terms") != "charge-terms" {
		t.Fatal("charge-terms must resolve as a known tab")
	}
}

// The tab action renders the read-only terms of the requested subscription.
func TestTabActionServesChargeTerms(t *testing.T) {
	var sawSub string
	deps := &DetailViewDeps{
		Routes: subscription.DefaultRoutes(), Labels: subscription.DefaultLabels(), TableLabels: types.TableLabels{},
		ReadSubscription: func(_ context.Context, r *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error) {
			return &subscriptionpb.ReadSubscriptionResponse{Data: []*subscriptionpb.Subscription{{Id: r.GetData().GetId(), PricePlanId: "pp1"}}}, nil
		},
		ChargeTermLabels: alt.DefaultLabels(),
		ChargeTerms: &alt.UseCases{ListAgreementLineTerms: func(_ context.Context, r *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
			sawSub = r.GetSubscriptionId()
			return &agreementlinetermpb.ListAgreementLineTermsResponse{Data: []*agreementlinetermpb.AgreementLineTerm{{Id: "t1", SubscriptionId: "s1", EffectiveFrom: "2026-01-01"}}}, nil
		}},
	}
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.SetPathValue("id", "s1")
	r.SetPathValue("tab", "charge-terms")
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"subscription:read"}))
	res := NewTabAction(deps).Handle(ctx, &view.ViewContext{Request: r})
	if res.Template != "subscription-tab-charge-terms" {
		t.Fatalf("template = %q (%d %v)", res.Template, res.StatusCode, res.Error)
	}
	pd := res.Data.(*PageData)
	if sawSub != "s1" || pd.ChargeTerms == nil || len(pd.ChargeTerms.Table.Rows) != 1 {
		t.Fatalf("sub=%q data=%+v", sawSub, pd.ChargeTerms)
	}
	// without subscription:read the tab is refused
	r2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	r2.SetPathValue("id", "s1")
	r2.SetPathValue("tab", "charge-terms")
	res = NewTabAction(deps).Handle(view.WithUserPermissions(context.Background(), types.NewUserPermissions(nil)), &view.ViewContext{Request: r2})
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("unauthorised status = %d", res.StatusCode)
	}
}
