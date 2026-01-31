<svelte:head>
	<title>Settings - WireGuard Manager</title>
</svelte:head>

<div class="mx-auto max-w-7xl">
	<!-- Page header -->
	<div class="mb-8">
		<h1 class="mb-2 text-3xl font-bold">Settings</h1>
		<p class="text-gray-400">Configure your WireGuard Manager preferences</p>
	</div>

	<!-- Coming Soon notice -->
	<div class="glass-card mb-8 p-8 text-center">
		<span class="mb-4 block text-6xl">🚧</span>
		<h2 class="mb-2 text-2xl font-bold">Coming Soon</h2>
		<p class="mb-6 text-gray-400">
			Settings functionality is under development. Check back later for configuration options!
		</p>
	</div>

	<!-- Mock settings panels (future features) -->
	<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
		<!-- Interface Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="mb-4 flex items-center gap-2 text-xl font-semibold">
				<span>🔧</span>
				Interface Settings
			</h3>
			<p class="mb-4 text-sm text-gray-400">Configure WireGuard interface parameters</p>
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Interface Name</span>
					<span class="text-sm text-gray-500">wg0</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Listen Port</span>
					<span class="text-sm text-gray-500">51820</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">MTU</span>
					<span class="text-sm text-gray-500">1420</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Server Address</span>
					<span class="text-sm text-gray-500">10.0.0.1/24</span>
				</div>
			</div>
		</div>

		<!-- Appearance Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="mb-4 flex items-center gap-2 text-xl font-semibold">
				<span>🎨</span>
				Appearance
			</h3>
			<p class="mb-4 text-sm text-gray-400">Customize the interface appearance</p>
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Theme</span>
					<span class="text-sm text-gray-500">Dark</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Glassmorphism</span>
					<span class="text-sm text-gray-500">Enabled</span>
				</div>
			</div>
		</div>

		<!-- Security Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="mb-4 flex items-center gap-2 text-xl font-semibold">
				<span>🔐</span>
				Security
			</h3>
			<p class="mb-4 text-sm text-gray-400">Manage authentication and access control</p>
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Authentication</span>
					<span class="text-sm text-gray-500">Not configured</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Session Timeout</span>
					<span class="text-sm text-gray-500">30 minutes</span>
				</div>
			</div>
		</div>

		<!-- Advanced Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="mb-4 flex items-center gap-2 text-xl font-semibold">
				<span>⚡</span>
				Advanced
			</h3>
			<p class="mb-4 text-sm text-gray-400">Advanced configuration options</p>
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Auto-refresh</span>
					<span class="text-sm text-gray-500">5 seconds</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-gray-300">Log Level</span>
					<span class="text-sm text-gray-500">Info</span>
				</div>
			</div>
		</div>
	</div>
</div>
