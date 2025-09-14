import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
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
    },
    deleteItem: async ({ request, fetch }) => {
        const data = await request.formData();
        const itemId = data.get('itemId');

        const response = await fetch(`/api/setlist/item/${itemId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            return fail(response.status, { error: 'Failed to delete item.' });
        }

        return { deleted: true, itemId: Number(itemId) };
    },
    updateSongNotes: async ({ request, fetch }) => {
        const data = await request.formData();
        const itemId = data.get('itemId');
        const notes = data.get('notes');

        const response = await fetch(`/api/setlist/item/${itemId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ notes })
        });

        if (!response.ok) return fail(response.status, { error: 'Failed to update song notes.' });
        const updatedItem = await response.json();
        return { updatedSong: true, item: updatedItem };
    },
    updateInterlude: async ({ request, fetch }) => {
        const data = await request.formData();
        const interludeId = data.get('interludeId');

        const payload = {
            title: data.get('title'),
            speaker: data.get('speaker') || null,
            script: data.get('script') || null,
            duration_seconds: Number(data.get('duration')) || null
        };

        const response = await fetch(`/api/interlude/${interludeId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) return fail(response.status, { error: 'Failed to update interlude.' });
        const updatedInterlude = await response.json();
        return { updatedInterlude: true, interlude: updatedInterlude };
    }
};