<script lang="ts">
    import type { ActionData, PageData } from './$types';
    import SongForm from '$lib/components/song/SongForm.svelte';

    let { data, form }: { data: PageData; form: ActionData } = $props();
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Modifier: <span class="text-indigo-500">{data.song.title}</span>
    </h1>

    <SongForm song={data.song} {form} cancelHref="/song" isEditing={true} />
</div>