<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    let { data, form } = $props<{ data: any, form: any }>();
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Edit Setlist
        </h1>
        <p class="mt-1 text-lg text-slate-600 dark:text-slate-400">
            Update the details for "<span class="font-semibold text-indigo-500">{data.setlist.name}</span>".
        </p>
    </header>

    <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <form class="space-y-6" method="POST" use:enhance>
            <Input label="Setlist Name" id="name" name="name" bind:value={data.setlist.name} required />

            <div>
                <label for="color" class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">
                    Label Color
                </label>
                <div class="mt-2 flex items-center gap-4">
                    <input
                            type="color"
                            id="color"
                            name="color"
                            bind:value={data.setlist.color}
                            class="h-10 w-16 cursor-pointer rounded-md border-0 bg-transparent p-0"
                    />
                </div>
            </div>

            {#if form?.error}
                <p class="text-sm text-red-500">{form.error}</p>
            {/if}

            <div class="flex items-center justify-between gap-4 pt-4">
                <a href="/setlist/{data.setlist.id}" class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                    Cancel
                </a>
                <Button isLoading={$navigating?.type === 'form'} autoWidth>
                    Save Changes
                </Button>
            </div>
        </form>
    </div>
</div>