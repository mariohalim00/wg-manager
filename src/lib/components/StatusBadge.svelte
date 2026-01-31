<script lang="ts">
	type Props = {
		status: 'online' | 'offline';
	};

	let { status }: Props = $props();
</script>

{#if status === 'online'}
	<span
		class="inline-flex items-center gap-1.5 rounded-lg bg-green-400/10 px-2 py-1 text-xs font-bold text-green-400"
	>
		<span class="pulse-online h-1.5 w-1.5 rounded-full bg-green-400"></span>
		Active
	</span>
{:else}
	<span
		class="inline-flex items-center gap-1.5 rounded-lg bg-slate-400/10 px-2 py-1 text-xs font-bold text-slate-400"
	>
		<span class="h-1.5 w-1.5 rounded-full bg-slate-400"></span>
		Inactive
	</span>
{/if}
