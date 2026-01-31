# Tasks: SvelteKit Frontend Implementation

**Input**: Design documents from `/specs/001-frontend-implementation/`  
**Prerequisites**: plan.md (✅), spec.md (✅), design-analysis.md (✅)

**Tests**: 
- **Backend (Go)**: Tests are MANDATORY per Constitution Principle I (Backend Testing Discipline). Write tests FIRST before implementation.
- **Frontend (Svelte)**: Tests are NOT required per Constitution Principle II (Frontend User Experience First). Focus on UX and performance instead.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for SvelteKit + TypeScript + TailwindCSS

- [X] T001 Install approved dependencies: `svelte-qrcode@^1.0.0`, `@tailwindcss/forms@^0.5.0`
- [X] T002 Configure TypeScript strict mode in `tsconfig.json` (no `any` types allowed, `strict: true`)
- [X] T003 [P] Update `tailwind.config.js` with glassmorphism theme (backdrop-blur utilities, custom colors, opacity scales)
- [X] T004 [P] Configure ESLint with TypeScript and Svelte rules in `eslint.config.js`
- [X] T005 [P] Configure Prettier with Svelte and Tailwind plugins in `.prettierrc`
- [X] T006 Create `.env.example` with `VITE_API_BASE_URL=http://localhost:8080` template
- [X] T007 Copy `.env.example` to `.env` for local development

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T008 Create TypeScript interfaces in `src/lib/types/peer.ts` (Peer, PeerFormData, PeerCreateResponse)
- [X] T009 [P] Create TypeScript interfaces in `src/lib/types/stats.ts` (InterfaceStats)
- [X] T010 [P] Create TypeScript interfaces in `src/lib/types/api.ts` (APIError, APIResponse)
- [X] T011 [P] Create TypeScript interfaces in `src/lib/types/notification.ts` (Notification)
- [X] T012 Implement base HTTP client in `src/lib/api/client.ts` (fetch wrapper with error handling, base URL from env)
- [X] T013 [P] Implement peer API client in `src/lib/api/peers.ts` (listPeers, addPeer, removePeer functions)
- [X] T014 [P] Implement stats API client in `src/lib/api/stats.ts` (getStats function)
- [X] T015 Create peers Svelte store in `src/lib/stores/peers.ts` (writable store with load, add, remove methods)
- [X] T016 [P] Create stats Svelte store in `src/lib/stores/stats.ts` (writable store with load method)
- [X] T017 [P] Create notifications Svelte store in `src/lib/stores/notifications.ts` (writable store with add, remove, auto-dismiss logic)
- [X] T018 Implement CIDR validation utility in `src/lib/utils/validation.ts` (validateCIDR function with regex)
- [X] T019 [P] Implement byte formatting utility in `src/lib/utils/formatting.ts` (formatBytes function: B, KB, MB, GB, TB)
- [X] T020 [P] Implement date formatting utility in `src/lib/utils/formatting.ts` (formatLastHandshake: relative time or "Never")
- [X] T021 [P] Define app constants in `src/lib/utils/constants.ts` (API base URL, handshake timeout threshold)
- [X] T022 Create global glassmorphism styles in `src/app.css` (backdrop-filter utilities, panel base styles, gradient backgrounds)
- [X] T023 Create root layout in `src/routes/+layout.svelte` (glassmorphism wrapper, gradient background, notification stack)
- [X] T024 Create Sidebar component in `src/lib/components/Sidebar.svelte` (navigation links for Dashboard, Peers, Stats, Settings with active state highlighting; usage widget placeholder at bottom; glassmorphism styling with backdrop-blur; responsive: full sidebar ≥1024px, simplified <1024px)
- [X] T025 Create LoadingSpinner component in `src/lib/components/LoadingSpinner.svelte` (centered spinner with glassmorphism)
- [X] T026 [P] Create Notification component in `src/lib/components/Notification.svelte` (toast notification with auto-dismiss, type-based styling)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - View Peers (Priority: P1) 🎯 MVP

