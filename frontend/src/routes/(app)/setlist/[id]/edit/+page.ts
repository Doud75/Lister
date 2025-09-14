import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const { id } = params;
    const res = await fetch(`/api/setlist/${id}`);

    if (!res.ok) {
        throw error(res.status, 'Failed to load setlist data for editing.');
    }

    const setlist = await res.json();
    return { setlist };
};