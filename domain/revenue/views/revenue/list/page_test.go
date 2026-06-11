package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies that buildTableRows on the revenue (invoice) list applies the
// correct Disabled flag to every interactive row action across the
// {viewer, editor, admin} matrix. The revenue surface uses "invoice:*"
// as its perm-code prefix (alias decision documented in phase1b-report).

import (
	"testing"

	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
	"github.com/erniealice/pyeza-golang/types"

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
)

// strPtr returns a pointer to its argument — protobuf fields with the
// `oneof` option (ReferenceNumber, FulfillmentStatus, ...) are *string
// in the generated Go code.
func strPtr(s string) *string { return &s }

func testRevenueLabels() revenuedomain.RevenueLabels {
	return revenuedomain.RevenueLabels{
		Errors: revenuedomain.RevenueErrorLabels{
			PermissionDenied:        "Missing permission",
			HasPaymentsCannotCancel: "Has payments — cannot cancel",
		},
		Actions: revenuedomain.RevenueActionLabels{
			View:               "View",
			Edit:               "Edit",
			Complete:           "Complete",
			Cancel:             "Cancel",
			DownloadInvoice:    "Download invoice",
			SendEmail:          "Send email",
			ReclassifyToDraft:  "Reclassify to draft",
			Delete:             "Delete",
		},
	}
}

func testRevenueRoutes() revenuedomain.RevenueRoutes {
	return revenuedomain.DefaultRevenueRoutes()
}

func findRevAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildTableRows_RevenuePermissionMatrix exercises the
// {viewer, editor, admin} matrix against the draft-status row actions:
// edit, check (Complete), delete (Cancel), download, mail (Send email).
func TestBuildTableRows_RevenuePermissionMatrix(t *testing.T) {
	t.Parallel()

	revenues := []*revenuepb.Revenue{
		{Id: "rev-1", ReferenceNumber: strPtr("INV-001"), Name: "Acme Sale", Status: "draft", Active: true},
	}

	cases := []struct {
		name              string
		perms             []string
		wantEditDisabled  bool
		wantCheckDisabled bool
		wantDelDisabled   bool
		wantDownloadDisabled bool
		wantMailDisabled  bool
	}{
		{
			name:                 "viewer (no invoice:read) — every mutating + download/mail disabled",
			perms:                []string{"invoice:list"},
			wantEditDisabled:     true,
			wantCheckDisabled:    true,
			wantDelDisabled:      true,
			wantDownloadDisabled: true,
			wantMailDisabled:     true,
		},
		{
			name:                 "viewer (invoice:list+read) — mutating disabled, download/mail enabled",
			perms:                []string{"invoice:list", "invoice:read"},
			wantEditDisabled:     true,
			wantCheckDisabled:    true,
			wantDelDisabled:      true,
			wantDownloadDisabled: false,
			wantMailDisabled:     false,
		},
		{
			name:                 "editor (no delete) — but delete on draft is status-change (uses :update)",
			perms:                []string{"invoice:list", "invoice:read", "invoice:create", "invoice:update"},
			wantEditDisabled:     false,
			wantCheckDisabled:    false,
			wantDelDisabled:      false, // delete on draft is a status transition gated on :update
			wantDownloadDisabled: false,
			wantMailDisabled:     false,
		},
		{
			name:                 "admin — every action enabled",
			perms:                []string{"invoice:list", "invoice:read", "invoice:create", "invoice:update", "invoice:delete"},
			wantEditDisabled:     false,
			wantCheckDisabled:    false,
			wantDelDisabled:      false,
			wantDownloadDisabled: false,
			wantMailDisabled:     false,
		},
	}

	l := testRevenueLabels()
	routes := testRevenueRoutes()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			rows := buildTableRows(revenues, "draft", l, routes, perms)
			if len(rows) != 1 {
				t.Fatalf("rows = %d, want 1", len(rows))
			}
			actions := rows[0].Actions

			if edit := findRevAction(actions, "edit"); edit == nil {
				t.Fatalf("edit action not found")
			} else if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}

			if check := findRevAction(actions, "check"); check == nil {
				t.Fatalf("check (Complete) action not found")
			} else if check.Disabled != tc.wantCheckDisabled {
				t.Errorf("check.Disabled = %v, want %v", check.Disabled, tc.wantCheckDisabled)
			}

			if del := findRevAction(actions, "delete"); del == nil {
				t.Fatalf("delete (Cancel) action not found")
			} else if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}

			if dl := findRevAction(actions, "download"); dl == nil {
				t.Fatalf("download action not found")
			} else if dl.Disabled != tc.wantDownloadDisabled {
				t.Errorf("download.Disabled = %v, want %v", dl.Disabled, tc.wantDownloadDisabled)
			}

			if mail := findRevAction(actions, "mail"); mail == nil {
				t.Fatalf("mail (Send email) action not found")
			} else if mail.Disabled != tc.wantMailDisabled {
				t.Errorf("mail.Disabled = %v, want %v", mail.Disabled, tc.wantMailDisabled)
			}
		})
	}
}

// TestBuildTableRows_CompleteRow_StatusGatedUndo verifies that the undo
// action on a complete revenue with payments stays disabled regardless
// of perms (status-gated DisabledTooltip takes priority over perm).
func TestBuildTableRows_CompleteRow_StatusGatedUndo(t *testing.T) {
	t.Parallel()

	revenues := []*revenuepb.Revenue{
		{
			Id:                "rev-paid",
			ReferenceNumber:   strPtr("INV-002"),
			Name:              "Paid Sale",
			Status:            "complete",
			Active:            true,
			FulfillmentStatus: strPtr("has_collection"),
		},
	}
	l := testRevenueLabels()
	routes := testRevenueRoutes()

	// Admin perms but the revenue has a collection → undo must stay disabled.
	perms := types.NewUserPermissions([]string{"invoice:list", "invoice:read", "invoice:update", "invoice:delete"})
	rows := buildTableRows(revenues, "complete", l, routes, perms)

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	undo := findRevAction(rows[0].Actions, "undo")
	if undo == nil {
		t.Fatalf("undo action not found")
	}
	if !undo.Disabled {
		t.Error("undo should be disabled when revenue has payments (state-gated)")
	}
	if undo.DisabledTooltip != l.Errors.HasPaymentsCannotCancel {
		t.Errorf("undo.DisabledTooltip = %q, want %q (state explanation, not perm)",
			undo.DisabledTooltip, l.Errors.HasPaymentsCannotCancel)
	}
}

// TestBuildTableRows_CancelledRow_NoMutatingActions verifies that a
// cancelled revenue only exposes a view action (no edit/delete).
func TestBuildTableRows_CancelledRow_NoMutatingActions(t *testing.T) {
	t.Parallel()

	revenues := []*revenuepb.Revenue{
		{Id: "rev-x", ReferenceNumber: strPtr("INV-X"), Name: "Cancelled", Status: "cancelled", Active: true},
	}
	l := testRevenueLabels()
	routes := testRevenueRoutes()

	perms := types.NewUserPermissions([]string{"invoice:list", "invoice:read", "invoice:update", "invoice:delete"})
	rows := buildTableRows(revenues, "cancelled", l, routes, perms)

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	// Cancelled rows only expose the view action.
	if len(rows[0].Actions) != 1 {
		t.Errorf("cancelled row actions = %d, want 1 (view only)", len(rows[0].Actions))
	}
	if findRevAction(rows[0].Actions, "view") == nil {
		t.Error("view action missing on cancelled row")
	}
}
