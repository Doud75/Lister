<script lang="ts">
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';

    let {
        setlistName,
        close
    } = $props<{
        setlistName: string;
        close: () => void;
    }>();

    let isSubmitting = $state(false);
</script>

<h3 class="text-lg font-semibold text-slate-900 dark:text-white">Supprimer la setlist</h3>
<p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
    Êtes-vous sûr de vouloir supprimer la setlist <span class="font-bold">"{setlistName}"</span> ? Cette
    action est irréversible et supprimera tous les éléments qu'elle contient.
</p>
<form
        method="POST"
        action="?/deleteSetlist"
        class="mt-6"
        use:enhance={() => {
		isSubmitting = true;
		return async ({ update }) => {
			await update();
			isSubmitting = false;
			close();
		};
	}}
>
    <div class="flex justify-end gap-3">
        <Button type="button" variant="secondary" onclick={close} autoWidth>Annuler</Button>

        <button
                type="submit"
                disabled={isSubmitting}
                class="flex w-auto justify-center rounded-md bg-red-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-red-500 disabled:cursor-not-allowed disabled:opacity-60"
        >
            {#if isSubmitting}
                <svg
                        class="h-5 w-5 animate-spin text-white"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                >
                    <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                    ></circle>
                    <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                </svg>
            {:else}
                Confirmer la suppression
            {/if}
        </button>
    </div>
</form>