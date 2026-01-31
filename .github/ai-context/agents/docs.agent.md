# Documentation-Focused Agent Rules

**Purpose**: Guidelines for AI agents working on documentation tasks.

## Documentation Philosophy

Documentation in this project serves three audiences:

1. **AI agents** — Context for code generation and understanding
2. **Developers** — Setup, troubleshooting, and contribution guides
3. **Users** — Feature guides and API references

## Documentation Types

### 1. AI Context Documentation (`.github/ai-context/`)

**Purpose**: Help AI understand project architecture, decisions, conventions

**Style**:

- Structured Markdown with clear sections
- Bullet points for scanability
- Code examples for patterns
- Mermaid diagrams for architecture
- Decision rationale (why, not just what)

**Maintenance**:

- Update when architecture changes
- Keep in sync with constitution
- Add to `chatmem.md` when new patterns emerge

### 2. API Documentation (`backend/API.md`)

**Purpose**: Source of truth for backend endpoints

**Style**:

- One section per endpoint
- Request/response schemas with examples
- HTTP status codes documented
- Error responses explained

**Requirements**:

- **Must update** when API changes
- JSON examples must be valid
- Keep in sync with implementation
- Include validation rules

### 3. Constitution (`.specify/memory/constitution.md`)

**Purpose**: Non-negotiable project principles

**Style**:

- Principle-based structure
- Rationale for each principle
- Version tracking
- Sync impact report at top

**Requirements**:

- Never modify without consensus
- Version bump on changes
- Update templates when principles change

### 4. User-Facing Documentation (README, guides)

**Purpose**: Help users understand and use the system

**Style**:

- Clear step-by-step instructions
- Troubleshooting sections
- Screenshots/diagrams where helpful
- Command examples (copy-paste ready)

## Documentation Tasks

### Task: Update API Documentation

**When**: Any time an endpoint changes

**Process**:

1. Locate endpoint in `backend/API.md`
2. Update request/response schemas
3. Update status codes if changed
4. Add/update error responses
5. Verify examples are valid JSON
6. Test endpoint to confirm accuracy

**Template**:

````markdown
### X. Endpoint Name

Returns description of what it does.

- **URL**: `/path/{param}`
- **Method**: `GET`/`POST`/`DELETE`
- **Request Body** (if applicable):
  ```json
  { "field": "value" }
  ```
- **Response Body**: `Type`
  ```json
  { "field": "value" }
  ```
- **Error Responses**:
  - `400 Bad Request`: Description
````

### Task: Document New Architectural Decision

**When**: Major design choice is made

**Process**:

1. Add to `.github/ai-context/knowledge/decisions.md`
2. Use template: DNNN: [Title], Decision, Context, Trade-offs, Rationale, Impact
3. Include table of alternatives considered
4. Explain why chosen approach wins
5. Note impact on codebase

### Task: Update Constitution

**When**: Core principle changes or new principle added

**Process**:

1. Update sync impact report at top
2. Modify/add principle section
3. Update version number
4. Update ratified/amended dates
5. Update all templates in `.specify/templates/` to align
6. Verify no contradictions with existing principles

### Task: Add to Project History

**When**: Major milestone reached or significant refactor

**Process**:

1. Update `.github/ai-context/knowledge/memory.md`
2. Add to timeline under appropriate phase
3. Document what changed and why
4. Note any abandoned approaches
5. Record lessons learned

### Task: Update Session Context

**When**: Completing major work session

**Process**:

1. Update `.github/ai-context/knowledge/chatmem.md`
2. Move completed items from "in progress" to "complete"
3. Add new "in progress" items if work continues
4. Update "Recent Sessions" with summary
5. Add new pitfalls/discoveries if encountered

## Documentation Standards

### Markdown Formatting

