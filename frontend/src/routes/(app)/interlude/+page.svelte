<script lang="ts">
    import { enhance } from '$app/forms';
    let { data } = $props();
</script>

<div class="container mx-auto px-4 sm:px-6">
    <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">Interlude Library</h1>
        <a
                href="/interlude/new"
                class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
        >
            + Add New Interlude
        </a>
    </div>

    <div class="mt-8 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        {#if data.interludes && data.interludes.length > 0}
            <ul class="divide-y divide-slate-200 dark:divide-slate-700">
                {#each data.interludes as interlude (interlude.id)}
                    <li class="flex items-center justify-between py-4">
                        <span class="font-medium text-slate-800 dark:text-slate-100">{interlude.title}</span>
                        <div class="flex items-center gap-4">
                            <a href="/interlude/{interlude.id}/edit" class="text-sm font-semibold text-indigo-600 hover:underline dark:text-indigo-400">Edit</a>
                            <form method="POST" action="?/deleteInterlude" use:enhance={({form}) => {
                                if (!confirm(`Are you sure you want to delete "${form.get('interludeTitle')}"?`)) {
                                    return ({cancel}) => cancel();
                                }
                            }}>
                                <input type="hidden" name="interludeId" value={interlude.id} />
                                <input type="hidden" name="interludeTitle" value={interlude.title} />
                                <button type="submit" class="text-sm font-semibold text-red-600 hover:underline dark:text-red-400">Delete</button>
                            </form>
                        </div>
                    </li>
                {/each}
            </ul>
        {:else}
            <p class="py-10 text-center text-sm text-slate-500 dark:text-slate-400">You have no saved interludes.</p>
        {/if}
    </div>
</div>