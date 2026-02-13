<script lang="ts">
    import Input from '$lib/components/ui/Input.svelte';
    import Button from '$lib/components/ui/Button.svelte';
    import { enhance } from '$app/forms';
    import type { ActionData } from './$types';

    let { form }: { form: ActionData } = $props();

    let isSubmitting = $state(false);
    let newPassword = $state('');
    let confirmPassword = $state('');
    let validationError = $state('');

    const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$/;

    function validateForm() {
        if (!passwordRegex.test(newPassword)) {
            validationError = "Le nouveau mot de passe doit contenir au moins 8 caractères, une majuscule, un chiffre et un caractère spécial.";
            return false;
        }
        if (newPassword !== confirmPassword) {
            validationError = "Les mots de passe ne correspondent pas.";
            return false;
        }
        validationError = '';
        return true;
    }
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Mon Compte
        </h1>
        <div class="mt-2 flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400">
            <a href="/" class="hover:underline">&larr; Back to Home</a>
        </div>
    </header>

    <div class="rounded-xl mx-auto max-w-3xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Changer de mot de passe</h2>
        <form
                method="POST"
                class="mt-6 space-y-6"
                use:enhance={({ cancel }) => {
                if (!validateForm()) {
                    cancel();
                    return;
                }
                isSubmitting = true;
                return async ({ update }) => {
                    await update();
                    isSubmitting = false;
                };
            }}
        >
            <fieldset disabled={isSubmitting} class="contents">
                <Input
                        label="Ancien mot de passe"
                        id="current_password"
                        name="current_password"
                        type="password"
                        required
                        togglePasswordVisibility={true}
                />
                <Input
                        label="Nouveau mot de passe"
                        id="new_password"
                        name="new_password"
                        type="password"
                        required
                        togglePasswordVisibility={true}
                        bind:value={newPassword}
                        oninput={() => validationError = ''}
                />
                <Input
                        label="Confirmer le nouveau mot de passe"
                        id="confirm_password"
                        name="confirm_password"
                        type="password"
                        required
                        togglePasswordVisibility={true}
                        bind:value={confirmPassword}
                        oninput={() => validationError = ''}
                />
            </fieldset>

            {#if validationError}
                <p class="text-sm text-red-500">{validationError}</p>
            {/if}

            {#if form?.error}
                <p class="text-sm text-red-500">{form.error}</p>
            {/if}

            <div class="flex justify-end pt-4">
                <Button isLoading={isSubmitting} autoWidth>
                    Sauvegarder
                </Button>
            </div>
        </form>
    </div>
</div>