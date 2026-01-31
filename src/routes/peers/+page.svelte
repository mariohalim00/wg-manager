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
	import { Plus } from 'lucide-svelte';

	// Loading state
	let loading = $state(true);

	// Modal state
	let showAddModal = $state(false);
	let showQRModal = $state(false);
	let showConfirmDelete = $state(false);

	// QR modal data
	let qrConfig = $state('');
	let qrPeerName = $state('');

	// Delete confirmation data
	let peerToDelete: Peer | null = $state(null);

	// Load peers on mount
	onMount(async () => {
		await peers.load();
		loading = false;
	});

	// Handle download config (for existing peers - not yet supported by API)
	function handleDownloadConfig() {
		addNotification({
			type: 'warning',
			message: `Config download for existing peers not yet supported. Please save the config when adding a new peer.`,
			duration: 5000
		});
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
	function handlePeerAdded(response: PeerCreateResponse) {
		qrConfig = generateWireGuardConfig(response);
		qrPeerName = response.name;
		showQRModal = true;
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
		<button onclick={handleAddPeer} class="glass-btn-primary flex items-center gap-2 px-6 py-3 text-lg font-semibold">
			<Plus size={24} />
			Add Peer
		</button>
	</div>

	<!-- Peer table or loading state -->
	{#if loading}
		<div class="glass-card">
			<LoadingSpinner size="lg" />
		</div>
	{:else}
		<PeerTable peers={$peers} onDownloadConfig={handleDownloadConfig} onRemove={handleRemove} />
	{/if}
</div>

<!-- Add Peer Modal -->
{#if showAddModal}
	<PeerModal onClose={() => (showAddModal = false)} onSuccess={handlePeerAdded} />
{/if}

<!-- QR Code Modal -->
{#if showQRModal}
	<QRCodeDisplay
		config={qrConfig}
		peerName={qrPeerName}
		onClose={() => (showQRModal = false)}
		onDownload={handleConfigDownload}
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
