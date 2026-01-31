<script lang="ts">
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Notification from '$lib/components/Notification.svelte';
	import { notifications } from '$lib/stores/notifications';

	let { children } = $props();
</script>

<svelte:head>
	<title>WireGuard Manager</title>
</svelte:head>

<!-- Main layout with gradient background -->
<div class="min-h-screen flex">
	<!-- Sidebar navigation -->
	<Sidebar />

	<!-- Main content area -->
	<main class="flex-1 p-6 md:p-8 overflow-y-auto">
		{@render children()}
	</main>

	<!-- Notification stack (fixed position, top-right) -->
	<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 w-80 max-w-[calc(100vw-2rem)]">
		{#each $notifications as notification (notification.id)}
			<Notification {notification} />
		{/each}
	</div>
</div>

