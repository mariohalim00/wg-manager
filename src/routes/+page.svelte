<script lang="ts">
	import { onMount } from 'svelte';
	import StatsCard from '$lib/components/StatsCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import { stats } from '$lib/stores/stats';
	import { peers } from '$lib/stores/peers';
	import { formatBytes } from '$lib/utils/formatting';

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
	<div class="mb-8">
		<h1 class="mb-2 text-3xl font-bold">Dashboard</h1>
		<p class="text-gray-400">Monitor your WireGuard VPN network at a glance</p>
	</div>

	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $stats}
		<!-- FR-008a: 3-card horizontal grid layout for main stats -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-3">
			<StatsCard
				title="Total Peers"
				value={$stats.peerCount}
				icon="👥"
				color="purple"
				subtitle={`${onlinePeersCount} online`}
			/>
			<StatsCard
				title="Data Received"
				value={formatBytes($stats.totalRx)}
				icon="📥"
				color="green"
				subtitle="Total RX"
			/>
			<StatsCard
				title="Data Transmitted"
				value={formatBytes($stats.totalTx)}
				icon="📤"
				color="yellow"
				subtitle="Total TX"
			/>
		</div>

		<!-- Quick actions panel -->
		<div class="glass-card mb-8 p-6">
			<h2 class="mb-4 text-xl font-semibold">Quick Actions</h2>
			<div class="flex flex-wrap gap-4">
				<a href="/peers" data-sveltekit-noscroll class="glass-btn-primary px-6 py-3">
					👥 Manage Peers
				</a>
				<a href="/stats" data-sveltekit-noscroll class="glass-btn-secondary px-6 py-3">
					📊 View Statistics
				</a>
				<a href="/settings" data-sveltekit-noscroll class="glass-btn-secondary px-6 py-3">
					⚙️ Settings
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
						class="text-sm text-blue-400 hover:text-blue-300"
					>
						View all →
					</a>
				</div>
			</div>
		</div>
	{:else}
		<div class="glass-card p-12 text-center">
			<span class="mb-4 block text-4xl">⚠️</span>
			<p class="text-gray-400">Unable to load dashboard data. Please try again later.</p>
		</div>
	{/if}
</div>
