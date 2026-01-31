<script lang="ts">
	import type { Notification } from '$lib/types/notification';
	import { notifications } from '$lib/stores/notifications';
	import { CheckCircle2, XCircle, AlertTriangle, Info, X } from 'lucide-svelte';

	type Props = {
		notification: Notification;
	};

	let { notification }: Props = $props();

	// Type-based styling
	const typeClasses = {
		success: 'border-green-500 bg-green-500/20 text-green-400',
		error: 'border-red-500 bg-red-500/20 text-red-400',
		warning: 'border-yellow-500 bg-yellow-500/20 text-yellow-400',
		info: 'border-blue-500 bg-blue-500/20 text-blue-400'
	};

	const typeIcons = {
		success: CheckCircle2,
		error: XCircle,
		warning: AlertTriangle,
		info: Info
	};

	let Icon = $derived(typeIcons[notification.type]);

	// Handle close
	function handleClose() {
		notifications.remove(notification.id);
	}
</script>

<div
	class="glass-card border-l-4 p-4 {typeClasses[
		notification.type
	]} animate-slide-up flex items-start gap-3"
	role="alert"
>
	<!-- Icon -->
	<div class="shrink-0 mt-0.5">
		<Icon size={20} />
	</div>

	<!-- Message -->
	<p class="flex-1 text-sm text-white/90">{notification.message}</p>

	<!-- Close button -->
	<button
		onclick={handleClose}
		class="shrink-0 text-white/50 transition-colors hover:text-white"
		aria-label="Close notification"
	>
		<X size={18} />
	</button>
</div>
