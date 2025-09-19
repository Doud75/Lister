import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type {SetlistDetails} from "$lib/types";

export const load: PageLoad = async ({ fetch, params }) => {
    const { id } = params;

    try {
        const res = await fetch(`/api/setlist/${id}`);
        if (!res.ok) {
            throw error(res.status, 'Setlist not found');
        }
        const setlistDetails: SetlistDetails = await res.json();
        return {
            setlistDetails
        };
    } catch (err: any) {
        throw error(err.status || 500, err.body?.message || 'Could not load setlist.');
    }
};