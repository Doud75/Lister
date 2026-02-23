import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.locator('#password').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

test.describe('Interlude Behavior in Setlists', () => {
    let setlistA_Id: number;
    let setlistB_Id: number;
    let interludeId: number;
    const interludeTitle = `Test Interlude ${Date.now()}`;
    const originalScript = 'This is the original default script.';

    test.beforeEach(async ({ page }) => {
        await login(page);

        const interludeRes = await page.request.post('/api/interlude', {
            data: {
                title: interludeTitle,
                speaker: 'Original Speaker',
                duration_seconds: 60,
                script: originalScript
            }
        });
        expect(interludeRes.ok()).toBeTruthy();
        interludeId = (await interludeRes.json()).id;

        const setlistARes = await page.request.post('/api/setlist', {
            data: { name: `Setlist A ${Date.now()}`, color: '#ff0000' }
        });
        expect(setlistARes.ok()).toBeTruthy();
        setlistA_Id = (await setlistARes.json()).id;

        const setlistBRes = await page.request.post('/api/setlist', {
            data: { name: `Setlist B ${Date.now()}`, color: '#0000ff' }
        });
        expect(setlistBRes.ok()).toBeTruthy();
        setlistB_Id = (await setlistBRes.json()).id;

        const addItemRes = await page.request.post(`/api/setlist/${setlistA_Id}/items`, {
            data: { item_type: 'interlude', item_id: interludeId }
        });
        expect(addItemRes.ok()).toBeTruthy();
    });

    test('should modify script for one setlist only', async ({ page }) => {
        await page.goto(`/setlist/${setlistA_Id}`);

        const interludeItemA = page.locator('li').filter({ hasText: interludeTitle });

        await expect(interludeItemA).toBeVisible();
        await expect(interludeItemA).toContainText(originalScript);

        await interludeItemA.getByRole('button', { name: 'Edit item' }).click({ force: true});

        const modal = page.getByRole('dialog');
        await expect(modal).toBeVisible();
        const scriptForSetlistA = 'This is the script ONLY for Setlist A.';
        await modal.getByLabel('Script').fill(scriptForSetlistA);
        await modal.getByRole('button', { name: 'Save Interlude' }).click();
        await expect(modal).toBeHidden();

        await expect(interludeItemA).toContainText(scriptForSetlistA);
        await expect(interludeItemA).not.toContainText(originalScript);

        await page.request.post(`/api/setlist/${setlistB_Id}/items`, {
            data: { item_type: 'interlude', item_id: interludeId }
        });

        await page.goto(`/setlist/${setlistB_Id}`);
        const interludeItemB = page.locator('li').filter({ hasText: interludeTitle });
        await expect(interludeItemB).toBeVisible();
        await expect(interludeItemB).toContainText(originalScript);
        await expect(interludeItemB).not.toContainText(scriptForSetlistA);
    });

    test('should modify global metadata and reflect in both setlists', async ({ page }) => {
        await page.request.post(`/api/setlist/${setlistB_Id}/items`, {
            data: { item_type: 'interlude', item_id: interludeId }
        });

        await page.goto(`/setlist/${setlistA_Id}`);
        const interludeItemA = page.locator('li').filter({ hasText: interludeTitle });
        await expect(interludeItemA).toBeVisible();

        await interludeItemA.getByRole('button', { name: 'Edit item' }).click({ force: true});

        const modal = page.getByRole('dialog');
        await expect(modal).toBeVisible();
        const newSpeaker = 'Global New Speaker';
        await modal.getByLabel('Speaker').fill(newSpeaker);
        await modal.getByRole('button', { name: 'Save Interlude' }).click();
        await expect(modal).toBeHidden();

        await expect(interludeItemA).toContainText(`Intervenant : ${newSpeaker}`);

        await page.goto(`/setlist/${setlistB_Id}`);
        const interludeItemB = page.locator('li').filter({ hasText: interludeTitle });
        await expect(interludeItemB).toBeVisible();
        await expect(interludeItemB).toContainText(`Intervenant : ${newSpeaker}`);
    });

    test('should delete an item from a setlist without affecting other setlists', async ({ page }) => {
        await page.request.post(`/api/setlist/${setlistB_Id}/items`, {
            data: { item_type: 'interlude', item_id: interludeId }
        });

        await page.goto(`/setlist/${setlistA_Id}`);
        const interludeItemA = page.locator('li').filter({ hasText: interludeTitle });
        await expect(interludeItemA).toBeVisible();

        await interludeItemA.getByRole('button', { name: 'Remove item' }).click({ force: true});

        await expect(interludeItemA).toBeHidden();
        await expect(page.getByText('This setlist is empty')).toBeVisible();

        await page.goto(`/setlist/${setlistB_Id}`);
        const interludeItemB = page.locator('li').filter({ hasText: interludeTitle });
        await expect(interludeItemB).toBeVisible();
        await expect(interludeItemB).toContainText(originalScript);
    });
});