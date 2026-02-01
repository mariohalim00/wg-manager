# Research: Dashboard UI Polish

**Date**: 2026-02-01  
**Scope**: Frontend-only dashboard UI refinement

## Decisions

### Header controls behavior
- **Decision**: Keep search, notifications, and user menu as visual-only placeholders.
- **Rationale**: Avoid introducing new behaviors per constraints; match design without backend or state changes.
- **Alternatives considered**: Implement search filtering (rejected: adds behavior and state).

### Missing fields (public key, listening port, subnet)
- **Decision**: Show placeholders ("—" or "Not available") and keep cards visible.
- **Rationale**: Preserves layout and hierarchy while avoiding API changes.
- **Alternatives considered**: Hide cards; use mock values (rejected: misleading or disruptive layout).

### Traffic charts
- **Decision**: Use static decorative charts structured for future data binding.
- **Rationale**: Visual polish without requiring timeseries data; keeps implementation ready for backend data later.
- **Alternatives considered**: Generate mock points; hide charts (rejected: less faithful to design).

### Peers table data
- **Decision**: Use real peer data only; show empty-state styling when list is empty.
- **Rationale**: Avoids mock data and respects constraint of no business logic changes.
- **Alternatives considered**: Mock rows; hide table (rejected: misleading or less useful).
