import path from 'node:path';
import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
    plugins: [svelte({ hot: !process.env.VITEST })],
    test: {
        globals: true,
        environment: 'jsdom',
        include: ['tests/**/*.test.{js,ts}'],
    },
    resolve: {
        alias: {
            '$lib': path.resolve(__dirname, './src/lib')
        }
    }
});