# Variant Subpackage Extraction — Design Plan

**Date:** 2026-02-22
**Branch:** `dev/20260222-variant-subpackage`
**Status:** Draft
**App/Package:** centymo-golang-ryta (primary), retail-admin (consumer)

---

## Overview

Extract the variant-related files (`variants.go`, `variant_page.go`) from the `detail` package into a new `detail/variant/` subpackage. This separates variant CRUD, variant detail page, and variant stock tab logic from the product detail page, reducing file count in the flat `detail/` directory and establishing a pattern for future subpackage extractions (e.g., `detail/option/`).

---

## Motivation

The `views/product/detail/` directory currently holds 8 files in a single `package detail`. The variant files (`variants.go` at 372 lines, `variant_page.go` at 440 lines) account for ~800 lines of variant-specific logic that has a clean one-way dependency on `page.go`. Extracting them into `detail/variant/` will:

1. **Reduce cognitive load** — variant logic is self-contained and easier to navigate as its own package
2. **Enable independent ownership** — variant package can evolve without touching the detail package
3. **Establish a pattern** — same approach can later be applied to `detail/option/` if the directory grows further
4. **Improve test isolation** — variant logic can get its own `_test.go` files without competing with detail tests

---

## Architecture

### Current structure (all `package detail`)

```
views/product/detail/
├── page.go              (370 lines) — product detail page, unexported: formatPrice, statusVariant, buildVariantsTable
├── variants.go          (372 lines) — variant CRUD, types: VariantDeps, VariantFormData, OptionSelection
├── variant_page.go      (440 lines) — variant detail page, stock tab, types: VariantPageData, OptionEntry
├── options.go           (331 lines) — options table builders, type: OptionsDeps
├── option_action.go     (229 lines) — option CRUD handlers
├── option_value_action.go (218 lines) — option value CRUD handlers
├── option_page.go       (138 lines) — option detail page, type: Breadcrumb, OptionPageData
└── attributes.go        (166 lines) — attribute CRUD handlers, type: AttributeDeps
```

### Proposed structure

```
views/product/detail/
├── page.go              — stays (EXPORT: FormatPrice, StatusVariant, BuildVariantsTable)
├── helpers.go           — NEW: extract HtmxSuccess, HtmxError from variants.go
├── options.go           — stays, update htmxSuccess→HtmxSuccess calls
├── option_action.go     — stays, update htmxSuccess→HtmxError calls
├── option_value_action.go — stays, update calls
├── option_page.go       — stays (Breadcrumb already exported)
├── attributes.go        — stays, update calls
└── variant/             — NEW Go subpackage (package variant)
    ├── deps.go          — Deps (was VariantDeps), VariantFormData, OptionSelection, etc.
    ├── page.go          — variant detail page + tab action + stock tab (from variant_page.go)
    └── action.go        — variant CRUD: table refresh, assign, edit, remove (from variants.go)
```

### Cross-file dependency map

```
variants.go → page.go:         Deps, buildVariantsTable, statusVariant (read-only)
variant_page.go → page.go:     formatPrice, statusVariant (read-only)
variant_page.go → option_page.go: Breadcrumb (read-only)
variant_page.go → variants.go: VariantDeps (type only)
option_action.go → variants.go:      htmxSuccess, htmxError
option_value_action.go → variants.go: htmxSuccess, htmxError
attributes.go → variants.go:         htmxSuccess, htmxError
page.go → (nothing in variant files)  ← NO circular dependency
```

### Export plan for `page.go`

| Current (unexported) | New (exported) | Used by |
|---|---|---|
| `formatPrice` | `FormatPrice` | variant/page.go |
| `statusVariant` | `StatusVariant` | variant/page.go, variant/action.go |
| `buildVariantsTable` | `BuildVariantsTable` | variant/action.go |

### Export plan for helpers.go (extracted from `variants.go:355-371`)

| Current (unexported) | New (exported) | Used by |
|---|---|---|
| `htmxSuccess` | `HtmxSuccess` | option_action.go, option_value_action.go, attributes.go, variant/action.go |
| `htmxError` | `HtmxError` | option_action.go, option_value_action.go, attributes.go, variant/action.go |

### Rename plan for `variant/` (package name provides context)

| Current (`productdetail.X`) | New (`productvariant.X`) |
|---|---|
| `VariantDeps` | `Deps` |
| `NewVariantsTableView` | `NewTableView` |
| `NewVariantAssignView` | `NewAssignView` |
| `NewVariantEditView` | `NewEditView` |
| `NewVariantRemoveView` | `NewRemoveView` |
| `NewVariantPageView` | `NewPageView` |
| `NewVariantTabAction` | `NewTabAction` |

---

## Implementation Steps

### Phase 1: Extract helpers.go

Move `htmxSuccess` and `htmxError` out of `variants.go` into a new `helpers.go`, exporting them as `HtmxSuccess` and `HtmxError`.

