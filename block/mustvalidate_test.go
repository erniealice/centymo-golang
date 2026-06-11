package block

// MustValidate — FAIL-CLOSED wiring guard (architecture-roast burn #1).
//
// RequireFor returns an error; MustValidate adds the posture: in dev/test
// (testing.Testing() is true here) a missing REQUIRED closure PANICS — loud,
// stack-traced, uncatchable-by-accident — so a nil-closure wiring gap can never
// be silently dropped into an empty-state render. OPTIONAL nils never trip it.
//
// These four tests mirror the fayna reference impl
// (packages/fayna-golang/block/block_test.go), using centymo's first REQUIRED
// module (Inventory: the five InventoryItem CRUD closures) and an OPTIONAL
// module (Resource: has a wantResource() flag but is NOT asserted in
// RequireFor) to prove required-vs-optional survives the fail-closed wrapper.

import (
	"context"
	"strings"
	"testing"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
)

// wireInventoryRequired sets every closure RequireFor checks for the Inventory
// module: the five InventoryItem CRUD closures.
func wireInventoryRequired(uc *UseCases) {
	inv := &uc.Inventory
	inv.ListInventoryItems = func(context.Context, *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error) {
		return nil, nil
	}
	inv.CreateInventoryItem = func(context.Context, *inventoryitempb.CreateInventoryItemRequest) (*inventoryitempb.CreateInventoryItemResponse, error) {
		return nil, nil
	}
	inv.ReadInventoryItem = func(context.Context, *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error) {
		return nil, nil
	}
	inv.UpdateInventoryItem = func(context.Context, *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error) {
		return nil, nil
	}
	inv.DeleteInventoryItem = func(context.Context, *inventoryitempb.DeleteInventoryItemRequest) (*inventoryitempb.DeleteInventoryItemResponse, error) {
		return nil, nil
	}
}

// TestMustValidate_NilRequiredClosure_Panics is the core burn-#1 proof: with
// the Inventory module enabled but one REQUIRED closure (ListInventoryItems)
// left nil, MustValidate must PANIC under test — not return an empty render, not
// silently degrade. This is the loud failure the bare-return path lacked.
func TestMustValidate_NilRequiredClosure_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireInventoryRequired(uc)
	uc.Inventory.ListInventoryItems = nil // drop exactly one REQUIRED closure

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustValidate(Inventory enabled, ListInventoryItems nil) should PANIC in dev/test, but did not")
		}
		msg, _ := r.(string)
		if !strings.Contains(msg, "ListInventoryItems") {
			t.Fatalf("panic message should name the missing field; got %q", msg)
		}
	}()

	// Should not reach the next line — MustValidate panics first.
	_ = uc.MustValidate(&blockConfig{inventory: true})
	t.Fatal("MustValidate returned instead of panicking on a nil REQUIRED closure")
}

// TestMustValidate_EmptyUseCases_EnableAll_Panics: a fully empty UseCases with
// every module enabled (the "permanently nil dashboard" trap) must panic loudly
// in dev/test rather than register a wall of empty views.
func TestMustValidate_EmptyUseCases_EnableAll_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	defer func() {
		if recover() == nil {
			t.Fatal("MustValidate(empty UseCases, enableAll) should PANIC in dev/test")
		}
	}()
	_ = uc.MustValidate(&blockConfig{enableAll: true})
	t.Fatal("MustValidate returned instead of panicking on an empty enableAll wiring")
}

// TestMustValidate_NilOptionalClosure_OK proves the required-vs-optional
// discrimination survives the fail-closed wrapper: the OPTIONAL Resource module
// (a wantResource() flag NOT asserted in RequireFor) with nil closures must pass
// MustValidate with NO panic and NO error — disabled/optional features stay
// legitimately nil.
func TestMustValidate_NilOptionalClosure_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	// Optional module enabled, its closures left nil.
	cfg := &blockConfig{resource: true}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(optional nil closures) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(cfg); err != nil {
		t.Fatalf("MustValidate(optional nil closures) should be nil, got %v", err)
	}
}

// TestMustValidate_FullyWired_OK: a completely wired REQUIRED set passes with no
// panic and no error (happy path — guard is silent when wiring is complete).
func TestMustValidate_FullyWired_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireInventoryRequired(uc)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(fully wired Inventory) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(&blockConfig{inventory: true}); err != nil {
		t.Fatalf("MustValidate(fully wired Inventory) should be nil, got %v", err)
	}
}
