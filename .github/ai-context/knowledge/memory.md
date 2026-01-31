# Project History & Memory

**Last Updated**: 2026-01-31

This document tracks the evolution of WireGuard Manager, major milestones, and key decisions that shaped the current codebase.

## Project Timeline

### Phase 1: Initial Concept (Estimated: Q4 2025)

**Goal**: Create a web-based management interface for WireGuard VPN peers.

**Decisions made**:
- Target Linux servers with WireGuard kernel module
- REST API + SPA architecture (backend + frontend)
- Go for backend (standard library, performance)
- Svelte for frontend (lightweight, reactive)

### Phase 2: Backend MVP (Estimated: Late 2025)

**Implemented**:
- Go HTTP server with 4 core endpoints (List, Add, Remove, Stats)
- WireGuard control via `wgctrl` library
- File-based metadata storage (`peers.json`)
- Structured JSON logging via `slog`
- CORS middleware for frontend integration
- Graceful shutdown handling

**Key files created**:
- `backend/cmd/server/main.go` — Server entry point
- `backend/internal/handlers/handlers.go` — Request handlers
- `backend/internal/wireguard/service.go` — WireGuard abstraction
- `backend/internal/config/config.go` — Configuration management

**Challenges encountered**:
- Root permissions required for WireGuard operations
- Fallback strategy: Implemented mock service for development
- Metadata vs. real-time state: Solved with hybrid model (JSON + kernel)

**Testing approach**:
- TDD adopted for all backend code
- Unit tests for handlers (using mock service)
- Integration tests with real WireGuard (when available)

### Phase 3: Frontend Development (Current: Early 2026)

**Branch**: `feat/frontend` (in progress)

**Implemented**:
- SvelteKit project scaffolding
- File-based routing (`/`, `/peers`, `/stats`)
- Component library (PeerTable, PeerModal, StatusBadge, QRCodeDisplay)
- Svelte stores for state management (peers, stats)
- TailwindCSS + DaisyUI integration
- API integration layer

**Key files created**:
- `src/routes/+page.svelte` — Dashboard
- `src/routes/peers/+page.svelte` — Peer management
- `src/lib/components/` — Reusable components
- `src/lib/stores/peers.ts` — Peer state management
- `src/lib/types.ts` — TypeScript interfaces

**Frontend philosophy** (Constitution II):
- UX and performance prioritized over test coverage
- No automated tests for frontend (manual testing + type safety)
- Performance budgets: TTI <3s, bundle <200KB, Lighthouse ≥90

**Current status**:
- Core functionality implemented
- Visual design based on DaisyUI components
- API integration complete
- Ready for integration testing and deployment

### Phase 4: Constitution & Spec-Driven Development (2026-01-31)

**Major addition**: Established project constitution and SpecKit integration

**Constitution principles** (v1.0.0):
- **Principle I**: Backend TDD mandatory (Red-Green-Refactor)
- **Principle II**: Frontend UX-first (no tests required)
- **Principle III**: API contract stability (breaking changes = major version)
- **Principle IV**: Twelve-Factor configuration
- **Principle V**: Performance budgets (backend <100ms, frontend TTI <3s)
- **Principle VI**: Structured logging and observability

**SpecKit setup**:
- `.specify/memory/constitution.md` — Core principles
- `.specify/templates/` — Spec, plan, tasks templates
- Integration with AI-driven development workflow

**Impact**:
- Codifies existing practices (TDD, performance focus)
- Clear separation: backend (test-heavy) vs. frontend (UX-heavy)
- Future features will use spec-driven workflow

### Phase 5: AI Context System (2026-01-31)

**Major addition**: Unified AI context for all AI agents

**Created**:
- `.github/ai-context/` — Centralized AI knowledge
- Agent-specific rules (AGENTS.md, CLAUDE.md, GEMINI.md)
- Skills files (go.skills.md, svelte.skills.md)
- Knowledge base (architecture.md, decisions.md, domains.md, memory.md, chatmem.md)
- Global copilot-instructions.md

