/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

declare const self: ServiceWorkerGlobalScope;

declare global {
	interface ServiceWorkerGlobalScope {
		__WB_MANIFEST: Array<{ url: string; revision: string | null }>;
	}
}

const STATIC_CACHE = 'static-v1';
const PAGES_CACHE = 'pages-v1';
const API_CACHE = 'api-v1';

const ALL_CACHES = [STATIC_CACHE, PAGES_CACHE, API_CACHE];

self.addEventListener('install', (event: ExtendableEvent) => {
	console.log('[SW] Installing...');
	event.waitUntil(
		caches
			.open(STATIC_CACHE)
			.then(async (cache) => {
				const urls = [...new Set(self.__WB_MANIFEST.map((e) => e.url))];
				console.log('[SW] Pre-caching', urls.length, 'assets');
				await Promise.allSettled(
					urls.map((url) =>
						cache.add(url).catch((err) => console.warn('[SW] Failed to cache:', url, err))
					)
				);
			})
			.then(() => {
				console.log('[SW] Install complete, skipping waiting');
				return self.skipWaiting();
			})
	);
});

self.addEventListener('activate', (event: ExtendableEvent) => {
	console.log('[SW] Activating...');
	event.waitUntil(
		caches
			.keys()
			.then((keys) =>
				Promise.all(keys.filter((k) => !ALL_CACHES.includes(k)).map((k) => caches.delete(k)))
			)
			.then(() => {
				console.log('[SW] Activated, claiming clients');
				return self.clients.claim();
			})
	);
});

self.addEventListener('fetch', (event: FetchEvent) => {
	const { request } = event;
	if (request.method !== 'GET') return;

	const url = new URL(request.url);
	if (url.origin !== self.location.origin) return;

	const path = url.pathname;

	// Version check: réseau uniquement, échoue silencieusement hors-ligne
	if (path === '/_app/version.json') return;

	// Chunks JS/CSS pré-cachés : cache-first (contenu adressé)
	if (path.startsWith('/_app/')) {
		event.respondWith(cacheFirst(request));
		return;
	}

	// Auth et routes publiques : toujours réseau, jamais de cache
	if (
		path.startsWith('/api/auth') ||
		path.startsWith('/login') ||
		path.startsWith('/signup') ||
		path === '/logout'
	)
		return;

	// console.log('[SW] fetch', request.mode, path); // décommenter pour debug verbeux

	if (request.mode === 'navigate') {
		event.respondWith(networkFirstNavigation(request));
		return;
	}

	if (path.includes('/__data.json')) {
		event.respondWith(networkFirstData(request));
		return;
	}

	if (path.startsWith('/api/')) {
		event.respondWith(networkFirstApi(request));
		return;
	}
});

async function cacheFirst(request: Request): Promise<Response> {
	const cached = await caches.match(request);
	if (cached) return cached;
	const response = await fetch(request);
	if (response.ok) {
		const cache = await caches.open(STATIC_CACHE);
		await cache.put(request, response.clone());
	}
	return response;
}

async function networkFirstNavigation(request: Request): Promise<Response> {
	console.log('[SW] navigate ->', request.url);
	const cache = await caches.open(PAGES_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		console.log('[SW] navigate network ok');
		return response;
	} catch {
		const cached = await cache.match(request);
		console.warn('[SW] navigate offline, cache hit:', !!cached, request.url);
		if (cached) return cached;
		return Response.error();
	}
}

async function networkFirstData(request: Request): Promise<Response> {
	console.log('[SW] data ->', request.url);
	const cache = await caches.open(PAGES_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await cache.match(request);
		console.warn('[SW] data offline, cache hit:', !!cached, request.url);
		return (
			cached ??
			new Response('{"offline":true}', {
				status: 503,
				headers: { 'Content-Type': 'application/json' }
			})
		);
	}
}

async function networkFirstApi(request: Request): Promise<Response> {
	console.log('[SW] api ->', request.url);
	const cache = await caches.open(API_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await cache.match(request);
		console.warn('[SW] api offline, cache hit:', !!cached, request.url);
		return (
			cached ??
			new Response('{"error":"Hors ligne"}', {
				status: 503,
				headers: { 'Content-Type': 'application/json' }
			})
		);
	}
}
