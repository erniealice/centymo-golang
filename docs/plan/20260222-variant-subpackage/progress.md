# Variant Subpackage Extraction — Progress Log

**Plan:** [plan.md](./plan.md)
**Started:** 2026-02-22
**Branch:** `dev/20260222-variant-subpackage`

---

## Phase 1: Extract helpers.go — DONE

- [x] Create `packages/centymo-golang/views/product/detail/helpers.go` with `HtmxSuccess` and `HtmxError`
- [x] Remove `htmxSuccess`/`htmxError` from `variants.go:355-371`
- [x] Update `variants.go` — replace all `htmxSuccess`/`htmxError` calls with `HtmxSuccess`/`HtmxError`
- [x] Update `option_action.go` — replace all `htmxSuccess`/`htmxError` calls
- [x] Update `option_value_action.go` — replace all `htmxSuccess`/`htmxError` calls
- [x] Update `attributes.go` — replace all `htmxSuccess`/`htmxError` calls
- [x] Verify build compiles with no errors

---

## Phase 2: Export functions in page.go — DONE

- [x] Rename `formatPrice` to `FormatPrice` in `page.go`
- [x] Rename `statusVariant` to `StatusVariant` in `page.go`
- [x] Rename `buildVariantsTable` to `BuildVariantsTable` in `page.go`
- [x] Update all call sites within `page.go`
- [x] Update call sites in `options.go` for `statusVariant` → `StatusVariant`
- [x] Update call sites in `variant_page.go` for `formatPrice`/`statusVariant`
- [x] Update call sites in `variants.go` for `buildVariantsTable`
- [x] Verify build compiles with no errors

---

## Phase 3: Create variant/ subpackage — DONE

- [x] Create directory `packages/centymo-golang/views/product/detail/variant/`
- [x] Create `variant/deps.go` — Deps type, VariantFormLabels, OptionValueChoice, OptionSelection, VariantFormData, unexported helpers (loadOptionSelections, loadVariantOptionSelections, saveVariantOptions, deleteVariantOptions)
- [x] Create `variant/action.go` — NewTableView, NewAssignView, NewEditView, NewRemoveView (import detail for BuildVariantsTable, HtmxSuccess, HtmxError)
- [x] Create `variant/page.go` — OptionEntry, VariantPageData, NewPageView, NewTabAction, buildVariantTabItems, loadVariantOptionEntries, buildStockTable (import detail for FormatPrice, StatusVariant, Breadcrumb)
- [x] Verify variant/ compiles as standalone package (no import cycles)

---

## Phase 4: Delete old files — DONE

- [x] Delete `packages/centymo-golang/views/product/detail/variants.go`
- [x] Delete `packages/centymo-golang/views/product/detail/variant_page.go`
- [x] Verify build compiles after deletion

---

## Phase 5: Update module.go imports and wiring — DONE

- [x] Add import `productvariant "github.com/erniealice/centymo-golang/views/product/detail/variant"` to `module.go`
- [x] Change `&productdetail.VariantDeps{...}` to `&productvariant.Deps{...}`
- [x] Change `productdetail.NewVariantsTableView` → `productvariant.NewTableView`
- [x] Change `productdetail.NewVariantAssignView` → `productvariant.NewAssignView`
- [x] Change `productdetail.NewVariantEditView` → `productvariant.NewEditView`
- [x] Change `productdetail.NewVariantRemoveView` → `productvariant.NewRemoveView`
- [x] Change `productdetail.NewVariantPageView` → `productvariant.NewPageView`
- [x] Change `productdetail.NewVariantTabAction` → `productvariant.NewTabAction`
- [x] Verify build compiles with no errors

---

## Phase 6: Build verification — DONE

- [x] Run `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"` in `apps/retail-admin/` — PASS
- [x] Run `go vet` with build tags — PASS
- [ ] Run E2E tests: `cd apps/retail-admin/tests && pnpm test` (deferred — not run in this session)

---

## Summary

- **Phases complete:** 6 / 6
- **Files modified:** 5 (page.go, options.go, option_action.go, option_value_action.go, attributes.go)
- **Files created:** 4 (helpers.go, variant/deps.go, variant/action.go, variant/page.go)
- **Files deleted:** 2 (variants.go, variant_page.go)
- **Module updated:** 1 (apps/retail-admin/internal/presentation/product/module.go)

---

## Skipped / Deferred (update as you work)

| Item | Reason |
|------|--------|
| E2E tests | Not run in this session — requires running server |

---

## How to Resume

All phases are complete. Only E2E verification remains:
1. Start the retail-admin server: `cd apps/retail-admin && powershell -ExecutionPolicy Bypass -File scripts/run.ps1`
2. Run E2E tests: `cd apps/retail-admin/tests && pnpm test`
3. Specifically check: product-options and inventory-detail suites

**Key files created:**
- `packages/centymo-golang/views/product/detail/helpers.go` — HtmxSuccess, HtmxError
- `packages/centymo-golang/views/product/detail/variant/deps.go` — Deps, types, helpers
- `packages/centymo-golang/views/product/detail/variant/action.go` — CRUD views
- `packages/centymo-golang/views/product/detail/variant/page.go` — detail page + tabs
