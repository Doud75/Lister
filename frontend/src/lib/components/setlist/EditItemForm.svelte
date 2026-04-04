<script lang="ts">
    import type { SetlistItem } from '$lib/types';
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import Textarea from '$lib/components/ui/Textarea.svelte';
    import { untrack } from 'svelte';

    let { item, close } = $props<{
        item: SetlistItem;
        close: () => void;
    }>();

    let modalTitle = $state(untrack(() => item.title ?? ''));
    let modalSpeaker = $state(untrack(() => item.speaker ?? ''));
    let modalDurMin = $state(untrack(() => item.duration_seconds != null ? Math.floor(item.duration_seconds / 60).toString() : ''));
    let modalDurSec = $state(untrack(() => item.duration_seconds != null ? (item.duration_seconds % 60).toString() : ''));
    let modalNotes = $state(untrack(() => item.notes ?? ''));

    const commonEnhance = () => {
        /* eslint-disable  @typescript-eslint/no-explicit-any */
        return async ({ update }: { update: any }) => {
            await update();
            close();
        };
    };
</script>

{#if item.item_type === 'song'}
    <h3 class="text-lg font-semibold text-slate-900 dark:text-white">
        Edit Note for {item.title ?? ''}
    </h3>
    <form method="POST" action="?/updateSongNotes" use:enhance={commonEnhance} class="mt-4 space-y-4">
        <input type="hidden" name="itemId" value={item.id} />
        <Textarea
                label="Notes"
                id="notes"
                name="notes"
                bind:value={modalNotes}
                placeholder="Add a comment..."
        />
        <div class="flex justify-end gap-3">
            <Button type="button" variant="secondary" onclick={close} autoWidth>Cancel</Button>
            <Button type="submit" autoWidth>Save Note</Button>
        </div>
    </form>
{/if}

{#if item.item_type === 'interlude'}
    <h3 class="text-lg font-semibold text-slate-900 dark:text-white">Edit Interlude</h3>
    <form
            method="POST"
            action="?/updateInterlude"
            use:enhance={commonEnhance}
            class="mt-4 space-y-4"
    >
        <input type="hidden" name="interludeId" value={item.interlude_id} />
        <input type="hidden" name="itemId" value={item.id} />
        <Input label="Title" id="title" name="title" bind:value={modalTitle} required />
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Input label="Speaker" id="speaker" name="speaker" bind:value={modalSpeaker} />
            <div>
                <label class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">
                    Duration
                </label>
                <div class="mt-2 flex items-center gap-2">
                    <input
                            type="number"
                            id="dur_min"
                            name="dur_min"
                            min="0"
                            placeholder="0"
                            bind:value={modalDurMin}
                            class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                    />
                    <span class="text-slate-500 dark:text-slate-400 shrink-0">min</span>
                    <input
                            type="number"
                            id="dur_sec"
                            name="dur_sec"
                            min="0"
                            max="59"
                            placeholder="00"
                            bind:value={modalDurSec}
                            class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                    />
                    <span class="text-slate-500 dark:text-slate-400 shrink-0">sec</span>
                </div>
            </div>
        </div>
        <Textarea label="Script" id="script" name="script" bind:value={modalNotes} />
        <div class="flex justify-end gap-3 pt-2">
            <Button type="button" variant="secondary" onclick={close} autoWidth>Cancel</Button>
            <Button type="submit" autoWidth>Save Interlude</Button>
        </div>
    </form>
{/if}
