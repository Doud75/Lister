import type { Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

type UserPayload = {
	user_id: number;
	exp: number;
	username: string;
	band_name: string;
	role: string;
};

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('jwt_token');
	const activeBandId = event.cookies.get('active_band_id');

	event.locals.token = token || null;
	event.locals.user = null;
	event.locals.activeBandId = activeBandId;

	if (token) {
		try {
			const decoded = jwtDecode<UserPayload>(token);
			if (decoded.exp * 1000 > Date.now()) {
				event.locals.user = {
					id: decoded.user_id,
					username: decoded.username,
					band_name: decoded.band_name,
					role: decoded.role
				};
			}
		} catch {
			// Le token est invalide, l'utilisateur et le token restent null
		}
	}

	return resolve(event);
};
