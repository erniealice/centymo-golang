package detail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	sibProductPricePlan "github.com/erniealice/centymo-golang/domain/subscription/product_price_plan"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	chargepolicypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

type pickerFn = func(context.Context, *chargepolicypb.ListPickerChargePoliciesRequest) (*chargepolicypb.ListPickerChargePoliciesResponse, error)

func pickerOf(policies ...*chargepolicypb.ChargePolicy) pickerFn {
	return func(context.Context, *chargepolicypb.ListPickerChargePoliciesRequest) (*chargepolicypb.ListPickerChargePoliciesResponse, error) {
		return &chargepolicypb.ListPickerChargePoliciesResponse{Data: policies, Success: true}, nil
	}
}

func failingPicker() pickerFn {
	return func(context.Context, *chargepolicypb.ListPickerChargePoliciesRequest) (*chargepolicypb.ListPickerChargePoliciesResponse, error) {
		return nil, errors.New("permission denied")
	}
}

var testFormLabels = sibProductPricePlan.FormLabels{ChargePolicyNone: "None", ChargePolicyUnavailable: "Unavailable charge policy"}

func policy(id, name string) *chargepolicypb.ChargePolicy {
	return &chargepolicypb.ChargePolicy{Id: id, Name: name}
}

// AC-CP-06 (view half): the section is offered only for RECURRING/CONTRACT,
// non-TOTAL_PACKAGE parents and only when the picker is wired.
func TestLoadChargePolicyOptionsGate(t *testing.T) {
	deps := &DetailViewDeps{ListPickerChargePolicies: pickerOf(policy("cp1", "Utilities"))}
	cases := []struct {
		name   string
		kind   priceplanpb.BillingKind
		basis  priceplanpb.AmountBasis
		parent bool
		want   bool
	}{
		{"recurring", priceplanpb.BillingKind_BILLING_KIND_RECURRING, priceplanpb.AmountBasis_AMOUNT_BASIS_PER_CYCLE, true, true},
		{"contract", priceplanpb.BillingKind_BILLING_KIND_CONTRACT, priceplanpb.AmountBasis_AMOUNT_BASIS_PER_CYCLE, true, true},
		{"one_time", priceplanpb.BillingKind_BILLING_KIND_ONE_TIME, priceplanpb.AmountBasis_AMOUNT_BASIS_PER_CYCLE, true, false},
		{"milestone", priceplanpb.BillingKind_BILLING_KIND_MILESTONE, priceplanpb.AmountBasis_AMOUNT_BASIS_PER_CYCLE, true, false},
		{"total_package", priceplanpb.BillingKind_BILLING_KIND_RECURRING, priceplanpb.AmountBasis_AMOUNT_BASIS_TOTAL_PACKAGE, true, false},
		{"no_parent", 0, 0, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var parent *priceplanpb.PricePlan
			if tc.parent {
				parent = &priceplanpb.PricePlan{BillingKind: tc.kind, AmountBasis: tc.basis}
			}
			got, _, _ := loadChargePolicyOptions(context.Background(), deps, parent, "", testFormLabels)
			if got != tc.want {
				t.Fatalf("enabled = %v, want %v", got, tc.want)
			}
		})
	}
	t.Run("unwired", func(t *testing.T) {
		parent := &priceplanpb.PricePlan{BillingKind: priceplanpb.BillingKind_BILLING_KIND_RECURRING}
		if got, _, _ := loadChargePolicyOptions(context.Background(), &DetailViewDeps{}, parent, "", testFormLabels); got {
			t.Fatal("section must be omitted when the picker is not wired")
		}
	})
}

func TestLoadChargePolicyOptionsNoneAndPreserved(t *testing.T) {
	parent := &priceplanpb.PricePlan{BillingKind: priceplanpb.BillingKind_BILLING_KIND_RECURRING}
	deps := &DetailViewDeps{ListPickerChargePolicies: pickerOf(policy("cp1", "Utilities"))}

	_, opts, _ := loadChargePolicyOptions(context.Background(), deps, parent, "cp1", testFormLabels)
	if len(opts) != 2 || opts[0].Value != "" || opts[0].Label != "None" || opts[0].Selected {
		t.Fatalf("expected leading unselected None option, got %+v", opts)
	}
	if opts[1].Value != "cp1" || opts[1].Label != "Utilities" || !opts[1].Selected {
		t.Fatalf("stored selection must be preselected: %+v", opts[1])
	}

	// A stored policy that later retired is absent from the picker; editing must
	// keep it rather than silently clear it.
	_, opts, _ = loadChargePolicyOptions(context.Background(), deps, parent, "retired-1", testFormLabels)
	last := opts[len(opts)-1]
	if last.Value != "retired-1" || !last.Selected || last.Label != "Unavailable charge policy" {
		t.Fatalf("retired stored policy must be preserved as selected: %+v", opts)
	}

	_, opts, _ = loadChargePolicyOptions(context.Background(), deps, parent, "", testFormLabels)
	if !opts[0].Selected {
		t.Fatalf("None must be selected when no policy is stored: %+v", opts)
	}
}

