package subscription_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"

	revenue "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	subform "github.com/erniealice/centymo-golang/domain/subscription/subscription/form"
)

func renderDrawer(t *testing.T, data *subform.Data) string {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, revenue.TemplatesFS, subscription.TemplatesFS)
	if err := r.Init(); err != nil {
		t.Fatalf("template set does not parse: %v", err)
	}
	w := httptest.NewRecorder()
	if err := r.Render(w, "subscription-drawer-form", data); err != nil {
		t.Fatalf("render subscription-drawer-form: %v", err)
	}
	return w.Body.String()
}

// openingTag returns the opening tag of the element carrying marker.
func openingTag(t *testing.T, out, marker string) string {
	t.Helper()
	i := strings.Index(out, marker)
	if i < 0 {
		t.Fatalf("%s missing", marker)
	}
	start := strings.LastIndex(out[:i], "<")
	return out[start : i+strings.Index(out[i:], ">")+1]
}

// The drawer is HTMX-swapped into a page whose CSP nonce differs from the partial's, so any inline
// <script> in it is blocked by the enforcing CSP (MMIS prod 2026-10-23: Save did nothing, the typed
// start date was ignored). Its wiring must come from the external centymo script + data-* hooks.
func TestSubscriptionDrawerHasNoInlineScripts(t *testing.T) {
	for _, suggest := range []bool{false, true} {
		out := renderDrawer(t, &subform.Data{
			SuggestPlanTerm: suggest, ClientID: "client-1", ClientLabel: `A & "B"`, ClientBillingCurrency: "PHP",
			DefaultTZ: "Asia/Manila", CommonLabels: pyeza.CommonLabels{},
		})
		if strings.Contains(out, "<script nonce") || strings.Count(out, "<script") != strings.Count(out, "<script src=") {
			t.Errorf("suggest=%v: drawer still carries an inline <script>", suggest)
		}
		if !strings.Contains(out, `<script src="/assets/js/centymo/subscription-drawer.js"></script>`) {
			t.Errorf("suggest=%v: subscription-drawer.js not included", suggest)
		}
		form := openingTag(t, out, "data-subscription-drawer")
		for _, want := range []string{`data-default-tz="Asia/Manila"`, `data-billing-currency="PHP"`, `data-client-id="client-1"`, `data-client-name="A &amp; &#34;B&#34;"`} {
			if !strings.Contains(form, want) {
				t.Errorf("suggest=%v: form tag missing %s: %s", suggest, want, form)
			}
		}
		if got := strings.Contains(form, "data-subscription-date-sync"); got == suggest {
			t.Errorf("suggest=%v: data-subscription-date-sync present=%v (the commencement-pricing form has no date sync)", suggest, got)
		}
	}
}

// Without JS state the server already renders the right visibility, so a plan-default / NONE
// enrollment never shows the required escalation inputs.
func TestSubscriptionDrawerEscalationInitialVisibility(t *testing.T) {
	cases := []struct {
		mode, scope                   string
		fixedHidden, intervalsHidden bool
	}{
		{"", "ESCALATION_SCOPE_WITHIN_AGREEMENT", true, false},
		{"ESCALATION_MODE_NONE", "ESCALATION_SCOPE_WITHIN_AGREEMENT", true, false},
		{"ESCALATION_MODE_FIXED_PERCENTAGE", "ESCALATION_SCOPE_WITHIN_AGREEMENT", false, false},
		{"ESCALATION_MODE_FIXED_PERCENTAGE", "ESCALATION_SCOPE_ON_RENEWAL", false, true},
	}
	for _, c := range cases {
		out := renderDrawer(t, &subform.Data{EscalationMode: c.mode, EscalationScope: c.scope, CommonLabels: pyeza.CommonLabels{}})
		fixed := openingTag(t, out, `data-testid="subscription-escalation-fixed-fields"`)
		if !strings.Contains(fixed, `data-lf-show-when="subscription-escalation-mode"`) || !strings.Contains(fixed, `data-lf-show-values="ESCALATION_MODE_FIXED_PERCENTAGE"`) {
			t.Errorf("mode %q: fixed fields lack the show-when hooks: %s", c.mode, fixed)
		}
		if got := strings.Contains(fixed, " hidden"); got != c.fixedHidden {
			t.Errorf("mode %q: fixed hidden=%v want %v", c.mode, got, c.fixedHidden)
		}
		intervals := openingTag(t, out, `data-testid="subscription-escalation-interval-fields"`)
		if got := strings.Contains(intervals, " hidden"); got != c.intervalsHidden {
			t.Errorf("mode %q scope %q: intervals hidden=%v want %v", c.mode, c.scope, got, c.intervalsHidden)
		}
	}
}
