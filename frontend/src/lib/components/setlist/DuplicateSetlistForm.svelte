<script lang="ts">
    import { enhance } from '$app/forms';
    import Button from '$lib/components/ui/Button.svelte';
    import Input from '$lib/components/ui/Input.svelte';

    let {
        setlistName,
        setColor,
        close
    } = $props<{
        setlistName: string;
        setColor: string;
        close: () => void;
    }>();

    let newSetName = $state(`Copie de ${setlistName}`);
    let newSetColor = $state(setColor);

    $effect(() => {
        newSetName = `Copie de ${setlistName}`;
        newSetColor = setColor;
    });

    let isSubmitting = $state(false);
</script>

<h3 class="text-lg font-semibold text-slate-900 dark:text-white">Dupliquer la setlist</h3>
<p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
    Créez une copie de "{setlistName}".
</p>
<form
        method="POST"
        action="?/duplicateSetlist"
        class="mt-4 space-y-4"
        use:enhance={() => {
            isSubmitting = true;

            return async ({ update }) => {
                await update();
                isSubmitting = false;
                close();
            };
	    }}
>
	<Input label="Nouveau nom" id="name" name="name" bind:value={newSetName} required />
	<div>
		<label for="color" class="block text-sm font-medium leading-6 text-slate-900 dark:text-slate-200">
			Nouvelle couleur
		</label>
		<div class="mt-2">
			<input
				type="color"
				id="color"
				name="color"
				bind:value={newSetColor}
				class="h-10 w-16 cursor-pointer rounded-md border-0 bg-transparent p-0"
			/>
		</div>
	</div>
	<div class="flex justify-end gap-3 pt-4">
		<Button type="button" variant="secondary" onclick={close} autoWidth>Annuler</Button>

		<Button type="submit" autoWidth isLoading={isSubmitting}>
			Créer la copie
		</Button>
	</div>
</form>