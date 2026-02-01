# Frontend Blockers & Backend Requirements Matrix

**Date**: 2026-02-01  
**Purpose**: Track what the frontend needs from the backend and current status

---

## Blocker Summary

| # | Feature | Required By | Backend Status | Frontend Status | Action |
|---|---------|------------|-----------------|-----------------|--------|
| 1 | `/stats` endpoint with new properties | Dashboard | ✅ DONE | Ready | Deploy |
| 2 | CI/CD workflow validation | PR merge | ⚠️ NEEDS FIX | Blocked | URGENT: Enable CI |
| 3 | Peer regeneration | QR modal | ❌ Not started | Waiting | Design needed |
| 4 | Peer edit endpoint | Peer table | ❌ Not started | Waiting | Design needed |
| 5 | Real-time updates | Dashboard/Stats | ❌ Not started | Nice-to-have | Low priority |

---

## Detailed Blocker List

### 1. ✅ Stats Endpoint Properties (CLEARED)

**Required For**: Dashboard UI (`src/routes/+page.svelte`)

**What Frontend Needs**:
- `/stats` endpoint returning:
  - `publicKey` (Server's WireGuard public key)
  - `listenPort` (WireGuard listening port)
  - `subnet` (VPN subnet CIDR)
  - `peerCount` (Total connected peers)
  - `totalRx` / `totalTx` (Cumulative data transfer)

**Backend Status**: ✅ IMPLEMENTED
- `Stats` struct defined in `backend/internal/wireguard/wireguard.go`
- `realService.GetStats()` populates all fields
- `mockService.GetStats()` returns mock data
- Handler correctly returns JSON response

**Frontend Impact**: RESOLVED
- Dashboard can display interface configuration
- Stats cards show real data
- No frontend changes needed

**Action**: None - Use immediately

**Effort**: N/A (Already done)

---

### 2. ⚠️ CI/CD Workflow Validation (BLOCKING)

**Required For**: PR validation and automated testing

**Current Problem**:
- `.github/workflows/ci.yml` is fully commented out
- Cannot validate PRs
- Cannot ensure tests pass on merge
- Backend changes not validated before frontend integration

**What Needs to be Done**:
1. Uncomment all sections in `.github/workflows/ci.yml`
2. Validate Go version `1.25.6` (may need to update)
3. Test workflow locally or on feature branch
4. Verify backend tests run successfully
5. Ensure frontend build passes

**Frontend Impact**: CRITICAL
- Frontend can't be tested in CI
- PRs can't be automatically validated
- Risk of regressions

**Action**: HIGH PRIORITY
```bash
# Fix CI workflow
1. Edit .github/workflows/ci.yml
2. Uncomment all lines
3. Update Go version if needed
4. Test on feature branch
5. Commit and push
```

**Effort**: 15 minutes

**Acceptance Criteria**:
- ✅ CI workflow runs on PR creation
- ✅ Backend tests pass in CI
- ✅ Frontend build succeeds in CI
- ✅ Go version resolves correctly

---

### 3. ❌ Peer Regeneration Endpoint (OPTIONAL)

**Required For**: QR modal "Regenerate Keys" button

**Current State**: Placeholder button with no functionality

**What Frontend Needs**:
```
POST /peers/{id}/regenerate-keys
Response: {
  id: string,
  publicKey: string,
  name: string,
  privateKey: string,
  config: string  // WireGuard config with new keys
}
```

**Frontend Implementation**:
```svelte
// In QRCodeDisplay.svelte
async function regenerateKeys() {
  const response = await fetch(`/api/peers/${peerId}/regenerate-keys`, {
    method: 'POST'
  });
  // Update QR code with new config
}
```

**Backend Changes Needed**:
1. Add `RegeneratePeer(id string) (PeerResponse, error)` to Service interface
2. Implement in realService and mockService
3. Add handler method `Regenerate()`
4. Register route: `POST /peers/{id}/regenerate-keys`
5. Add tests

**Priority**: LOW (Feature not in current MVP)
- No UI designed for this yet
- Workaround: Delete peer and re-add

**Effort**: 2 hours (backend only)

**Decision**: Defer to Phase 2

---

### 4. ❌ Peer Edit Endpoint (OPTIONAL)

**Required For**: Edit peer metadata (name, allowed IPs)

**Current State**: No edit UI exists

**What Frontend Needs**:
```
PATCH /peers/{id}
Request: {
  name?: string,
  allowedIPs?: string[]
}
Response: {
  id: string,
  publicKey: string,
  name: string,
  allowedIPs: string[],
  ...
}
```

**Frontend Implementation**:
```svelte
// In PeerTable.svelte or PeerModal.svelte
async function updatePeer(id, updates) {
  const response = await fetch(`/api/peers/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(updates)
  });
  // Update peer in store
}
```

**Backend Changes Needed**:
1. Create `PeerUpdate` struct with optional fields
2. Add `UpdatePeer(id string, updates PeerUpdate) (Peer, error)` to Service interface
3. Implement in realService and mockService
4. Add handler method `Update()`
5. Register route: `PATCH /peers/{id}`
6. Add tests with CIDR validation

**Priority**: LOW (Workaround exists: delete + re-add)
- Not in current UI design
- Can defer to future release

**Effort**: 2.5 hours (backend only)

**Decision**: Defer to Phase 2

---

### 5. ❌ Real-Time Updates (OPTIONAL)

**Required For**: Live stats refresh without manual button click

**Current State**: Frontend polls `/stats` endpoint manually

**Options**:
1. **Polling** (Current): Frontend polls every N seconds
   - Effort: Frontend only (30 minutes)
   - Complexity: Low
   - Cost: More HTTP requests

2. **WebSocket** (Future): Real-time push updates
   - Effort: Backend (2-3 hours) + Frontend (1 hour)
   - Complexity: High
   - Cost: Single persistent connection

3. **Server-Sent Events (SSE)** (Future): HTTP streaming
   - Effort: Backend (1.5 hours) + Frontend (1 hour)
   - Complexity: Medium
   - Cost: One-way connection

**Priority**: LOW (Not required for MVP)
- Polling sufficient for management interface
- Can implement if UX demands real-time feedback

**Decision**: Defer to Phase 2

**Recommendation**: If needed, start with polling frontend-side (no backend changes needed).

---

## Implementation Roadmap

### Phase 1: Enable Production (THIS WEEK)
```
TASK-002: Enable CI/CD workflow                    [15 min]
├─ Uncomment .github/workflows/ci.yml
├─ Update Go version if needed
└─ Test on feature branch

