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
    updateItem: async ({ request, fetch }) => {
        const data = await request.formData();
        const itemId = data.get('itemId');
        const notes = data.get('notes');

        const response = await fetch(`/api/setlist/item/${itemId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ notes })
        });

        if (!response.ok) {
            return fail(response.status, { error: 'Failed to update item notes.' });
        }

        const updatedItem = await response.json();
        return { updated: true, item: updatedItem };
    }
};