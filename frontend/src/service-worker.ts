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
const OFFLINE_URL = '/offline';

const ALL_CACHES = [STATIC_CACHE, PAGES_CACHE, API_CACHE];

self.addEventListener('install', (event: ExtendableEvent) => {
	event.waitUntil(
		caches
			.open(STATIC_CACHE)
			.then(async (cache) => {
				// Pre-cache offline fallback first so it's always available
				await cache.add(OFFLINE_URL).catch((err) => console.warn('[SW] Failed to cache offline page:', err));
				const urls = [...new Set(self.__WB_MANIFEST.map((e) => e.url))];
				await Promise.allSettled(
					urls.map((url) =>
						cache.add(url).catch((err) => console.warn('[SW] Failed to cache:', url, err))
					)
				);
			})
			.then(() => {
				return self.skipWaiting();
			})
	);
});

self.addEventListener('activate', (event: ExtendableEvent) => {
	event.waitUntil(
		caches
			.keys()
			.then((keys) =>
				Promise.all(keys.filter((k) => !ALL_CACHES.includes(k)).map((k) => caches.delete(k)))
			)
			.then(() => {
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

	// Fallback : assets pré-cachés non couverts (manifest.webmanifest, favicon, etc.)
	event.respondWith(cacheFirst(request));
});

async function cacheFirst(request: Request): Promise<Response> {
	const cached = await caches.match(request);
	if (cached) return cached;
	try {
		const response = await fetch(request);
		if (response.ok) {
			const cache = await caches.open(STATIC_CACHE);
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		return Response.error();
	}
}

async function networkFirstNavigation(request: Request): Promise<Response> {
	const cache = await caches.open(PAGES_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await cache.match(request);
		if (cached) return cached;
		const offline = await caches.match(OFFLINE_URL);
		return offline ?? Response.error();
	}
}

async function networkFirstData(request: Request): Promise<Response> {
	const cache = await caches.open(PAGES_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await cache.match(request, { ignoreSearch: true });
		return cached ?? Response.error();
	}
}

async function networkFirstApi(request: Request): Promise<Response> {
	const cache = await caches.open(API_CACHE);
	try {
		const response = await fetch(request);
		if (response.ok) {
			await cache.put(request, response.clone());
		}
		return response;
	} catch {
		const cached = await cache.match(request);
		return (
			cached ??
			new Response('{"error":"Hors ligne"}', {
				status: 503,
				headers: { 'Content-Type': 'application/json' }
			})
		);
	}
}
