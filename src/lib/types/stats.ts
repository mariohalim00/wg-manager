// Interface statistics types

export interface InterfaceStats {
	interfaceName: string; // e.g., "wg0"
	peerCount: number; // Total peers
	totalRx: number; // Total bytes received
	totalTx: number; // Total bytes transmitted
}
