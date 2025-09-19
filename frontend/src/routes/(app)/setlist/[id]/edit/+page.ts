import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { SetlistSummary } from '$lib/types';

export const load: PageLoad = async ({ fetch, params }) => {
    const { id } = params;
    const res = await fetch(`/api/setlist/${id}`);

    if (!res.ok) {
        throw error(res.status, 'Failed to load setlist data for editing.');
    }

    const setlist: SetlistSummary = await res.json();
    return { setlist };
};