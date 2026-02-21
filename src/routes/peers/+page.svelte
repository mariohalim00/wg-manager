<script lang="ts">
	import { onMount } from 'svelte';
	import PeerTable from '$lib/components/PeerTable.svelte';
	import PeerModal from '$lib/components/PeerModal.svelte';
	import QRCodeDisplay from '$lib/components/QRCodeDisplay.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import { peers } from '$lib/stores/peers';
	import { addNotification } from '$lib/stores/notifications';
	import { generateWireGuardConfig, downloadConfigFile } from '$lib/utils/config';
	import type { Peer, PeerCreateResponse } from '$lib/types/peer';

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

	// Search state
	let searchQuery = $state('');

	const filteredPeers = $derived(
		$peers.filter(
			(p) =>
				p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.publicKey.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.allowedIPs.some((ip) => ip.includes(searchQuery))
		)
	);

	// Load peers on mount and start polling
	onMount(() => {
		(async () => {
			await peers.load();
			peers.startPolling();
			loading = false;
		})();

		return () => {
			peers.stopPolling();
		};
	});

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

	// Handle remove peer - show confirmation dialog
	function handleRemove(peer: Peer) {
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

	// Handle peer added successfully - show QR code modal
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
			// Update current peer reference if needed
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
</script>

<svelte:head>
	<title>Peers - WireGuard Manager</title>
</svelte:head>

<div class="mx-auto max-w-7xl">
	<!-- Page header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="mb-2 text-3xl font-bold">WireGuard Peers</h1>
			<p class="text-gray-400">Manage VPN clients and connections</p>
		</div>
	</div>

	<!-- Peer table or loading state -->
	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else}
		<PeerTable
			peers={filteredPeers}
			bind:searchQuery
			onAdd={handleAddPeer}
			onDownloadConfig={handleDownloadConfig}
			onRemove={handleRemove}
			onShowQR={handleShowDetails}
			onEdit={handleEdit}
		/>
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
		onClose={() => ((showEditModal = false), (peerToEdit = null))}
		onSuccess={() => ((showEditModal = false), (peerToEdit = null))}
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
