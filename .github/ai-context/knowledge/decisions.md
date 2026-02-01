# Architectural Decisions

**Last Updated**: 2026-02-01

This document records important architectural decisions, trade-offs, and why they were made.

## D001: Hybrid Persistence Model (Metadata + Kernel State)

**Decision**: Store peer _names_ in JSON file (`peers.json`), but fetch real-time peer data (handshakes, rx/tx) from WireGuard kernel.

**Context**:

- WireGuard kernel keeps peer statistics in memory
- Kernel doesn't persist custom metadata (names, descriptions)
- Frontend needs both: human-readable names AND real-time stats

**Trade-offs**:

| Option              | Pros                                   | Cons                                     |
| ------------------- | -------------------------------------- | ---------------------------------------- |
| **Hybrid (chosen)** | Best UX: real-time data + custom names | Extra file I/O, sync complexity          |
| Database            | Persistent, queryable                  | Adds dependency, overkill for peer names |
| JSON only           | Simple, no deps                        | Stale stats, manual refresh needed       |
| Kernel names        | Single source of truth                 | Kernel API limited, no name field        |

**Rationale**: For a management interface, real-time stats are essential for visibility. Metadata in file is acceptable (low update frequency). This keeps complexity low while providing good UX.

**Impact on code**:

- `service.ListPeers()` queries kernel + reads JSON
- `service.AddPeer()` writes to both WireGuard and JSON
- Storage layer abstracts JSON I/O

## D002: No External Web Framework (net/http Standard Library)

**Decision**: Use Go's standard library `net/http` instead of frameworks like Echo, Gin, or Fiber.

**Context**:

- Initial MVP needed fast delivery
- Small surface area (4 endpoints only)
- Team familiar with standard library
- No complex routing or middleware needs initially

**Trade-offs**:

| Aspect         | Standard Library    | Framework                     |
| -------------- | ------------------- | ----------------------------- |
| Dependencies   | Minimal             | Adds deps (vendoring)         |
| Boilerplate    | More verbose        | Less boilerplate              |
| Routing        | Pattern matching    | Declarative                   |
| Middleware     | Manual composition  | Built-in chains               |
| Performance    | Fast                | Usually slower (extra layers) |
| Learning curve | Familiar to Go devs | Framework-specific            |

**Rationale**: MVP simplicity and reduced dependencies. Standard library is battle-tested and performant. Can migrate to framework later if needed.

**Impact on code**:

- Handlers are manual `http.HandlerFunc` functions
- Middleware is manual composition (applied in `main.go`)
- Routing via `mux.HandleFunc()` patterns

## D003: Mock Service for Graceful Fallback

**Decision**: Implement `mockService` in addition to `realService`. If WireGuard initialization fails, fall back to in-memory mock.

**Context**:

- Development environments may not have WireGuard kernel module
- Testing requires isolation from kernel
- Allows safe development on non-Linux machines (via WSL2, VMs)

**Trade-offs**:

| Approach                   | Pros                                | Cons                              |
| -------------------------- | ----------------------------------- | --------------------------------- |
| **Mock fallback (chosen)** | Dev without WireGuard, easy testing | Complexity (two implementations)  |
| Fail fast                  | Clear error                         | Can't develop without WireGuard   |
| Docker always              | Consistent env                      | Extra setup, not always available |

**Rationale**: Developers can work without WireGuard setup; testing doesn't require root. Warning log makes it obvious when running in mock mode.

**Impact on code**:

- Both `realService` and `mockService` implement `Service` interface
- `main.go` tries real service, falls back to mock with warning

## D004: Twelve-Factor Configuration Pattern

**Decision**: Configuration from `config.json` defaults + environment variable overrides.

**Context**:

- Application needs to run in development, staging, production
- Different environments have different endpoints, ports
- Avoid committing sensitive values (API keys, endpoints)

**Pattern**:

1. Load defaults from `config.json` (version-controlled)
2. Override with environment variables (per-environment)
3. Support `.env` file for local development (via `godotenv`)

**Why not**:

- Flags only: Too many to remember
- Vault only: Adds complexity, not needed for this scale
- All env vars: Hard to see what's configurable

**Impact on code**:

- `config.LoadConfig()` handles the cascade
- Environment variables: `WG_SERVER_PORT`, `WG_INTERFACE_NAME`, etc.
- CI/CD sets these via deployment tools

