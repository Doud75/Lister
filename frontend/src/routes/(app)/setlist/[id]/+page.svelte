<script lang="ts">
    import { page } from '$app/stores';
    import { formatDuration } from '$lib/utils/utils';
    import { dndzone } from 'svelte-dnd-action';
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';
    import type { ActionData } from './$types';

    let { data, form }: { data: any; form: ActionData } = $props();
    const setlistId = $page.params.id;

    let items = $state(data.setlistDetails.items);

    let isModalOpen = $state(false);
    let editingItem = $state<any>(null);
    let notesContent = $state('');

    $effect(() => {
        if (form?.deleted) {
            const index = items.findIndex((item) => item.id === form.itemId);
            if (index !== -1) {
                items.splice(index, 1);
            }
        }
        if (form?.updated) {
            const index = items.findIndex((item) => item.id === form.item.id);
            if (index !== -1) {
                items[index].notes = form.item.notes;
            }
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
        document.getElementById('order-form')?.requestSubmit();
    }

    function openEditModal(item: any) {
        editingItem = item;
        notesContent = item.notes?.String || '';
        isModalOpen = true;
    }

    function closeEditModal() {
        isModalOpen = false;
        editingItem = null;
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
                    >Total Duration: <span class="font-semibold"
                    >{formatDuration(totalDurationSeconds)}</span
                    ></span
                    >
                </div>
            </div>
            <div class="flex items-center gap-4">
                <a
                        href="/setlist/{setlistId}/edit"
                        class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
                >
                    Edit Info
                </a>
                <a
                        href="/setlist/{setlistId}/add"
                        class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
                >
                    + Add Item
                </a>
            </div>
        </div>
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        {#if items && items.length > 0}
            <form id="order-form" method="POST" action="?/updateOrder" use:enhance>
                <input type="hidden" name="itemIds" value={JSON.stringify(items.map((item) => item.id))} />
            </form>

            <ul
                    class="divide-y divide-slate-200 dark:divide-slate-700"
                    use:dndzone={{ items: items, flipDurationMs: 300 }}
                    onconsider={handleDndConsider}
                    onfinalize={handleDndFinalize}
            >
                {#each items as item (item.id)}
                    <li class="flex items-center justify-between gap-3 py-3">
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
                            <div class="flex-grow">
                                <p class="font-semibold text-slate-800 dark:text-slate-100">
                                    {item.title.String}
                                </p>
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
                                {#if item.notes?.Valid && item.notes.String}
                                    <p class="mt-1 text-xs italic text-slate-500 dark:text-slate-400">
                                        {item.notes.String}
                                    </p>
                                {/if}
                            </div>
                        </div>

                        <div class="flex items-center gap-2">
                            <button
                                    onclick={() => openEditModal(item)}
                                    type="button"
                                    class="rounded-md p-2 text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-700 dark:hover:text-slate-200"
                                    aria-label="Edit notes"
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
                                        d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10"
                                /></svg
                                >
                            </button>
                            <form method="POST" action="?/deleteItem" use:enhance>
                                <input type="hidden" name="itemId" value={item.id} />
                                <button
                                        type="submit"
                                        class="rounded-md p-2 text-slate-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-500/10 dark:hover:text-red-400"
                                        aria-label="Remove item"
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
                                            d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.134-2.033-2.134H8.033C6.91 2.75 6 3.704 6 4.874v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                                    /></svg
                                    >
                                </button>
                            </form>
                        </div>
                    </li>
                {/each}
            </ul>
        {:else}
            <div class="py-12 text-center">
                <svg
                        class="mx-auto h-12 w-12 text-slate-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                >
                    <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z"
                    />
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

{#if isModalOpen && editingItem}
    <div
            class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            onclick={closeEditModal}
            role="dialog"
            aria-modal="true"
    >
        <div
                class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl dark:bg-slate-800"
                onclick={(event) => event.stopPropagation()}
        >
            <h3 class="text-lg font-semibold text-slate-900 dark:text-white">
                Edit Notes for {editingItem.title.String}
            </h3>
            <form
                    method="POST"
                    action="?/updateItem"
                    use:enhance={() => {
					return async ({ update }) => {
						await update();
						closeEditModal();
					};
				}}
                    class="mt-4 space-y-4"
            >
                <input type="hidden" name="itemId" value={editingItem.id} />
                <div>
                    <label for="notes" class="sr-only">Notes</label>
                    <textarea
                            id="notes"
                            name="notes"
                            bind:value={notesContent}
                            rows="4"
                            class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                            placeholder="Add a comment, a reminder, etc."
                    ></textarea>
                </div>
                <div class="flex justify-end gap-3">
                    <Button type="button" variant="secondary" onclick={closeEditModal} autoWidth>Cancel</Button>
                    <Button type="submit" autoWidth>Save Notes</Button>
                </div>
            </form>
        </div>
    </div>
{/if}