<script lang="ts">
	import { Activity, TrendingUp, TrendingDown, Minus } from 'lucide-svelte';
	import type { ComponentType } from 'svelte';

	type Props = {
		title: string;
		value: string | number;
		icon: ComponentType;
		subtitle?: string;
		trend?: 'up' | 'down' | 'neutral';
		color?: 'blue' | 'green' | 'purple' | 'yellow';
	};

	let { title, value, icon: IconComponent, subtitle, trend, color = 'blue' }: Props = $props();

	// Color classes for icon background
	const colorClasses = {
		blue: 'bg-blue-500/10 text-blue-400',
		green: 'bg-green-500/10 text-green-400',
		purple: 'bg-purple-500/10 text-purple-400',
		yellow: 'bg-yellow-500/10 text-yellow-400'
	};

	// Trend components
	const trendComponents = {
		up: TrendingUp,
		down: TrendingDown,
		neutral: Minus
	};

	const trendColors = {
		up: 'text-green-400',
		down: 'text-red-400',
		neutral: 'text-gray-400'
	};
</script>

<div class="glass-card group relative overflow-hidden p-6 transition-all hover:bg-white/5">
	<!-- Glow effect on hover -->
	<div
		class="absolute -top-4 -right-4 h-20 w-20 rounded-full bg-{color}-500/10 blur-2xl transition-all group-hover:bg-{color}-500/20"
	></div>

	<!-- Icon -->
	<div class="relative mb-4 flex items-start justify-between">
		<div class="rounded-lg {colorClasses[color]} flex h-12 w-12 items-center justify-center">
			<IconComponent class="h-6 w-6" />
		</div>
		{#if trend}
			{@const TrendIcon = trendComponents[trend]}
			<span class="flex items-center gap-1 text-sm font-semibold {trendColors[trend]}">
				<TrendIcon class="h-4 w-4" />
			</span>
		{/if}
	</div>

	<!-- Content -->
	<div class="relative">
		<p class="mb-1 text-sm font-medium text-slate-400">{title}</p>
		<p class="mb-1 text-3xl font-bold">{value}</p>
		{#if subtitle}
			<p class="text-xs text-gray-500">{subtitle}</p>
		{/if}
	</div>
</div>
