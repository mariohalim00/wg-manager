<script lang="ts">
	type Props = {
		title: string;
		value: string | number;
		icon: string;
		subtitle?: string;
		trend?: 'up' | 'down' | 'neutral';
		color?: 'blue' | 'green' | 'purple' | 'yellow';
	};

	let { title, value, icon, subtitle, trend, color = 'blue' }: Props = $props();

	// Color classes for icon background
	const colorClasses = {
		blue: 'bg-blue-500/20 text-blue-400',
		green: 'bg-green-500/20 text-green-400',
		purple: 'bg-purple-500/20 text-purple-400',
		yellow: 'bg-yellow-500/20 text-yellow-400'
	};

	// Trend indicators
	const trendIcons = {
		up: '↗',
		down: '↘',
		neutral: '→'
	};

	const trendColors = {
		up: 'text-green-400',
		down: 'text-red-400',
		neutral: 'text-gray-400'
	};
</script>

<div class="glass-card p-6 transition-transform duration-200 hover:scale-105">
	<!-- Icon -->
	<div class="mb-4 flex items-start justify-between">
		<div
			class="h-12 w-12 rounded-lg {colorClasses[color]} flex items-center justify-center text-2xl"
		>
			{icon}
		</div>
		{#if trend}
			<span class="text-xl {trendColors[trend]}">{trendIcons[trend]}</span>
		{/if}
	</div>

	<!-- Content -->
	<div>
		<p class="mb-1 text-sm text-gray-400">{title}</p>
		<p class="mb-1 text-3xl font-bold">{value}</p>
		{#if subtitle}
			<p class="text-xs text-gray-500">{subtitle}</p>
		{/if}
	</div>
</div>
