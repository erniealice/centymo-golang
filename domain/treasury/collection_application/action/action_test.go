package action_test

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/erniealice/pyeza-golang/view"

	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
	raaction "github.com/erniealice/centymo-golang/domain/treasury/collection_application/action"
	ratest "github.com/erniealice/centymo-golang/domain/treasury/collection_application/internal/ratest"
)

func TestDrawerRendersFieldsAndPreselectsClient(t *testing.T) {
	f := &ratest.Fake{}
	h := raaction.NewReceiveApplyAction(ratest.ActionDeps(f))
	if res := h.Handle(ratest.CtxWith("collection_application:reverse"), ratest.Request(http.MethodGet, "/x", "")); res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Fatalf("denied => %v", res.Headers)
	}
	res := h.Handle(ratest.CtxWith("collection_application:create"), ratest.Request(http.MethodGet, "/action/collection/receive-apply?client_id=c2", ""))
	if res.Template != "receive-apply-drawer" {
		t.Fatalf("template = %q %v", res.Template, res.Headers)
	}
	html := ratest.RenderTo(t, res.Template, res.Data)
	ratest.MustContain(t, html,
		`data-testid="receive-apply-form"`, `id="receive-apply-client"`, `id="receive-apply-amount"`, `id="receive-apply-date"`,
		`id="receive-apply-reference"`, `data-testid="receive-apply-preview-button"`, `data-testid="receive-apply-submit"`,
		`hx-get="/action/collection/receive-apply/preview"`, `id="receive-apply-method"`, `Bank transfer`)
	if !regexpHas(html, `<option value="c2"[^>]*selected`) {
		t.Error("?client_id must preselect the client")
	}
}

func regexpHas(s, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(s)
}

func TestPreviewRendersPlanAndTotals(t *testing.T) {
	f := &ratest.Fake{}
	h := raaction.NewPreviewAction(ratest.ActionDeps(f))
	ctx := ratest.CtxWith("collection_application:create")
	q := url.Values{"client_id": {"c1"}, "amount": {"2000"}, "currency": {"PHP"}}
	res := h.Handle(ctx, ratest.Request(http.MethodGet, "/x?"+q.Encode(), ""))
	if res.Template != "receive-apply-preview" {
		t.Fatalf("template = %q %v", res.Template, res.Headers)
	}
	if len(f.PreviewedReqs) != 1 || f.PreviewedReqs[0].GetAmount() != 200000 || f.PreviewedReqs[0].GetClientId() != "c1" || f.PreviewedReqs[0].GetCurrency() != "PHP" {
		t.Fatalf("preview request = %v", f.PreviewedReqs)
	}
	html := ratest.RenderTo(t, res.Template, res.Data)
	ratest.MustContain(t, html, `data-testid="receive-apply-preview-table"`, "INV-1", "REC-000001",
		`data-testid="receive-apply-preview-applied">PHP 1,500.00`, `data-testid="receive-apply-preview-unapplied">PHP 500.00`)
	// N9 order is the plan's rank order: rent (rank 1) before the statement.
	if strings.Index(html, "INV-1") > strings.Index(html, "REC-000001") {
		t.Error("plan rows must render in rank order")
	}
}

func TestPreviewValidationAndRefusalsRenderInline(t *testing.T) {
	f := &ratest.Fake{}
	h := raaction.NewPreviewAction(ratest.ActionDeps(f))
	ctx := ratest.CtxWith("collection_application:create")
	l := ra.DefaultLabels()
	for _, tc := range []struct{ query, want string }{
		{"client_id=&amount=10&currency=PHP", l.Errors.ClientRequired},
		{"client_id=c1&amount=0&currency=PHP", l.Errors.AmountInvalid},
		{"client_id=c1&amount=abc&currency=PHP", l.Errors.AmountInvalid},
	} {
		res := h.Handle(ctx, ratest.Request(http.MethodGet, "/x?"+tc.query, ""))
		pd := res.Data.(*raaction.PreviewData)
		if pd.Error != tc.want {
			t.Errorf("%s => %q, want %q", tc.query, pd.Error, tc.want)
		}
	}
	if len(f.PreviewedReqs) != 0 {
		t.Error("invalid input must not reach the use case")
	}
	f.PreviewErr = ratest.CodedErr("currency_mismatch")
	res := h.Handle(ctx, ratest.Request(http.MethodGet, "/x?client_id=c1&amount=10&currency=USD", ""))
	if res.Data.(*raaction.PreviewData).Error != l.Errors.CurrencyMismatch {
		t.Errorf("currency mismatch => %q", res.Data.(*raaction.PreviewData).Error)
	}
	html := ratest.RenderTo(t, res.Template, res.Data)
	ratest.MustContain(t, html, `data-testid="receive-apply-preview-error"`)
	if res := h.Handle(ratest.CtxWith(), ratest.Request(http.MethodGet, "/x?client_id=c1&amount=10", "")); res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Errorf("denied => %v", res.Headers)
	}
}

