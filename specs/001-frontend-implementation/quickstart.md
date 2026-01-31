# WireGuard Manager - Frontend Quickstart Guide

This guide provides step-by-step instructions for setting up, running, and testing the WireGuard Manager frontend application.

## Prerequisites

### Required Software

- **Node.js**: v18 or higher (v20 recommended)
- **npm**: v9 or higher (comes with Node.js)
- **Git**: For version control
- **Modern Browser**: Chrome, Firefox, Safari, or Edge (latest version)
- **Backend API**: WireGuard Manager backend running on http://localhost:8080

### System Requirements

- Operating System: Linux, macOS, or Windows (WSL2 recommended for Windows)
- RAM: Minimum 4GB
- Disk Space: ~500MB for dependencies

## Installation

### 1. Clone the Repository (if not already done)

```bash
git clone <repository-url>
cd wg-manager
```

### 2. Install Dependencies

```bash
npm install
```

This will install:
- SvelteKit 2.x
- Svelte 5.x
- TypeScript 5.x
- TailwindCSS 4.x + DaisyUI
- Vite 7.x
- Development tools (ESLint, Prettier, svelte-check)

**Expected output**: ~350 packages installed, completion time ~30-60 seconds

### 3. Verify Installation

```bash
npm run check
```

**Expected output**: "svelte-check found 0 errors" (warnings are acceptable)

## Development Server

### Start Development Server

```bash
npm run dev
```

**Default URL**: http://localhost:5173  
**Auto-reload**: Enabled (changes auto-refresh the browser)

### Alternative: Expose to Network

```bash
npm run dev -- --host
```

**Access from other devices**: http://<your-ip>:5173

### Stop Development Server

Press `Ctrl+C` in the terminal

## Project Structure

```
src/
├── routes/                      # File-based routing (SvelteKit)
│   ├── +page.svelte            # Dashboard (/)
│   ├── +layout.svelte          # Global layout & navigation
│   ├── +error.svelte           # Error page
│   ├── peers/+page.svelte      # Peer management (/peers)
│   ├── stats/+page.svelte      # Statistics (/stats)
│   └── settings/+page.svelte   # Settings (/settings - mock UI)
├── lib/
│   ├── components/             # Reusable Svelte components
│   │   ├── PeerTable.svelte    # Peer list table
│   │   ├── PeerModal.svelte    # Add/edit peer modal
│   │   ├── QRCodeDisplay.svelte # QR code display modal
│   │   ├── ConfirmDialog.svelte # Confirmation dialog
│   │   ├── StatusBadge.svelte  # Online/offline indicator
│   │   ├── StatsCard.svelte    # Statistics card
│   │   ├── LoadingSpinner.svelte # Loading indicator
│   │   └── Notification.svelte # Toast notifications
│   ├── stores/                 # State management
│   │   ├── peers.ts            # Peer list + CRUD operations
│   │   └── stats.ts            # Statistics data
│   ├── utils/                  # Utility functions
│   │   ├── api.ts              # API client
│   │   ├── config.ts           # Config generation
│   │   └── formatting.ts       # Data formatters
│   └── types.ts                # TypeScript interfaces
├── app.css                     # Global styles (Tailwind)
└── app.html                    # HTML template
```

## Backend Integration

### Configure Backend URL

**Default**: API calls go to http://localhost:8080

**Change backend URL**: Edit `src/lib/utils/api.ts`

```typescript
const API_BASE_URL = 'http://localhost:8080'; // Change this
```

### Start Backend Server

Navigate to backend directory and run:

```bash
cd backend
go run cmd/server/main.go
```

**Verify backend is running**:
```bash
curl http://localhost:8080/peers
```

**Expected output**: JSON array of peers (may be empty initially)

### CORS Configuration

Backend must allow frontend origin. In `backend/internal/config/config.json`:

```json
{
  "corsAllowedOrigins": ["http://localhost:5173"]
}
```

Or set environment variable:
```bash
export CORS_ALLOWED_ORIGINS=http://localhost:5173
```

## Manual Testing Checklist

### Phase 3: View Peers (US1)

**T032: Peer list displays correctly**

