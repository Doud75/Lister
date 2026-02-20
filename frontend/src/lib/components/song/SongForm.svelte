<script lang="ts">
    import { untrack } from 'svelte';
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import Textarea from '$lib/components/ui/Textarea.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';
    import type { ActionData } from '../../../routes/(app)/song/new/$types';
    import type { Song } from '$lib/types';

    type FormSong = {
        title: string;
        album_name: string;
        song_key: string;
        duration_seconds: string;
        tempo: string;
        lyrics: string;
        links: string;
    };

    let {
        song = {} as Partial<Song>,
        form,
        cancelHref,
        isEditing = false,
        hiddenFrom = ''
    }: {
        song?: Partial<Song>;
        form: ActionData;
        cancelHref: string;
        isEditing?: boolean;
        hiddenFrom?: string;
    } = $props();


    function getFormValues(s: Partial<Song>): FormSong {
        return {
            title: s?.title ?? '',
            album_name: s?.album_name?.String ?? '',
            song_key: s?.song_key?.String ?? '',
            duration_seconds: s?.duration_seconds?.Int32?.toString() ?? '',
            tempo: s?.tempo?.Int32?.toString() ?? '',
            lyrics: s?.lyrics?.String ?? '',
            links: s?.links?.String ?? ''
        };
    }

    let formData = $state<FormSong>(untrack(() => getFormValues(song)));

    $effect(() => {
        Object.assign(formData, getFormValues(song));
    });
</script>

<form method="POST" use:enhance class="mt-8 space-y-6 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
    {#if hiddenFrom}
        <input type="hidden" name="from" value={hiddenFrom} />
    {/if}
    <Input label="Song Title" id="title" name="title" bind:value={formData.title} required />

    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <Input
                label="Album Name (optional)"
                id="album_name"
                name="album_name"
                placeholder="e.g., Abbey Road"
                bind:value={formData.album_name}
        />
        <Input
                label="Key (optional)"
                id="song_key"
                name="song_key"
                placeholder="e.g., Am"
                bind:value={formData.song_key}
        />
    </div>

    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <Input
                label="Duration (seconds)"
                id="duration_seconds"
                name="duration_seconds"
                type="number"
                bind:value={formData.duration_seconds}
        />
        <Input label="Tempo (BPM)" id="tempo" name="tempo" type="number" bind:value={formData.tempo} />
    </div>

    <Input
            label="Link (optional)"
            id="links"
            name="links"
            type="text"
            placeholder="https://youtube.com/..."
            bind:value={formData.links}
    />

    <div>
        <div class="mt-2">
            <Textarea label="Lyrics (optional)" id="lyrics" name="lyrics" rows={8} bind:value={formData.lyrics} />
        </div>
    </div>

    {#if form?.error}
        <p class="text-sm text-red-500">{form.error}</p>
    {/if}

    <div class="flex items-center justify-between gap-4 pt-4">
        <a
                href={cancelHref}
                class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
        >
            {#if isEditing}Annuler{:else}Cancel{/if}
        </a>
        <Button isLoading={$navigating?.type === 'form'} autoWidth>
            {#if isEditing}Sauvegarder{:else}Create Song{/if}
        </Button>
    </div>
</form>