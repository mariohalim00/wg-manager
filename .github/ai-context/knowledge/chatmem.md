# Chat Memory & Session Context

**Last Updated**: 2026-01-31

This file helps AI agents "pick up where they left off" by tracking active development areas, common pitfalls, and important invariants.

## Active Development Areas

### Current Branch: `feat/frontend`

**Status**: Frontend implementation in progress

**What's being built**:
- SvelteKit frontend (dashboard, peer management, stats)
- Integration with Go backend API
- Component library (PeerTable, PeerModal, StatusBadge, QRCodeDisplay, Notification)
- State management via Svelte stores

**What's complete**:
- ✅ Project scaffolding (SvelteKit + TypeScript + TailwindCSS + DaisyUI)
- ✅ Routing structure (`/`, `/peers`, `/stats`)
- ✅ Core components (PeerTable, PeerModal, StatusBadge)
- ✅ Svelte stores (peers, stats)
- ✅ Type definitions (types.ts)
- ✅ Backend API fully functional

**What's in progress**:
- 🔄 API integration (fetch calls in stores)
- 🔄 Error handling and loading states
- 🔄 QR code generation for peer config
- 🔄 Performance optimization (bundle size, lazy loading)

**What's next**:
- Integration testing (frontend + backend together)
- Performance audit (Lighthouse)
- Deployment setup (build + serve)

### Recent Sessions

#### Session: Constitution Establishment (2026-01-31)

**What was done**:
- Created `.specify/memory/constitution.md` (v1.0.0)
- Established 6 core principles (Backend TDD, Frontend UX, API Stability, Config, Performance, Observability)
- Updated spec/plan/tasks templates to align with constitution
- Codified testing philosophy: Backend = mandatory TDD, Frontend = no tests required

**Key takeaway**: Future features must verify alignment with constitution principles.

#### Session: AI Context System Creation (2026-01-31)

**What was done**:
- Created `.github/ai-context/` structure
- Agent rules (AGENTS.md, CLAUDE.md, GEMINI.md)
- Skills files (go.skills.md, svelte.skills.md)
- Knowledge base (architecture.md, decisions.md, domains.md, memory.md, chatmem.md)

**Key takeaway**: AI agents now have comprehensive project context.

## Critical Invariants (MUST NOT VIOLATE)

### Backend Invariants

1. **PublicKey is immutable identifier**
   - Never allow changing PublicKey of existing peer
   - Always use PublicKey as unique ID in API URLs
   - Format: Base64 string, ~44 chars

2. **AllowedIPs must be valid CIDR**
   - Always validate with `net.ParseCIDR()` before accepting
   - Reject invalid input with 400 Bad Request
   - Example valid: `10.0.0.2/32`, `10.0.0.0/24`

3. **TDD is mandatory for backend**
   - Write tests FIRST, then implementation
   - All handlers must have unit tests
   - All service methods must have contract tests
   - Never commit untested backend code

4. **Logging must be structured JSON**
   - Use `slog.Info()`, `slog.Error()`, etc.
   - Always include context: `slog.Error("failed to add peer", "peer_id", id, "error", err)`
   - Never log sensitive data (private keys, passwords)

5. **Configuration follows Twelve-Factor**
   - Defaults in `config.json`
   - Overrides via environment variables
   - Never hardcode secrets or endpoints in code

6. **API contract stability**
   - Breaking changes require MAJOR version bump
   - Keep `backend/API.md` updated with all endpoint changes
   - Backward compatibility is non-negotiable

### Frontend Invariants

1. **Type safety via TypeScript**
   - Always use `lang="ts"` in Svelte script blocks
   - All props, function params, returns must be typed
   - No `any` types without justification

2. **Performance budgets are strict**
   - TTI < 3 seconds on 3G
   - FCP < 1.5 seconds
   - Bundle size < 200KB gzipped
   - Lighthouse ≥ 90
   - Never commit performance regressions without justification

3. **No automated testing required**
   - Manual testing + type safety is sufficient
   - Focus on UX and performance, not test coverage
   - This is intentional per Constitution Principle II

4. **State via Svelte stores**
   - Use stores for shared state (peers, stats)
   - No prop drilling for global state
   - Components subscribe with `$store` syntax

5. **Styling via TailwindCSS + DaisyUI**
   - Never write custom CSS except for scoped component styles
   - Use Tailwind utilities for layout, spacing, colors
   - Use DaisyUI components for buttons, modals, tables

## Common Pitfalls & Solutions

### Pitfall 1: Type Mismatch Between Frontend and Backend

**Problem**: Backend uses Go structs (snake_case JSON tags), frontend uses TypeScript (camelCase)

**Example**:
- Backend: `"allowed_ips"`
- Frontend: `allowedIps`

**Solution**: 
- Keep `src/lib/types.ts` in sync with backend response schemas
- Document in `backend/API.md`
- Use camelCase in TypeScript, ensure backend JSON tags match

**How to verify**: Test API integration, check browser DevTools Network tab

### Pitfall 2: WireGuard Permissions

**Problem**: Backend fails to start with "permission denied" on WireGuard operations

**Symptom**: `Failed to initialize real service, falling back to mock`

**Solution**:
- Run backend with `sudo` or grant `CAP_NET_ADMIN` capability
- Or accept fallback to mock service for development
- Mock service is intentional for non-root development

**How to verify**: Check logs for "using mock" warning

### Pitfall 3: CORS Errors in Frontend

**Problem**: Frontend can't make API requests due to CORS policy

**Symptom**: Browser console shows CORS error

**Solution**:
- Ensure backend CORS middleware is enabled
- Set `CORS_ALLOWED_ORIGINS` to frontend origin (e.g., `http://localhost:5173`)
- In production, use reverse proxy to avoid CORS (serve frontend + backend from same origin)

