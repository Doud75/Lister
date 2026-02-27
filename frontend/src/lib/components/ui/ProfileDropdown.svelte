<script lang="ts">
    import type { PageData } from '../../../routes/(app)/$types';
    import { enhance } from '$app/forms';

    let { user }: { user: PageData['user'] } = $props();
    let isOpen = $state(false);

    function close() {
        isOpen = false;
    }
</script>

<div class="relative">
    <button
            onclick={() => (isOpen = !isOpen)}
            class="flex items-center justify-center h-10 w-10 rounded-full bg-slate-200 text-sm font-semibold uppercase text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
            aria-haspopup="true"
            aria-expanded={isOpen}
            aria-label="Ouvrir le menu du profil"
    >
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-6 w-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
        </svg>
    </button>

    {#if isOpen}
        <button
            type="button"
            onclick={close}
            class="fixed inset-0 z-10 h-full w-full cursor-default"
            tabindex="-1"
            aria-hidden="true"
        ></button>

        <div
                class="absolute right-0 z-20 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-slate-800 dark:ring-slate-700"
                role="menu"
                aria-orientation="vertical"
                aria-labelledby="user-menu-button"
        >
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-700">
                <p class="text-sm text-slate-500 dark:text-slate-400">Connecté en tant que</p>
                <p class="truncate text-sm font-medium text-slate-900 dark:text-white">{user?.username}</p>
            </div>
            <div class="py-1">
                {#if user?.role === 'admin'}
                    <a
                            href="/settings/members"
                            class="block px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                            role="menuitem"
                            onclick={close}>Gérer les membres</a
                    >
                {/if}
                <a
                        href="/settings/account"
                        class="block px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                        role="menuitem"
                        onclick={close}>Mon compte</a
                >
                <a
                        href="/dashboard"
                        class="block px-4 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                        role="menuitem"
                        onclick={close}>Mes Groupes</a
                >
            </div>
            <div class="border-t border-slate-200 py-1 dark:border-slate-700">
                <form action="/logout" method="POST" use:enhance>
                    <button
                            type="submit"
                            class="block w-full px-4 py-2 text-left text-sm text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                            role="menuitem"
                    >
                        Déconnexion
                    </button>
                </form>
            </div>
        </div>
    {/if}
</div>