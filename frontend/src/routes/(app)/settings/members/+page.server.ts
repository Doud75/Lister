import { error, fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { BandMember } from '$lib/types';

export const load: PageServerLoad = async ({ fetch, locals }) => {
    if (locals.user?.role !== 'admin') {
        throw redirect(303, '/');
    }

    const bandId = locals.activeBandId;
    if (!bandId) {
        throw error(400, 'Active band not selected');
    }
    const res = await fetch(`/api/bands/${bandId}/members`);

    if (!res.ok) {
        throw error(res.status, 'Failed to fetch members.');
    }

    const members: BandMember[] = await res.json();
    return { members };
};

export const actions: Actions = {
    inviteMember: async ({ request, fetch, locals }) => {
        const bandId = locals.activeBandId;
        const data = await request.formData();
        const username = data.get('username') as string;
        const password = data.get('password') as string; // Peut être une chaîne vide

        if (!username) {
            return fail(400, { error: "Le nom d'utilisateur est requis." });
        }

        const payload = {
            username,
            password: password || null
        };

        const response = await fetch(`/api/bands/${bandId}/members`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error || "Échec de l'ajout du membre.", username });
        }

        return { inviteSuccess: true };
    },

    removeMember: async ({ request, fetch, locals }) => {
        const bandId = locals.activeBandId;
        const data = await request.formData();
        const userId = data.get('userId');

        if (!userId) {
            return fail(400, { error: 'User ID manquant.' });
        }

        const response = await fetch(`/api/bands/${bandId}/members/${userId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            const result = await response.json().catch(() => ({ error: 'Une erreur est survenue lors de la suppression.' }));
            return fail(response.status, { error: result.error });
        }

        return { removeSuccess: true, removedUserId: Number(userId) };
    }
};