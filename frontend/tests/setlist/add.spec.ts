import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.getByLabel('Mot de passe').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

async function createEmptySetlist(page: Page): Promise<number> {
    const response = await page.request.post('/api/setlist', {
        data: {
            name: `Test Setlist ${Date.now()}`,
            color: '#00ff00'
        }
    });
    expect(response.ok()).toBeTruthy();
    const setlist = await response.json();
    return setlist.id;
}

const getLibraryContainer = (page: Page) => page.locator('div.lg\\:order-first');
const getSetlistContainer = (page: Page) => page.locator('div.lg\\:order-last');
const getSetlistItems = (page: Page) => getSetlistContainer(page).locator('ul[data-testid="setlist-items"] > li');


test.describe('Setlist Add Item Page', () => {
    let setlistId: number;
    let addUrl: string;

    test.beforeEach(async ({ page }) => {
        await login(page);
        setlistId = await createEmptySetlist(page);
        addUrl = `/setlist/${setlistId}/add`;
        await page.goto(addUrl);
        await expect(page.getByRole('heading', { name: /Build Setlist:/ })).toBeVisible();
    });

    test('should display an empty setlist and the song library', async ({ page }) => {
        await expect(getSetlistContainer(page).getByText('Add items from your library to get started.')).toBeVisible();
        await expect(getLibraryContainer(page).getByRole('button', { name: /Songs \(\d+\)/ })).toBeVisible();
        await expect(getLibraryContainer(page).getByText('Song Title 1')).toBeVisible();
    });

    test('should add a song, then an interlude, and the order should be correct', async ({ page }) => {
        const library = getLibraryContainer(page);
        const setlistItems = getSetlistItems(page);

        await library.locator('li').filter({ hasText: 'Song Title 1' }).getByRole('button').click();
        await expect(setlistItems).toHaveCount(1);
        await expect(setlistItems.nth(0)).toContainText('Song Title 1');

        await library.locator('li').filter({ hasText: 'Another Song To Add' }).getByRole('button').click();
        await expect(setlistItems).toHaveCount(2);
        await expect(setlistItems.nth(1)).toContainText('Another Song To Add');

        await library.getByRole('button', { name: /Interludes/ }).click();
        await library.locator('li').filter({ hasText: 'Interlude To Add' }).getByRole('button').click();

        await expect(setlistItems).toHaveCount(3);
        await expect(setlistItems.nth(0)).toContainText('Song Title 1');
        await expect(setlistItems.nth(1)).toContainText('Another Song To Add');
        await expect(setlistItems.nth(2)).toContainText('Interlude To Add');
    });

    test('should add items and allow reordering', async ({ page }) => {
        const library = getLibraryContainer(page);
        const setlistItems = getSetlistItems(page);

        await library.locator('li').filter({ hasText: 'Song Title 1' }).getByRole('button').click();
        await library.locator('li').filter({ hasText: 'Another Song To Add' }).getByRole('button').click();
        await library.getByRole('button', { name: /Interludes/ }).click();
        await library.locator('li').filter({ hasText: 'Interlude To Add' }).getByRole('button').click();
        await expect(setlistItems).toHaveCount(3);

        const handle = setlistItems.filter({ hasText: 'Song Title 1' }).locator('[aria-label="Drag to reorder"]');
        const dropTarget = setlistItems.filter({ hasText: 'Interlude To Add' }).locator('[aria-label="Drag to reorder"]');;

        await handle.hover();
        await page.mouse.down();
        await dropTarget.hover();
        await dropTarget.hover();
        await page.waitForTimeout(1000);
        await page.mouse.up();

        await expect(setlistItems.nth(0)).toContainText('Another Song To Add');
        await expect(setlistItems.nth(1)).toContainText('Interlude To Add');
        await expect(setlistItems.nth(2)).toContainText('Song Title 1');
    });

    test('should navigate to new song page and back via cancel', async ({ page }) => {
        await getLibraryContainer(page).getByRole('link', { name: '+ Create New Song' }).click();

        await page.waitForURL(`/song/new?redirectTo=${addUrl}`);
        await expect(page.getByRole('heading', { name: 'Add a New Song to Your Library' })).toBeVisible();

        await page.getByRole('link', { name: 'Cancel' }).click();

        await page.waitForURL(addUrl);
        await expect(page.getByRole('heading', { name: /Build Setlist:/ })).toBeVisible();
    });

    test('should navigate to new interlude page and back via cancel', async ({ page }) => {
        const library = getLibraryContainer(page);

        await library.getByRole('button', { name: /Interludes/ }).click();
        await library.getByRole('link', { name: '+ Create New Interlude' }).click();

        await page.waitForURL(`/interlude/new?redirectTo=${addUrl}`);
        await expect(page.getByRole('heading', { name: 'Add a New Interlude' })).toBeVisible();

        await page.getByRole('link', { name: 'Cancel' }).click();

        await page.waitForURL(addUrl);
        await expect(page.getByRole('heading', { name: /Build Setlist:/ })).toBeVisible();
    });
});