<script lang="ts">
	import { page } from '$app/stores';
</script>

<svelte:head>
	<title>Error - WireGuard Manager</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center p-4">
	<div class="glass-card w-full max-w-md p-8 text-center">
		<span class="mb-4 block text-6xl">⚠️</span>
		<h1 class="mb-2 text-3xl font-bold">Something Went Wrong</h1>

		{#if $page.error}
			<p class="mb-6 text-gray-400">
				{$page.error.message || 'An unexpected error occurred'}
			</p>
		{:else}
			<p class="mb-6 text-gray-400">An unexpected error occurred while loading this page.</p>
		{/if}

		<div class="flex justify-center gap-3">
			<a href="/" data-sveltekit-noscroll class="btn btn-primary"> Go Home </a>
			<button class="btn btn-outline" onclick={() => window.location.reload()}>
				Reload Page
			</button>
		</div>

		{#if $page.status}
			<p class="mt-6 text-sm text-gray-500">Error Code: {$page.status}</p>
		{/if}
	</div>
</div>
