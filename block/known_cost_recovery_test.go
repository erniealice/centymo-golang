package block

import (
	"context"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	"sort"
	"strings"
	"testing"

	"github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/consumer/compose"
	"github.com/erniealice/pyeza-golang/view"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

var knownCostKeys = []string{"expenditure.cost_source_component", "expenditure.allocation_batch", "subscription.agreement_line_term"}

func unitKeys(us []compose.Unit) []string {
	var k []string
	for _, u := range us {
		k = append(k, u.Key)
	}
	return k
}

// The known-cost recovery surface is opt-in: without the option (or with false) the unit set
// is unchanged, so service-admin / school-admin compositions stay byte-identical; with it,
// exactly the three units are appended after the default order.
func TestKnownCostRecoveryIsOptIn(t *testing.T) {
	base := AllUnits(nil, nil)
	for _, u := range base {
		for _, k := range knownCostKeys {
			if u.Key == k {
				t.Fatalf("unit %s present without WithKnownCostRecovery", k)
			}
		}
	}
	if off := AllUnits(nil, nil, WithKnownCostRecovery(false)); len(off) != len(base) {
		t.Fatalf("WithKnownCostRecovery(false) changed the unit set: %d vs %d", len(off), len(base))
	}
	on := AllUnits(nil, nil, WithKnownCostRecovery(true))
	if len(on) != len(base)+3 {
		t.Fatalf("WithKnownCostRecovery(true) added %d units, want 3", len(on)-len(base))
	}
	added := unitKeys(on[len(base):])
	sort.Strings(added)
	want := append([]string(nil), knownCostKeys...)
	sort.Strings(want)
	if strings.Join(added, ",") != strings.Join(want, ",") {
		t.Fatalf("appended units = %v, want %v", added, want)
	}
	for i := range base {
		if base[i].Key != on[i].Key {
			t.Fatalf("default unit order changed at %d: %s vs %s", i, base[i].Key, on[i].Key)
		}
	}
}

// R4: every unit validates its closures in Mount and returns an error when a required one
// is nil (a Mount error stops boot; RequireFor/MustValidate do not run on the compose path).
func TestKnownCostRecoveryUnitsFailClosedInMount(t *testing.T) {
	for _, u := range AllUnits(&UseCases{}, nil, WithKnownCostRecovery(true)) {
		isKnown := false
		for _, k := range knownCostKeys {
			isKnown = isKnown || u.Key == k
		}
		if !isKnown {
			continue
		}
		if u.Mount == nil {
			t.Errorf("%s has no Mount", u.Key)
			continue
		}
		err := u.Mount(&compose.MountContext{})
		if err == nil {
			t.Errorf("%s: Mount accepted an empty UseCases", u.Key)
			continue
		}
		if !strings.Contains(err.Error(), "missing use cases") {
			t.Errorf("%s: error does not name the missing closures: %v", u.Key, err)
		}
	}
}

func fullUseCases() *UseCases {
	u := &UseCases{}
	u.Expenditure.ReadExpenditure = func(context.Context, *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error) {
		return nil, nil
	}
	c := u
	c.CostSourceComponent.CreateCostSourceComponent = func(context.Context, *costsourcecomponentpb.CreateCostSourceComponentRequest) (*costsourcecomponentpb.CreateCostSourceComponentResponse, error) {
		return nil, nil
	}
	c.CostSourceComponent.ReadCostSourceComponent = func(context.Context, *costsourcecomponentpb.ReadCostSourceComponentRequest) (*costsourcecomponentpb.ReadCostSourceComponentResponse, error) {
		return nil, nil
	}
	c.CostSourceComponent.UpdateCostSourceComponent = func(context.Context, *costsourcecomponentpb.UpdateCostSourceComponentRequest) (*costsourcecomponentpb.UpdateCostSourceComponentResponse, error) {
		return nil, nil
	}
	c.CostSourceComponent.DeleteCostSourceComponent = func(context.Context, *costsourcecomponentpb.DeleteCostSourceComponentRequest) (*costsourcecomponentpb.DeleteCostSourceComponentResponse, error) {
		return nil, nil
	}
	c.CostSourceComponent.ListCostSourceComponents = func(context.Context, *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error) {
		return nil, nil
	}
	a := &c.AllocationBatch
	a.CreateAllocationBatch = func(context.Context, *allocationbatchpb.CreateAllocationBatchRequest) (*allocationbatchpb.CreateAllocationBatchResponse, error) {
		return nil, nil
	}
	a.UpdateAllocationBatchShares = func(context.Context, *allocationbatchpb.UpdateAllocationBatchSharesRequest) (*allocationbatchpb.UpdateAllocationBatchSharesResponse, error) {
		return nil, nil
	}
	a.PublishAllocationBatch = func(context.Context, *allocationbatchpb.PublishAllocationBatchRequest) (*allocationbatchpb.PublishAllocationBatchResponse, error) {
		return nil, nil
	}
	a.GetAllocationBatchListPageData = func(context.Context, *allocationbatchpb.GetAllocationBatchListPageDataRequest) (*allocationbatchpb.GetAllocationBatchListPageDataResponse, error) {
		return nil, nil
	}
	a.GetAllocationShareListPageData = func(context.Context, *allocationsharepb.GetAllocationShareListPageDataRequest) (*allocationsharepb.GetAllocationShareListPageDataResponse, error) {
		return nil, nil
	}
	a.ReadCostSourceComponent = c.CostSourceComponent.ReadCostSourceComponent
	a.ListAgreementLineTerms = func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error) {
		return nil, nil
	}
	a.GetSubscriptionItemPageData = func(context.Context, *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error) {
		return nil, nil
	}
	c.AgreementLineTerm.ListAgreementLineTerms = a.ListAgreementLineTerms
	return u
}