1. Navigate to http://localhost:5173/peers
2. **Verify**:
   - Table displays with columns: Name, Public Key, Status, Allowed IPs, Handshake, Actions
   - Status badges show "Online" (green) or "Offline" (gray)
   - Empty state shows when no peers exist
   - Loading spinner appears during initial load

**T033: Peer details are accurate**

1. Add a test peer via backend or UI
2. **Verify**:
   - Name matches input
   - Public Key displays (truncated with "..." if long)
   - AllowedIPs show in comma-separated list
   - Last Handshake shows timestamp or "Never"
   - TX/RX bytes display (formatted as KB/MB/GB)

### Phase 4: Add Peer (US2)

**T041: Peer creation workflow**

1. Click "+ Add Peer" button on /peers page
2. **Verify**:
   - Modal opens with glassmorphism styling
   - Form has fields: Name, Public Key (optional), Allowed IPs
   - Validation errors show for invalid inputs:
     - Empty name → "Name is required"
     - Invalid CIDR → "Invalid CIDR notation" (e.g., "10.0.0.2" without /32)
   - "Generate Key" button creates random key pair
   - "Add Peer" button disabled until valid input

3. Fill valid data:
   - Name: "Test Peer"
   - Public Key: (leave empty or use generated)
   - Allowed IPs: "10.0.0.2/32"

4. Click "Add Peer"
5. **Verify**:
   - Success notification appears ("Peer added successfully")
   - Modal closes
   - New peer appears in table immediately
   - QR code modal appears automatically (see Phase 5)

**T042: Form validation edge cases**

1. Try adding peer with duplicate Public Key
2. **Verify**: Error notification shows ("Peer with this public key already exists")

3. Try CIDR without /prefix: "10.0.0.2"
4. **Verify**: Error message shows before submission

5. Try multiple IPs: "10.0.0.2/32, 10.0.0.3/32"
6. **Verify**: Accepted and displayed correctly in table

### Phase 5: Download Config (US5)

**T047: QR code displays after peer creation**

1. Add a new peer (see Phase 4)
2. **Verify** QR code modal appears automatically:
   - QR code renders (scannable image)
   - Security warning banner shows (yellow/orange background)
   - "View configuration text" collapsible section present
   - "Download Config" and "Close" buttons visible

3. Scan QR code with WireGuard mobile app
4. **Verify**: Config imports successfully on mobile device

**T048: Config download functionality**

1. In QR code modal, click "Download Config"
2. **Verify**:
   - File downloads as `peer-name.conf` (name sanitized, no spaces/special chars)
   - File contains valid WireGuard config:
     ```ini
     [Interface]
     PrivateKey = <base64>
     Address = 10.0.0.2/32
     DNS = 1.1.1.1

     [Peer]
     PublicKey = <server-public-key>
     Endpoint = <server-ip>:51820
     AllowedIPs = 0.0.0.0/0
     PersistentKeepalive = 25
     ```

3. Import config file into WireGuard desktop client
4. **Verify**: Connection succeeds

### Phase 6: Remove Peer (US3)

**T055: Delete confirmation workflow**

1. Navigate to /peers page
2. Click "Delete" button next to any peer
3. **Verify**:
   - Confirmation dialog appears (glassmorphism modal)
   - Dialog shows peer name in message
   - "Cancel" and "Delete" buttons present
   - "Delete" button has red/danger styling

4. Click "Cancel"
5. **Verify**: Dialog closes, peer remains in table

6. Click "Delete" again, then "Delete" button in dialog
7. **Verify**:
   - Success notification shows ("Peer removed successfully")
   - Peer disappears from table immediately
   - No page reload required

**T056: Delete error handling**

1. Stop backend server
2. Try deleting a peer
3. **Verify**:
   - Error notification shows ("Failed to remove peer")
   - Peer remains in table
   - Dialog closes

4. Restart backend, retry delete
5. **Verify**: Delete succeeds

### Phase 7: View Stats (US4)

**T060: Stats page displays correctly**

1. Navigate to http://localhost:5173/stats
2. **Verify**:
   - 4 stat cards in grid (2x2):
     - Total Peers (purple icon)
     - Data Received (green icon)
     - Data Transmitted (yellow icon)
     - Online Peers (blue icon)
   - Interface details panel shows:
     - Interface name (e.g., "wg0")
     - Total peers count
     - Online peers count
   - Traffic summary panel shows:
     - Total RX (formatted bytes)
     - Total TX (formatted bytes)
   - Loading spinner shows initially
   - Data loads within 1-2 seconds

