# Implementation Plan - Frontend Build
## WireGuard Manager SvelteKit Frontend

**Plan Date**: 2026-02-01  
**Based On**: Updated specification with 3 clarifications  
**Focus**: Unfinished business + clarified features  
**Current Status**: Ready to code

---

## Phase Overview

| Phase | Effort | Focus | Dependencies |
|-------|--------|-------|--------------|
| **0: Setup & Verification** | 30 min | Backend validation, build config | Backend API ready |
| **1: Core UI & Store Setup** | 2-3h | Pages, layout, stores, API client | Phase 0 complete |
| **2: Peer Management (CRUD)** | 3-4h | List, Add (with clarifications), Delete | Phase 1 complete |
| **3: Error Handling & Optimistic Updates** | 2h | Context-aware errors, rollback logic | Phase 2 complete |
| **4: Polish & Performance** | 1-2h | Styling refinement, Lighthouse audit | Phase 3 complete |

**Total Estimated Time**: 8.5-11.5 hours (spread across sessions)

---

## Phase 0: Setup & Verification (30 minutes)

### Task 0.1: Fix CI/CD Workflow [15 min] ⚠️ BLOCKING
**Current Status**: `.github/workflows/ci.yml` is commented out  
**Requirement**: Enable PR validation

**Checklist**:
- [ ] Read `.github/workflows/ci.yml` 
- [ ] Uncomment all lines
- [ ] Validate Go version (update if 1.25.6 doesn't exist)
- [ ] Commit and push to feature branch
- [ ] Verify GitHub Actions runs successfully

**Expected Outcome**: CI/CD validates all backend and frontend changes automatically

---

### Task 0.2: Verify Backend /stats Endpoint [10 min]
**Current Status**: Code implementation complete, tests need verification  
**Requirement**: Ensure all 7 stats fields return correctly

**Checklist**:
- [ ] Run `cd backend && go test ./... -v`
- [ ] Confirm test for `/stats` endpoint exists
- [ ] Verify test includes all fields: interfaceName, publicKey, listenPort, subnet, peerCount, totalRx, totalTx
- [ ] All tests pass locally
- [ ] (Commit result if test added)

**Expected Outcome**: Backend verified ready for frontend integration

---

### Task 0.3: Frontend Build Config Verification [5 min]
**Current Status**: Switched to adapter-static, SPA mode enabled  
**Requirement**: Confirm build succeeds and serves correctly

**Checklist**:
- [ ] Run `npm run build`
- [ ] Verify build succeeds (output in `build/` directory)
- [ ] Check bundle size: `du -sh build/`
- [ ] Verify output includes index.html with fallback

**Expected Outcome**: Frontend build pipeline validated

---

## Phase 1: Core UI & Store Setup (2-3 hours)

### Task 1.1: API Client Setup [30 min]
**Based On**: FR-011 (API utilities)  
**Files**: `src/lib/api/client.ts`

**Implementation**:
```typescript
// API base URL from environment (default: http://localhost:8080)
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Fetch wrapper with error handling
async function apiCall(endpoint: string, options?: RequestInit) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options
  });
  
  if (!response.ok) {
    throw new ApiError(response.status, await response.text());
  }
  return response.json();
}
```

**Checklist**:
- [ ] Create `src/lib/api/client.ts`
- [ ] Export `apiCall()` function with error handling
- [ ] Create `ApiError` class for consistent error handling
- [ ] Support timeout (10s per spec)
- [ ] Update `src/lib/api/peers.ts` and `stats.ts` to use client

**Expected Outcome**: Centralized API communication layer ready

---

### Task 1.2: Svelte Stores Setup [45 min]
**Based On**: FR-010 (State management), clarifications on optimistic updates  
**Files**: `src/lib/stores/peers.ts`, `src/lib/stores/stats.ts`, `src/lib/stores/notifications.ts`

**Peers Store Requirements**:
- ✅ Writable store with peer array
- ✅ Load function: `loadPeers()` - calls GET /peers
- ✅ Add function: `addPeer(name, allowedIPs)` - **OPTIMISTIC UPDATE**
  - Update store immediately
  - Call API in background
  - Rollback if API fails
- ✅ Delete function: `deletePeer(id)` - **OPTIMISTIC UPDATE**
  - Remove from store immediately
  - Call API in background
  - Restore if API fails
- ✅ Error handling with store for error state
- ✅ Previous state snapshot for rollback

**Stats Store Requirements**:
- ✅ Writable store with stats object
- ✅ Load function: `loadStats()` - calls GET /stats
- ✅ Auto-refresh on interval (or manual)

**Notifications Store Requirements**:
- ✅ Notification queue (array of notifications)
- ✅ Add, remove, clear functions
- ✅ Auto-dismiss capability

**Code Pattern**:
```typescript
// Optimistic add pattern
export async function addPeer(name: string, allowedIPs: string[]) {
  const tempPeer = { id: generateTempId(), name, allowedIPs, ... };
  peers.update(p => [...p, tempPeer]); // Optimistic
  
  try {
    const response = await apiCall('/peers', { 
      method: 'POST',
      body: JSON.stringify({ name, allowedIPs })
    });
    // Update store with server response (has real ID, timestamps, etc.)
    peers.update(p => p.map(x => x.id === tempPeer.id ? response : x));
    return response;
  } catch (error) {
    // Rollback on error
    peers.update(p => p.filter(x => x.id !== tempPeer.id));
    throw error;
  }
}
```

**Checklist**:
- [ ] Create/update `src/lib/stores/peers.ts`
- [ ] Implement optimistic add with rollback
- [ ] Implement optimistic delete with rollback
- [ ] Implement previous state snapshot for rollback
- [ ] Create/update `src/lib/stores/stats.ts`
- [ ] Create/update `src/lib/stores/notifications.ts`
- [ ] Test stores with mock API calls (manual testing in dev)

**Expected Outcome**: State management with optimistic updates ready

---

### Task 1.3: Pages & Layout [45 min]
**Based On**: FR-017 (Sidebar), FR-016 (Responsive)  
**Files**: `src/routes/+layout.svelte`, `src/routes/+page.svelte`, `src/routes/peers/+page.svelte`, `src/routes/stats/+page.svelte`

**Layout Implementation**:
- [ ] Sidebar navigation (persistent, left side)
- [ ] Navigation links: Dashboard (/), Peers (/peers), Settings (/settings), Logs (/logs - disabled)
- [ ] Responsive: desktop sidebar, tablet simplified, mobile bottom nav
- [ ] Usage widget at bottom (mock data)
- [ ] Load stats on mount and display in layout

**Dashboard Page** (`src/routes/+page.svelte`):
- [ ] Display quick stats cards (FR-008a)
  - Peer count (with icon)
  - Online peers count (derived from lastHandshake)
  - Total RX formatted (FR-009: formatBytes)
  - Total TX formatted
- [ ] 3-card horizontal grid layout
- [ ] Load stats on mount
- [ ] Show loading state while fetching

**Peers Page** (`src/routes/peers/+page.svelte`):
- [ ] Display peer list (FR-001)
- [ ] "Add New Peer" button
- [ ] Peer table/list implementation (delegated to Task 2.1)
- [ ] Load peers on mount
- [ ] Show loading state

**Stats Page** (`src/routes/stats/+page.svelte`):
- [ ] Display interface name, peer count, total RX/TX (FR-008)
- [ ] Show detailed stats (no charts yet)
- [ ] Auto-refresh or manual refresh button

**Checklist**:
- [ ] Update/create all layout and page files
- [ ] Add responsive CSS/Tailwind classes
- [ ] Wire up store subscriptions with `$store` syntax
- [ ] Handle loading/empty states
- [ ] Test responsive design at 320px, 768px, 1024px

**Expected Outcome**: Page structure and basic layout complete

---

## Phase 2: Peer Management CRUD (3-4 hours)

### Task 2.1: PeerTable Component [60 min]
**Based On**: FR-001, FR-001a, FR-002 (Status badge)  
**Files**: `src/lib/components/PeerTable.svelte`

**Features**:
- [ ] Display peer list in table format
- [ ] Columns: Name, Status (badge), Allowed IPs, Last Handshake, RX/TX bytes
- [ ] Status badge: online (green) if lastHandshake < 120s, offline (gray) otherwise
- [ ] Action buttons: View Config (QR), Download Config, Delete
  - Desktop (≥1024px): hover-reveal (opacity-0 → group-hover:opacity-100)
  - Mobile (<1024px): always visible
- [ ] Empty state: "No peers configured" with "Add Peer" button
- [ ] Subscribe to `$peers` store and reactively update

**Code Structure**:
```svelte
<script lang="ts">
  import { peers } from '$lib/stores/peers';
  import StatusBadge from './StatusBadge.svelte';
  
  function isOnline(lastHandshake: string): boolean {
    const time = new Date(lastHandshake).getTime();
    return Date.now() - time < 120_000; // 120 seconds
  }
</script>

<table>
  {#each $peers as peer (peer.id)}
    <tr class="group hover:bg-white/10">
      <td>{peer.name}</td>
      <td><StatusBadge online={isOnline(peer.lastHandshake)} /></td>
      <!-- ... -->
      <td class="opacity-0 group-hover:opacity-100 transition-opacity">
        <!-- Action buttons -->
      </td>
    </tr>
  {/each}
</table>
```

**Checklist**:
- [ ] Create component with store subscription
- [ ] Implement responsive action button visibility
- [ ] Format bytes and timestamps (use utilities from Task 1.4)
- [ ] Implement empty state
- [ ] Test with mock peer data

**Expected Outcome**: Peer table displays and reactive to store changes

---

### Task 2.2: PeerModal Component - Add Form [90 min]
**Based On**: FR-003, FR-004, FR-005, FR-006a, clarification on form staying open  
**Files**: `src/lib/components/PeerModal.svelte`

**Features**:
- [ ] Modal with form for adding peer
- [ ] Form fields: Name (required), Allowed IPs (array, CIDR notation)
- [ ] Client-side validation (FR-004, FR-014):
  - Name required and non-empty
  - At least one Allowed IP required
  - Each Allowed IP must be valid CIDR (use utility from Task 1.4)
  - Show inline error messages
- [ ] Submit button: "Add Peer"
- [ ] **Form stays open after success** (clarification B)
- [ ] Form resets/clears after successful submission
- [ ] Show loading state during API call
- [ ] Handle submission error: keep form open, show error message, preserve input

**Optimistic Update Integration** (Clarification A + FR-013a):
- [ ] Call `peers.addPeer()` which optimistically updates store
- [ ] Show success notification
- [ ] Trigger config modal display with response (FR-006a)
- [ ] On error: error message shown, form stays open with input preserved
- [ ] Implement rollback if API fails

**Code Pattern**:
```svelte
<script lang="ts">
  import { addPeer } from '$lib/stores/peers';
  
  let formData = { name: '', allowedIPs: [''] };
  let errors = {};
  
  async function handleSubmit() {
    errors = validateForm(formData); // FR-014
    if (Object.keys(errors).length > 0) return;
    
    try {
      const response = await addPeer(formData.name, formData.allowedIPs);
      // Success: form stays open, show config modal
      showConfigModal(response);
      formData = { name: '', allowedIPs: [''] }; // Reset form
    } catch (error) {
      // Error: form stays open, input preserved, error shown
      showNotification('error', error.message);
    }
  }
</script>
```

**Checklist**:
- [ ] Create modal with form
- [ ] Implement validation with error display
- [ ] Wire up optimistic add with error recovery
- [ ] Form stays open after success
- [ ] Form resets after success
- [ ] Inline error messages for each field
- [ ] Loading state during submission
- [ ] Test with valid/invalid inputs

**Expected Outcome**: Add peer form functional with clarified behavior

---

### Task 2.3: QRCodeDisplay Component [60 min]
**Based On**: FR-006, FR-006a  
**Files**: `src/lib/components/QRCodeDisplay.svelte`

**Features**:
- [ ] Display QR code of peer config (use qrcode library)
- [ ] Display WireGuard config as text (.conf format)
- [ ] "Download Config" button to save .conf file
- [ ] Peer details: name, public key (partial)
- [ ] Connection string/endpoint
- [ ] **Separate modal from add form** (clarification B - FR-006a)
- [ ] Close button to dismiss config modal
- [ ] Form modal stays open when config modal is closed

**Integration with PeerModal**:
```svelte
<!-- Parent: PeerModal.svelte -->
{#if showAddForm}
  <PeerFormComponent on:success={handleAddSuccess} />
{/if}

{#if showConfigModal}
  <QRCodeDisplay config={configData} on:close={() => showConfigModal = false} />
{/if}
```

**Checklist**:
- [ ] Create component with QR code generation
- [ ] Display config as readable text
- [ ] Implement download button
- [ ] Style as separate overlay/modal
- [ ] Handle sensitive key data (display once, clear on close)
- [ ] Test QR code scanability

**Expected Outcome**: QR code and config display in separate modal

---

### Task 2.4: Delete Peer Confirmation & Action [45 min]
**Based On**: FR-007, clarification on context-aware errors  
**Files**: Modify `PeerTable.svelte` with ConfirmDialog

**Features**:
- [ ] Delete button triggers confirmation dialog
- [ ] Dialog shows: "Are you sure you want to remove [Peer Name]?"
- [ ] Confirm/Cancel options
- [ ] **Optimistic delete** (Clarification A):
  - Remove peer from store immediately on confirmation
  - Call DELETE /peers/{id} in background
  - If fails: restore peer to list, show error notification (FR-012, context-aware)
- [ ] Handle 404 idempotently (if peer already deleted, treat as success)
- [ ] Show loading state during deletion

**Error Handling** (Clarification D - FR-012):
```typescript
try {
  await deletePeer(peerId);
  // Success: peer already removed optimistically
  showNotification('success', 'Peer removed');
} catch (error) {
  // Error: restore peer to list
  rollbackPeerDeletion(previousState);
  showNotification('error', 'Failed to remove peer. Please try again.');
}
```

**Checklist**:
- [ ] Create/use ConfirmDialog component
- [ ] Implement optimistic delete with rollback
- [ ] Store previous state for rollback
- [ ] Handle 404 as success
- [ ] Show error notification on failure
- [ ] Test delete flow and error recovery

**Expected Outcome**: Delete functionality with optimistic updates and error recovery

---

## Phase 3: Error Handling & Optimistic Updates (2 hours)

### Task 3.1: Error Handling System [60 min]
**Based On**: FR-012, FR-012a, Clarification D (Context-aware)  
**Files**: `src/lib/stores/notifications.ts`, error utilities

**Implementation**:
- [ ] Notification store with queue
- [ ] Notification component with toast UI
- [ ] Error notification patterns:
  - **Form errors** (FR-012a): Keep form open, display inline + toast
  - **Delete errors** (FR-012a): Restore peer to list, show toast with retry option
  - **Fetch errors** (FR-012): Show error message with retry button
  - **Timeout** (FR-012): "Request timed out" with retry option
- [ ] Auto-dismiss after 5s (configurable)
- [ ] Support error, warning, success, info types

**Error Recovery Patterns**:
```typescript
// Form error: keep form open
function validateForm(data) {
  const errors = {};
  if (!data.name) errors.name = 'Name is required';
  if (data.allowedIPs.length === 0) errors.allowedIPs = 'At least one IP required';
  return errors; // Display inline in form
}

// Delete error: restore from previous state
const previousState = JSON.parse(JSON.stringify($peers));
try {
  await deletePeer(peerId); // Optimistic remove
} catch (error) {
  peers.set(previousState); // Rollback
  showNotification('error', 'Failed to delete. Try again?');
}

// Fetch error: show retry button
async function retryFetch() {
  try {
    await loadPeers();
  } catch (error) {
    showNotification('error', error.message, { hasRetry: true });
  }
}
```

**Checklist**:
- [ ] Create notification store
- [ ] Create Notification.svelte component with toast UI
- [ ] Implement error utilities for consistent error handling
- [ ] Add retry mechanism to fetch errors
- [ ] Implement context-aware error responses (form vs delete vs fetch)
- [ ] Test all error scenarios

**Expected Outcome**: Comprehensive error handling with context-aware responses

---

### Task 3.2: Optimistic Updates & Rollback [60 min]
**Based On**: Clarification A, FR-013a, SC-005  
**Files**: Enhance `src/lib/stores/peers.ts`

**Implementation**:
- [ ] Peer store with previous state tracking
- [ ] Add peer: optimistic update + rollback on failure
- [ ] Delete peer: optimistic removal + rollback on failure
- [ ] State snapshot mechanism for rollback
- [ ] Handle race conditions:
  - Multiple adds in quick succession
  - Add while delete pending
  - Concurrent admin actions

**State Machine Pattern**:
```typescript
// Add peer: before → optimistic → success/rollback
peers.update(p => [...p, tempPeer]); // Optimistic (before)

fetch(...).then(() => {
  peers.update(p => p.map(x => x.id === tempPeer.id ? response : x)); // Confirmed
}).catch(() => {
  peers.update(p => p.filter(x => x.id !== tempPeer.id)); // Rollback
});
```

**Checklist**:
- [ ] Implement previous state snapshot in store
- [ ] Add peer with optimistic update + rollback
- [ ] Delete peer with optimistic removal + rollback
- [ ] Handle API delays and race conditions
- [ ] Test: add successful, add failed, delete successful, delete failed
- [ ] Test: quick successive adds/deletes

**Expected Outcome**: Optimistic updates with reliable rollback mechanism

---

## Phase 4: Polish & Performance (1-2 hours)

### Task 4.1: Validation Utilities [30 min]
**Based On**: FR-014  
**Files**: `src/lib/utils/validation.ts`

**Implement**:
- [ ] CIDR validation: `validateCIDR(ip: string): boolean`
- [ ] Required field validation
- [ ] String sanitization (no SQL injection, XSS)
- [ ] IP array validation (at least one valid, reject duplicates)

```typescript
export function validateCIDR(cidr: string): { valid: boolean; error?: string } {
  const parts = cidr.split('/');
  if (parts.length !== 2) return { valid: false, error: 'Format: 10.0.0.2/32' };
  
  const [ip, prefix] = parts;
  if (!isValidIP(ip)) return { valid: false, error: 'Invalid IP address' };
  if (!/^\d+$/.test(prefix) || +prefix < 0 || +prefix > 32) {
    return { valid: false, error: 'Prefix must be 0-32' };
  }
  return { valid: true };
}
```

**Checklist**:
- [ ] Create validation utilities
- [ ] Export for use in PeerModal component
- [ ] Add inline error messages for each validation rule
- [ ] Test with edge cases (10.0.0.5, 10.0.0.5/32, 10.0.0.5/33, etc.)

---

### Task 4.2: Formatting Utilities [20 min]
**Based On**: FR-009  
**Files**: `src/lib/utils/formatting.ts`

**Implement**:
- [ ] `formatBytes(bytes: number): string` → "1.2 MB", "5.8 GB", etc.
- [ ] `formatTimestamp(timestamp: string): string` → "2 minutes ago", "Yesterday", etc.
- [ ] `formatCIDR(ip: string): string` → "10.0.0.2/32" (clean display)

**Checklist**:
- [ ] Create formatting utilities
- [ ] Use in PeerTable and Stats displays
- [ ] Handle edge cases (0 bytes, very large numbers, null timestamps)
- [ ] Test formatting for readability

---

### Task 4.3: Performance Audit [40 min]
**Based On**: FR and Constitution V (Performance budgets)  
**Targets**: TTI <3s, FCP <1.5s, bundle <200KB, Lighthouse ≥90

**Checklist**:
- [ ] Run `npm run build`
- [ ] Check bundle size: `du -sh build/_app/`
- [ ] Run Lighthouse audit: DevTools → Lighthouse
- [ ] Verify metrics:
  - [ ] TTI <3s
  - [ ] FCP <1.5s  
  - [ ] Performance score ≥90
- [ ] Identify large chunks if needed (remove unused deps, lazy-load routes)
- [ ] Check animation smoothness (60fps) in DevTools Performance

**Checklist**:
- [ ] Bundle size within budget
- [ ] Lighthouse score ≥90
- [ ] No performance regressions vs baseline

---

### Task 4.4: UI Polish & Styling [30 min]
**Based On**: FR-019, FR-020 (Glassmorphism design)  
**Files**: `src/app.css`, component styling

**Checklist**:
- [ ] Verify glassmorphism applied to panels (blur, opacity, border)
- [ ] Verify radial gradient backgrounds
- [ ] Check responsive design at key breakpoints (320px, 768px, 1024px)
- [ ] Verify action button visibility (hover on desktop, always visible mobile)
- [ ] Check empty state messaging and styling
- [ ] Verify loading spinners and animations

---

## Unfinished Business Tracker

### From Previous Session
| Item | Status | Impact | Priority |
|------|--------|--------|----------|
| UI mockup visual alignment | ⏳ In Progress | Dashboard/components don't match Figma | HIGH |
| CI/CD workflow disabled | ⏳ To Do | Blocks PR validation | CRITICAL |
| Backend stats test coverage | ⏳ To Do | Verify /stats returns all fields | MEDIUM |
| Backend peer regeneration | ❌ Not started | Nice-to-have (not blocking) | LOW |
| Backend peer edit endpoint | ❌ Not started | Nice-to-have (not blocking) | LOW |

### Updated with Clarifications
| Clarification | Implementation Task | Priority |
|---------------|-------------------|----------|
| Context-aware error handling | Task 3.1: Error Handling System | HIGH |
| Optimistic UI updates | Task 3.2: Optimistic Updates & Rollback | HIGH |
| Form modal behavior (stay open) | Task 2.2: PeerModal + Task 2.3: QRCodeDisplay | HIGH |

---

## Dependencies & Blockers

### Must Complete Before Starting
- ✅ Phase 0 (Backend verification, CI/CD fix)
- ✅ Phase 1 (Stores, API client, pages)

### Phase 2 Blockers
- Phase 1 must be complete
- Stores fully functional

### Phase 3 Blockers
- Phase 2 must be complete
- Optimistic updates need working stores

### Phase 4 Blockers
- Phase 3 must be complete
- All core functionality working

---

## Testing Strategy (Manual - per Constitution II)

### Form Validation Testing
- [ ] Empty name → error message displayed
- [ ] Empty allowed IPs → error message displayed
- [ ] Invalid CIDR (10.0.0.5) → error message, suggestion shown
- [ ] Valid form → submission enabled
- [ ] Form submission → optimistic update visible

### Error Recovery Testing
- [ ] Add peer, API fails → form stays open, input preserved, error shown, can retry
- [ ] Delete peer, API fails → peer restored to list, error notification shown
- [ ] Fetch peers timeout → error message with retry button visible
- [ ] Fetch succeeds → data displayed correctly

### Optimistic Update Testing
- [ ] Add peer → appears in list immediately (before API response)
- [ ] API response 500 → peer removed from list, error shown
- [ ] Delete peer → removed from list immediately (before API response)
- [ ] API response 404 → treat as success (peer already gone)

### Responsive Design Testing
- [ ] Desktop (1024px+): Sidebar visible, hover-reveal action buttons
- [ ] Tablet (768px): Simplified sidebar, always-visible buttons
- [ ] Mobile (320px): Bottom nav, stacked cards, action icons visible

### Performance Testing
- [ ] Initial load TTI <3s (measure with DevTools)
- [ ] Add peer → appears in UI <500ms (optimistic)
- [ ] Delete peer → removed from UI <500ms (optimistic)
- [ ] List refresh on mount <1s (API response)
- [ ] Bundle size <200KB gzipped

---

## File Structure After Implementation

```
src/
├── routes/
│   ├── +layout.svelte          (Sidebar, navigation)
│   ├── +page.svelte             (Dashboard with stats cards)
│   ├── peers/
│   │   └── +page.svelte         (Peer list page)
│   └── stats/
│       └── +page.svelte         (Stats page)
├── lib/
│   ├── api/
│   │   ├── client.ts            (API utilities)
│   │   ├── peers.ts             (Peer API functions)
│   │   └── stats.ts             (Stats API functions)
│   ├── components/
│   │   ├── PeerTable.svelte      (Peer list table)
│   │   ├── PeerModal.svelte      (Add peer form modal)
│   │   ├── QRCodeDisplay.svelte  (Config + QR display)
│   │   ├── StatusBadge.svelte    (Online/offline indicator)
│   │   ├── Notification.svelte   (Toast notifications)
│   │   ├── Sidebar.svelte        (Navigation sidebar)
│   │   └── ConfirmDialog.svelte  (Confirmation dialog)
│   ├── stores/
│   │   ├── peers.ts             (Peer list store + optimistic updates)
│   │   ├── stats.ts             (Stats store)
│   │   └── notifications.ts      (Notification queue)
│   ├── types/
│   │   ├── peer.ts
│   │   ├── stats.ts
│   │   ├── api.ts
│   │   └── notification.ts
│   └── utils/
│       ├── validation.ts         (CIDR, required fields, sanitization)
│       ├── formatting.ts         (bytes, timestamps)
│       └── constants.ts
└── app.css                       (Global styles + glassmorphism)
```

---

## Success Criteria & Acceptance

### Must Have (Phase 2 Complete)
- ✅ View all peers (FR-001)
- ✅ Status badge (online/offline) (FR-002)
- ✅ Add peer form (FR-003, FR-004, FR-005)
- ✅ Delete peer with confirmation (FR-007)
- ✅ View stats (FR-008)
- ✅ Config/QR code display (FR-006)

### Must Have (Phase 3 Complete)
- ✅ Context-aware error handling (FR-012, Clarification D)
- ✅ Form stays open after success (Clarification B)
- ✅ Optimistic updates with rollback (Clarification A)
- ✅ Error recovery & retry (FR-012a)

### Nice to Have (Phase 4)
- ✅ Performance: TTI <3s, bundle <200KB, Lighthouse ≥90
- ✅ Responsive: 320px mobile to 1920px+ desktop
- ✅ Styling: Glassmorphism, consistent design

---

## Timeline

**Optimal:** 8.5-11.5 hours over 2-3 sessions
- Session 1 (This): Phase 0 (30 min) + Phase 1 (2-3h)
- Session 2: Phase 2 (3-4h)
- Session 3: Phase 3 (2h) + Phase 4 (1-2h)

---

**Ready to begin Phase 0 & Phase 1?** ✅

Or would you like to:
- Review this plan in detail?
- Adjust timeline/priorities?
- Start with Phase 2 (assuming Phase 1 is done)?