- Create `packages/centymo-golang-ryta/views/product/detail/helpers.go` with `HtmxSuccess` and `HtmxError`
- Remove `htmxSuccess`/`htmxError` functions from `packages/centymo-golang-ryta/views/product/detail/variants.go:355-371`
- Update all callers in `detail/` to use exported names:
  - `packages/centymo-golang-ryta/views/product/detail/variants.go` — 6 call sites (lines 228, 244, 251, 265, 302, 317, 323, 336, 344, 347)
  - `packages/centymo-golang-ryta/views/product/detail/option_action.go` — 6 call sites (lines 69, 95, 121, 136, 177, 204, 217, 224, 226)
  - `packages/centymo-golang-ryta/views/product/detail/option_value_action.go` — 6 call sites (lines 84, 108, 111, 127, 166, 193, 209, 215)
  - `packages/centymo-golang-ryta/views/product/detail/attributes.go` — 4 call sites (lines 105, 110, 138, 141, 155, 163)

### Phase 2: Export functions in page.go

Export the three functions that the variant subpackage will need.

- `packages/centymo-golang-ryta/views/product/detail/page.go:334` — rename `formatPrice` to `FormatPrice`
- `packages/centymo-golang-ryta/views/product/detail/page.go:360` — rename `statusVariant` to `StatusVariant`
- `packages/centymo-golang-ryta/views/product/detail/page.go:201` — rename `buildVariantsTable` to `BuildVariantsTable`
- Update all call sites within `page.go` itself (lines 110, 116-118, 177, 285)

### Phase 3: Create variant/ subpackage

Create three new files in `packages/centymo-golang-ryta/views/product/detail/variant/`.

**deps.go** — Types extracted from `variants.go:17-58`:
- `Deps` (was `VariantDeps`) — lines 53-58
- `VariantFormLabels` — lines 18-21
- `OptionValueChoice` — lines 24-27
- `OptionSelection` — lines 30-36
- `VariantFormData` — lines 39-50
- Unexported helpers: `loadOptionSelections`, `loadVariantOptionSelections`, `saveVariantOptions`, `deleteVariantOptions` (lines 61-188). Note: `deleteVariantOptions` now calls `db.HardDelete()` (changed by variant-tabs-fix plan)

**action.go** — View constructors extracted from `variants.go:190-349`:
- `NewTableView` (was `NewVariantsTableView`) — calls `detail.BuildVariantsTable`
- `NewAssignView` (was `NewVariantAssignView`) — calls `loadOptionSelections`, `detail.HtmxSuccess`, `detail.HtmxError`
- `NewEditView` (was `NewVariantEditView`) — calls `loadOptionSelections`, `loadVariantOptionSelections`, `detail.HtmxSuccess`, `detail.HtmxError`
- `NewRemoveView` (was `NewVariantRemoveView`) — calls `deleteVariantOptions`, `detail.HtmxSuccess`, `detail.HtmxError`

**page.go** — Views extracted from `variant_page.go:1-440`:
- `OptionEntry` type — lines 18-21
- `VariantPageData` type — lines 24-42 (references `detail.Breadcrumb`)
- `NewPageView` (was `NewVariantPageView`) — calls `detail.FormatPrice`, `detail.StatusVariant`
- `NewTabAction` (was `NewVariantTabAction`) — calls `detail.FormatPrice`, `detail.StatusVariant`
- Unexported helpers: `buildVariantTabItems`, `loadVariantOptionEntries`, `buildStockTable`

### Phase 4: Delete old files

- Delete `packages/centymo-golang-ryta/views/product/detail/variants.go` (all code moved to variant/ or helpers.go)
- Delete `packages/centymo-golang-ryta/views/product/detail/variant_page.go` (all code moved to variant/page.go)

### Phase 5: Update module.go imports and wiring

- `apps/retail-admin/internal/presentation/product/module.go:12` — add import: `productvariant "github.com/erniealice/centymo-golang/views/product/detail/variant"`
- `module.go:78` — change `&productdetail.VariantDeps{...}` to `&productvariant.Deps{...}`
- `module.go:107-113` — change all `productdetail.NewVariant*` to `productvariant.New*`:
  - `productdetail.NewVariantsTableView` → `productvariant.NewTableView`
  - `productdetail.NewVariantAssignView` → `productvariant.NewAssignView`
  - `productdetail.NewVariantEditView` → `productvariant.NewEditView`
  - `productdetail.NewVariantRemoveView` → `productvariant.NewRemoveView`
  - `productdetail.NewVariantPageView` → `productvariant.NewPageView`
  - `productdetail.NewVariantTabAction` → `productvariant.NewTabAction`

### Phase 6: Build verification

- Run `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"` in `apps/retail-admin/`
- Fix any compilation errors (missing imports, wrong types, unexported references)
- Run existing E2E tests: `cd apps/retail-admin/tests && pnpm test` (product-options + inventory-detail suites)

---

## File References

