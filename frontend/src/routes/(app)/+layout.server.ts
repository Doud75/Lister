import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ cookies, locals, url }) => {
    if (!locals.user) {
        throw redirect(303, '/login');
    }

    const userBandsCookie = cookies.get('user_bands');
    const userBands = userBandsCookie ? JSON.parse(userBandsCookie) : [];
    const activeBandId = locals.activeBandId;

    if (!activeBandId && url.pathname !== '/dashboard' && !url.pathname.startsWith('/settings')) {
        throw redirect(303, '/dashboard');
    }

    return {
        user: locals.user,
        userBands,
        activeBandId
    };
};