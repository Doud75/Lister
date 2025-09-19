<script lang="ts">
    import type { SetlistItem } from '$lib/types';
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import Textarea from '$lib/components/ui/Textarea.svelte';

    let { item, close } = $props<{
        item: SetlistItem;
        close: () => void;
    }>();

    let modalTitle = $state(item.title.String || '');
    let modalSpeaker = $state(item.speaker?.String || '');
    let modalDuration = $state(item.duration_seconds?.Int32?.toString() || '');
    let modalScript = $state(item.script?.String || '');
    let modalNotes = $state(item.notes?.String || '');

    const commonEnhance = () => {
        return async ({ update }: { update: any }) => {
            await update();
            close();
        };
    };
</script>

{#if item.item_type === 'song'}
    <h3 class="text-lg font-semibold text-slate-900 dark:text-white">
        Edit Note for {item.title.String}
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
        <input type="hidden" name="interludeId" value={item.interlude_id.Int32} />
        <Input label="Title" id="title" name="title" bind:value={modalTitle} required />
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Input label="Speaker" id="speaker" name="speaker" bind:value={modalSpeaker} />
            <Input
                    label="Duration (s)"
                    id="duration"
                    name="duration"
                    type="number"
                    bind:value={modalDuration}
            />
        </div>
        <Textarea label="Script" id="script" name="script" bind:value={modalScript} />
        <div class="flex justify-end gap-3 pt-2">
            <Button type="button" variant="secondary" onclick={close} autoWidth>Cancel</Button>
            <Button type="submit" autoWidth>Save Interlude</Button>
        </div>
    </form>
{/if}