## D005: Structured Logging via slog (Not Text Logs)

**Decision**: Use Go's `log/slog` for JSON structured logging instead of text logs.

**Context**:

- Need to debug issues in production
- JSON logs are queryable (if using log aggregation)
- Structured logs improve observability

**Trade-offs**:

| Format            | Pros                                    | Cons                              |
| ----------------- | --------------------------------------- | --------------------------------- |
| **JSON (chosen)** | Queryable, structured, machine-readable | Harder to read in terminal        |
| Text              | Human-readable                          | Not queryable, parsing is fragile |

**Workaround**: Development: pipe through JSON pretty-printer if needed

```bash
go run . | jq .  # pretty-print JSON logs
```

**Impact on code**:

- All logs via `slog.Info()`, `slog.Error()`, etc.
- Structured key-value pairs: `slog.Error("failed to add peer", "peer_id", id, "error", err)`

## D006: Frontend UX/Performance Over Test Coverage (Constitution II)

**Decision**: Frontend development prioritizes user experience and performance. No automated testing required.

**Context**:

- SvelteKit with TypeScript provides type safety
- ESLint + Prettier enforce code quality
- Visual components are best tested manually
- Performance budgets (TTI < 3s, bundle < 200KB) are explicit goals

**Trade-offs**:

| Approach              | Pros                                      | Cons                          |
| --------------------- | ----------------------------------------- | ----------------------------- |
| **UX-first (chosen)** | Fast iteration, better UX, lower overhead | Less test coverage            |
| TDD (like backend)    | Regression detection, safety              | Slows iteration, TDD overhead |

**Rationale**: UI testing is expensive and brittle. Manual testing + TypeScript catch most bugs. Performance optimization delivers more value than test suites.

**Impact on code**:

- No jest/vitest configuration
- Focus on component reusability, clean props
- Manual testing in dev mode
- Lighthouse audits for performance

## D007: SvelteKit for Full-Stack Framework

**Decision**: Use SvelteKit instead of separate frontend (React) + backend.

**Context**:

- Team already familiar with Svelte
- Monorepo simplifies local development (frontend + backend in one workspace)
- SvelteKit provides routing, SSR, static export

**Trade-offs**:

| Framework              | Pros                           | Cons                             |
| ---------------------- | ------------------------------ | -------------------------------- |
| **SvelteKit (chosen)** | Cohesive, good DX, lightweight | New tool if unfamiliar           |
| React + Vite           | Mature ecosystem               | Extra scaffolding, two codebases |
| Vue + Nuxt             | Similar to SvelteKit           | Heavier than Svelte              |
| Angular                | Enterprise-grade               | Overly complex for this scale    |

**Rationale**: Svelte's reactive model is intuitive. SvelteKit's file-based routing aligns with how team thinks about pages. No separate backend needed for routing/rendering.

**Impact on code**:

- Routes in `src/routes/`
- Stores for state management
- SvelteKit adapter (currently auto) handles deployment

## D008: Backend Testing is Mandatory (Constitution I)

**Decision**: Go backend follows Test-Driven Development (TDD). All handlers and services must have tests.

**Context**:

- Backend controls critical infrastructure (WireGuard configuration)
- Bugs can cause connectivity loss, security issues
- TDD forces clarity on requirements and error handling

**Trade-offs**:

| Approach         | Pros                                              | Cons                                |
| ---------------- | ------------------------------------------------- | ----------------------------------- |
| **TDD (chosen)** | Fewer bugs, clear contracts, regression detection | Slower initial development          |
| Ad-hoc testing   | Faster development                                | Higher bug risk, harder to refactor |

**Rationale**: Backend is infrastructure code. The cost of bugs (downtime, security) far outweighs test writing overhead. TDD also improves API design clarity.

**Impact on code**:

- All handlers have test cases in `*_test.go`
- Mock services enable unit testing without kernel
- Integration tests validate WireGuard operations

## D009: API-First Contract (Not API Generated from Code)

**Decision**: API contract documented in `backend/API.md` (human-written) before implementation.

**Context**:

- Frontend needs clear, stable endpoint contracts
- Prevents breaking changes
- Enables parallel frontend/backend development

**Why not**:

- Swagger/OpenAPI: Overkill for 4 endpoints, adds tooling
- Code-generated docs: Doesn't guide design, just documents implementation

