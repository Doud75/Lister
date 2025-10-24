<script lang="ts">
    import { onMount } from 'svelte';
    import { registerSW } from 'virtual:pwa-register';

    onMount(() => {
        registerSW({
            onRegisteredSW(swUrl, r) {
                if (r) {
                    console.log(`Service Worker enregistré : ${swUrl}`);
                }
            },
            onRegisterError(error) {
                console.error("Erreur lors de l'enregistrement du Service Worker:", error);
            }
        });
    });
</script>

