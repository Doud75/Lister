import { test, expect, type Page } from '@playwright/test';

async function signup(page: Page, username: string, password: string) {
    await page.goto('/signup');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.getByLabel('Mot de passe', { exact: true }).fill(password);
    await page.getByRole('button', { name: 'Créer mon compte' }).click();
    await page.waitForURL('/dashboard');
}

async function createBand(page: Page, bandName: string) {
    await page.goto('/dashboard');
    await page.getByRole('button', { name: 'Créer un groupe' }).click();
    await page.locator('#band-name').fill(bandName);
    await page.getByRole('button', { name: 'Créer le groupe', exact: true }).click();
    await page.waitForURL('/');
    await page.goto('/dashboard');
}

async function logout(page: Page) {
    await page.goto('/');
    await page.getByRole('button', { name: 'Ouvrir le menu du profil' }).click();
    await page.getByRole('menuitem', { name: 'Déconnexion' }).click();
    await page.waitForURL('/login');
}

async function login(page: Page, username: string, password: string) {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.locator('#password').fill(password);
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

async function generateInviteLink(page: Page): Promise<string> {
    await page.goto('/settings/members');
    await page.getByRole('button', { name: 'Générer un lien' }).click();
    await expect(page.locator('input[readonly]').first()).toBeVisible({ timeout: 5000 });
    return page.locator('input[readonly]').first().inputValue();
}

test.describe('Default Band', () => {
    // ── Cas 1 : cliquer l'étoile change le groupe par défaut ──────────────
    test('clicking star on a non-default band makes it the default', async ({ page }) => {
        const timestamp = Date.now();
        const username = `star_user_${timestamp}`;
        const bandA = `Star Band A ${timestamp}`;
        const bandB = `Star Band B ${timestamp}`;

        await signup(page, username, 'StrongPass1!');
        await createBand(page, bandA);
        await createBand(page, bandB);

        await page.goto('/dashboard');

        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true }).first()).toBeVisible();

        await page.getByRole('button', { name: `Définir ${bandB} comme groupe par défaut`, exact: true }).click();

        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
        await expect(page.getByRole('button', { name: `Définir ${bandA} comme groupe par défaut`, exact: true })).toBeVisible();
    });

    // ── Cas 2 : après logout/login → redirigé sur le groupe par défaut ────
    test('after login, active band is the default band', async ({ page }) => {
        const timestamp = Date.now();
        const username = `default_login_user_${timestamp}`;
        const bandA = `Default Login Band A ${timestamp}`;
        const bandB = `Default Login Band B ${timestamp}`;

        await signup(page, username, 'StrongPass1!');
        await createBand(page, bandA);
        await createBand(page, bandB);

        await page.goto('/dashboard');
        await page.getByRole('button', { name: `Définir ${bandB} comme groupe par défaut`, exact: true }).click();
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();

        await logout(page);
        await login(page, username, 'StrongPass1!');

        await page.goto('/dashboard');
        await expect(page.getByText('✓ Groupe actif').first()).toBeVisible();
        await expect(page.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
    });

    // ── Cas 3 : quitter le groupe par défaut → un autre devient défaut ────
    test('leaving default band auto-reassigns default to another band', async ({ page, browser }) => {
        const timestamp = Date.now();
        const adminUsername = `admin_reassign_${timestamp}`;
        const memberUsername = `member_reassign_${timestamp}`;
        const bandA = `Reassign Band A ${timestamp}`;
        const bandB = `Reassign Band B ${timestamp}`;

        await signup(page, adminUsername, 'StrongPass1!');
        await createBand(page, bandA);
        const joinUrl = await generateInviteLink(page);

        const memberCtx = await browser.newContext();
        const memberPage = await memberCtx.newPage();

        await signup(memberPage, memberUsername, 'StrongPass1!');
        await createBand(memberPage, bandB);

        await memberPage.goto(joinUrl);
        await memberPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await memberPage.waitForURL('/');

        await memberPage.goto('/dashboard');
        await memberPage.getByRole('button', { name: `Définir ${bandA} comme groupe par défaut`, exact: true }).click();
        await expect(memberPage.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();

        await memberPage.getByRole('button', { name: `Quitter ${bandA}` }).click();
        await memberPage.getByRole('button', { name: 'Quitter le groupe' }).click();
        await memberPage.waitForURL(/\/dashboard/);

        await expect(memberPage.getByRole('button', { name: 'Groupe par défaut', exact: true })).toBeVisible();
        await expect(memberPage.getByRole('heading', { name: bandA, level: 2 })).not.toBeVisible();

        await memberCtx.close();
    });
});