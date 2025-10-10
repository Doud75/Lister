import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import type {Interlude, SetlistItem} from "$lib/types";

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
        const updatedItem: SetlistItem = await response.json();
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
        const updatedInterlude: Interlude = await response.json();
        return { updatedInterlude: true, interlude: updatedInterlude };
    },
    duplicateSetlist: async ({ request, params, fetch }) => {
        const data = await request.formData();
        const name = data.get('name');
        const color = data.get('color');

        if (!name || !color) {
            return fail(400, { error: 'New name and color are required.' });
        }

        const response = await fetch(`/api/setlist/${params.id}/duplicate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, color })
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error || 'Failed to duplicate setlist.' });
        }

        const newSetlist = await response.json();
        throw redirect(303, `/setlist/${newSetlist.id}`);
    }
};