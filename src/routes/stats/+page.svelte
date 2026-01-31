<script lang="ts">
	import { onMount } from 'svelte';
	import StatsCard from '$lib/components/StatsCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import { stats } from '$lib/stores/stats';
	import { formatBytes } from '$lib/utils/formatting';

	// Loading state
	let loading = $state(true);

	// Load stats on mount
	onMount(async () => {
		await stats.load();
		loading = false;
	});
</script>

<svelte:head>
	<title>Statistics - WireGuard Manager</title>
</svelte:head>

<div class="max-w-7xl mx-auto">
	<!-- Page header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold mb-2">Interface Statistics</h1>
		<p class="text-gray-400">Monitor VPN network performance and usage</p>
	</div>

	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $stats}
		<!-- Stats cards grid -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<StatsCard
				title="Interface"
				value={$stats.interfaceName}
				icon="🔧"
				color="blue"
				subtitle="WireGuard interface"
			/>
			<StatsCard
				title="Total Peers"
				value={$stats.peerCount}
				icon="👥"
				color="purple"
				subtitle={`${$stats.peerCount === 1 ? 'peer' : 'peers'} configured`}
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

		<!-- Additional info panel -->
		<div class="glass-card p-6">
			<h2 class="text-xl font-semibold mb-4">Network Overview</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<h3 class="text-sm font-medium text-gray-400 mb-2">Interface Details</h3>
					<dl class="space-y-2">
						<div class="flex justify-between">
							<dt class="text-gray-400">Name:</dt>
							<dd class="font-medium">{$stats.interfaceName}</dd>
						</div>
						<div class="flex justify-between">
							<dt class="text-gray-400">Status:</dt>
							<dd class="text-green-400 font-medium">Active</dd>
						</div>
					</dl>
				</div>
				<div>
					<h3 class="text-sm font-medium text-gray-400 mb-2">Traffic Summary</h3>
					<dl class="space-y-2">
						<div class="flex justify-between">
							<dt class="text-gray-400">Total Data:</dt>
							<dd class="font-medium">{formatBytes($stats.totalRx + $stats.totalTx)}</dd>
						</div>
						<div class="flex justify-between">
							<dt class="text-gray-400">Configured Peers:</dt>
							<dd class="font-medium">{$stats.peerCount}</dd>
						</div>
					</dl>
				</div>
			</div>
		</div>
	{:else}
		<div class="glass-card p-12 text-center">
			<span class="text-4xl mb-4 block">⚠️</span>
			<p class="text-gray-400">Unable to load statistics. Please try again later.</p>
		</div>
	{/if}
</div>

