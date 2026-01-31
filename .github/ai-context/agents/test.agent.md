# Testing-Focused Agent Rules

**Purpose**: Guidelines for AI agents working on tests and test-related code.

## Testing Philosophy by Component

### Backend (Go) — TDD Mandatory

**Per Constitution Principle I**: Backend MUST follow Test-Driven Development.

**Workflow**:
1. **Write test first** — Before any implementation
2. **Verify test fails** — Red (test should fail initially)
3. **Implement minimal code** — Make test pass
4. **Verify test passes** — Green (test should pass now)
5. **Refactor** — Clean up code while keeping tests green
6. **Repeat** — Next test

**Non-negotiable**: All backend handlers, services, and middleware must have tests.

### Frontend (Svelte) — No Tests Required

**Per Constitution Principle II**: Frontend prioritizes UX and performance over test coverage.

**Rationale**: 
- TypeScript provides type safety
- ESLint catches common errors
- Manual testing catches UI bugs faster than automated tests
- Test brittleness for UI changes outweighs benefits

**What to do instead**:
- Manual testing in dev mode (`npm run dev`)
- Type safety via TypeScript
- Performance audits (Lighthouse)
- Visual QA in browser

## Backend Testing Patterns

### Unit Tests for Handlers

**Location**: `backend/cmd/server/*_test.go`

**Pattern**:
```go
func TestHandlerName(t *testing.T) {
    // Arrange: Setup mock service, test data
    mockService := wireguard.NewMockService()
    handler := handlers.NewPeerHandler(mockService)
    
    // Act: Create request, call handler
    req := httptest.NewRequest("GET", "/peers", nil)
    rr := httptest.NewRecorder()
    handler.List(rr, req)
    
    // Assert: Verify response
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }
    
    var peers []wireguard.Peer
    if err := json.NewDecoder(rr.Body).Decode(&peers); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }
}
```

**Key points**:
- Use `httptest.NewRequest()` and `httptest.NewRecorder()`
- Test both success and error cases
- Verify HTTP status codes
- Verify response body structure
- Use mock service (no actual WireGuard operations)

### Contract Tests for Services

**Location**: `backend/internal/wireguard/*_test.go`

**Purpose**: Verify service interface behavior

**Pattern**:
```go
func TestServiceAddPeer(t *testing.T) {
    // Test both realService and mockService implement interface correctly
    services := []wireguard.Service{
        wireguard.NewMockService(),
        // realService requires WireGuard, skip in CI
    }
    
    for _, svc := range services {
        peer, err := svc.AddPeer("Test", "", []string{"10.0.0.2/32"})
        if err != nil {
            t.Fatalf("AddPeer failed: %v", err)
        }
        if peer.Name != "Test" {
            t.Errorf("expected name 'Test', got %s", peer.Name)
        }
    }
}
```

### Integration Tests

**Purpose**: Verify end-to-end functionality (API → Service → WireGuard)

**When**: Testing with real WireGuard kernel module (optional, not in CI)

**Pattern**:
```go
func TestIntegrationAddRemovePeer(t *testing.T) {
    if os.Getenv("INTEGRATION_TESTS") != "true" {
        t.Skip("Skipping integration test (set INTEGRATION_TESTS=true)")
    }
    
    // Real WireGuard service
    svc, err := wireguard.NewRealService("wg0", "./test_peers.json", "server:51820", "pubkey")
    if err != nil {
        t.Fatalf("Failed to init real service: %v", err)
    }
    defer svc.Close()
    
    // Test add/remove peer flow
    peer, err := svc.AddPeer("IntegrationTest", "", []string{"10.0.0.99/32"})
    if err != nil {
        t.Fatalf("AddPeer failed: %v", err)
    }
    
    // Verify peer exists
    peers, _ := svc.ListPeers()
    found := false
    for _, p := range peers {
        if p.ID == peer.ID {
            found = true
            break
        }
    }
    if !found {
        t.Error("Peer not found after adding")
    }
    
    // Clean up
    if err := svc.RemovePeer(peer.ID); err != nil {
        t.Fatalf("RemovePeer failed: %v", err)
    }
}
```

## Test Coverage Requirements

### Backend

**Required coverage**:
- ✅ All HTTP handlers (List, Add, Remove, Stats)
- ✅ All service methods (AddPeer, RemovePeer, ListPeers, GetStats)
- ✅ Input validation logic
- ✅ Error handling paths

