# Task Breakdown - Frontend Implementation
## WireGuard Manager SvelteKit Frontend

**Plan Reference**: `IMPLEMENTATION_PLAN.md`  
**Date**: 2026-02-01  
**Total Tasks**: 27  
**Estimated Time**: 8.5-11.5 hours

---

## Phase 0: Setup & Verification (30 minutes)

### Task 0.1: Fix CI/CD Workflow [15 min] ⚠️ CRITICAL BLOCKER
- [ ] Read `.github/workflows/ci.yml` file
- [ ] Uncomment all lines (currently all commented out)
- [ ] Check GitHub Actions setup-go versions: verify Go 1.25.6 exists
- [ ] If 1.25.6 not found: update to latest 1.25.x version
- [ ] Verify workflow YAML is valid syntax
- [ ] Commit: `git add .github/workflows/ci.yml && git commit -m "fix: enable CI workflow"`
- [ ] Push to feature branch: `git push origin feat/backend-improvement`
- [ ] Verify GitHub Actions runs automatically on push
- [ ] Confirm all jobs pass (backend tests + frontend build)

**Acceptance Criteria**:
- ✅ CI workflow uncommented and valid
- ✅ Go version resolved (1.25.6 or valid alternative)
- ✅ GitHub Actions runs on PR/push
- ✅ All backend tests pass in CI
- ✅ Frontend build succeeds in CI

**Files Modified**:
- `.github/workflows/ci.yml`

---

### Task 0.2: Verify Backend /stats Endpoint [10 min]
- [ ] Navigate to backend directory: `cd backend`
- [ ] Run tests: `go test ./... -v`
- [ ] Check for `/stats` endpoint test in output
- [ ] Verify test covers all 7 response fields:
  - interfaceName ✓
  - publicKey ✓
  - listenPort ✓
  - subnet ✓
  - peerCount ✓
  - totalRx ✓
  - totalTx ✓
- [ ] Confirm all tests pass
- [ ] Return to root: `cd ..`

**Acceptance Criteria**:
- ✅ All backend tests pass
- ✅ `/stats` endpoint test exists
- ✅ Test includes all 7 fields
- ✅ No errors in test output

**Files Verified**:
- `backend/cmd/server/main_test.go`
- `backend/internal/wireguard/wireguard.go`

---

### Task 0.3: Frontend Build Config Verification [5 min]
- [ ] Run build: `npm run build`
- [ ] Verify build succeeds without errors
- [ ] Check output directory: `ls -la build/`
- [ ] Verify `build/index.html` exists (SPA fallback)
- [ ] Check bundle size: `du -sh build/`
- [ ] Record baseline bundle size for performance comparison

**Acceptance Criteria**:
- ✅ Build succeeds with no errors
- ✅ `build/` directory created with files
- ✅ `build/index.html` exists for SPA mode
- ✅ Bundle size < 200KB (target)

**Files Verified**:
- `svelte.config.js`
- `src/routes/+layout.ts`
- `build/index.html`

---

## Phase 1: Core UI & Store Setup (2-3 hours)

### Task 1.1: Create API Client [30 min]
- [ ] Create file: `src/lib/api/client.ts`
- [ ] Export `ApiError` class with status code and message
- [ ] Export `apiCall()` function with:
  - Base URL from `VITE_API_BASE_URL` env variable (default: `http://localhost:8080`)
  - Automatic `Content-Type: application/json` header
  - Error handling for non-OK responses
  - 10-second timeout
- [ ] Throw `ApiError` on failure with status and message
- [ ] Update `src/lib/api/peers.ts`:
  - Import and use `apiCall()` instead of direct fetch
  - Export `listPeers()`, `addPeer()`, `deletePeer()` functions
- [ ] Update `src/lib/api/stats.ts`:
  - Import and use `apiCall()`
  - Export `getStats()` function

**Code Requirements**:
```typescript
// ApiError class
export class ApiError extends Error {
  constructor(public status: number, message: string) { ... }
}

// apiCall function
export async function apiCall(
  endpoint: string,
  options?: RequestInit
): Promise<any> { ... }
```

**Acceptance Criteria**:
- ✅ `ApiError` class properly throws and catches
- ✅ `apiCall()` handles timeout after 10s
- ✅ Non-200 responses throw error with status
- ✅ All API functions use centralized client
- ✅ Base URL configurable via env variable

**Files Modified**:
- `src/lib/api/client.ts` (new)
- `src/lib/api/peers.ts`
- `src/lib/api/stats.ts`

---

### Task 1.2: Create Peers Store with Optimistic Updates [45 min]
- [ ] Create file: `src/lib/stores/peers.ts`
- [ ] Export writable store: `peers` (array of Peer)
- [ ] Export function: `loadPeers()` - fetches from GET /peers
  - Handles errors gracefully
  - Updates `peers` store on success
