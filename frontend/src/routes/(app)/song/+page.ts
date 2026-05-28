import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch }) => {
    const res = await fetch('/api/song');
    if (!res.ok) throw error(res.status, 'Failed to fetch song library.');
    return { songs: await res.json() };
};
