package table

import (
	"testing"

	"github.com/erniealice/pyeza-golang/types"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
)

func input(codes ...string) Input {
	return Input{
		ExpenditureID: "e1", Currency: "PHP", Labels: csc.DefaultLabels(), Routes: csc.DefaultRoutes(),
		AllocateURL: "/action/allocation-batch/allocate/{id}", ViewAllocationURL: "/action/allocation-batch/view/{id}",
		Perms: types.NewUserPermissions(codes), NoPermission: "no permission",
	}
}

func actionByID(row types.TableRow, testid string) *types.TableAction {
	for i := range row.Actions {
		if row.Actions[i].TestID == testid {
			return &row.Actions[i]
		}
	}
	return nil
}

// Edit/delete are disabled with the claimed notice once claimed; unclaimed rows
// offer Allocate, allocation-claimed rows offer View allocation.
func TestClaimedRowsLockEditDeleteAndSwitchAllocationAction(t *testing.T) {
	alloc := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION
	recog := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_RECOGNITION
	rows := []*costsourcecomponentpb.CostSourceComponent{
		{Id: "open", Amount: 1000, Currency: "PHP"},
		{Id: "alloc", Amount: 2000, Currency: "PHP", ClaimKind: &alloc},
		{Id: "recog", Amount: 3000, Currency: "PHP", ClaimKind: &recog},
	}
	tc := Build(rows, input("expenditure:update", "allocation_batch:create", "allocation_batch:read"))
	if tc.ID != TableID || len(tc.Rows) != 3 {
		t.Fatalf("table = %s rows=%d", tc.ID, len(tc.Rows))
	}
	open, claimed, recognized := tc.Rows[0], tc.Rows[1], tc.Rows[2]
	if a := actionByID(open, "cost-source-component-edit-open"); a == nil || a.Disabled {
		t.Errorf("open edit = %+v", a)
	}
	if a := actionByID(open, "cost-source-component-allocate-open"); a == nil || a.Disabled || a.HxGet != "/action/allocation-batch/allocate/open" {
		t.Errorf("open allocate = %+v", a)
	}
	for _, id := range []string{"edit", "delete"} {
		a := actionByID(claimed, "cost-source-component-"+id+"-alloc")
		if a == nil || !a.Disabled || a.DisabledTooltip != csc.DefaultLabels().Detail.ClaimedNotice {
			t.Errorf("claimed %s = %+v", id, a)
		}
	}
	if a := actionByID(claimed, "cost-source-component-view-allocation-alloc"); a == nil || a.HxGet != "/action/allocation-batch/view/alloc" {
		t.Errorf("claimed view allocation = %+v", a)
	}
	if actionByID(claimed, "cost-source-component-allocate-alloc") != nil {
		t.Error("a claimed component must not offer Allocate")
	}
	if actionByID(recognized, "cost-source-component-allocate-recog") != nil || actionByID(recognized, "cost-source-component-view-allocation-recog") != nil {
		t.Error("a recognition-claimed component must offer neither Allocate nor View allocation")
	}
}

func TestRowActionsDisabledWithoutUpdatePermission(t *testing.T) {
	tc := Build([]*costsourcecomponentpb.CostSourceComponent{{Id: "open", Amount: 1000}}, input("expenditure:read"))
	for _, id := range []string{"edit", "delete", "allocate"} { // allocate: NoPermission fallback while MissingPermission is unset
		a := actionByID(tc.Rows[0], "cost-source-component-"+id+"-open")
		if a == nil || !a.Disabled || a.DisabledTooltip != "no permission" {
			t.Errorf("%s = %+v", id, a)
		}
	}
}

func TestCellsCarryQuantityPeriodAndClaim(t *testing.T) {
	qty, scale, unit := int64(12505), int32(1), "kWh"
	from, to := "2026-09-01", "2026-10-01"
	c := &costsourcecomponentpb.CostSourceComponent{Id: "c", Amount: 125075, Currency: "PHP", BasisQuantityScaled: &qty, BasisScale: &scale, BasisUnit: &unit, ServiceFrom: &from, ServiceTo: &to}
	tc := Build([]*costsourcecomponentpb.CostSourceComponent{c}, input("expenditure:read"))
	cells := tc.Rows[0].Cells
	if cells[2].Value != "1250.5 kWh" || cells[5].Value != "2026-09-01 to 2026-10-01" || cells[6].Value != csc.DefaultLabels().Enums.ClaimNone {
		t.Fatalf("cells = %+v", cells)
	}
}

// R5 M2/m3: the allocation actions are gated on allocation_batch (not expenditure) and their
// tooltips name the missing code.
func TestAllocationActionsGatedOnAllocationBatchPermissions(t *testing.T) {
	alloc := costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION
	rows := []*costsourcecomponentpb.CostSourceComponent{{Id: "open", Amount: 1}, {Id: "done", Amount: 2, ClaimKind: &alloc}}
	build := func(codes ...string) *types.TableConfig {
		in := input(codes...)
		in.MissingPermission = "Missing permission: %s"
		return Build(rows, in)
	}
	// expenditure:update alone does NOT enable Allocate, and View allocation needs allocation_batch:read
	tc := build("expenditure:update")
	if a := actionByID(tc.Rows[0], "cost-source-component-allocate-open"); a == nil || !a.Disabled || a.DisabledTooltip != "Missing permission: allocation_batch:create" {
		t.Errorf("allocate without allocation_batch = %+v", a)
	}
	if a := actionByID(tc.Rows[1], "cost-source-component-view-allocation-done"); a == nil || !a.Disabled || a.DisabledTooltip != "Missing permission: allocation_batch:read" {
		t.Errorf("view allocation without read = %+v", a)
	}
	for _, code := range []string{"allocation_batch:create", "allocation_batch:update"} {
		if a := actionByID(build(code).Rows[0], "cost-source-component-allocate-open"); a == nil || a.Disabled {
			t.Errorf("allocate with %s = %+v", code, a)
		}
	}
	if a := actionByID(build("allocation_batch:read").Rows[1], "cost-source-component-view-allocation-done"); a == nil || a.Disabled {
		t.Errorf("view allocation with read = %+v", a)
	}
}