// m4: an unreadable picker (e.g. no charge_policy:list) is logged, not swallowed: the section
// still renders (disabled), the stored selection keeps a label (never the raw id), and the table
// shows the label instead of the UUID.
func TestChargePolicyPickerUnreadable(t *testing.T) {
	parent := &priceplanpb.PricePlan{BillingKind: priceplanpb.BillingKind_BILLING_KIND_RECURRING}
	deps := &DetailViewDeps{ListPickerChargePolicies: failingPicker()}
	enabled, opts, unreadable := loadChargePolicyOptions(context.Background(), deps, parent, "cp9", testFormLabels)
	if !enabled || !unreadable {
		t.Fatalf("section must render disabled: enabled=%v unreadable=%v", enabled, unreadable)
	}
	if len(opts) != 2 || opts[1].Label != "Unavailable charge policy" || opts[1].Value != "cp9" {
		t.Fatalf("stored selection must keep a plain label: %+v", opts)
	}
	deps.CommonLabels.Errors.MissingPermission = "Missing permission: %s"
	if got := chargePolicyDisabledReason(deps, true); got != "Missing permission: charge_policy:list" {
		t.Fatalf("tooltip = %q", got)
	}
	if chargePolicyDisabledReason(deps, false) != "" {
		t.Fatal("no reason when readable")
	}

	cp := "cp9"
	tbl := buildProductPricesTable(context.Background(), &DetailViewDeps{
		ListProductPricePlans: func(context.Context, *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error) {
			return &productpriceplanpb.ListProductPricePlansResponse{Data: []*productpriceplanpb.ProductPricePlan{{Id: "a", PricePlanId: "pp1", ChargePolicyId: &cp}}}, nil
		},
		ListPickerChargePolicies: failingPicker(),
		ProductPricePlanLabels:   sibProductPricePlan.Labels{Form: testFormLabels},
	}, "pp1", "plan1")
	if got := tbl.Rows[0].Cells[len(tbl.Rows[0].Cells)-1].Value; got != "Unavailable charge policy" {
		t.Fatalf("table cell = %q, want the plain label (never the UUID)", got)
	}
}

func TestMarkupPercentDisplay(t *testing.T) {
	five := int32(500)
	if markupPercentDisplay(nil) != "0" || markupPercentDisplay(&five) != "5" {
		t.Fatal("markup percent display mismatch")
	}
}

func actionDeps(kind priceplanpb.BillingKind, createErr error, captured **productpriceplanpb.ProductPricePlan) *DetailViewDeps {
	return &DetailViewDeps{
		ReadPricePlan: func(_ context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error) {
			return &priceplanpb.ReadPricePlanResponse{Data: []*priceplanpb.PricePlan{{
				Id: req.GetData().GetId(), BillingKind: kind, BillingCurrency: "PHP",
			}}}, nil
		},
		CreateProductPricePlan: func(_ context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error) {
			*captured = req.GetData()
			return &productpriceplanpb.CreateProductPricePlanResponse{}, createErr
		},
	}
}

func postAdd(t *testing.T, deps *DetailViewDeps, form url.Values) view.ViewResult {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "pp1")
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"price_plan:read", "product_price_plan:create"}))
	return NewProductPriceAddAction(deps).Handle(ctx, &view.ViewContext{Request: req})
}

