package list

import (
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	collection "github.com/erniealice/centymo-golang/domain/treasury/collection"
)

func testDeps(receiveApplyURL string) *ListViewDeps {
	var cl pyeza.CommonLabels
	cl.Errors.MissingPermission = "Missing permission: %s"
	return &ListViewDeps{Routes: collection.DefaultRoutes(), Labels: collection.DefaultLabels(), CommonLabels: cl, ReceiveApplyURL: receiveApplyURL}
}

// Without the receive-and-apply flow the primary action is Add Collection,
// unchanged for every app that does not mount it.
func TestPrimaryActionDefaultsToAddCollection(t *testing.T) {
	deps := testDeps("")
	pa := primaryAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection:create"}))
	if pa.ActionURL != deps.Routes.AddURL || pa.Label != deps.Labels.Buttons.AddCollection || pa.Disabled || pa.TestID != "" {
		t.Errorf("primary action = %+v", pa)
	}
}

// With the flow mounted, Receive and apply becomes the primary action and is
// disabled (never hidden) without collection_application:create.
func TestPrimaryActionIsReceiveApplyWhenMounted(t *testing.T) {
	deps := testDeps("/action/collection/receive-apply")
	pa := primaryAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection_application:create"}))
	if pa.ActionURL != "/action/collection/receive-apply" || pa.Label != deps.Labels.Buttons.ReceiveApply || pa.Disabled || pa.TestID != "collection-receive-apply" {
		t.Errorf("primary action = %+v", pa)
	}
	pa = primaryAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection:create"}))
	if !pa.Disabled || pa.DisabledTooltip != "Missing permission: collection_application:create" {
		t.Errorf("without the permission the action must be disabled with the code: %+v", pa)
	}
}

// Option off: no secondary button (page unchanged).
func TestSecondaryAddAbsentWithoutReceiveApply(t *testing.T) {
	deps := testDeps("")
	if got := secondaryAddAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection:create"})); got != nil {
		t.Errorf("secondary = %+v, want nil", got)
	}
}

// Option on: Add Collection returns as a secondary button, disabled (never
// hidden) without collection:create.
func TestSecondaryAddPresentWithReceiveApply(t *testing.T) {
	deps := testDeps("/action/collection/receive-apply")
	got := secondaryAddAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection:create"}))
	if got == nil || got.ActionURL != deps.Routes.AddURL || got.Label != deps.Labels.Buttons.AddCollection || got.Disabled {
		t.Fatalf("secondary = %+v", got)
	}
	got = secondaryAddAction(deps, deps.Labels, types.NewUserPermissions([]string{"collection_application:create"}))
	if got == nil || !got.Disabled || got.DisabledTooltip != deps.Labels.Errors.PermissionDenied {
		t.Errorf("without collection:create it must be disabled with tooltip: %+v", got)
	}
}

// A receipt's mutating row actions are disabled with the coded reason; view stays enabled.
func TestReceiptRowActionsAreDisabledWithTheCodedReason(t *testing.T) {
	l := collection.DefaultLabels()
	in := []types.TableAction{{Action: "view"}, {Action: "edit"}, {Action: "deactivate"}, {Action: "activate"}, {Action: "delete"}}
	out := lockReceiptActions(in, l)
	for _, a := range out {
		wantLocked := a.Action != "view"
		if a.Disabled != wantLocked {
			t.Errorf("action %q disabled = %v, want %v", a.Action, a.Disabled, wantLocked)
		}
		if wantLocked && a.DisabledTooltip != l.Errors.ReceiptHasApplications {
			t.Errorf("action %q tooltip = %q", a.Action, a.DisabledTooltip)
		}
	}
	if l.Errors.ReceiptHasApplications == "" {
		t.Fatal("default label for receipt_has_applications is empty")
	}
}