**Rationale**: 
- Help AI agents understand project structure and conventions
- Reduce redundant documentation
- Consolidate context for multiple AI tools (Copilot, Claude, Gemini, Cursor)

## Major Refactors & Migrations

### Refactor 1: Handler Extraction (Estimated: Late 2025)

**What changed**: Moved handler logic from `main.go` into dedicated `handlers/` package

**Before**:
```go
// main.go (bloated)
mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
    // 50+ lines of handler logic inline
})
```

**After**:
```go
// main.go (clean)
peerHandler := handlers.NewPeerHandler(app.WireGuard)
mux.HandleFunc("GET /peers", peerHandler.List)

// handlers/handlers.go
func (h *PeerHandler) List(w http.ResponseWriter, r *http.Request) {
    // Handler logic
}
```

**Why**: Separation of concerns, testability, cleaner main.go

**Impact**: All tests moved to handler package; main.go simplified to 100 lines

### Refactor 2: Mock Service Introduction (Estimated: Late 2025)

**What changed**: Added `mockService` implementation alongside `realService`

**Rationale**: Enable development without WireGuard kernel, support testing

**Pattern**:
```go
// wireguard/service.go
type Service interface {
    ListPeers() ([]Peer, error)
    // ...
}

// wireguard/service.go (real)
type realService struct { /* ... */ }

// wireguard/mock.go (new)
type mockService struct { /* ... */ }
```

**Impact**: 
- Tests no longer require root permissions
- Development works on any machine
- `main.go` has fallback logic

## Abandoned Approaches

### Abandoned 1: Database for Metadata

**Tried**: SQLite for peer metadata storage

**Why abandoned**: 
- Overkill for small peer counts (<100 typically)
- Adds dependency and complexity
- File-based JSON is sufficient (low update frequency)
- Migrations would be overhead

**Current approach**: JSON file (`peers.json`)

**Lesson**: Start simple, add complexity only when needed

### Abandoned 2: Frontend Testing with Vitest

**Tried**: Setting up Vitest for Svelte component tests

**Why abandoned**:
- High setup cost for minimal value
- UI tests brittle (false positives on style changes)
- Manual testing caught same bugs faster
- TypeScript + ESLint provide sufficient safety

**Current approach**: Manual testing + performance audits

**Codified in**: Constitution Principle II (Frontend UX First)

### Abandoned 3: gRPC API

**Considered**: Using gRPC instead of REST for backend API

**Why abandoned**:
- Overkill for simple CRUD operations
- Browser support requires grpc-web (extra complexity)
- REST is simpler for small API surface
- No strong-typing benefit over TypeScript + JSON

**Current approach**: REST API with JSON (documented in API.md)

## Long-Term Constraints

### Constraint 1: Linux-Only Backend

**Reality**: WireGuard kernel module is Linux-only (natively)

**Implications**:
- Backend must run on Linux server
- Development on macOS/Windows requires VM/Docker/WSL
- Mock service mitigates this for development

**No plans to change**: WireGuard's Windows/macOS clients work differently; not worth supporting

### Constraint 2: Single Interface Management

**Reality**: System manages one WireGuard interface (e.g., `wg0`)

**Implications**:
- Multi-interface setups not supported
- Interface name hardcoded in config (can be changed via env var)

**Future consideration**: Could add multi-interface support if needed (low priority)

### Constraint 3: No Authentication

**Reality**: API has no built-in auth (open to anyone with network access)

**Implications**:
- Rely on network-level security (VPN, firewall, reverse proxy auth)
- Suitable for private networks only
- NOT suitable for public internet exposure

**Future consideration**: Add basic auth or OAuth if needed

### Constraint 4: File-Based Metadata

**Reality**: Peer metadata stored in JSON file

