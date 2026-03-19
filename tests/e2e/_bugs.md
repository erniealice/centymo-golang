# Centymo E2E Test Bugs

Bugs discovered during E2E testing of centymo package views against service-admin (localhost:8081, professional business type).

---

## BUG-1: Revenue drawer form shows "Failed to render page"

**Date:** 2026-03-19
**Severity:** P0 (blocks revenue CRUD)
**Route:** `/action/revenue/add` (GET), `/action/revenue/edit/{id}` (GET)
**Affected tests:** CEN-REV-002 (drawer form renders completely, creates revenue), CEN-REV-003 (saves edit)

**Problem:** When the revenue add/edit drawer opens, the form template renders only the first two fields (Reference Number and Invoice Date), then shows "Failed to render page" text next to the date field. The remaining form fields (currency, status, location_id, name/customer, notes) do not render. The sheet footer with Save/Cancel buttons also does not appear.

**Screenshot evidence:** `test-results/revenue-revenue-crud-CEN-R-a9e53-wer-form-renders-completely-msedge/test-failed-1.png`

**Root cause (likely):** The `sales-drawer-form.html` template references `.Labels.Form.CurrencyPlaceholder`, `.Labels.Form.StatusOngoing`, `.Labels.Form.StatusComplete`, `.Labels.Form.StatusCancelled`, and `.Labels.Form.CustomerNamePlaceholder`. One or more of these nested label fields may not be populated in the `FormLabels` struct, causing a Go template execution error at render time. The `formLabels()` function in `views/revenue/action/action.go` maps top-level label keys like `sales.form.customer` but the template also accesses sub-fields on `.Labels.Form.*` that may not exist in the struct.

**Impact:** Revenue creation and editing are completely blocked. The drawer cannot submit because:
1. Required fields (#name, #location_id) are not rendered
2. The submit button in the sheet footer is not rendered

**Workaround:** None (drawer form is non-functional).

**Test behavior:** Tests that depend on the full form are skipped with `test.skip(true, 'BUG: ...')`.

---

## BUG-2: Revenue seed data has no ongoing records

**Date:** 2026-03-19
**Severity:** P2 (test coverage gap)
**Route:** `/app/revenue/list/ongoing`

**Problem:** The service-admin seed data (`service1` database) contains zero revenue records with status "ongoing". The revenue list page at `/app/revenue/list/ongoing` shows an empty table with the message "No ongoing bookings — Create your first booking to get started." (Showing 0 to 0 of 0 entries).

**Screenshot evidence:** `test-results/revenue-revenue-crud-CEN-R-f8bc3-ws-data-rows-or-empty-state-msedge/test-failed-1.png`

**Impact:** Tests that require existing revenue rows for edit/view/action-button verification must be skipped. Specifically:
- CEN-REV-001: row action buttons, view link navigation
- CEN-REV-003: edit drawer pre-fill, save edit

**Workaround:** Add seed data with ongoing revenue records, or fix BUG-1 so the create test can generate test data first.

**Test behavior:** Tests that require rows are skipped with `test.skip(true, 'No ongoing revenue rows in seed data')`.

---

## Summary

| Bug | Area | Severity | Tests Skipped |
|-----|------|----------|---------------|
| BUG-1 | Revenue drawer form rendering | P0 | 3 (form complete, create, edit save) |
| BUG-2 | Revenue seed data empty | P2 | 3 (row actions, view link, edit pre-fill) |

**Total tests:** 35
**Passing:** 29
**Skipped (due to bugs):** 6
**Failed:** 0
