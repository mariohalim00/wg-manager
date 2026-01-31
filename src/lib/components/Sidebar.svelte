<script lang="ts">
	import { page } from '$app/stores';
	import { stats } from '$lib/stores/stats';
	import { formatBytes } from '$lib/utils/formatting';
	import { LayoutDashboard, Users, BarChart3, Settings, ShieldCheck, Menu, X } from 'lucide-svelte';

	// Mobile menu state
	let mobileMenuOpen = $state(false);

	// Navigation items with Lucide icons
	const navItems = [
		{ path: '/', label: 'Dashboard', icon: LayoutDashboard },
		{ path: '/peers', label: 'Peers', icon: Users },
		{ path: '/stats', label: 'Statistics', icon: BarChart3 },
		{ path: '/settings', label: 'Settings', icon: Settings }
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
	class="fixed top-4 left-4 z-50 rounded-lg bg-white/5 p-2 backdrop-blur-md md:hidden"
	onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
	aria-label="Toggle menu"
>
	{#if mobileMenuOpen}
		<X class="h-6 w-6" />
	{:else}
		<Menu class="h-6 w-6" />
	{/if}
</button>

<!-- Sidebar (responsive: hidden on mobile unless menu open) -->
<aside
	class="sidebar-glass fixed z-40 flex h-screen w-64 flex-col transition-transform duration-300 md:static {mobileMenuOpen
		? 'translate-x-0'
		: '-translate-x-full md:translate-x-0'}"
>
	<!-- Logo/Title -->
	<div class="p-6">
		<div class="mb-10 flex items-center gap-3">
			<div class="flex items-center justify-center rounded-lg bg-[#137fec] p-2">
				<ShieldCheck class="h-6 w-6 text-white" />
			</div>
			<div>
				<h1 class="text-lg font-bold tracking-tight">WireGuard</h1>
				<p class="text-xs font-medium uppercase tracking-wider text-[#137fec]">Manager</p>
			</div>
		</div>

		<!-- Navigation links -->
		<nav class="flex flex-col gap-2">
			{#each navItems as item (item.path)}
				{@const IconComponent = item.icon}
				<a
					href={item.path}
					data-sveltekit-noscroll
					class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-colors {isActive(
						item.path
					)
						? 'border border-[#137fec]/20 bg-[#137fec]/10 text-[#137fec]'
						: 'text-slate-300 hover:bg-white/5'}"
					onclick={() => (mobileMenuOpen = false)}
				>
					<IconComponent class="h-5 w-5" />
					<span>{item.label}</span>
				</a>
			{/each}
		</nav>
	</div>

	<!-- Usage widget (bottom of sidebar) -->
	<div class="mt-auto flex flex-col gap-4 p-6">
		<div class="glass-card rounded-xl p-4">
			<p class="mb-2 text-[10px] font-bold uppercase text-slate-500">Interface Usage</p>
			{#if $stats}
				<div class="space-y-2 text-sm">
					<div class="flex justify-between">
						<span class="text-slate-400">Interface:</span>
						<span class="font-medium">{$stats.interfaceName}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-slate-400">Peers:</span>
						<span class="font-medium">{$stats.peerCount}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-slate-400">RX:</span>
						<span class="font-medium text-green-400">{formatBytes($stats.totalRx)}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-slate-400">TX:</span>
						<span class="font-medium text-blue-400">{formatBytes($stats.totalTx)}</span>
					</div>
				</div>
			{:else}
				<p class="text-sm text-slate-400">Loading stats...</p>
			{/if}
		</div>
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