**How to verify**: Check browser DevTools Network tab for CORS headers

### Pitfall 4: Bundle Size Creep

**Problem**: Frontend bundle size exceeds 200KB budget

**Symptom**: Lighthouse performance score drops, slow load times

**Solution**:
- Check `npm run build` output for bundle size
- Identify large imports (use rollup-plugin-visualizer)
- Lazy-load heavy components
- Remove unused dependencies

**How to verify**: Run `npm run build`, check `.svelte-kit/output/client` bundle sizes

### Pitfall 5: Breaking API Changes Without Version Bump

**Problem**: Backend API changes break frontend integration

**Symptom**: Frontend shows errors, type mismatches, null values

**Solution**:
- NEVER change existing API response schemas without MAJOR version bump
- Add new optional fields instead of changing existing ones
- Update `backend/API.md` with all changes
- Verify frontend still works after backend changes

**How to verify**: Run frontend against backend, check browser console for errors

## Recent Discoveries

### Discovery 1: Mock Service is Development Lifesaver

**What we learned**: Mock service allows development without WireGuard kernel

**Evidence**: Team members on macOS/Windows can develop using mock

**Impact**: Development velocity increased, no VM/Docker required for frontend work

**Actionable**: Always use mock service in development, only test with real WireGuard before production

### Discovery 2: DaisyUI Simplifies UI Development

**What we learned**: DaisyUI components are consistent, accessible, responsive

**Evidence**: PeerTable, PeerModal, StatusBadge built in <1 hour with DaisyUI

**Impact**: UI development is fast, visually consistent

**Actionable**: Always check DaisyUI docs first before writing custom components

### Discovery 3: Svelte Stores Scale Well

**What we learned**: Stores handle API integration cleanly

**Evidence**: `peers.ts` encapsulates all peer operations (load, add, remove)

**Impact**: Components are simple, store handles complexity

**Actionable**: Put API calls in stores, not in components

## Frequently Asked Questions (Agents)

### Q: Should I write tests for this frontend component?

**A**: No. Per Constitution Principle II, frontend tests are not required. Focus on TypeScript type safety and manual testing.

### Q: Should I write tests for this backend handler?

**A**: Yes. Per Constitution Principle I, backend TDD is mandatory. Write tests FIRST, then implementation.

### Q: How do I run the full stack locally?

**A**:
```bash
# Terminal 1: Backend (from project root)
cd backend
go run cmd/server/main.go  # or sudo if you want real WireGuard

# Terminal 2: Frontend (from project root)
npm run dev
```
Frontend: http://localhost:5173 (or 5174 if 5173 taken)
Backend: http://localhost:8080

### Q: Where do I document new API endpoints?

**A**: Update `backend/API.md` with endpoint details (URL, method, request/response schemas, status codes).

### Q: How do I add a new Svelte component?

**A**:
1. Create `src/lib/components/[Name].svelte`
2. Define props with TypeScript types
3. Use TailwindCSS + DaisyUI for styling
4. Import and use in pages

### Q: How do I add a new backend handler?

**A**:
1. Add method to `PeerHandler` in `handlers/handlers.go`
2. Write tests FIRST in `main_test.go`
3. Implement handler method
4. Register route in `main.go`: `mux.HandleFunc("METHOD /path", handler.Method)`
5. Update `backend/API.md`

### Q: What's the difference between `realService` and `mockService`?

**A**: 
- `realService`: Actually controls WireGuard kernel (requires Linux + root)
- `mockService`: In-memory implementation (for development/testing)
- Fallback logic in `main.go` tries real, falls back to mock on error

### Q: How do I check performance?

**A**:
```bash
npm run build  # Check bundle sizes
npm run preview  # Serve production build
# Run Lighthouse audit in browser DevTools
```

## Things to Remember When Proposing Changes

1. **Backend changes**: TDD workflow (write tests first, verify they fail, implement, verify they pass)
2. **Frontend changes**: Type safety + manual testing (no automated tests needed)
3. **API changes**: Update `backend/API.md`, verify frontend still works
4. **Performance changes**: Measure before/after (Lighthouse, bundle size)
5. **Configuration changes**: Document in `backend/API.md` or README
6. **Breaking changes**: Justify explicitly, requires major version bump

## Current Blockers & Dependencies

### Blocker 1: None currently

**Status**: Development proceeding smoothly

### Dependency 1: WireGuard Kernel Module (Production Only)

**What**: Backend requires WireGuard kernel module for real operations

**Impact**: Production deployment must be on Linux with WireGuard installed

**Workaround**: Mock service works for development

## Session Recovery Checklist

When an AI agent starts a new session, verify:

- [ ] Understand current branch: `feat/frontend` (frontend development)
- [ ] Read constitution: `.specify/memory/constitution.md`
- [ ] Check architecture: `.github/ai-context/knowledge/architecture.md`
- [ ] Review decisions: `.github/ai-context/knowledge/decisions.md`
- [ ] Know domains: `.github/ai-context/knowledge/domains.md`
- [ ] Read skills: `.github/ai-context/skills/go.skills.md` and `svelte.skills.md`
- [ ] Check active work: This file (chatmem.md) → "Active Development Areas"
- [ ] Verify no blockers: This file (chatmem.md) → "Current Blockers"

## Notes for Future Sessions

**Next time an agent works on this project**:

1. If working on **backend**: Remember TDD is mandatory, write tests first
2. If working on **frontend**: Focus on UX/performance, no tests needed
3. If touching **API**: Update `backend/API.md`, verify frontend compatibility
4. If adding **features**: Use SpecKit workflow (spec → plan → tasks → implement)
5. If unsure: Read constitution, check decisions.md for precedent

**Keep this file updated**: When major work is done, update "Active Development Areas" and "Recent Sessions".
