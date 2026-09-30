package action_test

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	bcaction "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/action"
	bctest "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/internal/bctest"
)

func TestIssueDeniedWithoutPermission(t *testing.T) {
	f := bctest.Seed()
	res := bcaction.NewIssueAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("billable_charge:list"), bctest.Request(http.MethodPost, "/x", "charge_id=o1&document_series_id=s-active&issuance_key=k"))
	if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Fatalf("result = %d %v", res.StatusCode, res.Headers)
	}
	if len(f.IssuedReqs) != 0 {
		t.Error("use case must not run without the permission")
	}
}

func TestIssueGetMintsKeyPreselectsAndFiltersSeries(t *testing.T) {
	f := bctest.Seed()
	deps := bctest.ActionDeps(f)
	res := bcaction.NewIssueAction(deps).Handle(bctest.CtxWith("recovery_document:issue"), bctest.Request(http.MethodGet, "/action/billable-charge/issue?id=o2", ""))
	if res.Template != "billable-charge-issue-drawer" {
		t.Fatalf("template = %q status %d hdr %v", res.Template, res.StatusCode, res.Headers)
	}
	html := renderIssue(t, res.Data)
	bctest.MustContain(t, html,
		`name="issuance_key" value="key-fixed"`,
		`data-testid="billable-charge-select-o2"`,
		`data-testid="billable-charge-issue-submit"`,
		`value="s-active"`)
	if strings.Contains(html, `value="s-invoice"`) || strings.Contains(html, `value="s-retired"`) {
		t.Error("series picker must offer only active recovery-document series")
	}
	// ?id=o2 preselects o2 only.
	tag := func(id string) string {
		return regexp.MustCompile(`<input type="checkbox"[^>]*id="billable-charge-select-` + id + `"[^>]*>`).FindString(html)
	}
	if !strings.Contains(tag("o2"), "checked") || tag("o1") == "" || strings.Contains(tag("o1"), "checked") {
		t.Errorf("only the requested charge is preselected: o1=%q o2=%q", tag("o1"), tag("o2"))
	}
}

func renderIssue(t *testing.T, data any) string {
	t.Helper()
	return bctest.RenderTo(t, bctest.Renderer(t), "billable-charge-issue-drawer", bctest.InjectCommon(data))
}

func TestIssuePostSendsKeyAndChargesAndClosesSheet(t *testing.T) {
	f := bctest.Seed()
	form := url.Values{"charge_id": {"o1", "o2", "o1"}, "document_series_id": {"s-active"}, "issuance_key": {"key-1"}, "issue_date": {"2026-02-01"}, "due_date": {"2026-02-15"}}
	res := bcaction.NewIssueAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("recovery_document:issue"), bctest.Request(http.MethodPost, "/x", form.Encode()))
	if res.StatusCode != http.StatusOK || !strings.Contains(res.Headers["HX-Trigger"], `"refreshTable":"billable-charges-table"`) {
		t.Fatalf("result = %d %v", res.StatusCode, res.Headers)
	}
	if len(f.IssuedReqs) != 1 {
		t.Fatalf("issue calls = %d", len(f.IssuedReqs))
	}
	got := f.IssuedReqs[0]
	if got.GetIssuanceKey() != "key-1" || got.GetDocumentSeriesId() != "s-active" || strings.Join(got.GetBillableChargeIds(), ",") != "o1,o2" ||
		got.GetIssueDate() != "2026-02-01" || got.GetDueDate() != "2026-02-15" {
		t.Errorf("request = %v", got)
	}
}

