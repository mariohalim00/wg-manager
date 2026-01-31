<svelte:head>
	<title>Settings - WireGuard Manager</title>
</svelte:head>

<div class="max-w-7xl mx-auto">
	<!-- Page header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold mb-2">Settings</h1>
		<p class="text-gray-400">Configure your WireGuard Manager preferences</p>
	</div>

	<!-- Coming Soon notice -->
	<div class="glass-card p-8 text-center mb-8">
		<span class="text-6xl mb-4 block">🚧</span>
		<h2 class="text-2xl font-bold mb-2">Coming Soon</h2>
		<p class="text-gray-400 mb-6">
			Settings functionality is under development. Check back later for configuration options!
		</p>
	</div>

	<!-- Mock settings panels (future features) -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<!-- Interface Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<span>🔧</span>
				Interface Settings
			</h3>
			<p class="text-sm text-gray-400 mb-4">Configure WireGuard interface parameters</p>
			<div class="space-y-3">
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Interface Name</span>
					<span class="text-sm text-gray-500">wg0</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Listen Port</span>
					<span class="text-sm text-gray-500">51820</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">MTU</span>
					<span class="text-sm text-gray-500">1420</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Server Address</span>
					<span class="text-sm text-gray-500">10.0.0.1/24</span>
				</div>
			</div>
		</div>

		<!-- Appearance Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<span>🎨</span>
				Appearance
			</h3>
			<p class="text-sm text-gray-400 mb-4">Customize the interface appearance</p>
			<div class="space-y-3">
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Theme</span>
					<span class="text-sm text-gray-500">Dark</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Glassmorphism</span>
					<span class="text-sm text-gray-500">Enabled</span>
				</div>
			</div>
		</div>

		<!-- Security Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<span>🔐</span>
				Security
			</h3>
			<p class="text-sm text-gray-400 mb-4">Manage authentication and access control</p>
			<div class="space-y-3">
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Authentication</span>
					<span class="text-sm text-gray-500">Not configured</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Session Timeout</span>
					<span class="text-sm text-gray-500">30 minutes</span>
				</div>
			</div>
		</div>

		<!-- Advanced Settings (mock) -->
		<div class="glass-card p-6 opacity-50">
			<h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
				<span>⚡</span>
				Advanced
			</h3>
			<p class="text-sm text-gray-400 mb-4">Advanced configuration options</p>
			<div class="space-y-3">
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Auto-refresh</span>
					<span class="text-sm text-gray-500">5 seconds</span>
				</div>
				<div class="flex justify-between items-center">
					<span class="text-sm text-gray-300">Log Level</span>
					<span class="text-sm text-gray-500">Info</span>
				</div>
			</div>
		</div>
	</div>
</div>
