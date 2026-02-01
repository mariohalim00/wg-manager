# Quickstart: Dashboard UI Polish

## Prerequisites

- Node.js 20+
- Backend running (optional for UI-only layout work)

## Run frontend

```bash
npm install
npm run dev
```

Open http://localhost:5173 and navigate to the dashboard.

## Before Checklist

- Header layout and controls are visible.
- Status card grid is present.
- Traffic cards and charts are present.
- Peers table is present.

## Manual Validation Checklist

- Sidebar, header, status cards, charts, and peers table align to design spec.
- Dark mode is default; light mode remains consistent.
- Mobile: sidebar becomes drawer; cards and charts stack vertically.
- Tablet: 2-column card grid; peers table horizontally scrolls if needed.
- Placeholder values shown for missing public key, port, and subnet.
- Peers table shows empty-state when no peers exist.
