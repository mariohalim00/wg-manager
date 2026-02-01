# Implementation Plan: Dashboard UI Polish

**Branch**: `002-dashboard-ui-polish` | **Date**: 2026-02-01 | **Spec**: [specs/002-dashboard-ui-polish/spec.md](specs/002-dashboard-ui-polish/spec.md)
**Input**: Feature specification from `specs/002-dashboard-ui-polish/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Polish the dashboard UI to match a modern infrastructure dashboard: app shell layout, status cards, traffic cards with decorative charts, and a refined peers table. This is a frontend-only visual enhancement with no backend or business logic changes. All missing data fields will display placeholders, and charts remain static but structured for future data binding.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: TypeScript 5.x, Svelte 5.x (SvelteKit 2.x)  
**Primary Dependencies**: SvelteKit, TailwindCSS 4.x, DaisyUI  
**Storage**: N/A (UI-only changes)  
**Testing**: Manual UI testing only (per Constitution II)  
**Target Platform**: Web (modern browsers)
**Project Type**: Web application (SvelteKit frontend + Go backend, no backend changes)  
**Performance Goals**: 60 fps interactions, TTI <3s on 3G  
**Constraints**: Bundle <200KB gzipped, no new behavior or API changes  
**Scale/Scope**: Single dashboard page polish and supporting components

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Backend Testing Discipline (Principle I)**:
- [x] No backend changes in scope

**Frontend User Experience First (Principle II)**:
- [x] Frontend features focus on UX and performance (no test requirements)
- [x] Performance budgets defined (TTI, FCP, bundle size)
- [x] Manual testing approach documented
- [x] Accessibility considerations noted

**API Contract Stability (Principle III)**:
- [x] No API changes required

**Configuration & Environment (Principle IV)**:
- [x] No configuration changes required

**Performance Budgets (Principle V)**:
- [x] Backend: N/A
- [x] Frontend: <3s TTI, <1.5s FCP, <200KB bundle (gzipped)
- [x] Lighthouse target: ≥90

**Observability & Structured Logging (Principle VI)**:
- [x] No backend logging changes required

**Complexity Justification**:
- [x] No principle violations

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
backend/
└── cmd/server/

src/
├── routes/
│   └── +page.svelte
├── lib/
│   ├── components/
│   │   ├── StatsCard.svelte
│   │   ├── PeerTable.svelte
│   │   └── Sidebar.svelte
│   └── utils/
└── app.css
```

**Structure Decision**: Web app. Changes are confined to the SvelteKit frontend (`src/`) with no backend modifications.

## Phase 0: Research (completed)

- Output: [specs/002-dashboard-ui-polish/research.md](specs/002-dashboard-ui-polish/research.md)
- All clarifications resolved (placeholders, header controls, charts, peers data).

## Phase 1: Design & Contracts (completed)

- Data model: [specs/002-dashboard-ui-polish/data-model.md](specs/002-dashboard-ui-polish/data-model.md)
- Contracts: [specs/002-dashboard-ui-polish/contracts/README.md](specs/002-dashboard-ui-polish/contracts/README.md)
- Quickstart: [specs/002-dashboard-ui-polish/quickstart.md](specs/002-dashboard-ui-polish/quickstart.md)

### Constitution Check (post-design)

- Principle I (Backend TDD): N/A, no backend changes.
- Principle II (UX-first): UI polish focus, manual testing only, performance budgets honored.
- Principle III (API stability): No API changes.
- Principle IV (Config): No config changes.
- Principle V (Performance): Budget constraints maintained.
- Principle VI (Observability): No logging changes.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
