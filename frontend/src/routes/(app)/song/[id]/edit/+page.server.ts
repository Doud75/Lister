import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch, params }) => {
        const { id } = params;
        const data = await request.formData();

        const payload = {
            title: data.get('title'),
            album_name: data.get('album_name') || null,
            song_key: data.get('song_key') || null,
            duration_seconds: Number(data.get('duration_seconds')) || null,
            tempo: Number(data.get('tempo')) || null,
            lyrics: data.get('lyrics') || null
        };

        const response = await fetch(`/api/song/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        // Si la mise à jour réussit, on redirige vers la liste des chansons
        throw redirect(303, '/song');
    }
};