- [ ] Export function: `addPeer(name: string, allowedIPs: string[])` - **OPTIMISTIC**
  - Creates temporary peer with unique ID
  - Updates store IMMEDIATELY (optimistic)
  - Calls POST /peers API in background
  - On success: replaces temp peer with server response
  - On failure: removes temp peer from store, throws error
  - Returns response on success
- [ ] Export function: `deletePeer(id: string)` - **OPTIMISTIC**
  - Saves previous state snapshot
  - Removes peer from store IMMEDIATELY (optimistic)
  - Calls DELETE /peers/{id} API in background
  - On success: peer stays removed
  - On failure: restores previous state, throws error
- [ ] Export error store: `peersError` (writable, stores error message)
- [ ] Export loading store: `peersLoading` (writable, boolean)

**Code Pattern**:
```typescript
export async function addPeer(name: string, allowedIPs: string[]) {
  const tempId = `temp-${Date.now()}`;
  const tempPeer = { id: tempId, name, allowedIPs, ... };
  
  // 1. Optimistic update
  peers.update(p => [...p, tempPeer]);
  
  try {
    // 2. API call
    const response = await apiCall('/peers', {
      method: 'POST',
      body: JSON.stringify({ name, allowedIPs })
    });
    
    // 3. Replace temp with real peer
    peers.update(p => p.map(x => 
      x.id === tempId ? response : x
    ));
    return response;
  } catch (error) {
    // 4. Rollback
    peers.update(p => p.filter(x => x.id !== tempId));
    peersError.set(error.message);
    throw error;
  }
}
```

**Acceptance Criteria**:
- ✅ Optimistic add updates UI before API response
- ✅ Optimistic delete removes from UI immediately
- ✅ Failed add rolls back peer from list
- ✅ Failed delete restores peer to list
- ✅ Error message available via peersError store
- ✅ Loading state tracked via peersLoading

**Files Created/Modified**:
- `src/lib/stores/peers.ts` (new)

---

### Task 1.3: Create Stats Store [20 min]
- [ ] Create file: `src/lib/stores/stats.ts` (or update existing)
- [ ] Export writable store: `stats` (Stats object)
- [ ] Export function: `loadStats()` - fetches from GET /stats
  - Updates store on success
  - Sets error on failure
- [ ] Export error store: `statsError` (writable)
- [ ] Export loading store: `statsLoading` (writable)
- [ ] Initialize stores with default values (0 counts, empty strings)

**Code Pattern**:
```typescript
export const stats = writable<Stats>({
  interfaceName: '',
  publicKey: '',
  listenPort: 0,
  subnet: '',
  peerCount: 0,
  totalRx: 0,
  totalTx: 0
});

export async function loadStats() {
  statsLoading.set(true);
  try {
    const data = await apiCall('/stats');
    stats.set(data);
  } catch (error) {
    statsError.set(error.message);
  } finally {
    statsLoading.set(false);
  }
}
```

**Acceptance Criteria**:
- ✅ Stats store updates on successful load
- ✅ Error state tracked
- ✅ Loading state tracked
- ✅ Default values initialized

**Files Created/Modified**:
- `src/lib/stores/stats.ts` (new/updated)

---

### Task 1.4: Create Notifications Store [20 min]
- [ ] Create file: `src/lib/stores/notifications.ts` (or update existing)
- [ ] Export writable store: `notifications` (array of Notification)
- [ ] Export function: `showNotification(type, message, duration=5000)`
  - Creates notification with ID
  - Adds to queue
  - Auto-dismisses after duration (if duration > 0)
- [ ] Export function: `removeNotification(id)`
  - Removes from queue
- [ ] Notification types: 'success' | 'error' | 'warning' | 'info'

**Code Pattern**:
```typescript
export function showNotification(
  type: 'success' | 'error' | 'warning' | 'info',
  message: string,
  duration: number = 5000
) {
  const id = Date.now().toString();
  notifications.update(n => [...n, { id, type, message }]);
  
  if (duration > 0) {
    setTimeout(() => removeNotification(id), duration);
  }
  return id;
}
```

**Acceptance Criteria**:
- ✅ Notifications can be added and removed
- ✅ Auto-dismiss works after specified duration
- ✅ Multiple notifications can be queued

**Files Created/Modified**:
- `src/lib/stores/notifications.ts` (new/updated)

---

### Task 1.5: Update +layout.svelte - Sidebar Navigation [45 min]
- [ ] Update `src/routes/+layout.svelte`
- [ ] Import Sidebar component
- [ ] Create layout structure: sidebar (left) + main content (right)
- [ ] Responsive layout:
  - Desktop (≥1024px): Full sidebar visible + main
  - Tablet (768-1023px): Simplified sidebar + main
  - Mobile (<768px): Bottom nav + main (sidebar collapsible)
- [ ] Add `<slot />` for child routes
- [ ] Add Notifications component at top/corner for toast display
- [ ] On mount: call `loadStats()` to populate stats
- [ ] Subscribe to `$stats` and `$statsLoading`

