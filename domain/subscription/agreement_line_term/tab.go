package agreement_line_term

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/erniealice/pyeza-golang/types"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// TableID is the DOM id of the charge-terms table card.
const TableID = "agreement-line-terms-table"

// TabData is the template data of "subscription-tab-charge-terms".
type TabData struct {
	Labels Labels
	Table  *types.TableConfig
	Failed bool // the read failed: the tab shows no rows and no empty-state claim
}

// OriginLabel returns the display label of a term origin.
func OriginLabel(l Labels, o agreementlinetermpb.AgreementLineTermOrigin) string {
	switch o {
	case agreementlinetermpb.AgreementLineTermOrigin_AGREEMENT_LINE_TERM_ORIGIN_COPIED:
		return l.Enums.OriginCopied
	case agreementlinetermpb.AgreementLineTermOrigin_AGREEMENT_LINE_TERM_ORIGIN_NEGOTIATED:
		return l.Enums.OriginNegotiated
	}
	return "—"
}

// FormatMarkup renders basis points as a percentage with two decimals (250 -> "2.50%").
func FormatMarkup(bps int32) string {
	return fmt.Sprintf("%d.%02d%%", bps/100, bps%100)
}

// maxPages bounds the page walk (R5 m7): past it Terms fails closed.
const maxPages = 1000

// Terms lists every agreement term of one subscription (paged, C11), ordered by
// effective_from then id.
func (u *UseCases) Terms(ctx context.Context, subscriptionID string) ([]*agreementlinetermpb.AgreementLineTerm, error) {
	if u == nil || u.ListAgreementLineTerms == nil {
		return nil, fmt.Errorf("agreement_line_term: list unavailable")
	}
	var out []*agreementlinetermpb.AgreementLineTerm
	done := false
	for page := int32(1); page <= maxPages; page++ {
		resp, err := u.ListAgreementLineTerms(ctx, &agreementlinetermpb.ListAgreementLineTermsRequest{
			SubscriptionId: &subscriptionID,
			Pagination:     &commonpb.PaginationRequest{Limit: 100, Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}}},
		})
		if err != nil {
			return nil, err
		}
		for _, t := range resp.GetData() {
			if t != nil && t.GetSubscriptionId() == subscriptionID {
				out = append(out, t)
			}
		}
		if !resp.GetPagination().GetHasNext() {
			done = true
			break
		}
	}
	if !done {
		return nil, fmt.Errorf("agreement_line_term: page bound (%d pages) exceeded", maxPages)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].GetEffectiveFrom() != out[j].GetEffectiveFrom() {
			return out[i].GetEffectiveFrom() < out[j].GetEffectiveFrom()
		}
		return out[i].GetId() < out[j].GetId()
	})
	return out, nil
}

// lineNames maps product_price_plan id -> a display name for the price plan's lines
// (empty map when the closure is unbound or the read fails: ids are shown instead).
func (u *UseCases) lineNames(ctx context.Context, pricePlanID string) map[string]string {
	names := map[string]string{}
	if u == nil || u.ListProductPricePlans == nil || pricePlanID == "" {
		return names
	}
	resp, err := u.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{
		Filters: &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
			Field: "price_plan_id",
			FilterType: &commonpb.TypedFilter_StringFilter{StringFilter: &commonpb.StringFilter{
				Value: pricePlanID, Operator: commonpb.StringOperator_STRING_EQUALS, CaseSensitive: true,
			}},
		}}},
	})
	if err != nil {
		log.Printf("agreement_line_term tab: list product price plans of %s: %v", pricePlanID, err)
		return names
	}
	for _, p := range resp.GetData() {
		if p == nil {
			continue
		}
		name := p.GetProductPlan().GetProduct().GetName()
		if name == "" {
			name = p.GetProductPlan().GetName()
		}
		if name != "" {
			names[p.GetId()] = name
		}
	}
	return names
}

type policyRef struct{ name, version string }

