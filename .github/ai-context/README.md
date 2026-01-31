# AI Context System

**Version**: 1.0.0  
**Last Updated**: 2026-01-31

## Overview

This directory contains a **unified AI context system** designed to help all AI coding agents (GitHub Copilot, Claude, Gemini, Cursor, etc.) work effectively and safely in this repository.

**Purpose**: Provide AI agents with comprehensive project knowledge, conventions, and guardrails to generate high-quality code that aligns with project principles.

## Directory Structure

```
.github/ai-context/
│
├─ agents/                      ← AI agent-specific rules
│   ├─ AGENTS.md                ← Universal rules (all AI tools)
│   ├─ CLAUDE.md                ← Claude-specific optimizations
│   ├─ GEMINI.md                ← Gemini-specific patterns
│   ├─ docs.agent.md            ← Documentation-focused rules
│   ├─ security.agent.md        ← Security-focused rules
│   └─ test.agent.md            ← Testing-focused rules
│
├─ skills/                      ← Domain/tool-specific knowledge
│   ├─ go.skills.md             ← Go backend conventions
│   └─ svelte.skills.md         ← Svelte frontend conventions
│
├─ knowledge/                   ← Persistent project knowledge
│   ├─ architecture.md          ← System architecture & data flows
│   ├─ decisions.md             ← Architectural decision records (ADRs)
│   ├─ domains.md               ← Domain language & WireGuard concepts
│   ├─ memory.md                ← Project history & evolution
│   └─ chatmem.md               ← Session context & active work
│
└─ README.md                    ← This file
```

## Quick Start for AI Agents

### 1. Read Universal Rules
Start with [`agents/AGENTS.md`](agents/AGENTS.md) for baseline behavior expectations.

### 2. Check Tool-Specific Rules
- **Claude**: Read [`agents/CLAUDE.md`](agents/CLAUDE.md)
- **Gemini**: Read [`agents/GEMINI.md`](agents/GEMINI.md)
- **Other tools**: Use AGENTS.md as fallback

### 3. Understand the Stack
- **Backend**: Read [`skills/go.skills.md`](skills/go.skills.md)
- **Frontend**: Read [`skills/svelte.skills.md`](skills/svelte.skills.md)

### 4. Learn the Architecture
Read [`knowledge/architecture.md`](knowledge/architecture.md) for system overview, component boundaries, and data flows.

### 5. Understand Domain Terms
Read [`knowledge/domains.md`](knowledge/domains.md) for WireGuard concepts and terminology.

### 6. Check Active Work
Read [`knowledge/chatmem.md`](knowledge/chatmem.md) for current development status and common pitfalls.

## File Purposes

### `agents/` — Behavior Rules

**Purpose**: Define HOW AI agents should behave when generating code.

**AGENTS.md** — Universal rules for all AI tools:
- Decision-making workflow (analyze → propose → implement → verify)
- Safety rules (never violate)
- Cross-stack awareness
- Task execution patterns

**CLAUDE.md** — Claude-specific optimizations:
- Leverage extended reasoning
- Deep architectural analysis
- Complex debugging strategies

**GEMINI.md** — Gemini-specific patterns:
- Rapid pattern recognition
- Code generation from templates
- Focused scope

**docs.agent.md** — Documentation rules:
- Documentation standards
- Markdown formatting
- Mermaid diagrams
- API documentation maintenance

**security.agent.md** — Security rules:
- Input validation requirements
- Secret management
- Logging security
- Threat model

**test.agent.md** — Testing rules:
- TDD workflow (backend mandatory)
- Frontend testing approach (manual)
- Test patterns and anti-patterns

### `skills/` — Technical Facts

**Purpose**: Provide domain/tool-specific knowledge that AI agents should treat as facts.

**go.skills.md** — Go backend knowledge:
- Project structure
- Handler/service/middleware patterns
- Configuration (Twelve-Factor)
- Testing conventions (TDD)
- WireGuard library usage
- Error handling patterns