func TestReceivePostCommitsAndRedirectsToTheCollection(t *testing.T) {
	f := &ratest.Fake{}
	h := raaction.NewReceiveApplyAction(ratest.ActionDeps(f))
	ctx := ratest.CtxWith("collection_application:create")
	form := url.Values{"client_id": {"c1"}, "amount": {"2,000.50"}, "currency": {"PHP"}, "payment_date": {"2026-02-01"}, "collection_method_id": {"m1"}, "reference_number": {"CHK-9"}}
	res := h.Handle(ctx, ratest.Request(http.MethodPost, "/x", form.Encode()))
	if res.StatusCode != http.StatusOK || res.Headers["HX-Redirect"] != "/collections/detail/col-1" {
		t.Fatalf("success => %d %v", res.StatusCode, res.Headers)
	}
	got := f.ReceivedReqs[0]
	if got.GetClientId() != "c1" || got.GetAmount() != 200050 || got.GetCurrency() != "PHP" || got.GetPaymentDate() != "2026-02-01" ||
		got.GetCollectionMethodId() != "m1" || got.GetReferenceNumber() != "CHK-9" {
		t.Errorf("request = %v", got)
	}
}

func TestReceivePostValidationAndRefusalMapping(t *testing.T) {
	l := ra.DefaultLabels()
	f := &ratest.Fake{}
	h := raaction.NewReceiveApplyAction(ratest.ActionDeps(f))
	ctx := ratest.CtxWith("collection_application:create")
	if res := h.Handle(ctx, ratest.Request(http.MethodPost, "/x", "amount=10&currency=PHP")); res.Headers["HX-Error-Message"] != l.Errors.ClientRequired {
		t.Errorf("no client => %v", res.Headers)
	}
	if res := h.Handle(ctx, ratest.Request(http.MethodPost, "/x", "client_id=c1&amount=-5&currency=PHP")); res.Headers["HX-Error-Message"] != l.Errors.AmountInvalid {
		t.Errorf("bad amount => %v", res.Headers)
	}
	if len(f.ReceivedReqs) != 0 {
		t.Error("invalid input must not reach the use case")
	}
	for code, want := range map[string]string{
		"currency_mismatch":    l.Errors.CurrencyMismatch,
		"missing_posting":      l.Errors.MissingPosting,
		"transaction_required": l.Errors.TransactionRequired,
	} {
		f := &ratest.Fake{ReceiveErr: ratest.CodedErr(code)}
		res := raaction.NewReceiveApplyAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", "client_id=c1&amount=10&currency=PHP"))
		if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != want {
			t.Errorf("%s => %d %q, want %q", code, res.StatusCode, res.Headers["HX-Error-Message"], want)
		}
	}
	f2 := &ratest.Fake{ReceiveErr: ratest.CodedErr("brand_new")}
	if res := raaction.NewReceiveApplyAction(ratest.ActionDeps(f2)).Handle(ctx, ratest.Request(http.MethodPost, "/x", "client_id=c1&amount=10&currency=PHP")); res.Headers["HX-Error-Message"] != "Something went wrong" {
		t.Errorf("unknown code => %v", res.Headers)
	}
}

func TestReverseRefreshesApplicationsAndMapsRefusals(t *testing.T) {
	f := &ratest.Fake{}
	h := raaction.NewReverseAction(ratest.ActionDeps(f))
	if res := h.Handle(ratest.CtxWith("collection_application:create"), ratest.Request(http.MethodPost, "/x", "", "id", "a1")); res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Fatalf("denied => %v", res.Headers)
	}
	ctx := ratest.CtxWith("collection_application:reverse")
	res := h.Handle(ctx, ratest.Request(http.MethodPost, "/x", "", "id", "a1"))
	if res.StatusCode != http.StatusOK || !strings.Contains(res.Headers["HX-Trigger"], `"refreshTable":"`+ra.ApplicationsTableID+`"`) {
		t.Fatalf("success => %d %v", res.StatusCode, res.Headers)
	}
	if len(f.ReversedReqs) != 1 || f.ReversedReqs[0].GetCollectionApplicationId() != "a1" {
		t.Errorf("reverse request = %v", f.ReversedReqs)
	}
	f.ReverseErr = ratest.CodedErr("already_reversed")
	res = raaction.NewReverseAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", "", "id", "a1"))
	if res.Headers["HX-Error-Message"] != ra.DefaultLabels().Errors.AlreadyReversed {
		t.Errorf("refusal => %v", res.Headers)
	}
}

