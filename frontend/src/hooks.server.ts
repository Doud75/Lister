import type { Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';
import { env } from '$env/dynamic/private';

type UserPayload = {
    user_id: number;
    exp: number;
};

const BACKEND_URL = env.BACKEND_INTERNAL_URL || 'http://backend:8089/api';

export const handle: Handle = async ({ event, resolve }) => {
    if (event.url.pathname === '/logout') {
        return resolve(event);
    }

    const token = event.cookies.get('jwt_token');
    const refreshToken = event.cookies.get('refresh_token');
    const activeBandId = event.cookies.get('active_band_id');

    event.locals.token = token || null;
    event.locals.user = null;
    event.locals.activeBandId = activeBandId;

    let decoded: UserPayload | null = null;
    let needsRefresh = false;

    if (token) {
        try {
            decoded = jwtDecode<UserPayload>(token);
            const expiresIn = decoded.exp * 1000 - Date.now();
            needsRefresh = expiresIn < 0;
        } catch (error) {
            console.log('[AUTH] JWT decode failed:', error);
            needsRefresh = true;
        }
    } else {
        console.log('[AUTH] No JWT token found');
    }

    if ((!token || needsRefresh) && refreshToken) {
        try {
            const refreshResponse = await fetch(`${BACKEND_URL}/auth/refresh`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken })
            });

            if (refreshResponse.ok) {
                const refreshData = await refreshResponse.json();
                
                const cookieOptions = {
                    path: '/',
                    httpOnly: true,
                    secure: process.env.NODE_ENV === 'production',
                    sameSite: 'lax' as const
                };
                
                event.cookies.set('jwt_token', refreshData.token, {
                    ...cookieOptions,
                    maxAge: 60 * 60 * 24 * 30
                });
                
                event.cookies.set('refresh_token', refreshData.refresh_token, {
                    ...cookieOptions,
                    maxAge: 60 * 60 * 24 * 30
                });
                
                event.locals.token = refreshData.token;
                decoded = jwtDecode<UserPayload>(refreshData.token);
            } else {
                console.error('[AUTH] Refresh failed with status:', refreshResponse.status);
            }
        } catch (error) {
            console.error('[AUTH] Refresh error:', error);
        }
    }

    if (decoded && decoded.exp * 1000 > Date.now() && activeBandId) {
        try {
            const userInfoUrl = `${BACKEND_URL}/user/info`;
            const userInfoRes = await fetch(userInfoUrl, {
                headers: {
                    'Authorization': `Bearer ${event.locals.token || token}`,
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
        } catch (error) {
            console.error('[AUTH] Failed to fetch user info:', error);
        }
    }

    return resolve(event);
};