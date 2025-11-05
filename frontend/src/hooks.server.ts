import type { Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';
import { env } from '$env/dynamic/private';

type UserPayload = {
	user_id: number;
	exp: number;
};

const BACKEND_URL = env.BACKEND_INTERNAL_URL || 'http://backend:8089/api';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('jwt_token');
	const activeBandId = event.cookies.get('active_band_id');

	event.locals.token = token || null;
	event.locals.user = null;
	event.locals.activeBandId = activeBandId;

	if (token && activeBandId) {
		try {
			const decoded = jwtDecode<UserPayload>(token);
			if (decoded.exp * 1000 > Date.now()) {
				const userInfoUrl = `${BACKEND_URL}/user/info`;
				const userInfoRes = await fetch(userInfoUrl, {
					headers: {
						'Authorization': `Bearer ${token}`,
						'X-Band-ID': activeBandId
					}
				});

				if (userInfoRes.ok) {
					const userInfo = await userInfoRes.json();
					event.locals.user = {
						id: decoded.user_id,
						username: userInfo.username,
						band_name: userInfo.band_name,
						role: userInfo.role
					};
				}
			}
		} catch {
			// Le token est invalide, l'utilisateur et le token restent null
		}
	}

	return resolve(event);
};