# Go Backend Skills & Conventions

**Location**: `backend/` directory

## Language Version & Setup

- **Go Version**: 1.25.6 (see `go.mod`)
- **Module**: `wg-manager/backend`
- **Dependencies**:
  - `golang.zx2c4.com/wireguard/wgctrl` — WireGuard control interface
  - `github.com/joho/godotenv` — Environment variable loading
- **Logging**: Standard library `log/slog` (JSON structured logging)
- **Web**: Standard library `net/http` (no external framework)

## Project Structure Conventions

```
backend/
├── cmd/server/
│   ├── main.go              # Server entry point, router setup
│   └── main_test.go         # Integration tests for handlers
├── internal/
│   ├── config/              # Configuration loading (Twelve-Factor App)
│   │   ├── config.go        # Config struct and LoadConfig()
│   │   ├── config.json      # Default configuration
│   │   └── ...
│   ├── handlers/            # HTTP request handlers (business logic dispatch)
│   │   ├── handlers.go      # PeerHandler methods (List, Add, Remove, Stats)
│   │   └── ...
│   ├── middleware/          # HTTP middleware (cross-cutting concerns)
│   │   ├── cors.go          # CORS header setup
│   │   ├── logging.go       # Request/response logging
│   │   └── ...
│   └── wireguard/           # WireGuard abstraction layer
│       ├── service.go       # Service interface and realService implementation
│       ├── mock.go          # Mock service for testing/development
│       ├── storage.go       # Peer metadata persistence (JSON file)
│       ├── keys.go          # Key generation utilities
│       ├── config_gen.go    # WireGuard config template generation
│       └── ...
└── data/
    └── peers.json           # Persistent peer metadata
```

## Handler Pattern

All handlers follow this structure:

```go
type PeerHandler struct {
    Service wireguard.Service
}

func (h *PeerHandler) MethodName(w http.ResponseWriter, r *http.Request) {
    // Validation
    // Business logic delegation to Service
    // Response serialization
    // Error handling with slog
}
```

**Key Points**:

- Validation happens in handler (CIDR parsing, required fields)
- All business logic delegated to `wireguard.Service` interface
- Error logging uses `slog.Error()` with context
- HTTP status codes: 201 (Created), 204 (No Content), 400 (Bad Request), 500 (Server Error)

## Service Interface Pattern

The `wireguard.Service` interface abstracts WireGuard operations:

```go
type Service interface {
    ListPeers() ([]Peer, error)
    AddPeer(name, publicKey string, allowedIPs []string) (*Peer, error)
    RemovePeer(peerID string) error
    GetStats() (*Stats, error)
}
```

**Implementations**:

- `realService` — Actual WireGuard interface (requires kernel module + permissions)
- `mockService` — Test/development fallback

## Middleware Pattern

Middleware wraps the `http.Handler`:

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-request logic
        slog.Info("request", "method", r.Method, "path", r.URL.Path)
        next.ServeHTTP(w, r)
        // Post-request logic
    })
}
```

Applied in `main.go`:

```go
wrappedMux := middleware.LoggingMiddleware(mux)
wrappedMux = middleware.CORSMiddleware(wrappedMux)
```

## Configuration Pattern (Twelve-Factor)

Configuration loading hierarchy:

1. Load defaults from `config.json`
2. Override with environment variables (if set)
3. Support `.env` file via `godotenv.Load()`

**Environment variables**:

- `WG_SERVER_PORT` — Server listen address (default `:8080`)
- `WG_INTERFACE_NAME` — WireGuard interface (default `wg0`)
- `WG_STORAGE_PATH` — Peer metadata file (default `./data/peers.json`)
- `WG_SERVER_ENDPOINT` — Public server endpoint for peer config (default from config.json)
- `WG_SERVER_PUBKEY` — Server's public key (default from config.json)
- `CORS_ALLOWED_ORIGINS` — Comma-separated origins (reflective by default in dev)

## Testing Conventions (TDD Mandatory)

**Test file locations**: `*_test.go` in the same package

**Test pattern**:

```go
func TestFunctionName(t *testing.T) {
    // Arrange: Setup test data, mocks
    // Act: Call function under test
    // Assert: Verify output and side effects
}
```

**Key testing principles** (per Constitution I):

- Write tests FIRST, then implementation (Red-Green-Refactor)
- All handlers MUST have unit tests
- Services MUST have contract tests
- Integration tests for WireGuard operations (using mockService)

## Logging Conventions

**Structured logging via `slog`**:

```go
slog.Info("operation_name",
    "key1", value1,
    "key2", value2)

slog.Error("failed operation",
    "error", err,
    "context_key", context_value)
```

**Log levels**:

- `ERROR` — Operation failures, unrecoverable errors
- `INFO` — Important operations (server start, peer added, removed)
- `WARN` — Degraded functionality (fallback to mock service)
- `DEBUG` — Diagnostic info (rarely used in production)

**What to log**:

- Every API request (method, path, status, duration)
- WireGuard operations (add, remove peer; success/failure)
- Configuration loading
- Server startup/shutdown
- NO sensitive data (private keys, passwords)

## Error Handling

**Style**:

```go
if err != nil {
    slog.Error("operation failed", "error", err, "context", value)
    http.Error(w, "User-facing message", http.StatusInternalServerError)
    return
}
```

**Key points**:

- Log with context (what operation failed, what data involved)
- Return generic user-facing error messages (don't leak internals)
- Use appropriate HTTP status codes
- Always include error wrapping: `fmt.Errorf("failed to X: %w", err)`

## WireGuard Domain Knowledge

- **Peer**: A client connected to the WireGuard interface
  - Identified by `PublicKey` (immutable)
  - Has optional metadata (name) stored in `peers.json`
  - Has allowed IPs (CIDR notation, e.g., `10.0.0.2/32`)
  - Tracks handshake time, rx/tx bytes in real-time
- **Interface**: The WireGuard virtual network interface (e.g., `wg0`)
  - Has one public key, multiple peers
  - Managed by kernel module, accessed via `wgctrl` library
- **Real vs. Mock**:
  - Real service requires root permissions and active WireGuard kernel module
  - Mock service stores peers in memory; useful for development/testing

## Common Tasks

### Adding a new endpoint:

1. Add handler method to `PeerHandler`
2. Register route in `main.go`: `mux.HandleFunc("METHOD /path/{param}", handler.Method)`
3. Write tests FIRST in `main_test.go` (TDD)
4. Implement handler method
5. Update `API.md` with endpoint docs

### Modifying existing handler:

1. Write failing test for new behavior
2. Implement the handler change
3. Run tests to verify
4. Update API documentation

### Adding new service logic:

1. Define new method in `Service` interface
2. Implement in both `realService` and `mockService`
3. Write contract tests for the interface
4. Use in handler
