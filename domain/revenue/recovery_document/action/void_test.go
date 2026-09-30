package action_test

import (
	"net/http"
	"testing"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rdaction "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/action"
	rdtest "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/internal/rdtest"
)

func TestVoidDrawerAndPost(t *testing.T) {
	f := rdtest.Seed()
	h := rdaction.NewVoidAction(rdtest.ActionDeps(f))
	if res := h.Handle(rdtest.CtxWith("recovery_document:read"), rdtest.Request(http.MethodGet, "/x", "", "id", "d1")); res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Fatalf("denied => %v", res.Headers)
	}
	ctx := rdtest.CtxWith("recovery_document:void")
	res := h.Handle(ctx, rdtest.Request(http.MethodGet, "/x", "", "id", "d1"))
	if res.Template != "recovery-document-void-drawer" {
		t.Fatalf("drawer template = %q %v", res.Template, res.Headers)
	}
	html := rdtest.RenderTo(t, res.Template, rdtest.InjectCommon(res.Data))
	rdtest.MustContain(t, html, `data-testid="recovery-document-void-form"`, `id="recovery-document-void-reason"`, `data-testid="recovery-document-void-submit"`)

	// reason required
	res = h.Handle(ctx, rdtest.Request(http.MethodPost, "/x", "reason=+", "id", "d1"))
	if res.Headers["HX-Error-Message"] != rd.DefaultLabels().Errors.VoidReasonRequired || len(f.VoidedReqs) != 0 {
		t.Fatalf("empty reason => %v (calls %d)", res.Headers, len(f.VoidedReqs))
	}
	// success redirects to the detail page
	res = h.Handle(ctx, rdtest.Request(http.MethodPost, "/x", "reason=Duplicate", "id", "d1"))
	if res.StatusCode != http.StatusOK || res.Headers["HX-Redirect"] != "/revenue/recovery-documents/detail/d1" {
		t.Fatalf("success => %d %v", res.StatusCode, res.Headers)
	}
	if len(f.VoidedReqs) != 1 || f.VoidedReqs[0].GetRecoveryDocumentId() != "d1" || f.VoidedReqs[0].GetReason() != "Duplicate" {
		t.Errorf("void request = %v", f.VoidedReqs)
	}
}

// Refusal-code mapping: void refused while applied / credit-noted.
func TestVoidPostMapsRefusalCodes(t *testing.T) {
	l := rd.DefaultLabels()
	for code, want := range map[string]string{
		"has_applications": l.Errors.HasApplications,
		"has_credit_notes": l.Errors.HasCreditNotes,
		"already_void":     l.Errors.AlreadyVoid,
	} {
		f := rdtest.Seed()
		f.VoidErr = rdtest.CodedErr(code)
		res := rdaction.NewVoidAction(rdtest.ActionDeps(f)).Handle(rdtest.CtxWith("recovery_document:void"), rdtest.Request(http.MethodPost, "/x", "reason=x", "id", "d1"))
		if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != want {
			t.Errorf("%s => %d %q, want %q", code, res.StatusCode, res.Headers["HX-Error-Message"], want)
		}
	}
	f := rdtest.Seed()
	f.VoidErr = rdtest.CodedErr("brand_new")
	res := rdaction.NewVoidAction(rdtest.ActionDeps(f)).Handle(rdtest.CtxWith("recovery_document:void"), rdtest.Request(http.MethodPost, "/x", "reason=x", "id", "d1"))
	if res.Headers["HX-Error-Message"] != "Something went wrong" {
		t.Errorf("unknown code => %q", res.Headers["HX-Error-Message"])
	}
}

func TestVoidFailsClosedWhenUnwired(t *testing.T) {
	deps := rdtest.ActionDeps(rdtest.Seed())
	deps.UseCases = &rd.UseCases{}
	res := rdaction.NewVoidAction(deps).Handle(rdtest.CtxWith("recovery_document:void"), rdtest.Request(http.MethodGet, "/x", "", "id", "d1"))
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unwired => %d", res.StatusCode)
	}
}

// A strict-gate denial after the view's permission check maps to the shared
// PermissionDenied message, not the generic error (R5 M1, C20).
func TestVoidStrictGateDenialMapsToPermissionDenied(t *testing.T) {
	f := rdtest.Seed()
	f.VoidErr = rdtest.CodedErr("permission_denied")
	res := rdaction.NewVoidAction(rdtest.ActionDeps(f)).Handle(rdtest.CtxWith("recovery_document:void"), rdtest.Request(http.MethodPost, "/x", "reason=x", "id", "d1"))
	if got := res.Headers["HX-Error-Message"]; got != "Permission denied" {
		t.Errorf("denial => %q, want the shared PermissionDenied message", got)
	}
}
