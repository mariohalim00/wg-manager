# Domain Knowledge & Terminology

**Last Updated**: 2026-01-31

This document clarifies domain-specific terms and concepts used in this codebase. Understanding these prevents confusion and helps implement features correctly.

## WireGuard Concepts

### Peer

In WireGuard terminology, a **peer** is a client or remote endpoint connected to the VPN.

**Key properties**:
- **PublicKey** (immutable identifier): The peer's WireGuard public key (base64 string, ~44 chars)
  - Unique identifier in the system
  - Used in API URLs: `DELETE /peers/{publicKey}`
  - Cannot be changed once added
- **AllowedIPs** (CIDR notation): IP ranges the peer is allowed to route
  - Example: `10.0.0.2/32` (single IP) or `10.0.0.0/24` (subnet)
  - Peer can only send/receive traffic from these IPs
  - Multiple IPs per peer allowed
- **Endpoint** (optional): Peer's public IP:port (where packets originate from)
  - Set by WireGuard kernel when peer connects
  - May be `null` if peer never connected
- **LastHandshake**: Timestamp of last successful WireGuard protocol handshake
  - Indicates peer is "active" if recent
  - `0` or missing = peer never connected
- **RX/TX bytes**: Real-time data transfer statistics from kernel

**In this system**:
- Stored in WireGuard kernel interface
- Metadata (names) stored in `peers.json` for user-friendly display
- Status (online/offline) inferred from LastHandshake time

### WireGuard Interface

A **WireGuard interface** is a virtual network device on the server.

**Key properties**:
- **Name**: Interface name in Linux (default: `wg0`)
  - Configured via `WG_INTERFACE_NAME` env var
  - May differ on different systems
- **PublicKey**: Server's public key (required for clients to connect)
- **PrivateKey**: Server's private key (secret, never shared)
- **ListenPort**: UDP port (usually 51820)

**In this system**:
- One interface configured per server
- Peers connect TO this interface
- Interface is managed via kernel module + `wgctrl` library

### Handshake

A **handshake** is a WireGuard protocol exchange that establishes a secure channel.

**When it happens**:
- Client initiates connection
- Server responds
- Both sides agree on encryption keys

**In this system**:
- `LastHandshake` timestamp indicates peer is active
- If LastHandshake > some threshold (e.g., 2 minutes), peer is considered "offline"
- Status badge logic depends on this

### CIDR Notation

**CIDR** = Classless Inter-Domain Routing (IP + prefix length)

**Format**: `IP/PREFIX`
- `10.0.0.2/32` = Single IP (host, /32 means all 32 bits are network)
- `10.0.0.0/24` = Network (first 24 bits fixed, last 8 bits variable)
  - Includes `10.0.0.1` through `10.0.0.255`
- `192.168.0.0/16` = Large network (2^16 addresses)

**Validation**: Backend checks that AllowedIPs are valid CIDR using `net.ParseCIDR()`

**In this system**:
- Each peer typically has `AllowedIPs: ["10.0.0.X/32"]` (single IP)
- Could support subnets if needed

## System Concepts

### Real vs. Mock Service

The system has two implementations of the `wireguard.Service` interface:

**realService**:
- Actually controls WireGuard kernel module via `wgctrl` library
- Requires: Linux, WireGuard kernel module, root/CAP_NET_ADMIN permissions
- Used in production and when WireGuard is available
- Lives in `backend/internal/wireguard/`

**mockService**:
- Stores peers in memory (HashMap)
- No kernel interaction
- Used for: Development, testing, fallback if WireGuard unavailable
- Lives in `backend/internal/wireguard/`

**Fallback logic** in `main.go`:
```go
wgService, err := wireguard.NewRealService(...)
if err != nil {
    slog.Warn("Failed to initialize real service, using mock", "error", err)
    wgService = wireguard.NewMockService()
}
```

**Key implication**: Testing doesn't require root or kernel module. Development works on any machine.

### Handler vs. Service

**Handler** (HTTP layer):
- Parses incoming HTTP request
- Validates input (CIDR, required fields)
- Calls service method
- Returns HTTP response

**Service** (business logic layer):
- Performs actual WireGuard operations
- Persists metadata
- Returns domain objects (Peer, Stats)
- Does not know about HTTP

**Separation benefit**: Service is reusable, testable without HTTP, swappable (real vs. mock).

### Metadata vs. Real-Time State

**Metadata** (persistent, from `peers.json`):
- Peer name (user-chosen label)
- Stored on disk
- Changes only when user adds/removes peer

**Real-time state** (from kernel):
- LastHandshake (when peer last connected)
- RX/TX bytes (cumulative since peer added)
- Endpoint (current IP:port)
- Extracted fresh from kernel on every `ListPeers()` call

**Why separate?**:
- WireGuard kernel doesn't store names (only PublicKey + config)
- Metadata needs persistence across reboots
- Real-time stats need freshness

## UI/Frontend Concepts

### Peer Status

A peer can be "online" or "offline" based on LastHandshake time.

