import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
    return { user: locals.user ?? null };
};

export const actions: Actions = {
    default: async ({ params, fetch, locals, cookies }) => {
        const { token } = params;

        if (!locals.token) {
            return fail(401, { error: 'Vous devez être connecté pour rejoindre un groupe.' });
        }

        const res = await fetch(`/api/invitations/${token}/accept`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${locals.token}`,
                'Content-Type': 'application/json'
            }
        });

        if (res.status === 409) {
            throw redirect(303, '/');
        }

        if (res.status === 410) {
            return fail(410, { error: "Cette invitation a expiré." });
        }

        if (!res.ok) {
            const body = await res.json().catch(() => ({ error: "Une erreur s'est produite." }));
            return fail(res.status, { error: body.error ?? "Une erreur s'est produite." });
        }

        const data = await res.json();
        const bandId: number = data.band_id;
        const bandName: string = data.band_name;

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 30,
            sameSite: 'lax' as const
        };

        cookies.set('active_band_id', bandId.toString(), cookieOptions);

        const userBandsCookie = cookies.get('user_bands');
        const userBands: { id: number; name: string }[] = userBandsCookie ? JSON.parse(userBandsCookie) : [];
        const alreadyInList = userBands.some((b) => b.id === bandId);
        if (!alreadyInList) {
            userBands.push({ id: bandId, name: bandName });
            cookies.set('user_bands', JSON.stringify(userBands), cookieOptions);
        }

        throw redirect(303, '/');
    }
};