func TestHandlersFailClosedWhenUnwired(t *testing.T) {
	d := ratest.ActionDeps(&ratest.Fake{})
	d.UseCases = &ra.UseCases{}
	for name, v := range map[string]view.View{
		"receive": raaction.NewReceiveApplyAction(d), "preview": raaction.NewPreviewAction(d), "reverse": raaction.NewReverseAction(d),
	} {
		res := v.Handle(ratest.CtxWith("collection_application:create", "collection_application:reverse"), ratest.Request(http.MethodGet, "/x", "", "id", "a1"))
		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("%s unwired => %d", name, res.StatusCode)
		}
	}
}

// R5 M6: the amount goes through the shared pyeza parser; "1.-5", "1.+5", "-1"
// and every other malformed or non-positive value never reaches the use case.
func TestAmountUsesTheSharedCentavosParser(t *testing.T) {
	ctx := ratest.CtxWith("collection_application:create")
	for in, want := range map[string]int64{"10": 1000, "10.5": 1050, "1,234.05": 123405} {
		f := &ratest.Fake{}
		res := raaction.NewReceiveApplyAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", url.Values{"client_id": {"c1"}, "amount": {in}, "currency": {"PHP"}}.Encode()))
		if res.StatusCode != http.StatusOK || len(f.ReceivedReqs) != 1 || f.ReceivedReqs[0].GetAmount() != want {
			t.Errorf("%q => %d %v, want %d centavos", in, res.StatusCode, f.ReceivedReqs, want)
		}
	}
	for _, in := range []string{"", "0", "-1", "1.-5", "1.+5", ".5", "abc", "1.234", "99999999999999999999"} {
		f := &ratest.Fake{}
		res := raaction.NewReceiveApplyAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", url.Values{"client_id": {"c1"}, "amount": {in}, "currency": {"PHP"}}.Encode()))
		if res.StatusCode != http.StatusUnprocessableEntity || len(f.ReceivedReqs) != 0 {
			t.Errorf("%q => %d, received %v; want a refusal without a use-case call", in, res.StatusCode, f.ReceivedReqs)
		}
	}
}

// A strict-gate denial after the view's permission check maps to the shared
// PermissionDenied message on every handler, inline in the preview (R5 M1, C20).
func TestStrictGateDenialMapsToPermissionDenied(t *testing.T) {
	ctx := ratest.CtxWith("collection_application:create", "collection_application:reverse")
	f := &ratest.Fake{ReceiveErr: ratest.CodedErr("permission_denied")}
	res := raaction.NewReceiveApplyAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", "client_id=c1&amount=10&currency=PHP"))
	if res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Errorf("receive denial => %v", res.Headers)
	}
	f = &ratest.Fake{PreviewErr: ratest.CodedErr("permission_denied")}
	res = raaction.NewPreviewAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodGet, "/x?client_id=c1&amount=10&currency=PHP", ""))
	if res.Data.(*raaction.PreviewData).Error != "Permission denied" {
		t.Errorf("preview denial => %q", res.Data.(*raaction.PreviewData).Error)
	}
	f = &ratest.Fake{ReverseErr: ratest.CodedErr("permission_denied")}
	res = raaction.NewReverseAction(ratest.ActionDeps(f)).Handle(ctx, ratest.Request(http.MethodPost, "/x", "", "id", "a1"))
	if res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Errorf("reverse denial => %v", res.Headers)
	}
}

// R5 M4: with the async client search bound, the drawer never lists clients: it
// reads only the preselected one (ReadClient) so client #101+ stays reachable.
func TestDrawerUsesAsyncClientSearchAndReadsOnlyThePreselect(t *testing.T) {
	f := &ratest.Fake{}
	deps := ratest.ActionDeps(f)
	deps.SearchClientURL = "/action/revenue/search/clients"
	res := raaction.NewReceiveApplyAction(deps).Handle(ratest.CtxWith("collection_application:create"),
		ratest.Request(http.MethodGet, "/action/collection/receive-apply?client_id=c2", ""))
	html := ratest.RenderTo(t, res.Template, res.Data)
	ratest.MustContain(t, html, `id="receive-apply-client"`, `/action/revenue/search/clients`, `Tenant Two`)
	if f.ListClientPages != 0 || len(f.ReadClients) != 1 || f.ReadClients[0] != "c2" {
		t.Errorf("list pages=%d reads=%v; want 0 list calls and one ReadClient(c2)", f.ListClientPages, f.ReadClients)
	}
}
