# Variant Detail Tabs Fix -- Design Plan

**Date:** 2026-02-22
**Branch:** `dev/20260222-variant-tabs-fix`
**Status:** In Progress
**App/Package:** centymo-golang (primary), retail-admin seeder (secondary)

---

## Overview

Fix two bugs on the variant detail page: (1) the Info tab's Options section not rendering for seeded variants, and (2) the Stock tab showing "coming soon" despite inventory items having `product_variant_id` set in the seeder code. Both bugs share the same root cause: the seeder uses `ON CONFLICT (id) DO NOTHING`, which means re-running the seeder after schema additions or data corrections has no effect on pre-existing rows.

---

## Motivation

The variant detail page is a key navigation target from the product detail's Variants tab. When Options and Stock data fail to render, the page looks empty and broken, undermining confidence in the inventory system. These bugs block the inventory detail and product options E2E test suites.

---

## Root Cause Analysis

### Bug 1: Info tab missing Options section

**Symptom:** Variant `pv-007-001` (AirPods Pro 2, SKU-007-001) has `product_variant_option` records in the DB, but the Options section does not appear on the Info tab.

**Root cause chain:**

1. `loadVariantOptionEntries()` (`variant_page.go:241`) calls `db.ListSimple(ctx, "product_variant_option")`.
2. `ListSimple` delegates to `PostgresOperations.List()` (`core/operations.go:327`).
3. `List()` hard-codes `WHERE active = true` as the first filter condition (line 333).
4. The two `product_variant_option` records for `pv-007-001` have `active = FALSE` in the database.
5. Therefore `ListSimple` returns zero matching rows, `loadVariantOptionEntries` returns `nil`, and the template's `{{if .OptionEntries}}` block is skipped.

**Why are the records `active = FALSE`?**

Two possible scenarios:
- **Scenario A (soft-delete cycle):** A previous edit of variant `pv-007-001` triggered `deleteVariantOptions()` (`variants.go:170`), which calls `db.Delete()` -- a soft delete (`UPDATE ... SET active = false`). The original seeder records (`pvo-007-001-charging`, etc.) got soft-deleted. If `saveVariantOptions` was called afterward, it created NEW records with auto-generated IDs and `active = true`. But if the user didn't select any options in the form (empty submission), no new records were created.
- **Scenario B (stale data):** The seeder `seedProductVariantOptions` inserts with `active: true` BUT uses `ON CONFLICT (id) DO NOTHING`. If these records were previously created (by an older seeder run or manual insert) with `active = false`, re-running the seeder does not update them.

**The fix requires two changes:**
1. **Seeder fix:** Change `product_variant_option` INSERT to `ON CONFLICT (id) DO UPDATE SET active = true, product_option_value_id = $3, date_modified = NOW()`.
2. **Delete behavior fix:** Change `deleteVariantOptions()` to use hard-delete instead of soft-delete for junction table records. Junction tables (`product_variant_option`) are association records, not business entities -- soft-delete creates orphaned inactive rows that block the seeder's `ON CONFLICT DO NOTHING`. The `DatabaseOperation` interface already provides `HardDelete()`.

### Bug 2: Stock tab empty

**Symptom:** Inventory items `inv-001` through `inv-022` have `product_variant_id` values in the seeder code, but the Stock tab shows the "coming soon" empty state.

**Root cause chain:**

1. `buildStockTable()` (`variant_page.go:330`) calls `db.ListSimple(ctx, "inventory_item")`.
2. `ListSimple` returns all active inventory items with all columns (uses `SELECT *`).
3. The function filters by `item["product_variant_id"].(string) == variantID` (line 367-369).
4. The filter produces zero matches, returning `nil` (line 410-412), triggering the "coming soon" panel.

**Why does the filter fail?**

The `inventory_item` table may not have had the `product_variant_id` column when the records were first created. The column was added later via `ALTER TABLE ADD COLUMN IF NOT EXISTS` in the test seed SQL. The seeder's INSERT uses `ON CONFLICT (id) DO NOTHING`, so if the rows already existed before `product_variant_id` was added to the INSERT statement, those rows have `NULL` for `product_variant_id`. The type assertion `item["product_variant_id"].(string)` on a `NULL` value yields `""`, which never matches any variant ID.

**The fix:**
1. **Seeder fix:** Change inventory_item INSERT to `ON CONFLICT (id) DO UPDATE SET product_variant_id = EXCLUDED.product_variant_id, ...` so re-running the seeder updates existing rows with the variant ID.