**Implications**:
- Not suitable for high-frequency updates (file I/O on every write)
- No concurrent write safety (single process only)
- Works fine for typical use case (<100 peers, infrequent adds/removes)

**Future consideration**: Migrate to database if peer count grows significantly

## Key Learnings

### Learning 1: TDD is Worth It (Backend)

**Observation**: TDD initially felt slow, but paid off during refactors

**Evidence**: Refactoring handlers was safe because tests caught regressions

**Lesson**: For infrastructure code (WireGuard operations), testing prevents costly bugs

**Codified in**: Constitution Principle I (Backend Testing Discipline)

### Learning 2: Frontend Tests Have Diminishing Returns

**Observation**: Attempted UI testing was time-consuming, brittle, low value

**Evidence**: Manual testing caught more bugs faster; tests broke on style changes

**Lesson**: For UI-heavy code, type safety + manual testing > automated tests

**Codified in**: Constitution Principle II (Frontend UX First)

### Learning 3: Simple Configuration Wins

**Observation**: Tried complex config management, ended up with simpler approach

**Evidence**: JSON defaults + env var overrides cover all use cases

**Lesson**: Twelve-Factor pattern is sufficient for most projects

**Codified in**: Constitution Principle IV (Configuration & Environment)

### Learning 4: Documentation Prevents Rework

**Observation**: Undocumented API decisions led to frontend/backend mismatches

**Evidence**: Early iterations had type mismatches (snake_case vs. camelCase)

**Lesson**: API.md as source of truth prevents integration bugs

**Codified in**: Decision D009 (API-First Contract)

## Evolution of Code Style

### Backend Go Style

**Early**: Inline handlers, minimal error handling
**Now**: 
- Handler → Service separation
- Comprehensive error logging with context
- Struct-based handlers (PeerHandler with dependencies)
- Interface-driven design (Service interface + implementations)

### Frontend Svelte Style

**Early**: Monolithic components, props drilling
**Now**:
- Small, focused components (PeerTable, PeerModal, StatusBadge)
- Svelte stores for shared state (no prop drilling)
- TypeScript interfaces for all data (types.ts)
- DaisyUI components for consistent styling

## Known Technical Debt

### Debt 1: No Peer Edit Functionality

**What's missing**: Can add/remove peers, but not edit existing peers

**Why**: MVP scope limitation

**Impact**: Users must delete + re-add to change AllowedIPs or name

**Priority**: Medium (nice-to-have, workaround exists)

### Debt 2: No Bulk Operations

**What's missing**: Can't add/remove multiple peers at once

**Why**: Not needed for typical use case (few peers per session)

**Impact**: Manual work for large imports

**Priority**: Low (edge case)

### Debt 3: No Real-Time Updates

**What's missing**: Frontend polls manually, no WebSocket/SSE

**Why**: Polling sufficient for current use case

**Impact**: Stats are slightly stale (until user refreshes)

**Priority**: Low (polling works, real-time is nice-to-have)

### Debt 4: No Backup/Restore

**What's missing**: Can't export/import full peer list

**Why**: Not in MVP scope

**Impact**: Manual backup of peers.json required

**Priority**: Medium (useful for disaster recovery)

## Migration Path from Legacy

**N/A**: This is a greenfield project, no legacy migration needed.

## Future Roadmap Considerations

**Not committed, but considered**:

1. **Authentication layer** — Basic auth or reverse proxy integration
2. **Database backend** — Replace peers.json for larger deployments
3. **Multi-interface support** — Manage multiple WireGuard interfaces
4. **Real-time updates** — WebSocket for live stats
5. **Peer groups** — Organize peers into categories
6. **Bulk operations** — Add/remove multiple peers at once
7. **Backup/restore** — Export/import peer configurations
8. **Metrics export** — Prometheus endpoint for monitoring
9. **Mobile app** — Native iOS/Android app (long-term)

**Priority**: Validate current system in production before expanding scope.
