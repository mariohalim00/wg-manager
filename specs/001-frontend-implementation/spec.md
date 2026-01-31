# Feature Specification: SvelteKit Frontend Implementation

**Feature Branch**: `1-frontend-implementation`  
**Created**: 2026-02-01  
**Status**: Draft  
**Input**: User description: "I have a WireGuard manager backend and I want a SvelteKit + Tailwind frontend. The project is already scaffolded, some routes are prepared, and I have a design ready. Auth is not needed for now. Plan a simple, high-level roadmap for the frontend, including: pages (/peers, /stats), API integration to backend endpoints, state management with Svelte stores, forms for add/edit peers and peer config download, logging/validation utilities, testing strategy, and possible future extensions like real-time stats or theming."

## Clarifications

### Session 2026-02-01

- Q: Navigation structure (sidebar vs. top nav vs. hybrid) → A: Sidebar navigation (persistent left sidebar with Dashboard, Peers, Settings, Logs links)
- Q: Interface Settings page scope (exclude, require new API, mock UI, or hide link) → A: Mock/Visual only (implement UI with static/fake data, mark as "Preview" or "Coming Soon", functional implementation deferred)
- Q: Mobile responsiveness priority (desktop only, tablet+basic mobile, full mobile parity, tablet only) → A: Tablet + Basic Mobile (responsive down to 320px with simplified mobile layout: bottom nav, stacked cards, no complex touch gestures)
- Q: Dashboard stats display format (full charts, stats cards only, separate page only, mock placeholders) → A: Dashboard stats cards (no charts) - show numeric stats cards (peer count, total RX/TX) without visual graphs, charts deferred to future
- Q: Peer table hover actions (always visible, responsive visibility, swipe gestures, tap to expand) → A: Always-visible icons (mobile), hover (desktop) - action icons always show on <1024px, hover-reveal on ≥1024px

## User Scenarios & Testing

### User Story 1 - View and Manage Peers (Priority: P1)

As a WireGuard administrator, I need to view all configured peers and their connection status in a dashboard, so I can monitor the VPN network health at a glance.

**Why this priority**: Core functionality that provides immediate value. Without the ability to view peers, the application cannot serve its primary purpose. This is the foundation for all other peer management features.

**Independent Test**: Administrator can navigate to /peers page, see a list of all configured peers with their names, status (online/offline), allowed IPs, and data transfer statistics. Success is measured by the page loading within 3 seconds and displaying accurate data from the backend API.

**Acceptance Scenarios**:

1. **Given** no peers are configured, **When** administrator visits /peers page, **Then** an empty state message is displayed indicating "No peers configured" with a clear call-to-action to add the first peer
2. **Given** multiple peers exist with various connection states, **When** page loads, **Then** all peers are displayed in a table/list format with name, status badge (online/offline), allowed IPs, last handshake time, and transfer statistics (RX/TX bytes)
3. **Given** a peer's status changes from offline to online, **When** page refreshes or updates, **Then** the status badge updates to show "online" with visual indicators (color, icon)
4. **Given** administrator is viewing peer list, **When** clicking on a peer's delete/remove button, **Then** a confirmation dialog appears to prevent accidental deletion

---

### User Story 2 - Add New Peer (Priority: P1)

As a WireGuard administrator, I need to add new peers to the VPN network by providing a name and allowed IPs, so I can grant access to new users or devices.

**Why this priority**: Essential for network growth and onboarding. Without the ability to add peers, the system is read-only and limited. This is required immediately after the view functionality.

**Independent Test**: Administrator can click "Add Peer" button, fill in peer name and allowed IPs in a form, submit it, and see the new peer appear in the peer list. The system should also provide a downloadable configuration file or QR code for the client.

**Acceptance Scenarios**:

