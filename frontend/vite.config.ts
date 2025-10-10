import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		SvelteKitPWA({
			registerType: 'autoUpdate',

			includeAssets: ['favicon.svg', 'apple-touch-icon.png', 'web-app-manifest-192x192.png', 'web-app-manifest-512x512.png'],

			manifest: {
				name: 'lister',
				short_name: 'Lister',
				description: 'Une application pour gérer les setlists de votre groupe',
				theme_color: '#333333',
				background_color: '#ffffff',
				display: 'standalone',
				icons: [
					{
						src: 'web-app-manifest-192x192.png',
						sizes: '192x192',
						type: 'image/png'
					},
					{
						src: 'web-app-manifest-512x512.png',
						sizes: '512x512',
						type: 'image/png',
					},
					{
						src: 'web-app-manifest-512x512.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'any maskable'
					}
				]
			},
			workbox: {
				runtimeCaching: [
					{
						urlPattern: ({ url }) => {
							return url.pathname.startsWith('/login') || url.pathname.startsWith('/signup');
						},
						handler: 'NetworkOnly',
					},
				],
			},
			devOptions: {
				enabled: true,
				type: 'module',
			},
		})
	]
});