**CSS Requirements**:
```svelte
<div class="flex h-screen">
  <Sidebar /> {/* Left sidebar */}
  <main class="flex-1 overflow-auto">
    <Notifications /> {/* Toast at corner */}
    <slot />
  </main>
</div>
```

**Acceptance Criteria**:
- ✅ Sidebar displays navigation links
- ✅ Stats loaded on mount
- ✅ Responsive at all breakpoints
- ✅ Notifications visible at top/corner
- ✅ Child routes render in <slot />

**Files Modified**:
- `src/routes/+layout.svelte`

---

### Task 1.6: Create Sidebar Component [30 min]
- [ ] Create file: `src/lib/components/Sidebar.svelte`
- [ ] Display navigation sections:
  - Dashboard (/)
  - Peers (/peers)
  - Settings (/settings)
  - Logs (/logs) - disabled/grayed out
- [ ] Display usage widget at bottom (mock data for now)
- [ ] Active link indicator (current page highlighted)
- [ ] Icons for each section (use Lucide Svelte)
- [ ] Responsive styling:
  - Desktop: Full width, all text visible
  - Tablet: Narrower, text labels
  - Mobile: Collapsed/bottom nav (optional, can be hidden)

**Acceptance Criteria**:
- ✅ Navigation links work (use `<a>` or SvelteKit navigation)
- ✅ Active page highlighted
- ✅ Usage widget displays (mock: "5% used")
- ✅ Responsive at key breakpoints
- ✅ Icons render correctly

**Files Created**:
- `src/lib/components/Sidebar.svelte` (new)

---

### Task 1.7: Create Dashboard Page [45 min]
- [ ] Update `src/routes/+page.svelte`
- [ ] Display page title: "Dashboard"
- [ ] Create stats cards section with 3 cards:
  1. **Peer Count** - icon + "Total Peers" + number (from `$stats.peerCount`)
  2. **Online Count** - icon + "Online Peers" + derived count
  3. **Total Received** - icon + "Total RX" + formatted bytes
  4. **Total Sent** - icon + "Total TX" + formatted bytes
- [ ] Cards layout: horizontal grid, responsive (stack on mobile)
- [ ] Show loading skeleton while stats loading
- [ ] Subscribe to `$stats` and `$statsLoading`
- [ ] Use formatBytes utility for RX/TX (create in Task 4.2)
- [ ] Apply glassmorphism styling (Task 4.4)

**Code Pattern**:
```svelte
<script>
  import { stats, statsLoading, loadStats } from '$lib/stores/stats';
  import { formatBytes } from '$lib/utils/formatting';
  
  onMount(() => loadStats());
</script>

<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
  <StatsCard icon="users" label="Total Peers" value={$stats.peerCount} />
  <StatsCard icon="wifi" label="Online" value={onlineCount} />
  <StatsCard icon="download" label="Total RX" value={formatBytes($stats.totalRx)} />
  <StatsCard icon="upload" label="Total TX" value={formatBytes($stats.totalTx)} />
</div>
```

**Acceptance Criteria**:
- ✅ Dashboard displays 3-4 stats cards
- ✅ Cards show correct data from store
- ✅ Loading state shows skeleton
- ✅ Responsive layout on all breakpoints
- ✅ Bytes formatted as human-readable

**Files Modified**:
- `src/routes/+page.svelte`

---

### Task 1.8: Create Peers Page [20 min]
- [ ] Update `src/routes/peers/+page.svelte`
- [ ] Display page title: "Peers"
- [ ] Add "Add New Peer" button at top
- [ ] Render PeerTable component (delegated to Task 2.1)
- [ ] Load peers on mount: call `loadPeers()`
- [ ] Show loading state while fetching
- [ ] Show empty state if no peers

**Code Pattern**:
```svelte
<script>
  import { peers, peersLoading, loadPeers } from '$lib/stores/peers';
  import PeerTable from '$lib/components/PeerTable.svelte';
  import PeerModal from '$lib/components/PeerModal.svelte';
  
  let showAddModal = false;
  onMount(() => loadPeers());
</script>

<div class="space-y-4">
  <button on:click={() => showAddModal = true}>Add New Peer</button>
  {#if $peersLoading}
    <LoadingSpinner />
  {:else if $peers.length === 0}
    <EmptyState message="No peers configured" />
  {:else}
    <PeerTable />
  {/if}
  {#if showAddModal}
    <PeerModal on:close={() => showAddModal = false} />
  {/if}
</div>
```

**Acceptance Criteria**:
- ✅ Page loads peers on mount
- ✅ Loading state displayed while fetching
- ✅ Empty state shown if no peers
- ✅ "Add New Peer" button visible
- ✅ PeerTable rendered with peer data

**Files Modified**:
- `src/routes/peers/+page.svelte`

---

### Task 1.9: Create Stats Page [20 min]
- [ ] Create/update `src/routes/stats/+page.svelte`
- [ ] Display page title: "Interface Statistics"
- [ ] Show detailed stats:
  - Interface Name: `$stats.interfaceName`
  - Listen Port: `$stats.listenPort`
  - Subnet: `$stats.subnet`
  - Peer Count: `$stats.peerCount`
  - Total RX: formatted bytes
  - Total TX: formatted bytes
