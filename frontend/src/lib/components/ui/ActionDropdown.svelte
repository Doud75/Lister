<script lang="ts">
    import type { Snippet } from 'svelte';

    let { children } = $props<{
        children?: Snippet<[ { close: () => void } ]>;
    }>();

    let isOpen = $state(false);

    function close() {
        isOpen = false;
    }
</script>

<div class="relative">
    <button
            onclick={() => (isOpen = !isOpen)}
            class="flex h-10 w-10 items-center justify-center rounded-full text-sm font-semibold text-slate-500 hover:bg-slate-200 dark:text-slate-400 dark:hover:bg-slate-700"
            aria-haspopup="true"
            aria-expanded={isOpen}
            aria-label="Ouvrir le menu d'actions"
    >
        <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                class="h-5 w-5"
        >
            <path
                    d="M10 3a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM10 8.5a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM11.5 15.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0Z"
            />
        </svg>
    </button>

    {#if isOpen}
        <div onclick={close} class="fixed inset-0 z-10 h-full w-full"></div>

        <div
                class="absolute right-0 z-20 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-slate-800 dark:ring-slate-700"
                role="menu"
                aria-orientation="vertical"
        >
            {@render children?.({ close })}
        </div>
    {/if}
</div>