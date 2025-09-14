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
            <a href="/setlist/new" class="flex w-auto justify-center rounded-md bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500">
                + Create New Setlist
            </a>
        </div>
    </header>

    <div class="space-y-8">
        <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
            <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Your Setlists</h2>

            {#if data.setlists.length > 0}
                <div class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    {#each data.setlists as setlist (setlist.id)}
                        <a href="/setlist/{setlist.id}" class="group block">
                            <div class="h-full rounded-lg border border-slate-200 bg-slate-50 p-5 shadow-sm transition-all duration-200 group-hover:-translate-y-1 group-hover:shadow-md dark:border-slate-700 dark:bg-slate-800/50 dark:group-hover:border-slate-600">
                                <div class="flex items-start gap-4">
                                    <span class="mt-1 block h-4 w-4 flex-shrink-0 rounded-full" style="background-color: {setlist.color};"></span>
                                    <div>
                                        <h3 class="font-semibold text-slate-900 dark:text-white">{setlist.name}</h3>
                                        <p class="text-sm text-slate-500 dark:text-slate-400">
                                            Created: {new Date(setlist.created_at).toLocaleDateString()}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </a>
                    {/each}
                </div>
            {:else}
                <div class="mt-6 rounded-lg border-2 border-dashed border-slate-300 p-12 text-center dark:border-slate-700">
                    <svg class="mx-auto h-12 w-12 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                        <path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
                    </svg>
                    <h3 class="mt-2 text-sm font-semibold text-slate-900 dark:text-white">No setlists</h3>
                    <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">Get started by creating a new setlist.</p>
                </div>
            {/if}
        </div>
        <div class="rounded-xl bg-white p-6 shadow-lg dark:bg-slate-800">
            <h2 class="text-xl font-semibold text-slate-800 dark:text-slate-100">Bibliothèque</h2>
            <div class="mt-4">
                <a
                        href="/song"
                        class="group flex items-center justify-between rounded-lg border border-slate-200 bg-slate-50 p-4 transition-colors hover:bg-slate-100 dark:border-slate-700 dark:bg-slate-800/50 dark:hover:bg-slate-700/70"
                >
                    <div class="flex items-center gap-4">
                        <div
                                class="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-100 dark:bg-indigo-500/20"
                        >
                            <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke-width="1.5"
                                    stroke="currentColor"
                                    class="h-6 w-6 text-indigo-600 dark:text-indigo-400"
                            >
                                <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z"
                                />
                            </svg>
                        </div>
                        <div>
                            <h3 class="font-semibold text-slate-900 dark:text-white">Gérer les chansons</h3>
                            <p class="text-sm text-slate-500 dark:text-slate-400">
                                Voir et modifier tous les morceaux
                            </p>
                        </div>
                    </div>
                    <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            class="h-5 w-5 text-slate-400 transition-transform group-hover:translate-x-1"
                    >
                        <path
                                fill-rule="evenodd"
                                d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
                                clip-rule="evenodd"
                        />
                    </svg>
                </a>
            </div>
        </div>
    </div>
</div>