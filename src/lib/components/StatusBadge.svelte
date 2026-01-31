<script lang="ts">
	type Props = {
		status: 'online' | 'offline';
	};

	let { status }: Props = $props();
</script>

{#if status === 'online'}
	<span
		class="inline-flex items-center gap-2 rounded-full border border-green-500/40 bg-green-500/20 px-3 py-1 text-sm font-medium text-green-400"
	>
		<span class="h-2 w-2 animate-pulse rounded-full bg-green-400"></span>
		Online
	</span>
{:else}
	<span
		class="inline-flex items-center gap-2 rounded-full border border-gray-500/40 bg-gray-500/20 px-3 py-1 text-sm font-medium text-gray-400"
	>
		<span class="h-2 w-2 rounded-full bg-gray-400"></span>
		Offline
	</span>
{/if}
