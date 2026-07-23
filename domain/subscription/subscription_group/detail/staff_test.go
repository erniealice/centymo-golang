package detail

// staff_test.go — fail-closed unit coverage for the Teaching-Staff tab (§6.1).
//
// This is the coverage the E2E spec's Case-7 skip cited (as living in
// page_test.go) but that did not exist: page_test.go tests only the
// Subscriptions/roster tab. The tab's permission gates — the tab-body read
// gate (sgppsEntity:list), the per-row action edit gate (create/update), the
// drawer render gate (CanAny create/update), and the assign/clear write gates
// (create/update vs delete) — are security-critical fail-closed paths and are
// asserted here so a regression cannot silently open them.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// staffStrPtr returns a pointer to s. Kept local to staff_test.go so the C1
// Teaching-Staff-tab commit is self-contained; the sibling strp helper lives in
// page_test.go, which lands with separate (roster/gender-bands) work.
func staffStrPtr(s string) *string { return &s }

const (
	sgppsList   = sgppsEntity + ":list"
	sgppsCreate = sgppsEntity + ":create"
	sgppsUpdate = sgppsEntity + ":update"
	sgppsDelete = sgppsEntity + ":delete"
)

// staffFixtureDeps wires a section (sec1 / plan1) with one unassigned offering
// (pp1, "Mathematics") whose eligible pool holds one staff (staff1). A non-empty
// pool means the per-row action is NOT disabled by the empty-pool branch, so the
// CanEdit permission gate is what the disabled/enabled assertions actually
// exercise. assign records every AssignGroupServicer call so the write-gate tests
// can prove the closure is (or is not) reached.
func staffFixtureDeps(assign *assignRecorder) *DetailViewDeps {
	deps := &DetailViewDeps{
		Routes:       subscription_group.DefaultRoutes(),
		Labels:       subscription_group.DefaultLabels(),
		CommonLabels: pyeza.CommonLabels{Buttons: pyeza.ButtonLabels{Cancel: "Cancel", Save: "Save", Update: "Update"}},
		TableLabels:  types.TableLabels{},
		ReadSubscriptionGroup: func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error) {
			return &subscriptiongrouppb.ReadSubscriptionGroupResponse{Data: []*subscriptiongrouppb.SubscriptionGroup{{Id: "sec1", PlanId: staffStrPtr("plan1")}}}, nil
		},
		ListProductPlans: func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error) {
			return &productplanpb.ListProductPlansResponse{Data: []*productplanpb.ProductPlan{{Id: "pp1", Name: "Mathematics", Active: true}}}, nil
		},
		ListSubscriptionGroupProductPlanStaffs: func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error) {
			return &sgppspb.ListSubscriptionGroupProductPlanStaffsResponse{Data: nil}, nil // pp1 unassigned
		},
		ListProductPlanStaffs: func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error) {
			return &productplanstaffpb.ListProductPlanStaffsResponse{Data: []*productplanstaffpb.ProductPlanStaff{{Id: "pps1", ProductPlanId: "pp1", StaffId: "staff1", Role: "primary", Active: true}}}, nil
		},
		ListStaffNames: func(ctx context.Context) map[string]string { return map[string]string{"staff1": "Angeline Daoang"} },
	}
	if assign != nil {
		deps.AssignGroupServicer = func(ctx context.Context, groupID, productPlanID, staffID, role string) (string, error) {
			assign.calls++
			assign.groupID, assign.productPlanID, assign.staffID, assign.role = groupID, productPlanID, staffID, role
			return "edge1", nil
		}
	}
	return deps
}

// assignRecorder captures the arguments the assign closure was invoked with.
type assignRecorder struct {
	calls                           int
	groupID, productPlanID, staffID string
	role                            string
}

func staffCtx(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func firstAction(t *testing.T, cfg *types.TableConfig) types.TableAction {
	t.Helper()
	if cfg == nil || len(cfg.Rows) == 0 || len(cfg.Rows[0].Actions) == 0 {
		t.Fatalf("expected a table with one row + one action, got %+v", cfg)
	}
	return cfg.Rows[0].Actions[0]
}

func postForm(target string, form url.Values) *view.ViewContext {
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "sec1")
	return &view.ViewContext{Request: req}
}

// TestBuildStaffTabData_NoListPerm_FailClosed: without sgpps:list the tab must
// render the read-only note — Authorized=false, no table, no rows, CanEdit=false
// (the template's {{else}} branch shows Staff.Unauthorized, not a table/drawer).
func TestBuildStaffTabData_NoListPerm_FailClosed(t *testing.T) {
	deps := staffFixtureDeps(nil)
	sg := &subscriptiongrouppb.SubscriptionGroup{Id: "sec1", PlanId: staffStrPtr("plan1")}
	data := buildStaffTabData(staffCtx(), deps, sg, deps.Labels) // empty perms

	if data.Authorized {
		t.Errorf("Authorized = true, want false without %s", sgppsList)
	}
	if data.Table != nil {
		t.Errorf("Table = %+v, want nil (no table without :list)", data.Table)
	}
	if len(data.Rows) != 0 {
		t.Errorf("Rows = %d, want 0 without :list", len(data.Rows))
	}
	if data.CanEdit {
		t.Errorf("CanEdit = true, want false without :list")
	}
}

