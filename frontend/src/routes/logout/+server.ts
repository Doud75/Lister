import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ cookies }) => {
    cookies.delete('jwt_token', { path: '/' });
    cookies.delete('active_band_id', { path: '/' });
    cookies.delete('user_bands', { path: '/' });

    return json({ success: true, message: "Déconnexion réussie" });
};