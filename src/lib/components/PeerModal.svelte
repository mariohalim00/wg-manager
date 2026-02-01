<script lang="ts">
	import { peers } from '$lib/stores/peers';
	import { validateCIDRList } from '$lib/utils/validation';
	import type { PeerFormData, PeerCreateResponse } from '$lib/types/peer';
	import { X, User, Server, Shield } from 'lucide-svelte';

	type Props = {
		onClose: () => void;
		onSuccess?: (response: PeerCreateResponse) => void;
	};

	let { onClose, onSuccess }: Props = $props();

	// Form state
	let name = $state('');
	let allowedIPsInput = $state('');
	let loading = $state(false);

	// Validation errors
	let nameError = $state('');
	let allowedIPsError = $state('');

	// Validate form
	function validateForm(): boolean {
		let valid = true;

		// Validate name
		if (!name.trim()) {
			nameError = 'Name is required';
			valid = false;
		} else {
			nameError = '';
		}

		// Validate allowed IPs (CIDR notation)
		const cidrsValidation = validateCIDRList(allowedIPsInput);
		if (!cidrsValidation.valid) {
			allowedIPsError = cidrsValidation.error || 'Invalid CIDR notation. Example: 10.0.0.2/32';
			valid = false;
		} else {
			allowedIPsError = '';
		}

		return valid;
	}

	// Handle form submit
	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		loading = true;

		// Parse allowed IPs into array
		const allowedIPs = allowedIPsInput
			.split(/[,\n]/)
			.map((s) => s.trim())
			.filter((s) => s.length > 0);

		const formData: PeerFormData = {
			name: name.trim(),
			allowedIPs
		};

		const response = await peers.add(formData);

		loading = false;

		if (response) {
			// Success - call onSuccess callback if provided, then close
			if (onSuccess) {
				onSuccess(response);
			}
			onClose();
		}
		// Error notification is handled by peers store
	}

	// Handle overlay click to close
	function handleOverlayClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	// Handle keyboard events for accessibility
	function handleOverlayKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	function handleContentKeydown(event: KeyboardEvent) {
		// Prevent propagation to avoid closing on Escape within content
		event.stopPropagation();
	}
</script>

<!-- Modal overlay -->
<div
	class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	onclick={handleOverlayClick}
	onkeydown={handleOverlayKeydown}
	role="dialog"
	aria-modal="true"
	aria-labelledby="modal-title"
	tabindex="-1"
>
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<!-- Modal content -->
	<div
		class="glass-card animate-slide-up w-full max-w-md overflow-hidden"
		onclick={(e) => e.stopPropagation()}
		onkeydown={handleContentKeydown}
		role="document"
	>
		<!-- Header -->
		<div class="flex items-center justify-between border-b border-white/5 bg-white/5 px-6 py-4">
			<h2 id="modal-title" class="text-xl font-bold tracking-tight text-white">Add New Peer</h2>
			<button
				onclick={onClose}
				class="rounded-lg p-1 text-slate-400 transition-colors hover:bg-white/10 hover:text-white"
				aria-label="Close modal"
			>
				<X size={20} />
			</button>
		</div>

		<!-- Form -->
		<form class="p-6" onsubmit={(e) => (e.preventDefault(), handleSubmit())}>
			<!-- Name field -->
			<div class="mb-5">
				<label for="peer-name" class="mb-2 block text-sm font-medium text-slate-300">
					Peer Name <span class="text-red-400">*</span>
				</label>
				<div class="relative">
					<User class="absolute top-1/2 left-3 -translate-y-1/2 text-slate-400" size={18} />
					<input
						id="peer-name"
						type="text"
						bind:value={name}
						placeholder="e.g., iPhone 15 Pro"
						class="glass-input w-full pl-10"
						disabled={loading}
						required
					/>
				</div>
				{#if nameError}
					<p class="mt-1 text-sm text-red-400">{nameError}</p>
				{/if}
			</div>

			<!-- Allowed IPs field -->
			<div class="mb-6">
				<label for="allowed-ips" class="mb-2 block text-sm font-medium text-slate-300">
					Allowed IPs (CIDR) <span class="text-red-400">*</span>
				</label>
				<div class="relative">
					<Server class="absolute top-3 left-3 text-slate-400" size={18} />
					<textarea
						id="allowed-ips"
						bind:value={allowedIPsInput}
						placeholder="10.0.0.2/32&#10;10.0.1.0/24"
						rows="3"
						class="glass-input w-full resize-none pl-10"
						disabled={loading}
						required
					></textarea>
				</div>
				<p class="mt-2 text-xs text-slate-400">
					Enter one or more CIDR notations (comma or newline separated).
				</p>
				{#if allowedIPsError}
					<p class="mt-1 text-sm text-red-400">{allowedIPsError}</p>
				{:else}
					<p class="mt-1 text-xs text-slate-500">
						Example: <code class="rounded bg-white/5 px-1 py-0.5 text-blue-400">10.0.0.5/32</code>
					</p>
				{/if}
			</div>

			<!-- Actions -->
			<div class="flex items-center justify-end gap-3 pt-2">
				<button type="button" onclick={onClose} class="glass-btn-secondary" disabled={loading}>
					Cancel
				</button>
				<button type="submit" class="glass-btn-primary flex items-center gap-2" disabled={loading}>
					{#if loading}
						<span class="h-4 w-4 animate-spin rounded-full border-2 border-white/20 border-t-white"
						></span>
						Adding...
					{:else}
						<Shield size={18} />
						Add Peer
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
