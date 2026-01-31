<script lang="ts">
	import { peers } from '../stores/peers';
	import { validateCIDRList } from '../utils/validation';
	import type { PeerFormData, PeerCreateResponse } from '../types/peer';

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
</script>

<!-- Modal overlay -->
<div
	class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="modal-title"
>
	<!-- Modal content -->
	<div class="glass-card animate-slide-up w-full max-w-md p-6" onclick={(e) => e.stopPropagation()}>
		<!-- Header -->
		<div class="mb-6 flex items-center justify-between">
			<h2 id="modal-title" class="text-2xl font-bold">Add New Peer</h2>
			<button
				onclick={onClose}
				class="text-gray-400 transition-colors hover:text-white"
				aria-label="Close modal"
			>
				<span class="text-2xl">✕</span>
			</button>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => (e.preventDefault(), handleSubmit())}>
			<!-- Name field -->
			<div class="mb-4">
				<label for="peer-name" class="mb-2 block text-sm font-medium">
					Peer Name <span class="text-red-400">*</span>
				</label>
				<input
					id="peer-name"
					type="text"
					bind:value={name}
					placeholder="e.g., John's Laptop"
					class="glass-input w-full"
					disabled={loading}
					required
				/>
				{#if nameError}
					<p class="mt-1 text-sm text-red-400">{nameError}</p>
				{/if}
			</div>

			<!-- Allowed IPs field -->
			<div class="mb-6">
				<label for="allowed-ips" class="mb-2 block text-sm font-medium">
					Allowed IPs (CIDR) <span class="text-red-400">*</span>
				</label>
				<textarea
					id="allowed-ips"
					bind:value={allowedIPsInput}
					placeholder="10.0.0.2/32&#10;10.0.1.0/24"
					rows="3"
					class="glass-input w-full resize-none"
					disabled={loading}
					required
				></textarea>
				<p class="mt-1 text-xs text-gray-400">
					Enter one or more CIDR notations (comma or newline separated)
				</p>
				{#if allowedIPsError}
					<p class="mt-1 text-sm text-red-400">{allowedIPsError}</p>
					<p class="mt-1 text-xs text-gray-400">
						Examples: <code class="text-blue-400">10.0.0.5/32</code>,
						<code class="text-blue-400">192.168.1.0/24</code>
					</p>
				{/if}
			</div>

			<!-- Actions -->
			<div class="flex justify-end gap-3">
				<button
					type="button"
					onclick={onClose}
					class="glass-btn-secondary px-6 py-2"
					disabled={loading}
				>
					Cancel
				</button>
				<button type="submit" class="glass-btn-primary px-6 py-2" disabled={loading}>
					{#if loading}
						<span class="flex items-center gap-2">
							<span
								class="h-4 w-4 animate-spin rounded-full border-2 border-t-white border-r-transparent border-b-transparent border-l-transparent"
							></span>
							Adding...
						</span>
					{:else}
						Add Peer
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