**Definition** (in frontend):
- **Online**: LastHandshake < ~2 minutes ago (configurable threshold)
- **Offline**: LastHandshake > 2 minutes ago or never connected (LastHandshake = 0)

**Why this matters**:
- Users see visual indicator (StatusBadge component)
- Online = peer can route traffic
- Offline = peer may need reconnection

**Note**: WireGuard automatically reconnects; "offline" doesn't mean connection broken, just dormant.

### Transfer Stats

Data transfer metrics displayed in UI.

**RX** = Received bytes (data from peer to server)
**TX** = Transmitted bytes (data from server to peer)

**Characteristics**:
- Cumulative since peer was added (never resets)
- From kernel's perspective (updated in real-time)
- Displayed in PeerTable and stats pages

## Data Model

### Peer Interface (Frontend TypeScript)

```typescript
interface Peer {
    id: string;              // PublicKey (unique identifier)
    name: string;            // User-chosen name from metadata
    publicKey: string;       // WireGuard public key (same as id)
    status: 'online' | 'offline'; // Inferred from LastHandshake
    allowedIps: string[];    // CIDR notation array
    latestHandshake: string; // ISO timestamp or "Never"
    transfer: {
        received: number;    // RX bytes
        sent: number;        // TX bytes
    };
}
```

**Naming conventions**:
- TypeScript uses camelCase: `allowedIps`, `latestHandshake`
- Backend uses snake_case in JSON: `allowed_ips`, `latest_handshake` (if needed)
- Keep in sync across frontend/backend (type definitions in `types.ts`)

### Stats Object

```typescript
interface Stats {
    interfaceName: string;   // e.g., "wg0"
    peerCount: number;       // Total peers connected
    totalRx: number;         // Total bytes received across all peers
    totalTx: number;         // Total bytes transmitted across all peers
}
```

## Operational Concepts

### API Idempotency

- **GET /peers** — Always returns current state (idempotent)
- **GET /stats** — Always returns current state (idempotent)
- **POST /peers** — Creates new peer (idempotent if same PublicKey, but errors on duplicate)
- **DELETE /peers/{id}** — Removes peer (idempotent if peer already deleted, returns 204)

**Implication**: Frontend can retry failed requests safely (with caveats for POST).

### Configuration & Deployment

**WireGuard configuration scope**:
- This system manages ONE WireGuard interface (e.g., `wg0`)
- Multiple interfaces not currently supported
- Interface name configurable via `WG_INTERFACE_NAME`

**Peer configuration scope**:
- Each peer is independent
- Adding peer: Just adds to WireGuard (doesn't connect client automatically)
- Client must be configured with server's endpoint + key to connect
- Server endpoint is `WG_SERVER_ENDPOINT` (public IP:port for clients to connect to)

**Bootstrapping a peer**:
1. Server adds peer (POST /peers)
2. Server generates/displays QR code or config file (frontend responsibility)
3. Client scans QR or imports config
4. Client connects; handshake occurs
5. Status becomes "online" in dashboard

## Common Confusion Points

### Q: What's the difference between `publicKey` and `id` in Peer?

**A**: They're the same value (both PublicKey). `id` is the unique identifier, `publicKey` is the actual key. Kept for clarity in both frontend and API responses.

### Q: Why does adding a peer NOT require a private key?

**A**: 
- PublicKey can be provided by client (they generated their own keypair)
- Backend can generate a keypair if not provided
- PrivateKey is never stored or transmitted (security)
- Client keeps their PrivateKey secret; server only needs PublicKey

### Q: What happens if two peers have the same PublicKey?

**A**: Error (not allowed). PublicKey must be unique. The system prevents duplicates.

### Q: Why is AllowedIPs an array, not a single IP?

**A**: Flexibility. A peer might legitimately route multiple IP ranges. Common case: `["10.0.0.2/32"]` (one IP), but `["10.0.0.2/32", "10.0.0.3/32"]` is valid (peer controls two IPs).

### Q: What's the difference between "offline" and "disconnected"?

**A**: In this system, **offline** is the term used (status = 'offline'). It means no recent handshake. The peer is not actively routing traffic, but the WireGuard config still exists. "Disconnected" isn't used formally.

### Q: If I delete a peer, can I add it back with the same PublicKey?

**A**: Yes. Deletion only removes from WireGuard and metadata. The PublicKey can be reused.

## Glossary

| Term | Definition |
|------|-----------|
| **Peer** | A client/endpoint in WireGuard (identified by PublicKey) |
| **Handshake** | WireGuard protocol exchange (LastHandshake indicates peer activity) |
| **CIDR** | IP notation (e.g., `10.0.0.2/32`) for allowed IP ranges |
| **RX/TX** | Bytes received/transmitted (from kernel stats) |
| **Interface** | Virtual WireGuard network device (e.g., `wg0`) |
| **Metadata** | Peer names, stored in `peers.json` |
| **Real-time state** | Current kernel state (LastHandshake, RX/TX, endpoint) |
| **Status** | Online/offline (based on LastHandshake freshness) |
| **Endpoint** | Peer's public IP:port (set by kernel when peer connects) |