### Bug 3: Variant save issue (low priority, likely resolved by Bug 1 fix)

When editing a variant that previously had no options, the delete+recreate cycle in `saveVariantOptions` handles this correctly (new records get `active: true`). The reported issue was likely caused by the soft-deleted ghost records from a previous edit cycle. Once `deleteVariantOptions` uses hard-delete, this issue is resolved.

---

## Architecture

### Data flow: loadVariantOptionEntries

```
variant_page.go:loadVariantOptionEntries(ctx, db, variantID)
  -> db.ListSimple(ctx, "product_variant_option")
     -> PostgresOperations.List() -> SELECT * WHERE active = true
  -> filter by product_variant_id == variantID
  -> lookup product_option_value (label) and product_option (name)
  -> return []OptionEntry{Name, Value}
```

### Data flow: buildStockTable

```
variant_page.go:buildStockTable(ctx, deps, productID, variantID)
  -> db.ListSimple(ctx, "inventory_item")
     -> PostgresOperations.List() -> SELECT * WHERE active = true
  -> filter by product_variant_id == variantID
  -> build TableConfig rows
  -> return nil if no rows (triggers "coming soon" panel)
```

### Delete behavior comparison

```
Current (soft-delete):
  deleteVariantOptions -> db.Delete(id) -> UPDATE SET active=false
  Result: ghost rows with active=false that block ON CONFLICT DO NOTHING

Fixed (hard-delete):
  deleteVariantOptions -> db.HardDelete(id) -> DELETE FROM ... WHERE id=$1
  Result: rows removed, seeder can re-insert cleanly
```

---

## Implementation Steps

### Phase 1: Fix deleteVariantOptions to use hard-delete

The `deleteVariantOptions` function in `variants.go` currently uses `db.Delete()` (soft-delete). Junction table records like `product_variant_option` should be hard-deleted because they are pure association records with no business significance as inactive entities.

- Change `deleteVariantOptions` to call a hard-delete operation instead of `db.Delete()`: `packages/centymo-golang/views/product/detail/variants.go:170-188`
- The `centymo.DataSource` interface (`packages/centymo-golang/datasource.go`) currently only exposes `Delete()`. We need to either:
  - **Option A:** Add `HardDelete(ctx, collection, id)` to the `DataSource` interface
  - **Option B:** Use `db.Query()` or a raw SQL approach via the existing interface
  - **Option A is preferred** since `HardDelete` already exists on `PostgresOperations` and the `DatabaseAdapter` wraps it

Steps:
1. Add `HardDelete(ctx context.Context, collection string, id string) error` to `centymo.DataSource` interface: `packages/centymo-golang/datasource.go`
2. Ensure `espyna.DatabaseAdapter` exposes `HardDelete`: `packages/espyna-golang/consumer/adapter_database.go`
3. Check that `DatabaseOperation` interface already has `HardDelete`: `packages/espyna-golang/internal/infrastructure/adapters/secondary/database/common/interface/operations.go`
4. Update `deleteVariantOptions` to use `db.HardDelete()`: `packages/centymo-golang/views/product/detail/variants.go:183`

### Phase 2: Fix seeder ON CONFLICT clauses

Two seeder functions need `ON CONFLICT DO UPDATE` instead of `DO NOTHING`:

- **`seedProductVariantOptions`:** Change to `ON CONFLICT (id) DO UPDATE SET product_option_value_id = EXCLUDED.product_option_value_id, active = true, date_modified = NOW()`: `apps/retail-admin/cmd/seeder/main.go:1964-1970`
- **`seedInventoryItems`:** Change to `ON CONFLICT (id) DO UPDATE SET product_variant_id = EXCLUDED.product_variant_id, name = EXCLUDED.name, ...`: `apps/retail-admin/cmd/seeder/main.go:966-976`

### Phase 3: Add product_variant_id column migration

The `inventory_item` table's `product_variant_id` column is only added by the test seed SQL (`tests/seed/products-inventory-seed.sql`). For production, we need an `ALTER TABLE ADD COLUMN IF NOT EXISTS` in the seeder or a new migration. Since the seeder already creates tables when needed, adding the ALTER to the seeder is the pragmatic approach.

- Add `ALTER TABLE IF EXISTS inventory_item ADD COLUMN IF NOT EXISTS product_variant_id TEXT` at the top of `seedInventoryItems` (or as a separate migration step): `apps/retail-admin/cmd/seeder/main.go:926`

### Phase 4: Verify with Playwright E2E tests

