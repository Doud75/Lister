<script lang="ts">
    import { enhance } from '$app/forms';
    import { page } from '$app/stores';
    import Modal from '$lib/components/ui/Modal.svelte';
    import type { ActionData, PageData } from './$types';

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let isModalOpen = $state(false);
    let bandName = $state('');
    let isCreating = $state(false);

    let isLeaveModalOpen = $state(false);
    let leaveBandId = $state('');
    let leaveBandName = $state('');
    let isLeaving = $state(false);

    function openModal() {
        bandName = '';
        isModalOpen = true;
    }

    function openLeaveModal(id: number, name: string, e: MouseEvent) {
        e.preventDefault();
        e.stopPropagation();
        leaveBandId = id.toString();
        leaveBandName = name;
        isLeaveModalOpen = true;
    }

    const leftGroupName = $derived($page.url.searchParams.get('left'));
</script>

<svelte:head>
    <title>Mes Groupes — Lister</title>
</svelte:head>

<div class="container mx-auto max-w-4xl px-4 py-8 sm:px-6">
    <div class="mb-8 flex flex-wrap items-center justify-between gap-4">
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">Mes Groupes</h1>
            <p class="mt-1 text-slate-500 dark:text-slate-400">
                {data.bands.length === 0
                    ? "Vous ne faites partie d'aucun groupe pour le moment."
                    : `Vous faites partie de ${data.bands.length} groupe${data.bands.length > 1 ? 's' : ''}.`}
            </p>
        </div>
        <button
            onclick={openModal}
            class="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
        >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5">
                <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"/>
            </svg>
            Créer un groupe
        </button>
    </div>

    {#if leftGroupName}
        <div class="mb-6 rounded-lg bg-teal-50 px-4 py-3 text-sm text-teal-700 dark:bg-teal-900/30 dark:text-teal-300">
            Vous avez quitté le groupe <strong>{leftGroupName}</strong>.
        </div>
    {/if}

    {#if form?.leaveError}
        <div class="mb-6 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-400">
            {form.leaveError}
        </div>
    {/if}

    {#if data.bands.length === 0}
        <div class="rounded-xl border-2 border-dashed border-slate-300 p-16 text-center dark:border-slate-700">
            <svg class="mx-auto h-12 w-12 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z" />
            </svg>
            <h3 class="mt-4 text-base font-semibold text-slate-900 dark:text-white">Bienvenue !</h3>
            <p class="mt-2 text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto">
                Pour commencer, créez votre premier groupe ou demandez un lien d'invitation à votre administrateur.
            </p>
        </div>
    {:else}
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {#each data.bands as band (band.id)}
                <div class="group relative rounded-xl border bg-white shadow-sm transition-all duration-200 hover:-translate-y-1 hover:shadow-md dark:bg-slate-800
                    {data.activeBandId === band.id.toString()
                        ? 'border-indigo-400 ring-2 ring-indigo-400 dark:border-indigo-500 dark:ring-indigo-500'
                        : 'border-slate-200 dark:border-slate-700 dark:hover:border-slate-600'}">
                    <div class="flex items-center">
                        <form method="POST" action="?/switchBand" use:enhance class="flex-1 min-w-0">
                            <input type="hidden" name="bandId" value={band.id} />
                            <button
                                type="submit"
                                class="w-full p-5 text-left"
                            >
                                <div class="flex items-start justify-between pr-2">
                                    <div class="min-w-0 flex-1">
                                        <h2 class="truncate text-base font-semibold text-slate-900 dark:text-white">{band.name}</h2>
                                        <span class="mt-1.5 inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium
                                            {band.role === 'admin'
                                                ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300'
                                                : 'bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-300'}">
                                            {band.role === 'admin' ? 'Admin' : 'Membre'}
                                        </span>
                                    </div>
                                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                                        class="h-5 w-5 flex-shrink-0 text-slate-400 transition-transform group-hover:translate-x-1">
                                        <path fill-rule="evenodd" d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                                {#if data.activeBandId === band.id.toString()}
                                    <p class="mt-3 text-xs text-indigo-600 dark:text-indigo-400">✓ Groupe actif</p>
                                {/if}
                            </button>
                        </form>

                        <div class="flex-shrink-0 flex items-center pr-3">
                            <button
                                type="button"
                                onclick={(e) => openLeaveModal(band.id, band.name, e)}
                                class="rounded-md p-2 text-slate-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-500/10 dark:hover:text-red-400"
                                aria-label="Quitter {band.name}"
                                title="Quitter ce groupe"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<Modal isOpen={isModalOpen} onClose={() => (isModalOpen = false)}>
    <div class="space-y-4">
        <h2 class="text-lg font-semibold text-slate-900 dark:text-white">Créer un nouveau groupe</h2>

        {#if form?.error}
            <p class="rounded-lg bg-red-100 px-4 py-2 text-sm text-red-700 dark:bg-red-900/30 dark:text-red-400">
                {form.error}
            </p>
        {/if}

        <form
            method="POST"
            action="?/createBand"
            use:enhance={() => {
                isCreating = true;
                return async ({ update }) => {
                    isCreating = false;
                    isModalOpen = false;
                    await update();
                };
            }}
            class="space-y-4"
        >
            <div>
                <label for="band-name" class="block text-sm font-medium text-slate-700 dark:text-slate-300">
                    Nom du groupe
                </label>
                <input
                    id="band-name"
                    name="name"
                    type="text"
                    bind:value={bandName}
                    placeholder="Ex : Les Rockers du Dimanche"
                    required
                    maxlength="100"
                    class="mt-1 block w-full rounded-lg border border-slate-300 bg-white px-4 py-2.5 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 dark:border-slate-600 dark:bg-slate-700 dark:text-white dark:placeholder-slate-400"
                />
            </div>

            <div class="flex justify-end gap-3 pt-2">
                <button
                    type="button"
                    onclick={() => (isModalOpen = false)}
                    class="rounded-lg border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-slate-700"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    disabled={isCreating || !bandName.trim()}
                    class="rounded-lg bg-indigo-600 px-5 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500 disabled:opacity-50"
                >
                    {isCreating ? 'Création...' : 'Créer le groupe'}
                </button>
            </div>
        </form>
    </div>
</Modal>

<Modal isOpen={isLeaveModalOpen} onClose={() => (isLeaveModalOpen = false)}>
    <div class="space-y-4">
        <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5 text-red-600 dark:text-red-400">
                    <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495ZM10 5a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 10 5Zm0 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z" clip-rule="evenodd" />
                </svg>
            </div>
            <h2 class="text-lg font-semibold text-slate-900 dark:text-white">Quitter le groupe</h2>
        </div>

        <p class="text-sm text-slate-600 dark:text-slate-400">
            Êtes-vous sûr de vouloir quitter <strong>{leaveBandName}</strong> ?
            Vous n'aurez plus accès à ce groupe ni à ses données.
        </p>

        <form
            method="POST"
            action="?/leaveBand"
            use:enhance={() => {
                isLeaving = true;
                return async ({ update }) => {
                    isLeaving = false;
                    isLeaveModalOpen = false;
                    await update();
                };
            }}
        >
            <input type="hidden" name="bandId" value={leaveBandId} />
            <input type="hidden" name="bandName" value={leaveBandName} />
            <div class="flex justify-end gap-3 pt-2">
                <button
                    type="button"
                    onclick={() => (isLeaveModalOpen = false)}
                    class="rounded-lg border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-slate-700"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    disabled={isLeaving}
                    class="rounded-lg bg-red-600 px-5 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-red-500 disabled:opacity-50"
                >
                    {isLeaving ? 'Départ...' : 'Quitter le groupe'}
                </button>
            </div>
        </form>
    </div>
</Modal>
