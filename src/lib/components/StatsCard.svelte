<script lang="ts">
	import { TrendingUp, TrendingDown, Minus } from 'lucide-svelte';
	import type { ComponentType } from 'svelte';

	type Props = {
		title: string;
		value: string | number;
		icon?: ComponentType;
		subtitle?: string;
		trend?: 'up' | 'down' | 'neutral';
		trendValue?: string;
		color?: 'blue' | 'green' | 'purple' | 'yellow';
	};

	let {
		title,
		value,
		icon: IconComponent,
		subtitle,
		trend,
		trendValue,
		color = 'blue'
	}: Props = $props();

	// Color classes matching mockup design
	const glowColors = {
		blue: 'bg-[#137fec]/10 group-hover:bg-[#137fec]/20',
		green: 'bg-green-500/10 group-hover:bg-green-500/20',
		purple: 'bg-purple-500/10 group-hover:bg-purple-500/20',
		yellow: 'bg-yellow-500/10 group-hover:bg-yellow-500/20'
	};

	const iconColors = {
		blue: 'bg-[#137fec]/20 text-[#137fec]',
		green: 'bg-green-500/20 text-green-400',
		purple: 'bg-purple-500/20 text-purple-400',
		yellow: 'bg-yellow-500/20 text-yellow-400'
	};

	// Trend components and colors
	const trendComponents = {
		up: TrendingUp,
		down: TrendingDown,
		neutral: Minus
	};

	const trendColors = {
		up: 'text-green-400 bg-green-400/10',
		down: 'text-red-400 bg-red-400/10',
		neutral: 'text-gray-400 bg-gray-400/10'
	};
</script>

<div class="glass group relative flex flex-col gap-2 overflow-hidden rounded-2xl p-6">
	<!-- Glow effect on hover -->
	<div
		class="absolute -top-4 -right-4 h-20 w-20 rounded-full {glowColors[
			color
		]} blur-2xl transition-all"
	></div>

	<!-- Header with optional icon and trend -->
	<div class="relative flex items-start justify-between">
		{#if IconComponent}
			<div class="flex h-10 w-10 items-center justify-center rounded-lg {iconColors[color]}">
				<IconComponent class="h-5 w-5" />
			</div>
		{/if}
		{#if trend && trendValue}
			{@const TrendIcon = trendComponents[trend]}
			<div
				class="flex items-center gap-1 rounded-lg px-2 py-1 text-sm font-bold {trendColors[trend]}"
			>
				<TrendIcon class="h-3 w-3" />
				{trendValue}
			</div>
		{/if}
	</div>

	<!-- Content -->
	<div class="relative mt-2">
		<p class="text-sm font-medium text-slate-400">{title}</p>
		<p class="text-2xl font-bold tracking-tight">{value}</p>
		{#if subtitle}
			<p class="mt-1 text-xs text-slate-500">{subtitle}</p>
		{/if}
	</div>
</div>
