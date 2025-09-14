<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    let { data, form } = $props();
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Edit Song: <span class="text-indigo-500">{data.song.title}</span>
    </h1>

    <form method="POST" use:enhance class="mt-8 space-y-6 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <Input label="Song Title" id="title" name="title" required value={data.song.title} />
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
            <Input label="Duration (seconds)" id="duration_seconds" name="duration_seconds" type="number" value={data.song.duration_seconds?.Int32}/>
            <Input label="Tempo (BPM)" id="tempo" name="tempo" type="number" value={data.song.tempo?.Int32}/>
            <Input label="Key" id="song_key" name="song_key" placeholder="e.g., Am" value={data.song.song_key?.String}/>
        </div>

        {#if form?.error}
            <p class="text-sm text-red-500">{form.error}</p>
        {/if}

        <div class="flex items-center justify-between gap-4 pt-4">
            <a href="/song" class="rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm dark:bg-slate-700 dark:text-slate-200">Cancel</a>
            <Button isLoading={$navigating?.type === 'form'} autoWidth>Save Changes</Button>
        </div>
    </form>
</div>