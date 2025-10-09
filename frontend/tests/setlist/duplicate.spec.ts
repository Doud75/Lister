import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel('Username').fill('testuser');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

const SETLIST_TO_DUPLICATE_ID = 1;
const SETLIST_URL = `/setlist/${SETLIST_TO_DUPLICATE_ID}`;

test.describe('Setlist Duplication', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(SETLIST_URL);
        await expect(page.getByRole('heading', { name: 'Test Setlist' })).toBeVisible();
    });

    test('should open the duplication modal and pre-fill fields', async ({ page }) => {
        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Dupliquer' }).click();

        await expect(page.getByRole('dialog')).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Dupliquer la setlist' })).toBeVisible();

        await expect(page.getByLabel('Nouveau nom')).toHaveValue('Copie de Test Setlist');
    });

    test('should successfully duplicate a setlist and redirect', async ({ page }) => {
        const newName = `Copied Setlist - ${Date.now()}`;

        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Dupliquer' }).click();

        await expect(page.getByRole('dialog')).toBeVisible();
        await page.getByLabel('Nouveau nom').fill(newName);
        await page.getByRole('button', { name: 'Créer la copie' }).click();

        await page.waitForURL(/\/setlist\/\d+/);
        await expect(page).not.toHaveURL(SETLIST_URL);

        await expect(page.getByRole('heading', { name: newName })).toBeVisible();

        const items = page.locator('ul[data-testid="setlist-items"] > li');
        await expect(items).toHaveCount(3);
        await expect(items.nth(0)).toContainText('Song Title 1');
        await expect(items.nth(1)).toContainText('Test Interlude');
        await expect(items.nth(2)).toContainText('Song Title 2');
    });

    test('should close the modal on cancel', async ({ page }) => {
        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Dupliquer' }).click();

        await expect(page.getByRole('dialog')).toBeVisible();
        await page.getByRole('button', { name: 'Annuler' }).click();
        await expect(page.getByRole('dialog')).toBeHidden();
    });
});