**Optional coverage**:
- Middleware (logging, CORS) — Nice to have, not critical
- Configuration loading — Already tested by usage
- Storage layer — Covered by service tests

**Target**: >80% code coverage for `handlers/` and `wireguard/` packages

### Frontend

**Required coverage**: None (per Constitution Principle II)

**Alternative quality measures**:
- TypeScript compilation without errors
- ESLint passes without warnings
- Manual testing of all user flows
- Lighthouse performance score ≥90

## Writing Good Tests

### Principle 1: Tests Should Be Readable

**Good test names**:
```go
func TestPeerHandlerAdd_ValidInput_ReturnsCreated(t *testing.T) {}
func TestPeerHandlerAdd_MissingName_ReturnsBadRequest(t *testing.T) {}
func TestPeerHandlerAdd_InvalidCIDR_ReturnsBadRequest(t *testing.T) {}
```

**Bad test names**:
```go
func TestAdd(t *testing.T) {}  // Too vague
func Test1(t *testing.T) {}    // Meaningless
```

### Principle 2: Tests Should Be Independent

**Each test should**:
- Set up its own data
- Not depend on other tests' side effects
- Clean up after itself
- Be runnable in any order

**Example**:
```go
func TestExample(t *testing.T) {
    // Arrange: Fresh state for this test
    mockService := wireguard.NewMockService()
    
    // Act: Perform operation
    // ...
    
    // Assert: Verify outcome
    // ...
    
    // (No cleanup needed for mock service, but would for real resources)
}
```

### Principle 3: Test Behavior, Not Implementation

**Good**: Test what the function does (API contract)
```go
func TestAddPeer_ValidInput_ReturnsPeer(t *testing.T) {
    // Test: Given valid input, function returns peer with expected properties
}
```

**Bad**: Test how the function does it (internal details)
```go
func TestAddPeer_CallsWgctrlConfigure(t *testing.T) {
    // ❌ Tests internal implementation, brittle
}
```

### Principle 4: Test Edge Cases

**Always test**:
- Happy path (valid input, expected output)
- Invalid input (missing fields, malformed data)
- Boundary conditions (empty lists, max values)
- Error conditions (network failures, permission denied)

**Example**:
```go
func TestPeerHandlerAdd_EdgeCases(t *testing.T) {
    tests := []struct {
        name       string
        input      AddPeerRequest
        wantStatus int
        wantErr    string
    }{
        {"Valid", AddPeerRequest{Name: "Test", AllowedIPs: []string{"10.0.0.2/32"}}, 201, ""},
        {"MissingName", AddPeerRequest{AllowedIPs: []string{"10.0.0.2/32"}}, 400, "Name is required"},
        {"EmptyAllowedIPs", AddPeerRequest{Name: "Test"}, 400, "At least one AllowedIP is required"},
        {"InvalidCIDR", AddPeerRequest{Name: "Test", AllowedIPs: []string{"invalid"}}, 400, "Invalid AllowedIP CIDR"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test each case
        })
    }
}
```

## Test Organization

### File Naming

- Test file: `*_test.go` (same package as code under test)
- Example: `handlers.go` → `handlers_test.go`

### Test Function Naming

**Pattern**: `Test<Function>_<Scenario>_<ExpectedOutcome>`

**Examples**:
- `TestPeerHandlerList_EmptyList_ReturnsEmptyArray`
- `TestPeerHandlerAdd_ValidInput_ReturnsCreated`
- `TestPeerHandlerRemove_NonexistentPeer_ReturnsNoContent`

### Table-Driven Tests

**When**: Testing multiple scenarios for same function

