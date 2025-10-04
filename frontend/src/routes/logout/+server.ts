import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ cookies }) => {
    cookies.delete('jwt_token', { path: '/' });
    cookies.delete('active_band_id', { path: '/' });
    cookies.delete('user_bands', { path: '/' });
    throw redirect(303, '/login');
};