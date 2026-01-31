<script lang="ts">
	import QRCode from 'svelte-qrcode';
	import { TriangleAlert, Download, X, CheckCircle2 } from 'lucide-svelte';

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
	class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="qr-modal-title"
>
	<!-- Modal content -->
	<div class="glass-card animate-slide-up w-full max-w-lg p-6" onclick={(e) => e.stopPropagation()}>
		<!-- Header -->
		<div class="mb-6 flex items-center justify-between">
			<h2 id="qr-modal-title" class="text-2xl font-bold">Peer Configuration</h2>
			<button
				onclick={onClose}
				class="text-gray-400 transition-colors hover:text-white"
				aria-label="Close modal"
			>
				<X size={24} />
			</button>
		</div>

		<!-- Success message -->
		<div class="mb-6 flex items-center gap-2 rounded-lg border border-green-500/40 bg-green-500/20 p-4">
			<CheckCircle2 class="text-green-400 flex-shrink-0" size={20} />
			<p class="text-sm text-green-400">
				Peer <strong>{peerName}</strong> added successfully!
			</p>
		</div>

		<!-- QR Code display -->
		<div class="mb-6 flex flex-col items-center">
			<div class="mb-4 rounded-lg bg-white p-4">
				<QRCode value={config} size={256} />
			</div>
			<p class="text-center text-sm text-gray-400">
				Scan with WireGuard mobile app to import configuration
			</p>
		</div>

		<!-- Security warning -->
		<div class="mb-6 flex items-start gap-2 rounded-lg border border-yellow-500/40 bg-yellow-500/20 p-3">
			<TriangleAlert class="text-yellow-400 flex-shrink-0 mt-0.5" size={16} />
			<p class="text-xs text-yellow-400">
				<strong>Security Note:</strong> This configuration contains a private key. Save it now - it
				won't be shown again after closing this window.
			</p>
		</div>

		<!-- Actions -->
		<div class="flex flex-col gap-3">
			<button onclick={onDownload} class="glass-btn-primary flex items-center justify-center gap-2 w-full px-6 py-3">
				<Download size={20} /> Download Config File
			</button>
			<button onclick={onClose} class="glass-btn-secondary w-full px-6 py-3">Close</button>
		</div>

		<!-- Config preview (collapsible) -->
		<details class="mt-6">
			<summary class="cursor-pointer text-sm text-gray-400 transition-colors hover:text-white">
				View configuration text
			</summary>
			<pre class="glass-card mt-3 overflow-x-auto p-4 text-xs text-gray-300">{config}</pre>
		</details>
	</div>
</div>
