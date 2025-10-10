<script lang="ts">
    import type { SetlistItem } from '$lib/types';
    import { enhance } from '$app/forms';
    import { dragHandle } from 'svelte-dnd-action';
    import { formatItemDuration } from '$lib/utils/utils';

    let { item, index, onEdit } = $props<{
        item: SetlistItem;
        index: number;
        onEdit: (item: SetlistItem) => void;
    }>();
</script>

<li class="flex items-center justify-between gap-3 py-4">
    <div class="flex min-w-0 flex-grow items-start gap-4">
        <div class="flex-shrink-0 pt-1">
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
                ><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5" /></svg
                >
            </div>
        </div>

        <div class="min-w-0 flex-grow">
            {#if item.item_type === 'song'}
                <div class="flex items-center gap-3">
                    <span class="text-lg font-bold text-slate-400 dark:text-slate-500">{index + 1}.</span>
                    <p class="truncate font-semibold text-slate-800 dark:text-slate-100">
                        {item.title.String}
                    </p>
                </div>
                <div
                        class="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 pl-8 text-xs text-slate-500 dark:text-slate-400"
                >
                    {#if item.duration_seconds?.Valid}
                        <span>Durée: {formatItemDuration(item.duration_seconds.Int32)}</span>
                    {/if}
                    {#if item.tempo?.Valid}
                        <span class="hidden sm:inline">&bull;</span>
                        <span>Tempo: {item.tempo.Int32} BPM</span>
                    {/if}
                    {#if item.song_key?.Valid}
                        <span class="hidden sm:inline">&bull;</span>
                        <span>Tonalité: {item.song_key.String}</span>
                    {/if}
                    {#if item.links?.Valid}
                        <span class="hidden sm:inline">&bull;</span>
                        <a href={item.links.String} target="_blank" rel="noopener noreferrer" class="hover:underline"
                        >Lien</a
                        >
                    {/if}
                </div>
                {#if item.notes?.Valid && item.notes.String}
                    <p class="mt-2 whitespace-pre-wrap pl-8 text-xs italic text-slate-500 dark:text-slate-400">
                        {item.notes.String}
                    </p>
                {/if}
            {:else if item.item_type === 'interlude'}
                <div
                        class="w-full rounded-md border-l-4 border-teal-500 bg-teal-50 p-4 dark:border-teal-400 dark:bg-slate-700/50"
                >
                    <div class="flex items-center gap-3">
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
                        <p class="truncate font-semibold text-teal-900 dark:text-teal-200">
                            {item.title.String}
                        </p>
                    </div>
                    <div
                            class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 pl-8 text-xs text-teal-700 dark:text-teal-300"
                    >
                        {#if item.speaker?.Valid && item.speaker.String}
                            <span>Speaker: <span class="font-medium">{item.speaker.String}</span></span>
                            <span class="hidden sm:inline">&bull;</span>
                        {/if}
                        <span>Duration: {formatItemDuration(item.duration_seconds.Int32)}</span>
                    </div>
                    {#if item.notes?.Valid && item.notes.String}
                        <p class="mt-2 whitespace-pre-wrap pl-8 text-xs italic text-teal-800 dark:text-teal-200">
                            {item.notes.String}
                        </p>
                    {/if}
                </div>
            {/if}
        </div>
    </div>

    <div class="flex flex-shrink-0 items-center gap-2 pl-4">
        <button
                onclick={() => onEdit(item)}
                type="button"
                class="rounded-md p-2 text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-700 dark:hover:text-slate-200"
                aria-label="Edit item"
        ><svg
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
        ></button
        >
        <form method="POST" action="?/deleteItem" use:enhance>
            <input type="hidden" name="itemId" value={item.id} />
            <button
                    type="submit"
                    class="rounded-md p-2 text-slate-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-500/10 dark:hover:text-red-400"
                    aria-label="Remove item"
            ><svg
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
            ></button
            >
        </form>
    </div>
</li>