| File | Change | Phase |
|------|--------|-------|
| `packages/centymo-golang-ryta/views/product/detail/helpers.go` | **New file** — HtmxSuccess, HtmxError | 1 |
| `packages/centymo-golang-ryta/views/product/detail/variants.go` | Remove htmxSuccess/htmxError, update calls → HtmxSuccess/HtmxError | 1 |
| `packages/centymo-golang-ryta/views/product/detail/option_action.go` | Update htmxSuccess→HtmxSuccess, htmxError→HtmxError | 1 |
| `packages/centymo-golang-ryta/views/product/detail/option_value_action.go` | Update htmxSuccess→HtmxSuccess, htmxError→HtmxError | 1 |
| `packages/centymo-golang-ryta/views/product/detail/attributes.go` | Update htmxSuccess→HtmxSuccess, htmxError→HtmxError | 1 |
| `packages/centymo-golang-ryta/views/product/detail/page.go` | Export: FormatPrice, StatusVariant, BuildVariantsTable | 2 |
| `packages/centymo-golang-ryta/views/product/detail/variant/deps.go` | **New file** — Deps, types, unexported helpers | 3 |
| `packages/centymo-golang-ryta/views/product/detail/variant/action.go` | **New file** — NewTableView, NewAssignView, NewEditView, NewRemoveView | 3 |
| `packages/centymo-golang-ryta/views/product/detail/variant/page.go` | **New file** — NewPageView, NewTabAction, VariantPageData, stock tab | 3 |
| `packages/centymo-golang-ryta/views/product/detail/variants.go` | **Delete** — all code extracted | 4 |
| `packages/centymo-golang-ryta/views/product/detail/variant_page.go` | **Delete** — all code extracted | 4 |
| `apps/retail-admin/internal/presentation/product/module.go` | New import, update deps + view constructors | 5 |

---

## Context & Sub-Agent Strategy

**Estimated files to read:** 10 (8 detail/*.go + module.go + centymo labels)
**Estimated files to modify:** 8
**Estimated files to create:** 4 (helpers.go + 3 in variant/)
**Estimated files to delete:** 2 (variants.go + variant_page.go)
**Estimated context usage:** Low (<30 files)

No sub-agents needed. Single session is sufficient. All phases are sequential (each depends on the prior phase compiling cleanly).

---

## Risk & Dependencies

| Risk | Impact | Mitigation |
|------|--------|------------|
| Unexported helpers used across packages | Build failure — won't compile | Phase 1 exports helpers FIRST, verified before Phase 3 |
| Import cycle detail ↔ variant | Build failure — Go rejects cycles | Dependency is one-way: variant → detail. Verified by code audit: page.go never calls variant functions |
| Stale `go.work` or module cache | Build picks up old package | Run `go mod tidy` after creating variant/ subpackage |
| E2E tests break due to wiring change | Variant CRUD stops working | Phase 6 runs full E2E suite to catch regressions |

**Dependencies:**
- Phase 2 depends on Phase 1 (htmx helpers must be exported before variant/ uses them)
- Phase 3 depends on Phase 2 (exported functions must exist for variant/ imports)
- Phase 4 depends on Phase 3 (can only delete old files after new ones compile)
- Phase 5 depends on Phase 4 (module.go update happens after extraction is complete)
- Phase 6 depends on Phase 5 (build verification after all changes)

---

## Acceptance Criteria

- [ ] `go build` passes in `apps/retail-admin/` with all build tags
- [ ] `detail/variant/` is a separate Go package (`package variant`) with 3 files
- [ ] No code duplication between `detail/` and `detail/variant/`
- [ ] `variants.go` and `variant_page.go` are deleted from `detail/`
- [ ] `helpers.go` exports `HtmxSuccess` and `HtmxError`, used by all CRUD handlers
- [ ] `module.go` uses `productvariant` import alias for all variant views
- [ ] Existing E2E tests pass (product-options, inventory-detail suites)
- [ ] No import cycle between `detail` and `detail/variant`

---

## Design Decisions

**Why extract helpers.go instead of keeping htmxSuccess/htmxError in variants.go?**
These functions are used by `option_action.go`, `option_value_action.go`, `attributes.go`, and the new `variant/` subpackage. They are generic HTMX response helpers, not variant-specific. Extracting to `helpers.go` makes them a first-class shared utility within the `detail` package and avoids the variant subpackage needing to re-implement them.

**Why not extract options into detail/option/ at the same time?**
Scope discipline. The variant extraction is a clean, self-contained change. Options can follow the same pattern later if the `detail/` directory grows further. Doing both simultaneously would double the risk surface.

**Why use detail.Breadcrumb instead of duplicating the type?**
`Breadcrumb` is already exported from `option_page.go` and represents a generic UI concept shared across detail sub-pages (option detail, variant detail). Importing it from `detail` keeps the type canonical and avoids drift.

**Why rename VariantDeps to Deps in the new package?**
Go convention: the package name (`variant`) already provides the context, so `variant.Deps` reads better than `variant.VariantDeps`. Same logic applies to dropping the "Variant" prefix from all view constructors (`NewPageView` instead of `NewVariantPageView`).
