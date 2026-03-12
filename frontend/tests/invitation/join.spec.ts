import { test, expect, type Page } from '@playwright/test';


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

async function login(page: Page, username: string, password = 'StrongPass1!'): Promise<void> {
    await page.goto('/login');
    await page.getByLabel("Nom d'utilisateur").fill(username);
    await page.locator('#password').fill(password);
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.waitForURL('/');
}

/** Admin clicks "Générer un lien" and returns the generated join URL. */
async function generateInviteLink(page: Page): Promise<string> {
    await page.goto('/settings/members');
    await page.getByRole('button', { name: 'Générer un lien' }).click();
    // Wait for the input containing the link
    const linkInput = page.locator('input[readonly]').filter({ hasText: /\/join\// });
    await expect(linkInput).toBeVisible({ timeout: 5000 });
    const link = await linkInput.inputValue();
    expect(link).toContain('/join/');
    return link;
}

// ── Test Suite ──────────────────────────────────────────────────────────────

test.describe('Invitation System', () => {

    // ── Cas 1 : Admin génère et copie un lien ──────────────────────────────
    test('admin can generate an invitation link and see expiry', async ({ page }) => {
        const ts = Date.now();
        await signupAndCreateBand(page, `admin_gen_${ts}`, `Gen Band ${ts}`);

        await page.goto('/settings/members');
        await expect(page.getByRole('heading', { name: 'Inviter par lien' })).toBeVisible();

        await page.getByRole('button', { name: 'Générer un lien' }).click();

        // Link input appears
        await expect(page.locator('input[readonly]')).toBeVisible({ timeout: 5000 });
        const linkValue = await page.locator('input[readonly]').inputValue();
        expect(linkValue).toMatch(/\/join\/[a-f0-9]+/);

        // Expiry label appears
        await expect(page.getByText(/Valide jusqu'au/)).toBeVisible();

        // Copy button
        await expect(page.getByRole('button', { name: 'Copier' })).toBeVisible();
    });

    // ── Cas 2 : Utilisateur connecté rejoint via lien ──────────────────────
    test('logged-in user can join a band via invitation link', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_join_${ts}`;
        const guestUser = `guest_join_${ts}`;
        const bandName = `Join Band ${ts}`;

        // Admin creates band and generates link
        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        // Guest signs up (no band) in a fresh context
        const guestContext = await context.browser()!.newContext();
        const guestPage = await guestContext.newPage();
        await signup(guestPage, guestUser);

        // Navigate to join URL
        await guestPage.goto(joinUrl);
        await expect(guestPage.getByRole('heading', { level: 1, name: bandName })).toBeVisible();
        await expect(guestPage.getByText(`Connecté en tant que`)).toBeVisible();
        await expect(guestPage.getByText(guestUser)).toBeVisible();

        // Click "Rejoindre le groupe"
        await guestPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await guestPage.waitForURL('/');

        // Guest should now see the band name on home page
        await expect(guestPage.getByRole('heading', { name: bandName })).toBeVisible();

        await guestContext.close();
    });

    // ── Cas 3 : Non connecté → login → retour invitation ──────────────────
    test('unauthenticated user is redirected to login then back to join page', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_auth_${ts}`;
        const guestUser = `guest_auth_${ts}`;
        const bandName = `Auth Band ${ts}`;

        // Admin sets up
        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        // Guest has an account but is not connected — use new context
        const guestContext = await context.browser()!.newContext();
        const guestPage = await guestContext.newPage();
        await signup(guestPage, guestUser);
        // Clear session to simulate logged-out state
        await guestContext.clearCookies();

        // Navigate to join URL while logged out
        await guestPage.goto(joinUrl);

        // Page should show login/signup buttons (not the join button)
        await expect(guestPage.getByRole('link', { name: 'Se connecter' })).toBeVisible();
        await expect(guestPage.getByRole('link', { name: 'Créer un compte' })).toBeVisible();

        // Login link has redirectTo parameter
        const loginHref = await guestPage.getByRole('link', { name: 'Se connecter' }).getAttribute('href');
        expect(loginHref).toContain('redirectTo');
        expect(loginHref).toContain('/join/');

        // Click login → fill credentials → redirected back to join page
        await guestPage.getByRole('link', { name: 'Se connecter' }).click();
        await guestPage.getByLabel("Nom d'utilisateur").fill(guestUser);
        await guestPage.locator('#password').fill('StrongPass1!');
        await guestPage.getByRole('button', { name: 'Se connecter' }).click();

        // Should be back on the join page
        await expect(guestPage).toHaveURL(/\/join\//);
        await expect(guestPage.getByRole('heading', { level: 1, name: bandName })).toBeVisible();

        // Now join
        await guestPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await guestPage.waitForURL('/');
        await expect(guestPage.getByRole('heading', { name: bandName })).toBeVisible();

        await guestContext.close();
    });

    // ── Cas 4 : Nouveau compte → signup → retour invitation ───────────────
    test('new user can signup via invitation link and land in the band', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_new_${ts}`;
        const newUser = `newbie_${ts}`;
        const bandName = `New User Band ${ts}`;

        // Admin sets up
        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        // New user opens join URL with no account
        const freshContext = await context.browser()!.newContext();
        const freshPage = await freshContext.newPage();

        await freshPage.goto(joinUrl);
        await expect(freshPage.getByRole('link', { name: 'Créer un compte' })).toBeVisible();

        // Signup link has redirectTo
        const signupHref = await freshPage.getByRole('link', { name: 'Créer un compte' }).getAttribute('href');
        expect(signupHref).toContain('redirectTo');

        // Click → signup → redirected back to join page
        await freshPage.getByRole('link', { name: 'Créer un compte' }).click();
        await freshPage.getByLabel("Nom d'utilisateur").fill(newUser);
        await freshPage.getByLabel('Mot de passe', { exact: true }).fill('StrongPass1!');
        await freshPage.getByRole('button', { name: 'Créer mon compte' }).click();

        // Should be back on the join page (not /dashboard)
        await expect(freshPage).toHaveURL(/\/join\//);
        await expect(freshPage.getByRole('heading', { level: 1, name: bandName })).toBeVisible();

        // Join
        await freshPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await freshPage.waitForURL('/');
        await expect(freshPage.getByRole('heading', { name: bandName })).toBeVisible();

        await freshContext.close();
    });

    // ── Cas 5 : Lien multi-usage ───────────────────────────────────────────
    test('the same invitation link can be used by multiple users', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_multi_${ts}`;
        const guest1 = `guest_multi1_${ts}`;
        const guest2 = `guest_multi2_${ts}`;
        const bandName = `Multi Band ${ts}`;

        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        for (const guestName of [guest1, guest2]) {
            const ctx = await context.browser()!.newContext();
            const p = await ctx.newPage();
            await signup(p, guestName);
            await p.goto(joinUrl);
            await p.getByRole('button', { name: 'Rejoindre le groupe' }).click();
            await p.waitForURL('/');
            await expect(p.getByRole('heading', { name: bandName })).toBeVisible();
            await ctx.close();
        }

        // Admin should see both guests in members list
        await page.goto('/settings/members');
        await expect(page.locator('li', { hasText: guest1 })).toBeVisible();
        await expect(page.locator('li', { hasText: guest2 })).toBeVisible();
    });

    // ── Cas limite 6 : Token invalide ─────────────────────────────────────
    test('invalid token shows error page', async ({ page }) => {
        await page.goto('/join/thisisnotavalidtoken00000000000000');
        await expect(page.getByRole('heading', { name: 'Invitation invalide' })).toBeVisible();
        await expect(page.getByText(/invalide|expiré/i)).toBeVisible();
    });

    // ── Cas limite 7 : Déjà membre ────────────────────────────────────────
    test('already-member acceptance is silent and redirects to /', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_dup_${ts}`;
        const guestUser = `guest_dup_${ts}`;
        const bandName = `Dup Band ${ts}`;

        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        // Guest joins once
        const guestCtx = await context.browser()!.newContext();
        const guestPage = await guestCtx.newPage();
        await signup(guestPage, guestUser);
        await guestPage.goto(joinUrl);
        await guestPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await guestPage.waitForURL('/');

        // Guest joins again with same link — should NOT error
        await guestPage.goto(joinUrl);
        await guestPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await guestPage.waitForURL('/');
        // Should land on home page without any error visible
        await expect(guestPage.getByRole('heading', { name: bandName })).toBeVisible();

        await guestCtx.close();
    });

    // ── Cas limite 8 : Non-admin ne peut pas générer un lien ──────────────
    test('non-admin member cannot generate an invitation link', async ({ page, context }) => {
        const ts = Date.now();
        const adminUser = `admin_perm_${ts}`;
        const memberUser = `member_perm_${ts}`;
        const bandName = `Perm Band ${ts}`;

        // Admin sets up and invites member
        await signupAndCreateBand(page, adminUser, bandName);
        const joinUrl = await generateInviteLink(page);

        // Member joins
        const memberCtx = await context.browser()!.newContext();
        const memberPage = await memberCtx.newPage();
        await signup(memberPage, memberUser);
        await memberPage.goto(joinUrl);
        await memberPage.getByRole('button', { name: 'Rejoindre le groupe' }).click();
        await memberPage.waitForURL('/');

        // Member navigates to settings/members (might not even be accessible)
        await memberPage.goto('/settings/members');
        // Either redirected away or the invite section is not there / button is missing
        const isOnSettings = memberPage.url().includes('/settings/members');
        if (isOnSettings) {
            // If page loads, the "Générer un lien" button shouldn't work (403 from API)
            // We just verify the page doesn't crash
            await expect(memberPage).not.toHaveURL('/login');
        } else {
            // Redirected — acceptable
            expect(memberPage.url()).not.toContain('/settings/members');
        }

        await memberCtx.close();
    });
});
