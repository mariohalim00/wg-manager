# GitHub Copilot Instructions for WireGuard Manager

**Project**: WireGuard Manager — Full-stack VPN peer management system  
**Stack**: Go (backend) + SvelteKit (frontend) + WireGuard  
**Last Updated**: 2026-01-31

## Quick Start for Copilot

### Project Philosophy

This is a **dual-nature codebase** with different quality approaches for backend vs. frontend:

- **Backend (Go)**: Test-Driven Development, high reliability, infrastructure code
- **Frontend (Svelte)**: UX-first, performance-focused, manual testing

**Golden Rule**: Respect the [Constitution](./../.specify/memory/constitution.md) — it defines non-negotiable principles.

### Core Workflow

When generating code:

1. **Analyze context** — Read relevant files first (use grep, semantic search)
2. **Follow patterns** — Match existing code style (see Skills section below)
3. **Respect principles** — Backend = TDD mandatory, Frontend = no tests required
4. **Verify integration** — Changes must work across frontend + backend

## Agent Behavior & Rules

**See**: [ai-context/agents/AGENTS.md](ai-context/agents/AGENTS.md) for universal rules

**Summary**:
- **Propose before acting** on significant changes
- **Type safety always** (TypeScript frontend, Go backend)
- **Cross-stack awareness** (API changes affect both sides)
- **Performance budgets** non-negotiable (TTI <3s, bundle <200KB, Lighthouse ≥90)
- **Never hardcode secrets** (use env vars, see [config.go](../../backend/internal/config/config.go))

### Safety Rules (Never Violate)

1. **Backend**: Never commit untested Go code (handlers, services, middleware)
2. **Frontend**: Never exceed performance budgets (bundle <200KB, TTI <3s)
3. **API**: Never break contracts without MAJOR version bump
4. **Types**: Never use `any` in TypeScript or skip type annotations
5. **Security**: Never log sensitive data (private keys, passwords)
6. **Constitution**: Never violate principles without explicit justification

## Skills Reference

### Backend (Go)

**See**: [ai-context/skills/go.skills.md](ai-context/skills/go.skills.md)

**Quick facts**:
- Module: `wg-manager/backend`, Go 1.25.6
- Web: Standard library `net/http` (no framework)
- Logging: `log/slog` (JSON structured)
- Testing: TDD mandatory (write tests first)
- WireGuard: `golang.zx2c4.com/wireguard/wgctrl` library

**Handler pattern**:
```go
func (h *PeerHandler) MethodName(w http.ResponseWriter, r *http.Request) {
    // 1. Validate input
    // 2. Call service method
    // 3. Serialize response
    // 4. Handle errors with slog
}
```

**Service pattern**: Interface-driven (`Service` interface, `realService` + `mockService` implementations)

**Configuration**: Twelve-Factor (JSON defaults + env var overrides)

### Frontend (SvelteKit)

**See**: [ai-context/skills/svelte.skills.md](ai-context/skills/svelte.skills.md)

**Quick facts**:
- SvelteKit 2.x, Svelte 5.x, TypeScript 5.x
- Routing: File-based (`src/routes/`)
- State: Svelte stores (`$store` syntax)
- Styling: TailwindCSS + DaisyUI (utility-first, no custom CSS)
- Build: Vite 7.x

**Component pattern**:
```svelte
<script lang="ts">
    import type { Peer } from '$lib/types';
    export let peer: Peer;
    export let onDelete: (id: string) => void = () => {};
</script>

<div class="card bg-base-100">
    <!-- Use TailwindCSS + DaisyUI -->
</div>
```

**No testing required**: Manual testing + TypeScript type safety sufficient (per Constitution II)

## Knowledge Base

### System Architecture

**See**: [ai-context/knowledge/architecture.md](ai-context/knowledge/architecture.md)

**Summary**:
- Frontend (SvelteKit) → Backend API (Go) → WireGuard Kernel
- Hybrid persistence: Metadata (peers.json) + Real-time state (kernel)
- REST API: 4 endpoints (List, Add, Remove, Stats)
- Data flow: Frontend store → HTTP → Handler → Service → WireGuard/Storage

