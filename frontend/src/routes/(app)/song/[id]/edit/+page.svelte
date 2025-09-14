<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';
    import type { ActionData, PageData } from './$types';

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let title = $state(data.song.title ?? '');
    let album_name = $state(data.song.album_name?.String ?? '');
    let song_key = $state(data.song.song_key?.String ?? '');
    let duration_seconds = $state(data.song.duration_seconds?.Int32?.toString() ?? '');
    let tempo = $state(data.song.tempo?.Int32?.toString() ?? '');
    let lyrics = $state(data.song.lyrics?.String ?? '');
    let links = $state(data.song.links?.String ?? '');
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Modifier: <span class="text-indigo-500">{data.song.title}</span>
    </h1>

    <form
            method="POST"
            use:enhance
            class="mt-8 space-y-6 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800"
    >
        <Input label="Song Title" id="title" name="title" bind:value={title} required />

        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input
                    label="Album Name (optional)"
                    id="album_name"
                    name="album_name"
                    placeholder="e.g., Abbey Road"
                    bind:value={album_name}
            />
            <Input
                    label="Key (optional)"
                    id="song_key"
                    name="song_key"
                    placeholder="e.g., Am"
                    bind:value={song_key}
            />
        </div>

        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input
                    label="Duration (seconds)"
                    id="duration_seconds"
                    name="duration_seconds"
                    type="number"
                    bind:value={duration_seconds}
            />
            <Input label="Tempo (BPM)" id="tempo" name="tempo" type="number" bind:value={tempo} />
        </div>

        <Input label="Link (optional)" id="links" name="links" type="text" placeholder="https://youtube.com/..." bind:value={links} />

        <div>
            <label for="lyrics" class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200"
            >Lyrics (optional)</label
            >
            <div class="mt-2">
				<textarea
                        id="lyrics"
                        name="lyrics"
                        rows="8"
                        bind:value={lyrics}
                        class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                ></textarea>
            </div>
        </div>

        {#if form?.error}
            <p class="text-sm text-red-500">{form.error}</p>
        {/if}

        <div class="flex items-center justify-between gap-4 pt-4">
            <a
                    href="/song"
                    class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
            >
                Annuler
            </a>
            <Button isLoading={$navigating?.type === 'form'} autoWidth>Sauvegarder</Button>
        </div>
    </form>
</div>