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
			allowedIPsError =
				cidrsValidation.error || 'Invalid CIDR notation. Example: 10.0.0.2/32';
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
	class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-fade-in"
	onclick={handleOverlayClick}
	role="dialog"
	aria-modal="true"
	aria-labelledby="modal-title"
>
	<!-- Modal content -->
	<div class="glass-card p-6 w-full max-w-md animate-slide-up" onclick={(e) => e.stopPropagation()}>
		<!-- Header -->
		<div class="flex items-center justify-between mb-6">
			<h2 id="modal-title" class="text-2xl font-bold">Add New Peer</h2>
			<button
				onclick={onClose}
				class="text-gray-400 hover:text-white transition-colors"
				aria-label="Close modal"
			>
				<span class="text-2xl">✕</span>
			</button>
		</div>

		<!-- Form -->
		<form onsubmit={(e) => (e.preventDefault(), handleSubmit())}>
			<!-- Name field -->
			<div class="mb-4">
				<label for="peer-name" class="block text-sm font-medium mb-2">
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
					<p class="text-red-400 text-sm mt-1">{nameError}</p>
				{/if}
			</div>

			<!-- Allowed IPs field -->
			<div class="mb-6">
				<label for="allowed-ips" class="block text-sm font-medium mb-2">
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
				<p class="text-xs text-gray-400 mt-1">
					Enter one or more CIDR notations (comma or newline separated)
				</p>
				{#if allowedIPsError}
					<p class="text-red-400 text-sm mt-1">{allowedIPsError}</p>
					<p class="text-xs text-gray-400 mt-1">
						Examples: <code class="text-blue-400">10.0.0.5/32</code>,
						<code class="text-blue-400">192.168.1.0/24</code>
					</p>
				{/if}
			</div>

			<!-- Actions -->
			<div class="flex gap-3 justify-end">
				<button
					type="button"
					onclick={onClose}
					class="glass-btn-secondary px-6 py-2"
					disabled={loading}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="glass-btn-primary px-6 py-2"
					disabled={loading}
				>
					{#if loading}
						<span class="flex items-center gap-2">
							<span
								class="w-4 h-4 border-2 border-t-white border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin"
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

