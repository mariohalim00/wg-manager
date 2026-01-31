# Design Analysis & API Gap Report

**Date**: 2026-02-01  
**Feature**: SvelteKit Frontend Implementation  
**Purpose**: Compare design mockups with current spec and identify API gaps

## Design Analysis Summary

Based on the provided HTML design files in `specs/001-frontend-implementation/design/`, the visual design implements a **glassmorphism aesthetic** with:

### Core Design Characteristics

1. **Glassmorphism Styling**
   - `background: rgba(255, 255, 255, 0.03)` with `backdrop-filter: blur(12px)`
   - Semi-transparent panels with subtle borders (`border: 1px solid rgba(255, 255, 255, 0.08)`)
   - Layered depth with varying blur intensities (12px, 16px, 20px, 24px)
   - Radial gradient backgrounds: `radial-gradient(circle at top left, #1a2a3a 0%, #101922 100%)`

2. **Color Palette**
   - Primary: `#137fec` (blue)
   - Background Dark: `#101922`, `#0a0f14`, `#0b1117` (dark navy variations)
   - Status Colors: Green `#22c55e` (online), Red (offline/delete), Slate gray (offline state)
   - Accent overlays: `bg-white/5`, `bg-white/10` for hover states

3. **Typography**
   - Font: Inter (display text), JetBrains Mono (monospace for IPs/codes)
   - Bold headings with tight tracking (`tracking-tight`, `tracking-tighter`)
   - Uppercase labels with wider tracking (`tracking-wider`, `tracking-widest`)

4. **Visual Components**
   - Animated status indicators with pulse effect (`@keyframes pulse`)
   - Material Symbols Outlined icons
   - SVG charts for traffic visualization
   - QR code display with gradient container backgrounds

### Design Files Breakdown

#### 1. **vpn_management_dashboard_1/code.html** (Main Dashboard)

**Content**:

- **Sidebar navigation** (left side, 64px width)
  - Logo + version badge
  - Navigation links: Dashboard, Peers, Settings, Logs
  - Usage limit widget at bottom
- **Main content area**:
  - Quick stats cards (Active Peers, Public Key, Subnet)
  - Traffic charts (Total Received, Total Sent) with SVG line graphs
  - **Peers table** with columns:
    - Status (online/offline badge with pulse animation)
    - Peer Name (with device icon)
    - Internal IP (monospace)
    - Transfer U/D (upload/download with icons)
    - Last Handshake (relative time)
    - Actions (hover-reveal: download, QR, delete icons)

**Gaps vs. Current Spec**:

- ✅ Peer list: Covered in spec (FR-001, User Story 1)
- ✅ Status badges: Covered (FR-002)
- ✅ Transfer stats: Covered (FR-009)
- ❌ **Sidebar navigation**: Spec mentions "top navigation bar or sidebar" (FR-017) but doesn't detail sidebar layout
- ❌ **Quick stats cards**: Not explicitly in spec (Public Key, Subnet display)
- ❌ **Traffic charts**: Spec mentions stats page, but not inline charts on dashboard
- ❌ **Usage limit widget**: No API endpoint or spec requirement for quota tracking

#### 2. **create_new_peer/code.html** (Add Peer Modal/Page)

**Content**:

- Full-page centered form with glassmorphism card
- Top navigation bar (alternative layout, not sidebar)
- Form fields:
  - Peer Name (required)
  - Allowed IPs (CIDR notation with placeholder example)
  - Optional: Endpoint field, Pre-shared Key toggle
- Submit button with loading state
- Accent glow effects on focus

**Gaps vs. Current Spec**:

- ✅ Add peer form: Covered (FR-003, FR-004, User Story 2)
- ❌ **Endpoint field**: Not in current API `/peers POST` request (API.md shows only name, publicKey, allowedIPs)
- ❌ **Pre-shared Key option**: Not in API
- ✅ Validation and error handling: Covered (FR-004, FR-012)

#### 3. **peer_connection_details_1/code.html** (QR Code Modal)

**Content**:

- Full-screen overlay with backdrop blur
- Centered glassmorphism modal
- Peer name + online status indicator
- Large QR code display (white background with gradient container)
- Descriptive text for scanning instructions
- Close button (top-right)

