# Backend Blockers Analysis & Implementation Plan

**Date**: 2026-02-01  
**Current Phase**: Frontend Integration & Backend Completion  
**Status**: Planning Phase

---

## Executive Summary

The backend API is **70% complete**. Core peer management endpoints are implemented (`GET /peers`, `POST /peers`, `DELETE /peers`, `GET /stats`). However, the `/stats` endpoint implementation has been completed in code but **may need verification testing** to ensure it returns the documented properties correctly. Additionally, secondary features are not yet implemented.

### Backend Implementation Status

| Component | Status | Notes |
|-----------|--------|-------|
| Core `/peers` endpoints | ✅ Implemented | List, Add, Remove peers |
| `/stats` endpoint | ✅ Code Ready | Includes publicKey, listenPort, subnet |
| `/stats` endpoint | ⚠️ Untested | Needs verification in tests |
| QR Code generation | ✅ Implemented | Part of add peer response |
| Peer regeneration | ❌ Not Started | [NEEDS IMPL] |
| Real-time WebSocket | ❌ Not Started | Low priority |
| Bulk operations | ❌ Not Started | Low priority |
| CI/CD workflow | ⚠️ Broken | Comments out entirely |

---

## Critical Blockers (Must Fix for Frontend Testing)

### 1. ✅ Stats Endpoint Properties Implementation

**Status**: ✅ COMPLETED IN CODE

**What's Done**:
- `Stats` struct includes: `publicKey`, `listenPort`, `subnet` (wireguard.go)
- `realService.GetStats()` populates all fields from `wgctrl.Device()` and config
- Handler correctly returns Stats as JSON

**What's Missing**:
- **Test verification**: No test in `main_test.go` verifies `/stats` endpoint returns correct properties
- **Mock service sync**: `mockService.GetStats()` may not include new properties

**Task**: [TASK-001] Verify Stats endpoint tests and update mock service if needed

**Effort**: 30 minutes

---

### 2. ⚠️ CI/CD Workflow Disabled

**Status**: BROKEN - BLOCKING PR MERGE

**Problem**:
- `.github/workflows/ci.yml` is fully commented out
- Cannot run automated tests
- Cannot validate PR merges

**Blockers**:
- Go version `1.25.6` specified in ci.yml may not exist (verify against official releases)
- Backend tests need to run before frontend integration

**Task**: [TASK-002] Fix CI workflow - validate Go version and enable tests

**Effort**: 45 minutes

---

## Secondary Blockers (Needed for Full Feature Parity)

### 3. ❌ Peer Regeneration Endpoint

**Status**: NOT IMPLEMENTED

**What's Needed**:
- New endpoint: `POST /peers/{id}/regenerate-keys`
- Generate new keypair for existing peer
- Return new configuration with QR code
- Update peer metadata if name changed

**Backend Implementation**:
1. Add `RegeneratePeer()` method to `Service` interface
2. Implement in both `realService` and `mockService`
3. Create handler method `Regenerate()`
4. Add route in `main.go`
5. Add tests in `main_test.go`

**Frontend Integration**:
- Wire up to QRCodeDisplay "Regenerate Keys" button in `src/lib/components/QRCodeDisplay.svelte`

**Task**: [TASK-003] Implement peer key regeneration endpoint

**Effort**: 2 hours

---

### 4. ❌ Peer Edit Endpoint

**Status**: NOT IMPLEMENTED

**What's Needed**:
- Endpoint to update peer metadata (name, allowed IPs)
- `PATCH /peers/{id}` or `PUT /peers/{id}`
- Validate CIDR notation for allowed IPs
- Update both metadata and WireGuard config

**Backend Implementation**:
1. Add `UpdatePeer()` method to `Service` interface
2. Implement in both `realService` and `mockService`
3. Create handler method `Update()`
4. Add route in `main.go`
5. Add tests in `main_test.go`

**Frontend Integration**:
- Add "Edit" button in PeerTable
- Create edit modal in `src/lib/components/PeerModal.svelte`

**Priority**: MEDIUM (workaround: delete + re-add)

**Task**: [TASK-004] Implement peer edit/update endpoint

**Effort**: 2.5 hours

---

## Testing & Verification Gaps

### 5. Tests for Stats Endpoint

**What's Missing**:
- Test case for `/stats` endpoint with new properties
- Verify `publicKey` matches device public key
- Verify `listenPort` is correctly populated
- Verify `subnet` loads from config

