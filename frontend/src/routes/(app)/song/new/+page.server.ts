import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch, url }) => {
        const data = await request.formData();

        const payload = {
            title: data.get('title'),
            album_name: data.get('album_name') || null,
            song_key: data.get('song_key') || null,
            duration_seconds: Number(data.get('duration_seconds')) || null,
            tempo: Number(data.get('tempo')) || null,
            lyrics: data.get('lyrics') || null,
            links: data.get('links') || null
        };

        const response = await fetch('/api/song', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        const redirectTo = url.searchParams.get('redirectTo') || '/song';
        return { success: true, redirectTo: redirectTo };
    }
};