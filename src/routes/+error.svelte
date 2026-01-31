<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
</script>

<svelte:head>
	<title>Error - WireGuard Manager</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4">
	<div class="glass-card p-8 max-w-md w-full text-center">
		<span class="text-6xl mb-4 block">⚠️</span>
		<h1 class="text-3xl font-bold mb-2">Something Went Wrong</h1>
		
		{#if $page.error}
			<p class="text-gray-400 mb-6">
				{$page.error.message || 'An unexpected error occurred'}
			</p>
		{:else}
			<p class="text-gray-400 mb-6">
				An unexpected error occurred while loading this page.
			</p>
		{/if}

		<div class="flex gap-3 justify-center">
			<button
				class="btn btn-primary"
				onclick={() => goto('/')}
			>
				Go Home
			</button>
			<button
				class="btn btn-outline"
				onclick={() => window.location.reload()}
			>
				Reload Page
			</button>
		</div>

		{#if $page.status}
			<p class="text-sm text-gray-500 mt-6">Error Code: {$page.status}</p>
		{/if}
	</div>
</div>
