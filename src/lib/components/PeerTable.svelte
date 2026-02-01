<script lang="ts">
	import type { Peer } from '$lib/types';
	import { formatBytes, formatLastHandshake } from '$lib/utils/formatting';
	import {
		Download,
		Trash2,
		QrCode,
		ArrowUp,
		ArrowDown,
		Laptop,
		Smartphone,
		Router,
		Monitor,
		Search
	} from 'lucide-svelte';

	type Props = {
		peers: Peer[];
		onDownloadConfig: (peer: Peer) => void;
		onRemove: (peer: Peer) => void;
		onShowQR?: (peer: Peer) => void;
	};

	let { peers, onDownloadConfig, onRemove, onShowQR }: Props = $props();

	// Get device icon based on peer name (simple heuristic)
	function getDeviceIcon(name: string) {
		const lowerName = name.toLowerCase();
		if (
			lowerName.includes('iphone') ||
			lowerName.includes('android') ||
			lowerName.includes('pixel') ||
			lowerName.includes('phone')
		) {
			return Smartphone;
		}
		if (lowerName.includes('mac') || lowerName.includes('laptop') || lowerName.includes('book')) {
			return Laptop;
		}
		if (
			lowerName.includes('router') ||
			lowerName.includes('gateway') ||
			lowerName.includes('openwrt')
		) {
			return Router;
		}
		return Monitor;
	}
</script>

<div class="dashboard-surface mb-12 overflow-hidden rounded-2xl">
	<!-- Header -->
	<div class="flex items-center justify-between border-b border-white/5 px-6 py-5">
		<div class="flex items-center gap-2">
			<h3 class="text-xl font-bold tracking-tight">Active Peers</h3>
			<span class="rounded-full bg-[#137fec]/10 px-2.5 py-0.5 text-sm font-bold text-[#137fec]">
				{peers.length}
			</span>
		</div>
		<div class="flex gap-3">
			<button
				class="focus-ring rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold transition-all hover:border-white/20 hover:bg-white/10"
			>
				Filter
			</button>
			<button
				class="focus-ring rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold transition-all hover:border-white/20 hover:bg-white/10"
			>
				Export
			</button>
		</div>
	</div>

	<!-- Table -->
	<div class="overflow-x-auto">
		<table class="w-full text-left">
			<thead>
				<tr class="text-[11px] font-bold tracking-widest text-slate-500 uppercase">
					<th class="px-6 py-4">Status</th>
					<th class="px-6 py-4">Peer Name</th>
					<th class="px-6 py-4">Internal IP</th>
					<th class="px-6 py-4">Transfer (U/D)</th>
					<th class="px-6 py-4">Last Handshake</th>
					<th class="px-6 py-4 text-right">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-white/5">
				{#each peers as peer (peer.id)}
					{@const DeviceIcon = getDeviceIcon(peer.name)}
					{@const isOnline = peer.status === 'online'}
					<tr class="group transition-colors hover:bg-white/[0.02]">
						<!-- Status -->
						<td class="px-6 py-5">
							<div class="flex items-center gap-2">
								{#if isOnline}
									<div class="pulse-online h-2 w-2 rounded-full bg-green-500"></div>
									<span class="text-xs font-bold tracking-tighter text-green-500 uppercase"
										>Online</span
									>
								{:else}
									<div class="h-2 w-2 rounded-full bg-slate-600"></div>
									<span class="text-xs font-bold tracking-tighter text-slate-500 uppercase"
										>Offline</span
									>
								{/if}
							</div>
						</td>

						<!-- Peer Name with Icon -->
						<td class="px-6 py-5">
							<div class="flex items-center gap-3">
								<div
									class="flex h-8 w-8 items-center justify-center rounded-lg {isOnline
										? 'bg-[#137fec]/20 text-[#137fec]'
										: 'bg-white/10 text-slate-400'}"
								>
									<DeviceIcon size={18} />
								</div>
								<div>
									<p class="text-sm font-bold">{peer.name || 'Unnamed Peer'}</p>
									<p class="text-xs text-slate-500">{peer.publicKey.slice(0, 8)}...</p>
								</div>
							</div>
						</td>

						<!-- Internal IP -->
						<td class="px-6 py-5">
							{#if peer.allowedIPs.length > 0}
								<span
									class="text-tabular rounded bg-white/5 px-2 py-1 font-mono text-xs text-slate-300"
								>
									{peer.allowedIPs[0].replace('/32', '')}
								</span>
							{:else}
								<span class="text-slate-500">—</span>
							{/if}
						</td>

						<!-- Transfer -->
						<td class="px-6 py-5">
							<div class="flex flex-col gap-1">
								<div
									class="flex items-center gap-1.5 text-xs {isOnline
										? 'font-medium text-green-500'
										: 'text-slate-500'}"
								>
									<ArrowUp size={14} />
									<span class="text-tabular">{formatBytes(peer.transmitBytes)}</span>
								</div>
								<div
									class="flex items-center gap-1.5 text-xs {isOnline
										? 'font-medium text-[#137fec]'
										: 'text-slate-500'}"
								>
									<ArrowDown size={14} />
									<span class="text-tabular">{formatBytes(peer.receiveBytes)}</span>
								</div>
							</div>
						</td>

						<!-- Last Handshake -->
						<td class="px-6 py-5">
							<span class="text-xs text-slate-400">{formatLastHandshake(peer.lastHandshake)}</span>
						</td>

						<!-- Actions -->
						<td class="px-6 py-5">
							<div
								class="flex items-center justify-end gap-2 opacity-0 transition-opacity group-hover:opacity-100"
							>
								<button
									onclick={() => onDownloadConfig(peer)}
									class="focus-ring rounded-lg p-2 text-slate-400 transition-colors hover:bg-white/5 hover:text-white"
									title="Download Config"
									aria-label="Download Config"
								>
									<Download size={18} />
								</button>
								{#if onShowQR}
									<button
										onclick={() => onShowQR(peer)}
										class="focus-ring rounded-lg p-2 text-slate-400 transition-colors hover:bg-white/5 hover:text-white"
										title="View QR Code"
									>
										<QrCode size={18} />
									</button>
								{/if}
								<button
									onclick={() => onRemove(peer)}
									class="focus-ring rounded-lg p-2 text-slate-400 transition-colors hover:bg-red-400/5 hover:text-red-400"
									title="Delete Peer"
								>
									<Trash2 size={18} />
								</button>
							</div>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="6" class="px-6 py-12 text-center text-slate-400">
							<div class="flex flex-col items-center gap-4">
								<div class="rounded-full bg-white/5 p-4">
									<Search size={32} />
								</div>
								<div>
									<p class="font-medium text-white">No peers found</p>
									<p class="mt-1 text-sm">Add a new peer to get started</p>
								</div>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Footer -->
	{#if peers.length > 0}
		<div class="flex items-center justify-center border-t border-white/5 p-4">
			<a
				href="/peers"
				class="px-6 py-2 text-sm font-bold text-slate-400 transition-colors hover:text-white"
			>
				View All Peers
			</a>
		</div>
	{/if}
</div>
