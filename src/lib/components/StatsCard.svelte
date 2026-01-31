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

<div class="glass-card p-6 hover:scale-105 transition-transform duration-200">
	<!-- Icon -->
	<div class="flex items-start justify-between mb-4">
		<div class="w-12 h-12 rounded-lg {colorClasses[color]} flex items-center justify-center text-2xl">
			{icon}
		</div>
		{#if trend}
			<span class="text-xl {trendColors[trend]}">{trendIcons[trend]}</span>
		{/if}
	</div>

	<!-- Content -->
	<div>
		<p class="text-sm text-gray-400 mb-1">{title}</p>
		<p class="text-3xl font-bold mb-1">{value}</p>
		{#if subtitle}
			<p class="text-xs text-gray-500">{subtitle}</p>
		{/if}
	</div>
</div>
