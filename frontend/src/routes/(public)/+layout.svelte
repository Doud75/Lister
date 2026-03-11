<script lang="ts">
    let { children } = $props();
</script>

<div class="min-h-screen bg-slate-50 dark:bg-slate-950">
    {@render children()}
</div>
