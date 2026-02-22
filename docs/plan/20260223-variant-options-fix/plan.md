# Variant Options Fix — Design Plan

**Date:** 2026-02-23
**Branch:** `dev/20260223-variant-options-fix`
**Status:** Draft
**App/Package:** centymo-golang-ryta (primary), esqyma-ryta (proto), retail-admin (wiring)

---

## Overview

Fix three issues with variant option management: (1) saving options during edit fails silently due to unique constraint + soft delete conflict, (2) info tab only shows options that have assigned values — should show ALL product options even without a selection, (3) edit drawer should enforce required options when `product_option.required = true`.

---

## Motivation

When editing a variant's options (e.g., selecting Color = Black), the save appears to succeed (sheet closes) but the options are not persisted. The root cause: `deleteVariantOptions` soft-deletes old rows (sets `active=false`), then `saveVariantOptions` tries to INSERT new rows with the same `(product_variant_id, product_option_value_id)` pair — hitting the `unique_together` constraint on the `product_variant_option` table. The INSERT fails, error is logged but not surfaced to the user.

Additionally, the info tab completely hides the "Options" section when no variant options are assigned, giving no indication of which options are configurable. And there's no way to mark product options as required vs optional.

---

## Implementation Steps

### Phase 1: Proto Change — Add `required` field to ProductOption

- Add `bool required = 15` with default `false` to `product_option.proto`
- Rebuild proto: `cd packages/esqyma-ryta && pnpm generate`
- Run DDL generation: `cd packages/esqyma-ryta && pnpm run generate:ddl`
- Manually `ALTER TABLE product_option ADD COLUMN required BOOLEAN DEFAULT false` (or include in seed)

### Phase 2: Fix Variant Option Save — Hard Delete Junction Rows

The core bug: `deleteVariantOptions` uses proto-based soft delete, leaving ghost rows that block re-INSERT.

- In `deps.go:deleteVariantOptions`, replace `deps.DeleteProductVariantOption(...)` with `deps.DB.HardDelete(ctx, "product_variant_option", voID)` — junction rows have no meaningful history worth preserving
- Remove the `DeleteProductVariantOption` dependency from `Deps` struct (no longer needed)
- Remove wiring of `DeleteProductVariantOption` in `module.go` and `views.go`/`container.go`
- Also fix `NewRemoveView` in `action.go` — currently calls `deleteVariantOptions` before deleting the variant; this should also use hard delete

### Phase 3: Info Tab — Show All Product Options

Currently `loadVariantOptionEntries` returns only options that have an assigned value. Should show ALL product options for the product, with "—" for unassigned ones.

- Refactor `loadVariantOptionEntries` in `page.go`:
  1. Load ALL active product options for the product (via `ListProductOptions`, filter by `product_id`)
  2. Load ALL product option values (for label lookup)
  3. Load variant option records for this variant (existing logic)
  4. For each product option: if variant has a selection → show value label; if not → show "—"
- Update template `variant-tab-info` to always show the "Options" section header when there are product options (remove `{{if .OptionEntries}}` guard, or change to check if product HAS options regardless of assignment)

### Phase 4: Edit Drawer — Required Option Validation

- Add `Required bool` field to `OptionSelection` struct in `deps.go`
- In `loadOptionSelections`, populate `Required` from `o.GetRequired()`
- Update template `variant-drawer-form.html`:
  - Add `required` attribute to `<select>` when option is required
  - Add visual indicator (asterisk or "Required" label) for required options
- Add server-side validation in `NewEditView` and `NewAssignView` POST handlers:
  - After parsing form, check each required option has a non-empty value
  - Return `HtmxError` if any required option is missing

### Phase 5: Fix HtmxSuccess Target for Variant Detail Page

The edit response sends `HtmxSuccess("product-variants-table")` which targets the parent product detail page's table — not the variant detail page. When editing from within the variant detail page, the info tab doesn't refresh.

