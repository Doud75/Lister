import type { PageLoad } from './$types';
import { error, redirect } from '@sveltejs/kit';
import type { BandMember } from '$lib/types';

export const load: PageLoad = async ({ fetch, parent }) => {
    const { user, activeBandId } = await parent();

    if (user?.role !== 'admin') {
        throw redirect(303, '/');
    }

    if (!activeBandId) {
        throw error(400, 'Active band not selected');
    }

    const res = await fetch(`/api/bands/${activeBandId}/members`);
    if (!res.ok) throw error(res.status, 'Failed to fetch members.');

    const members: BandMember[] = await res.json();
    return { members };
};
