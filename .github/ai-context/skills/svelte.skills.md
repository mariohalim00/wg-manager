# Svelte/SvelteKit Frontend Skills & Conventions

**Location**: `src/` directory

## Framework & Tooling Stack

- **SvelteKit**: 2.x (full-stack framework with routing)
- **Svelte**: 5.x (reactive component framework)
- **TypeScript**: 5.x (strict type checking required)
- **CSS**: TailwindCSS 4.x + DaisyUI components
- **Build**: Vite 7.x
- **Code Quality**: ESLint + Prettier (with Svelte and Tailwind plugins)
- **Type Checking**: `svelte-check`

## Project Structure Conventions

```
src/
├── routes/                      # File-based routing (SvelteKit)
│   ├── +layout.svelte           # Root layout (nav, global structure)
│   ├── +page.svelte             # Dashboard (home page)
│   ├── peers/
│   │   └── +page.svelte         # Peer management page
│   └── stats/
│       └── +page.svelte         # Statistics page
├── lib/
│   ├── components/              # Reusable Svelte components
│   │   ├── PeerTable.svelte     # Renders peer list in table
│   │   ├── PeerModal.svelte     # Add/edit peer modal
│   │   ├── QRCodeDisplay.svelte # QR code rendering
│   │   ├── StatusBadge.svelte   # Online/offline indicator
│   │   └── Notification.svelte  # Toast notifications
│   ├── stores/                  # Reactive state management (Svelte stores)
│   │   ├── peers.ts             # Peer list state + API integration
│   │   └── stats.ts             # Statistics state + API integration
│   ├── assets/                  # Static assets (images, fonts)
│   ├── types.ts                 # TypeScript interfaces (shared with backend)
│   └── index.ts                 # Barrel export for convenient imports
├── app.html                     # HTML template (Vite entry point)
├── app.css                      # Global styles (Tailwind directives)
└── app.d.ts                     # SvelteKit type definitions
```

## Routing Conventions

SvelteKit uses file-based routing:
- `src/routes/+page.svelte` → `/` (home/dashboard)
- `src/routes/peers/+page.svelte` → `/peers`
- `src/routes/stats/+page.svelte` → `/stats`
- `+layout.svelte` applies to the directory and all child routes

**Navigation** in Svelte:
```svelte
<a href="/peers">Go to Peers</a>
<!-- or programmatically: -->
<script>
    import { goto } from '$app/navigation';
    goto('/peers');
</script>
```

## Component Pattern

**Standard Svelte component structure**:

```svelte
<script lang="ts">
    import type { Peer } from '$lib/types';
    
    // Props (inputs)
    export let peer: Peer;
    export let onDelete: (id: string) => void = () => {};
    
    // Component logic
    let isLoading = false;
    const handleClick = async () => {
        isLoading = true;
        try {
            // Async logic
        } finally {
            isLoading = false;
        }
    };
</script>

<div class="card">
    {#if isLoading}
        <div class="loading">Loading...</div>
    {:else}
        <p>{peer.name}</p>
    {/if}
</div>

<style>
    /* Optional scoped styles (CSS Module-like) */
    .card {
        /* ... */
    }
</style>
```

**Key conventions**:
- Always use `lang="ts"` in script blocks
- Props exported at top of script
- Export handler props with default no-ops (`() => {}`)
- Use reactive statements with `$:` for derived state
- Use `#if`/`{:else}`/`{/if}` for conditionals
- Use `#each` with unique key for lists

## State Management (Svelte Stores)

**Writable stores** (mutable state):
```typescript
// stores/peers.ts
import { writable } from 'svelte/store';
import type { Peer } from '$lib/types';

export const peers = writable<Peer[]>([]);

// Usage in components:
<script>
    import { peers } from '$lib/stores/peers';
    
    $peers.forEach(p => console.log(p.name)); // $ prefix auto-subscribes
    peers.set([...]); // Update store
    peers.update(p => [...p, newPeer]); // Update with function
</script>
```

**Derived stores** (computed from other stores):
```typescript
import { derived } from 'svelte/store';

export const onlinePeers = derived(peers, p => p.filter(x => x.status === 'online'));
```

**Store conventions** in this project:
- `peers.ts` — List of all peers, API integration (load, add, remove)
- `stats.ts` — Overall interface statistics, API integration
- Stores handle API calls, error state, loading state
- Components subscribe with `$store` syntax or use `subscribe()` method

## API Integration

**Backend API base**: `http://localhost:8080` (or via `+page.server.ts` server-side proxy)

**Endpoints** (see `backend/API.md`):
- `GET /peers` → List all peers
- `POST /peers` → Add peer
- `DELETE /peers/{id}` → Remove peer
- `GET /stats` → Get statistics

