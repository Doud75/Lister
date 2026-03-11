import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ data, url }) => {
    const { user } = data;
    const { pathname } = url;

    const isLoginOrSignup = pathname.startsWith('/login') || pathname.startsWith('/signup');
    const isPublicRoute = isLoginOrSignup || pathname.startsWith('/join');

    if (!user && !isPublicRoute) {
        throw redirect(307, '/login');
    }

    if (user && isLoginOrSignup) {
        throw redirect(307, '/');
    }

    return { user };
};