**svelte.skills.md** — Svelte frontend knowledge:
- SvelteKit routing (file-based)
- Component patterns
- State management (stores)
- TypeScript conventions
- TailwindCSS + DaisyUI styling
- Performance budgets

### `knowledge/` — Context & History

**Purpose**: Provide persistent knowledge about the system, domain, and project evolution.

**architecture.md** — System architecture:
- High-level overview
- Component responsibilities
- Data flows
- Integration points
- Deployment architecture
- Mermaid diagrams

**decisions.md** — Architectural decision records:
- Major design decisions
- Trade-offs considered
- Rationale for choices
- Impact on codebase
- Decision template for future use

**domains.md** — Domain language:
- WireGuard concepts (Peer, PublicKey, AllowedIPs, Handshake)
- System concepts (Real vs. Mock service, Handler vs. Service)
- Common confusion points
- Glossary

**memory.md** — Project history:
- Timeline of major phases
- Refactors and migrations
- Abandoned approaches (and why)
- Long-term constraints
- Key learnings
- Known technical debt

**chatmem.md** — Session context:
- Active development areas
- Current branch status
- Critical invariants
- Common pitfalls & solutions
- Recent discoveries
- Session recovery checklist

## Usage by AI Tool

### GitHub Copilot

**Entry point**: [`../copilot-instructions.md`](../copilot-instructions.md)

**Workflow**:
1. Copilot reads `copilot-instructions.md` (global instructions)
2. References `agents/AGENTS.md` for behavior rules
3. Consults `skills/*.skills.md` for technical patterns
4. Uses `knowledge/*` for architectural context

### Claude

**Entry point**: [`agents/CLAUDE.md`](agents/CLAUDE.md)

**Workflow**:
1. Read `AGENTS.md` (universal rules)
2. Apply `CLAUDE.md` (Claude-specific optimizations)
3. Use `skills/` and `knowledge/` as needed
4. Leverage extended reasoning for complex tasks

### Gemini

**Entry point**: [`agents/GEMINI.md`](agents/GEMINI.md)

**Workflow**:
1. Read `AGENTS.md` (universal rules)
2. Apply `GEMINI.md` (Gemini-specific patterns)
3. Use `skills/` for rapid pattern matching
4. Consult `knowledge/` for context

### Other Tools (Cursor, Windsurf, Cline, etc.)

**Fallback**: Use [`agents/AGENTS.md`](agents/AGENTS.md) as universal rules.

## Key Principles

### 1. Consolidation (No Redundancy)