// policyOf resolves a pinned version id to (policy name, "v<n>"); cached per call.
func (u *UseCases) policyOf(ctx context.Context, l Labels, cache map[string]policyRef, versionID string) policyRef {
	if ref, ok := cache[versionID]; ok {
		return ref
	}
	ref := policyRef{name: versionID}
	if u != nil && u.ReadChargePolicyVersion != nil {
		vr, err := u.ReadChargePolicyVersion(ctx, &versionpb.ReadChargePolicyVersionRequest{Data: &versionpb.ChargePolicyVersion{Id: versionID}})
		if err != nil {
			log.Printf("agreement_line_term tab: read charge policy version %s: %v", versionID, err)
		} else if len(vr.GetData()) > 0 && vr.GetData()[0] != nil {
			v := vr.GetData()[0]
			ref.version = strings.Replace(l.Detail.VersionLabel, "{0}", fmt.Sprint(v.GetVersionNumber()), 1)
			if u.ReadChargePolicy != nil {
				pr, perr := u.ReadChargePolicy(ctx, &policypb.ReadChargePolicyRequest{Data: &policypb.ChargePolicy{Id: v.GetChargePolicyId()}})
				if perr != nil {
					log.Printf("agreement_line_term tab: read charge policy %s: %v", v.GetChargePolicyId(), perr)
				} else if len(pr.GetData()) > 0 && pr.GetData()[0] != nil && pr.GetData()[0].GetName() != "" {
					ref.name = pr.GetData()[0].GetName()
				}
			}
		}
	}
	cache[versionID] = ref
	return ref
}

// Load builds the read-only Charge terms tab of one subscription.
func (u *UseCases) Load(ctx context.Context, l Labels, tableLabels types.TableLabels, subscriptionID, pricePlanID string) *TabData {
	terms, err := u.Terms(ctx, subscriptionID)
	if err != nil {
		log.Printf("agreement_line_term tab: terms of subscription %s: %v", subscriptionID, err)
		return &TabData{Labels: l, Failed: true, Table: Build(nil, l, tableLabels, nil, nil)}
	}
	lines := u.lineNames(ctx, pricePlanID)
	cache := map[string]policyRef{}
	return &TabData{Labels: l, Table: Build(terms, l, tableLabels, lines, func(versionID string) (string, string) {
		ref := u.policyOf(ctx, l, cache, versionID)
		return ref.name, ref.version
	})}
}

// Build renders the terms as a table (columns: line, policy, version, markup,
// from, to, origin). policyOf may be nil (ids are shown).
func Build(terms []*agreementlinetermpb.AgreementLineTerm, l Labels, tableLabels types.TableLabels, lines map[string]string, policyOf func(versionID string) (name, version string)) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "line", Label: l.Columns.Line, NoSort: true},
		{Key: "policy", Label: l.Columns.Policy, NoSort: true},
		{Key: "version", Label: l.Columns.Version, NoSort: true, WidthClass: "col-md"},
		{Key: "markup", Label: l.Columns.Markup, NoSort: true, WidthClass: "col-xl"},
		{Key: "effective_from", Label: l.Columns.EffectiveFrom, NoSort: true, WidthClass: "col-3xl"},
		{Key: "effective_to", Label: l.Columns.EffectiveTo, NoSort: true, WidthClass: "col-3xl"},
		{Key: "origin", Label: l.Columns.Origin, NoSort: true, WidthClass: "col-xl"},
	}
	rows := make([]types.TableRow, 0, len(terms))
	for _, t := range terms {
		line := lines[t.GetProductPricePlanId()]
		if line == "" {
			line = t.GetProductPricePlanId()
		}
		policy, version := t.GetChargePolicyVersionId(), "—"
		if policyOf != nil {
			if n, v := policyOf(t.GetChargePolicyVersionId()); n != "" {
				policy = n
				if v != "" {
					version = v
				}
			}
		}
		markup := "—"
		if t.MarkupBps != nil {
			markup = FormatMarkup(t.GetMarkupBps())
		}
		to := l.Detail.OpenEnded
		if t.EffectiveTo != nil && t.GetEffectiveTo() != "" {
			to = t.GetEffectiveTo()
		}
		rows = append(rows, types.TableRow{
			ID: t.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: line},
				{Type: "text", Value: policy},
				{Type: "text", Value: version},
				{Type: "text", Value: markup},
				{Type: "text", Value: t.GetEffectiveFrom()},
				{Type: "text", Value: to},
				{Type: "text", Value: OriginLabel(l, t.GetOrigin())},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)
	return &types.TableConfig{
		ID: TableID, Columns: columns, Rows: rows, Labels: tableLabels,
		EmptyState: types.TableEmptyState{Title: l.Empty.Title, Message: l.Empty.Message},
	}
}