1. **Given** administrator clicks "Add Peer" button, **When** the form modal/dialog opens, **Then** fields for peer name (text) and allowed IPs (CIDR notation, text input or multi-input) are displayed with clear labels and examples
2. **Given** administrator fills valid data (name: "John's Laptop", allowedIPs: "10.0.0.5/32"), **When** form is submitted, **Then** peer is added via POST /peers API, backend generates keypair if needed, and frontend displays success notification
3. **Given** peer is successfully added, **When** response is received from backend, **Then** a downloadable configuration file (.conf format) and QR code are presented for easy client setup
4. **Given** administrator enters invalid CIDR notation (e.g., "10.0.0.5" without /32), **When** attempting to submit, **Then** inline validation error appears explaining correct CIDR format with examples
5. **Given** administrator leaves peer name empty, **When** attempting to submit, **Then** validation error appears indicating "Name is required"
6. **Given** backend API returns error (e.g., duplicate public key, invalid data), **When** submission fails, **Then** user-friendly error message is displayed in a notification/alert

---

### User Story 3 - Remove Peer (Priority: P2)

As a WireGuard administrator, I need to remove peers from the VPN network when access should be revoked, so I can maintain security and control over who can connect.

**Why this priority**: Important for security and access management, but not required for initial setup. Can be added after view and add functionality are stable.

**Independent Test**: Administrator can click a delete/remove button on a peer entry, confirm the deletion in a dialog, and see the peer disappear from the list immediately.

**Acceptance Scenarios**:

1. **Given** peer list is displayed, **When** administrator clicks delete/remove icon on a peer, **Then** a confirmation dialog appears with peer name and warning message ("Are you sure you want to remove [Peer Name]? This action cannot be undone.")
2. **Given** confirmation dialog is open, **When** administrator confirms deletion, **Then** DELETE /peers/{id} API is called, peer is removed, and success notification is shown
3. **Given** peer is successfully removed, **When** page updates, **Then** peer no longer appears in the list and total peer count decreases
4. **Given** backend API fails to remove peer (e.g., peer not found, permission error), **When** deletion request completes, **Then** error notification is displayed without removing peer from UI

---

### User Story 4 - View Interface Statistics (Priority: P2)

As a WireGuard administrator, I need to view aggregate statistics for the VPN interface (total peers, total RX/TX data), so I can monitor overall network usage and performance.

**Why this priority**: Valuable for monitoring and reporting, but not critical for basic peer management. Can be implemented after core CRUD operations are functional.

**Independent Test**: Administrator can navigate to /stats page and see aggregate metrics (total peers, total bytes received/transmitted) fetched from GET /stats endpoint.

**Acceptance Scenarios**:

1. **Given** administrator navigates to /stats page, **When** page loads, **Then** interface name (e.g., "wg0"), total peer count, total received bytes, and total transmitted bytes are displayed prominently
2. **Given** data transfer statistics are large numbers, **When** displayed, **Then** values are formatted in human-readable units (KB, MB, GB, TB) rather than raw bytes
3. **Given** administrator is viewing /stats page, **When** no peers are configured, **Then** peer count shows "0" and transfer stats show "0 bytes" with appropriate messaging

---

### User Story 5 - Download Peer Configuration (Priority: P3)

As a WireGuard administrator, I need to download peer configuration files after adding a peer, so I can easily distribute VPN credentials to users via multiple methods (file download, QR code).

**Why this priority**: Enhances user experience and makes onboarding easier, but the core add functionality can work without this if credentials are displayed as text initially.

**Independent Test**: After adding a peer, administrator can click "Download Config" button to save a .conf file or display QR code that can be scanned by WireGuard mobile clients.

**Acceptance Scenarios**:

1. **Given** peer is successfully added, **When** response includes config and private key, **Then** a modal/section displays both "Download Config" button and QR code image
2. **Given** administrator clicks "Download Config", **When** button is clicked, **Then** a .conf file is downloaded with proper WireGuard client configuration format (see backend API response for config content)
3. **Given** QR code is displayed, **When** scanned with WireGuard mobile app, **Then** peer can successfully import the configuration and connect to VPN
4. **Given** configuration contains sensitive private key, **When** modal is closed or peer list is refreshed, **Then** private key is no longer accessible (security best practice)

---

### Edge Cases

