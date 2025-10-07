<script lang="ts">
    import { navigating } from '$app/stores';
    import { browser } from '$app/environment';
    import type { PageData } from './$types';
    import type {Snippet} from "svelte";

    let { children, data }: { children: Snippet; data: PageData } = $props();

    const SWIPE_THRESHOLD = 70;
    const SWIPE_EDGE_WIDTH = 50;

    let isSwiping = $state(false);
    let startX = $state(0);
    let translateX = $state(0);
    let canNavigateBack = $state(browser && window.history.length > 1);
    let containerEl: HTMLElement | undefined = $state();

    function updateNavigationState() {
        canNavigateBack = window.history.length > 1;
    }

    function handleTouchStart(e: TouchEvent) {
        if (e.touches.length !== 1) return;
        const touchX = e.touches[0].clientX;
        if (touchX > SWIPE_EDGE_WIDTH && touchX < window.innerWidth - SWIPE_EDGE_WIDTH) return;
        startX = touchX;
        isSwiping = true;
        if (containerEl) containerEl.style.transition = 'none';
    }

    function handleTouchMove(e: TouchEvent) {
        if (!isSwiping || e.touches.length !== 1) return;
        const currentX = e.touches[0].clientX;
        let deltaX = currentX - startX;
        if (deltaX > 0 && !canNavigateBack) deltaX = 0;
        translateX = deltaX;
    }

    function handleTouchEnd() {
        if (!isSwiping) return;
        if (containerEl) containerEl.style.transition = 'transform 0.3s ease';
        if (translateX > SWIPE_THRESHOLD && canNavigateBack) {
            history.back();
            setTimeout(() => { updateNavigationState(); }, 300);
        } else if (translateX < -SWIPE_THRESHOLD) {
            history.forward();
            setTimeout(updateNavigationState, 300);
        }
        isSwiping = false;
        translateX = 0;
    }

    $effect(() => {
        if ($navigating) {
            isSwiping = false;
            translateX = 0;
            if (containerEl) {
                containerEl.style.transition = 'none';
                containerEl.style.transform = 'translateX(0)';
            }
            setTimeout(updateNavigationState, 50);
        }
    });
</script>

<div class="min-h-screen bg-slate-100 dark:bg-slate-900">
    <header class="bg-white shadow-sm dark:bg-slate-800">
        <nav class="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6">
            <div class="flex items-center gap-4">
                <a href="/" class="flex-shrink-0" aria-label="Go to Home">
                    <svg class="h-8 w-auto text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z" />
                    </svg>
                </a>

                {#if data.userBands && data.userBands.length > 1}
                    <div class="ml-4">
                        <form method="POST" action="/switch-band">
                            <select
                                    name="bandId"
                                    class="p-2 rounded-md border-slate-300 text-sm font-medium shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:border-slate-600 dark:bg-slate-700 dark:text-white"
                                    onchange={(e) => e.currentTarget.form?.requestSubmit()}
                            >
                                {#each data.userBands as band (band.id)}
                                    <option value={band.id} selected={band.id.toString() === data.activeBandId}>
                                        {band.name}
                                    </option>
                                {/each}
                            </select>
                        </form>
                    </div>
                {/if}
            </div>
            <div class="flex items-center gap-4">
                {#if data.user?.role === 'admin'}
                    <a href="/settings/members" class="rounded-md bg-slate-200 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                        Membres
                    </a>
                {/if}
                <a href="/settings/account" class="rounded-md bg-slate-200 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                    Account
                </a>
                <a href="/logout" class="rounded-md bg-slate-200 px-3 py-1.5 text-sm font-medium text-slate-700 hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600">
                    Logout
                </a>
            </div>
        </nav>
    </header>

    <main
            class="py-10"
            ontouchstart={handleTouchStart}
            ontouchmove={handleTouchMove}
            ontouchend={handleTouchEnd}
            bind:this={containerEl}
            style:transform="translateX({translateX}px)"
            style:transition={isSwiping ? 'none' : 'transform 0.3s ease'}
            style:touch-action="pan-y"
    >
        {@render children?.()}
    </main>
</div>

<style>
    main {
        touch-action: pan-y;
    }
</style>