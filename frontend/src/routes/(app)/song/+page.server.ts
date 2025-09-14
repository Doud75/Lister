import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
    const fetchSongs = async () => {
        const res = await fetch('/api/song');
        if (!res.ok) {
            throw error(res.status, 'Failed to fetch song library.');
        }
        return res.json();
    };

    try {
        return {
            songs: await fetchSongs()
        };
    } catch (err: any) {
        throw error(err.status || 500, err.body?.message || 'Could not load song library.');
    }
};

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