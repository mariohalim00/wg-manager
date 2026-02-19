<script lang="ts">
	import { peers } from '$lib/stores/peers';
	import { validateCIDRList } from '$lib/utils/validation';
	import type { Peer, PeerFormData, PeerCreateResponse, PeerUpdateRequest } from '$lib/types/peer';
	import { X, User, Server, Shield } from 'lucide-svelte';

	type Props = {
		mode?: 'add' | 'edit';
		peer?: Peer;
		onClose: () => void;
		onSuccess?: (response: PeerCreateResponse | Peer) => void;
	};

	let { mode = 'add', peer, onClose, onSuccess }: Props = $props();

	// Form state
	// Note: disabled warning because we want the form to be initialized with a value for the first ime only
	// svelte-ignore state_referenced_locally
		let name = $state(peer?.name || '');
	// svelte-ignore state_referenced_locally
	let allowedIPsInput = $state(peer?.allowedIPs.join('\n') || '');
	let dns = $state(''); // Loaded from settings if empty
	let mtu = $state(1420);
	let keepalive = $state(25);
	// svelte-ignore state_referenced_locally
	let interfaceAddress = $state(peer?.interfaceAddress || '');
	let preSharedKey = $state(false);
	let loading = $state(false);

	let isEdit = $derived(mode === 'edit');

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

		if (isEdit && peer) {
			const updateData: PeerUpdateRequest = {
				name: name.trim(),
				allowedIPs,
				interfaceAddress: interfaceAddress.trim() || undefined
			};
			const success = await peers.update(peer.id, updateData);
			loading = false;
			if (success) {
				if (onSuccess) {
					// For update, we might not have a full PeerCreateResponse but we have the updated info
					onSuccess({ ...peer, ...updateData });
				}
				onClose();
			}
		} else {
			const formData: PeerFormData = {
				name: name.trim(),
				allowedIPs,
				dns: dns.trim() || undefined,
				mtu: mtu > 0 ? mtu : undefined,
				persistentKeepalive: keepalive > 0 ? keepalive : undefined,
				preSharedKey,
				interfaceAddress: interfaceAddress.trim() || undefined
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
			<h2 id="modal-title" class="text-xl font-bold tracking-tight text-white">
				{isEdit ? 'Edit Peer' : 'Add New Peer'}
			</h2>
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
				{/if}
			</div>

			<!-- Advanced settings toggle -->
			<div class="mb-4">
				<details class="group">
					<summary
						class="flex cursor-pointer list-none items-center gap-2 text-sm font-medium text-slate-400 transition-colors hover:text-slate-300"
					>
						<span
							class="flex h-5 w-5 items-center justify-center rounded bg-white/5 transition-transform group-open:rotate-90"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="12"
								height="12"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="3"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<path d="m9 18 6-6-6-6" />
							</svg>
						</span>
						Advanced Configuration
					</summary>
					<div class="mt-4 space-y-4 rounded-lg bg-white/5 p-4">
						<!-- DNS -->
						<div>
							<label for="dns" class="mb-1 block text-xs font-semibold text-slate-500 uppercase">
								DNS Servers
							</label>
							<input
								id="dns"
								type="text"
								bind:value={dns}
								placeholder="8.8.8.8, 1.1.1.1"
								class="glass-input text-sm"
								disabled={loading}
							/>
						</div>

						<div class="grid grid-rows-2 gap-4">
							<!-- MTU -->
							<div>
								<label for="mtu" class="mb-1 block text-xs font-semibold text-slate-500 uppercase">
									MTU
								</label>
								<input
									id="mtu"
									type="number"
									bind:value={mtu}
									placeholder="1420"
									class="glass-input text-sm"
									disabled={loading}
								/>
							</div>
							<!-- Keepalive -->
							<div>
								<label
									for="keepalive"
									class="mb-1 block text-xs font-semibold text-slate-500 uppercase"
								>
									Keepalive
								</label>
								<input
									id="keepalive"
									type="number"
									bind:value={keepalive}
									placeholder="25"
									class="glass-input text-sm"
									disabled={loading}
								/>
							</div>
						</div>

						<!-- Interface Address -->
						<div>
							<label
								for="interface-address"
								class="mb-1 block text-xs font-semibold text-slate-500 uppercase"
							>
								Interface Address (Optional)
							</label>
							<input
								id="interface-address"
								type="text"
								bind:value={interfaceAddress}
								placeholder="e.g. 10.0.0.5/32"
								class="glass-input text-sm"
								disabled={loading}
							/>
							<p class="mt-1 text-[10px] text-slate-500">
								Manually set the [Interface] Address field in the client config.
							</p>
						</div>

						{#if !isEdit}
							<!-- PresharedKey -->
							<div class="flex items-center gap-3">
								<input
									id="psk"
									type="checkbox"
									bind:checked={preSharedKey}
									class="h-4 w-4 rounded border-white/10 bg-white/10 text-blue-500 focus:ring-0 focus:ring-offset-0"
									disabled={loading}
								/>
								<label for="psk" class="text-sm text-slate-300"> Generate Preshared Key </label>
							</div>
						{/if}
					</div>
				</details>
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
						{isEdit ? 'Saving...' : 'Adding...'}
					{:else}
						<Shield size={18} />
						{isEdit ? 'Save Changes' : 'Add Peer'}
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
