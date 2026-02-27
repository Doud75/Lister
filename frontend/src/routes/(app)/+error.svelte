<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';

    onMount(() => {
        if ($page.status === 403) {
            goto('/dashboard', { replaceState: true });
        }
    });
</script>

{#if $page.status !== 403}
    <div class="flex min-h-screen flex-col items-center justify-center bg-slate-100 dark:bg-slate-900">
        <div class="rounded-xl bg-white p-10 shadow-lg dark:bg-slate-800 text-center max-w-md">
            <p class="text-6xl font-bold text-slate-300 dark:text-slate-600">{$page.status}</p>
            <h1 class="mt-4 text-2xl font-semibold text-slate-900 dark:text-white">
                {$page.status === 404 ? 'Page introuvable' : 'Une erreur est survenue'}
            </h1>
            <p class="mt-2 text-slate-500 dark:text-slate-400">
                {$page.error?.message ?? 'Veuillez réessayer ou retourner à l\'accueil.'}
            </p>
            <a
                href="/"
                class="mt-6 inline-block rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500 transition-colors"
            >
                Retour à l'accueil
            </a>
        </div>
    </div>
{/if}
