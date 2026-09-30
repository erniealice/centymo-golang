package allocation_batch

import (
	"context"
	"fmt"
	"sort"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Participant is one subscription that may share a cost: it has an opted-in
// agreement term that covers the cost's whole service period. The list only
// offers choices; the use case re-proves every term when the allocation is saved.
type Participant struct {
	SubscriptionID   string
	ClientID         string
	SubscriptionName string
	ClientName       string
}

// Covers reports whether a term's half-open [effective_from, effective_to) covers
// the service interval [from, to) (ISO dates compare as strings).
func Covers(t *agreementlinetermpb.AgreementLineTerm, from, to string) bool {
	if t == nil || !t.GetActive() || from == "" || to == "" {
		return false
	}
	if t.GetEffectiveFrom() > from {
		return false
	}
	return t.EffectiveTo == nil || t.GetEffectiveTo() == "" || t.GetEffectiveTo() >= to
}

// ReadComponent reads one cost source component.
func (u *UseCases) ReadComponent(ctx context.Context, id string) (*costsourcecomponentpb.CostSourceComponent, error) {
	if u == nil || u.ReadCostSourceComponent == nil {
		return nil, fmt.Errorf("allocation_batch: component read unavailable")
	}
	resp, err := u.ReadCostSourceComponent(ctx, &costsourcecomponentpb.ReadCostSourceComponentRequest{Data: &costsourcecomponentpb.CostSourceComponent{Id: id}})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 || resp.GetData()[0] == nil {
		return nil, fmt.Errorf("allocation_batch: component %s not found", id)
	}
	return resp.GetData()[0], nil
}

// Participants returns the subscriptions with an opted-in agreement term that
// covers [from, to), one entry per subscription, ordered by subscription name.
// It pages through every agreement term (C11) and reads each subscription once
// for its display names.
func (u *UseCases) Participants(ctx context.Context, from, to string) ([]Participant, error) {
	if u == nil || u.ListAgreementLineTerms == nil || u.GetSubscriptionItemPageData == nil {
		return nil, fmt.Errorf("allocation_batch: participant reads unavailable")
	}
	bySub := map[string]*Participant{}
	var order []string
	done := false
	for page := int32(1); page <= maxPages; page++ {
		resp, err := u.ListAgreementLineTerms(ctx, &agreementlinetermpb.ListAgreementLineTermsRequest{
			Sort:       &commonpb.SortRequest{Fields: []*commonpb.SortField{{Field: "id", Direction: commonpb.SortDirection_ASC}}},
			Pagination: pageRequest(page),
		})
		if err != nil {
			return nil, err
		}
		for _, t := range resp.GetData() {
			if !Covers(t, from, to) || t.GetSubscriptionId() == "" || t.GetClientId() == "" {
				continue
			}
			if _, seen := bySub[t.GetSubscriptionId()]; !seen {
				bySub[t.GetSubscriptionId()] = &Participant{SubscriptionID: t.GetSubscriptionId(), ClientID: t.GetClientId()}
				order = append(order, t.GetSubscriptionId())
			}
		}
		if !resp.GetPagination().GetHasNext() {
			done = true
			break
		}
	}
	if !done {
		return nil, errPageBound
	}
	out := make([]Participant, 0, len(order))
	for _, id := range order {
		p := bySub[id]
		sresp, err := u.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{SubscriptionId: id})
		if err != nil {
			return nil, err
		}
		sub := sresp.GetSubscription()
		p.SubscriptionName = sub.GetName()
		if c := sub.GetClient(); c != nil {
			p.ClientName = c.GetName()
			if p.ClientName == "" {
				p.ClientName = joinName(c.GetUser().GetFirstName(), c.GetUser().GetLastName())
			}
		}
		out = append(out, *p)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].SubscriptionName != out[j].SubscriptionName {
			return out[i].SubscriptionName < out[j].SubscriptionName
		}
		return out[i].SubscriptionID < out[j].SubscriptionID
	})
	return out, nil
}

func joinName(a, b string) string {
	switch {
	case a != "" && b != "":
		return a + " " + b
	case a != "":
		return a
	}
	return b
}

// BatchesOf returns every revision of a component's allocation, newest revision first.
func (u *UseCases) BatchesOf(ctx context.Context, componentID string) ([]*allocationbatchpb.AllocationBatch, error) {
	if u == nil || u.GetAllocationBatchListPageData == nil {
		return nil, fmt.Errorf("allocation_batch: batch list unavailable")
	}
	var out []*allocationbatchpb.AllocationBatch
	done := false
	for page := int32(1); page <= maxPages; page++ {
		resp, err := u.GetAllocationBatchListPageData(ctx, &allocationbatchpb.GetAllocationBatchListPageDataRequest{
			Filters:    stringEquals("cost_source_component_id", componentID),
			Pagination: pageRequest(page),
		})
		if err != nil {
			return nil, err
		}
		for _, b := range resp.GetAllocationBatchList() {
			if b != nil && b.GetCostSourceComponentId() == componentID {
				out = append(out, b)
			}
		}
		if !resp.GetPagination().GetHasNext() {
			done = true
			break
		}
	}
	if !done {
		return nil, errPageBound
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].GetRevision() > out[j].GetRevision() })
	return out, nil
}

// SharesOf returns the shares of one batch in sequence order.
func (u *UseCases) SharesOf(ctx context.Context, batchID string) ([]*allocationsharepb.AllocationShare, error) {
	if u == nil || u.GetAllocationShareListPageData == nil {
		return nil, fmt.Errorf("allocation_batch: share list unavailable")
	}
	var out []*allocationsharepb.AllocationShare
	done := false
	for page := int32(1); page <= maxPages; page++ {
		resp, err := u.GetAllocationShareListPageData(ctx, &allocationsharepb.GetAllocationShareListPageDataRequest{
			Filters:    stringEquals("allocation_batch_id", batchID),
			Pagination: pageRequest(page),
		})
		if err != nil {
			return nil, err
		}
		for _, s := range resp.GetAllocationShareList() {
			if s != nil && s.GetAllocationBatchId() == batchID {
				out = append(out, s)
			}
		}
		if !resp.GetPagination().GetHasNext() {
			done = true
			break
		}
	}
	if !done {
		return nil, errPageBound
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].GetSequenceOrder() < out[j].GetSequenceOrder() })
	return out, nil
}
