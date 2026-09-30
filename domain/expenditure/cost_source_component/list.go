package cost_source_component

import (
	"context"
	"fmt"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
)

// PageSize is the page size of every paged read in this package.
const PageSize = 100

// MaxPages bounds every page walker in this package (R5 m7; recovery_document/ports.go:74 precedent):
// past it the walk fails closed instead of looping or returning a silently truncated list.
const MaxPages = 1000

var errPageBound = fmt.Errorf("cost_source_component: page bound (%d pages) exceeded", MaxPages)

// ListForExpenditure returns every recoverable cost line of one expenditure,
// paging through ListCostSourceComponents until a short page (C11: no
// silent cap; the response carries no pagination block). Rows of another expenditure are dropped defensively.
func (u *UseCases) ListForExpenditure(ctx context.Context, expenditureID string) ([]*costsourcecomponentpb.CostSourceComponent, error) {
	if u == nil || u.ListCostSourceComponents == nil {
		return nil, fmt.Errorf("cost_source_component: list unavailable")
	}
	var out []*costsourcecomponentpb.CostSourceComponent
	for page := int32(1); page <= MaxPages; page++ {
		resp, err := u.ListCostSourceComponents(ctx, &costsourcecomponentpb.ListCostSourceComponentsRequest{
			Filters: &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
				Field: "expenditure_id",
				FilterType: &commonpb.TypedFilter_StringFilter{StringFilter: &commonpb.StringFilter{
					Value: expenditureID, Operator: commonpb.StringOperator_STRING_EQUALS, CaseSensitive: true,
				}},
			}}},
			Sort:       &commonpb.SortRequest{Fields: []*commonpb.SortField{{Field: "date_created", Direction: commonpb.SortDirection_ASC}, {Field: "id", Direction: commonpb.SortDirection_ASC}}},
			Pagination: &commonpb.PaginationRequest{Limit: PageSize, Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}}},
		})
		if err != nil {
			return nil, err
		}
		for _, c := range resp.GetData() {
			if c != nil && c.GetExpenditureId() == expenditureID {
				out = append(out, c)
			}
		}
		if len(resp.GetData()) < PageSize {
			return out, nil
		}
	}
	return nil, errPageBound
}
