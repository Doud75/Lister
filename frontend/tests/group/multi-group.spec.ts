import { test, expect, type Page } from '@playwright/test';

async function loginMultiGroupUser(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('multiGroupUser');
    await page.getByLabel('Mot de passe').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

test.describe('Multi-Group Functionality', () => {
    test('should allow a user to switch between bands and see isolated data', async ({ page }) => {
        await loginMultiGroupUser(page);

        await expect(page.getByRole('heading', { name: 'Band A' })).toBeVisible();

        await expect(page.getByText('Setlist A')).toBeVisible();
        await expect(page.getByText('Setlist B')).toBeHidden();

        await page.getByRole('link', { name: 'Gérer les chansons' }).click();
        await page.waitForURL('/song');
        await expect(page.getByText('Chanson A1')).toBeVisible();
        await expect(page.getByText('Chanson B1')).toBeHidden();

        await page.goto('/');

        const bandSelector = page.locator('select[name="bandId"]');
        await bandSelector.selectOption({ label: 'Band B' });

        await page.waitForURL('/');

        await expect(page.getByRole('heading', { name: 'Band B' })).toBeVisible();
        await expect(page.getByText('Setlist B')).toBeVisible();
        await expect(page.getByText('Setlist A')).toBeHidden();

        await page.getByRole('link', { name: 'Gérer les chansons' }).click();
        await page.waitForURL('/song');
        await expect(page.getByText('Chanson B1')).toBeVisible();
        await expect(page.getByText('Chanson A1')).toBeHidden();

        await page.reload();

        await expect(page.getByText('Chanson B1')).toBeVisible();
        await expect(page.getByText('Chanson A1')).toBeHidden();

        await page.goto('/');
        await expect(page.getByRole('heading', { name: 'Band B' })).toBeVisible();
    });
});