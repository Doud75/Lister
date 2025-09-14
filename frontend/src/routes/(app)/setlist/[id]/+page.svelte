<script lang="ts">
    import { page } from '$app/stores';
    import { formatDuration } from '$lib/utils/utils';
    import { dndzone } from 'svelte-dnd-action';
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';

    let { data, form } = $props();
    const setlistId = $page.params.id;
    let items = $state(data.setlistDetails.items);

    $effect(() => {
        if (form?.success && form?.deletedItemId) {
            items = items.filter((item) => item.id !== parseInt(form.deletedItemId, 10));
        }
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
        // CORRECTION CLÉ : On soumet le formulaire manuellement
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
                    <span>Total Duration: <span class="font-semibold">{formatDuration(totalDurationSeconds)}</span></span
                    >
                </div>
            </div>
            <div class="flex items-center gap-2">
                <a
                        href="/setlist/{setlistId}/add"
                        class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
                >
                    + Add Item
                </a>
                <form
                        method="POST"
                        action="?/deleteSetlist"
                        use:enhance={() => {
						return ({ cancel }) => {
							if (
								!confirm('Are you sure you want to delete this setlist? This action cannot be undone.')
							) {
								cancel();
							}
						};
					}}
                >
                    <Button type="submit" variant="secondary" autoWidth>Delete Setlist</Button>
                </form>
            </div>
        </div>
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        {#if items && items.length > 0}
            <!-- CORRECTION CLÉ : Ajout de l'id et changement du bouton en div -->
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
                                    <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke-width="1.5"
                                            stroke="currentColor"
                                            class="h-5 w-5"
                                    >
                                        <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                d="M3.75 9h16.5m-16.5 6.75h16.5"
                                        />
                                    </svg>
                                </div>
                                <span class="w-8 text-lg font-bold text-slate-400 dark:text-slate-500"
                                >{items.findIndex((i) => i.id === item.id) + 1}.</span
                                >
                                <div>
                                    <p class="font-semibold text-slate-800 dark:text-slate-100">{item.title.String}</p>
                                    {#if item.item_type === 'song'}
                                        <div
                                                class="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-slate-500 dark:text-slate-400"
                                        >
                                            <span>Duration: {formatItemDuration(item.duration_seconds.Int32)}</span>
                                            {#if item.tempo.Valid}
                                                <span class="hidden sm:inline">&bull;</span>
                                                <span>Tempo: {item.tempo.Int32} BPM</span>
                                            {/if}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                            <form method="POST" action="?/deleteItem" use:enhance>
                                <input type="hidden" name="itemId" value={item.id} />
                                <button
                                        type="submit"
                                        class="rounded-full p-2 text-slate-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-slate-700"
                                        aria-label="Delete item"
                                >
                                    <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M8.75 1A2.75 2.75 0 006 3.75v.443c-.795.077-1.58.22-2.365.468a.75.75 0 10.23 1.482l.149-.022.841 10.518A2.75 2.75 0 007.596 19h4.807a2.75 2.75 0 002.742-2.576l.84-10.518.149.022a.75.75 0 10.23-1.482A41.31 41.31 0 0014 4.193v-.443A2.75 2.75 0 0011.25 1h-2.5zM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4zM8.58 7.72a.75.75 0 00-1.5.06l.3 7.5a.75.75 0 101.5-.06l-.3-7.5zm4.34.06a.75.75 0 10-1.5-.06l-.3 7.5a.75.75 0 101.5.06l.3-7.5z" clip-rule="evenodd" /></svg>
                                </button>
                            </form>
                        </li>
                    {/each}
                </ul>
            </form>
        {:else}
            <div class="py-12 text-center">
                <svg
                        class="mx-auto h-12 w-12 text-slate-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z"
                /></svg
                >
                <h3 class="mt-2 text-sm font-semibold text-slate-900 dark:text-white">This setlist is empty</h3>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    Get started by adding your first item.
                </p>
            </div>
        {/if}
    </div>
</div>