**Pattern**:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {"Scenario1", input1, output1, false},
        {"Scenario2", input2, output2, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Running Tests

### All Tests

```bash
cd backend
go test ./...
```

### Specific Package

```bash
go test ./internal/handlers
```

### Specific Test

```bash
go test -run TestPeerHandlerAdd
```

### With Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
```

### Verbose Output

```bash
go test -v ./...
```

### Integration Tests (Manual)

```bash
INTEGRATION_TESTS=true go test ./...
```

## Mocking Strategy

### Mock Service

**Purpose**: Test handlers without real WireGuard

**Location**: `backend/internal/wireguard/mock.go`

**Usage**:
```go
mockService := wireguard.NewMockService()
handler := handlers.NewPeerHandler(mockService)
// Test handler with mock service
```

**Characteristics**:
- In-memory peer storage (no file I/O)
- No kernel interaction
- Predictable behavior (ideal for testing)

### When NOT to Mock

**Don't mock**:
- Standard library functions (e.g., `net.ParseCIDR`)
- Simple value objects (e.g., `Peer` struct)
- Things you want to actually test

**Mock only**:
- External dependencies (WireGuard kernel)
- I/O operations (file, network)
- Complex collaborators

## Test Maintenance

### When Code Changes

**After modifying handlers/services**:
1. Run existing tests: `go test ./...`
2. If tests fail, update tests to match new behavior
3. If adding new functionality, add new tests (TDD)
4. Verify coverage hasn't dropped: `go test -cover ./...`

### When API Changes

**If endpoint changes**:
1. Update tests to reflect new request/response schemas
2. Add tests for new validation rules
3. Update `backend/API.md` documentation
4. Verify frontend still works (manual integration test)

### Refactoring

**When refactoring code**:
1. Run tests before refactoring: `go test ./...` (should pass)
2. Refactor code (don't change tests)
3. Run tests after refactoring: `go test ./...` (should still pass)
4. If tests fail, fix code (not tests)

**Goal**: Tests should not change during refactoring (they test behavior, not implementation)

## Common Testing Mistakes

### Mistake 1: Testing Implementation, Not Behavior

**Bad**:
```go
func TestHandlerCallsServiceMethod(t *testing.T) {
    // Verifies internal call was made (brittle)
}
```

**Good**:
```go
func TestHandlerReturnsExpectedResponse(t *testing.T) {
    // Verifies public behavior (robust)
}
```

### Mistake 2: Not Testing Error Cases

**Bad**: Only test happy path

**Good**: Test both success and failure paths
```go
func TestAdd_Success(t *testing.T) { /* ... */ }
func TestAdd_InvalidInput(t *testing.T) { /* ... */ }
func TestAdd_ServiceError(t *testing.T) { /* ... */ }
```

### Mistake 3: Fragile Tests

**Symptoms**:
- Tests break on minor refactors
- Tests depend on specific implementation details
- Tests fail intermittently

**Solutions**:
- Test public interfaces, not internals
- Use mocks for external dependencies
- Make tests deterministic (no random data, fixed timestamps)

### Mistake 4: Not Running Tests Before Commit

**Always**:
```bash
go test ./...  # Verify tests pass
go test -race ./...  # Check for race conditions
```

**Before pushing**:
```bash
go test -cover ./...  # Verify coverage
```

## TDD Workflow Example

**Task**: Add endpoint to edit peer name

### Step 1: Write Test (Red)

```go
func TestPeerHandlerUpdate_ValidInput_ReturnsOK(t *testing.T) {
    mockService := wireguard.NewMockService()
    handler := handlers.NewPeerHandler(mockService)
    
    // Add a peer first
    mockService.AddPeer("OldName", "publickey123", []string{"10.0.0.2/32"})
    
    // Update peer name
    body := `{"name":"NewName"}`
    req := httptest.NewRequest("PATCH", "/peers/publickey123", strings.NewReader(body))
    rr := httptest.NewRecorder()
    
    handler.Update(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
    
    // Verify name was updated
    peers, _ := mockService.ListPeers()
    if peers[0].Name != "NewName" {
        t.Errorf("expected name 'NewName', got %s", peers[0].Name)
    }
}
```

**Run test**: `go test ./...` → **FAILS** (handler.Update doesn't exist yet)

### Step 2: Implement (Green)

```go
// handlers.go
func (h *PeerHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    var req struct{ Name string `json:"name"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    if err := h.Service.UpdatePeerName(id, req.Name); err != nil {
        http.Error(w, "Failed to update", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}
```

**Run test**: `go test ./...` → **PASSES**

### Step 3: Refactor

Clean up code, add logging, improve error handling. Tests should still pass.

### Step 4: Repeat

Add more tests (edge cases, error handling), implement features.

## Test Checklist

Before considering backend feature "done":

- [ ] Tests written FIRST (TDD workflow)
- [ ] All happy path scenarios tested
- [ ] All error cases tested (invalid input, service errors)
- [ ] Edge cases tested (empty lists, boundary values)
- [ ] Tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Coverage acceptable: `go test -cover ./...`
- [ ] Tests are readable (good names, clear assertions)
- [ ] Tests are independent (no shared state)
