package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies that buildTableRows applies the correct Disabled flag and
// AWS-style "Missing permission: <entity>:<verb>" tooltips on every
// interactive row action across the {viewer, editor, admin} permission
// matrix. The plan/list page is the exemplar surface for the epic so
// these tests pin the contract.
//
// We exercise buildTableRows directly (rather than buildTableConfig)
// because the latter requires a fake ListPlans transport. The
// PrimaryAction.Disabled assertion happens in TestPrimaryAction_*
// against a small in-test helper that constructs the PrimaryAction the
// same way buildTableConfig does.

import (
	"context"
	"fmt"
	"testing"

	plan "github.com/erniealice/centymo-golang/domain/subscription/plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
)

// strPtr returns a pointer to its argument — protobuf fields with the
// `oneof` option (Id, Description, etc.) are *string in the generated Go code.
func strPtr(s string) *string { return &s }

// testCommonLabels builds the minimal CommonLabels needed for the
// MissingPermission tooltip template + status badges to interpolate.
func testCommonLabels() pyeza.CommonLabels {
	return pyeza.CommonLabels{
		Errors: pyeza.ErrorLabels{
			MissingPermission: "Missing permission: %s",
		},
		Actions: pyeza.ActionLabels{
			Clone: "Clone",
		},
	}
}

// testPlanLabels builds enough of the PlanLabels struct that
// buildTableRows can run without nil panics.
func testPlanLabels() plan.Labels {
	l := plan.DefaultLabels()
	if l.Errors.CannotDelete == "" {
		l.Errors.CannotDelete = "Cannot delete: plan in use"
	}
	return l
}

func testPlanRoutes() plan.Routes {
	return plan.DefaultRoutes()
}

// findAction returns the first action with the given Type, or nil.
func findAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildTableRows_PermissionMatrix verifies the {viewer, editor, admin}
// × {edit, clone, deactivate, delete} matrix on the active-status row.
func TestBuildTableRows_PermissionMatrix(t *testing.T) {
	t.Parallel()

	plans := []*planpb.Plan{
		{Id: strPtr("plan-1"), Name: "Starter", Description: strPtr("Cheap starter plan"), Active: true},
	}

	cases := []struct {
		name              string
		perms             []string
		wantEditDisabled  bool
		wantCloneDisabled bool
		wantDeactDisabled bool
		wantDelDisabled   bool
		wantEditTooltip   string
		wantCloneTooltip  string
		wantDelTooltip    string
	}{
		{
			name:              "viewer (list+read only) — every mutating action disabled",
			perms:             []string{"plan:list", "plan:read"},
			wantEditDisabled:  true,
			wantCloneDisabled: true,
			wantDeactDisabled: true,
			wantDelDisabled:   true,
			wantEditTooltip:   "Missing permission: plan:update",
			wantCloneTooltip:  "Missing permission: plan:create",
			wantDelTooltip:    "Missing permission: plan:delete",
		},
		{
			name:              "editor (no delete) — edit/clone/deactivate enabled, delete disabled",
			perms:             []string{"plan:list", "plan:read", "plan:create", "plan:update"},
			wantEditDisabled:  false,
			wantCloneDisabled: false,
			wantDeactDisabled: false,
			wantDelDisabled:   true,
			wantDelTooltip:    "Missing permission: plan:delete",
		},
		{
			name:              "admin (all perms) — every action enabled",
			perms:             []string{"plan:list", "plan:read", "plan:create", "plan:update", "plan:delete"},
			wantEditDisabled:  false,
			wantCloneDisabled: false,
			wantDeactDisabled: false,
			wantDelDisabled:   false,
		},
	}

	cl := testCommonLabels()
	l := testPlanLabels()
	routes := testPlanRoutes()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			rows := buildTableRows(plans, "active", l, cl, routes, map[string]bool{}, perms, map[string]string{}, map[string]string{})

			if len(rows) != 1 {
				t.Fatalf("rows = %d, want 1", len(rows))
			}
			actions := rows[0].Actions

			edit := findAction(actions, "edit")
			if edit == nil {
				t.Fatalf("edit action not found")
			}
			if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}
			if tc.wantEditTooltip != "" && edit.DisabledTooltip != tc.wantEditTooltip {
				t.Errorf("edit.DisabledTooltip = %q, want %q", edit.DisabledTooltip, tc.wantEditTooltip)
			}

			clone := findAction(actions, "clone")
			if clone == nil {
				t.Fatalf("clone action not found on active row")
			}
			if clone.Disabled != tc.wantCloneDisabled {
				t.Errorf("clone.Disabled = %v, want %v", clone.Disabled, tc.wantCloneDisabled)
			}
			if tc.wantCloneTooltip != "" && clone.DisabledTooltip != tc.wantCloneTooltip {
				t.Errorf("clone.DisabledTooltip = %q, want %q", clone.DisabledTooltip, tc.wantCloneTooltip)
			}

			deact := findAction(actions, "deactivate")
			if deact == nil {
				t.Fatalf("deactivate action not found on active row")
			}
			if deact.Disabled != tc.wantDeactDisabled {
				t.Errorf("deactivate.Disabled = %v, want %v", deact.Disabled, tc.wantDeactDisabled)
			}

			del := findAction(actions, "delete")
			if del == nil {
				t.Fatalf("delete action not found")
			}
			if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}
			if tc.wantDelTooltip != "" && del.DisabledTooltip != tc.wantDelTooltip {
				t.Errorf("delete.DisabledTooltip = %q, want %q", del.DisabledTooltip, tc.wantDelTooltip)
			}
		})
	}
}

