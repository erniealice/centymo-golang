# Product Variant Detail Fixes — Design Plan (Round 2)

**Date:** 2026-02-22
**Branch:** `dev/20260222-product-detail-restructure`
**Status:** Draft
**Package:** centymo-golang-ryta + retail-admin seeder

---

## Overview

Fix four post-implementation issues on the product variant detail page: (1) swap page title so variant name is on top and SKU is the subtitle, (2) widen spacing between info tab sections, (3) fix the non-functional "Edit Variant" button (missing `Sheet.open()` call), and (4) fix the Go seeder so the stock tab shows variant-linked inventory items.

---

## Motivation

After completing the initial product variant detail page implementation, user testing revealed:
- The page title shows SKU as primary title, which is less intuitive than showing the variant/product name first
- The "Edit Variant" button has HTMX attributes but no `onclick` handler to open the sheet
- Section spacing on the info tab is too tight between Variant Information, Options, and Actions
- The stock tab is always empty because the Go seeder never sets `product_variant_id` on inventory items

---

## Implementation Steps

### Phase 1: Fix Page Title — Name on Top, SKU as Subtitle

Currently in `variant_page.go`:
```go
HeaderTitle:    sku          // SKU as primary title
HeaderSubtitle: productName  // Product name as subtitle
```

Swap to:
```go
HeaderTitle:    productName  // Variant/product name as primary title
HeaderSubtitle: sku          // SKU as subtitle
```

- **Step 1.1:** In `packages/centymo-golang-ryta/views/product/detail/variant_page.go`, find where `HeaderTitle` and `HeaderSubtitle` are set. Swap them so `HeaderTitle = productName` and `HeaderSubtitle = sku`
- **Step 1.2:** Update `Title` field (browser tab title) to use the product name instead of SKU, or a composite like `"productName — sku"`

### Phase 2: Widen Info Tab Section Spacing

Currently sections use inline `style="margin-top: var(--spacing-lg);"` between Variant Information, Options, and Actions. The gap needs to be wider.

- **Step 2.1:** In `packages/centymo-golang-ryta/assets/css/variant-detail.css`, add a spacing override targeting section titles after content:
  ```css
  /* Widen gap between info tab sections */
  .tab-scroll .detail-section-title ~ .detail-section-title {
      margin-top: 2.5rem;
  }

  .tab-scroll .detail-actions {
      margin-top: 2.5rem;
  }
  ```
- **Step 2.2:** In `packages/centymo-golang-ryta/templates/product/variant-detail.html`, remove the inline `style="margin-top: var(--spacing-lg);"` from the Options section title and the `.detail-actions` div — let CSS handle it
- **Step 2.3:** Copy the updated CSS to `apps/retail-admin/assets/css/centymo/variant-detail.css`

### Phase 3: Fix "Edit Variant" Button

The button at `variant-detail.html:77-83` uses `data-sheet-title` which no JS processes. The working pattern uses `onclick="Sheet.open('...')"`.

Current (broken):
```html
<a class="btn btn--ghost btn--sm"
   hx-get="..."
   hx-target="#sheetContent"
   hx-swap="innerHTML"
   data-sheet-title="{{.Labels.Variant.Edit}}">
```

Fix to match working pattern (e.g., toolbar primary action in `table.html`):
```html
<a class="btn btn--ghost btn--sm"
   hx-get="..."
   hx-target="#sheetContent"
   hx-swap="innerHTML"
   hx-push-url="false"
   onclick="Sheet.open('{{.Labels.Variant.Edit}}')">
```

- **Step 3.1:** In `packages/centymo-golang-ryta/templates/product/variant-detail.html`, replace `data-sheet-title="{{.Labels.Variant.Edit}}"` with `onclick="Sheet.open('{{.Labels.Variant.Edit}}')"` on the Edit Variant button
- **Step 3.2:** Add `hx-push-url="false"` to prevent URL change when loading sheet content

### Phase 4: Fix Seeder — Link Inventory Items to Product Variants

The Go seeder `seedInventoryItems()` in `apps/retail-admin/cmd/seeder/main.go` (lines 926-983) creates inventory items with `product_id` but never sets `product_variant_id`. The stock tab at `variant_page.go:354-357` filters by `product_variant_id`, resulting in an empty table.

The SQL seed file (`products-inventory-seed.sql`) does it correctly — it sets `product_variant_id` on each inventory item. The Go seeder needs to match this pattern.

- **Step 4.1:** In `apps/retail-admin/cmd/seeder/main.go`, extend the `invRow` struct to include a `variantID string` field
- **Step 4.2:** Update the INSERT query (lines 967-976) to include `product_variant_id` column and its `$N` placeholder
- **Step 4.3:** Map each inventory item to an appropriate product variant — the variants are created in `seedProductVariants()`, so capture returned variant IDs or query them during inventory seeding
- **Step 4.4:** Ensure at least 2-3 inventory items per variant (for different locations) so the stock tab has visible data

---

## File References

| File | Change | Phase |
|------|--------|-------|
| `packages/centymo-golang-ryta/views/product/detail/variant_page.go` | Swap HeaderTitle ↔ HeaderSubtitle, update Title | 1 |
| `packages/centymo-golang-ryta/assets/css/variant-detail.css` | Add section spacing overrides | 2 |
| `apps/retail-admin/assets/css/centymo/variant-detail.css` | Mirror CSS changes from package | 2 |
| `packages/centymo-golang-ryta/templates/product/variant-detail.html` | Remove inline margin styles (Phase 2), fix Edit button onclick (Phase 3) | 2, 3 |
| `apps/retail-admin/cmd/seeder/main.go` | Add product_variant_id to inventory item seeding | 4 |

---

## Context & Sub-Agent Strategy

**Estimated files to read:** 5
**Estimated files to modify:** 5
**Estimated context usage:** Low (<30 files)

No sub-agents needed. All four phases are small, targeted fixes. Single session is sufficient.

---

## Risk & Dependencies

| Risk | Impact | Mitigation |
|------|--------|------------|
| Seeder variant ID references stale | Low — dev-only seeder | Query variant IDs by product_id during inventory seeding |
| CSS selector specificity conflicts | Low | Use `.tab-scroll` scoping to avoid bleeding to other detail pages |

**Dependencies:**
- Phases 1, 2, 3 are fully independent — can be done in any order
- Phase 4 is independent (different file entirely)

---

## Acceptance Criteria

- [ ] Variant detail page shows product name as primary title (top), SKU as subtitle (below)
- [ ] Info tab sections have visibly wider spacing between Variant Information, Options, and Actions
- [ ] "Edit Variant" button opens the sheet drawer when clicked
- [ ] Stock tab shows inventory items after re-running the seeder
- [ ] Build passes: `go build -tags "google_uuidv7,mock_auth,mock_storage,noop,postgresql,vanilla"`
- [ ] No regressions on product detail page or other variant tabs

---

## Design Decisions

**Title content:** Using the parent product name as HeaderTitle (e.g., "iPhone 17 Pro Max") since the option values are already visible in the info tab's Options section. The SKU subtitle provides the unique variant identifier at a glance.

**Section spacing approach:** Using CSS class-based spacing (`variant-detail.css`) rather than inline styles for maintainability. The selector `.tab-scroll .detail-section-title ~ .detail-section-title` targets subsequent section titles without affecting the first one, creating natural visual separation.

**Sheet.open pattern:** Following the established pattern from `pyeza-golang-ryta/components/table.html` where the toolbar primary action button uses `onclick="Sheet.open('...')"`. The `data-sheet-title` attribute was never wired to any JS handler and should be removed.
