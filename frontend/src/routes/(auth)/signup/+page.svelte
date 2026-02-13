<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    type ActionData = {
        error?: string;
        errors?: Record<string, string>;
        data?: {
            bandName?: string;
            username?: string;
        };
    } | null;

    let { form }: { form: ActionData } = $props();

    let bandName = $state('');
    let username = $state('');
    let password = $state('');

    $effect(() => {
        if (form?.data?.bandName) bandName = form.data.bandName;
        if (form?.data?.username) username = form.data.username;
    });

</script>

<div class="space-y-6">
    <div>
        <h2 class="text-center text-2xl font-bold leading-9 tracking-tight text-slate-900 dark:text-white">
            Créer un nouveau groupe
        </h2>
        <p class="mt-2 text-center text-sm text-slate-600 dark:text-slate-400">
            Commencez par donner un nom à votre groupe et créer votre compte.
        </p>
    </div>

    <form method="POST" use:enhance class="space-y-6">
        <Input
                label="Nom du groupe"
                id="bandName"
                name="bandName"
                placeholder="Les Rolling Scones"
                required
                bind:value={bandName}
        />

        <div class="space-y-1">
            <Input
                    label="Votre nom d'utilisateur"
                    id="username"
                    name="username"
                    placeholder="votre_pseudo"
                    required
                    bind:value={username}
            />
            {#if form?.errors?.username}
                <p class="text-sm text-red-500 font-medium">{form.errors.username}</p>
            {/if}
            <p class="text-xs text-slate-500">
                3-50 caractères, alphanumérique & underscore uniquement.
            </p>
        </div>

        <div class="space-y-1">
            <Input
                    label="Mot de passe"
                    id="password"
                    name="password"
                    type="password"
                    required
                    togglePasswordVisibility={true}
                    bind:value={password}
            />
            {#if form?.errors?.password}
                <p class="text-sm text-red-500 font-medium">{form.errors.password}</p>
            {/if}
            
            <ul class="text-xs text-slate-500 list-disc ml-4 space-y-1 mt-2">
                <li class={password.length >= 8 ? 'text-teal-600 dark:text-teal-400' : ''}>Minimum 8 caractères</li>
                <li class={/[A-Z]/.test(password) ? 'text-teal-600 dark:text-teal-400' : ''}>Au moins 1 majuscule</li>
                <li class={/[0-9]/.test(password) ? 'text-teal-600 dark:text-teal-400' : ''}>Au moins 1 chiffre</li>
                <li class={/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password) ? 'text-teal-600 dark:text-teal-400' : ''}>Au moins 1 caractère spécial</li>
            </ul>
        </div>

        {#if form?.error}
            <p class="text-center text-sm text-red-500 font-bold bg-red-50 dark:bg-red-900/10 p-2 rounded">{form.error}</p>
        {/if}

        <Button isLoading={$navigating?.type === 'form'}>
            {#if $navigating?.type === 'form'}
                Création...
            {:else}
                Créer le groupe et le compte
            {/if}
        </Button>
    </form>

    <p class="mt-8 text-center text-sm text-slate-500 dark:text-slate-400">
        Déjà un compte ? <a href="/login" class="font-semibold leading-6 text-indigo-500 hover:text-indigo-400">Se connecter</a>.
    </p>
</div>