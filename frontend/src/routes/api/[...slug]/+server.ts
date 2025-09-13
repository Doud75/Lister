import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_URL = 'http://backend:8089/api';

const handleProxy: RequestHandler = async ({ request, params, fetch, locals }) => {
    const url = `${BACKEND_URL}/${params.slug}`;

    const headers = new Headers(request.headers);

    if (locals.token) {
        headers.set('Authorization', `Bearer ${locals.token}`);
    }

    try {
        const response = await fetch(url, {
            method: request.method,
            headers: headers, // Utilise les nouveaux headers
            body: request.method !== 'GET' && request.method !== 'HEAD' ? request.body : null,
            duplex: 'half'
        });
        return response;
    } catch (e) {
        console.error('API proxy error:', e);
        throw error(500, 'Could not connect to the backend API.');
    }
};

export const GET = handleProxy;
export const POST = handleProxy;
export const PUT = handleProxy;
export const PATCH = handleProxy;
export const DELETE = handleProxy;