**Gaps vs. Current Spec**:

- ✅ QR code display: Covered (FR-006, User Story 5)
- ✅ Modal overlay: Covered implicitly
- ✅ Peer details: Covered (name, status)

#### 4. **interface*settings*&\_logs/code.html** (Settings Page)

**Content**:

- Top navigation bar
- Breadcrumbs: System / Interface Configuration
- Page heading: "wg0 Interface" with Active badge
- Action buttons: Restart Service, Save Changes
- **Interface Toggle** (on/off switch for WireGuard service)
- **Network Configuration panel**:
  - Listen Port (input: 51820)
  - MTU Size (input: 1420)
  - Interface Addresses (CIDR input: 10.0.0.1/24, fd00::1/64)
- **Security Credentials panel**:
  - Server Public Key (read-only, copy button)
  - Server Private Key (hidden, reveal toggle)

**Gaps vs. Current Spec & API**:

- ❌ **NO API ENDPOINT** for interface settings (GET/PUT /interface or /settings)
- ❌ **NO API ENDPOINT** for service control (POST /interface/restart)
- ❌ Not in spec: Interface configuration is entirely missing from functional requirements
- ❌ **Server keys display**: API.md shows `WG_SERVER_PUBKEY` in config, but no GET endpoint to fetch dynamically

#### 5. **Mobile Views** (mobile\_\*.html files)

**Content**:

- Responsive layouts for peer creation, dashboard, QR display
- Bottom navigation instead of sidebar
- Stacked card layouts
- Touch-optimized buttons

**Gaps vs. Current Spec**:

- ⚠️ Spec assumption #5 states "Desktop-first design; mobile responsiveness is secondary"
- Mobile views suggest more investment in responsive design than spec indicates

## API Gap Analysis

### Existing API Endpoints (from backend/API.md)

| Endpoint      | Method | Purpose              | Design Usage              |
| ------------- | ------ | -------------------- | ------------------------- |
| `/peers`      | GET    | List all peers       | ✅ Main dashboard table   |
| `/peers`      | POST   | Add new peer         | ✅ Create peer form       |
| `/peers/{id}` | DELETE | Remove peer          | ✅ Delete action in table |
| `/stats`      | GET    | Interface statistics | ✅ Stats page             |

### Missing API Endpoints (Required by Design)

| Feature (Design)       | Required Endpoint                               | Priority    | Notes                                           |
| ---------------------- | ----------------------------------------------- | ----------- | ----------------------------------------------- |
| **Interface Settings** | `GET /interface` or `/config`                   | P3 (Future) | Fetch listen port, MTU, addresses, server keys  |
| **Interface Settings** | `PUT /interface` or `/config`                   | P3 (Future) | Update listen port, MTU, addresses              |
| **Service Control**    | `POST /interface/restart` or `/service/restart` | P3 (Future) | Restart WireGuard service                       |
| **Service Control**    | `POST /interface/toggle` or `/service/toggle`   | P3 (Future) | Start/stop WireGuard service                    |
| **Usage Quota**        | `GET /usage` or part of `/stats`                | P3 (Future) | Current usage vs. limit (sidebar widget)        |
| **Server Keys**        | `GET /keys` or part of `/stats`                 | P2 (Medium) | Fetch server public key dynamically for display |

### Data Gaps in Existing API

| API Response             | Current Schema                                     | Design Requirement                              | Gap?        |
| ------------------------ | -------------------------------------------------- | ----------------------------------------------- | ----------- |
| `POST /peers` response   | Includes `config` (string) and `privateKey`        | ✅ Sufficient for QR + download                 | No          |
| `GET /peers` response    | Includes `receiveBytes`, `transmitBytes`           | ✅ Sufficient for transfer stats                | No          |
| `GET /peers` response    | Includes `endpoint`, `lastHandshake`               | ✅ Sufficient for status + display              | No          |
| `GET /stats` response    | `interfaceName`, `peerCount`, `totalRx`, `totalTx` | ❌ Missing server public key                    | Yes (P2)    |
| `GET /stats` response    | Basic totals only                                  | ❌ Missing per-peer RX/TX history for charts    | Yes (P3)    |
| **NEW**: Server metadata | N/A                                                | Design shows server public key, subnet, version | Yes (P2-P3) |

