import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel('Username').fill('testuser');
    await page.getByLabel('Password').fill('Password123!');
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

const NEW_SONG_URL = '/song/new';

test.describe('Create New Song Page', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(NEW_SONG_URL);
        await expect(page.getByRole('heading', { name: 'Add a New Song to Your Library' })).toBeVisible();
    });

    test('should successfully create a new song and redirect to library', async ({ page }) => {
        const newSongTitle = `My New Test Song ${Date.now()}`;

        await page.getByLabel('Song Title').fill(newSongTitle);
        await page.getByLabel('Album Name (optional)').fill('Test Creations');
        await page.getByLabel('Tempo (BPM)').fill('120');
        await page.getByRole('button', { name: 'Create Song' }).click();

        await page.waitForURL('/song');

        await expect(page.getByRole('heading', { name: 'Test Creations' })).toBeVisible();
        await expect(page.getByText(newSongTitle)).toBeVisible();
    });

    test('should show validation error if title is empty', async ({ page }) => {
        await page.getByRole('button', { name: 'Create Song' }).click();

        await page.waitForTimeout(500);
        await expect(page).toHaveURL(NEW_SONG_URL);

        const titleInput = page.getByLabel('Song Title');
        const isInvalid = await titleInput.evaluate(
            (element) => !(element as HTMLInputElement).checkValidity()
        );
        expect(isInvalid).toBe(true);
    });

    test('should navigate back to song library when cancel is clicked', async ({ page }) => {
        await page.getByRole('link', { name: 'Cancel' }).click();

        await page.waitForURL('/song');
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
    });
});