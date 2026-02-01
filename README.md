# WireGuard Manager

A full-stack WireGuard peer management system with a SvelteKit dashboard and a Go backend. Manage peers, view stats, and generate client configs with a fast, modern UI. Suitable for server planning to run native wireguard but wants an easy way to create and manage profiles.

## Features

- Peer management: list, add, remove
- Real-time statistics (handshakes, RX/TX totals)
- QR/config display for client setup
- Responsive dashboard UI 
- Structured logging and robust API layer

## Architecture

Frontend (SvelteKit) → Backend API (Go/net\http) → WireGuard kernel

Metadata is stored in JSON, while real-time stats are read from the WireGuard interface. The backend uses a mock service when WireGuard is unavailable for local development.

## Tech Stack

- Frontend: SvelteKit 2.x, Svelte 5.x, TypeScript 5.x
- Backend: Go 1.25.6, net\http, slog
- Styling: TailwindCSS 4.x + DaisyUI
- WireGuard: wgctrl

## Requirements

- Node.js 20+ (frontend)
- Go 1.25.6+ (backend)
- WireGuard kernel module (for real service)

## Getting Started

### Backend

```bash
cd backend
go test ./...
go run ./cmd/server (run with sudo)
```

The backend runs on :8080 by default.

### Frontend

```bash
npm install
npm run dev
```

The frontend runs on http://localhost:5173.

## Configuration

Backend configuration uses defaults in config.json with environment overrides.

Common environment variables:

- WG_SERVER_PORT (default :8080)
- WG_INTERFACE_NAME (default wg0)
- WG_STORAGE_PATH (default ./data/peers.json)
- WG_SERVER_ENDPOINT (public endpoint for clients)
- WG_SERVER_PUBKEY (server public key)
- CORS_ALLOWED_ORIGINS (comma-separated)

## API Endpoints

- GET /peers
- POST /peers
- DELETE /peers/{publicKey}
- GET /stats

See backend/API.md for full request and response schemas.

## Scripts

Frontend:

- npm run dev
- npm run build
- npm run preview
- npm run lint
- npm run check

Backend:

- go test ./...
- sudo go run ./cmd/server (run with elevated permissions, you will need it)

## Deployment

The frontend builds as a static SPA via adapter-static. Serve the build output with any static host and proxy API requests to the Go backend.

Example flow:

1. npm run build
2. Serve build/ with Nginx/Caddy
3. Run backend on :8080 (or configured port)

## License

MIT