**Integration points**:
- API contract: `backend/API.md` (source of truth)
- Types: Keep `src/lib/types.ts` in sync with Go structs
- CORS: Backend middleware allows frontend origin

### Domain Language

**See**: [ai-context/knowledge/domains.md](ai-context/knowledge/domains.md)

**Critical terms**:
- **Peer**: WireGuard client (identified by PublicKey, not name)
- **PublicKey**: Immutable identifier (base64, ~44 chars)
- **AllowedIPs**: CIDR notation (e.g., `10.0.0.2/32`)
- **Handshake**: WireGuard protocol exchange (LastHandshake = activity indicator)
- **Status**: `online` (recent handshake) or `offline` (stale/never)
- **Real vs. Mock**: Real = actual WireGuard, Mock = in-memory (dev/test)

### Architectural Decisions

**See**: [ai-context/knowledge/decisions.md](ai-context/knowledge/decisions.md)

**Key decisions**:
- D002: No web framework (standard library `net/http`)
- D003: Mock service for graceful fallback
- D006: Frontend UX/performance over test coverage
- D008: Backend TDD is mandatory

### Project History

**See**: [ai-context/knowledge/memory.md](ai-context/knowledge/memory.md)

**Current phase**: Frontend development (branch: `feat/frontend`)

**Recent work**:
- Constitution established (v1.0.0)
- AI context system created
- Frontend components implemented (PeerTable, PeerModal, StatusBadge)

### Session Context

**See**: [ai-context/knowledge/chatmem.md](ai-context/knowledge/chatmem.md)

**Active work**:
- Frontend integration (API calls in stores)
- Performance optimization (bundle size, Lighthouse)
- Deployment preparation

**Critical invariants**:
- PublicKey is immutable
- AllowedIPs must be valid CIDR
- TDD mandatory for backend
- Performance budgets strict for frontend

## Common Tasks & How to Execute

### Task: Add New API Endpoint

**Process**:
1. **Document first**: Add to `backend/API.md` (endpoint, request/response, status codes)
2. **Write tests**: Create test in `backend/cmd/server/main_test.go` (TDD)
3. **Implement handler**: Add method to `PeerHandler` in `backend/internal/handlers/handlers.go`
4. **Register route**: Update `main.go`: `mux.HandleFunc("METHOD /path", handler.Method)`
5. **Verify tests pass**: Run `go test ./...`
6. **Update frontend types**: If new data, update `src/lib/types.ts`

### Task: Add New Svelte Component

**Process**:
1. **Create file**: `src/lib/components/[Name].svelte`
2. **Define props**: `export let prop: Type;` (TypeScript types required)
3. **Use TailwindCSS + DaisyUI**: Style with utility classes
4. **Import in page**: `import Component from '$lib/components/Component.svelte';`
5. **Test manually**: Run `npm run dev`, verify UI

### Task: Modify Backend Handler

**Process**:
1. **Read existing tests**: Understand expected behavior
2. **Add new test**: Write failing test for new behavior (TDD)
3. **Modify handler**: Update logic in `handlers/handlers.go`
4. **Verify tests pass**: Run `go test ./...`
5. **Update API.md**: Document changes

### Task: Integrate New API in Frontend

**Process**:
1. **Update types**: Add/modify interface in `src/lib/types.ts`
2. **Add store function**: Create API call in `src/lib/stores/[name].ts`
3. **Use in component**: Subscribe with `$store` or call store function
4. **Handle loading/error states**: Show spinners, error messages
5. **Test manually**: Verify in browser DevTools Network tab

### Task: Optimize Frontend Performance

**Process**:
1. **Baseline**: Run `npm run build`, note bundle sizes
2. **Identify heavy imports**: Check build output, use visualizer
3. **Optimize**: Lazy-load routes, remove unused deps, tree-shake
4. **Rebuild**: Run `npm run build`, verify size decreased
5. **Audit**: Run Lighthouse, verify score ≥90

## Code Patterns to Follow

### Backend: Error Handling

```go
if err != nil {
    slog.Error("operation failed", 
        "error", err, 
        "context_key", contextValue)
    http.Error(w, "User-facing message", http.StatusInternalServerError)
    return
}
```

