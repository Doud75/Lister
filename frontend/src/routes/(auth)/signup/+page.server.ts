import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, cookies, fetch }) => {
        const data = await request.formData();
        const bandName = data.get('bandName')?.toString() || '';
        const username = data.get('username')?.toString() || '';
        const password = data.get('password')?.toString() || '';

        const errors: Record<string, string> = {};

        const { validateUsername, validatePassword, sanitizeText } = await import('$lib/validation');

        const usernameValidation = validateUsername(username);
        if (!usernameValidation.success) {
            errors.username = usernameValidation.error!;
        }

        const passwordValidation = validatePassword(password);
        if (!passwordValidation.success) {
            errors.password = passwordValidation.error!;
        }

        if (Object.keys(errors).length > 0) {
            return fail(400, {
                data: { bandName, username },
                errors
            });
        }

        const sanitizedUsername = sanitizeText(username);
        const sanitizedBandName = sanitizeText(bandName);

        const response = await fetch('/api/auth/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                band_name: sanitizedBandName,
                username: sanitizedUsername,
                password
            })
        });

        const result = await response.json();
        if (!response.ok) {
            if (result.error?.includes('username')) {
                return fail(409, { data: { bandName, username }, errors: { username: result.error } });
            }
            return fail(response.status, { data: { bandName, username }, error: result.error });
        }

        const cookieOptions = {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 30,
            sameSite: 'lax' as const
        };

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