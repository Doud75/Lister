import { test, expect, type Page } from '@playwright/test';

async function login(page: Page, user: string, pass: string) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill(user);
    await page.locator('#password').fill(pass);
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
        await page.locator('#password').fill('password123');
        await page.getByRole('button', { name: 'Se connecter' }).click();
        await expect(page.getByText('Identifiant ou mot de passe incorrect.')).toBeVisible();

        await page.locator('#password').fill(newPassword);
        await page.getByRole('button', { name: 'Se connecter' }).click();
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Testers' })).toBeVisible();
    });
});

test.describe('Settings - Account Page (Orphan User)', () => {
    test('should allow an orphan user to access settings and change password', async ({ page, context }) => {
        const ts = Date.now();
        const username = `orphan_account_${ts}`;
        const password = 'StrongPass1!';
        const newPassword = 'NewOrphanPass1!';

        // Signup creates a user with a band
        await page.goto('/signup');
        await page.getByLabel('Nom du groupe').fill(`Orphan Band ${ts}`);
        await page.getByLabel("Votre nom d'utilisateur").fill(username);
        await page.getByLabel('Mot de passe', { exact: true }).fill(password);
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await page.waitForURL('/');

        // Remove active_band_id to simulate orphan state
        await context.clearCookies({ name: 'active_band_id' });

        // / should redirect to /dashboard (no band)
        await page.goto('/');
        await page.waitForURL('/dashboard');

        // Settings account should still be accessible without a band
        await page.goto('/settings/account');
        await expect(page.getByRole('heading', { name: 'Mon Compte' })).toBeVisible();

        // Should be able to change the password
        await page.getByLabel('Ancien mot de passe').fill(password);
        await page.getByLabel('Nouveau mot de passe', { exact: true }).fill(newPassword);
        await page.getByLabel('Confirmer le nouveau mot de passe').fill(newPassword);
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        // After password change, the account page server redirects to '/'
        // Since no active band, we should end up at /dashboard
        await page.waitForURL('/dashboard');
    });
});