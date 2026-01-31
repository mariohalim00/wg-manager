# Universal AI Agent Rules

**Applies to**: All AI agents (Copilot, Claude, Gemini, Cursor, etc.)

## Core Agent Behavior

### Decision-Making
- **Analyze first**: Read relevant code and documentation before suggesting changes
- **Propose before acting**: For significant changes, outline the approach and wait for confirmation
- **Verify assumptions**: Ask for clarification when ambiguity exists
- **Consider impact**: Think about side effects across the full stack (frontend + backend)

### Code Quality
- **Type safety**: TypeScript frontend, Go backend—leverage type systems fully
- **Follow constitution**: Respect [Constitution.md](../../../.specify/memory/constitution.md)
  - **Backend**: TDD mandatory (write tests first for all handler/service logic)
  - **Frontend**: UX/performance focused (no tests required, but optimize for speed and usability)
- **Documentation**: Update comments and docs alongside code changes
- **Performance**: Frontend performance budgets are non-negotiable (TTI <3s, bundle <200KB)

### Cross-Stack Awareness
- **API contracts matter**: Changes to `backend/internal/handlers` or API endpoints must update [API.md](../../../backend/API.md)
- **Frontend integration**: Always verify API calls in frontend stores match backend response schemas
- **Shared data model**: Keep [types.ts](../../../src/lib/types.ts) and backend Peer struct in sync
- **CORS configuration**: Backend CORS settings must match frontend origin

## Task Execution Workflow

1. **Understand context**: Read relevant files, understand existing patterns
2. **Plan**: Outline changes, identify affected files
3. **Implement**: Make changes following project conventions
4. **Verify**: Ensure code builds, follows conventions, aligns with constitution
5. **Document**: Update comments, types, README if needed

## Safety Rules (Non-Negotiable)

- **Never modify tests without understanding TDD workflow** (backend)
- **Never commit untested Go code** (handlers, services, middleware)
- **Never remove or modify types without checking all usages**
- **Never hardcode sensitive data** (use environment variables, see [config.go](../../../backend/internal/config/config.go))
- **Never break API contracts** without major version bump and migration plan
- **Never commit performance regressions** without justification (see Constitution V)

## Tool-Specific Details

See agent-specific rules:
- [CLAUDE.md](CLAUDE.md) — Claude-specific optimizations
- [GEMINI.md](GEMINI.md) — Gemini-specific patterns
- Copilot uses `copilot-instructions.md` globally
  
---

You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

## Available MCP Tools:

### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.