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

			includeAssets: ['favicon.svg', 'apple-touch-icon.png', 'pwa-192x192.png', 'pwa-512x512.png'],

			manifest: {
				name: 'lister',
				short_name: 'Lister',
				description: 'Une application pour gérer les setlists de votre groupe',
				theme_color: '#333333',
				background_color: '#ffffff',
				display: 'standalone',
				icons: [
					{
						src: 'pwa-192x192.png',
						sizes: '192x192',
						type: 'image/png'
					},
					{
						src: 'pwa-512x512.png',
						sizes: '512x512',
						type: 'image/png'
					}
				]
			},

			workbox: {
				navigateFallback: '/',
				globPatterns: ['client/**/*.{js,css,html,ico,png,svg,webp}', 'server/manifest.webmanifest']
			},

			devOptions: {
				enabled: true,
				type: 'module',
			},
		})
	]
});