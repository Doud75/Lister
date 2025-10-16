<script lang="ts">
    let {
        tabs,
        activeTab = $bindable()
    } = $props<{
        tabs: string[];
        activeTab?: string;
    }>();
</script>

<div class="border-b border-slate-200 dark:border-slate-700">
    <nav class="-mb-px flex space-x-6" aria-label="Tabs">
        {#each tabs as tab (tab)}
            <button
                    onclick={() => (activeTab = tab)}
                    class="whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium {activeTab === tab
					? 'border-indigo-500 text-indigo-600 dark:text-indigo-400'
					: 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-200'}"
            >
                {tab}
            </button>
        {/each}
    </nav>
</div>