**Integration in stores**:
```typescript
// peers.ts example
async function loadPeers() {
    const response = await fetch('/api/peers'); // or 'http://localhost:8080/peers'
    if (!response.ok) throw new Error('Failed to load peers');
    const data = await response.json();
    peers.set(data);
}

export async function addPeer(name: string, publicKey: string, allowedIPs: string[]) {
    const response = await fetch('/api/peers', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, publicKey, allowedIPs })
    });
    if (!response.ok) throw new Error('Failed to add peer');
    const newPeer = await response.json();
    peers.update(p => [...p, newPeer]);
    return newPeer;
}
```

## TypeScript Conventions

**Strict mode** (tsconfig.json enforces strict checks):
- No `any` types without justification
- All function parameters and return types typed
- Components use `lang="ts"` in script blocks

**Type sharing** with backend:
- `src/lib/types.ts` defines `Peer`, `Stats` interfaces
- Keep these in sync with backend Go structs
- Use camelCase in TypeScript (e.g., `allowedIPs`, `latestHandshake`)

**Example**:
```typescript
// src/lib/types.ts
export interface Peer {
    id: string;
    name: string;
    publicKey: string;
    status: 'online' | 'offline';
    allowedIps: string[];
    latestHandshake: string;
    transfer: {
        received: number;
        sent: number;
    };
}
```

## Styling Conventions

**TailwindCSS + DaisyUI**:
- Use Tailwind utility classes for layout, spacing, colors
- Use DaisyUI components for buttons, modals, tables, forms
- Never write custom CSS except for scoped component styles (very rare)
- Tailwind plugin `prettier-plugin-tailwindcss` auto-sorts class names

**Example**:
```svelte
<div class="card bg-base-100 shadow-md">
    <div class="card-body">
        <h2 class="card-title">Title</h2>
        <p class="text-sm text-base-content/70">Description</p>
        <div class="card-actions justify-end">
            <button class="btn btn-primary">Action</button>
        </div>
    </div>
</div>
```

**DaisyUI commonly used components**:
- `button`, `btn` — Buttons
- `modal` — Modal dialogs
- `table` — Tables
- `form-control`, `input`, `select` — Form inputs
- `alert` — Alert boxes
- `badge` — Status badges
- `loading` — Loading spinners
- `navbar` — Navigation bar

## Code Quality Rules

### Type Safety
- Always type component props: `export let data: Peer`
- Always type function returns: `function doSomething(): Promise<Peer>`
- Use generics where appropriate

### ESLint + Prettier
```bash
npm run lint          # Check for issues
npm run format        # Auto-format code
```

**Key rules**:
- Semicolons required
- Single quotes for strings
- Indentation: 2 spaces
- Max line length: typically 100-120 chars (Prettier handles)

### Performance Considerations (Constitution II)

**Performance budgets** (strict targets):
- **TTI (Time to Interactive)**: < 3 seconds on 3G
- **FCP (First Contentful Paint)**: < 1.5 seconds
- **Bundle size**: < 200KB gzipped
- **Lighthouse score**: ≥ 90

**Optimization tips**:
- Lazy-load routes where possible
- Minimize JavaScript bundle (avoid large dependencies)
- Optimize images (use WebP, responsive sizes)
- Use Svelte's reactivity efficiently (avoid unnecessary re-renders)
- Cache API responses where appropriate

## No Frontend Tests (Constitution II)

**Frontend development is UX/performance-focused, not test-driven.**
- No jest/vitest/playwright required for frontend
- Manual testing and user feedback drive quality
- Code quality relies on TypeScript, ESLint, and careful review
- Performance profiling takes precedence over unit tests

## Common Tasks

### Adding a new page:
1. Create `src/routes/[name]/+page.svelte`
2. Import components and stores
3. Use Tailwind + DaisyUI for styling
4. Wire up API calls via stores
5. Test manually in dev (`npm run dev`)

### Adding a new component:
1. Create `src/lib/components/[Name].svelte`
2. Define props with types
3. Define event handlers
4. Use TailwindCSS + DaisyUI
5. Import and use in pages

### Integrating new API endpoint:
1. Add to `src/lib/types.ts` if new data type
2. Create store function in `src/lib/stores/` for API integration
3. Use store in component with `$store` syntax
4. Handle loading/error states in UI

### Debugging performance:
1. Run `npm run build` and check bundle size
2. Use DevTools Lighthouse audit
3. Check Core Web Vitals in DevTools Performance tab
4. Identify large imports or render bottlenecks
5. Optimize or defer non-critical code
