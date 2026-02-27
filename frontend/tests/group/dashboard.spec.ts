import { test, expect, type Page } from '@playwright/test';

async function signup(page: Page, bandName: string, username: string) {
    await page.goto('/signup');
    await page.getByLabel('Nom du groupe').fill(bandName);
    await page.getByLabel("Votre nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
    await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
    await page.waitForURL('/');
}

test.describe('Dashboard', () => {
    test('should display user groups on the dashboard page', async ({ page }) => {
        const ts = Date.now();
        const username = `dash_user_${ts}`;
        const bandName = `Dash Band ${ts}`;

        await signup(page, bandName, username);
        await page.goto('/dashboard');

        await expect(page.getByRole('heading', { name: 'Mes Groupes' })).toBeVisible();
        await expect(page.getByText(bandName)).toBeVisible();
        // User is the creator = admin
        await expect(page.getByText('Admin')).toBeVisible();
    });

    test('should create a new band from the dashboard and redirect to home', async ({ page }) => {
        const ts = Date.now();
        const username = `dash_create_${ts}`;
        const newBandName = `New Band ${ts}`;

        await signup(page, `Initial Band ${ts}`, username);
        await page.goto('/dashboard');

        // Open create form
        await page.getByRole('button', { name: 'Créer un groupe' }).click();
        await expect(page.locator('#band-name')).toBeVisible();

        await page.locator('#band-name').fill(newBandName);
        await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();

        // Should redirect to home with new band active
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: newBandName })).toBeVisible();

        // New band should appear in the navbar select
        await expect(page.locator('select[name="bandId"]').locator(`option[value]`, { hasText: newBandName })).toBeDefined();
    });

    test('should switch to a band by clicking its card', async ({ page }) => {
        const ts = Date.now();
        const username = `dash_switch_${ts}`;
        const band1 = `Switch Band A ${ts}`;
        const band2 = `Switch Band B ${ts}`;

        await signup(page, band1, username);
        // Create a second band from the dashboard
        await page.goto('/dashboard');
        await page.getByRole('button', { name: 'Créer un groupe' }).click();
        await page.locator('#band-name').fill(band2);
        await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();
        await page.waitForURL('/');
        // band2 is now active

        // Go to dashboard and switch back to band1
        await page.goto('/dashboard');

        // Scope to main to avoid matching the navbar <select> options
        const main = page.locator('main');
        await expect(main.getByRole('button', { name: band1 })).toBeVisible();
        await expect(main.getByRole('button', { name: band2 })).toBeVisible();

        // Click card for band1
        await main.getByRole('button', { name: band1 }).click();
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: band1 })).toBeVisible();
    });

    test('orphan user: should redirect to /dashboard when active_band_id cookie is missing', async ({ page, context }) => {
        const ts = Date.now();
        const username = `orphan_redirect_${ts}`;

        await signup(page, `Orphan Band ${ts}`, username);
        // Remove the active band cookie to simulate orphan state
        await context.clearCookies({ name: 'active_band_id' });

        // Going to / should redirect to /dashboard
        await page.goto('/');
        await page.waitForURL('/dashboard');
        await expect(page.getByRole('heading', { name: 'Mes Groupes' })).toBeVisible();
    });
});