- [ ] Add refresh button to reload stats
- [ ] Show loading state during fetch
- [ ] Load stats on mount

**Acceptance Criteria**:
- ✅ All stat fields displayed
- ✅ Refresh button works
- ✅ Data updates on refresh
- ✅ Loading state shown

**Files Modified/Created**:
- `src/routes/stats/+page.svelte` (new/updated)

---

## Phase 2: Peer Management CRUD (3-4 hours)

### Task 2.1: Create PeerTable Component [60 min]
- [ ] Create file: `src/lib/components/PeerTable.svelte`
- [ ] Subscribe to `$peers` store
- [ ] Create table with columns:
  - Name
  - Status (badge: online/offline)
  - Allowed IPs
  - Last Handshake
  - RX bytes (formatted)
  - TX bytes (formatted)
  - Actions (hover-reveal on desktop, always visible on mobile)
- [ ] Status badge logic:
  - Online: if `lastHandshake` < 120 seconds ago
  - Offline: otherwise
  - Calculate with: `Date.now() - new Date(lastHandshake) < 120000`
- [ ] Action buttons (in actions column):
  - View QR Code
  - Download Config
  - Delete (with confirmation)
- [ ] Responsive visibility (FR-001a):
  - Desktop (≥1024px): use `opacity-0 group-hover:opacity-100` Tailwind
  - Mobile (<1024px): always visible
- [ ] Empty state if no peers

**Code Pattern**:
```svelte
<script>
  import { peers, deletePeer } from '$lib/stores/peers';
  import StatusBadge from './StatusBadge.svelte';
  
  function isOnline(lastHandshake: string): boolean {
    return Date.now() - new Date(lastHandshake).getTime() < 120_000;
  }
  
  function onDelete(peerId: string) {
    if (confirm(`Remove ${peerId}?`)) {
      deletePeer(peerId).catch(error => {
        showNotification('error', error.message);
      });
    }
  }
</script>

<table>
  {#each $peers as peer (peer.id)}
    <tr class="group hover:bg-white/10">
      <td>{peer.name}</td>
      <td><StatusBadge online={isOnline(peer.lastHandshake)} /></td>
      <!-- ... -->
      <td class="opacity-0 group-hover:opacity-100 transition-opacity">
        <button on:click={() => showQRCode(peer)}>QR</button>
        <button on:click={() => downloadConfig(peer)}>DL</button>
        <button on:click={() => onDelete(peer.id)}>Delete</button>
      </td>
    </tr>
  {/each}
</table>
```

**Acceptance Criteria**:
- ✅ Table displays all peer data
- ✅ Status badge shows online/offline correctly
- ✅ Action buttons hover-reveal on desktop
- ✅ Action buttons always visible on mobile
- ✅ Bytes formatted as human-readable
- ✅ Empty state shown if no peers
- ✅ Table updates reactively when peers store changes

**Files Created**:
- `src/lib/components/PeerTable.svelte` (new)

---

### Task 2.2: Create PeerModal Component - Add Form [90 min]
- [ ] Create file: `src/lib/components/PeerModal.svelte`
- [ ] Modal with form for adding peer:
  - Name input (required)
  - Allowed IPs (array of inputs, at least 1 required)
  - Add/Cancel buttons
- [ ] Form validation (FR-004, FR-014):
  - Name: required, non-empty
  - Allowed IPs: required, at least 1
  - Each IP: valid CIDR notation (use validateCIDR utility)
  - Show inline error messages under each field
- [ ] Client-side validation on blur or submit
- [ ] Submit button disabled while loading or validation errors
- [ ] **Clarification B**: Form stays open after success
  - Form resets/clears after successful submission
  - Does NOT close modal
- [ ] **Clarification A**: Optimistic updates
  - Call `addPeer()` which updates store immediately
  - If API succeeds: config modal opens
  - If API fails: error shown, form stays open with input preserved
- [ ] Show loading spinner during submission
- [ ] Config modal appears on success (separate component, Task 2.3)

