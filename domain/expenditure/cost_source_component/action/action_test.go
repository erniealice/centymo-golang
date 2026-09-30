package action

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	"github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/form"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
)

type coded struct{ code string }

func (c coded) Error() string     { return "refused: " + c.code }
func (c coded) ErrorCode() string { return c.code }

func common() pyeza.CommonLabels {
	c := pyeza.CommonLabels{}
	c.Errors.PermissionDenied = "no permission"
	return c
}

func ctx(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func req(method, target, body string, pv ...string) *view.ViewContext {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		r = httptest.NewRequest(method, target, nil)
	}
	for i := 0; i+1 < len(pv); i += 2 {
		r.SetPathValue(pv[i], pv[i+1])
	}
	return &view.ViewContext{Request: r, CurrentPath: r.URL.Path}
}

type fake struct {
	created, updated *costsourcecomponentpb.CostSourceComponent
	deleted          string
	current          *costsourcecomponentpb.CostSourceComponent
	err              error
	parentReads      int

	// bill-line reads (B3): every line of every expenditure, filtered only when the request filters
	lineItems []*expenditurelineitempb.ExpenditureLineItem
	lineReq   *expenditurelineitempb.ListExpenditureLineItemsRequest
	// workspace is the set of expenditure ids the actor's workspace owns (ParentExists); nil = all
	workspace map[string]bool
}

func (f *fake) uc() *csc.UseCases {
	return &csc.UseCases{
		CreateCostSourceComponent: func(_ context.Context, r *costsourcecomponentpb.CreateCostSourceComponentRequest) (*costsourcecomponentpb.CreateCostSourceComponentResponse, error) {
			f.created = r.GetData()
			return &costsourcecomponentpb.CreateCostSourceComponentResponse{Success: true}, f.err
		},
		UpdateCostSourceComponent: func(_ context.Context, r *costsourcecomponentpb.UpdateCostSourceComponentRequest) (*costsourcecomponentpb.UpdateCostSourceComponentResponse, error) {
			f.updated = r.GetData()
			return &costsourcecomponentpb.UpdateCostSourceComponentResponse{Success: true}, f.err
		},
		DeleteCostSourceComponent: func(_ context.Context, r *costsourcecomponentpb.DeleteCostSourceComponentRequest) (*costsourcecomponentpb.DeleteCostSourceComponentResponse, error) {
			f.deleted = r.GetData().GetId()
			return &costsourcecomponentpb.DeleteCostSourceComponentResponse{Success: true}, f.err
		},
		ListExpenditureLineItems: func(_ context.Context, r *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error) {
			f.lineReq = r
			var out []*expenditurelineitempb.ExpenditureLineItem
			for _, li := range f.lineItems {
				// like the adapter: a Filters expenditure_id narrows the list server-side
				if fl := r.GetFilters().GetFilters(); len(fl) == 1 && li.GetExpenditureId() != fl[0].GetStringFilter().GetValue() {
					continue
				}
				out = append(out, li)
			}
			return &expenditurelineitempb.ListExpenditureLineItemsResponse{Data: out}, nil
		},
		ReadCostSourceComponent: func(_ context.Context, r *costsourcecomponentpb.ReadCostSourceComponentRequest) (*costsourcecomponentpb.ReadCostSourceComponentResponse, error) {
			if f.current == nil {
				return &costsourcecomponentpb.ReadCostSourceComponentResponse{}, nil
			}
			return &costsourcecomponentpb.ReadCostSourceComponentResponse{Data: []*costsourcecomponentpb.CostSourceComponent{f.current}}, nil
		},
	}
}

func deps(f *fake) *Deps {
	return &Deps{Routes: csc.DefaultRoutes(), Labels: csc.DefaultLabels(), CommonLabels: common(), UseCases: f.uc(),
		ReadParent: func(_ context.Context, id string) (string, bool) {
			f.parentReads++
			return "PHP", f.workspace == nil || f.workspace[id]
		}}
}

const validForm = "component_kind=COST_SOURCE_COMPONENT_KIND_ENERGY&description=Common+electricity&basis_quantity=1250.5&basis_unit=kWh&amount=1250.75&currency=PHP&tax_fact=COST_TAX_FACT_VATABLE&service_from=2026-09-01&service_to=2026-10-01"

func TestAddRefusedWithoutPermission(t *testing.T) {
	f := &fake{}
	for _, m := range []string{http.MethodGet, http.MethodPost} {
		res := NewAddAction(deps(f)).Handle(ctx("expenditure:read"), req(m, "/x", validForm, "id", "e1"))
		if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != "no permission" {
			t.Fatalf("%s: status %d headers %v", m, res.StatusCode, res.Headers)
		}
	}
	if f.created != nil {
		t.Fatal("create ran without permission")
	}
}

