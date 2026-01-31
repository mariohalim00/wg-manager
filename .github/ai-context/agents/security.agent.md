# Security-Focused Agent Rules

**Purpose**: Guidelines for AI agents working on security-sensitive code.

## Security Philosophy

This project manages **critical network infrastructure** (WireGuard VPN configuration). Security is paramount.

**Principle**: Defense in depth. Multiple layers of protection, fail securely.

## Security Priorities

1. **Input validation** — Never trust user input
2. **Secret management** — No hardcoded credentials
3. **Access control** — Assume hostile network
4. **Logging security** — No sensitive data in logs
5. **Secure defaults** — Opt-in for risky features

## Critical Security Rules

### Rule 1: Input Validation (Backend)

**Always validate user input before processing.**

**CIDR notation** (AllowedIPs):
```go
// REQUIRED validation
for _, ip := range req.AllowedIPs {
    if _, _, err := net.ParseCIDR(ip); err != nil {
        http.Error(w, fmt.Sprintf("Invalid AllowedIP CIDR: %s", ip), http.StatusBadRequest)
        return
    }
}
```

**PublicKey format**:
```go
// WireGuard public keys are base64-encoded, ~44 chars
// wgctrl library validates format internally
```

**Required fields**:
```go
req.Name = strings.TrimSpace(req.Name)
if req.Name == "" {
    http.Error(w, "Name is required", http.StatusBadRequest)
    return
}
```

**Never**:
- Trust input lengths (check before allocating large buffers)
- Allow SQL injection (not applicable here, but principle applies)
- Execute unvalidated commands
- Parse untrusted data without bounds checking

### Rule 2: Secret Management

**NO hardcoded secrets in code.**

**Bad**:
```go
const ADMIN_PASSWORD = "supersecret"  // ❌ NEVER
const API_KEY = "abc123"              // ❌ NEVER
```

**Good**:
```go
// Load from environment variables
adminPassword := os.Getenv("ADMIN_PASSWORD")
if adminPassword == "" {
    slog.Error("ADMIN_PASSWORD not set")
    os.Exit(1)
}
```

**Configuration**:
- Use `.env` file (gitignored) for local development
- Use environment variables in production
- Document all required env vars in `backend/API.md`

**WireGuard private keys**:
- Generated on demand (never stored by backend)
- Returned to user once, never logged
- User responsible for securing their private key

### Rule 3: Logging Security

**Never log sensitive data.**

**Sensitive**:
- Private keys (WireGuard PrivateKey)
- Passwords, API keys, tokens
- Full request bodies (may contain secrets)
- Personal identifiable information (PII)

**Safe to log**:
- Public keys (PublicKey)
- IP addresses (already public in WireGuard)
- Peer names (user-chosen, non-sensitive)
- Timestamps, operation types, error messages

**Example**:
```go
// ❌ BAD: Logs private key
slog.Info("peer added", "privateKey", privateKey)

// ✅ GOOD: Logs public key only
slog.Info("peer added", "publicKey", publicKey, "name", name)
```

### Rule 4: Access Control

**Assumption**: API is exposed on private network or behind authentication proxy.

**Current state**:
- No built-in authentication
- Relies on network security (firewall, VPN, reverse proxy)
- **NOT suitable for public internet exposure**

**If adding authentication**:
- Use standard auth mechanisms (OAuth, JWT, Basic Auth over HTTPS)
- Never roll custom crypto
- Hash passwords with bcrypt or argon2
- Use HTTPS in production (TLS/SSL)

**API design**:
- No authorization bypass via path manipulation
- No IDOR (Insecure Direct Object References)
  - Currently mitigated: peer ID = public key (user can't guess others' public keys)

### Rule 5: CORS Security

**CORS is a trust boundary.**

**Development**:
```go
// Reflective CORS (allows any origin) — ONLY for development
if cfg.CORSAllowedOrigins == "" {
    w.Header().Set("Access-Control-Allow-Origin", origin)
}
```

**Production**:
```go
// Explicit whitelist
allowedOrigins := strings.Split(cfg.CORSAllowedOrigins, ",")
for _, allowed := range allowedOrigins {
    if origin == strings.TrimSpace(allowed) {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        break
    }
}
```

**Never**:
- `Access-Control-Allow-Origin: *` in production
- Allow credentials (`Access-Control-Allow-Credentials: true`) with wildcard origin

### Rule 6: Error Messages

**Don't leak implementation details in error responses.**

**Bad**:
```go
http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
// ❌ Leaks internal structure
```

**Good**:
```go
slog.Error("database operation failed", "error", err, "operation", "add_peer")
http.Error(w, "Internal server error", http.StatusInternalServerError)
// ✅ Generic user message, detailed logs internally
```

**Validation errors** (400 Bad Request):
- **OK to be specific**: "Invalid CIDR notation: 10.0.0.256/32"
- Helps user correct input
- Doesn't reveal internal state

**Server errors** (500 Internal Server Error):
- **Be generic**: "Internal server error"
- Log details with `slog.Error()`
- Don't expose stack traces, file paths, database schema

## Security Anti-Patterns

### Anti-Pattern 1: Trusting User Input

**Symptom**: Directly using user input without validation

**Risk**: Injection attacks, buffer overflows, crashes

**Solution**: Validate all input (CIDR, required fields, lengths)

### Anti-Pattern 2: Logging Secrets

**Symptom**: `slog.Info("key", privateKey)`

**Risk**: Secrets exposed in logs, attackable surface area

**Solution**: Audit all log statements, never log private keys or passwords

