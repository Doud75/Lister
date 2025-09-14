<script lang="ts">
    let { data } = $props<{
        data: {
            userInfo: { username: string; band_name: string };
            setlists: Array<{ id: number; name: string; color: string; created_at: string }>;
        };
    }>();
</script>

<div class="container mx-auto px-4 sm:px-6">
    <header class="mb-8">
        <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
                <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                    {data.userInfo.band_name}
                </h1>
                <p class="mt-1 text-lg text-slate-600 dark:text-slate-400">
                    Welcome back, <span class="font-semibold">{data.userInfo.username}</span>!
                </p>
            </div>
            <a
                    href="/setlist/new"
                    class="flex w-auto items-center gap-2 justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-slate-900"
            >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5">
                    <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                </svg>
                Create New Setlist
            </a>
        </div>
    </header>

    <div class="space-y-8">
        <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
            <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Your Setlists</h2>

            {#if data.setlists.length > 0}
                <div class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    {#each data.setlists as setlist (setlist.id)}
                        <a
                                href="/setlist/{setlist.id}"
                                class="group relative block overflow-hidden rounded-lg border border-slate-200 shadow-sm transition-all duration-200 hover:shadow-md hover:-translate-y-1 dark:border-slate-700"
                        >
                            <div class="absolute left-0 top-0 h-full w-1.5" style="background-color: {setlist.color};"></div>
                            <div class="p-5 pl-6">
                                <h3 class="font-semibold text-slate-900 dark:text-white">{setlist.name}</h3>
                                <p class="text-sm text-slate-500 dark:text-slate-400">
                                    Created: {new Date(setlist.created_at).toLocaleDateString()}
                                </p>
                            </div>
                        </a>
                    {/each}
                </div>
            {:else}
                <div class="mt-6 rounded-lg border-2 border-dashed border-slate-300 p-12 text-center dark:border-slate-700">
                    <svg class="mx-auto h-12 w-12 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                        <path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
                    </svg>
                    <h3 class="mt-2 text-sm font-semibold text-slate-900 dark:text-white">No setlists yet</h3>
                    <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">Get started by creating your first setlist.</p>
                </div>
            {/if}
        </div>
    </div>
</div>