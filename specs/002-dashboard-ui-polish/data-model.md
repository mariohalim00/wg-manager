# Data Model: Dashboard UI Polish

This feature introduces no new backend data. It refines presentation of existing data with placeholders for missing fields.

## Entities (UI-only)

### DashboardView
- **Purpose**: Composition of status cards, traffic summaries, and peers table.
- **Fields**: N/A (derived from stores and static UI content).

### StatusCard
- **Fields**:
  - `label` (string)
  - `value` (string)
  - `indicator` (optional: status dot or icon)
  - `action` (optional: copy affordance; visual-only)
- **Validation**:
  - If data missing, `value` uses placeholder "—" or "Not available".

### TrafficSummary
- **Fields**:
  - `title` (string)
  - `value` (string)
  - `delta` (string, optional)
  - `chart` (decorative SVG or background)
- **Validation**:
  - Chart is static but structured for future data binding.

### PeerRow
- **Fields** (from existing `Peer` type):
  - `status` (online/offline)
  - `name` (string)
  - `allowedIps` (string[])
  - `transfer.received`, `transfer.sent` (numbers)
  - `latestHandshake` (string)
- **Validation**:
  - Render only real peers; show empty-state UI when list is empty.

## Relationships

- `DashboardView` composes multiple `StatusCard` items, two `TrafficSummary` cards, and a list of `PeerRow` entries.
