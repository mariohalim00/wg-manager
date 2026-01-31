<script lang="ts">
	type Props = {
		status: 'online' | 'offline';
	};

	let { status }: Props = $props();
</script>

{#if status === 'online'}
	<span
		class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium bg-green-500/20 text-green-400 border border-green-500/40"
	>
		<span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
		Online
	</span>
{:else}
	<span
		class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium bg-gray-500/20 text-gray-400 border border-gray-500/40"
	>
		<span class="w-2 h-2 rounded-full bg-gray-400"></span>
		Offline
	</span>
{/if}

