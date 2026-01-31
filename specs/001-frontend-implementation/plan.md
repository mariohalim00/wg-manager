# Implementation Plan: SvelteKit Frontend Implementation

**Branch**: `feat/frontend` | **Date**: 2026-02-01 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/001-frontend-implementation/spec.md`

## Summary

Build a glassmorphism-styled VPN management dashboard using SvelteKit, TypeScript, and TailwindCSS to provide peer management, statistics visualization, and interface configuration (mock UI). The frontend connects to an existing Go backend API (GET/POST/DELETE /peers, GET /stats) with responsive layouts from desktop (1024px+) to mobile (320px+), sidebar navigation, and strict performance budgets (TTI <3s, bundle <200KB).

**Primary Technical Approach**:

- **SvelteKit 2.x** with TypeScript (strict, no `any`)
- **TailwindCSS 4.x** for utility-first styling with glassmorphism design system
- **Svelte Stores** for state management (peers, stats, UI state)
- **File-based routing** for pages (/, /peers, /stats, /settings)
- **QR code generation** for peer configuration sharing
- **Manual testing** strategy (no automated frontend tests per Constitution Principle II)

## Technical Context

**Language/Version**: TypeScript 5.x (strict mode, no `any` types)  
**Framework**: SvelteKit 2.x, Svelte 5.x  
**Primary Dependencies**:

- **TailwindCSS 4.x** (utility-first CSS)
- **@tailwindcss/forms** (form styling plugin)
- **Material Symbols Outlined** (icon library, CDN or package)
- **qrcode** or **svelte-qrcode** (QR code generation)
- **date-fns** (optional, for timestamp formatting)

**Storage**: None (frontend only, API-driven)  
**Testing**: Manual testing + TypeScript type checking (no Jest/Vitest per Constitution Principle II)  
**Target Platform**: Modern browsers (Chrome, Firefox, Safari, Edge - latest 2 versions)  
**Project Type**: Web application (frontend only, connects to existing backend)  
**Performance Goals**:

- Time to Interactive (TTI): <3 seconds on 3G
- First Contentful Paint (FCP): <1.5 seconds
- Bundle size: <200KB gzipped
- Lighthouse Performance score: ≥90
- 60fps animations for modals, transitions, status badges

**Constraints**:

- No `any` types in TypeScript (strict typing)
- Backend API is read-only (cannot modify API contract)
- Settings page is UI-only (mock data, no functional backend)
- Must match glassmorphism design aesthetic from provided mockups

**Scale/Scope**:

- 4 main pages (Dashboard, Peers, Stats, Settings)
- ~15-20 Svelte components
- 3 Svelte stores (peers, stats, UI/notifications)
- ~10-15 API integration functions
- Responsive: 320px (mobile) to 1920px+ (desktop)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

**Principle I: Backend Testing Discipline (Go)**:

- [x] N/A - This is frontend-only implementation
- [x] Backend API already tested and documented (backend/API.md)
- [x] No backend code modifications required

**Principle II: Frontend User Experience First (Svelte)**:

- [x] Frontend features focus on UX and performance (no automated test requirements per Constitution)
- [x] Performance budgets defined: TTI <3s on 3G, FCP <1.5s, bundle <200KB gzipped, Lighthouse ≥90
- [x] Manual testing approach documented (visual QA in quickstart.md, API integration testing)
- [x] Accessibility considerations: ARIA labels, keyboard navigation for modals/forms

**Principle III: API Contract Stability**:

- [x] No API changes required (consuming existing endpoints: GET/POST /peers, DELETE /peers/{id}, GET /stats)
- [x] API contract documented in backend/API.md (stable, no breaking changes)
- [x] Frontend TypeScript types match backend response schemas

**Principle IV: Configuration & Environment**:

- [x] Backend API base URL configurable via VITE_API_BASE_URL environment variable (.env file)
- [x] No hardcoded credentials or secrets in source code
- [x] Twelve-Factor configuration pattern followed

**Principle V: Performance Budgets**:

- [x] Frontend: TTI <3s on 3G, FCP <1.5s, bundle <200KB gzipped, Lighthouse ≥90
- [x] Backend: API response <100ms p95 (already validated by existing backend)

**Principle VI: Observability & Structured Logging**:

- [x] Frontend error logging to browser console (structured, with context)
- [x] API error responses captured and displayed to user (no sensitive data exposed)
- [x] No sensitive data (private keys, passwords) logged client-side

**Complexity Justification**:

- [x] No principle violations

**GATE STATUS**: ✅ PASS - All constitution principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/001-frontend-implementation/
├── spec.md                 # Feature specification (COMPLETE)
├── design-analysis.md      # Design vs. API gap analysis (COMPLETE)
├── plan.md                 # This file (implementation plan - COMPLETE)
├── tasks.md                # Task breakdown (COMPLETE)
├── design/                 # Design mockups (PROVIDED)
│   └── stitch_vpn_management_dashboard/
└── checklists/
    └── requirements.md     # Spec validation (COMPLETE)
```

