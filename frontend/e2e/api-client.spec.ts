import { test, expect } from '@playwright/test';

test('should handle API errors gracefully', async ({ page }) => {
    // Intercept all API calls and return errors
    await page.route('**/api/v1/**', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    // Navigate to the dashboard
    await page.goto('/');

    // Wait for the page to stabilize
    await page.waitForTimeout(2000);

    // Check console for errors
    const errors = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });

    // Verify the page doesn't crash
    await expect(page.locator('body')).toBeVisible();

    // Verify connection status shows disconnected
    await expect(page.locator('text=Disconnected')).toBeVisible({ timeout: 10000 });
  });

  test('should show fallback UI when backend is unavailable', async ({ page }) => {
    // Mock the status endpoint to simulate backend down
    await page.route('**/api/v1/status', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    // Navigate to the dashboard
    await page.goto('/');

    // Verify fallback UI elements are displayed
    await expect(page.locator('text=Disconnected')).toBeVisible({ timeout: 10000 });

    // Verify portfolio value shows fallback
    await expect(page.locator('text=$0.00')).toBeVisible();
  });
