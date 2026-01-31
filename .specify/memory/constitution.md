<!--
═══════════════════════════════════════════════════════════════════════════════
SYNC IMPACT REPORT
═══════════════════════════════════════════════════════════════════════════════
Version Change: [NOT SET] → 1.0.0

Modified Principles:
  - All principles defined for first time

Added Sections:
  - I. Backend Testing Discipline (Go)
  - II. Frontend User Experience First (Svelte)
  - III. API Contract Stability
  - IV. Configuration & Environment
  - V. Performance Budgets
  - VI. Observability & Structured Logging
  - Technology Stack
  - Development Workflow

Removed Sections:
  - None

Templates Status:
  ✅ plan-template.md - Updated constitution check gates
  ✅ spec-template.md - Updated to reflect frontend UX priority and no frontend tests
  ✅ tasks-template.md - Updated to remove frontend test requirements

Follow-up TODOs:
  - None

Change Rationale:
  Initial constitution establishment defining core principles for WireGuard Manager,
  emphasizing backend testing discipline, frontend UX/performance, and API stability.
═══════════════════════════════════════════════════════════════════════════════
-->

# WireGuard Manager Constitution

## Core Principles

### I. Backend Testing Discipline (Go)

**Backend code MUST follow Test-Driven Development (TDD)**:
- Write tests FIRST → User approval → Tests fail → Then implement
- Red-Green-Refactor cycle strictly enforced
- All API handlers MUST have unit tests
- Service layer logic MUST have contract tests
- Integration tests required for WireGuard interface interactions

**Rationale**: Backend manages critical infrastructure (WireGuard peers, network configuration). Bugs can cause connectivity loss or security vulnerabilities. Comprehensive testing ensures reliability and prevents regressions.

### II. Frontend User Experience First (Svelte)

**Frontend development prioritizes user experience and performance over test coverage**:
- **NO automated testing required** for frontend components or pages
- Development focus MUST be on:
  - Fast initial page load (<2s on 3G)
  - Smooth interactions (60fps animations)
  - Responsive design (mobile-first approach)
  - Intuitive UX with clear visual feedback
  - Accessibility (ARIA labels, keyboard navigation)
- Manual testing and user feedback drive quality
- Performance profiling takes precedence over unit tests

**Rationale**: For a management interface with limited complexity, excellent UX and snappy performance deliver more user value than extensive test suites. Manual validation during development is sufficient for UI correctness.

### III. API Contract Stability

**Public API endpoints MUST maintain backward compatibility**:
- Breaking changes require MAJOR version bump
- All request/response schemas documented in API.md
- Changes to existing endpoints require migration plan
- New optional fields allowed; removing/renaming fields prohibited without version bump

**Rationale**: API contract stability ensures frontend and potential external integrations don't break unexpectedly. Clear versioning communicates impact.

### IV. Configuration & Environment

**Environment-based configuration following Twelve-Factor principles**:
- Defaults in `backend/internal/config/config.json`
- Runtime overrides via environment variables or `.env`
- NO hardcoded credentials or sensitive values
- All configurable values MUST be documented in API.md

**Rationale**: Enables flexible deployment (development, staging, production) without code changes and prevents credential leaks.

### V. Performance Budgets

**Strict performance requirements**:

**Backend**:
- API response time: <100ms p95 for peer operations
- /stats endpoint: <50ms p95
- WireGuard configuration updates: <200ms

**Frontend**:
- Time to Interactive (TTI): <3s on 3G
- First Contentful Paint (FCP): <1.5s
- JavaScript bundle size: <200KB (gzipped)
- Lighthouse Performance score: ≥90

**Rationale**: Management interfaces must feel instant. Slow responses degrade user trust and productivity.

### VI. Observability & Structured Logging

**All backend operations MUST be observable**:
- Structured JSON logging via `slog`
- Log levels: ERROR for failures, INFO for key operations, DEBUG for diagnostics
- Every API request logged with: method, path, status, duration
- WireGuard operations logged with: action, peer ID, outcome

**Frontend observability**:
- Browser console errors logged
- Failed API calls captured with context
- NO sensitive data (private keys, passwords) in logs

**Rationale**: Debugging production issues requires visibility. Structured logs enable efficient troubleshooting and monitoring.

## Technology Stack

**Backend**:
- Language: Go 1.21+
- Framework: Standard library `net/http`
- WireGuard: `golang.zx2c4.com/wireguard/wgctrl`
- Testing: Go standard `testing` package
- Logging: `log/slog`

**Frontend**:
- Framework: SvelteKit 2.x
- Language: TypeScript 5.x
- Styling: TailwindCSS 4.x + DaisyUI
- Build: Vite 7.x

**Deployment**:
- Backend: Linux server with WireGuard kernel module
- Frontend: Static build served via SvelteKit adapter

## Development Workflow

**Feature Development**:
1. Specification created in `/specs/[###-feature-name]/spec.md`
2. Implementation plan generated via `/speckit.plan` command
3. Backend features: TDD workflow (tests → implementation)
4. Frontend features: Direct implementation with performance focus
5. Manual integration testing before merge

**Code Review Requirements**:
- Backend: All tests passing + new tests for new functionality
- Frontend: Performance check (Lighthouse), visual QA in dev environment
- API changes: API.md documentation updated
- Breaking changes: Version bump + migration notes

**Deployment Gates**:
- Backend: All tests passing (CI/CD enforcement)
- Frontend: Build succeeds + bundle size within budget
- Integration: Manual smoke test on staging environment

## Governance

This constitution supersedes all other development practices for the WireGuard Manager project.

**Amendments**:
- Constitution changes require documentation of rationale
- Version bump according to semantic versioning rules
- All dependent templates/docs MUST be updated for consistency
- Use `/speckit.constitution` command for amendments

**Compliance**:
- All feature specifications and plans MUST verify alignment with principles
- Pull requests MUST demonstrate principle adherence
- Complexity or principle violations MUST be justified explicitly

**Version**: 1.0.0 | **Ratified**: 2026-01-31 | **Last Amended**: 2026-01-31
