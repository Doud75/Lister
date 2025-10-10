<script lang="ts">
    import type {Snippet} from "svelte";

    let {
        isOpen,
        onClose,
        children
    } = $props<{
        isOpen: boolean;
        onClose: () => void;
        children: Snippet;
    }>();
</script>

{#if isOpen}
    <div
            class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            role="dialog"
            aria-modal="true"
            tabindex="0"
            onclick={onClose}
            onkeydown={(e) => (e.key === "Escape") && onClose()}
    >
        <div
                class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl dark:bg-slate-800 max-h-[90vh] overflow-y-auto"
                role="presentation"
                onclick={(event) => event.stopPropagation()}
        >
            {@render children()}
        </div>
    </div>

{/if}