func TestIssuePostRejectsEmptySelectionAndMissingKey(t *testing.T) {
	f := bctest.Seed()
	h := bcaction.NewIssueAction(bctest.ActionDeps(f))
	res := h.Handle(bctest.CtxWith("recovery_document:issue"), bctest.Request(http.MethodPost, "/x", "document_series_id=s-active&issuance_key=k"))
	if res.Headers["HX-Error-Message"] != rd.DefaultLabels().Errors.NothingSelected {
		t.Errorf("no selection => %v", res.Headers)
	}
	res = h.Handle(bctest.CtxWith("recovery_document:issue"), bctest.Request(http.MethodPost, "/x", "charge_id=o1&document_series_id=s-active"))
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("missing key => %d", res.StatusCode)
	}
	if len(f.IssuedReqs) != 0 {
		t.Error("use case must not run on invalid input")
	}
}

// Refusal-code mapping: each ErrorCode() resolves to its Lyngua message.
func TestIssuePostMapsRefusalCodes(t *testing.T) {
	labels := rd.DefaultLabels()
	for code, want := range map[string]string{
		"series_retired":    labels.Errors.SeriesRetired,
		"issuance_conflict": labels.Errors.IssuanceConflict,
		"missing_posting":   labels.Errors.MissingPosting,
		"not_open":          labels.Errors.NotOpen,
	} {
		f := bctest.Seed()
		f.IssueErr = bctest.CodedErr(code)
		res := bcaction.NewIssueAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("recovery_document:issue"),
			bctest.Request(http.MethodPost, "/x", "charge_id=o1&document_series_id=s-active&issuance_key=k"))
		if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != want {
			t.Errorf("%s => %d %q, want %q", code, res.StatusCode, res.Headers["HX-Error-Message"], want)
		}
	}
	f := bctest.Seed()
	f.IssueErr = bctest.CodedErr("brand_new_code")
	res := bcaction.NewIssueAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("recovery_document:issue"),
		bctest.Request(http.MethodPost, "/x", "charge_id=o1&document_series_id=s-active&issuance_key=k"))
	if res.Headers["HX-Error-Message"] != "Something went wrong" {
		t.Errorf("unknown code => %q", res.Headers["HX-Error-Message"])
	}
}

func TestAdjustDrawerOnlyForIssuedOriginal(t *testing.T) {
	f := bctest.Seed()
	h := bcaction.NewAdjustAction(bctest.ActionDeps(f))
	ctx := bctest.CtxWith("billable_charge:adjust")
	res := h.Handle(ctx, bctest.Request(http.MethodGet, "/x", "", "id", "i1"))
	if res.Template != "billable-charge-adjust-drawer" {
		t.Fatalf("issued original => %q %v", res.Template, res.Headers)
	}
	html := bctest.RenderTo(t, bctest.Renderer(t), "billable-charge-adjust-drawer", bctest.InjectCommon(res.Data))
	bctest.MustContain(t, html, `data-testid="billable-charge-adjust-submit"`, `id="billable-charge-new-amount"`, `value="1000.00"`)
	for _, id := range []string{"i2", "o1", "missing"} {
		res = h.Handle(ctx, bctest.Request(http.MethodGet, "/x", "", "id", id))
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("%s => %d, want a refusal", id, res.StatusCode)
		}
	}
}

