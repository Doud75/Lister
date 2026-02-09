import { redirect } from '@sveltejs/kit';
import type { Actions } from '@sveltejs/kit';

export const actions: Actions = {
    default: async ({ cookies, fetch }) => {
        const refreshToken = cookies.get('refresh_token');
        const jwtToken = cookies.get('jwt_token');
        
        cookies.delete('jwt_token', { path: '/' });
        cookies.delete('refresh_token', { path: '/' });
        cookies.delete('active_band_id', { path: '/' });
        cookies.delete('user_bands', { path: '/' });
        
        if (refreshToken && jwtToken) {
            try {
                await fetch('/api/auth/logout', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${jwtToken}`
                    },
                    body: JSON.stringify({ refresh_token: refreshToken })
                });
            } catch {
                // Ignore errors, user is already logged out client-side
            }
        }
        
        throw redirect(303, '/login');
    }
};
