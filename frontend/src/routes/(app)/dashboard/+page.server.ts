import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, locals }) => {
    if (!locals.user) {
        throw redirect(303, '/login');
    }

    const res = await fetch('/api/user/bands');
    if (!res.ok) {
        return { bands: [] };
    }
    const bands = await res.json();

    return {
        bands,
        activeBandId: locals.activeBandId
    };
};

export const actions: Actions = {
    createBand: async ({ request, fetch, cookies }) => {
        const formData = await request.formData();
        const name = formData.get('name')?.toString()?.trim() ?? '';

        if (!name) {
            return fail(400, { error: 'Le nom du groupe est requis.' });
        }

        const res = await fetch('/api/bands', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        });

        const result = await res.json();

        if (!res.ok) {
            return fail(res.status, { error: result.error ?? 'Erreur lors de la création du groupe.' });
        }

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 30,
            sameSite: 'lax' as const
        };

        cookies.set('active_band_id', result.id.toString(), cookieOptions);

        const existingBands = JSON.parse(cookies.get('user_bands') ?? '[]');
        existingBands.push({ id: result.id, name: result.name });
        cookies.set('user_bands', JSON.stringify(existingBands), cookieOptions);

        throw redirect(303, '/');
    },

    switchBand: async ({ request, cookies }) => {
        const formData = await request.formData();
        const bandId = formData.get('bandId')?.toString();

        if (bandId) {
            cookies.set('active_band_id', bandId, {
                path: '/',
                httpOnly: true,
                secure: process.env.NODE_ENV === 'production',
                maxAge: 60 * 60 * 24 * 7,
                sameSite: 'lax'
            });
        }

        throw redirect(303, '/');
    },

    setDefault: async ({ request, fetch }) => {
        const formData = await request.formData();
        const bandId = formData.get('bandId')?.toString();

        if (!bandId) return fail(400, { defaultError: 'Identifiant manquant.' });

        const res = await fetch('/api/user/default-band', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ band_id: parseInt(bandId) })
        });

        if (!res.ok) {
            const result = await res.json().catch(() => ({}));
            return fail(res.status, { defaultError: result.error ?? 'Erreur.' });
        }

        return { setDefaultSuccess: true, bandId: parseInt(bandId) };
    },

    leaveBand: async ({ request, fetch, cookies }) => {
        const formData = await request.formData();
        const bandId = formData.get('bandId')?.toString();
        const bandName = formData.get('bandName')?.toString() ?? 'ce groupe';

        if (!bandId) {
            return fail(400, { leaveError: 'Identifiant du groupe manquant.' });
        }

        const res = await fetch(`/api/bands/${bandId}/members/me`, { method: 'DELETE' });

        if (!res.ok) {
            const result = await res.json().catch(() => ({}));
            if (res.status === 409) {
                return fail(409, {
                    leaveError: result.error ?? 'Impossible de quitter : vous êtes le dernier administrateur.'
                });
            }
            return fail(res.status, { leaveError: result.error ?? "Erreur lors du départ du groupe." });
        }

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            sameSite: 'lax' as const
        };

        const userBandsCookie = cookies.get('user_bands');
        if (userBandsCookie) {
            const userBands: { id: number; name: string }[] = JSON.parse(userBandsCookie);
            const updated = userBands.filter((b) => b.id.toString() !== bandId);
            cookies.set('user_bands', JSON.stringify(updated), { ...cookieOptions, maxAge: 60 * 60 * 24 * 30 });
        }

        if (cookies.get('active_band_id') === bandId) {
            cookies.delete('active_band_id', { path: '/' });
        }

        throw redirect(303, `/dashboard?left=${encodeURIComponent(bandName)}`);
    }
};
