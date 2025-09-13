<script lang="ts">
    import { enhance } from '$app/forms';
    import type { PageData } from './$types';
    import { navigating } from '$app/stores';
    import { invalidateAll } from '$app/navigation';
    import { formatDuration } from '$lib/utils/utils';
    import { dndzone } from 'svelte-dnd-action';

    let { data } = $props<{ data: PageData }>();

    let items = $state(data.setlist.items);
    const { setlist, songs, interludes } = data;

    let activeTab = $state<'songs' | 'interludes'>('songs');
    let addingItemId = $state<number | null>(null);

    $effect(() => {
        items = data.setlist.items;
    });

    const totalDurationSeconds = $derived(
        items.reduce((total, item) => {
            const duration = item.item_type === 'song' ? item.duration_seconds?.Int32 ?? 0 : 0;
            return total + duration;
        }, 0)
    );

    function handleDndConsider(e: CustomEvent) {
        items = e.detail.items;
    }

    function handleDndFinalize(e: CustomEvent) {
        items = e.detail.items;
        document.getElementById('order-form-add-page')?.requestSubmit();
    }

    function createOptimisticUpdater(itemType: 'song' | 'interlude') {
        return ({ formData }: { formData: FormData }) => {
            const id = Number(formData.get(itemType === 'song' ? 'songId' : 'interludeId'));
            const title = formData.get('title')?.toString() || '...';
            const duration = Number(formData.get('duration')) || 0;

            addingItemId = id;

            const optimisticItem = {
                id: Date.now(),
                item_type: itemType,
                title: { String: title, Valid: true },
                position: items.length,
                duration_seconds: { Int32: duration, Valid: true }
            };

            items.push(optimisticItem as any);

            return async ({ result }: { result: any }) => {
                if (result.type === 'failure') {
                    items = items.filter((item) => item.id !== optimisticItem.id);
                }
                await invalidateAll();
                addingItemId = null;
            };
        };
    }
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Build Setlist: <span class="text-indigo-500">{setlist.name}</span>
        </h1>
        <a
                href="/setlist/{setlist.id}"
                class="mt-2 inline-block text-sm text-indigo-500 hover:underline dark:text-indigo-400"
        >
            &larr; Done Editing
        </a>
    </header>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
        <div
                class="flex h-[85vh] flex-col rounded-xl bg-white shadow-lg dark:bg-slate-800 lg:order-last"
        >
            <div class="flex items-baseline justify-between p-6 pb-4">
                <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">
                    Current Setlist
                </h2>
                <span class="text-sm font-medium text-slate-500 dark:text-slate-400">
					Total: {formatDuration(totalDurationSeconds)}
				</span>
            </div>
            <div class="flex-grow overflow-y-auto px-6 pb-6">
                {#if items.length > 0}
                    <form id="order-form-add-page" method="POST" action="?/updateOrder" use:enhance>
                        <input type="hidden" name="itemIds" value={JSON.stringify(items.map((item) => item.id))} />
                        <ul
                                class="space-y-3"
                                use:dndzone={{ items: items, flipDurationMs: 300 }}
                                onconsider={handleDndConsider}
                                onfinalize={handleDndFinalize}
                        >
                            {#each items as item (item.id)}
                                <li class="flex items-center gap-2">
                                    <div
                                            class="cursor-grab rounded-md p-2 text-slate-400 hover:bg-slate-100 active:cursor-grabbing dark:hover:bg-slate-700"
                                            aria-label="Drag to reorder"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/></svg>
                                    </div>
                                    <div class="flex-grow">
                                        {#if item.item_type === 'song'}
                                            <div class="flex items-center gap-4 rounded-md bg-slate-100 p-3 dark:bg-slate-700">
                                                <span class="font-bold text-slate-400 dark:text-slate-500">{items.findIndex(i => i.id === item.id) + 1}.</span>
                                                <span class="font-medium text-slate-800 dark:text-slate-100">{item.title.String}</span>
                                            </div>
                                        {:else}
                                            <div class="flex items-center gap-3 rounded-md border-l-4 border-teal-500 bg-teal-50 p-3 dark:border-teal-400 dark:bg-slate-700">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5 flex-shrink-0 text-teal-600 dark:text-teal-300"><path fill-rule="evenodd" d="M18 5v8a2 2 0 01-2 2h-5l-5 4v-4H4a2 2 0 01-2-2V5a2 2 0 012-2h12a2 2 0 012 2zM9 10a1 1 0 11-2 0 1 1 0 012 0zm5 0a1 1 0 11-2 0 1 1 0 012 0z" clip-rule="evenodd"/></svg>
                                                <span class="font-medium italic text-teal-800 dark:text-teal-200">{item.title.String}</span>
                                            </div>
                                        {/if}
                                    </div>
                                </li>
                            {/each}
                        </ul>
                    </form>
                {:else}
                    <div class="flex h-full items-center justify-center">
                        <p class="py-10 text-center text-sm text-slate-500 dark:text-slate-400">
                            Add items from your library to get started.
                        </p>
                    </div>
                {/if}
            </div>
        </div>

        <div class="flex h-[85vh] flex-col rounded-xl bg-white shadow-lg dark:bg-slate-800 lg:order-first">
            <div class="border-b border-slate-200 px-6 dark:border-slate-700">
                <nav class="-mb-px flex space-x-6" aria-label="Tabs">
                    <button onclick={() => (activeTab = 'songs')} class="whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-200" class:!border-indigo-500={activeTab === 'songs'} class:!text-indigo-600={activeTab === 'songs'} class:dark:!text-indigo-400={activeTab === 'songs'}>
                        Songs ({songs.length})
                    </button>
                    <button onclick={() => (activeTab = 'interludes')} class="whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-200" class:!border-indigo-500={activeTab === 'interludes'} class:!text-indigo-600={activeTab === 'interludes'} class:dark:!text-indigo-400={activeTab === 'interludes'}>
                        Interludes ({interludes.length})
                    </button>
                </nav>
            </div>
            <div class="border-b border-slate-200 p-4 dark:border-slate-700">
                {#if activeTab === 'songs'}
                    <a href="/song/new?redirectTo=/setlist/{setlist.id}/add" class="block w-full rounded-md bg-indigo-50 py-2.5 text-center text-sm font-semibold text-indigo-700 transition-colors hover:bg-indigo-100 dark:bg-indigo-500/10 dark:text-indigo-300 dark:hover:bg-indigo-500/20">
                        + Create New Song
                    </a>
                {:else}
                    <a href="/interlude/new?redirectTo=/setlist/{setlist.id}/add" class="block w-full rounded-md bg-indigo-50 py-2.5 text-center text-sm font-semibold text-indigo-700 transition-colors hover:bg-indigo-100 dark:bg-indigo-500/10 dark:text-indigo-300 dark:hover:bg-indigo-500/20">
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
                                    <form method="POST" action="?/addSong" use:enhance={createOptimisticUpdater('song')}>
                                        <input type="hidden" name="songId" value={song.id} />
                                        <input type="hidden" name="title" value={song.title} />
                                        <input type="hidden" name="duration" value={song.duration_seconds?.Int32 ?? 0} />
                                        <button type="submit" disabled={$navigating || addingItemId !== null} class="flex w-full items-center justify-between rounded-md p-3 text-left text-slate-800 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 dark:text-slate-200 dark:hover:bg-slate-700">
                                            <span>{song.title}</span>
                                            {#if addingItemId === song.id}
                                                <span class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-indigo-500"></span>
                                            {:else}
                                                <span class="text-xl text-slate-400 transition-colors group-hover:text-indigo-500">+</span>
                                            {/if}
                                        </button>
                                    </form>
                                </li>
                            {/each}
                        </ul>
                    {:else}
                        <p class="p-4 text-center text-sm text-slate-500 dark:text-slate-400">Your song library is empty.</p>
                    {/if}
                {/if}
                {#if activeTab === 'interludes'}
                    {#if interludes.length > 0}
                        <ul class="space-y-2">
                            {#each interludes as interlude (interlude.id)}
                                <li>
                                    <form method="POST" action="?/addInterlude" use:enhance={createOptimisticUpdater('interlude')}>
                                        <input type="hidden" name="interludeId" value={interlude.id} />
                                        <input type="hidden" name="title" value={interlude.title} />
                                        <input type="hidden" name="duration" value={interlude.duration_seconds?.Int32 ?? 0} />
                                        <button type="submit" disabled={$navigating || addingItemId !== null} class="flex w-full items-center justify-between rounded-md p-3 text-left text-slate-800 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 dark:text-slate-200 dark:hover:bg-slate-700">
                                            <span>{interlude.title}</span>
                                            {#if addingItemId === interlude.id}
                                                <span class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-indigo-500"></span>
                                            {:else}
                                                <span class="text-xl text-slate-400 transition-colors group-hover:text-indigo-500">+</span>
                                            {/if}
                                        </button>
                                    </form>
                                </li>
                            {/each}
                        </ul>
                    {:else}
                        <p class="p-4 text-center text-sm text-slate-500 dark:text-slate-400">You have no saved interludes.</p>
                    {/if}
                {/if}
            </div>
        </div>
    </div>
</div>