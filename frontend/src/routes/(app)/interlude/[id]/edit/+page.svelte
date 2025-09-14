<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    let { data, form } = $props();
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Edit Interlude: <span class="text-indigo-500">{data.interlude.title}</span>
    </h1>

    <form method="POST" use:enhance class="mt-8 space-y-6 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <Input label="Interlude Title" id="title" name="title" required value={data.interlude.title} />
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input label="Speaker (optional)" id="speaker" name="speaker" placeholder="e.g., Lead Singer" value={data.interlude.speaker?.String}/>
            <Input label="Duration (seconds)" id="duration_seconds" name="duration_seconds" type="number" value={data.interlude.duration_seconds?.Int32}/>
        </div>
        <div>
            <label for="script" class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">Script / Notes (optional)</label>
            <div class="mt-2">
                <textarea id="script" name="script" rows="4" class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500">{data.interlude.script?.String ?? ''}</textarea>
            </div>
        </div>

        {#if form?.error}
            <p class="text-sm text-red-500">{form.error}</p>
        {/if}

        <div class="flex items-center justify-between gap-4 pt-4">
            <a href="/interlude" class="rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm dark:bg-slate-700 dark:text-slate-200">Cancel</a>
            <Button isLoading={$navigating?.type === 'form'} autoWidth>Save Changes</Button>
        </div>
    </form>
</div>