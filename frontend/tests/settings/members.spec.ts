// frontend/tests/settings/members.spec.ts

import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.locator('#password').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

test.describe.serial('Settings - Members Page (Admin)', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto('/settings/members');
    });

    test('should display the list of current members', async ({ page }) => {
        await expect(page.getByRole('heading', { name: 'Gérer les membres du groupe' })).toBeVisible();

        const adminUserRow = page.locator('li', { hasText: 'testuser' });
        await expect(adminUserRow).toBeVisible();
        await expect(adminUserRow.getByText('admin')).toBeVisible();

        const memberUserRow = page.locator('li', { hasText: 'memberuser' });
        await expect(memberUserRow).toBeVisible();
        await expect(memberUserRow.getByText('member', { exact: true })).toBeVisible();
    });

    test('should successfully remove a member', async ({ page }) => {
        const memberUserRow = page.locator('li', { hasText: 'memberuser' });
        await expect(memberUserRow).toBeVisible();

        await memberUserRow.getByRole('button', { name: 'Supprimer memberuser' }).click();

        await expect(memberUserRow).toBeHidden();

        await page.reload();
        await expect(page.locator('li', { hasText: 'memberuser' })).toBeHidden();
        await expect(page.locator('li', { hasText: 'testuser' })).toBeVisible();
    });

    test('should display the invitation link section', async ({ page }) => {
        await expect(page.getByRole('heading', { name: 'Inviter par lien' })).toBeVisible();
        await expect(page.getByRole('button', { name: 'Générer un lien' })).toBeVisible();
    });
});