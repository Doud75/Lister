import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import {extractSongData} from "$lib/server/songActions";

export const actions: Actions = {
    default: async (event) => {
        const { fetch, url } = event;
        const payload = await extractSongData(event);

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