### Anti-Pattern 3: Hardcoded Credentials

**Symptom**: `const API_KEY = "secret"`

**Risk**: Committed to git, exposed in binary, rotation impossible

**Solution**: Environment variables, secret management systems

### Anti-Pattern 4: Reflective CORS in Production

**Symptom**: `Access-Control-Allow-Origin: <request origin>`

**Risk**: Any site can make requests, CSRF potential

**Solution**: Whitelist specific origins in production

### Anti-Pattern 5: Overly Detailed Error Messages

**Symptom**: `"Failed to connect to database at 10.0.0.5:5432"`

**Risk**: Leaks internal network topology, attackable

**Solution**: Generic error to user, detailed error in logs

## Secure Coding Patterns

### Pattern 1: Fail Securely

**Principle**: On error, deny access by default.

```go
if err := validateInput(req); err != nil {
    http.Error(w, "Invalid input", http.StatusBadRequest)
    return  // ✅ Fail closed: don't proceed
}
// Proceed with validated input
```

### Pattern 2: Principle of Least Privilege

**Principle**: Run with minimum required permissions.

**Current**:
- Backend needs `CAP_NET_ADMIN` for WireGuard operations
- Fall back to mock service if lacking permissions (safe degradation)

**Future**:
- Use systemd capabilities instead of full root
- Run backend as non-root user with specific capabilities

### Pattern 3: Defense in Depth

**Principle**: Multiple layers of protection.

**Layers in this project**:
1. **Network**: Firewall, VPN, private network
2. **Input validation**: CIDR parsing, required fields
3. **Service interface**: Abstracts WireGuard operations (prevents direct kernel manipulation)
4. **Logging**: Audit trail for all operations
5. **Error handling**: Generic user messages, detailed internal logs

### Pattern 4: Secure Defaults

**Principle**: Default configuration should be secure.

**Examples**:
- CORS in production: Whitelist only (not reflective)
- Logging: JSON structured (queryable, no sensitive data)
- Configuration: No secrets in `config.json` (use env vars)

## Security Tasks

### Task: Review New Endpoint for Security

**Checklist**:
- [ ] Input validation implemented (CIDR, required fields)
- [ ] No hardcoded secrets or credentials
- [ ] Error messages are generic (no internal details leaked)
- [ ] Logging does not include sensitive data
- [ ] CORS headers appropriate for environment
- [ ] Authorization checks if applicable (future: auth layer)

### Task: Add Authentication (Future)

**Considerations**:
- Use well-tested library (e.g., OAuth, JWT)
- HTTPS required (no plain HTTP with auth)
- Hash passwords with bcrypt or argon2
- Implement rate limiting (prevent brute force)
- Use secure session management (HTTP-only cookies, CSRF tokens)

### Task: Secure Logging Changes

**Before logging new data**:
- [ ] Is this data sensitive? (keys, passwords, PII)
- [ ] Could this leak internal structure? (file paths, database schema)
- [ ] Is this data necessary for debugging? (justify inclusion)

**If sensitive**:
- Hash or redact before logging
- Or omit entirely

## Threat Model

### Threat 1: Malicious API User

**Scenario**: Attacker has network access to API

**Mitigations**:
- Input validation (prevents malformed requests)
- No secrets in responses (prevents key extraction)
- Audit logging (detects suspicious activity)

**Residual risk**: Attacker can add/remove peers if no authentication layer

**Recommendation**: Deploy behind authentication proxy or firewall

### Threat 2: Log Injection

**Scenario**: Attacker injects malicious data into logs

**Mitigations**:
- Structured JSON logging (parseable, no injection)
- Log sanitization (strip control characters)
- No user input directly in log messages (use key-value pairs)

**Example**:
```go
// ❌ BAD: Log injection possible
slog.Info(fmt.Sprintf("User %s logged in", username))

// ✅ GOOD: Structured, injection-safe
slog.Info("user login", "username", username)
```

### Threat 3: Private Key Exposure

**Scenario**: User's private key leaked

**Mitigations**:
- Private key returned only once (POST /peers response)
- Never stored by backend
- Never logged
- User responsible for securing key

**Residual risk**: If user doesn't secure key, attacker can impersonate peer

**Recommendation**: Educate users to store keys securely

### Threat 4: CORS Misconfiguration

**Scenario**: Malicious website makes API requests on behalf of user

**Mitigations**:
- Whitelist CORS origins in production
- No credentials with wildcard origin
- Reflective CORS only in development

**Residual risk**: If CORS misconfigured, CSRF attacks possible

**Recommendation**: Verify CORS settings before production deployment

## Security Checklist for Code Review

Before approving code:

- [ ] No hardcoded secrets or credentials
- [ ] Input validation for all user-supplied data
- [ ] No sensitive data logged (private keys, passwords)
- [ ] Error messages are generic (no internal details)
- [ ] CORS configuration appropriate for environment
- [ ] Authentication/authorization checks if applicable
- [ ] Secure defaults (no insecure opt-out required)
- [ ] Dependencies are up-to-date (no known vulnerabilities)

## Security Resources

- **OWASP Top 10**: https://owasp.org/www-project-top-ten/
- **Go Security Checklist**: https://github.com/guardrailsio/awesome-golang-security
- **WireGuard Protocol**: https://www.wireguard.com/papers/wireguard.pdf

## Reporting Security Issues

**DO NOT** open public GitHub issues for security vulnerabilities.

**Instead**:
- Email project maintainer directly
- Include: vulnerability description, steps to reproduce, impact assessment
- Allow time for patch before public disclosure
