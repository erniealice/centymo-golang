package cost_source_component

import (
	"context"
	"fmt"
	"testing"

	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
)

// C11: the list pages until a short page and never silently caps at one page.
func TestListForExpenditurePagesUntilShortPage(t *testing.T) {
	calls := 0
	u := &UseCases{ListCostSourceComponents: func(_ context.Context, req *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error) {
		calls++
		page := req.GetPagination().GetOffset().GetPage()
		n := PageSize
		if page == 3 {
			n = 7
		}
		var rows []*costsourcecomponentpb.CostSourceComponent
		for i := 0; i < n; i++ {
			rows = append(rows, &costsourcecomponentpb.CostSourceComponent{Id: fmt.Sprintf("c-%d-%d", page, i), ExpenditureId: "e1"})
		}
		rows = append(rows[:0:0], rows...)
		if page == 1 {
			rows[0].ExpenditureId = "other" // a foreign row is dropped defensively
		}
		return &costsourcecomponentpb.ListCostSourceComponentsResponse{Data: rows}, nil
	}}
	got, err := u.ListForExpenditure(context.Background(), "e1")
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 || len(got) != 2*PageSize+7-1 {
		t.Fatalf("calls=%d rows=%d", calls, len(got))
	}
	if _, err := (&UseCases{}).ListForExpenditure(context.Background(), "e1"); err == nil {
		t.Fatal("unwired list must fail closed")
	}
}

// R5 m7: a source that never returns a short page fails closed at MaxPages.
func TestListForExpenditureFailsClosedAtPageBound(t *testing.T) {
	calls := 0
	u := &UseCases{ListCostSourceComponents: func(context.Context, *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error) {
		calls++
		rows := make([]*costsourcecomponentpb.CostSourceComponent, PageSize)
		for i := range rows {
			rows[i] = &costsourcecomponentpb.CostSourceComponent{Id: fmt.Sprintf("c-%d-%d", calls, i), ExpenditureId: "e1"}
		}
		return &costsourcecomponentpb.ListCostSourceComponentsResponse{Data: rows}, nil
	}}
	got, err := u.ListForExpenditure(context.Background(), "e1")
	if err == nil || got != nil || calls != MaxPages {
		t.Fatalf("got=%d err=%v calls=%d", len(got), err, calls)
	}
}
