import { test, expect, type Page } from '@playwright/test';

async function login(page: Page, user: string, pass: string) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill(user);
    await page.getByLabel('Mot de passe').fill(pass);
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

async function logout(page: Page) {
    await page.goto('/');
    await page.getByRole('button', { name: "Ouvrir le menu du profil" }).click();
    await page.getByRole('menuitem', { name: 'Déconnexion' }).click();
    await page.waitForURL('/login');
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
        const newPassword = 'NewPassword123!';

        await page.getByLabel('Ancien mot de passe').fill('password123');
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill(newPassword);
        await page.getByLabel('Confirmer le nouveau mot de passe').fill(newPassword);
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();

        await logout(page);

        await page.getByLabel("Nom d'utilisateur").fill(userForTest);
        await page.getByLabel('Mot de passe').fill('password123');
        await page.getByRole('button', { name: 'Se connecter' }).click();
        await expect(page.getByText('invalid credentials')).toBeVisible();

        await page.getByLabel('Mot de passe').fill(newPassword);
        await page.getByRole('button', { name: 'Se connecter' }).click();
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();
    });
});