**Note**: research.md, data-model.md, quickstart.md, and contracts/ are optional artifacts. Technology decisions are documented in Phase 0 above. TypeScript interfaces defined inline in Phase 1. Developer guide will be created as task T083 during implementation.

### Source Code (repository root)

```text
wg-manager/
├── src/                              # SvelteKit source
│   ├── routes/                       # File-based routing
│   │   ├── +layout.svelte           # Root layout (sidebar, glassmorphism base)
│   │   ├── +layout.ts               # Root load function (fetch stats for sidebar)
│   │   ├── +page.svelte             # Dashboard (/) - stats cards + peer table
│   │   ├── +error.svelte            # Global error page
│   │   ├── peers/
│   │   │   └── +page.svelte         # Peers page (/peers) - peer management
│   │   ├── stats/
│   │   │   └── +page.svelte         # Stats page (/stats) - detailed statistics
│   │   └── settings/
│   │       └── +page.svelte         # Settings page (/settings) - mock UI
│   ├── lib/
│   │   ├── components/              # Reusable Svelte components
│   │   │   ├── Sidebar.svelte       # Sidebar navigation
│   │   │   ├── PeerTable.svelte     # Peer list table
│   │   │   ├── PeerModal.svelte     # Add peer form modal
│   │   │   ├── QRCodeDisplay.svelte # QR code modal
│   │   │   ├── StatusBadge.svelte   # Online/offline indicator
│   │   │   ├── Notification.svelte  # Toast notification
│   │   │   ├── StatsCard.svelte     # Dashboard stats card
│   │   │   ├── ConfirmDialog.svelte # Delete confirmation
│   │   │   └── LoadingSpinner.svelte # Loading state
│   │   ├── stores/                  # Svelte stores (state management)
│   │   │   ├── peers.ts             # Peer list + CRUD operations
│   │   │   ├── stats.ts             # Interface statistics
│   │   │   └── notifications.ts     # Toast notifications queue
│   │   ├── api/                     # API client functions
│   │   │   ├── client.ts            # Base HTTP client (fetch wrapper)
│   │   │   ├── peers.ts             # Peer API calls
│   │   │   └── stats.ts             # Stats API calls
│   │   ├── utils/                   # Utility functions
│   │   │   ├── validation.ts        # CIDR validation, form validators
│   │   │   ├── formatting.ts        # Byte formatting, date formatting
│   │   │   └── constants.ts         # App constants (API base URL, etc.)
│   │   ├── types/                   # TypeScript types
│   │   │   ├── peer.ts              # Peer interface
│   │   │   ├── stats.ts             # Stats interface
│   │   │   └── api.ts               # API response types
│   │   └── index.ts                 # Barrel export
│   ├── app.html                     # HTML template
│   ├── app.css                      # Global styles (Tailwind directives + glassmorphism)
│   └── app.d.ts                     # SvelteKit type declarations
├── static/                           # Static assets
│   ├── robots.txt
│   └── favicon.ico
├── svelte.config.js                 # SvelteKit configuration
├── tailwind.config.js               # Tailwind configuration (glassmorphism theme)
├── vite.config.ts                   # Vite configuration
├── tsconfig.json                    # TypeScript strict configuration
├── eslint.config.js                 # ESLint + Svelte + TypeScript rules
├── .prettierrc                      # Prettier configuration
├── package.json                     # Dependencies
└── .env.example                     # Environment variable template
```

**Structure Decision**: Web application (frontend only). Uses SvelteKit's conventional file-based routing structure with clear separation of routes, components, stores, and API clients. TypeScript types are co-located with implementation files for better maintainability.

## Complexity Tracking

> No constitution violations - this section is empty.

---

## Phase 0: Research & Technology Decisions

**Status**: Research completed prior to planning. Decisions documented below (no separate research.md artifact required).

### Technology Decisions (Pre-Made)

