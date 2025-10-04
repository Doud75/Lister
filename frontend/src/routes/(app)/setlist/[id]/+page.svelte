<script lang="ts">
    import { page } from '$app/stores';
    import {calculateTotalDuration, formatDuration} from '$lib/utils/utils';
    import { dragHandleZone } from 'svelte-dnd-action';
    import { enhance } from '$app/forms';
    import type { ActionData, PageData } from './$types';
    import { generateSetlistPdf } from '$lib/utils/pdfGenerator';
    import SetlistItem from '$lib/components/setlist/SetlistItem.svelte';
    import Modal from '$lib/components/ui/Modal.svelte';
    import EditItemForm from '$lib/components/setlist/EditItemForm.svelte';
    import type { SetlistItem as SetlistItemType } from '$lib/types';

    let { data, form }: { data: PageData; form: ActionData } = $props();
    const setlistId = $page.params.id;

    let items = $state<SetlistItemType[]>(data.setlistDetails.items);
    let isModalOpen = $state(false);
    let editingItem = $state<SetlistItemType | null>(null);

    $effect(() => {
        if (form?.deleted) {
            const index = items.findIndex((item: SetlistItemType) => item.id === form.itemId);
            if (index !== -1) {
                items.splice(index, 1);
            }
        }
        if (form?.updatedSong) {
            const index = items.findIndex((item: SetlistItemType) => item.id === form.item.id);
            if (index !== -1) {
                items[index].notes = form.item.notes;
            }
        }
        if (form?.updatedInterlude) {
            const index = items.findIndex(
                (item) => item.item_type === 'interlude' && item.interlude_id?.Int32 === form.interlude.id
            );
            if (index !== -1) {
                const updatedItem = items[index] as SetlistItemType;
                if (updatedItem.item_type === 'interlude') {
                    updatedItem.title.String = form.interlude.title;
                    updatedItem.speaker = form.interlude.speaker;
                    updatedItem.script = form.interlude.script;
                    updatedItem.duration_seconds = form.interlude.duration_seconds;
                    items[index] = updatedItem;
                }
            }
        }
    });

    const totalDurationSeconds = $derived(calculateTotalDuration(items));

    function handleDndConsider(e: CustomEvent) {
        items = e.detail.items;
    }

    function handleDndFinalize(e: CustomEvent) {
        items = e.detail.items;
        document.getElementById('order-form')?.requestSubmit();
    }

    function openEditModal(item: SetlistItemType) {
        editingItem = item;
        isModalOpen = true;
    }

    function closeEditModal() {
        isModalOpen = false;
        editingItem = null;
    }

    function downloadPdf() {
        generateSetlistPdf({ ...data.setlistDetails, items }, totalDurationSeconds);
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
                    <a href="/" class="hover:underline">&larr; Back to Home</a>
                    <span>&bull;</span>
                    <span
                    >Total Duration: <span class="font-semibold"
                    >{formatDuration(totalDurationSeconds)}</span
                    ></span
                    >
                </div>
            </div>
            <div class="flex items-center gap-4">
                <button
                        onclick={downloadPdf}
                        type="button"
                        aria-label="Télécharger le PDF"
                        class="flex w-auto items-center gap-2 justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                        <path
                                d="M10.75 2.75a.75.75 0 0 0-1.5 0v8.614L6.295 8.235a.75.75 0 1 0-1.09 1.03l4.25 4.5a.75.75 0 0 0 1.09 0l4.25-4.5a.75.75 0 0 0-1.09-1.03l-2.955 3.129V2.75Z"
                        />
                        <path
                                d="M3.5 12.75a.75.75 0 0 0-1.5 0v2.5A2.75 2.75 0 0 0 4.75 18h10.5A2.75 2.75 0 0 0 18 15.25v-2.5a.75.75 0 0 0-1.5 0v2.5c0 .69-.56 1.25-1.25 1.25H4.75c-.69 0-1.25-.56-1.25-1.25v-2.5Z"
                        />
                    </svg>
                </button>
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
                    data-testid="setlist-items"
                    class="divide-y divide-slate-200 dark:divide-slate-700"
                    use:dragHandleZone={{ items: items, flipDurationMs: 300 }}
                    onconsider={handleDndConsider}
                    onfinalize={handleDndFinalize}
            >
                {#each items as item, index (item.id)}
                    <SetlistItem {item} {index} onEdit={openEditModal} />
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

<Modal isOpen={isModalOpen} onClose={closeEditModal}>
    {#if editingItem}
        <EditItemForm item={editingItem} close={closeEditModal} />
    {/if}
</Modal>