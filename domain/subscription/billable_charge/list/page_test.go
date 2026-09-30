package list_test

import (
	"net/http"
	"strings"
	"testing"

	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	bctest "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/internal/bctest"
	bclist "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/list"
)

func TestListDeniedWithoutPermission(t *testing.T) {
	res := bclist.NewView(bctest.ListDeps(bctest.Seed())).Handle(bctest.CtxWith(), bctest.Request(http.MethodGet, "/x", "", "status", "open"))
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", res.StatusCode)
	}
}

func TestListFailsClosedWhenUnwired(t *testing.T) {
	deps := bctest.ListDeps(bctest.Seed())
	deps.UseCases = &bc.UseCases{}
	res := bclist.NewView(deps).Handle(bctest.CtxWith("billable_charge:list"), bctest.Request(http.MethodGet, "/x", "", "status", "open"))
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", res.StatusCode)
	}
}

func TestListOpenAndIssuedRowsAndActions(t *testing.T) {
	deps := bctest.ListDeps(bctest.Seed())
	ctx := bctest.CtxWith("billable_charge:list", "billable_charge:adjust", "recovery_document:issue")

	res := bclist.NewView(deps).Handle(ctx, bctest.Request(http.MethodGet, "/revenue/billable-charges/list/open", "", "status", "open"))
	if res.Template != "billable-charge-list" {
		t.Fatalf("template = %q (%v)", res.Template, res.Error)
	}
	pd := res.Data.(*bclist.PageData)
	if len(pd.Table.Rows) != 2 || pd.Table.Rows[0].ID != "o1" {
		t.Fatalf("open rows = %+v", pd.Table.Rows)
	}
	if pd.Table.PrimaryAction == nil || pd.Table.PrimaryAction.TestID != "billable-charge-issue" || pd.Table.PrimaryAction.Disabled {
		t.Errorf("open tab primary action = %+v", pd.Table.PrimaryAction)
	}
	if len(pd.Table.Rows[0].Actions) != 1 || pd.Table.Rows[0].Actions[0].Action != "issue" {
		t.Errorf("open row actions = %+v", pd.Table.Rows[0].Actions)
	}
	// Row drawers need hx-get: pyeza table-actions emits data-*-url only for the
	// built-in edit/clone/delete/... actions, so a URL-only custom action renders
	// a dead button (W5 B-W5-3).
	if a := pd.Table.Rows[0].Actions[0]; a.HxGet == "" || a.HxTarget != "#sheetContent" {
		t.Errorf("issue row action must open the sheet via hx-get: %+v", a)
	}

	res = bclist.NewView(deps).Handle(ctx, bctest.Request(http.MethodGet, "/x", "", "status", "issued"))
	pd = res.Data.(*bclist.PageData)
	if len(pd.Table.Rows) != 2 {
		t.Fatalf("issued rows = %d", len(pd.Table.Rows))
	}
	// Adjust only on issued ORIGINAL charges; the correction row has no action.
	for _, row := range pd.Table.Rows {
		wantAdjust := row.ID == "i1"
		if got := len(row.Actions) == 1 && row.Actions[0].Action == "adjust"; got != wantAdjust {
			t.Errorf("row %s adjust action = %v, want %v", row.ID, got, wantAdjust)
		}
		if wantAdjust && (row.Actions[0].HxGet == "" || row.Actions[0].HxTarget != "#sheetContent") {
			t.Errorf("adjust row action must open the sheet via hx-get: %+v", row.Actions[0])
		}
	}
	if pd.Table.PrimaryAction != nil {
		t.Error("issued tab must not offer Issue")
	}
}

func TestListDisablesIssueWithoutPermission(t *testing.T) {
	res := bclist.NewView(bctest.ListDeps(bctest.Seed())).Handle(bctest.CtxWith("billable_charge:list"), bctest.Request(http.MethodGet, "/x", "", "status", "open"))
	pd := res.Data.(*bclist.PageData)
	if !pd.Table.PrimaryAction.Disabled || !strings.Contains(pd.Table.PrimaryAction.DisabledTooltip, "recovery_document:issue") {
		t.Errorf("primary action not disabled with the missing permission: %+v", pd.Table.PrimaryAction)
	}
}

func TestTemplatesRenderListWithViewData(t *testing.T) {
	f := bctest.Seed()
	res := bclist.NewView(bctest.ListDeps(f)).Handle(bctest.CtxWith("billable_charge:list", "recovery_document:issue", "billable_charge:adjust"),
		bctest.Request(http.MethodGet, "/x", "", "status", "open"))
	pd := res.Data.(*bclist.PageData)
	html := bctest.RenderTo(t, bctest.Renderer(t), "billable-charge-list-content", pd)
	bctest.MustContain(t, html, `data-testid="billable-charge-list"`, `data-testid="billable-charge-status-tabs"`, `data-testid="billable-charge-row-o1"`, `data-testid="billable-charge-issue"`)
}
