package action

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	"github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch/form"
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
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

func req(method, body string, pv ...string) *view.ViewContext {
	r := httptest.NewRequest(method, "/x", strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for i := 0; i+1 < len(pv); i += 2 {
		r.SetPathValue(pv[i], pv[i+1])
	}
	return &view.ViewContext{Request: r, CurrentPath: r.URL.Path}
}

func all(codes ...string) []string {
	return append([]string{"allocation_batch:read", "allocation_batch:create", "allocation_batch:update", "allocation_batch:publish"}, codes...)
}

// fake is an in-memory espyna stand-in: the component under allocation, the agreement
// terms, an optional existing draft, and the recorded writes. Amounts of stored shares
// come from a plain proportional split (the real largest-remainder split is espyna's).
type fake struct {
	component *costsourcecomponentpb.CostSourceComponent
	terms     []*agreementlinetermpb.AgreementLineTerm
	batches   []*allocationbatchpb.AllocationBatch
	shares    map[string][]*allocationsharepb.AllocationShare

	created   *allocationbatchpb.CreateAllocationBatchRequest
	updated   *allocationbatchpb.UpdateAllocationBatchSharesRequest
	published []string
	writeErr  error
	pubErr    error
}

func newFake() *fake {
	from, to := "2026-09-01", "2026-10-01"
	return &fake{
		component: &costsourcecomponentpb.CostSourceComponent{Id: "c1", ExpenditureId: "e1", Amount: 100000, Currency: "PHP", ServiceFrom: &from, ServiceTo: &to,
			ComponentKind: costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_ENERGY},
		terms: []*agreementlinetermpb.AgreementLineTerm{
			{Id: "t1", SubscriptionId: "s1", ClientId: "cl1", EffectiveFrom: "2026-01-01", Active: true},
			{Id: "t2", SubscriptionId: "s2", ClientId: "cl2", EffectiveFrom: "2026-01-01", Active: true},
			{Id: "t3", SubscriptionId: "s3", ClientId: "cl3", EffectiveFrom: "2026-09-15", Active: true}, // starts mid-period: excluded
		},
		shares: map[string][]*allocationsharepb.AllocationShare{},
	}
}

func stored(in []*allocationsharepb.AllocationShare, total int64) []*allocationsharepb.AllocationShare {
	out := make([]*allocationsharepb.AllocationShare, len(in))
	for i, s := range in {
		den := s.GetBasisDenominator()
		out[i] = &allocationsharepb.AllocationShare{Id: "sh", ShareKind: s.GetShareKind(), SubscriptionId: s.SubscriptionId, ClientId: s.ClientId,
			BasisNumerator: s.GetBasisNumerator(), BasisDenominator: den, Amount: total * s.GetBasisNumerator() / den, SequenceOrder: int32(i + 1)}
	}
	return out
}

func (f *fake) uc() *ab.UseCases {
	return &ab.UseCases{
		CreateAllocationBatch: func(_ context.Context, r *allocationbatchpb.CreateAllocationBatchRequest) (*allocationbatchpb.CreateAllocationBatchResponse, error) {
			f.created = r
			if f.writeErr != nil {
				return nil, f.writeErr
			}
			b := &allocationbatchpb.AllocationBatch{Id: "b1", CostSourceComponentId: "c1", Revision: 1, Status: allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_DRAFT}
			return &allocationbatchpb.CreateAllocationBatchResponse{Data: []*allocationbatchpb.AllocationBatch{b}, Shares: stored(r.GetShares(), f.component.GetAmount()), Success: true}, nil
		},
		UpdateAllocationBatchShares: func(_ context.Context, r *allocationbatchpb.UpdateAllocationBatchSharesRequest) (*allocationbatchpb.UpdateAllocationBatchSharesResponse, error) {
			f.updated = r
			if f.writeErr != nil {
				return nil, f.writeErr
			}
			return &allocationbatchpb.UpdateAllocationBatchSharesResponse{Data: &allocationbatchpb.AllocationBatch{Id: r.GetAllocationBatchId(), Revision: 1}, Shares: stored(r.GetShares(), f.component.GetAmount()), Success: true}, nil
		},
		PublishAllocationBatch: func(_ context.Context, r *allocationbatchpb.PublishAllocationBatchRequest) (*allocationbatchpb.PublishAllocationBatchResponse, error) {
			f.published = append(f.published, r.GetAllocationBatchId())
			if f.pubErr != nil {
				return nil, f.pubErr
			}
			return &allocationbatchpb.PublishAllocationBatchResponse{Success: true}, nil
		},
		GetAllocationBatchListPageData: func(_ context.Context, r *allocationbatchpb.GetAllocationBatchListPageDataRequest) (*allocationbatchpb.GetAllocationBatchListPageDataResponse, error) {
			return &allocationbatchpb.GetAllocationBatchListPageDataResponse{AllocationBatchList: f.batches, Pagination: &commonpb.PaginationResponse{}}, nil
		},
		GetAllocationShareListPageData: func(_ context.Context, r *allocationsharepb.GetAllocationShareListPageDataRequest) (*allocationsharepb.GetAllocationShareListPageDataResponse, error) {
			id := r.GetFilters().GetFilters()[0].GetStringFilter().GetValue()
			return &allocationsharepb.GetAllocationShareListPageDataResponse{AllocationShareList: f.shares[id], Pagination: &commonpb.PaginationResponse{}}, nil
		},
		ReadCostSourceComponent: func(_ context.Context, r *costsourcecomponentpb.ReadCostSourceComponentRequest) (*costsourcecomponentpb.ReadCostSourceComponentResponse, error) {
			return &costsourcecomponentpb.ReadCostSourceComponentResponse{Data: []*costsourcecomponentpb.CostSourceComponent{f.component}}, nil
		},
		ListAgreementLineTerms: func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
			return &agreementlinetermpb.ListAgreementLineTermsResponse{Data: f.terms, Pagination: &commonpb.PaginationResponse{}}, nil
		},
		GetSubscriptionItemPageData: func(_ context.Context, r *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error) {
			return &subscriptionpb.GetSubscriptionItemPageDataResponse{Subscription: &subscriptionpb.Subscription{Id: r.GetSubscriptionId(), Name: "Lease " + r.GetSubscriptionId()}}, nil
		},
	}
}

