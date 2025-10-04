import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ cookies, locals }) => {
    if (!locals.user) {
        throw redirect(303, '/login');
    }

    const userBandsCookie = cookies.get('user_bands');
    const userBands = userBandsCookie ? JSON.parse(userBandsCookie) : [];

    return {
        userBands,
        activeBandId: locals.activeBandId
    };
};