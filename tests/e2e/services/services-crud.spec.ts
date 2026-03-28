import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * CEN-PRD-001: Services/Products List (professional route: /app/services/...)
 * CEN-PRD-002: Services/Products Add via Drawer
 * CEN-PRD-003: Services/Products Edit via Drawer
 *
 * Table ID: #products-table
 * Form fields: #name, #description, #price, #currency, #active (toggle)
 * Routes: /app/services/list/active, /action/services/add, /action/services/edit/{id}
 */

test.describe('CEN-PRD-001: Services List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();
  });

  test('displays products table with correct column headers', async ({ page }) => {
    const headers = page.locator('#products-table thead th');
    const count = await headers.count();
    // Checkbox + Name + Description + Price + Status + Actions = 6 minimum
    expect(count).toBeGreaterThanOrEqual(5);

    const headerTexts = await page.locator('#products-table thead th .column-label').allTextContents();
    // Column labels come from i18n — check structural count instead
    expect(headerTexts.length).toBeGreaterThanOrEqual(4);
  });

  test('shows data rows in the table', async ({ page }) => {
    const rows = page.locator('#products-table tbody tr');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(1);

    // First row should have cells
    const firstRowCells = page.locator('#products-table tbody tr:first-child td');
    const cellCount = await firstRowCells.count();
    // Checkbox + Name + Description + Price + Status + Actions
    expect(cellCount).toBeGreaterThanOrEqual(5);
  });

  test('has primary action button in toolbar', async ({ page }) => {
    const primaryAction = page.locator('.toolbar-primary-action');
    await expect(primaryAction).toBeVisible();
    await expect(primaryAction).toBeEnabled();
  });

  test('row has action buttons (view, edit, delete)', async ({ page }) => {
    const firstRow = page.locator('#products-table tbody tr:first-child');
    const viewLink = firstRow.locator('a.action-btn.view');
    const editBtn = firstRow.locator('.action-btn.edit');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(viewLink).toBeVisible();
    await expect(editBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });

  test('view link navigates to service detail', async ({ page }) => {
    const viewLink = page.locator('#products-table tbody tr:first-child a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toContain('/app/services/detail/');
  });

  test('shows pagination with entry count', async ({ page }) => {
    const pagination = page.locator('.table-footer, .pagination-info');
    await expect(pagination).toBeVisible();
  });
});

test.describe('CEN-PRD-002: Services Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();
  });

  test('opens drawer when primary action clicked', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Verify form fields exist by ID
    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#price')).toBeVisible();
  });

  test('drawer has all required form fields', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#description')).toBeVisible();
    await expect(page.locator('#price')).toBeVisible();
    await expect(page.locator('#currency')).toBeVisible();
  });

  test('creates service via drawer form', async ({ page }) => {
    const ts = Date.now();

    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill required fields
    await page.locator('#name').fill(`TestService${ts}`);
    await page.locator('#description').fill(`E2E test service created at ${ts}`);
    await page.locator('#price').fill('100.00');

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Verify drawer closes (waits for HTMX settle + sheet close animation)
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 15000 });
  });

  test('cancel closes drawer without creating', async ({ page }) => {
    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Fill something
    await page.locator('#name').fill('ShouldNotSave');

    // Cancel — use the secondary button in sheet footer
    await page.locator('#sheet .sheet-footer .btn-secondary').click();

    // Drawer should close
    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });
  });
});