**T061: Dashboard displays correctly (FR-008a)**

1. Navigate to http://localhost:5173/
2. **Verify**:
   - Page title: "Dashboard"
   - **3-card horizontal grid** (FR-008a requirement):
     - Card 1: Total Peers (with "X online" subtitle)
     - Card 2: Data Received (RX)
     - Card 3: Data Transmitted (TX)
   - Quick Actions panel with 3 buttons:
     - "Manage Peers" → navigates to /peers
     - "View Statistics" → navigates to /stats
     - "Settings" → navigates to /settings
   - Network Status panel shows:
     - Interface name with "Active" badge
     - Peer count with online count
     - "View all →" link to /peers

### Phase 8: Settings Page (FR-021)

**Manual test**: Navigate to http://localhost:5173/settings

**Verify**:
- "Coming Soon" banner with construction emoji
- 4 mock panels displayed (2x2 grid):
  - Interface Settings (wg0, port, MTU, address)
  - Appearance (theme, glassmorphism)
  - Security (auth, session timeout)
  - Advanced (auto-refresh, log level)
- All panels have `opacity-50` (disabled appearance)
- Consistent glassmorphism styling

## Performance Validation

### Lighthouse Audit

1. Open browser DevTools (F12)
2. Navigate to "Lighthouse" tab
3. Select:
   - Mode: Navigation
   - Device: Desktop
   - Categories: Performance
4. Click "Analyze page load"

**Target scores** (Constitution V):
- Performance: ≥90
- First Contentful Paint (FCP): <1.5s
- Time to Interactive (TTI): <3s

### Bundle Size Check

```bash
npm run build
```

**Check output**:
```
vite v7.x.x building for production...
✓ built in Xs
✓ 5 modules transformed.

dist/_app/immutable/... (gzipped size should be <200KB total)
```

**Target**: JavaScript bundle <200KB gzipped (Constitution V)

### Network Throttling Test

1. Open DevTools → Network tab
2. Set throttling: "Fast 3G"
3. Reload page
4. **Verify**: TTI <3s, page usable within 3 seconds

## Responsive Testing

### Breakpoint Testing

Test at these viewport sizes:

**Mobile** (320px width):
- Open DevTools → Device Toolbar (Ctrl+Shift+M)
- Select "iPhone SE" or custom 320px width
- **Verify**:
  - Navigation collapses to hamburger menu
  - Tables scroll horizontally if needed
  - Cards stack vertically
  - Buttons remain accessible
  - Text remains readable

**Tablet** (768px width):
- Select "iPad Mini" or custom 768px width
- **Verify**:
  - 2-column grid for stat cards
  - Navigation shows full links
  - Tables fit within viewport

**Desktop** (1024px+):
- Select "Laptop" or custom 1024px width
- **Verify**:
  - 3-column grid for dashboard cards
  - Sidebar navigation visible
  - All features accessible
  - Hover effects work

### Cross-Browser Testing

Test in:
- **Chrome**: Latest version (primary target)
- **Firefox**: Latest version
- **Safari**: Latest version (macOS)
- **Edge**: Latest version

**Verify**:
- Glassmorphism effects render correctly
- Modals open/close smoothly
- API calls succeed
- No console errors

## Troubleshooting

### Common Issues

**Issue**: "Cannot GET /peers" in browser

**Solution**: Start development server with `npm run dev`

---

**Issue**: "Network error" notifications

**Solution**: 
1. Verify backend is running: `curl http://localhost:8080/peers`
2. Check CORS configuration in backend
3. Verify API_BASE_URL in `src/lib/utils/api.ts`

---

**Issue**: "Module not found: svelte-qrcode"

**Solution**: Reinstall dependencies: `npm install`

---

**Issue**: Blank page or white screen

**Solution**:
1. Check browser console for errors
2. Verify backend is running
3. Clear browser cache and reload

---

**Issue**: TypeScript errors in editor

**Solution**: Run `npm run check` and fix reported issues

---

**Issue**: Styles not applied

