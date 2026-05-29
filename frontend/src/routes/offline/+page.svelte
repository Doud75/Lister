<script lang="ts">
    function goBack() {
        if (window.history.length > 1) {
            history.back();
        } else {
            window.location.href = '/';
        }
    }
</script>

<svelte:head>
    <title>Hors ligne — Setly</title>
</svelte:head>

<div class="flex min-h-screen flex-col items-center justify-center bg-slate-100 dark:bg-slate-900 px-4">
    <div class="text-center max-w-md">
        <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-16 w-16 text-slate-300 dark:text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 3l18 18M8.111 8.111A3.5 3.5 0 0 0 7 10.5c0 .98.402 1.865 1.05 2.5M16.5 10.5a4.5 4.5 0 0 0-4.95-4.484M12 20.25v.008" />
        </svg>
        <h1 class="mt-4 text-2xl font-semibold text-slate-900 dark:text-white">
            Vous êtes hors ligne
        </h1>
        <p class="mt-2 text-slate-500 dark:text-slate-400">
            Cette page n'est pas disponible sans connexion Internet.
        </p>
        <div class="mt-8 flex flex-col sm:flex-row gap-3 justify-center">
            <button
                onclick={goBack}
                class="inline-flex items-center justify-center gap-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-5 py-2.5 text-sm font-semibold text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 transition-colors"
            >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 5l-7 7 7 7"/></svg>
                Retour
            </button>
            <button
                onclick={() => window.location.reload()}
                class="inline-flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500 transition-colors"
            >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M8 16H3v5"/></svg>
                Réessayer
            </button>
        </div>
    </div>
</div>