test.describe('CEN-PRD-003: Services Edit via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();
  });

  test('opens edit drawer with pre-filled data', async ({ page }) => {
    const editBtn = page.locator('#products-table tbody tr:first-child .action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Name should be pre-filled
    const name = await page.locator('#name').inputValue();
    expect(name.length).toBeGreaterThan(0);
  });

  test('saves edit and closes drawer', async ({ page }) => {
    const editBtn = page.locator('#products-table tbody tr:first-child .action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();

    // Modify a field
    const ts = Date.now();
    await page.locator('#name').fill(`Updated${ts}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Verify drawer closes (waits for HTMX settle + sheet close animation)
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 15000 });
  });
});

test.describe('CEN-PRD-004: Service Detail Page', () => {
  test('detail page loads and renders correctly', async ({ page }) => {
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();

    const viewLink = page.locator('#products-table tbody tr:first-child a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toBeTruthy();

    await page.goto(href!);

    const h1 = page.locator('h1').first();
    await expect(h1).toBeVisible({ timeout: 10000 });
    const h1Text = await h1.textContent();
    expect(h1Text!.trim().length).toBeGreaterThan(0);

    const bodyText = await page.textContent('body');
    expect(bodyText).not.toContain('Page content not available');

    const detailLayout = page.locator('.detail-header, .detail-layout, .info-grid');
    await expect(detailLayout.first()).toBeVisible({ timeout: 5000 });
  });
});

test.describe('CEN-PRD-LIFECYCLE: Service Full Lifecycle', () => {
  test('creates, edits, views detail, and deletes a service', async ({ page }) => {
    const ts = Date.now();

    // 1. Navigate to list page
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();

    // 2. Add new record via drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();
    await waitForHtmxSettle(page);

    await page.locator('#name').fill(`E2EService${ts}`);
    await page.locator('#description').fill(`Service created by lifecycle test at ${ts}`);
    await page.locator('#price').fill('250.00');

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();
    await waitForHtmxSettle(page);
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 15000 });

    // 3. Find the newly created record
    await page.waitForTimeout(500);
    await page.reload();
    await expect(page.locator('#products-table')).toBeVisible();

    const rows = page.locator('#products-table tbody tr[data-id]');
    const rowCount = await rows.count();
    expect(rowCount).toBeGreaterThan(0);

    // Service names may be normalized to different case by the backend (e.g. "E2eservice")
    const serviceNameLower = `e2eservice${ts}`;
    let targetRowIndex = -1;
    for (let i = 0; i < rowCount; i++) {
      const rowText = await rows.nth(i).textContent();
      if (rowText?.toLowerCase().includes(serviceNameLower)) {
        targetRowIndex = i;
        break;
      }
    }
    expect(targetRowIndex).toBeGreaterThanOrEqual(0);
    const targetRow = rows.nth(targetRowIndex);

    // 4. Edit the record
    await targetRow.locator('.action-btn.edit').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible();
    await waitForHtmxSettle(page);

    const nameValue = await page.locator('#name').inputValue();
    expect(nameValue.length).toBeGreaterThan(0);

    await page.locator('#description').fill(`Edited by lifecycle test at ${ts}`);
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();
    await waitForHtmxSettle(page);
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 15000 });

    // 5. View detail page
    await page.reload();
    await expect(page.locator('#products-table')).toBeVisible();

    const rowsAfterEdit = page.locator('#products-table tbody tr[data-id]');
    let detailRowIndex = -1;
    for (let i = 0; i < await rowsAfterEdit.count(); i++) {
      const rowText = await rowsAfterEdit.nth(i).textContent();
      if (rowText?.toLowerCase().includes(serviceNameLower)) {
        detailRowIndex = i;
        break;
      }
    }
    expect(detailRowIndex).toBeGreaterThanOrEqual(0);

    const viewLink = rowsAfterEdit.nth(detailRowIndex).locator('a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toBeTruthy();

    await page.goto(href!);

    // 6. Verify detail page renders
    const h1 = page.locator('h1').first();
    await expect(h1).toBeVisible({ timeout: 10000 });
    const h1Text = await h1.textContent();
    expect(h1Text!.trim().length).toBeGreaterThan(0);

    const bodyText = await page.textContent('body');
    expect(bodyText).not.toContain('Page content not available');

    const detailLayout = page.locator('.detail-header, .detail-layout, .info-grid');
    await expect(detailLayout.first()).toBeVisible({ timeout: 5000 });

    // 7. Navigate back and delete the test record
    await page.goto('/app/services/list/active');
    await expect(page.locator('#products-table')).toBeVisible();

    const rowsForDelete = page.locator('#products-table tbody tr[data-id]');
    for (let i = 0; i < await rowsForDelete.count(); i++) {
      const rowText = await rowsForDelete.nth(i).textContent();
      if (rowText?.toLowerCase().includes(serviceNameLower)) {
        const deleteBtn = rowsForDelete.nth(i).locator('.action-btn.delete');
        if (await deleteBtn.isVisible()) {
          await deleteBtn.click();
          const confirmBtn = page.locator('#dialog.visible .dialog-btn-confirm');
          await expect(confirmBtn).toBeVisible({ timeout: 5000 });
          await confirmBtn.click();
          await waitForHtmxSettle(page);
        }
        break;
      }
    }
  });
});
