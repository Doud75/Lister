import { test, expect, type Page, type Browser } from '@playwright/test';

// ── Helpers ───────────────────────────────────────────────────────────────

async function signup(page: Page, username: string): Promise<void> {
    await page.goto('/signup');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
    await page.getByRole('button', { name: 'Créer mon compte' }).click();
    await page.waitForURL('/dashboard');
}

async function signupAndCreateBand(page: Page, username: string, bandName: string): Promise<void> {
    await signup(page, username);
    await page.getByRole('button', { name: 'Créer un groupe' }).click();
    await page.locator('#band-name').fill(bandName);
    await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();
    await page.waitForURL('/');
}

async function signupInNewContext(browser: Browser, username: string): Promise<Page> {
    const ctx = await browser.newContext();
    const p = await ctx.newPage();
    await signup(p, username);
    return p;
}

async function generateInviteLink(page: Page): Promise<string> {
    await page.goto('/settings/members');
    await page.getByRole('button', { name: 'Générer un lien' }).click();
    await expect(page.locator('input[readonly]').first()).toBeVisible({ timeout: 5000 });
    return page.locator('input[readonly]').first().inputValue();
}

// ── Tests ────────────────────────────────────────────────────────────────

test.describe('Leave Band', () => {

    // ── Cas 1 : Membre (non-admin) quitte avec succès ────────────────────
    test('member can leave a band and it disappears from dashboard', async ({ page, browser }) => {
        const ts = Date.now();
        const bandName = `Leave Band ${ts}`;

        // Admin crée le groupe et génère un lien
        await signupAndCreateBand(page, `admin_leave_${ts}`, bandName);
        const joinUrl = await generateInviteLink(page);

        // Membre rejoint via le lien
        const memberPage = await signupInNewContext(browser, `member_leave_${ts}`);
        await memberPage.goto(joinUrl);
        await memberPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await memberPage.waitForURL('/');

        // Membre va sur le dashboard
        await memberPage.goto('/dashboard');
        await expect(memberPage.getByRole('button', { name: bandName })).toBeVisible();

        // Membre clique sur la poubelle → modale de confirmation s'ouvre
        await memberPage.getByRole('button', { name: `Quitter ${bandName}` }).click();
        await expect(memberPage.getByText('Êtes-vous sûr de vouloir quitter')).toBeVisible();
        await expect(memberPage.getByRole('strong', { name: bandName })).toBeVisible();

        // Confirmation
        await memberPage.getByRole('button', { name: 'Quitter le groupe' }).click();

        // Redirigé vers /dashboard avec flash message
        await memberPage.waitForURL(/\/dashboard/);
        await expect(memberPage.getByText(`Vous avez quitté le groupe`)).toBeVisible();
        await expect(memberPage.getByText(bandName)).toBeVisible();

        // Le groupe n'apparaît plus dans la liste
        await expect(memberPage.getByRole('button', { name: bandName })).toBeHidden();

        await memberPage.context().close();
    });

    // ── Cas 2 : Seul admin → bloqué avec message 409 ────────────────────
    test('last admin cannot leave and sees an error message', async ({ page }) => {
        const ts = Date.now();
        const bandName = `Solo Admin Band ${ts}`;

        await signupAndCreateBand(page, `solo_admin_${ts}`, bandName);
        await page.goto('/dashboard');

        // Ouvre la modale de quitter
        await page.getByRole('button', { name: `Quitter ${bandName}` }).click();
        await expect(page.getByText('Êtes-vous sûr de vouloir quitter')).toBeVisible();

        await page.getByRole('button', { name: 'Quitter le groupe' }).click();

        // Reste sur /dashboard, message d'erreur visible
        await expect(page).toHaveURL(/\/dashboard/);
        await expect(page.getByText('dernier administrateur')).toBeVisible();

        // Le groupe est toujours présent
        await expect(page.getByRole('button', { name: bandName })).toBeVisible();
    });

    // ── Cas 3 : Admin avec co-admin peut quitter ─────────────────────────
    test('admin with another admin can leave the band', async ({ page, browser }) => {
        const ts = Date.now();
        const bandName = `Co Admin Band ${ts}`;

        // Premier admin crée le groupe
        await signupAndCreateBand(page, `admin_co1_${ts}`, bandName);
        const joinUrl = await generateInviteLink(page);

        // Co-admin rejoint et est promu admin par l'API directement
        const coAdminPage = await signupInNewContext(browser, `admin_co2_${ts}`);
        await coAdminPage.goto(joinUrl);
        await coAdminPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await coAdminPage.waitForURL('/');

        // Promouvoir le co-admin via settings/members (l'admin 1 le fait)
        await page.goto('/settings/members');
        const memberRow = page.locator('li', { hasText: `admin_co2_${ts}` });
        await expect(memberRow).toBeVisible();
        // Cliquer sur "Promouvoir admin"
        await memberRow.getByRole('button', { name: 'Promouvoir admin' }).click();

        // Admin 1 quitte maintenant
        await page.goto('/dashboard');
        await page.getByRole('button', { name: `Quitter ${bandName}` }).click();
        await expect(page.getByText('Êtes-vous sûr de vouloir quitter')).toBeVisible();
        await page.getByRole('button', { name: 'Quitter le groupe' }).click();

        await page.waitForURL(/\/dashboard/);
        await expect(page.getByText('Vous avez quitté le groupe')).toBeVisible();
        await expect(page.getByRole('button', { name: bandName })).toBeHidden();

        // Le groupe existe encore (co-admin toujours dedans)
        await expect(coAdminPage.getByRole('heading', { name: bandName })).toBeVisible();

        await coAdminPage.context().close();
    });
});
