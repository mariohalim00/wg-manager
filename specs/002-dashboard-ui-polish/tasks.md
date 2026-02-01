---

description: "Task list for Dashboard UI Polish"
---

# Tasks: Dashboard UI Polish

**Input**: Design documents from `specs/002-dashboard-ui-polish/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**:
- **Backend (Go)**: Not applicable (no backend changes)
- **Frontend (Svelte)**: Manual UX validation only (per Constitution II)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 Capture baseline dashboard structure notes in specs/002-dashboard-ui-polish/quickstart.md (add a “Before” checklist with 4 bullets: header, status cards, traffic cards, peers table)

---

## Phase 2: Foundational (Blocking Prerequisites)

- [x] T002 Define dashboard UI tokens and shared utility classes in src/app.css (dark/light surfaces, borders, typography scale, focus rings)
- [x] T003 [P] Add reusable status/metric card class patterns in src/app.css (card shells, value typography, tabular numerals)
- [x] T004 [P] Add reusable table header/footer styling helpers in src/app.css (section headers, count badges, footer links)

---

## Phase 3: User Story 1 - Confident Dashboard Overview (Priority: P1) 🎯 MVP

**Goal**: Deliver a scannable dashboard with app shell, status cards, traffic summaries, and peers table.

**Independent Test**: Open the dashboard and verify header, status cards, traffic cards, and peers table are visible, readable, and aligned to the spec.

### Implementation for User Story 1

- [x] T005 [US1] Implement the dashboard header layout and status card grid in src/routes/+page.svelte (interface context, visual-only header controls, 4-card grid)
- [x] T006 [US1] Apply placeholder values for Public Key, Listening Port, Subnet in src/routes/+page.svelte (keep cards visible)
- [x] T007 [US1] Create static traffic cards with decorative charts in src/routes/+page.svelte (structure ready for future data binding)
- [x] T008 [US1] Refine status card styling in src/lib/components/StatsCard.svelte to match hierarchy and palette
- [x] T009 [US1] Refine peers table visual hierarchy in src/lib/components/PeerTable.svelte (row states, device subtitle, transfer arrows)

**Checkpoint**: User Story 1 is visually complete and matches the dashboard layout goals.

---

## Phase 4: User Story 2 - Responsive Operations (Priority: P2)

**Goal**: Ensure the dashboard remains usable on mobile and tablet breakpoints.

**Independent Test**: Resize viewport to mobile and tablet; verify sidebar behavior, card layouts, and table scrolling.

### Implementation for User Story 2

- [x] T010 [US2] Update responsive grid behavior in src/routes/+page.svelte (4/2/1 card grid, stacked charts on small screens)
- [x] T011 [US2] Ensure peers table horizontal scrolling on tablet in src/lib/components/PeerTable.svelte
- [x] T012 [US2] Confirm sidebar drawer behavior on mobile in src/lib/components/Sidebar.svelte (no behavior changes)

**Checkpoint**: User Story 2 is responsive across mobile/tablet/desktop.

---

## Phase 5: User Story 3 - Consistent Visual Semantics (Priority: P3)

**Goal**: Ensure consistent palette, typography, and interaction feedback in dark and light modes.

**Independent Test**: Toggle dark/light mode and confirm consistent semantics, focus/hover states, and text hierarchy.

### Implementation for User Story 3

- [x] T013 [US3] Align color usage and hover/focus states in src/app.css (primary/secondary palettes, neutral grays, soft rings)
- [x] T014 [US3] Standardize typography scale in src/routes/+page.svelte and src/lib/components/StatsCard.svelte
- [x] T015 [US3] Add empty-state styling in src/lib/components/PeerTable.svelte when no peers exist

**Checkpoint**: User Story 3 satisfies visual consistency and empty-state requirements.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T016 [P] Run manual UI validation using specs/002-dashboard-ui-polish/quickstart.md and record pass/fail against all 6 checklist bullets
- [ ] T017 [P] Audit for unintended behavior changes in src/routes/+page.svelte and src/lib/components/Sidebar.svelte (confirm visual-only controls)
- [ ] T018 [P] Verify performance budget impact by reviewing bundle output after build (no new dependencies)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Stories (Phases 3–5)**: Depend on Foundational completion
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational
- **US2 (P2)**: Starts after Foundational
- **US3 (P3)**: Starts after Foundational

### Parallel Opportunities

- T003 and T004 can run in parallel (different shared styles)
- US2 and US3 tasks can proceed in parallel after US1 is stable
- Phase 6 tasks can run in parallel

---

## Parallel Example: User Story 1

- T005 (header + status grid) can be done alongside T009 (peer table styling)
- T007 (traffic cards) can be done alongside T008 (StatsCard styling)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 (Setup)
2. Complete Phase 2 (Foundational)
3. Complete Phase 3 (US1)
4. Validate US1 visually per quickstart

### Incremental Delivery

1. Add US2 for responsive behaviors
2. Add US3 for semantic consistency and empty state
3. Finish Polish phase
