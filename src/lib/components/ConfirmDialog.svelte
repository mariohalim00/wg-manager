<script lang="ts">
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
	class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-fade-in"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="confirm-dialog-title"
>
	<!-- Modal content -->
	<div
		class="glass-card p-6 w-full max-w-md animate-slide-up"
		onclick={(e) => e.stopPropagation()}
	>
		<!-- Header -->
		<div class="flex items-center gap-3 mb-4">
			{#if danger}
				<span class="text-3xl">⚠️</span>
			{:else}
				<span class="text-3xl">❓</span>
			{/if}
			<h2 id="confirm-dialog-title" class="text-xl font-bold">{title}</h2>
		</div>

		<!-- Message -->
		<p class="text-gray-300 mb-6">{message}</p>

		<!-- Actions -->
		<div class="flex gap-3 justify-end">
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
