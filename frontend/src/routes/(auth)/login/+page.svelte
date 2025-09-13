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
            Log in to your account
        </h2>
    </div>

    <form method="POST" use:enhance class="space-y-6">
        <Input label="Username" id="username" name="username" required />
        <Input label="Password" id="password" name="password" type="password" required />

        {#if form?.error}
            <p class="text-center text-sm text-red-500">{form.error}</p>
        {/if}

        <Button isLoading={$navigating?.type === 'form'}>
            {#if $navigating?.type === 'form'}
                Logging in...
            {:else}
                Log In
            {/if}
        </Button>
    </form>

    <p class="mt-8 text-center text-sm text-slate-500 dark:text-slate-400">
        No account? <a href="/signup" class="font-semibold leading-6 text-indigo-500 hover:text-indigo-400">Create a band</a> or <a href="/join" class="font-semibold leading-6 text-indigo-500 hover:text-indigo-400">join one</a>.
    </p>
</div>