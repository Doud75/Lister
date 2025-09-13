import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ data, url }) => {
    const { user } = data;
    const { pathname } = url;

    const isAuthRoute = pathname.startsWith('/login') || pathname.startsWith('/signup') || pathname.startsWith('/join');

    if (!user && !isAuthRoute) {
        throw redirect(307, '/login');
    }

    if (user && isAuthRoute) {
        throw redirect(307, '/');
    }

    return { user };
};