type recorder struct{ paths []string }

func (r *recorder) GET(p string, _ view.View, _ ...string)  { r.paths = append(r.paths, "GET "+p) }
func (r *recorder) POST(p string, _ view.View, _ ...string) { r.paths = append(r.paths, "POST "+p) }

// With every closure bound the units mount and register their routes.
func TestKnownCostRecoveryUnitsMountRoutes(t *testing.T) {
	uc := fullUseCases()
	rec := &recorder{}
	for _, u := range knownCostRecoveryUnits(uc, nil) {
		if err := u.Mount(&compose.MountContext{Routes: rec}); err != nil {
			t.Fatalf("%s: %v", u.Key, err)
		}
	}
	got := strings.Join(rec.paths, "\n")
	for _, want := range []string{
		"POST /action/cost-source-component/add/{id}", "GET /action/cost-source-component/edit/{id}", "POST /action/cost-source-component/delete/{id}", "GET /action/cost-source-component/table/{id}",
		"GET /action/allocation-batch/allocate/{id}", "POST /action/allocation-batch/allocate/{id}", "POST /action/allocation-batch/preview/{id}", "POST /action/allocation-batch/publish/{id}", "GET /action/allocation-batch/view/{id}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("route not registered: %s", want)
		}
	}
	// cost lines: add GET/POST, edit GET/POST, delete, table = 6; allocation: allocate GET/POST,
	// preview, publish, view = 5.
	if len(rec.paths) != 11 {
		t.Errorf("registered %d routes: %v", len(rec.paths), rec.paths)
	}
}

// C13: the bind is nil-safe and binds nothing from an empty aggregate.
func TestBindBindsNothingFromAnEmptyAggregate(t *testing.T) {
	dst := buildCentymoUseCases(&consumer.UseCases{}, nil)
	if dst.CostSourceComponent.CreateCostSourceComponent != nil || dst.AllocationBatch.PublishAllocationBatch != nil ||
		dst.AgreementLineTerm.ListAgreementLineTerms != nil {
		t.Fatal("bind invented closures from an empty aggregate")
	}
}

// The S1 units are compose-v2 only: RequireFor never demands their closures.
func TestRequireForDoesNotCoverKnownCostRecovery(t *testing.T) {
	if err := (&UseCases{}).RequireFor(&blockConfig{}); err != nil && strings.Contains(err.Error(), "CostSourceComponent") {
		t.Errorf("surface closures required by RequireFor: %v", err)
	}
}

// An app that did not opt in mounts no sibling unit: the expenditure / subscription
// wiring helpers leave their deps untouched (the tabs are absent).
func TestNonOptInAppsGetNoTabs(t *testing.T) {
	uc := fullUseCases()
	deps := &expendituredomain.ExpenditureModuleDeps{}
	wireRecoverableCosts(&compose.MountContext{}, deps, uc)
	if deps.CostSourceComponents != nil {
		t.Fatal("recoverable costs tab wired without the opt-in unit")
	}
	if chargeTermsWiringOf(&compose.MountContext{}, uc) != nil {
		t.Fatal("charge terms tab wired without the opt-in unit")
	}
}

// The cost line unit needs the parent read (the Add drawer's existence proof): an unbound
// ReadExpenditure fails Mount naming the field instead of 404-ing every drawer at runtime.
func TestCostSourceComponentUnitRequiresTheParentRead(t *testing.T) {
	uc := fullUseCases()
	uc.Expenditure.ReadExpenditure = nil
	err := requireCostSourceComponent(uc)
	if err == nil || !strings.Contains(err.Error(), "UseCases.Expenditure.ReadExpenditure") {
		t.Fatalf("err = %v", err)
	}
}
