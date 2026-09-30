package agreement_line_term

import (
	"context"
	"errors"
	"testing"

	"github.com/erniealice/pyeza-golang/types"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

func TestFormatMarkupAndOrigin(t *testing.T) {
	for bps, want := range map[int32]string{0: "0.00%", 250: "2.50%", 1: "0.01%", 12345: "123.45%"} {
		if got := FormatMarkup(bps); got != want {
			t.Errorf("FormatMarkup(%d) = %q, want %q", bps, got, want)
		}
	}
	l := DefaultLabels()
	if OriginLabel(l, agreementlinetermpb.AgreementLineTermOrigin_AGREEMENT_LINE_TERM_ORIGIN_NEGOTIATED) != l.Enums.OriginNegotiated || OriginLabel(l, 0) != "—" {
		t.Error("origin labels")
	}
}

// Terms pages every page (C11), keeps only the subscription's rows and orders by effective_from then id.
func TestTermsPagesAndOrders(t *testing.T) {
	pages := [][]*agreementlinetermpb.AgreementLineTerm{
		{{Id: "b", SubscriptionId: "s1", EffectiveFrom: "2026-06-01"}, {Id: "z", SubscriptionId: "other", EffectiveFrom: "2026-01-01"}},
		{{Id: "a", SubscriptionId: "s1", EffectiveFrom: "2026-06-01"}, {Id: "c", SubscriptionId: "s1", EffectiveFrom: "2026-01-01"}},
	}
	var sawSub string
	u := &UseCases{ListAgreementLineTerms: func(_ context.Context, r *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
		sawSub = r.GetSubscriptionId()
		p := int(r.GetPagination().GetOffset().GetPage())
		return &agreementlinetermpb.ListAgreementLineTermsResponse{Data: pages[p-1], Pagination: &commonpb.PaginationResponse{HasNext: p < len(pages)}}, nil
	}}
	got, err := u.Terms(context.Background(), "s1")
	if err != nil {
		t.Fatal(err)
	}
	var ids string
	for _, x := range got {
		ids += x.GetId()
	}
	if sawSub != "s1" || ids != "cab" {
		t.Fatalf("sub=%q order=%q", sawSub, ids)
	}
	if _, err := (&UseCases{}).Terms(context.Background(), "s1"); err == nil {
		t.Fatal("unwired list must fail closed")
	}
}

func TestLoadResolvesNamesAndDegrades(t *testing.T) {
	l := DefaultLabels()
	bps := int32(500)
	terms := []*agreementlinetermpb.AgreementLineTerm{{Id: "t1", SubscriptionId: "s1", ProductPricePlanId: "ppp1", ChargePolicyVersionId: "v1", MarkupBps: &bps, EffectiveFrom: "2026-01-01"}}
	u := &UseCases{
		ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
			return &agreementlinetermpb.ListAgreementLineTermsResponse{Data: terms}, nil
		},
		ListProductPricePlans: func(context.Context, *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error) {
			return &productpriceplanpb.ListProductPricePlansResponse{Data: []*productpriceplanpb.ProductPricePlan{
				{Id: "ppp1", ProductPlan: &productplanpb.ProductPlan{Product: &productpb.Product{Name: "Electricity"}}},
			}}, nil
		},
		ReadChargePolicyVersion: func(context.Context, *versionpb.ReadChargePolicyVersionRequest) (*versionpb.ReadChargePolicyVersionResponse, error) {
			return &versionpb.ReadChargePolicyVersionResponse{Data: []*versionpb.ChargePolicyVersion{{Id: "v1", ChargePolicyId: "p1", VersionNumber: 3}}}, nil
		},
		ReadChargePolicy: func(context.Context, *policypb.ReadChargePolicyRequest) (*policypb.ReadChargePolicyResponse, error) {
			return &policypb.ReadChargePolicyResponse{Data: []*policypb.ChargePolicy{{Id: "p1", Name: "Utility recovery"}}}, nil
		},
	}
	tab := u.Load(context.Background(), l, types.TableLabels{}, "s1", "pp1")
	if tab.Failed || len(tab.Table.Rows) != 1 {
		t.Fatalf("tab = %+v", tab)
	}
	c := tab.Table.Rows[0].Cells
	if c[0].Value != "Electricity" || c[1].Value != "Utility recovery" || c[2].Value != "v3" || c[3].Value != "5.00%" || c[4].Value != "2026-01-01" || c[5].Value != l.Detail.OpenEnded {
		t.Fatalf("cells = %+v", c)
	}

	// name reads unbound: ids are shown, the term rows still render
	bare := &UseCases{ListAgreementLineTerms: u.ListAgreementLineTerms}
	c = bare.Load(context.Background(), l, types.TableLabels{}, "s1", "pp1").Table.Rows[0].Cells
	if c[0].Value != "ppp1" || c[1].Value != "v1" || c[2].Value != "—" {
		t.Fatalf("degraded cells = %+v", c)
	}

	// a failed read is flagged, never an empty-state claim of "no terms"
	failing := &UseCases{ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
		return nil, errors.New("boom")
	}}
	if tab := failing.Load(context.Background(), l, types.TableLabels{}, "s1", "pp1"); !tab.Failed {
		t.Fatal("a failed read must set Failed")
	}
}

// R5 m7: an endless HasNext stream fails closed at the page bound instead of looping.
func TestTermsFailsClosedAtPageBound(t *testing.T) {
	calls := 0
	u := &UseCases{ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
		calls++
		return &agreementlinetermpb.ListAgreementLineTermsResponse{Pagination: &commonpb.PaginationResponse{HasNext: true}}, nil
	}}
	got, err := u.Terms(context.Background(), "s1")
	if err == nil || got != nil || calls != maxPages {
		t.Fatalf("got=%v err=%v calls=%d", got, err, calls)
	}
}
