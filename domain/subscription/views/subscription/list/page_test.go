package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies that buildTableRows applies the correct Disabled flag on every
// interactive row action for the subscription list across the
// {viewer, editor, admin} permission matrix. Phase 2 of the epic added
// gating to this surface; these tests pin the contract.

import (
	"context"
	"testing"

	subscription "github.com/erniealice/centymo-golang/domain/subscription"
	"github.com/erniealice/pyeza-golang/types"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

func testSubLabels() subscription.SubscriptionLabels {
	l := subscription.DefaultSubscriptionLabels()
	if l.Errors.NoPermission == "" {
		l.Errors.NoPermission = "Missing permission"
	}
	if l.Errors.InUse == "" {
		l.Errors.InUse = "Cannot delete: subscription in use"
	}
	return l
}

func testSubRoutes() subscription.SubscriptionRoutes {
	return subscription.DefaultSubscriptionRoutes()
}

func findSubAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildTableRows_SubscriptionPermissionMatrix exercises the
// {viewer, editor, admin} matrix against the {edit, deactivate, delete}
// row actions on an active subscription.
func TestBuildTableRows_SubscriptionPermissionMatrix(t *testing.T) {
	t.Parallel()

	subs := []*subscriptionpb.Subscription{
		{Id: "sub-1", Name: "Acme Monthly", Active: true},
	}

	cases := []struct {
		name              string
		perms             []string
		wantEditDisabled  bool
		wantDeactDisabled bool
		wantDelDisabled   bool
	}{
		{
			name:              "viewer — every mutating action disabled",
			perms:             []string{"subscription:list", "subscription:read"},
			wantEditDisabled:  true,
			wantDeactDisabled: true,
			wantDelDisabled:   true,
		},
		{
			name:              "editor without delete — edit/deactivate enabled, delete disabled",
			perms:             []string{"subscription:list", "subscription:read", "subscription:create", "subscription:update"},
			wantEditDisabled:  false,
			wantDeactDisabled: false,
			wantDelDisabled:   true,
		},
		{
			name:              "admin — every action enabled",
			perms:             []string{"subscription:list", "subscription:read", "subscription:create", "subscription:update", "subscription:delete"},
			wantEditDisabled:  false,
			wantDeactDisabled: false,
			wantDelDisabled:   false,
		},
	}

	l := testSubLabels()
	routes := testSubRoutes()
	ctx := context.Background()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			rows := buildTableRows(ctx, subs, "active", l, routes, map[string]bool{}, perms)

			if len(rows) != 1 {
				t.Fatalf("rows = %d, want 1", len(rows))
			}
			actions := rows[0].Actions

			if edit := findSubAction(actions, "edit"); edit == nil {
				t.Fatalf("edit action not found")
			} else if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}

			if deact := findSubAction(actions, "deactivate"); deact == nil {
				t.Fatalf("deactivate action not found")
			} else if deact.Disabled != tc.wantDeactDisabled {
				t.Errorf("deactivate.Disabled = %v, want %v", deact.Disabled, tc.wantDeactDisabled)
			}

			if del := findSubAction(actions, "delete"); del == nil {
				t.Fatalf("delete action not found")
			} else if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}
		})
	}
}

// TestBuildTableRows_InactiveRow_ActivateInsteadOfDeactivate verifies that
// inactive subscription rows surface an activate (not deactivate) action.
func TestBuildTableRows_SubscriptionInactiveRow(t *testing.T) {
	t.Parallel()

	subs := []*subscriptionpb.Subscription{
		{Id: "sub-2", Name: "Old Plan", Active: false},
	}
	l := testSubLabels()
	routes := testSubRoutes()

	perms := types.NewUserPermissions([]string{"subscription:list", "subscription:read"})
	rows := buildTableRows(context.Background(), subs, "inactive", l, routes, nil, perms)

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if findSubAction(rows[0].Actions, "deactivate") != nil {
		t.Error("deactivate should not appear on inactive row")
	}
	act := findSubAction(rows[0].Actions, "activate")
	if act == nil {
		t.Fatalf("activate action missing on inactive row")
	}
	if !act.Disabled {
		t.Error("activate should be disabled for viewer-only user (missing subscription:update)")
	}
}
