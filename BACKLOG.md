# WireGuard Manager - Backlog

**Last Updated**: 2026-02-01

This document tracks pending work, known issues, and future enhancements.

---

## Blocking Issues

### UI Does Not Match Mockup Design

**Priority**: HIGH (Blocking)  
**Status**: In Progress  
**Affects**: All frontend pages

**Problem**: 
Current UI implementation does not visually match the Figma-exported HTML mockups. Testing cannot proceed until the UI achieves acceptable visual fidelity.

**Mockup Location**: `specs/001-frontend-implementation/design/stitch_vpn_management_dashboard/`

**Specific Gaps**:
- Traffic charts are SVG placeholders (not real data visualizations)
- Spacing and padding don't match mockup exactly
- Shadows and gradients need fine-tuning
- Font weights and sizes may differ
- Glass effects need polish

**Next Steps**:
1. Compare each component side-by-side with mockup HTML
2. Extract exact CSS values from mockup files
3. Update Tailwind classes to match
4. Consider using mockup CSS directly where appropriate

---

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

## Infrastructure / DevOps

### CI Workflow Disabled

**Priority**: Medium  
**File**: `.github/workflows/ci.yml`

The entire CI workflow is commented out. Before enabling:
1. Verify Go version 1.25.6 exists (or update to valid version)
2. Test each job individually
3. Ensure frontend build passes in CI environment

### Go Version Validation

**Priority**: Medium  
**Affects**: Backend CI

CI specifies Go version `1.25.6` which may not exist. Verify and update to a valid Go version before enabling CI.

---

## Component-Specific Issues

### PeerModal Styling

**Priority**: Medium  
**Component**: `src/lib/components/PeerModal.svelte`

Modal styling needs alignment with `create_new_peer/code.html` mockup:
- Input field styling
- Button placement
- Form layout

### Settings Page

**Priority**: Medium  
**Page**: `/settings`

Settings page needs alignment with `interface_settings_&_logs/code.html` mockup.

### Peers Page

**Priority**: Medium  
**Page**: `/peers`

Peer management page may need updates to match overall design system.

---

*Last Updated: 2026-02-01*
