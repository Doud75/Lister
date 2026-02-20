import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel('Username').fill('testuser');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

const SONG_DETAIL_URL = '/song/1';
const SETLIST_DETAIL_URL = '/setlist/1';

test.describe('Song Detail Page', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(SONG_DETAIL_URL);
        await expect(page.getByRole('heading', { name: 'Song Title 1' })).toBeVisible();
    });

    test('should display song metadata', async ({ page }) => {
        await expect(page.getByText('Test Album')).toBeVisible();
        await expect(page.getByText('120 BPM')).toBeVisible();
        await expect(page.getByText('3m 05s')).toBeVisible();
    });

    test('should display lyrics placeholder when no lyrics are set', async ({ page }) => {
        await expect(page.getByText('Aucune parole renseignée.')).toBeVisible();
    });

    test('should navigate back to song library via back button', async ({ page }) => {
        await page.getByRole('link', { name: '← Retour' }).click();

        await page.waitForURL('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
    });

    test('should navigate to edit page via Modifier button', async ({ page }) => {
        await page.getByRole('link', { name: 'Modifier' }).click();

        await page.waitForURL('/song/1/edit**');
        await expect(page.getByRole('heading', { name: 'Modifier: Song Title 1' })).toBeVisible();
    });

    test('should navigate to detail page from song library', async ({ page }) => {
        await page.goto('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();

        const songItem = page.locator('li', { hasText: 'Song Title 1' });
        const songLink = songItem.getByRole('link', { name: 'Song Title 1', exact: true });
        await expect(songLink).toBeVisible();
        await songLink.click();

        await page.waitForURL('/song/1');
        await expect(page.getByRole('heading', { name: 'Song Title 1' })).toBeVisible();
    });

    test('should use from param to navigate back to setlist', async ({ page }) => {
        await page.goto(`/song/1?from=${SETLIST_DETAIL_URL}`);

        const backLink = page.getByRole('link', { name: '← Retour' });
        await expect(backLink).toHaveAttribute('href', SETLIST_DETAIL_URL);

        await backLink.click();
        await page.waitForURL(SETLIST_DETAIL_URL);
        await expect(page.getByRole('heading', { name: 'Test Setlist' })).toBeVisible();
    });
});

test.describe('Song Detail - Navigation depuis une setlist', () => {
    test('should navigate to song detail from a setlist with correct back link', async ({ page }) => {
        await login(page);
        await page.goto(SETLIST_DETAIL_URL);

        await page.getByRole('link', { name: 'Song Title 1' }).click();

        await page.waitForURL(`/song/1?from=${SETLIST_DETAIL_URL}`);
        await expect(page.getByRole('heading', { name: 'Song Title 1' })).toBeVisible();

        const backLink = page.getByRole('link', { name: '← Retour' });
        await expect(backLink).toHaveAttribute('href', SETLIST_DETAIL_URL);
    });
});

test.describe('Edit Song - Redirect dynamique via from', () => {
    test('should redirect to song detail after save when coming from detail page', async ({ page }) => {
        await login(page);
        await page.goto('/song/1/edit?from=/song/1');

        await expect(page.getByRole('link', { name: 'Annuler' })).toHaveAttribute('href', '/song/1');

        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await page.waitForURL('/song/1');
        await expect(page.getByRole('heading', { name: 'Song Title 1' })).toBeVisible();
    });
});
