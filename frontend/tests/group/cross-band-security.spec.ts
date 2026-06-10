import { test, expect, type Page } from '@playwright/test';

async function loginMultiGroupUser(page: Page) {
	await page.goto('/login');
	await page.getByLabel("Nom d'utilisateur").fill('multiGroupUser');
	await page.locator('#password').fill('password123');
	await page.getByRole('button', { name: 'Se connecter' }).click();
	await page.waitForURL('/');
}

async function switchToBand(page: Page, bandName: string) {
	await page.goto('/');
	const bandSelector = page.locator('select[name="bandId"]');
	await bandSelector.selectOption({ label: bandName });
	await page.waitForURL('/');
	await expect(page.getByRole('heading', { name: bandName })).toBeVisible();
}

async function getActiveBandData(page: Page) {
	const setlistResponse = await page.request.get('/api/setlist');
	expect(setlistResponse.ok()).toBeTruthy();
	const setlists = await setlistResponse.json();

	const songResponse = await page.request.get('/api/song');
	expect(songResponse.ok()).toBeTruthy();
	const songs = await songResponse.json();

	return { setlists, songs };
}

test.describe('Cross-Band Security', () => {
	test('should prevent a user from modifying setlists or using songs of another band', async ({
		page
	}) => {
		await loginMultiGroupUser(page);

		// Collect Band B resource IDs while Band B is active.
		await switchToBand(page, 'Band B');
		const bandB = await getActiveBandData(page);
		const setlistB = bandB.setlists.find((s: { name: string }) => s.name === 'Setlist B');
		const songB = bandB.songs.find((s: { title: string }) => s.title === 'Chanson B1');
		expect(setlistB).toBeDefined();
		expect(songB).toBeDefined();

		// Switch back to Band A: every request below carries Band A as active band.
		await switchToBand(page, 'Band A');
		const bandA = await getActiveBandData(page);
		const setlistA = bandA.setlists.find((s: { name: string }) => s.name === 'Setlist A');
		const songA = bandA.songs.find((s: { title: string }) => s.title === 'Chanson A1');
		expect(setlistA).toBeDefined();
		expect(songA).toBeDefined();

		// Cannot add an item to another band's setlist.
		const addToForeignSetlist = await page.request.post(`/api/setlist/${setlistB.id}/items`, {
			data: { item_type: 'song', item_id: songA.id }
		});
		expect(addToForeignSetlist.status()).toBe(404);

		// Cannot add another band's song to an owned setlist.
		const addForeignSong = await page.request.post(`/api/setlist/${setlistA.id}/items`, {
			data: { item_type: 'song', item_id: songB.id }
		});
		expect(addForeignSong.status()).toBe(404);

		// Cannot reorder the items of another band's setlist.
		const reorderForeignSetlist = await page.request.put(
			`/api/setlist/${setlistB.id}/items/order`,
			{
				data: { item_ids: [1] }
			}
		);
		expect(reorderForeignSetlist.status()).toBe(404);

		// Sanity check: adding an owned song to an owned setlist still works.
		const addOwnSong = await page.request.post(`/api/setlist/${setlistA.id}/items`, {
			data: { item_type: 'song', item_id: songA.id }
		});
		expect(addOwnSong.status()).toBe(201);
	});
});
