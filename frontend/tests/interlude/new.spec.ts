import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.getByLabel('Mot de passe').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

async function createSetlistForRedirect(page: Page): Promise<string> {
    const response = await page.request.post('/api/setlist', {
        data: {
            name: `Test Setlist for Redirect ${Date.now()}`,
            color: '#00ff00'
        }
    });
    expect(response.ok()).toBeTruthy();
    const setlist = await response.json();
    return `/setlist/${setlist.id}/add`;
}

test.describe('Create New Interlude Page', () => {
    let addUrl: string;

    test.beforeEach(async ({ page }) => {
        await login(page);
        addUrl = await createSetlistForRedirect(page);
    });

    test('should successfully create a new interlude', async ({page}) => {
        await page.goto(`/interlude/new?redirectTo=${addUrl}`);

        const newInterludeTitle = `My New Test Interlude ${Date.now()}`;
        await page.getByLabel('Interlude Title').fill(newInterludeTitle);
        await page.getByLabel('Speaker (optional)').fill('Test Speaker');
        await page.getByLabel('Duration (seconds)').fill('45');
        await page.getByLabel('Script / Notes (optional)').fill('This is the default script.');
        await page.getByRole('button', { name: 'Create Interlude' }).click();

        await page.waitForURL(addUrl);

        await page.getByRole('button', { name: /Interludes/ }).click();
        await expect(page.locator('li').filter({ hasText: newInterludeTitle })).toBeVisible();
    });

    test('should show validation error if title is empty', async ({ page }) => {
        await page.goto('/interlude/new');
        await page.getByRole('button', { name: 'Create Interlude' }).click();

        await page.waitForTimeout(500);
        await expect(page).toHaveURL('/interlude/new');

        const titleInput = page.getByLabel('Interlude Title');
        const isInvalid = await titleInput.evaluate(
            (element) => !(element as HTMLInputElement).checkValidity()
        );
        expect(isInvalid).toBe(true);
    });

    test('should navigate back when cancel is clicked', async ({ page }) => {
        await page.goto(`/interlude/new?redirectTo=${addUrl}`);
        await page.getByRole('link', { name: 'Cancel' }).click();

        await page.waitForURL(addUrl);
        await expect(page.getByRole('heading', { name: /Build Setlist:/ })).toBeVisible();
    });
});