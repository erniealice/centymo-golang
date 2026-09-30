package detail

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	"github.com/erniealice/centymo-golang/domain/expenditure/expenditure"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
)

func keysOf(deps *DetailViewDeps) []string {
	var out []string
	for _, it := range buildTabItems(expenditure.Labels{}, "e1", expenditure.DefaultRoutes(), recoverableCostsEnabled(deps)) {
		out = append(out, it.Key)
	}
	return out
}

func wiredDeps(rows ...*costsourcecomponentpb.CostSourceComponent) *DetailViewDeps {
	return &DetailViewDeps{
		Routes: expenditure.DefaultRoutes(), CostSourceComponentRoutes: csc.DefaultRoutes(), CostSourceComponentLabels: csc.DefaultLabels(),
		ReadExpenditure: func(_ context.Context, r *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error) {
			return &expenditurepb.ReadExpenditureResponse{Data: []*expenditurepb.Expenditure{{Id: r.GetData().GetId(), TotalAmount: 10000, Currency: "PHP"}}}, nil
		},
		CostSourceComponents: &csc.UseCases{ListCostSourceComponents: func(context.Context, *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error) {
			return &costsourcecomponentpb.ListCostSourceComponentsResponse{Data: rows}, nil
		}},
	}
}

// The Recoverable costs tab exists only where the app wired it (opt-in).
func TestRecoverableCostsTabIsOptIn(t *testing.T) {
	off := keysOf(&DetailViewDeps{})
	for _, k := range off {
		if k == "recoverable-costs" {
			t.Fatalf("unwired app offers the tab: %v", off)
		}
	}
	on := keysOf(wiredDeps())
	if len(on) != len(off)+1 {
		t.Fatalf("wired tabs = %v, unwired = %v", on, off)
	}
	found := false
	for _, k := range on {
		found = found || k == "recoverable-costs"
	}
	if !found {
		t.Fatalf("wired tabs = %v", on)
	}
	if recoverableCostsEnabled(&DetailViewDeps{CostSourceComponents: &csc.UseCases{}}) {
		t.Fatal("a use-case set without the list closure must not enable the tab")
	}
}

func TestTabActionServesRecoverableCosts(t *testing.T) {
	deps := wiredDeps(
		&costsourcecomponentpb.CostSourceComponent{Id: "c1", ExpenditureId: "e1", Amount: 3000, Currency: "PHP"},
		&costsourcecomponentpb.CostSourceComponent{Id: "c2", ExpenditureId: "e1", Amount: 2500, Currency: "PHP"},
	)
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.SetPathValue("id", "e1")
	r.SetPathValue("tab", "recoverable-costs")
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"expenditure:read", "expenditure:update"}))
	res := NewTabAction(deps).Handle(ctx, &view.ViewContext{Request: r})
	if res.Template != "expense-tab-recoverable-costs" {
		t.Fatalf("template = %q (%d %v)", res.Template, res.StatusCode, res.Error)
	}
	rc := res.Data.(*PageData).RecoverableCosts
	if rc == nil || len(rc.Table.Rows) != 2 || rc.ComponentsSum != "55.00" || rc.BillTotal != "100.00" || rc.Difference != "45.00" || !rc.CanAdd {
		t.Fatalf("data = %+v", rc)
	}
	if rc.AddURL != "/action/cost-source-component/add/e1" {
		t.Fatalf("add url = %q", rc.AddURL)
	}
	// read-only user: the add button is disabled, the rows' actions are disabled
	ctx = view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"expenditure:read"}))
	rc = NewTabAction(deps).Handle(ctx, &view.ViewContext{Request: r}).Data.(*PageData).RecoverableCosts
	if rc.CanAdd {
		t.Fatal("add must be disabled without expenditure:update")
	}
	// a failing list is flagged (never a silent empty table)
	deps.CostSourceComponents = &csc.UseCases{ListCostSourceComponents: func(context.Context, *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error) {
		return nil, fmt.Errorf("boom")
	}}
	if rc := NewTabAction(deps).Handle(ctx, &view.ViewContext{Request: r}).Data.(*PageData).RecoverableCosts; !rc.Failed {
		t.Fatal("a failed list must set Failed")
	}
	// an unwired app: the tab body is empty
	if got := NewTabAction(&DetailViewDeps{Routes: deps.Routes, ReadExpenditure: deps.ReadExpenditure}).Handle(ctx, &view.ViewContext{Request: r}).Data.(*PageData).RecoverableCosts; got != nil {
		t.Fatal("unwired tab must carry no data")
	}
}
