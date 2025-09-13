import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

// La fonction 'load' s'exécute sur le serveur avant le rendu de la page.
export const load: PageServerLoad = async ({ fetch, params }) => {
    const { id } = params;

    const fetchSetlistInfo = async () => {
        const res = await fetch(`/api/setlist/${id}`);
        if (!res.ok) throw error(res.status, 'Setlist not found');
        return res.json();
    };

    const fetchSongs = async () => {
        const res = await fetch('/api/song');
        if (!res.ok) throw error(res.status, 'Failed to fetch song library');
        return res.json();
    };

    try {
        const [setlist, songs] = await Promise.all([fetchSetlistInfo(), fetchSongs()]);
        return { setlist, songs };
    } catch (err: any) {
        throw error(err.status || 500, err.body?.message || 'Could not load required data.');
    }
};

export const actions: Actions = {
    addSong: async ({ request, params, fetch }) => {
        const data = await request.formData();
        const songId = data.get('songId');

        const response = await fetch(`/api/setlist/${params.id}/items`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                item_type: 'song',
                item_id: Number(songId)
            })
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error || 'Failed to add song.' });
        }

        return { success: true, addedSongId: songId };
    }
};