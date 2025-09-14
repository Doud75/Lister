<script lang="ts">
    import { enhance } from '$app/forms';
    let { data } = $props();
</script>

<div class="container mx-auto px-4 sm:px-6">
    <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">Song Library</h1>
        <a
                href="/song/new"
                class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500"
        >
            + Add New Song
        </a>
    </div>

    <div class="mt-8 rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
        <ul class="divide-y divide-slate-200 dark:divide-slate-700">
            {#each data.songs as song (song.id)}
                <li class="flex items-center justify-between py-4">
                    <span class="font-medium text-slate-800 dark:text-slate-100">{song.title}</span>
                    <div class="flex items-center gap-4">
                        <a href="/song/{song.id}/edit" class="text-sm font-semibold text-indigo-600 hover:underline dark:text-indigo-400">Edit</a>
                        <form method="POST" action="?/deleteSong" use:enhance={({form}) => {
							if (!confirm(`Are you sure you want to delete "${form.get('songTitle')}"?`)) {
								return ({cancel}) => cancel();
							}
						}}>
                            <input type="hidden" name="songId" value={song.id} />
                            <input type="hidden" name="songTitle" value={song.title} />
                            <button type="submit" class="text-sm font-semibold text-red-600 hover:underline dark:text-red-400">Delete</button>
                        </form>
                    </div>
                </li>
            {/each}
        </ul>
    </div>
</div>