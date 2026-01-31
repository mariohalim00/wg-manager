# Claude-Specific Agent Rules

**Purpose**: Claude excels at deep reasoning and context synthesis. Use this to maximize quality.

## Behavior Optimization for Claude

### Leverage Extended Reasoning

- Use your superior reasoning capability to understand _why_ code is structured this way
- When proposing changes, explain the architectural reasoning behind them
- For complex refactors, outline the risk analysis and migration strategy
- Synthesize information across 5+ files to identify hidden dependencies

### Documentation and Explanation

- Provide detailed comments for non-obvious logic (WireGuard service interactions, store reactivity)
- Explain architectural decisions when introducing new patterns
- When suggesting changes, explain why existing code might be fragile

### Code Review Mindset

- Think like a reviewer: what would catch bugs or regressions?
- Identify edge cases and potential failure modes
- Suggest defensive programming patterns (especially for WireGuard operations)

## Claude-Specific Task Types

**Architecture analysis**: Ideal for understanding data flows and proposing major refactors
**Complex debugging**: Synthesize logs, code, and error messages to find root causes
**Documentation writing**: Use reasoning to create comprehensive guides and decision records
**Test strategy**: Design comprehensive test suites for backend (following TDD)

## When to Ask for Clarification

- Architecture decisions not documented in [decisions.md](../knowledge/decisions.md)
- Trade-offs between performance, complexity, and maintainability
- Design patterns that differ from common conventions
- Potential impact on other components

## Integration with SpecKit

Claude should read and respect the Constitution and SpecKit outputs:

- Check [constitution.md](../../../.specify/memory/constitution.md) before proposing changes
- For new features, follow spec-driven development workflow
- Synthesize feature specs with architectural knowledge
