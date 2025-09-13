<script lang="ts">
    let { children } = $props();
</script>

<div class="min-h-screen bg-slate-100 dark:bg-slate-900">
    <header class="bg-white shadow-sm dark:bg-slate-800">
        <nav class="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6">
            <div class="flex items-center gap-4">
                <a href="/" class="flex-shrink-0" aria-label="Go to Dashboard">
                    <svg class="h-8 w-auto text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z" />
                    </svg>
                </a>
                <div class="hidden md:block">
                    <a href="/" class="rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-700">Dashboard</a>
                </div>
            </div>
            <div class="flex items-center">
                <a href="/logout" class="rounded-md bg-slate-200 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                    Logout
                </a>
            </div>
        </nav>
    </header>

    <main class="py-10">
        {@render children?.()}
    </main>
</div>