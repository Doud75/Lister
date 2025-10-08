// frontend/tests/settings/account.spec.ts

import { test, expect, type Page } from '@playwright/test';

async function login(page: Page, user: string, pass: string) {
    await page.goto('/login');
    await page.getByLabel('Username').fill(user);
    await page.getByLabel('Password').fill(pass);
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

test.describe.serial('Settings - Account Page', () => {
    // On se connecte en tant que notre nouvel utilisateur dédié
    const userForTest = 'passwordChangeUser';

    test.beforeEach(async ({ page }) => {
        await login(page, userForTest, 'password123');
        await page.goto('/settings/account');
    });

    test('should display the change password form', async ({ page }) => {
        await expect(page.getByRole('heading', { name: 'Mon Compte' })).toBeVisible();
        await expect(page.getByLabel('Ancien mot de passe')).toBeVisible();
        await expect(page.getByLabel('Nouveau mot de passe', { exact: true })).toBeVisible();
        await expect(page.getByLabel('Confirmer le nouveau mot de passe')).toBeVisible();
    });

    test('should show an error for incorrect current password', async ({ page }) => {
        await page.getByLabel('Ancien mot de passe').fill('wrongpassword');
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill('newpassword123');
        await page.getByLabel('Confirmer le nouveau mot de passe').fill('newpassword123');
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await expect(page.getByText('Votre ancien mot de passe est incorrect.')).toBeVisible();
        await expect(page).toHaveURL('/settings/account');
    });

    test('should successfully change password and allow re-login', async ({ page }) => {
        const newPassword = 'newpassword456';

        await page.getByLabel('Ancien mot de passe').fill('password123');
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill(newPassword);
        await page.getByLabel('Confirmer le nouveau mot de passe').fill(newPassword);
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();

        // Logout
        await page.goto('/logout');
        await page.waitForURL('/login');

        // Try to login with old password (should fail)
        await page.getByLabel('Username').fill(userForTest);
        await page.getByLabel('Password').fill('password123');
        await page.getByRole('button', { name: 'Log In' }).click();
        await expect(page.getByText('invalid credentials')).toBeVisible();

        // Login with new password (should succeed)
        await page.getByLabel('Password').fill(newPassword);
        await page.getByRole('button', { name: 'Log In' }).click();
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();
    });
});