// TestBuildStaffTabData_ListOnly_RowActionDisabled: with :list but no create/
// update the table renders but every per-row assign action is disabled with the
// unauthorized tooltip (the write affordance is fail-closed on read-only perms).
func TestBuildStaffTabData_ListOnly_RowActionDisabled(t *testing.T) {
	deps := staffFixtureDeps(nil)
	sg := &subscriptiongrouppb.SubscriptionGroup{Id: "sec1", PlanId: staffStrPtr("plan1")}
	data := buildStaffTabData(staffCtx(sgppsList), deps, sg, deps.Labels)

	if !data.Authorized {
		t.Fatalf("Authorized = false, want true with %s", sgppsList)
	}
	if data.CanEdit {
		t.Errorf("CanEdit = true, want false without create/update")
	}
	act := firstAction(t, data.Table)
	if !act.Disabled {
		t.Errorf("row action Disabled = false, want true without create/update")
	}
	if act.DisabledTooltip != deps.Labels.Staff.Unauthorized {
		t.Errorf("DisabledTooltip = %q, want %q", act.DisabledTooltip, deps.Labels.Staff.Unauthorized)
	}
}

// TestBuildStaffTabData_CreatePerm_RowActionEnabled: with :list + :create the
// action is enabled AND is a real, focusable control — a URL-driven edit button
// (Action="edit", URL set, HxGet empty), NOT the href-less <a> that the HxGet
// branch emits (P2-H1 / WCAG 2.1.1 regression guard).
func TestBuildStaffTabData_CreatePerm_RowActionEnabled(t *testing.T) {
	deps := staffFixtureDeps(nil)
	sg := &subscriptiongrouppb.SubscriptionGroup{Id: "sec1", PlanId: staffStrPtr("plan1")}
	data := buildStaffTabData(staffCtx(sgppsList, sgppsCreate), deps, sg, deps.Labels)

	if !data.CanEdit {
		t.Errorf("CanEdit = false, want true with :create")
	}
	act := firstAction(t, data.Table)
	if act.Disabled {
		t.Errorf("row action Disabled = true, want enabled with :create and a non-empty pool")
	}
	if act.Action != "edit" {
		t.Errorf("row action Action = %q, want \"edit\" (native-button convention)", act.Action)
	}
	if act.URL == "" {
		t.Errorf("row action URL empty, want the assign-drawer URL (focusable edit button)")
	}
	if act.HxGet != "" {
		t.Errorf("row action HxGet = %q, want empty (the HxGet branch renders an <a> without href)", act.HxGet)
	}
}

// TestBuildStaffTable_StandardTableParity locks the W2f standard-table parity:
// the staff table carries the same client-side toolbar/footer chrome the
// Subscriptions detail-tab enables (buildSubscriptionsTable) — search / columns /
// density / entries — plus the actions column the assign row action needs, and is
// NOT the embedded Minimal variant. A regression back to Minimal (no toolbar/
// footer) or a dropped chrome flag fails here. The tab-scoped empty copy
// (P2-LBL-3) must also survive standardization, NOT the generic l.Empty.* the
// Subscriptions tab uses.
func TestBuildStaffTable_StandardTableParity(t *testing.T) {
	deps := staffFixtureDeps(nil)
	sg := &subscriptiongrouppb.SubscriptionGroup{Id: "sec1", PlanId: staffStrPtr("plan1")}
	data := buildStaffTabData(staffCtx(sgppsList, sgppsCreate), deps, sg, deps.Labels)

	if data.Table == nil {
		t.Fatalf("Table = nil, want a rendered table with :list")
	}
	if data.Table.Minimal {
		t.Errorf("Table.Minimal = true, want false (standard table, not the embedded/minimal variant)")
	}
	if !data.Table.ShowSearch || !data.Table.ShowColumns || !data.Table.ShowDensity || !data.Table.ShowEntries {
		t.Errorf("Table chrome = {search:%v columns:%v density:%v entries:%v}, want all true (Subscriptions-tab standard)",
			data.Table.ShowSearch, data.Table.ShowColumns, data.Table.ShowDensity, data.Table.ShowEntries)
	}
	if !data.Table.ShowActions {
		t.Errorf("Table.ShowActions = false, want true (the assign row action needs the actions column)")
	}
	if data.Table.EmptyState.Title != deps.Labels.Staff.EmptyTitle {
		t.Errorf("EmptyState.Title = %q, want the tab-scoped %q", data.Table.EmptyState.Title, deps.Labels.Staff.EmptyTitle)
	}
}

