import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ data, url }) => {
    const { user } = data;
    const { pathname } = url;

    const isLoginOrSignup = pathname.startsWith('/login') || pathname.startsWith('/signup');
    const isPublicRoute = isLoginOrSignup || pathname.startsWith('/join') || pathname === '/offline';

    if (!user && !isPublicRoute) {
        // En offline, le server load ne peut pas résoudre l'user.
        // Rediriger vers /offline au lieu de /login pour éviter une boucle
        if (browser && !navigator.onLine) {
            throw redirect(307, '/offline');
        }
        throw redirect(307, '/login');
    }

    if (user && isLoginOrSignup) {
        throw redirect(307, '/');
    }

    return { user };
};