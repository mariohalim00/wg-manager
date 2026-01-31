# Project: wg-manager (SvelteKit Frontend)

## Project Overview

This is a SvelteKit project named "wg-manager", designed as a frontend application for managing WireGuard configurations. It leverages TypeScript for type safety, Tailwind CSS for styling, and DaisyUI for pre-built UI components. The project is structured with dedicated routes for a dashboard, peer management, and statistics, along with reusable Svelte components and reactive Svelte stores for state management.

## Key Technologies

- **Framework:** SvelteKit
- **Language:** TypeScript
- **Styling:** Tailwind CSS, DaisyUI
- **Build Tool:** Vite
- **Code Quality:** ESLint, Prettier
- **Node.js Version Management:** Mise (configured for Node.js LTS)

## Project Structure

The core application logic and UI are located in the `src/` directory:

- **`src/routes/`**: Contains the main application pages:
  - `+page.svelte`: The main dashboard page, displaying overall statistics.
  - `peers/+page.svelte`: Page for listing and managing WireGuard peers, including an "Add Peer" modal.
  - `stats/+page.svelte`: Page dedicated to displaying detailed statistics.
- **`src/lib/components/`**: Reusable Svelte components:
  - `PeerTable.svelte`: Displays a table of WireGuard peers.
  - `PeerModal.svelte`: A modal for adding or editing peer information.
  - `QRCodeDisplay.svelte`: Component to display QR codes (e.g., for peer configuration).
  - `StatusBadge.svelte`: A visual indicator for peer status (online/offline).
  - `Notification.svelte`: Displays simple notifications (success/error).
- **`src/lib/stores/`**: Svelte stores for reactive state management:
  - `peers.ts`: Manages the list of WireGuard peers.
  - `stats.ts`: Manages application-wide statistics.
- **`src/lib/types.ts`**: TypeScript interfaces and types, such as the `Peer` interface.
- **`src/app.html`**: The main HTML template for the SvelteKit application, including a responsive navigation bar using DaisyUI.
- **`src/app.css`**: Global CSS file, including Tailwind CSS base, components, and utilities directives.

## Building and Running

This project uses `npm` for package management and script execution.

- **Install Dependencies:**
  ```bash
  npm install
  ```
- **Run Development Server:**
  Starts the development server with hot-reloading.
  ```bash
  npm run dev
  ```
  To open the app in a new browser tab automatically:
  ```bash
  npm run dev -- --open
  ```
- **Build for Production:**
  Creates an optimized production build of the application.
  ```bash
  npm run build
  ```
- **Preview Production Build:**
  Serves the production build locally for testing.
  ```bash
  npm run preview
  ```
- **Check Code for Errors:**
  Runs Svelte-check and TypeScript checks.
  ```bash
  npm run check
  ```
- **Lint Code:**
  Checks code for style and potential errors using ESLint.
  ```bash
  npm run lint
  ```
- **Format Code:**
  Formats code using Prettier.
  ```bash
  npm run format
  ```

## Development Conventions

- **TypeScript:** All new Svelte components and JavaScript files should use TypeScript for improved type safety.
- **ESLint:** Adhere to the ESLint rules defined in `eslint.config.js` for consistent code quality.
- **Prettier:** Code formatting is enforced using Prettier. Run `npm run format` to automatically format your code. `prettier-plugin-tailwindcss` is used for Tailwind CSS class sorting.
- **Styling:** Utilize Tailwind CSS classes for styling, complemented by DaisyUI components.
- **Mise:** The Node.js version for development is managed by `mise`, as defined in `.mise.toml`.

---

# Project: wg-manager (Go Backend)

## Project Overview

This is a minimal Go backend application (`wg-manager/backend`) designed to serve as an API for the WireGuard manager frontend. It uses the standard `net/http` library for web serving, `slog` for structured logging, and follows Go best practices for a clean and extendable project structure.

## Key Technologies

- **Language:** Go
- **Web Framework:** `net/http` (standard library)
- **Logging:** `slog`
- **Configuration:** Custom JSON-based configuration loading
- **Node.js Version Management:** Mise (configured for Go "latest")

## Project Structure

The Go backend code resides in the `backend/` directory:

- **`backend/cmd/server/main.go`**: Application entry point, handles configuration loading (JSON + .env), dependency injection, modern route registration, and graceful shutdown.
- **`backend/internal/handlers/`**: Contains modular HTTP handlers for Peers and Stats, including robust input validation.
- **`backend/internal/config/`**: Manages configuration via JSON and environment variable overrides (Twelve-Factor App).
- **`backend/internal/middleware/`**: Implements logging and CORS with dynamic origin support.
- **`backend/internal/wireguard/`**: Interface with `wgctrl`, persistent storage for metadata, and key/config generation.

## API Endpoints

- **`GET /peers`**: Returns a JSON array of WireGuard peers with real-time stats and persistent names.
- **`POST /peers`**: Adds/configures a peer. Supports auto-key generation and returns client config.
- **`DELETE /peers/{id}`**: Removes a peer and its metadata using path parameters.
- **`GET /stats`**: Returns aggregate interface statistics.

## Building and Running

- **Configuration**: Use `backend/internal/config/config.json` for defaults and `.env` for environment-specific overrides (see `.env.example`).
- **Run Server**:
  ```bash
  cd backend
  go run ./cmd/server/main.go
  ```
- **Testing**:
  ```bash
  cd backend/cmd/server
  go test .
  ```

## Development Conventions

- **Graceful Shutdown**: Always close the `wireguard.Service` to release system resources.
- **Configuration**: Prioritize environment variables for production secrets and endpoint settings.
- **Validation**: Validate all user input at the handler level before processing.
