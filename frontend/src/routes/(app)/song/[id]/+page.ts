import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const { id } = params;
    const res = await fetch(`/api/song/${id}`);

    if (!res.ok) {
        throw error(res.status, 'Cette chanson est introuvable.');
    }

    const song = await res.json();
    return { song };
};
