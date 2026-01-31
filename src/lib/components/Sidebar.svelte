<script lang="ts">
	import { page } from '$app/stores';
	import { stats } from '$lib/stores/stats';
	import { formatBytes } from '$lib/utils/formatting';

	// Mobile menu state
	let mobileMenuOpen = $state(false);

	// Navigation items
	const navItems = [
		{ path: '/', label: 'Dashboard', icon: '📊' },
		{ path: '/peers', label: 'Peers', icon: '👥' },
		{ path: '/stats', label: 'Statistics', icon: '📈' },
		{ path: '/settings', label: 'Settings', icon: '⚙️' }
	];

	// Check if route is active
	function isActive(path: string): boolean {
		if (path === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(path);
	}
</script>

<!-- Mobile hamburger button -->
<button
	class="glass-btn-primary fixed top-4 left-4 z-50 p-2 md:hidden"
	onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
	aria-label="Toggle menu"
>
	{#if mobileMenuOpen}
		<span class="text-2xl">✕</span>
	{:else}
		<span class="text-2xl">☰</span>
	{/if}
</button>

<!-- Sidebar (responsive: hidden on mobile unless menu open) -->
<aside
	class="glass-card fixed z-40 flex h-screen w-64 flex-col gap-6 p-6 transition-transform duration-300 md:static {mobileMenuOpen
		? 'translate-x-0'
		: '-translate-x-full md:translate-x-0'}"
>
	<!-- Logo/Title -->
	<div class="flex items-center gap-3">
		<span class="text-3xl">🔐</span>
		<div>
			<h1 class="text-xl font-bold">WireGuard</h1>
			<p class="text-sm text-gray-400">Manager</p>
		</div>
	</div>

	<!-- Navigation links -->
	<nav class="flex flex-col gap-2">
		{#each navItems as item (item.path)}
			<a
				href={item.path}
				data-sveltekit-noscroll
				class="flex items-center gap-3 rounded-lg px-4 py-3 transition-all {isActive(item.path)
					? 'bg-glass-hover text-white'
					: 'hover:bg-glass-bg text-gray-300 hover:text-white'}"
				onclick={() => (mobileMenuOpen = false)}
			>
				<span class="text-xl">{item.icon}</span>
				<span class="font-medium">{item.label}</span>
			</a>
		{/each}
	</nav>

	<!-- Usage widget (bottom of sidebar) -->
	<div class="glass-card mt-auto p-4">
		<h3 class="mb-3 text-sm font-semibold text-gray-300">Interface Usage</h3>
		{#if $stats}
			<div class="space-y-2 text-sm">
				<div class="flex justify-between">
					<span class="text-gray-400">Interface:</span>
					<span class="font-medium">{$stats.interfaceName}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-gray-400">Peers:</span>
					<span class="font-medium">{$stats.peerCount}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-gray-400">RX:</span>
					<span class="font-medium text-green-400">{formatBytes($stats.totalRx)}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-gray-400">TX:</span>
					<span class="font-medium text-blue-400">{formatBytes($stats.totalTx)}</span>
				</div>
			</div>
		{:else}
			<p class="text-sm text-gray-400">Loading stats...</p>
		{/if}
	</div>
</aside>

<!-- Overlay for mobile menu (click outside to close) -->
{#if mobileMenuOpen}
	<button
		class="fixed inset-0 z-30 bg-black/50 md:hidden"
		onclick={() => (mobileMenuOpen = false)}
		aria-label="Close menu"
	></button>
{/if}
