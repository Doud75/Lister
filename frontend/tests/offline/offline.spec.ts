import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.locator('#password').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

async function waitForServiceWorker(page: Page) {
    await page.evaluate(() => navigator.serviceWorker.ready);
}

test.describe('Mode hors-ligne', () => {
    test('affiche la pilule "Hors ligne" quand le réseau est coupé et la retire au retour', async ({ page, context }) => {
        await login(page);
        await waitForServiceWorker(page);

        const indicator = page.getByTestId('offline-indicator');
        await expect(indicator).not.toBeVisible();

        await context.setOffline(true);
        await expect(indicator).toBeVisible();
        await expect(indicator).toContainText('Hors ligne');

        await context.setOffline(false);
        await expect(indicator).not.toBeVisible();
    });

    test('sert une page déjà visitée depuis le cache quand hors-ligne', async ({ page, context }) => {
        await login(page);
        await page.goto('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
        await waitForServiceWorker(page);

        await context.setOffline(true);
        await page.reload();

        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
    });

    test('affiche un toast "pas de connexion réseau" pour une page jamais visitée', async ({ page, context }) => {
        await login(page);
        await page.goto('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
        await waitForServiceWorker(page);

        await context.setOffline(true);

        // Navigation client-side vers un détail de chanson jamais visité
        // (le détail utilise un client load qui appelle /api/song/:id → SW retourne 503 → toast)
        // On exclut /song/new qui n'a pas de load et ne déclencherait pas d'erreur
        const songLink = page.locator('a[href^="/song/"]:not([href="/song/new"])').first();
        await songLink.click();

        await expect(page.getByText('Pas de connexion réseau')).toBeVisible();
    });
});
