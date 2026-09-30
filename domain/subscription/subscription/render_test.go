package subscription_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	revenue "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	alt "github.com/erniealice/centymo-golang/domain/subscription/agreement_line_term"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	subdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription/detail"
	subform "github.com/erniealice/centymo-golang/domain/subscription/subscription/form"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
)

func render(t *testing.T, name string, data any) string {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, revenue.TemplatesFS, subscription.TemplatesFS)
	if err := r.Init(); err != nil {
		t.Fatalf("template set does not parse: %v", err)
	}
	w := httptest.NewRecorder()
	if err := r.Render(w, name, data); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return w.Body.String()
}

func wantAll(t *testing.T, out string, subs ...string) {
	t.Helper()
	for _, s := range subs {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestChargeTermsTabRenders(t *testing.T) {
	l := alt.DefaultLabels()
	to := "2026-12-31"
	bps := int32(250)
	tbl := alt.Build([]*agreementlinetermpb.AgreementLineTerm{
		{Id: "t1", ProductPricePlanId: "ppp1", ChargePolicyVersionId: "v1", MarkupBps: &bps, EffectiveFrom: "2026-01-01", Origin: agreementlinetermpb.AgreementLineTermOrigin_AGREEMENT_LINE_TERM_ORIGIN_COPIED},
		{Id: "t2", ProductPricePlanId: "ppp1", ChargePolicyVersionId: "v2", EffectiveFrom: "2026-06-01", EffectiveTo: &to},
	}, l, types.TableLabels{}, map[string]string{"ppp1": "Electricity"}, func(id string) (string, string) { return "Utility recovery", "v" + id[1:] })
	out := render(t, "subscription-tab-charge-terms", &subdetail.PageData{ChargeTerms: &alt.TabData{Labels: l, Table: tbl}})
	wantAll(t, out, `data-testid="charge-terms-tab"`, `data-id="t1"`, "Electricity", "Utility recovery", "2.50%", l.Detail.OpenEnded, "2026-12-31", l.Enums.OriginCopied, `data-testid="charge-terms-pinned-notice"`)
	out = render(t, "subscription-tab-charge-terms", &subdetail.PageData{ChargeTerms: &alt.TabData{Labels: l, Failed: true, Table: alt.Build(nil, l, types.TableLabels{}, nil, nil)}})
	wantAll(t, out, `data-testid="charge-terms-error"`, l.Detail.LoadFailed)
	if strings.Contains(out, `data-testid="charge-terms-pinned-notice"`) {
		t.Error("a failed read must not show the table notice")
	}
}

// W5 finding: the agreement drawer's fixed-escalation inputs are `required`; unless the section is
// hidden (and, via the pyeza data-lf-show-when helper, disabled) for plan-default / NONE, an
// agreement without escalation cannot be submitted.
func TestAgreementEscalationFieldsOnlyForFixedPercentage(t *testing.T) {
	section := func(out string) string {
		i := strings.Index(out, `data-testid="subscription-escalation-fixed-fields"`)
		if i < 0 {
			t.Fatal("escalation fixed-fields section missing")
		}
		start := strings.LastIndex(out[:i], "<div")
		return out[start : i+strings.Index(out[i:], ">")+1]
	}
	for mode, wantHidden := range map[string]bool{"": true, "ESCALATION_MODE_NONE": true, "ESCALATION_MODE_FIXED_PERCENTAGE": false} {
		tag := section(render(t, "subscription-drawer-form", &subform.Data{EscalationMode: mode, CommonLabels: pyeza.CommonLabels{}}))
		wantAll(t, tag, `data-lf-show-when="subscription-escalation-mode"`, `data-lf-show-values="ESCALATION_MODE_FIXED_PERCENTAGE"`)
		if got := strings.Contains(tag, " hidden"); got != wantHidden {
			t.Errorf("mode %q: hidden=%v want %v (%s)", mode, got, wantHidden, tag)
		}
	}
}
