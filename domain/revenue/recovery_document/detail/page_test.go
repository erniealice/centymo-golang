package detail_test

import (
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	"net/http"
	"strings"
	"testing"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rddetail "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/detail"
	rdtest "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/internal/rdtest"
)

func TestDetailDeniedNotFoundAndTabs(t *testing.T) {
	f := rdtest.Seed()
	deps := rdtest.DetailDeps(f)
	if res := rddetail.NewView(deps).Handle(rdtest.CtxWith(), rdtest.Request(http.MethodGet, "/x", "", "id", "d1")); res.StatusCode != http.StatusForbidden {
		t.Fatalf("denied => %d", res.StatusCode)
	}
	ctx := rdtest.CtxWith("recovery_document:read", "recovery_document:void", "collection_application:create", "collection_application:reverse")
	if res := rddetail.NewView(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "nope")); res.StatusCode != http.StatusNotFound {
		t.Fatalf("missing => %d", res.StatusCode)
	}

	res := rddetail.NewView(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "d1"))
	pd := res.Data.(*rddetail.PageData)
	// total 2000.00, credit note -300.00, applied 1200.00 (reversed row not counted):
	// balance = total + credited - applied = 500.00 (R5 M5, same rule as aging).
	if pd.Applied != "PHP 1,200.00" || pd.Balance != "PHP 500.00" || pd.Total != "PHP 2,000.00" {
		t.Errorf("amounts total=%q applied=%q balance=%q", pd.Total, pd.Applied, pd.Balance)
	}
	if len(pd.TabItems) != 4 {
		t.Fatalf("tabs = %d", len(pd.TabItems))
	}
	html := rdtest.RenderTo(t, "recovery-document-detail-content", pd)
	rdtest.MustContain(t, html, `data-testid="recovery-document-detail"`, `data-testid="recovery-document-void"`, `data-testid="recovery-document-apply-payment"`, `data-testid="recovery-document-info-balance"`)

	for _, tab := range []struct{ name, tid string }{
		{"lines", "recovery-document-line-l1"},
		{"applications", "collection-application-reverse-a1"},
		{"credit-notes", "recovery-document-credit-note-d3"},
	} {
		res := rddetail.NewTabAction(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "d1", "tab", tab.name))
		if res.Template != "recovery-document-tab-"+tab.name {
			t.Fatalf("%s template = %q (%d %v)", tab.name, res.Template, res.StatusCode, res.Error)
		}
		html := rdtest.RenderTo(t, res.Template, res.Data)
		rdtest.MustContain(t, html, tab.tid)
	}
	// Only the APPLIED original is reversible; the REVERSED row has no action.
	res = rddetail.NewTabAction(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "d1", "tab", "applications"))
	html = rdtest.RenderTo(t, "recovery-document-tab-applications", res.Data)
	if strings.Contains(html, `collection-application-reverse-a2`) {
		t.Error("a reversed application must not offer Reverse")
	}
}

func TestDetailCreditNoteShowsCorrectedDocumentAndVoidDisabled(t *testing.T) {
	deps := rdtest.DetailDeps(rdtest.Seed())
	ctx := rdtest.CtxWith("recovery_document:read")
	res := rddetail.NewView(deps).Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "d3"))
	pd := res.Data.(*rddetail.PageData)
	if pd.CorrectsNumber != "REC-000001" || pd.BalanceLabel != rd.DefaultLabels().Detail.CreditBalance {
		t.Errorf("corrects=%q balanceLabel=%q", pd.CorrectsNumber, pd.BalanceLabel)
	}
	if pd.CanVoid || !strings.Contains(pd.VoidDisabledTip, "recovery_document:void") {
		t.Errorf("void must be disabled with the missing permission: can=%v tip=%q", pd.CanVoid, pd.VoidDisabledTip)
	}
	html := rdtest.RenderTo(t, "recovery-document-detail-content", pd)
	rdtest.MustContain(t, html, `data-testid="recovery-document-void"`, `disabled`)
}

// R5 M5 (click path 2): a 100.00 statement corrected by a 20.00 credit note owes
// 80.00 on the detail, the same as preview and aging.
func TestDetailBalanceIncludesIssuedCreditNotes(t *testing.T) {
	f := rdtest.Seed()
	f.Docs = []*recoverydocumentpb.RecoveryDocument{
		rdtest.Doc("s1", "REC-1", rdtest.StIssued, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_STATEMENT, 10000, ""),
		rdtest.Doc("n1", "CN-1", rdtest.StIssued, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_CREDIT_NOTE, -2000, "s1"),
		rdtest.Doc("n2", "CN-2", rdtest.StVoid, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_CREDIT_NOTE, -5000, "s1"),
	}
	f.Apps = nil
	res := rddetail.NewView(rdtest.DetailDeps(f)).Handle(rdtest.CtxWith("recovery_document:read"), rdtest.Request(http.MethodGet, "/x", "", "id", "s1"))
	pd := res.Data.(*rddetail.PageData)
	if pd.Total != "PHP 100.00" || pd.Balance != "PHP 80.00" {
		t.Errorf("total=%q balance=%q, want PHP 100.00 / PHP 80.00 (the void note must not count)", pd.Total, pd.Balance)
	}
}

// R5 M4: the detail resolves its one client name through ReadClient, never a
// full client list.
func TestDetailReadsOneClientNotTheWholeList(t *testing.T) {
	name := "Tenant One"
	f := rdtest.Seed()
	f.Clients = []*clientpb.Client{{Id: "c1", Name: &name}}
	res := rddetail.NewView(rdtest.DetailDeps(f)).Handle(rdtest.CtxWith("recovery_document:read"), rdtest.Request(http.MethodGet, "/x", "", "id", "d1"))
	pd := res.Data.(*rddetail.PageData)
	if pd.Client != name || len(f.ReadClients) != 1 || f.ListClientPages != 0 {
		t.Errorf("client=%q reads=%v listPages=%d", pd.Client, f.ReadClients, f.ListClientPages)
	}
}