func TestAddDrawerAndCreateConvertsToCentavos(t *testing.T) {
	f := &fake{}
	res := NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "e1"))
	if res.Template != "cost-source-component-drawer-form" {
		t.Fatalf("GET template = %q", res.Template)
	}
	d := res.Data.(*form.Data)
	if d.ExpenditureID != "e1" || d.Currency != "PHP" || len(d.Kinds) != 4 || len(d.TaxFacts) != 5 {
		t.Fatalf("drawer data = %+v", d)
	}
	res = NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "e1"))
	if res.StatusCode != http.StatusOK || !strings.Contains(res.Headers["HX-Trigger"], "cost-source-components-table") {
		t.Fatalf("POST status %d headers %v", res.StatusCode, res.Headers)
	}
	c := f.created
	if c == nil || c.GetAmount() != 125075 || c.GetExpenditureId() != "e1" || c.GetCurrency() != "PHP" ||
		c.GetBasisQuantityScaled() != 12505 || c.GetBasisScale() != 1 || c.GetServiceFrom() != "2026-09-01" ||
		c.GetComponentKind() != costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_ENERGY ||
		c.GetTaxFact() != costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_VATABLE {
		t.Fatalf("created = %v", c)
	}
}

func TestAddRejectsBadForm(t *testing.T) {
	f := &fake{}
	for _, body := range []string{
		"component_kind=&amount=10&service_from=2026-09-01&service_to=2026-10-01",
		"component_kind=COST_SOURCE_COMPONENT_KIND_WATER&amount=0&service_from=2026-09-01&service_to=2026-10-01",
		"component_kind=COST_SOURCE_COMPONENT_KIND_WATER&amount=10&service_from=&service_to=2026-10-01",
	} {
		res := NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", body, "id", "e1"))
		if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != csc.DefaultLabels().Errors.FormInvalid {
			t.Errorf("body %q: status %d headers %v", body, res.StatusCode, res.Headers)
		}
	}
	if f.created != nil {
		t.Fatal("create ran on an invalid form")
	}
}

// The use-case refusal code maps to its label (claimed / reference_invalid).
func TestRefusalCodesMapToLabels(t *testing.T) {
	l := csc.DefaultLabels()
	for code, want := range map[string]string{csc.ErrClaimed: l.Errors.Claimed, csc.ErrReferenceInvalid: l.Errors.ReferenceInvalid} {
		f := &fake{err: coded{code}}
		res := NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "e1"))
		if res.Headers["HX-Error-Message"] != want {
			t.Errorf("code %s: message %q, want %q", code, res.Headers["HX-Error-Message"], want)
		}
	}
}

func TestEditRefusesClaimedAndDeniedAndUnwired(t *testing.T) {
	claim := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION
	f := &fake{current: &costsourcecomponentpb.CostSourceComponent{Id: "c1", ExpenditureId: "e1", Amount: 100, ClaimKind: &claim}}
	res := NewEditAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "c1"))
	if res.Headers["HX-Error-Message"] != csc.DefaultLabels().Errors.Claimed {
		t.Fatalf("claimed GET = %d %v", res.StatusCode, res.Headers)
	}
	res = NewEditAction(deps(f)).Handle(ctx("expenditure:read"), req(http.MethodPost, "/x", validForm, "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" || f.updated != nil {
		t.Fatalf("denied POST = %d %v (updated=%v)", res.StatusCode, res.Headers, f.updated)
	}
	d := deps(f)
	d.UseCases = &csc.UseCases{}
	res = NewEditAction(d).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "c1"))
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unwired POST status = %d", res.StatusCode)
	}
}

func TestEditPrefillsAndUpdates(t *testing.T) {
	qty, scale := int64(12505), int32(1)
	desc := "Common electricity"
	from, to := "2026-09-01", "2026-10-01"
	f := &fake{current: &costsourcecomponentpb.CostSourceComponent{
		Id: "c1", ExpenditureId: "e1", Amount: 125075, Currency: "PHP", Description: &desc,
		BasisQuantityScaled: &qty, BasisScale: &scale, ServiceFrom: &from, ServiceTo: &to,
		ComponentKind: costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_ENERGY,
	}}
	res := NewEditAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "c1"))
	d, ok := res.Data.(*form.Data)
	if !ok || !d.IsEdit || d.Amount != "1250.75" || d.Quantity != "1250.5" || d.ServiceFrom != from {
		t.Fatalf("edit drawer = %+v", res.Data)
	}
	res = NewEditAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "c1"))
	if res.StatusCode != http.StatusOK || f.updated.GetId() != "c1" || f.updated.GetAmount() != 125075 {
		t.Fatalf("update = %d %v", res.StatusCode, f.updated)
	}
}

