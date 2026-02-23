import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.getByLabel('Mot de passe').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

const NEW_SETLIST_URL = '/setlist/new';

test.describe('Create New Setlist Page', () => {

    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(NEW_SETLIST_URL);
        await expect(page.getByRole('heading', { name: 'Create a New Setlist' })).toBeVisible();
    });

    test('should successfully create a new setlist and redirect', async ({ page }) => {
        const newSetlistName = `My Awesome Setlist ${Date.now()}`;

        await page.getByLabel('Setlist Name').fill(newSetlistName);
        await page.getByRole('button', { name: 'Create Setlist' }).click();

        await page.waitForURL(/\/setlist\/\d+/);

        await expect(page.getByRole('heading', { name: newSetlistName })).toBeVisible();
    });

    test('should not submit and should stay on the page if the setlist name is empty', async ({ page }) => {
        await page.getByRole('button', { name: 'Create Setlist' }).click();

        await page.waitForTimeout(500);

        await expect(page).toHaveURL(NEW_SETLIST_URL);

        const nameInput = page.getByLabel('Setlist Name');
        const isInvalid = await nameInput.evaluate(element => !(element as HTMLInputElement).checkValidity());
        expect(isInvalid).toBe(true);
    });

    test('should navigate back to home when cancel is clicked', async ({ page }) => {
        await page.getByRole('link', { name: 'Cancel' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();
    });
});