**Goal**: Display list of WireGuard peers with real-time status and transfer statistics

**Independent Test**: Navigate to `/peers` page, see peer list with names, status badges (online/offline), allowed IPs, transfer stats, and last handshake time

### Tests for User Story 1 ⚠️

> **FRONTEND**: No tests required per Constitution Principle II. Manual testing checklist in quickstart.md.

### Implementation for User Story 1

- [X] T027 [P] [US1] Create StatusBadge component in `src/lib/components/StatusBadge.svelte` (online/offline indicator with pulse animation)
- [X] T028 [P] [US1] Create PeerTable component in `src/lib/components/PeerTable.svelte` (table layout with peer rows, action buttons with responsive visibility)
- [X] T029 [US1] Create Peers page in `src/routes/peers/+page.svelte` (load peers store, render PeerTable, handle loading state)
- [X] T030 [US1] Add status derivation logic in `src/lib/stores/peers.ts` (compute online/offline from lastHandshake timestamp)
- [X] T031 [US1] Add responsive action button logic in `src/lib/components/PeerTable.svelte` (implements FR-001a: always-visible with display:flex <1024px, hover-reveal ≥1024px using Tailwind group and group-hover:opacity-100 transition)
- [ ] T032 [US1] Manual Test: Peer list rendering with sample data (verify status badges show correct online/offline, byte formatting uses KB/MB/GB, date formatting shows relative time)
- [ ] T033 [US1] Manual Test: Empty state when no peers exist (display "No peers configured" message with "Add Peer" CTA button)

**Checkpoint**: At this point, User Story 1 should be fully functional - peer list displays with correct data, status, and formatting

---

## Phase 4: User Story 2 - Add Peer (Priority: P1) 🎯 MVP

**Goal**: Add new WireGuard peers with name and allowed IPs via modal form

**Independent Test**: Click "Add Peer" button, fill form with name and CIDR notation IPs, submit, see new peer in list with success notification

### Tests for User Story 2 ⚠️

> **FRONTEND**: No tests required per Constitution Principle II. Manual testing checklist in quickstart.md.

### Implementation for User Story 2

- [X] T034 [P] [US2] Create PeerModal component in `src/lib/components/PeerModal.svelte` (glassmorphism modal with form fields, validation, loading state)
- [X] T035 [US2] Add form validation logic in PeerModal (validate name required, validate CIDR notation for each allowed IP)
- [X] T036 [US2] Wire PeerModal submit to peers store addPeer method (POST /peers API call)
- [X] T037 [US2] Display inline validation errors in PeerModal (show CIDR format examples on error)
- [X] T038 [US2] Add success notification after peer creation (use notifications store)
- [X] T039 [US2] Add error handling for API failures in PeerModal (400 validation errors, 500 server errors)
- [X] T040 [US2] Update Peers page to include "Add Peer" button that opens PeerModal
- [ ] T041 [US2] Manual Test: Form validation with invalid CIDR notation (e.g., "10.0.0.5" without /32, "10.0.0.300/32"; verify inline error messages with examples)
- [ ] T042 [US2] Manual Test: Successful peer addition flow end-to-end (form submit, API call, list refresh, success notification appears)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - view peer list AND add new peers

---

## Phase 5: User Story 5 - Download Peer Configuration (Priority: P3)

**Goal**: Display QR code and download config file after adding peer

**Independent Test**: After adding peer, see QR code modal with scannable code and "Download Config" button that saves .conf file

### Tests for User Story 5 ⚠️

> **FRONTEND**: No tests required per Constitution Principle II. Manual testing checklist in quickstart.md.

### Implementation for User Story 5

