import { test, expect, type Page } from '@playwright/test';

async function logout(page: Page) {
    await page.goto('/');
    await page.getByRole('button', { name: 'Ouvrir le menu du profil' }).click();
    await page.getByRole('menuitem', { name: 'Déconnexion' }).click();
    await page.waitForURL('/login');
}

async function signupAndCreateBand(page: Page, username: string, bandName: string) {
    await page.goto('/signup');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
    await page.getByRole('button', { name: 'Créer mon compte' }).click();
    await page.waitForURL('/dashboard');
    await page.getByRole('button', { name: 'Créer un groupe' }).click();
    await page.locator('#band-name').fill(bandName);
    await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();
    await page.waitForURL('/');
}

test.describe('Duplicate Band Name', () => {
    test('should allow creating two bands with the same name', async ({ page }) => {
        const bandName = 'The Beatles';
        const timestamp = Date.now();

        // --- Premier compte + groupe ---
        await signupAndCreateBand(page, `beatles_user_1_${timestamp}`, bandName);
        await expect(page.getByRole('heading', { name: bandName })).toBeVisible();

        // --- Se déconnecter via le menu profil ---
        await logout(page);

        // --- Deuxième compte + groupe avec le MÊME nom de groupe ---
        await signupAndCreateBand(page, `beatles_user_2_${timestamp}`, bandName);

        // Doit réussir : pas d'erreur "nom de groupe déjà pris"
        await expect(page.getByRole('heading', { name: bandName })).toBeVisible();
    });
});