**Solution**:
1. Check Tailwind config includes correct content paths
2. Restart dev server: `Ctrl+C` then `npm run dev`

## Production Build

### Build for Production

```bash
npm run build
```

**Output**: `build/` directory with static files

### Preview Production Build

```bash
npm run preview
```

**Default URL**: http://localhost:4173

### Deploy Production Build

The `build/` directory contains static files. Deploy to:
- **Nginx**: Serve from `build/` directory
- **Caddy**: Serve with SPA fallback
- **Netlify/Vercel**: Connect Git repo (auto-deploy)

**Example Nginx config**:
```nginx
server {
    listen 80;
    server_name wg-manager.example.com;
    root /path/to/build;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
    }
}
```

## Code Quality

### Run ESLint

```bash
npm run lint
```

**Fix issues automatically**:
```bash
npm run lint -- --fix
```

### Format with Prettier

```bash
npm run format
```

### TypeScript Check

```bash
npm run check
```

**Watch mode** (auto-check on file changes):
```bash
npm run check:watch
```

## Development Tips

### Hot Module Replacement (HMR)

- Changes to `.svelte` files auto-reload without losing state
- Changes to stores may require full page reload
- CSS changes apply instantly

### Debugging

**Browser DevTools**:
- Use React DevTools for component inspection (Svelte adapters available)
- Network tab shows API calls
- Console logs errors and warnings

**Svelte-specific debugging**:
```svelte
<script>
  $inspect(variableName); // Logs reactive changes
</script>
```

### State Management

- Use stores for shared state: `import { peers } from '$lib/stores/peers'`
- Subscribe with `$store` syntax in components
- Avoid prop drilling; stores are global

### Performance Tips

- Use `{#key}` blocks to force re-renders
- Lazy-load heavy components with `import()`
- Optimize images (WebP format, responsive sizes)
- Minimize bundle with tree-shaking (Vite handles this)

## Functional Requirements Verification

### All FR Requirements (FR-001 through FR-021)

Run through all acceptance scenarios from `specs/001-frontend-implementation/spec.md`:

**FR-001** — ✓ Dashboard shows peer count, RX, TX, online peers  
**FR-002** — ✓ Peer list displays in table with status badges  
**FR-003** — ✓ Add peer modal with form validation  
**FR-004** — ✓ Key generation button works  
**FR-005** — ✓ CIDR validation enforced  
**FR-006** — ✓ Delete confirmation dialog  
**FR-007** — ✓ QR code display after peer creation  
**FR-008** — ✓ Config file download  
**FR-008a** — ✓ 3-card dashboard layout (Constitution requirement)  
**FR-009** — ✓ Stats page with detailed metrics  
**FR-010** — ✓ Responsive layout (mobile, tablet, desktop)  
**FR-011** — ✓ Navigation sidebar (desktop) and bottom nav (mobile)  
**FR-012** — ✓ Loading states for all API calls  
**FR-013** — ✓ Error notifications  
**FR-014** — ✓ Success notifications  
**FR-015** — ✓ Empty state messages  
**FR-016** — ✓ TypeScript strict mode (no `any` types)  
**FR-017** — ✓ ESLint + Prettier configured  
**FR-018** — ✓ Glassmorphism design system  
**FR-019** — ✓ Blur values: 16px (cards), 40px (modals)  
**FR-020** — ✓ Consistent spacing and colors  
**FR-021** — ✓ Settings page mock UI with "Coming Soon"

## Next Steps

After completing this quickstart:

1. **Explore codebase**: Read component source code to understand patterns
2. **Modify styling**: Edit `tailwind.config.js` to customize theme
3. **Add features**: Follow existing component patterns for new features
4. **Backend integration**: Test with live WireGuard backend
5. **Deploy**: Build and deploy to production environment

## Support & Documentation

- **API Documentation**: `backend/API.md`
- **Component Library**: Explore `src/lib/components/`
- **Type Definitions**: See `src/lib/types.ts`
- **Spec Document**: `specs/001-frontend-implementation/spec.md`
- **Tasks Breakdown**: `specs/001-frontend-implementation/tasks.md`

---

**Last Updated**: 2026-01-31  
**Frontend Version**: v1.0.0 (MVP)  
**Constitution Compliance**: ✓ All principles followed
