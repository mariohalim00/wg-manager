<script lang="ts">
	import { onMount } from 'svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import PeerTable from '$lib/components/PeerTable.svelte';
	import PeerModal from '$lib/components/PeerModal.svelte';
	import QRCodeDisplay from '$lib/components/QRCodeDisplay.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { stats } from '$lib/stores/stats';
	import { peers } from '$lib/stores/peers';
	import { addNotification } from '$lib/stores/notifications';
	import { formatBytes } from '$lib/utils/formatting';
	import { generateWireGuardConfig, downloadConfigFile } from '$lib/utils/config';
	import type { Peer, PeerCreateResponse } from '$lib/types/peer';
	import { Search, Bell, TrendingUp, Copy, Check, Plus } from 'lucide-svelte';

	// Loading state
	let loading = $state(true);

	// Modal state
	let showAddModal = $state(false);
	let showEditModal = $state(false);
	let showQRModal = $state(false);
	let showConfirmDelete = $state(false);

	// QR modal data
	let qrConfig = $state('');
	let qrPeerName = $state('');
	let qrAllowedIPs = $state<string[]>([]);
	let qrEndpoint = $state('');
	let qrPublicKey = $state('');
	let currentPeerForQR = $state<Peer | null>(null);

	// Edit/Delete confirmation data
	let peerToEdit: Peer | null = $state(null);
	let peerToDelete: Peer | null = $state(null);

	// Load data on mount
	onMount(() => {
		(async () => {
			await Promise.all([stats.load(), peers.load()]);
			peers.startPolling();
			loading = false;
		})();

		return () => {
			peers.stopPolling();
		};
	});

	// Derived: online peers count
	let onlinePeersCount = $derived($peers.filter((p) => p.status === 'online').length);

	const placeholderValue = '—';

	// Handle download config
	async function handleDownloadConfig(peer: Peer) {
		const config = await peers.getConfig(peer.id);

		if (!config) {
			return; // Error already handled by store
		}

		downloadConfigFile(config, peer.name);
		addNotification({
			type: 'success',
			message: `Configuration file downloaded for ${peer.name}`,
			duration: 3000
		});
	}

	// Handle show details/QR
	async function handleShowDetails(peer: Peer) {
		currentPeerForQR = peer;
		qrPeerName = peer.name;
		qrAllowedIPs = peer.allowedIPs;
		qrEndpoint = peer.endpoint || '';
		qrPublicKey = peer.publicKey;
		
		const config = await peers.getConfig(peer.id);
		qrConfig = config || '';
		
		showQRModal = true;
	}

	// Handle edit peer
	function handleEdit(peer: Peer) {
		peerToEdit = peer;
		showEditModal = true;
	}

	// Handle remove peer
	function handleRemovePeer(peer: Peer) {
		peerToDelete = peer;
		showConfirmDelete = true;
	}

	// Confirm peer removal
	async function confirmRemove() {
		if (peerToDelete) {
			await peers.remove(peerToDelete.id, peerToDelete.name);
			peerToDelete = null;
		}
		showConfirmDelete = false;
	}

	// Cancel peer removal
	function cancelRemove() {
		peerToDelete = null;
		showConfirmDelete = false;
	}

	// Handle add peer
	function handleAddPeer() {
		showAddModal = true;
	}

	// Handle peer added successfully
	function handlePeerAdded(response: PeerCreateResponse | Peer) {
		if ('config' in response && response.config) {
			qrConfig = response.config;
		} else {
			qrConfig = generateWireGuardConfig(response as PeerCreateResponse);
		}
		qrPeerName = response.name;
		qrAllowedIPs = response.allowedIPs;
		qrPublicKey = response.publicKey;
		showQRModal = true;
	}

	// Handle key regeneration
	async function handleRegenerate() {
		if (!currentPeerForQR) return;

		const response = await peers.regenerateKeys(currentPeerForQR.id);
		if (response) {
			qrConfig = response.config || generateWireGuardConfig(response);
			qrPublicKey = response.publicKey;
			const updated = $peers.find((p) => p.id === response.id);
			if (updated) currentPeerForQR = updated;
		}
	}

	// Handle config download from QR modal
	function handleConfigDownload() {
		downloadConfigFile(qrConfig, qrPeerName);
		addNotification({
			type: 'success',
			message: `Configuration file downloaded for ${qrPeerName}`,
			duration: 3000
		});
	}

	// Copied state
	let copied = $state(false);

	// Copy public key to clipboard
	function handleCopyPublicKey() {
		if (!$stats?.publicKey) return;
		try {
			navigator.clipboard.writeText($stats.publicKey);
			copied = true;

			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch (error) {
			console.error('Failed to copy public key:', error);
		}
	}
</script>

<svelte:head>
	<title>Dashboard - WireGuard Manager</title>
</svelte:head>

<div class="relative p-6">
	<!-- Header/Top Nav matching mockup -->
	<header
		class="sticky top-0 z-10 mb-8 flex flex-wrap items-center justify-between gap-4 border-b border-white/5 bg-[#101922]/40 px-0 py-5 backdrop-blur-md md:-mx-8 md:px-8"
	>
		<div class="flex flex-col gap-1">
			<p class="text-xs font-semibold tracking-[0.2em] text-slate-500 uppercase">
				WireGuard Interface
			</p>
			<h2 class="text-2xl font-semibold tracking-tight">
				{$stats?.interfaceName || 'wg0'}
				<span class="text-[#137fec]"> · Active</span>
			</h2>
		</div>
		<div class="hidden items-center gap-4 md:flex">
			<div class="relative">
				<Search class="absolute top-1/2 left-3 -translate-y-1/2 text-xl text-slate-400" size={20} />
				<input
					class="focus-ring w-60 rounded-xl border border-white/10 bg-white/5 py-2 pr-4 pl-10 text-sm text-white"
					placeholder="Search peers..."
					type="text"
					aria-label="Search peers"
				/>
			</div>
			<button
				class="focus-ring rounded-lg border border-white/10 bg-white/5 p-2 text-slate-400 transition-colors hover:text-white"
				aria-label="Notifications"
			>
				<Bell size={20} />
			</button>
			<div class="h-10 w-10 overflow-hidden rounded-full ring-2 ring-[#137fec]/20">
				<div
					class="flex h-full w-full items-center justify-center bg-[#137fec] text-sm font-semibold text-white"
					aria-label="User profile"
				>
					AD
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
			<!-- Status Cards -->
			<div class="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-4">
				<!-- Status Card -->
				<div class="metric-card relative overflow-hidden">
					<div
						class="absolute -top-8 -right-8 h-28 w-28 rounded-full bg-green-500/10 blur-3xl"
					></div>
					<p class="metric-label">Status</p>
					<div class="relative z-10 mt-3 flex items-center gap-2.5">
						<div
							class="pulse-online h-3 w-3 rounded-full bg-green-500 shadow-lg shadow-green-500/50"
						></div>
						<p class="metric-value text-2xl">Active</p>
					</div>
					<p class="metric-subtext mt-2">
						<span class="font-semibold text-white">{onlinePeersCount}</span> online peers
					</p>
				</div>

				<!-- Public Key Card -->
				<div class="metric-card">
					<p class="metric-label">Public Key</p>
					<div class="mt-3 flex items-center gap-2">
						<p class="metric-value font-mono text-lg leading-tight break-all">
							{$stats.publicKey
								? `${$stats.publicKey.slice(0, 8)}...${$stats.publicKey.slice(-5)}`
								: placeholderValue}
						</p>
						{#if $stats.publicKey}
							<button
								class="focus-ring shrink-0 rounded-lg border border-white/10 bg-white/5 p-1.5 text-slate-400 transition-colors hover:bg-white/10 hover:text-white"
								aria-label="Copy public key"
								onclick={handleCopyPublicKey}
							>
								{#if copied}
									<Check size={14} />
								{:else}
									<Copy size={14} />
								{/if}
							</button>
						{/if}
					</div>
					<p class="metric-subtext mt-2">Server identity</p>
				</div>

				<!-- Listening Port Card -->
				<div class="metric-card">
					<p class="metric-label">Listening Port</p>
					<p class="metric-value text-tabular mt-3 text-3xl">
						{$stats.listenPort && $stats.listenPort > 0 ? $stats.listenPort : placeholderValue}
					</p>
					<p class="metric-subtext mt-2">UDP port</p>
				</div>

				<!-- Subnet Card -->
				<div class="metric-card">
					<p class="metric-label">Subnet</p>
					<p class="metric-value text-tabular mt-3 font-mono text-2xl">
						{$stats.subnet ? $stats.subnet : placeholderValue}
					</p>
					<p class="metric-subtext mt-2">VPN address pool</p>
				</div>
			</div>

			<!-- Traffic Charts Section -->
			<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
				<!-- Total Received -->
				<div class="dashboard-surface overflow-hidden rounded-2xl">
					<div class="p-6 pb-4">
						<div class="mb-2 flex items-start justify-between">
							<div>
								<p class="mb-1 text-xs font-semibold tracking-wider text-slate-500 uppercase">
									Total Received
								</p>
								<h3 class="text-tabular text-4xl font-black tracking-tighter">
									{formatBytes($stats.totalRx)}
								</h3>
								<p class="mt-1 text-sm text-slate-400">
									{$stats.totalRx > 0 ? 'Total data received' : 'No data'}
								</p>
							</div>
							<div
								class="flex items-center gap-1.5 rounded-xl bg-green-500/10 px-3 py-1.5 text-sm font-bold text-green-400 ring-1 ring-green-500/20"
							>
								<TrendingUp size={16} />
								<span>12%</span>
							</div>
						</div>
					</div>
					<!-- SVG Chart Placeholder -->
					<div class="traffic-chart px-6">
						<svg class="h-full w-full" viewBox="0 0 500 150" preserveAspectRatio="none">
							<path
								d="M0 130 Q 50 120, 100 140 T 200 80 T 300 100 T 400 40 T 500 70"
								fill="none"
								stroke="#137fec"
								stroke-linecap="round"
								stroke-width="2.5"
							></path>
							<path
								d="M0 130 Q 50 120, 100 140 T 200 80 T 300 100 T 400 40 T 500 70 V 150 H 0 Z"
								fill="url(#grad-blue)"
								opacity="0.15"
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
				<div class="dashboard-surface overflow-hidden rounded-2xl">
					<div class="p-6 pb-4">
						<div class="mb-2 flex items-start justify-between">
							<div>
								<p class="mb-1 text-xs font-semibold tracking-wider text-slate-500 uppercase">
									Total Sent
								</p>
								<h3 class="text-tabular text-4xl font-black tracking-tighter">
									{formatBytes($stats.totalTx)}
								</h3>
								<p class="mt-1 text-sm text-slate-400">
									{$stats.totalTx > 0 ? formatBytes($stats.totalTx) : 'No data'}
								</p>
							</div>
							<div
								class="flex items-center gap-1.5 rounded-xl bg-green-500/10 px-3 py-1.5 text-sm font-bold text-green-400 ring-1 ring-green-500/20"
							>
								<TrendingUp size={16} />
								<span>5%</span>
							</div>
						</div>
					</div>
					<!-- SVG Chart Placeholder -->
					<div class="traffic-chart px-6">
						<svg class="h-full w-full" viewBox="0 0 500 150" preserveAspectRatio="none">
							<path
								d="M0 110 Q 70 120, 150 90 T 280 130 T 400 60 T 500 100"
								fill="none"
								stroke="#94a3b8"
								stroke-linecap="round"
								stroke-width="2.5"
							></path>
							<path
								d="M0 110 Q 70 120, 150 90 T 280 130 T 400 60 T 500 100 V 150 H 0 Z"
								fill="url(#grad-gray)"
								opacity="0.15"
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
				onShowQR={handleShowDetails}
				onEdit={handleEdit}
			/>
		</div>
	{:else}
		<div class="glass-card flex flex-col items-center justify-center p-12 text-center">
			<p class="text-gray-400">Unable to load dashboard data. Please try again later.</p>
		</div>
	{/if}
</div>

<!-- Add Peer Modal -->
{#if showAddModal}
	<PeerModal onClose={() => (showAddModal = false)} onSuccess={handlePeerAdded} />
{/if}

<!-- Edit Peer Modal -->
{#if showEditModal && peerToEdit}
	<PeerModal
		mode="edit"
		peer={peerToEdit}
		onClose={() => (showEditModal = false, peerToEdit = null)}
		onSuccess={() => (showEditModal = false, peerToEdit = null)}
	/>
{/if}

<!-- QR Code Modal / Details -->
{#if showQRModal}
	<QRCodeDisplay
		config={qrConfig}
		peerName={qrPeerName}
		allowedIPs={qrAllowedIPs}
		endpoint={qrEndpoint}
		publicKey={qrPublicKey}
		onClose={() => (showQRModal = false)}
		onDownload={handleConfigDownload}
		onRegenerate={handleRegenerate}
	/>
{/if}

<!-- Confirm Delete Modal -->
{#if showConfirmDelete && peerToDelete}
	<ConfirmDialog
		title="Remove Peer"
		message="Are you sure you want to remove '{peerToDelete.name}'? This action cannot be undone."
		confirmText="Remove"
		cancelText="Cancel"
		danger={true}
		onConfirm={confirmRemove}
		onCancel={cancelRemove}
	/>
{/if}

<!-- Background Decorative Elements (matching mockup) -->
<div
	class="pointer-events-none fixed top-[-10%] left-[-10%] z-[-1] h-[40%] w-[40%] rounded-full bg-[#137fec]/20 blur-[120px]"
></div>
<div
	class="pointer-events-none fixed right-[-10%] bottom-[-10%] z-[-1] h-[30%] w-[30%] rounded-full bg-blue-900/10 blur-[100px]"
></div>