func TestAddActionPostsChargePolicyID(t *testing.T) {
	var got *productpriceplanpb.ProductPricePlan
	deps := actionDeps(priceplanpb.BillingKind_BILLING_KIND_RECURRING, nil, &got)
	res := postAdd(t, deps, url.Values{
		"product_plan_id": {"pl1"}, "price": {"10.00"}, "currency": {"PHP"},
		"billing_treatment": {"BILLING_TREATMENT_USAGE_BASED"}, "charge_policy_id": {" cp1 "},
	})
	if res.StatusCode != 0 && res.StatusCode >= 400 {
		t.Fatalf("unexpected failure: %+v", res.Headers)
	}
	if got == nil || got.GetChargePolicyId() != "cp1" {
		t.Fatalf("charge_policy_id not forwarded: %v", got)
	}
	if got.MarkupBps != nil {
		t.Fatal("markup_bps must never be posted in S1")
	}
}

func TestAddActionOmitsChargePolicyWhenBlank(t *testing.T) {
	var got *productpriceplanpb.ProductPricePlan
	deps := actionDeps(priceplanpb.BillingKind_BILLING_KIND_RECURRING, nil, &got)
	postAdd(t, deps, url.Values{"product_plan_id": {"pl1"}, "price": {"10"}, "currency": {"PHP"}, "charge_policy_id": {""}})
	if got == nil || got.ChargePolicyId != nil {
		t.Fatalf("blank selection must leave charge_policy_id unset: %v", got)
	}
}

// AC-CP-06: a forged POST carrying charge_policy_id on a ONE_TIME plan is
// forwarded untouched (the use-case guard is the single authority) and its
// refusal is surfaced to the operator, not swallowed.
func TestAddActionForgedChargePolicySurfacesGuardRefusal(t *testing.T) {
	var got *productpriceplanpb.ProductPricePlan
	guard := errors.New("product_price_plan: charge_policy requires a RECURRING or CONTRACT price plan")
	deps := actionDeps(priceplanpb.BillingKind_BILLING_KIND_ONE_TIME, guard, &got)
	res := postAdd(t, deps, url.Values{
		"product_plan_id": {"pl1"}, "price": {"10"}, "currency": {"PHP"},
		"billing_treatment": {"BILLING_TREATMENT_USAGE_BASED"}, "charge_policy_id": {"cp1"},
	})
	if got == nil || got.GetChargePolicyId() != "cp1" {
		t.Fatalf("view must not silently drop the forged id: %v", got)
	}
	if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != guard.Error() {
		t.Fatalf("guard refusal not surfaced: status=%d headers=%v", res.StatusCode, res.Headers)
	}
}

func TestProductPricesTableChargePolicyColumn(t *testing.T) {
	cp := "cp1"
	list := func(_ context.Context, _ *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error) {
		return &productpriceplanpb.ListProductPricePlansResponse{Data: []*productpriceplanpb.ProductPricePlan{
			{Id: "a", PricePlanId: "pp1", ChargePolicyId: &cp},
			{Id: "b", PricePlanId: "pp1"},
		}}, nil
	}
	off := buildProductPricesTable(context.Background(), &DetailViewDeps{ListProductPricePlans: list}, "pp1", "plan1")
	for _, c := range off.Columns {
		if c.Key == "charge_policy" {
			t.Fatal("column must be absent when the picker is not wired (UC-PRES)")
		}
	}
	if len(off.Rows) != 2 || len(off.Rows[0].Cells) != 3 {
		t.Fatalf("unwired table must keep 3 cells: %+v", off.Rows)
	}

	on := buildProductPricesTable(context.Background(), &DetailViewDeps{
		ListProductPricePlans:    list,
		ListPickerChargePolicies: pickerOf(policy("cp1", "Utilities")),
		ProductPricePlanLabels:   sibProductPricePlan.Labels{Columns: sibProductPricePlan.ColumnLabels{ChargePolicy: "Charge policy"}},
	}, "pp1", "plan1")
	if on.Columns[len(on.Columns)-1].Key != "charge_policy" || on.Columns[len(on.Columns)-1].Label != "Charge policy" {
		t.Fatalf("charge_policy column missing: %+v", on.Columns)
	}
	byID := map[string]string{}
	for _, r := range on.Rows {
		byID[r.ID] = r.Cells[len(r.Cells)-1].Value
	}
	if byID["a"] != "Utilities" || byID["b"] != "—" {
		t.Fatalf("cells = %v", byID)
	}
}

