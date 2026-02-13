// frontend/tests/settings/members.spec.ts

import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel('Username').fill('testuser');
    await page.getByLabel('Password').fill('Password123!');
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

// AJOUT: On utilise .serial pour garantir l'ordre d'exécution
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
        // CORRECTION: Utilisation de { exact: true } pour cibler uniquement le rôle "member" et pas le nom "memberuser"
        await expect(memberUserRow.getByText('member', { exact: true })).toBeVisible();
    });

    test('should successfully remove a member', async ({ page }) => {
        const memberUserRow = page.locator('li', { hasText: 'memberuser' });
        await expect(memberUserRow).toBeVisible();

        await memberUserRow.getByRole('button', { name: 'Supprimer memberuser' }).click();

        await expect(memberUserRow).toBeHidden();

        // Verify persistence
        await page.reload();
        await expect(page.locator('li', { hasText: 'memberuser' })).toBeHidden();
        await expect(page.locator('li', { hasText: 'testuser' })).toBeVisible();
    });

    test('should invite an existing user not in the band', async ({ page }) => {
        await page.getByLabel("Nom d'utilisateur").fill('multiGroupUser');
        await page.getByRole('button', { name: 'Ajouter / Inviter' }).click();

        await expect(page.locator('li', { hasText: 'multiGroupUser' })).toBeVisible();
    });

    test('should create and invite a new user', async ({ page }) => {
        const newUser = `newUser_${Date.now()}`;

        await page.getByLabel("Nom d'utilisateur").fill(newUser);
        await page.getByRole('button', { name: 'Ajouter / Inviter' }).click();

        await expect(page.getByText('Utilisateur non trouvé. Veuillez définir un mot de passe pour créer son compte.')).toBeVisible();

        await page.getByLabel('Mot de passe temporaire').fill('Password123!');
        await page.getByRole('button', { name: 'Créer et Inviter' }).click();

        await expect(page.locator('li', { hasText: newUser })).toBeVisible();
    });
});