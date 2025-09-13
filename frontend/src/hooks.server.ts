import type { Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

type UserPayload = {
    user_id: number;
    band_id: number;
    exp: number;
};

export const handle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get('jwt_token');

    event.locals.user = null;
    event.locals.token = null;

    if (token) {
        event.locals.token = token;
        try {
            const decoded = jwtDecode<UserPayload>(token);
            if (decoded.exp * 1000 > Date.now()) {
                event.locals.user = { id: decoded.user_id, bandId: decoded.band_id };
            }
        } catch {
            // Le token est invalide, l'utilisateur et le token restent null
        }
    }

    return resolve(event);
};