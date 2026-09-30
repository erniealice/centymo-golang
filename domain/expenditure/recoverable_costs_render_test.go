package expenditure_test

import (
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	abform "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch/form"
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	cscform "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/form"
	csctable "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/table"
	epkg "github.com/erniealice/centymo-golang/domain/expenditure/expenditure"
	expdetail "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/detail"
	revenue "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
)

var inlineScript = regexp.MustCompile(`<script[^>]*>`)

// stubSigner is a deterministic pyeza.WorkspaceFormSigner: the signature names the bound path.
type stubSigner struct{}

func (stubSigner) SignFields(workspaceID, actionPath string) (string, error) {
	return "sig-for-" + workspaceID + "-on-" + actionPath, nil
}

// renderRaw renders without the C15 fragment checks (for pre-existing templates outside this plan).
func renderRaw(t *testing.T, name string, data any) string {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, revenue.TemplatesFS, epkg.TemplatesFS, csc.TemplatesFS, ab.TemplatesFS, subscription.TemplatesFS)
	r.SetWorkspaceFormSigner(stubSigner{})
	if err := r.Init(); err != nil {
		t.Fatalf("template set does not parse: %v", err)
	}
	w := httptest.NewRecorder()
	if err := r.Render(w, name, data); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return w.Body.String()
}

func render(t *testing.T, name string, data any) string {
	t.Helper()
	out := renderRaw(t, name, data)
	// C15: no inline script / style in HTMX fragments. The external centymo delegation
	// <script src=...> include is the established precedent (CSP-safe) and is allowed.
	for _, m := range inlineScript.FindAllString(out, -1) {
		if !strings.Contains(m, "src=") {
			t.Errorf("%s renders an inline script: %s", name, m)
		}
	}
	// The pyeza table component (framework-owned) emits its own cell styles; only the
	// hand-written fragments are held to the no-inline-style rule.
	if !strings.Contains(out, "data-table") && !strings.Contains(out, `class="table`) && strings.Contains(out, " style=") {
		t.Errorf("%s renders an inline style", name)
	}
	if strings.Contains(out, "onclick=") {
		t.Errorf("%s renders onclick", name)
	}
	return out
}

