import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const { id } = params;
    const res = await fetch(`/api/song/${id}`);

    if (!res.ok) {
        throw error(res.status, 'Failed to load song data for editing.');
    }

    const song = await res.json();
    return { song };
};