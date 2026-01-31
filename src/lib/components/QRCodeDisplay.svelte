<script lang="ts">
	import QRCode from 'svelte-qrcode';

	type Props = {
		config: string;
		peerName: string;
		onClose: () => void;
		onDownload: () => void;
	};

	let { config, peerName, onClose, onDownload }: Props = $props();

	// Handle overlay click to close
	function handleOverlayClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}
</script>

<!-- Modal overlay -->
<div
	class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-fade-in"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="qr-modal-title"
>
	<!-- Modal content -->
	<div
		class="glass-card p-6 w-full max-w-lg animate-slide-up"
		onclick={(e) => e.stopPropagation()}
	>
		<!-- Header -->
		<div class="flex items-center justify-between mb-6">
			<h2 id="qr-modal-title" class="text-2xl font-bold">Peer Configuration</h2>
			<button
				onclick={onClose}
				class="text-gray-400 hover:text-white transition-colors"
				aria-label="Close modal"
			>
				<span class="text-2xl">✕</span>
			</button>
		</div>

		<!-- Success message -->
		<div class="bg-green-500/20 border border-green-500/40 rounded-lg p-4 mb-6">
			<p class="text-green-400 text-sm">
				✓ Peer <strong>{peerName}</strong> added successfully!
			</p>
		</div>

		<!-- QR Code display -->
		<div class="flex flex-col items-center mb-6">
			<div class="bg-white p-4 rounded-lg mb-4">
				<QRCode value={config} size={256} />
			</div>
			<p class="text-sm text-gray-400 text-center">
				Scan with WireGuard mobile app to import configuration
			</p>
		</div>

		<!-- Security warning -->
		<div class="bg-yellow-500/20 border border-yellow-500/40 rounded-lg p-3 mb-6">
			<p class="text-yellow-400 text-xs">
				⚠️ <strong>Security Note:</strong> This configuration contains a private key. Save it now - it
				won't be shown again after closing this window.
			</p>
		</div>

		<!-- Actions -->
		<div class="flex flex-col gap-3">
			<button onclick={onDownload} class="glass-btn-primary px-6 py-3 w-full">
				📥 Download Config File
			</button>
			<button onclick={onClose} class="glass-btn-secondary px-6 py-3 w-full">Close</button>
		</div>

		<!-- Config preview (collapsible) -->
		<details class="mt-6">
			<summary class="cursor-pointer text-sm text-gray-400 hover:text-white transition-colors">
				View configuration text
			</summary>
			<pre class="mt-3 glass-card p-4 text-xs text-gray-300 overflow-x-auto">{config}</pre>
		</details>
	</div>
</div>

