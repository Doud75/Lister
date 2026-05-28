import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    deleteSong: async ({ request, fetch }) => {
        const data = await request.formData();
        const songId = data.get('songId');

        if (!songId) {
            return fail(400, { error: 'Song ID is missing.' });
        }

        const response = await fetch(`/api/song/${songId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            return fail(response.status, { error: 'Failed to delete the song.' });
        }

        return { deleted: true, songId: Number(songId) };
    }
};
