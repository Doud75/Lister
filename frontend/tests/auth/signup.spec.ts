import { test, expect } from '@playwright/test';

test.describe('Signup Page', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/signup');
    });

    test('should validate username Requirements', async ({ page }) => {
        await page.getByLabel('Nom du groupe').fill('Test Band');
        await page.getByLabel('Mot de passe', { exact: true }).fill('Valid1Password!');

        await page.getByLabel('Votre nom d\'utilisateur').fill('yo');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le nom d\'utilisateur doit contenir au moins 3 caractères.')).toBeVisible();

        await page.getByLabel('Votre nom d\'utilisateur').fill('bad user!');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le nom d\'utilisateur ne peut contenir que des lettres, des chiffres et des underscores.')).toBeVisible();
    });

    test('should validate password Requirements', async ({ page }) => {
        await page.getByLabel('Nom du groupe').fill('Test Band');
        await page.getByLabel('Votre nom d\'utilisateur').fill('ValidUser');

        await page.getByLabel('Mot de passe', { exact: true }).fill('weak');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le mot de passe doit contenir au moins 8 caractères.')).toBeVisible();

        await page.getByLabel('Mot de passe', { exact: true }).fill('nocaps1!');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le mot de passe doit contenir au moins une majuscule.')).toBeVisible();

        await page.getByLabel('Mot de passe', { exact: true }).fill('NoNumbers!');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le mot de passe doit contenir au moins un chiffre.')).toBeVisible();
        
        await page.getByLabel('Mot de passe', { exact: true }).fill('NoSpecial1');
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();
        await expect(page.getByText('Le mot de passe doit contenir au moins un caractère spécial.')).toBeVisible();
    });

    test('should signup successfully with valid data', async ({ page }) => {
        const username = `TestUser_${Date.now()}`;
        await page.getByLabel('Nom du groupe').fill('The Test Band');
        await page.getByLabel('Votre nom d\'utilisateur').fill(username);
        await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
        
        await page.getByRole('button', { name: 'Créer le groupe et le compte' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('heading', { name: 'The Test Band' })).toBeVisible();
    });
});
