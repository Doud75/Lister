<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    let { form } = $props();
</script>

<div class="space-y-6">
    <div>
        <h2 class="text-center text-2xl font-bold leading-9 tracking-tight text-slate-900 dark:text-white">
            Join an Existing Band
        </h2>
        <p class="mt-2 text-center text-sm text-slate-600 dark:text-slate-400">
            Enter the exact band name and create your personal account.
        </p>
    </div>

    <form method="POST" use:enhance class="space-y-6">
        <Input label="Exact Band Name" id="bandName" name="bandName" required />
        <Input label="Your Username" id="username" name="username" required />
        <Input label="Password" id="password" name="password" type="password" required />

        {#if form?.error}
            <p class="text-center text-sm text-red-500">{form.error}</p>
        {/if}

        <Button isLoading={$navigating?.type === 'form'} variant="secondary">
            {#if $navigating?.type === 'form'}
                Joining...
            {:else}
                Join Band
            {/if}
        </Button>
    </form>

    <p class="mt-8 text-center text-sm text-slate-500 dark:text-slate-400">
        Need to create a new band? <a href="/signup" class="font-semibold leading-6 text-indigo-500 hover:text-indigo-400">Create one here</a>.
    </p>
</div>