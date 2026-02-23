import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.locator('#password').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

const EDIT_SONG_URL = '/song/1/edit'; // On édite la chanson avec l'ID 1

test.describe('Edit Song Page', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(EDIT_SONG_URL);
        await expect(page.getByRole('heading', { name: 'Modifier: Song Title 1' })).toBeVisible();
    });

    test('should successfully update a song and redirect to library', async ({ page }) => {
        const updatedTempo = '135';

        await page.getByLabel('Tempo (BPM)').fill(updatedTempo);
        await page.getByRole('button', { name: 'Sauvegarder' }).click();

        await page.waitForURL('/song');

        await expect(page.getByText("Song Title 1")).toBeVisible();

        const updatedSongItem = page.locator('li', { hasText: "Song Title 1" });
        await expect(updatedSongItem.getByText(`Tempo: ${updatedTempo} BPM`)).toBeVisible();
    });

    test('should navigate back to song library when cancel is clicked', async ({ page }) => {
        await page.getByRole('link', { name: 'Annuler' }).click();

        await page.waitForURL('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
    });
});