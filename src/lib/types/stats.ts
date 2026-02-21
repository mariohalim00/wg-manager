// Interface statistics types

export interface InterfaceStats {
	interfaceName: string; // e.g., "wg0"
	publicKey: string; // Server's WireGuard public key
	listenPort: number; // WireGuard listening port
	subnet: string; // VPN subnet CIDR
	peerCount: number; // Total peers
	totalRx: number; // Total bytes received
	totalTx: number; // Total bytes transmitted
}

export interface StatsHistoryItem {
	timestamp: number;
	totalRx: number;
	totalTx: number;
}
