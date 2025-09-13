import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch }) => {
    const fetchUserInfo = async () => {
        const res = await fetch('/api/user/info');
        if (!res.ok) throw error(res.status, 'Failed to fetch user info');
        return res.json();
    };

    const fetchSetlists = async () => {
        const res = await fetch('/api/setlist');
        if (!res.ok) throw error(res.status, 'Failed to fetch setlists');
        return res.json();
    };

    try {
        const [userInfo, setlists] = await Promise.all([fetchUserInfo(), fetchSetlists()]);
        return { userInfo, setlists };
    } catch (err: any) {
        throw error(err.status || 500, err.body?.message || 'Could not load dashboard data.');
    }
};