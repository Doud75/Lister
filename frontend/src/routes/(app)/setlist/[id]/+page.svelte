<script lang="ts">
    import { page } from '$app/stores';
    import { formatDuration } from '$lib/utils/utils';
    import { dndzone } from 'svelte-dnd-action';
    import { enhance } from '$app/forms';

    let { data } = $props();
    const setlistId = $page.params.id;

    let items = $state(data.setlistDetails.items);

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
        document.getElementById('order-form')?.requestSubmit();
    }

    function formatItemDuration(seconds: number | null | undefined): string {
        if (seconds === null || seconds === undefined) return '-';
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
    }
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
                <div class="flex items-center gap-3">
					<span
                            class="block h-5 w-5 flex-shrink-0 rounded-full"
                            style="background-color: {data.setlistDetails.color};"
                    ></span>
                    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                        {data.setlistDetails.name}
                    </h1>
                </div>
                <div class="mt-2 flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400">
                    <a href="/" class="hover:underline">&larr; Back to Dashboard</a>
                    <span>&bull;</span>
                    <span
                    >Total Duration: <span class="font-semibold">{formatDuration(totalDurationSeconds)}</span></span
                    >
                </div>
            </div>
            <div class="flex items-center gap-4">
                <a href="/setlist/{setlistId}/edit" class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                    Edit Info
                </a>
                <a href="/setlist/{setlistId}/add" class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500">
                    + Add Item
                </a>
            </div>
        </div>
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        {#if items && items.length > 0}
            <form id="order-form" method="POST" action="?/updateOrder" use:enhance>
                <input type="hidden" name="itemIds" value={JSON.stringify(items.map((item) => item.id))} />

                <ul
                        class="divide-y divide-slate-200 dark:divide-slate-700"
                        use:dndzone={{ items: items, flipDurationMs: 300 }}
                        onconsider={handleDndConsider}
                        onfinalize={handleDndFinalize}
                >
                    {#each items as item (item.id)}
                        <li class="flex items-center justify-between py-3">
                            <div class="flex flex-grow items-center gap-4">
                                <div
                                        class="cursor-grab rounded-md p-2 text-slate-400 hover:bg-slate-100 active:cursor-grabbing dark:hover:bg-slate-700"
                                        aria-label="Drag to reorder"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/>
                                    </svg>
                                </div>

                                <span class="w-8 text-lg font-bold text-slate-400 dark:text-slate-500">{items.findIndex(i => i.id === item.id) + 1}.</span>
                                <div>
                                    <p class="font-semibold text-slate-800 dark:text-slate-100">{item.title.String}</p>
                                    {#if item.item_type === 'song'}
                                        <div class="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-slate-500 dark:text-slate-400">
                                            <span>Duration: {formatItemDuration(item.duration_seconds.Int32)}</span>
                                            {#if item.tempo.Valid}
                                                <span class="hidden sm:inline">&bull;</span>
                                                <span>Tempo: {item.tempo.Int32} BPM</span>
                                            {/if}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </li>
                    {/each}
                </ul>
            </form>
        {:else}
            <div class="py-12 text-center">
                <svg class="mx-auto h-12 w-12 text-slate-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z" />
                </svg>
                <h3 class="mt-2 text-sm font-semibold text-slate-900 dark:text-white">
                    This setlist is empty
                </h3>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    Get started by adding your first item.
                </p>
            </div>
        {/if}
    </div>
</div>