func TestAdjustPostAmountsAndRefusals(t *testing.T) {
	f := bctest.Seed()
	h := bcaction.NewAdjustAction(bctest.ActionDeps(f))
	ctx := bctest.CtxWith("billable_charge:adjust")
	res := h.Handle(ctx, bctest.Request(http.MethodPost, "/x", "new_amount=750.50&reason=Metering+error", "id", "i1"))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("success => %d %v", res.StatusCode, res.Headers)
	}
	if len(f.AdjustedReqs) != 1 || f.AdjustedReqs[0].GetNewAmount() != 75050 || f.AdjustedReqs[0].GetReason() != "Metering error" || f.AdjustedReqs[0].GetBillableChargeId() != "i1" {
		t.Errorf("adjust request = %v", f.AdjustedReqs)
	}
	for _, bad := range []string{"new_amount=&reason=x", "new_amount=-1&reason=x", "new_amount=1&reason=", "new_amount=1.234&reason=x", "new_amount=abc&reason=x"} {
		res = h.Handle(ctx, bctest.Request(http.MethodPost, "/x", bad, "id", "i1"))
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("%q => %d, want invalid", bad, res.StatusCode)
		}
	}
	if len(f.AdjustedReqs) != 1 {
		t.Error("invalid input must not reach the use case")
	}
	f.AdjustErr = bctest.CodedErr("adjust_not_downward")
	res = h.Handle(ctx, bctest.Request(http.MethodPost, "/x", "new_amount=2000&reason=x", "id", "i1"))
	if res.Headers["HX-Error-Message"] != bc.DefaultLabels().Errors.AdjustNotDownward {
		t.Errorf("refusal => %q", res.Headers["HX-Error-Message"])
	}
	res = h.Handle(bctest.CtxWith("billable_charge:list"), bctest.Request(http.MethodPost, "/x", "new_amount=1&reason=x", "id", "i1"))
	if res.Headers["HX-Error-Message"] != "Permission denied" {
		t.Errorf("denied => %v", res.Headers)
	}
}

// R5 M6: the amount goes through the shared pyeza parser; malformed and
// non-positive input (including "1.-5" / "1.+5") never reaches the use case.
func TestAdjustAmountUsesTheSharedCentavosParser(t *testing.T) {
	ctx := bctest.CtxWith("billable_charge:adjust")
	for in, want := range map[string]int64{"12": 1200, "12.5": 1250, "12.50": 1250, "1,234.05": 123405} {
		f := bctest.Seed()
		res := bcaction.NewAdjustAction(bctest.ActionDeps(f)).Handle(ctx, bctest.Request(http.MethodPost, "/x", url.Values{"new_amount": {in}, "reason": {"r"}}.Encode(), "id", "i1"))
		if res.StatusCode != http.StatusOK || len(f.AdjustedReqs) != 1 || f.AdjustedReqs[0].GetNewAmount() != want {
			t.Errorf("%q => %d %v, want %d centavos", in, res.StatusCode, f.AdjustedReqs, want)
		}
	}
	for _, in := range []string{"", "0", "-1", "1.-5", "1.+5", ".5", "abc", "1.234", "1e3", "99999999999999999999"} {
		f := bctest.Seed()
		res := bcaction.NewAdjustAction(bctest.ActionDeps(f)).Handle(ctx, bctest.Request(http.MethodPost, "/x", url.Values{"new_amount": {in}, "reason": {"r"}}.Encode(), "id", "i1"))
		if res.StatusCode != http.StatusUnprocessableEntity || len(f.AdjustedReqs) != 0 {
			t.Errorf("%q => %d, adjusted %v; want a refusal without a use-case call", in, res.StatusCode, f.AdjustedReqs)
		}
	}
}

// A strict-gate denial after the view's permission check maps to the shared
// PermissionDenied message on both writes (R5 M1, C20).
func TestStrictGateDenialMapsToPermissionDenied(t *testing.T) {
	f := bctest.Seed()
	f.IssueErr = bctest.CodedErr("permission_denied")
	res := bcaction.NewIssueAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("recovery_document:issue"),
		bctest.Request(http.MethodPost, "/x", "charge_id=o1&document_series_id=s-active&issuance_key=k"))
	if got := res.Headers["HX-Error-Message"]; got != "Permission denied" {
		t.Errorf("issue denial => %q", got)
	}
	f = bctest.Seed()
	f.AdjustErr = bctest.CodedErr("permission_denied")
	res = bcaction.NewAdjustAction(bctest.ActionDeps(f)).Handle(bctest.CtxWith("billable_charge:adjust"),
		bctest.Request(http.MethodPost, "/x", "new_amount=1&reason=x", "id", "i1"))
	if got := res.Headers["HX-Error-Message"]; got != "Permission denied" {
		t.Errorf("adjust denial => %q", got)
	}
}
