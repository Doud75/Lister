import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    deleteInterlude: async ({ request, fetch }) => {
        const data = await request.formData();
        const interludeId = data.get('interludeId');
        const response = await fetch(`/api/interlude/${interludeId}`, { method: 'DELETE' });

        if (!response.ok) {
            return fail(response.status, { error: 'Failed to delete interlude.' });
        }
        return { success: true };
    }
};