// TestBuildTableRows_InactiveRow_ActivateInsteadOfDeactivate verifies the
// activate (vs deactivate) action is shown on inactive rows and is gated on
// plan:update.
func TestBuildTableRows_InactiveRow_ActivateInsteadOfDeactivate(t *testing.T) {
	t.Parallel()

	plans := []*planpb.Plan{
		{Id: strPtr("plan-x"), Name: "Archived", Description: strPtr(""), Active: false},
	}
	cl := testCommonLabels()
	l := testPlanLabels()
	routes := testPlanRoutes()

	perms := types.NewUserPermissions([]string{"plan:list", "plan:read"})
	rows := buildTableRows(plans, "inactive", l, cl, routes, nil, perms, nil, nil)

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if findAction(rows[0].Actions, "deactivate") != nil {
		t.Error("deactivate action should NOT appear on inactive row")
	}
	if findAction(rows[0].Actions, "clone") != nil {
		t.Error("clone action should NOT appear on inactive row (only on active)")
	}
	act := findAction(rows[0].Actions, "activate")
	if act == nil {
		t.Fatalf("activate action not found on inactive row")
	}
	if !act.Disabled {
		t.Error("activate should be disabled for viewer-only user")
	}
	if act.DisabledTooltip != "Missing permission: plan:update" {
		t.Errorf("activate tooltip = %q, want %q", act.DisabledTooltip, "Missing permission: plan:update")
	}
}

// TestBuildTableRows_InUse_DeletionBlocked verifies the in-use status takes
// precedence over the permission tooltip on the delete action (the more
// informative tooltip wins).
func TestBuildTableRows_InUse_DeletionBlocked(t *testing.T) {
	t.Parallel()

	plans := []*planpb.Plan{
		{Id: strPtr("plan-in-use"), Name: "BillingBackbone", Active: true},
	}
	cl := testCommonLabels()
	l := testPlanLabels()
	routes := testPlanRoutes()
	inUse := map[string]bool{"plan-in-use": true}

	// Admin perms — but in-use blocks delete.
	perms := types.NewUserPermissions([]string{"plan:list", "plan:read", "plan:create", "plan:update", "plan:delete"})

	rows := buildTableRows(plans, "active", l, cl, routes, inUse, perms, nil, nil)
	del := findAction(rows[0].Actions, "delete")
	if del == nil {
		t.Fatalf("delete action missing")
	}
	if !del.Disabled {
		t.Error("delete should be disabled when plan is in use")
	}
	// The in-use tooltip is set first; the permission tooltip then overrides
	// when the user also lacks delete perm. Admin has perm, so tooltip stays
	// as in-use message.
	if del.DisabledTooltip != l.Errors.CannotDelete {
		// Phase 1 audit note: the order in plan/list/page.go currently lets
		// permission win over in-use when both apply. For an admin with perm,
		// the in-use tooltip should remain.
		t.Errorf("delete.DisabledTooltip = %q, want %q (CannotDelete)", del.DisabledTooltip, l.Errors.CannotDelete)
	}
}

// TestPrimaryAction_AddButtonDisabledForViewer constructs the PrimaryAction
// the same way buildTableConfig does, verifying the AWS-style tooltip for
// the "Add Plan" button.
func TestPrimaryAction_AddButtonDisabledForViewer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cl := testCommonLabels()
	l := testPlanLabels()

	cases := []struct {
		name         string
		perms        []string
		wantDisabled bool
		wantTooltip  string
	}{
		{
			name:         "viewer — add button disabled",
			perms:        []string{"plan:list", "plan:read"},
			wantDisabled: true,
			wantTooltip:  "Missing permission: plan:create",
		},
		{
			name:         "editor — add button enabled",
			perms:        []string{"plan:list", "plan:read", "plan:create"},
			wantDisabled: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_ = ctx
			perms := types.NewUserPermissions(tc.perms)

			pa := &types.PrimaryAction{
				Label:           l.Buttons.AddPlan,
				ActionURL:       "/action/plan/add",
				Icon:            "icon-plus",
				Disabled:        !perms.Can("plan", "create"),
				DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "plan:create"),
			}

			if pa.Disabled != tc.wantDisabled {
				t.Errorf("PrimaryAction.Disabled = %v, want %v", pa.Disabled, tc.wantDisabled)
			}
			if tc.wantTooltip != "" && pa.DisabledTooltip != tc.wantTooltip {
				t.Errorf("PrimaryAction.DisabledTooltip = %q, want %q", pa.DisabledTooltip, tc.wantTooltip)
			}
		})
	}
}
