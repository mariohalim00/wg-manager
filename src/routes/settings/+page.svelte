<script lang="ts">
	import { onMount } from 'svelte';
	import { settings } from '$lib/stores/settings';
	import { Save, RefreshCw, Globe, Shield, Zap, Server } from 'lucide-svelte';

	let loading = $state(true);
	let saving = $state(false);

	// Local state for form fields
	let serverAddress = $state('');
	let dns = $state('');
	let mtu = $state(1420);
	let keepalive = $state(25);
	let endpoint = $state('');

	onMount(async () => {
		await settings.load();
		const current = $settings;
		if (current) {
			serverAddress = current.serverAddress;
			dns = current.dns;
			mtu = current.mtu || 1420;
			keepalive = current.keepalive || 25;
			endpoint = current.endpoint;
		}
		loading = false;
	});

	async function handleSave() {
		saving = true;
		await settings.save({
			serverAddress,
			dns,
			mtu,
			keepalive,
			endpoint
		});
		saving = false;
	}

	async function handleRefresh() {
		loading = true;
		await settings.load();
		const current = $settings;
		if (current) {
			serverAddress = current.serverAddress;
			dns = current.dns;
			mtu = current.mtu || 1420;
			keepalive = current.keepalive || 25;
			endpoint = current.endpoint;
		}
		loading = false;
	}
</script>

<svelte:head>
	<title>Settings - WireGuard Manager</title>
</svelte:head>

<div class="mx-auto max-w-4xl pb-12">
	<!-- Page header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="mb-1 text-3xl font-bold tracking-tight text-white">Settings</h1>
			<p class="text-slate-400">Configure global WireGuard parameters and server defaults</p>
		</div>
		<button
			onclick={handleRefresh}
			class="glass-btn-secondary flex items-center gap-2"
			disabled={loading || saving}
		>
			<RefreshCw size={18} class={loading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	{#if loading}
		<div class="flex h-64 items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-4 border-blue-500/20 border-t-blue-500"
			></div>
		</div>
	{:else}
		<div class="space-y-6">
			<!-- Server Configuration -->
			<div class="glass-card overflow-hidden">
				<div class="flex items-center gap-3 border-b border-white/5 bg-white/5 px-6 py-4">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-500/10 text-blue-400"
					>
						<Server size={20} />
					</div>
					<h2 class="text-xl font-bold text-white">Server Configuration</h2>
				</div>
				<div class="grid grid-cols-1 gap-6 p-6 md:grid-cols-2">
					<div class="space-y-4">
						<div>
							<label for="server-address" class="mb-1.5 block text-sm font-medium text-slate-300">
								Server Internal Address
							</label>
							<input
								id="server-address"
								type="text"
								bind:value={serverAddress}
								placeholder="10.0.0.1/24"
								class="glass-input w-full"
							/>
							<p class="mt-1.5 text-xs text-slate-500">
								The internal VPN IP address of this server.
							</p>
						</div>
						<div>
							<label for="endpoint" class="mb-1.5 block text-sm font-medium text-slate-300">
								Public Endpoint
							</label>
							<input
								id="endpoint"
								type="text"
								bind:value={endpoint}
								placeholder="vpn.example.com:51820"
								class="glass-input w-full"
							/>
							<p class="mt-1.5 text-xs text-slate-500">
								The public IP or hostname and port used by clients to connect.
							</p>
						</div>
					</div>
					<div class="rounded-xl border border-blue-500/10 bg-blue-500/5 p-4">
						<h4 class="mb-2 text-sm font-semibold text-blue-400">Info</h4>
						<p class="text-xs leading-relaxed text-slate-400">
							These values are used as defaults when generating client configurations. Changing them
							will not automatically update existing peer configurations, but new configurations
							will use these values.
						</p>
					</div>
				</div>
			</div>

			<!-- Client Defaults -->
			<div class="glass-card overflow-hidden">
				<div class="flex items-center gap-3 border-b border-white/5 bg-white/5 px-6 py-4">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-500/10 text-purple-400"
					>
						<Shield size={20} />
					</div>
					<h2 class="text-xl font-bold text-white">Client Defaults</h2>
				</div>
				<div class="grid grid-cols-1 gap-6 p-6 md:grid-cols-2">
					<div>
						<label for="dns-default" class="mb-1.5 block text-sm font-medium text-slate-300">
							Default DNS Servers
						</label>
						<div class="relative">
							<Globe class="absolute top-1/2 left-3 -translate-y-1/2 text-slate-400" size={16} />
							<input
								id="dns-default"
								type="text"
								bind:value={dns}
								placeholder="1.1.1.1, 8.8.8.8"
								class="glass-input w-full pl-10"
							/>
						</div>
						<p class="mt-1.5 text-xs text-slate-500">
							Comma-separated list of DNS servers for clients.
						</p>
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="mtu-default" class="mb-1.5 block text-sm font-medium text-slate-300">
								Default MTU
							</label>
							<div class="relative">
								<Zap class="absolute top-1/2 left-3 -translate-y-1/2 text-slate-400" size={16} />
								<input
									id="mtu-default"
									type="number"
									bind:value={mtu}
									class="glass-input w-full pl-10"
								/>
							</div>
						</div>
						<div>
							<label
								for="keepalive-default"
								class="mb-1.5 block text-sm font-medium text-slate-300"
							>
								Keepalive (sec)
							</label>
							<input
								id="keepalive-default"
								type="number"
								bind:value={keepalive}
								class="glass-input w-full"
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Save Actions -->
			<div class="flex justify-end gap-3 pt-4">
				<button
					onclick={handleSave}
					class="glass-btn-primary flex items-center gap-2 px-8 py-3"
					disabled={saving}
				>
					{#if saving}
						<span class="h-4 w-4 animate-spin rounded-full border-2 border-white/20 border-t-white"
						></span>
						Saving...
					{:else}
						<Save size={20} />
						Save All Settings
					{/if}
				</button>
			</div>
		</div>
	{/if}
</div>
