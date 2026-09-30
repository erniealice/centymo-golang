package allocation_batch

import (
	"context"
	"testing"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

func term(id, sub, client, from, to string, active bool) *agreementlinetermpb.AgreementLineTerm {
	t := &agreementlinetermpb.AgreementLineTerm{Id: id, SubscriptionId: sub, ClientId: client, EffectiveFrom: from, Active: active}
	if to != "" {
		t.EffectiveTo = &to
	}
	return t
}

// Half-open [from, to) coverage: the term must start on/before the cost period and end on/after it.
func TestCoversWholeServiceInterval(t *testing.T) {
	cases := []struct {
		name string
		t    *agreementlinetermpb.AgreementLineTerm
		want bool
	}{
		{"open-ended covering", term("a", "s", "c", "2026-01-01", "", true), true},
		{"ends exactly at service_to", term("b", "s", "c", "2026-01-01", "2026-10-01", true), true},
		{"ends before service_to", term("c", "s", "c", "2026-01-01", "2026-09-30", true), false},
		{"starts after service_from", term("d", "s", "c", "2026-09-02", "", true), false},
		{"inactive", term("e", "s", "c", "2026-01-01", "", false), false},
	}
	for _, c := range cases {
		if got := Covers(c.t, "2026-09-01", "2026-10-01"); got != c.want {
			t.Errorf("%s: Covers = %v, want %v", c.name, got, c.want)
		}
	}
	if Covers(nil, "a", "b") || Covers(term("x", "s", "c", "2026-01-01", "", true), "", "2026-10-01") {
		t.Error("nil term / blank period must not cover")
	}
}

// Participants pages every term (C11), keeps one entry per subscription, drops
// non-covering terms, and resolves the display names once per subscription.
func TestParticipantsPagesAndFilters(t *testing.T) {
	pages := [][]*agreementlinetermpb.AgreementLineTerm{
		{term("1", "sub-b", "cli-b", "2026-01-01", "", true), term("2", "sub-x", "cli-x", "2026-09-15", "", true)},
		{term("3", "sub-a", "cli-a", "2026-01-01", "2026-12-31", true), term("4", "sub-b", "cli-b", "2026-01-01", "", true)},
	}
	listCalls, readCalls := 0, 0
	u := &UseCases{
		ListAgreementLineTerms: func(_ context.Context, r *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
			listCalls++
			p := int(r.GetPagination().GetOffset().GetPage())
			return &agreementlinetermpb.ListAgreementLineTermsResponse{Data: pages[p-1], Pagination: &commonpb.PaginationResponse{HasNext: p < len(pages)}}, nil
		},
		GetSubscriptionItemPageData: func(_ context.Context, r *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error) {
			readCalls++
			id := r.GetSubscriptionId()
			cn := "Tenant " + id
			return &subscriptionpb.GetSubscriptionItemPageDataResponse{Subscription: &subscriptionpb.Subscription{Id: id, Name: "Lease " + id, Client: &clientpb.Client{Name: &cn}}}, nil
		},
	}
	got, err := u.Participants(context.Background(), "2026-09-01", "2026-10-01")
	if err != nil {
		t.Fatal(err)
	}
	if listCalls != 2 || readCalls != 2 {
		t.Fatalf("listCalls=%d readCalls=%d", listCalls, readCalls)
	}
	if len(got) != 2 || got[0].SubscriptionID != "sub-a" || got[1].SubscriptionID != "sub-b" {
		t.Fatalf("participants = %+v", got)
	}
	if got[0].ClientID != "cli-a" || got[0].SubscriptionName != "Lease sub-a" || got[0].ClientName != "Tenant sub-a" {
		t.Fatalf("participant = %+v", got[0])
	}
	if _, err := (&UseCases{}).Participants(context.Background(), "a", "b"); err == nil {
		t.Fatal("unwired participants must fail closed")
	}
}

// R5 m7: every walker fails closed at the page bound (an endless HasNext stream never loops or truncates).
func TestWalkersFailClosedAtPageBound(t *testing.T) {
	endless := &commonpb.PaginationResponse{HasNext: true}
	u := &UseCases{
		ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
			return &agreementlinetermpb.ListAgreementLineTermsResponse{Pagination: endless}, nil
		},
		GetSubscriptionItemPageData: func(context.Context, *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error) {
			return &subscriptionpb.GetSubscriptionItemPageDataResponse{}, nil
		},
		GetAllocationBatchListPageData: func(context.Context, *allocationbatchpb.GetAllocationBatchListPageDataRequest) (*allocationbatchpb.GetAllocationBatchListPageDataResponse, error) {
			return &allocationbatchpb.GetAllocationBatchListPageDataResponse{Pagination: endless}, nil
		},
		GetAllocationShareListPageData: func(context.Context, *allocationsharepb.GetAllocationShareListPageDataRequest) (*allocationsharepb.GetAllocationShareListPageDataResponse, error) {
			return &allocationsharepb.GetAllocationShareListPageDataResponse{Pagination: endless}, nil
		},
	}
	ctx := context.Background()
	if got, err := u.Participants(ctx, "2026-09-01", "2026-10-01"); err == nil || got != nil {
		t.Errorf("Participants = %v, %v", got, err)
	}
	if got, err := u.BatchesOf(ctx, "c1"); err == nil || got != nil {
		t.Errorf("BatchesOf = %v, %v", got, err)
	}
	if got, err := u.SharesOf(ctx, "b1"); err == nil || got != nil {
		t.Errorf("SharesOf = %v, %v", got, err)
	}
}
