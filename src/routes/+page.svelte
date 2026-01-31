<script lang="ts">
	import { onMount } from 'svelte';
	import StatsCard from '$lib/components/StatsCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import { stats } from '$lib/stores/stats';
	import { peers } from '$lib/stores/peers';
	import { formatBytes } from '$lib/utils/formatting';
	import { Users, Download, Upload, Search, Bell, BarChart3, Settings, TriangleAlert, ArrowRight } from 'lucide-svelte';

	// Loading state
	let loading = $state(true);

	// Load data on mount
	onMount(async () => {
		await Promise.all([stats.load(), peers.load()]);
		loading = false;
	});

	// Derived: online peers count
	let onlinePeersCount = $derived($peers.filter((p) => p.status === 'online').length);
</script>

<svelte:head>
	<title>Dashboard - WireGuard Manager</title>
</svelte:head>

<div class="mx-auto max-w-7xl">
	<!-- Page header -->
	<header class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="mb-1 text-2xl font-black tracking-tight">Dashboard</h1>
			<p class="text-sm font-medium text-slate-400">
				Monitor your WireGuard VPN network at a glance
			</p>
		</div>

		<!-- Search & User (Desktop) -->
		<div class="hidden items-center gap-6 md:flex">
			<div class="relative">
				<Search class="absolute top-1/2 left-3 -translate-y-1/2 text-slate-400" size={18} />
				<input
					class="glass-input w-64 pl-10"
					placeholder="Search peers..."
					type="text"
				/>
			</div>
			<div class="flex items-center gap-2">
				<button class="rounded-lg bg-white/5 p-2 text-slate-400 transition-colors hover:text-white">
					<Bell size={20} />
				</button>
				<div class="h-10 w-10 overflow-hidden rounded-full ring-2 ring-[#137fec]/20">
					<img
						src="https://ui-avatars.com/api/?name=Admin&background=137fec&color=fff"
						alt="User profile"
						class="h-full w-full object-cover"
					/>
				</div>
			</div>
		</div>
	</header>

	{#if loading}
		<div class="glass-card flex h-64 items-center justify-center">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $stats}
		<!-- FR-008a: 3-card horizontal grid layout for main stats -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-3">
			<StatsCard
				title="Total Peers"
				value={$stats.peerCount}
				icon={Users}
				color="purple"
				subtitle={`${onlinePeersCount} online`}
			/>
			<StatsCard
				title="Data Received"
				value={formatBytes($stats.totalRx)}
				icon={Download}
				color="green"
				subtitle="Total incoming traffic"
				trend="up"
			/>
			<StatsCard
				title="Data Sent"
				value={formatBytes($stats.totalTx)}
				icon={Upload}
				color="blue"
				subtitle="Total outgoing traffic"
				trend="up"
			/>
		</div>

		<!-- Quick actions panel -->
		<div class="glass-card mb-8 p-6">
			<h2 class="mb-4 text-xl font-semibold">Quick Actions</h2>
			<div class="flex flex-wrap gap-4">
				<a href="/peers" data-sveltekit-noscroll class="glass-btn-primary flex items-center gap-2 px-6 py-3">
					<Users size={20} /> Manage Peers
				</a>
				<a href="/stats" data-sveltekit-noscroll class="glass-btn-secondary flex items-center gap-2 px-6 py-3">
					<BarChart3 size={20} /> View Statistics
				</a>
				<a href="/settings" data-sveltekit-noscroll class="glass-btn-secondary flex items-center gap-2 px-6 py-3">
					<Settings size={20} /> Settings
				</a>
			</div>
		</div>

		<!-- Network status panel -->
		<div class="glass-card p-6">
			<h2 class="mb-4 text-xl font-semibold">Network Status</h2>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-400">Interface</p>
						<p class="font-medium">{$stats.interfaceName}</p>
					</div>
					<span class="rounded-full bg-green-500/20 px-3 py-1 text-sm text-green-400">Active</span>
				</div>
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-400">Configured Peers</p>
						<p class="font-medium">{$stats.peerCount} peers ({onlinePeersCount} online)</p>
					</div>
					<a
						href="/peers"
						data-sveltekit-noscroll
						class="flex items-center gap-1 text-sm text-blue-400 hover:text-blue-300"
					>
						View all <ArrowRight size={14} />
					</a>
				</div>
			</div>
		</div>
	{:else}
		<div class="glass-card flex flex-col items-center justify-center p-12 text-center">
			<TriangleAlert class="mb-4 text-yellow-500" size={48} />
			<p class="text-gray-400">Unable to load dashboard data. Please try again later.</p>
		</div>
	{/if}
</div>
