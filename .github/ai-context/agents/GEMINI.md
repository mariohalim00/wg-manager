# Gemini-Specific Agent Rules

**Purpose**: Gemini is effective at rapid code generation and pattern matching. Use this to optimize output.

## Behavior Optimization for Gemini

### Rapid Pattern Recognition
- Identify recurring patterns and apply them consistently (e.g., middleware wrapping, handler structure)
- Use existing examples as templates for new code
- Generate boilerplate efficiently following established conventions

### Code Generation
- When generating Go code, follow the handler/service/middleware pattern established in backend/
- When generating Svelte code, use existing components as templates (PeerTable, PeerModal pattern)
- Include type hints and interfaces—don't skip TypeScript or Go signatures

### Focused Scope
- Work on well-defined, specific tasks (e.g., "add a new API endpoint" vs. "refactor logging")
- For complex changes, break into smaller subtasks and handle sequentially
- Ask for clarification on ambiguous requirements rather than guessing

## Gemini-Specific Task Types

**Code generation**: Creating new components, handlers, or services from patterns
**Boilerplate**: Setting up new files following project conventions
**Refactoring**: Applying consistent patterns across similar files
**Testing**: Generating test cases for Go handlers (TDD-driven)

## Key Constraints

- Always include proper error handling (Go: return `error`, HTTP handlers: log and respond with status)
- Maintain type safety in TypeScript and Go signatures
- Never generate code without understanding the data flow (see [architecture.md](../knowledge/architecture.md))

## When to Escalate to Claude or Humans

- Architectural decisions with trade-offs
- Major refactors affecting multiple subsystems
- Edge cases or error scenarios that need special handling
- Performance-critical code requiring benchmarking
