import { test, expect } from '@playwright/test';

test('should load dashboard with fallback data when API fails', async ({ page }) => {
    // Mock API responses to simulate backend failures
    await page.route('**/api/v1/analytics/balance-history', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    await page.route('**/api/v1/newcoins', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    // Navigate to the dashboard
    await page.goto('/');

    // Verify page title is visible
    await expect(page.locator('header')).toContainText('DASHBOARD');

    // Verify fallback data is displayed
    await expect(page.locator('text=$0.00')).toBeVisible();

    // Verify connection status shows disconnected
    await expect(page.locator('text=Disconnected')).toBeVisible({ timeout: 10000 });

    // Verify chart is displayed even with no data
    await expect(page.locator('canvas')).toBeVisible();
  });

  test('should display error toast when API fails', async ({ page }) => {
    // Navigate to the dashboard
    await page.goto('/');

    // Wait for toast to appear
    const toast = page.locator('div[role="status"]');

    // Verify toast is displayed
    await expect(toast).toBeVisible({ timeout: 10000 });

    // Verify toast contains error message
    await expect(toast).toContainText('Error');
  });

  test('should show connection status indicator', async ({ page }) => {
    // Navigate to the dashboard
    await page.goto('/');

    // Verify connection status component is visible
    const connectionStatus = page.locator('header').getByText(/Connected|Disconnected/);

    // Verify connection status is displayed
    await expect(connectionStatus).toBeVisible({ timeout: 10000 });
  });
