<script lang="ts">
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';
    import type { PageData } from './$types';
    import { navigating } from '$app/stores';

    let { data, form } = $props<{ data: PageData, form?: { success?: boolean, addedSongId?: string, error?: string } }>();
    const { setlist, songs } = data;
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Add to "{setlist.name}"
        </h1>
        <a href="/setlist/{setlist.id}" class="mt-2 inline-block text-sm text-indigo-500 hover:underline dark:text-indigo-400">
            &larr; Back to setlist
        </a>
    </header>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
        <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
            <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Song Library</h2>
            {#if form?.error}
                <p class="mt-2 text-sm text-red-500">{form.error}</p>
            {/if}
            <div class="mt-4 max-h-[60vh] overflow-y-auto space-y-2">
                {#if songs.length > 0}
                    {#each songs as song (song.id)}
                        <form method="POST" action="?/addSong" use:enhance>
                            <input type="hidden" name="songId" value={song.id} />
                            <!-- CORRECTION ICI: Ajout des classes dark:text-slate-200 -->
                            <button type="submit" disabled={$navigating?.type === 'form'} class="flex w-full items-center justify-between rounded-md p-3 text-left text-slate-800 transition-colors hover:bg-slate-100 disabled:opacity-50 dark:text-slate-200 dark:hover:bg-slate-700">
                                <span>{song.title}</span>
                                {#if form?.success && form?.addedSongId == song.id.toString()}
                                    <span class="text-lg text-green-500">✓</span>
                                {:else}
                                    <span class="text-xl text-slate-400">+</span>
                                {/if}
                            </button>
                        </form>
                    {/each}
                {:else}
                    <div class="p-4 text-center">
                        <p class="text-sm text-slate-500 dark:text-slate-400">Your song library is empty.</p>
                        <a href="/song/new?redirectTo=/setlist/{setlist.id}/add" class="mt-2 inline-block text-sm font-semibold text-indigo-500 hover:underline dark:text-indigo-400">
                            Create your first song
                        </a>
                    </div>
                {/if}
            </div>
        </div>

        <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
            <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">
                Or Create a New Song
            </h2>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                This will add it to your library and you can then add it to the setlist.
            </p>
            <div class="mt-6">
                <a href="/song/new?redirectTo=/setlist/{setlist.id}/add">
                    <Button>Create New Song</Button>
                </a>
            </div>
        </div>
    </div>
</div>