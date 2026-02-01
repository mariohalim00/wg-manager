<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { stats } from '$lib/stores/stats';
	import { peers } from '$lib/stores/peers';
	import PeerTable from '$lib/components/PeerTable.svelte';
	import PeerModal from '$lib/components/PeerModal.svelte';
	import QRCodeDisplay from '$lib/components/QRCodeDisplay.svelte';
	import Notification from '$lib/components/Notification.svelte';
	import { notifications } from '$lib/stores/notifications';
	import { getConfigUrl } from '$lib/api/peers';
	import type { Peer, PeerCreateResponse } from '$lib/types/peer';
	import type { StatsHistoryItem } from '$lib/types/stats';
	import {
		Users,
		Zap,
		Activity,
		Plus,
		ArrowUpRight,
		ArrowDownLeft,
		Search,
		X
	} from 'lucide-svelte';
	import Chart from 'chart.js/auto';

	let pollingInterval: number;
	let showAddModal = $state(false);
	let showDetailsModal = $state(false);
	let selectedPeer = $state<Peer | null>(null);
	let searchQuery = $state('');

	const history = stats.history;

	// Chart refs
	let rxChartCanvas: HTMLCanvasElement | undefined = $state();
	let txChartCanvas: HTMLCanvasElement | undefined = $state();
	let rxChart: Chart | undefined;
	let txChart: Chart | undefined;

	onMount(async () => {
		await Promise.all([stats.load(), peers.load(), stats.loadHistory()]);

		// Initialize charts
		if (rxChartCanvas && txChartCanvas) {
			initCharts();
		}

		pollingInterval = window.setInterval(async () => {
			await Promise.all([stats.load(), peers.load(), stats.loadHistory()]);
			updateCharts();
		}, 60000); // 1 minute polling
	});

	onDestroy(() => {
		if (pollingInterval) clearInterval(pollingInterval);
		if (rxChart) rxChart.destroy();
		if (txChart) txChart.destroy();
	});

	function initCharts() {
		const historyData = $history;
		const labels = historyData.map((h: StatsHistoryItem) =>
			new Date(h.timestamp * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
		);

		const rxData = historyData.map((h: StatsHistoryItem) => h.totalRx / (1024 * 1024)); // MB
		const txData = historyData.map((h: StatsHistoryItem) => h.totalTx / (1024 * 1024)); // MB

		const chartOptions = {
			responsive: true,
			maintainAspectRatio: false,
			plugins: { legend: { display: false } },
			scales: {
				x: { display: false },
				y: {
					display: true,
					grid: { color: 'rgba(255, 255, 255, 0.05)' },
					ticks: { color: '#64748b', font: { size: 10 } }
				}
			},
			elements: {
				line: { tension: 0.4 },
				point: { radius: 0 }
			}
		};

		if (rxChartCanvas) {
			rxChart = new Chart(rxChartCanvas, {
				type: 'line',
				data: {
					labels,
					datasets: [
						{
							data: rxData,
							borderColor: '#3b82f6',
							backgroundColor: 'rgba(59, 130, 246, 0.1)',
							fill: true,
							borderWidth: 2
						}
					]
				},
				options: chartOptions as any
			});
		}

		if (txChartCanvas) {
			txChart = new Chart(txChartCanvas, {
				type: 'line',
				data: {
					labels,
					datasets: [
						{
							data: txData,
							borderColor: '#a855f7',
							backgroundColor: 'rgba(168, 85, 247, 0.1)',
							fill: true,
							borderWidth: 2
						}
					]
				},
				options: chartOptions as any
			});
		}
	}

	function updateCharts() {
		const historyData = $history;
		const labels = historyData.map((h: StatsHistoryItem) =>
			new Date(h.timestamp * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
		);

		if (rxChart) {
			rxChart.data.labels = labels;
			rxChart.data.datasets[0].data = historyData.map(
				(h: StatsHistoryItem) => h.totalRx / (1024 * 1024)
			);
			rxChart.update('none');
		}

		if (txChart) {
			txChart.data.labels = labels;
			txChart.data.datasets[0].data = historyData.map(
				(h: StatsHistoryItem) => h.totalTx / (1024 * 1024)
			);
			txChart.update('none');
		}
	}

	const filteredPeers = $derived(
		$peers.filter(
			(p) =>
				p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.publicKey.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.allowedIPs.some((ip) => ip.includes(searchQuery))
		)
	);

	function formatBytes(bytes: number) {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function handleDownloadConfig(peer: Peer) {
		const url = getConfigUrl(peer.id);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${peer.name || peer.id}.conf`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
	}

	function handleShowDetails(peer: Peer) {
		selectedPeer = peer;
		showDetailsModal = true;
	}

	function handleEditPeer(peer: Peer) {
		selectedPeer = peer;
		showAddModal = true;
	}

	function handleRemovePeer(peer: Peer) {
		if (confirm(`Are you sure you want to remove peer "${peer.name || peer.publicKey}"?`)) {
			peers.remove(peer.id, peer.name);
		}
	}

	async function handleRegenerate() {
		if (!selectedPeer) return;
		const response = await peers.regenerateKeys(selectedPeer.id);
		if (response) {
			selectedPeer = { ...selectedPeer, ...response };
		}
	}

	function handleConfigDownload() {
		if (selectedPeer && selectedPeer.config) {
			handleDownloadConfig(selectedPeer);
		}
	}

	function handleAddSuccess(response: PeerCreateResponse | Peer) {
		if ('config' in response) {
			selectedPeer = response as any;
			showDetailsModal = true;
		}
	}
</script>

<svelte:head>
	<title>Dashboard - WireGuard Manager</title>
</svelte:head>

<div class="animate-fade-in mx-auto max-w-7xl">
	<!-- Summary Stats -->
	<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
		<div class="glass-card flex items-center p-6 transition-all hover:scale-[1.02]">
			<div
				class="mr-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-blue-500/10 text-blue-400"
			>
				<Users size={24} />
			</div>
			<div>
				<p class="text-sm font-medium text-slate-400">Total Peers</p>
				<p class="text-2xl font-bold text-white">{$stats?.peerCount || 0}</p>
			</div>
		</div>

		<div class="glass-card flex items-center p-6 transition-all hover:scale-[1.02]">
			<div
				class="mr-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-green-500/10 text-green-400"
			>
				<Activity size={24} />
			</div>
			<div>
				<p class="text-sm font-medium text-slate-400">Active Now</p>
				<p class="text-2xl font-bold text-white">
					{$peers.filter((p) => p.status === 'online').length}
				</p>
			</div>
		</div>

		<div class="glass-card flex items-center p-6 transition-all hover:scale-[1.02]">
			<div
				class="mr-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-purple-500/10 text-purple-400"
			>
				<ArrowDownLeft size={24} />
			</div>
			<div>
				<p class="text-sm font-medium text-slate-400">Total Received</p>
				<p class="text-2xl font-bold text-white">{formatBytes($stats?.totalRx || 0)}</p>
			</div>
		</div>

		<div class="glass-card flex items-center p-6 transition-all hover:scale-[1.02]">
			<div
				class="mr-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-orange-500/10 text-orange-400"
			>
				<ArrowUpRight size={24} />
			</div>
			<div>
				<p class="text-sm font-medium text-slate-400">Total Transmitted</p>
				<p class="text-2xl font-bold text-white">{formatBytes($stats?.totalTx || 0)}</p>
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
		<!-- Main Content - Peers Table -->
		<div class="lg:col-span-2">
			<div class="glass-card mb-8 overflow-hidden">
				<div
					class="flex flex-col gap-4 border-b border-white/5 bg-white/5 px-6 py-4 md:flex-row md:items-center md:justify-between"
				>
					<h2 class="text-xl font-bold text-white">WireGuard Peers</h2>
					<div class="flex items-center gap-3">
						<div class="relative">
							<Search class="absolute top-1/2 left-3 -translate-y-1/2 text-slate-400" size={16} />
							<input
								type="text"
								bind:value={searchQuery}
								placeholder="Search peers..."
								class="glass-input w-48 pl-10 text-sm"
							/>
						</div>
						<button
							onclick={() => {
								selectedPeer = null;
								showAddModal = true;
							}}
							class="glass-btn-primary flex h-10 items-center gap-2 px-4 text-sm"
						>
							<Plus size={18} />
							Add Peer
						</button>
					</div>
				</div>

				<div class="p-0">
					<PeerTable
						peers={filteredPeers}
						onEdit={handleEditPeer}
						onRemove={handleRemovePeer}
						onDownloadConfig={handleDownloadConfig}
						onShowQR={handleShowDetails}
					/>
					{#if filteredPeers.length === 0}
						<div class="flex flex-col items-center justify-center py-12 text-slate-500">
							<p class="mb-2">No peers found</p>
							{#if searchQuery}
								<button
									onclick={() => (searchQuery = '')}
									class="text-sm text-blue-400 hover:underline"
								>
									Clear search
								</button>
							{/if}
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Sidebar - Interface Stats & Charts -->
		<div class="space-y-8">
			<div class="glass-card p-6">
				<h3 class="mb-4 flex items-center gap-2 font-bold text-white">
					<Zap size={18} class="text-blue-400" />
					Interface: {$stats?.interfaceName || 'wg0'}
				</h3>
				<div class="space-y-4">
					<div class="flex flex-col gap-1">
						<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
							>Public Key</span
						>
						<span class="truncate rounded bg-white/5 p-2 font-mono text-xs text-slate-300">
							{$stats?.publicKey || 'Loading...'}
						</span>
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div class="flex flex-col gap-1">
							<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase">Port</span
							>
							<span class="text-sm font-medium text-white">{$stats?.listenPort || '---'}</span>
						</div>
						<div class="flex flex-col gap-1">
							<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
								>Subnet</span
							>
							<span class="text-sm font-medium text-white">{$stats?.subnet || '---'}</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Traffic Chart - RX -->
			<div class="glass-card overflow-hidden">
				<div class="flex items-center justify-between border-b border-white/5 bg-white/5 px-6 py-4">
					<span class="text-sm font-bold text-white">Incoming (MB)</span>
					<ArrowDownLeft size={16} class="text-blue-400" />
				</div>
				<div class="h-40 p-4">
					<canvas bind:this={rxChartCanvas}></canvas>
				</div>
			</div>

			<!-- Traffic Chart - TX -->
			<div class="glass-card overflow-hidden">
				<div class="flex items-center justify-between border-b border-white/5 bg-white/5 px-6 py-4">
					<span class="text-sm font-bold text-white">Outgoing (MB)</span>
					<ArrowUpRight size={16} class="text-purple-400" />
				</div>
				<div class="h-40 p-4">
					<canvas bind:this={txChartCanvas}></canvas>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Modals -->
{#if showAddModal}
	<PeerModal
		mode={selectedPeer ? 'edit' : 'add'}
		peer={selectedPeer || undefined}
		onClose={() => {
			showAddModal = false;
			selectedPeer = null;
		}}
		onSuccess={handleAddSuccess}
	/>
{/if}

{#if showDetailsModal && selectedPeer}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
		onclick={() => (showDetailsModal = false)}
		onkeydown={(e) => e.key === 'Escape' && (showDetailsModal = false)}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div
			class="glass-card w-full max-w-lg overflow-hidden"
			onclick={(e) => e.stopPropagation()}
			role="document"
		>
			<div class="flex items-center justify-between border-b border-white/5 bg-white/5 px-6 py-4">
				<h3 class="text-xl font-bold text-white">Peer Configuration</h3>
				<button
					onclick={() => (showDetailsModal = false)}
					class="rounded-lg p-1 text-slate-400 hover:bg-white/10 hover:text-white"
				>
					<X size={20} />
				</button>
			</div>
			<div class="p-6">
				<QRCodeDisplay
					config={selectedPeer.config || ''}
					peerName={selectedPeer.name}
					allowedIPs={selectedPeer.allowedIPs}
					endpoint={selectedPeer.endpoint}
					publicKey={selectedPeer.publicKey}
					onClose={() => (showDetailsModal = false)}
					onDownload={handleConfigDownload}
					onRegenerate={handleRegenerate}
				/>
				<div class="mt-8 flex justify-end">
					<button onclick={() => (showDetailsModal = false)} class="glass-btn-primary px-8">
						Done
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Notifications -->
<div class="fixed right-6 bottom-6 z-[100] flex flex-col gap-3">
	{#each $notifications as notification (notification.id)}
		<Notification {notification} />
	{/each}
</div>
