# WireGuard Manager - Backlog

## API Enhancements

### Missing Stats Endpoint Properties

**Priority**: Medium  
**Affects**: Dashboard UI (`src/routes/+page.svelte`)

The mockup design shows additional interface statistics that are not currently returned by the `/stats` endpoint.

#### Properties to Add to `/stats` Response

| Property | Type | Description | Source |
|----------|------|-------------|--------|
| `publicKey` | `string` | Server's WireGuard public key | `wgctrl` Device.PublicKey |
| `listenPort` | `number` | WireGuard listening port (e.g., 51820) | `wgctrl` Device.ListenPort |
| `subnet` | `string` | VPN subnet CIDR (e.g., "10.0.0.0/24") | Configuration or derived |

#### Implementation Notes

1. **Backend Changes** (`backend/internal/wireguard/service.go`):
   - Update `Stats` struct to include new fields
   - Populate `publicKey` and `listenPort` from `wgctrl.Device()`
   - `subnet` may need to be stored in config or derived from peer AllowedIPs

2. **API.md Updates**:
   - Update `/stats` response schema documentation

3. **Frontend Changes** (`src/lib/types/stats.ts`):
   - Add new fields to `InterfaceStats` interface

#### Current Workaround

Dashboard displays available data (interface name, peer count, traffic stats) with alternative cards:
- Interface Name
- Total Peers  
- Online Peers count

---

## Future Enhancements

### QR Code Modal - Regenerate Keys

**Priority**: Low  
**Component**: `QRCodeDisplay.svelte`

The "Regenerate Keys" button in the QR modal is currently a placeholder. Needs:
- Backend endpoint to regenerate keypair for a peer
- Update peer configuration
- Refresh QR code with new config

### Real-time Updates

**Priority**: Low  
**Affects**: Dashboard, Peer list

Currently uses manual refresh. Consider:
- WebSocket for live stats updates
- Polling interval for peer status
- SSE (Server-Sent Events) alternative

---

*Last Updated: 2026-02-01*
