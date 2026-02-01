<script lang="ts">
	import { peers } from '$lib/stores/peers';
	import { getQrUrl } from '$lib/api/peers';
	import {
		X,
		Copy,
		Download,
		Globe,
		Server,
		Shield,
		Check,
		RefreshCw,
		ExternalLink
	} from 'lucide-svelte';
	import { fade, scale, fly } from 'svelte/transition';

	type Props = {
		config: string;
		peerName: string;
		allowedIPs?: string[];
		endpoint?: string;
		publicKey?: string;
		onClose: () => void;
		onDownload: () => void;
		onRegenerate?: () => void;
	};

	let {
		config,
		peerName,
		allowedIPs = [],
		endpoint = '',
		publicKey = '',
		onClose,
		onDownload,
		onRegenerate
	}: Props = $props();

	let copying = $state(false);
	let regenerating = $state(false);
	let qrTimestamp = $state(Date.now());

	async function handleRegenerate() {
		if (!onRegenerate) return;
		regenerating = true;
		try {
			await onRegenerate();
			qrTimestamp = Date.now();
		} finally {
			regenerating = false;
		}
	}

	// Handle overlay click to close
	function handleOverlayClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	// Copy config to clipboard
	function copyConfig() {
		navigator.clipboard.writeText(config);
	}

	// Copy public key to clipboard
	function copyPublicKey() {
		if (publicKey) {
			navigator.clipboard.writeText(publicKey);
		}
	}
</script>

<!-- Modal overlay matching mockup design -->
<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 px-4 backdrop-blur-sm"
	onclick={handleOverlayClick}
	onkeydown={(e) => e.key === 'Escape' && onClose()}
	role="dialog"
	aria-modal="true"
	aria-labelledby="qr-modal-title"
	tabindex="-1"
>
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<!-- Modal content - glass-panel style matching mockup -->
	<div
		class="glass-panel animate-slide-up flex w-full max-w-[540px] flex-col overflow-hidden rounded-2xl border border-white/10 shadow-[0_32px_64px_-16px_rgba(0,0,0,0.6)]"
		onclick={(e) => e.stopPropagation()}
		onkeydown={() => {}}
		role="document"
	>
		<!-- Header -->
		<div class="flex items-center justify-between px-10 pt-10 pb-4">
			<div class="flex flex-col">
				<div class="flex items-center gap-2">
					<h1 id="qr-modal-title" class="text-2xl font-bold tracking-tight text-white">
						{peerName}
					</h1>
					<span class="flex h-2 w-2 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]"
					></span>
				</div>
				<p class="text-sm font-normal text-[#92adc9]/60">Peer connection details and setup</p>
			</div>
			<button
				onclick={onClose}
				class="rounded-lg bg-white/5 p-2 text-white/40 transition-colors hover:bg-white/10 hover:text-white"
				aria-label="Close modal"
			>
				<X size={20} />
			</button>
		</div>

		<!-- QR Code Section -->
		<div class="flex flex-col items-center px-10 py-6">
			<div class="relative rounded-2xl border border-[#137fec]/40 p-1">
				<div
					class="qr-gradient-container relative flex h-80 w-80 items-center justify-center overflow-hidden rounded-xl"
				>
					<div class="relative z-10 rounded-lg bg-white p-6 shadow-[0_20px_30px_rgba(0,0,0,0.4)]">
						<img src="{getQrUrl(publicKey)}?t={qrTimestamp}" alt="WireGuard QR Code" class="h-[200px] w-[200px]" />
					</div>
				</div>
			</div>
			<p
				class="mt-8 max-w-[360px] text-center text-[13px] leading-relaxed font-normal text-[#92adc9]/70"
			>
				Scan this QR code with the WireGuard app on your mobile device to import the configuration
				instantly.
			</p>
		</div>

		<!-- Action Buttons -->
		<div class="flex gap-3 px-10 py-4">
			<button
				onclick={copyConfig}
				class="flex h-12 flex-1 items-center justify-center gap-3 rounded-xl bg-[#137fec] text-[15px] font-bold text-white shadow-lg shadow-[#137fec]/20 transition-all hover:bg-[#137fec]/90 active:scale-[0.98]"
			>
				<Copy size={20} />
				<span>Copy Config</span>
			</button>
			<button
				onclick={onDownload}
				class="flex h-12 w-12 items-center justify-center rounded-xl border border-white/5 bg-[#1e2a37] text-white/80 transition-all hover:bg-[#253444] hover:text-white"
				title="Download .conf file"
			>
				<Download size={20} />
			</button>
		</div>

		<!-- Connection Details -->
		<div class="px-10 pt-4 pb-10">
			<div class="space-y-5 rounded-xl border border-white/5 bg-black/30 p-6">
				<!-- Allowed IPs -->
				<div class="flex flex-col gap-1.5">
					<span class="text-[10px] font-bold tracking-[0.1em] text-[#92adc9]/50 uppercase"
						>Allowed IPs</span
					>
					<code class="font-mono text-[14px] text-[#137fec]"
						>{allowedIPs.join(', ') || '0.0.0.0/0, ::/0'}</code
					>
				</div>

				<!-- Endpoint & Keepalive -->
				<div class="grid grid-cols-2 gap-8">
					<div class="flex flex-col gap-1.5">
						<span class="text-[10px] font-bold tracking-[0.1em] text-[#92adc9]/50 uppercase"
							>Endpoint</span
						>
						<code class="font-mono text-[13px] text-white/90"
							>{endpoint || 'vpn.example.com:51820'}</code
						>
					</div>
					<div class="flex flex-col gap-1.5">
						<span class="text-[10px] font-bold tracking-[0.1em] text-[#92adc9]/50 uppercase"
							>Keepalive</span
						>
						<code class="font-mono text-[13px] text-white/90">25 seconds</code>
					</div>
				</div>

				<!-- Public Key -->
				{#if publicKey}
					<div class="flex flex-col gap-1.5 pt-1">
						<div class="flex items-center justify-between">
							<span class="text-[10px] font-bold tracking-[0.1em] text-[#92adc9]/50 uppercase"
								>Public Key</span
							>
							<button
								onclick={copyPublicKey}
								class="text-[16px] text-white/30 transition-colors hover:text-white"
							>
								<Copy size={16} />
							</button>
						</div>
						<code class="truncate font-mono text-[12px] tracking-tight text-white/50"
							>{publicKey}</code
						>
					</div>
				{/if}
			</div>

			<!-- Regenerate Keys Button -->
			<div class="mt-8 flex justify-center">
				<button
					onclick={handleRegenerate}
					disabled={regenerating}
					class="flex items-center gap-2 rounded-lg px-4 py-2 text-[11px] font-bold tracking-[0.15em] text-red-400/60 uppercase transition-all hover:bg-red-400/5 hover:text-red-400 disabled:opacity-50"
				>
					<RefreshCw size={18} class={regenerating ? 'animate-spin' : ''} />
					{regenerating ? 'Regenerating...' : 'Regenerate Keys'}
				</button>
			</div>
		</div>
	</div>
</div>

<!-- Footer hint -->
<div class="fixed bottom-8 left-1/2 z-50 -translate-x-1/2">
	<p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">
		Esc to Close • WireGuard Secure Tunnel
	</p>
</div>
