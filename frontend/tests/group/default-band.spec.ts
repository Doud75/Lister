import { test, expect, type Page } from '@playwright/test';

async function signup(page: Page, username: string, password: string) {
    await page.goto('/signup');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe', { exact: true }).fill(password);
    await page.getByRole('button', { name: 'Créer mon compte' }).click();
    await page.waitForURL('/dashboard');
}

async function createBand(page: Page, bandName: string) {
    await page.goto('/dashboard');
    await page.getByRole('button', { name: 'Créer un groupe' }).click();
    await page.locator('#band-name').fill(bandName);
    await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();
    await page.waitForURL('/');
    await page.goto('/dashboard');
}

async function logout(page: Page) {
    await page.goto('/logout');
}

async function login(page: Page, username: string, password: string) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe').fill(password);
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

test.describe('Default Band', () => {
    // ── Cas 1 : cliquer l'étoile change le groupe par défaut ──────────────
    test('clicking star on a non-default band makes it the default', async ({ page }) => {
        const timestamp = Date.now();
        const username = `star_user_${timestamp}`;
        const bandA = `Star Band A ${timestamp}`;
        const bandB = `Star Band B ${timestamp}`;

        await signup(page, username, 'StrongPass1!');
        await createBand(page, bandA);
        await createBand(page, bandB);

        await page.goto('/dashboard');

        // bandA devrait avoir l'étoile pleine (défaut)
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true }).first()).toBeVisible();

        // Cliquer sur l'étoile vide de bandB
        await page.getByRole('button', { name: `Définir ${bandB} comme groupe par défaut`, exact: true }).click();

        // Maintenant bandB a l'étoile pleine, bandA a l'étoile vide
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
        await expect(page.getByRole('button', { name: `Définir ${bandA} comme groupe par défaut`, exact: true })).toBeVisible();
    });

    // ── Cas 2 : après logout/login → redirigé sur le groupe par défaut ────
    test('after login, active band is the default band', async ({ page }) => {
        const timestamp = Date.now();
        const username = `default_login_user_${timestamp}`;
        const bandA = `Default Login Band A ${timestamp}`;
        const bandB = `Default Login Band B ${timestamp}`;

        await signup(page, username, 'StrongPass1!');
        await createBand(page, bandA);
        await createBand(page, bandB);

        // Définir bandB comme défaut
        await page.goto('/dashboard');
        await page.getByRole('button', { name: `Définir ${bandB} comme groupe par défaut`, exact: true }).click();
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();

        // Logout et reconnexion
        await logout(page);
        await login(page, username, 'StrongPass1!');

        // Le groupe actif doit être bandB (défaut)
        await page.goto('/dashboard');
        await expect(page.getByText('✓ Groupe actif').first()).toBeVisible();
        // bandB doit avoir l'étoile pleine
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
    });

    // ── Cas 3 : quitter le groupe par défaut → un autre devient défaut ────
    test('leaving default band auto-reassigns default to another band', async ({ page }) => {
        const timestamp = Date.now();
        const username = `leave_default_user_${timestamp}`;
        const bandA = `Leave Default Band A ${timestamp}`;
        const bandB = `Leave Default Band B ${timestamp}`;

        await signup(page, username, 'StrongPass1!');
        await createBand(page, bandA);
        await createBand(page, bandB);

        // S'assurer que bandA est le défaut — si ce n'est pas déjà le cas, le définir
        await page.goto('/dashboard');
        const setBandADefault = page.getByRole('button', { name: `Définir ${bandA} comme groupe par défaut`, exact: true });
        if (await setBandADefault.isVisible()) {
            await setBandADefault.click();
            await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
        }

        // Quitter bandA (le groupe par défaut)
        await page.getByRole('button', { name: `Quitter ${bandA}` }).click();
        await page.getByRole('button', { name: 'Quitter le groupe' }).click();
        await page.waitForURL(/\/dashboard/);

        // bandB doit maintenant avoir l'étoile pleine (auto-réassignée)
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
        // bandA ne doit plus apparaître
        await expect(page.getByRole('heading', { name: bandA, level: 2 })).not.toBeVisible();
    });
});