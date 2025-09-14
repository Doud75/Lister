import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch, params }) => {
        const data = await request.formData();
        const name = data.get('name');
        const color = data.get('color');

        const response = await fetch(`/api/setlist/${params.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, color })
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        throw redirect(303, `/setlist/${params.id}`);
    }
};