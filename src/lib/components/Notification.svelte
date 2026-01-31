<script lang="ts">
	import type { Notification } from '$lib/types/notification';
	import { notifications } from '$lib/stores/notifications';

	type Props = {
		notification: Notification;
	};

	let { notification }: Props = $props();

	// Type-based styling
	const typeClasses = {
		success: 'border-green-500 bg-green-500/20',
		error: 'border-red-500 bg-red-500/20',
		warning: 'border-yellow-500 bg-yellow-500/20',
		info: 'border-blue-500 bg-blue-500/20'
	};

	const typeIcons = {
		success: '✓',
		error: '✕',
		warning: '⚠',
		info: 'ℹ'
	};

	// Handle close
	function handleClose() {
		notifications.remove(notification.id);
	}
</script>

<div
	class="glass-card p-4 border-l-4 {typeClasses[
		notification.type
	]} animate-slide-up flex items-start gap-3"
	role="alert"
>
	<!-- Icon -->
	<span class="text-xl flex-shrink-0">{typeIcons[notification.type]}</span>

	<!-- Message -->
	<p class="flex-1 text-sm">{notification.message}</p>

	<!-- Close button -->
	<button
		onclick={handleClose}
		class="text-gray-400 hover:text-white transition-colors flex-shrink-0"
		aria-label="Close notification"
	>
		<span class="text-lg">✕</span>
	</button>
</div>

