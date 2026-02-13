import { test, expect, type Page } from '@playwright/test';

async function login(page: Page, user: string, pass: string) {
    await page.goto('/login');
    await page.getByLabel('Username').fill(user);
    await page.getByLabel('Password').fill(pass);
    await page.getByRole('button', { name: 'Log In' }).click();
    await page.waitForURL('/');
}

async function createSetlist(page: Page, name: string): Promise<{ id: number; name: string }> {
    const setlistRes = await page.request.post('/api/setlist', {
        data: { name, color: '#00ffff' }
    });
    expect(setlistRes.ok()).toBeTruthy();
    const setlist = await setlistRes.json();
    return { id: setlist.id, name: setlist.name };
}

test.describe.serial('Setlist Admin Actions [As Admin]', () => {
    let setlistId: number;
    let setlistName: string;

    test.beforeEach(async ({ page }) => {
        await login(page, 'testuser', 'Password123!');
        const uniqueName = `Setlist for Actions Test ${Date.now()}`;
        const createdSetlist = await createSetlist(page, uniqueName);
        setlistId = createdSetlist.id;
        setlistName = createdSetlist.name;
    });

    test('should allow an admin to archive and unarchive a setlist', async ({ page }) => {
        await page.goto(`/setlist/${setlistId}`);
        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Archiver' }).click();

        await expect(page.getByText('Archivée')).toBeVisible();

        await page.goto('/');
        await expect(page.getByRole('link', { name: setlistName })).toBeHidden();
        await page.getByRole('button', { name: 'Archivées' }).click();
        await expect(page.getByRole('link', { name: setlistName })).toBeVisible();

        await page.goto(`/setlist/${setlistId}`);
        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Désarchiver' }).click();
        await expect(page.getByText('Archivée')).toBeHidden();

        await page.goto('/');
        await expect(page.getByRole('link', { name: setlistName })).toBeVisible();
        await page.getByRole('button', { name: 'Archivées' }).click();
        await expect(page.getByRole('link', { name: setlistName })).toBeHidden();
    });

    test('should allow an admin to delete a setlist after confirmation', async ({ page }) => {
        await page.goto(`/setlist/${setlistId}`);

        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Supprimer' }).click();

        const modal = page.getByRole('dialog');
        await expect(modal).toBeVisible();
        await expect(modal).toContainText(`Êtes-vous sûr de vouloir supprimer la setlist "${setlistName}" ?`);

        await modal.getByRole('button', { name: 'Annuler' }).click();
        await expect(modal).toBeHidden();
        await expect(page.url()).toContain(`/setlist/${setlistId}`);

        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await page.getByRole('menuitem', { name: 'Supprimer' }).click();
        await page.getByRole('button', { name: 'Confirmer la suppression' }).click();

        await page.waitForURL('/');
        await expect(page.getByRole('link', { name: setlistName })).toBeHidden();
        await page.getByRole('button', { name: 'Archivées' }).click();
        await expect(page.getByRole('link', { name: setlistName })).toBeHidden();
    });
});

test.describe('Setlist Admin Actions [As Member]', () => {
    let setlistId: number;

    test.beforeAll(async ({ browser }) => {
        const page = await browser.newPage();
        await login(page, 'testuser', 'Password123!');
        const createdSetlist = await createSetlist(page, `Setlist for Member View ${Date.now()}`);
        setlistId = createdSetlist.id;
        await page.close();
    });

    test('should not show admin actions for a non-admin user', async ({ page }) => {
        await login(page, 'memberuser', 'Password123!');
        await page.goto(`/setlist/${setlistId}`);

        await page.getByRole('button', { name: "Ouvrir le menu d'actions" }).click();
        await expect(page.getByRole('menuitem', { name: 'Archiver' })).toBeHidden();
        await expect(page.getByRole('menuitem', { name: 'Désarchiver' })).toBeHidden();
        await expect(page.getByRole('menuitem', { name: 'Supprimer' })).toBeHidden();
    });
});