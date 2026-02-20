<script lang="ts">
    import type { ActionData, PageData } from './$types';
    import SongForm from '$lib/components/song/SongForm.svelte';
    import { page } from '$app/stores';

    let { data, form }: { data: PageData; form: ActionData } = $props();

    const backHref = $derived($page.url.searchParams.get('from') ?? '/song');
    const fromParam = $derived($page.url.searchParams.get('from') ?? '');
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Modifier: <span class="text-indigo-500">{data.song.title}</span>
    </h1>

    <SongForm song={data.song} {form} cancelHref={backHref} isEditing={true} hiddenFrom={fromParam} />
</div>