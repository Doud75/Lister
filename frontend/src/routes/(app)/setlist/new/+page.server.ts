import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch }) => {
        const data = await request.formData();
        const name = data.get('name');
        const color = data.get('color');

        const response = await fetch('/api/setlist', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, color })
        });

        const result = await response.json();

        if (!response.ok) {
            return fail(response.status, { error: result.error });
        }

        throw redirect(303, `/setlist/${result.id}`);
    }
};