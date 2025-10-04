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
            return fail(response.status, { error: result.error });
        }

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 7,
            sameSite: 'lax' as const
        };

        cookies.set('jwt_token', result.token, cookieOptions);

        if (result.bands && result.bands.length > 0) {
            cookies.set('user_bands', JSON.stringify(result.bands), cookieOptions);
            cookies.set('active_band_id', result.bands[0].id.toString(), cookieOptions);
        }

        throw redirect(303, '/');
    }
};