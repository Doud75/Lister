<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import PWABadge from '$lib/components/PWABadge.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<PWABadge/>

{@render children?.()}
