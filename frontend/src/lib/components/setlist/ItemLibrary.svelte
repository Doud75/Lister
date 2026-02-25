<script lang="ts">
    import { navigating } from '$app/stores';
    import type { Song, Interlude } from '$lib/types';
    import type { SubmitFunction } from '@sveltejs/kit';
    import { enhance } from '$app/forms';

    let {
        songs,
        interludes,
        setlistId,
        createOptimisticUpdater,
        addingItemId = $bindable()
    } = $props<{
        songs: Song[];
        interludes: Interlude[];
        setlistId: number;
        createOptimisticUpdater: (itemType: 'song' | 'interlude') => SubmitFunction;
        addingItemId?: number | null;
    }>();

    let activeTab = $state<'songs' | 'interludes'>('songs');
</script>

<div class="flex h-[85vh] flex-col rounded-xl bg-white shadow-lg dark:bg-slate-800 lg:order-first">
    <div class="border-b border-slate-200 px-6 dark:border-slate-700">
        <nav class="-mb-px flex space-x-6" aria-label="Tabs">
            <button
                    onclick={() => (activeTab = 'songs')}
                    class="whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-200"
                    class:!border-indigo-500={activeTab === 'songs'}
                    class:!text-indigo-600={activeTab === 'songs'}
                    class:dark:!text-indigo-400={activeTab === 'songs'}
            >
                Songs ({songs.length})
            </button>
            <button
                    onclick={() => (activeTab = 'interludes')}
                    class="whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-200"
                    class:!border-indigo-500={activeTab === 'interludes'}
                    class:!text-indigo-600={activeTab === 'interludes'}
                    class:dark:!text-indigo-400={activeTab === 'interludes'}
            >
                Interludes ({interludes.length})
            </button>
        </nav>
    </div>
    <div class="border-b border-slate-200 p-4 dark:border-slate-700">
        {#if activeTab === 'songs'}
            <a
                    href="/song/new?redirectTo=/setlist/{setlistId}/add"
                    class="block w-full rounded-md bg-indigo-50 py-2.5 text-center text-sm font-semibold text-indigo-700 transition-colors hover:bg-indigo-100 dark:bg-indigo-500/10 dark:text-indigo-300 dark:hover:bg-indigo-500/20"
            >
                + Create New Song
            </a>
        {:else}
            <a
                    href="/interlude/new?redirectTo=/setlist/{setlistId}/add"
                    class="block w-full rounded-md bg-indigo-50 py-2.5 text-center text-sm font-semibold text-indigo-700 transition-colors hover:bg-indigo-100 dark:bg-indigo-500/10 dark:text-indigo-300 dark:hover:bg-indigo-500/20"
            >
                + Create New Interlude
            </a>
        {/if}
    </div>
    <div class="flex-grow overflow-y-auto p-6">
        {#if activeTab === 'songs'}
            {#if songs.length > 0}
                <ul class="space-y-2">
                    {#each songs as song (song.id)}
                        <li>
                            <form
                                    method="POST"
                                    action="?/addSong"
                                    use:enhance={createOptimisticUpdater('song')}
                            >
                                <input type="hidden" name="songId" value={song.id} />
                                <input type="hidden" name="title" value={song.title} />
                                <input
                                        type="hidden"
                                        name="duration"
                                        value={song.duration_seconds ?? 0}
                                />
                                <button
                                        type="submit"
                                        disabled={!!$navigating || addingItemId !== null}
                                        class="flex w-full items-center justify-between rounded-md p-3 text-left text-slate-800 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 dark:text-slate-200 dark:hover:bg-slate-700"
                                >
                                    <span>{song.title}</span>
                                    {#if addingItemId === song.id}
										<span
                                                class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-indigo-500"
                                        ></span>
                                    {:else}
										<span
                                                class="text-xl text-slate-400 transition-colors group-hover:text-indigo-500"
                                        >+</span
                                        >
                                    {/if}
                                </button>
                            </form>
                        </li>
                    {/each}
                </ul>
            {:else}
                <p class="p-4 text-center text-sm text-slate-500 dark:text-slate-400">
                    Your song library is empty.
                </p>
            {/if}
        {/if}
        {#if activeTab === 'interludes'}
            {#if interludes.length > 0}
                <ul class="space-y-2">
                    {#each interludes as interlude (interlude.id)}
                        <li>
                            <form
                                    method="POST"
                                    action="?/addInterlude"
                                    use:enhance={createOptimisticUpdater('interlude')}
                            >
                                <input type="hidden" name="interludeId" value={interlude.id} />
                                <input type="hidden" name="title" value={interlude.title} />
                                <input
                                        type="hidden"
                                        name="duration"
                                        value={interlude.duration_seconds ?? 0}
                                />
                                <button
                                        type="submit"
                                        disabled={!!$navigating || addingItemId !== null}
                                        class="flex w-full items-center justify-between rounded-md p-3 text-left text-slate-800 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 dark:text-slate-200 dark:hover:bg-slate-700"
                                >
                                    <span>{interlude.title}</span>
                                    {#if addingItemId === interlude.id}
										<span
                                                class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-indigo-500"
                                        ></span>
                                    {:else}
										<span
                                                class="text-xl text-slate-400 transition-colors group-hover:text-indigo-500"
                                        >+</span
                                        >
                                    {/if}
                                </button>
                            </form>
                        </li>
                    {/each}
                </ul>
            {:else}
                <p class="p-4 text-center text-sm text-slate-500 dark:text-slate-400">
                    You have no saved interludes.
                </p>
            {/if}
        {/if}
    </div>
</div>