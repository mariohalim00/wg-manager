<script lang="ts">
	import { onMount } from 'svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import PeerTable from '$lib/components/PeerTable.svelte';
	import { stats } from '$lib/stores/stats';
	import { peers } from '$lib/stores/peers';
	import { formatBytes } from '$lib/utils/formatting';
	import { Search, Bell, TrendingUp, Users } from 'lucide-svelte';

	// Loading state
	let loading = $state(true);

	// Load data on mount
	onMount(async () => {
		await Promise.all([stats.load(), peers.load()]);
		loading = false;
	});

	// Derived: online peers count
	let onlinePeersCount = $derived($peers.filter((p) => p.status === 'online').length);

	// Placeholder handlers for peer table
	function handleDownloadConfig(peer: any) {
		console.log('Download config for', peer.name);
	}

	function handleRemovePeer(peer: any) {
		console.log('Remove peer', peer.name);
	}
</script>

<svelte:head>
	<title>Dashboard - WireGuard Manager</title>
</svelte:head>

<div class="relative">
	<!-- Header/Top Nav matching mockup -->
	<header class="sticky top-0 z-10 mb-8 flex items-center justify-between border-b border-white/5 bg-[#101922]/40 px-0 py-6 backdrop-blur-md md:-mx-8 md:px-8">
		<div class="flex items-center gap-4">
			<h2 class="text-2xl font-black tracking-tight">
				Interface: <span class="text-[#137fec]">{$stats?.interfaceName || 'wg0'}</span>
			</h2>
		</div>
		<div class="hidden items-center gap-6 md:flex">
			<div class="relative">
				<Search class="absolute top-1/2 left-3 -translate-y-1/2 text-xl text-slate-400" size={20} />
				<input
					class="w-64 rounded-xl border border-white/10 bg-white/5 py-2 pr-4 pl-10 text-sm text-white focus:border-[#137fec] focus:ring-[#137fec]"
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
		<div class="space-y-8">
			<!-- Interface Stats Grid (matching mockup design) -->
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
				<!-- Status Card -->
				<div class="glass group relative flex flex-col gap-2 overflow-hidden rounded-2xl p-6">
					<div class="absolute -top-4 -right-4 h-20 w-20 rounded-full bg-green-500/10 blur-2xl transition-all group-hover:bg-green-500/20"></div>
					<p class="text-sm font-medium text-slate-400">Status</p>
					<div class="flex items-center gap-2">
						<div class="pulse-online h-2.5 w-2.5 rounded-full bg-green-500"></div>
						<p class="text-2xl font-bold tracking-tight">Active</p>
					</div>
				</div>

				<!-- Interface Name Card -->
				<div class="glass relative flex flex-col gap-2 overflow-hidden rounded-2xl p-6">
					<p class="text-sm font-medium text-slate-400">Interface</p>
					<p class="font-mono text-2xl font-bold tracking-tight">{$stats.interfaceName}</p>
				</div>

				<!-- Peer Count Card -->
				<div class="glass group relative flex flex-col gap-2 overflow-hidden rounded-2xl p-6">
					<div class="absolute -top-4 -right-4 h-20 w-20 rounded-full bg-[#137fec]/10 blur-2xl transition-all group-hover:bg-[#137fec]/20"></div>
					<p class="text-sm font-medium text-slate-400">Total Peers</p>
					<div class="flex items-center gap-3">
						<p class="text-2xl font-bold tracking-tight">{$stats.peerCount}</p>
						<div class="flex items-center gap-1 text-sm text-green-400">
							<Users size={14} />
							<span>{onlinePeersCount} online</span>
						</div>
					</div>
				</div>

				<!-- Online Peers Card -->
				<div class="glass relative flex flex-col gap-2 overflow-hidden rounded-2xl p-6">
					<p class="text-sm font-medium text-slate-400">Online Peers</p>
					<p class="text-2xl font-bold tracking-tight">{onlinePeersCount} / {$stats.peerCount}</p>
				</div>
			</div>

			<!-- Traffic Charts Section -->
			<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
				<!-- Total Received -->
				<div class="glass rounded-2xl p-6">
					<div class="mb-6 flex items-start justify-between">
						<div>
							<p class="text-sm font-medium text-slate-400">Total Received</p>
							<h3 class="text-3xl font-black tracking-tight">
								{formatBytes($stats.totalRx)}
							</h3>
						</div>
						<div class="flex items-center gap-1 rounded-lg bg-green-400/10 px-2 py-1 text-sm font-bold text-green-400">
							<TrendingUp size={14} />
							12%
						</div>
					</div>
					<!-- SVG Chart Placeholder -->
					<div class="h-32 w-full">
						<svg class="h-full w-full" viewBox="0 0 500 150">
							<path
								d="M0 130 Q 50 120, 100 140 T 200 80 T 300 100 T 400 40 T 500 70"
								fill="none"
								stroke="#137fec"
								stroke-linecap="round"
								stroke-width="3"
							></path>
							<path
								d="M0 130 Q 50 120, 100 140 T 200 80 T 300 100 T 400 40 T 500 70 V 150 H 0 Z"
								fill="url(#grad-blue)"
								opacity="0.1"
							></path>
							<defs>
								<linearGradient id="grad-blue" x1="0%" x2="0%" y1="0%" y2="100%">
									<stop offset="0%" style="stop-color: #137fec; stop-opacity: 1"></stop>
									<stop offset="100%" style="stop-color: #137fec; stop-opacity: 0"></stop>
								</linearGradient>
							</defs>
						</svg>
					</div>
				</div>

				<!-- Total Sent -->
				<div class="glass rounded-2xl p-6">
					<div class="mb-6 flex items-start justify-between">
						<div>
							<p class="text-sm font-medium text-slate-400">Total Sent</p>
							<h3 class="text-3xl font-black tracking-tight">
								{formatBytes($stats.totalTx)}
							</h3>
						</div>
						<div class="flex items-center gap-1 rounded-lg bg-green-400/10 px-2 py-1 text-sm font-bold text-green-400">
							<TrendingUp size={14} />
							5%
						</div>
					</div>
					<!-- SVG Chart Placeholder -->
					<div class="h-32 w-full">
						<svg class="h-full w-full" viewBox="0 0 500 150">
							<path
								d="M0 110 Q 70 120, 150 90 T 280 130 T 400 60 T 500 100"
								fill="none"
								stroke="#94a3b8"
								stroke-linecap="round"
								stroke-width="3"
							></path>
							<path
								d="M0 110 Q 70 120, 150 90 T 280 130 T 400 60 T 500 100 V 150 H 0 Z"
								fill="url(#grad-gray)"
								opacity="0.1"
							></path>
							<defs>
								<linearGradient id="grad-gray" x1="0%" x2="0%" y1="0%" y2="100%">
									<stop offset="0%" style="stop-color: #94a3b8; stop-opacity: 1"></stop>
									<stop offset="100%" style="stop-color: #94a3b8; stop-opacity: 0"></stop>
								</linearGradient>
							</defs>
						</svg>
					</div>
				</div>
			</div>

			<!-- Peers Table -->
			<PeerTable
				peers={$peers}
				onDownloadConfig={handleDownloadConfig}
				onRemove={handleRemovePeer}
			/>
		</div>
	{:else}
		<div class="glass-card flex flex-col items-center justify-center p-12 text-center">
			<p class="text-gray-400">Unable to load dashboard data. Please try again later.</p>
		</div>
	{/if}
</div>

<!-- Background Decorative Elements (matching mockup) -->
<div class="pointer-events-none fixed top-[-10%] left-[-10%] z-[-1] h-[40%] w-[40%] rounded-full bg-[#137fec]/20 blur-[120px]"></div>
<div class="pointer-events-none fixed right-[-10%] bottom-[-10%] z-[-1] h-[30%] w-[30%] rounded-full bg-blue-900/10 blur-[100px]"></div>
