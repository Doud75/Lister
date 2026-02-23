import { test, expect } from '@playwright/test';

test.describe('Login Page', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/login');
    });

    test('should display an error message with invalid credentials', async ({ page }) => {
        await page.getByLabel("Nom d'utilisateur").fill('utilisateur_inexistant');
        await page.getByLabel('Mot de passe').fill('mauvaisMotDePasse1!');

        await page.getByRole('button', { name: 'Se connecter' }).click();

        await expect(
            page.getByText('Identifiant ou mot de passe incorrect.')
        ).toBeVisible();
    });

    test('should display an error message with wrong password', async ({ page }) => {
        await page.getByLabel("Nom d'utilisateur").fill('testuser');
        await page.getByLabel('Mot de passe').fill('WrongPassword9!');

        await page.getByRole('button', { name: 'Se connecter' }).click();

        await expect(
            page.getByText('Identifiant ou mot de passe incorrect.')
        ).toBeVisible();
    });

    test('should redirect to home on successful login', async ({ page }) => {
        await page.getByLabel("Nom d'utilisateur").fill('testuser');
        await page.getByLabel('Mot de passe').fill('password123');

        await page.getByRole('button', { name: 'Se connecter' }).click();

        await page.waitForURL('/');
    });

    test('should have a link to the signup page', async ({ page }) => {
        const signupLink = page.getByRole('link', { name: 'Créez un groupe' });
        await expect(signupLink).toBeVisible();
        await expect(signupLink).toHaveAttribute('href', '/signup');
    });
});
