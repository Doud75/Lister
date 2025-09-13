import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

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

    const fetchInterludes = async () => {
        const res = await fetch('/api/interlude');
        if (!res.ok) throw error(res.status, 'Failed to fetch interludes');
        return res.json();
    };

    try {
        const [setlist, songs, interludes] = await Promise.all([
            fetchSetlistInfo(),
            fetchSongs(),
            fetchInterludes()
        ]);
        return { setlist, songs, interludes };
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
        return { success: true };
    },

    addInterlude: async ({ request, params, fetch }) => {
        const data = await request.formData();
        const interludeId = data.get('interludeId');

        const response = await fetch(`/api/setlist/${params.id}/items`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                item_type: 'interlude',
                item_id: Number(interludeId)
            })
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error || 'Failed to add interlude.' });
        }

        return { success: true };
    },

    updateOrder: async ({ request, params, fetch }) => {
        const data = await request.formData();
        const itemIdsStr = data.get('itemIds');

        if (!itemIdsStr) {
            return fail(400, { error: 'No item IDs provided.' });
        }

        const item_ids = JSON.parse(itemIdsStr.toString());

        const response = await fetch(`/api/setlist/${params.id}/items/order`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ item_ids })
        });

        if (!response.ok) {
            return fail(response.status, { error: 'Failed to save the new order.' });
        }

        return { success: true };
    }
};