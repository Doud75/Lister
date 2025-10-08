<script lang="ts">
    import { enhance } from '$app/forms';
    import { invalidateAll } from '$app/navigation';
    import type { ActionData, PageData } from './$types';
    import Input from '$lib/components/ui/Input.svelte';
    import Button from '$lib/components/ui/Button.svelte';
    import { onMount } from 'svelte';

    type UserSearchResult = {
        id: number;
        username: string;
    };

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let members = $derived(data.members);
    let isInviting = $state(false);
    let usernameInput = $state('');
    let passwordInput = $state('');
    let showPassword = $state(false);

    let searchResults = $state<UserSearchResult[]>([]);
    let isSearching = $state(false);
    let showResults = $state(false);
    let debounceTimer: number;

    async function searchUsers(query: string) {
        if (query.length < 3) {
            searchResults = [];
            showResults = false;
            return;
        }
        isSearching = true;
        try {
            const res = await fetch(`/api/user/search?q=${encodeURIComponent(query)}`);
            if (res.ok) {
                searchResults = await res.json();
                showResults = true;
            }
        } catch (e) {
            console.error('Search failed', e);
        } finally {
            isSearching = false;
        }
    }

    function handleUsernameInput(event: Event) {
        const input = event.target as HTMLInputElement;
        usernameInput = input.value;
        clearTimeout(debounceTimer);
        debounceTimer = window.setTimeout(() => {
            searchUsers(usernameInput);
        }, 300);
    }

    function selectUser(username: string) {
        usernameInput = username;
        showResults = false;
        searchResults = [];
    }

    let containerEl: HTMLElement;
    onMount(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerEl && !containerEl.contains(event.target as Node)) {
                showResults = false;
            }
        };
        document.addEventListener('click', handleClickOutside);
        return () => document.removeEventListener('click', handleClickOutside);
    });

    $effect(() => {
        if (form?.removeSuccess) {
            const index = members.findIndex((m) => m.id === form.removedUserId);
            if (index !== -1) {
                members.splice(index, 1);
            }
        }

        if (form?.inviteSuccess) {
            usernameInput = '';
            passwordInput = '';
            showPassword = false;
            invalidateAll();
        }

        if (form?.error && form.error.includes('password is required')) {
            showPassword = true;
            usernameInput = form.username || '';
        } else if (form?.error) {
            showPassword = false;
        }
    });
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Gérer les membres du groupe
        </h1>
        <div class="mt-2 flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400">
            <a href="/" class="hover:underline">&larr; Back to Home</a>
        </div>
    </header>

    <div class="grid grid-cols-1 gap-12 lg:grid-cols-3">
        <div class="lg:col-span-2">
            <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
                <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">
                    Membres actuels ({members.length})
                </h2>
                {#if members.length > 0}
                    <ul class="mt-4 divide-y divide-slate-200 dark:divide-slate-700">
                        {#each members as member (member.id)}
                            <li class="flex items-center justify-between gap-3 py-4">
                                <div>
                                    <p class="font-semibold text-slate-800 dark:text-slate-100">
                                        {member.username}
                                    </p>
                                    <span
                                            class="mt-1 inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset {member.role ===
										'admin'
											? 'bg-blue-50 text-blue-700 ring-blue-600/20 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20'
											: 'bg-slate-50 text-slate-600 ring-slate-500/20 dark:bg-slate-500/10 dark:text-slate-400 dark:ring-slate-500/20'}"
                                    >
										{member.role}
									</span>
                                </div>

                                <form method="POST" action="?/removeMember" use:enhance>
                                    <input type="hidden" name="userId" value={member.id} />
                                    <button
                                            type="submit"
                                            class="rounded-md p-2 text-slate-400 hover:bg-red-50 hover:text-red-600 disabled:cursor-not-allowed disabled:opacity-50 dark:text-slate-400 dark:hover:bg-red-500/10 dark:hover:text-red-400"
                                            aria-label="Supprimer {member.username}"
                                            disabled={member.role === 'admin' &&
											members.filter((m) => m.role === 'admin').length <= 1}
                                    >
                                        <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke-width="1.5"
                                                stroke="currentColor"
                                                class="h-5 w-5"
                                        ><path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.134-2.033-2.134H8.033C6.91 2.75 6 3.704 6 4.874v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                                        /></svg
                                        >
                                    </button>
                                </form>
                            </li>
                        {/each}
                    </ul>
                {:else}
                    <p class="mt-4 text-sm text-slate-500 dark:text-slate-400">Aucun membre dans ce groupe.</p>
                {/if}
            </div>
        </div>

        <div class="lg:col-span-1">
            <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800" bind:this={containerEl}>
                <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Ajouter un membre</h2>
                <form
                        method="POST"
                        action="?/inviteMember"
                        class="mt-6 space-y-6"
                        use:enhance={() => {
						isInviting = true;
						return async ({ update }) => {
                            showResults = false;
							await update();
							isInviting = false;
						};
					}}
                >
                    <fieldset disabled={isInviting} class="contents">
                        <div class="relative">
                            <Input
                                    label="Nom d'utilisateur"
                                    id="username"
                                    name="username"
                                    bind:value={usernameInput}
                                    oninput={handleUsernameInput}
                                    autocomplete="off"
                                    required
                            />
                            {#if showResults && searchResults.length > 0}
                                <ul class="absolute z-10 mt-1 w-full rounded-md border border-slate-300 bg-white py-1 text-sm shadow-lg dark:border-slate-600 dark:bg-slate-700">
                                    {#each searchResults as user (user.id)}
                                        <li>
                                            <button type="button" onclick={() => selectUser(user.username)} class="w-full px-4 py-2 text-left text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-600">
                                                {user.username}
                                            </button>
                                        </li>
                                    {/each}
                                </ul>
                            {:else if isSearching}
                                <div class="absolute z-10 mt-1 w-full rounded-md border border-slate-300 bg-white py-2 px-4 text-sm text-slate-500 shadow-lg dark:border-slate-600 dark:bg-slate-700 dark:text-slate-400">
                                    Recherche...
                                </div>
                            {/if}
                        </div>

                        {#if showPassword}
                            <Input
                                    label="Mot de passe temporaire"
                                    id="password"
                                    name="password"
                                    type="password"
                                    bind:value={passwordInput}
                                    required
                            />
                        {/if}
                    </fieldset>

                    {#if form?.error}
                        <p class="text-sm text-red-500">
                            {#if form.error.includes('password is required')}
                                Utilisateur non trouvé. Veuillez définir un mot de passe pour créer son compte.
                            {:else}
                                {form.error}
                            {/if}
                        </p>
                    {/if}

                    <div class="pt-2">
                        <Button isLoading={isInviting} autoWidth>
                            {#if showPassword}
                                Créer et Inviter
                            {:else}
                                Ajouter / Inviter
                            {/if}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>