func wantAll(t *testing.T, out string, subs ...string) {
	t.Helper()
	for _, s := range subs {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestCostSourceComponentDrawerRenders(t *testing.T) {
	l := csc.DefaultLabels()
	out := render(t, "cost-source-component-drawer-form", &cscform.Data{
		FormAction: "/action/cost-source-component/add/e1", ExpenditureID: "e1", Currency: "PHP",
		Kinds:       []cscform.Option{{Value: "COST_SOURCE_COMPONENT_KIND_ENERGY", Label: "Energy", Selected: true}},
		TaxFacts:    []cscform.Option{{Value: "COST_TAX_FACT_VATABLE", Label: "Taxable"}},
		BillLines:   []cscform.Option{{Value: "li1", Label: "Line 1"}},
		SubmitLabel: l.Buttons.Add, Labels: l, CommonLabels: pyeza.CommonLabels{},
	})
	wantAll(t, out, `data-testid="cost-source-component-form"`, `data-testid="cost-source-component-kind"`, `data-testid="cost-source-component-amount"`,
		`data-testid="cost-source-component-service-from"`, `data-testid="cost-source-component-line-item"`, `data-testid="cost-source-component-save"`,
		`hx-post="/action/cost-source-component/add/e1"`, l.Buttons.Add)
	// m4: the parent is pinned by the path, never posted back
	if strings.Contains(out, `name="expenditure_id"`) {
		t.Error("the drawer must not post a hidden expenditure_id")
	}
}

func TestAllocationDrawerAndPreviewRender(t *testing.T) {
	l := ab.DefaultLabels()
	c := pyeza.CommonLabels{}
	c.Errors.PermissionDenied = "no permission"
	preview := &abform.Preview{Currency: "PHP", Total: "1000.00", Source: "1000.00", Balanced: true, Lines: []abform.PreviewLine{
		{KindLabel: "Recoverable", Name: "Lease s1", Weight: "600", Percent: "60.00%", Amount: "600.00"},
	}}
	data := &abform.Data{
		FormAction: "/action/allocation-batch/allocate/c1", PreviewURL: "/action/allocation-batch/preview/c1", PublishURL: "/action/allocation-batch/publish/c1",
		ComponentID: "c1", ComponentLabel: "Energy", ServicePeriod: "2026-09-01 to 2026-10-01", SourceAmount: "1000.00", Currency: "PHP",
		Denominator: "1000", CanPublish: true, HasDraft: true, Revision: "Revision 1", Preview: preview, Labels: l, CommonLabels: c,
		Rows: []abform.Row{
			{Kind: "ALLOCATION_SHARE_KIND_RECOVERABLE", KindLabel: "Recoverable", SubscriptionID: "s1", ClientID: "cl1", Name: "Lease s1", Slug: "s1", Numerator: "600"},
			{Kind: "ALLOCATION_SHARE_KIND_OWN_USE", KindLabel: "Own use", Name: "Own use", Slug: "own-use", Numerator: "0"},
		},
	}
	out := render(t, "allocation-batch-drawer-form", data)
	wantAll(t, out, `data-testid="allocation-numerator-s1"`, `data-testid="allocation-numerator-own-use"`, `id="allocation-denominator"`, `data-testid="allocation-denominator-group"`,
		`data-testid="allocation-preview"`, `data-testid="allocation-save-draft"`, `data-testid="allocation-publish"`,
		`hx-post="/action/allocation-batch/publish/c1"`, `hx-include="#allocation-weights-form"`, `data-testid="allocation-publish-form"`, `hx-confirm=`, `data-confirm-title=`, `data-testid="allocation-preview-total"`,
		`name="share_kind" value="ALLOCATION_SHARE_KIND_RECOVERABLE"`, `name="subscription_id" value="s1"`, `1000.00`)
	// no publish permission: the button is disabled and carries no hx-post
	data.CanPublish = false
	out = render(t, "allocation-batch-drawer-form", data)
	if strings.Contains(out, `hx-post="/action/allocation-batch/publish/c1"`) || !strings.Contains(out, `disabled aria-disabled="true"`) {
		t.Error("publish must be disabled without the permission")
	}
	// the preview panel renders a refusal in place of the table (plain button request: always a 200 fragment)
	out = render(t, "allocation-batch-preview", &abform.PreviewPage{Labels: l, Error: "The shares do not add up."})
	wantAll(t, out, `data-testid="allocation-preview-error"`, "The shares do not add up.")
	if strings.Contains(out, "allocation-preview-table") {
		t.Error("an error panel must not render the table")
	}
	// nobody can share the cost: empty state, no weight inputs
	data.Rows = nil
	out = render(t, "allocation-batch-drawer-form", data)
	wantAll(t, out, `data-testid="allocation-no-participants"`)
	if strings.Contains(out, `data-testid="allocation-participants-table"`) {
		t.Error("participants table rendered with no participants")
	}
	// read-only view
	out = render(t, "allocation-batch-view", &abform.ViewData{ComponentLabel: "Energy", ServicePeriod: "p", Revision: "Revision 1", Status: "Published", Preview: preview, Labels: l, CommonLabels: c})
	wantAll(t, out, `data-testid="allocation-view"`, `data-testid="allocation-view-status"`, `data-testid="allocation-preview-line"`)
}

func TestRecoverableCostsTabRenders(t *testing.T) {
	l := csc.DefaultLabels()
	perms := types.NewUserPermissions([]string{"expenditure:update"})
	claim := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION
	tbl := csctable.Build([]*costsourcecomponentpb.CostSourceComponent{
		{Id: "open", Amount: 1000, Currency: "PHP"}, {Id: "done", Amount: 2000, Currency: "PHP", ClaimKind: &claim},
	}, csctable.Input{ExpenditureID: "e1", Currency: "PHP", Labels: l, Routes: csc.DefaultRoutes(), AllocateURL: ab.AllocateURL, ViewAllocationURL: ab.ViewURL, Perms: perms, NoPermission: "no permission"})
	out := render(t, "expense-tab-recoverable-costs", &expdetail.PageData{RecoverableCosts: &expdetail.RecoverableCostsData{
		Labels: l, Table: tbl, AddURL: "/action/cost-source-component/add/e1", CanAdd: true, Currency: "PHP", ComponentsSum: "30.00", BillTotal: "100.00", Difference: "70.00",
	}})
	wantAll(t, out, `data-testid="cost-source-component-add"`, `hx-get="/action/cost-source-component/add/e1"`, `data-testid="recoverable-costs-difference"`,
		`data-testid="cost-source-component-allocate-open"`, `data-testid="cost-source-component-view-allocation-done"`, `data-id="open"`)
	// a failed read shows an alert, never an empty-state claim
	out = render(t, "expense-tab-recoverable-costs", &expdetail.PageData{RecoverableCosts: &expdetail.RecoverableCostsData{Labels: l, Failed: true, CanAdd: false, AddTooltip: "no permission"}})
	wantAll(t, out, `data-testid="recoverable-costs-error"`, `disabled aria-disabled="true"`)
	// tab not wired: renders nothing
	if out := render(t, "expense-tab-recoverable-costs", &expdetail.PageData{}); strings.TrimSpace(out) != "" {
		t.Errorf("unwired tab rendered %q", out)
	}
}

// R5 B1: the Preview trigger sits in its OWN form carrying a signature bound to PreviewURL. A trigger
// outside a form that hx-includes the weights form would send that form's FormAction-bound signature
// and the action guard would answer 409.
func TestAllocationPreviewTriggerCarriesItsOwnPathBoundSignature(t *testing.T) {
	l := ab.DefaultLabels()
	c := pyeza.CommonLabels{}
	c.Errors.MissingPermission = "Missing permission: %s"
	data := &abform.Data{
		FormAction: "/action/allocation-batch/allocate/c1", PreviewURL: "/action/allocation-batch/preview/c1", PublishURL: "/action/allocation-batch/publish/c1",
		WorkspaceID: "ws1", ComponentID: "c1", Denominator: "1000", CanPublish: false, PublishTooltip: "Missing permission: allocation_batch:publish",
		Labels: l, CommonLabels: c,
		Rows: []abform.Row{{Kind: "ALLOCATION_SHARE_KIND_RECOVERABLE", KindLabel: "Recoverable", SubscriptionID: "s1", ClientID: "cl1", Name: "Lease s1", Slug: "s1", Numerator: "600"}},
	}
	out := render(t, "allocation-batch-drawer-form", data)
	start := strings.LastIndex(out[:strings.Index(out, `data-testid="allocation-preview-form"`)+1], "<form")
	if start < 0 {
		t.Fatal("the Preview trigger must live in its own <form>")
	}
	end := strings.Index(out[start:], "</form>")
	formHTML := out[start : start+end]
	wantAll(t, formHTML, `hx-post="/action/allocation-batch/preview/c1"`, `hx-include="#allocation-weights-form"`,
		`value="sig-for-ws1-on-/action/allocation-batch/preview/c1"`, `data-testid="allocation-preview"`)
	// W5 B-W5-2: the request must be issued by the button, not the <form> — pyeza sheet.js closes the
	// drawer on any successful request whose element is a FORM inside the sheet.
	formTag := formHTML[:strings.Index(formHTML, ">")+1]
	if strings.Contains(formTag, "hx-post") {
		t.Errorf("preview <form> tag must not carry hx-post (drawer would close): %s", formTag)
	}
	// M2/m3: the disabled Publish tooltip names the missing code
	wantAll(t, out, `title="Missing permission: allocation_batch:publish"`)
}

// R5 M3: the expenditure detail Add button opens the sheet through data-hx-on (the page loads no
// centymo delegation script).
func TestRecoverableCostsAddButtonUsesSheetOpenNotDelegation(t *testing.T) {
	l := csc.DefaultLabels()
	out := render(t, "expense-tab-recoverable-costs", &expdetail.PageData{RecoverableCosts: &expdetail.RecoverableCostsData{
		Labels: l, Currency: "PHP", AddURL: "/action/cost-source-component/add/e1", CanAdd: true,
		Table: csctable.Build(nil, csctable.Input{ExpenditureID: "e1", Labels: l, Routes: csc.DefaultRoutes(), Perms: types.NewUserPermissions(nil)}),
	}})
	i := strings.Index(out, `data-testid="cost-source-component-add"`)
	if i < 0 {
		t.Fatal("add button missing")
	}
	btn := out[i : i+strings.Index(out[i:], ">")]
	if !strings.Contains(btn, `data-hx-on="sheet-open"`) || strings.Contains(btn, "data-lf-sheet=") {
		t.Errorf("add button = %s", btn)
	}
}