// TestRenderAssignDrawer_NoEditPerm_FailClosed: opening the drawer requires
// create OR update; :list alone gets a 422 HTMXError with no drawer template.
func TestRenderAssignDrawer_NoEditPerm_FailClosed(t *testing.T) {
	deps := staffFixtureDeps(nil)
	req := httptest.NewRequest(http.MethodGet, "/action/subscription-group/sec1/tab/staff-assign?pp=pp1", nil)
	res := renderAssignDrawer(staffCtx(sgppsList), deps, &view.ViewContext{Request: req}, "sec1")

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 without create/update", res.StatusCode)
	}
	if res.Template != "" {
		t.Errorf("Template = %q, want empty (no drawer rendered)", res.Template)
	}
	if got := res.Headers["HX-Error-Message"]; got != deps.Labels.Errors.Unauthorized {
		t.Errorf("HX-Error-Message = %q, want %q", got, deps.Labels.Errors.Unauthorized)
	}
}

// TestNewAssignAction_SetWithoutEditPerm_FailClosed: a set (staff present) with
// only :list must 422 and never reach the assign closure.
func TestNewAssignAction_SetWithoutEditPerm_FailClosed(t *testing.T) {
	rec := &assignRecorder{}
	deps := staffFixtureDeps(rec)
	vc := postForm("/action/subscription-group/assign/sec1", url.Values{
		"product_plan_id": {"pp1"}, "staff_id": {"staff1"}, "role": {"primary"},
	})
	res := NewAssignAction(deps).Handle(staffCtx(sgppsList), vc)

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for a set without create/update", res.StatusCode)
	}
	if got := res.Headers["HX-Error-Message"]; got != deps.Labels.Errors.Unauthorized {
		t.Errorf("HX-Error-Message = %q, want %q", got, deps.Labels.Errors.Unauthorized)
	}
	if rec.calls != 0 {
		t.Errorf("assign closure called %d times, want 0 (gate must short-circuit)", rec.calls)
	}
}

// TestNewAssignAction_ClearWithoutDeletePerm_FailClosed: a clear (blank staff)
// needs :delete; create+update alone must 422 and never reach the closure.
func TestNewAssignAction_ClearWithoutDeletePerm_FailClosed(t *testing.T) {
	rec := &assignRecorder{}
	deps := staffFixtureDeps(rec)
	vc := postForm("/action/subscription-group/assign/sec1", url.Values{
		"product_plan_id": {"pp1"}, "staff_id": {""}, "role": {"primary"},
	})
	res := NewAssignAction(deps).Handle(staffCtx(sgppsCreate, sgppsUpdate), vc)

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for a clear without :delete", res.StatusCode)
	}
	if got := res.Headers["HX-Error-Message"]; got != deps.Labels.Errors.Unauthorized {
		t.Errorf("HX-Error-Message = %q, want %q", got, deps.Labels.Errors.Unauthorized)
	}
	if rec.calls != 0 {
		t.Errorf("assign closure called %d times, want 0", rec.calls)
	}
}

// TestNewAssignAction_ClearFlagRoutesToDeleteGate: the drawer's explicit Clear
// control (clear=1) blanks staff_id even when the hidden input still carries one,
// so it is gated by :delete — create+update alone must 422 and not call assign.
func TestNewAssignAction_ClearFlagRoutesToDeleteGate(t *testing.T) {
	rec := &assignRecorder{}
	deps := staffFixtureDeps(rec)
	vc := postForm("/action/subscription-group/assign/sec1", url.Values{
		"product_plan_id": {"pp1"}, "staff_id": {"staff1"}, "role": {"primary"}, "clear": {"1"},
	})
	res := NewAssignAction(deps).Handle(staffCtx(sgppsCreate, sgppsUpdate), vc)

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 (clear=1 → :delete gate) without :delete", res.StatusCode)
	}
	if rec.calls != 0 {
		t.Errorf("assign closure called %d times, want 0 (clear must gate on :delete)", rec.calls)
	}
}

// TestNewAssignAction_SetWithCreatePerm_Succeeds proves the gates are not
// over-blocking: a set with :create reaches the closure with the section id from
// the signed path (never the body) and returns the staff-tab refresh trigger — so
// the fail-closed tests above fail for the right reason.
func TestNewAssignAction_SetWithCreatePerm_Succeeds(t *testing.T) {
	rec := &assignRecorder{}
	deps := staffFixtureDeps(rec)
	vc := postForm("/action/subscription-group/assign/sec1", url.Values{
		"product_plan_id": {"pp1"}, "staff_id": {"staff1"}, "role": {"primary"},
	})
	res := NewAssignAction(deps).Handle(staffCtx(sgppsCreate), vc)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 with :create", res.StatusCode)
	}
	if trig := res.Headers["HX-Trigger"]; !strings.Contains(trig, staffTabID) {
		t.Errorf("HX-Trigger = %q, want it to refresh %q", trig, staffTabID)
	}
	if rec.calls != 1 {
		t.Fatalf("assign closure called %d times, want 1", rec.calls)
	}
	if rec.groupID != "sec1" {
		t.Errorf("assign groupID = %q, want sec1 (from the signed path)", rec.groupID)
	}
	if rec.productPlanID != "pp1" || rec.staffID != "staff1" || rec.role != "primary" {
		t.Errorf("assign args = (%q,%q,%q), want (pp1,staff1,primary)", rec.productPlanID, rec.staffID, rec.role)
	}
}