### Backend: Adding Route

```go
// main.go
mux.HandleFunc("GET /path", handler.Method)
mux.HandleFunc("POST /path", handler.Method)
mux.HandleFunc("DELETE /path/{id}", handler.Method)
```

### Frontend: Svelte Store

```typescript
// stores/example.ts
import { writable } from 'svelte/store';
export const data = writable<Type[]>([]);

export async function loadData() {
    const response = await fetch('/api/endpoint');
    const json = await response.json();
    data.set(json);
}
```

### Frontend: Component with Props

```svelte
<script lang="ts">
    import type { Peer } from '$lib/types';
    export let peer: Peer;
    export let onClick: () => void = () => {};
</script>

<button class="btn btn-primary" on:click={onClick}>
    {peer.name}
</button>
```

## When to Ask for Clarification

- Architecture decisions not in [decisions.md](ai-context/knowledge/decisions.md)
- Trade-offs between performance, complexity, maintainability
- Design patterns that differ from established conventions
- Potential impact on cross-stack integration
- Performance regression justification

## Integration with SpecKit

This project uses **Spec-Driven Development** via SpecKit.

**Workflow**:
1. Feature specification: `/specs/[###-feature]/spec.md`
2. Implementation plan: `/speckit.plan` command generates `plan.md`
3. Task decomposition: `/speckit.tasks` command generates `tasks.md`
4. Implementation: Follow tasks, verify against spec

**Constitution check**: All specs/plans must verify alignment with [constitution.md](./../.specify/memory/constitution.md)

**Templates**: Located in `.specify/templates/` (already updated to reflect constitution)

## Debugging & Troubleshooting

### Backend Won't Start

**Symptom**: "Failed to initialize WireGuard service"

**Solution**: 
- Check if WireGuard kernel module loaded: `lsmod | grep wireguard`
- Run with sudo: `sudo go run cmd/server/main.go`
- Or accept mock service (warning logged)

### Frontend CORS Errors

**Symptom**: Browser console shows CORS policy error

**Solution**:
- Verify backend CORS middleware enabled
- Set `CORS_ALLOWED_ORIGINS=http://localhost:5173`
- Or use reverse proxy (serve frontend + backend from same origin)

### Type Mismatch (Frontend/Backend)

**Symptom**: TypeScript errors, unexpected null values

**Solution**:
- Verify `src/lib/types.ts` matches backend response schemas
- Check `backend/API.md` for correct response format
- Test API in browser DevTools Network tab

### Bundle Size Exceeded

**Symptom**: Build output shows >200KB bundle

**Solution**:
- Check for large imports (use `rollup-plugin-visualizer`)
- Lazy-load heavy components
- Remove unused dependencies
- Tree-shake with proper imports

## Version Control & Branching

**Current branch**: `feat/frontend` (frontend development)  
**Main branch**: `main` (stable releases)

**Branching strategy**: Feature branches (e.g., `feat/auth`, `fix/bug-name`)

**Commit messages**: Conventional commits (e.g., `feat: add peer modal`, `fix: CORS headers`)

## Additional Resources

- **Constitution**: [.specify/memory/constitution.md](../../.specify/memory/constitution.md) — Core principles (v1.0.0)
- **API Docs**: [backend/API.md](../../backend/API.md) — Endpoint specifications
- **Go Docs**: [ai-context/skills/go.skills.md](ai-context/skills/go.skills.md) — Backend patterns
- **Svelte Docs**: [ai-context/skills/svelte.skills.md](ai-context/skills/svelte.skills.md) — Frontend patterns
- **Architecture**: [ai-context/knowledge/architecture.md](ai-context/knowledge/architecture.md) — System design
- **Decisions**: [ai-context/knowledge/decisions.md](ai-context/knowledge/decisions.md) — Why things are this way
- **Domain Terms**: [ai-context/knowledge/domains.md](ai-context/knowledge/domains.md) — WireGuard concepts

---

**Remember**: Backend = TDD mandatory | Frontend = UX-first, no tests | API = documented contracts | Performance = strict budgets