1. **Glassmorphism with TailwindCSS**
   - **Decision**: TailwindCSS utilities (`backdrop-blur-md`, `bg-white/10`) + custom CSS in app.css for complex gradients
   - **Rationale**: Native Tailwind support, no plugin needed, performant

2. **QR Code Library**
   - **Decision**: `svelte-qrcode@^1.0.0` (approved in dependencies)
   - **Rationale**: Svelte-native, lightweight (~15KB), matches ecosystem

3. **Icon Library**
   - **Decision**: Material Symbols Outlined via CDN (Google Fonts)
   - **Rationale**: Zero build overhead, matches design mockups, fallback to npm package if offline needed

4. **State Management Pattern**
   - **Decision**: Writable stores (peers, stats, notifications) with derived stores for computed state
   - **Rationale**: Native Svelte pattern, reactive, no external library needed

5. **Form Validation Strategy**
   - **Decision**: Custom regex validation in `src/lib/utils/validation.ts`
   - **Rationale**: Lightweight, no dependencies, sufficient for CIDR validation

6. **Responsive Navigation Pattern**
   - **Decision**: Conditional rendering with Tailwind breakpoints (`lg:block`, `lg:hidden`)
   - **Rationale**: Declarative, no JavaScript state management needed

7. **Date/Time Formatting**
   - **Decision**: `Intl.RelativeTimeFormat` (native browser API)
   - **Rationale**: Zero dependencies, standard API, sufficient for "2 minutes ago" formatting

8. **Environment Configuration**
   - **Decision**: Vite env vars (`VITE_API_BASE_URL`) with `.env` file, fallback to `http://localhost:8080`
   - **Rationale**: Standard Vite pattern, build-time replacement, supports all deployment targets

---

## Phase 1: Design & Contracts

### Data Model (data-model.md)

**Entities** (TypeScript interfaces):

```typescript
// src/lib/types/peer.ts
export interface Peer {
	id: string; // PublicKey (unique identifier)
	publicKey: string; // WireGuard public key (base64)
	name: string; // User-friendly name
	endpoint?: string; // Peer's public IP:port (optional, set by kernel)
	allowedIPs: string[]; // CIDR notation array (e.g., ["10.0.0.2/32"])
	lastHandshake: string; // ISO timestamp or "0" (never connected)
	receiveBytes: number; // Total bytes received
	transmitBytes: number; // Total bytes transmitted
	status: 'online' | 'offline'; // Derived from lastHandshake (client-side)
}

export interface PeerFormData {
	name: string;
	allowedIPs: string[]; // CIDR strings (validated before submit)
	publicKey?: string; // Optional (backend generates if omitted)
}

export interface PeerCreateResponse {
	id: string;
	publicKey: string;
	name: string;
	allowedIPs: string[];
	config: string; // WireGuard .conf file content
	privateKey?: string; // Only if backend generated keypair
}

// src/lib/types/stats.ts
export interface InterfaceStats {
	interfaceName: string; // e.g., "wg0"
	peerCount: number; // Total peers
	totalRx: number; // Total bytes received
	totalTx: number; // Total bytes transmitted
}

// src/lib/types/api.ts
export interface APIError {
	error: string; // User-facing error message
	details?: string; // Optional detailed error
}

export interface APIResponse<T> {
	data?: T;
	error?: APIError;
	status: number;
}

// src/lib/types/notification.ts
export interface Notification {
	id: string; // Unique ID for dismissal
	type: 'success' | 'error' | 'warning' | 'info';
	message: string;
	duration?: number; // Auto-dismiss duration (ms), undefined = manual dismiss
}
```

**Validation Rules**:

- `Peer.name`: Required, non-empty string (trimmed)
- `Peer.allowedIPs`: Required, array of valid CIDR strings (e.g., `10.0.0.2/32`)
- `Peer.status`: Derived from `lastHandshake` (online if within last 2-3 minutes)

### API Contracts (contracts/)

**API Client Interface** (`contracts/api-client.ts`):

```typescript
// Base HTTP client wrapper
export interface HTTPClient {
	get<T>(url: string): Promise<APIResponse<T>>;
	post<T>(url: string, body: unknown): Promise<APIResponse<T>>;
	delete<T>(url: string): Promise<APIResponse<T>>;
}

// Peer API
export interface PeerAPI {
	listPeers(): Promise<APIResponse<Peer[]>>;
	addPeer(data: PeerFormData): Promise<APIResponse<PeerCreateResponse>>;
	removePeer(peerId: string): Promise<APIResponse<void>>;
}

// Stats API
export interface StatsAPI {
	getStats(): Promise<APIResponse<InterfaceStats>>;
}
```

