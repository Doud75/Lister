import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, cookies, fetch }) => {
        const data = await request.formData();
        const username = data.get('username');
        const password = data.get('password');

        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const result = await response.json();
        if (!response.ok) {
            return fail(response.status, { error: result.error, code: result.code });
        }

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 30,
            sameSite: 'lax' as const
        };

        cookies.delete('jwt_token', { path: '/' });
        cookies.delete('refresh_token', { path: '/' });
        cookies.delete('user_bands', { path: '/' });
        cookies.delete('active_band_id', { path: '/' });

        cookies.set('jwt_token', result.token, cookieOptions);
        
        if (result.refresh_token) {
            cookies.set('refresh_token', result.refresh_token, {
                ...cookieOptions,
                maxAge: 60 * 60 * 24 * 30
            });
        }

        if (result.bands && result.bands.length > 0) {
            cookies.set('user_bands', JSON.stringify(result.bands), cookieOptions);
            cookies.set('active_band_id', result.bands[0].id.toString(), cookieOptions);
        }

        throw redirect(303, '/');
    }
};