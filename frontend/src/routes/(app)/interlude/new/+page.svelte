<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';
    import { goto } from '$app/navigation';

    type FormState = {
        error?: string;
        success?: boolean;
        redirectTo?: string;
    } | undefined;

    let { form } = $props<{ form: FormState }>();

    $effect(() => {
        if (form?.success && form?.redirectTo) {
            goto(form.redirectTo);
        }
    });
</script>

<div class="container mx-auto max-w-2xl px-4 sm:px-6">
    <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Add a New Interlude
    </h1>
    <p class="mt-1 text-lg text-slate-600 dark:text-slate-400">
        This will add the interlude to your global library.
    </p>

    <form
            method="POST"
            use:enhance
            class="mt-8 space-y-6 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800"
    >
        <Input label="Interlude Title" id="title" name="title" required />
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <Input label="Speaker (optional)" id="speaker" name="speaker" placeholder="e.g., Lead Singer" />
            <Input
                    label="Duration (seconds)"
                    id="duration_seconds"
                    name="duration_seconds"
                    type="number"
            />
        </div>

        <div>
            <label
                    for="script"
                    class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200"
            >Script / Notes (optional)</label
            >
            <div class="mt-2">
				<textarea
                        id="script"
                        name="script"
                        rows="4"
                        class="block w-full rounded-md border-0 bg-white/5 py-2 px-3 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 dark:bg-white/5 dark:text-white dark:ring-white/10 dark:focus:ring-indigo-500"
                ></textarea>
            </div>
        </div>

        {#if form?.error}
            <p class="text-sm text-red-500">{form.error}</p>
        {/if}

        <div class="flex items-center justify-between gap-4 pt-4">
            <a
                    href="/"
                    class="flex w-auto justify-center rounded-md bg-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
            >
                Cancel
            </a>
            <Button isLoading={$navigating?.type === 'form'} autoWidth>Create Interlude</Button>
        </div>
    </form>
</div>