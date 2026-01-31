<script lang="ts">
	import type { Peer } from '../types/peer';
	import StatusBadge from './StatusBadge.svelte';
	import { formatBytes, formatLastHandshake } from '../utils/formatting';

	type Props = {
		peers: Peer[];
		onDownloadConfig: (peer: Peer) => void;
		onRemove: (peer: Peer) => void;
	};

	let { peers, onDownloadConfig, onRemove }: Props = $props();
</script>

<div class="glass-card overflow-hidden">
	<div class="overflow-x-auto">
		<table class="w-full">
			<thead class="bg-glass-hover border-glass-border border-b">
				<tr class="text-left text-sm text-gray-300">
					<th class="px-6 py-4 font-semibold">Name</th>
					<th class="px-6 py-4 font-semibold">Status</th>
					<th class="px-6 py-4 font-semibold">Allowed IPs</th>
					<th class="hidden px-6 py-4 font-semibold md:table-cell">Last Handshake</th>
					<th class="hidden px-6 py-4 font-semibold lg:table-cell">Transfer</th>
					<th class="px-6 py-4 text-right font-semibold">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-glass-border divide-y">
				{#each peers as peer (peer.id)}
					<!-- FR-001a: Responsive action buttons (group for hover reveal ≥1024px) -->
					<tr class="group hover:bg-glass-hover transition-colors">
						<td class="px-6 py-4">
							<div>
								<p class="font-medium text-white">{peer.name}</p>
								<p class="max-w-xs truncate text-xs text-gray-400" title={peer.publicKey}>
									{peer.publicKey}
								</p>
							</div>
						</td>
						<td class="px-6 py-4">
							<StatusBadge status={peer.status} />
						</td>
						<td class="px-6 py-4">
							<div class="text-sm text-gray-300">
								{#each peer.allowedIPs as ip (ip)}
									<div>{ip}</div>
								{/each}
							</div>
						</td>
						<td class="hidden px-6 py-4 text-sm text-gray-400 md:table-cell">
							{formatLastHandshake(peer.lastHandshake)}
						</td>
						<td class="hidden px-6 py-4 lg:table-cell">
							<div class="text-sm">
								<div class="text-green-400">
									↓ {formatBytes(peer.receiveBytes)}
								</div>
								<div class="text-blue-400">
									↑ {formatBytes(peer.transmitBytes)}
								</div>
							</div>
						</td>
						<td class="px-6 py-4 text-right">
							<!-- FR-001a: Always visible <1024px, hover-reveal ≥1024px -->
							<div
								class="flex justify-end gap-2 opacity-100 transition-opacity duration-200 lg:opacity-0 lg:group-hover:opacity-100"
							>
								<button
									onclick={() => onDownloadConfig(peer)}
									class="glass-btn-secondary px-3 py-1 text-sm"
									title="Download config"
								>
									📥 Config
								</button>
								<button
									onclick={() => onRemove(peer)}
									class="glass-btn-secondary px-3 py-1 text-sm text-red-400 hover:text-red-300"
									title="Remove peer"
								>
									🗑️ Remove
								</button>
							</div>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="6" class="px-6 py-12 text-center text-gray-400">
							<div class="flex flex-col items-center gap-4">
								<span class="text-4xl">📭</span>
								<p>No peers found. Add your first peer to get started!</p>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