**Impact on code**:

- `API.md` is source of truth for endpoint behavior
- Tests verify code matches API.md
- Frontend builds against API.md spec (not code-reading)

## D010: SvelteKit adapter-static with SSR Prerendering (2026-02-01) - CORRECTED

**Decision**: Use `@sveltejs/adapter-static` for production builds WITH SSR and prerendering enabled (not SPA mode).

**Context**:

- `adapter-auto` failed to detect a supported production environment
- Build was failing with "Could not detect a supported production environment"
- SvelteKit's SSR capability is a core requirement ("non-negotiable per project requirements")
- Need to produce static site but with server-side rendering during build

**Trade-offs**:

| Adapter + Config           | Pros                              | Cons                                   |
| -------------------------- | --------------------------------- | -------------------------------------- |
| **Static + SSR (chosen)**  | SSR works, static output, reliable| Requires proper prerender/fallback config |
| Static SPA mode            | Simple, pure client-side          | No SSR (violates core requirement)     |
| Auto                       | Automatic detection               | Failed in our environment              |
| Node adapter               | Full SSR support                  | Requires Node.js server (not static)   |

**Rationale**: 

Constitution Principle II prioritizes frontend UX/performance but doesn't preclude SSR. In fact, SvelteKit was chosen specifically for its SSR capability (as user stated: "SSR is non-negotiable (this is why i picked svelte kit)"). With `adapter-static`, we can:
- Pre-render routes at build time using `export const prerender = true;`
- Serve as a static site for fast deployment
- Still leverage SvelteKit's SSR benefits
- Support client-side hydration for interactivity

**Impact on code**:

- `src/routes/+layout.ts`: `export const prerender = true;`
- `svelte.config.js`: `adapter-static` with `fallback: 'index.html'` for SPA-like behavior
- Build output is static files in `build/` directory
- All routes pre-rendered at build time
- Client-side navigation works via SvelteKit's router
| Node                | SSR support                       | Requires Node.js server in production  |
| Vercel/Netlify      | Platform-optimized                | Vendor lock-in                         |

**Rationale**: Frontend is purely client-side rendered (communicates with Go backend API). Static adapter is simpler and appropriate for SPA deployment.

**Impact on code**:

- `svelte.config.js` uses `adapter-static` with `fallback: 'index.html'`
- Created `src/routes/+layout.ts` with `ssr = false`, `prerender = false`
- Build output goes to `build/` directory
- Can be served by any static file server (Nginx, Caddy, etc.)

## D011: API.md as Source of Truth for Frontend Data (2026-02-01)

**Decision**: Frontend must only use properties documented in `backend/API.md`. Missing properties should be added to backlog, not hardcoded.

**Context**:

- Mockup designs showed UI elements (Public Key, Listen Port, Subnet) that required data not in API
- Frontend was using non-existent properties (`$stats.publicKey`, `$stats.listenPort`, `$stats.subnet`)
- Build succeeded but runtime would fail with undefined values

**Trade-offs**:

| Approach                          | Pros                      | Cons                          |
| --------------------------------- | ------------------------- | ----------------------------- |
| **API.md as truth (chosen)**      | Type-safe, no surprises   | May need to adapt UI design   |
| Hardcode placeholder values       | Quick fix                 | Misleading UI, tech debt      |
| Expand API without documentation  | Faster UI development     | Contract drift, breaking changes |

**Rationale**: Constitution III mandates API contract stability. Adding undocumented properties creates implicit contracts that can break. Better to document gaps in BACKLOG.md.

**Impact on code**:

- `+page.svelte` dashboard cards show available data (interface name, peer count, online count)
- `BACKLOG.md` documents missing API properties for future backend work
- `src/lib/types/stats.ts` matches API.md exactly

## Decision Methodology

When making future architectural decisions:

1. **Write the decision in this file** (copy template below)
2. **Include trade-offs** (pros/cons of alternatives)
3. **Document the rationale** (why chosen over alternatives)
4. **Note impacts** (what code/processes must change)
5. **Review with team** (seek consensus)

### Template

```
## DNNN: [Decision Title]

**Decision**: [One sentence decision]

**Context**: [Why this decision now?]

**Trade-offs**: [Pros/cons of options as table]

**Rationale**: [Why this wins?]

**Impact on code**: [What changes in codebase?]
```