- **Headings**: Use ATX-style (`#`, `##`, `###`)
- **Code blocks**: Always specify language (` ```go`, ` ```typescript`, ` ```bash`)
- **Lists**: Use `-` for unordered, `1.` for ordered
- **Tables**: Use Markdown tables, align with `|---|`
- **Links**: Use relative paths for internal docs

### Mermaid Diagrams

**When to use**:

- System architecture
- Data flows
- Sequence diagrams
- State machines
- ER diagrams

**Style**:

```markdown
\`\`\`mermaid
graph TB
A[Component A] --> B[Component B]
B --> C[Component C]
\`\`\`
```

### Code Examples

**Requirements**:

- Must be valid, runnable code
- Include necessary imports
- Show both success and error cases
- Annotate non-obvious logic with comments
- Use realistic variable names (not `foo`, `bar`)

**Example**:

```go
// Good example
func (h *PeerHandler) Add(w http.ResponseWriter, r *http.Request) {
    var req AddPeerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("failed to decode request", "error", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    // ... rest of handler
}
```

### Tone & Voice

- **Technical accuracy** over brevity
- **Active voice** ("The system does X" not "X is done")
- **Present tense** ("The handler validates input")
- **Clear and direct** (avoid marketing language)
- **Assume technical audience** (developers, not end-users for AI context docs)

## Documentation Maintenance

### Regular Updates

**Weekly**:

- Review `chatmem.md` — Update active work status

**Per Release**:

- Update version numbers in constitution, API.md
- Review architecture.md for accuracy
- Add to project history (memory.md)

**On Major Changes**:

- Update decisions.md with new architectural decisions
- Update constitution if principles change
- Update all affected AI context files

### Consistency Checks

Before committing documentation:

- [ ] Code examples are valid and tested
- [ ] Links are not broken (use relative paths)
- [ ] Mermaid diagrams render correctly
- [ ] No contradictions with constitution
- [ ] API.md matches actual implementation
- [ ] Types in examples match `types.ts` and Go structs

## Special Documentation Rules

### For AI Context Files

- **Be detailed**: AI agents benefit from comprehensive context
- **Include rationale**: Explain _why_, not just _what_
- **Provide examples**: Show patterns in code snippets
- **Cross-reference**: Link to related docs

### For User-Facing Docs

- **Be concise**: Users want quick answers
- **Include troubleshooting**: Anticipate common issues
- **Step-by-step**: Use numbered lists for procedures
- **Visual aids**: Screenshots/diagrams where helpful

### For API Documentation

- **Accuracy is critical**: Frontend depends on this
- **Include all edge cases**: Error responses, validation
- **JSON examples must be valid**: Test them
- **Keep synchronized**: Update with every API change

## Tools & Commands

### Preview Documentation

```bash
# Markdown preview (VSCode built-in)
Cmd/Ctrl + Shift + V

# Mermaid preview (VSCode with mermaid extension)
Install: mermaid-preview or markdown-preview-mermaid-support
```

### Validate Links

```bash
# Check for broken links (markdown-link-check)
npx markdown-link-check README.md
```

### Format Markdown

```bash
# Prettier for consistent formatting
npm run format  # formats all Markdown files
```

## When to Create New Documentation

### New AI Context File

**When**: New domain/skill area emerges (e.g., `kubernetes.skills.md` if K8s added)

**Where**: `.github/ai-context/skills/` or `.github/ai-context/knowledge/`

**Process**: Follow existing file structure, add to copilot-instructions.md

### New User Guide

**When**: New user-facing feature added

**Where**: Root `README.md` or separate `docs/` folder

**Process**: Write step-by-step guide with examples

### New Architecture Diagram

**When**: System structure changes significantly

**Where**: `.github/ai-context/knowledge/architecture.md`

**Process**: Use Mermaid, explain components and data flows

## Documentation Anti-Patterns

**Avoid**:

- ❌ Documenting what's obvious from code (e.g., "This function returns a value")
- ❌ Copy-pasting code without explanation (show _why_, not just _what_)
- ❌ Outdated examples (keep in sync with implementation)
- ❌ Vague descriptions ("The system handles this")
- ❌ Marketing language ("Our amazing feature...")
- ❌ Incomplete examples (missing imports, invalid code)
- ❌ Broken links (use relative paths, verify they work)

**Do**:

- ✅ Explain _why_ architecture is designed this way
- ✅ Document trade-offs and alternatives considered
- ✅ Show realistic, runnable examples
- ✅ Include error cases and edge cases
- ✅ Cross-reference related documentation
- ✅ Keep synchronized with code changes
- ✅ Provide context for AI agents and developers