**Code Pattern**:
```svelte
<script>
  import { addPeer } from '$lib/stores/peers';
  import { showNotification } from '$lib/stores/notifications';
  
  let formData = { name: '', allowedIPs: [''] };
  let errors = {};
  let isLoading = false;
  let configData = null;
  
  function validateForm() {
    errors = {};
    if (!formData.name.trim()) errors.name = 'Name required';
    if (formData.allowedIPs.length === 0) errors.allowedIPs = 'At least 1 IP';
    
    formData.allowedIPs.forEach((ip, i) => {
      const validation = validateCIDR(ip);
      if (!validation.valid) {
        errors[`ip_${i}`] = validation.error;
      }
    });
    return Object.keys(errors).length === 0;
  }
  
  async function handleSubmit() {
    if (!validateForm()) return;
    
    isLoading = true;
    try {
      const response = await addPeer(formData.name, formData.allowedIPs);
      showNotification('success', 'Peer added');
      configData = response; // Trigger config modal
      formData = { name: '', allowedIPs: [''] }; // Reset form
    } catch (error) {
      showNotification('error', error.message);
      // Form stays open with input preserved
    } finally {
      isLoading = false;
    }
  }
</script>

<dialog open>
  <form on:submit|preventDefault={handleSubmit}>
    <input bind:value={formData.name} placeholder="Peer name" />
    {#if errors.name}<span class="error">{errors.name}</span>{/if}
    
    {#each formData.allowedIPs as ip, i (i)}
      <input bind:value={ip} placeholder="10.0.0.2/32" />
      {#if errors[`ip_${i}`]}<span class="error">{errors[`ip_${i}`]}</span>{/if}
    {/each}
    
    <button on:click={() => formData.allowedIPs = [...formData.allowedIPs, '']}>
      Add IP
    </button>
    
    <button type="submit" disabled={isLoading || Object.keys(errors).length > 0}>
      {isLoading ? 'Adding...' : 'Add Peer'}
    </button>
    <button type="button" on:click>Cancel</button>
  </form>
  
  {#if configData}
    <QRCodeDisplay config={configData} on:close={() => configData = null} />
  {/if}
</dialog>
```

**Acceptance Criteria**:
- ✅ Form validates name and IPs
- ✅ Inline error messages displayed
- ✅ Submit button disabled on errors
- ✅ Form calls `addPeer()` on submit (optimistic)
- ✅ Form stays open after success
- ✅ Form resets after success
- ✅ Error shown on failure, form stays open
- ✅ Input preserved if error
- ✅ Config modal shown on success (Task 2.3)
- ✅ Multiple IPs can be added

**Files Created**:
- `src/lib/components/PeerModal.svelte` (new)

---

### Task 2.3: Create QRCodeDisplay Component [60 min]
- [ ] Create file: `src/lib/components/QRCodeDisplay.svelte`
- [ ] **Separate modal** from PeerModal (Clarification B - FR-006a)
- [ ] Display in overlay/side panel, not blocking form
- [ ] Show QR code (use qrcode library):
  - Generate from peer config string
  - Display as image/canvas
- [ ] Show WireGuard config as readable text:
  - Format: standard .conf file format
  - Include peer's private key (if available)
  - Include server details
- [ ] Display peer info summary:
  - Peer name
  - Public key (partial, for reference)
- [ ] "Download Config" button:
  - Create .conf file with proper format
  - Trigger browser download
- [ ] "Close" button to dismiss modal
- [ ] Clear sensitive data (private key) when modal closes

**Code Pattern**:
```svelte
<script>
  import QRCode from 'qrcode';
  
  export let config: string;
  export let peer = {};
  
  let qrCanvas;
  
  onMount(async () => {
    await QRCode.toCanvas(qrCanvas, config, { width: 200 });
  });
  
  function downloadConfig() {
    const element = document.createElement('a');
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + config);
    element.setAttribute('download', `${peer.name}.conf`);
    element.click();
  }
</script>

<div class="modal-overlay">
  <div class="modal-content">
    <canvas bind:this={qrCanvas}></canvas>
    <pre>{config}</pre>
    <button on:click={downloadConfig}>Download Config</button>
    <button on:click>Close</button>
  </div>
</div>
```

**Dependencies**:
- Install qrcode library: `npm install qrcode`