- [X] T043 [P] [US5] Create QRCodeDisplay component in `src/lib/components/QRCodeDisplay.svelte` (full-screen glassmorphism modal with QR code from svelte-qrcode)
- [X] T044 [US5] Add download config function in QRCodeDisplay (create .conf file blob, trigger browser download)
- [X] T045 [US5] Update PeerModal to show QRCodeDisplay on successful peer creation (pass config and private key from API response)
- [X] T046 [US5] Add security note in QRCodeDisplay (warn that private key is sensitive, won't be shown again after modal close)
- [ ] T047 [US5] Manual Test: QR code generation with real peer config (verify scannable by WireGuard mobile app, QR displays correctly)
- [ ] T048 [US5] Manual Test: Config file download (verify .conf format, valid WireGuard syntax, file downloads with correct name)
- [ ] T048 [US5] Manual Test: Config file download (verify .conf format, valid WireGuard syntax, file downloads with correct name)

**Checkpoint**: At this point, User Stories 1, 2, AND 5 should work - view peers, add peers, and download configs

---

## Phase 6: User Story 3 - Remove Peer (Priority: P2)

**Goal**: Remove peers from VPN network with confirmation dialog

**Independent Test**: Click delete button on peer row, confirm in dialog, see peer removed from list with success notification

### Tests for User Story 3 ⚠️

> **FRONTEND**: No tests required per Constitution Principle II. Manual testing checklist in quickstart.md.

### Implementation for User Story 3

- [X] T049 [P] [US3] Create ConfirmDialog component in `src/lib/components/ConfirmDialog.svelte` (glassmorphism confirmation modal with cancel/confirm buttons)
- [X] T050 [US3] Add delete button to PeerTable component (trash icon with responsive visibility logic)
- [X] T051 [US3] Wire delete button to open ConfirmDialog with peer name
- [X] T052 [US3] Wire ConfirmDialog confirm action to peers store removePeer method (DELETE /peers/{id} API call)
- [X] T053 [US3] Add success notification after peer deletion (use notifications store)
- [X] T054 [US3] Add error handling for deletion failures (404 peer not found, 500 server error)
- [ ] T055 [US3] Manual Test: Deletion flow end-to-end (click delete, confirm dialog appears, click confirm, API call succeeds, list refreshes, success notification shows)
- [ ] T056 [US3] Manual Test: Cancellation flow (click delete, confirm dialog appears, click cancel, dialog closes, peer remains in list)

**Checkpoint**: At this point, User Stories 1, 2, 3, AND 5 should work - full CRUD operations for peers

---

## Phase 7: User Story 4 - View Interface Statistics (Priority: P2)

**Goal**: Display aggregate VPN interface statistics (peer count, total RX/TX)

**Independent Test**: Navigate to `/stats` page, see interface name, total peer count, total received bytes, and total transmitted bytes formatted in human-readable units

### Tests for User Story 4 ⚠️

> **FRONTEND**: No tests required per Constitution Principle II. Manual testing checklist in quickstart.md.

### Implementation for User Story 4

- [X] T057 [P] [US4] Create StatsCard component in `src/lib/components/StatsCard.svelte` (glassmorphism card with icon, label, value, formatted units)
- [X] T058 [US4] Create Stats page in `src/routes/stats/+page.svelte` (load stats store, render multiple StatsCard components)
- [X] T059 [US4] Add stats cards to Dashboard page in `src/routes/+page.svelte` (show active peers, total RX, total TX - numeric only per FR-008a)
- [ ] T060 [US4] Manual Test: Stats display with sample data (verify byte formatting uses appropriate units: KB, MB, GB, TB; values are human-readable)
- [ ] T061 [US4] Manual Test: Stats page with zero peers (show "0" peer count, "0 bytes" transfer stats with appropriate messaging)

**Checkpoint**: At this point, all P1 and P2 user stories are complete - full peer management + statistics

---

## Phase 8: Settings Page (Mock UI - FR-021)

**Goal**: Display Settings page UI with read-only interface configuration (no functional backend)

**Independent Test**: Navigate to `/settings` page, see interface configuration panel with listen port, MTU, addresses, server public key, and "Coming Soon" indicators

### Implementation for Settings Page

- [X] T062 [P] [SETTINGS] Create Settings page in `src/routes/settings/+page.svelte` (glassmorphism panels with static mock data)
- [X] T063 [SETTINGS] Add mock interface configuration display (listen port: 51820, MTU: 1420, addresses: 10.0.0.1/24, server public key: mock base64 string; all fields read-only/disabled)
- [X] T064 [SETTINGS] Add visual indicators for "Coming Soon" features (add badge/chip with "Coming Soon" text, disable service control buttons with tooltip, show placeholder usage quota widget with mock data and "Preview" label)
- [X] T065 [SETTINGS] Style Settings page with glassmorphism panels matching design mockups

**Checkpoint**: Settings page displays with mock data, visual consistency with other pages

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories, performance optimization, documentation

- [X] T066 [P] Add Material Symbols Outlined icons via CDN in `src/app.html` (or install npm package if offline support needed)
- [X] T067 [P] Add favicon to `static/favicon.ico`
- [X] T068 [P] Update `static/robots.txt` for production deployment
- [X] T069 Optimize Tailwind config for production (purge unused styles, minimize bundle size)
- [X] T070 [P] Add loading states to all API calls (show LoadingSpinner in pages during fetch)
- [X] T071 [P] Add error boundaries to catch React-like errors in Svelte (use `+error.svelte` global error page)
- [X] T072 Create global error page in `src/routes/+error.svelte` (glassmorphism error display with back button)
- [ ] T073 Test responsive layouts at breakpoints: 320px (mobile), 768px (tablet), 1024px (desktop), 1920px (large desktop)
- [ ] T074 Test sidebar navigation on all pages (verify active link highlighting)
- [ ] T075 Test bottom navigation on mobile (<1024px) (verify icon-only nav, active state)
- [ ] T076 Test hover-reveal actions on desktop (≥1024px) (verify group-hover pattern works)
- [ ] T077 Test always-visible actions on mobile (<1024px) (verify buttons always visible)
- [X] T078 Run Lighthouse performance audit (target: ≥90 score, <200KB bundle, TTI <3s, FCP <1.5s)
- [X] T079 Optimize bundle size if needed (lazy-load heavy components, tree-shake unused code)
- [X] T080 [P] Verify TypeScript strict mode compliance (no `any` types, all parameters/returns typed)
- [X] T081 [P] Run ESLint and fix any warnings/errors across all files
- [X] T082 [P] Run Prettier to format all source files
- [X] T083 Create/update quickstart.md with developer setup instructions (prerequisites, install, dev server, build, manual testing checklist)
- [ ] T084 Test all user stories end-to-end against live backend API (verify API integration works)
- [ ] T085 Verify all FR requirements from spec.md are implemented (FR-001 through FR-021)
- [ ] T086 Verify all acceptance scenarios from user stories pass manual testing
- [ ] T087 Document any deviations from plan.md or spec.md (if any)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Settings (Phase 8)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1) - View Peers**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1) - Add Peer**: Can start after Foundational (Phase 2) - No dependencies (but naturally integrates with US1 for list refresh)
- **User Story 5 (P3) - Download Config**: Depends on User Story 2 (Add Peer) - Extends PeerModal with QR code display
- **User Story 3 (P2) - Remove Peer**: Can start after Foundational (Phase 2) - No dependencies (but integrates with US1 for list refresh)
- **User Story 4 (P2) - View Stats**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Within Each User Story