- **One source of truth** for each topic
- Cross-reference with links (don't duplicate)
- Example: API contract lives in `backend/API.md`, not repeated elsewhere

### 2. Purpose-Driven Organization

- **agents/**: HOW to behave
- **skills/**: WHAT to know (technical facts)
- **knowledge/**: WHY things are this way (context, history)

### 3. Scannable & Concise

- Bullet points over paragraphs
- Code examples for clarity
- Mermaid diagrams for architecture
- Tables for comparisons

### 4. Tool-Agnostic (Mostly)

- `AGENTS.md` applies to all AI tools
- Tool-specific files (`CLAUDE.md`, `GEMINI.md`) enhance, don't replace
- Frontend/backend skills are tool-agnostic

## Integration with Project

### Relationship to Constitution

**Constitution** ([`../../.specify/memory/constitution.md`](../../.specify/memory/constitution.md)) defines **non-negotiable principles**.

**AI Context** provides **implementation guidance** for those principles.

**Example**:
- Constitution Principle I: "Backend TDD mandatory"
- AI Context `test.agent.md`: HOW to do TDD (patterns, workflow, examples)

### Relationship to SpecKit

**SpecKit** ([`../../.specify/`](../../.specify/)) provides **spec-driven development workflow**.

**AI Context** helps agents **implement specs correctly** by providing:
- Technical knowledge (Go/Svelte patterns)
- Architectural context (where new code fits)
- Historical context (why things are structured this way)

**Workflow**:
1. Spec created (via SpecKit)
2. AI agent reads spec + AI context
3. AI agent generates code following patterns and principles
4. Code reviewed against constitution and spec

## Maintenance

### When to Update

**agents/** — Update when:
- AI agent behavior needs adjustment
- New safety rules identified
- Tool-specific optimizations discovered

**skills/** — Update when:
- New technology added to stack
- Conventions change
- New patterns emerge

**knowledge/** — Update when:
- Architecture changes significantly
- Major decisions made
- New domain concepts introduced
- Active work shifts (update chatmem.md frequently)

### How to Update

1. **Identify file** to update based on content type
2. **Edit file** with new information
3. **Cross-reference** related files (update links if needed)
4. **Verify consistency** with constitution and other docs
5. **Update "Last Updated" date** at top of file

### Consistency Checks

Before committing:
- [ ] No contradictions with constitution
- [ ] Links are not broken
- [ ] Code examples are valid
- [ ] Mermaid diagrams render correctly
- [ ] No redundant information across files

## Benefits of This System

### For AI Agents

- **Immediate productivity**: Agents understand project from first interaction
- **Reduced errors**: Safety rules and patterns prevent common mistakes
- **Context awareness**: Historical decisions explain current structure
- **Consistent output**: All agents follow same conventions

### For Developers

- **Onboarding**: New developers read same docs as AI agents
- **Knowledge preservation**: Project knowledge not lost with team changes
- **Quality assurance**: AI-generated code follows established patterns
- **Debugging**: AI agents understand architecture when helping debug

### For Project

- **Code consistency**: All code (human + AI) follows same conventions
- **Faster development**: AI agents accelerate feature implementation
- **Lower maintenance**: Well-documented decisions prevent rework
- **Easier evolution**: Clear architecture enables safe refactoring

## Examples

### Example 1: Adding a New Backend Endpoint

**Agent reads**:
1. `agents/AGENTS.md` → Universal behavior (propose first, verify integration)
2. `skills/go.skills.md` → Handler pattern, TDD workflow
3. `knowledge/architecture.md` → Where new endpoint fits in system
4. `agents/test.agent.md` → TDD workflow details

**Agent generates**:
1. Test file (TDD red phase)
2. Handler implementation (TDD green phase)
3. Route registration in `main.go`
4. Updates to `backend/API.md`

**Result**: New endpoint follows conventions, has tests, is documented.

### Example 2: Creating a New Svelte Component

**Agent reads**:
1. `agents/AGENTS.md` → Universal behavior
2. `skills/svelte.skills.md` → Component pattern, styling conventions
3. `knowledge/architecture.md` → How frontend integrates with backend
4. `knowledge/chatmem.md` → Current active work (avoid conflicts)

**Agent generates**:
1. Component file with TypeScript props
2. TailwindCSS + DaisyUI styling
3. Integration with existing stores

**Result**: Component follows conventions, uses correct styling, integrates cleanly.

## Evolution

This AI context system will evolve as the project grows:

**Version 1.0.0** (2026-01-31):
- Initial structure
- Backend (Go) and Frontend (Svelte) skills
- Architecture, decisions, domain, memory, chatmem knowledge
- Agent rules for Copilot, Claude, Gemini
- Specialized agents: docs, security, test

**Future additions** (as needed):
- More specialized agent files (performance.agent.md, refactor.agent.md)
- Additional skills files (docker.skills.md, kubernetes.skills.md)
- Expanded knowledge (more ADRs, detailed data flows)
- Tool-specific optimizations

## Contributing to AI Context

When adding to this system:

1. **Identify correct location**: agents/ vs. skills/ vs. knowledge/
2. **Follow existing patterns**: Use same formatting, structure
3. **Be concise**: Bullet points, code examples, diagrams
4. **Cross-reference**: Link to related docs
5. **Test with AI agent**: Verify agent uses new context correctly
6. **Update README**: If structure changes significantly

## Questions?

**For agents**: If unclear, ask for clarification rather than guessing.  
**For developers**: Review this README and explore individual files.  
**For maintainers**: Follow maintenance guidelines above.

---

**Remember**: This system helps AI agents work WITH the project, not AGAINST it. Keep it updated, concise, and useful.
