<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
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

<div class="max-w-7xl mx-auto">
	<!-- Page header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold mb-2">Dashboard</h1>
		<p class="text-gray-400">Monitor your WireGuard VPN network at a glance</p>
	</div>

	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $stats}
		<!-- FR-008a: 3-card horizontal grid layout for main stats -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
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
		<div class="glass-card p-6 mb-8">
			<h2 class="text-xl font-semibold mb-4">Quick Actions</h2>
			<div class="flex flex-wrap gap-4">
				<button onclick={() => goto('/peers')} class="glass-btn-primary px-6 py-3">
					👥 Manage Peers
				</button>
				<button onclick={() => goto('/stats')} class="glass-btn-secondary px-6 py-3">
					📊 View Statistics
				</button>
				<button onclick={() => goto('/settings')} class="glass-btn-secondary px-6 py-3">
					⚙️ Settings
				</button>
			</div>
		</div>

		<!-- Network status panel -->
		<div class="glass-card p-6">
			<h2 class="text-xl font-semibold mb-4">Network Status</h2>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-400">Interface</p>
						<p class="font-medium">{$stats.interfaceName}</p>
					</div>
					<span class="px-3 py-1 rounded-full bg-green-500/20 text-green-400 text-sm">Active</span>
				</div>
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-400">Configured Peers</p>
						<p class="font-medium">{$stats.peerCount} peers ({onlinePeersCount} online)</p>
					</div>
					<button onclick={() => goto('/peers')} class="text-blue-400 hover:text-blue-300 text-sm">
						View all →
					</button>
				</div>
			</div>
		</div>
	{:else}
		<div class="glass-card p-12 text-center">
			<span class="text-4xl mb-4 block">⚠️</span>
			<p class="text-gray-400">Unable to load dashboard data. Please try again later.</p>
		</div>
	{/if}
</div>

                <h2 class="card-title">Online Peers</h2>
                <p class="text-4xl font-bold">{$stats.onlinePeers}</p>
            </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Data Sent</h2>
                <p class="text-4xl font-bold">{$stats.totalDataUsage.sent} KB</p>
            </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title">Data Received</h2>
                <p class="text-4xl font-bold">{$stats.totalDataUsage.received} KB</p>
            </div>
        </div>
    </div>
</div>