- Components before pages
- Store integration before component wiring
- Form validation before API submission
- Success/error handling after API integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel within Phase 2 (types, API clients, stores, utilities, components)
- Once Foundational phase completes, user stories can start in parallel:
  - **US1 (View Peers)** and **US2 (Add Peer)** can be worked on simultaneously (different components)
  - **US3 (Remove Peer)** and **US4 (View Stats)** can be worked on simultaneously (different pages)
  - **Settings Page (Phase 8)** can be worked on in parallel with any user story
- Components within a story marked [P] can run in parallel (e.g., StatusBadge + PeerTable in US1)
- Polish tasks marked [P] in Phase 9 can run in parallel (icons, favicon, robots.txt, type checking, linting, formatting)

---

## Parallel Example: Foundational Phase

```bash
# Launch all type definition tasks together:
Task: "Create TypeScript interfaces in src/lib/types/peer.ts"
Task: "Create TypeScript interfaces in src/lib/types/stats.ts"
Task: "Create TypeScript interfaces in src/lib/types/api.ts"
Task: "Create TypeScript interfaces in src/lib/types/notification.ts"

# Launch all API client tasks together (after base client is done):
Task: "Implement peer API client in src/lib/api/peers.ts"
Task: "Implement stats API client in src/lib/api/stats.ts"

# Launch all store tasks together:
Task: "Create stats Svelte store in src/lib/stores/stats.ts"
Task: "Create notifications Svelte store in src/lib/stores/notifications.ts"

# Launch all utility tasks together:
Task: "Implement byte formatting utility in src/lib/utils/formatting.ts"
Task: "Implement date formatting utility in src/lib/utils/formatting.ts"
Task: "Define app constants in src/lib/utils/constants.ts"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup (7 tasks)
2. Complete Phase 2: Foundational (19 tasks) - CRITICAL blocking phase
3. Complete Phase 3: User Story 1 - View Peers (7 tasks)
4. Complete Phase 4: User Story 2 - Add Peer (9 tasks)
5. **STOP and VALIDATE**: Test peer list + add peer independently
6. Deploy/demo if ready (Core functionality: ~42 tasks total)

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready (26 tasks)
2. Add User Story 1 → Test independently → Deploy/Demo (MVP: view peers)
3. Add User Story 2 → Test independently → Deploy/Demo (MVP: view + add peers)
4. Add User Story 5 → Test independently → Deploy/Demo (view + add + download)
5. Add User Story 3 → Test independently → Deploy/Demo (full CRUD)
6. Add User Story 4 → Test independently → Deploy/Demo (CRUD + stats)
7. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (26 tasks)
2. Once Foundational is done:
   - **Developer A**: User Story 1 (View Peers) + User Story 4 (View Stats)
   - **Developer B**: User Story 2 (Add Peer) + User Story 5 (Download Config)
   - **Developer C**: User Story 3 (Remove Peer) + Settings Page (Phase 8)
   - **Developer D**: Polish & Performance (Phase 9)
3. Stories complete and integrate independently

---

## Task Summary

- **Phase 1 (Setup)**: 7 tasks
- **Phase 2 (Foundational)**: 19 tasks ⚠️ CRITICAL - blocks all user stories
- **Phase 3 (US1 - View Peers)**: 7 tasks 🎯 MVP
- **Phase 4 (US2 - Add Peer)**: 9 tasks 🎯 MVP
- **Phase 5 (US5 - Download Config)**: 6 tasks
- **Phase 6 (US3 - Remove Peer)**: 8 tasks
- **Phase 7 (US4 - View Stats)**: 5 tasks
- **Phase 8 (Settings Page)**: 4 tasks
- **Phase 9 (Polish)**: 22 tasks

**Total Tasks**: 87

**MVP Tasks (Setup + Foundational + US1 + US2)**: 42 tasks

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story for traceability (US1, US2, US3, US4, US5)
- [SETTINGS] label for Settings page tasks (mock UI, FR-021)
- Each user story should be independently completable and testable
- **No frontend tests required** per Constitution Principle II - focus on UX and performance
- Manual testing checklist in quickstart.md covers all acceptance scenarios
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Performance validation at Phase 9 (Lighthouse audit, bundle size check)
- TypeScript strict mode enforced (no `any` types allowed)
