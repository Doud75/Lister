// frontend/tests/settings/account.spec.ts

import { test, expect, type Page } from '@playwright/test';

async function login(page: Page, user: string, pass: string) {
    await page.goto('/login');
    await page.getByLabel('Username').fill(user);
    await page.getByLabel('Password').fill(pass);
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

// AJOUT: On utilise .serial pour forcer les tests de ce fichier à s'exécuter l'un après l'autre
test.describe.serial('Settings - Account Page', () => {
    test.beforeEach(async ({ page }) => {
        // Le seeder doit être réinitialisé avant chaque *fichier* de test, ce que fait `make test-*`.
        // Ici, on se connecte avant chaque test du fichier.
        await login(page, 'testuser', 'password123');
        await page.goto('/settings/account');
    });

    test('should display the change password form', async ({ page }) => {
        await expect(page.getByRole('heading', { name: 'Mon Compte' })).toBeVisible();
        await expect(page.getByLabel('Ancien mot de passe')).toBeVisible();
        // CORRECTION: Utilisation de { exact: true } pour éviter l'ambiguïté
        await expect(page.getByLabel('Nouveau mot de passe', { exact: true })).toBeVisible();
        await expect(page.getByLabel('Confirmer le nouveau mot de passe')).toBeVisible();
    });

    test('should show an error for incorrect current password', async ({ page }) => {
        await page.getByLabel('Ancien mot de passe').fill('wrongpassword');
        // CORRECTION: Utilisation de { exact: true }
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill('newpassword123');
        await page.getByLabel('Confirmer le nouveau mot de passe').fill('newpassword123');
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await expect(page.getByText('Votre ancien mot de passe est incorrect.')).toBeVisible();
        await expect(page).toHaveURL('/settings/account');
    });

    test('should successfully change password and allow re-login', async ({ page }) => {
        const newPassword = 'newpassword456';

        await page.getByLabel('Ancien mot de passe').fill('password123');
        // CORRECTION: Utilisation de { exact: true }
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill(newPassword);
        await page.getByLabel('Confirmer le nouveau mot de passe').fill(newPassword);
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();

        // Logout
        await page.goto('/logout');
        await page.waitForURL('/login');

        // Try to login with old password (should fail)
        await page.getByLabel('Username').fill('testuser');
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