**Acceptance Criteria**:
- ✅ QR code generates from config
- ✅ Config displayed as readable text
- ✅ Download button saves .conf file
- ✅ Modal displays as overlay (doesn't block form)
- ✅ Close button dismisses modal
- ✅ Form visible behind modal

**Files Created**:
- `src/lib/components/QRCodeDisplay.svelte` (new)

---

### Task 2.4: Create Delete Confirmation & Optimistic Removal [45 min]
- [ ] Update `src/lib/components/PeerTable.svelte`:
  - Delete button triggers confirmation dialog (Task 2.5)
- [ ] Create/use `ConfirmDialog.svelte` (or modal)
- [ ] Dialog shows: "Remove [Peer Name]?"
- [ ] Confirm/Cancel options
- [ ] **Clarification A**: Optimistic delete
  - Call `deletePeer(id)` which removes from store immediately
  - API called in background
  - If success: peer stays removed
  - If failure: restore peer to list
- [ ] **Clarification D**: Context-aware error handling
  - If delete fails: show error notification with retry option
  - Peer restored to list
- [ ] Handle 404 idempotently (if peer already deleted, treat as success)
- [ ] Show loading state during deletion

**Code Pattern** (in Task 1.2, already implemented):
```typescript
export async function deletePeer(id: string) {
  const previousState = JSON.parse(JSON.stringify(get(peers)));
  
  // Optimistic: remove immediately
  peers.update(p => p.filter(x => x.id !== id));
  
  try {
    // API call
    await apiCall(`/peers/${id}`, { method: 'DELETE' });
    // Success: already removed
  } catch (error) {
    // Rollback
    peers.set(previousState);
    peersError.set(error.message);
    throw error;
  }
}
```

**Acceptance Criteria**:
- ✅ Confirmation dialog shown before delete
- ✅ Peer removed from list immediately (optimistic)
- ✅ Delete API called
- ✅ If success: peer stays removed
- ✅ If failure: peer restored to list
- ✅ Error notification shown on failure
- ✅ 404 treated as success

**Files Modified**:
- `src/lib/components/PeerTable.svelte`
- `src/lib/components/ConfirmDialog.svelte` (updated/created in Task 2.5)

---

### Task 2.5: Create ConfirmDialog Component [20 min]
- [ ] Create/update file: `src/lib/components/ConfirmDialog.svelte`
- [ ] Modal dialog with message, confirm/cancel buttons
- [ ] Accept props:
  - title: string (e.g., "Remove Peer?")
  - message: string (e.g., "Are you sure you want to remove Alice?")
  - onConfirm: callback function
  - onCancel: callback function
- [ ] Confirm button triggers `onConfirm()`
- [ ] Cancel button triggers `onCancel()`

**Acceptance Criteria**:
- ✅ Dialog displays title and message
- ✅ Confirm button calls callback
- ✅ Cancel button closes dialog
- ✅ Styling matches glassmorphism theme

**Files Created/Modified**:
- `src/lib/components/ConfirmDialog.svelte` (new/updated)

---

## Phase 3: Error Handling & Optimistic Updates (2 hours)

### Task 3.1: Implement Error Handling System [60 min]
- [ ] Already partially done in Task 1.2 (store error handling)
- [ ] Enhance error handling patterns:
  - **Form errors** (FR-012a): 
    - Validation errors: inline display in form
    - API errors: form stays open, input preserved, error toast shown
    - User can fix and retry without re-entering data
  - **Delete errors** (FR-012a):
    - Peer restored to list immediately
    - Error notification shown
    - "Retry" button available (calls deletePeer again)
  - **Fetch errors** (FR-012):
    - "Failed to load peers"
    - "Retry" button reloads data
  - **Timeout** (FR-012):
    - "Request timed out (10s)"
    - "Retry" button (handled by API client timeout)
- [ ] Create error notification with retry capability:
  - Notification component supports action button
  - Retry button re-triggers the failed operation
- [ ] Standardize error messages (user-friendly, no technical details)

**Error Handling Map**:
| Scenario | UI | State | Recovery |
|----------|-----|-------|----------|
| Add validation error | Form: inline errors | Form open | Fix + retry |
| Add API error | Toast: error msg | Form open | Retry or close |
| Delete API error | Toast: error + retry | Peer restored | Retry delete |
| Fetch timeout | Toast: "Timed out" + retry | Empty state | Retry load |

**Acceptance Criteria**:
- ✅ Form errors shown inline
- ✅ API errors shown in toast with retry
- ✅ Delete errors restore peer to list
- ✅ Retry functionality works
- ✅ No technical error messages exposed
- ✅ User can recover from all error states

**Files Modified**:
- `src/lib/stores/peers.ts` (error handling)
- `src/lib/stores/notifications.ts` (retry support)
- `src/lib/components/Notification.svelte` (action button)

---

### Task 3.2: Implement Optimistic Updates & Rollback [60 min]
- [ ] Already partially done in Task 1.2
- [ ] Verify optimistic update flow:
  1. User action → UI updates immediately (store change)
  2. API call in background
  3. On success → no additional change (already updated)
  4. On failure → rollback to previous state
- [ ] Previous state snapshot mechanism:
  - Before any mutation: `const prev = JSON.parse(JSON.stringify(get(store)))`
  - On failure: `store.set(prev)`
- [ ] Handle edge cases:
  - Multiple adds in quick succession (each with unique temp ID)
  - Add while delete pending (independent operations)
  - Race conditions (last-write-wins or queue operations)
- [ ] Test optimistic flow:
  - Add peer → appears immediately
  - API succeeds → peer stays with server data
  - API fails → peer removed, error shown
  - Delete peer → removed immediately
  - API fails → peer restored

**Optimistic Update Flow**:
```
User clicks "Add"
    ↓
Validation passes
    ↓
Store: add temp peer (UI updates instantly)
    ↓
API: POST /peers (in background)
    ↓
Success: Replace temp with server peer
    ↓
Failure: Remove temp peer, show error, form stays open
```

**Acceptance Criteria**:
- ✅ UI updates immediately on user action
- ✅ API call happens in background
- ✅ Rollback works: peer removed on add failure
- ✅ Rollback works: peer restored on delete failure
- ✅ Multiple operations don't conflict
- ✅ Previous state correctly restored

**Files Modified**:
- `src/lib/stores/peers.ts` (snapshot + rollback)

---

## Phase 4: Polish & Performance (1-2 hours)

### Task 4.1: Create Validation Utilities [30 min]
- [ ] Create file: `src/lib/utils/validation.ts`
- [ ] Export function: `validateCIDR(cidr: string): { valid: boolean; error?: string }`
  - Check format: must be "IP/PREFIX"
  - Validate IP: must be valid IP address (0.0.0.0 - 255.255.255.255)
  - Validate prefix: must be 0-32 (for IPv4)
  - Return: `{ valid: true }` or `{ valid: false, error: "Description" }`
  - Example errors:
    - "Format: 10.0.0.2/32"
    - "Invalid IP address"
    - "Prefix must be 0-32"
- [ ] Export function: `validateRequired(value: string): boolean`
  - Trim and check length > 0
- [ ] Export function: `sanitizeString(input: string): string`
  - Trim whitespace
  - Remove potentially dangerous characters (optional, basic sanitization)

**Code Pattern**:
```typescript
export function validateCIDR(cidr: string): { valid: boolean; error?: string } {
  const parts = cidr.split('/');
  if (parts.length !== 2) {
    return { valid: false, error: 'Format: 10.0.0.2/32' };
  }
  
  const [ip, prefix] = parts;
  
  // Validate IP
  const ipParts = ip.split('.');
  if (ipParts.length !== 4 || !ipParts.every(p => {
    const num = parseInt(p, 10);
    return !isNaN(num) && num >= 0 && num <= 255;
  })) {
    return { valid: false, error: 'Invalid IP address' };
  }
  
  // Validate prefix
  const prefixNum = parseInt(prefix, 10);
  if (isNaN(prefixNum) || prefixNum < 0 || prefixNum > 32) {
    return { valid: false, error: 'Prefix must be 0-32' };
  }
  
  return { valid: true };
}
```

**Acceptance Criteria**:
- ✅ CIDR validation catches invalid formats
- ✅ CIDR validation catches invalid IPs
- ✅ CIDR validation catches invalid prefixes
- ✅ Error messages are helpful
- ✅ Used in PeerModal form

**Files Created**:
- `src/lib/utils/validation.ts` (new)

---

### Task 4.2: Create Formatting Utilities [20 min]
- [ ] Create file: `src/lib/utils/formatting.ts`
- [ ] Export function: `formatBytes(bytes: number): string`
  - 0 → "0 B"
  - 1024 → "1 KB"
  - 1048576 → "1 MB"
  - 1073741824 → "1 GB"
  - 1099511627776 → "1 TB"
  - 1125899906842624 → "1 PB"
  - Use proper unit scaling
- [ ] Export function: `formatTimestamp(timestamp: string): string`
  - Parse ISO timestamp
  - Return relative time: "2 minutes ago", "Yesterday", "2 days ago"
  - Or return "Never" if timestamp is "0" or null
- [ ] Export function: `formatCIDR(ip: string): string` (optional)
  - Clean display of CIDR notation

**Code Pattern**:
```typescript
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const index = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = (bytes / Math.pow(1024, index)).toFixed(1);
  return `${value} ${units[index]}`;
}

export function formatTimestamp(timestamp: string): string {
  if (!timestamp || timestamp === '0') return 'Never';
  
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  
  if (diffSecs < 60) return 'Just now';
  if (diffSecs < 3600) return `${Math.floor(diffSecs / 60)} min ago`;
  if (diffSecs < 86400) return `${Math.floor(diffSecs / 3600)} hours ago`;
  if (diffSecs < 604800) return `${Math.floor(diffSecs / 86400)} days ago`;
  
  return date.toLocaleDateString();
}
```

**Acceptance Criteria**:
- ✅ Bytes formatted as KB, MB, GB, TB, PB
- ✅ Timestamps formatted as relative ("2 min ago")
- ✅ "Never" shown for null/0 timestamps
- ✅ Used in PeerTable and Dashboard

**Files Created**:
- `src/lib/utils/formatting.ts` (new)

---

### Task 4.3: Performance Audit [40 min]
- [ ] Run production build: `npm run build`
- [ ] Check bundle sizes:
  - `du -sh build/` (total)
  - `du -sh build/_app/immutable/chunks/` (chunk sizes)
  - Goal: < 200KB gzipped
- [ ] Run Lighthouse audit:
  - Open DevTools → Lighthouse
  - Run Performance audit
  - Check scores:
    - Performance: ≥90
    - Time to Interactive (TTI): <3s
    - First Contentful Paint (FCP): <1.5s
- [ ] If bundle too large:
  - Identify large dependencies with: `npm list`
  - Consider lazy-loading routes
  - Check for unused imports
- [ ] Record baseline metrics for comparison

**Acceptance Criteria**:
- ✅ Bundle size < 200KB gzipped
- ✅ Lighthouse Performance ≥90
- ✅ TTI < 3s
- ✅ FCP < 1.5s
- ✅ Metrics documented

**Files to Check**:
- `build/` directory structure
- `package.json` for dependency sizes

---

### Task 4.4: UI Polish & Glassmorphism Styling [30 min]
- [ ] Review `src/app.css` for global styles
- [ ] Verify glassmorphism applied to panels:
  - `backdrop-filter: blur(16px)`
  - `background: rgba(255, 255, 255, 0.1)` (primary)
  - `background: rgba(255, 255, 255, 0.05)` (secondary)
  - `border: 1px solid rgba(255, 255, 255, 0.1)`
  - `box-shadow` for depth
- [ ] Verify radial gradient background:
  - `background: radial-gradient(circle at 20% 20%, #1a2a3a 0%, #101922 100%)`
- [ ] Check responsive design:
  - Desktop (1024px+): Full sidebar, hover effects
  - Tablet (768px): Simplified layout
  - Mobile (320px): Stacked layout, bottom nav option
- [ ] Verify action button visibility:
  - Desktop: `opacity-0 group-hover:opacity-100`
  - Mobile: Always visible
- [ ] Check animations and transitions:
  - Smooth 60fps animations
  - Modal slide-in
  - Button hover states
- [ ] Verify empty state messaging:
  - "No peers configured" with clear call-to-action
  - Proper styling and spacing

**Glassmorphism Utility Classes** (should exist in app.css):
```css
.glass {
  backdrop-filter: blur(16px);
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
}

.glass-primary {
  background: rgba(255, 255, 255, 0.1);
}

.glass-secondary {
  background: rgba(255, 255, 255, 0.05);
}

.glass-hover:hover {
  background: rgba(255, 255, 255, 0.15);
}
```

**Acceptance Criteria**:
- ✅ Glassmorphism applied to all panels/modals
- ✅ Radial gradient background visible
- ✅ Responsive layout at all breakpoints
- ✅ Action buttons responsive (hover on desktop, visible on mobile)
- ✅ Animations smooth and 60fps
- ✅ Empty states styled properly

**Files Modified**:
- `src/app.css`
- Component files (add glass classes)

---

### Task 4.5: Final Testing Checklist [30 min]
- [ ] **Form Validation**:
  - Empty name → error shown, submit disabled
  - Empty IPs → error shown, submit disabled
  - Invalid CIDR → error shown with suggestion
  - Valid form → submit enabled
- [ ] **Error Recovery**:
  - Add fails → form stays open, input preserved
  - Delete fails → peer restored to list, error shown
  - Fetch fails → error message + retry button
- [ ] **Optimistic Updates**:
  - Add peer → appears immediately
  - Delete peer → removed immediately
  - Add failure → peer removed from list
  - Delete failure → peer restored to list
- [ ] **Responsive Design**:
  - Desktop (1024px+): Sidebar visible, hover effects work
  - Tablet (768px): Layout simplified, always-visible buttons
  - Mobile (320px): Stacked layout, navigation accessible
- [ ] **Performance**:
  - Initial load < 3s (TTI)
  - Peer operations < 500ms (optimistic + UI)
  - Lighthouse ≥90
- [ ] **Data Display**:
  - Bytes formatted (1.2 MB, 5.8 GB, etc.)
  - Timestamps formatted ("2 min ago", "Never")
  - Status badges correct (online < 120s, offline otherwise)

**Acceptance Criteria**:
- ✅ All manual tests pass
- ✅ No console errors
- ✅ Performance targets met
- ✅ Responsive at all breakpoints
- ✅ Ready for deployment

---

## Task Summary

| Phase | Tasks | Effort | Status |
|-------|-------|--------|--------|
| 0 | 0.1, 0.2, 0.3 | 30m | ⏳ |
| 1 | 1.1-1.9 | 2-3h | ⏳ |
| 2 | 2.1-2.5 | 3-4h | ⏳ |
| 3 | 3.1-3.2 | 2h | ⏳ |
| 4 | 4.1-4.5 | 1-2h | ⏳ |

**Total**: 27 tasks, 8.5-11.5 hours

---

## Dependency Graph

```
Phase 0 (Setup)
    ↓
Phase 1 (Core UI & Stores)
    ↓
Phase 2 (CRUD Operations)
    ↓
Phase 3 (Error Handling & Optimistic Updates)
    ↓
Phase 4 (Polish & Performance)
```

**Parallel within phases**:
- Phase 1: Tasks 1.1-1.9 mostly independent (after 1.1 API client)
- Phase 2: Tasks 2.1-2.5 mostly independent (depend on stores)
- Phase 3: Tasks 3.1-3.2 dependent (error handling feeds into optimistic)
- Phase 4: Tasks 4.1-4.5 independent (can start once Phases 2-3 complete)

---

## Ready to Code?

All tasks documented and ready. Choose your next step:

1. **Start Phase 0** (30 min) - Fix CI/CD, verify backend
2. **Start Phase 1** (2-3h) - Build core UI and stores
3. **Start specific task** - Let me know which one

**What's your preference?**
