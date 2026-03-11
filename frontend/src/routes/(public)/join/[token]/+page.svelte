<script lang="ts">
    import { enhance } from '$app/forms';
    import type { ActionData, PageData } from './$types';

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let isJoining = $state(false);

    const currentToken = $derived((data as { token?: string }).token ?? '');
    const redirectTo = $derived(typeof window !== 'undefined' ? window.location.pathname : `/join/${currentToken}`);
    const loginUrl = $derived(`/login?redirectTo=${encodeURIComponent(redirectTo)}`);
    const signupUrl = $derived(`/signup?redirectTo=${encodeURIComponent(redirectTo)}`);
</script>

<svelte:head>
    <title>
        {data.valid ? `Rejoindre ${data.band_name} — Setlist` : 'Invitation invalide — Setlist'}
    </title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center px-4 py-16">
    <div class="w-full max-w-md">

        {#if !data.valid}
            <div class="rounded-2xl bg-white p-8 shadow-xl dark:bg-slate-800 text-center">
                <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-red-500" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" />
                    </svg>
                </div>
                <h1 class="text-xl font-bold text-slate-800 dark:text-white">Invitation invalide</h1>
                <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
                    {data.error ?? "Ce lien d'invitation est invalide ou a expiré."}
                </p>
                <a href="/" class="mt-6 inline-block text-sm font-medium text-indigo-500 hover:text-indigo-400">
                    Retour à l'accueil →
                </a>
            </div>

        {:else}
            <div class="rounded-2xl bg-white p-8 shadow-xl dark:bg-slate-800">

                <div class="mb-8 text-center">
                    <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-indigo-100 dark:bg-indigo-900/30">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m12 0a5.971 5.971 0 0 0-.941-3.197m0 0A5.995 5.995 0 0 0 12 12.75a5.995 5.995 0 0 0-5.058 2.772m0 0a3 3 0 0 0-4.681 2.72 8.986 8.986 0 0 0 3.74.477m.94-3.197a5.971 5.971 0 0 0-.94 3.197M15 6.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Zm-13.5 0a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Z" />
                        </svg>
                    </div>
                    <p class="text-sm text-slate-500 dark:text-slate-400">Vous avez été invité à rejoindre</p>
                    <h1 class="mt-1 text-2xl font-bold text-slate-900 dark:text-white">{data.band_name}</h1>
                    <span class="mt-2 inline-flex items-center rounded-full bg-indigo-50 px-3 py-1 text-xs font-medium text-indigo-700 ring-1 ring-inset ring-indigo-600/20 dark:bg-indigo-500/10 dark:text-indigo-400 dark:ring-indigo-500/20">
                        Rôle : {data.role}
                    </span>
                </div>

                {#if form?.error}
                    <div class="mb-4 rounded-lg bg-red-50 p-3 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
                        {form.error}
                    </div>
                {/if}

                {#if !data.user}
                    <div class="space-y-3">
                        <p class="text-center text-sm text-slate-600 dark:text-slate-400">
                            Connectez-vous ou créez un compte pour rejoindre ce groupe.
                        </p>
                        <a
                            href={loginUrl}
                            class="flex w-full items-center justify-center rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >
                            Se connecter
                        </a>
                        <a
                            href={signupUrl}
                            class="flex w-full items-center justify-center rounded-lg border border-slate-300 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-50 dark:border-slate-600 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
                        >
                            Créer un compte
                        </a>
                    </div>

                {:else}
                    <div class="space-y-4">
                        <div class="flex items-center gap-3 rounded-lg bg-slate-50 px-4 py-3 dark:bg-slate-700/50">
                            <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-indigo-100 text-sm font-bold text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400">
                                {data.user.username[0].toUpperCase()}
                            </div>
                            <div>
                                <p class="text-sm font-medium text-slate-800 dark:text-slate-100">Connecté en tant que</p>
                                <p class="text-sm text-slate-600 dark:text-slate-400">{data.user.username}</p>
                            </div>
                        </div>

                        <form method="POST" use:enhance={() => {
                            isJoining = true;
                            return async ({ update }) => {
                                await update();
                                isJoining = false;
                            };
                        }}>
                            <button
                                type="submit"
                                disabled={isJoining}
                                class="flex w-full items-center justify-center rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                            >
                                {#if isJoining}
                                    <svg class="mr-2 h-4 w-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                                    </svg>
                                    Rejoindre...
                                {:else}
                                    Rejoindre le groupe
                                {/if}
                            </button>
                        </form>

                        <p class="text-center text-xs text-slate-500 dark:text-slate-400">
                            Ce n'est pas vous ?
                            <a href="/logout" class="font-medium text-indigo-500 hover:text-indigo-400">Se déconnecter</a>
                        </p>
                    </div>
                {/if}
            </div>
        {/if}
    </div>
</div>
