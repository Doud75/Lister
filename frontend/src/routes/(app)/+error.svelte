<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { toastStore } from '$lib/stores/toastStore';

    onMount(() => {
        if ($page.status === 403) {
            goto('/dashboard', { replaceState: true });
            return;
        }

        if ($page.status !== 404) {
            toastStore.error('Une erreur est survenue. Veuillez réessayer.');
            if (window.history.length > 1) {
                history.back();
            } else {
                goto('/', { replaceState: true });
            }
        }
    });
</script>

{#if $page.status === 404}
    <div class="flex min-h-[calc(100vh-4rem)] flex-col items-center justify-center px-4">
        <div class="text-center max-w-md">
            <p class="text-8xl font-black text-slate-200 dark:text-slate-700 select-none">404</p>
            <h1 class="mt-4 text-2xl font-semibold text-slate-900 dark:text-white">
                Page introuvable
            </h1>
            <p class="mt-2 text-slate-500 dark:text-slate-400">
                La page que vous cherchez n'existe pas ou a été déplacée.
            </p>
            <div class="mt-8 flex flex-col sm:flex-row gap-3 justify-center">
                <button
                    onclick={() => history.back()}
                    class="inline-flex items-center justify-center gap-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-5 py-2.5 text-sm font-semibold text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 5l-7 7 7 7"/></svg>
                    Retour
                </button>
                <a
                    href="/"
                    class="inline-flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500 transition-colors"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
                    Accueil
                </a>
            </div>
        </div>
    </div>
{/if}