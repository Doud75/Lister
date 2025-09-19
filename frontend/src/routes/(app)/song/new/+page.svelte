<script lang="ts">
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import type { ActionData } from './$types';
    import SongForm from '$lib/components/song/SongForm.svelte';

    let { form }: { form: ActionData } = $props();

    const cancelHref = $page.url.searchParams.get('redirectTo') || '/song';

    $effect(() => {
        if (form?.success && form?.redirectTo) {
            goto(form.redirectTo);
        }
    });
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Add a New Song to Your Library
    </h1>
    <p class="mt-1 text-lg text-slate-600 dark:text-slate-400">
        This will add the song to your global library.
    </p>

    <SongForm {form} {cancelHref} />
</div>