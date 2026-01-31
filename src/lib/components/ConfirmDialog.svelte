<script lang="ts">
	import { TriangleAlert, HelpCircle } from 'lucide-svelte';

	type Props = {
		title: string;
		message: string;
		confirmText?: string;
		cancelText?: string;
		danger?: boolean;
		onConfirm: () => void;
		onCancel: () => void;
	};

	let {
		title,
		message,
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		danger = false,
		onConfirm,
		onCancel
	}: Props = $props();

	// Handle overlay click to close
	function handleOverlayClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onCancel();
		}
	}
</script>

<!-- Modal overlay -->
<div
	class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="confirm-dialog-title"
>
	<!-- Modal content -->
	<div class="glass-card animate-slide-up w-full max-w-md p-6" onclick={(e) => e.stopPropagation()}>
		<!-- Header -->
		<div class="mb-4 flex items-center gap-3">
			{#if danger}
				<TriangleAlert class="text-red-500" size={28} />
			{:else}
				<HelpCircle class="text-blue-500" size={28} />
			{/if}
			<h2 id="confirm-dialog-title" class="text-xl font-bold">{title}</h2>
		</div>

		<!-- Message -->
		<p class="mb-6 text-gray-300">{message}</p>

		<!-- Actions -->
		<div class="flex justify-end gap-3">
			<button onclick={onCancel} class="glass-btn-secondary px-6 py-2">{cancelText}</button>
			<button
				onclick={onConfirm}
				class="glass-btn-primary px-6 py-2 {danger ? 'text-red-400 hover:text-red-300' : ''}"
			>
				{confirmText}
			</button>
		</div>
	</div>
</div>
