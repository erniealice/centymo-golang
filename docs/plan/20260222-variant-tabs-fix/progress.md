# Variant Detail Tabs Fix -- Progress Log

**Plan:** [plan.md](./plan.md)
**Started:** 2026-02-22
**Branch:** `dev/20260222-variant-tabs-fix`

---

## Phase 1: Fix deleteVariantOptions to use hard-delete -- COMPLETE

- [x] Add `HardDelete(ctx context.Context, collection string, id string) error` to `centymo.DataSource` interface: `packages/centymo-golang-ryta/datasource.go`
- [x] Add `HardDelete` to `DatabaseOperation` interface: `packages/espyna-golang-ryta/internal/infrastructure/adapters/secondary/database/common/interface/operations.go` (was missing, added between Delete and List)
- [x] Add `HardDelete` wrapper to `espyna.DatabaseAdapter`: `packages/espyna-golang-ryta/consumer/adapter_database.go`
- [x] Check for mock implementations — added `HardDelete` to `MockOperations` (delegates to `Delete`): `packages/espyna-golang-ryta/internal/infrastructure/adapters/secondary/database/mock/core/operations.go`
- [x] Change `deleteVariantOptions` to call `db.HardDelete()` instead of `db.Delete()`: `packages/centymo-golang-ryta/views/product/detail/variants.go:183`
- [x] Verify build: `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"` — passed

---

## Phase 2: Fix seeder ON CONFLICT clauses -- COMPLETE

- [x] Change `seedProductVariantOptions` to use `ON CONFLICT (id) DO UPDATE SET product_option_value_id = EXCLUDED.product_option_value_id, active = true, date_modified = NOW()`: `apps/retail-admin/cmd/seeder/main.go`
- [x] Change `seedInventoryItems` to use `ON CONFLICT (id) DO UPDATE SET product_variant_id = EXCLUDED.product_variant_id, name = EXCLUDED.name, ...`: `apps/retail-admin/cmd/seeder/main.go`
- [ ] Re-run the seeder locally: `cd apps/retail-admin && go run cmd/seeder/main.go`
- [ ] Verify via SQL: `SELECT id, active FROM product_variant_option WHERE product_variant_id = 'pv-007-001'`
- [ ] Verify via SQL: `SELECT id, product_variant_id FROM inventory_item WHERE product_variant_id = 'pv-007-001'`

---

## Phase 3: Add product_variant_id column migration -- COMPLETE

- [x] Add `ALTER TABLE IF EXISTS "inventory_item" ADD COLUMN IF NOT EXISTS "product_variant_id" TEXT` at the top of `seedInventoryItems`: `apps/retail-admin/cmd/seeder/main.go`

---

## Phase 4: Verify with E2E tests -- NOT STARTED

- [ ] Navigate to variant detail page for pv-007-001 and verify Info tab shows Options section
- [ ] Verify Stock tab shows inventory items table (not "coming soon")
- [ ] Run existing E2E tests to confirm no regressions: `cd apps/retail-admin/tests && pnpm test`

---

## Phase 5: Re-run seeder and manual verification -- NOT STARTED

- [ ] Re-run seeder after all code changes
- [ ] Browse to `/app/products/detail/prod-007/variant/pv-007-001`
- [ ] Confirm Info tab Options section renders
- [ ] Confirm Stock tab shows table with inventory items
- [ ] Confirm editing variant options works (delete + recreate cycle)

---

## Summary

- **Phases complete:** 3 / 5 (code changes done, verification pending)
- **Files modified:** 6 / 7
  - `packages/espyna-golang-ryta/internal/infrastructure/adapters/secondary/database/common/interface/operations.go` — added HardDelete to interface
  - `packages/espyna-golang-ryta/internal/infrastructure/adapters/secondary/database/mock/core/operations.go` — added HardDelete mock
  - `packages/espyna-golang-ryta/consumer/adapter_database.go` — added HardDelete wrapper
  - `packages/centymo-golang-ryta/datasource.go` — added HardDelete to DataSource
  - `packages/centymo-golang-ryta/views/product/detail/variants.go` — changed deleteVariantOptions to use HardDelete
  - `apps/retail-admin/cmd/seeder/main.go` — fixed ON CONFLICT clauses + added ALTER TABLE migration

---

## Skipped / Deferred (update as you work)

| Item | Reason |
|------|--------|
| -- | -- |

---

## How to Resume

All code changes for Phases 1-3 are complete and build-verified. To continue:
1. Re-run the seeder: `cd apps/retail-admin && go run cmd/seeder/main.go`
2. Verify DB state with SQL queries (Phase 2 remaining items)
3. Start the server and navigate to variant detail page for manual verification (Phase 5)
4. Run E2E tests for regression check (Phase 4)