**Task**: [TASK-001a] Add test cases for Stats endpoint

**Effort**: 30 minutes

---

### 6. Tests for Error Cases

**What's Missing**:
- Tests for edge cases (invalid CIDR, duplicate peer, etc.)
- Tests for permission errors (when running without root)
- Fallback to mock service behavior

**Task**: [TASK-005] Add comprehensive error case tests

**Effort**: 1.5 hours

---

## Infrastructure Tasks

### 7. CI/CD Pipeline

**What's Needed**:
1. Validate Go version (1.25.6 or update to valid version)
2. Uncomment and test each job in ci.yml
3. Ensure backend tests run before frontend build
4. Add frontend build step if not present
5. Add Docker build for deployment verification

**Task**: [TASK-002] Fix and enable CI workflow

**Effort**: 45 minutes (validation + fixes)

---

## Implementation Priority & Order

### Phase 1: Make Frontend Testing Possible (THIS SESSION)
```
TASK-001: Verify Stats endpoint & update mock service    [30 min]
TASK-002: Fix CI workflow and enable tests               [45 min]
────────────────────────────────────────────────────────
Total Phase 1: 1 hour 15 minutes
```

**Outcome**: Frontend can test with real data from `/stats` endpoint. CI/CD validates all changes.

---

### Phase 2: Complete Core Features (NEXT SESSION)
```
TASK-003: Implement peer regeneration endpoint           [2 hours]
TASK-004: Implement peer edit/update endpoint            [2.5 hours]
────────────────────────────────────────────────────────
Total Phase 2: 4.5 hours
```

**Outcome**: All CRUD operations complete. Frontend has full feature parity with design.

---

### Phase 3: Testing & Polish (OPTIONAL)
```
TASK-005: Add comprehensive error case tests             [1.5 hours]
```

**Outcome**: Robust error handling and edge case coverage.

---

## Detailed Task Breakdown

### TASK-001: Verify Stats Endpoint Implementation

**Description**: Ensure `/stats` endpoint returns all documented properties

**Steps**:
1. Run backend tests: `cd backend && go test ./...`
2. Check test output for `/stats` endpoint tests
3. Add test case if missing:
   ```go
   func TestStatsEndpoint(t *testing.T) {
       // GET /stats
       // Verify response contains: publicKey, listenPort, subnet, peerCount, totalRx, totalTx
   }
   ```
4. Update `mockService.GetStats()` to match Stats struct with all fields
5. Run tests again to verify

**Files to Modify**:
- `backend/cmd/server/main_test.go` (add/verify test)
- `backend/internal/wireguard/mock.go` (update GetStats method if needed)

**Acceptance Criteria**:
- ✅ Test for `/stats` endpoint exists and passes
- ✅ All 7 fields returned in response (interfaceName, publicKey, listenPort, subnet, peerCount, totalRx, totalTx)
- ✅ Mock service returns valid stats

---

### TASK-002: Fix CI/CD Workflow

**Description**: Enable GitHub Actions CI workflow for automated testing

**Steps**:
1. Check Go version availability: Visit https://github.com/actions/setup-go
2. Verify if Go 1.25.6 exists or update to latest stable version
3. Uncomment all lines in `.github/workflows/ci.yml`
4. Test workflow locally (or push to feature branch)
5. Verify backend tests pass in CI
6. Verify frontend build succeeds in CI

**Files to Modify**:
- `.github/workflows/ci.yml` (uncomment all sections)

**Acceptance Criteria**:
- ✅ CI workflow uncommented and valid YAML
- ✅ Go version resolved (1.25.6 or valid alternative)
- ✅ CI runs on PR creation
- ✅ All backend tests pass in CI
- ✅ Frontend build succeeds in CI

---

### TASK-003: Implement Peer Regeneration Endpoint

**Description**: Allow users to regenerate keypair for existing peer

**Backend API Design**:
```
POST /peers/{id}/regenerate-keys
Response: PeerResponse (with new privateKey and config)
```

**Implementation Steps**:
1. Add `RegeneratePeer(id string) (PeerResponse, error)` to `Service` interface
2. Implement in `realService`:
   - Remove old peer from WireGuard
   - Generate new keypair
   - Add peer back with same metadata (name, allowedIPs)
   - Return new config
3. Implement in `mockService` with same logic
4. Add handler method `Regenerate()` in handlers.go
5. Register route in main.go: `mux.HandleFunc("POST /peers/{id}/regenerate-keys", handler.Regenerate)`
6. Add comprehensive tests in main_test.go

