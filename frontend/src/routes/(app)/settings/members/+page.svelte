<script lang="ts">
    import { enhance } from '$app/forms';
    import { invalidateAll } from '$app/navigation';
    import type { ActionData, PageData } from './$types';
    import Input from '$lib/components/ui/Input.svelte';
    import Button from '$lib/components/ui/Button.svelte';

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let members = $derived(data.members);
    let isInviting = $state(false);
    let usernameInput = $state('');
    let passwordInput = $state('');

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
            invalidateAll();
        }
    });
</script>

<div class="container mx-auto max-w-4xl px-4 sm:px-6">
    <header class="mb-8">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
            Gérer les membres du groupe
        </h1>
        <p class="mt-1 text-lg text-slate-600 dark:text-slate-400">
            Invitez de nouveaux membres et gérez les accès.
        </p>
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

                                <!-- enhance simple, sans callback. Le $effect s'occupe de tout. -->
                                <form method="POST" action="?/removeMember" use:enhance>
                                    <input type="hidden" name="userId" value={member.id} />
                                    <button
                                            type="submit"
                                            class="rounded-md p-2 text-sm font-medium text-slate-500 hover:bg-red-50 hover:text-red-600 disabled:cursor-not-allowed disabled:opacity-50 dark:text-slate-400 dark:hover:bg-red-500/10 dark:hover:text-red-400"
                                            aria-label="Supprimer {member.username}"
                                            disabled={member.role === 'admin' &&
											members.filter((m) => m.role === 'admin').length <= 1}
                                    >
                                        Supprimer
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
            <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
                <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Inviter un membre</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                    Créez un compte pour un nouveau membre. Il pourra changer son mot de passe plus tard.
                </p>
                <form
                        method="POST"
                        action="?/inviteMember"
                        class="mt-6 space-y-6"
                        use:enhance={() => {
						isInviting = true;
						return async ({ update }) => {
							await update(); // On a juste besoin de `update` pour mettre à jour la prop `form`
							isInviting = false;
						};
					}}
                >
                    <fieldset disabled={isInviting} class="contents">
                        <Input
                                label="Nom d'utilisateur"
                                id="username"
                                name="username"
                                bind:value={usernameInput}
                                required
                        />
                        <Input
                                label="Mot de passe temporaire"
                                id="password"
                                name="password"
                                type="password"
                                bind:value={passwordInput}
                                required
                        />
                    </fieldset>

                    {#if form?.error}
                        <p class="text-sm text-red-500">{form.error}</p>
                    {/if}

                    <div class="pt-2">
                        <Button isLoading={isInviting} autoWidth> Inviter </Button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>