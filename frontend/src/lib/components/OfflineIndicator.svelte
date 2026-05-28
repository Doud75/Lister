<script lang="ts">
	import { browser } from '$app/environment';

	let isOnline = $state(browser ? navigator.onLine : true);

	$effect(() => {
		if (!browser) return;

		const handleOnline = () => (isOnline = true);
		const handleOffline = () => (isOnline = false);

		window.addEventListener('online', handleOnline);
		window.addEventListener('offline', handleOffline);

		return () => {
			window.removeEventListener('online', handleOnline);
			window.removeEventListener('offline', handleOffline);
		};
	});
</script>

{#if !isOnline}
	<div
		class="fixed top-4 left-1/2 z-50 -translate-x-1/2 flex items-center gap-2 rounded-full bg-slate-800 px-4 py-2 text-sm font-medium text-white shadow-lg dark:bg-slate-700"
		role="status"
		aria-live="polite"
		data-testid="offline-indicator"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="h-4 w-4 flex-shrink-0"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
			stroke-width="2"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M8.288 15.038a5.25 5.25 0 0 1 7.424 0M5.106 11.856c3.807-3.808 9.98-3.808 13.788 0M1.924 8.674c5.565-5.565 14.587-5.565 20.152 0M12.53 18.22l-.53.53-.53-.53a.75.75 0 0 1 1.06 0Z"
			/>
			<line x1="2" y1="2" x2="22" y2="22" stroke-linecap="round" stroke-width="2" />
		</svg>
		Hors ligne
	</div>
{/if}
