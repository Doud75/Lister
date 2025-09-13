<script lang="ts">
    import { onMount } from 'svelte';
    import { registerSW } from 'virtual:pwa-register';

    let offlineReady = false;
    let needRefresh = false;
    let updateServiceWorker: (reloadPage?: boolean) => Promise<void>;

    onMount(() => {
        updateServiceWorker = registerSW({
            onOfflineReady() {
                offlineReady = true;
            },
            onNeedRefresh() {
                needRefresh = true;
            },
            onRegistered(r: ServiceWorkerRegistration | undefined) {
                console.log('Service Worker registered:', r);
            },
            onRegisterError(error: unknown) {
                console.log('Service Worker registration error:', error);
            }
        });
    });

    function handleRefresh() {
        if (needRefresh) {
            updateServiceWorker(true);
        }
    }

    function close() {
        offlineReady = false;
        needRefresh = false;
    }
</script>

{#if offlineReady || needRefresh}
    <div class="pwa-toast" role="alert">
        <div class="message">
            {#if needRefresh}
                <span>New content available, click on reload button to update.</span>
            {/if}
        </div>
        {#if needRefresh}
            <button onclick={handleRefresh}>Reload</button>
        {/if}
        <button onclick={close}>Close</button>
    </div>
{/if}

<style>
    .pwa-toast {
        position: fixed;
        right: 16px;
        bottom: 16px;
        background: #222;
        color: white;
        padding: 12px;
        border-radius: 4px;
        z-index: 1000;
        display: flex;
        align-items: center;
        gap: 12px;
    }
    button {
        border: 1px solid white;
        background: transparent;
        color: white;
        border-radius: 2px;
        padding: 4px 8px;
        cursor: pointer;
    }
</style>