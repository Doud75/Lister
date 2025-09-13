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
    }
};