## Recommendations

### Immediate Actions (Current Spec/MVP)

1. **Add FR-019**: System MUST display glassmorphism visual design with backdrop blur, semi-transparent panels, and radial gradient backgrounds
2. **Add FR-020**: System MUST implement sidebar navigation with logo, navigation links (Dashboard, Peers, Settings, Logs), and usage widget (data can be mocked initially)
3. **Update FR-017**: Clarify navigation structure—sidebar on desktop (design-preferred) vs. top nav (spec-mentioned alternative)
4. **Update Assumption #5**: Spec says "mobile responsiveness is secondary," but design includes full mobile views—clarify priority
5. **Add to Out of Scope**: Interface configuration page (Settings), service control (restart/toggle), usage quota tracking

### Future Feature Specifications (Deferred)

Create separate specs for:

1. **Interface Configuration Management** (P3):
   - API endpoints: GET/PUT `/interface` or `/config`
   - Settings page UI for listen port, MTU, addresses
   - Server key display and regeneration
   - Service control (restart, start/stop)

2. **Usage Quota Tracking** (P3):
   - API endpoint: GET `/usage` or extend `/stats`
   - Sidebar widget showing current usage vs. limit
   - Configurable quota thresholds

3. **Historical Traffic Charts** (P3):
   - API endpoint: GET `/stats/history` or `/peers/{id}/history`
   - Time-series data for RX/TX
   - SVG/Canvas chart components

### Design Implementation Notes

**For developers**:

- Use TailwindCSS classes matching design: `bg-white/5`, `backdrop-blur-lg`, `border-white/10`
- Material Symbols Outlined icons (already in design CDN links)
- Pulse animation keyframes for online status badges
- Hover-reveal pattern for peer table actions (`group-hover:opacity-100`)
- QR code library suggestion: `svelte-qrcode` (already in spec dependencies)

## Missing Features Report (for User)

### Features in Design WITHOUT Backend API Support

These features are **visually designed** but cannot be implemented without new backend endpoints:

1. **Interface Settings Page** (`interface_settings_&_logs/code.html`)
   - **Impact**: Cannot configure listen port, MTU, or interface addresses
   - **Workaround**: Mock data initially; defer to future spec
   - **Recommendation**: Add to Out of Scope section; create separate spec later

2. **Service Control** (Interface toggle, Restart button)
   - **Impact**: Cannot start/stop or restart WireGuard service from UI
   - **Workaround**: Show UI as read-only or hide controls
   - **Recommendation**: Add to Out of Scope; requires privileged backend API

3. **Usage Quota Widget** (Sidebar bottom)
   - **Impact**: Cannot track or display usage limits (6.5 GB of 10 GB)
   - **Workaround**: Hide widget or show static "Usage tracking coming soon"
   - **Recommendation**: Add to Out of Scope; requires backend quota system

4. **Server Public Key Display** (Dashboard quick stat, Settings page)
   - **Impact**: Backend config has `WG_SERVER_PUBKEY`, but no API endpoint to fetch dynamically
   - **Workaround**: Hardcode in frontend config or fetch from extended `/stats` response
   - **Recommendation**: **Enhance `/stats` API response** to include `serverPublicKey` field (low effort, high value)

### Features in Design WITH API Support (Ready to Implement)

✅ Peer list table with status, name, IP, transfer, handshake  
✅ Add peer form with name and allowedIPs validation  
✅ QR code display after peer creation  
✅ Delete peer with confirmation  
✅ Stats page with totals (can enhance with charts using same data)

## Conclusion

**Current spec covers 80% of design functionality.** The main gaps are:

- **Settings page** (no API endpoints)
- **Service control** (no API endpoints)
- **Usage quota** (no API endpoints)
- **Server public key display** (minor API enhancement needed)

**Recommended next steps**:

1. Run `/speckit.clarify` to refine navigation structure, visual design requirements, and mobile priority
2. Update spec with glassmorphism design system details (FR-019, FR-020)
3. Add Interface Settings, Service Control, and Usage Quota to Out of Scope
4. Optionally: Submit backend API enhancement request for server public key in `/stats` response
5. Proceed to `/speckit.plan` after clarifications integrated
