<!--
Sync Impact Report:
Version change: none → 0.1.0
Modified principles:
  - Clean Code
  - Don't Repeat Yourself (DRY)
  - Frontend UI/UX Design Adherence
  - Frontend Technology Stack
  - Backend Technology Stack & Practices
  - Adherence to Language & Framework Best Practices
Added sections:
  - Core Principles
  - Project Conventions
  - Development Workflow
  - Governance
Removed sections: None
Templates requiring updates:
  - .specify/templates/plan-template.md (✅ reviewed, no changes needed)
  - .specify/templates/spec-template.md (✅ reviewed, no changes needed)
  - .specify/templates/tasks-template.md (✅ reviewed, no changes needed)
Follow-up TODOs: None
-->
# wg-manager Constitution

## Core Principles

### I. Clean Code
Code must be clear, readable, maintainable, and understandable by others. It should be self-documenting where possible.
Rationale: Reduces technical debt, improves collaboration, and simplifies debugging and and feature development.

### II. Don't Repeat Yourself (DRY)
Avoid duplication of knowledge or logic across the codebase. Implement logic in one, authoritative place.
Rationale: Ensures consistency, reduces bugs, and simplifies maintenance by requiring changes in only one location.

### III. Frontend UI/UX Design Adherence
The user interface and user experience must strictly follow the provided design specifications or existing design language established for the project.
Rationale: Ensures a consistent and high-quality user experience, aligns with brand guidelines, and reduces rework.

### IV. Frontend Technology Stack
The frontend application is built using SvelteKit as the framework and Tailwind CSS for all styling.
Rationale: Standardizes development, leverages efficient tooling, and ensures consistency in component architecture and styling.

### V. Backend Technology Stack & Practices
The backend application is developed using the latest stable version of Go. It must use the `net/http` package for web services and `slog` for JSON-formatted logging. Development must strictly adhere to Test-Driven Development (TDD).
Rationale: Ensures robustness, performance, maintainability, and early bug detection, while standardizing backend development.

### VI. Adherence to Language & Framework Best Practices
All code (frontend and backend) must follow the established best practices, idiomatic patterns, and style guides for their respective languages (TypeScript, Go), frameworks (SvelteKit), and libraries (Tailwind CSS, `net/http`, `slog`).
Rationale: Promotes high-quality, maintainable, and performant code, facilitating collaboration and onboarding.

## Project Conventions

*   **Version Control:** Git is used for version control. All changes must be made via pull requests (PRs).
*   **Dependency Management:** Frontend dependencies are managed with `npm`. Backend dependencies are managed with Go Modules.
*   **Environment Management:** `mise` is used to manage Node.js and Go versions.

## Development Workflow

*   **Feature Branches:** All new features or bug fixes must be developed on dedicated feature branches.
*   **Code Review:** All pull requests require at least one approving review before merging.
*   **Testing:** New features and bug fixes must include corresponding tests (unit, integration, end-to-end where applicable).
*   **Continuous Integration:** PRs must pass all CI checks (linting, testing, building) before merging.

## Governance

*   **Amendment Procedure:** Any proposed changes to this Constitution must be submitted as a pull request, discussed, and approved by the core development team.
*   **Versioning Policy:** The `CONSTITUTION_VERSION` follows Semantic Versioning (MAJOR.MINOR.PATCH).
    *   `MAJOR`: Backward incompatible governance/principle removals or redefinitions.
    *   `MINOR`: New principle/section added or materially expanded guidance.
    *   `PATCH`: Clarifications, wording, typo fixes, non-semantic refinements.
*   **Compliance Review:** Code reviews and CI checks will enforce adherence to the principles outlined in this Constitution.

**Version**: 0.1.0 | **Ratified**: 2026-01-31 | **Last Amended**: 2026-01-31