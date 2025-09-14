import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, params, fetch }) => {
        const data = await request.formData();
        const payload = {
            title: data.get('title'),
            duration_seconds: Number(data.get('duration_seconds')) || null,
            tempo: Number(data.get('tempo')) || null,
            song_key: data.get('song_key') || null
        };

        const response = await fetch(`/api/song/${params.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        throw redirect(303, '/song');
    }
};