// Template gate: the section renders only when ChargePolicyEnabled, starts
// hidden unless USAGE_BASED, and ships the toggle script. Pyeza partials are
// stubbed; only the drawer's own markup is under test.
func TestDrawerTemplateChargePolicySection(t *testing.T) {
	src, err := os.ReadFile("../../product_price_plan/templates/product-price-plan-drawer-form.html")
	if err != nil {
		t.Fatal(err)
	}
	stubs := `{{define "form-section"}}<h3>{{.Title}}</h3>{{end}}` +
		`{{define "form-group"}}<x name="{{.Name}}" testid="{{.TestId}}" {{if .Disabled}}disabled{{end}}></x>{{end}}` +
		`{{define "ppp-parent-context"}}{{end}}{{define "ppp-fields"}}{{end}}` +
		`{{define "sheet-form-footer"}}{{end}}{{define "icon-info"}}{{end}}`
	tmpl, err := template.New("t").Funcs(template.FuncMap{
		"dict": func(kv ...any) map[string]any {
			m := map[string]any{}
			for i := 0; i+1 < len(kv); i += 2 {
				m[kv[i].(string)] = kv[i+1]
			}
			return m
		},
		"actionForm": func(...any) template.HTML { return "" },
	}).Parse(stubs + string(src))
	if err != nil {
		t.Fatal(err)
	}
	render := func(d *ProductPricePlanFormData) string {
		var b bytes.Buffer
		if err := tmpl.ExecuteTemplate(&b, "product-price-plan-drawer-form", d); err != nil {
			t.Fatalf("render: %v", err)
		}
		return b.String()
	}

	off := render(&ProductPricePlanFormData{})
	if strings.Contains(off, "ppp-charge-policy") || strings.Contains(off, "billing_treatment") {
		t.Fatalf("drawer must be unchanged when the section is not enabled:\n%s", off)
	}

	hidden := render(&ProductPricePlanFormData{ChargePolicyEnabled: true})
	for _, want := range []string{`id="ppp-charge-policy-section"`, `data-lf-show-when="billing_treatment"`, `data-lf-show-values="BILLING_TREATMENT_USAGE_BASED" hidden>`, `testid="ppp-charge-policy"`, `testid="ppp-markup"`} {
		if !strings.Contains(hidden, want) {
			t.Errorf("missing %q in\n%s", want, hidden)
		}
	}
	if strings.Contains(hidden, "<script") || strings.Contains(hidden, "style=") {
		t.Fatalf("section must be CSP-safe (no inline script/style):\n%s", hidden)
	}
	shown := render(&ProductPricePlanFormData{ChargePolicyEnabled: true, ChargePolicyVisible: true})
	if strings.Contains(shown, `USAGE_BASED" hidden`) {
		t.Fatalf("USAGE_BASED line must render the section visible:\n%s", shown)
	}
}

// codedErr mimics the use case's guard refusals: a stable ErrorCode() and an English message.
type codedErr struct{ code, msg string }

func (e codedErr) Error() string     { return e.msg }
func (e codedErr) ErrorCode() string { return e.code }

// Every guard sentinel code maps to its localized label; the English message is never rendered.
func TestAddActionGuardRefusalRendersLocalizedLabel(t *testing.T) {
	labels := sibProductPricePlan.Labels{Errors: sibProductPricePlan.ErrorLabels{
		ChargePolicyNotUsageBased:    "L:not_usage_based",
		ChargePolicyBillingKind:      "L:billing_kind",
		ChargePolicyUnavailable:      "L:unavailable",
		ChargePolicyNotFound:         "L:not_found",
		ChargePolicyMarkupNotAllowed: "L:markup",
		ChargePolicyUnverifiable:     "L:unverifiable",
	}}
	for code, want := range map[string]string{
		"charge_policy_not_usage_based":    "L:not_usage_based",
		"charge_policy_billing_kind":       "L:billing_kind",
		"charge_policy_unavailable":        "L:unavailable",
		"charge_policy_not_found":          "L:not_found",
		"charge_policy_markup_not_allowed": "L:markup",
		"charge_policy_unverifiable":       "L:unverifiable",
	} {
		t.Run(code, func(t *testing.T) {
			var got *productpriceplanpb.ProductPricePlan
			// wrapped like the transactional update path does
			refusal := fmt.Errorf("update failed: %w", codedErr{code: code, msg: "ENGLISH " + code})
			deps := actionDeps(priceplanpb.BillingKind_BILLING_KIND_RECURRING, refusal, &got)
			deps.ProductPricePlanLabels = labels
			res := postAdd(t, deps, url.Values{"product_plan_id": {"pl1"}, "price": {"10"}, "currency": {"PHP"}, "charge_policy_id": {"cp1"}})
			if res.Headers["HX-Error-Message"] != want {
				t.Fatalf("want %q, got %q", want, res.Headers["HX-Error-Message"])
			}
		})
	}
	// unlabeled code and uncoded errors fall back to the raw text
	if m := (sibProductPricePlan.Labels{}).GuardErrorMessage(codedErr{code: "charge_policy_not_found", msg: "raw"}); m != "raw" {
		t.Fatalf("unlabeled fallback: %q", m)
	}
}

