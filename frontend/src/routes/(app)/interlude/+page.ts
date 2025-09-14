import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const res = await fetch('/api/interlude');
    if (!res.ok) {
        throw error(res.status, 'Failed to fetch interludes');
    }
    const interludes = await res.json();
    return { interludes };
};