package list_test

import (
	"net/http"
	"testing"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rdtest "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/internal/rdtest"
	rdlist "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/list"
)

func TestListDeniedAndFailClosed(t *testing.T) {
	f := rdtest.Seed()
	if res := rdlist.NewView(rdtest.ListDeps(f)).Handle(rdtest.CtxWith(), rdtest.Request(http.MethodGet, "/x", "", "status", "issued")); res.StatusCode != http.StatusForbidden {
		t.Fatalf("no permission => %d", res.StatusCode)
	}
	deps := rdtest.ListDeps(f)
	deps.UseCases = &rd.UseCases{}
	if res := rdlist.NewView(deps).Handle(rdtest.CtxWith("recovery_document:list"), rdtest.Request(http.MethodGet, "/x", "", "status", "issued")); res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unwired => %d", res.StatusCode)
	}
}

func TestListTabsFilterByStatus(t *testing.T) {
	deps := rdtest.ListDeps(rdtest.Seed())
	ctx := rdtest.CtxWith("recovery_document:list")
	res := rdlist.NewView(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "status", "issued"))
	pd := res.Data.(*rdlist.PageData)
	if len(pd.Table.Rows) != 2 || pd.Table.Rows[0].ID != "d1" {
		t.Fatalf("issued rows = %+v", pd.Table.Rows)
	}
	if len(pd.StatusTabs) != 2 || pd.StatusTabs[0].Key != "issued" || pd.StatusTabs[1].Key != "void" {
		t.Errorf("tabs = %+v", pd.StatusTabs)
	}
	res = rdlist.NewView(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "status", "void"))
	pd = res.Data.(*rdlist.PageData)
	if len(pd.Table.Rows) != 1 || pd.Table.Rows[0].ID != "d2" {
		t.Fatalf("void rows = %+v", pd.Table.Rows)
	}
	html := rdtest.RenderTo(t, "recovery-document-list-content", pd)
	rdtest.MustContain(t, html, `data-testid="recovery-document-list"`, `data-testid="recovery-document-row-d2"`)
}
