<script lang="ts">
	// Optional size prop
	type Props = {
		size?: 'sm' | 'md' | 'lg';
	};

	let { size = 'md' }: Props = $props();

	// Size classes
	const sizeClasses = {
		sm: 'w-8 h-8 border-2',
		md: 'w-12 h-12 border-4',
		lg: 'w-16 h-16 border-4'
	};
</script>

<div class="flex items-center justify-center p-8">
	<div
		class="animate-spin rounded-full border-t-blue-400 border-r-transparent border-b-transparent border-l-transparent {sizeClasses[
			size
		]}"
		role="status"
		aria-label="Loading"
	></div>
</div>
