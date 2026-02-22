# Product Variant Detail Fixes — Progress Log (Round 2)

**Plan:** [plan.md](./plan.md)
**Started:** 2026-02-22
**Branch:** `dev/20260222-product-detail-restructure`

---

## Phase 1: Fix Page Title — COMPLETE

- [x] Swap HeaderTitle (→ product name) and HeaderSubtitle (→ SKU) in `variant_page.go`
- [x] Update Title field for browser tab

---

## Phase 2: Widen Info Tab Section Spacing — COMPLETE

- [x] Add section spacing CSS override in `packages/centymo-golang-ryta/assets/css/variant-detail.css`
- [x] Remove inline `style="margin-top: ..."` from `variant-detail.html`
- [x] Mirror CSS to `apps/retail-admin/assets/css/centymo/variant-detail.css`

---

## Phase 3: Fix "Edit Variant" Button — COMPLETE

- [x] Replace `data-sheet-title` with `onclick="Sheet.open('...')"` in `variant-detail.html`
- [x] Add `hx-push-url="false"` to match standard pattern

---

## Phase 4: Fix Seeder — COMPLETE

- [x] Extend `invRow` struct with `variantID` field in `seeder/main.go`
- [x] Add `product_variant_id` to INSERT statement
- [x] Map inventory items to seeded product variants (mixed across variants for realistic data)
- [x] Multiple items share variants across locations (e.g., pv-001-001 at ACB + RG)

---

## Summary

- **Phases complete:** 4 / 4
- **Files modified:** 5 / 5
- **Build status:** Passing

---

## Skipped / Deferred (update as you work)

| Item | Reason |
|------|--------|
| — | — |

---

## How to Resume

All implementation is complete. Remaining:
1. Re-run the seeder to populate `product_variant_id` on inventory items
2. After package changes: `touch apps/retail-admin/internal/composition/container.go` for Air rebuild
3. Verify visually:
   - Variant detail page title shows product name (top) + SKU (subtitle)
   - Info tab sections have wider spacing
   - "Edit Variant" button opens the sheet drawer
   - Stock tab shows inventory items for the variant
