import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const res = await fetch(`/api/interlude/${params.id}`);
    if (!res.ok) {
        throw error(res.status, 'Interlude not found');
    }
    const interlude = await res.json();
    return { interlude };
};