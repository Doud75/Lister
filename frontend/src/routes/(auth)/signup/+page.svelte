<script lang="ts">
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';
    import { enhance } from '$app/forms';
    import { navigating } from '$app/stores';

    let { form } = $props();

    let username = $state('');
    let password = $state('');
    let validationError = $state('');

    const usernameRegex = /^[a-zA-Z0-9_]{3,50}$/;
    const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$/;

    function validateForm() {
        if (!usernameRegex.test(username)) {
            validationError = "Le nom d'utilisateur doit contenir entre 3 et 50 caractères alphanumériques ou underscore.";
            return false;
        }
        if (!passwordRegex.test(password)) {
            validationError = "Le mot de passe doit contenir au moins 8 caractères, une majuscule, un chiffre et un caractère spécial.";
            return false;
        }
        validationError = '';
        return true;
    }
</script>

<div class="space-y-6">
    <div>
        <h2 class="text-center text-2xl font-bold leading-9 tracking-tight text-slate-900 dark:text-white">
            Créer un nouveau groupe
        </h2>
        <p class="mt-2 text-center text-sm text-slate-600 dark:text-slate-400">
            Commencez par donner un nom à votre groupe et créez votre compte.
        </p>
    </div>

    <form
            method="POST"
            use:enhance={({ cancel }) => {
            if (!validateForm()) {
                cancel();
            }
            return async ({ update }) => {
                await update();
            };
        }}
            class="space-y-6"
    >
        <Input label="Nom du groupe" id="bandName" name="bandName" placeholder="The Rolling Scones" required />
        <Input
                label="Votre nom d'utilisateur"
                id="username"
                name="username"
                placeholder="votre_pseudo"
                required
                bind:value={username}
                oninput={() => validationError = ''}
        />
        <Input
                label="Mot de passe"
                id="password"
                name="password"
                type="password"
                required
                togglePasswordVisibility={true}
                bind:value={password}
                oninput={() => validationError = ''}
        />

        {#if validationError}
            <p class="text-center text-sm text-red-500">{validationError}</p>
        {/if}

        {#if form?.error}
            <p class="text-center text-sm text-red-500">{form.error}</p>
        {/if}

        <Button isLoading={$navigating?.type === 'form'}>
            {#if $navigating?.type === 'form'}
                Création en cours...
            {:else}
                Créer le groupe et s'inscrire
            {/if}
        </Button>
    </form>

    <p class="mt-8 text-center text-sm text-slate-500 dark:text-slate-400">
        Déjà un compte ? <a href="/login" class="font-semibold leading-6 text-indigo-500 hover:text-indigo-400">Login</a>.
    </p>
</div>