**Files to Modify**:
- `backend/internal/wireguard/wireguard.go` (add to Service interface)
- `backend/internal/wireguard/service.go` (implement in realService)
- `backend/internal/wireguard/mock.go` (implement in mockService)
- `backend/internal/handlers/handlers.go` (add Regenerate method)
- `backend/cmd/server/main.go` (register route)
- `backend/cmd/server/main_test.go` (add tests)
- `backend/API.md` (document endpoint)

**Acceptance Criteria**:
- ✅ Endpoint accepts POST request
- ✅ Returns new keypair and config
- ✅ Old keypair is replaced in WireGuard
- ✅ Metadata preserved (name, allowedIPs)
- ✅ Tests verify new config is different from old
- ✅ Error handling for invalid peer ID

---

### TASK-004: Implement Peer Edit Endpoint

**Description**: Allow users to update peer metadata and configuration

**Backend API Design**:
```
PATCH /peers/{id}
Request: { name?: string, allowedIPs?: string[] }
Response: Peer (updated)
```

**Implementation Steps**:
1. Add `UpdatePeer(id string, updates PeerUpdate) (Peer, error)` to `Service` interface
2. Create `PeerUpdate` struct with optional fields
3. Implement in `realService`:
   - Validate new allowedIPs (CIDR notation)
   - Update WireGuard config
   - Update metadata in storage
   - Return updated peer
4. Implement in `mockService` with same logic
5. Add handler method `Update()` in handlers.go
6. Register route in main.go: `mux.HandleFunc("PATCH /peers/{id}", handler.Update)`
7. Add comprehensive tests in main_test.go

**Files to Modify**:
- `backend/internal/wireguard/wireguard.go` (add PeerUpdate struct and to Service interface)
- `backend/internal/wireguard/service.go` (implement in realService)
- `backend/internal/wireguard/mock.go` (implement in mockService)
- `backend/internal/handlers/handlers.go` (add Update method)
- `backend/cmd/server/main.go` (register route)
- `backend/cmd/server/main_test.go` (add tests)
- `backend/API.md` (document endpoint)

**Acceptance Criteria**:
- ✅ Endpoint accepts PATCH request with optional fields
- ✅ Can update peer name
- ✅ Can update allowedIPs (with CIDR validation)
- ✅ Unspecified fields remain unchanged
- ✅ Metadata and WireGuard config stay in sync
- ✅ Tests verify individual and combined updates
- ✅ Error handling for invalid peer ID, invalid CIDR

---

### TASK-005: Add Comprehensive Error Case Tests

**Description**: Test edge cases and error conditions

**Test Cases to Add**:
- Duplicate peer (same publicKey)
- Invalid CIDR notation
- Empty allowedIPs array
- Whitespace-only name
- Missing required fields
- Permission errors (if running without CAP_NET_ADMIN)
- Nonexistent peer operations

**Files to Modify**:
- `backend/cmd/server/main_test.go` (add test cases)

**Acceptance Criteria**:
- ✅ All error cases return appropriate HTTP status codes
- ✅ Error messages are user-friendly
- ✅ No panics or unhandled errors
- ✅ 90%+ test coverage for handlers

---

## Known Issues & Notes

1. **Mock Service Sync**: Mock service may not include new Stats properties. Check and update if needed.
2. **Go Version**: Verify Go 1.25.6 availability or use latest stable (1.25.x or 1.26.x).
3. **Root Permissions**: Tests may fail if not running with CAP_NET_ADMIN. Mock service fallback handles this.
4. **Performance**: Current file-based storage acceptable for <100 peers. May need database migration for larger deployments.

---

## Summary Table

| Task | Effort | Priority | Status | Blocker |
|------|--------|----------|--------|---------|
| TASK-001: Verify Stats | 30m | CRITICAL | TODO | YES |
| TASK-002: Fix CI | 45m | CRITICAL | TODO | YES |
| TASK-003: Peer Regeneration | 2h | HIGH | TODO | NO |
| TASK-004: Peer Edit | 2.5h | HIGH | TODO | NO |
| TASK-005: Error Tests | 1.5h | MEDIUM | TODO | NO |

**Total Effort**: ~7 hours (1.25h critical, 4.5h features, 1.5h testing)

**This Session Goal**: Complete TASK-001 & TASK-002 (1h 15m) to unblock frontend testing.

---

*Last Updated: 2026-02-01*
