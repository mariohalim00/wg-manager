<script lang="ts">
	// Optional size prop
	type Props = {
		size?: 'sm' | 'md' | 'lg';
	};

	let { size = 'md' }: Props = $props();

	// Size classes
	const sizeClasses = {
		sm: 'h-6 w-6 border-2',
		md: 'h-10 w-10 border-[3px]',
		lg: 'h-16 w-16 border-4'
	};
</script>

<div class="flex items-center justify-center p-4">
	<div
		class="animate-spin rounded-full border-white/20 border-t-[#137fec] {sizeClasses[size]}"
		role="status"
		aria-label="Loading"
	></div>
</div>