STATUS: Backend ✅ Done | Frontend ⏳ Waiting for CI
```

### Phase 2: Complete Core Features (NEXT WEEK)
```
TASK-003: Peer regeneration endpoint              [2 hours]
TASK-004: Peer edit endpoint                      [2.5 hours]
└─ Design UI for edit modal
└─ Wire frontend components

STATUS: Defer until Phase 2 planning
```

### Phase 3: Polish & Optimization (FUTURE)
```
Real-time updates (WebSocket/SSE)                [3-4 hours]
Bulk operations (import/export)                  [TBD]
Backup/restore functionality                     [TBD]
```

---

## Current Dependencies

```
Frontend  →  Backend API
   ↓              ↓
PeerModal    POST /peers          ✅ Ready
PeerTable    GET /peers           ✅ Ready
             DELETE /peers/{id}   ✅ Ready
Dashboard    GET /stats           ✅ Ready
QRCodeDisplay needs /regenerate   ❌ Not started
PeerTable    needs PATCH /peers    ❌ Not started
Real-time UI needs WebSocket      ❌ Not started
```

---

## Recommendation

**Immediate Action** (Do Today - 15 minutes):
1. Enable CI/CD workflow
2. Verify backend tests pass
3. Start frontend integration testing with `/stats` endpoint

**Next Sprint** (If Time Permits - 4.5 hours):
1. Implement peer regeneration endpoint
2. Implement peer edit endpoint
3. Add UI for both features

**Blocked On**: CI/CD workflow only. Everything else is optional Phase 2 features.

---

## Notes

- Backend is **production-ready** for current frontend requirements
- All 4 core endpoints fully implemented and working
- Stats endpoint already includes requested properties (publicKey, listenPort, subnet)
- No hidden blockers - frontend can proceed with integration

---

*Last Updated: 2026-02-01*  
*Status: Ready for Phase 1 (CI/CD fix) → Phase 2 (optional features)*
