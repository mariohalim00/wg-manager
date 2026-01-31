<script lang="ts">
	import type { Peer } from '$lib/types';
	import StatusBadge from './StatusBadge.svelte';
	import { formatBytes, formatLastHandshake } from '$lib/utils/formatting';
	import { Download, Upload, Trash2, Copy, Search, Filter } from 'lucide-svelte';

	type Props = {
		peers: Peer[];
		onDownloadConfig: (peer: Peer) => void;
		onRemove: (peer: Peer) => void;
	};

	let { peers, onDownloadConfig, onRemove }: Props = $props();
</script>

<div class="glass mb-12 overflow-hidden rounded-2xl">
	<!-- Header with Filter -->
	<div class="flex items-center justify-between border-b border-white/5 p-6">
		<h3 class="text-lg font-bold">Active Peers ({peers.length})</h3>
		<div class="flex gap-2">
			<button
				class="gap-2 rounded-xl bg-white/5 px-4 py-2 text-sm font-medium transition-all hover:bg-white/10 flex items-center"
			>
				<Filter size={16} />
				Filter
			</button>
		</div>
	</div>

	<!-- Table -->
	<div class="overflow-x-auto">
		<table class="w-full">
			<thead>
				<tr class="border-b border-white/5 text-left text-xs font-semibold text-slate-400 uppercase">
					<th class="px-6 py-4">Peer Name / Public Key</th>
					<th class="px-6 py-4">Status</th>
					<th class="px-6 py-4">Allowed IPs</th>
					<th class="px-6 py-4">Transfer</th>
					<th class="px-6 py-4">Last Handshake</th>
					<th class="px-6 py-4 text-right">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-white/5">
				{#each peers as peer (peer.id)}
					<tr class="group hover:bg-white/5 transition-colors">
						<td class="px-6 py-4">
							<div class="flex flex-col">
								<span class="font-bold text-white max-w-[200px] truncate" title={peer.name}
									>{peer.name || 'Unnamed Peer'}</span
								>
								<div class="flex items-center gap-1 text-slate-400 text-xs">
									<span class="font-mono max-w-[120px] truncate">{peer.publicKey}</span>
									<button class="hover:text-white" title="Copy Public Key">
										<Copy size={12} />
									</button>
								</div>
							</div>
						</td>
						<td class="px-6 py-4">
							<StatusBadge status={peer.status} />
						</td>
						<td class="px-6 py-4">
							<div class="flex flex-wrap gap-1">
								{#each peer.allowedIPs as ip (ip)}
									<span class="rounded bg-white/5 px-1.5 py-0.5 text-xs font-mono text-slate-300"
										>{ip}</span
									>
								{/each}
							</div>
						</td>
						<td class="px-6 py-4">
							<div class="flex flex-col text-xs font-medium">
								<span class="text-green-400 flex items-center gap-1">
									<Download size={12} /> {formatBytes(peer.receiveBytes)}
								</span>
								<span class="text-blue-400 flex items-center gap-1">
									<Upload size={12} /> {formatBytes(peer.transmitBytes)}
								</span>
							</div>
						</td>
						<td class="px-6 py-4">
							<span class="text-sm text-slate-400 font-mono"
								>{formatLastHandshake(peer.lastHandshake)}</span
							>
						</td>
						<td class="px-6 py-4 text-right">
							<div class="flex justify-end gap-2 text-slate-400 opacity-60 group-hover:opacity-100 transition-opacity">
								<button
									onclick={() => onDownloadConfig(peer)}
									class="rounded-lg p-2 hover:bg-white/10 hover:text-white transition-colors"
									title="Download Configuration"
								>
									<Download size={18} />
								</button>
								<button
									onclick={() => onRemove(peer)}
									class="rounded-lg p-2 hover:bg-red-500/20 hover:text-red-400 transition-colors"
									title="Remove Peer"
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
								<div class="p-4 rounded-full bg-white/5">
									<Search size={32} />
								</div>
								<div>
									<p class="font-medium text-white">No peers found</p>
									<p class="text-sm mt-1">Add a new peer to get started</p>
								</div>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
