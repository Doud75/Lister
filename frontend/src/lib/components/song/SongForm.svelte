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
        dur_min: string;
        dur_sec: string;
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
        const totalSec = s?.duration_seconds ?? null;
        return {
            title: s?.title ?? '',
            album_name: s?.album_name ?? '',
            song_key: s?.song_key ?? '',
            dur_min: totalSec != null ? Math.floor(totalSec / 60).toString() : '',
            dur_sec: totalSec != null ? (totalSec % 60).toString() : '',
            tempo: s?.tempo?.toString() ?? '',
            lyrics: s?.lyrics ?? '',
            links: s?.links ?? ''
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
    <Input label="Titre de la chanson" id="title" name="title" bind:value={formData.title} required />

    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <Input
                label="Nom de l'album (optionnel)"
                id="album_name"
                name="album_name"
                placeholder="ex : Abbey Road"
                bind:value={formData.album_name}
        />
        <Input
                label="Tonalité (optionnelle)"
                id="song_key"
                name="song_key"
                placeholder="ex : Lam"
                bind:value={formData.song_key}
        />
    </div>

    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
            <label class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">
                Durée (optionnelle)
            </label>
            <div class="mt-2 flex items-center gap-2">
                <input
                        type="number"
                        id="dur_min"
                        name="dur_min"
                        min="0"
                        placeholder="0"
                        bind:value={formData.dur_min}
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
                        bind:value={formData.dur_sec}
                        class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                />
                <span class="text-slate-500 dark:text-slate-400 shrink-0">sec</span>
            </div>
        </div>
        <Input label="Tempo (BPM)" id="tempo" name="tempo" type="number" bind:value={formData.tempo} />
    </div>

    <Input
            label="Lien (optionnel)"
            id="links"
            name="links"
            type="text"
            placeholder="https://youtube.com/..."
            bind:value={formData.links}
    />

    <div>
        <div class="mt-2">
            <Textarea label="Paroles (optionnelles)" id="lyrics" name="lyrics" rows={8} bind:value={formData.lyrics} />
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
            Annuler
        </a>
        <Button isLoading={$navigating?.type === 'form'} autoWidth>
            {#if isEditing}Sauvegarder{:else}Créer la chanson{/if}
        </Button>
    </div>
</form>