func postEdit(t *testing.T, deps *DetailViewDeps, form url.Values) view.ViewResult {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "pp1")
	req.SetPathValue("ppid", "line1")
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"price_plan:read", "product_price_plan:update"}))
	return NewProductPriceEditAction(deps).Handle(ctx, &view.ViewContext{Request: req})
}

func editDeps(stored string, captured **productpriceplanpb.ProductPricePlan) *DetailViewDeps {
	d := actionDeps(priceplanpb.BillingKind_BILLING_KIND_RECURRING, nil, new(*productpriceplanpb.ProductPricePlan))
	line := &productpriceplanpb.ProductPricePlan{Id: "line1", PricePlanId: "pp1", ProductPlanId: "pl1", BillingCurrency: "PHP",
		BillingTreatment: productpriceplanpb.BillingTreatment_BILLING_TREATMENT_USAGE_BASED}
	if stored != "" {
		line.ChargePolicyId = &stored
	}
	d.ListProductPricePlans = func(context.Context, *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error) {
		return &productpriceplanpb.ListProductPricePlansResponse{Data: []*productpriceplanpb.ProductPricePlan{line}}, nil
	}
	d.UpdateProductPricePlan = func(_ context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error) {
		*captured = req.GetData()
		return &productpriceplanpb.UpdateProductPricePlanResponse{}, nil
	}
	d.ListPickerChargePolicies = pickerOf(policy("cp1", "Utilities"))
	return d
}

func TestEditActionChargePolicyClear(t *testing.T) {
	base := url.Values{"product_plan_id": {"pl1"}, "price": {"10"}, "currency": {"PHP"}}
	with := func(kv map[string]string) url.Values {
		v := url.Values{}
		for k, x := range base {
			v[k] = x
		}
		for k, x := range kv {
			v.Set(k, x)
		}
		return v
	}
	t.Run("None on a policy-bound line posts the typed clear (present and empty)", func(t *testing.T) {
		var got *productpriceplanpb.ProductPricePlan
		postEdit(t, editDeps("cp1", &got), with(map[string]string{"billing_treatment": "BILLING_TREATMENT_USAGE_BASED", "charge_policy_id": ""}))
		if got == nil || got.ChargePolicyId == nil || got.GetChargePolicyId() != "" {
			t.Fatalf("want present-and-empty clear, got %+v", got)
		}
	})
	t.Run("switching off USAGE_BASED sends no view-side clear (the use case owns it, C12)", func(t *testing.T) {
		var got *productpriceplanpb.ProductPricePlan
		postEdit(t, editDeps("cp1", &got), with(map[string]string{"billing_treatment": "BILLING_TREATMENT_RECURRING"}))
		if got == nil || got.ChargePolicyId != nil {
			t.Fatalf("view must not derive the clear, got %+v", got)
		}
	})
	t.Run("a chosen policy is forwarded, an untouched line stays untouched", func(t *testing.T) {
		var got *productpriceplanpb.ProductPricePlan
		postEdit(t, editDeps("cp1", &got), with(map[string]string{"billing_treatment": "BILLING_TREATMENT_USAGE_BASED", "charge_policy_id": "cp2"}))
		if got.GetChargePolicyId() != "cp2" {
			t.Fatalf("chosen policy: %+v", got)
		}
		postEdit(t, editDeps("cp1", &got), with(map[string]string{"billing_treatment": "BILLING_TREATMENT_USAGE_BASED"}))
		if got.ChargePolicyId != nil {
			t.Fatalf("policy must be left alone when nothing changed: %+v", got)
		}
		postEdit(t, editDeps("", &got), with(map[string]string{"billing_treatment": "BILLING_TREATMENT_RECURRING", "charge_policy_id": ""}))
		if got.ChargePolicyId != nil {
			t.Fatalf("no clear when the line has no policy: %+v", got)
		}
	})
}
