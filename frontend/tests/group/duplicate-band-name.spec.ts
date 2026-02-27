import { test, expect, type Page } from '@playwright/test';

async function logout(page: Page) {
    await page.goto('/');
    await page.getByRole('button', { name: 'Ouvrir le menu du profil' }).click();
    await page.getByRole('menuitem', { name: 'Déconnexion' }).click();
    await page.waitForURL('/login');
}

test.describe('Duplicate Band Name', () => {
    test('should allow creating two bands with the same name', async ({ page }) => {
        const bandName = 'The Beatles';
        const timestamp = Date.now();

        // --- Premier compte + groupe ---
        await page.goto('/signup');
        await page.getByLabel('Nom du groupe').fill(bandName);
        await page.getByLabel("Votre nom d'utilisateur").fill(`beatles_user_1_${timestamp}`);
        await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: bandName })).toBeVisible();

        // --- Se déconnecter via le menu profil ---
        await logout(page);

        // --- Deuxième compte + groupe avec le MÊME nom de groupe ---
        await page.goto('/signup');
        await page.getByLabel('Nom du groupe').fill(bandName);
        await page.getByLabel("Votre nom d'utilisateur").fill(`beatles_user_2_${timestamp}`);
        await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();

        // Doit réussir : pas d'erreur "nom de groupe déjà pris"
        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: bandName })).toBeVisible();
    });
});
