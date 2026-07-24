package subscription_group_product_plan

import (
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
)

// StringEqFilter builds a single field == value FilterRequest. Shared by
// list/, action/ and detail/ so each sub-package doesn't hand-roll its own
// copy (mirrors the read-contract's batched LIST_IN discipline, espyna.md
// §1b).
func StringEqFilter(field, value string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: field,
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{Value: value, Operator: commonpb.StringOperator_STRING_EQUALS},
		},
	}}}
}

// StringInFilter builds a single field IN (values) FilterRequest.
func StringInFilter(field string, values []string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: field,
		FilterType: &commonpb.TypedFilter_ListFilter{
			ListFilter: &commonpb.ListFilter{Values: values, Operator: commonpb.ListOperator_LIST_IN},
		},
	}}}
}

// BoolEqFilter builds a single field == value boolean FilterRequest (used to
// scope active-only reads).
func BoolEqFilter(field string, value bool) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field:      field,
		FilterType: &commonpb.TypedFilter_BooleanFilter{BooleanFilter: &commonpb.BooleanFilter{Value: value}},
	}}}
}

// AndFilters merges filter requests into one (all conditions AND-ed — mirrors
// how the postgres adapter treats a Filters slice).
func AndFilters(reqs ...*commonpb.FilterRequest) *commonpb.FilterRequest {
	out := &commonpb.FilterRequest{}
	for _, r := range reqs {
		if r == nil {
			continue
		}
		out.Filters = append(out.Filters, r.Filters...)
	}
	return out
}