func deps(f *fake) *Deps {
	return &Deps{Routes: ab.DefaultRoutes(), Labels: ab.DefaultLabels(), ComponentLabels: csc.DefaultLabels(), CommonLabels: common(), UseCases: f.uc()}
}

const rowKind = "ALLOCATION_SHARE_KIND_"

// body builds the drawer POST: parallel rows kind|sub|client|weight|label.
func body(denominator string, rows ...[5]string) string {
	v := url.Values{}
	for _, r := range rows {
		v.Add("share_kind", rowKind+r[0])
		v.Add("subscription_id", r[1])
		v.Add("client_id", r[2])
		v.Add("numerator", r[3])
		v.Add("row_label", r[4])
	}
	v.Set("denominator", denominator)
	return v.Encode()
}

var twoTenantsAndOwnUse = body("1000",
	[5]string{"RECOVERABLE", "s1", "cl1", "600", "Lease s1"},
	[5]string{"RECOVERABLE", "s2", "cl2", "300", "Lease s2"},
	[5]string{"OWN_USE", "", "", "100", "Own use"},
)

func TestDrawerOffersOnlyCoveredParticipantsPlusFixedRows(t *testing.T) {
	f := newFake()
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodGet, "", "id", "c1"))
	if res.Template != "allocation-batch-drawer-form" {
		t.Fatalf("template = %q (%d %v)", res.Template, res.StatusCode, res.Headers)
	}
	d := res.Data.(*form.Data)
	var slugs []string
	for _, r := range d.Rows {
		slugs = append(slugs, r.Slug)
	}
	if got := strings.Join(slugs, ","); got != "s1,s2,own-use,vacancy,common-loss" {
		t.Fatalf("rows = %s", got)
	}
	if d.SourceAmount != "1000.00" || d.Currency != "PHP" || !d.CanPublish || d.HasDraft {
		t.Fatalf("data = %+v", d)
	}
}

func TestDrawerPrefillsExistingDraft(t *testing.T) {
	f := newFake()
	f.batches = []*allocationbatchpb.AllocationBatch{{Id: "b1", CostSourceComponentId: "c1", Revision: 2, Status: allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_DRAFT}}
	sid := "s1"
	f.shares["b1"] = []*allocationsharepb.AllocationShare{{AllocationBatchId: "b1", ShareKind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_RECOVERABLE, SubscriptionId: &sid, BasisNumerator: 250, BasisDenominator: 1000, Amount: 25000}}
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodGet, "", "id", "c1"))
	d := res.Data.(*form.Data)
	if !d.HasDraft || d.Denominator != "1000" || d.Rows[0].Numerator != "250" || d.Preview == nil || d.Preview.Lines[0].Name != "Lease s1" {
		t.Fatalf("draft prefill = %+v preview=%+v", d, d.Preview)
	}
}

