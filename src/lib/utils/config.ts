// WireGuard configuration utilities

import type { PeerCreateResponse } from '../types/peer';

/**
 * Generate WireGuard client configuration file content
 * @param response - Peer creation response from API
 * @returns Configuration file content as string
 */
export function generateWireGuardConfig(response: PeerCreateResponse): string {
	// Use the config from API response if available
	if (response.config) {
		return response.config;
	}

	console.log(response)
	// Fallback: generate config manually (should not normally happen)
	// This assumes the API provides the full config in the response
	console.warn('No config in API response, generating fallback config');

	const config = `[Interface]
PrivateKey = ${response.privateKey || '<PRIVATE_KEY>'}
Address = ${response.allowedIPs.join(', ')}
DNS = 1.1.1.1

[Peer]
PublicKey = ${response.publicKey}
Endpoint = <SERVER_ENDPOINT>:51820
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25`;

	return config;
}

/**
 * Trigger browser download of configuration file
 * @param config - Configuration file content
 * @param peerName - Peer name for filename
 */
export function downloadConfigFile(config: string, peerName: string): void {
	// Sanitize filename (remove special characters)
	const filename = `${peerName.replace(/[^a-zA-Z0-9-_]/g, '_')}.conf`;

	// Create blob and download
	const blob = new Blob([config], { type: 'text/plain' });
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
}
