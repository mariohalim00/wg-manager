// Peer types for WireGuard peer management

export interface Peer {
	id: string; // PublicKey (unique identifier)
	publicKey: string; // WireGuard public key (base64)
	name: string; // User-friendly name
	endpoint?: string; // Peer's public IP:port (optional, set by kernel)
	allowedIPs: string[]; // CIDR notation array (e.g., ["10.0.0.2/32"])
	lastHandshake: string; // ISO timestamp or "0" (never connected)
	receiveBytes: number; // Total bytes received
	transmitBytes: number; // Total bytes transmitted
	status: 'online' | 'offline'; // Derived from lastHandshake (client-side)
	config?: string;
	dns?: string;
	mtu?: number;
	persistentKeepalive?: number;
	preSharedKey?: boolean;
}

export interface PeerFormData {
	name: string; // Required user-friendly name
	allowedIPs: string[]; // CIDR strings (validated before submit)
	publicKey?: string; // Optional (backend generates if omitted)
	dns?: string;
	mtu?: number;
	persistentKeepalive?: number;
	preSharedKey?: boolean;
}

export interface PeerCreateResponse {
	id: string;
	publicKey: string;
	name: string;
	allowedIPs: string[];
	config: string; // WireGuard .conf file content
	privateKey?: string; // Only if backend generated keypair
	presharedKey?: string;
}

export interface PeerUpdateRequest {
	name?: string;
	allowedIPs?: string[];
	dns?: string;
	mtu?: number;
	persistentKeepalive?: number;
}
