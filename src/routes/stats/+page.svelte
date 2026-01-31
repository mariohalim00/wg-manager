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

<div class="mx-auto max-w-7xl">
	<!-- Page header -->
	<div class="mb-8">
		<h1 class="mb-2 text-3xl font-bold">Interface Statistics</h1>
		<p class="text-gray-400">Monitor VPN network performance and usage</p>
	</div>

	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else if $stats}
		<!-- Stats cards grid -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
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
			<h2 class="mb-4 text-xl font-semibold">Network Overview</h2>
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				<div>
					<h3 class="mb-2 text-sm font-medium text-gray-400">Interface Details</h3>
					<dl class="space-y-2">
						<div class="flex justify-between">
							<dt class="text-gray-400">Name:</dt>
							<dd class="font-medium">{$stats.interfaceName}</dd>
						</div>
						<div class="flex justify-between">
							<dt class="text-gray-400">Status:</dt>
							<dd class="font-medium text-green-400">Active</dd>
						</div>
					</dl>
				</div>
				<div>
					<h3 class="mb-2 text-sm font-medium text-gray-400">Traffic Summary</h3>
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
			<span class="mb-4 block text-4xl">⚠️</span>
			<p class="text-gray-400">Unable to load statistics. Please try again later.</p>
		</div>
	{/if}
</div>