- What happens when backend API is unreachable or returns 500 errors? (Display user-friendly error notification, avoid app crash, allow retry)
- How does system handle invalid CIDR notation input variations (missing /32, incorrect format like "10.0.0.300/32")? (Client-side validation with clear error messages before submission)
- What if peer handshake timestamp is very old (years ago) or null/undefined? (Display "Never" or "Offline" status, handle gracefully without breaking UI)
- How to handle concurrent administrators removing the same peer? (Backend returns 404/error, frontend shows notification that peer was already removed)
- What if allowed IPs array is empty in API response? (Display "None" or "-" in UI, flag as configuration issue)
- How to format very large transfer statistics (petabytes)? (Use appropriate unit scaling: PB, EB)
- What happens when browser local storage or state is corrupted? (Graceful degradation, refetch from backend API on page load)

## Requirements

### Functional Requirements

- **FR-001**: System MUST display a paginated or scrollable list of all WireGuard peers fetched from GET /peers endpoint
- **FR-001a**: System MUST display peer action buttons (Download Config, View QR Code, Delete) with responsive visibility: always visible on screens <1024px, hover-reveal (group-hover pattern) on screens ≥1024px
- **FR-002**: System MUST show real-time peer status (online/offline) based on lastHandshake timestamp (online if within last 2-3 minutes)
- **FR-003**: System MUST provide a form to add new peers with fields for name (required) and allowed IPs (required, CIDR notation)
- **FR-004**: System MUST validate CIDR notation client-side before submitting to POST /peers endpoint (reject invalid formats like "10.0.0.5" without prefix)
- **FR-005**: System MUST call POST /peers API with name and allowedIPs, handle both successful creation (201) and error responses (400/500)
- **FR-006**: System MUST display peer configuration (WireGuard .conf format) and QR code after successful peer addition for easy client setup
- **FR-007**: System MUST provide delete/remove functionality for peers via DELETE /peers/{id} endpoint with confirmation dialog
- **FR-008**: System MUST fetch and display aggregate interface statistics from GET /stats endpoint (interface name, peer count, total RX/TX)
- **FR-008a**: System MUST display quick stats cards on main dashboard showing active peer count, total received bytes, and total sent bytes (numeric display only, no charts)
- **FR-009**: System MUST format data transfer statistics in human-readable units (bytes, KB, MB, GB, TB) rather than raw byte counts
- **FR-010**: System MUST use Svelte stores for state management (peers store, stats store) to share data across components and pages
- **FR-011**: System MUST implement API client utilities/functions to encapsulate backend communication (base URL configuration, error handling, response parsing)
- **FR-012**: System MUST display user-friendly error notifications for API failures (network errors, 400/500 responses) without exposing technical details
- **FR-013**: System MUST provide loading states (spinners, skeleton screens) during API requests to improve perceived performance
- **FR-014**: System MUST implement input validation utilities for common patterns (CIDR notation, required fields, string sanitization)
- **FR-015**: System MUST use TailwindCSS utility classes for styling and maintain consistent design system (spacing, colors, typography)
- **FR-016**: System MUST be responsive from mobile (320px) to desktop (1920px+) with breakpoints: mobile (<768px) uses bottom navigation and stacked cards, tablet (768px-1023px) uses simplified sidebar, desktop (1024px+) uses full sidebar with all features
- **FR-017**: System MUST include persistent sidebar navigation on left side with sections: Dashboard, Peers, Settings (future), Logs (future), and usage widget at bottom
- **FR-018**: System MUST handle empty states gracefully (no peers configured, no stats available) with clear messaging and calls-to-action
- **FR-019**: System MUST implement glassmorphism visual design with backdrop blur effects (blur(12-24px)), semi-transparent panels (rgba backgrounds with 0.03-0.4 alpha), and subtle borders (rgba(255,255,255,0.08-0.1))
- **FR-020**: System MUST use radial gradient backgrounds (e.g., radial-gradient(circle at top left, #1a2a3a, #101922)) and layered depth with varying panel opacities
- **FR-021**: System MUST implement Settings page UI (read-only/preview mode) displaying static interface configuration (listen port, MTU, addresses, server public key) with visual indicators that functional editing is coming soon

### Key Entities

- **Peer**: Represents a WireGuard client/device with attributes: id (public key), publicKey, name, endpoint, allowedIPs (array), lastHandshake (timestamp), receiveBytes, transmitBytes. Status (online/offline) is derived from lastHandshake freshness.
- **Stats**: Aggregate interface statistics with attributes: interfaceName, peerCount, totalRx, totalTx. Represents overall VPN network health.
- **PeerForm**: User input data for adding peers with attributes: name (string, required), allowedIPs (string array, CIDR notation, required). Used for form validation and submission.
- **Notification**: User feedback message with attributes: type (success/error/warning), message (string), duration (auto-dismiss timeout). Used for API response feedback.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Administrator can view all peers and their status within 3 seconds of navigating to /peers page (Time to Interactive)
- **SC-002**: Administrator can add a new peer in under 1 minute (including form fill, submission, and viewing QR code/config)
- **SC-003**: 95% of add peer operations complete successfully on first attempt (minimal validation errors, clear form labels)
- **SC-004**: All API error scenarios display user-friendly notifications without crashing the application
- **SC-005**: Peer list refreshes within 1 second after adding or removing a peer (optimistic updates or fast API response)
- **SC-006**: Statistics page loads and displays aggregate data within 2 seconds

### Performance Requirements

**Frontend (per Constitution Principle II & V)**:
- Time to Interactive (TTI): <3s on 3G connection
- First Contentful Paint (FCP): <1.5s on initial page load
- Bundle size impact: Total bundle size <200KB gzipped (SvelteKit + TailwindCSS + components)
- Lighthouse Performance score: ≥90 for both /peers and /stats pages
- Animation smoothness: 60fps for modal transitions, status badge updates, and list rendering

**Backend Integration**:
- API response handling: Display loading state for requests >500ms, timeout after 10s with error notification
- State updates: Svelte store updates propagate to UI within 100ms (reactive updates)

## Assumptions

1. **Backend API availability**: Backend server is running and accessible at a configured base URL (default: http://localhost:8080, configurable via environment variables)
2. **API contract stability**: Backend API endpoints (/peers, /stats) follow the documented schema in backend/API.md without breaking changes
3. **SvelteKit scaffolding**: Project is already initialized with SvelteKit, TailwindCSS, and basic routing structure as indicated by user
4. **No authentication**: No login, session management, or role-based access control required for initial release (all users are administrators)
5. **Responsive design priority**: Desktop-optimized experience with tablet and basic mobile support; mobile layouts simplified (bottom nav, stacked cards) without advanced touch gestures
6. **Modern browser support**: Target latest 2 versions of Chrome, Firefox, Safari, Edge (no IE11 support)
7. **CIDR notation knowledge**: Administrators understand basic networking concepts and CIDR notation (provide examples in form, but no in-depth tutorial)
8. **Single interface management**: Backend manages one WireGuard interface (wg0 or configured), frontend assumes single interface context
9. **Real-time updates**: Polling strategy acceptable for status updates (no WebSocket/SSE required initially); future enhancement for real-time stats
10. **QR code library**: Use existing QR code generation library (e.g., qrcode.js, svelte-qrcode) for peer config QR codes

## Testing Strategy

**Frontend Testing Approach (per Constitution Principle II)**:
- **Manual testing prioritized**: No automated unit/integration tests required for frontend components (align with Constitution: Frontend UX First)
- **Type safety via TypeScript**: Leverage TypeScript strict mode to catch type errors at compile time (primary quality gate)
- **Visual QA**: Manual testing of all user flows in development environment (add peer, remove peer, view stats, form validation)
- **Performance auditing**: Use Lighthouse and browser DevTools to verify performance budgets (TTI, FCP, bundle size)
- **API integration testing**: Manually test against real backend API to verify request/response handling and error scenarios

**Backend Testing (existing, no changes)**:
- Backend already has comprehensive test coverage per Constitution Principle I (TDD mandatory)
- Frontend assumes backend API contract is stable and tested

**Exploratory Testing Focus Areas**:
- Form validation edge cases (empty inputs, invalid CIDR, special characters in names)
- API error handling (network failures, 400/500 responses, timeout scenarios)
- UI responsiveness on different screen sizes (desktop, tablet)
- Browser compatibility (latest Chrome, Firefox, Safari, Edge)
- Performance under load (large peer lists, frequent refreshes)

## Dependencies

- **Backend API**: Requires backend server running with endpoints: GET /peers, POST /peers, DELETE /peers/{id}, GET /stats
- **SvelteKit**: Framework and routing already scaffolded (confirmed by user)
- **TailwindCSS**: Styling framework already integrated (confirmed by user)
- **TypeScript**: Type safety and interfaces (assumed present in SvelteKit scaffold)
- **QR Code Library**: To be added as npm dependency (e.g., qrcode, svelte-qrcode) for peer config QR codes
- **Date/Time Utilities**: For formatting lastHandshake timestamps (use built-in Intl.DateTimeFormat or date-fns library)
- **CORS Configuration**: Backend must allow frontend origin (localhost:5173 for dev, production domain for deployment)

## Out of Scope (Future Extensions)

The following features are explicitly excluded from this specification and may be considered for future releases:

1. **Authentication & Authorization**: No login, user management, or role-based access control in this release
2. **Real-time Updates**: No WebSocket/Server-Sent Events for live peer status updates; polling or manual refresh acceptable
3. **Theming & Dark Mode**: Single theme/color scheme for initial release; theme switcher is future enhancement (glassmorphism dark theme is default)
4. **Peer Editing**: No ability to edit peer details (name, allowed IPs) after creation; must delete and re-add
5. **Functional Interface Configuration**: Settings page displays static/mock data only; actual configuration changes (listen port, MTU, addresses) require backend API endpoints (deferred to future release)
6. **Service Control**: Interface toggle (start/stop) and restart service button are UI-only; functional implementation requires backend API (deferred)
7. **Usage Quota Tracking**: Sidebar usage widget displays mock data; actual quota tracking requires backend API (deferred)
8. **Traffic Visualization Charts**: Dashboard shows numeric stats cards only; SVG/Canvas line graphs for RX/TX trends deferred to future release
9. **Historical Stats Data**: Real-time stats only; no time-series history or trend analysis (requires backend API enhancement)
5. **Advanced Filtering & Search**: No search bar or filters for peer list; simple list display only
6. **Peer Groups/Organization**: No ability to group peers by tags, categories, or custom metadata
7. **Activity Logs/Audit Trail**: No logging of administrator actions (who added/removed peers, when)
8. **Email/Notification Integration**: No automatic notifications for peer status changes or system events
9. **Multi-Interface Support**: Single WireGuard interface management only; no support for managing multiple interfaces
10. **Backup/Export**: No ability to export peer list or configuration backups
11. **Internationalization (i18n)**: English language only; no multi-language support
12. **Accessibility (a11y) enhancements**: Basic accessibility compliance, but no WCAG 2.1 AA certification effort
13. **Peer Connection Diagnostics**: No ping tests, traceroute, or advanced troubleshooting tools
14. **Bandwidth Limiting**: No ability to set transfer rate limits per peer
15. **Mobile App**: Web-only interface; no native iOS/Android app

These items may be prioritized in subsequent specifications based on user feedback and business priorities.

## Constitution Alignment Check

This specification aligns with the WireGuard Manager Constitution (v1.0.0) as follows:

- **Principle I (Backend Testing Discipline)**: Backend API (already implemented) follows TDD; frontend consumes tested endpoints
- **Principle II (Frontend UX First)**: This spec prioritizes UX and performance over automated testing; manual testing strategy documented
- **Principle III (API Contract Stability)**: Relies on documented backend API endpoints without breaking changes
- **Principle IV (Configuration & Environment)**: Backend API base URL configurable via environment variables
- **Principle V (Performance Budgets)**: Explicit performance requirements defined (TTI <3s, FCP <1.5s, bundle <200KB, Lighthouse ≥90)
- **Principle VI (Observability)**: Frontend will log errors to browser console; backend already has structured logging

**No constitution violations identified.**
