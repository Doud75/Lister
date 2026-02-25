<script lang="ts">
    import type { PageData } from './$types';
    import { page } from '$app/stores';
    import { formatDuration } from '$lib/utils/utils';

    let { data }: { data: PageData } = $props();

    const song = $derived(data.song);
    const backHref = $derived($page.url.searchParams.get('from') ?? '/song');
    const editHref = $derived(`/song/${song.id}/edit?from=${$page.url.pathname}${$page.url.search}`);
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
                <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                    {song.title}
                </h1>
                <div class="mt-2 text-sm text-slate-500 dark:text-slate-400">
                    <a href={backHref} class="hover:underline">&larr; Retour</a>
                </div>
            </div>
            <a
                href={editHref}
                class="flex w-auto items-center gap-2 rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
            >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-4 w-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                </svg>
                Modifier
            </a>
        </div>

        <div class="mt-4 flex flex-wrap items-center gap-2">
            {#if song.album_name}
                <span class="rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700 dark:bg-slate-700 dark:text-slate-200">
                    {song.album_name}
                </span>
            {/if}
            {#if song.song_key}
                <span class="rounded-full bg-indigo-100 px-3 py-1 text-sm font-medium text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                    Tonalité : {song.song_key}
                </span>
            {/if}
            {#if song.tempo !== null}
                <span class="rounded-full bg-indigo-100 px-3 py-1 text-sm font-medium text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                    {song.tempo} BPM
                </span>
            {/if}
            {#if song.duration_seconds !== null}
                <span class="rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700 dark:bg-slate-700 dark:text-slate-200">
                    {formatDuration(song.duration_seconds)}
                </span>
            {/if}
        </div>

        {#if song.links}
            <div class="mt-3">
                <a
                    href={song.links}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-sm text-indigo-600 hover:underline dark:text-indigo-400"
                >
                    {song.links} &nearr;
                </a>
            </div>
        {/if}
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
            Paroles
        </h2>
        {#if song.lyrics}
            <p class="whitespace-pre-wrap leading-relaxed text-slate-800 dark:text-slate-100">
                {song.lyrics}
            </p>
        {:else}
            <p class="italic text-slate-400 dark:text-slate-500">
                Aucune parole renseignée.
            </p>
        {/if}
    </div>
</div>
