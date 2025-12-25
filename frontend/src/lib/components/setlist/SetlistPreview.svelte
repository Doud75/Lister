<script lang="ts">
    import { enhance } from '$app/forms';
    import { dragHandle, dragHandleZone } from 'svelte-dnd-action';
    import { formatDuration } from '$lib/utils/utils';
    import type { SetlistItem } from '$lib/types';

    let {
        items: initialItems,
        totalDurationSeconds
    } = $props<{
        items: SetlistItem[];
        totalDurationSeconds: number;
    }>();

    let items = $derived(initialItems);

    function handleDndConsider(e: CustomEvent) {
        items = e.detail.items;
    }

    function handleDndFinalize(e: CustomEvent) {
        items = e.detail.items;
        document.getElementById('order-form-add-page')?.requestSubmit();
    }
</script>

<div class="flex h-[85vh] flex-col rounded-xl bg-white shadow-lg dark:bg-slate-800 lg:order-last">
    <div class="flex items-baseline justify-between p-6 pb-4">
        <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Current Setlist</h2>
        <span class="text-sm font-medium text-slate-500 dark:text-slate-400">
			Total: {formatDuration(totalDurationSeconds)}
		</span>
    </div>
    <div class="flex-grow overflow-y-auto px-6 pb-6">
        {#if items.length > 0}
            <form id="order-form-add-page" method="POST" action="?/updateOrder" use:enhance>
                <input
                        type="hidden"
                        name="itemIds"
                        value={JSON.stringify(items.map((item: SetlistItem) => item.id))}
                />
                <ul
                        data-testid="setlist-items"
                        class="space-y-3"
                        use:dragHandleZone={{ items: items, flipDurationMs: 300 }}
                        onconsider={handleDndConsider}
                        onfinalize={handleDndFinalize}
                >
                    {#each items as item (item.id)}
                        <li class="flex items-center gap-2">
                            <div
                                    use:dragHandle
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
                                ><path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        d="M3.75 9h16.5m-16.5 6.75h16.5"
                                /></svg
                                >
                            </div>
                            <div class="flex-grow">
                                {#if item.item_type === 'song'}
                                    <div class="flex items-center gap-4 rounded-md bg-slate-100 p-3 dark:bg-slate-700">
										<span class="font-bold text-slate-400 dark:text-slate-500"
                                        >{items.findIndex((i) => i.id === item.id) + 1}.</span
                                        >
                                        <span class="font-medium text-slate-800 dark:text-slate-100"
                                        >{item.title.String}</span
                                        >
                                    </div>
                                {:else}
                                    <div
                                            class="flex items-center gap-3 rounded-md border-l-4 border-teal-500 bg-teal-50 p-3 dark:border-teal-400 dark:bg-slate-700"
                                    >
                                        <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                viewBox="0 0 20 20"
                                                fill="currentColor"
                                                class="h-5 w-5 flex-shrink-0 text-teal-600 dark:text-teal-300"
                                        ><path
                                                fill-rule="evenodd"
                                                d="M18 5v8a2 2 0 01-2 2h-5l-5 4v-4H4a2 2 0 01-2-2V5a2 2 0 012-2h12a2 2 0 012 2zM9 10a1 1 0 11-2 0 1 1 0 012 0zm5 0a1 1 0 11-2 0 1 1 0 012 0z"
                                                clip-rule="evenodd"
                                        /></svg
                                        >
                                        <span class="font-medium italic text-teal-800 dark:text-teal-200"
                                        >{item.title.String}</span
                                        >
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