- Add a new helper or modify the edit POST response to trigger an info tab reload when editing from the variant detail page
- Option A: After successful edit, return `HX-Redirect` to the variant detail URL to force a full page reload
- Option B: Use `HX-Trigger: refreshVariantInfo` and add an HTMX listener on the info tab that reloads its content

---

## File References

| File | Change | Phase |
|------|--------|-------|
| `packages/esqyma-ryta/proto/v1/domain/product/product_option/product_option.proto` | Add `bool required = 15` | 1 |
| `packages/esqyma-ryta/pkg/schema/v1/domain/product/product_option/product_option.pb.go` | Regenerated | 1 |
| `packages/centymo-golang-ryta/views/product/detail/variant/deps.go` | Hard delete in `deleteVariantOptions`, add `Required` to `OptionSelection` | 2, 4 |
| `packages/centymo-golang-ryta/views/product/detail/variant/action.go` | Update `NewRemoveView` to hard delete, add required validation | 2, 4 |
| `packages/centymo-golang-ryta/views/product/detail/variant/page.go` | Refactor `loadVariantOptionEntries` to show all options | 3 |
| `packages/centymo-golang-ryta/templates/product/variant-detail.html` | Remove `{{if .OptionEntries}}` guard, show "—" for unassigned | 3 |
| `packages/centymo-golang-ryta/templates/product/variant-drawer-form.html` | Add `required` attribute and visual indicator | 4 |
| `apps/retail-admin/internal/presentation/product/module.go` | Remove `DeleteProductVariantOption` wiring | 2 |
| `apps/retail-admin/internal/composition/views.go` | Remove `DeleteProductVariantOption` wiring | 2 |
| `apps/retail-admin/internal/composition/container.go` | Remove `DeleteProductVariantOption` wiring | 2 |

---

## Context & Sub-Agent Strategy

**Estimated files to read:** ~12
**Estimated files to modify:** ~10
**Estimated context usage:** Low (<30 files)

No sub-agents needed. Single session is sufficient.

---

## Risk & Dependencies

| Risk | Impact | Mitigation |
|------|--------|------------|
| `required` column missing from DB | Medium — GetRequired() returns false | ALTER TABLE or seed SQL adds the column |
| Hard delete removes audit trail for junction rows | Low — junction rows have no meaningful history | Acceptable trade-off for correctness |

**Dependencies:**
- Phase 2-5 can proceed without Phase 1 (required field only needed for Phase 4)
- Phase 3 and Phase 5 are independent of each other
- Phase 4 depends on Phase 1 (proto field must exist)

---

## Acceptance Criteria

- [ ] Editing a variant's option selections persists correctly (verified via page reload)
- [ ] Info tab shows ALL product options for the product, with "—" for unassigned ones
- [ ] Required options show asterisk/indicator in the edit drawer
- [ ] Submitting the form without a required option shows an error
- [ ] Optional options can be left unselected and the form saves successfully
- [ ] Proto rebuild generates correct Go code with `Required` field
- [ ] Build passes: `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"`

---

## Design Decisions

**Hard delete vs. soft delete for junction rows:** `product_variant_option` is a pure junction table linking variants to option values. There's no business reason to preserve soft-deleted junction rows — they have no status, no history, and their only purpose is the relationship. Soft delete + unique constraint creates an irreconcilable conflict (ghost rows block re-creation). Hard delete is the correct choice here, matching how other junction tables (e.g., role_permission) should behave.

**Show all options on info tab vs. only assigned:** Showing all product options provides context about what's configurable, even before values are assigned. The "—" placeholder clearly indicates "not set" vs. "hidden". This is standard in commerce admin UIs (Shopify, WooCommerce).

**Option A (HX-Redirect) vs Option B (HX-Trigger) for post-edit refresh:** Option A (redirect) is simpler and guarantees fresh data but causes a full page reload. Option B (trigger) is smoother UX but requires adding HTMX listener wiring. Recommend Option A for simplicity unless the user prefers the smoother experience.
