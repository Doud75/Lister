<script lang="ts">
	import { onMount } from 'svelte';
	import { registerSW } from 'virtual:pwa-register';

	onMount(() => {
		registerSW({ immediate: true });
	});
</script>