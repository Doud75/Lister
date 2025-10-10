<script lang="ts">
    import {page} from '$app/stores';
    import {calculateTotalDuration, formatDuration} from '$lib/utils/utils';
    import {dragHandleZone} from 'svelte-dnd-action';
    import {enhance} from '$app/forms';
    import type {ActionData, PageData} from './$types';
    import {generateSetlistPdf} from '$lib/utils/pdfGenerator';
    import SetlistItem from '$lib/components/setlist/SetlistItem.svelte';
    import Modal from '$lib/components/ui/Modal.svelte';
    import EditItemForm from '$lib/components/setlist/EditItemForm.svelte';
    import type {SetlistItem as SetlistItemType} from '$lib/types';
    import DuplicateSetlistForm from '$lib/components/setlist/DuplicateSetlistForm.svelte';
    import {beforeNavigate} from "$app/navigation";
    import ActionDropdown from '$lib/components/ui/ActionDropdown.svelte';

    let {data, form}: { data: PageData; form: ActionData } = $props();
    const setlistId = $page.params.id;

    let items = $state<SetlistItemType[]>(data.setlistDetails.items);
    let isModalOpen = $state(false);
    let editingItem = $state<SetlistItemType | null>(null);
    let isDuplicateModalOpen = $state(false);

    beforeNavigate(() => {
        isModalOpen = false;
        isDuplicateModalOpen = false;
    });

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
                    updatedItem.duration_seconds = form.interlude.duration_seconds;
                    items[index] = updatedItem;
                }
            }

            if (form.item) {
                const index = items.findIndex((i: SetlistItemType) => i.id === form.item.id);
                if (index !== -1) {
                    items[index].notes = form.item.notes;
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
        generateSetlistPdf({...data.setlistDetails, items}, totalDurationSeconds);
    }

</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <div>
            <div class="flex flex-wrap items-center justify-between gap-4">
                <div class="flex items-center gap-3">
					<span
                            class="block h-5 w-5 flex-shrink-0 rounded-full"
                            style="background-color: {data.setlistDetails.color};"
                    ></span>
                    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                        {data.setlistDetails.name}
                    </h1>
                </div>
                <ActionDropdown>
                    {#snippet children({ close })}
                        <div class="py-1" role="none">
                            <button
                                onclick={() => {
                                    isDuplicateModalOpen = true;
                                    close();
                                }}
                                class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                                role="menuitem"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                                    <path d="M7 3.5A1.5 1.5 0 0 1 8.5 2h3.879a1.5 1.5 0 0 1 1.06.44l3.122 3.121A1.5 1.5 0 0 1 17 6.621V16.5a1.5 1.5 0 0 1-1.5 1.5h-7A1.5 1.5 0 0 1 7 16.5v-13Z" />
                                    <path d="M5 5.5A1.5 1.5 0 0 1 6.5 4h1V3H6.5A2.5 2.5 0 0 0 4 5.5v11A2.5 2.5 0 0 0 6.5 19h7a2.5 2.5 0 0 0 2.5-2.5v-1h1v1A3.5 3.5 0 0 1 13.5 20h-7A3.5 3.5 0 0 1 3 16.5v-11A3.5 3.5 0 0 1 6.5 2h1V4H5V5.5Z" />
                                </svg>
                                Dupliquer
                            </button>
                            <a
                                href="/setlist/{setlistId}/edit"
                                onclick={close}
                                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                                role="menuitem"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                                    <path d="M15.364 2.636a2 2 0 0 1 2.828 2.828l-9.9 9.9-3.182.354a.5.5 0 0 1-.556-.556l.354-3.182 9.9-9.9Zm-2.12 2.122-8.607 8.606-.202 1.818 1.818-.202 8.607-8.606-1.616-1.616Z" />
                                </svg>
                                Modifier les infos
                            </a>
                            <button
                                onclick={() => {
                                    downloadPdf();
                                    close();
                                }}
                                class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                                role="menuitem"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                                    <path d="M10.75 2.75a.75.75 0 0 0-1.5 0v8.614L6.295 8.235a.75.75 0 1 0-1.09 1.03l4.25 4.5a.75.75 0 0 0 1.09 0l4.25-4.5a.75.75 0 0 0-1.09-1.03l-2.955 3.129V2.75Z" />
                                    <path d="M3.5 12.75a.75.75 0 0 0-1.5 0v2.5A2.75 2.75 0 0 0 4.75 18h10.5A2.75 2.75 0 0 0 18 15.25v-2.5a.75.75 0 0 0-1.5 0v2.5c0 .69-.56 1.25-1.25 1.25H4.75c-.69 0-1.25-.56-1.25-1.25v-2.5Z" />
                                </svg>
                                Télécharger en PDF
                            </button>
                            <a
                                href="/setlist/{setlistId}/add"
                                class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                                role="menuitem"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
                                    <path d="M10 4a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 10 4Z" />
                                </svg>
                                Ajouter un item
                            </a>
                        </div>
                    {/snippet}
                </ActionDropdown>

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
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        {#if items && items.length > 0}
            <form id="order-form" method="POST" action="?/updateOrder" use:enhance>
                <input type="hidden" name="itemIds" value={JSON.stringify(items.map((item) => item.id))}/>
            </form>

            <ul
                    data-testid="setlist-items"
                    class="divide-y divide-slate-200 dark:divide-slate-700"
                    use:dragHandleZone={{ items: items, flipDurationMs: 300 }}
                    onconsider={handleDndConsider}
                    onfinalize={handleDndFinalize}
            >
                {#each items as item, index (item.id)}
                    <SetlistItem {item} {index} onEdit={openEditModal}/>
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

<Modal isOpen={isDuplicateModalOpen} onClose={() => (isDuplicateModalOpen = false)}>
    <DuplicateSetlistForm
            setlistName={data.setlistDetails.name}
            setColor={data.setlistDetails.color}
            close={() => (isDuplicateModalOpen = false)}
    />
</Modal>
<Modal isOpen={isModalOpen} onClose={closeEditModal}>
    {#if editingItem}
        <EditItemForm item={editingItem} close={closeEditModal}/>
    {/if}
</Modal>