Write or update E2E tests to verify both tabs render correctly for variant `pv-007-001`:

- Verify Info tab shows Options section with at least one option entry: `apps/retail-admin/tests/e2e/inventory-detail/` (or `product-options/`)
- Verify Stock tab shows the stock table (not the "coming soon" panel): same test directory

### Phase 5: Re-run seeder and manual verification

- Re-run the seeder to update existing records
- Navigate to `/app/products/detail/prod-007/variant/pv-007-001` and verify:
  - Info tab shows Options section with "Charging Case" -> "USB-C" (or similar)
  - Stock tab shows inventory items `inv-006` and `inv-011` in the table

---

## File References

| File | Change | Phase |
|------|--------|-------|
| `packages/centymo-golang/datasource.go` | Add `HardDelete` method to `DataSource` interface | 1 |
| `packages/espyna-golang/consumer/adapter_database.go` | Add `HardDelete` wrapper method | 1 |
| `packages/espyna-golang/internal/infrastructure/adapters/secondary/database/common/interface/operations.go` | Verify `HardDelete` exists on `DatabaseOperation` interface (may need adding) | 1 |
| `packages/centymo-golang/views/product/detail/variants.go:170-188` | Change `db.Delete()` to `db.HardDelete()` in `deleteVariantOptions` | 1 |
| `apps/retail-admin/cmd/seeder/main.go:1964-1970` | Change `seedProductVariantOptions` `ON CONFLICT` to `DO UPDATE` | 2 |
| `apps/retail-admin/cmd/seeder/main.go:966-976` | Change `seedInventoryItems` `ON CONFLICT` to `DO UPDATE` | 2 |
| `apps/retail-admin/cmd/seeder/main.go:926` | Add `ALTER TABLE` for `product_variant_id` column before inventory seeding | 3 |

---

## Context & Sub-Agent Strategy

**Estimated files to read:** 12
**Estimated files to modify:** 5-7
**Estimated context usage:** Low (<30 files)

No sub-agents needed. Single session is sufficient. The changes are focused and the root causes are well-understood.

---

## Risk & Dependencies

| Risk | Impact | Mitigation |
|------|--------|------------|
| Adding `HardDelete` to `DataSource` interface breaks mock implementations | Medium -- any mock or test double implementing `DataSource` must add the method | Check for mock implementations and update them |
| Seeder `DO UPDATE` may overwrite user-modified data | Low -- seeder is only run in development/staging | Document that seeder is destructive for seeded IDs |
| `deleteVariantOptions` hard-delete removes records permanently | Low -- junction records have no audit value; the variant edit creates fresh records anyway | No mitigation needed; this is the correct behavior |

**Dependencies:**
- Phase 2 depends on Phase 1 (seeder fix is ineffective if soft-deleted ghost rows remain)
- Phase 3 is independent and can run in parallel with Phase 1
- Phase 4 depends on Phases 1-3
- Phase 5 depends on Phase 2

---

## Acceptance Criteria

- [ ] Variant `pv-007-001` Info tab shows the Options section with option name-value pairs
- [ ] Variant `pv-007-001` Stock tab shows inventory items (inv-006, inv-011) in a table, not "coming soon"
- [ ] Editing a variant's options correctly deletes old junction records (hard-delete, not soft-delete)
- [ ] Re-running the seeder updates existing `product_variant_option` records to `active = true`
- [ ] Re-running the seeder updates existing `inventory_item` records with `product_variant_id`
- [ ] Build passes with `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"`
- [ ] Existing E2E tests continue to pass

---

## Design Decisions

**Hard-delete for junction tables:** Junction tables (`product_variant_option`) represent pure many-to-many associations. Unlike business entities (products, inventory items), they have no meaningful "inactive" state. Soft-deleting them creates ghost rows that interfere with `ON CONFLICT` inserts, `ListSimple` queries, and data integrity. Hard-delete is the correct approach for junction tables.

**Seeder `ON CONFLICT DO UPDATE` vs `DO NOTHING`:** The original `DO NOTHING` was a conservative choice to preserve user data. However, for development seeders that define canonical test data, `DO UPDATE` is correct -- the seeder should be idempotent and always produce the expected state. The test seed SQL (`products-inventory-seed.sql`) already uses `DO UPDATE` as the established pattern.

**Adding `HardDelete` to `DataSource` interface:** This is a small interface expansion that adds an important capability. The alternative (using raw SQL via `Query`) would break the abstraction layer. Since `HardDelete` already exists on `PostgresOperations`, exposing it through the interface chain is the clean approach.
