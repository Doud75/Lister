import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel('Username').fill('testuser');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

const SETLIST_ID = 1;
const SETLIST_URL = `/setlist/${SETLIST_ID}`;
const getSetlistItems = (page: Page) => page.locator('ul[data-testid="setlist-items"] > li');


test.describe('Setlist Detail Page', () => {

    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(SETLIST_URL);
        await page.waitForTimeout(1000);
    });

    test('should display the setlist items in the correct initial order', async ({ page }) => {
        const items = getSetlistItems(page);

        await expect(items).toHaveCount(3);
        await expect(items.nth(0)).toContainText('Song Title 1');
        await expect(items.nth(1)).toContainText('Test Interlude');
        await expect(items.nth(2)).toContainText('Song Title 2');
    });

    test('should reorder items using drag and drop', async ({ page }) => {
        const handle = page.locator('li').filter({ hasText: 'Song Title 1' }).locator('[aria-label="Drag to reorder"]');
        const dropTarget = page.locator('li').filter({ hasText: 'Song Title 2' });

        await handle.hover();
        await page.mouse.down();
        await dropTarget.hover();
        await dropTarget.hover();
        await page.waitForTimeout(1000);

        await page.mouse.up();

        const items = getSetlistItems(page);
        await expect(items.first()).toContainText('Test Interlude');

        const reorderedItems = page.locator('ul[data-testid="setlist-items"] > li');
        await expect(reorderedItems.nth(0)).toContainText('Test Interlude');
        await expect(reorderedItems.nth(1)).toContainText('Song Title 2');
        await expect(reorderedItems.nth(2)).toContainText('Song Title 1');
    });

    test('should delete an item from the setlist', async ({ page }) => {
        const interludeItem = page.locator('li').filter({ hasText: 'Test Interlude' });

        await expect(interludeItem).toBeVisible();
        await interludeItem.getByRole('button', { name: 'Remove item' }).click({ force: true });

        await expect(interludeItem).toBeHidden();
        await expect(page.locator('ul[data-testid="setlist-items"] > li')).toHaveCount(2);
    });

    test("should edit a song's notes via the modal", async ({ page }) => {
        const songItem = page.locator('li').filter({ hasText: 'Song Title 1' });
        const newNote = 'Start with the guitar riff.';

        await songItem.getByRole('button', { name: 'Edit item' }).click({ force: true });
        await expect(page.getByRole('dialog')).toBeVisible();

        const noteTextarea = page.getByRole('dialog').getByPlaceholder('Add a comment...');
        await noteTextarea.fill(newNote);
        await page.getByRole('button', { name: 'Save Note' }).click();

        await expect(page.getByRole('dialog')).toBeHidden();
        await expect(page.getByText(newNote)).toBeVisible();
    });
});