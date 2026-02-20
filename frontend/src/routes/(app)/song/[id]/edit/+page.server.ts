import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import {extractSongData} from "$lib/server/songActions";

export const actions: Actions = {
    default: async (event) => {
        const { fetch, params, request } = event;
        const { id } = params;
        const formData = await request.formData();
        const from = formData.get('from')?.toString() ?? '';
        const payload = await extractSongData(event);

        const response = await fetch(`/api/song/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const result = await response.json();
            return fail(response.status, { error: result.error });
        }

        const redirectTo = from && from.startsWith('/') ? from : '/song';
        throw redirect(303, redirectTo);
    }
};