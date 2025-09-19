<script lang="ts">
    import type { PageData } from './$types';
    import { invalidateAll } from '$app/navigation';
    import type { SetlistItem } from '$lib/types';
    import SetlistPreview from '$lib/components/setlist/SetlistPreview.svelte';
    import ItemLibrary from '$lib/components/setlist/ItemLibrary.svelte';
    import type { SubmitFunction } from '@sveltejs/kit';
    import {calculateTotalDuration} from "$lib/utils/utils";

    let { data } = $props<{ data: PageData }>();
    let { setlist, songs, interludes } = data;

    let items = $derived(data.setlist.items);
    let addingItemId = $state<number | null>(null);

    const totalDurationSeconds = $derived(calculateTotalDuration(items));

    function createOptimisticUpdater(itemType: 'song' | 'interlude'): SubmitFunction {
        return ({ formData }) => {
            const id = Number(formData.get(itemType === 'song' ? 'songId' : 'interludeId'));
            const title = formData.get('title')?.toString() || '...';
            const duration = Number(formData.get('duration')) || 0;

            addingItemId = id;

            const optimisticItem: Partial<SetlistItem> = {
                id: Date.now(),
                item_type: itemType,
                title: { String: title, Valid: true },
                position: items.length,
                duration_seconds: { Int32: duration, Valid: true }
            };

            items.push(optimisticItem as SetlistItem);

            return async ({ result }) => {
                if (result.type === 'failure') {
                    items = items.filter((item: SetlistItem) => item.id !== optimisticItem.id);
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
        <ItemLibrary
                {songs}
                {interludes}
                setlistId={setlist.id}
                {createOptimisticUpdater}
                bind:addingItemId
        />
        <SetlistPreview items={items} {totalDurationSeconds} />
    </div>
</div>