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
	class="fixed top-4 left-4 z-50 md:hidden glass-btn-primary p-2"
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
	class="glass-card w-64 h-screen p-6 flex flex-col gap-6 fixed md:static transition-transform duration-300 z-40 {mobileMenuOpen
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
		{#each navItems as item}
			<a
				href={item.path}
				class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all {isActive(item.path)
					? 'bg-glass-hover text-white'
					: 'text-gray-300 hover:bg-glass-bg hover:text-white'}"
				onclick={() => (mobileMenuOpen = false)}
			>
				<span class="text-xl">{item.icon}</span>
				<span class="font-medium">{item.label}</span>
			</a>
		{/each}
	</nav>

	<!-- Usage widget (bottom of sidebar) -->
	<div class="mt-auto glass-card p-4">
		<h3 class="text-sm font-semibold mb-3 text-gray-300">Interface Usage</h3>
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
		class="fixed inset-0 bg-black/50 z-30 md:hidden"
		onclick={() => (mobileMenuOpen = false)}
		aria-label="Close menu"
	></button>
{/if}
