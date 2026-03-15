<script lang="ts">
    import { page } from '$app/stores';
</script>

<svelte:head>
    <title>Erreur {$page.status}</title>
</svelte:head>

<div class="flex min-h-screen flex-col items-center justify-center bg-slate-100 px-4">
    <div class="rounded-xl bg-white p-10 shadow-lg text-center max-w-md w-full">
        <p class="text-7xl font-black text-slate-200 select-none">{$page.status}</p>
        <h1 class="mt-4 text-xl font-semibold text-slate-900">
            {#if $page.status === 404}
                Page introuvable
            {:else if $page.status === 403}
                Accès refusé
            {:else}
                Une erreur est survenue
            {/if}
        </h1>
        <p class="mt-2 text-sm text-slate-500">
            {$page.error?.message ?? 'Veuillez réessayer ou retourner à l\'accueil.'}
        </p>
        <a
            href="/login"
            class="mt-6 inline-block rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500 transition-colors"
        >
            Retour à la connexion
        </a>
    </div>
</div>