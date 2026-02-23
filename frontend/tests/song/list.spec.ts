import { test, expect, type Page } from '@playwright/test';

async function login(page: Page) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill('testuser');
    await page.getByLabel('Mot de passe').fill('password123');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

const SONG_LIST_URL = '/song';

test.describe('Song Library Page', () => {
    test.beforeEach(async ({ page }) => {
        await login(page);
        await page.goto(SONG_LIST_URL);
        await expect(page.getByRole('heading', { name: 'Bibliothèque de chansons' })).toBeVisible();
    });

    test('should display songs grouped by album', async ({ page }) => {
        // Vérifie la présence des titres d'albums
        await expect(page.getByRole('heading', { name: 'Test Album' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Sans Album' })).toBeVisible();

        // Vérifie la présence de chansons spécifiques dans les bons groupes
        const testAlbumSection = page.locator('div.rounded-xl', { hasText: 'Test Album' });
        const noAlbumSection = page.locator('div.rounded-xl', { hasText: 'Sans Album' });

        await expect(testAlbumSection.getByText('Song Title 1')).toBeVisible();
        await expect(noAlbumSection.getByText('Song Title 2')).toBeVisible();
    });

    test('should delete a song from the library', async ({ page }) => {
        // On s'assure que la chanson à supprimer est bien présente
        const songToDelete = page.locator('li', { hasText: 'Song To Delete' });
        await expect(songToDelete).toBeVisible();

        // On clique sur le bouton de suppression
        await songToDelete.getByRole('button', { name: 'Supprimer Song To Delete' }).click();

        // On vérifie que la chanson a bien disparu de la liste grâce à l'UI optimiste
        await expect(songToDelete).toBeHidden();

        // Recharger la page pour confirmer que la suppression est persistante
        await page.reload();
        await expect(page.locator('li', { hasText: 'Song To Delete' })).toBeHidden();
    });

    test('should navigate to the edit page for a song', async ({ page }) => {
        const songToEdit = page.locator('li', { hasText: 'Song Title 1' });

        // Clique sur le lien d'édition
        await songToEdit.getByRole('link', { name: 'Modifier Song Title 1' }).click();

        // Vérifie qu'on est sur la bonne page
        await page.waitForURL('/song/1/edit');
        await expect(page.getByRole('heading', { name: 'Modifier: Song Title 1' })).toBeVisible();
    });

    test('should navigate to the new song page', async ({ page }) => {
        await page.getByRole('link', { name: '+ Ajouter une chanson' }).click();

        await page.waitForURL('/song/new');
        await expect(page.getByRole('heading', { name: 'Add a New Song to Your Library' })).toBeVisible();
    });
});