func TestDrawerRefusesClaimedAndUnauthorised(t *testing.T) {
	f := newFake()
	claim := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_RECOGNITION
	f.component.ClaimKind = &claim
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodGet, "", "id", "c1"))
	if res.Headers["HX-Error-Message"] != ab.DefaultLabels().Errors.SourceClaimedByRecognition {
		t.Fatalf("claimed = %d %v", res.StatusCode, res.Headers)
	}
	f = newFake()
	res = NewAllocateAction(deps(f)).Handle(ctx("allocation_batch:read"), req(http.MethodGet, "", "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" {
		t.Fatalf("denied = %v", res.Headers)
	}
}

func TestSaveDraftCreatesThenUpdates(t *testing.T) {
	f := newFake()
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.StatusCode != http.StatusOK || f.created == nil || f.updated != nil {
		t.Fatalf("create: %d created=%v updated=%v", res.StatusCode, f.created, f.updated)
	}
	if f.created.GetData().GetCostSourceComponentId() != "c1" || len(f.created.GetShares()) != 3 {
		t.Fatalf("create request = %v", f.created)
	}
	for _, s := range f.created.GetShares() {
		if s.GetBasisDenominator() != 1000 {
			t.Fatalf("shared denominator = %d", s.GetBasisDenominator())
		}
	}
	if s := f.created.GetShares()[2]; s.GetShareKind() != allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_OWN_USE || s.SubscriptionId != nil || s.ClientId != nil {
		t.Fatalf("own use share = %v", s)
	}

	f.created = nil
	f.batches = []*allocationbatchpb.AllocationBatch{{Id: "b9", CostSourceComponentId: "c1", Revision: 1, Status: allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_DRAFT}}
	res = NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.StatusCode != http.StatusOK || f.created != nil || f.updated.GetAllocationBatchId() != "b9" {
		t.Fatalf("update: %d created=%v updated=%v", res.StatusCode, f.created, f.updated)
	}
}

// A blank total weight defaults to the sum of the weights; a bad form is refused before any write.
func TestSaveDraftDenominatorDefaultAndBadForm(t *testing.T) {
	f := newFake()
	NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, body("", [5]string{"RECOVERABLE", "s1", "cl1", "40", "x"}, [5]string{"OWN_USE", "", "", "10", "y"}), "id", "c1"))
	if f.created == nil || f.created.GetShares()[0].GetBasisDenominator() != 50 {
		t.Fatalf("default denominator: %v", f.created)
	}
	f = newFake()
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, body("100", [5]string{"RECOVERABLE", "s1", "cl1", "-5", "x"}), "id", "c1"))
	if res.Headers["HX-Error-Message"] != ab.DefaultLabels().Errors.FormInvalid || f.created != nil {
		t.Fatalf("negative weight: %v created=%v", res.Headers, f.created)
	}
}

func TestPreviewShowsExactStoredAmounts(t *testing.T) {
	f := newFake()
	res := NewPreviewAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.Template != "allocation-batch-preview" {
		t.Fatalf("template %q (%d %v)", res.Template, res.StatusCode, res.Headers)
	}
	p := res.Data.(*form.PreviewPage).Preview
	if len(p.Lines) != 3 || p.Lines[0].Name != "Lease s1" || p.Lines[0].Amount != "600.00" || p.Lines[1].Amount != "300.00" || p.Lines[2].Amount != "100.00" {
		t.Fatalf("preview = %+v", p)
	}
	if p.Total != "1000.00" || !p.Balanced {
		t.Fatalf("total %s balanced %v", p.Total, p.Balanced)
	}
}

