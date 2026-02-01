<script lang="ts">
	import { page } from '$app/stores';
	import { stats } from '$lib/stores/stats';
	import {
		LayoutDashboard,
		Users,
		Settings,
		ShieldCheck,
		Menu,
		X,
		PlusCircle
	} from 'lucide-svelte';

	// Mobile menu state
	let mobileMenuOpen = $state(false);

	// Navigation items matching mockup design
	const navItems = [
		{ path: '/', label: 'Dashboard', icon: LayoutDashboard },
		{ path: '/peers', label: 'Peers', icon: Users },
		{ path: '/settings', label: 'Settings', icon: Settings }
	];

	// Check if route is active
	function isActive(path: string): boolean {
		if (path === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(path);
	}

	// Calculate usage percentage (mock for now, can be replaced with real data)
	let usagePercentage = $derived(
		$stats
			? Math.min(
					100,
					Math.round((($stats.totalRx + $stats.totalTx) / (10 * 1024 * 1024 * 1024)) * 100)
				)
			: 0
	);
	let usageGB = $derived(
		$stats ? (($stats.totalRx + $stats.totalTx) / (1024 * 1024 * 1024)).toFixed(1) : '0'
	);
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

<!-- Sidebar matching mockup design -->
<aside
	class="sidebar-glass fixed z-40 flex h-screen w-64 flex-col transition-transform duration-300 md:static {mobileMenuOpen
		? 'translate-x-0'
		: '-translate-x-full md:translate-x-0'}"
>
	<!-- Logo/Brand Section -->
	<div class="p-6">
		<div class="mb-10 flex items-center gap-3">
			<div class="flex items-center justify-center rounded-lg bg-[#137fec] p-2">
				<ShieldCheck class="h-6 w-6 text-white" />
			</div>
			<div>
				<h1 class="text-lg font-bold tracking-tight">WireGuard</h1>
				<p class="text-xs font-medium tracking-wider text-[#137fec] uppercase">v2.4.0 Active</p>
			</div>
		</div>

		<!-- Navigation links -->
		<nav class="flex flex-col gap-2">
			{#each navItems as item (item.path)}
				{@const IconComponent = item.icon}
				<a
					href={item.path}
					data-sveltekit-noscroll
					class="flex items-center gap-3 rounded-xl px-4 py-3 transition-colors {isActive(item.path)
						? 'border border-[#137fec]/20 bg-[#137fec]/10 text-[#137fec]'
						: 'text-slate-300 hover:bg-white/5'}"
					onclick={() => (mobileMenuOpen = false)}
				>
					<IconComponent class="h-5 w-5" />
					<span class="text-sm font-medium">{item.label}</span>
				</a>
			{/each}
		</nav>
	</div>

	<!-- Bottom Section: Usage + Add Peer Button -->
	<div class="mt-auto flex flex-col gap-4 p-6">
		<!-- Usage Limit Widget -->
		<div class="glass rounded-xl p-4">
			<p class="mb-2 text-[10px] font-bold text-slate-500 uppercase">Usage Limit</p>
			<div class="h-1.5 w-full overflow-hidden rounded-full bg-white/10">
				<div
					class="h-full bg-[#137fec] transition-all duration-500"
					style="width: {Math.min(usagePercentage, 100)}%"
				></div>
			</div>
			<p class="mt-2 text-[11px] text-slate-400">{usageGB} GB of 10 GB used</p>
		</div>

		<!-- Add New Peer Button -->
		<a
			href="/peers"
			class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#137fec] py-3 text-sm font-bold text-white shadow-lg shadow-[#137fec]/20 transition-all hover:bg-[#137fec]/90"
		>
			<PlusCircle class="h-4 w-4" />
			Add New Peer
		</a>
	</div>
</aside>

<!-- Overlay for mobile menu -->
{#if mobileMenuOpen}
	<button
		class="fixed inset-0 z-30 bg-black/50 md:hidden"
		onclick={() => (mobileMenuOpen = false)}
		aria-label="Close menu"
	></button>
{/if}
