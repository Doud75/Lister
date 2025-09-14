import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const res = await fetch('/api/song');
    if (!res.ok) {
        throw error(res.status, 'Failed to fetch songs');
    }
    const songs = await res.json();
    return { songs };
};