func TestDeleteRefusesWithoutPermissionAndMapsClaimed(t *testing.T) {
	f := &fake{}
	res := NewDeleteAction(deps(f)).Handle(ctx("expenditure:read"), req(http.MethodPost, "/x", "", "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" || f.deleted != "" {
		t.Fatalf("denied delete = %v deleted=%q", res.Headers, f.deleted)
	}
	f.err = coded{csc.ErrClaimed}
	res = NewDeleteAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", "", "id", "c1"))
	if res.Headers["HX-Error-Message"] != csc.DefaultLabels().Errors.Claimed {
		t.Fatalf("claimed delete = %v", res.Headers)
	}
	f.err = nil
	res = NewDeleteAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", "", "id", "c9"))
	if res.StatusCode != http.StatusOK || f.deleted != "c9" {
		t.Fatalf("delete = %d %q", res.StatusCode, f.deleted)
	}
}

// R5 B3: the Add drawer GET proves the URL parent in the workspace first (404) and lists its bill lines
// with a server-side expenditure_id filter, never another expenditure's rows.
func TestAddDrawerReadsParentFirstAndFiltersBillLines(t *testing.T) {
	f := &fake{
		workspace: map[string]bool{"e1": true},
		lineItems: []*expenditurelineitempb.ExpenditureLineItem{
			{Id: "l1", ExpenditureId: "e1", Description: "Mine"},
			{Id: "l2", ExpenditureId: "foreign", Description: "Foreign tenant line"},
		},
	}
	res := NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "foreign"))
	if res.StatusCode != http.StatusNotFound || f.lineReq != nil {
		t.Fatalf("foreign parent = %d listed=%v", res.StatusCode, f.lineReq != nil)
	}
	d := deps(f)
	d.ReadParent = nil
	if res := NewAddAction(d).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "e1")); res.StatusCode != http.StatusNotFound {
		t.Fatalf("nil ReadParent = %d", res.StatusCode)
	}
	res = NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "e1"))
	got := res.Data.(*form.Data)
	fl := f.lineReq.GetFilters().GetFilters()
	if len(fl) != 1 || fl[0].GetField() != "expenditure_id" || fl[0].GetStringFilter().GetValue() != "e1" {
		t.Fatalf("filters = %v", fl)
	}
	if len(got.BillLines) != 2 || got.BillLines[1].Value != "l1" {
		t.Fatalf("bill lines = %+v", got.BillLines)
	}
}

// R5 M1: a coded permission_denied refusal maps to the permission message, not the generic error.
func TestPermissionDeniedRefusalMapsToPermissionMessage(t *testing.T) {
	f := &fake{err: coded{csc.ErrPermissionDenied}}
	for name, res := range map[string]view.ViewResult{
		"add":    NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "e1")),
		"delete": NewDeleteAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", "", "id", "c1")),
	} {
		if res.Headers["HX-Error-Message"] != "no permission" {
			t.Errorf("%s: %v", name, res.Headers)
		}
	}
	f = &fake{current: &costsourcecomponentpb.CostSourceComponent{Id: "c1", ExpenditureId: "e1", Amount: 100}, err: coded{csc.ErrPermissionDenied}}
	if res := NewEditAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm, "id", "c1")); res.Headers["HX-Error-Message"] != "no permission" {
		t.Errorf("edit: %v", res.Headers)
	}
}

// R5 m4: Edit never takes the parent from the form (update pins the stored expenditure).
func TestEditIgnoresPostedExpenditureID(t *testing.T) {
	f := &fake{current: &costsourcecomponentpb.CostSourceComponent{Id: "c1", ExpenditureId: "e1", Amount: 100}}
	NewEditAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodPost, "/x", validForm+"&expenditure_id=hijack", "id", "c1"))
	if f.updated == nil || f.updated.GetExpenditureId() != "" {
		t.Fatalf("updated = %v", f.updated)
	}
}

// The Add drawer GET reads the parent expenditure exactly once (existence + currency).
func TestAddDrawerReadsTheParentExpenditureOncePerRequest(t *testing.T) {
	f := &fake{}
	res := NewAddAction(deps(f)).Handle(ctx("expenditure:update"), req(http.MethodGet, "/x", "", "id", "e1"))
	if res.StatusCode != http.StatusOK || f.parentReads != 1 {
		t.Fatalf("status=%d parentReads=%d", res.StatusCode, f.parentReads)
	}
}
