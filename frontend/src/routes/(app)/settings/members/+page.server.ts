import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
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