func TestPublishSavesThenPublishesAndMapsRefusals(t *testing.T) {
	f := newFake()
	res := NewPublishAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.StatusCode != http.StatusOK || len(f.published) != 1 || f.published[0] != "b1" {
		t.Fatalf("publish = %d %v", res.StatusCode, f.published)
	}
	l := ab.DefaultLabels()
	for code, want := range map[string]string{
		ab.ErrTermBoundaryCrossed:        l.Errors.TermBoundaryCrossed,
		ab.ErrSourceClaimedByRecognition: l.Errors.SourceClaimedByRecognition,
		ab.ErrSharesTotalMismatch:        l.Errors.SharesTotalMismatch,
	} {
		f = newFake()
		f.pubErr = coded{code}
		res = NewPublishAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
		if res.Headers["HX-Error-Message"] != want {
			t.Errorf("code %s: %q want %q", code, res.Headers["HX-Error-Message"], want)
		}
	}
	// a refused save never reaches publish
	f = newFake()
	f.writeErr = coded{ab.ErrSharesTotalMismatch}
	res = NewPublishAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if len(f.published) != 0 || res.Headers["HX-Error-Message"] != l.Errors.SharesTotalMismatch {
		t.Fatalf("save refusal: published=%v headers=%v", f.published, res.Headers)
	}
}

func TestPublishDeniedWithoutPublishPermission(t *testing.T) {
	f := newFake()
	res := NewPublishAction(deps(f)).Handle(ctx("allocation_batch:create", "allocation_batch:update"), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" || f.created != nil || len(f.published) != 0 {
		t.Fatalf("denied publish = %v created=%v published=%v", res.Headers, f.created, f.published)
	}
	d := deps(f)
	d.UseCases = &ab.UseCases{}
	res = NewPublishAction(d).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unwired publish = %d", res.StatusCode)
	}
}

func TestViewShowsPublishedAllocation(t *testing.T) {
	f := newFake()
	claim := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION
	f.component.ClaimKind = &claim
	f.batches = []*allocationbatchpb.AllocationBatch{{Id: "b1", CostSourceComponentId: "c1", Revision: 1, Status: allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_PUBLISHED}}
	sid := "s1"
	f.shares["b1"] = []*allocationsharepb.AllocationShare{
		{AllocationBatchId: "b1", ShareKind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_RECOVERABLE, SubscriptionId: &sid, BasisNumerator: 900, BasisDenominator: 1000, Amount: 90000, SequenceOrder: 1},
		{AllocationBatchId: "b1", ShareKind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_OWN_USE, BasisNumerator: 100, BasisDenominator: 1000, Amount: 10000, SequenceOrder: 2},
	}
	res := NewViewAction(deps(f)).Handle(ctx("allocation_batch:read"), req(http.MethodGet, "", "id", "c1"))
	if res.Template != "allocation-batch-view" {
		t.Fatalf("template = %q (%d %v)", res.Template, res.StatusCode, res.Headers)
	}
	v := res.Data.(*form.ViewData)
	if len(v.Preview.Lines) != 2 || v.Preview.Lines[0].Name != "Lease s1" || v.Preview.Lines[0].Percent != "90.00%" || v.Preview.Total != "1000.00" {
		t.Fatalf("view = %+v", v.Preview)
	}
	if r := NewViewAction(deps(f)).Handle(ctx("expenditure:read"), req(http.MethodGet, "", "id", "c1")); r.Headers["HX-Error-Message"] != "no permission" {
		t.Fatalf("denied view = %v", r.Headers)
	}
}

// R5 M1: a coded permission_denied refusal (strict gate) maps to the permission message on every
// write path, not the generic error.
func TestPermissionDeniedRefusalMapsToPermissionMessage(t *testing.T) {
	f := newFake()
	f.writeErr = coded{ab.ErrPermissionDenied}
	res := NewAllocateAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" {
		t.Errorf("save: %v", res.Headers)
	}
	f = newFake()
	f.pubErr = coded{ab.ErrPermissionDenied}
	res = NewPublishAction(deps(f)).Handle(ctx(all()...), req(http.MethodPost, twoTenantsAndOwnUse, "id", "c1"))
	if res.Headers["HX-Error-Message"] != "no permission" {
		t.Errorf("publish: %v", res.Headers)
	}
}

// R5 m3: the drawer carries a tooltip naming the missing publish code.
func TestDrawerPublishTooltipNamesMissingCode(t *testing.T) {
	f := newFake()
	d := deps(f)
	d.CommonLabels.Errors.MissingPermission = "Missing permission: %s"
	res := NewAllocateAction(d).Handle(ctx("allocation_batch:read", "allocation_batch:create"), req(http.MethodGet, "", "id", "c1"))
	data, ok := res.Data.(*form.Data)
	if !ok || data.CanPublish || data.PublishTooltip != "Missing permission: allocation_batch:publish" {
		t.Fatalf("drawer = %+v", res.Data)
	}
}
