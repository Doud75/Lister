import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, fetch, url }) => {
        const data = await request.formData();

        const min = parseInt(data.get('dur_min')?.toString() || '0') || 0;
        const sec = parseInt(data.get('dur_sec')?.toString() || '0') || 0;
        const payload = {
            title: data.get('title'),
            speaker: data.get('speaker') || null,
            script: data.get('script') || null,
            duration_seconds: (min === 0 && sec === 0) ? null : min * 60 + sec
        };

        const response = await fetch('/api/interlude', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        const redirectTo = url.searchParams.get('redirectTo') || '/';
        return { success: true, redirectTo: redirectTo };
    }
};