**API Endpoints** (from backend/API.md):

| Endpoint      | Method | Request                            | Response             | Status Codes |
| ------------- | ------ | ---------------------------------- | -------------------- | ------------ |
| `/peers`      | GET    | -                                  | `Peer[]`             | 200          |
| `/peers`      | POST   | `{ name, publicKey?, allowedIPs }` | `PeerCreateResponse` | 201, 400     |
| `/peers/{id}` | DELETE | -                                  | -                    | 204, 404     |
| `/stats`      | GET    | -                                  | `InterfaceStats`     | 200          |

**Error Handling**:

- 400 Bad Request: Validation error (display inline in form)
- 404 Not Found: Peer not found (show notification)
- 500 Internal Server Error: Backend failure (show generic error notification)
- Network errors: Timeout, unreachable (show connection error notification)

### Component Architecture

**Component Hierarchy**:

```
App (+layout.svelte)
├── Sidebar
│   ├── Navigation links (Dashboard, Peers, Stats, Settings)
│   └── UsageWidget (mock data)
├── Routes
│   ├── Dashboard (+page.svelte)
│   │   ├── StatsCard (3x: Active Peers, Total RX, Total TX)
│   │   └── PeerTable
│   │       ├── StatusBadge
│   │       └── Actions (Download, QR, Delete)
│   ├── Peers (/peers/+page.svelte)
│   │   ├── PeerTable
│   │   ├── PeerModal (Add peer form)
│   │   └── QRCodeDisplay (Modal)
│   ├── Stats (/stats/+page.svelte)
│   │   └── StatsCard (Detailed metrics)
│   └── Settings (/settings/+page.svelte)
│       └── MockConfigurationPanel (Read-only)
├── Notification (Toast stack)
└── ConfirmDialog (Delete confirmation)
```

**Component Responsibilities**:

- **Sidebar**: Navigation, branding, usage widget (mock)
- **PeerTable**: Display peer list, row actions (hover-reveal on desktop, always-visible on mobile)
- **PeerModal**: Add peer form with validation, submit to API, display QR/config on success
- **QRCodeDisplay**: Full-screen modal showing QR code + download button
- **StatusBadge**: Online/offline indicator with pulse animation
- **Notification**: Toast notifications (auto-dismiss or manual close)
- **ConfirmDialog**: Generic confirmation dialog for destructive actions

### Quickstart Guide (quickstart.md)

**Developer Setup**:

