import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const res = await fetch(`/api/song/${params.id}`);
    if (!res.ok) {
        throw error(res.status, 'Song not found');
    }
    const song = await res.json();
    return { song };
};