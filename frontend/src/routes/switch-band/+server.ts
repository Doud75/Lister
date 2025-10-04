import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, cookies }) => {
    const formData = await request.formData();
    const bandId = formData.get('bandId')?.toString();

    if (bandId) {
        cookies.set('active_band_id', bandId, {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 7,
            sameSite: 'lax'
        });
    }

    throw redirect(303, '/');
};