1. **Prerequisites**: Node.js 18+, npm 9+
2. **Install dependencies**: `npm install`
3. **Configure backend URL**: Copy `.env.example` to `.env`, set `VITE_API_BASE_URL=http://localhost:8080`
4. **Start dev server**: `npm run dev` (opens http://localhost:5173)
5. **Build for production**: `npm run build` (output in `build/`)
6. **Preview production build**: `npm run preview`

**Key Commands**:

```bash
npm run dev           # Start dev server (http://localhost:5173)
npm run build         # Build for production (SSR/SPA)
npm run preview       # Preview production build
npm run lint          # Run ESLint
npm run format        # Run Prettier
npm run check         # TypeScript type checking
```

**Environment Variables**:

```bash
VITE_API_BASE_URL=http://localhost:8080  # Backend API base URL
```

**Manual Testing Checklist**:

- [ ] Dashboard loads within 3 seconds
- [ ] Peer list displays correctly (online/offline status, transfer stats)
- [ ] Add peer form validates CIDR notation
- [ ] QR code displays after successful peer creation
- [ ] Delete peer shows confirmation dialog
- [ ] Stats page displays aggregate metrics
- [ ] Settings page displays mock configuration (read-only)
- [ ] Responsive layout works on mobile (320px), tablet (768px), desktop (1024px+)
- [ ] Navigation switches to bottom nav on mobile (<1024px)
- [ ] Hover actions work on desktop, always visible on mobile
- [ ] Lighthouse Performance score ≥90
- [ ] Bundle size <200KB gzipped

---

## Phase 2: Task Breakdown

**Deferred to `/speckit.tasks` command** - Will generate detailed tasks after this plan is approved.

**High-level task categories** (preview):

1. **Project Setup** (SvelteKit, TypeScript, TailwindCSS, ESLint, Prettier)
2. **Design System** (Glassmorphism theme, Tailwind config, global styles)
3. **Layout & Navigation** (Sidebar, responsive breakpoints, bottom nav)
4. **Type Definitions** (Peer, Stats, API response types)
5. **API Client** (HTTP client, error handling, peer/stats endpoints)
6. **Svelte Stores** (peers store, stats store, notifications store)
7. **Utility Functions** (CIDR validation, byte formatting, date formatting)
8. **Core Components** (PeerTable, StatusBadge, StatsCard, Notification)
9. **Form & Modals** (PeerModal, QRCodeDisplay, ConfirmDialog)
10. **Pages** (Dashboard, Peers, Stats, Settings)
11. **Integration Testing** (Manual QA against backend API)
12. **Performance Optimization** (Bundle size, lazy loading, Lighthouse audit)

---

## Dependencies & Package Requirements

**Core Dependencies** (require approval):

```json
{
	"@sveltejs/kit": "^2.0.0",
	"svelte": "^5.0.0",
	"typescript": "^5.0.0",
	"tailwindcss": "^4.0.0",
	"@tailwindcss/forms": "^0.5.0",
	"svelte-qrcode": "^1.0.0",
	"vite": "^7.0.0"
}
```

**Dev Dependencies**:

```json
{
	"@sveltejs/adapter-auto": "^3.0.0",
	"@sveltejs/vite-plugin-svelte": "^4.0.0",
	"eslint": "^9.0.0",
	"eslint-plugin-svelte": "^2.0.0",
	"@typescript-eslint/eslint-plugin": "^8.0.0",
	"@typescript-eslint/parser": "^8.0.0",
	"prettier": "^3.0.0",
	"prettier-plugin-svelte": "^3.0.0",
	"prettier-plugin-tailwindcss": "^0.6.0"
}
```

**Optional Dependencies** (lightweight alternatives available):

- `date-fns` (for rich date formatting, can use native `Intl.RelativeTimeFormat` instead)
- `@material-symbols/svg-400` (for offline icon support, can use CDN instead)

**Total estimated bundle size**: ~150-180KB gzipped (within <200KB budget)

---

## Risks & Mitigation

| Risk                                             | Impact                              | Mitigation                                                                                     |
| ------------------------------------------------ | ----------------------------------- | ---------------------------------------------------------------------------------------------- |
| **Bundle size exceeds 200KB**                    | High (performance budget violation) | Use Vite bundle analyzer, lazy-load heavy components (QR modal), tree-shake unused TailwindCSS |
| **CIDR validation edge cases**                   | Medium (user confusion)             | Comprehensive regex testing, clear error messages with examples                                |
| **Backend API unavailable during dev**           | Medium (blocks integration testing) | Use mock data in stores until backend is available, feature flag for mock mode                 |
| **Glassmorphism performance on low-end devices** | Low (degraded UX)                   | Use `@supports` CSS queries to disable backdrop blur on unsupported browsers                   |
| **TypeScript strict mode friction**              | Low (dev velocity)                  | Use `satisfies` operator, utility types (`Partial<T>`, `Pick<T>`), avoid type assertions       |
| **Responsive layout complexity**                 | Medium (maintenance overhead)       | Use Tailwind responsive utilities (`lg:`, `md:`, `sm:`), test at breakpoint boundaries         |

---

## Success Criteria

**Implementation is complete when**:

- ✅ All functional requirements (FR-001 through FR-021) are implemented
- ✅ All user stories (P1, P2) have acceptance scenarios passing
- ✅ Performance budgets met (TTI <3s, FCP <1.5s, bundle <200KB, Lighthouse ≥90)
- ✅ Responsive design works on mobile (320px), tablet (768px), desktop (1024px+)
- ✅ Glassmorphism design matches provided mockups
- ✅ Manual testing checklist (quickstart.md) passes
- ✅ TypeScript compiles without errors (strict mode, no `any`)
- ✅ ESLint and Prettier pass

**Definition of Done**:

- Code merged to `1-frontend-implementation` branch
- Manual QA completed against backend API
- Performance audit (Lighthouse) results documented
- Quickstart guide validated by another developer
- Ready for production deployment

---

## Next Steps

1. **Approve this plan** - Review and confirm technical approach, dependencies, and structure
2. **Run `/speckit.tasks`** - Generate detailed task breakdown from this plan
3. **Setup project** - Initialize SvelteKit, install dependencies, configure Tailwind
4. **Implement Phase 1** - Build layout, navigation, type definitions, API client
5. **Implement Phase 2** - Build components, stores, pages
6. **Manual QA** - Test against backend, verify performance budgets
7. **Deploy** - Build for production, deploy to hosting platform
